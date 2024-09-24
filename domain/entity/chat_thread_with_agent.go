package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatThreadWithAgent struct {
	ID           uint      `db:"id" json:"id"`
	UUID         uuid.UUID `db:"uuid" json:"uuid"`                     // スレッドのuuid
	GroupID      uint      `db:"group_id" json:"group_id"`             // チャットグループID
	AgentStaffID uint      `db:"agent_staff_id" json:"agent_staff_id"` //エージェントスタッフのid
	StaffName    string    `db:"staff_name" json:"staff_name"`         // 担当者名
	Title        string    `db:"title" json:"title"`                   // スレッドタイトル
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	Messages []ChatMessageWithAgent `db:"messages" json:"messages"`

	// dbにはない
	UnwatchedCount    uint      `db:"unwatched_count" json:"unwatched_count"`         // 未読件数
	MessageCOunt      uint      `db:"message_count" json:"message_count"`             // メッセージ件数
	LatestMessageTime time.Time `db:"latest_message_time" json:"latest_message_time"` // 最新メッセージの時間(messagesテーブルのcreated_at、並び替え用)
}

func NewChatThreadWithAgent(
	groupID uint,
	agentStaffID uint,
	title string,
) *ChatThreadWithAgent {
	return &ChatThreadWithAgent{
		GroupID:      groupID,
		AgentStaffID: agentStaffID,
		Title:        title,
	}
}

type CreateChatThreadWithAgentParam struct {
	GroupID      uint                   `json:"group_id" validate:"required"`
	AgentStaffID uint                   `json:"agent_staff_id" validate:"required"`
	Title        string                 `json:"title" validate:"required"`
	Messages     []ChatMessageWithAgent `json:"messages"`
}
