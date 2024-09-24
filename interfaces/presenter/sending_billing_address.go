package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingBillingAddressJSONPresenter(resp responses.SendingBillingAddress) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingBillingAddressListJSONPresenter(resp responses.SendingBillingAddressList) Presenter {
	return NewJSONPresenter(200, resp)
}
