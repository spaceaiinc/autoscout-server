package entity

import (
	"time"
)

type EnterpriseReferenceMaterial struct {
	ID               uint      `db:"id" json:"id"`
	EnterpriseID     uint      `db:"enterprise_id" json:"enterprise_id"`
	Reference1PDFURL string    `db:"reference1_pdf_url" json:"reference1_pdf_url"`
	Reference2PDFURL string    `db:"reference2_pdf_url" json:"reference2_pdf_url"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

func NewEnterpriseReferenceMaterial(
	enterpriseID uint,
	reference1PDFURL string,
	reference2PDFURL string,
) *EnterpriseReferenceMaterial {
	return &EnterpriseReferenceMaterial{
		EnterpriseID:     enterpriseID,
		Reference1PDFURL: reference1PDFURL,
		Reference2PDFURL: reference2PDFURL,
	}
}

type CreateOrUpdateEnterpriseReferenceMaterialParam struct {
	EnterpriseID     uint   `db:"enterprise_id" json:"enterprise_id" validate:"required"`
	Reference1PDFURL string `db:"reference1_pdf_url" json:"reference1_pdf_url"`
	Reference2PDFURL string `db:"reference2_pdf_url" json:"reference2_pdf_url"`
}
