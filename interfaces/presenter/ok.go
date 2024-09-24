package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewOKJSONPresenter(resp responses.OK) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewIDListJSONPresenter(resp responses.IDList) Presenter {
	return NewJSONPresenter(200, resp)
}
