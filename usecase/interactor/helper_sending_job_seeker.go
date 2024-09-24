package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 求人の最大ページ数を返す（本番実装までは1ページあたり5件）
func getSendingJobSeekerListMaxPage(sendingJobSeekerList []*entity.SendingJobSeeker) uint {
	var maxPage = len(sendingJobSeekerList) / 20

	if 0 < (len(sendingJobSeekerList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 求人リストからIDリストを取得する
func getSendingJobSeekerIDList(sendingJobSeekerList []*entity.SendingJobSeeker) []uint {
	var idListUint []uint

	if len(sendingJobSeekerList) < 1 {
		return idListUint
	}

	for _, sendingJobSeeker := range sendingJobSeekerList {
		idListUint = append(idListUint, sendingJobSeeker.ID)
	}

	return idListUint
}

// 送客求職者の最大ページ数を返す（本番実装までは1ページあたり20件）
func getSendingManagementListMaxPage(sendingJobSeekerList []*entity.SendingJobSeekerTable) uint {
	var maxPage = len(sendingJobSeekerList) / 20

	if 0 < (len(sendingJobSeekerList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの送客求職者一覧を返す（本番実装までは1ページあたり20件）
func getSendingManagementListWithPage(sendingJobSeekerList []*entity.SendingJobSeekerTable, page uint) []*entity.SendingJobSeekerTable {
	var (
		perPage uint = 20
		listLen uint = uint(len(sendingJobSeekerList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return sendingJobSeekerList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.SendingJobSeekerTable{}
	}

	if (listLen - first) <= perPage {
		return sendingJobSeekerList[first:]
	}
	return sendingJobSeekerList[first:last]
}
