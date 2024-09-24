package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type Schedule struct {
	Schedule *entity.Schedule `json:"schedule"`
}

func NewSchedule(schedule *entity.Schedule) Schedule {
	return Schedule{
		Schedule: schedule,
	}
}

type ScheduleList struct {
	ScheduleList []*entity.Schedule `json:"schedule_list"`
}

func NewScheduleList(schedules []*entity.Schedule) ScheduleList {
	return ScheduleList{
		ScheduleList: schedules,
	}
}
