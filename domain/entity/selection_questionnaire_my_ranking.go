package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SelectionQuestionnaireMyRanking struct {
	ID                       uint      `db:"id" json:"id"`
	SelectionQuestionnaireID uint      `db:"selection_questionnaire_id" json:"selection_questionnaire_id"`
	Rank                     null.Int  `db:"rank" json:"rank"`
	CompanyName              string    `db:"company_name" json:"company_name"`
	Phase                    null.Int  `db:"phase" json:"phase"`
	SelectionDate            string    `db:"selection_date" json:"selection_date"`
	CreatedAt                time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time `db:"updated_at" json:"updated_at"`
}

func NewSelectionQuestionnaireMyRanking(
	selectionQuestionnaireID uint,
	rank null.Int,
	companyName string,
	phase null.Int,
	selectionDate string,
) *SelectionQuestionnaireMyRanking {
	return &SelectionQuestionnaireMyRanking{
		SelectionQuestionnaireID: selectionQuestionnaireID,
		Rank:                     rank,
		CompanyName:              companyName,
		Phase:                    phase,
		SelectionDate:            selectionDate,
	}
}

type CreateOrUpdateSelectionQuestionnaireMyRankingParam struct {
	SelectionQuestionnaireID uint     `db:"selection_questionnaire_id" json:"selection_questionnaire_id"`
	Rank                     null.Int `db:"rank" json:"rank"`
	CompanyName              string   `db:"company_name" json:"company_name"`
	Phase                    null.Int `db:"phase" json:"phase"`
	SelectionDate            string   `db:"selection_date" json:"selection_date"`
}

type DeleteSelectionQuestionnaireMyRankingParam struct {
	ID uint `json:"id" validate:"required"`
}
