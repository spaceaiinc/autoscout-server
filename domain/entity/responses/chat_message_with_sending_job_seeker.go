package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ChatMessageWithSendingJobSeeker struct {
	ChatMessageWithSendingJobSeeker *entity.ChatMessageWithSendingJobSeeker `json:"chat_message_with_sending_job_seeker"`
}

func NewChatMessageWithSendingJobSeeker(chatMessageWithSendingJobSeeker *entity.ChatMessageWithSendingJobSeeker) ChatMessageWithSendingJobSeeker {
	return ChatMessageWithSendingJobSeeker{
		ChatMessageWithSendingJobSeeker: chatMessageWithSendingJobSeeker,
	}
}

type ChatMessageWithSendingJobSeekerList struct {
	ChatMessageWithSendingJobSeekerList []*entity.ChatMessageWithSendingJobSeeker `json:"chat_message_with_sending_job_seeker_list"`
}

func NewChatMessageWithSendingJobSeekerList(chatMessageWithSendingJobSeekers []*entity.ChatMessageWithSendingJobSeeker) ChatMessageWithSendingJobSeekerList {
	return ChatMessageWithSendingJobSeekerList{
		ChatMessageWithSendingJobSeekerList: chatMessageWithSendingJobSeekers,
	}
}
