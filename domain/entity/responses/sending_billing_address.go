package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SendingBillingAddress struct {
	SendingBillingAddress *entity.SendingBillingAddress `json:"sending_billing_address"`
}

func NewSendingBillingAddress(sendingBillingAddress *entity.SendingBillingAddress) SendingBillingAddress {
	return SendingBillingAddress{
		SendingBillingAddress: sendingBillingAddress,
	}
}

type SendingBillingAddressList struct {
	SendingBillingAddressList []*entity.SendingBillingAddress `json:"sending_billing_address_list"`
}

func NewSendingBillingAddressList(sendingBillingAddresses []*entity.SendingBillingAddress) SendingBillingAddressList {
	return SendingBillingAddressList{
		SendingBillingAddressList: sendingBillingAddresses,
	}
}
