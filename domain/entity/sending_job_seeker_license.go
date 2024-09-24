package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerLicense struct {
	ID                 uint      `db:"id" json:"id"`
	SendingJobSeekerID uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	LicenseType        null.Int  `db:"license_type" json:"license_type"`
	AcquisitionTime    string    `db:"acquisition_time" json:"acquisition_time"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerLicense(
	sendingJobSeekerID uint,
	licenseType null.Int,
	acquisitionTime string,
) *SendingJobSeekerLicense {
	return &SendingJobSeekerLicense{
		SendingJobSeekerID: sendingJobSeekerID,
		LicenseType:        licenseType,
		AcquisitionTime:    acquisitionTime,
	}
}
