// Copyright 2020 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package traceflow

import (
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	"github.com/vmware-tanzu/antrea/pkg/agent/config"
	"github.com/vmware-tanzu/antrea/pkg/agent/interfacestore"
	"github.com/vmware-tanzu/antrea/pkg/agent/openflow"
	traceflowv1 "github.com/vmware-tanzu/antrea/pkg/apis/traceflow/v1"
	clientsetversioned "github.com/vmware-tanzu/antrea/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/antrea/pkg/ovs/ovsconfig"
)

const (
	controllerName = "AntreaAgentTraceflowController"
	// How long to wait before retrying the processing of a traceflow
	minRetryDelay = 5 * time.Second
	maxRetryDelay = 300 * time.Second
	// Default number of workers processing traceflow request
	defaultWorkers = 1
)

// Controller is responsible for setting up Openflow entries and inject traceflow packet into
// switch for traceflow request.
type Controller struct {
	traceflowClient clientsetversioned.Interface
	ovsBridgeClient ovsconfig.OVSBridgeClient
	ofClient        openflow.Client
	interfaceStore  interfacestore.InterfaceStore
	networkConfig   *config.NetworkConfig
	nodeConfig      *config.NodeConfig
	nodeLister      corelisters.NodeLister
	namespaceLister corelisters.NamespaceLister
	podLister       corelisters.PodLister
	serviceLister   corelisters.ServiceLister
	queue           workqueue.RateLimitingInterface
	running         map[uint8]*traceflowv1.Traceflow // tag->tf if tf.Status.Phase is INITIAL or RUNNING
	senders         map[uint8]*traceflowv1.Traceflow // tag->tf if tf is RUNNING and this Node is sender
}

// NewTraceflowController instantiates a new Controller object which will process Traceflow
// events.
func NewTraceflowController(
	traceflowClient clientsetversioned.Interface,
	informerFactory informers.SharedInformerFactory,
	client openflow.Client,
	ovsBridgeClient ovsconfig.OVSBridgeClient,
	interfaceStore interfacestore.InterfaceStore,
	networkConfig *config.NetworkConfig,
	nodeConfig *config.NodeConfig) *Controller {

	nodeInformer := informerFactory.Core().V1().Nodes()
	namespaceInformer := informerFactory.Core().V1().Namespaces()
	podInformer := informerFactory.Core().V1().Pods()
	serviceInformer := informerFactory.Core().V1().Services()
	running := make(map[uint8]*traceflowv1.Traceflow)
	senders := make(map[uint8]*traceflowv1.Traceflow)

	controller := &Controller{
		traceflowClient: traceflowClient,
		ovsBridgeClient: ovsBridgeClient,
		ofClient:        client,
		interfaceStore:  interfaceStore,
		networkConfig:   networkConfig,
		nodeConfig:      nodeConfig,
		nodeLister:      nodeInformer.Lister(),
		namespaceLister: namespaceInformer.Lister(),
		podLister:       podInformer.Lister(),
		serviceLister:   serviceInformer.Lister(),
		queue:           workqueue.NewNamedRateLimitingQueue(workqueue.NewItemExponentialFailureRateLimiter(minRetryDelay, maxRetryDelay), "traceflow"),
		running:         running,
		senders:         senders}
	return controller
}

// enqueueTraceflow adds an object to the controller work queue.
func (c *Controller) enqueueTraceflow(tf *traceflowv1.Traceflow) {
	c.queue.Add(tf.Name)
}

// Run will create defaultWorkers workers (go routines) which will process the Traceflow events from the
// workqueue.
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer c.queue.ShutDown()

	// If agent is running noancap mode, traceflow is not supported.
	if c.networkConfig.TrafficEncapMode.SupportsNoEncap() {
		<-stopCh
		return
	}

	klog.Infof("Starting %s", controllerName)
	defer klog.Infof("Shutting down %s", controllerName)

	go c.ReceivePacketIn(stopCh)
	for i := 0; i < defaultWorkers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}
	wait.PollUntil(time.Second, func() (done bool, err error) {
		list, err := c.traceflowClient.AntreaV1().Traceflows().List(v1.ListOptions{})
		if err != nil {
			klog.Info("Fail to list all Antrea Traceflow CRD: %v", err)
			return false, err
		}
		for _, tf := range list.Items {
			p := tf.Status.Phase
			tag := tf.Status.CrossNodeTag
			klog.Infof("DEBUG: Get traceflow %s phase %s tag %s", tf.Name, p, tag)
			switch p {
			case
				traceflowv1.INITIAL,
				traceflowv1.RUNNING:
				if tag != 0 {
					klog.Infof("DEBUG: Add traceflow %s to queue", tf.Name)
					c.running[tag] = &tf
					c.enqueueTraceflow(&tf)
				}
			}
		}
		return false, nil
	}, stopCh)
	<-stopCh
}

// worker is a long-running function that will continually call the processTraceflowItem function
// in order to read and process a message on the workqueue.
func (c *Controller) worker() {
	for c.processTraceflowItem() {
	}
}

