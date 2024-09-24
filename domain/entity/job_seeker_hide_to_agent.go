package entity

import "time"

type JobSeekerHideToAgent struct {
	ID          uint      `db:"id" json:"id"`
	JobSeekerID uint      `db:"job_seeker_id" json:"job_seeker_id"`
	AgentID     uint      `db:"agent_id" json:"agent_id"`
	AgentName   string    `db:"agent_name" json:"agent_name"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerHideToAgent(
	jobSeekerID uint,
	agentID uint,
) *JobSeekerHideToAgent {
	return &JobSeekerHideToAgent{
		JobSeekerID: jobSeekerID,
		AgentID:     agentID,
	}
}
