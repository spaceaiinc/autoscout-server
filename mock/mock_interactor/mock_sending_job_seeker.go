// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/interactor/sending_job_seeker.go
//
// Generated by this command:
//
//	mockgen -source ./usecase/interactor/sending_job_seeker.go -destination ./mock/mock_interactor/mock_sending_job_seeker.go
//

// Package mock_interactor is a generated GoMock package.
package mock_interactor

import (
	reflect "reflect"

	entity "github.com/spaceaiinc/autoscout-server/domain/entity"
	interactor "github.com/spaceaiinc/autoscout-server/usecase/interactor"
	gomock "go.uber.org/mock/gomock"
)

// MockSendingJobSeekerInteractor is a mock of SendingJobSeekerInteractor interface.
type MockSendingJobSeekerInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockSendingJobSeekerInteractorMockRecorder
}

// MockSendingJobSeekerInteractorMockRecorder is the mock recorder for MockSendingJobSeekerInteractor.
type MockSendingJobSeekerInteractorMockRecorder struct {
	mock *MockSendingJobSeekerInteractor
}

// NewMockSendingJobSeekerInteractor creates a new mock instance.
func NewMockSendingJobSeekerInteractor(ctrl *gomock.Controller) *MockSendingJobSeekerInteractor {
	mock := &MockSendingJobSeekerInteractor{ctrl: ctrl}
	mock.recorder = &MockSendingJobSeekerInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendingJobSeekerInteractor) EXPECT() *MockSendingJobSeekerInteractorMockRecorder {
	return m.recorder
}

