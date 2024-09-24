package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ChatGroupWithSendingJobSeeker struct {
	ChatGroupWithSendingJobSeeker *entity.ChatGroupWithSendingJobSeeker `json:"chat_group_with_sending_job_seeker"`
}

func NewChatGroupWithSendingJobSeeker(chatGroupWithSendingJobSeeker *entity.ChatGroupWithSendingJobSeeker) ChatGroupWithSendingJobSeeker {
	return ChatGroupWithSendingJobSeeker{
		ChatGroupWithSendingJobSeeker: chatGroupWithSendingJobSeeker,
	}
}

type ChatGroupWithSendingJobSeekerList struct {
	ChatGroupWithSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker `json:"chat_group_with_sending_job_seeker_list"`
}

func NewChatGroupWithSendingJobSeekerList(chatGroupWithSendingJobSeekers []*entity.ChatGroupWithSendingJobSeeker) ChatGroupWithSendingJobSeekerList {
	return ChatGroupWithSendingJobSeekerList{
		ChatGroupWithSendingJobSeekerList: chatGroupWithSendingJobSeekers,
	}
}

type ChatGroupWithSendingJobSeekerListAndMaxPage struct {
	ChatGroupWithSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker `json:"chat_group_with_sending_job_seeker_list"`
	MaxPageNumber                     uint                                    `json:"max_page_number"`
}

func NewChatGroupWithSendingJobSeekerListAndMaxPage(chatGroupWithSendingJobSeekers []*entity.ChatGroupWithSendingJobSeeker, maxPageNumber uint) ChatGroupWithSendingJobSeekerListAndMaxPage {
	return ChatGroupWithSendingJobSeekerListAndMaxPage{
		ChatGroupWithSendingJobSeekerList: chatGroupWithSendingJobSeekers,
		MaxPageNumber:                     maxPageNumber,
	}
}

type ChatGroupWithSendingJobSeekerUnWatchedCount struct {
	UnWatchedCount uint `json:"unwatched_count"`
}

func NewChatGroupWithSendingJobSeekerUnWatchedCount(unWatchedCount uint) ChatGroupWithSendingJobSeekerUnWatchedCount {
	return ChatGroupWithSendingJobSeekerUnWatchedCount{
		UnWatchedCount: unWatchedCount,
	}
}
