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
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"

	"github.com/vmware-tanzu/antrea/pkg/agent/flowexporter"
	"github.com/vmware-tanzu/antrea/pkg/agent/interfacestore"
	"github.com/vmware-tanzu/antrea/pkg/agent/metrics"
	"github.com/vmware-tanzu/antrea/pkg/agent/proxy"
)

type ConnectionStore interface {
	// GetConnByKey gets the connection in connection map given the connection key.
	GetConnByKey(connKey flowexporter.ConnectionKey) (*flowexporter.Connection, bool)
	// ForAllConnectionsDo execute the callback for each connection in connection map.
	ForAllConnectionsDo(callback flowexporter.ConnectionMapCallBack) error
	// DeleteConnWithoutLock deletes the connection from the connection map given
	// the connection key without grabbing the lock. Caller is expected to grab lock.
	DeleteConnWithoutLock(connKey flowexporter.ConnectionKey) error
	// AddOrUpdateConn updates the connection if it is already present, i.e., update timestamp, counters etc.,
	// or adds a new connection with the resolved K8s metadata.
	AddOrUpdateConn(conn *flowexporter.Connection)
}

type connectionStore struct {
	connections   map[flowexporter.ConnectionKey]*flowexporter.Connection
	ifaceStore    interfacestore.InterfaceStore
	antreaProxier proxy.Proxier
	mutex         sync.Mutex
	isDenyConn    bool
}

func NewConnectionStore(
	ifaceStore interfacestore.InterfaceStore,
	proxier proxy.Proxier,
	isDenyConn bool,
) *connectionStore {
	return &connectionStore{
		connections:   make(map[flowexporter.ConnectionKey]*flowexporter.Connection),
		ifaceStore:    ifaceStore,
		antreaProxier: proxier,
		isDenyConn:    isDenyConn,
	}
}

// GetConnByKey gets the connection in connection map given the connection key.
func (cs *connectionStore) GetConnByKey(connKey flowexporter.ConnectionKey) (*flowexporter.Connection, bool) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	conn, found := cs.connections[connKey]
	return conn, found
}

// ForAllConnectionsDo execute the callback for each connection in connection map.
func (cs *connectionStore) ForAllConnectionsDo(callback flowexporter.ConnectionMapCallBack) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	for k, v := range cs.connections {
		err := callback(k, v)
		if err != nil {
			klog.Errorf("Callback execution failed for flow with key: %v, conn: %v, k, v: %v", k, v, err)
			return err
		}
	}
	return nil
}

// DeleteConnWithoutLock deletes the connection from the connection map given
// the connection key without grabbing the lock. Caller is expected to grab lock.
func (cs *connectionStore) DeleteConnWithoutLock(connKey flowexporter.ConnectionKey) error {
	_, exists := cs.connections[connKey]
	if !exists {
		return fmt.Errorf("connection with key %v doesn't exist in map", connKey)
	}
	delete(cs.connections, connKey)
	if !cs.isDenyConn {
		metrics.TotalAntreaConnectionsInConnTrackTable.Dec()
	}
	return nil
}

func (cs *connectionStore) fillPodInfo(conn *flowexporter.Connection) {
	if cs.ifaceStore == nil {
		klog.Warning("Interface store is not available to retrieve local Pods information.")
		return
	}
	// sourceIP/destinationIP are mapped only to local pods and not remote pods.
	sIface, srcFound := cs.ifaceStore.GetInterfaceByIP(conn.FlowKey.SourceAddress.String())
	dIface, dstFound := cs.ifaceStore.GetInterfaceByIP(conn.FlowKey.DestinationAddress.String())
	if !srcFound && !dstFound {
		klog.Warningf("Cannot map any of the IP %s or %s to a local Pod", conn.FlowKey.SourceAddress.String(), conn.FlowKey.DestinationAddress.String())
	}
	if srcFound && sIface.Type == interfacestore.ContainerInterface {
		conn.SourcePodName = sIface.ContainerInterfaceConfig.PodName
		conn.SourcePodNamespace = sIface.ContainerInterfaceConfig.PodNamespace
	}
	if dstFound && dIface.Type == interfacestore.ContainerInterface {
		conn.DestinationPodName = dIface.ContainerInterfaceConfig.PodName
		conn.DestinationPodNamespace = dIface.ContainerInterfaceConfig.PodNamespace
	}
}

func (cs *connectionStore) fillServiceInfo(conn *flowexporter.Connection, serviceStr string) {
	// resolve destination Service information
	if cs.antreaProxier != nil {
		servicePortName, exists := cs.antreaProxier.GetServiceByIP(serviceStr)
		if !exists {
			klog.Warningf("Could not retrieve the Service info from antrea-agent-proxier for the serviceStr: %s", serviceStr)
		} else {
			conn.DestinationServicePortName = servicePortName.String()
		}
	}
}

// LookupServiceProtocol returns the corresponding Service protocol string for a given protocol identifier
func lookupServiceProtocol(protoID uint8) (corev1.Protocol, error) {
	serviceProto, found := serviceProtocolMap[protoID]
	if !found {
		return "", fmt.Errorf("unknown protocol identifier: %d", protoID)
	}
	return serviceProto, nil
}