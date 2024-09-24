// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/interactor/sending_sale.go
//
// Generated by this command:
//
//	mockgen -source ./usecase/interactor/sending_sale.go -destination ./mock/mock_interactor/mock_sending_sale.go
//

// Package mock_interactor is a generated GoMock package.
package mock_interactor

import (
	reflect "reflect"

	interactor "github.com/spaceaiinc/autoscout-server/usecase/interactor"
	gomock "go.uber.org/mock/gomock"
)

// MockSendingSaleInteractor is a mock of SendingSaleInteractor interface.
type MockSendingSaleInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockSendingSaleInteractorMockRecorder
}

// MockSendingSaleInteractorMockRecorder is the mock recorder for MockSendingSaleInteractor.
type MockSendingSaleInteractorMockRecorder struct {
	mock *MockSendingSaleInteractor
}

// NewMockSendingSaleInteractor creates a new mock instance.
func NewMockSendingSaleInteractor(ctrl *gomock.Controller) *MockSendingSaleInteractor {
	mock := &MockSendingSaleInteractor{ctrl: ctrl}
	mock.recorder = &MockSendingSaleInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendingSaleInteractor) EXPECT() *MockSendingSaleInteractorMockRecorder {
	return m.recorder
}

// CreateSendingSale mocks base method.
func (m *MockSendingSaleInteractor) CreateSendingSale(input interactor.CreateSendingSaleInput) (interactor.CreateSendingSaleOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSendingSale", input)
	ret0, _ := ret[0].(interactor.CreateSendingSaleOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSendingSale indicates an expected call of CreateSendingSale.
func (mr *MockSendingSaleInteractorMockRecorder) CreateSendingSale(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSendingSale", reflect.TypeOf((*MockSendingSaleInteractor)(nil).CreateSendingSale), input)
}

// GetAllSendingSaleForMonthly mocks base method.
func (m *MockSendingSaleInteractor) GetAllSendingSaleForMonthly(input interactor.GetAllSendingSaleForMonthlyInput) (interactor.GetAllSendingSaleForMonthlyOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSendingSaleForMonthly", input)
	ret0, _ := ret[0].(interactor.GetAllSendingSaleForMonthlyOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSendingSaleForMonthly indicates an expected call of GetAllSendingSaleForMonthly.
func (mr *MockSendingSaleInteractorMockRecorder) GetAllSendingSaleForMonthly(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSendingSaleForMonthly", reflect.TypeOf((*MockSendingSaleInteractor)(nil).GetAllSendingSaleForMonthly), input)
}

// GetSendingSaleByID mocks base method.
func (m *MockSendingSaleInteractor) GetSendingSaleByID(input interactor.GetSendingSaleByIDInput) (interactor.GetSendingSaleByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingSaleByID", input)
	ret0, _ := ret[0].(interactor.GetSendingSaleByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingSaleByID indicates an expected call of GetSendingSaleByID.
func (mr *MockSendingSaleInteractorMockRecorder) GetSendingSaleByID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingSaleByID", reflect.TypeOf((*MockSendingSaleInteractor)(nil).GetSendingSaleByID), input)
}

// GetSendingSaleByJobSeekerIDAndEnterpriseID mocks base method.
func (m *MockSendingSaleInteractor) GetSendingSaleByJobSeekerIDAndEnterpriseID(input interactor.GetSendingSaleByJobSeekerIDAndEnterpriseIDInput) (interactor.GetSendingSaleByJobSeekerIDAndEnterpriseIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingSaleByJobSeekerIDAndEnterpriseID", input)
	ret0, _ := ret[0].(interactor.GetSendingSaleByJobSeekerIDAndEnterpriseIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingSaleByJobSeekerIDAndEnterpriseID indicates an expected call of GetSendingSaleByJobSeekerIDAndEnterpriseID.
func (mr *MockSendingSaleInteractorMockRecorder) GetSendingSaleByJobSeekerIDAndEnterpriseID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingSaleByJobSeekerIDAndEnterpriseID", reflect.TypeOf((*MockSendingSaleInteractor)(nil).GetSendingSaleByJobSeekerIDAndEnterpriseID), input)
}

// GetSendingSaleListByAgentIDForMonthly mocks base method.
func (m *MockSendingSaleInteractor) GetSendingSaleListByAgentIDForMonthly(input interactor.GetSendingSaleListByAgentIDForMonthlyInput) (interactor.GetSendingSaleListByAgentIDForMonthlyOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingSaleListByAgentIDForMonthly", input)
	ret0, _ := ret[0].(interactor.GetSendingSaleListByAgentIDForMonthlyOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingSaleListByAgentIDForMonthly indicates an expected call of GetSendingSaleListByAgentIDForMonthly.
func (mr *MockSendingSaleInteractorMockRecorder) GetSendingSaleListByAgentIDForMonthly(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingSaleListByAgentIDForMonthly", reflect.TypeOf((*MockSendingSaleInteractor)(nil).GetSendingSaleListByAgentIDForMonthly), input)
}

// GetSendingSaleListBySenderAgentIDForMonthly mocks base method.
func (m *MockSendingSaleInteractor) GetSendingSaleListBySenderAgentIDForMonthly(input interactor.GetSendingSaleListBySenderAgentIDForMonthlyInput) (interactor.GetSendingSaleListBySenderAgentIDForMonthlyOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingSaleListBySenderAgentIDForMonthly", input)
	ret0, _ := ret[0].(interactor.GetSendingSaleListBySenderAgentIDForMonthlyOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingSaleListBySenderAgentIDForMonthly indicates an expected call of GetSendingSaleListBySenderAgentIDForMonthly.
func (mr *MockSendingSaleInteractorMockRecorder) GetSendingSaleListBySenderAgentIDForMonthly(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingSaleListBySenderAgentIDForMonthly", reflect.TypeOf((*MockSendingSaleInteractor)(nil).GetSendingSaleListBySenderAgentIDForMonthly), input)
}

// UpdateSendingSale mocks base method.
func (m *MockSendingSaleInteractor) UpdateSendingSale(input interactor.UpdateSendingSaleInput) (interactor.UpdateSendingSaleOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSendingSale", input)
	ret0, _ := ret[0].(interactor.UpdateSendingSaleOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSendingSale indicates an expected call of UpdateSendingSale.
func (mr *MockSendingSaleInteractorMockRecorder) UpdateSendingSale(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSendingSale", reflect.TypeOf((*MockSendingSaleInteractor)(nil).UpdateSendingSale), input)
}
