package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationIndustry struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	Industry                null.Int  `db:"industry" json:"industry"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationIndustry(
	sendingJobInformationID uint,
	industry null.Int,
) *SendingJobInformationIndustry {
	return &SendingJobInformationIndustry{
		SendingJobInformationID: sendingJobInformationID,
		Industry:                industry,
	}
}
