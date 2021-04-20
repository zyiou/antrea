// Copyright 2021 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package denyconnections

import (
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/vmware-tanzu/antrea/pkg/agent/flowexporter"
)

var tuple = flowexporter.Tuple{
	SourceAddress:      net.IP{1, 2, 3, 4},
	DestinationAddress: net.IP{5, 6, 7, 8},
	Protocol:           6,
	SourcePort:         30000,
	DestinationPort:    80,
}

var refTime = time.Now()
var denyConn = &flowexporter.DenyConnection{
	FlowKey:                        tuple,
	Bytes:                          uint64(60),
	IngressNetworkPolicyRuleAction: "Reject",
	IsIPv6:                         false,
	TimeSeen:                       refTime,
}

func TestDenyConnectionStore_AddOrUpdateConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Create two flows with same tuple, add them into denyConnectionStore.
	refTime := time.Now()
	denyConn1 := denyConn
	denyConn2 := &flowexporter.DenyConnection{
		FlowKey:                        tuple,
		Bytes:                          uint64(60),
		IngressNetworkPolicyRuleAction: "Reject",
		IsIPv6:                         false,
		TimeSeen:                       refTime.Add(+(time.Second)),
	}
	connKey := flowexporter.GetConnKey(tuple)
	denyConnStore := NewDenyConnectionStore()
	denyConnStore.AddOrUpdateConnection(denyConn1)
	assert.Equal(t, 1, len(denyConnStore.connections))
	assert.NotNil(t, denyConnStore.connections[connKey])
	assert.Equal(t, uint64(60), denyConnStore.connections[connKey].Bytes)
	denyConnStore.AddOrUpdateConnection(denyConn2)
	assert.Equal(t, 1, len(denyConnStore.connections))
	assert.NotNil(t, denyConnStore.connections[connKey])
	assert.Equal(t, uint64(120), denyConnStore.connections[connKey].Bytes)
}

func TestDenyConnectionStore_DeleteConnWithoutLock(t *testing.T) {
	connKey := flowexporter.GetConnKey(tuple)
	denyConnStore := NewDenyConnectionStore()
	denyConnStore.connections[connKey] = denyConn
	assert.Equal(t, 1, len(denyConnStore.connections))
	denyConnStore.DeleteConnWithoutLock(connKey)
	assert.Equal(t, 0, len(denyConnStore.connections))
}

func TestDenyConnectionStore_GetNumOfConnections(t *testing.T) {
	connKey := flowexporter.GetConnKey(tuple)
	denyConnStore := NewDenyConnectionStore()
	assert.Equal(t, 0, denyConnStore.GetNumOfConnections())
	denyConnStore.connections[connKey] = denyConn
	assert.Equal(t, 1, denyConnStore.GetNumOfConnections())
}

func TestDenyConnectionStore_ForAllConnectionsDo(t *testing.T) {
	connKey := flowexporter.GetConnKey(tuple)
	denyConnStore := NewDenyConnectionStore()
	denyConnStore.connections[connKey] = denyConn

	resetPacketLength := func(key flowexporter.ConnectionKey, conn *flowexporter.DenyConnection) error {
		conn.Bytes = 0
		return nil
	}
	denyConnStore.ForAllConnectionsDo(resetPacketLength)
	assert.Equal(t, 1, len(denyConnStore.connections))
	conn, exist := denyConnStore.connections[connKey]
	assert.True(t, exist)
	assert.Equal(t, uint64(0), conn.Bytes, "packeLength should be reset.")
}
