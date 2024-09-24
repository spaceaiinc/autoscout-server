package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type InterviewTaskHandler interface {
	// 汎用系 API
	CreateNextInterviewTask(param entity.NextInterviewTaskParam) (presenter.Presenter, error)
	GetLatestAdjustmentTaskListByAgentID(agentID uint) (presenter.Presenter, error)   // 面談調整タスク（phase_category = (0 or 1)）の取得
	GetLatestConfirmationTaskListByAgentID(agentID uint) (presenter.Presenter, error) // 参加確認（phase_category = (2 or 3)）タスクの取得
	GetInterviewTaskListByGroupID(groupID uint) (presenter.Presenter, error)          // GroupIDから面談調整タスク一覧を取得 *アクティビティ表示に使用

	UpdateInterviewTaskGroupLastWatched(groupID uint) (presenter.Presenter, error) // 最終閲覧時間を更新
	UpdateInterviewTaskCAStaffID(param entity.UpdateCAStaffIDParam) (presenter.Presenter, error)
	UpdateInterviewTaskInterviewDate(param entity.UpdateInterviewDateParam) (presenter.Presenter, error)

	DeleteLatestInterviewTask(param entity.DeleteLatestInterviewTaskParam) (presenter.Presenter, error)
}

type InterviewTaskHandlerImpl struct {
	interviewTaskInteractor interactor.InterviewTaskInteractor
}

func NewInterviewTaskHandlerImpl(itI interactor.InterviewTaskInteractor) InterviewTaskHandler {
	return &InterviewTaskHandlerImpl{
		interviewTaskInteractor: itI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *InterviewTaskHandlerImpl) CreateNextInterviewTask(param entity.NextInterviewTaskParam) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.CreateNextInterviewTask(interactor.CreateNextInterviewTaskInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *InterviewTaskHandlerImpl) GetLatestAdjustmentTaskListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.GetLatestAdjustmentTaskListByAgentID(interactor.GetLatestAdjustmentTaskListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewTaskListJSONPresenter(responses.NewInterviewTaskList(output.InterviewTaskList)), nil
}

func (h *InterviewTaskHandlerImpl) GetLatestConfirmationTaskListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.GetLatestConfirmationTaskListByAgentID(interactor.GetLatestConfirmationTaskListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewTaskListJSONPresenter(responses.NewInterviewTaskList(output.InterviewTaskList)), nil
}

// GroupIDから面談調整タスク一覧を取得 *アクティビティ表示に使用
func (h *InterviewTaskHandlerImpl) GetInterviewTaskListByGroupID(groupID uint) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.GetInterviewTaskListByGroupID(interactor.GetInterviewTaskListByGroupIDInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewTaskListJSONPresenter(responses.NewInterviewTaskList(output.InterviewTaskList)), nil
}

func (h *InterviewTaskHandlerImpl) UpdateInterviewTaskGroupLastWatched(groupID uint) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.UpdateInterviewTaskGroupLastWatched(interactor.UpdateInterviewTaskGroupLastWatchedInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *InterviewTaskHandlerImpl) UpdateInterviewTaskCAStaffID(param entity.UpdateCAStaffIDParam) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.UpdateInterviewTaskCAStaffID(interactor.UpdateInterviewTaskCAStaffIDInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewTaskJSONPresenter(responses.NewInterviewTask(output.InterviewTask)), nil
}

func (h *InterviewTaskHandlerImpl) UpdateInterviewTaskInterviewDate(param entity.UpdateInterviewDateParam) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.UpdateInterviewTaskInterviewDate(interactor.UpdateInterviewTaskInterviewDateInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewInterviewTaskJSONPresenter(responses.NewInterviewTask(output.InterviewTask)), nil
}

// 面談調整タスクを削除
func (h *InterviewTaskHandlerImpl) DeleteLatestInterviewTask(param entity.DeleteLatestInterviewTaskParam) (presenter.Presenter, error) {
	output, err := h.interviewTaskInteractor.DeleteLatestInterviewTask(interactor.DeleteLatestInterviewTaskInput{
		Param: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
