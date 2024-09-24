package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewChatMessageWithJobSeekerJSONPresenter(resp responses.ChatMessageWithJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatMessageWithJobSeekerListJSONPresenter(resp responses.ChatMessageWithJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}
