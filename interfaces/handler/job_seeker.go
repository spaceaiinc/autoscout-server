package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type JobSeekerHandler interface {
	// 汎用系 API
	CreateJobSeeker(param entity.CreateOrUpdateJobSeekerParam, agentStaffID uint) (presenter.Presenter, error)
	UpdateJobSeeker(param entity.CreateOrUpdateJobSeekerParam, jobSeekerID, agentStaffID uint) (presenter.Presenter, error)
	DeleteJobSeeker(param entity.DeleteJobSeekerParam) (presenter.Presenter, error)
	GetJobSeekerByID(jobSeekerID uint) (presenter.Presenter, error)
	GetJobSeekerByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	GetJobSeekerDocumentByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	GetJobSeekerByTaskGroupUUID(taskGroupUUID uuid.UUID) (presenter.Presenter, error)
	GetJobSeekerListByIDList(idList []uint, agentID uint) (presenter.Presenter, error)
	GetDuplicateJobSeekerList(agentID uint, lastName, firstName, lastFurigana, firstFurigana, email, phoneNumber string) (presenter.Presenter, error)
	GetSelectListForCreateOrUpdateJobSeekerByAgentID(token string, agentID uint) (presenter.Presenter, error)
	UpdateActivityMemoByJobSeekerID(param entity.ActivityMemoParam, jobSeekerID uint) (presenter.Presenter, error)
	UpdateCanViewMatchingJob(param entity.UpdateJobSeekerCanViewMatchingJobParam) (presenter.Presenter, error)

	// 求人検索→求職者検索
	GetSearchJobSeekerListByAgentIDAndType(agentID, pageNumber uint, searchParam entity.SearchJobSeeker, searchType entity.JobSeekerType) (presenter.Presenter, error) // 求人検索→求職者検索(絞り込み)

	// シェア求人検索→自社求職者検索
	GetSearchPublicJobSeekerListByAgentIDAndPage(agentID, pageNumber uint, jobInformationIDList []uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error) // シェア求人検索→自社求職者検索(絞り込み)

	// 絞り込み検索
	GetSearchJobSeekerListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error)
	GetSearchActiveJobSeekerListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error)
	GetSearchAllianceJobSeekerListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error)

	// 求職者資料関連API
	CreateJobSeekerDocument(param entity.CreateOrUpdateJobSeekerDocumentParam) (presenter.Presenter, error)
	UpdateJobSeekerDocument(param entity.CreateOrUpdateJobSeekerDocumentParam) (presenter.Presenter, error)
	UpdateJobSeekerDocumentForTask(param entity.CreateOrUpdateJobSeekerDocumentParam, jobSeekerUUID, jobInformationUUID uuid.UUID) (presenter.Presenter, error)

	GetJobSeekerDocumentByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error)

	// CSV操作系 API
	ImportJobSeekerCSV(param []*entity.JobSeeker, agentID uint) (presenter.Presenter, error) // csvから求職者を登録する関数
	ExportJobSeekerCSV(agentID uint) (string, error)                                         // 求職者一覧をCSVに出力する関数

	// LINE関連API
	UpdateJobSeekerLineID(param entity.UpdateJobSeekerLineIDParam) (presenter.Presenter, error) // 求職者のLINEIDを更新する関数

	// 面談前アンケート関連 API
	CreateInitialQuestionnaire(param entity.CreateInitialQuestionnaireParam) (presenter.Presenter, error) // 面談前アンケートを登録 (求職者の同意項目、業界、職種、勤務地、質問要望、ファイル（履歴書（原本）、職務経歴書（原本）））

	// Admin API
	GetAllJobSeeker(pageNumber uint) (presenter.Presenter, error)

	DeleteJobSeekerResumePDFURL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerResumeOriginURL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerCVPDFURL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerCVOriginURL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerRecommendationPDFURL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerRecommendationOriginURL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerIDPhotoURL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerOtherDocument1URL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerOtherDocument2URL(jobSeekerID uint) (presenter.Presenter, error)
	DeleteJobSeekerOtherDocument3URL(jobSeekerID uint) (presenter.Presenter, error)

	// ゲストページ用　API
	GetJobSeekerForInitialStepByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	GetGuestJobSeekerForByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	GetJobSeekerDesiredForGuestByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	GetJobSeekerAgentIDByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	CheckJobSeekerByUUIDAndName(param entity.CheckJobSeekerByUUIDAndNameParam) (presenter.Presenter, error)
	UpdateJobSeekerPassword(param entity.UpdateJobSeekerPasswordParam) (presenter.Presenter, error) // ゲスト求職者のパスワードを更新
	SendJobSeekerResetPasswordEmail(param entity.SendJobSeekerResetPasswordEmailParam) (presenter.Presenter, error)
	SendJobSeekerContact(param entity.SendJobSeekerContactParam) (presenter.Presenter, error)
	UpdateInterviewDateByJobSeekerID(param entity.UpdateJobSeekerInterviewDateFromGestPageParam) (presenter.Presenter, error)

	// LP
	CreateJobSeekerFromLP(param entity.CreateJobSeekerFromLPParam) (presenter.Presenter, error)
	UpdateJobSeekerPhoneFromLP(param entity.UpdateJobSeekerPhoneFromLPParam) (presenter.Presenter, error)
	UpdateJobSeekerDesiredFromLP(param entity.UpdateJobSeekerDesiredFromLPParam) (presenter.Presenter, error)
	GetJobSeekerLPRegisterStatusByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	SendLPContact(param entity.SendContactFromLPParam) (presenter.Presenter, error)
	SendJobSeekerResetPasswordEmailForLP(param entity.SendJobSeekerResetPasswordEmailFromLPParam) (presenter.Presenter, error)
	ResetPasswordForLP(param entity.ResetPasswordFromLPParam) (presenter.Presenter, error)
	CheckResetPasswordToken(resetPasswordToken string) (presenter.Presenter, error)
}

