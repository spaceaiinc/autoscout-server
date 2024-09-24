package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type BillingAddress struct {
	BillingAddress *entity.BillingAddress `json:"billing_address"`
}

func NewBillingAddress(billingAddress *entity.BillingAddress) BillingAddress {
	return BillingAddress{
		BillingAddress: billingAddress,
	}
}

type BillingAddressList struct {
	BillingAddressList []*entity.BillingAddress `json:"billing_address_list"`
}

func NewBillingAddressList(billingAddresses []*entity.BillingAddress) BillingAddressList {
	return BillingAddressList{
		BillingAddressList: billingAddresses,
	}
}

type BillingAddressListAndMaxPageAndIDList struct {
	MaxPageNumber      uint                     `json:"max_page_number"`
	IDList             []uint                   `json:"id_list"`
	BillingAddressList []*entity.BillingAddress `json:"billing_address_list"`
}

func NewBillingAddressListAndMaxPageAndIDList(billingAddresses []*entity.BillingAddress, maxPageNumber uint, idList []uint) BillingAddressListAndMaxPageAndIDList {
	return BillingAddressListAndMaxPageAndIDList{
		MaxPageNumber:      maxPageNumber,
		IDList:             idList,
		BillingAddressList: billingAddresses,
	}
}
