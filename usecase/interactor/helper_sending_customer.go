package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 求人の最大ページ数を返す（本番実装までは1ページあたり5件）
func getSendingCustomerListMaxPage(sendingCustomerList []*entity.SendingCustomer) uint {
	var maxPage = len(sendingCustomerList) / 20

	if 0 < (len(sendingCustomerList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの求人一覧を返す（本番実装までは1ページあたり5件）
func getSendingCustomerListWithPage(sendingCustomerList []*entity.SendingCustomer, page uint) []*entity.SendingCustomer {
	var (
		perPage uint = 20
		listLen uint = uint(len(sendingCustomerList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return sendingCustomerList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.SendingCustomer{}
	}

	if (listLen - first) <= perPage {
		return sendingCustomerList[first:]
	}
	return sendingCustomerList[first:last]
}
