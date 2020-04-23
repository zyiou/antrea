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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Phase uint
type ComponentType uint
type ResourceType uint
type DropReason uint
type PacketResouceType uint
type PacketTransportType uint

const (
	INITIAL Phase = iota
	RUNNING
	SUCCESS
	TIMEOUT
	ERROR
)

const (
	SPOOFGUARD ComponentType = iota
	LB
	ROUTING
	DFW
	FORWARDING
)

const (
	DELIVERED ResourceType = iota
	RECEIVED
	FORWARDED
	DROPPED
)

const (
	FieldsPacketData PacketResouceType = iota
)

const (
	UNICAST PacketTransportType = iota
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Traceflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	SrcNamespace string `json:"srcNamespace,omitempty"`
	SrcPod       string `json:"srcPod,omitempty"`
	DstNamespace string `json:"dstNamespace,omitempty"`
	DstPod       string `json:"dstPod,omitempty"`
	DstService   string `json:"dstService,omitempty"`
	RoundID      string `json:"roundID,omitempty"`

	Packet `json:",inline"`
	Status `json:",inline"`
}

type IPHeader struct {
	SrcIP    string `json:"srcIP,omitempty"`
	DstIP    string `json:"dstIP,omitempty"`
	Protocol int    `json:"protocol,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
	Flags    int    `json:"flags,omitempty"`
}

type IPv6Header struct {
	SrcIP      string `json:"srcIPv6IP,omitempty"`
	DstIP      string `json:"dstIPv6IP,omitempty"`
	NextHeader int    `json:"nextHeader,omitempty"`
	HopLimit   int    `json:"hopLimit,omitempty"`
}

type TransportHeader struct {
	ICMPEchoRequestHeader `json:",inline"`
	UDPHeader             `json:",inline"`
	TCPHeader             `json:",inline"`
}

type ICMPEchoRequestHeader struct {
	ID       int `json:"id,omitempty"`
	Sequence int `json:"sequence,omitempty"`
}

type UDPHeader struct {
	SrcPort int `json:"srcUDPPort,omitempty"`
	DstPort int `json:"dstUDPPort,omitempty"`
}

type TCPHeader struct {
	SrcPort  int `json:"srcTCPPort,omitempty"`
	DstPort  int `json:"dstTCPPort,omitempty"`
	TCPFlags int `json:"tcpFlags,omitempty"`
}

type Packet struct {
	ResouceType   PacketResouceType   `json:"resouceType,omitempty"`
	FrameSize     int                 `json:"frameSize,omitempty"`
	TransportType PacketTransportType `json:"transportType,omitempty"`
	Payload       string              `json:"payload,omitempty"`

	IPHeader        `json:",inline"`
	IPv6Header      `json:",inline"`
	TransportHeader `json:",inline"`
}

type Status struct {
	Phase        Phase         `json:"phase,omitempty"`
	CrossNodeTag uint8         `json:"crossNodeTag,omitempty"`
	NodeSender   []Observation `json:"nodeSender,omitempty"`
	NodeReceiver []Observation `json:"nodeReceiver,omitempty"`
}

type Observation struct {
	ComponentType    ComponentType `json:"componentType,omitempty"`
	ComponentSubType string        `json:"componentSubType,omitempty"`
	ComponentName    string        `json:"componentName,omitempty"`
	ResourceType     ResourceType  `json:"resourceType,omitempty"`
	RoundID          string        `json:"roundID,omitempty"`
	NodeUUID         string        `json:"nodeUUID,omitempty"`
	PodUUID          string        `json:"podUUID,omitempty"`
	DstMAC           string        `json:"dstMAC,omitempty"`
	RuleID           string        `json:"ruleID,omitempty"`
	Rule             string        `json:"rule,omitempty"`
	TTL              int           `json:"ttl,omitempty"`
	TranslatedSrcIP  string        `json:"translatedSrcIP,omitempty"`
	TranslatedDstIP  string        `json:"translatedDstIP,omitempty"`
	DropReason       DropReason    `json:"dropReason,omitempty"`
	Timestamp        int           `json:"timestamp,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TraceflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Traceflow `json:"items"`
}
