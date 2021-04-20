// Copyright 2020 Antrea Authors
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

package connections

import (
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"antrea.io/antrea/pkg/agent/flowexporter"
	interfacestoretest "antrea.io/antrea/pkg/agent/interfacestore/testing"
)

const testPollInterval = 0 // Not used in these tests, hence 0.

func makeTuple(srcIP *net.IP, dstIP *net.IP, protoID uint8, srcPort uint16, dstPort uint16) (flowexporter.Tuple, flowexporter.Tuple) {
	tuple := flowexporter.Tuple{
		SourceAddress:      *srcIP,
		DestinationAddress: *dstIP,
		Protocol:           protoID,
		SourcePort:         srcPort,
		DestinationPort:    dstPort,
	}
	revTuple := flowexporter.Tuple{
		SourceAddress:      *dstIP,
		DestinationAddress: *srcIP,
		Protocol:           protoID,
		SourcePort:         dstPort,
		DestinationPort:    srcPort,
	}
	return tuple, revTuple
}

func TestConnectionStore_ForAllConnectionsDo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Create two flows; one is already in connectionStore and other one is new
	testFlows := make([]*flowexporter.Connection, 2)
	testFlowKeys := make([]*flowexporter.ConnectionKey, 2)
	refTime := time.Now()
	// Flow-1, which is already in connectionStore
	tuple1, revTuple1 := makeTuple(&net.IP{1, 2, 3, 4}, &net.IP{4, 3, 2, 1}, 6, 65280, 255)
	testFlows[0] = &flowexporter.Connection{
		StartTime:       refTime.Add(-(time.Second * 50)),
		StopTime:        refTime,
		OriginalPackets: 0xffff,
		OriginalBytes:   0xbaaaaa0000000000,
		ReversePackets:  0xff,
		ReverseBytes:    0xbaaa,
		TupleOrig:       tuple1,
		TupleReply:      revTuple1,
		IsPresent:       true,
	}
	// Flow-2, which is not in connectionStore
	tuple2, revTuple2 := makeTuple(&net.IP{5, 6, 7, 8}, &net.IP{8, 7, 6, 5}, 6, 60001, 200)
	testFlows[1] = &flowexporter.Connection{
		StartTime:       refTime.Add(-(time.Second * 20)),
		StopTime:        refTime,
		OriginalPackets: 0xbb,
		OriginalBytes:   0xcbbb,
		ReversePackets:  0xbbbb,
		ReverseBytes:    0xcbbbb0000000000,
		TupleOrig:       tuple2,
		TupleReply:      revTuple2,
		IsPresent:       true,
	}
	for i, flow := range testFlows {
		connKey := flowexporter.NewConnectionKey(flow)
		testFlowKeys[i] = &connKey
	}
	// Create connectionStore
	mockIfaceStore := interfacestoretest.NewMockInterfaceStore(ctrl)
	connStore := NewConnectionStore(mockIfaceStore, nil)
	// Add flows to the Connection store
	for i, flow := range testFlows {
		connStore.connections[*testFlowKeys[i]] = flow
	}

	resetTwoFields := func(key flowexporter.ConnectionKey, conn *flowexporter.Connection) error {
		conn.IsPresent = false
		conn.OriginalPackets = 0
		return nil
	}
	connStore.ForAllConnectionsDo(resetTwoFields)
	// Check isActive and OriginalPackets, if they are reset or not.
	for i := 0; i < len(testFlows); i++ {
		conn, ok := connStore.GetConnByKey(*testFlowKeys[i])
		assert.Equal(t, ok, true, "connection should be there in connection store")
		assert.Equal(t, conn.IsPresent, false, "isActive flag should be reset")
		assert.Equal(t, conn.OriginalPackets, uint64(0), "OriginalPackets should be reset")
	}
}
