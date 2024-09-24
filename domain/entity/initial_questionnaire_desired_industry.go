package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type InitialQuestionnaireDesiredIndustry struct {
	ID                     uint      `db:"id" json:"id"`
	UUID                   uuid.UUID `db:"uuid" json:"uuid"`
	InitialQuestionnaireID uint      `db:"initial_questionnaire_id" json:"initial_questionnaire_id"`
	DesiredIndustry        null.Int  `db:"desired_industry" json:"desired_industry"`
	DesiredRank            null.Int  `db:"desired_rank" json:"desired_rank"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}

func NewInitialQuestionnaireDesiredIndustry(
	initialQuestionnaireID uint,
	desiredIndustry null.Int,
	desiredRank null.Int,
) *InitialQuestionnaireDesiredIndustry {
	return &InitialQuestionnaireDesiredIndustry{
		InitialQuestionnaireID: initialQuestionnaireID,
		DesiredIndustry:        desiredIndustry,
		DesiredRank:            desiredRank,
	}
}
