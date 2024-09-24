// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/interactor/sending_customer.go
//
// Generated by this command:
//
//	mockgen -source ./usecase/interactor/sending_customer.go -destination ./mock/mock_interactor/mock_sending_customer.go
//

// Package mock_interactor is a generated GoMock package.
package mock_interactor

import (
	reflect "reflect"

	interactor "github.com/spaceaiinc/autoscout-server/usecase/interactor"
	gomock "go.uber.org/mock/gomock"
)

// MockSendingCustomerInteractor is a mock of SendingCustomerInteractor interface.
type MockSendingCustomerInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockSendingCustomerInteractorMockRecorder
}

// MockSendingCustomerInteractorMockRecorder is the mock recorder for MockSendingCustomerInteractor.
type MockSendingCustomerInteractorMockRecorder struct {
	mock *MockSendingCustomerInteractor
}

// NewMockSendingCustomerInteractor creates a new mock instance.
func NewMockSendingCustomerInteractor(ctrl *gomock.Controller) *MockSendingCustomerInteractor {
	mock := &MockSendingCustomerInteractor{ctrl: ctrl}
	mock.recorder = &MockSendingCustomerInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendingCustomerInteractor) EXPECT() *MockSendingCustomerInteractorMockRecorder {
	return m.recorder
}

// GetSearchSendingCustomerListByPageAndTabAndAgentID mocks base method.
func (m *MockSendingCustomerInteractor) GetSearchSendingCustomerListByPageAndTabAndAgentID(input interactor.GetSearchSendingCustomerListByPageAndTabAndAgentIDInput) (interactor.GetSearchSendingCustomerListByPageAndTabAndAgentIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSearchSendingCustomerListByPageAndTabAndAgentID", input)
	ret0, _ := ret[0].(interactor.GetSearchSendingCustomerListByPageAndTabAndAgentIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSearchSendingCustomerListByPageAndTabAndAgentID indicates an expected call of GetSearchSendingCustomerListByPageAndTabAndAgentID.
func (mr *MockSendingCustomerInteractorMockRecorder) GetSearchSendingCustomerListByPageAndTabAndAgentID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSearchSendingCustomerListByPageAndTabAndAgentID", reflect.TypeOf((*MockSendingCustomerInteractor)(nil).GetSearchSendingCustomerListByPageAndTabAndAgentID), input)
}

// GetSendingCustomerByID mocks base method.
func (m *MockSendingCustomerInteractor) GetSendingCustomerByID(input interactor.GetSendingCustomerByIDInput) (interactor.GetSendingCustomerByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingCustomerByID", input)
	ret0, _ := ret[0].(interactor.GetSendingCustomerByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingCustomerByID indicates an expected call of GetSendingCustomerByID.
func (mr *MockSendingCustomerInteractorMockRecorder) GetSendingCustomerByID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingCustomerByID", reflect.TypeOf((*MockSendingCustomerInteractor)(nil).GetSendingCustomerByID), input)
}

// GetSendingCustomerListByAgentID mocks base method.
func (m *MockSendingCustomerInteractor) GetSendingCustomerListByAgentID(input interactor.GetSendingCustomerListByAgentIDInput) (interactor.GetSendingCustomerListByAgentIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingCustomerListByAgentID", input)
	ret0, _ := ret[0].(interactor.GetSendingCustomerListByAgentIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingCustomerListByAgentID indicates an expected call of GetSendingCustomerListByAgentID.
func (mr *MockSendingCustomerInteractorMockRecorder) GetSendingCustomerListByAgentID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingCustomerListByAgentID", reflect.TypeOf((*MockSendingCustomerInteractor)(nil).GetSendingCustomerListByAgentID), input)
}
