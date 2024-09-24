package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type Sale struct {
	Sale *entity.Sale `json:"sale"`
}

func NewSale(sale *entity.Sale) Sale {
	return Sale{
		Sale: sale,
	}
}

type SaleList struct {
	SaleList []*entity.Sale `json:"sale_list"`
}

func NewSaleList(sales []*entity.Sale) SaleList {
	return SaleList{
		SaleList: sales,
	}
}

type SaleListAndMaxPageAndIDList struct {
	MaxPageNumber uint           `json:"max_page_number"`
	IDList        []uint         `json:"id_list"`
	SaleList      []*entity.Sale `json:"sale_list"`
}

func NewSaleListAndMaxPageAndIDList(sales []*entity.Sale, maxPageNumber uint, idList []uint) SaleListAndMaxPageAndIDList {
	return SaleListAndMaxPageAndIDList{
		MaxPageNumber: maxPageNumber,
		IDList:        idList,
		SaleList:      sales,
	}
}
