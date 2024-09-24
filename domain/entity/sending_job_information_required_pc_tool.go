package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredPCTool struct {
	ID               uint      `db:"id" json:"id"`
	ConditionID      uint      `db:"condition_id" json:"condition_id"`
	Tool             null.Int  `db:"tool" json:"tool"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationRequiredPCTool(
	conditionID uint,
	tool null.Int,
) *SendingJobInformationRequiredPCTool {
	return &SendingJobInformationRequiredPCTool{
		ConditionID: conditionID,
		Tool:        tool,
	}
}
