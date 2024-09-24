package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 求職者チャットの最大ページ数を返す（本番実装までは1ページあたり200件）
func getChatGroupSendingJobSeekerListMaxPage(chatGroupSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker) uint {
	var maxPage = len(chatGroupSendingJobSeekerList) / 200

	if 0 < (len(chatGroupSendingJobSeekerList) % 200) {
		maxPage++
	}

	// MaxPageは最大を「1」とする(チャットのページングのみの仕様)
	if maxPage == 0 {
		maxPage = 1
	}

	return uint(maxPage)
}

// 指定ページの求職者チャット一覧を返す（本番実装までは1ページあたり200件）
func getChatGroupSendingJobSeekerListWithPage(chatGroupSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker, page uint) []*entity.ChatGroupWithSendingJobSeeker {
	var (
		perPage uint = 200
		listLen uint = uint(len(chatGroupSendingJobSeekerList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return chatGroupSendingJobSeekerList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.ChatGroupWithSendingJobSeeker{}
	}

	if (listLen - first) <= perPage {
		return chatGroupSendingJobSeekerList[first:]
	}
	return chatGroupSendingJobSeekerList[first:last]
}
