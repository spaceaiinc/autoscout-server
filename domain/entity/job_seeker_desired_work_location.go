package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerDesiredWorkLocation struct {
	ID                  uint      `db:"id" json:"id"`
	JobSeekerID         uint      `db:"job_seeker_id" json:"job_seeker_id"`
	DesiredWorkLocation null.Int  `db:"desired_work_location" json:"desired_work_location"`
	DesiredRank         null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerDesiredWorkLocation(
	jobSeekerID uint,
	desiredWorkLocation null.Int,
	desiredRank null.Int,
) *JobSeekerDesiredWorkLocation {
	return &JobSeekerDesiredWorkLocation{
		JobSeekerID:         jobSeekerID,
		DesiredWorkLocation: desiredWorkLocation,
		DesiredRank:         desiredRank,
	}
}
