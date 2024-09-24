package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredExperienceDevelopmentType struct {
	ID                      uint      `db:"id" json:"id"`
	ExperienceDevelopmentID uint      `db:"experience_development_id" json:"experience_development_id"`
	DevelopmentType         null.Int  `db:"development_type" json:"development_type"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationRequiredExperienceDevelopmentType(
	experienceDevelopmentID uint,
	developmentType null.Int,
) *JobInformationRequiredExperienceDevelopmentType {
	return &JobInformationRequiredExperienceDevelopmentType{
		ExperienceDevelopmentID: experienceDevelopmentID,
		DevelopmentType:         developmentType,
	}
}
