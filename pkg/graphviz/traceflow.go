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

var resourceString = map[v1.ResourceType]string{
	v1.DELIVERED: "Delivered",
	v1.RECEIVED:  "Received",
	v1.FORWARDED: "Forwarded",
	v1.DROPPED:   "Dropped",
}

func getObservation(target *v1.Observation, obs []v1.Observation) *v1.Observation {
	for i := range obs {
		if target.ComponentType == v1.DFW && target.ComponentType == obs[i].ComponentType {
			return &obs[i]
		}
		if target.ComponentType == obs[i].ComponentType && target.ResourceType == obs[i].ResourceType {
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
		nodeName := fmt.Sprintf("%s_%d", graph.Name(), i)
		node, _ := graph.CreateNode(nodeName)
		node.SetLabel(componentString[o.ComponentType])
		node.SetShape(cgraph.BoxShape)
		node.SetStyle(cgraph.RoundedNodeStyle)
		nodes = append(nodes, node)
		var edge *cgraph.Edge
		edgeName := fmt.Sprintf("%s_%d", graph.Name(), i)
		if i > 0 {
			edge, _ = graph.CreateEdge(edgeName, nodes[i-1], nodes[i])
		} else {
			edge, _ = graph.CreateEdge(edgeName, graph.FirstNode(), nodes[i])
		}
		edge.SetDir(dir)
		actualObj := getObservation(&o, obs)
		if actualObj != nil && len(actualObj.NodeUUID) > 0 {
			graph.SetLabel(actualObj.NodeUUID)
		}
		switch {
		case actualObj == nil:
			edge.SetStyle("invis")
		case actualObj.ComponentType == v1.DFW && actualObj.ResourceType == v1.DROPPED:
			node.SetColor("red")
			edge.SetStyle("invis")
			node.SetLabel(componentString[actualObj.ComponentType]+"\n"+resourceString[actualObj.ResourceType])
		default:
			node.SetColor("blue")
			edge.SetColor("blue")
			node.SetLabel(componentString[actualObj.ComponentType]+"\n"+resourceString[actualObj.ResourceType])
		}
	}
	return nodes
}

var expectedSenderObservations = []v1.Observation{
	{
		ComponentType: v1.SPOOFGUARD,
		ResourceType:  v1.FORWARDED,
	},
	{
		ComponentType: v1.DFW,
		ResourceType:  v1.FORWARDED,
	},
	{
		ComponentType: v1.FORWARDING,
		ResourceType:  v1.FORWARDED,
	},
}

var expectedReceiverObservations = []v1.Observation{
	{
		ComponentType: v1.FORWARDING,
		ResourceType:  v1.RECEIVED,
	},
	{
		ComponentType: v1.DFW,
		ResourceType:  v1.FORWARDED,
	},
	{
		ComponentType: v1.FORWARDING,
		ResourceType:  v1.DELIVERED,
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
	cluster1.CreateNode(tf.SrcNamespace + "/" + tf.SrcPod)
	nodes1 := genSubGraph(cluster1, expectedSenderObservations, tf.Status.NodeSender, cgraph.ForwardDir)

	cluster2 := graph.SubGraph("cluster_dest", 1)
	cluster2.SetStyle(cgraph.FilledGraphStyle)
	cluster2.CreateNode(tf.DstNamespace + "/" + tf.DstPod)
	nodes2 := genSubGraph(cluster2, expectedReceiverObservations, tf.Status.NodeReceiver, cgraph.BackDir)

	edge, _ := graph.CreateEdge("cross_node", nodes1[len(nodes1)-1], nodes2[len(nodes2)-1])
	edge.SetConstraint(false)
	switch {
	case len(tf.Status.NodeReceiver) == 0 && len(tf.Status.NodeSender) == len(expectedSenderObservations):
		edge.SetColor("red")
		edge.SetStyle(cgraph.DashedEdgeStyle)
	case len(tf.Status.NodeSender) < len(expectedSenderObservations):
		edge.SetStyle("invis")
	default:
		edge.SetColor("blue")
	}

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
