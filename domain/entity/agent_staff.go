package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type AgentStaff struct {
	ID                              uint      `db:"id" json:"id"`
	UUID                            uuid.UUID `db:"uuid" json:"uuid"`
	AgentID                         uint      `db:"agent_id" json:"agent_id"`
	AgentUUID                       uuid.UUID `db:"agent_uuid" json:"agent_uuid"`
	FirebaseID                      string    `db:"firebase_id" json:"firebase_id"`
	Authority                       null.Int  `db:"authority" json:"authority"`
	StaffName                       string    `db:"staff_name" json:"staff_name"`
	Furigana                        string    `db:"furigana" json:"furigana"`
	Email                           string    `db:"email" json:"email"`
	StaffPhoneNumber                string    `db:"staff_phone_number" json:"staff_phone_number"`
	Department                      string    `db:"department" json:"department"`
	Position                        string    `db:"position" json:"position"`
	Remarks                         string    `db:"remarks" json:"remarks"`
	UsageStatus                     null.Int  `db:"usage_status" json:"usage_status"`
	Notification                    null.Int  `db:"notification" json:"notification"`                       // 未使用
	NotificationJobSeeker           bool      `db:"notification_job_seeker" json:"notification_job_seeker"` // メール通知（求職者）
	NotificationUnwatched           bool      `db:"notification_unwatched" json:"notification_unwatched"`   // メール通知（未処理・未読）
	LastLogin                       time.Time `db:"last_login" json:"last_login"`
	UsageStartDate                  time.Time `db:"usage_start_date" json:"usage_start_date"`
	UsageEndDate                    time.Time `db:"usage_end_date" json:"usage_end_date"`
	IsDeleted                       bool      `db:"is_deleted" json:"is_deleted"` // 論理削除フラグ false: 有効, true: 削除済み
	CreatedAt                       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                       time.Time `db:"updated_at" json:"updated_at"`
	AgentName                       string    `db:"agent_name" json:"agent_name"`
	LineBotID                       string    `db:"line_bot_id" json:"-"`
	LineMessagingChannelSecret      string    `db:"line_messaging_channel_secret" json:"line_messaging_channel_secret"`             // LINE Messaging APIのチャネルシークレット
	LineMessagingChannelAccessToken string    `db:"line_messaging_channel_access_token" json:"line_messaging_channel_access_token"` // LINE Messaging APIのチャネルアクセストークン
	IsCRMActive                     bool      `db:"is_crm_active" json:"is_crm_active"`                                             // CRM機能の有無
	IsAllianceActive                bool      `db:"is_alliance_active" json:"is_alliance_active"`                                   // アライアンス機能の有無
	IsSendingActive                 bool      `db:"is_sending_active" json:"is_sending_active"`                                     // 送客機能の有無
	SendingType                     null.Int  `db:"sending_type" json:"sending_type"`                                               // 送客のタイプ（0: 通常, 1: 送客管理(アンドイーズ仕様)）
}

func NewAgentStaff(
	agentID uint,
	firebaseID string,
	authority null.Int,
	staff_name string,
	furigana string,
	email string,
	staffPhoneNumber string,
	department string,
	position string,
	remarks string,
	usageStatus null.Int,
	notification null.Int,
	notificationJobSeeker bool,
	notificationUnwatched bool,
) *AgentStaff {
	return &AgentStaff{
		AgentID:               agentID,
		FirebaseID:            firebaseID,
		Authority:             authority,
		StaffName:             staff_name,
		Furigana:              furigana,
		Email:                 email,
		StaffPhoneNumber:      staffPhoneNumber,
		Department:            department,
		Position:              position,
		Remarks:               remarks,
		UsageStatus:           usageStatus,
		Notification:          notification,
		NotificationJobSeeker: notificationJobSeeker,
		NotificationUnwatched: notificationUnwatched,
	}
}

type AgentStaffSignUpForAdminParam struct {
	// 担当者情報
	FirebaseID            string    `json:"firebase_id"`
	Authority             null.Int  `json:"authority" validate:"required"`
	StaffName             string    `json:"staff_name" validate:"required"`
	Furigana              string    `json:"furigana" validate:"required"`
	Email                 string    `json:"email" validate:"required"`
	Password              string    `json:"password" validate:"required"`
	StaffPhoneNumber      string    `json:"staff_phone_number" validate:"required"`
	Department            string    `json:"department"`
	Position              string    `json:"position"`
	Remarks               string    `json:"remarks"`
	UsageStatus           null.Int  `json:"usage_status" validate:"required"`
	Notification          null.Int  `json:"notification" validate:"required"`
	NotificationJobSeeker bool      `json:"notification_job_seeker"` // メール通知（求職者）
	NotificationUnwatched bool      `json:"notification_unwatched"`  // メール通知（未処理・未読）
	LastLogin             time.Time `json:"last_login"`
}

type AgentStaffUpdateParam struct {
	Authority             null.Int `json:"authority" validate:"required"`
	StaffName             string   `json:"staff_name" validate:"required"`
	Furigana              string   `json:"furigana" validate:"required"`
	StaffPhoneNumber      string   `json:"staff_phone_number" validate:"required"`
	Department            string   `json:"department"`
	Position              string   `json:"position"`
	Remarks               string   `json:"remarks"`
	UsageStatus           null.Int `json:"usage_status" validate:"required"`
	Notification          null.Int `json:"notification" validate:"required"`
	NotificationJobSeeker bool     `json:"notification_job_seeker"` // メール通知（求職者）
	NotificationUnwatched bool     `json:"notification_unwatched"`  // メール通知（未処理・未読）
}

type AgentStaffEmailUpdateParam struct {
	Email string `json:"email" validate:"required"`
}

type AgentStaffPasswordUpdateParam struct {
	Password string `json:"password" validate:"required"`
}

// 管理者or一般
type AgentStaffAuthority int64

const (
	AuthorityAdmin AgentStaffAuthority = iota
	AuthorityGeneral
)

// 利用可能or利用不可
type AgentStaffUsageStatus int64

const (
	UsageStatusAvailable AgentStaffUsageStatus = iota
	UsageStatusNotAvailable
)

// 担当者削除　*firebaseのアイパスを削除&DBのis_deletedをtrueにする.　求人求職者を引き継ぎ
type DeleteAgentStaffParam struct {
	AgentStaffID uint     `json:"agent_staff_id" validate:"required"`
	NewRAStaffID null.Int `json:"new_ra_staff_id"`
	NewCAStaffID null.Int `json:"new_ca_staff_id"`
}

type UpdateAgentStaffNotificationJobSeekerParam struct {
	AgentStaffID          uint `json:"agent_staff_id" validate:"required"`
	NotificationJobSeeker bool `json:"notification_job_seeker"`
}

type UpdateAgentStaffNotificationUnwatchedParam struct {
	AgentStaffID          uint `json:"agent_staff_id" validate:"required"`
	NotificationUnwatched bool `json:"notification_unwatched"`
}

type UpdateAgentStaffAuthorityParam struct {
	AgentStaffID uint `json:"agent_staff_id" validate:"required"`
	Authority    uint `json:"authority"`
}
