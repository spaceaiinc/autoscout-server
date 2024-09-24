package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ChatGroupWithJobSeeker struct {
	ChatGroupWithJobSeeker *entity.ChatGroupWithJobSeeker `json:"chat_group_with_job_seeker"`
}

func NewChatGroupWithJobSeeker(chatGroupWithJobSeeker *entity.ChatGroupWithJobSeeker) ChatGroupWithJobSeeker {
	return ChatGroupWithJobSeeker{
		ChatGroupWithJobSeeker: chatGroupWithJobSeeker,
	}
}

type ChatGroupWithJobSeekerList struct {
	ChatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker `json:"chat_group_with_job_seeker_list"`
}

func NewChatGroupWithJobSeekerList(chatGroupWithJobSeekers []*entity.ChatGroupWithJobSeeker) ChatGroupWithJobSeekerList {
	return ChatGroupWithJobSeekerList{
		ChatGroupWithJobSeekerList: chatGroupWithJobSeekers,
	}
}

type ChatGroupWithJobSeekerListAndMaxPage struct {
	ChatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker `json:"chat_group_with_job_seeker_list"`
	MaxPageNumber              uint                             `json:"max_page_number"`
}

func NewChatGroupWithJobSeekerListAndMaxPage(chatGroupWithJobSeekers []*entity.ChatGroupWithJobSeeker, maxPageNumber uint) ChatGroupWithJobSeekerListAndMaxPage {
	return ChatGroupWithJobSeekerListAndMaxPage{
		ChatGroupWithJobSeekerList: chatGroupWithJobSeekers,
		MaxPageNumber:              maxPageNumber,
	}
}

type ChatGroupWithJobSeekerUnWatchedCount struct {
	UnWatchedCount uint `json:"unwatched_count"`
}

func NewChatGroupWithJobSeekerUnWatchedCount(unWatchedCount uint) ChatGroupWithJobSeekerUnWatchedCount {
	return ChatGroupWithJobSeekerUnWatchedCount{
		UnWatchedCount: unWatchedCount,
	}
}
