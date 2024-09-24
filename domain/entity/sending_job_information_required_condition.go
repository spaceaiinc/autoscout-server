package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredCondition struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobInformationID uint      `db:"sending_job_information_id" json:"sending_job_information_id"`
	IsCommon                bool      `db:"is_common" json:"is_common"` // 共通条件orパターン条件 {true: 共通条件, false: パターン}
	RequiredManagement      null.Int  `db:"required_management" json:"required_management"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`

	// 関連テーブル
	RequiredLicenses               []SendingJobInformationRequiredLicense               `db:"-" json:"required_licenses"`                // 必要資格　複数
	RequiredPCTools                []SendingJobInformationRequiredPCTool                `db:"-" json:"required_pc_tools"`                // 必要PC業務ツール　複数
	RequiredExperienceDevelopments []SendingJobInformationRequiredExperienceDevelopment `db:"-" json:"required_experience_developments"` //必要開発経験　言語・OS各1つずつ
	RequiredLanguages              SendingJobInformationRequiredLanguage                `db:"-" json:"required_languages"`               // 必要言語 単数
	RequiredExperienceJobs         SendingJobInformationRequiredExperienceJob           `db:"-" json:"required_experience_jobs"`         // 必要業職種経験　単数

}

func NewSendingJobInformationRequiredCondition(
	sendingJobInformationID uint,
	isCommon bool,
	requiredManagement null.Int,
) *SendingJobInformationRequiredCondition {
	return &SendingJobInformationRequiredCondition{
		SendingJobInformationID: sendingJobInformationID,
		IsCommon:                isCommon,
		RequiredManagement:      requiredManagement,
	}
}
