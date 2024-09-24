package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationEmploymentStatus struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	EmploymentStatus        null.Int  `db:"employment_status" json:"employment_status"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationEmploymentStatus(
	sendingJobInformationID uint,
	employmentStatus null.Int,
) *SendingJobInformationEmploymentStatus {
	return &SendingJobInformationEmploymentStatus{
		SendingJobInformationID: sendingJobInformationID,
		EmploymentStatus:        employmentStatus,
	}
}
