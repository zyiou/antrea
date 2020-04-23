package graphviz

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"k8s.io/klog"

	"github.com/vmware-tanzu/antrea/pkg/apis/traceflow/v1"
)

var componentString = map[v1.ComponentType]string{
	v1.SPOOFGUARD: "SPOOFGUARD",
	v1.LB:         "LB",
	v1.ROUTING:    "ROUTING",
	v1.DFW:        "DFW",
	v1.FORWARDING: "FORWARDING",
}

func getObservation(target *v1.Observation, obs []v1.Observation) *v1.Observation {
	for i := range obs {
		if target.ComponentType == obs[i].ComponentType {
			return &obs[i]
		}
	}
	return nil
}

func genSubGraph(graph *cgraph.Graph, expectedObs []v1.Observation, obs []v1.Observation, dir cgraph.DirType) []*cgraph.Node {
	var nodes []*cgraph.Node
	if dir == cgraph.BackDir {
		for i := len(expectedObs)/2 - 1; i >= 0; i-- {
			opp := len(expectedObs) - 1 - i
			expectedObs[i], expectedObs[opp] = expectedObs[opp], expectedObs[i]
		}
	}
	for i, o := range expectedObs {
		if len(o.NodeUUID) > 0 {
			graph.SetLabel(o.NodeUUID)
		}
		nodeName := fmt.Sprintf("%s_%d", graph.Name(), i)
		node, _ := graph.CreateNode(nodeName)
		node.SetLabel(componentString[o.ComponentType])
		node.SetShape(cgraph.BoxShape)
		node.SetStyle(cgraph.RoundedNodeStyle)
		nodes = append(nodes, node)
		if i == 0 {
			continue
		}
		edgeName := fmt.Sprintf("%s_%d_%d", graph.Name(), i-1, i)
		edge, _ := graph.CreateEdge(edgeName, nodes[i-1], nodes[i])
		edge.SetDir(dir)
		actualObj := getObservation(&o, obs)
		if actualObj == nil {
			node.SetColor("red")
			edge.SetStyle("invis")
		}
	}
	return nodes
}

var expectedSenderObservations = []v1.Observation{
	{
		ComponentType: v1.SPOOFGUARD,
	},
	{
		ComponentType: v1.DFW,
	},
	{
		ComponentType: v1.ROUTING,
	},
}

var expectedReceiverObservations = []v1.Observation{
	{
		ComponentType: v1.ROUTING,
	},
	{
		ComponentType: v1.DFW,
	},
	{
		ComponentType: v1.FORWARDING,
	},
}

func GenGraph(tf *v1.Traceflow) string {
	var err error
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		klog.Error(err)
	}
	cluster1 := graph.SubGraph("cluster_source", 1)
	cluster1.SetStyle(cgraph.FilledGraphStyle)
	srcPod, _ := cluster1.CreateNode(tf.SrcNamespace + "/" + tf.SrcPod)
	nodes1 := genSubGraph(cluster1, expectedSenderObservations, tf.Status.NodeSender, cgraph.ForwardDir)
	cluster1.CreateEdge("src-pod-out", srcPod, nodes1[0])

	cluster2 := graph.SubGraph("cluster_dest", 1)
	cluster2.SetStyle(cgraph.FilledGraphStyle)
	dstPod, _ := cluster2.CreateNode(tf.DstNamespace + "/" + tf.DstPod)
	nodes2 := genSubGraph(cluster2, expectedReceiverObservations, tf.Status.NodeReceiver, cgraph.BackDir)
	dstPodIn, _ := cluster2.CreateEdge("dst-pod-in", dstPod, nodes2[0])
	dstPodIn.SetDir(cgraph.BackDir)
	//dstPodIn.SetStyle("invis")

	edge, _ := graph.CreateEdge("cross_node", nodes1[len(nodes1)-1], nodes2[len(nodes2)-1])
	edge.SetConstraint(false)

	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		klog.Fatal(err)
	}
	klog.Infof("Graphz:\n%v", buf.String())
	if err := graph.Close(); err != nil {
		klog.Fatal(err)
	}
	g.Close()
	return buf.String()
}
