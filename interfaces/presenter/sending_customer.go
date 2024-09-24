package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingCustomerJSONPresenter(resp responses.SendingCustomer) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingCustomerListJSONPresenter(resp responses.SendingCustomerList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingCustomerListAndMaxPageAndIDListJSONPresenter(resp responses.SendingCustomerListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingCustomerListAndMaxPageAndCountJSONPresenter(resp responses.SendingCustomerListAndMaxPageAndCount) Presenter {
	return NewJSONPresenter(200, resp)
}
