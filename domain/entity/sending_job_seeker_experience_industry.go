package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerExperienceIndustry struct {
	ID            uint      `db:"id" json:"id"`
	WorkHistoryID uint      `db:"work_history_id" json:"work_history_id"`
	Industry      null.Int  `db:"industry" json:"industry"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerExperienceIndustry(
	workHistoryID uint,
	industry null.Int,
) *SendingJobSeekerExperienceIndustry {
	return &SendingJobSeekerExperienceIndustry{
		WorkHistoryID: workHistoryID,
		Industry:      industry,
	}
}
