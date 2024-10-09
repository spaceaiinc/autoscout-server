package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/google/uuid"
)

type Agent struct {
	ID                              uint      `db:"id" json:"id"`
	UUID                            uuid.UUID `db:"uuid" json:"uuid"`
	AgentName                       string    `db:"agent_name" json:"agent_name"`                                                   // エージェント名
	PermissionCode                  string    `db:"permission_code" json:"permission_code"`                                         // 職業紹介許可番号
	OfficeLocation                  string    `db:"office_location" json:"office_location"`                                         // 住所
	Representative                  string    `db:"representative" json:"representative"`                                           // 代表者
	Establish                       string    `db:"establish" json:"establish"`                                                     // 設立（年月）
	CorporateSiteURL                string    `db:"corporate_site_url" json:"corporate_site_url"`                                   // ホームページ
	PermissionYear                  string    `db:"permission_year" json:"permission_year"`                                         // 紹介事業許可年数
	WorkersCount                    null.Int  `db:"workers_count" json:"workers_count"`                                             // 従業員数
	PhoneNumber                     string    `db:"phone_number" json:"phone_number"`                                               // 電話番号
	InterviewAdjustmentEmail        string    `db:"interview_adjustment_email" json:"interview_adjustment_email"`                   // メールアドレス（面談調整用）
	AgreementFileURL                string    `db:"agreement_file_url" json:"agreement_file_url"`                                   // 同意書ファイルのURL
	LineBotID                       string    `db:"line_bot_id" json:"-"`                                                           // LINE BotアカウントのID
	LineMessagingChannelSecret      string    `db:"line_messaging_channel_secret" json:"line_messaging_channel_secret"`             // LINE Messaging APIのチャネルシークレット
	LineMessagingChannelAccessToken string    `db:"line_messaging_channel_access_token" json:"line_messaging_channel_access_token"` // LINE Messaging APIのチャネルアクセストークン
	LineLoginChannelID              string    `db:"line_loging_channel_id" json:"line_loging_channel_id"`                           // LINE LoginのチャネルID
	LineLoginChannelSecret          string    `db:"line_login_channel_secret" json:"line_login_channel_secret"`                     // LINE Loginのチャネルシークレット
	SendingAgreementFileURL         string    `db:"sending_agreement_file_url" json:"sending_agreement_file_url"`                   // 送客用同意書ファイルのURL
	IsCRMActive                     bool      `db:"is_crm_active" json:"is_crm_active"`                                             // CRM機能の有無
	IsAllianceActive                bool      `db:"is_alliance_active" json:"is_alliance_active"`                                   // アライアンス機能の有無
	IsSendingActive                 bool      `db:"is_sending_active" json:"is_sending_active"`                                     // 送客機能の有無
	SendingType                     null.Int  `db:"sending_type" json:"sending_type"`                                               // 送客のタイプ（0: 通常, 1: 送客管理(アンドイーズ仕様)）
	CreatedAt                       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                       time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	AgentStaffList    []AgentStaff `db:"agent_staff_list" json:"agent_staff"`
	ClaimAccountCount int          `db:"claim_account_count" json:"claim_account_count"`
}

func NewAgent(
	agentName string,
	officeLocation string,
	representative string,
	establish string,
	corporateSiteURL string,
	permissionCode string,
	permissionYear string,
	workersCount null.Int,
	phoneNumber string,
	interviewAdjustmentEmail string,
	agreementFileURL string,
	lineBotID string,
	lineMessagingChannelSecret string,
	lineMessagingChannelAccessToken string,
	lineLoginChannelID string,
	lineLoginChannelSecret string,
	sendingAgreementFileURL string,
	isCRMActive bool,
	isAllianceActive bool,
	isSendingActive bool,
	sendingType null.Int,
) *Agent {
	return &Agent{
		AgentName:                       agentName,
		OfficeLocation:                  officeLocation,
		Representative:                  representative,
		Establish:                       establish,
		CorporateSiteURL:                corporateSiteURL,
		PermissionCode:                  permissionCode,
		PermissionYear:                  permissionYear,
		WorkersCount:                    workersCount,
		PhoneNumber:                     phoneNumber,
		InterviewAdjustmentEmail:        interviewAdjustmentEmail,
		AgreementFileURL:                agreementFileURL,
		LineBotID:                       lineBotID,
		LineMessagingChannelSecret:      lineMessagingChannelSecret,
		LineMessagingChannelAccessToken: lineMessagingChannelAccessToken,
		LineLoginChannelID:              lineLoginChannelID,
		LineLoginChannelSecret:          lineLoginChannelSecret,
		// IsTrial:                         isTrial,
		// TrialEndTime:                    trialEndTime,
		SendingAgreementFileURL: sendingAgreementFileURL,
		IsCRMActive:             isCRMActive,
		IsAllianceActive:        isAllianceActive,
		IsSendingActive:         isSendingActive,
		SendingType:             sendingType,
	}
}

