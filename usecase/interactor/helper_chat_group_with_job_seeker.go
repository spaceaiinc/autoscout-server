package interactor

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 求職者チャットの最大ページ数を返す（本番実装までは1ページあたり100件）
func getChatGroupJobSeekerListMaxPage(chatGroupJobSeekerList []*entity.ChatGroupWithJobSeeker) uint {
	var maxPage = len(chatGroupJobSeekerList) / 100

	if 0 < (len(chatGroupJobSeekerList) % 100) {
		maxPage++
	}

	// MaxPageは最大を「1」とする(チャットのページングのみの仕様)
	if maxPage == 0 {
		maxPage = 1
	}

	return uint(maxPage)
}

// 指定ページの求職者チャット一覧を返す（本番実装までは1ページあたり100件）
func getChatGroupJobSeekerListWithPage(chatGroupJobSeekerList []*entity.ChatGroupWithJobSeeker, page uint) []*entity.ChatGroupWithJobSeeker {
	var (
		perPage uint = 100
		listLen uint = uint(len(chatGroupJobSeekerList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return chatGroupJobSeekerList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.ChatGroupWithJobSeeker{}
	}

	if (listLen - first) <= perPage {
		return chatGroupJobSeekerList[first:]
	}
	return chatGroupJobSeekerList[first:last]
}
