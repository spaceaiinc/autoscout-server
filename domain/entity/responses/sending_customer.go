package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SendingCustomer struct {
	SendingCustomer *entity.SendingCustomer `json:"sending_customer"`
}

func NewSendingCustomer(sendingCustomer *entity.SendingCustomer) SendingCustomer {
	return SendingCustomer{
		SendingCustomer: sendingCustomer,
	}
}

type SendingCustomerList struct {
	SendingCustomerList []*entity.SendingCustomer `json:"sending_customer_list"`
}

func NewSendingCustomerList(sendingCustomers []*entity.SendingCustomer) SendingCustomerList {
	return SendingCustomerList{
		SendingCustomerList: sendingCustomers,
	}
}

type SendingCustomerListAndMaxPageAndIDList struct {
	MaxPageNumber   uint                      `json:"max_page_number"`
	IDList          []uint                    `json:"id_list"`
	SendingCustomer []*entity.SendingCustomer `json:"sending_customer_list"`
}

func NewSendingCustomerListAndMaxPageAndIDList(sendingCustomers []*entity.SendingCustomer, maxPageNumber uint, idList []uint) SendingCustomerListAndMaxPageAndIDList {
	return SendingCustomerListAndMaxPageAndIDList{
		MaxPageNumber:   maxPageNumber,
		IDList:          idList,
		SendingCustomer: sendingCustomers,
	}
}

type SendingCustomerListAndMaxPageAndCount struct {
	MaxPageNumber   uint                      `json:"max_page_number"`
	SendingCustomer []*entity.SendingCustomer `json:"sending_customer_list"`
	AllCount        uint                      `json:"all_count"`
	SendingCount    uint                      `json:"sending_count"`
	CompleteCount   uint                      `json:"complete_count"`
	CloseCount      uint                      `json:"close_count"`
}

func NewSendingCustomerListAndMaxPageAndCount(
	sendingCustomers []*entity.SendingCustomer,
	maxPageNumber uint,
	allCount uint,
	sendingCount uint,
	completeCount uint,
	closeCount uint,
) SendingCustomerListAndMaxPageAndCount {
	return SendingCustomerListAndMaxPageAndCount{
		MaxPageNumber:   maxPageNumber,
		SendingCustomer: sendingCustomers,
		AllCount:        allCount,
		SendingCount:    sendingCount,
		CompleteCount:   completeCount,
		CloseCount:      closeCount,
	}
}
