package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredExperienceOccupation struct {
	ID                   uint      `db:"id" json:"id"`
	ExperienceJobID      uint      `db:"experience_job_id" json:"experience_job_id"`
	ExperienceOccupation null.Int  `db:"experience_occupation" json:"experience_occupation"`
	CreatedAt            time.Time `db:"created_at" json:"-"`
	UpdatedAt            time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationRequiredExperienceOccupation(
	experienceJobID uint,
	experienceOccupation null.Int,
) *SendingJobInformationRequiredExperienceOccupation {
	return &SendingJobInformationRequiredExperienceOccupation{
		ExperienceJobID:      experienceJobID,
		ExperienceOccupation: experienceOccupation,
	}
}
