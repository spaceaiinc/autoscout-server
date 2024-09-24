package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type AgentAllianceNeed struct {
	ID           uint      `db:"id" json:"id"`
	AgentID      uint      `db:"agent_id" json:"agent_id"`
	AllianceNeed null.Int  `db:"alliance_need" json:"alliance_need"` //アライアンスのニーズ （求人を提供したい,求職者を提供したい,求人を提供して欲しい,求職者を提供して欲しい）
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

func NewAgentAllianceNeed(
	agentID uint,
	allianceNeed null.Int,
) *AgentAllianceNeed {
	return &AgentAllianceNeed{
		AgentID:      agentID,
		AllianceNeed: allianceNeed,
	}
}
