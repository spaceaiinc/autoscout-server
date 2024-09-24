package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type JobSeekerSchedule struct {
	JobSeekerSchedule *entity.JobSeekerSchedule `json:"job_seeker_schedule"`
}

func NewJobSeekerSchedule(jobSeeker *entity.JobSeekerSchedule) JobSeekerSchedule {
	return JobSeekerSchedule{
		JobSeekerSchedule: jobSeeker,
	}
}

type JobSeekerScheduleList struct {
	JobSeekerScheduleList []*entity.JobSeekerSchedule `json:"job_seeker_schedule_list"`
}

func NewJobSeekerScheduleList(jobSeekers []*entity.JobSeekerSchedule) JobSeekerScheduleList {
	return JobSeekerScheduleList{
		JobSeekerScheduleList: jobSeekers,
	}
}