type CreateOrUpdateAgentParam struct {
	AgentName                       string   `db:"agent_name" json:"agent_name" validate:"required"`                               // エージェント名
	OfficeLocation                  string   `db:"office_location" json:"office_location"`                                         // 住所
	Representative                  string   `db:"representative" json:"representative"`                                           // 代表者
	Establish                       string   `db:"establish" json:"establish"`                                                     // 設立（年月）
	CorporateSiteURL                string   `db:"corporate_site_url" json:"corporate_site_url"`                                   // ホームページ
	PermissionCode                  string   `db:"permission_code" json:"permission_code" validate:"required"`                     // 職業紹介許可番号
	PermissionYear                  string   `db:"permission_year" json:"permission_year"`                                         // 紹介事業許可年数
	WorkersCount                    null.Int `db:"workers_count" json:"workers_count"`                                             // 職業紹介の業務に従事する者の数
	PhoneNumber                     string   `db:"phone_number" json:"phone_number"`                                               // 連絡先
	InterviewAdjustmentEmail        string   `db:"interview_adjustment_email" json:"interview_adjustment_email"`                   // メールアドレス（面談調整用）
	AgreementFileURL                string   `db:"agreement_file_url" json:"agreement_file_url"`                                   // 同意書ファイルのURL
	LineBotID                       string   `db:"line_bot_id" json:"-"`                                                           // LINE BotアカウントのID
	LineMessagingChannelSecret      string   `db:"line_messaging_channel_secret" json:"line_messaging_channel_secret"`             // LINE Messaging APIのチャネルシークレット
	LineMessagingChannelAccessToken string   `db:"line_messaging_channel_access_token" json:"line_messaging_channel_access_token"` // LINE Messaging APIのチャネルアクセストークン
	LineLoginChannelID              string   `db:"line_loging_channel_id" json:"line_loging_channel_id"`                           // LINE LoginのチャネルID
	LineLoginChannelSecret          string   `db:"line_login_channel_secret" json:"line_login_channel_secret"`                     // LINE Loginのチャネルシークレット
	SendingAgreementFileURL         string   `db:"sending_agreement_file_url" json:"sending_agreement_file_url"`                   // 送客用同意書ファイルのURL
	IsCRMActive                     bool     `db:"is_crm_active" json:"is_crm_active"`                                             // CRM機能の有無
	IsAllianceActive                bool     `db:"is_alliance_active" json:"is_alliance_active"`                                   // アライアンス機能の有無
	IsSendingActive                 bool     `db:"is_sending_active" json:"is_sending_active"`                                     // 送客機能の有無
	SendingType                     null.Int `db:"sending_type" json:"sending_type"`                                               // 送客のタイプ（0: 通常, 1: 送客管理(アンドイーズ仕様)）

	// 関連テーブル
}

// 同意書ファイルのURLを更新
type AgentAgreementFileURLParam struct {
	AgentID                 uint   `db:"agent_id" json:"agent_id" validate:"required"` // エージェントID
	AgreementFileURL        string `json:"agreement_file_url"`                         // 同意書ファイルのURL
	SendingAgreementFileURL string `json:"sending_agreement_file_url"`                 // 送客用同意書ファイルのURL
}

// CRM機能の有無を更新
type AgentForAdminParam struct {
	AgentID          uint     `db:"agent_id" json:"agent_id" validate:"required"` // エージェントID
	IsCRMActive      bool     `db:"is_crm_active" json:"is_crm_active"`           // CRM機能の有無
	IsAllianceActive bool     `db:"is_alliance_active" json:"is_alliance_active"` // アライアンス機能の有無
	IsSendingActive  bool     `db:"is_sending_active" json:"is_sending_active"`   // 送客機能の有無
	SendingType      null.Int `db:"sending_type" json:"sending_type"`             // 送客のタイプ（0: 通常, 1: 送客管理(アンドイーズ仕様)）
}

