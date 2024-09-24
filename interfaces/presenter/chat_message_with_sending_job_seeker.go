package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewChatMessageWithSendingJobSeekerJSONPresenter(resp responses.ChatMessageWithSendingJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatMessageWithSendingJobSeekerListJSONPresenter(resp responses.ChatMessageWithSendingJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}
