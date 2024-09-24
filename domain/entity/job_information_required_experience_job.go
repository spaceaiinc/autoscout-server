package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredExperienceJob struct {
	ID              uint      `db:"id" json:"id"`
	ConditionID     uint      `db:"condition_id" json:"condition_id"`
	ExperienceYear  null.Int  `db:"experience_year" json:"experience_year"`
	ExperienceMonth null.Int  `db:"experience_month" json:"experience_month"`
	CreatedAt       time.Time `db:"created_at" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"-"`

	// 関連テーブル
	ExperienceIndustries  []JobInformationRequiredExperienceIndustry   `db:"experience_industries" json:"experience_industries"`
	ExperienceOccupations []JobInformationRequiredExperienceOccupation `db:"experience_occupations" json:"experience_occupations"`
}

func NewJobInformationRequiredExperienceJob(
	conditionID uint,
	experienceYear null.Int,
	experienceMonth null.Int,
) *JobInformationRequiredExperienceJob {
	return &JobInformationRequiredExperienceJob{
		ConditionID:     conditionID,
		ExperienceYear:  experienceYear,
		ExperienceMonth: experienceMonth,
	}
}
