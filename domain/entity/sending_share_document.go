package entity

import "time"

type SendingShareDocument struct {
	ID                             uint      `db:"id" json:"id"`
	SendingJobSeekerID             uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`                         // 送客求職者のID
	SendingEnterpriseID            uint      `db:"sending_enterprise_id" json:"sending_enterprise_id"`                         // 送客先エージェントのID
	IsShareUploadResume            bool      `db:"is_share_upload_resume" json:"is_share_upload_resume"`                       // アップロードした履歴書のシェア
	IsShareUploadCV                bool      `db:"is_share_upload_cv" json:"is_share_upload_cv"`                               // アップロードした職務経歴書のシェア
	IsShareUploadRecommendation    bool      `db:"is_share_upload_recommendation" json:"is_share_upload_recommendation"`       // アップロードした推薦状のシェア
	IsShareGeneratedResume         bool      `db:"is_share_generated_resume" json:"is_share_generated_resume"`                 // 自動生成した履歴書のシェア
	IsShareGeneratedCV             bool      `db:"is_share_generated_cv" json:"is_share_generated_cv"`                         // 自動生成した職務経歴書のシェア
	IsShareGeneratedRecommendation bool      `db:"is_share_generated_recommendation" json:"is_share_generated_recommendation"` // 自動生成した推薦状のシェア
	CreatedAt                      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                      time.Time `db:"updated_at" json:"updated_at"`
}

func NewSendingShareDocument(
	sendingJobSeekerID uint,
	sendingEnterpriseID uint,
	isShareUploadResume bool,
	isShareUploadCV bool,
	isShareUploadRecommendation bool,
	isShareGeneratedResume bool,
	isShareGeneratedCV bool,
	isShareGeneratedRecommendation bool,
) *SendingShareDocument {
	return &SendingShareDocument{
		SendingJobSeekerID:             sendingJobSeekerID,
		SendingEnterpriseID:            sendingEnterpriseID,
		IsShareUploadResume:            isShareUploadResume,
		IsShareUploadCV:                isShareUploadCV,
		IsShareUploadRecommendation:    isShareUploadRecommendation,
		IsShareGeneratedResume:         isShareGeneratedResume,
		IsShareGeneratedCV:             isShareGeneratedCV,
		IsShareGeneratedRecommendation: isShareGeneratedRecommendation,
	}
}
