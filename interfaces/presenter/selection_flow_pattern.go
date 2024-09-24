package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSelectionFlowPatternJSONPresenter(resp responses.SelectionFlowPattern) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSelectionFlowPatternListJSONPresenter(resp responses.SelectionFlowPatternList) Presenter {
	return NewJSONPresenter(200, resp)
}
