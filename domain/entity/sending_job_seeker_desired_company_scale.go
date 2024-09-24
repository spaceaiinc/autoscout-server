package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerDesiredCompanyScale struct {
	ID                  uint      `db:"id" json:"id"`
	SendingJobSeekerID  uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	DesiredCompanyScale null.Int  `db:"desired_company_scale" json:"desired_company_scale"`
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerDesiredCompanyScale(
	sendingJobSeekerID uint,
	desiredCompanyScale null.Int,
) *SendingJobSeekerDesiredCompanyScale {
	return &SendingJobSeekerDesiredCompanyScale{
		SendingJobSeekerID:  sendingJobSeekerID,
		DesiredCompanyScale: desiredCompanyScale,
	}
}
