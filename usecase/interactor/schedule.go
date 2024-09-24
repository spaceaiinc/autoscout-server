package interactor

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ScheduleInteractor interface {
	// 汎用系
	GetScheduleListWithInPeriod(input GetScheduleListWithInPeriodInput) (GetScheduleListWithInPeriodOutput, error)
	GetScheduleListWithInPeriodByStaffIDList(input GetScheduleListWithInPeriodByStaffIDListInput) (GetScheduleListWithInPeriodByStaffIDListOutput, error)
}

type ScheduleInteractorImpl struct {
	firebase                     usecase.Firebase
	sendgrid                     config.Sendgrid
	oneSignal                    config.OneSignal
	agentRepository              usecase.AgentRepository
	agentStaffRepository         usecase.AgentStaffRepository
	interviewTaskGroupRepository usecase.InterviewTaskGroupRepository
	interviewTaskRepository      usecase.InterviewTaskRepository
	taskGroupRepository          usecase.TaskGroupRepository
	taskRepository               usecase.TaskRepository
	jobSeekerScheduleRepository  usecase.JobSeekerScheduleRepository
}

// ScheduleInteractorImpl is an implementation of ScheduleInteractor
func NewScheduleInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	itgR usecase.InterviewTaskGroupRepository,
	itR usecase.InterviewTaskRepository,
	tgR usecase.TaskGroupRepository,
	tR usecase.TaskRepository,
	jssR usecase.JobSeekerScheduleRepository,
) ScheduleInteractor {
	return &ScheduleInteractorImpl{
		firebase:                     fb,
		sendgrid:                     sg,
		oneSignal:                    os,
		agentRepository:              aR,
		agentStaffRepository:         asR,
		interviewTaskGroupRepository: itgR,
		interviewTaskRepository:      itR,
		taskGroupRepository:          tgR,
		taskRepository:               tR,
		jobSeekerScheduleRepository:  jssR,
	}
}

/****************************************************************************************/
/// 汎用系
//

type GetScheduleListWithInPeriodInput struct {
	AgentStaffID uint
	StartDate    string
	EndDate      string
}

type GetScheduleListWithInPeriodOutput struct {
	ScheduleList []*entity.Schedule
}

func (i *ScheduleInteractorImpl) GetScheduleListWithInPeriod(input GetScheduleListWithInPeriodInput) (GetScheduleListWithInPeriodOutput, error) {
	var (
		output       GetScheduleListWithInPeriodOutput
		err          error
		scheduleList []*entity.Schedule
	)

	layout := "2006-01-02T15:04:05.000Z" // 文字列のフォーマットを指定

	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		fmt.Println("変換エラー:", err)
		return output, err
	}

	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		fmt.Println("変換エラー:", err)
		return output, err
	}

	jobSeekerScheduleList, err := i.jobSeekerScheduleRepository.GetByStaffIDAndPeriod(input.AgentStaffID, uint(entity.ScheduleTypeByEnterprise), startDate.Add(time.Duration(9)*time.Hour), endDate.Add(time.Duration(9)*time.Hour))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	interviewTaskGroupList, err := i.interviewTaskGroupRepository.GetByStaffIDAndPeriod(input.AgentStaffID, startDate.Add(time.Duration(9)*time.Hour), endDate.Add(time.Duration(9)*time.Hour))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, itg := range interviewTaskGroupList {

		// 1時間後を取得します。
		addOneHourDateTime := itg.InterviewDate.Add(time.Duration(1) * time.Hour)

		schedule := entity.NewSchedule(
			itg.ID,
			input.AgentStaffID,
			itg.StaffName,
			itg.AgentName,
			"",
			"",
			fmt.Sprintf("面談/%s%s", itg.LastName, itg.FirstName),
			itg.InterviewDate,
			addOneHourDateTime,
			"",
			"",
			"interview",
		)
		scheduleList = append(scheduleList, schedule)
	}

	for _, js := range jobSeekerScheduleList {
		fmt.Println(
			js.CAStaffName,
			js.CAAgentName,
			js.RAStaffName,
			js.RAAgentName,
		)

		startTime, err := time.Parse("2006-01-02T15:04", js.StartTime)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		startTime = startTime.Add(time.Duration(-9) * time.Hour)

		endTime, err := time.Parse("2006-01-02T15:04", js.EndTime)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		endTime = endTime.Add(time.Duration(-9) * time.Hour)

		schedule := entity.NewSchedule(
			js.ID,
			input.AgentStaffID,
			js.CAStaffName,
			js.CAAgentName,
			js.RAStaffName,
			js.RAAgentName,
			fmt.Sprintf("%s / %s / %s", getStrTaskPhase(js.PhaseCategory), js.CompanyName, (js.LastName+js.FirstName)),
			startTime,
			endTime,
			js.SeekerDescription,
			js.StaffDescription,
			"selection",
		)

		scheduleList = append(scheduleList, schedule)
	}

	output.ScheduleList = scheduleList

	return output, nil
}

