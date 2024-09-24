package entity

import (
	"time"
)

type SendingEnterpriseReferenceMaterial struct {
	ID                  uint      `db:"id" json:"id"`
	SendingEnterpriseID uint      `db:"sending_enterprise_id" json:"sending_enterprise_id"`
	Reference1PDFURL    string    `db:"reference1_pdf_url" json:"reference1_pdf_url"`
	Reference2PDFURL    string    `db:"reference2_pdf_url" json:"reference2_pdf_url"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

func NewSendingEnterpriseReferenceMaterial(
	sendingEnterpriseID uint,
	reference1PDFURL string,
	reference2PDFURL string,
) *SendingEnterpriseReferenceMaterial {
	return &SendingEnterpriseReferenceMaterial{
		SendingEnterpriseID: sendingEnterpriseID,
		Reference1PDFURL:    reference1PDFURL,
		Reference2PDFURL:    reference2PDFURL,
	}
}

type CreateOrUpdateSendingEnterpriseReferenceMaterialParam struct {
	SendingEnterpriseID uint   `db:"sending_enterprise_id" json:"sending_enterprise_id" validate:"required"`
	Reference1PDFURL    string `db:"reference1_pdf_url" json:"reference1_pdf_url"`
	Reference2PDFURL    string `db:"reference2_pdf_url" json:"reference2_pdf_url"`
	ImageURL            string `db:"image_url" json:"image_url"` // sending_enterprise_specialityのimage_url
}

// 送客先エージェントの参考資料タイプ
var SendingEnterpriseFileType = []string{
	"送客先エージェント画像",
	"参考資料1",
	"参考資料2",
}
