package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobInformationRequiredLanguageType struct {
	ID           uint      `db:"id" json:"id"`
	LanguageID   uint      `db:"language_id" json:"language_id"`
	LanguageType null.Int  `db:"language_type" json:"language_type"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobInformationRequiredLanguageType(
	languageID uint,
	languageType null.Int,
) *SendingJobInformationRequiredLanguageType {
	return &SendingJobInformationRequiredLanguageType{
		LanguageID:   languageID,
		LanguageType: languageType,
	}
}
