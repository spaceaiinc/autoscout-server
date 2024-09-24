package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type ChatGroupWithSendingJobSeeker struct {
	ID                            uint      `db:"id" json:"id"`
	SendingJobSeekerID            uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	AgentLastSendAt               time.Time `db:"agent_last_send_at" json:"agent_last_send_at"`                                 // エージェントの最終送信時間
	AgentLastWatchedAt            time.Time `db:"agent_last_watched_at" json:"agent_last_watched_at"`                           // エージェントの最終閲覧時間
	SendingJobSeekerLastSendAt    time.Time `db:"sending_job_seeker_last_send_at" json:"sending_job_seeker_last_send_at"`       // 求職者の最終送信時間
	SendingJobSeekerLastWatchedAt time.Time `db:"sending_job_seeker_last_watched_at" json:"sending_job_seeker_last_watched_at"` // 求職者の最終閲覧時間
	LineActive                    bool      `db:"line_active" json:"line_active"`                                               // LINEがブロックされているか
	CreatedAt                     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                     time.Time `db:"updated_at" json:"updated_at"`

	// db外
	LineID        string   `db:"line_id" json:"-"`
	LastName      string   `db:"last_name" json:"last_name"`
	FirstName     string   `db:"first_name" json:"first_name"`
	LastFurigana  string   `db:"last_furigana" json:"last_furigana"`
	FirstFurigana string   `db:"first_furigana" json:"first_furigana"`
	Email         string   `db:"email" json:"email"`
	CAStaffID     null.Int `db:"ca_staff_id" json:"ca_staff_id"`
	StaffName     string   `db:"staff_name" json:"staff_name"`
	Phase         null.Int `db:"phase" json:"phase"` // 求職者の面談フェーズ
	IsBlocked     bool     `json:"is_blocked"`
}

func NewChatGroupWithSendingJobSeeker(
	sendingJobSeekerID uint,
	lineActive bool,
) *ChatGroupWithSendingJobSeeker {
	return &ChatGroupWithSendingJobSeeker{
		SendingJobSeekerID: sendingJobSeekerID,
		LineActive:         lineActive,
	}
}

type CreateChatGroupWithSendingJobSeekerParam struct {
	SendingJobSeekerID uint `json:"sending_job_seeker_id" validate:"required"`
}

type SearchChatSendingJobSeeker struct {
	FreeWord     string
	AgentStaffID string
	PhaseTypes   []null.Int
}

func NewSearchChatSendingJobSeeker(
	freeword string,
	agentStaffID string,
	phaseTypes []null.Int,
) *SearchChatSendingJobSeeker {
	return &SearchChatSendingJobSeeker{
		FreeWord:     freeword,
		AgentStaffID: agentStaffID,
		PhaseTypes:   phaseTypes,
	}
}
