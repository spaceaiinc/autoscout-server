package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerDesiredWorkLocation struct {
	ID                  uint      `db:"id" json:"id"`
	SendingJobSeekerID  uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	DesiredWorkLocation null.Int  `db:"desired_work_location" json:"desired_work_location"`
	DesiredRank         null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerDesiredWorkLocation(
	sendingJobSeekerID uint,
	desiredWorkLocation null.Int,
	desiredRank null.Int,
) *SendingJobSeekerDesiredWorkLocation {
	return &SendingJobSeekerDesiredWorkLocation{
		SendingJobSeekerID:  sendingJobSeekerID,
		DesiredWorkLocation: desiredWorkLocation,
		DesiredRank:         desiredRank,
	}
}
