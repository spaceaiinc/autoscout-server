package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type ChatGroupWithJobSeeker struct {
	ID                     uint      `db:"id" json:"id"`
	AgentID                uint      `db:"agent_id" json:"agent_id"`
	CAStaffID              null.Int  `db:"ca_staff_id" json:"ca_staff_id"`
	JobSeekerID            uint      `db:"job_seeker_id" json:"job_seeker_id"`
	LineID                 string    `db:"line_id" json:"-"`
	LastName               string    `db:"last_name" json:"last_name"`
	FirstName              string    `db:"first_name" json:"first_name"`
	LastFurigana           string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana          string    `db:"first_furigana" json:"first_furigana"`
	Email                  string    `db:"email" json:"email"`
	StaffName              string    `db:"staff_name" json:"staff_name"`
	AgentLastSendAt        time.Time `db:"agent_last_send_at" json:"agent_last_send_at"`                 // エージェントの最終送信時間
	AgentLastWatchedAt     time.Time `db:"agent_last_watched_at" json:"agent_last_watched_at"`           // エージェントの最終閲覧時間
	JobSeekerLastSendAt    time.Time `db:"job_seeker_last_send_at" json:"job_seeker_last_send_at"`       // 求職者の最終送信時間
	JobSeekerLastWatchedAt time.Time `db:"job_seeker_last_watched_at" json:"job_seeker_last_watched_at"` // 求職者の最終閲覧時間
	LineActive             bool      `db:"line_active" json:"line_active"`                               // LINE登録しているか
	Phase                  null.Int  `db:"phase" json:"phase"`                                           // 求職者の面談フェーズ
	JobSeekerUUID          uuid.UUID `db:"job_seeker_uuid" json:"job_seeker_uuid"`
	IDPhotoURL             string    `db:"id_photo_url" json:"id_photo_url"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`

	// db外
	IsBlocked              bool      `json:"is_blocked"`                                 // ブロックされているか
}

func NewChatGroupWithJobSeeker(
	agentID uint,
	JobInformationEmploymentStatusobSeekerID uint,
	lineActive bool,
) *ChatGroupWithJobSeeker {
	return &ChatGroupWithJobSeeker{
		AgentID:     agentID,
		JobSeekerID: JobInformationEmploymentStatusobSeekerID,
		LineActive:  lineActive,
	}
}

type CreateChatGroupWithJobSeekerParam struct {
	AgentID     uint `json:"agent_id" validate:"required"`
	JobSeekerID uint `json:"job_seeker_id" validate:"required"`
}

type SearchChatJobSeeker struct {
	FreeWord     string
	AgentStaffID string
	PhaseTypes   []null.Int
}

func NewSearchChatJobSeeker(
	freeword string,
	agentStaffID string,
	phaseTypes []null.Int,
) *SearchChatJobSeeker {
	return &SearchChatJobSeeker{
		FreeWord:     freeword,
		AgentStaffID: agentStaffID,
		PhaseTypes:   phaseTypes,
	}
}
