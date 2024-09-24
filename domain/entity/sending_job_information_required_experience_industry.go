package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredExperienceIndustry struct {
	ID                 uint      `db:"id" json:"id"`
	ExperienceJobID    uint      `db:"experience_job_id" json:"experience_job_id"`
	ExperienceIndustry null.Int  `db:"experience_industry" json:"experience_industry"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationRequiredExperienceIndustry(
	experienceJobID uint,
	experienceIndustry null.Int,
) *SendingJobInformationRequiredExperienceIndustry {
	return &SendingJobInformationRequiredExperienceIndustry{
		ExperienceJobID:    experienceJobID,
		ExperienceIndustry: experienceIndustry,
	}
}
