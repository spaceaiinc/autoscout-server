package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobInformationSelectionFlowPattern struct {
	ID               uint      `json:"id" db:"id"`
	JobInformationID uint      `json:"job_information_id" db:"job_information_id"`
	AgentID          uint      `db:"agent_id" json:"agent_id,omitempty"`
	PublicStatus     null.Int  `json:"public_status" db:"public_status"`
	FlowTitle        string    `json:"flow_title" db:"flow_title"`
	FlowPattern      null.Int  `json:"flow_pattern" db:"flow_pattern"`
	IsDeleted        bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`

	// 他テーブル
	SelectionInformations []JobInformationSelectionInformation `json:"selection_informations" db:"selection_informations"`

	SelectionInformationID null.Int `db:"selection_information_id" json:"-"`
	IsQuestionnairy        bool     `db:"is_questionnairy" json:"-"`
	SelectionType          null.Int `db:"selection_type" json:"-"`
}

func NewJobInformationSelectionFlowPattern(
	jobInformationID uint,
	publicStatus null.Int,
	flowTitle string,
	flowPattern null.Int,
	isDeleted bool,
) *JobInformationSelectionFlowPattern {
	return &JobInformationSelectionFlowPattern{
		JobInformationID: jobInformationID,
		PublicStatus:     publicStatus,
		FlowTitle:        flowTitle,
		FlowPattern:      flowPattern,
		IsDeleted:        isDeleted,
	}
}

type CreateAndUpdateSelectionFlowPatternParam struct {
	JobInformationID uint     `json:"job_information_id" validate:"required"`
	PublicStatus     null.Int `json:"public_status"  validate:"required"`
	FlowTitle        string   `json:"flow_title" validate:"required"`
	FlowPattern      null.Int `json:"flow_pattern" validate:"required"`

	SelectionInformations []JobInformationSelectionInformation `json:"selection_informations"`
}
