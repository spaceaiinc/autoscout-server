package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerEndStatus struct {
	ID                 uint      `json:"id" db:"id"`
	SendingJobSeekerID uint      `json:"sending_job_seeker_id" db:"sending_job_seeker_id"`
	EndReason          string    `json:"end_reason" db:"end_reason"` // 終了理由
	EndStatus          null.Int  `json:"end_status" db:"end_status"` // 終了ステータス
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

func NewSendingJobSeekerEndStatus(
	sendingJobSeekerID uint,
	endReason string,
	endStatus null.Int,
) *SendingJobSeekerEndStatus {
	return &SendingJobSeekerEndStatus{
		SendingJobSeekerID: sendingJobSeekerID,
		EndReason:          endReason,
		EndStatus:          endStatus,
	}
}

type CreateSendingJobSeekerEndStatusParam struct {
	SendingJobSeekerID uint     `json:"sending_job_seeker_id" validate:"required"`
	EndReason          string   `json:"end_reason" validate:"required"` // 終了理由
	EndStatus          null.Int `json:"end_status"`                    // 終了ステータス
}
