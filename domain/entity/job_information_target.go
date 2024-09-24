package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationTarget struct {
	ID               uint      `db:"id" json:"-"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	Target           null.Int  `db:"target" json:"target"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationTarget(
	jobInformationID uint,
	target null.Int,
) *JobInformationTarget {
	return &JobInformationTarget{
		JobInformationID: jobInformationID,
		Target:           target,
	}
}
