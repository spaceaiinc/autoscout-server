package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerLicense struct {
	ID              uint      `db:"id" json:"id"`
	JobSeekerID     uint      `db:"job_seeker_id" json:"job_seeker_id"`
	LicenseType     null.Int  `db:"license_type" json:"license_type"`
	AcquisitionTime string    `db:"acquisition_time" json:"acquisition_time"`
	CreatedAt       time.Time `db:"created_at" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerLicense(
	jobSeekerID uint,
	licenseType null.Int,
	acquisitionTime string,
) *JobSeekerLicense {
	return &JobSeekerLicense{
		JobSeekerID:     jobSeekerID,
		LicenseType:     licenseType,
		AcquisitionTime: acquisitionTime,
	}
}
