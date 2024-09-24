package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type JobInformationHandler interface {
	// 汎用系 API
	CreateJobInformation(param entity.CreateJobInformationParam, billingAddressID uint) (presenter.Presenter, error)
	UpdateJobInformation(param entity.UpdateJobInformationParam, jobInformationID uint) (presenter.Presenter, error)
	DeleteJobInformation(jobInformationID uint) (presenter.Presenter, error)
	GetJobInformationByID(jobInformationID uint) (presenter.Presenter, error)
	GetJobListingByJobInformationUUID(uuid uuid.UUID) (presenter.Presenter, error)
	GetJobInformationByUUID(uuid uuid.UUID) (presenter.Presenter, error)
	GetJobListingForJobSeeker(jobInformationUUID, jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	GetJobInformationListByBillingAddressID(billingAddressID uint) (presenter.Presenter, error)
	GetJobInformationListByEnterpriseID(enterpriseID uint) (presenter.Presenter, error)
	GetJobInformationListByAgentID(agentID uint) (presenter.Presenter, error)
	GetSelectionFlowPatternListByJobInformationID(jobInformationID uint) (presenter.Presenter, error)
	GetOpenSelectionFlowPatternListByJobInformationID(jobInformationID uint) (presenter.Presenter, error)
	GetSelectionFlowPatternByID(selectionFlowID uint) (presenter.Presenter, error)
	CreateSelectionFlowPattern(param entity.CreateAndUpdateSelectionFlowPatternParam) (presenter.Presenter, error)
	UpdateSelectionFlowPattern(param entity.CreateAndUpdateSelectionFlowPatternParam, selectionFlowID uint) (presenter.Presenter, error)
	DeltedSelectionFlowPattern(selectionFlowID uint) (presenter.Presenter, error)
	GetJobInformationListByIDList(idList []uint) (presenter.Presenter, error)
	GetJobListingListByJobSeekerUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)

	// 求職者検索→求人検索
	GetJobInformationListByAgentIDAndType(agentID, pageNumber uint, searchType entity.JobInformationType) (presenter.Presenter, error) // 求職者検索→求人検索
	GetSearchJobInformationListByAgentIDAndType(agentID, pageNumber uint, searchParam entity.SearchJobInformation, searchType entity.JobInformationType) (presenter.Presenter, error)

	// シェア求職者検索→自社求人検索
	GetSearchPublicJobInformationListByAgentIDAndPage(agentID, pageNumber uint, jobSeekerIDList []uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error) // シェア求職者検索→自社求人検索(絞り込み検索)

	// 求人の絞り込み検索
	GetSearchActiveJobInformationListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error)
	GetSearchJobInformationListByOtherAgentID(agentID, pageNumber uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error) // アライアンス求人の絞り込み検索を取得する

	// LP用 API
	GetSearchJobInformationCountByLPDiagnosis(searchParam entity.DiagnosisParam) (presenter.Presenter, error)
	GetSearchJobListingListByJobSeekerUUID(searchParam entity.SearchMatchingJobListParam) (presenter.Presenter, error)
	GetJobInformationListForDiagnosis() (presenter.Presenter, error)
	GetJobListingListAndJobSeekerDesiredForDiagnosis(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	// 求職者のエントリー希望と興味あり求人の取得
	GetJobListingListByJobSeekerUUIDAndInterestedType(param entity.InterestedTypeJobListParam) (presenter.Presenter, error)

	// admin API
	GetAllJobInformation(pageNumber uint) (presenter.Presenter, error)
}

type JobInformationHandlerImpl struct {
	jobInformationInteractor interactor.JobInformationInteractor
}

