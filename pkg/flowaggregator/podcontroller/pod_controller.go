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

package podcontroller

import (
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

const (
	controllerName = "FlowAggregatorPodController"
	// Set resyncPeriod to 0 to disable resyncing.
	resyncPeriod time.Duration = 0
)

type PodInfo struct {
	Namespace string
	Name      string
	NodeName  string
}

type PodController struct {
	podInformer     coreinformers.PodInformer
	podLister       corelisters.PodLister
	podListerSynced cache.InformerSynced
	podIPStore      map[string]PodInfo
	mutex           sync.Mutex
}

func NewPodController(
	informerFactory informers.SharedInformerFactory,
) *PodController {
	podInformer := informerFactory.Core().V1().Pods()
	controller := &PodController{
		podInformer:     podInformer,
		podLister:       podInformer.Lister(),
		podListerSynced: podInformer.Informer().HasSynced,
		podIPStore:      make(map[string]PodInfo),
	}
	podInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    controller.addPod,
			UpdateFunc: controller.updatePod,
			DeleteFunc: controller.deletePod,
		},
		resyncPeriod,
	)
	return controller
}

func (p *PodController) GetPodInfoByIP(podIP string) (PodInfo, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	podInfo, exist := p.podIPStore[podIP]
	if !exist {
		klog.V(4).Infof("No pod is cached with the corresponding ip %s.", podIP)
	}
	return podInfo, exist
}

func (p *PodController) syncPodStore(action, podIP string, podInfo PodInfo) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, exists := p.podIPStore[podIP]; !exists && action == "DELETE" {
		klog.V(4).Infof("Cannot find the old ip %s in store.", podIP)
		return
	}
	if action == "UPDATE" {
		p.podIPStore[podIP] = podInfo
	} else {
		delete(p.podIPStore, podIP)
	}
}

func (p *PodController) addPod(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		klog.V(4).Info("Object does not have correct Pod interface")
		return
	}

	klog.V(2).Infof("Processing Pod %s/%s ADD event", pod.Namespace, pod.Name)
	if pod.Status.PodIP == "" {
		// Can't add pod with IP as key when IP is unset.
		klog.V(4).Infof("Pod %s/%s has not been assigned IP yet. Defer processing", pod.Namespace, pod.Name)
		return
	}
	p.syncPodStore("UPDATE", pod.Status.PodIP, PodInfo{pod.Namespace, pod.Name, pod.Spec.NodeName})
}

func (p *PodController) updatePod(oldObj, curObj interface{}) {
	oldPod := oldObj.(*v1.Pod)
	curPod := curObj.(*v1.Pod)
	klog.V(2).Infof("Processing Pod %s/%s UPDATE event", curPod.Namespace, curPod.Name)
	if curPod.Status.PodIP == oldPod.Status.PodIP {
		return
	}
	p.syncPodStore("UPDATE", curPod.Status.PodIP, PodInfo{curPod.Namespace, curPod.Name, curPod.Spec.NodeName})
}

func (p *PodController) deletePod(old interface{}) {
	pod, ok := old.(*v1.Pod)
	if !ok {
		tombstone, ok := old.(cache.DeletedFinalStateUnknown)
		if !ok {
			klog.Errorf("Error decoding object when deleting Pod, invalid type: %v", old)
			return
		}
		pod, ok = tombstone.Obj.(*v1.Pod)
		if !ok {
			klog.Errorf("Error decoding object tombstone when deleting Pod, invalid type: %v", tombstone.Obj)
			return
		}
	}
	klog.V(2).Infof("Processing Pod %s/%s DELETE event", pod.Namespace, pod.Name)
	p.syncPodStore("DELETE", pod.Status.PodIP, PodInfo{pod.Namespace, pod.Name, pod.Spec.NodeName})
}

func (p *PodController) Run(stopCh <-chan struct{}) {
	klog.Infof("Starting %s", controllerName)
	defer klog.Infof("Shutting down %s", controllerName)

	klog.Infof("Waiting for caches to sync for %s", controllerName)
	if !cache.WaitForCacheSync(stopCh, p.podListerSynced) {
		klog.Errorf("Unable to sync caches for %s", controllerName)
		return
	}
	klog.Infof("Caches are synced for %s", controllerName)
	<-stopCh
}
