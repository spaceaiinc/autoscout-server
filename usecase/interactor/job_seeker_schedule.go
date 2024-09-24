package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerScheduleInteractor interface {
	// 汎用系 API
	CreateJobSeekerSchedule(input CreateJobSeekerScheduleInput) (CreateJobSeekerScheduleOutput, error)
	UpdateJobSeekerSchedule(input UpdateJobSeekerScheduleInput) (UpdateJobSeekerScheduleOutput, error)
	ShareScheduleToCAStaff(input ShareScheduleToCAStaffInput) (ShareScheduleToCAStaffOutput, error)
	GetJobSeekerScheduleListByJobSeekerID(input GetJobSeekerScheduleListByJobSeekerIDInput) (GetJobSeekerScheduleListByJobSeekerIDOutput, error)
	GetJobSeekerScheduleTypeListByJobSeekerID(input GetJobSeekerScheduleTypeListByJobSeekerIDInput) (GetJobSeekerScheduleTypeListByJobSeekerIDOutput, error)
	DeleteJobSeekerSchedule(input DeleteJobSeekerScheduleInput) (DeleteJobSeekerScheduleOutput, error)
}

type JobSeekerScheduleInteractorImpl struct {
	firebase                    usecase.Firebase
	sendgrid                    config.Sendgrid
	jobSeekerScheduleRepository usecase.JobSeekerScheduleRepository
	jobSeekerRepository         usecase.JobSeekerRepository
	agentStaffRepository        usecase.AgentStaffRepository
}

func NewJobSeekerScheduleInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	jssR usecase.JobSeekerScheduleRepository,
	jsR usecase.JobSeekerRepository,
	asR usecase.AgentStaffRepository,
) JobSeekerScheduleInteractor {
	return &JobSeekerScheduleInteractorImpl{
		firebase:                    fb,
		sendgrid:                    sg,
		jobSeekerScheduleRepository: jssR,
		jobSeekerRepository:         jsR,
		agentStaffRepository:        asR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//求職者の作成
type CreateJobSeekerScheduleInput struct {
	CreateParam entity.CreateOrUpdateJobSeekerScheduleParam
}

type CreateJobSeekerScheduleOutput struct {
	JobSeekerSchedule *entity.JobSeekerSchedule
}

func (i *JobSeekerScheduleInteractorImpl) CreateJobSeekerSchedule(input CreateJobSeekerScheduleInput) (CreateJobSeekerScheduleOutput, error) {
	var (
		output CreateJobSeekerScheduleOutput
		err    error
	)

	jobSeekerSchedule := entity.NewJobSeekerSchedule(
		input.CreateParam.JobSeekerID,
		input.CreateParam.TaskID,
		input.CreateParam.ScheduleType,
		input.CreateParam.Title,
		input.CreateParam.StartTime,
		input.CreateParam.EndTime,
		input.CreateParam.SeekerDescription,
		input.CreateParam.StaffDescription,
		input.CreateParam.IsShare,
		input.CreateParam.RepetitionCount,
	)

	err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerSchedule = jobSeekerSchedule

	return output, nil
}

// 求職者のスケジュール情報を更新
type UpdateJobSeekerScheduleInput struct {
	JobSeekerScheduleID uint
	UpdateParam         entity.CreateOrUpdateJobSeekerScheduleParam
}

type UpdateJobSeekerScheduleOutput struct {
	JobSeekerSchedule *entity.JobSeekerSchedule
}

func (i *JobSeekerScheduleInteractorImpl) UpdateJobSeekerSchedule(input UpdateJobSeekerScheduleInput) (UpdateJobSeekerScheduleOutput, error) {
	var (
		output UpdateJobSeekerScheduleOutput
		err    error
	)

	jobSeekerSchedule := entity.NewJobSeekerSchedule(
		input.UpdateParam.JobSeekerID,
		input.UpdateParam.TaskID,
		input.UpdateParam.ScheduleType,
		input.UpdateParam.Title,
		input.UpdateParam.StartTime,
		input.UpdateParam.EndTime,
		input.UpdateParam.SeekerDescription,
		input.UpdateParam.StaffDescription,
		input.UpdateParam.IsShare,
		input.UpdateParam.RepetitionCount,
	)

	err = i.jobSeekerScheduleRepository.Update(input.JobSeekerScheduleID, jobSeekerSchedule)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobSeekerSchedule.ID = input.JobSeekerScheduleID
	output.JobSeekerSchedule = jobSeekerSchedule

	return output, nil
}

// 担当CAに
type ShareScheduleToCAStaffInput struct {
	JobSeekerID uint
	Param       entity.ShareScheduleParam
}

type ShareScheduleToCAStaffOutput struct {
	JobSeekerScheduleList []*entity.JobSeekerSchedule
}

func (i *JobSeekerScheduleInteractorImpl) ShareScheduleToCAStaff(input ShareScheduleToCAStaffInput) (ShareScheduleToCAStaffOutput, error) {
	var (
		output     ShareScheduleToCAStaffOutput
		err        error
		idUintList []uint

		scheduleList []*entity.JobSeekerSchedule
		scheduleStr  string
	)

	for _, schedule := range input.Param.SharedList {
		idUintList = append(idUintList, schedule.ID)
	}

	if len(idUintList) > 0 {
		err = i.jobSeekerScheduleRepository.UpdateIsSharedByIDList(idUintList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, schedule := range input.Param.SharedList {
			// IsShareをtrueにする
			jobSeekerSchedule := entity.NewJobSeekerSchedule(
				schedule.JobSeekerID,
				schedule.TaskID,
				schedule.ScheduleType,
				schedule.Title,
				schedule.StartTime,
				schedule.EndTime,
				schedule.SeekerDescription,
				schedule.StaffDescription,
				true,
				schedule.RepetitionCount,
			)

			jobSeekerSchedule.ID = schedule.ID

			scheduleList = append(scheduleList, jobSeekerSchedule)

			timePastStr := fmt.Sprintf(
				"%s 〜 %s\n",
				timeStrFormat(schedule.StartTime), timeStrFormat(schedule.EndTime),
			)

			scheduleStr = scheduleStr + timePastStr
		}
	}

	jobSeeker, err := i.jobSeekerRepository.FindByID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentStaff, err := i.agentStaffRepository.FindByID(uint(jobSeeker.AgentStaffID.Int64))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 通知がオフの場合はメール送信しない
	if agentStaff.NotificationJobSeeker {
		output.JobSeekerScheduleList = scheduleList
		return output, nil
	}

	messageText := fmt.Sprintf(
		"%s様がスケジュールを共有しました。\n求職者詳細ページからご確認ください。\n\n■共有されたスケジュール\n%s",
		jobSeeker.LastName+jobSeeker.FirstName, scheduleStr,
	)

	// メール送信
	err = utility.SendMailToSingle(
		i.sendgrid.APIKey,
		"求職者が希望日時を共有しました",
		messageText,
		entity.EmailUser{
			Name:  "autoscout事務局",
			Email: "info@spaceai.jp",
		},
		entity.EmailUser{
			Name:  agentStaff.StaffName,
			Email: agentStaff.Email,
		},
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerScheduleList = scheduleList

	return output, nil
}

// 求職者IDでスケジュール情報を取得する
type GetJobSeekerScheduleListByJobSeekerIDInput struct {
	JobSeekerID uint
}

type GetJobSeekerScheduleListByJobSeekerIDOutput struct {
	JobSeekerScheduleList []*entity.JobSeekerSchedule
}

func (i *JobSeekerScheduleInteractorImpl) GetJobSeekerScheduleListByJobSeekerID(input GetJobSeekerScheduleListByJobSeekerIDInput) (GetJobSeekerScheduleListByJobSeekerIDOutput, error) {
	var (
		output GetJobSeekerScheduleListByJobSeekerIDOutput
		err    error
	)

	jobSeekerScheduleList, err := i.jobSeekerScheduleRepository.GetByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerScheduleList = jobSeekerScheduleList

	return output, nil
}

// 求職者IDでスケジュール情報を取得する
type GetJobSeekerScheduleTypeListByJobSeekerIDInput struct {
	JobSeekerID  uint
	ScheduleType uint
}

type GetJobSeekerScheduleTypeListByJobSeekerIDOutput struct {
	JobSeekerScheduleList []*entity.JobSeekerSchedule
}

func (i *JobSeekerScheduleInteractorImpl) GetJobSeekerScheduleTypeListByJobSeekerID(input GetJobSeekerScheduleTypeListByJobSeekerIDInput) (GetJobSeekerScheduleTypeListByJobSeekerIDOutput, error) {
	var (
		output GetJobSeekerScheduleTypeListByJobSeekerIDOutput
		err    error
	)

	jobSeekerScheduleList, err := i.jobSeekerScheduleRepository.GetByJobSeekerIDAndScheuldType(input.JobSeekerID, input.ScheduleType)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerScheduleList = jobSeekerScheduleList

	return output, nil
}

// スケジュールの削除
type DeleteJobSeekerScheduleInput struct {
	JobSeekerScheduleID uint
}

type DeleteJobSeekerScheduleOutput struct {
	OK bool
}

func (i *JobSeekerScheduleInteractorImpl) DeleteJobSeekerSchedule(input DeleteJobSeekerScheduleInput) (DeleteJobSeekerScheduleOutput, error) {
	var (
		output DeleteJobSeekerScheduleOutput
		err    error
	)

	err = i.jobSeekerScheduleRepository.Delete(input.JobSeekerScheduleID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}
