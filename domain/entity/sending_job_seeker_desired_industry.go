package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerDesiredIndustry struct {
	ID                 uint      `db:"id" json:"id"`
	SendingJobSeekerID uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	DesiredIndustry    null.Int  `db:"desired_industry" json:"desired_industry"`
	DesiredRank        null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerDesiredIndustry(
	sendingJobSeekerID uint,
	desiredIndustry null.Int,
	desiredRank null.Int,
) *SendingJobSeekerDesiredIndustry {
	return &SendingJobSeekerDesiredIndustry{
		SendingJobSeekerID: sendingJobSeekerID,
		DesiredIndustry:    desiredIndustry,
		DesiredRank:        desiredRank,
	}
}
