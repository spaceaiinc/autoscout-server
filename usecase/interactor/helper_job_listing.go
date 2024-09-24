package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 求人の最大ページ数を返す（本番実装までは1ページあたり20件）
func getJobListingListMaxPage(jobListingList []*entity.JobListing) uint {
	var maxPage = len(jobListingList) / 20
	if 0 < (len(jobListingList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの求人一覧を返す（本番実装までは1ページあたり20件）
func getJobListingListWithPage(jobListingList []*entity.JobListing, page uint) []*entity.JobListing {
	var (
		perPage uint = 20
		listLen uint = uint(len(jobListingList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return jobListingList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.JobListing{}
	}

	if (listLen - first) <= perPage {
		return jobListingList[first:]
	}
	return jobListingList[first:last]
}
