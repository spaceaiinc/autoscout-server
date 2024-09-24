package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerPCTool struct {
	ID          uint      `json:"id" db:"id"`
	SendingJobSeekerID uint      `json:"sending_job_seeker_id" db:"sending_job_seeker_id"`
	Tool        null.Int  `json:"tool" db:"tool"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewSendingJobSeekerPCTool(
	sendingJobSeekerID uint,
	tool null.Int,
) *SendingJobSeekerPCTool {
	return &SendingJobSeekerPCTool{
		SendingJobSeekerID: sendingJobSeekerID,
		Tool:        tool,
	}
}
