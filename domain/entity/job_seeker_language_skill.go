package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerLanguageSkill struct {
	ID                      uint      `db:"id" json:"id"`
	JobSeekerID             uint      `db:"job_seeker_id" json:"job_seeker_id"`
	LanguageType            null.Int  `db:"language_type" json:"language_type"`
	LanguageLevel           null.Int  `db:"language_level" json:"language_level"` // 語学レベル {0:日常会話, 1:ビジネス 99:不問} 追加
	// TalkingSkill            null.Int  `db:"talking_skill" json:"talking_skill"` // 廃止
	// ReadingSkill            null.Int  `db:"reading_skill" json:"reading_skill"` // 廃止
	// WritingSkill            null.Int  `db:"writing_skill" json:"writing_skill"` // 廃止
	Toeic                   null.Int  `db:"toeic" json:"toeic"`
	ToeicExaminationYear    string    `db:"toeic_examination_year" json:"toeic_examination_year"`
	ToeflIBT                null.Int  `db:"toefl_ibt" json:"toefl_ibt"`
	ToeflIBTExaminationYear string    `db:"toefl_ibt_examination_year" json:"toefl_ibt_examination_year"`
	ToeflPBT                null.Int  `db:"toefl_pbt" json:"toefl_pbt"`
	ToeflPBTExaminationYear string    `db:"toefl_pbt_examination_year" json:"toefl_pbt_examination_year"`
	// BusinessExperienceYear  null.Int  `db:"business_experience_year" json:"business_experience_year"` // 廃止 
	// BusinessExperienceMonth null.Int  `db:"business_experience_month" json:"business_experience_month"` // 廃止
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerLanguageSkill(
	jobSeekerID uint,
	languageType null.Int,
	languageLevel null.Int,
	toeic null.Int,
	toeicExaminationYear string,
	toeflIBT null.Int,
	toeflIBTExaminationYear string,
	toeflPBT null.Int,
	toeflPBTExaminationYear string,
) *JobSeekerLanguageSkill {
	return &JobSeekerLanguageSkill{
		JobSeekerID:             jobSeekerID,
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
