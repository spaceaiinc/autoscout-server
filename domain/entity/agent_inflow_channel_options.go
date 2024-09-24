package entity

import (
	"time"
)

// エージェントの流入経路マスタを管理するテーブル
type AgentInflowChannelOption struct {
	ID          uint      `db:"id" json:"id"`
	AgentID     uint      `db:"agent_id" json:"agent_id"`
	ChannelName string    `db:"channel_name" json:"channel_name"` // 流入経路の名前
	IsOpen      bool      `db:"is_open" json:"is_open"`           // Open / Close
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func NewAgentInflowChannelOption(
	agentID uint,
	ChannelName string,
	IsOpen bool,
) *AgentInflowChannelOption {
	return &AgentInflowChannelOption{
		AgentID:       agentID,
		ChannelName:   ChannelName,
		IsOpen:        IsOpen,
	}
}
