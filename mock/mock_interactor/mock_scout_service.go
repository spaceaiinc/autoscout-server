// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/interactor/scout_service.go
//
// Generated by this command:
//
//	mockgen -source ./usecase/interactor/scout_service.go -destination ./mock/mock_interactor/mock_scout_service.go
//

// Package mock_interactor is a generated GoMock package.
package mock_interactor

import (
	reflect "reflect"

	interactor "github.com/spaceaiinc/autoscout-server/usecase/interactor"
	gomock "go.uber.org/mock/gomock"
)

// MockScoutServiceInteractor is a mock of ScoutServiceInteractor interface.
type MockScoutServiceInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockScoutServiceInteractorMockRecorder
}

// MockScoutServiceInteractorMockRecorder is the mock recorder for MockScoutServiceInteractor.
type MockScoutServiceInteractorMockRecorder struct {
	mock *MockScoutServiceInteractor
}

// NewMockScoutServiceInteractor creates a new mock instance.
func NewMockScoutServiceInteractor(ctrl *gomock.Controller) *MockScoutServiceInteractor {
	mock := &MockScoutServiceInteractor{ctrl: ctrl}
	mock.recorder = &MockScoutServiceInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScoutServiceInteractor) EXPECT() *MockScoutServiceInteractorMockRecorder {
	return m.recorder
}

