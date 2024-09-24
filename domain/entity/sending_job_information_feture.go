package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationFeature struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	Feature                 null.Int  `db:"feature" json:"feature"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationFeature(
	sendingJobInformationID uint,
	feature null.Int,
) *SendingJobInformationFeature {
	return &SendingJobInformationFeature{
		SendingJobInformationID: sendingJobInformationID,
		Feature:                 feature,
	}
}