type GetScheduleListWithInPeriodByStaffIDListInput struct {
	StaffIDList []uint
	StartDate   string
	EndDate     string
}

type GetScheduleListWithInPeriodByStaffIDListOutput struct {
	ScheduleList []*entity.Schedule
}

func (i *ScheduleInteractorImpl) GetScheduleListWithInPeriodByStaffIDList(input GetScheduleListWithInPeriodByStaffIDListInput) (GetScheduleListWithInPeriodByStaffIDListOutput, error) {
	var (
		output       GetScheduleListWithInPeriodByStaffIDListOutput
		err          error
		scheduleList []*entity.Schedule
	)

	layout := "2006-01-02T15:04:05.000Z" // 文字列のフォーマットを指定

	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		fmt.Println("変換エラー:", err)
		return output, err
	}

	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		fmt.Println("変換エラー:", err)
		return output, err
	}

	jobSeekerScheduleList, err := i.jobSeekerScheduleRepository.GetByStaffIDAndPeriodByStaffIDList(input.StaffIDList, uint(entity.ScheduleTypeByEnterprise), startDate.Add(time.Duration(9)*time.Hour), endDate.Add(time.Duration(9)*time.Hour))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	interviewTaskGroupList, err := i.interviewTaskGroupRepository.GetByStaffIDAndPeriodByStaffIDList(input.StaffIDList, startDate.Add(time.Duration(9)*time.Hour), endDate.Add(time.Duration(9)*time.Hour))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, itg := range interviewTaskGroupList {
		// 1時間後を取得します。
		addOneHourDateTime := itg.InterviewDate.Add(time.Duration(1) * time.Hour)

		schedule := entity.NewSchedule(
			itg.ID,
			itg.AgentStaffID,
			itg.StaffName,
			itg.AgentName,
			"",
			"",
			fmt.Sprintf("面談/%s%s", itg.LastName, itg.FirstName),
			itg.InterviewDate,
			addOneHourDateTime,
			"",
			"",
			"interview",
		)
		scheduleList = append(scheduleList, schedule)
	}

	for _, js := range jobSeekerScheduleList {
		startTime, err := time.Parse("2006-01-02T15:04", js.StartTime)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		startTime = startTime.Add(time.Duration(-9) * time.Hour)

		endTime, err := time.Parse("2006-01-02T15:04", js.EndTime)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		endTime = endTime.Add(time.Duration(-9) * time.Hour)

		schedule := entity.NewSchedule(
			js.ID,
			js.CAStaffID,
			js.CAStaffName,
			js.CAAgentName,
			js.RAStaffName,
			js.RAAgentName,
			fmt.Sprintf("%s / %s / %s", getStrTaskPhase(js.PhaseCategory), js.CompanyName, (js.LastName+js.FirstName)),
			startTime,
			endTime,
			js.SeekerDescription,
			js.StaffDescription,
			"selection",
		)

		scheduleList = append(scheduleList, schedule)
	}

	output.ScheduleList = scheduleList

	return output, nil
}
