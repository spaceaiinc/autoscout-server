package entity

import (
	"time"
)

// 送客先エージェントの「エージェントの特徴」を管理するテーブル
type SendingEnterpriseSpeciality struct {
	ID                       uint      `db:"id" json:"id"`
	SendingEnterpriseID      uint      `db:"sending_enterprise_id" json:"sending_enterprise_id"`
	ImageURL                 string    `db:"image_url" json:"image_url"`                                     // エージェント画像
	JobInformationCount      int       `db:"job_information_count" json:"job_information_count"`             // 保有求人数
	SpecializedOccupation    string    `db:"specialized_occupation" json:"specialized_occupation"`           // 得意な職種
	SpecializedIndustry      string    `db:"specialized_industry" json:"specialized_industry"`               // 得意な業種
	SpecializedArea          string    `db:"specialized_area" json:"specialized_area"`                       // 得意エリア
	SpecializedCompanyType   string    `db:"specialized_company_type" json:"specialized_company_type"`       // 得意な企業タイプ
	SpecializedJobSeekerType string    `db:"specialized_job_seeker_type" json:"specialized_job_seeker_type"` // 得意な求職者タイプ
	ConsultingStrengths      string    `db:"consulting_strengths" json:"consulting_strengths"`               // コンサルティングの強み
	SupportContent           string    `db:"support_content" json:"support_content"`                         // サポート内容
	PRPoint                  string    `db:"pr_point" json:"pr_point"`                                       // PRポイント
	CreatedAt                time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time `db:"updated_at" json:"updated_at"`
}

func NewSendingEnterpriseSpeciality(
	sendingEnterpriseID uint,
	imageURL string,
	jobInformationCount int,
	specializedOccupation string,
	specializedIndustry string,
	specializedArea string,
	specializedCompanyType string,
	specializedJobSeekerType string,
	consultingStrengths string,
	supportContent string,
	prPoint string,
) *SendingEnterpriseSpeciality {
	return &SendingEnterpriseSpeciality{
		SendingEnterpriseID:      sendingEnterpriseID,
		ImageURL:                 imageURL,
		JobInformationCount:      jobInformationCount,
		SpecializedOccupation:    specializedOccupation,
		SpecializedIndustry:      specializedIndustry,
		SpecializedArea:          specializedArea,
		SpecializedCompanyType:   specializedCompanyType,
		SpecializedJobSeekerType: specializedJobSeekerType,
		ConsultingStrengths:      consultingStrengths,
		SupportContent:           supportContent,
		PRPoint:                  prPoint,
	}
}

type CreateOrUpdateSendingEnterpriseSpecialityParam struct {
	SendingEnterpriseID      uint   `db:"sending_enterprise_id" json:"sending_enterprise_id" validate:"required"`
	ImageURL                 string `db:"image_url" json:"image_url"`                                     // エージェント画像
	JobInformationCount      int    `db:"job_information_count" json:"job_information_count"`             // 保有求人数
	SpecializedOccupation    string `db:"specialized_occupation" json:"specialized_occupation"`           // 得意な職種
	SpecializedIndustry      string `db:"specialized_industry" json:"specialized_industry"`               // 得意な業種
	SpecializedArea          string `db:"specialized_area" json:"specialized_area"`                       // 得意エリア
	SpecializedCompanyType   string `db:"specialized_company_type" json:"specialized_company_type"`       // 得意な企業タイプ
	SpecializedJobSeekerType string `db:"specialized_job_seeker_type" json:"specialized_job_seeker_type"` // 得意な求職者タイプ
	ConsultingStrengths      string `db:"consulting_strengths" json:"consulting_strengths"`               // コンサルティングの強み
	SupportContent           string `db:"support_content" json:"support_content"`                         // サポート内容
	PRPoint                  string `db:"pr_point" json:"pr_point"`                                       // PRポイント
}
