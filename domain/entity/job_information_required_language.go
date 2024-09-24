package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationRequiredLanguage struct {
	ID            uint      `db:"id" json:"id"`
	ConditionID   uint      `db:"condition_id" json:"condition_id"`
	LanguageLevel null.Int  `db:"language_level" json:"language_level"` // 語学レベル {0:日常会話, 1:ビジネス 99:不問} 追加
	Toeic         null.Int  `db:"toeic" json:"toeic"`
	ToeflIBT      null.Int  `db:"toefl_ibt" json:"toefl_ibt"`
	ToeflPBT      null.Int  `db:"toefl_pbt" json:"toefl_pbt"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"-"`

	// 関連テーブル
	LanguageTypes []JobInformationRequiredLanguageType `db:"language_types" json:"language_types"`
}

func NewJobInformationRequiredLanguage(
	conditionID uint,
	languageLevel null.Int,
	toeic null.Int,
	toeflIBT null.Int,
	toeflPBT null.Int,
) *JobInformationRequiredLanguage {
	return &JobInformationRequiredLanguage{
		ConditionID:   conditionID,
		LanguageLevel: languageLevel,
		Toeic:         toeic,
		ToeflIBT:      toeflIBT,
		ToeflPBT:      toeflPBT,
	}
}