// BatchEntry mocks base method.
func (m *MockScoutServiceInteractor) BatchEntry(input interactor.BatchEntryInput) (interactor.BatchEntryOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchEntry", input)
	ret0, _ := ret[0].(interactor.BatchEntryOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchEntry indicates an expected call of BatchEntry.
func (mr *MockScoutServiceInteractorMockRecorder) BatchEntry(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchEntry", reflect.TypeOf((*MockScoutServiceInteractor)(nil).BatchEntry), input)
}

// BatchScout mocks base method.
func (m *MockScoutServiceInteractor) BatchScout(input interactor.BatchScoutInput) (interactor.BatchScoutOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchScout", input)
	ret0, _ := ret[0].(interactor.BatchScoutOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchScout indicates an expected call of BatchScout.
func (mr *MockScoutServiceInteractorMockRecorder) BatchScout(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchScout", reflect.TypeOf((*MockScoutServiceInteractor)(nil).BatchScout), input)
}

// CreateScoutService mocks base method.
func (m *MockScoutServiceInteractor) CreateScoutService(input interactor.CreateScoutServiceInput) (interactor.CreateScoutServiceByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateScoutService", input)
	ret0, _ := ret[0].(interactor.CreateScoutServiceByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateScoutService indicates an expected call of CreateScoutService.
func (mr *MockScoutServiceInteractorMockRecorder) CreateScoutService(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateScoutService", reflect.TypeOf((*MockScoutServiceInteractor)(nil).CreateScoutService), input)
}

// DeleteScoutService mocks base method.
func (m *MockScoutServiceInteractor) DeleteScoutService(input interactor.DeleteScoutServiceInput) (interactor.DeleteScoutServiceByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteScoutService", input)
	ret0, _ := ret[0].(interactor.DeleteScoutServiceByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteScoutService indicates an expected call of DeleteScoutService.
func (mr *MockScoutServiceInteractorMockRecorder) DeleteScoutService(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteScoutService", reflect.TypeOf((*MockScoutServiceInteractor)(nil).DeleteScoutService), input)
}

// EntryOnAmbi mocks base method.
func (m *MockScoutServiceInteractor) EntryOnAmbi(input interactor.EntryOnAmbiInput) (interactor.EntryOnAmbiOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EntryOnAmbi", input)
	ret0, _ := ret[0].(interactor.EntryOnAmbiOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EntryOnAmbi indicates an expected call of EntryOnAmbi.
func (mr *MockScoutServiceInteractorMockRecorder) EntryOnAmbi(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EntryOnAmbi", reflect.TypeOf((*MockScoutServiceInteractor)(nil).EntryOnAmbi), input)
}

// EntryOnMynaviAgentScout mocks base method.
func (m *MockScoutServiceInteractor) EntryOnMynaviAgentScout(input interactor.EntryOnMynaviAgentScoutInput) (interactor.EntryOnMynaviAgentScoutOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EntryOnMynaviAgentScout", input)
	ret0, _ := ret[0].(interactor.EntryOnMynaviAgentScoutOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EntryOnMynaviAgentScout indicates an expected call of EntryOnMynaviAgentScout.
func (mr *MockScoutServiceInteractorMockRecorder) EntryOnMynaviAgentScout(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EntryOnMynaviAgentScout", reflect.TypeOf((*MockScoutServiceInteractor)(nil).EntryOnMynaviAgentScout), input)
}

// EntryOnMynaviScouting mocks base method.
func (m *MockScoutServiceInteractor) EntryOnMynaviScouting(input interactor.EntryOnMynaviScoutingInput) (interactor.EntryOnMynaviScoutingOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EntryOnMynaviScouting", input)
	ret0, _ := ret[0].(interactor.EntryOnMynaviScoutingOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EntryOnMynaviScouting indicates an expected call of EntryOnMynaviScouting.
func (mr *MockScoutServiceInteractorMockRecorder) EntryOnMynaviScouting(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EntryOnMynaviScouting", reflect.TypeOf((*MockScoutServiceInteractor)(nil).EntryOnMynaviScouting), input)
}

// EntryOnRan mocks base method.
func (m *MockScoutServiceInteractor) EntryOnRan(input interactor.EntryOnRanInput) (interactor.EntryOnRanOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EntryOnRan", input)
	ret0, _ := ret[0].(interactor.EntryOnRanOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EntryOnRan indicates an expected call of EntryOnRan.
func (mr *MockScoutServiceInteractorMockRecorder) EntryOnRan(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EntryOnRan", reflect.TypeOf((*MockScoutServiceInteractor)(nil).EntryOnRan), input)
}

// GetByID mocks base method.
func (m *MockScoutServiceInteractor) GetByID(input interactor.ScoutServiceGetByIDInput) (interactor.ScoutServiceGetByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", input)
	ret0, _ := ret[0].(interactor.ScoutServiceGetByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockScoutServiceInteractorMockRecorder) GetByID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockScoutServiceInteractor)(nil).GetByID), input)
}

// GetListByAgentID mocks base method.
func (m *MockScoutServiceInteractor) GetListByAgentID(input interactor.GetListByAgentIDInput) (interactor.GetListByAgentIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListByAgentID", input)
	ret0, _ := ret[0].(interactor.GetListByAgentIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListByAgentID indicates an expected call of GetListByAgentID.
func (mr *MockScoutServiceInteractorMockRecorder) GetListByAgentID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListByAgentID", reflect.TypeOf((*MockScoutServiceInteractor)(nil).GetListByAgentID), input)
}

// GmailWebHook mocks base method.
func (m *MockScoutServiceInteractor) GmailWebHook(input interactor.GmailWebHookInput) (interactor.GmailWebHookOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GmailWebHook", input)
	ret0, _ := ret[0].(interactor.GmailWebHookOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GmailWebHook indicates an expected call of GmailWebHook.
func (mr *MockScoutServiceInteractorMockRecorder) GmailWebHook(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GmailWebHook", reflect.TypeOf((*MockScoutServiceInteractor)(nil).GmailWebHook), input)
}

// ScoutOnAmbi mocks base method.
func (m *MockScoutServiceInteractor) ScoutOnAmbi(input interactor.ScoutOnAmbiInput) (interactor.ScoutOnAmbiOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScoutOnAmbi", input)
	ret0, _ := ret[0].(interactor.ScoutOnAmbiOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScoutOnAmbi indicates an expected call of ScoutOnAmbi.
func (mr *MockScoutServiceInteractorMockRecorder) ScoutOnAmbi(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScoutOnAmbi", reflect.TypeOf((*MockScoutServiceInteractor)(nil).ScoutOnAmbi), input)
}

// ScoutOnMynaviAgentScout mocks base method.
func (m *MockScoutServiceInteractor) ScoutOnMynaviAgentScout(input interactor.ScoutOnMynaviAgentScoutInput) (interactor.ScoutOnMynaviAgentScoutOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScoutOnMynaviAgentScout", input)
	ret0, _ := ret[0].(interactor.ScoutOnMynaviAgentScoutOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScoutOnMynaviAgentScout indicates an expected call of ScoutOnMynaviAgentScout.
func (mr *MockScoutServiceInteractorMockRecorder) ScoutOnMynaviAgentScout(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScoutOnMynaviAgentScout", reflect.TypeOf((*MockScoutServiceInteractor)(nil).ScoutOnMynaviAgentScout), input)
}

// ScoutOnMynaviScouting mocks base method.
func (m *MockScoutServiceInteractor) ScoutOnMynaviScouting(input interactor.ScoutOnMynaviScoutingInput) (interactor.ScoutOnMynaviScoutingOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScoutOnMynaviScouting", input)
	ret0, _ := ret[0].(interactor.ScoutOnMynaviScoutingOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScoutOnMynaviScouting indicates an expected call of ScoutOnMynaviScouting.
func (mr *MockScoutServiceInteractorMockRecorder) ScoutOnMynaviScouting(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScoutOnMynaviScouting", reflect.TypeOf((*MockScoutServiceInteractor)(nil).ScoutOnMynaviScouting), input)
}

// ScoutOnRan mocks base method.
func (m *MockScoutServiceInteractor) ScoutOnRan(input interactor.ScoutOnRanInput) (interactor.ScoutOnRanOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScoutOnRan", input)
	ret0, _ := ret[0].(interactor.ScoutOnRanOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScoutOnRan indicates an expected call of ScoutOnRan.
func (mr *MockScoutServiceInteractorMockRecorder) ScoutOnRan(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScoutOnRan", reflect.TypeOf((*MockScoutServiceInteractor)(nil).ScoutOnRan), input)
}

// UpdateScoutService mocks base method.
func (m *MockScoutServiceInteractor) UpdateScoutService(input interactor.UpdateScoutServiceInput) (interactor.UpdateScoutServiceByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateScoutService", input)
	ret0, _ := ret[0].(interactor.UpdateScoutServiceByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateScoutService indicates an expected call of UpdateScoutService.
func (mr *MockScoutServiceInteractorMockRecorder) UpdateScoutService(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateScoutService", reflect.TypeOf((*MockScoutServiceInteractor)(nil).UpdateScoutService), input)
}

// UpdateScoutServicePassword mocks base method.
func (m *MockScoutServiceInteractor) UpdateScoutServicePassword(input interactor.UpdateScoutServicePasswordInput) (interactor.UpdateScoutServicePasswordByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateScoutServicePassword", input)
	ret0, _ := ret[0].(interactor.UpdateScoutServicePasswordByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateScoutServicePassword indicates an expected call of UpdateScoutServicePassword.
func (mr *MockScoutServiceInteractorMockRecorder) UpdateScoutServicePassword(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateScoutServicePassword", reflect.TypeOf((*MockScoutServiceInteractor)(nil).UpdateScoutServicePassword), input)
}
