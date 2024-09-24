package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type TaskGroup struct {
	ID                          uint      `db:"id" json:"id"`
	UUID                        uuid.UUID `db:"uuid" json:"uuid"`
	JobSeekerID                 uint      `db:"job_seeker_id" json:"job_seeker_id"`
	JobInformationID            uint      `db:"job_information_id" json:"job_information_id"`
	SelectionFlowPatternID      null.Int  `db:"selection_flow_pattern_id" json:"selection_flow_pattern_id"`
	RALastRequestAt             time.Time `db:"ra_last_request_at" json:"ra_last_request_at"` // RAの最終依頼時間
	RALastWatchedAt             time.Time `db:"ra_last_watched_at" json:"ra_last_watched_at"` // RAの最終閲覧時間
	CALastRequestAt             time.Time `db:"ca_last_request_at" json:"ca_last_request_at"` // CAの最終依頼時間
	CALastWatchedAt             time.Time `db:"ca_last_watched_at" json:"ca_last_watched_at"` // CAの最終閲覧時間
	JoiningDate                 string    `db:"joining_date" json:"joining_date"`             // 入社日
	IsDoubleSided               bool      `db:"is_double_sided" json:"is_double_sided"`
	ExternalJobInformationTitle string    `db:"external_job_information_title" json:"external_job_information_title"`
	ExternalCompanyName         string    `db:"external_company_name" json:"external_company_name"`
	ExternalJobListingURL       string    `db:"external_job_listing_url" json:"external_job_listing_url"`
	IsSelfApplication           bool      `db:"is_self_application" json:"is_self_application"` // 自己応募フラグ（true: 自己応募, false: エージェント応募）
	RAAgentID                   uint      `db:"ra_agent_id" json:"ra_agent_id"`
	CAAgentID                   uint      `db:"ca_agent_id" json:"ca_agent_id"`
	CreatedAt                   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at" json:"updated_at"`

	// 求職者テーブル
	LastName      string `db:"last_name" json:"last_name"`
	FirstName     string `db:"first_name" json:"first_name"`
	LastFurigana  string `db:"last_furigana" json:"last_furigana"`
	FirstFurigana string `db:"first_furigana" json:"first_furigana"`
	Birthday      string `db:"birthday" json:"birthday"`

	// 求人テーブル
	Occupation null.Int `db:"occupation" json:"occupation"`
	Title      string   `db:"title" json:"title"`

	// 企業テーブル
	CompanyName string `db:"company_name" json:"company_name"`

	// 担当者テーブル
	RAStaffID   uint   `db:"ra_staff_id" json:"ra_staff_id"`
	CAStaffID   uint   `db:"ca_staff_id" json:"ca_staff_id"`
	CAStaffName string `db:"ca_staff_name" json:"ca_staff_name"`
	RAStaffName string `db:"ra_staff_name" json:"ra_staff_name"`

	// 面談項目
	InterviewDate time.Time `db:"interview_date" json:"interview_date"`

	PhaseCategory    null.Int `db:"phase_category" json:"phase_category"`
	PhaseSubCategory null.Int `db:"phase_sub_category" json:"phase_sub_category"`

	// 他テーブル
	LatestTask Task `db:"latest_task" json:"latest_task"` // 最新のタスク

	// グループ内のタスク
	Tasks    []Task            `db:"tasks" json:"tasks"`
	Document TaskGroupDocument `db:"document" json:"document"`
}

func NewTaskGroup(
	jobSeekerID uint,
	jobInformationID uint,
	isDoubleSided bool,
	externalJobInformationTitle string,
	externalCompanyName string,
	externalJobListingURL string,
	raAgentID uint,
	caAgentID uint,
) *TaskGroup {
	return &TaskGroup{
		JobSeekerID:                 jobSeekerID,
		JobInformationID:            jobInformationID,
		IsDoubleSided:               isDoubleSided,
		ExternalJobInformationTitle: externalJobInformationTitle,
		ExternalCompanyName:         externalCompanyName,
		ExternalJobListingURL:       externalJobListingURL,
		RAAgentID:                   raAgentID,
		CAAgentID:                   caAgentID,
	}
}

