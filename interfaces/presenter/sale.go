package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSaleJSONPresenter(resp responses.Sale) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSaleListJSONPresenter(resp responses.SaleList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSaleListAndMaxPageAndIDListJSONPresenter(resp responses.SaleListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}
