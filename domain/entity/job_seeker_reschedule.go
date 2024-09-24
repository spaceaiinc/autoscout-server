package entity

import "time"

type JobSeekerReschedule struct {
	ID           uint      `db:"id" json:"id"`
	RescheduleID uint      `db:"reschedule_id" json:"reschedule_id"`
	TaskID       uint      `db:"task_id" json:"task_id"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerReschedule(
	rescheduleID uint,
	taskID uint,
) *JobSeekerReschedule {
	return &JobSeekerReschedule{
		RescheduleID: rescheduleID,
		TaskID:       taskID,
	}
}
