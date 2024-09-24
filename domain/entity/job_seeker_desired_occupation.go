package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerDesiredOccupation struct {
	ID                uint      `db:"id" json:"id"`
	JobSeekerID       uint      `db:"job_seeker_id" json:"job_seeker_id"`
	DesiredOccupation null.Int  `db:"desired_occupation" json:"desired_occupation"`
	DesiredRank       null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt         time.Time `db:"created_at" json:"-"`
	UpdatedAt         time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerDesiredOccupation(
	jobSeekerID uint,
	desiredOccupation null.Int,
	desiredRank null.Int,
) *JobSeekerDesiredOccupation {
	return &JobSeekerDesiredOccupation{
		JobSeekerID:       jobSeekerID,
		DesiredOccupation: desiredOccupation,
		DesiredRank:       desiredRank,
	}
}
