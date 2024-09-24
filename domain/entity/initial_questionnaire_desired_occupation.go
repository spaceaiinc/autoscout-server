package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type InitialQuestionnaireDesiredOccupation struct {
	ID                     uint      `db:"id" json:"id"`
	UUID                   uuid.UUID `db:"uuid" json:"uuid"`
	InitialQuestionnaireID uint      `db:"initial_questionnaire_id" json:"initial_questionnaire_id"`
	DesiredOccupation      null.Int  `db:"desired_occupation" json:"desired_occupation"`
	DesiredRank            null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}

func NewInitialQuestionnaireDesiredOccupation(
	initialQuestionnaireID uint,
	desiredOccupation null.Int,
	desiredRank null.Int,
) *InitialQuestionnaireDesiredOccupation {
	return &InitialQuestionnaireDesiredOccupation{
		InitialQuestionnaireID: initialQuestionnaireID,
		DesiredOccupation:      desiredOccupation,
		DesiredRank:            desiredRank,
	}
}
