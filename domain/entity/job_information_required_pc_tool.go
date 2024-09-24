package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredPCTool struct {
	ID               uint      `db:"id" json:"-"`
	ConditionID      uint      `db:"condition_id" json:"condition_id"`
	Tool             null.Int  `db:"tool" json:"tool"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"` 
}

func NewJobInformationRequiredPCTool(
	conditionID uint,
	tool null.Int,
) *JobInformationRequiredPCTool {
	return &JobInformationRequiredPCTool{
		ConditionID: conditionID,
		Tool:        tool,
	}
}
