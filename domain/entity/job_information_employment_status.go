package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationEmploymentStatus struct {
	ID               uint      `db:"id" json:"id"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	EmploymentStatus null.Int  `db:"employment_status" json:"employment_status"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationEmploymentStatus(
	jobInformationID uint,
	employmentStatus null.Int,
) *JobInformationEmploymentStatus {
	return &JobInformationEmploymentStatus{
		JobInformationID: jobInformationID,
		EmploymentStatus: employmentStatus,
	}
}
