// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/interactor/message_template.go
//
// Generated by this command:
//
//	mockgen -source ./usecase/interactor/message_template.go -destination ./mock/mock_interactor/mock_message_template.go
//

// Package mock_interactor is a generated GoMock package.
package mock_interactor

import (
	reflect "reflect"

	interactor "github.com/spaceaiinc/autoscout-server/usecase/interactor"
	gomock "go.uber.org/mock/gomock"
)

// MockMessageTemplateInteractor is a mock of MessageTemplateInteractor interface.
type MockMessageTemplateInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockMessageTemplateInteractorMockRecorder
}

// MockMessageTemplateInteractorMockRecorder is the mock recorder for MockMessageTemplateInteractor.
type MockMessageTemplateInteractorMockRecorder struct {
	mock *MockMessageTemplateInteractor
}

// NewMockMessageTemplateInteractor creates a new mock instance.
func NewMockMessageTemplateInteractor(ctrl *gomock.Controller) *MockMessageTemplateInteractor {
	mock := &MockMessageTemplateInteractor{ctrl: ctrl}
	mock.recorder = &MockMessageTemplateInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageTemplateInteractor) EXPECT() *MockMessageTemplateInteractorMockRecorder {
	return m.recorder
}

// CreateMessageTemplate mocks base method.
func (m *MockMessageTemplateInteractor) CreateMessageTemplate(input interactor.CreateMessageTemplateInput) (interactor.CreateMessageTemplateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessageTemplate", input)
	ret0, _ := ret[0].(interactor.CreateMessageTemplateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMessageTemplate indicates an expected call of CreateMessageTemplate.
func (mr *MockMessageTemplateInteractorMockRecorder) CreateMessageTemplate(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessageTemplate", reflect.TypeOf((*MockMessageTemplateInteractor)(nil).CreateMessageTemplate), input)
}

// DeleteMessageTemplate mocks base method.
func (m *MockMessageTemplateInteractor) DeleteMessageTemplate(input interactor.DeleteMessageTemplateInput) (interactor.DeleteMessageTemplateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessageTemplate", input)
	ret0, _ := ret[0].(interactor.DeleteMessageTemplateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMessageTemplate indicates an expected call of DeleteMessageTemplate.
func (mr *MockMessageTemplateInteractorMockRecorder) DeleteMessageTemplate(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessageTemplate", reflect.TypeOf((*MockMessageTemplateInteractor)(nil).DeleteMessageTemplate), input)
}

// GetMessageTemplateByID mocks base method.
func (m *MockMessageTemplateInteractor) GetMessageTemplateByID(input interactor.GetMessageTemplateByIDInput) (interactor.GetMessageTemplateByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageTemplateByID", input)
	ret0, _ := ret[0].(interactor.GetMessageTemplateByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageTemplateByID indicates an expected call of GetMessageTemplateByID.
func (mr *MockMessageTemplateInteractorMockRecorder) GetMessageTemplateByID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageTemplateByID", reflect.TypeOf((*MockMessageTemplateInteractor)(nil).GetMessageTemplateByID), input)
}

// GetMessageTemplateListByAgentID mocks base method.
func (m *MockMessageTemplateInteractor) GetMessageTemplateListByAgentID(input interactor.GetMessageTemplateListByAgentIDInput) (interactor.GetMessageTemplateListByAgentIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageTemplateListByAgentID", input)
	ret0, _ := ret[0].(interactor.GetMessageTemplateListByAgentIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageTemplateListByAgentID indicates an expected call of GetMessageTemplateListByAgentID.
func (mr *MockMessageTemplateInteractorMockRecorder) GetMessageTemplateListByAgentID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageTemplateListByAgentID", reflect.TypeOf((*MockMessageTemplateInteractor)(nil).GetMessageTemplateListByAgentID), input)
}

// GetMessageTemplateListByAgentStaffID mocks base method.
func (m *MockMessageTemplateInteractor) GetMessageTemplateListByAgentStaffID(input interactor.GetMessageTemplateListByAgentStaffIDInput) (interactor.GetMessageTemplateListByAgentStaffIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageTemplateListByAgentStaffID", input)
	ret0, _ := ret[0].(interactor.GetMessageTemplateListByAgentStaffIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageTemplateListByAgentStaffID indicates an expected call of GetMessageTemplateListByAgentStaffID.
func (mr *MockMessageTemplateInteractorMockRecorder) GetMessageTemplateListByAgentStaffID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageTemplateListByAgentStaffID", reflect.TypeOf((*MockMessageTemplateInteractor)(nil).GetMessageTemplateListByAgentStaffID), input)
}

// GetMessageTemplateListByAgentStaffIDAndSendScene mocks base method.
func (m *MockMessageTemplateInteractor) GetMessageTemplateListByAgentStaffIDAndSendScene(input interactor.GetMessageTemplateListByAgentStaffIDAndSendSceneInput) (interactor.GetMessageTemplateListByAgentStaffIDAndSendSceneOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageTemplateListByAgentStaffIDAndSendScene", input)
	ret0, _ := ret[0].(interactor.GetMessageTemplateListByAgentStaffIDAndSendSceneOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageTemplateListByAgentStaffIDAndSendScene indicates an expected call of GetMessageTemplateListByAgentStaffIDAndSendScene.
func (mr *MockMessageTemplateInteractorMockRecorder) GetMessageTemplateListByAgentStaffIDAndSendScene(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageTemplateListByAgentStaffIDAndSendScene", reflect.TypeOf((*MockMessageTemplateInteractor)(nil).GetMessageTemplateListByAgentStaffIDAndSendScene), input)
}

// UpdateMessageTemplate mocks base method.
func (m *MockMessageTemplateInteractor) UpdateMessageTemplate(input interactor.UpdateMessageTemplateInput) (interactor.UpdateMessageTemplateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessageTemplate", input)
	ret0, _ := ret[0].(interactor.UpdateMessageTemplateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMessageTemplate indicates an expected call of UpdateMessageTemplate.
func (mr *MockMessageTemplateInteractorMockRecorder) UpdateMessageTemplate(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessageTemplate", reflect.TypeOf((*MockMessageTemplateInteractor)(nil).UpdateMessageTemplate), input)
}
