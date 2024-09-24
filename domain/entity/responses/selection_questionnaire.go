package responses

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

type SelectionQuestionnaire struct {
	SelectionQuestionnaire *entity.SelectionQuestionnaire `json:"selection_questionnaire"`
}

func NewSelectionQuestionnaire(selectionQuestionnaire *entity.SelectionQuestionnaire) SelectionQuestionnaire {
	return SelectionQuestionnaire{
		SelectionQuestionnaire: selectionQuestionnaire,
	}
}

type SelectionQuestionnaireList struct {
	SelectionQuestionnaireList []*entity.SelectionQuestionnaire `json:"selection_questionnaire_list"`
}

func NewSelectionQuestionnaireList(selectionQuestionnaires []*entity.SelectionQuestionnaire) SelectionQuestionnaireList {
	return SelectionQuestionnaireList{
		SelectionQuestionnaireList: selectionQuestionnaires,
	}
}

type SelectionQuestionnaireUUID struct {
	SelectionQuestionnaireUUID uuid.UUID `json:"selection_questionnaire_uuid"`
}

func NewSelectionQuestionnaireUUID(selectionQuestionnaireUUID uuid.UUID) SelectionQuestionnaireUUID {
	return SelectionQuestionnaireUUID{
		SelectionQuestionnaireUUID: selectionQuestionnaireUUID,
	}
}
