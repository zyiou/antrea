package main

import (
	v1 "github.com/vmware-tanzu/antrea/pkg/apis/traceflow/v1"
	"github.com/vmware-tanzu/antrea/pkg/graphviz"
	"k8s.io/component-base/logs"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	GenGraphTest()
}

func GenGraphTest() string {
	tf := newTestSuccess()
	return graphviz.GenGraph(tf)
}

func newTestSuccess() *v1.Traceflow {
	tf := &v1.Traceflow{
		TypeMeta:     metav1.TypeMeta{},
		ObjectMeta:   metav1.ObjectMeta{},
		SrcNamespace: "ns1",
		SrcPod:       "client",
		DstNamespace: "ns2",
		DstPod:       "server",
		DstService:   "",
		RoundID:      "",
		Packet:       v1.Packet{},
		Status:       v1.Status{},
	}
	ob1 := v1.Observation{
		ComponentType: v1.SPOOFGUARD,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond(),
		NodeUUID:      "node A",
	}
	ob2 := v1.Observation{
		ComponentType: v1.DFW,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond() + 1,
	}
	ob3 := v1.Observation{
		ComponentType: v1.FORWARDING,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond() + 2,
	}
	ob4 := v1.Observation{
		ComponentType: v1.FORWARDING,
		ResourceType: v1.RECEIVED,
		Timestamp:     time.Now().Nanosecond() + 3,
		NodeUUID:      "node B",
	}
	ob5 := v1.Observation{
		ComponentType: v1.DFW,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond() + 4,
	}
	ob6 := v1.Observation{
		ComponentType: v1.FORWARDING,
		Timestamp:     time.Now().Nanosecond() + 5,
	}
	tf.Status.NodeSender = append(tf.Status.NodeSender, ob1, ob2, ob3)
	tf.Status.NodeReceiver = append(tf.Status.NodeReceiver, ob4, ob5, ob6)
	return tf
}

func newTestTraceflowDisconnect() *v1.Traceflow {
	tf := &v1.Traceflow{
		TypeMeta:     metav1.TypeMeta{},
		ObjectMeta:   metav1.ObjectMeta{},
		SrcNamespace: "ns1",
		SrcPod:       "client",
		DstNamespace: "ns2",
		DstPod:       "server",
		DstService:   "",
		RoundID:      "",
		Packet:       v1.Packet{},
		Status:       v1.Status{},
	}
	ob1 := v1.Observation{
		ComponentType: v1.SPOOFGUARD,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond(),
		NodeUUID:      "node A",
	}
	ob2 := v1.Observation{
		ComponentType: v1.DFW,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond() + 1,
	}
	ob3 := v1.Observation{
		ComponentType: v1.FORWARDING,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond() + 2,
	}
	tf.Status.NodeSender = append(tf.Status.NodeSender, ob1, ob2, ob3)
	return tf
}

func newTestNetPol() *v1.Traceflow {
	tf := &v1.Traceflow{
		TypeMeta:     metav1.TypeMeta{},
		ObjectMeta:   metav1.ObjectMeta{},
		SrcNamespace: "ns1",
		SrcPod:       "client",
		DstNamespace: "ns2",
		DstPod:       "server",
		DstService:   "",
		RoundID:      "",
		Packet:       v1.Packet{},
		Status:       v1.Status{},
	}
	ob1 := v1.Observation{
		ComponentType: v1.SPOOFGUARD,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond(),
		NodeUUID:      "node A",
	}
	ob2 := v1.Observation{
		ComponentType: v1.DFW,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond() + 1,
	}
	ob3 := v1.Observation{
		ComponentType: v1.FORWARDING,
		ResourceType: v1.FORWARDED,
		Timestamp:     time.Now().Nanosecond() + 2,
	}
	ob4 := v1.Observation{
		ComponentType: v1.FORWARDING,
		ResourceType: v1.RECEIVED,
		Timestamp:     time.Now().Nanosecond() + 3,
		NodeUUID:      "node B",
	}
	ob5 := v1.Observation{
		ComponentType: v1.DFW,
		ResourceType: v1.DROPPED,
		Timestamp:     time.Now().Nanosecond() + 4,
	}
	//ob6 := v1.Observation{
	//	ComponentType: v1.FORWARDING,
	//	Timestamp:     time.Now().Nanosecond() + 5,
	//}
	tf.Status.NodeSender = append(tf.Status.NodeSender, ob1, ob2, ob3)
	tf.Status.NodeReceiver = append(tf.Status.NodeReceiver, ob4, ob5)
	return tf
}