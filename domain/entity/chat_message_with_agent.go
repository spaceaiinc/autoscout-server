package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessageWithAgent struct {
	ID           uint      `db:"id" json:"id"`
	UUID         uuid.UUID `db:"uuid" json:"uuid"`                     // メッセージのuuid
	GroupID      uint      `db:"group_id" json:"group_id"`             // チャットグループID
	FileURL      string    `db:"file_url" json:"file_url"`             //URL
	ThreadID     uint      `db:"thread_id" json:"thread_id"`           // スレッドID
	AgentStaffID uint      `db:"agent_staff_id" json:"agent_staff_id"` //エージェントスタッフのid
	StaffName    string    `db:"staff_name" json:"staff_name"`         // 担当者名
	Message      string    `db:"message" json:"message"`               // メッセージ
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

	NotWatched bool `db:"-" json:"not_watched"` // 未読の場合はtrue

	MessageToUsers []ChatMessageToUserWithAgent `db:"message_to_users" json:"message_to_users"`
}

func NewChatMessageWithAgent(
	threadID uint,
	agentStaffID uint,
	message string,
	fileURL string,
) *ChatMessageWithAgent {
	return &ChatMessageWithAgent{
		ThreadID:     threadID,
		AgentStaffID: agentStaffID,
		Message:      message,
		FileURL:      fileURL,
	}
}

type SendChatMessageWithAgentParam struct {
	// GroupID         uint   `json:"group_id" validate:"required"`       // チャットグループID
	ThreadID      uint         `json:"thread_id" validate:""`              // スレッドID
	AgentStaffID  uint         `json:"agent_staff_id" validate:"required"` // エージェントID
	Message       string       `json:"message" validate:"required"`        // メッセージ
	FileURL       string       `json:"file_url"`
	ToAgentStaffs []AgentStaff `json:"to_agent_staffs"` // 宛先に選択した担当者ID
	// ReAgentStaffIDs []uint `json:"re_agent_staff_ids"`                 // 返信先に選択した担当者ID
}

type UpdateWatchedAtParam struct {
	ThreadID     uint `json:"thread_id" validate:"required"`      // スレッドID
	AgentStaffID uint `json:"agent_staff_id" validate:"required"` // エージェントID
}
