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
//

// Code generated by MockGen. DO NOT EDIT.
// Source: antrea.io/antrea/pkg/ipfix (interfaces: IPFIXExportingProcess,IPFIXRegistry,IPFIXCollectingProcess,IPFIXAggregationProcess)

// Package testing is a generated GoMock package.
package testing

import (
	gomock "github.com/golang/mock/gomock"
	entities "github.com/vmware/go-ipfix/pkg/entities"
	intermediate "github.com/vmware/go-ipfix/pkg/intermediate"
	reflect "reflect"
	time "time"
)

// MockIPFIXExportingProcess is a mock of IPFIXExportingProcess interface
type MockIPFIXExportingProcess struct {
	ctrl     *gomock.Controller
	recorder *MockIPFIXExportingProcessMockRecorder
}

// MockIPFIXExportingProcessMockRecorder is the mock recorder for MockIPFIXExportingProcess
type MockIPFIXExportingProcessMockRecorder struct {
	mock *MockIPFIXExportingProcess
}

// NewMockIPFIXExportingProcess creates a new mock instance
func NewMockIPFIXExportingProcess(ctrl *gomock.Controller) *MockIPFIXExportingProcess {
	mock := &MockIPFIXExportingProcess{ctrl: ctrl}
	mock.recorder = &MockIPFIXExportingProcessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPFIXExportingProcess) EXPECT() *MockIPFIXExportingProcessMockRecorder {
	return m.recorder
}

// CloseConnToCollector mocks base method
func (m *MockIPFIXExportingProcess) CloseConnToCollector() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseConnToCollector")
}

// CloseConnToCollector indicates an expected call of CloseConnToCollector
func (mr *MockIPFIXExportingProcessMockRecorder) CloseConnToCollector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseConnToCollector", reflect.TypeOf((*MockIPFIXExportingProcess)(nil).CloseConnToCollector))
}

// NewTemplateID mocks base method
func (m *MockIPFIXExportingProcess) NewTemplateID() uint16 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTemplateID")
	ret0, _ := ret[0].(uint16)
	return ret0
}

// NewTemplateID indicates an expected call of NewTemplateID
func (mr *MockIPFIXExportingProcessMockRecorder) NewTemplateID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTemplateID", reflect.TypeOf((*MockIPFIXExportingProcess)(nil).NewTemplateID))
}

// SendSet mocks base method
func (m *MockIPFIXExportingProcess) SendSet(arg0 entities.Set) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSet", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendSet indicates an expected call of SendSet
func (mr *MockIPFIXExportingProcessMockRecorder) SendSet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSet", reflect.TypeOf((*MockIPFIXExportingProcess)(nil).SendSet), arg0)
}

// MockIPFIXRegistry is a mock of IPFIXRegistry interface
type MockIPFIXRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockIPFIXRegistryMockRecorder
}

// MockIPFIXRegistryMockRecorder is the mock recorder for MockIPFIXRegistry
type MockIPFIXRegistryMockRecorder struct {
	mock *MockIPFIXRegistry
}

// NewMockIPFIXRegistry creates a new mock instance
func NewMockIPFIXRegistry(ctrl *gomock.Controller) *MockIPFIXRegistry {
	mock := &MockIPFIXRegistry{ctrl: ctrl}
	mock.recorder = &MockIPFIXRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPFIXRegistry) EXPECT() *MockIPFIXRegistryMockRecorder {
	return m.recorder
}

// GetInfoElement mocks base method
func (m *MockIPFIXRegistry) GetInfoElement(arg0 string, arg1 uint32) (*entities.InfoElement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfoElement", arg0, arg1)
	ret0, _ := ret[0].(*entities.InfoElement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfoElement indicates an expected call of GetInfoElement
func (mr *MockIPFIXRegistryMockRecorder) GetInfoElement(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfoElement", reflect.TypeOf((*MockIPFIXRegistry)(nil).GetInfoElement), arg0, arg1)
}

// LoadRegistry mocks base method
func (m *MockIPFIXRegistry) LoadRegistry() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoadRegistry")
}

// LoadRegistry indicates an expected call of LoadRegistry
func (mr *MockIPFIXRegistryMockRecorder) LoadRegistry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadRegistry", reflect.TypeOf((*MockIPFIXRegistry)(nil).LoadRegistry))
}

// MockIPFIXCollectingProcess is a mock of IPFIXCollectingProcess interface
type MockIPFIXCollectingProcess struct {
	ctrl     *gomock.Controller
	recorder *MockIPFIXCollectingProcessMockRecorder
}

// MockIPFIXCollectingProcessMockRecorder is the mock recorder for MockIPFIXCollectingProcess
type MockIPFIXCollectingProcessMockRecorder struct {
	mock *MockIPFIXCollectingProcess
}

// NewMockIPFIXCollectingProcess creates a new mock instance
func NewMockIPFIXCollectingProcess(ctrl *gomock.Controller) *MockIPFIXCollectingProcess {
	mock := &MockIPFIXCollectingProcess{ctrl: ctrl}
	mock.recorder = &MockIPFIXCollectingProcessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPFIXCollectingProcess) EXPECT() *MockIPFIXCollectingProcessMockRecorder {
	return m.recorder
}

