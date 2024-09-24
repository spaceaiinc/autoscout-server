package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SendingSale struct {
	SendingSale *entity.SendingSale `json:"sending_sale"`
}

func NewSendingSale(sendingSale *entity.SendingSale) SendingSale {
	return SendingSale{
		SendingSale: sendingSale,
	}
}

type SendingSaleList struct {
	SendingSaleList []*entity.SendingSale `json:"sending_sale_list"`
}

func NewSendingSaleList(sendingSales []*entity.SendingSale) SendingSaleList {
	return SendingSaleList{
		SendingSaleList: sendingSales,
	}
}

type SendingSaleListAndMaxPageAndIDList struct {
	MaxPageNumber   uint                  `json:"max_page_number"`
	IDList          []uint                `json:"id_list"`
	SendingSaleList []*entity.SendingSale `json:"sending_sale_list"`
}

func NewSendingSaleListAndMaxPageAndIDList(sendingSales []*entity.SendingSale, maxPageNumber uint, idList []uint) SendingSaleListAndMaxPageAndIDList {
	return SendingSaleListAndMaxPageAndIDList{
		MaxPageNumber:   maxPageNumber,
		IDList:          idList,
		SendingSaleList: sendingSales,
	}
}
