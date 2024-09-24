package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSelectionQuestionnaireJSONPresenter(resp responses.SelectionQuestionnaire) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSelectionQuestionnaireListJSONPresenter(resp responses.SelectionQuestionnaireList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSelectionQuestionnaireUUIDJSONPresenter(resp responses.SelectionQuestionnaireUUID) Presenter {
	return NewJSONPresenter(200, resp)
}
