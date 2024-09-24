package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type InitialQuestionnaireDesiredWorkLocation struct {
	ID                     uint      `db:"id" json:"id"`
	UUID                   uuid.UUID `db:"uuid" json:"uuid"`
	InitialQuestionnaireID uint      `db:"initial_questionnaire_id" json:"initial_questionnaire_id"`
	DesiredWorkLocation    null.Int  `db:"desired_work_location" json:"desired_work_location"`
	DesiredRank            null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}

func NewInitialQuestionnaireDesiredWorkLocation(
	initialQuestionnaireID uint,
	desiredWorkLocation null.Int,
	desiredRank null.Int,
) *InitialQuestionnaireDesiredWorkLocation {
	return &InitialQuestionnaireDesiredWorkLocation{
		InitialQuestionnaireID: initialQuestionnaireID,
		DesiredWorkLocation:    desiredWorkLocation,
		DesiredRank:            desiredRank,
	}
}
