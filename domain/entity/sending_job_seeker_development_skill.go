package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerDevelopmentSkill struct {
	ID                  uint      `db:"id" json:"id"`
	SendingJobSeekerID  uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	DevelopmentCategory null.Int  `db:"development_category" json:"development_category"`
	DevelopmentType     null.Int  `db:"development_type" json:"development_type"`
	ExperienceYear      null.Int  `db:"experience_year" json:"experience_year"`
	ExperienceMonth     null.Int  `db:"experience_month" json:"experience_month"`
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerDevelopmentSkill(
	sendingJobSeekerID uint,
	developmentCategory null.Int,
	developmentType null.Int,
	experienceYear null.Int,
	experienceMonth null.Int,
) *SendingJobSeekerDevelopmentSkill {
	return &SendingJobSeekerDevelopmentSkill{
		SendingJobSeekerID:  sendingJobSeekerID,
		DevelopmentCategory: developmentCategory,
		DevelopmentType:     developmentType,
		ExperienceYear:      experienceYear,
		ExperienceMonth:     experienceMonth,
	}
}
