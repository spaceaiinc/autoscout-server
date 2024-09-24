package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingPhaseEndStatus struct {
	ID             uint      `json:"id" db:"id"`
	SendingPhaseID uint      `json:"sending_phase_id" db:"sending_phase_id"`
	EndReason      string    `json:"end_reason" db:"end_reason"` // 終了理由
	EndStatus      null.Int  `json:"end_status" db:"end_status"` // 終了ステータス
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func NewSendingPhaseEndStatus(
	sendingPhaseID uint,
	endReason string,
	endStatus null.Int,
) *SendingPhaseEndStatus {
	return &SendingPhaseEndStatus{
		SendingPhaseID: sendingPhaseID,
		EndReason:      endReason,
		EndStatus:      endStatus,
	}
}

type CreateSendingPhaseEndStatusParam struct {
	SendingPhaseID     uint     `json:"sending_phase_id" validate:"required"`
	EndReason          string   `json:"end_reason" validate:"required"`            // 終了理由
	EndStatus          null.Int `json:"end_status"`                                // 終了ステータス
}
