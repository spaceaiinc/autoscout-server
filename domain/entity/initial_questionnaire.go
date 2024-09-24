package entity

import (
	"time"

	"github.com/google/uuid"
)

type InitialQuestionnaire struct {
	ID          uint      `db:"id" json:"id"`
	UUID        uuid.UUID `db:"uuid" json:"uuid"`
	JobSeekerID uint      `db:"job_seeker_id" json:"job_seeker_id"`
	Question    string    `db:"question" json:"question"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	DesiredIndustries    []InitialQuestionnaireDesiredIndustry     `db:"desired_industries" json:"desired_industries"`
	DesiredOccupations   []InitialQuestionnaireDesiredOccupation   `db:"desired_occupations" json:"desired_occupations"`
	DesiredWorkLocations []InitialQuestionnaireDesiredWorkLocation `db:"desired_prefectures" json:"desired_work_locations"`
}

func NewInitialQuestionnaire(
	jobSeekerID uint,
	question string,
) *InitialQuestionnaire {
	return &InitialQuestionnaire{
		JobSeekerID: jobSeekerID,
		Question:    question,
	}
}

// 求職者の面談前アンケート登録
// 1. 求職者の同意項目アップデート
// 2. アンケート情報の登録（業界、職種、勤務地、質問要望）
// 3. 求職者情報の更新（業界、職種、勤務地）
// 4. 求職者情報の更新（ファイル（履歴書（原本）、商務経歴書（原本）））
type CreateInitialQuestionnaireParam struct {
	JobSeekerID uint   `db:"job_seeker_id" json:"job_seeker_id" validate:"required"`
	Question    string `db:"question" json:"question"`

	// 関連テーブル
	DesiredIndustries    []InitialQuestionnaireDesiredIndustry     `json:"desired_industries"`
	DesiredOccupations   []InitialQuestionnaireDesiredOccupation   `json:"desired_occupations"`
	DesiredWorkLocations []InitialQuestionnaireDesiredWorkLocation `json:"desired_work_locations"`

	// 求職者テーブル
	Agreement bool              `db:"agreement" json:"agreement"` // 求職者の同意
	Documents JobSeekerDocument `json:"documents"`                // 求職者の資料
}
