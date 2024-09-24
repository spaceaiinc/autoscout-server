package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 企業の最大ページ数を返す（本番実装までは1ページあたり20件）
func getEnterpriseListMaxPage(enterpriseList []*entity.EnterpriseProfile) uint {
	var maxPage = len(enterpriseList) / 20

	if 0 < (len(enterpriseList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの企業一覧を返す（本番実装までは1ページあたり20件）
func getEnterpriseListWithPage(enterpriseList []*entity.EnterpriseProfile, page uint) []*entity.EnterpriseProfile {
	var (
		perPage uint = 20
		listLen uint = uint(len(enterpriseList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return enterpriseList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.EnterpriseProfile{}
	}

	if (listLen - first) <= perPage {
		return enterpriseList[first:]
	}
	return enterpriseList[first:last]
}