// GetMsgChan mocks base method
func (m *MockIPFIXCollectingProcess) GetMsgChan() chan *entities.Message {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMsgChan")
	ret0, _ := ret[0].(chan *entities.Message)
	return ret0
}

// GetMsgChan indicates an expected call of GetMsgChan
func (mr *MockIPFIXCollectingProcessMockRecorder) GetMsgChan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMsgChan", reflect.TypeOf((*MockIPFIXCollectingProcess)(nil).GetMsgChan))
}

// Start mocks base method
func (m *MockIPFIXCollectingProcess) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockIPFIXCollectingProcessMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockIPFIXCollectingProcess)(nil).Start))
}

// Stop mocks base method
func (m *MockIPFIXCollectingProcess) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockIPFIXCollectingProcessMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockIPFIXCollectingProcess)(nil).Stop))
}

// MockIPFIXAggregationProcess is a mock of IPFIXAggregationProcess interface
type MockIPFIXAggregationProcess struct {
	ctrl     *gomock.Controller
	recorder *MockIPFIXAggregationProcessMockRecorder
}

// MockIPFIXAggregationProcessMockRecorder is the mock recorder for MockIPFIXAggregationProcess
type MockIPFIXAggregationProcessMockRecorder struct {
	mock *MockIPFIXAggregationProcess
}

// NewMockIPFIXAggregationProcess creates a new mock instance
func NewMockIPFIXAggregationProcess(ctrl *gomock.Controller) *MockIPFIXAggregationProcess {
	mock := &MockIPFIXAggregationProcess{ctrl: ctrl}
	mock.recorder = &MockIPFIXAggregationProcessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPFIXAggregationProcess) EXPECT() *MockIPFIXAggregationProcessMockRecorder {
	return m.recorder
}

// ForAllExpiredFlowRecordsDo mocks base method
func (m *MockIPFIXAggregationProcess) ForAllExpiredFlowRecordsDo(arg0 intermediate.FlowKeyRecordMapCallBack) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForAllExpiredFlowRecordsDo", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForAllExpiredFlowRecordsDo indicates an expected call of ForAllExpiredFlowRecordsDo
func (mr *MockIPFIXAggregationProcessMockRecorder) ForAllExpiredFlowRecordsDo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForAllExpiredFlowRecordsDo", reflect.TypeOf((*MockIPFIXAggregationProcess)(nil).ForAllExpiredFlowRecordsDo), arg0)
}

// GetExpiryFromExpirePriorityQueue mocks base method
func (m *MockIPFIXAggregationProcess) GetExpiryFromExpirePriorityQueue() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpiryFromExpirePriorityQueue")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetExpiryFromExpirePriorityQueue indicates an expected call of GetExpiryFromExpirePriorityQueue
func (mr *MockIPFIXAggregationProcessMockRecorder) GetExpiryFromExpirePriorityQueue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpiryFromExpirePriorityQueue", reflect.TypeOf((*MockIPFIXAggregationProcess)(nil).GetExpiryFromExpirePriorityQueue))
}

// IsMetadataFilled mocks base method
func (m *MockIPFIXAggregationProcess) IsMetadataFilled(arg0 intermediate.AggregationFlowRecord) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMetadataFilled", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsMetadataFilled indicates an expected call of IsMetadataFilled
func (mr *MockIPFIXAggregationProcessMockRecorder) IsMetadataFilled(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMetadataFilled", reflect.TypeOf((*MockIPFIXAggregationProcess)(nil).IsMetadataFilled), arg0)
}

// ResetStatElementsInRecord mocks base method
func (m *MockIPFIXAggregationProcess) ResetStatElementsInRecord(arg0 entities.Record) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetStatElementsInRecord", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetStatElementsInRecord indicates an expected call of ResetStatElementsInRecord
func (mr *MockIPFIXAggregationProcessMockRecorder) ResetStatElementsInRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetStatElementsInRecord", reflect.TypeOf((*MockIPFIXAggregationProcess)(nil).ResetStatElementsInRecord), arg0)
}

// SetMetadataFilled mocks base method
func (m *MockIPFIXAggregationProcess) SetMetadataFilled(arg0 intermediate.AggregationFlowRecord) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMetadataFilled", arg0)
}

// SetMetadataFilled indicates an expected call of SetMetadataFilled
func (mr *MockIPFIXAggregationProcessMockRecorder) SetMetadataFilled(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMetadataFilled", reflect.TypeOf((*MockIPFIXAggregationProcess)(nil).SetMetadataFilled), arg0)
}

// Start mocks base method
func (m *MockIPFIXAggregationProcess) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockIPFIXAggregationProcessMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockIPFIXAggregationProcess)(nil).Start))
}

// Stop mocks base method
func (m *MockIPFIXAggregationProcess) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockIPFIXAggregationProcessMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockIPFIXAggregationProcess)(nil).Stop))
}
