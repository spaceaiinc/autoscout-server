package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingJobInformationHandler interface {
	// 汎用系 API
	CreateSendingJobInformation(param entity.CreateSendingJobInformationParam) (presenter.Presenter, error)
	UpdateSendingJobInformation(param entity.UpdateSendingJobInformationParam, jobInformationID uint) (presenter.Presenter, error)
	DeleteSendingJobInformation(jobInformationID uint) (presenter.Presenter, error)
	GetSendingJobInformationByID(jobInformationID uint) (presenter.Presenter, error)
	GetSendingJobInformationByUUID(uuid uuid.UUID) (presenter.Presenter, error)
	GetJobListingBySendingJobInformationUUID(uuid uuid.UUID) (presenter.Presenter, error)
	GetSendingJobInformationListBySendingEnterpriseID(sendingEnterpriseID uint) (presenter.Presenter, error)

	// ページ用API
	GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(sendingJobSeekerID, sendingEnterpriseID, pageNumber uint) (presenter.Presenter, error)                                                // 送客先を探すページ
	GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(sendingJobSeekerID, sendingEnterpriseID, pageNumber uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error) // 送客先を探すページ

	// 求人の絞り込み検索
	// GetSearchActiveSendingJobInformationListByAgentID(agentID, pageNumber uint, searchParam entity.SearchSendingJobInformation) (presenter.Presenter, error)

	// CSV API
	ImportSendingJobInformationCSV(paramList []*entity.SendingJobInformation, missedRecords []uint) (presenter.Presenter, error)
}

type SendingJobInformationHandlerImpl struct {
	jobInformationInteractor interactor.SendingJobInformationInteractor
}

func NewSendingJobInformationHandlerImpl(jiI interactor.SendingJobInformationInteractor) SendingJobInformationHandler {
	return &SendingJobInformationHandlerImpl{
		jobInformationInteractor: jiI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *SendingJobInformationHandlerImpl) CreateSendingJobInformation(param entity.CreateSendingJobInformationParam) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.CreateSendingJobInformation(interactor.CreateSendingJobInformationInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobInformationJSONPresenter(responses.NewSendingJobInformation(output.SendingJobInformation)), nil
}

func (h *SendingJobInformationHandlerImpl) UpdateSendingJobInformation(param entity.UpdateSendingJobInformationParam, jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.UpdateSendingJobInformation(interactor.UpdateSendingJobInformationInput{
		SendingJobInformationID: jobInformationID,
		UpdateParam:             param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobInformationJSONPresenter(responses.NewSendingJobInformation(output.SendingJobInformation)), nil
}

func (h *SendingJobInformationHandlerImpl) DeleteSendingJobInformation(jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.DeleteSendingJobInformation(interactor.DeleteSendingJobInformationInput{
		SendingJobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 求人IDから求人情報を取得する
func (h *SendingJobInformationHandlerImpl) GetSendingJobInformationByID(jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSendingJobInformationByID(interactor.GetSendingJobInformationByIDInput{
		SendingJobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobInformationJSONPresenter(responses.NewSendingJobInformation(output.SendingJobInformation)), nil
}

// 求人のuuidから求人情報を取得する
func (h *SendingJobInformationHandlerImpl) GetSendingJobInformationByUUID(uuid uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSendingJobInformationByUUID(interactor.GetSendingJobInformationByUUIDInput{
		UUID: uuid,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobInformationJSONPresenter(responses.NewSendingJobInformation(output.SendingJobInformation)), nil
}

// 求人のuuidから求人情報を取得する
func (h *SendingJobInformationHandlerImpl) GetJobListingBySendingJobInformationUUID(uuid uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobListingBySendingJobInformationUUID(interactor.GetJobListingBySendingJobInformationUUIDInput{
		UUID: uuid,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobListingForSendingJSONPresenter(responses.NewJobListingForSending(output.JobListing)), nil
}

// 企業IDから求人一覧を取得する
func (h *SendingJobInformationHandlerImpl) GetSendingJobInformationListBySendingEnterpriseID(sendingEnterpriseID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSendingJobInformationListBySendingEnterpriseID(interactor.GetSendingJobInformationListBySendingEnterpriseIDInput{
		SendingEnterpriseID: sendingEnterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobInformationListJSONPresenter(responses.NewSendingJobInformationList(output.SendingJobInformationList)), nil
}

/****************************************************************************************/
// ページ用
//
// まだ送客していない送客先とその送客先が保有する求人数の一覧を取得するapi
func (h *SendingJobInformationHandlerImpl) GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(sendingJobSeekerID, sendingEnterpriseID, pageNumber uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(interactor.GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetInput{
		SendingJobSeekerID:  sendingJobSeekerID,
		SendingEnterpriseID: sendingEnterpriseID,
		PageNumber:          pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseListAndJobInformationListAndMaxPageJSONPresenter(responses.NewSendingEnterpriseListAndJobInformationListAndMaxPage(output.SendingEnterpriseList, output.SendingJobInformationList, output.MaxPageNumber)), nil
}

// 絞り込み まだ送客していない送客先とその送客先が保有する求人数の一覧を取得するapi
func (h *SendingJobInformationHandlerImpl) GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(sendingJobSeekerID, sendingEnterpriseID, pageNumber uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(interactor.GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetInput{
		SendingJobSeekerID:  sendingJobSeekerID,
		SendingEnterpriseID: sendingEnterpriseID,
		PageNumber:          pageNumber,
		SearchParam:         searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingEnterpriseListAndJobInformationListAndMaxPageJSONPresenter(responses.NewSendingEnterpriseListAndJobInformationListAndMaxPage(output.SendingEnterpriseList, output.SendingJobInformationList, output.MaxPageNumber)), nil
}

/****************************************************************************************/
// csv
//
// csvファイルを読み込む
func (h *SendingJobInformationHandlerImpl) ImportSendingJobInformationCSV(paramList []*entity.SendingJobInformation, missedRecords []uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.ImportSendingJobInformationCSV(interactor.ImportSendingJobInformationCSVInput{
		CreateParamList: paramList,
		MissedRecords:   missedRecords,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMissedRecordsAndOKJSONPresenter(responses.NewMissedRecordsAndOK(output.MissedRecords, output.OK)), nil
}
