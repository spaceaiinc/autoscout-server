package entity

import (
	"time"

	"github.com/google/uuid"
)

type AgentRobot struct {
	ID            uint      `db:"id" json:"id"`
	UUID          uuid.UUID `db:"uuid" json:"uuid"`
	AgentID       uint      `db:"agent_id" json:"agent_id"`
	Name          string    `db:"name" json:"name"`                       // ロボット名
	IsEntryActive bool      `db:"is_entry_active" json:"is_entry_active"` // アクティブかどうか/false:走らせない true:走る *falseの場合はエントリー作成できない
	IsScoutActive bool      `db:"is_scout_active" json:"is_scout_active"` // アクティブかどうか/false:走らせない true:走る *falseの場合はスカウト作成できない
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	ScoutServices []ScoutService `json:"scout_services"`

	//結合項目
	AgentName string `db:"agent_name" json:"agent_name"`
}

func NewAgentRobot(
	agentID uint,
	name string,
	IsEntryActive bool,
	IsScoutActive bool,
) *AgentRobot {
	return &AgentRobot{
		AgentID:       agentID,
		Name:          name,
		IsEntryActive: IsEntryActive,
		IsScoutActive: IsScoutActive,
	}
}

type CreateOrUpdateAgentRobotParam struct {
	AgentID       uint   `json:"agent_id"`
	Name          string `json:"name"`
	IsEntryActive bool   `json:"is_entry_active"` // アクティブかどうか/false:走らせない true:走る *falseの場合はエントリー作成できない
	IsScoutActive bool   `json:"is_scout_active"` // アクティブかどうか/false:走らせない true:走る *falseの場合はスカウト作成できない

	// 関連テーブル
	ScoutServices []ScoutService `json:"scout_services"`
}
