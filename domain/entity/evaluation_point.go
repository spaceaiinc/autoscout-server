package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type EvaluationPoint struct {
	ID                     uint      `db:"id" json:"id"`
	TaskID                 uint      `db:"task_id" json:"task_id"`
	SelectionInformationID null.Int  `db:"selection_information_id" json:"selection_information_id"`
	JobSeekerID            uint      `db:"job_seeker_id" json:"job_seeker_id"`
	JobInformationID       uint      `db:"job_information_id" json:"job_information_id"`
	GoodPoint              string    `db:"good_point" json:"good_point"`
	NGPoint                string    `db:"ng_point" json:"ng_point"`
	IsPassed               bool      `db:"is_passed" json:"is_passed"`             // 「true: 合格 or false: 不合格」
	IsReInterview          bool      `db:"is_re_interview" json:"is_re_interview"` // 再面接の有無
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`

	LastName      string `db:"last_name" json:"last_name"`
	FirstName     string `db:"first_name" json:"first_name"`
	LastFurigana  string `db:"last_furigana" json:"last_furigana"`
	FirstFurigana string `db:"first_furigana" json:"first_furigana"`
	CAAgentID     uint   `db:"ca_agent_id" json:"ca_agent_id"`
}

func NewEvaluationPoint(
	taskID uint,
	selectionInformationID null.Int,
	jobSeekerID uint,
	jobInformationID uint,
	goodPoint string,
	ngPoint string,
	isPassed bool,
	isReInterview bool,
) *EvaluationPoint {
	return &EvaluationPoint{
		TaskID:                 taskID,
		SelectionInformationID: selectionInformationID,
		JobSeekerID:            jobSeekerID,
		JobInformationID:       jobInformationID,
		GoodPoint:              goodPoint,
		NGPoint:                ngPoint,
		IsPassed:               isPassed,
		IsReInterview:          isReInterview,
	}
}

type CreateOrUpdateEvaluationPointParam struct {
	TaskID                 uint     `json:"task_id" validate:"required"`
	SelectionInformationID null.Int `json:"selection_information_id" validate:"required"`
	JobSeekerID            uint     `json:"job_seeker_id" validate:"required"`
	JobInformationID       uint     `json:"job_information_id" validate:"required"`
	GoodPoint              string   `json:"good_point"`
	NGPoint                string   `json:"ng_point"`
	IsPassed               bool     `json:"is_passed"`       // 「true: 合格 or false: 不合格」
	IsReInterview          bool     `json:"is_re_interview"` // 再面接の有無
}

type DeleteEvaluationPointParam struct {
	ID uint `json:"id" validate:"required"`
}
