package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerDocument struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobSeekerID      uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	ResumeOriginURL         string    `db:"resume_origin_url" json:"resume_origin_url"`
	ResumePDFURL            string    `db:"resume_pdf_url" json:"resume_pdf_url"`
	CVOriginURL             string    `db:"cv_origin_url" json:"cv_origin_url"`
	CVPDFURL                string    `db:"cv_pdf_url" json:"cv_pdf_url"`
	RecommendationOriginURL string    `db:"recommendation_origin_url" json:"recommendation_origin_url"`
	RecommendationPDFURL    string    `db:"recommendation_pdf_url" json:"recommendation_pdf_url"`
	IDPhotoURL              string    `db:"id_photo_url" json:"id_photo_url"`
	OtherDocument1URL       string    `db:"other_document1_url" json:"other_document1_url"`
	OtherDocument2URL       string    `db:"other_document2_url" json:"other_document2_url"`
	OtherDocument3URL       string    `db:"other_document3_url" json:"other_document3_url"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`

	// 求職者情報
	AgentID   uint   `db:"agent_id" json:"agent_id"`
	LastName  string `db:"last_name" json:"last_name"`
	FirstName string `db:"first_name" json:"first_name"`
}

func NewSendingJobSeekerDocument(
	sendingJobSeekerID uint,
	resumeOriginURL string,
	resumePDFURL string,
	cvOriginURL string,
	cvPDFURL string,
	recommendOriginURL string,
	recommendPDFURL string,
	idPhotoURL string,
	otherDocument1URL string,
	otherDocument2URL string,
	otherDocument3URL string,
) *SendingJobSeekerDocument {
	return &SendingJobSeekerDocument{
		SendingJobSeekerID:      sendingJobSeekerID,
		ResumeOriginURL:         resumeOriginURL,
		ResumePDFURL:            resumePDFURL,
		CVOriginURL:             cvOriginURL,
		CVPDFURL:                cvPDFURL,
		RecommendationOriginURL: recommendOriginURL,
		RecommendationPDFURL:    recommendPDFURL,
		IDPhotoURL:              idPhotoURL,
		OtherDocument1URL:       otherDocument1URL,
		OtherDocument2URL:       otherDocument2URL,
		OtherDocument3URL:       otherDocument3URL,
	}
}

type CreateOrUpdateSendingJobSeekerDocumentParam struct {
	SendingJobSeekerID      null.Int `db:"sending_job_seeker_id" json:"sending_job_seeker_id" validate:"required"`
	ResumeOriginURL         string   `db:"resume_origin_url" json:"resume_origin_url"`
	ResumePDFURL            string   `db:"resume_pdf_url" json:"resume_pdf_url"`
	CVOriginURL             string   `db:"cv_origin_url" json:"cv_origin_url"`
	CVPDFURL                string   `db:"cv_pdf_url" json:"cv_pdf_url"`
	RecommendationOriginURL string   `db:"recommendation_origin_url" json:"recommendation_origin_url"`
	RecommendationPDFURL    string   `db:"recommendation_pdf_url" json:"recommendation_pdf_url"`
	IDPhotoURL              string   `db:"id_photo_url" json:"id_photo_url"`
	OtherDocument1URL       string   `db:"other_document1_url" json:"other_document1_url"`
	OtherDocument2URL       string   `db:"other_document2_url" json:"other_document2_url"`
	OtherDocument3URL       string   `db:"other_document3_url" json:"other_document3_url"`
}
