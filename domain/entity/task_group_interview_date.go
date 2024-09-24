package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type TaskGroupInterviewDate struct {
	ID            uint      `db:"id" json:"id"`
	TaskGroupID   uint      `db:"task_group_id" json:"task_group_id"`
	PhaseCategory null.Int  `db:"phase_category" json:"phase_category"`
	InterviewDate string    `db:"interview_date" json:"interview_date"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

	JobSeekerID uint   `db:"job_seeker_id" json:"job_seeker_id"`
	CompanyName string `db:"company_name" json:"company_name"`
}

func NewTaskGroupInterviewDate(
	taskGroupID uint,
	phaseCategory null.Int,
	interviewDate string,
) *TaskGroupInterviewDate {
	return &TaskGroupInterviewDate{
		TaskGroupID:   taskGroupID,
		PhaseCategory: phaseCategory,
		InterviewDate: interviewDate,
	}
}