// エージェントのLINEアカウントと連携
type AgentLineChannelParam struct {
	AgentID                         uint      `db:"agent_id" json:"agent_id" validate:"required"` // エージェントID
	AgentUUID                       uuid.UUID `db:"agent_uuid" json:"agent_uuid"`
	LineBotID                       string    `db:"line_bot_id" json:"-"`                                                                               // LINE BotアカウントのID
	LineMessagingChannelSecret      string    `db:"line_messaging_channel_secret" json:"line_messaging_channel_secret" validate:"required"`             // LINE Messaging APIのチャネルシークレット
	LineMessagingChannelAccessToken string    `db:"line_messaging_channel_access_token" json:"line_messaging_channel_access_token" validate:"required"` // LINE Messaging APIのチャネルアクセストークン
	LineLoginChannelID              string    `db:"line_loging_channel_id" json:"line_loging_channel_id" validate:"required"`                           // LINE LoginのチャネルID
	LineLoginChannelSecret          string    `db:"line_login_channel_secret" json:"line_login_channel_secret" validate:"required"`                     // LINE Loginのチャネルシークレット
}

// リクエストボディを受け取る構造体
type Response struct {
	Body string `json:"RequestBody"`
}

// リクエストボディから特定のパラメータを受け取る構造体
type Event struct {
	Events []struct {
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
		Source struct {
			UserID string `json:"userID"`
		} `json:"source"`
	}
}

type Bot struct {
	BotID       string `json:"userID"`
	DisplayName string `json:"displayName"`
	ChatMode    string `json:"chatMode"`
}

// エージェントのLINEアカウントと連携
type LineTokenParam struct {
	LineChannelAccessToken string `db:"line_channel_access_token" json:"line_channel_access_token" validate:"required"` // LINEチャネルトークン
}

type CreateAgentAndAgentStaffParam struct {
	// エージェント情報
	AgentName                string   `db:"agent_name" json:"agent_name" validate:"required"`             // エージェント名
	PermissionCode           string   `db:"permission_code" json:"permission_code" validate:"required"`   // 職業紹介許可番号
	OfficeLocation           string   `db:"office_location" json:"office_location" validate:"required"`   // 住所
	Representative           string   `db:"representative" json:"representative" validate:"required"`     // 代表者
	Establish                string   `db:"establish" json:"establish"`                                   // 設立（年月）
	CorporateSiteURL         string   `db:"corporate_site_url" json:"corporate_site_url"`                 // ホームページ
	PermissionYear           string   `db:"permission_year" json:"permission_year"`                       // 紹介事業許可年数
	WorkersCount             null.Int `db:"workers_count" json:"workers_count"`                           // 従業員数
	PhoneNumber              string   `db:"phone_number" json:"phone_number"`                             // 連絡先
	InterviewAdjustmentEmail string   `db:"interview_adjustment_email" json:"interview_adjustment_email"` // メールアドレス（面談調整用）
	AgreementFileURL         string   `db:"agreement_file_url" json:"agreement_file_url"`                 // 連絡先
	IsCRMActive              bool     `db:"is_crm_active" json:"is_crm_active"`                           // CRM機能の有無
	IsAllianceActive         bool     `db:"is_alliance_active" json:"is_alliance_active"`                 // アライアンス機能の有無
	IsSendingActive          bool     `db:"is_sending_active" json:"is_sending_active"`                   // 送客機能の有無
	SendingType              null.Int `db:"sending_type" json:"sending_type"`                             // 送客のタイプ（0: 通常, 1: 送客管理(アンドイーズ仕様)）

	// 担当者情報
	Email                 string   `json:"email" validate:"required"`                      // メールアドレス
	Password              string   `json:"password" validate:"required"`                   // パスワード
	Authority             null.Int `db:"authority" json:"authority"`                       // 権限
	StaffName             string   `db:"staff_name" json:"staff_name" validate:"required"` // 担当者名
	Furigana              string   `db:"furigana" json:"furigana"`                         // 担当者名フリガナ
	StaffPhoneNumber      string   `db:"staff_phone_number" json:"staff_phone_number"`     // 電話番号
	Department            string   `db:"department" json:"department"`                     // 部署
	Position              string   `db:"position" json:"position"`                         // 役職
	Remarks               string   `db:"remarks" json:"remarks"`                           // 備考
	UsageStatus           null.Int `db:"usage_status" json:"usage_status"`                 // 使用状況
	Notification          null.Int `db:"notification" json:"notification"`                 // 通知
	NotificationJobSeeker bool     `json:"notification_job_seeker"`                        // 求職者通知設定
	NotificationUnwatched bool     `json:"notification_unwatched"`                         // 未読通知設定
}
