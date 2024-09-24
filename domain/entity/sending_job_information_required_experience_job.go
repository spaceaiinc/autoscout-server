package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredExperienceJob struct {
	ID              uint      `db:"id" json:"id"`
	ConditionID     uint      `db:"condition_id" json:"condition_id"`
	ExperienceYear  null.Int  `db:"experience_year" json:"experience_year"`
	ExperienceMonth null.Int  `db:"experience_month" json:"experience_month"`
	CreatedAt       time.Time `db:"created_at" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"-"`

	// 関連テーブル
	ExperienceIndustries  []SendingJobInformationRequiredExperienceIndustry   `db:"experience_industries" json:"experience_industries"`
	ExperienceOccupations []SendingJobInformationRequiredExperienceOccupation `db:"experience_occupations" json:"experience_occupations"`
}

func NewSendingJobInformationRequiredExperienceJob(
	conditionID uint,
	experienceYear null.Int,
	experienceMonth null.Int,
) *SendingJobInformationRequiredExperienceJob {
	return &SendingJobInformationRequiredExperienceJob{
		ConditionID:     conditionID,
		ExperienceYear:  experienceYear,
		ExperienceMonth: experienceMonth,
	}
}
