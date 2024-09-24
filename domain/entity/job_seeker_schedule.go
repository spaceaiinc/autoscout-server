package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerSchedule struct {
	ID                uint      `db:"id" json:"id"`
	JobSeekerID       uint      `db:"job_seeker_id" json:"job_seeker_id"`
	TaskID            null.Int  `db:"task_id" json:"task_id"`
	ScheduleType      null.Int  `db:"schedule_type" json:"schedule_type"`
	Title             string    `db:"title" json:"title"`
	StartTime         string    `db:"start_time" json:"start_time"`
	EndTime           string    `db:"end_time" json:"end_time"`
	SeekerDescription string    `db:"seeker_description" json:"seeker_description"`
	StaffDescription  string    `db:"staff_description" json:"staff_description"`
	IsShare           bool      `db:"is_share" json:"is_share"`
	RepetitionCount   null.Int  `db:"repetition_count" json:"repetition_count"`
	CreatedAt         time.Time `db:"created_at" json:"-"`
	UpdatedAt         time.Time `db:"updated_at" json:"-"`

	// タスクテーブル
	TaskGroupID   null.Int `db:"task_group_id" json:"task_group_id"`
	PhaseCategory null.Int `db:"phase_category" json:"phase_category"`

	// リスケ
	RescheduleID null.Int `db:"reschedule_id" json:"reschedule_id"`

	// 求職者テーブル
	LastName      string `db:"last_name" json:"last_name"`
	FirstName     string `db:"first_name" json:"first_name"`
	LastFurigana  string `db:"last_furigana" json:"last_furigana"`
	FirstFurigana string `db:"first_furigana" json:"first_furigana"`

	// 企業テーブル
	CompanyName string `db:"company_name" json:"company_name"`

	// 担当者
	CAStaffID   uint   `db:"ca_staff_id" json:"ca_staff_id"`
	CAStaffName string `db:"ca_staff_name" json:"ca_staff_name"`
	CAAgentName string `db:"ca_agent_name" json:"ca_agent_name"`
	RAStaffName string `db:"ra_staff_name" json:"ra_staff_name"`
	RAAgentName string `db:"ra_agent_name" json:"ra_agent_name"`

	// タスクグループ
	ExternalCompanyName string `db:"external_company_name" json:"external_company_name"`
}

func NewJobSeekerSchedule(
	jobSeekerID uint,
	taskID null.Int,
	scheduleType null.Int,
	title string,
	startTime string,
	endTime string,
	seekerDescription string,
	staffDescription string,
	isShare bool,
	repetitionCount null.Int,
) *JobSeekerSchedule {
	return &JobSeekerSchedule{
		JobSeekerID:       jobSeekerID,
		TaskID:            taskID,
		ScheduleType:      scheduleType,
		Title:             title,
		StartTime:         startTime,
		EndTime:           endTime,
		SeekerDescription: seekerDescription,
		StaffDescription:  staffDescription,
		IsShare:           isShare,
		RepetitionCount:   repetitionCount,
	}
}

type CreateOrUpdateJobSeekerScheduleParam struct {
	JobSeekerID       uint     `json:"job_seeker_id" validate:"required"`
	TaskID            null.Int `json:"task_id"`
	ScheduleType      null.Int `json:"schedule_type" validate:"required"`
	Title             string   `json:"title"`
	StartTime         string   `json:"start_time" validate:"required"`
	EndTime           string   `json:"end_time" validate:"required"`
	SeekerDescription string   `json:"seeker_description"`
	StaffDescription  string   `json:"staff_description"`
	IsShare           bool     `json:"is_share"`
	RepetitionCount   null.Int `db:"repetition_count" json:"repetition_count"`
}

type ShareScheduleParam struct {
	SharedList []JobSeekerSchedule `json:"shared_list" validate:"required"`
}

const (
	ScheduleTypeBySeeker     int64 = iota // 求職者調整（求職者の確保日時）
	ScheduleTypeByCA                      // CA調整（選考の候補日時）
	ScheduleTypeByEnterprise              // 企業調整（選考の確定日時）
	ScheduleTypeByReschedule              // リスケ
)
