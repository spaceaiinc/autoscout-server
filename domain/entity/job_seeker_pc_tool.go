package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerPCTool struct {
	ID          uint      `json:"id" db:"id"`
	JobSeekerID uint      `json:"job_seeker_id" db:"job_seeker_id"`
	Tool        null.Int  `json:"tool" db:"tool"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewJobSeekerPCTool(
	jobSeekerID uint,
	tool null.Int,
) *JobSeekerPCTool {
	return &JobSeekerPCTool{
		JobSeekerID: jobSeekerID,
		Tool:        tool,
	}
}
