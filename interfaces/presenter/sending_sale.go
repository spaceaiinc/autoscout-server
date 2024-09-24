package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingSaleJSONPresenter(resp responses.SendingSale) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingSaleListJSONPresenter(resp responses.SendingSaleList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingSaleListAndMaxPageAndIDListJSONPresenter(resp responses.SendingSaleListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}
