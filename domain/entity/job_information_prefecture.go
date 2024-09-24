package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationPrefecture struct {
	ID               uint      `db:"id" json:"id"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	Prefecture       null.Int  `db:"prefecture" json:"prefecture"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationPrefecture(
	jobInformationID uint,
	prefecture null.Int,
) *JobInformationPrefecture {
	return &JobInformationPrefecture{
		JobInformationID: jobInformationID,
		Prefecture:       prefecture,
	}
}
