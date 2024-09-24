package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID               uint      `db:"id" json:"id"`
	AgentID          uint      `db:"agent_id" json:"agent_id"`
	AgentUUID        uuid.UUID `db:"agent_uuid" json:"agent_uuid"`
	FirebaseID       string    `db:"firebase_id" json:"firebase_id"`
	Name             string    `db:"name" json:"name"`
	AgentName        string    `db:"agent_name" json:"agent_name"` // エージェント名
	LoginTime        time.Time `db:"login_time" json:"login_time"`
	Authority        null.Int  `db:"authority" json:"authority"`
	AllianceAgentIDs []uint    `db:"alliance_agent_ids" json:"alliance_agent_ids"` // アライアンスのエージェントIDリスト
	LastLogin        time.Time `db:"last_login" json:"last_login"`                 // ログイン状況
	UsageStatus      null.Int  `db:"usage_status" json:"usage_status"`             // 利用状況
	IsCRMActive      bool      `db:"is_crm_active" json:"is_crm_active"`           // CRM機能の有無
	IsAllianceActive bool      `db:"is_alliance_active" json:"is_alliance_active"` // アライアンス機能の有無
	IsSendingActive  bool      `db:"is_sending_active" json:"is_sending_active"`   // 送客機能の有無
	SendingType      null.Int  `db:"sending_type" json:"sending_type"`             // 送客のタイプ（0: 通常, 1: 送客管理(アンドイーズ仕様)）
	IsAccountAndes   bool      `db:"is_account_andes" json:"is_account_andes"`     // アンドイーズのアカウントの確認用の値
}

func NewUser(
	id uint,
	agentID uint,
	agentUUID uuid.UUID,
	firebaseID string,
	name string,
	agentName string,
	lastLogin time.Time,
	authority null.Int,
	usageStatus null.Int,
	isCRMActive bool,
	isAllianceActive bool,
	isSendingActive bool,
	sendingType null.Int,
	isAccountAndes bool,
) *User {
	return &User{
		ID:               id,
		AgentID:          agentID,
		AgentUUID:        agentUUID,
		FirebaseID:       firebaseID,
		Name:             name,
		AgentName:        agentName,
		LastLogin:        lastLogin,
		Authority:        authority,
		UsageStatus:      usageStatus,
		IsCRMActive:      isCRMActive,
		IsAllianceActive: isAllianceActive,
		IsSendingActive:  isSendingActive,
		SendingType:      sendingType,
		IsAccountAndes:   isAccountAndes,
	}
}

type GestEnterpriseUser struct {
	JobInformationUUID uuid.UUID `db:"job_information_uuid" json:"job_information_uuid"`
	CompanyName        string    `db:"company_name" json:"company_name"`
}

func NewGestEnterpriseUser(
	jobInformationUUID uuid.UUID,
	companyName string,
) *GestEnterpriseUser {
	return &GestEnterpriseUser{
		JobInformationUUID: jobInformationUUID,
		CompanyName:        companyName,
	}
}

type GestJobSeekerUser struct {
	ID                 uint      `db:"id" json:"id"`
	JobSeekerUUID      uuid.UUID `db:"job_seeker_uuid" json:"job_seeker_uuid"`
	LastName           string    `db:"last_name" json:"last_name"`
	FirstName          string    `db:"first_name" json:"first_name"`
	Email              string    `db:"email" json:"email"`
	AgentID            uint      `db:"agent_id" json:"agent_id"`
	Phase              null.Int  `db:"phase" json:"phase"`
	CanViewMatchingJob bool      `db:"can_view_matching_job" json:"can_view_matching_job"`
}

func NewGestJobSeekerUser(
	id uint,
	jobSeekerUUID uuid.UUID,
	lastName string,
	firstName string,
	email string,
	agentID uint,
	phase null.Int,
	canViewMatchingJob bool,
) *GestJobSeekerUser {
	return &GestJobSeekerUser{
		ID:                 id,
		JobSeekerUUID:      jobSeekerUUID,
		LastName:           lastName,
		FirstName:          firstName,
		Email:              email,
		AgentID:            agentID,
		Phase:              phase,
		CanViewMatchingJob: canViewMatchingJob,
	}
}
