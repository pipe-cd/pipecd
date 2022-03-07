// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/kubernetes (interfaces: Provider)

// Package providertest is a generated GoMock package.
package providertest

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	kubernetes "github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/kubernetes"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// Apply mocks base method.
func (m *MockProvider) Apply(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Apply", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Apply indicates an expected call of Apply.
func (mr *MockProviderMockRecorder) Apply(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockProvider)(nil).Apply), arg0)
}

// ApplyManifest mocks base method.
func (m *MockProvider) ApplyManifest(arg0 context.Context, arg1 kubernetes.Manifest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyManifest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyManifest indicates an expected call of ApplyManifest.
func (mr *MockProviderMockRecorder) ApplyManifest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyManifest", reflect.TypeOf((*MockProvider)(nil).ApplyManifest), arg0, arg1)
}

// CreateManifest mocks base method.
func (m *MockProvider) CreateManifest(arg0 context.Context, arg1 kubernetes.Manifest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateManifest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateManifest indicates an expected call of CreateManifest.
func (mr *MockProviderMockRecorder) CreateManifest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateManifest", reflect.TypeOf((*MockProvider)(nil).CreateManifest), arg0, arg1)
}

// Delete mocks base method.
func (m *MockProvider) Delete(arg0 context.Context, arg1 kubernetes.ResourceKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProviderMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProvider)(nil).Delete), arg0, arg1)
}

// LoadManifests mocks base method.
func (m *MockProvider) LoadManifests(arg0 context.Context) ([]kubernetes.Manifest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadManifests", arg0)
	ret0, _ := ret[0].([]kubernetes.Manifest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadManifests indicates an expected call of LoadManifests.
func (mr *MockProviderMockRecorder) LoadManifests(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadManifests", reflect.TypeOf((*MockProvider)(nil).LoadManifests), arg0)
}

// ReplaceManifest mocks base method.
func (m *MockProvider) ReplaceManifest(arg0 context.Context, arg1 kubernetes.Manifest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceManifest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceManifest indicates an expected call of ReplaceManifest.
func (mr *MockProviderMockRecorder) ReplaceManifest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceManifest", reflect.TypeOf((*MockProvider)(nil).ReplaceManifest), arg0, arg1)
}
