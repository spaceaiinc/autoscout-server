package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationTarget struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	Target                  null.Int  `db:"target" json:"target"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationTarget(
	sendingJobInformationID uint,
	target null.Int,
) *SendingJobInformationTarget {
	return &SendingJobInformationTarget{
		SendingJobInformationID: sendingJobInformationID,
		Target:                  target,
	}
}
