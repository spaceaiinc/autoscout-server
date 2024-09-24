package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// 面談時間タスクグループ
type InterviewAdjustmentTemplate struct {
	ID        uint      `db:"id" json:"id"`                 // 重複しないカラム毎のid
	AgentID   uint      `db:"agent_id" json:"agent_id"`     //エージェントのID
	SendScene null.Int  `db:"send_scene" json:"send_scene"` // 送信シーン（マスクレジュメ or 書類推薦 or 請求）
	Title     string    `db:"title" json:"title"`           // テンプレートのタイトル
	Subject   string    `db:"subject" json:"subject"`       // 件名
	Content   string    `db:"content" json:"content"`       // テンプレートの内容
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func NewInterviewAdjustmentTemplate(
	agentID uint,
	sendScene null.Int,
	title string,
	subject string,
	content string,
) *InterviewAdjustmentTemplate {
	return &InterviewAdjustmentTemplate{
		AgentID:   agentID,
		SendScene: sendScene,
		Title:     title,
		Subject:   subject,
		Content:   content,
	}
}

type CreateOrUpdateInterviewAdjustmentTemplateParam struct {
	ID        uint     `json:"id"`
	AgentID   uint     `json:"agent_id" validate:"required"`
	SendScene null.Int `json:"send_scene"`
	Title     string   `json:"title"`
	Subject   string   `json:"subject"`
	Content   string   `json:"content"`
}