func NewJobInformationHandlerImpl(jiI interactor.JobInformationInteractor) JobInformationHandler {
	return &JobInformationHandlerImpl{
		jobInformationInteractor: jiI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *JobInformationHandlerImpl) CreateJobInformation(param entity.CreateJobInformationParam, billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.CreateJobInformation(interactor.CreateJobInformationInput{
		BillingAddressID: billingAddressID,
		CreateParam:      param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationJSONPresenter(responses.NewJobInformation(output.JobInformation)), nil
}

func (h *JobInformationHandlerImpl) UpdateJobInformation(param entity.UpdateJobInformationParam, jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.UpdateJobInformation(interactor.UpdateJobInformationInput{
		JobInformationID: jobInformationID,
		UpdateParam:      param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationJSONPresenter(responses.NewJobInformation(output.JobInformation)), nil
}

func (h *JobInformationHandlerImpl) DeleteJobInformation(jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.DeleteJobInformation(interactor.DeleteJobInformationInput{
		JobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 求人IDから求人情報を取得する
func (h *JobInformationHandlerImpl) GetJobInformationByID(jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationByID(interactor.GetJobInformationByIDInput{
		JobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationJSONPresenter(responses.NewJobInformation(output.JobInformation)), nil
}

// 求人のuuidから求人情報を取得する
func (h *JobInformationHandlerImpl) GetJobInformationByUUID(uuid uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationByUUID(interactor.GetJobInformationByUUIDInput{
		UUID: uuid,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationJSONPresenter(responses.NewJobInformation(output.JobInformation)), nil
}

// 求人のuuidから求人情報を取得する
func (h *JobInformationHandlerImpl) GetJobListingByJobInformationUUID(uuid uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobListingByJobInformationUUID(interactor.GetJobListingByJobInformationUUIDInput{
		UUID: uuid,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobListingJSONPresenter(responses.NewJobListing(output.JobListing)), nil
}

// 求職者が確認する求人票情報を取得（求人票 + タスクに紐づいた選考情報）
func (h *JobInformationHandlerImpl) GetJobListingForJobSeeker(jobInformationUUID, jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobListingForJobSeeker(interactor.GetJobListingForJobSeekerInput{
		JobInformationUUID: jobInformationUUID,
		JobSeekerUUID:      jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobListingJSONPresenter(responses.NewJobListing(output.JobListing)), nil
}

// 請求先IDから求人情報一覧を取得する
func (h *JobInformationHandlerImpl) GetJobInformationListByBillingAddressID(billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationListByBillingAddressID(interactor.GetJobInformationListByBillingAddressIDInput{
		BillingAddressID: billingAddressID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListJSONPresenter(responses.NewJobInformationList(output.JobInformationList)), nil
}

// 企業IDから求人一覧を取得する
func (h *JobInformationHandlerImpl) GetJobInformationListByEnterpriseID(enterpriseID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationListByEnterpriseID(interactor.GetJobInformationListByEnterpriseIDInput{
		EnterpriseID: enterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListJSONPresenter(responses.NewJobInformationList(output.JobInformationList)), nil
}

// agentIDを使って求人情報一覧を取得する関数
func (h *JobInformationHandlerImpl) GetJobInformationListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationListByAgentID(interactor.GetJobInformationListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListJSONPresenter(responses.NewJobInformationList(output.JobInformationList)), nil
}

// 求人IDを使って選考フローパターンの一覧を取得する
func (h *JobInformationHandlerImpl) GetSelectionFlowPatternListByJobInformationID(jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSelectionFlowPatternListByJobInformationID(interactor.GetSelectionFlowPatternListByJobInformationIDInput{
		JobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionFlowPatternListJSONPresenter(responses.NewSelectionFlowPatternList(output.SelectionFlowPatternList)), nil
}

// 求人IDを使って選考フローパターンの一覧を取得する
func (h *JobInformationHandlerImpl) GetOpenSelectionFlowPatternListByJobInformationID(jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetOpenSelectionFlowPatternListByJobInformationID(interactor.GetOpenSelectionFlowPatternListByJobInformationIDInput{
		JobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionFlowPatternListJSONPresenter(responses.NewSelectionFlowPatternList(output.SelectionFlowPatternList)), nil
}

// 求人IDを使って選考フローパターンの一覧を取得する
func (h *JobInformationHandlerImpl) GetSelectionFlowPatternByID(selectionFlowID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSelectionFlowPatternByID(interactor.GetSelectionFlowPatternByIDInput{
		SelectionFlowID: selectionFlowID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionFlowPatternJSONPresenter(responses.NewSelectionFlowPattern(output.SelectionFlowPattern)), nil
}

// 選考フローパターンの作成
func (h *JobInformationHandlerImpl) CreateSelectionFlowPattern(param entity.CreateAndUpdateSelectionFlowPatternParam) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.CreateSelectionFlowPattern(interactor.CreateSelectionFlowPatternInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionFlowPatternJSONPresenter(responses.NewSelectionFlowPattern(output.SelectionFlowPattern)), nil
}

// 選考フローパターンの更新
func (h *JobInformationHandlerImpl) UpdateSelectionFlowPattern(param entity.CreateAndUpdateSelectionFlowPatternParam, selectionFlowID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.UpdateSelectionFlowPattern(interactor.UpdateSelectionFlowPatternInput{
		UpdateParam:     param,
		SelectionFlowID: selectionFlowID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectionFlowPatternJSONPresenter(responses.NewSelectionFlowPattern(output.SelectionFlowPattern)), nil
}

// 選考フローパターンの更新
func (h *JobInformationHandlerImpl) DeltedSelectionFlowPattern(selectionFlowID uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.DeltedSelectionFlowPattern(interactor.DeltedSelectionFlowPatternInput{
		SelectionFlowID: selectionFlowID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobInformationHandlerImpl) GetJobInformationListByIDList(idList []uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationListByIDList(interactor.GetJobInformationListByIDListInput{
		IDList: idList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListJSONPresenter(responses.NewJobInformationList(output.JobInformationList)), nil
}

func (h *JobInformationHandlerImpl) GetJobListingListByJobSeekerUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobListingListByJobSeekerUUID(interactor.GetJobListingListByJobSeekerUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobListingListForJobSeekerJSONPresenter(responses.NewJobListingListForJobSeeker(output.NotYetEntryJobListingList, output.AcceptJobOfferJobListingList, output.HoldJobOfferJobListingList, output.SelectionJobListingList, output.EndJobListingList)), nil
}

/****************************************************************************************/
// 求職者検索→求人検索 API
//
func (h *JobInformationHandlerImpl) GetJobInformationListByAgentIDAndType(agentID, pageNumber uint, searchType entity.JobInformationType) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationListByAgentIDAndType(interactor.GetJobInformationListByAgentIDAndTypeInput{
		AgentID:    agentID,
		PageNumber: pageNumber,
		Type:       searchType,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListAndMaxPageAndIDListAndListCountJSONPresenter(responses.NewJobInformationListAndMaxPageAndIDListAndListCount(
		output.JobInformationList, output.MaxPageNumber, output.IDList, output.AllCount, output.OwnCount, output.AllianceCount,
	)), nil
}

func (h *JobInformationHandlerImpl) GetSearchJobInformationListByAgentIDAndType(agentID, pageNumber uint, searchParam entity.SearchJobInformation, searchType entity.JobInformationType) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSearchJobInformationListByAgentIDAndType(interactor.GetSearchJobInformationListByAgentIDAndTypeInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
		Type:        searchType,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListAndMaxPageAndIDListAndListCountJSONPresenter(responses.NewJobInformationListAndMaxPageAndIDListAndListCount(
		output.JobInformationList, output.MaxPageNumber, output.IDList, output.AllCount, output.OwnCount, output.AllianceCount,
	)), nil
}

/****************************************************************************************/
// シェア求職者検索→自社求人検索(絞り込み)
//
func (h *JobInformationHandlerImpl) GetSearchPublicJobInformationListByAgentIDAndPage(agentID, pageNumber uint, jobSeekerIDList []uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSearchPublicJobInformationListByAgentIDAndPage(interactor.GetSearchPublicJobInformationListByAgentIDAndPageInput{
		AgentID:         agentID,
		PageNumber:      pageNumber,
		JobSeekerIDList: jobSeekerIDList,
		SearchParam:     searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListAndMaxPageAndIDListJSONPresenter(responses.NewJobInformationListAndMaxPageAndIDList(output.JobInformationList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
/// 求人の絞り込み検索
func (h *JobInformationHandlerImpl) GetSearchActiveJobInformationListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSearchActiveJobInformationListByAgentID(interactor.GetSearchActiveJobInformationListByAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListAndMaxPageAndIDListJSONPresenter(responses.NewJobInformationListAndMaxPageAndIDList(output.JobInformationList, output.MaxPageNumber, output.IDList)), nil
}

// アライアンスエージェント求人の絞り込み検索
func (h *JobInformationHandlerImpl) GetSearchJobInformationListByOtherAgentID(agentID, pageNumber uint, searchParam entity.SearchJobInformation) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSearchJobInformationListByOtherAgentID(interactor.GetSearchJobInformationListByOtherAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListAndMaxPageAndIDListJSONPresenter(responses.NewJobInformationListAndMaxPageAndIDList(output.JobInformationList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
// LP用 API
func (h *JobInformationHandlerImpl) GetSearchJobInformationCountByLPDiagnosis(searchParam entity.DiagnosisParam) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSearchJobInformationCountByLPDiagnosis(interactor.GetSearchJobInformationCountByLPDiagnosisInput{
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationCountAndIncomeJSONPresenter(responses.NewJobInformationCountAndIncome(output.Count, output.GuaranteedInterviewCount, output.ExpectedIncome)), nil
}

func (h *JobInformationHandlerImpl) GetSearchJobListingListByJobSeekerUUID(searchParam entity.SearchMatchingJobListParam) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetSearchJobListingListByJobSeekerUUID(interactor.GetSearchJobListingListByJobSeekerUUIDInput{
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobListingListAndCountJSONPresenter(responses.NewJobListingListAndCount(
		output.JobListingList,
		output.Count,
		output.CountOfGuaranteedInterview,
		output.ExpectedUnderIncome,
		output.ExpectedOverIncome,
		output.MaxPageNumber,
	)), nil
}

func (h *JobInformationHandlerImpl) GetJobInformationListForDiagnosis() (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobInformationListForDiagnosis()

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListForDiagnosisJSONPresenter(responses.NewJobInformationListForDiagnosis(output.JobInformationList)), nil
}

func (h *JobInformationHandlerImpl) GetJobListingListAndJobSeekerDesiredForDiagnosis(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobListingListAndJobSeekerDesiredForDiagnosis(interactor.GetJobListingListAndJobSeekerDesiredForDiagnosisInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobListingListAndJobSeekerDesiredJSONPresenter(responses.NewJobListingListAndJobSeekerDesired(output.JobListingList, output.JobSeekerDesired)), nil
}

// 求職者のエントリー希望と興味あり求人の取得
func (h *JobInformationHandlerImpl) GetJobListingListByJobSeekerUUIDAndInterestedType(param entity.InterestedTypeJobListParam) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetJobListingListByJobSeekerUUIDAndInterestedType(interactor.GetJobListingListByJobSeekerUUIDAndInterestedTypeInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobListingListAndMaxPageJSONPresenter(responses.NewJobListingListAndMaxPage(output.JobListingList, output.MaxPageNumber)), nil
}

/****************************************************************************************/
/// Admin API

// すべての求人情報を取得する
func (h *JobInformationHandlerImpl) GetAllJobInformation(pageNumber uint) (presenter.Presenter, error) {
	output, err := h.jobInformationInteractor.GetAllJobInformation(interactor.GetAllJobInformationInput{
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobInformationListAndMaxPageAndIDListJSONPresenter(responses.NewJobInformationListAndMaxPageAndIDList(output.JobInformationList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
