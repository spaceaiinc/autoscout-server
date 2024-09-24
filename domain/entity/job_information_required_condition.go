package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredCondition struct {
	ID                 uint      `json:"id" db:"id"`
	JobInformationID   uint      `json:"job_information_id" db:"job_information_id"`
	IsCommon           bool      `json:"is_common" db:"is_common"` // 共通条件orパターン条件 {true: 共通条件, false: パターン}
	RequiredManagement null.Int  `json:"required_management" db:"required_management"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`

	// 関連テーブル
	RequiredLicenses               []JobInformationRequiredLicense               `db:"-" json:"required_licenses"`                // 必要資格　複数
	RequiredPCTools                []JobInformationRequiredPCTool                `db:"-" json:"required_pc_tools"`                // 必要PC業務ツール　複数
	RequiredLanguages              JobInformationRequiredLanguage                `db:"-" json:"required_languages"`               // 必要言語 単数
	RequiredExperienceDevelopments []JobInformationRequiredExperienceDevelopment `db:"-" json:"required_experience_developments"` //必要開発経験　言語・OS各1つずつ
	RequiredExperienceJobs         JobInformationRequiredExperienceJob           `db:"-" json:"required_experience_jobs"`         // 必要業職種経験　単数
}

func NewJobInformationRequiredCondition(
	jobInformationID uint,
	isCommon bool,
	requiredManagement null.Int,
) *JobInformationRequiredCondition {
	return &JobInformationRequiredCondition{
		JobInformationID:   jobInformationID,
		IsCommon:           isCommon,
		RequiredManagement: requiredManagement,
	}
}