// CreateSendingInitialQuestionnaire mocks base method.
func (m *MockSendingJobSeekerInteractor) CreateSendingInitialQuestionnaire(input interactor.CreateSendingInitialQuestionnaireInput) (interactor.CreateSendingInitialQuestionnairdhutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSendingInitialQuestionnaire", input)
	ret0, _ := ret[0].(interactor.CreateSendingInitialQuestionnairdhutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSendingInitialQuestionnaire indicates an expected call of CreateSendingInitialQuestionnaire.
func (mr *MockSendingJobSeekerInteractorMockRecorder) CreateSendingInitialQuestionnaire(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSendingInitialQuestionnaire", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).CreateSendingInitialQuestionnaire), input)
}

// CreateSendingJobSeeker mocks base method.
func (m *MockSendingJobSeekerInteractor) CreateSendingJobSeeker(input interactor.CreateSendingJobSeekerInput) (interactor.CreateSendingJobSeekerOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSendingJobSeeker", input)
	ret0, _ := ret[0].(interactor.CreateSendingJobSeekerOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSendingJobSeeker indicates an expected call of CreateSendingJobSeeker.
func (mr *MockSendingJobSeekerInteractorMockRecorder) CreateSendingJobSeeker(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSendingJobSeeker", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).CreateSendingJobSeeker), input)
}

// CreateSendingJobSeekerEndStatus mocks base method.
func (m *MockSendingJobSeekerInteractor) CreateSendingJobSeekerEndStatus(input interactor.CreateSendingJobSeekerEndStatusInput) (interactor.CreateSendingJobSeekerEndStatusOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSendingJobSeekerEndStatus", input)
	ret0, _ := ret[0].(interactor.CreateSendingJobSeekerEndStatusOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSendingJobSeekerEndStatus indicates an expected call of CreateSendingJobSeekerEndStatus.
func (mr *MockSendingJobSeekerInteractorMockRecorder) CreateSendingJobSeekerEndStatus(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSendingJobSeekerEndStatus", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).CreateSendingJobSeekerEndStatus), input)
}

// CreateSendingJobSeekerFromJobSeeker mocks base method.
func (m *MockSendingJobSeekerInteractor) CreateSendingJobSeekerFromJobSeeker(input interactor.CreateSendingJobSeekerFromJobSeekerInput) (interactor.CreateSendingJobSeekerFromJobSeekerOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSendingJobSeekerFromJobSeeker", input)
	ret0, _ := ret[0].(interactor.CreateSendingJobSeekerFromJobSeekerOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSendingJobSeekerFromJobSeeker indicates an expected call of CreateSendingJobSeekerFromJobSeeker.
func (mr *MockSendingJobSeekerInteractorMockRecorder) CreateSendingJobSeekerFromJobSeeker(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSendingJobSeekerFromJobSeeker", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).CreateSendingJobSeekerFromJobSeeker), input)
}

// DeleteSendingJobSeeker mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeeker(input interactor.DeleteSendingJobSeekerInput) (interactor.DeleteSendingJobSeekerOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeeker", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeeker indicates an expected call of DeleteSendingJobSeeker.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeeker(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeeker", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeeker), input)
}

// DeleteSendingJobSeekerCVOriginURL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerCVOriginURL(input interactor.DeleteSendingJobSeekerCVOriginURLInput) (interactor.DeleteSendingJobSeekerCVOriginURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerCVOriginURL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerCVOriginURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerCVOriginURL indicates an expected call of DeleteSendingJobSeekerCVOriginURL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerCVOriginURL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerCVOriginURL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerCVOriginURL), input)
}

// DeleteSendingJobSeekerCVPDFURL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerCVPDFURL(input interactor.DeleteSendingJobSeekerCVPDFURLInput) (interactor.DeleteSendingJobSeekerCVPDFURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerCVPDFURL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerCVPDFURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerCVPDFURL indicates an expected call of DeleteSendingJobSeekerCVPDFURL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerCVPDFURL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerCVPDFURL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerCVPDFURL), input)
}

// DeleteSendingJobSeekerIDPhotoURL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerIDPhotoURL(input interactor.DeleteSendingJobSeekerIDPhotoURLInput) (interactor.DeleteSendingJobSeekerIDPhotoURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerIDPhotoURL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerIDPhotoURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerIDPhotoURL indicates an expected call of DeleteSendingJobSeekerIDPhotoURL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerIDPhotoURL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerIDPhotoURL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerIDPhotoURL), input)
}

// DeleteSendingJobSeekerOtherDocument1URL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerOtherDocument1URL(input interactor.DeleteSendingJobSeekerOtherDocument1URLInput) (interactor.DeleteSendingJobSeekerOtherDocument1URLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerOtherDocument1URL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerOtherDocument1URLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerOtherDocument1URL indicates an expected call of DeleteSendingJobSeekerOtherDocument1URL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerOtherDocument1URL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerOtherDocument1URL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerOtherDocument1URL), input)
}

// DeleteSendingJobSeekerOtherDocument2URL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerOtherDocument2URL(input interactor.DeleteSendingJobSeekerOtherDocument2URLInput) (interactor.DeleteSendingJobSeekerOtherDocument2URLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerOtherDocument2URL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerOtherDocument2URLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerOtherDocument2URL indicates an expected call of DeleteSendingJobSeekerOtherDocument2URL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerOtherDocument2URL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerOtherDocument2URL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerOtherDocument2URL), input)
}

// DeleteSendingJobSeekerOtherDocument3URL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerOtherDocument3URL(input interactor.DeleteSendingJobSeekerOtherDocument3URLInput) (interactor.DeleteSendingJobSeekerOtherDocument3URLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerOtherDocument3URL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerOtherDocument3URLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerOtherDocument3URL indicates an expected call of DeleteSendingJobSeekerOtherDocument3URL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerOtherDocument3URL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerOtherDocument3URL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerOtherDocument3URL), input)
}

// DeleteSendingJobSeekerRecommendationOriginURL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerRecommendationOriginURL(input interactor.DeleteSendingJobSeekerRecommendationOriginURLInput) (interactor.DeleteSendingJobSeekerRecommendationOriginURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerRecommendationOriginURL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerRecommendationOriginURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerRecommendationOriginURL indicates an expected call of DeleteSendingJobSeekerRecommendationOriginURL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerRecommendationOriginURL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerRecommendationOriginURL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerRecommendationOriginURL), input)
}

// DeleteSendingJobSeekerRecommendationPDFURL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerRecommendationPDFURL(input interactor.DeleteSendingJobSeekerRecommendationPDFURLInput) (interactor.DeleteSendingJobSeekerRecommendationPDFURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerRecommendationPDFURL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerRecommendationPDFURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerRecommendationPDFURL indicates an expected call of DeleteSendingJobSeekerRecommendationPDFURL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerRecommendationPDFURL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerRecommendationPDFURL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerRecommendationPDFURL), input)
}

