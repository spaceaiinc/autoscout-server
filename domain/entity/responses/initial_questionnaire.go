package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type InitialQuestionnaire struct {
	InitialQuestionnaire *entity.InitialQuestionnaire `json:"initial_questionnaire"`
}

func NewInitialQuestionnaire(initialQuestionnaire *entity.InitialQuestionnaire) InitialQuestionnaire {
	return InitialQuestionnaire{
		InitialQuestionnaire: initialQuestionnaire,
	}
}

type InitialQuestionnaireList struct {
	InitialQuestionnaireList []*entity.InitialQuestionnaire `json:"initial_questionnaire_list"`
}

func NewInitialQuestionnaireList(initialQuestionnaires []*entity.InitialQuestionnaire) InitialQuestionnaireList {
	return InitialQuestionnaireList{
		InitialQuestionnaireList: initialQuestionnaires,
	}
}
