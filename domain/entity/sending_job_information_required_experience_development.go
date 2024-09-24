package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredExperienceDevelopment struct {
	ID                  uint      `db:"id" json:"id"`
	ConditionID         uint      `db:"condition_id" json:"condition_id"`
	DevelopmentCategory null.Int  `db:"development_category" json:"development_category"`
	ExperienceYear      null.Int  `db:"experience_year" json:"experience_year"`
	ExperienceMonth     null.Int  `db:"experience_month" json:"experience_month"`
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`

	// 関連テーブル
	ExperienceDevelopmentTypes []SendingJobInformationRequiredExperienceDevelopmentType `db:"-" json:"experience_development_types"`
}

func NewSendingJobInformationRequiredExperienceDevelopment(
	conditionID uint,
	developmentCategory null.Int,
	experienceYear null.Int,
	experienceMonth null.Int,
) *SendingJobInformationRequiredExperienceDevelopment {
	return &SendingJobInformationRequiredExperienceDevelopment{
		ConditionID:         conditionID,
		DevelopmentCategory: developmentCategory,
		ExperienceYear:      experienceYear,
		ExperienceMonth:     experienceMonth,
	}
}
