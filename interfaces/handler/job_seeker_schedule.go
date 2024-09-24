package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type JobSeekerScheduleHandler interface {
	// 汎用系 API
	CreateJobSeekerSchedule(param entity.CreateOrUpdateJobSeekerScheduleParam) (presenter.Presenter, error)
	UpdateJobSeekerSchedule(jobSeekerScheduleID uint, param entity.CreateOrUpdateJobSeekerScheduleParam) (presenter.Presenter, error)
	ShareScheduleToCAStaff(jobSeekerID uint, param entity.ShareScheduleParam) (presenter.Presenter, error)
	GetJobSeekerScheduleListByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error)
	GetJobSeekerScheduleTypeListByJobSeekerID(jobSeekerID, scheduleType uint) (presenter.Presenter, error)
	DeleteJobSeekerSchedule(jobSeekerScheduleID uint) (presenter.Presenter, error)
}

type JobSeekerScheduleHandlerImpl struct {
	jobSeekerScheduleInteractor interactor.JobSeekerScheduleInteractor
}

func NewJobSeekerScheduleHandlerImpl(dI interactor.JobSeekerScheduleInteractor) JobSeekerScheduleHandler {
	return &JobSeekerScheduleHandlerImpl{
		jobSeekerScheduleInteractor: dI,
	}
}

func (h *JobSeekerScheduleHandlerImpl) CreateJobSeekerSchedule(param entity.CreateOrUpdateJobSeekerScheduleParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerScheduleInteractor.CreateJobSeekerSchedule(interactor.CreateJobSeekerScheduleInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerScheduleJSONPresenter(responses.NewJobSeekerSchedule(output.JobSeekerSchedule)), nil
}

func (h *JobSeekerScheduleHandlerImpl) UpdateJobSeekerSchedule(jobSeekerScheduleID uint, param entity.CreateOrUpdateJobSeekerScheduleParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerScheduleInteractor.UpdateJobSeekerSchedule(interactor.UpdateJobSeekerScheduleInput{
		JobSeekerScheduleID: jobSeekerScheduleID,
		UpdateParam:         param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerScheduleJSONPresenter(responses.NewJobSeekerSchedule(output.JobSeekerSchedule)), nil
}

// CA担当者に新規追加した日程を追加
func (h *JobSeekerScheduleHandlerImpl) ShareScheduleToCAStaff(jobSeekerID uint, param entity.ShareScheduleParam) (presenter.Presenter, error) {
	output, err := h.jobSeekerScheduleInteractor.ShareScheduleToCAStaff(interactor.ShareScheduleToCAStaffInput{
		JobSeekerID: jobSeekerID,
		Param:       param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerScheduleListJSONPresenter(responses.NewJobSeekerScheduleList(output.JobSeekerScheduleList)), nil
}

// 求職者IDでスケジュール情報を取得
func (h *JobSeekerScheduleHandlerImpl) GetJobSeekerScheduleListByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerScheduleInteractor.GetJobSeekerScheduleListByJobSeekerID(interactor.GetJobSeekerScheduleListByJobSeekerIDInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerScheduleListJSONPresenter(responses.NewJobSeekerScheduleList(output.JobSeekerScheduleList)), nil
}

// 求職者IDとスケジュールタイプでスケジュール情報を取得
func (h *JobSeekerScheduleHandlerImpl) GetJobSeekerScheduleTypeListByJobSeekerID(jobSeekerID, scheduleType uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerScheduleInteractor.GetJobSeekerScheduleTypeListByJobSeekerID(interactor.GetJobSeekerScheduleTypeListByJobSeekerIDInput{
		JobSeekerID:  jobSeekerID,
		ScheduleType: scheduleType,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerScheduleListJSONPresenter(responses.NewJobSeekerScheduleList(output.JobSeekerScheduleList)), nil
}

// スケジュールを削除
func (h *JobSeekerScheduleHandlerImpl) DeleteJobSeekerSchedule(jobSeekerScheduleID uint) (presenter.Presenter, error) {
	output, err := h.jobSeekerScheduleInteractor.DeleteJobSeekerSchedule(interactor.DeleteJobSeekerScheduleInput{
		JobSeekerScheduleID: jobSeekerScheduleID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
