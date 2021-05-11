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

package connections

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/component-base/metrics/testutil"

	"antrea.io/antrea/pkg/agent/flowexporter"
	interfacestoretest "antrea.io/antrea/pkg/agent/interfacestore/testing"
	"antrea.io/antrea/pkg/agent/metrics"
	proxytest "antrea.io/antrea/pkg/agent/proxy/testing"
	k8sproxy "antrea.io/antrea/third_party/proxy"
)

func TestDenyConnectionStore_AddOrUpdateConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metrics.InitializeConnectionMetrics()
	// Create two flows for testing adding and updating of same connection.
	refTime := time.Now()
	tuple, _ := makeTuple(&net.IP{1, 2, 3, 4}, &net.IP{4, 3, 2, 1}, 6, 65280, 255)
	servicePortName := k8sproxy.ServicePortName{
		NamespacedName: types.NamespacedName{
			Namespace: "serviceNS1",
			Name:      "service1",
		},
		Port:     "255",
		Protocol: v1.ProtocolTCP,
	}
	// flow-1 for testing adding
	testFlow1 := flowexporter.Connection{
		TimeSeen: refTime.Add(-(time.Second * 20)),
		FlowKey:  tuple,
		IsIPv6:   false,
		Bytes:    uint64(60),
	}
	// flow-2 for testing updating
	testFlow2 := flowexporter.Connection{
		TimeSeen: refTime.Add(-(time.Second * 10)),
		FlowKey:  tuple,
		IsIPv6:   false,
		Bytes:    uint64(60),
	}
	mockIfaceStore := interfacestoretest.NewMockInterfaceStore(ctrl)
	mockProxier := proxytest.NewMockProxier(ctrl)
	protocol, _ := lookupServiceProtocol(tuple.Protocol)
	serviceStr := fmt.Sprintf("%s:%d/%s", tuple.DestinationAddress.String(), tuple.DestinationPort, protocol)
	mockProxier.EXPECT().GetServiceByIP(serviceStr).Return(servicePortName, true)
	mockIfaceStore.EXPECT().GetInterfaceByIP(tuple.SourceAddress.String()).Return(nil, false)
	mockIfaceStore.EXPECT().GetInterfaceByIP(tuple.DestinationAddress.String()).Return(nil, false)

	denyConnStore := NewDenyConnectionStore(mockIfaceStore, mockProxier)

	denyConnStore.AddOrUpdateConn(&testFlow1)
	expConn := testFlow1
	expConn.DestinationServicePortName = servicePortName.String()
	actualConn, ok := denyConnStore.GetConnByKey(flowexporter.NewConnectionKey(&testFlow1))
	assert.Equal(t, ok, true, "deny connection should be there in deny connection store")
	assert.Equal(t, expConn, *actualConn, "deny connections should be equal")
	checkDenyConnectionMetrics(t, len(denyConnStore.connections))

	denyConnStore.AddOrUpdateConn(&testFlow2)
	expConn.Bytes = uint64(120)
	actualConn, ok = denyConnStore.GetConnByKey(flowexporter.NewConnectionKey(&testFlow1))
	assert.Equal(t, ok, true, "deny connection should be there in deny connection store")
	assert.Equal(t, expConn, *actualConn, "deny connections should be equal")
	checkDenyConnectionMetrics(t, len(denyConnStore.connections))
}

func TestDenyConnectionStore_DeleteConnWithoutLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metrics.InitializeConnectionMetrics()
	// Create denyConnectionStore
	mockIfaceStore := interfacestoretest.NewMockInterfaceStore(ctrl)
	connStore := NewDenyConnectionStore(mockIfaceStore, nil)
	refTime := time.Now()
	tuple1, tuple2 := makeTuple(&net.IP{1, 2, 3, 4}, &net.IP{4, 3, 2, 1}, 6, 65280, 255)
	conn := &flowexporter.Connection{
		TimeSeen: refTime.Add(-(time.Second * 20)),
		FlowKey:  tuple1,
		IsIPv6:   false,
		Bytes:    uint64(60),
	}
	connKey := flowexporter.NewConnectionKey(conn)
	connStore.connections[connKey] = conn
	// create invalid connection key to be deleted
	conn.FlowKey = tuple2
	invalidConnKey := flowexporter.NewConnectionKey(conn)

	// For testing purposes, set the metric
	metrics.TotalDenyConnections.Set(1)

	err := connStore.DeleteConnWithoutLock(connKey)
	assert.Nil(t, err, "DeleteConnWithoutLock should return nil")
	_, exists := connStore.GetConnByKey(connKey)
	assert.Equal(t, exists, false, "connection should be deleted in connection store")
	checkDenyConnectionMetrics(t, len(connStore.connections))
	assert.NotNil(t, connStore.DeleteConnWithoutLock(invalidConnKey))
	checkDenyConnectionMetrics(t, len(connStore.connections))
}

func checkDenyConnectionMetrics(t *testing.T, numConns int) {
	expectedDenyConnectionCount := `
	# HELP antrea_agent_deny_connection_count [ALPHA] Number of deny connections detected by Flow Exporter deny connections tracking. This metric gets updated when a flow is rejected/dropped by network policy.
	# TYPE antrea_agent_deny_connection_count gauge
	`
	expectedDenyConnectionCount = expectedDenyConnectionCount + fmt.Sprintf("antrea_agent_deny_connection_count %d\n", numConns)
	err := testutil.GatherAndCompare(legacyregistry.DefaultGatherer, strings.NewReader(expectedDenyConnectionCount), "antrea_agent_deny_connection_count")
	assert.NoError(t, err)
}
