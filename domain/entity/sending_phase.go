package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// 送客進捗
type SendingPhase struct {
	ID                  uint      `db:"id" json:"id"`
	UUID                uuid.UUID `db:"uuid" json:"uuid"`
	SendingJobSeekerID  uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"` // 求職者のID
	SendingEnterpriseID uint      `db:"sending_enterprise_id" json:"sending_enterprise_id"` // 送客先企業のID
	Phase               null.Int  `db:"phase" json:"phase"`                                 // 進捗
	SendingDate         time.Time `db:"sending_date" json:"sending_date"`                   // 送客予定日時
	IsAttended          bool      `db:"is_attended" json:"is_attended"`                     // 送客予定日時
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`

	// db外項目
	LastName              string    `db:"last_name" json:"last_name"`
	FirstName             string    `db:"first_name" json:"first_name"`
	LastFurigana          string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana         string    `db:"first_furigana" json:"first_furigana"`
	AgentStaffID          null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName             string    `db:"staff_name" json:"staff_name"`
	AgentName             string    `db:"agent_name" json:"agent_name"`
	SenderAgentName       string    `db:"sender_agent_name" json:"sender_agent_name"`
	InterviewDate         time.Time `db:"interview_date" json:"interview_date"`
	ScheduleAdjustmentURL string    `db:"schedule_adjustment_url" json:"schedule_adjustment_url"` // 日程調整URL
	Commission            null.Int  `db:"commission" json:"commission"`                           // 送客単価

	// 終了理由
	EndReason string   `json:"end_reason" db:"end_reason"` // 終了理由
	EndStatus null.Int `json:"end_status" db:"end_status"` // 終了ステータス
}

func NewSendingPhase(
	sendingJobSeekerID uint,
	sendingEnterpriseID uint,
	phase null.Int,
	sendingDate time.Time,
	isAttended bool,
) *SendingPhase {
	return &SendingPhase{
		SendingJobSeekerID:  sendingJobSeekerID,
		SendingEnterpriseID: sendingEnterpriseID,
		Phase:               phase,
		SendingDate:         sendingDate,
		IsAttended:          isAttended,
	}
}

type CreateSendingPhaseParam struct {
	SendingJobSeekerID  uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id" validate:"required"` // 求職者のID
	SendingEnterpriseID uint      `db:"sending_enterprise_id" json:"sending_enterprise_id" validate:"required"` // 送客先企業のID
	Phase               null.Int  `db:"phase" json:"phase"`                                                     // 進捗
	SendingDate         time.Time `db:"sending_date" json:"sending_date"`                                       // 送客予定日時
}

// 進捗一覧(管理者側)表示用
type SendingJobSeekerTable struct {
	ID                    uint      `db:"id" json:"id"`
	SendingJobSeekerID    uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	JobSeekerName         string    `db:"job_seeker_name" json:"job_seeker_name"`
	AgentStaffID          null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName             string    `db:"staff_name" json:"staff_name"`
	AgentName             string    `db:"agent_name" json:"agent_name"`
	SenderAgentName       string    `db:"sender_agent_name" json:"sender_agent_name"`
	SendingDate           time.Time `db:"sending_date" json:"sending_date"`
	InterviewDate         time.Time `db:"interview_date" json:"interview_date"`
	IsNotWaitingViewed    bool      `db:"is_not_waiting_viewed" json:"is_not_waiting_viewed"`
	IsNotUnregisterViewed bool      `db:"is_not_unregister_viewed" json:"is_not_unregister_viewed"`
}

func NewSendingJobSeekerTable(
	id uint, // 送客応諾の時はphaseのidでそれ以外はsendingJobSeekerIDを入れる
	sendingJobSeekerID uint,
	jobSeekerName string,
	agentStaffID null.Int,
	staffName string,
	agentName string,
	senderAgentName string,
	sendingDate time.Time,
	interviewDate time.Time,
	isNotWaitingViewed bool,
	isNotUnregisterViewed bool,
) *SendingJobSeekerTable {
	return &SendingJobSeekerTable{
		ID:                    id,
		SendingJobSeekerID:    sendingJobSeekerID,
		JobSeekerName:         jobSeekerName,
		AgentStaffID:          agentStaffID,
		StaffName:             staffName,
		AgentName:             agentName,
		SenderAgentName:       senderAgentName,
		SendingDate:           sendingDate,
		InterviewDate:         interviewDate,
		IsNotWaitingViewed:    isNotWaitingViewed,
		IsNotUnregisterViewed: isNotUnregisterViewed,
	}
}

// 進捗一覧のパネルナンバー
type SendingJobSeekerTableTabNumber int64

const (
	// 送客完了
	CompleteSendingTab SendingJobSeekerTableTabNumber = iota

	// 送客応諾
	AcceptSendingTab

	// 面談実施済み
	PreparingAfterInterviewTab

	// 面談実施待ち
	WaitingForInterviewTab

	// 詳細未登録 / 日程未登録
	UnregisterTab

	// 終了
	CloseSendingTab
)
