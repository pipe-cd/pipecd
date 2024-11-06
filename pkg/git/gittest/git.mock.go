// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pipe-cd/pipecd/pkg/git (interfaces: Repo)

// Package gittest is a generated GoMock package.
package gittest

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	git "github.com/pipe-cd/pipecd/pkg/git"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// ChangedFiles mocks base method.
func (m *MockRepo) ChangedFiles(arg0 context.Context, arg1, arg2 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangedFiles", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangedFiles indicates an expected call of ChangedFiles.
func (mr *MockRepoMockRecorder) ChangedFiles(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangedFiles", reflect.TypeOf((*MockRepo)(nil).ChangedFiles), arg0, arg1, arg2)
}

// Checkout mocks base method.
func (m *MockRepo) Checkout(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkout", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Checkout indicates an expected call of Checkout.
func (mr *MockRepoMockRecorder) Checkout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkout", reflect.TypeOf((*MockRepo)(nil).Checkout), arg0, arg1)
}

// CheckoutPullRequest mocks base method.
func (m *MockRepo) CheckoutPullRequest(arg0 context.Context, arg1 int, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckoutPullRequest", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckoutPullRequest indicates an expected call of CheckoutPullRequest.
func (mr *MockRepoMockRecorder) CheckoutPullRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckoutPullRequest", reflect.TypeOf((*MockRepo)(nil).CheckoutPullRequest), arg0, arg1, arg2)
}

// Clean mocks base method.
func (m *MockRepo) Clean() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clean")
	ret0, _ := ret[0].(error)
	return ret0
}

// Clean indicates an expected call of Clean.
func (mr *MockRepoMockRecorder) Clean() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clean", reflect.TypeOf((*MockRepo)(nil).Clean))
}

// CleanPartially mocks base method.
func (m *MockRepo) CleanPartially(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanPartially", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanPartially indicates an expected call of CleanPartially.
func (mr *MockRepoMockRecorder) CleanPartially(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanPartially", reflect.TypeOf((*MockRepo)(nil).CleanPartially), arg0, arg1)
}

// CommitChanges mocks base method.
func (m *MockRepo) CommitChanges(arg0 context.Context, arg1, arg2 string, arg3 bool, arg4 map[string][]byte, arg5 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitChanges", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitChanges indicates an expected call of CommitChanges.
func (mr *MockRepoMockRecorder) CommitChanges(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitChanges", reflect.TypeOf((*MockRepo)(nil).CommitChanges), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Copy mocks base method.
func (m *MockRepo) Copy(arg0 string) (git.Repo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy", arg0)
	ret0, _ := ret[0].(git.Repo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Copy indicates an expected call of Copy.
func (mr *MockRepoMockRecorder) Copy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockRepo)(nil).Copy), arg0)
}

// GetClonedBranch mocks base method.
func (m *MockRepo) GetClonedBranch() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClonedBranch")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetClonedBranch indicates an expected call of GetClonedBranch.
func (mr *MockRepoMockRecorder) GetClonedBranch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClonedBranch", reflect.TypeOf((*MockRepo)(nil).GetClonedBranch))
}

// GetCommitForRev mocks base method.
func (m *MockRepo) GetCommitForRev(arg0 context.Context, arg1 string) (git.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitForRev", arg0, arg1)
	ret0, _ := ret[0].(git.Commit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitForRev indicates an expected call of GetCommitForRev.
func (mr *MockRepoMockRecorder) GetCommitForRev(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitForRev", reflect.TypeOf((*MockRepo)(nil).GetCommitForRev), arg0, arg1)
}

// GetLatestCommit mocks base method.
func (m *MockRepo) GetLatestCommit(arg0 context.Context) (git.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestCommit", arg0)
	ret0, _ := ret[0].(git.Commit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestCommit indicates an expected call of GetLatestCommit.
func (mr *MockRepoMockRecorder) GetLatestCommit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestCommit", reflect.TypeOf((*MockRepo)(nil).GetLatestCommit), arg0)
}

// GetPath mocks base method.
func (m *MockRepo) GetPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPath indicates an expected call of GetPath.
func (mr *MockRepoMockRecorder) GetPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPath", reflect.TypeOf((*MockRepo)(nil).GetPath))
}

// ListCommits mocks base method.
func (m *MockRepo) ListCommits(arg0 context.Context, arg1 string) ([]git.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCommits", arg0, arg1)
	ret0, _ := ret[0].([]git.Commit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCommits indicates an expected call of ListCommits.
func (mr *MockRepoMockRecorder) ListCommits(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCommits", reflect.TypeOf((*MockRepo)(nil).ListCommits), arg0, arg1)
}

// MergeRemoteBranch mocks base method.
func (m *MockRepo) MergeRemoteBranch(arg0 context.Context, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MergeRemoteBranch", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// MergeRemoteBranch indicates an expected call of MergeRemoteBranch.
func (mr *MockRepoMockRecorder) MergeRemoteBranch(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MergeRemoteBranch", reflect.TypeOf((*MockRepo)(nil).MergeRemoteBranch), arg0, arg1, arg2, arg3)
}

// Pull mocks base method.
func (m *MockRepo) Pull(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pull", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Pull indicates an expected call of Pull.
func (mr *MockRepoMockRecorder) Pull(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pull", reflect.TypeOf((*MockRepo)(nil).Pull), arg0, arg1)
}

// Push mocks base method.
func (m *MockRepo) Push(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Push", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push.
func (mr *MockRepoMockRecorder) Push(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockRepo)(nil).Push), arg0, arg1)
}
