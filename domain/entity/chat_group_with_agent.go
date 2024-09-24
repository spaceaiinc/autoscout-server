package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatGroupWithAgent struct {
	ID         uint      `db:"id" json:"id"`
	UUID       uuid.UUID `db:"uuid" json:"uuid"` // グループのuuid
	Agent1ID   uint      `db:"agent1_id" json:"agent1_id"`
	Agent2ID   uint      `db:"agent2_id" json:"agent2_id"`
	LastSendAt time.Time `db:"last_send_at" json:"last_send_at"` // エージェントの最終送信時間
	// Agent1LastSendAt    time.Time `db:"agent1_last_send_at" json:"agent1_last_send_at"`       // エージェントの最終送信時間
	// Agent1LastWatchedAt time.Time `db:"agent1_last_watched_at" json:"agent1_last_watched_at"` // エージェントの最終閲覧時間
	// Agent2LastSendAt    time.Time `db:"agent2_last_send_at" json:"agent2_last_send_at"`       // エージェントの最終送信時間
	// Agent2LastWatchedAt time.Time `db:"agent2_last_watched_at" json:"agent2_last_watched_at"` // エージェントの最終閲覧時間
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	AgentID   uint      `db:"agent_id" json:"agent_id"` // チャット相手のIDが入る
	AgentName string    `db:"agent_name" json:"agent_name"` // チャット相手のエージェント名が入る

	// 未読通知で使用
	Agent1Name string `db:"agent1_name" json:"agent1_name"`
	Agent2Name string `db:"agent2_name" json:"agent2_name"`

	Threads []ChatThreadWithAgent `db:"threads" json:"threads"`

	UnwatchedCount uint `db:"unwatched_count" json:"unwatched_count"` // 未読件数
}

func NewChatGroupWithAgent(
	agent1ID uint,
	agent2ID uint,
) *ChatGroupWithAgent {
	return &ChatGroupWithAgent{
		Agent1ID: agent1ID,
		Agent2ID: agent2ID,
	}
}

type CreateChatGroupWithAgentParam struct {
	Agent1ID uint `json:"agent1_id" validate:"required"`
	Agent2ID uint `json:"agent2_id" validate:"required"`
}
