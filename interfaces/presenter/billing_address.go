package presenter

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
)

func NewBillingAddressJSONPresenter(resp responses.BillingAddress) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewBillingAddressListJSONPresenter(resp responses.BillingAddressList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewBillingAddressPageListAndMaxPageAndIDListJSONPresenter(resp responses.BillingAddressListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}
