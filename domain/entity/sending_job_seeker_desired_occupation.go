package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerDesiredOccupation struct {
	ID                 uint      `db:"id" json:"id"`
	SendingJobSeekerID uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	DesiredOccupation  null.Int  `db:"desired_occupation" json:"desired_occupation"`
	DesiredRank        null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerDesiredOccupation(
	sendingJobSeekerID uint,
	desiredOccupation null.Int,
	desiredRank null.Int,
) *SendingJobSeekerDesiredOccupation {
	return &SendingJobSeekerDesiredOccupation{
		SendingJobSeekerID: sendingJobSeekerID,
		DesiredOccupation:  desiredOccupation,
		DesiredRank:        desiredRank,
	}
}
