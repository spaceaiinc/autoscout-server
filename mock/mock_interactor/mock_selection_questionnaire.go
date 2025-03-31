// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/interactor/selection_questionnaire.go
//
// Generated by this command:
//
//	mockgen -source ./usecase/interactor/selection_questionnaire.go -destination ./mock/mock_interactor/mock_selection_questionnaire.go
//

// Package mock_interactor is a generated GoMock package.
package mock_interactor

import (
	reflect "reflect"

	interactor "github.com/spaceaiinc/autoscout-server/usecase/interactor"
	gomock "go.uber.org/mock/gomock"
)

// MockSelectionQuestionnaireInteractor is a mock of SelectionQuestionnaireInteractor interface.
type MockSelectionQuestionnaireInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockSelectionQuestionnaireInteractorMockRecorder
}

// MockSelectionQuestionnaireInteractorMockRecorder is the mock recorder for MockSelectionQuestionnaireInteractor.
type MockSelectionQuestionnaireInteractorMockRecorder struct {
	mock *MockSelectionQuestionnaireInteractor
}

// NewMockSelectionQuestionnaireInteractor creates a new mock instance.
func NewMockSelectionQuestionnaireInteractor(ctrl *gomock.Controller) *MockSelectionQuestionnaireInteractor {
	mock := &MockSelectionQuestionnaireInteractor{ctrl: ctrl}
	mock.recorder = &MockSelectionQuestionnaireInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSelectionQuestionnaireInteractor) EXPECT() *MockSelectionQuestionnaireInteractorMockRecorder {
	return m.recorder
}

// CreateSelectionQuestionnaire mocks base method.
func (m *MockSelectionQuestionnaireInteractor) CreateSelectionQuestionnaire(input interactor.CreateSelectionQuestionnaireInput) (interactor.CreateSelectionQuestionnaireOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSelectionQuestionnaire", input)
	ret0, _ := ret[0].(interactor.CreateSelectionQuestionnaireOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSelectionQuestionnaire indicates an expected call of CreateSelectionQuestionnaire.
func (mr *MockSelectionQuestionnaireInteractorMockRecorder) CreateSelectionQuestionnaire(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSelectionQuestionnaire", reflect.TypeOf((*MockSelectionQuestionnaireInteractor)(nil).CreateSelectionQuestionnaire), input)
}

// GenerateSelectionQuestionnaireByUUID mocks base method.
func (m *MockSelectionQuestionnaireInteractor) GenerateSelectionQuestionnaireByUUID() (interactor.GenerateSelectionQuestionnaireByUUIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSelectionQuestionnaireByUUID")
	ret0, _ := ret[0].(interactor.GenerateSelectionQuestionnaireByUUIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateSelectionQuestionnaireByUUID indicates an expected call of GenerateSelectionQuestionnaireByUUID.
func (mr *MockSelectionQuestionnaireInteractorMockRecorder) GenerateSelectionQuestionnaireByUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSelectionQuestionnaireByUUID", reflect.TypeOf((*MockSelectionQuestionnaireInteractor)(nil).GenerateSelectionQuestionnaireByUUID))
}

// GetQuestionnaireForJobSeeker mocks base method.
func (m *MockSelectionQuestionnaireInteractor) GetQuestionnaireForJobSeeker(input interactor.GetQuestionnaireForJobSeekerInput) (interactor.GetQuestionnaireForJobSeekerOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuestionnaireForJobSeeker", input)
	ret0, _ := ret[0].(interactor.GetQuestionnaireForJobSeekerOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuestionnaireForJobSeeker indicates an expected call of GetQuestionnaireForJobSeeker.
func (mr *MockSelectionQuestionnaireInteractorMockRecorder) GetQuestionnaireForJobSeeker(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuestionnaireForJobSeeker", reflect.TypeOf((*MockSelectionQuestionnaireInteractor)(nil).GetQuestionnaireForJobSeeker), input)
}

// GetSelectionQuestionnaireOrNullByUUID mocks base method.
func (m *MockSelectionQuestionnaireInteractor) GetSelectionQuestionnaireOrNullByUUID(input interactor.GetSelectionQuestionnaireOrNullByUUIDInput) (interactor.GetSelectionQuestionnaireOrNullByUUIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelectionQuestionnaireOrNullByUUID", input)
	ret0, _ := ret[0].(interactor.GetSelectionQuestionnaireOrNullByUUIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSelectionQuestionnaireOrNullByUUID indicates an expected call of GetSelectionQuestionnaireOrNullByUUID.
func (mr *MockSelectionQuestionnaireInteractorMockRecorder) GetSelectionQuestionnaireOrNullByUUID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelectionQuestionnaireOrNullByUUID", reflect.TypeOf((*MockSelectionQuestionnaireInteractor)(nil).GetSelectionQuestionnaireOrNullByUUID), input)
}

// GetUnansweredQuestionnaireListByJobSeekerUUID mocks base method.
func (m *MockSelectionQuestionnaireInteractor) GetUnansweredQuestionnaireListByJobSeekerUUID(input interactor.GetUnansweredQuestionnaireListByJobSeekerUUIDInput) (interactor.GetUnansweredQuestionnaireListByJobSeekerUUIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnansweredQuestionnaireListByJobSeekerUUID", input)
	ret0, _ := ret[0].(interactor.GetUnansweredQuestionnaireListByJobSeekerUUIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnansweredQuestionnaireListByJobSeekerUUID indicates an expected call of GetUnansweredQuestionnaireListByJobSeekerUUID.
func (mr *MockSelectionQuestionnaireInteractorMockRecorder) GetUnansweredQuestionnaireListByJobSeekerUUID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnansweredQuestionnaireListByJobSeekerUUID", reflect.TypeOf((*MockSelectionQuestionnaireInteractor)(nil).GetUnansweredQuestionnaireListByJobSeekerUUID), input)
}

// UpdateSelectionQuestionnaire mocks base method.
func (m *MockSelectionQuestionnaireInteractor) UpdateSelectionQuestionnaire(input interactor.UpdateSelectionQuestionnaireInput) (interactor.UpdateSelectionQuestionnaireOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSelectionQuestionnaire", input)
	ret0, _ := ret[0].(interactor.UpdateSelectionQuestionnaireOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSelectionQuestionnaire indicates an expected call of UpdateSelectionQuestionnaire.
func (mr *MockSelectionQuestionnaireInteractorMockRecorder) UpdateSelectionQuestionnaire(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSelectionQuestionnaire", reflect.TypeOf((*MockSelectionQuestionnaireInteractor)(nil).UpdateSelectionQuestionnaire), input)
}
