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
	"fmt"
	"sync"

	"k8s.io/klog"

	"github.com/vmware-tanzu/antrea/pkg/agent/flowexporter"
)

type DenyConnectionStore interface {
	// AddOrUpdateConnection adds or updates connection into deny connection map.
	AddOrUpdateConnection(conn *flowexporter.DenyConnection)
	// ForAllConnectionsDo executes the callback for each connection in deny connection map.
	ForAllConnectionsDo(callback flowexporter.DenyConnectionMapCallBack) error
	// ContainsConnection checks whether connection exists in deny connection store.
	ContainsConnection(flowKey flowexporter.Tuple) bool
	// DeleteConnWithoutLock deletes the deny connection from the connection map given
	// the connection key without grabbing the lock. Caller is expected to grab lock.
	DeleteConnWithoutLock(connKey flowexporter.ConnectionKey) error
	// GetNumOfConnections returns number of connections in deny connection map.
	GetNumOfConnections() int
}

type denyConnectionStore struct {
	connections map[flowexporter.ConnectionKey]*flowexporter.DenyConnection
	mutex       sync.Mutex
}

func NewDenyConnectionStore() *denyConnectionStore {
	return &denyConnectionStore{
		connections: make(map[flowexporter.ConnectionKey]*flowexporter.DenyConnection),
	}
}

// AddOrUpdateConnection adds or updates connection into deny connection map.
func (ds *denyConnectionStore) AddOrUpdateConnection(conn *flowexporter.DenyConnection) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	connKey := flowexporter.GetConnKey(conn.FlowKey)
	existingConn, exist := ds.connections[connKey]
	if exist {
		existingConn.Bytes += conn.Bytes
		klog.V(2).Infof("Deny connection with flowKey %v has been updated.", connKey)
		return
	}
	ds.connections[connKey] = conn
}

// ContainsConnection checks whether connection exists in deny connection store.
func (ds *denyConnectionStore) ContainsConnection(flowKey flowexporter.Tuple) bool {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	connKey := flowexporter.GetConnKey(flowKey)
	_, exist := ds.connections[connKey]
	return exist
}

// ForAllConnectionsDo execute the callback for each connection in deny connection map.
func (ds *denyConnectionStore) ForAllConnectionsDo(callback flowexporter.DenyConnectionMapCallBack) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	for k, v := range ds.connections {
		err := callback(k, v)
		if err != nil {
			klog.Errorf("Callback execution failed for flow with key: %v, conn: %v, k, v: %v", k, v, err)
			return err
		}
	}
	return nil
}

// DeleteConnWithoutLock deletes the deny connection from the connection map given
// the connection key without grabbing the lock. Caller is expected to grab lock.
func (ds *denyConnectionStore) DeleteConnWithoutLock(connKey flowexporter.ConnectionKey) error {
	_, exists := ds.connections[connKey]
	if !exists {
		return fmt.Errorf("deny connection with key %v doesn't exist in map", connKey)
	}
	delete(ds.connections, connKey)
	return nil
}

// GetNumOfConnections returns number of connections in deny connection map.
func (ds *denyConnectionStore) GetNumOfConnections() int {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	return len(ds.connections)
}
