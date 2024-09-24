package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredLanguageType struct {
	ID           uint      `db:"id" json:"id"`
	LanguageID   uint      `db:"language_id" json:"language_id"`
	LanguageType null.Int  `db:"language_type" json:"language_type"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

func NewJobInformationRequiredLanguageType(
	languageID uint,
	languageType null.Int,
) *JobInformationRequiredLanguageType {
	return &JobInformationRequiredLanguageType{
		LanguageID:   languageID,
		LanguageType: languageType,
	}
}
