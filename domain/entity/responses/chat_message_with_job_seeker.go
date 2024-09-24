package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ChatMessageWithJobSeeker struct {
	ChatMessageWithJobSeeker *entity.ChatMessageWithJobSeeker `json:"chat_message_with_job_seeker"`
}

func NewChatMessageWithJobSeeker(chatMessageWithJobSeeker *entity.ChatMessageWithJobSeeker) ChatMessageWithJobSeeker {
	return ChatMessageWithJobSeeker{
		ChatMessageWithJobSeeker: chatMessageWithJobSeeker,
	}
}

type ChatMessageWithJobSeekerList struct {
	ChatMessageWithJobSeekerList []*entity.ChatMessageWithJobSeeker `json:"chat_message_with_job_seeker_list"`
}

func NewChatMessageWithJobSeekerList(chatMessageWithJobSeekers []*entity.ChatMessageWithJobSeeker) ChatMessageWithJobSeekerList {
	return ChatMessageWithJobSeekerList{
		ChatMessageWithJobSeekerList: chatMessageWithJobSeekers,
	}
}