type JobSeekerHandlerImpl struct {
	jobSeekerInteractor interactor.JobSeekerInteractor
}

func NewJobSeekerHandlerImpl(jsI interactor.JobSeekerInteractor) JobSeekerHandler {
	return &JobSeekerHandlerImpl{
		jobSeekerInteractor: jsI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *JobSeekerHandlerImpl) CreateJobSeeker(param entity.CreateOrUpdateJobSeekerParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.CreateJobSeeker(interactor.CreateJobSeekerInput{
		CreateParam:  param,
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerJSONPresenter(responses.NewJobSeeker(output.JobSeeker)), nil
}

func (h *JobSeekerHandlerImpl) UpdateJobSeeker(param entity.CreateOrUpdateJobSeekerParam, jobSeekerID, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateJobSeeker(interactor.UpdateJobSeekerInput{
		UpdateParam:  param,
		JobSeekerID:  jobSeekerID,
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerJSONPresenter(responses.NewJobSeeker(output.JobSeeker)), nil
}

func (h *JobSeekerHandlerImpl) UpdateActivityMemoByJobSeekerID(param entity.ActivityMemoParam, jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateActivityMemoByJobSeekerID(interactor.UpdateActivityMemoByJobSeekerIDInput{
		Param:       param,
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) UpdateCanViewMatchingJob(param entity.UpdateJobSeekerCanViewMatchingJobParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateCanViewMatchingJob(interactor.UpdateCanViewMatchingJobInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) DeleteJobSeeker(param entity.DeleteJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeeker(interactor.DeleteJobSeekerInput{
		DeleteParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 求職者IDから企業情報を取得する
func (h *JobSeekerHandlerImpl) GetJobSeekerByID(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerByID(interactor.GetJobSeekerByIDInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerJSONPresenter(responses.NewJobSeeker(output.JobSeeker)), nil
}

// 求職者uuidから求職者情報を取得する
func (h *JobSeekerHandlerImpl) GetJobSeekerByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerByUUID(interactor.GetJobSeekerByUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerJSONPresenter(responses.NewJobSeeker(output.JobSeeker)), nil
}

// 求職者uuidから求職者の応募書類情報を取得する
func (h *JobSeekerHandlerImpl) GetJobSeekerDocumentByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerDocumentByUUID(interactor.GetJobSeekerDocumentByUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerDocumentJSONPresenter(responses.NewJobSeekerDocument(output.Document)), nil
}

// タスクグループのuuidから求職者情報を取得する
func (h *JobSeekerHandlerImpl) GetJobSeekerByTaskGroupUUID(taskGroupUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerByTaskGroupUUID(interactor.GetJobSeekerByTaskGroupUUIDInput{
		TaskGroupUUID: taskGroupUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerJSONPresenter(responses.NewJobSeeker(output.JobSeeker)), nil
}

func (h *JobSeekerHandlerImpl) GetJobSeekerListByIDList(idList []uint, agentID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerListByIDList(interactor.GetJobSeekerListByIDListInput{
		IDList:  idList,
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListJSONPresenter(responses.NewJobSeekerList(output.JobSeekerList)), nil
}

// クエリパラム（last_name, first_name, last_furigana, first_furigana, email, phone_number）に合致する求職者情報を取得
func (h *JobSeekerHandlerImpl) GetDuplicateJobSeekerList(agentID uint, lastName, firstName, lastFurigana, firstFurigana, email, phoneNumber string) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetDuplicateJobSeekerList(interactor.GetDuplicateJobSeekerListInput{
		AgentID:       agentID,
		LastName:      lastName,
		FirstName:     firstName,
		LastFurigana:  lastFurigana,
		FirstFurigana: firstFurigana,
		Email:         email,
		PhoneNumber:   phoneNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListJSONPresenter(responses.NewJobSeekerList(output.JobSeekerList)), nil
}

// 求職者編集ページで使用されるセレクトボックスに必要なデータを取得（自社の担当者一覧、アライアンスエージェント一覧、流入経路）
func (h *JobSeekerHandlerImpl) GetSelectListForCreateOrUpdateJobSeekerByAgentID(token string, agentID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetSelectListForCreateOrUpdateJobSeekerByAgentID(interactor.GetSelectListForCreateOrUpdateJobSeekerByAgentIDInput{
		Token:   token,
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSelectListForCreateOrUpdateJobSeekerJSONPresenter(responses.NewSelectListForCreateOrUpdateJobSeeker(output.AgentStaffList, output.AgentInflowChannelOptionList, output.AllianceAgentList)), nil
}

/****************************************************************************************/
/// 絞り込み検索 API

// エージェントIDとクエリパラムで求職者一覧を絞り込み
func (h *JobSeekerHandlerImpl) GetSearchJobSeekerListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetSearchJobSeekerListByAgentID(interactor.GetSearchJobSeekerListByAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDList(output.JobSeekerList, output.MaxPageNumber, output.IDList)), nil
}

func (h *JobSeekerHandlerImpl) GetSearchActiveJobSeekerListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetSearchActiveJobSeekerListByAgentID(interactor.GetSearchActiveJobSeekerListByAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDList(output.JobSeekerList, output.MaxPageNumber, output.IDList)), nil
}

// エージェントIDとクエリパラムで求職者一覧を絞り込み
func (h *JobSeekerHandlerImpl) GetSearchAllianceJobSeekerListByAgentID(agentID, pageNumber uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetSearchAllianceJobSeekerListByAgentID(interactor.GetSearchAllianceJobSeekerListByAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDList(output.JobSeekerList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
// 求人検索→求職者検索 API
//
func (h *JobSeekerHandlerImpl) GetSearchJobSeekerListByAgentIDAndType(agentID, pageNumber uint, searchParam entity.SearchJobSeeker, searchType entity.JobSeekerType) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetSearchJobSeekerListByAgentIDAndType(interactor.GetSearchJobSeekerListByAgentIDAndTypeInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
		Type:        searchType,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListAndListCountJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDListAndListCount(
		output.JobSeekerList, output.MaxPageNumber, output.IDList, output.AllCount, output.OwnCount, output.AllianceCount,
	)), nil
}

/****************************************************************************************/
// シェア求人検索→自社求職者検索(絞り込み)
//
func (h *JobSeekerHandlerImpl) GetSearchPublicJobSeekerListByAgentIDAndPage(agentID, pageNumber uint, jobInformationIDList []uint, searchParam entity.SearchJobSeeker) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetSearchPublicJobSeekerListByAgentIDAndPage(interactor.GetSearchPublicJobSeekerListByAgentIDAndPageInput{
		AgentID:              agentID,
		PageNumber:           pageNumber,
		JobInformationIDList: jobInformationIDList,
		SearchParam:          searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDList(output.JobSeekerList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
/// 資料関連API
//
// 求職者資料情報の登録
func (h *JobSeekerHandlerImpl) CreateJobSeekerDocument(param entity.CreateOrUpdateJobSeekerDocumentParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.CreateJobSeekerDocument(interactor.CreateJobSeekerDocumentInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerDocumentJSONPresenter(responses.NewJobSeekerDocument(output.JobSeekerDocument)), nil
}

// 求職者資料情報の更新
func (h *JobSeekerHandlerImpl) UpdateJobSeekerDocument(param entity.CreateOrUpdateJobSeekerDocumentParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateJobSeekerDocument(interactor.UpdateJobSeekerDocumentInput{
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerDocumentJSONPresenter(responses.NewJobSeekerDocument(output.JobSeekerDocument)), nil
}

func (h *JobSeekerHandlerImpl) UpdateJobSeekerDocumentForTask(param entity.CreateOrUpdateJobSeekerDocumentParam, jobSeekerUUID, jobInformationUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateJobSeekerDocumentForTask(interactor.UpdateJobSeekerDocumentForTaskInput{
		UpdateParam:        param,
		JobSeekerUUID:      jobSeekerUUID,
		JobInformationUUID: jobInformationUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerDocumentJSONPresenter(responses.NewJobSeekerDocument(output.JobSeekerDocument)), nil
}

// 求職者資料情報の取得
func (h *JobSeekerHandlerImpl) GetJobSeekerDocumentByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerDocumentByJobSeekerID(interactor.GetJobSeekerDocumentByJobSeekerIDInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerDocumentJSONPresenter(responses.NewJobSeekerDocument(output.JobSeekerDocument)), nil
}

/****************************************************************************************/
/// CSV操作 API
//
//csvファイルを読み込む
func (h *JobSeekerHandlerImpl) ImportJobSeekerCSV(param []*entity.JobSeeker, agentID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.ImportJobSeekerCSV(interactor.ImportJobSeekerCSVInput{
		CreateParam: param,
		AgentID:     agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewMissedRecordsAndOKJSONPresenter(responses.NewMissedRecordsAndOK(output.MissedRecords, output.OK)), nil
}

// csvファイルを出力する
func (h *JobSeekerHandlerImpl) ExportJobSeekerCSV(agentID uint) (string, error) {
	output, err := h.jobSeekerInteractor.ExportJobSeekerCSV(interactor.ExportJobSeekerCSVInput{
		AgentID: agentID,
	})

	if err != nil {
		return "", err
	}

	return output.FilePath.PathName, nil
}

/****************************************************************************************/
///LINE関連 API
//
// LINE連携情報の登録
func (h *JobSeekerHandlerImpl) UpdateJobSeekerLineID(param entity.UpdateJobSeekerLineIDParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateJobSeekerLineID(interactor.UpdateJobSeekerLineIDInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
// 面談前アンケート関連 API
//
// 面談前アンケートを登録 (求職者の同意項目、業界、職種、勤務地、質問要望、ファイル（履歴書（原本）、職務経歴書（原本）））
func (h *JobSeekerHandlerImpl) CreateInitialQuestionnaire(param entity.CreateInitialQuestionnaireParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.CreateInitialQuestionnaire(interactor.CreateInitialQuestionnaireInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
/// Admin API

// すべての企業情報を取得する
func (h *JobSeekerHandlerImpl) GetAllJobSeeker(pageNumber uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetAllJobSeeker(interactor.GetAllJobSeekerInput{
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDList(output.JobSeekerList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/

func (h *JobSeekerHandlerImpl) DeleteJobSeekerResumePDFURL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerResumePDFURL(interactor.DeleteJobSeekerResumePDFURLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) DeleteJobSeekerResumeOriginURL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerResumeOriginURL(interactor.DeleteJobSeekerResumeOriginURLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) DeleteJobSeekerCVPDFURL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerCVPDFURL(interactor.DeleteJobSeekerCVPDFURLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) DeleteJobSeekerCVOriginURL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerCVOriginURL(interactor.DeleteJobSeekerCVOriginURLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) DeleteJobSeekerRecommendationPDFURL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerRecommendationPDFURL(interactor.DeleteJobSeekerRecommendationPDFURLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) DeleteJobSeekerRecommendationOriginURL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerRecommendationOriginURL(interactor.DeleteJobSeekerRecommendationOriginURLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) DeleteJobSeekerIDPhotoURL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerIDPhotoURL(interactor.DeleteJobSeekerIDPhotoURLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// その他①を空文字で更新
func (h *JobSeekerHandlerImpl) DeleteJobSeekerOtherDocument1URL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerOtherDocument1URL(interactor.DeleteJobSeekerOtherDocument1URLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// その他②を空文字で更新
func (h *JobSeekerHandlerImpl) DeleteJobSeekerOtherDocument2URL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerOtherDocument2URL(interactor.DeleteJobSeekerOtherDocument2URLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// その他③を空文字で更新
func (h *JobSeekerHandlerImpl) DeleteJobSeekerOtherDocument3URL(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.DeleteJobSeekerOtherDocument3URL(interactor.DeleteJobSeekerOtherDocument3URLInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
// ゲストページ用 API
//
// 求職者uuidから求職者情報を取得する
func (h *JobSeekerHandlerImpl) GetJobSeekerForInitialStepByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerForInitialStepByUUID(interactor.GetJobSeekerForInitialStepByUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerForGuestJSONPresenter(responses.NewJobSeekerForGuest(output.JobSeekerForGuest)), nil
}

func (h *JobSeekerHandlerImpl) GetGuestJobSeekerForByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetGuestJobSeekerForByUUID(interactor.GetGuestJobSeekerForByUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewGestJobSeekerUserSessionJSONPresenter(responses.NewGestJobSeekerUserSession(output.User)), nil
}

func (h *JobSeekerHandlerImpl) GetJobSeekerDesiredForGuestByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerDesiredForGuestByUUID(interactor.GetJobSeekerDesiredForGuestByUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerDesiredForGuestJSONPresenter(responses.NewJobSeekerDesiredForGuest(output.JobSeekerDesired)), nil
}

// 求職者uuidからエージェントIDを取得する
func (h *JobSeekerHandlerImpl) GetJobSeekerAgentIDByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerAgentIDByUUID(interactor.GetJobSeekerAgentIDByUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerAgentIDJSONPresenter(responses.NewJobSeekerAgentID(output.AgentID)), nil
}

func (h *JobSeekerHandlerImpl) CheckJobSeekerByUUIDAndName(param entity.CheckJobSeekerByUUIDAndNameParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.CheckJobSeekerByUUIDAndName(interactor.CheckJobSeekerByUUIDAndNameInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) UpdateJobSeekerPassword(param entity.UpdateJobSeekerPasswordParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateJobSeekerPassword(interactor.UpdateJobSeekerPasswordInput{
		Param: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) SendJobSeekerResetPasswordEmail(param entity.SendJobSeekerResetPasswordEmailParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.SendJobSeekerResetPasswordEmail(interactor.SendJobSeekerResetPasswordEmailInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) SendJobSeekerContact(param entity.SendJobSeekerContactParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.SendJobSeekerContact(interactor.SendJobSeekerContactInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}
	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) UpdateInterviewDateByJobSeekerID(param entity.UpdateJobSeekerInterviewDateFromGestPageParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateInterviewDateByJobSeekerID(interactor.UpdateInterviewDateByJobSeekerIDInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}
	return presenter.NewGestJobSeekerUserSessionJSONPresenter(responses.NewGestJobSeekerUserSession(output.User)), nil
}

/****************************************************************************************/
// LP用 API
//
func (h *JobSeekerHandlerImpl) CreateJobSeekerFromLP(param entity.CreateJobSeekerFromLPParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.CreateJobSeekerFromLP(interactor.CreateJobSeekerFromLPInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerUUIDJSONPresenter(responses.NewJobSeekerUUID(output.UUID)), nil
}

// LPから求職者の電話番号を更新
func (h *JobSeekerHandlerImpl) UpdateJobSeekerPhoneFromLP(param entity.UpdateJobSeekerPhoneFromLPParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateJobSeekerPhoneFromLP(interactor.UpdateJobSeekerPhoneFromLPInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerLoginFromLPJSONPresenter(responses.NewJobSeekerLoginFromLP(output.UUID, output.LogintToken)), nil
}

// LPから求職者の希望条件を更新
func (h *JobSeekerHandlerImpl) UpdateJobSeekerDesiredFromLP(param entity.UpdateJobSeekerDesiredFromLPParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.UpdateJobSeekerDesiredFromLP(interactor.UpdateJobSeekerDesiredFromLPInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerLoginFromLPJSONPresenter(responses.NewJobSeekerLoginFromLP(output.UUID, output.LogintToken)), nil
}

// LPの登録状況を取得する
func (h *JobSeekerHandlerImpl) GetJobSeekerLPRegisterStatusByUUID(jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.GetJobSeekerLPRegisterStatusByUUID(interactor.GetJobSeekerLPRegisterStatusByUUIDInput{
		JobSeekerUUID: jobSeekerUUID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerRegisterStatusJSONPresenter(responses.NewJobSeekerRegisterStatus(output.JobSeekerRegisterStatus)), nil
}

func (h *JobSeekerHandlerImpl) SendLPContact(param entity.SendContactFromLPParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.SendLPContact(interactor.SendLPContactInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}
	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *JobSeekerHandlerImpl) SendJobSeekerResetPasswordEmailForLP(param entity.SendJobSeekerResetPasswordEmailFromLPParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.SendJobSeekerResetPasswordEmailForLP(interactor.SendJobSeekerResetPasswordEmailForLPInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// パスワードリセット
func (h *JobSeekerHandlerImpl) ResetPasswordForLP(param entity.ResetPasswordFromLPParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.ResetPasswordForLP(interactor.ResetPasswordForLPInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// LPからパスワードトークンの有効性を確認
func (h *JobSeekerHandlerImpl) CheckResetPasswordToken(resetPasswordToken string) (presenter.Presenter, error) {
	output, err := h.jobSeekerInteractor.CheckResetPasswordToken(interactor.CheckResetPasswordTokenInput{
		ResetPasswordToken: resetPasswordToken,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
