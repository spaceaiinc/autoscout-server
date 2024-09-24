package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationOccupation struct {
	ID               uint      `json:"-" db:"id"`
	JobInformationID uint      `json:"-" db:"job_information_id"`
	Occupation       null.Int  `json:"occupation" db:"occupation"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationOccupation(
	jobInformationID uint,
	occupation null.Int,
) *JobInformationOccupation {
	return &JobInformationOccupation{
		JobInformationID: jobInformationID,
		Occupation:       occupation,
	}
}
