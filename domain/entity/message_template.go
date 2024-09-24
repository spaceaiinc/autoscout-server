package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type MessageTemplate struct {
	ID           uint      `db:"id" json:"id"`
	AgentStaffID uint      `db:"agent_staff_id" json:"agent_staff_id"`
	SendScene    null.Int  `db:"send_scene" json:"send_scene"`
	Title        string    `db:"title" json:"title"`
	Subject      string    `db:"subject" json:"subject"` // 件名
	Content      string    `db:"content" json:"content"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func NewMessageTemplate(
	agentStaffID uint,
	sendScene null.Int,
	title string,
	subject string,
	content string,
) *MessageTemplate {
	return &MessageTemplate{
		AgentStaffID: agentStaffID,
		SendScene:    sendScene,
		Title:        title,
		Subject:      subject,
		Content:      content,
	}
}

type CreateOrUpdateMessageTemplateParam struct {
	AgentStaffID uint     `json:"agent_staff_id" validate:"required"`
	SendScene    null.Int `json:"send_scene" validate:"required"`
	Title        string   `json:"title" validate:"required"`
	Subject      string   `json:"subject"`
	Content      string   `json:"content" validate:"required"`
}
