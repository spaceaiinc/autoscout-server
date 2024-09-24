package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 企業の最大ページ数を返す（本番実装までは1ページあたり20件）
func getSendingEnterpriseListMaxPage(sendingEnterpriseList []*entity.SendingEnterprise) uint {
	var maxPage = len(sendingEnterpriseList) / 20

	if 0 < (len(sendingEnterpriseList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの企業一覧を返す（本番実装までは1ページあたり20件）
func getSendingEnterpriseListWithPage(sendingEnterpriseList []*entity.SendingEnterprise, page uint) []*entity.SendingEnterprise {
	var (
		perPage uint = 20
		listLen uint = uint(len(sendingEnterpriseList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return sendingEnterpriseList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.SendingEnterprise{}
	}

	if (listLen - first) <= perPage {
		return sendingEnterpriseList[first:]
	}
	return sendingEnterpriseList[first:last]
}
