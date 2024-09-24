package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerDevelopmentSkill struct {
	ID                  uint      `db:"id" json:"id"`
	JobSeekerID         uint      `db:"job_seeker_id" json:"job_seeker_id"`
	DevelopmentCategory null.Int  `db:"development_category" json:"development_category"`
	DevelopmentType     null.Int  `db:"development_type" json:"development_type"`
	ExperienceYear      null.Int  `db:"experience_year" json:"experience_year"`
	ExperienceMonth     null.Int  `db:"experience_month" json:"experience_month"`
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerDevelopmentSkill(
	jobSeekerID uint,
	developmentCategory null.Int,
	developmentType null.Int,
	experienceYear null.Int,
	experienceMonth null.Int,
) *JobSeekerDevelopmentSkill {
	return &JobSeekerDevelopmentSkill{
		JobSeekerID:         jobSeekerID,
		DevelopmentCategory: developmentCategory,
		DevelopmentType:     developmentType,
		ExperienceYear:      experienceYear,
		ExperienceMonth:     experienceMonth,
	}
}
