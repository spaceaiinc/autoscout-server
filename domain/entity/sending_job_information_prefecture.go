package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationPrefecture struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	Prefecture              null.Int  `db:"prefecture" json:"prefecture"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationPrefecture(
	sendingJobInformationID uint,
	prefecture null.Int,
) *SendingJobInformationPrefecture {
	return &SendingJobInformationPrefecture{
		SendingJobInformationID: sendingJobInformationID,
		Prefecture:              prefecture,
	}
}
