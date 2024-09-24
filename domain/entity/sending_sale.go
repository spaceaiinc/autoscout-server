package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingSale struct {
	ID                  uint      `db:"id" json:"id"`
	SendingJobSeekerID  uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"` // 送客求職者のID
	SendingEnterpriseID uint      `db:"sending_enterprise_id" json:"sending_enterprise_id"` // 送客先エージェントのID
	MotoyuiSales        null.Int  `db:"motoyui_sales" json:"motoyui_sales"`                 // Motoyuiの売上
	Kickback            null.Int  `db:"kickback" json:"kickback"`                           // キックバック（送客元の売上）
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`

	// db外
	SendingDate       time.Time `db:"sending_date" json:"sending_date"`               // 送客確定日(sending_phasesのsending_date)
	JobSeekerName     string    `db:"job_seeker_name" json:"job_seeker_name"`         // 送客求職者の名前
	ReceiverAgentName string    `db:"receiver_agent_name" json:"receiver_agent_name"` // 送客先エージェントの名前
	SenderAgentName   string    `db:"sender_agent_name" json:"sender_agent_name"`     // 送客元エージェントの名前
}

func NewSendingSale(
	sendingJobSeekerID uint,
	sendingEnterpriseID uint,
	motoyuiSales null.Int,
	kickback null.Int,
) *SendingSale {
	return &SendingSale{
		SendingJobSeekerID:  sendingJobSeekerID,
		SendingEnterpriseID: sendingEnterpriseID,
		MotoyuiSales:        motoyuiSales,
		Kickback:            kickback,
	}
}

type CreateSendingSaleParam struct {
	SendingJobSeekerID  uint     `db:"sending_job_seeker_id" json:"sending_job_seeker_id" validate:"required"` // 送客求職者のID
	SendingEnterpriseID uint     `db:"sending_enterprise_id" json:"sending_enterprise_id" validate:"required"` // 送客先エージェントのID
	MotoyuiSales        null.Int `db:"motoyui_sales" json:"motoyui_sales"`                                     // Motoyuiの売上
	Kickback            null.Int `db:"kickback" json:"kickback"`                                               // キックバック（送客元の売上）
}

type UpdateSendingSaleParam struct {
	MotoyuiSales null.Int `db:"motoyui_sales" json:"motoyui_sales"` // Motoyuiの売上
	Kickback     null.Int `db:"kickback" json:"kickback"`           // キックバック（送客元の売上）
}
