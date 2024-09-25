package handler

import (
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
	"gopkg.in/guregu/null.v4"
)

type TaskHandler interface {
	// タスクを始める関数（打診）
	// Memo: 旧版の為、新しいapiの動作確認後に削除（が対応）
	SoundOutForJobInformation(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error)
	SoundOutForSendJobListing(param entity.SoundOutForJobInformationForSendMessageParam) (presenter.Presenter, error)
	SoundOutForMaskResume(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error)
	SoundOutForRequestShareJobSeeker(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error)
	SoundOutForRequestShareJobInformation(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error)

	// タスクを始める関数（打診）
	SoundOutGroupForJobInformation(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error)
	SoundOutGroupForSendJobListing(param entity.SoundOutForJobInformationForSendMessageParam) (presenter.Presenter, error)

	GetLatestTaskListByAgentStaffID(agentID, agentStaffID, deadLine, staffType, partnerType, taskType, phase, jobSeekerID uint) (presenter.Presenter, error)
	GetLatestSameTaskListByJobSeekerID(jobSeekerID, taskID uint, phaseCategory, phaseSubCategory null.Int) (presenter.Presenter, error)
	GetLatestTaskByJobSeekerIDAndJobInformationID(jobSeekerID, jobInformationID uint) (presenter.Presenter, error)
	CreateNextTaskAfterEntryPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error)
	CreateNextTaskAfterDocumentSelectionPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error) // 書類選考
	CreateNextTaskAfterSelectionPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error)         // 選考(1次-最終)
	CreateNextTaskAfterDeclinePhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error)           //  内定辞退
	CreateNextTaskAfterHoldJobOfferPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error)      // 内定保留
	CreateNextTaskAfterAcceptJobOfferPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error)    // 内定承諾
	CreateNextSameTaskList(param entity.NextSameTaskListParam, agentStaffID uint) (presenter.Presenter, error)            // 同一タスクをまとめて処理
	CreateEntryTaskFromMatchingJob(param entity.CreateEntryTaskFromMatchingJobParam) (presenter.Presenter, error)         // マイページのマッチ求人からエントリー
	UpdateTaskGroupDocument(param entity.UpdateTaskGroupDocumentParam) (presenter.Presenter, error)
	UpdateRALastWatched(groupID uint) (presenter.Presenter, error)
	UpdateRALastRequest(groupID uint) (presenter.Presenter, error)
	UpdateCALastWatched(groupID uint) (presenter.Presenter, error)
	UpdateCALastRequest(groupID uint) (presenter.Presenter, error)
	UpdateLastRequest(groupID uint) (presenter.Presenter, error)
	UpdateExternalJob(groupID uint, param entity.ExternalJob) (presenter.Presenter, error)

	// タスクを順番に複数処理するapi
	CreateTaskInBatchProcessing(token string, param entity.CreateTaskInBatchProcessingParam) (presenter.Presenter, error)

	// 求人打診可能なグループと不可能なグループを取得
	GetSoundOutGroupList(agentID uint, jobSeekerID uint, jobInformationIDList []uint) (presenter.Presenter, error)

	// 汎用系 API（仕様変更前に作成した関数）
	GetTaskByID(taskID uint) (presenter.Presenter, error)
	GetTaskListByAgentIDAndPage(agentID, pageNumber uint) (presenter.Presenter, error)
	GetTaskListAfterEntryByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error)
	GetSearchTaskListByAgentIDAndPage(agentStaffID, pageNumber uint, param entity.SearchTask) (presenter.Presenter, error)
	GetTaskGroupByID(taskGroupID uint) (presenter.Presenter, error)
	GetActiveTaskCountByBillingAddressID(billingAddressID uint) (presenter.Presenter, error)
	GetActiveTaskCountByJobInformationID(jobInformationID uint) (presenter.Presenter, error)
	GetActiveTaskCountBySelectionID(selectionID uint) (presenter.Presenter, error)

	DeleteTask(param entity.DeleteTaskParam) (presenter.Presenter, error) // GoogleCalendarからタスクを取得する関数
	// Admin API

	// Batch API
	BatchNotifyUnwatched(now time.Time) (presenter.Presenter, error)

	// トップページのタスク一覧を取得するapi
	GetJobSeekerTaskListByAgentStaffID(agentStaffID uint) (presenter.Presenter, error)
}

