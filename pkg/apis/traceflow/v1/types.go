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

	SrcNamespace string
	SrcPod       string
	DstNamespace string
	DstPod       string
	DstService   string
	RoundID      string

	Packet
	Status
}

type IPHeader struct {
	SrcIP    string
	DstIP    string
	Protocol int
	TTL      int
	Flags    int
}

type IPv6Header struct {
	SrcIP      string
	DstIP      string
	NextHeader int
	HopLimit   int
}

type TransportHeader struct {
	ICMPEchoRequestHeader
	UDPHeader
	TCPHeader
}

type ICMPEchoRequestHeader struct {
	ID       int
	Sequence int
}

type UDPHeader struct {
	SrcPort int
	DstPort int
}

type TCPHeader struct {
	SrcPort  int
	DstPort  int
	TCPFlags int
}

type Packet struct {
	ResouceType   PacketResouceType
	FrameSize     int
	TransportType PacketTransportType
	Payload       string

	IPHeader
	IPv6Header
	TransportHeader
}

type Status struct {
	Phase        Phase
	CrossNodeTag uint8
	NodeSender   []Observation
	NodeReceiver []Observation
}

type Observation struct {
	ComponentType    ComponentType
	ComponentSubType string
	ComponentName    string
	ResourceType     ResourceType
	RoundID          string
	NodeUUID         string
	PodUUID          string
	DstMAC           string
	RuleID           string
	Rule             string
	TTL              int
	TranslatedSrcIP  string
	TranslatedDstIP  string
	DropReason       DropReason
	Timestamp        int
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TraceflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Traceflow
}
