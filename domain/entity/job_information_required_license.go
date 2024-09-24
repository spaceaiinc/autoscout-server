package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredLicense struct {
	ID          uint      `db:"id" json:"id"`
	ConditionID uint      `db:"condition_id" json:"condition_id"`
	License     null.Int  `db:"license" json:"license"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationRequiredLicense(
	conditionID uint,
	license null.Int,
) *JobInformationRequiredLicense {
	return &JobInformationRequiredLicense{
		ConditionID: conditionID,
		License:     license,
	}
}