type TaskHandlerImpl struct {
	taskInteractor interactor.TaskInteractor
}

func NewTaskHandlerImpl(tI interactor.TaskInteractor) TaskHandler {
	return &TaskHandlerImpl{
		taskInteractor: tI,
	}
}

/****************************************************************************************/
// タスク開始の関数（打診）
// Memo: 旧版の為、新しいapiの動作確認後に削除（が対応）
//
func (h *TaskHandlerImpl) SoundOutForJobInformation(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.SoundOutForJobInformation(interactor.SoundOutForJobInformationInput{
		SoundOutParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

func (h *TaskHandlerImpl) SoundOutForSendJobListing(param entity.SoundOutForJobInformationForSendMessageParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.SoundOutForSendJobListing(interactor.SoundOutForSendJobListingInput{
		SoundOutParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

func (h *TaskHandlerImpl) SoundOutForMaskResume(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.SoundOutForMaskResume(interactor.SoundOutForMaskResumeInput{
		SoundOutParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

// シェア依頼処理（求人検索ページから）
func (h *TaskHandlerImpl) SoundOutForRequestShareJobSeeker(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.SoundOutForRequestShareJobSeeker(interactor.SoundOutForRequestShareJobSeekerInput{
		SoundOutParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

// シェア依頼処理（求職者検索ページから）
func (h *TaskHandlerImpl) SoundOutForRequestShareJobInformation(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.SoundOutForRequestShareJobInformation(interactor.SoundOutForRequestShareJobInformationInput{
		SoundOutParam: param,
	})

	if err != nil {
		return nil, err
	}

	// return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

/****************************************************************************************/
// タスク開始の関数（打診）
//
func (h *TaskHandlerImpl) SoundOutGroupForJobInformation(param entity.SoundOutForJobInformationParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.SoundOutGroupForJobInformation(interactor.SoundOutGroupForJobInformationInput{
		SoundOutParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

func (h *TaskHandlerImpl) SoundOutGroupForSendJobListing(param entity.SoundOutForJobInformationForSendMessageParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.SoundOutGroupForSendJobListing(interactor.SoundOutGroupForSendJobListingInput{
		SoundOutParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

/****************************************************************************************/

func (h *TaskHandlerImpl) CreateTaskInBatchProcessing(token string, param entity.CreateTaskInBatchProcessingParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateTaskInBatchProcessing(interactor.CreateTaskInBatchProcessingInput{
		Token:       token,
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 求人打診可能なグループと不可能なグループを取得
func (h *TaskHandlerImpl) GetSoundOutGroupList(agentID uint, jobSeekerID uint, jobInformationIDList []uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetSoundOutGroupList(interactor.GetSoundOutGroupListInput{
		AgentID:              agentID,
		JobSeekerID:          jobSeekerID,
		JobInformationIDList: jobInformationIDList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSoundOutGroupListJSONPresenter(responses.NewSoundOutGroupList(output.SoundOutGroupList, output.AlreadyCreatedGroupList)), nil
}

func (h *TaskHandlerImpl) GetLatestTaskListByAgentStaffID(agentID, agentStaffID, deadLine, staffType, partnerType, taskType, phase, jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetLatestTaskListByAgentStaffID(interactor.GetLatestTaskListByAgentStaffIDInput{
		AgentID:      agentID,
		AgentStaffID: agentStaffID,
		DeadLine:     deadLine,
		StaffType:    staffType,
		PartnerType:  partnerType,
		TaskType:     taskType,
		Phase:        phase,
		JobSeekerID:  jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListAndJobSeekerListJSONPresenter(responses.NewTaskListAndJobSeekerList(output.TaskList, output.JobSeekerList)), nil
}

// 「同一求職者」で「同一フェーズ」のタスクのリストを取得
func (h *TaskHandlerImpl) GetLatestSameTaskListByJobSeekerID(jobSeekerID, taskID uint, phaseCategory, phaseSubCategory null.Int) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetLatestSameTaskListByJobSeekerID(interactor.GetLatestSameTaskListByJobSeekerIDInput{
		JobSeekerID: jobSeekerID,
		TaskID:      taskID,
		Phase:       phaseCategory,
		PhaseSub:    phaseSubCategory,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

func (h *TaskHandlerImpl) GetLatestTaskByJobSeekerIDAndJobInformationID(jobSeekerID, jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetLatestTaskByJobSeekerIDAndJobInformationID(interactor.GetLatestTaskByJobSeekerIDAndJobInformationIDInput{
		JobSeekerID:      jobSeekerID,
		JobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskJSONPresenter(responses.NewTask(output.Task)), nil
}

func (h *TaskHandlerImpl) CreateNextTaskAfterEntryPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateNextTaskAfterEntryPhase(interactor.CreateNextTaskAfterEntryPhaseInput{
		NextTaskParam: param,
		AgentStaffID:  agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 書類選考
func (h *TaskHandlerImpl) CreateNextTaskAfterDocumentSelectionPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateNextTaskAfterDocumentSelectionPhase(interactor.CreateNextTaskAfterDocumentSelectionPhaseInput{
		NextTaskParam: param,
		AgentStaffID:  agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 選考(1次-最終)
func (h *TaskHandlerImpl) CreateNextTaskAfterSelectionPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateNextTaskAfterSelectionPhase(interactor.CreateNextTaskAfterSelectionPhaseInput{
		NextTaskParam: param,
		AgentStaffID:  agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 辞退
func (h *TaskHandlerImpl) CreateNextTaskAfterDeclinePhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateNextTaskAfterDeclinePhase(interactor.CreateNextTaskAfterDeclinePhaseInput{
		NextTaskParam: param,
		AgentStaffID:  agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 内定保留
func (h *TaskHandlerImpl) CreateNextTaskAfterHoldJobOfferPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateNextTaskAfterHoldJobOfferPhase(interactor.CreateNextTaskAfterHoldJobOfferPhaseInput{
		NextTaskParam: param,
		AgentStaffID:  agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 内定承諾
func (h *TaskHandlerImpl) CreateNextTaskAfterAcceptJobOfferPhase(param entity.NextTaskParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateNextTaskAfterAcceptJobOfferPhase(interactor.CreateNextTaskAfterAcceptJobOfferPhaseInput{
		NextTaskParam: param,
		AgentStaffID:  agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 同一タスクをまとめて処理
func (h *TaskHandlerImpl) CreateNextSameTaskList(param entity.NextSameTaskListParam, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateNextSameTaskList(interactor.CreateNextSameTaskListInput{
		NextSameTaskListParam: param,
		AgentStaffID:          agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// マイページのマッチ求人からエントリー
func (h *TaskHandlerImpl) CreateEntryTaskFromMatchingJob(param entity.CreateEntryTaskFromMatchingJobParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.CreateEntryTaskFromMatchingJob(interactor.CreateEntryTaskFromMatchingJobInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *TaskHandlerImpl) UpdateTaskGroupDocument(param entity.UpdateTaskGroupDocumentParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.UpdateTaskGroupDocument(interactor.UpdateTaskGroupDocumentInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *TaskHandlerImpl) UpdateRALastWatched(groupID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.UpdateRALastWatched(interactor.UpdateRALastWatchedInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *TaskHandlerImpl) UpdateRALastRequest(groupID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.UpdateRALastRequest(interactor.UpdateRALastRequestInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *TaskHandlerImpl) UpdateCALastWatched(groupID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.UpdateCALastWatched(interactor.UpdateCALastWatchedInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *TaskHandlerImpl) UpdateCALastRequest(groupID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.UpdateCALastRequest(interactor.UpdateCALastRequestInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 自分から自分にタスクを作成する場合に使用
// RAとCAの最終依頼時間を更新
func (h *TaskHandlerImpl) UpdateLastRequest(groupID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.UpdateLastRequest(interactor.UpdateLastRequestInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 外部求人情報を更新する
func (h *TaskHandlerImpl) UpdateExternalJob(groupID uint, param entity.ExternalJob) (presenter.Presenter, error) {
	output, err := h.taskInteractor.UpdateExternalJob(interactor.UpdateExternalJobInput{
		GroupID: groupID,
		Param:   param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
/// 汎用系 API（仕様変更前に作成した関数）
//

// タスクIDからタスクを取得する
func (h *TaskHandlerImpl) GetTaskByID(taskID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetTaskByID(interactor.GetTaskByIDInput{
		TaskID: taskID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskJSONPresenter(responses.NewTask(output.Task)), nil
}

// タスクグループの一覧取得（エージェントが関わっているタスク）
func (h *TaskHandlerImpl) GetTaskListByAgentIDAndPage(agentID, pageNumber uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetTaskListByAgentIDAndPage(interactor.GetTaskListByAgentIDAndPageInput{
		AgentID:    agentID,
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskPageListAndMaxPageJSONPresenter(responses.NewTaskListAndMaxPage(output.TaskList, output.MaxPageNumber)), nil
}

// タスクグループの一覧取得（エージェントが関わっているタスク）
func (h *TaskHandlerImpl) GetTaskListAfterEntryByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetTaskListAfterEntryByJobSeekerID(interactor.GetTaskListAfterEntryByJobSeekerIDInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskListJSONPresenter(responses.NewTaskList(output.TaskList)), nil
}

func (h *TaskHandlerImpl) GetSearchTaskListByAgentIDAndPage(agentID, pageNumber uint, searchParam entity.SearchTask) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetSearchTaskListByAgentIDAndPage(interactor.GetSearchTaskListByAgentIDAndPageInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskPageListAndMaxPageJSONPresenter(responses.NewTaskListAndMaxPage(output.TaskList, output.MaxPageNumber)), nil
}

// タスクグループの一覧取得（自分が関わっているタスク）
func (h *TaskHandlerImpl) GetTaskGroupByID(taskGroupID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetTaskGroupByID(interactor.GetTaskGroupByIDInput{
		TaskGroupID: taskGroupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewTaskGroupJSONPresenter(responses.NewTaskGroup(output.TaskGroup)), nil
}

// アクティブなタスクの数を取得（請求先のID）
func (h *TaskHandlerImpl) GetActiveTaskCountByBillingAddressID(billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetActiveTaskCountByBillingAddressID(interactor.GetActiveTaskCountByBillingAddressIDInput{
		BillingAddressID: billingAddressID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewActiveTaskCountJSONPresenter(responses.NewActiveTaskCount(output.TaskCount)), nil
}

// アクティブなタスクの数を取得（求人のID）
func (h *TaskHandlerImpl) GetActiveTaskCountByJobInformationID(jobInformationID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetActiveTaskCountByJobInformationID(interactor.GetActiveTaskCountByJobInformationIDInput{
		JobInformationID: jobInformationID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewActiveTaskCountJSONPresenter(responses.NewActiveTaskCount(output.TaskCount)), nil
}

// アクティブなタスクの数を取得（選考フローのID）
func (h *TaskHandlerImpl) GetActiveTaskCountBySelectionID(selectionID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetActiveTaskCountBySelectionID(interactor.GetActiveTaskCountBySelectionIDInput{
		SelectionID: selectionID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewActiveTaskCountJSONPresenter(responses.NewActiveTaskCount(output.TaskCount)), nil
}

func (h *TaskHandlerImpl) DeleteTask(param entity.DeleteTaskParam) (presenter.Presenter, error) {
	output, err := h.taskInteractor.DeleteTask(interactor.DeleteTaskInput{
		DeleteParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
/// Admin API

/****************************************************************************************/
/****************************************************************************************/
// Batch API
//
// 未読or未処理を通知(エージェントとのチャット, タスク)
func (h *TaskHandlerImpl) BatchNotifyUnwatched(now time.Time) (presenter.Presenter, error) {
	output, err := h.taskInteractor.BatchNotifyUnwatched(interactor.BatchNotifyUnwatchedInput{
		Now: now,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/

// トップページのタスク一覧を取得するapi
func (h *TaskHandlerImpl) GetJobSeekerTaskListByAgentStaffID(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.taskInteractor.GetJobSeekerTaskListByAgentStaffID(interactor.GetJobSeekerTaskListByAgentStaffIDInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerTaskListJSONPresenter(responses.NewJobSeekerTaskList(output.TaskList)), nil
}
