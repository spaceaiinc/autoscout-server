package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingJobSeekerHandler interface {
	// 汎用系 API
	CreateSendingJobSeeker(param entity.CreateSendingJobSeekerParam) (presenter.Presenter, error)
	CreateSendingJobSeekerFromJobSeeker(param entity.JobSeeker) (presenter.Presenter, error)                            // CRM求職者から送客求職者を作成
	CreateSendingInitialQuestionnaire(param entity.CreateSendingInitialQuestionnaireParam) (presenter.Presenter, error) // 面談前アンケートを登録 (求職者の同意項目、業界、職種、勤務地、質問要望、ファイル（履歴書（原本）、職務経歴書（原本）））

	FirstUpdateSendingJobSeeker(sendingCustomerID uint, param entity.FirstUpdateSendingJobSeekerParam) (presenter.Presenter, error)
	UpdateSendingJobSeeker(sendingJobSeekerID uint, param entity.UpdateSendingJobSeekerParam) (presenter.Presenter, error)
	UpdateSendingJobSeekerPhase(sendingJobSeekerID uint, phase uint) (presenter.Presenter, error)
	UpdateSendingInterviewDateBySendingJobSeekerID(sendingJobSeekerID uint, interviewDate time.Time) (presenter.Presenter, error)
	UpdateSendingJobSeekerActivityMemo(sendingJobSeekerID uint, activityMemo string) (presenter.Presenter, error)
	UpdateIsVewForWating(sendingJobSeekerID uint) (presenter.Presenter, error)
	UpdateIsVewForUnregister(sendingJobSeekerID uint) (presenter.Presenter, error)

	DeleteSendingJobSeeker(sendingJobSeekerID uint) (presenter.Presenter, error)

	GetSendingJobSeekerByID(sendingJobSeekerID uint) (presenter.Presenter, error)
	GetSendingJobSeekerByUUID(sendingJobSeekerID uuid.UUID) (presenter.Presenter, error)
	GetIsNotViewSendingJobSeekerCountByAgentStaffID(agentStaffID uint) (presenter.Presenter, error)

	// { value: number; label: string }[]の形式でリストを取得する
	GetSearchListForSendingJobSeekerManagementByAgentID(agentID uint) (presenter.Presenter, error)

	// 書類 API
	UpdateSendingJobSeekerDocument(param entity.CreateOrUpdateSendingJobSeekerDocumentParam) (presenter.Presenter, error)
	GetSendingJobSeekerDocumentBySendingJobSeekerID(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerResumePDFURL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerResumeOriginURL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerCVPDFURL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerCVOriginURL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerRecommendationPDFURL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerRecommendationOriginURL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerIDPhotoURL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerOtherDocument1URL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerOtherDocument2URL(sendingJobSeekerID uint) (presenter.Presenter, error)
	DeleteSendingJobSeekerOtherDocument3URL(sendingJobSeekerID uint) (presenter.Presenter, error)

	// 送客終了理由　API
	CreateSendingJobSeekerEndStatus(param entity.CreateSendingJobSeekerEndStatusParam) (presenter.Presenter, error) // 送客終了理由を登録 (送客終了理由、送客終了理由詳細
}

type SendingJobSeekerHandlerImpl struct {
	sendingJobSeekerInteractor interactor.SendingJobSeekerInteractor
}

func NewSendingJobSeekerHandlerImpl(epI interactor.SendingJobSeekerInteractor) SendingJobSeekerHandler {
	return &SendingJobSeekerHandlerImpl{
		sendingJobSeekerInteractor: epI,
	}
}

/****************************************************************************************/
// 汎用系 API
//
func (h *SendingJobSeekerHandlerImpl) CreateSendingJobSeeker(param entity.CreateSendingJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.CreateSendingJobSeeker(interactor.CreateSendingJobSeekerInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingCustomerJSONPresenter(responses.NewSendingCustomer(output.SendingCustomer)), nil
}

// CRM求職者から送客求職者を作成
func (h *SendingJobSeekerHandlerImpl) CreateSendingJobSeekerFromJobSeeker(param entity.JobSeeker) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.CreateSendingJobSeekerFromJobSeeker(interactor.CreateSendingJobSeekerFromJobSeekerInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerJSONPresenter(responses.NewSendingJobSeeker(output.SendingJobSeeker)), nil
}

func (h *SendingJobSeekerHandlerImpl) FirstUpdateSendingJobSeeker(sendingCustomerID uint, param entity.FirstUpdateSendingJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.FirstUpdateSendingJobSeeker(interactor.FirstUpdateSendingJobSeekerInput{
		SendingCustomerID: sendingCustomerID,
		UpdateParam:       param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerJSONPresenter(responses.NewSendingJobSeeker(output.SendingJobSeeker)), nil
}

func (h *SendingJobSeekerHandlerImpl) UpdateSendingJobSeeker(sendingJobSeekerID uint, param entity.UpdateSendingJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.UpdateSendingJobSeeker(interactor.UpdateSendingJobSeekerInput{
		SendingJobSeekerID: sendingJobSeekerID,
		UpdateParam:        param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerJSONPresenter(responses.NewSendingJobSeeker(output.SendingJobSeeker)), nil
}

// 送客求職者のフェーズを更新
func (h *SendingJobSeekerHandlerImpl) UpdateSendingJobSeekerPhase(sendingJobSeekerID uint, phase uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.UpdateSendingJobSeekerPhase(interactor.UpdateSendingJobSeekerPhaseInput{
		SendingJobSeekerID: sendingJobSeekerID,
		Phase:              phase,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 送客求職者の面談日時を更新
func (h *SendingJobSeekerHandlerImpl) UpdateSendingInterviewDateBySendingJobSeekerID(SendingJobSeekerID uint, interviewDate time.Time) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.UpdateSendingInterviewDateBySendingJobSeekerID(interactor.UpdateSendingJobSeekerInterviewDateInput{
		SendingJobSeekerID: SendingJobSeekerID,
		InterviewDate:      interviewDate,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) UpdateSendingJobSeekerActivityMemo(sendingJobSeekerID uint, activityMemo string) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.UpdateSendingJobSeekerActivityMemo(interactor.UpdateSendingJobSeekerActivityMemoInput{
		SendingJobSeekerID: sendingJobSeekerID,
		ActivityMemo:       activityMemo,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 面談実施待ちの未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
func (h *SendingJobSeekerHandlerImpl) UpdateIsVewForWating(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.UpdateIsVewForWating(interactor.UpdateIsVewForWatingInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 未登録の未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
func (h *SendingJobSeekerHandlerImpl) UpdateIsVewForUnregister(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.UpdateIsVewForUnregister(interactor.UpdateIsVewForUnregisterInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 送客求職者の削除
func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeeker(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeeker(interactor.DeleteSendingJobSeekerInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) GetSendingJobSeekerByID(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.GetSendingJobSeekerByID(interactor.GetSendingJobSeekerByIDInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerJSONPresenter(responses.NewSendingJobSeeker(output.SendingJobSeeker)), nil
}

func (h *SendingJobSeekerHandlerImpl) GetSendingJobSeekerByUUID(sendingJobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.GetSendingJobSeekerByUUID(interactor.GetSendingJobSeekerByUUIDInput{
		SendingJobSeekerUUID: sendingJobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerJSONPresenter(responses.NewSendingJobSeeker(output.SendingJobSeeker)), nil
}

func (h *SendingJobSeekerHandlerImpl) GetIsNotViewSendingJobSeekerCountByAgentStaffID(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.GetIsNotViewSendingJobSeekerCountByAgentStaffID(interactor.GetIsNotViewSendingJobSeekerCountByAgentStaffIDInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewIsNotViewCountJSONPresenter(responses.NewIsNotViewCount(output.IsViewCount)), nil
}

/****************************************************************************************/
// { value: number; label: string }[]の形式でリストを取得する
//
func (h *SendingJobSeekerHandlerImpl) GetSearchListForSendingJobSeekerManagementByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.GetSearchListForSendingJobSeekerManagementByAgentID(interactor.GetSearchListForSendingJobSeekerManagementByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewLabelListForSendingJobSeekerManagementJSONPresenter(responses.NewLabelListForSendingJobSeekerManagement(output.StaffList, output.SendAgentList, output.SenderList)), nil
}

/****************************************************************************************/
// 書類 API
//
// 送客求職者書類の更新
func (h *SendingJobSeekerHandlerImpl) UpdateSendingJobSeekerDocument(param entity.CreateOrUpdateSendingJobSeekerDocumentParam) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.UpdateSendingJobSeekerDocument(interactor.UpdateSendingJobSeekerDocumentInput{
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerDocumentJSONPresenter(responses.NewSendingJobSeekerDocument(output.SendingJobSeekerDocument)), nil
}

// 送客求職者書類の取得
func (h *SendingJobSeekerHandlerImpl) GetSendingJobSeekerDocumentBySendingJobSeekerID(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.GetSendingJobSeekerDocumentBySendingJobSeekerID(interactor.GetSendingJobSeekerDocumentBySendingJobSeekerIDInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerDocumentJSONPresenter(responses.NewSendingJobSeekerDocument(output.SendingJobSeekerDocument)), nil
}

func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerResumePDFURL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerResumePDFURL(interactor.DeleteSendingJobSeekerResumePDFURLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerResumeOriginURL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerResumeOriginURL(interactor.DeleteSendingJobSeekerResumeOriginURLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerCVPDFURL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerCVPDFURL(interactor.DeleteSendingJobSeekerCVPDFURLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerCVOriginURL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerCVOriginURL(interactor.DeleteSendingJobSeekerCVOriginURLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerRecommendationPDFURL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerRecommendationPDFURL(interactor.DeleteSendingJobSeekerRecommendationPDFURLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerRecommendationOriginURL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerRecommendationOriginURL(interactor.DeleteSendingJobSeekerRecommendationOriginURLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerIDPhotoURL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerIDPhotoURL(interactor.DeleteSendingJobSeekerIDPhotoURLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// その他①を空文字で更新
func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerOtherDocument1URL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerOtherDocument1URL(interactor.DeleteSendingJobSeekerOtherDocument1URLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// その他②を空文字で更新
func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerOtherDocument2URL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerOtherDocument2URL(interactor.DeleteSendingJobSeekerOtherDocument2URLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// その他③を空文字で更新
func (h *SendingJobSeekerHandlerImpl) DeleteSendingJobSeekerOtherDocument3URL(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.DeleteSendingJobSeekerOtherDocument3URL(interactor.DeleteSendingJobSeekerOtherDocument3URLInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
// 送客終了理由 API
//
func (h *SendingJobSeekerHandlerImpl) CreateSendingJobSeekerEndStatus(param entity.CreateSendingJobSeekerEndStatusParam) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.CreateSendingJobSeekerEndStatus(interactor.CreateSendingJobSeekerEndStatusInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
// 面談前アンケート API
//
func (h *SendingJobSeekerHandlerImpl) CreateSendingInitialQuestionnaire(param entity.CreateSendingInitialQuestionnaireParam) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerInteractor.CreateSendingInitialQuestionnaire(interactor.CreateSendingInitialQuestionnaireInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
