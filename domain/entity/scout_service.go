package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type ScoutService struct {
	ID                            uint      `db:"id" json:"id"`
	UUID                          uuid.UUID `db:"uuid" json:"uuid"`
	AgentRobotID                  uint      `db:"agent_robot_id" json:"agent_robot_id"` // エージェントロボットID
	AgentStaffID                  uint      `db:"agent_staff_id" json:"agent_staff_id"` // エージェントスタッフID
	LoginID                       string    `db:"login_id" json:"login_id"`
	Password                      string    `db:"password" json:"password"`
	ServiceType                   null.Int  `db:"service_type" json:"service_type"`                                         // サービスタイプ(0: RAN, 1: マイナビ転職スカウト, 2:AMBI)
	IsActive                      bool      `db:"is_active" json:"is_active"`                                               // アクティブかどうか/false:走らせない true:走る(媒体共通)
	Memo                          string    `db:"memo" json:"memo"`                                                         // メモ
	TemplateTitleForEmployed      string    `db:"template_title_for_employed" json:"template_title_for_employed"`           // 面談調整メールのテンプレート ※就業中
	TemplateTitleForUnemployed    string    `db:"template_title_for_unemployed" json:"template_title_for_unemployed"`       // 面談調整メールのテンプレート ※離職中
	InterviewAdjustmentTemplateID null.Int  `db:"interview_adjustment_template_id" json:"interview_adjustment_template_id"` // 面談調整メールのテンプレートID
	InflowChannelID               null.Int  `db:"inflow_channel_id" json:"inflow_channel_id"`                               // インフローチャンネルID
	LastSendCount              null.Int  `db:"last_send_count" json:"last_send_count"`                             // 最終送信求職者のID
	CreatedAt                     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                     time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	Templates     []ScoutServiceTemplate     `json:"templates"`
	GetEntryTimes []ScoutServiceGetEntryTime `json:"get_entry_times"`

	// DBに存在しない項目
	StaffName string `db:"staff_name" json:"staff_name"`
	RobotName string `db:"robot_name" json:"robot_name"`
}

// ScoutServiceType サービスタイプ
const (
	ScoutServiceTypeRan int64 = iota
	ScoutServiceTypeMynaviScouting
	ScoutServiceTypeAmbi
	ScoutServiceTypeMynaviAgentScout
)

var ScoutServiceTypeLabel = map[int64]string{
	ScoutServiceTypeRan:              "RAN",
	ScoutServiceTypeMynaviScouting:   "マイナビスカウティング",
	ScoutServiceTypeAmbi:             "AMBI",
	ScoutServiceTypeMynaviAgentScout: "マイナビエージェントスカウト",
}

func NewScoutService(
	agentRobotID uint,
	agentStaffID uint,
	loginID string,
	password string,
	serviceType null.Int,
	isActive bool,
	memo string,
	templateTitleForEmployed string,
	templateTitleForUnemployed string,
	interviewAdjustmentTemplateID null.Int,
	inflowChannelID null.Int,
) *ScoutService {
	return &ScoutService{
		AgentRobotID:                  agentRobotID,
		AgentStaffID:                  agentStaffID,
		LoginID:                       loginID,
		Password:                      password,
		ServiceType:                   serviceType,
		IsActive:                      isActive,
		Memo:                          memo,
		TemplateTitleForEmployed:      templateTitleForEmployed,
		TemplateTitleForUnemployed:    templateTitleForUnemployed,
		InterviewAdjustmentTemplateID: interviewAdjustmentTemplateID,
		InflowChannelID:               inflowChannelID,
	}
}

type CreateOrUpdateScoutServiceParam struct {
	AgentRobotID                  uint     `db:"agent_robot_id" json:"agent_robot_id"` // エージェントロボットID
	AgentStaffID                  uint     `db:"agent_staff_id" json:"agent_staff_id"` // エージェントスタッフID
	LoginID                       string   `json:"login_id" validate:"required"`
	Password                      string   `json:"password"`
	ServiceType                   null.Int `json:"service_type" validate:"required"`
	IsActive                      bool     `json:"is_active"`
	Memo                          string   `json:"memo"`
	TemplateTitleForEmployed      string   `json:"template_title_for_employed"`   // 面談調整メールのテンプレート ※就業中
	TemplateTitleForUnemployed    string   `json:"template_title_for_unemployed"` // 面談調整メールのテンプレート ※離職中
	InterviewAdjustmentTemplateID null.Int `json:"interview_adjustment_template_id"`
	InflowChannelID               null.Int `json:"inflow_channel_id"`

	// 関連テーブル
	Templates     []ScoutServiceTemplate     `json:"templates"`
	GetEntryTimes []ScoutServiceGetEntryTime `json:"get_entry_times"`
}

type UpdateScoutServicePasswordParam struct {
	ScoutServiceID uint   `json:"scout_service_id" validate:"required"`
	Password       string `json:"password" validate:"required"`
}

type DeleteScoutServiceParam struct {
	ID uint `json:"id" validate:"required"`
}

type EntryOnScoutServiceParam struct {
	AgentID        uint `json:"agent_id" validate:"required"`
	ScoutServiceID uint `json:"scout_service_id" validate:"required"`
}
