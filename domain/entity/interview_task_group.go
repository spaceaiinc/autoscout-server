package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// 面談時間タスクグループ
type InterviewTaskGroup struct {
	ID                 uint      `db:"id" json:"id"`                                     // 重複しないカラム毎のid
	UUID               uuid.UUID `db:"uuid" json:"uuid"`                                 //重複しないカラム毎のUUID
	AgentID            uint      `db:"agent_id" json:"agent_id"`                         //エージェントのID
	JobSeekerID        uint      `db:"job_seeker_id" json:"job_seeker_id"`               //求職者のID
	InterviewDate      time.Time `db:"interview_date" json:"interview_date"`             //面談日時
	FirstInterviewDate time.Time `db:"first_interview_date" json:"first_interview_date"` //初回面談日時(KPIの計算の基準になる)
	LastRequestAt      time.Time `db:"last_request_at" json:"last_request_at"`           //最終タスク依頼時間
	LastWatchedAt      time.Time `db:"last_watched_at" json:"last_watched_at"`           //最終タスク閲覧時間
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`

	// 求職者テーブル
	LastName      string   `db:"last_name" json:"last_name"`
	FirstName     string   `db:"first_name" json:"first_name"`
	LastFurigana  string   `db:"last_furigana" json:"last_furigana"`
	FirstFurigana string   `db:"first_furigana" json:"first_furigana"`
	Phase         null.Int `db:"phase" json:"phase"`

	// 担当者テーブル
	AgentStaffID uint   `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName    string `db:"staff_name" json:"staff_name"`
	AgentName    string `db:"agent_name" json:"agent_name"`

	// 面談タスクテーブル
	PhaseCategory    null.Int `db:"phase_category" json:"phase_category"`
	PhaseSubCategory null.Int `db:"phase_sub_category" json:"phase_sub_category"`

	// グループ内のタスク
	InterviewTask []InterviewTask `db:"interview_tasks" json:"interview_tasks"`
}

func NewInterviewTaskGroup(
	agentID uint,
	jobSeekerID uint,
	interviewDate time.Time,
	firstInterviewDate time.Time,
) *InterviewTaskGroup {
	return &InterviewTaskGroup{
		AgentID:            agentID,
		JobSeekerID:        jobSeekerID,
		InterviewDate:      interviewDate,
		FirstInterviewDate: firstInterviewDate,
	}
}
