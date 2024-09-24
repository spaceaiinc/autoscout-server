package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 企業の最大ページ数を返す（本番実装までは1ページあたり20件）
func getBillingAddressListMaxPage(billingAddressList []*entity.BillingAddress) uint {
	var maxPage = len(billingAddressList) / 20

	if 0 < (len(billingAddressList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの企業一覧を返す（本番実装までは1ページあたり20件）
func getBillingAddressListWithPage(billingAddressList []*entity.BillingAddress, page uint) []*entity.BillingAddress {
	var (
		perPage uint = 20
		listLen uint = uint(len(billingAddressList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return billingAddressList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.BillingAddress{}
	}

	if (listLen - first) <= perPage {
		return billingAddressList[first:]
	}
	return billingAddressList[first:last]
}
