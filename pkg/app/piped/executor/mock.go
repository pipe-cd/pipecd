// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/app/piped/executor/executor.go

// Package executor is a generated GoMock package.
package executor

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	kubernetes "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	model "github.com/pipe-cd/pipe/pkg/model"
)

// MockExecutor is a mock of Executor interface
type MockExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockExecutorMockRecorder
}

// MockExecutorMockRecorder is the mock recorder for MockExecutor
type MockExecutorMockRecorder struct {
	mock *MockExecutor
}

// NewMockExecutor creates a new mock instance
func NewMockExecutor(ctrl *gomock.Controller) *MockExecutor {
	mock := &MockExecutor{ctrl: ctrl}
	mock.recorder = &MockExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExecutor) EXPECT() *MockExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockExecutor) Execute(sig StopSignal) model.StageStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", sig)
	ret0, _ := ret[0].(model.StageStatus)
	return ret0
}

// Execute indicates an expected call of Execute
func (mr *MockExecutorMockRecorder) Execute(sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockExecutor)(nil).Execute), sig)
}

// MockLogPersister is a mock of LogPersister interface
type MockLogPersister struct {
	ctrl     *gomock.Controller
	recorder *MockLogPersisterMockRecorder
}

// MockLogPersisterMockRecorder is the mock recorder for MockLogPersister
type MockLogPersisterMockRecorder struct {
	mock *MockLogPersister
}

// NewMockLogPersister creates a new mock instance
func NewMockLogPersister(ctrl *gomock.Controller) *MockLogPersister {
	mock := &MockLogPersister{ctrl: ctrl}
	mock.recorder = &MockLogPersisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogPersister) EXPECT() *MockLogPersisterMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockLogPersister) Write(log []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", log)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockLogPersisterMockRecorder) Write(log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockLogPersister)(nil).Write), log)
}

// Info mocks base method
func (m *MockLogPersister) Info(log string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", log)
}

// Info indicates an expected call of Info
func (mr *MockLogPersisterMockRecorder) Info(log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogPersister)(nil).Info), log)
}

// Infof mocks base method
func (m *MockLogPersister) Infof(format string, a ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof
func (mr *MockLogPersisterMockRecorder) Infof(format interface{}, a ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockLogPersister)(nil).Infof), varargs...)
}

// Success mocks base method
func (m *MockLogPersister) Success(log string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Success", log)
}

// Success indicates an expected call of Success
func (mr *MockLogPersisterMockRecorder) Success(log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Success", reflect.TypeOf((*MockLogPersister)(nil).Success), log)
}

// Successf mocks base method
func (m *MockLogPersister) Successf(format string, a ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	m.ctrl.Call(m, "Successf", varargs...)
}

// Successf indicates an expected call of Successf
func (mr *MockLogPersisterMockRecorder) Successf(format interface{}, a ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Successf", reflect.TypeOf((*MockLogPersister)(nil).Successf), varargs...)
}

// Error mocks base method
func (m *MockLogPersister) Error(log string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", log)
}

// Error indicates an expected call of Error
func (mr *MockLogPersisterMockRecorder) Error(log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogPersister)(nil).Error), log)
}

// Errorf mocks base method
func (m *MockLogPersister) Errorf(format string, a ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf
func (mr *MockLogPersisterMockRecorder) Errorf(format interface{}, a ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockLogPersister)(nil).Errorf), varargs...)
}

// MockMetadataStore is a mock of MetadataStore interface
type MockMetadataStore struct {
	ctrl     *gomock.Controller
	recorder *MockMetadataStoreMockRecorder
}

// MockMetadataStoreMockRecorder is the mock recorder for MockMetadataStore
type MockMetadataStoreMockRecorder struct {
	mock *MockMetadataStore
}

// NewMockMetadataStore creates a new mock instance
func NewMockMetadataStore(ctrl *gomock.Controller) *MockMetadataStore {
	mock := &MockMetadataStore{ctrl: ctrl}
	mock.recorder = &MockMetadataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetadataStore) EXPECT() *MockMetadataStoreMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockMetadataStore) Get(key string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockMetadataStoreMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMetadataStore)(nil).Get), key)
}

// Set mocks base method
func (m *MockMetadataStore) Set(ctx context.Context, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockMetadataStoreMockRecorder) Set(ctx, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockMetadataStore)(nil).Set), ctx, key, value)
}

// GetStageMetadata mocks base method
func (m *MockMetadataStore) GetStageMetadata(stageID string) (map[string]string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStageMetadata", stageID)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetStageMetadata indicates an expected call of GetStageMetadata
func (mr *MockMetadataStoreMockRecorder) GetStageMetadata(stageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStageMetadata", reflect.TypeOf((*MockMetadataStore)(nil).GetStageMetadata), stageID)
}

// SetStageMetadata mocks base method
func (m *MockMetadataStore) SetStageMetadata(ctx context.Context, stageID string, metadata map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStageMetadata", ctx, stageID, metadata)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStageMetadata indicates an expected call of SetStageMetadata
func (mr *MockMetadataStoreMockRecorder) SetStageMetadata(ctx, stageID, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStageMetadata", reflect.TypeOf((*MockMetadataStore)(nil).SetStageMetadata), ctx, stageID, metadata)
}

// MockCommandLister is a mock of CommandLister interface
type MockCommandLister struct {
	ctrl     *gomock.Controller
	recorder *MockCommandListerMockRecorder
}

// MockCommandListerMockRecorder is the mock recorder for MockCommandLister
type MockCommandListerMockRecorder struct {
	mock *MockCommandLister
}

// NewMockCommandLister creates a new mock instance
func NewMockCommandLister(ctrl *gomock.Controller) *MockCommandLister {
	mock := &MockCommandLister{ctrl: ctrl}
	mock.recorder = &MockCommandListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandLister) EXPECT() *MockCommandListerMockRecorder {
	return m.recorder
}

// ListCommands mocks base method
func (m *MockCommandLister) ListCommands() []model.ReportableCommand {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCommands")
	ret0, _ := ret[0].([]model.ReportableCommand)
	return ret0
}

// ListCommands indicates an expected call of ListCommands
func (mr *MockCommandListerMockRecorder) ListCommands() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCommands", reflect.TypeOf((*MockCommandLister)(nil).ListCommands))
}

// MockAppLiveResourceLister is a mock of AppLiveResourceLister interface
type MockAppLiveResourceLister struct {
	ctrl     *gomock.Controller
	recorder *MockAppLiveResourceListerMockRecorder
}

// MockAppLiveResourceListerMockRecorder is the mock recorder for MockAppLiveResourceLister
type MockAppLiveResourceListerMockRecorder struct {
	mock *MockAppLiveResourceLister
}

// NewMockAppLiveResourceLister creates a new mock instance
func NewMockAppLiveResourceLister(ctrl *gomock.Controller) *MockAppLiveResourceLister {
	mock := &MockAppLiveResourceLister{ctrl: ctrl}
	mock.recorder = &MockAppLiveResourceListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppLiveResourceLister) EXPECT() *MockAppLiveResourceListerMockRecorder {
	return m.recorder
}

// ListKubernetesResources mocks base method
func (m *MockAppLiveResourceLister) ListKubernetesResources() ([]kubernetes.Manifest, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListKubernetesResources")
	ret0, _ := ret[0].([]kubernetes.Manifest)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ListKubernetesResources indicates an expected call of ListKubernetesResources
func (mr *MockAppLiveResourceListerMockRecorder) ListKubernetesResources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListKubernetesResources", reflect.TypeOf((*MockAppLiveResourceLister)(nil).ListKubernetesResources))
}