// processTraceflowItem processes an item in the "traceflow" work queue, by calling syncTraceflow
// after casting the item to a string (Traceflow name). If syncTraceflow returns an error, this
// function handles it by updating phase to "ERROR". If syncTraceflow is successful, the Traceflow
// is removed from the queue until we get notified of a new change. This function returns false if
// and only if the work queue was shutdown (no more items will be processed).
func (c *Controller) processTraceflowItem() bool {
	obj, quit := c.queue.Get()
	if quit {
		return false
	}
	// We call Done here so the workqueue knows we have finished processing this item. We also
	// must remember to call Forget if we do not want this work item being re-queued. For
	// example, we do not call Forget if a transient error occurs, instead the item is put back
	// on the workqueue and attempted again after a back-off period.
	defer c.queue.Done(obj)

	// We expect strings (Node name) to come off the workqueue.
	if key, ok := obj.(string); !ok {
		// As the item in the workqueue is actually invalid, we call Forget here else we'd
		// go into a loop of attempting to process a work item that is invalid.
		// This should not happen: enqueueTraceflow only enqueues strings.
		c.queue.Forget(obj)
		klog.Errorf("Expected string in work queue but got %#v", obj)
		return true
	} else if err := c.syncTraceflow(key); err == nil {
		// If no error occurs we Forget this item so it does not get queued again until
		// another change happens.
		c.queue.Forget(key)
	} else {
		klog.Errorf("Error syncing Node %s, Aborting. Error: %v", key, err)
		c.queue.Forget(key)
	}
	return true
}

// syncTraceflow
func (c *Controller) syncTraceflow(traceflowName string) error {
	klog.Infof("DEBUG: Syncing traceflow %s", traceflowName)
	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("Finished syncing Traceflow for %s. (%v)", traceflowName, time.Since(startTime))
	}()

	tf, err := c.traceflowClient.AntreaV1().Traceflows().Get(traceflowName, v1.GetOptions{})
	if err != nil {
		return err
	}
	return c.startTraceflow(tf)
}

// startTraceflow
func (c *Controller) startTraceflow(tf *traceflowv1.Traceflow) error {
	// Deploy flow entries for traceflow
	err := c.deployFlowEntries(tf)
	defer func() {
		if err != nil {
			c.updateTraceflowCRD(tf, traceflowv1.ERROR)
		}
	}()
	if err != nil {
		return err
	}

	// Inject packet if this Node is sender
	return c.injectPacket(tf)
}

func (c *Controller) deployFlowEntries(tf *traceflowv1.Traceflow) error {
	klog.Infof("DEBUG: Installing traceflow %s entries", tf.Name)
	return c.ofClient.InstallTraceflowFlows(tf.CrossNodeTag)
}

func (c *Controller) injectPacket(tf *traceflowv1.Traceflow) error {
	namespaceName := tf.SrcNamespace
	podName := tf.SrcPod
	podInterface, ok := c.interfaceStore.GetContainerInterface(podName, namespaceName)
	// Pod not found in current Node, skip inject packet
	if !ok {
		return nil
	}

	// Update Traceflow phase to RUNNING
	klog.Infof("Injecting packet for Traceflow %s", tf.Name)
	c.senders[tf.CrossNodeTag] = tf
	tf, err := c.updateTraceflowCRD(tf, traceflowv1.RUNNING)
	if err != nil {
		klog.Errorf("Update Traceflow %s status failed", tf.Name)
		return err
	}

	// Calculate destination MAC/IP
	dstMAC := ""
	dstIP := tf.IPHeader.DstIP
	if dstIP == "" {
		dstPodInterface, ok := c.interfaceStore.GetContainerInterface(tf.DstPod, tf.DstNamespace)
		if ok {
			dstMAC = dstPodInterface.MAC.String()
			dstIP = dstPodInterface.IP.String()
		} else {
			dstPod, err := c.podLister.Pods(tf.DstNamespace).Get(tf.DstPod)
			if err != nil {
				return err
			}
			// dstMAC is "" for globalVirtualMAC
			dstIP = dstPod.Status.PodIP
		}
	}

	klog.Infof("DEBUG: Packet tag: %s, srcMAC: %s, dstMAC: %s, srcIP: %s, dstIP: %s, inPort: %s, TF: %+v",
		tf.CrossNodeTag, podInterface.MAC.String(), dstMAC, podInterface.IP.String(), dstIP, podInterface.OFPort, tf)

	return c.ofClient.SendTraceflowPacket(
		tf.CrossNodeTag,
		podInterface.MAC.String(),
		dstMAC,
		podInterface.IP.String(),
		dstIP,
		uint8(tf.Protocol),
		uint8(tf.TTL),
		uint16(tf.Flags),
		uint16(tf.TCPHeader.SrcPort),
		uint16(tf.TCPHeader.DstPort),
		uint8(tf.TCPFlags),
		uint16(tf.UDPHeader.SrcPort),
		uint16(tf.UDPHeader.DstPort),
		0,
		0,
		uint16(tf.ICMPEchoRequestHeader.ID),
		uint16(tf.ICMPEchoRequestHeader.Sequence),
		uint32(podInterface.OFPort),
		-1)
}

// Update traceflow CRD status
func (c *Controller) updateTraceflowCRD(
	tf *traceflowv1.Traceflow,
	phase traceflowv1.Phase) (*traceflowv1.Traceflow, error) {

	tf.Status.Phase = phase
	return c.traceflowClient.AntreaV1().Traceflows().Update(tf)
}
