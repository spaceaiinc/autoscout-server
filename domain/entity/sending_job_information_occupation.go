package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationOccupation struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	Occupation              null.Int  `db:"occupation" json:"occupation"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationOccupation(
	sendingJobInformationID uint,
	occupation null.Int,
) *SendingJobInformationOccupation {
	return &SendingJobInformationOccupation{
		SendingJobInformationID: sendingJobInformationID,
		Occupation:              occupation,
	}
}