type DeleteTaskGroupParam struct {
	ID uint `json:"id" validate:"required"`
}

type SoundOutGroup struct {
	JobSeekerID        uint      `db:"job_seeker_id" json:"job_seeker_id"`
	JobInformationID   uint      `db:"job_information_id" json:"job_information_id"`
	JobSeekerUUID      uuid.UUID `db:"job_seeker_uuid" json:"job_seeker_uuid"`
	JobInformationUUID uuid.UUID `db:"job_information_uuid" json:"job_information_uuid"`
	RAAgentID          uint      `db:"ra_agent_id" json:"ra_agent_id"`
	RAStaffID          uint      `db:"ra_staff_id" json:"ra_staff_id"`
	RAAgentName        string    `db:"ra_agent_name" json:"ra_agent_name"`
	RAStaffName        string    `db:"ra_staff_name" json:"ra_staff_name"`
	CAAgentID          uint      `db:"ca_agent_id" json:"ca_agent_id"`
	CAStaffID          uint      `db:"ca_staff_id" json:"ca_staff_id"`
	CAAgentName        string    `db:"ca_agent_name" json:"ca_agent_name"`
	CAStaffName        string    `db:"ca_staff_name" json:"ca_staff_name"`
	CompanyName        string    `db:"company_name" json:"company_name"`
	Title              string    `db:"title" json:"title"`
	LastName           string    `db:"last_name" json:"last_name"`
	FirstName          string    `db:"first_name" json:"first_name"`
	LastFurigana       string    `db:"last_furigana" json:"last_furigana"`
	FirstFurigana      string    `db:"first_furigana" json:"first_furigana"`
	SeekerEmail        string    `db:"seeker_email" json:"seeker_email"`
	LineID             string    `db:"line_id" json:"-"`
	LineActive         bool      `db:"line_active" json:"line_active"`
	CAEmail            string    `db:"ca_email" json:"ca_email"`
	CAPhoneNumber      string    `db:"ca_phone_number" json:"ca_phone_number"`
	IsExternal         bool      `db:"is_external" json:"is_external"` // 外部求人フラグ true: 外部求人, false: autoscout内求人
	JobListingURL      string    `db:"job_listing_url" json:"job_listing_url"`
}

func NewSoundOutGroup(
	jobSeekerID uint,
	jobInformationID uint,
	jobSeekerUUID uuid.UUID,
	jobInformationUUID uuid.UUID,
	raAgentID uint,
	raStaffID uint,
	raAgentName string,
	raStaffName string,
	caAgentID uint,
	caStaffID uint,
	caAgentName string,
	caStaffName string,
	companyName string,
	title string,
	lastName string,
	firstName string,
	lastFurigana string,
	firstFurigana string,
	seekerEmail string,
	lineID string,
	lineActive bool,
	caEmail string,
	caPhoneNumber string,
	isExternal bool,
) *SoundOutGroup {
	return &SoundOutGroup{
		JobSeekerID:        jobSeekerID,
		JobInformationID:   jobInformationID,
		JobSeekerUUID:      jobSeekerUUID,
		JobInformationUUID: jobInformationUUID,
		RAAgentID:          raAgentID,
		RAStaffID:          raStaffID,
		RAAgentName:        raAgentName,
		RAStaffName:        raStaffName,
		CAAgentID:          caAgentID,
		CAStaffID:          caStaffID,
		CAAgentName:        caAgentName,
		CAStaffName:        caStaffName,
		CompanyName:        companyName,
		Title:              title,
		LastName:           lastName,
		FirstName:          firstName,
		LastFurigana:       lastFurigana,
		FirstFurigana:      firstFurigana,
		SeekerEmail:        seekerEmail,
		LineID:             lineID,
		LineActive:         lineActive,
		CAEmail:            caEmail,
		CAPhoneNumber:      caPhoneNumber,
		IsExternal:         isExternal,
	}
}

// 外部求人更新用
type ExternalJob struct {
	ExternalJobInformationTitle string `json:"external_job_information_title" validate:"required"`
	ExternalCompanyName         string `json:"external_company_name" validate:"required"`
	ExternalJobListingURL       string `json:"external_job_listing_url" validate:"required"`
}
