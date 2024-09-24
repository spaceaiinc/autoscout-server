package entity

import (
	"time"
)

type TaskIsRecommendDocument struct {
	ID                                 uint      `db:"id" json:"id"`
	TaskID                             uint      `db:"task_id" json:"task_id"`
	IsRecommendUploadResume            bool      `db:"is_recommend_upload_resume" json:"is_recommend_upload_resume"`
	IsRecommendUploadCV                bool      `db:"is_recommend_upload_cv" json:"is_recommend_upload_cv"`
	IsRecommendUploadRecommendation    bool      `db:"is_recommend_upload_recommendation" json:"is_recommend_upload_recommendation"`
	IsRecommendGeneratedResume         bool      `db:"is_recommend_generated_resume" json:"is_recommend_generated_resume"`
	IsRecommendGeneratedCV             bool      `db:"is_recommend_generated_cv" json:"is_recommend_generated_cv"`
	IsRecommendGeneratedRecommendation bool      `db:"is_recommend_generated_recommendation" json:"is_recommend_generated_recommendation"`
	IsRecommendGeneratedMaskResume     bool      `db:"is_recommend_generated_mask_resume" json:"is_recommend_generated_mask_resume"`
	CreatedAt                          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                          time.Time `db:"updated_at" json:"updated_at"`
}

func NewTaskIsRecommendDocument(
	taskID uint,
	isRecommendUploadResume bool,
	isRecommendUploadCV bool,
	isRecommendUploadRecommendation bool,
	isRecommendGeneratedResume bool,
	isRecommendGeneratedCV bool,
	isRecommendGeneratedRecommendation bool,
	isRecommendGeneratedMaskResume bool,
) *TaskIsRecommendDocument {
	return &TaskIsRecommendDocument{
		TaskID:                             taskID,
		IsRecommendUploadResume:            isRecommendUploadResume,
		IsRecommendUploadCV:                isRecommendUploadCV,
		IsRecommendUploadRecommendation:    isRecommendUploadRecommendation,
		IsRecommendGeneratedResume:         isRecommendGeneratedResume,
		IsRecommendGeneratedCV:             isRecommendGeneratedCV,
		IsRecommendGeneratedRecommendation: isRecommendGeneratedRecommendation,
		IsRecommendGeneratedMaskResume:     isRecommendGeneratedMaskResume,
	}
}
