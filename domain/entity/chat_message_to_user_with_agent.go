package entity

import (
	"time"
)

type ChatMessageToUserWithAgent struct {
	ID           uint      `db:"id" json:"id"`
	MessageID    uint      `db:"message_id" json:"message_id"`         // メッセージID
	AgentStaffID uint      `db:"agent_staff_id" json:"agent_staff_id"` //エージェントスタッフのid
	StaffName    string    `db:"staff_name" json:"staff_name"`         // 担当者名
	SendAt       time.Time `db:"send_at" json:"send_at"`               // メッセージ送信日時
	WatchedAt    time.Time `db:"watched_at" json:"watched_at"`         // メッセージ閲覧日時
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

	// dbにはない
	// IsWatched bool `db:"-" json:"is_watched"` // メッセージが既読かどうか
	GroupID  uint `db:"group_id" json:"group_id"`   // チャットグループID 未読判定に使用
	ThreadID uint `db:"thread_id" json:"thread_id"` // スレッドID　未読判定に使用
}

func NewChatMessageToUserWithAgent(
	messageID uint,
	agentStaffID uint,
	// sendAt time.Time,
	// watchedAt time.Time,
) *ChatMessageToUserWithAgent {
	return &ChatMessageToUserWithAgent{
		MessageID:    messageID,
		AgentStaffID: agentStaffID,
		// SendAt:       sendAt,
		// WatchedAt:    watchedAt,
	}
}
