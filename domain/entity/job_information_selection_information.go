package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationSelectionInformation struct {
	ID              uint      `json:"id" db:"id"`
	SelectionFlowID uint      `json:"selection_flow_id" db:"selection_flow_id"`
	SelectionType   null.Int  `json:"selection_type" db:"selection_type"`
	SelectionPoint  string    `json:"selection_point" db:"selection_point"`
	PassedExample   string    `json:"passed_example" db:"passed_example"`
	FailExample     string    `json:"fail_example" db:"fail_example"`
	PassingRate     null.Int  `json:"passing_rate" db:"passing_rate"`
	IsQuestionnairy bool      `json:"is_questionnairy" db:"is_questionnairy"`
	CreatedAt       time.Time `db:"created_at" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"-"`

	// アンケート
	SelectionQuestionnaires []SelectionQuestionnaire `json:"selection_questionnaires" db:"-"`
	IsSelfIntroductionCount uint                     `db:"is_self_introduction_count" json:"is_self_introduction_count"`
	IsSelfPRCount           uint                     `db:"is_self_pr_count" json:"is_self_pr_count"`
	IsRetireReasonCount     uint                     `db:"is_retire_reason_count" json:"is_retire_reason_count"`
	IsJobChangeAxisCount    uint                     `db:"is_job_change_axis_count" json:"is_job_change_axis_count"`
	IsApplyingReasonCount   uint                     `db:"is_applying_reason_count" json:"is_applying_reason_count"`
	IsCareerVisionCount     uint                     `db:"is_career_vision_count" json:"is_career_vision_count"`
	IsReverseQuestionCount  uint                     `db:"is_reverse_question_count" json:"is_reverse_question_count"`

	// 評価点
	EvaluationPoints []EvaluationPoint `json:"evaluation_points" db:"-"`
	PassedExamples   []EvaluationPoint `json:"passed_examples" db:"-"`
	FailureExamples  []EvaluationPoint `json:"failure_examples" db:"-"`
}

func NewJobInformationSelectionInformation(
	selectionFlowID uint,
	selectionType null.Int,
	selectionPoint string,
	passedExample string,
	failExample string,
	passingRate null.Int,
	isQuestionnairy bool,
) *JobInformationSelectionInformation {
	return &JobInformationSelectionInformation{
		SelectionFlowID: selectionFlowID,
		SelectionType:   selectionType,
		SelectionPoint:  selectionPoint,
		PassedExample:   passedExample,
		FailExample:     failExample,
		PassingRate:     passingRate,
		IsQuestionnairy: isQuestionnairy,
	}
}
