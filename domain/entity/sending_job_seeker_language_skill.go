package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerLanguageSkill struct {
	ID                      uint      `db:"id" json:"id"`
	SendingJobSeekerID      uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	LanguageType            null.Int  `db:"language_type" json:"language_type"`
	LanguageLevel           null.Int  `db:"language_level" json:"language_level"` // 語学レベル {0:日常会話, 1:ビジネス 99:不問} 追加
	Toeic                   null.Int  `db:"toeic" json:"toeic"`
	ToeicExaminationYear    string    `db:"toeic_examination_year" json:"toeic_examination_year"`
	ToeflIBT                null.Int  `db:"toefl_ibt" json:"toefl_ibt"`
	ToeflIBTExaminationYear string    `db:"toefl_ibt_examination_year" json:"toefl_ibt_examination_year"`
	ToeflPBT                null.Int  `db:"toefl_pbt" json:"toefl_pbt"`
	ToeflPBTExaminationYear string    `db:"toefl_pbt_examination_year" json:"toefl_pbt_examination_year"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingJobSeekerLanguageSkill(
	sendingJobSeekerID uint,
	languageType null.Int,
	languageLevel null.Int,
	toeic null.Int,
	toeicExaminationYear string,
	toeflIBT null.Int,
	toeflIBTExaminationYear string,
	toeflPBT null.Int,
	toeflPBTExaminationYear string,
) *SendingJobSeekerLanguageSkill {
	return &SendingJobSeekerLanguageSkill{
		SendingJobSeekerID:      sendingJobSeekerID,
		LanguageType:            languageType,
		LanguageLevel:           languageLevel,
		Toeic:                   toeic,
		ToeicExaminationYear:    toeicExaminationYear,
		ToeflIBT:                toeflIBT,
		ToeflIBTExaminationYear: toeflIBTExaminationYear,
		ToeflPBT:                toeflPBT,
		ToeflPBTExaminationYear: toeflPBTExaminationYear,
	}
}
