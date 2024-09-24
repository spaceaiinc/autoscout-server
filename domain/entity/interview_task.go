package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// 面談調整タスク
type InterviewTask struct {
	ID                   uint      `db:"id" json:"id"`
	InterviewTaskGroupID uint      `db:"interview_task_group_id" json:"interview_task_group_id"`
	AgentStaffID         null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	CAStaffID            null.Int  `db:"ca_staff_id" json:"ca_staff_id"`
	CAStaffName          string    `db:"ca_staff_name" json:"ca_staff_name"`
	PhaseCategory        null.Int  `db:"phase_category" json:"phase_category"`
	PhaseSubCategory     null.Int  `db:"phase_sub_category" json:"phase_sub_category"`
	Remarks              string    `db:"remarks" json:"remarks"`
	DeadlineDay          string    `db:"deadline_day" json:"deadline_day"`
	DeadlineTime         null.Int  `db:"deadline_time" json:"deadline_time"`
	SelectActionLabel    string    `db:"select_action_label" json:"select_action_label"` // 面談調整タスクで選択したアクションのラベル
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`

	// エージェント
	StaffName                string    `db:"staff_name" json:"staff_name"`
	AgentID                  uint      `db:"agent_id" json:"agent_id"`
	AgentName                string    `db:"agent_name" json:"agent_name"`
	AgentUUID                uuid.UUID `db:"agent_uuid" json:"agent_uuid"`
	InterviewAdjustmentEmail string    `db:"interview_adjustment_email" json:"interview_adjustment_email"` // メールアドレス（面談調整用）

	// 面談項目
	InterviewDate      time.Time `db:"interview_date" json:"interview_date"`
	FirstInterviewDate time.Time `db:"first_interview_date" json:"first_interview_date"`

	// 求職者項目
	JobSeekerID   uint      `db:"job_seeker_id" json:"job_seeker_id"`
	JobSeekerUUID uuid.UUID `db:"job_seeker_uuid" json:"job_seeker_uuid"`
	LastName      string    `db:"last_name" json:"last_name"`
	FirstName     string    `db:"first_name" json:"first_name"`
	LastFurigana  string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana string    `db:"first_furigana" json:"first_furigana"`
	Phase         null.Int  `db:"phase" json:"phase"`
	Email         string    `db:"email" json:"email"`
	// LineID             string    `db:"line_id" json:"-"`
	LineActive         bool      `db:"line_active" json:"line_active"` // LINEがブロックされているか
	PhoneNumber        string    `db:"phone_number" json:"phone_number"`
	JobSeekerCreatedAt time.Time `db:"job_seeker_created_at" json:"job_seeker_created_at"`

	// タスクグループ
	LastRequestAt time.Time `db:"last_request_at" json:"last_request_at"` // 最終依頼時間
	LastWatchedAt time.Time `db:"last_watched_at" json:"last_watched_at"` // 最終閲覧時間
}

func NewInterviewTask(
	interviewTaskGroupID uint,
	agentStaffID null.Int,
	caStaffID null.Int,
	phaseCategory null.Int,
	phaseSubCategory null.Int,
	remarks string,
	deadlineDay string,
	deadlineTime null.Int,
	selectActionLabel string,
) *InterviewTask {
	return &InterviewTask{
		InterviewTaskGroupID: interviewTaskGroupID,
		AgentStaffID:         agentStaffID,
		CAStaffID:            caStaffID,
		PhaseCategory:        phaseCategory,
		PhaseSubCategory:     phaseSubCategory,
		Remarks:              remarks,
		DeadlineDay:          deadlineDay,
		DeadlineTime:         deadlineTime,
		SelectActionLabel:    selectActionLabel,
	}
}

// 次のタスク
type NextInterviewTaskParam struct {
	InterviewTaskGroupID uint      `db:"interview_task_group_id" json:"interview_task_group_id"`
	JobSeekerID          uint      `json:"job_seeker_id" validate:"required"`
	AgentID              uint      `json:"agent_id" validate:"required"`
	AgentStaffID         null.Int  `db:"agent_staff_id" json:"agent_staff_id"`
	CAStaffID            null.Int  `db:"ca_staff_id" json:"ca_staff_id"`
	PhaseCategory        null.Int  `db:"phase_category" json:"phase_category"`
	PhaseSubCategory     null.Int  `db:"phase_sub_category" json:"phase_sub_category"`
	Remarks              string    `db:"remarks" json:"remarks"`
	DeadlineDay          string    `db:"deadline_day" json:"deadline_day"`
	DeadlineTime         null.Int  `db:"deadline_time" json:"deadline_time"`
	SelectActionLabel    string    `db:"select_action_label" json:"select_action_label"` // 面談調整タスクで選択したアクションのラベル
	InterviewDate        time.Time `json:"interview_date"`
	FirstInterviewDate   time.Time `json:"first_interview_date"`
	ValidationType       string    `json:"validation_type"` // 'email' or 'interview' or 'emailAndInterview'

	// メール
	Mail Mail `json:"mail"`
}

type Mail struct {
	From    EmailUser `json:"from"`    // 送信元の情報（エージェント名 + エージェントメールアドレス（面談調整用））
	To      EmailUser `json:"to"`      // 送信先の情報（）
	Subject string    `json:"subject"` // メールの件名
	Content string    `json:"content"` // メールの本文
}

type UpdateCAStaffIDParam struct {
	JobSeekerID     uint     `json:"job_seeker_id" validate:"required"`
	CAStaffID       null.Int `json:"ca_staff_id" validate:"required"`
	InterviewTaskID uint     `json:"interview_task_id" validate:"required"`
}

type UpdateInterviewDateParam struct {
	InterviewTaskID      uint      `json:"interview_task_id" validate:"required"`
	InterviewTaskGroupID uint      `json:"interview_task_group_id" validate:"required"`
	InterviewDate        time.Time `json:"interview_date" validate:"required"`
}

type DeleteLatestInterviewTaskParam struct {
	InterviewTaskID      uint `json:"interview_task_id" validate:"required"`
	InterviewTaskGroupID uint `json:"interview_task_group_id" validate:"required"`
}
