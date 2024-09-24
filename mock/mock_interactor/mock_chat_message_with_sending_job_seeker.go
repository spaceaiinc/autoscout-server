// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/interactor/chat_message_with_sending_job_seeker.go
//
// Generated by this command:
//
//	mockgen -source ./usecase/interactor/chat_message_with_sending_job_seeker.go -destination ./mock/mock_interactor/mock_chat_message_with_sending_job_seeker.go
//

// Package mock_interactor is a generated GoMock package.
package mock_interactor

import (
	reflect "reflect"

	interactor "github.com/spaceaiinc/autoscout-server/usecase/interactor"
	gomock "go.uber.org/mock/gomock"
)

// MockChatMessageWithSendingJobSeekerInteractor is a mock of ChatMessageWithSendingJobSeekerInteractor interface.
type MockChatMessageWithSendingJobSeekerInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockChatMessageWithSendingJobSeekerInteractorMockRecorder
}

// MockChatMessageWithSendingJobSeekerInteractorMockRecorder is the mock recorder for MockChatMessageWithSendingJobSeekerInteractor.
type MockChatMessageWithSendingJobSeekerInteractorMockRecorder struct {
	mock *MockChatMessageWithSendingJobSeekerInteractor
}

// NewMockChatMessageWithSendingJobSeekerInteractor creates a new mock instance.
func NewMockChatMessageWithSendingJobSeekerInteractor(ctrl *gomock.Controller) *MockChatMessageWithSendingJobSeekerInteractor {
	mock := &MockChatMessageWithSendingJobSeekerInteractor{ctrl: ctrl}
	mock.recorder = &MockChatMessageWithSendingJobSeekerInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatMessageWithSendingJobSeekerInteractor) EXPECT() *MockChatMessageWithSendingJobSeekerInteractorMockRecorder {
	return m.recorder
}

// GetChatMessageWithSendingJobSeekerListByGroupID mocks base method.
func (m *MockChatMessageWithSendingJobSeekerInteractor) GetChatMessageWithSendingJobSeekerListByGroupID(input interactor.GetChatMessageWithSendingJobSeekerListByGroupIDInput) (interactor.GetChatMessageWithSendingJobSeekerListByGroupIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChatMessageWithSendingJobSeekerListByGroupID", input)
	ret0, _ := ret[0].(interactor.GetChatMessageWithSendingJobSeekerListByGroupIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChatMessageWithSendingJobSeekerListByGroupID indicates an expected call of GetChatMessageWithSendingJobSeekerListByGroupID.
func (mr *MockChatMessageWithSendingJobSeekerInteractorMockRecorder) GetChatMessageWithSendingJobSeekerListByGroupID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatMessageWithSendingJobSeekerListByGroupID", reflect.TypeOf((*MockChatMessageWithSendingJobSeekerInteractor)(nil).GetChatMessageWithSendingJobSeekerListByGroupID), input)
}

// SendChatMessageWithSendingJobSeekerLineImage mocks base method.
func (m *MockChatMessageWithSendingJobSeekerInteractor) SendChatMessageWithSendingJobSeekerLineImage(input interactor.SendChatMessageWithSendingJobSeekerLineImageInput) (interactor.SendChatMessageWithSendingJobSeekerLineImageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendChatMessageWithSendingJobSeekerLineImage", input)
	ret0, _ := ret[0].(interactor.SendChatMessageWithSendingJobSeekerLineImageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendChatMessageWithSendingJobSeekerLineImage indicates an expected call of SendChatMessageWithSendingJobSeekerLineImage.
func (mr *MockChatMessageWithSendingJobSeekerInteractorMockRecorder) SendChatMessageWithSendingJobSeekerLineImage(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendChatMessageWithSendingJobSeekerLineImage", reflect.TypeOf((*MockChatMessageWithSendingJobSeekerInteractor)(nil).SendChatMessageWithSendingJobSeekerLineImage), input)
}

// SendChatMessageWithSendingJobSeekerLineMessage mocks base method.
func (m *MockChatMessageWithSendingJobSeekerInteractor) SendChatMessageWithSendingJobSeekerLineMessage(input interactor.SendChatMessageWithSendingJobSeekerLineMessageInput) (interactor.SendChatMessageWithSendingJobSeekerLineMessageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendChatMessageWithSendingJobSeekerLineMessage", input)
	ret0, _ := ret[0].(interactor.SendChatMessageWithSendingJobSeekerLineMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendChatMessageWithSendingJobSeekerLineMessage indicates an expected call of SendChatMessageWithSendingJobSeekerLineMessage.
func (mr *MockChatMessageWithSendingJobSeekerInteractorMockRecorder) SendChatMessageWithSendingJobSeekerLineMessage(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendChatMessageWithSendingJobSeekerLineMessage", reflect.TypeOf((*MockChatMessageWithSendingJobSeekerInteractor)(nil).SendChatMessageWithSendingJobSeekerLineMessage), input)
}
