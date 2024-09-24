package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationFeature struct {
	ID               uint      `db:"id" json:"-"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	Feature          null.Int  `db:"feature" json:"feature"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationFeature(
	jobInformationID uint,
	feature null.Int,
) *JobInformationFeature {
	return &JobInformationFeature{
		JobInformationID: jobInformationID,
		Feature:          feature,
	}
}
