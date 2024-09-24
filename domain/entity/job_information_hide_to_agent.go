package entity

import "time"

type JobInformationHideToAgent struct {
	ID               uint      `db:"id" json:"id"`
	JobInformationID uint      `db:"job_information_id" json:"job_information_id"`
	AgentID          uint      `db:"agent_id" json:"agent_id"`
	AgentName        string    `db:"agent_name" json:"agent_name"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationHideToAgent(
	jobInformationID uint,
	agentID uint,
) *JobInformationHideToAgent {
	return &JobInformationHideToAgent{
		JobInformationID: jobInformationID,
		AgentID:          agentID,
	}
}