// DeleteSendingJobSeekerResumeOriginURL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerResumeOriginURL(input interactor.DeleteSendingJobSeekerResumeOriginURLInput) (interactor.DeleteSendingJobSeekerResumeOriginURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerResumeOriginURL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerResumeOriginURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerResumeOriginURL indicates an expected call of DeleteSendingJobSeekerResumeOriginURL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerResumeOriginURL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerResumeOriginURL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerResumeOriginURL), input)
}

// DeleteSendingJobSeekerResumePDFURL mocks base method.
func (m *MockSendingJobSeekerInteractor) DeleteSendingJobSeekerResumePDFURL(input interactor.DeleteSendingJobSeekerResumePDFURLInput) (interactor.DeleteSendingJobSeekerResumePDFURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSendingJobSeekerResumePDFURL", input)
	ret0, _ := ret[0].(interactor.DeleteSendingJobSeekerResumePDFURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSendingJobSeekerResumePDFURL indicates an expected call of DeleteSendingJobSeekerResumePDFURL.
func (mr *MockSendingJobSeekerInteractorMockRecorder) DeleteSendingJobSeekerResumePDFURL(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSendingJobSeekerResumePDFURL", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).DeleteSendingJobSeekerResumePDFURL), input)
}

// FirstUpdateSendingJobSeeker mocks base method.
func (m *MockSendingJobSeekerInteractor) FirstUpdateSendingJobSeeker(input interactor.FirstUpdateSendingJobSeekerInput) (interactor.FirstUpdateSendingJobSeekerOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstUpdateSendingJobSeeker", input)
	ret0, _ := ret[0].(interactor.FirstUpdateSendingJobSeekerOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FirstUpdateSendingJobSeeker indicates an expected call of FirstUpdateSendingJobSeeker.
func (mr *MockSendingJobSeekerInteractorMockRecorder) FirstUpdateSendingJobSeeker(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstUpdateSendingJobSeeker", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).FirstUpdateSendingJobSeeker), input)
}

// GetIsNotViewSendingJobSeekerCountByAgentStaffID mocks base method.
func (m *MockSendingJobSeekerInteractor) GetIsNotViewSendingJobSeekerCountByAgentStaffID(input interactor.GetIsNotViewSendingJobSeekerCountByAgentStaffIDInput) (interactor.GetIsNotViewSendingJobSeekerCountByAgentStaffIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIsNotViewSendingJobSeekerCountByAgentStaffID", input)
	ret0, _ := ret[0].(interactor.GetIsNotViewSendingJobSeekerCountByAgentStaffIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIsNotViewSendingJobSeekerCountByAgentStaffID indicates an expected call of GetIsNotViewSendingJobSeekerCountByAgentStaffID.
func (mr *MockSendingJobSeekerInteractorMockRecorder) GetIsNotViewSendingJobSeekerCountByAgentStaffID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIsNotViewSendingJobSeekerCountByAgentStaffID", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).GetIsNotViewSendingJobSeekerCountByAgentStaffID), input)
}

// GetSearchListForSendingJobSeekerManagementByAgentID mocks base method.
func (m *MockSendingJobSeekerInteractor) GetSearchListForSendingJobSeekerManagementByAgentID(input interactor.GetSearchListForSendingJobSeekerManagementByAgentIDInput) (interactor.GetSearchListForSendingJobSeekerManagementByAgentIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSearchListForSendingJobSeekerManagementByAgentID", input)
	ret0, _ := ret[0].(interactor.GetSearchListForSendingJobSeekerManagementByAgentIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSearchListForSendingJobSeekerManagementByAgentID indicates an expected call of GetSearchListForSendingJobSeekerManagementByAgentID.
func (mr *MockSendingJobSeekerInteractorMockRecorder) GetSearchListForSendingJobSeekerManagementByAgentID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSearchListForSendingJobSeekerManagementByAgentID", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).GetSearchListForSendingJobSeekerManagementByAgentID), input)
}

// GetSendingJobSeekerByID mocks base method.
func (m *MockSendingJobSeekerInteractor) GetSendingJobSeekerByID(input interactor.GetSendingJobSeekerByIDInput) (interactor.GetSendingJobSeekerByIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingJobSeekerByID", input)
	ret0, _ := ret[0].(interactor.GetSendingJobSeekerByIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingJobSeekerByID indicates an expected call of GetSendingJobSeekerByID.
func (mr *MockSendingJobSeekerInteractorMockRecorder) GetSendingJobSeekerByID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingJobSeekerByID", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).GetSendingJobSeekerByID), input)
}

// GetSendingJobSeekerByUUID mocks base method.
func (m *MockSendingJobSeekerInteractor) GetSendingJobSeekerByUUID(input interactor.GetSendingJobSeekerByUUIDInput) (interactor.GetSendingJobSeekerByUUIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingJobSeekerByUUID", input)
	ret0, _ := ret[0].(interactor.GetSendingJobSeekerByUUIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingJobSeekerByUUID indicates an expected call of GetSendingJobSeekerByUUID.
func (mr *MockSendingJobSeekerInteractorMockRecorder) GetSendingJobSeekerByUUID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingJobSeekerByUUID", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).GetSendingJobSeekerByUUID), input)
}

// GetSendingJobSeekerDocumentBySendingJobSeekerID mocks base method.
func (m *MockSendingJobSeekerInteractor) GetSendingJobSeekerDocumentBySendingJobSeekerID(input interactor.GetSendingJobSeekerDocumentBySendingJobSeekerIDInput) (interactor.GetSendingJobSeekerDocumentBySendingJobSeekerIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingJobSeekerDocumentBySendingJobSeekerID", input)
	ret0, _ := ret[0].(interactor.GetSendingJobSeekerDocumentBySendingJobSeekerIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingJobSeekerDocumentBySendingJobSeekerID indicates an expected call of GetSendingJobSeekerDocumentBySendingJobSeekerID.
func (mr *MockSendingJobSeekerInteractorMockRecorder) GetSendingJobSeekerDocumentBySendingJobSeekerID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingJobSeekerDocumentBySendingJobSeekerID", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).GetSendingJobSeekerDocumentBySendingJobSeekerID), input)
}

// GetSendingJobSeekerDocumentByUUID mocks base method.
func (m *MockSendingJobSeekerInteractor) GetSendingJobSeekerDocumentByUUID(input interactor.GetSendingJobSeekerDocumentByUUIDInput) (interactor.GetSendingJobSeekerDocumentByUUIDOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendingJobSeekerDocumentByUUID", input)
	ret0, _ := ret[0].(interactor.GetSendingJobSeekerDocumentByUUIDOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendingJobSeekerDocumentByUUID indicates an expected call of GetSendingJobSeekerDocumentByUUID.
func (mr *MockSendingJobSeekerInteractorMockRecorder) GetSendingJobSeekerDocumentByUUID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendingJobSeekerDocumentByUUID", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).GetSendingJobSeekerDocumentByUUID), input)
}

// UpdateIsVewForUnregister mocks base method.
func (m *MockSendingJobSeekerInteractor) UpdateIsVewForUnregister(input interactor.UpdateIsVewForUnregisterInput) (interactor.UpdateIsVewForUnregisterOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIsVewForUnregister", input)
	ret0, _ := ret[0].(interactor.UpdateIsVewForUnregisterOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateIsVewForUnregister indicates an expected call of UpdateIsVewForUnregister.
func (mr *MockSendingJobSeekerInteractorMockRecorder) UpdateIsVewForUnregister(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIsVewForUnregister", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).UpdateIsVewForUnregister), input)
}

// UpdateIsVewForWating mocks base method.
func (m *MockSendingJobSeekerInteractor) UpdateIsVewForWating(input interactor.UpdateIsVewForWatingInput) (interactor.UpdateIsVewForWatingOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIsVewForWating", input)
	ret0, _ := ret[0].(interactor.UpdateIsVewForWatingOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateIsVewForWating indicates an expected call of UpdateIsVewForWating.
func (mr *MockSendingJobSeekerInteractorMockRecorder) UpdateIsVewForWating(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIsVewForWating", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).UpdateIsVewForWating), input)
}

// UpdateSendingInterviewDateBySendingJobSeekerID mocks base method.
func (m *MockSendingJobSeekerInteractor) UpdateSendingInterviewDateBySendingJobSeekerID(input interactor.UpdateSendingJobSeekerInterviewDateInput) (interactor.UpdateSendingJobSeekerInterviewDateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSendingInterviewDateBySendingJobSeekerID", input)
	ret0, _ := ret[0].(interactor.UpdateSendingJobSeekerInterviewDateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSendingInterviewDateBySendingJobSeekerID indicates an expected call of UpdateSendingInterviewDateBySendingJobSeekerID.
func (mr *MockSendingJobSeekerInteractorMockRecorder) UpdateSendingInterviewDateBySendingJobSeekerID(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSendingInterviewDateBySendingJobSeekerID", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).UpdateSendingInterviewDateBySendingJobSeekerID), input)
}

// UpdateSendingJobSeeker mocks base method.
func (m *MockSendingJobSeekerInteractor) UpdateSendingJobSeeker(input interactor.UpdateSendingJobSeekerInput) (interactor.UpdateSendingJobSeekerOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSendingJobSeeker", input)
	ret0, _ := ret[0].(interactor.UpdateSendingJobSeekerOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSendingJobSeeker indicates an expected call of UpdateSendingJobSeeker.
func (mr *MockSendingJobSeekerInteractorMockRecorder) UpdateSendingJobSeeker(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSendingJobSeeker", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).UpdateSendingJobSeeker), input)
}

// UpdateSendingJobSeekerActivityMemo mocks base method.
func (m *MockSendingJobSeekerInteractor) UpdateSendingJobSeekerActivityMemo(input interactor.UpdateSendingJobSeekerActivityMemoInput) (interactor.UpdateSendingJobSeekerActivityMemoOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSendingJobSeekerActivityMemo", input)
	ret0, _ := ret[0].(interactor.UpdateSendingJobSeekerActivityMemoOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSendingJobSeekerActivityMemo indicates an expected call of UpdateSendingJobSeekerActivityMemo.
func (mr *MockSendingJobSeekerInteractorMockRecorder) UpdateSendingJobSeekerActivityMemo(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSendingJobSeekerActivityMemo", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).UpdateSendingJobSeekerActivityMemo), input)
}

// UpdateSendingJobSeekerDocument mocks base method.
func (m *MockSendingJobSeekerInteractor) UpdateSendingJobSeekerDocument(input interactor.UpdateSendingJobSeekerDocumentInput) (interactor.UpdateSendingJobSeekerDocumentOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSendingJobSeekerDocument", input)
	ret0, _ := ret[0].(interactor.UpdateSendingJobSeekerDocumentOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSendingJobSeekerDocument indicates an expected call of UpdateSendingJobSeekerDocument.
func (mr *MockSendingJobSeekerInteractorMockRecorder) UpdateSendingJobSeekerDocument(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSendingJobSeekerDocument", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).UpdateSendingJobSeekerDocument), input)
}

// UpdateSendingJobSeekerPhase mocks base method.
func (m *MockSendingJobSeekerInteractor) UpdateSendingJobSeekerPhase(input interactor.UpdateSendingJobSeekerPhaseInput) (interactor.UpdateSendingJobSeekerPhaseOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSendingJobSeekerPhase", input)
	ret0, _ := ret[0].(interactor.UpdateSendingJobSeekerPhaseOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSendingJobSeekerPhase indicates an expected call of UpdateSendingJobSeekerPhase.
func (mr *MockSendingJobSeekerInteractorMockRecorder) UpdateSendingJobSeekerPhase(input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSendingJobSeekerPhase", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).UpdateSendingJobSeekerPhase), input)
}

// setSendingJobSeekerChildTable mocks base method.
func (m *MockSendingJobSeekerInteractor) setSendingJobSeekerChildTable(sendingJobSeeker *entity.SendingJobSeeker) (*entity.SendingJobSeeker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "setSendingJobSeekerChildTable", sendingJobSeeker)
	ret0, _ := ret[0].(*entity.SendingJobSeeker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// setSendingJobSeekerChildTable indicates an expected call of setSendingJobSeekerChildTable.
func (mr *MockSendingJobSeekerInteractorMockRecorder) setSendingJobSeekerChildTable(sendingJobSeeker any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setSendingJobSeekerChildTable", reflect.TypeOf((*MockSendingJobSeekerInteractor)(nil).setSendingJobSeekerChildTable), sendingJobSeeker)
}
