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

	"k8s.io/klog"

	"antrea.io/antrea/pkg/agent/flowexporter"
	"antrea.io/antrea/pkg/agent/interfacestore"
	"antrea.io/antrea/pkg/agent/proxy"
	"antrea.io/antrea/pkg/util/ip"
)

type DenyConnectionStore struct {
	*connectionStore
}

func NewDenyConnectionStore(ifaceStore interfacestore.InterfaceStore,
	proxier proxy.Proxier) *DenyConnectionStore {
	return &DenyConnectionStore{
		connectionStore: NewConnectionStore(ifaceStore, proxier, true),
	}
}

// AddOrUpdateConn updates the connection if it is already present, i.e., update timestamp, counters etc.,
// or adds a new connection with the resolved K8s metadata.
func (ds *DenyConnectionStore) AddOrUpdateConn(conn *flowexporter.Connection) {
	connKey := flowexporter.NewConnectionKey(conn)
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	existingConn, exist := ds.connections[connKey]
	if exist {
		existingConn.Bytes += conn.Bytes
		klog.V(2).Infof("Deny connection with flowKey %v has been updated.", connKey)
		return
	}
	ds.fillPodInfo(conn)
	protocolStr := ip.IPProtocolNumberToString(conn.FlowKey.Protocol, "UnknownProtocol")
	serviceStr := fmt.Sprintf("%s:%d/%s", conn.FlowKey.DestinationAddress, conn.FlowKey.DestinationPort, protocolStr)
	ds.fillServiceInfo(conn, serviceStr)
	ds.connections[connKey] = conn
}
