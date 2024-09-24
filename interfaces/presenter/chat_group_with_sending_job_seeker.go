package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewChatGroupWithSendingJobSeekerJSONPresenter(resp responses.ChatGroupWithSendingJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatGroupWithSendingJobSeekerListJSONPresenter(resp responses.ChatGroupWithSendingJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatGroupWithSendingJobSeekerListAndMaxPageJSONPresenter(resp responses.ChatGroupWithSendingJobSeekerListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatGroupWithSendingJobSeekerUnwatchedCountJSONPresenter(resp responses.ChatGroupWithSendingJobSeekerUnWatchedCount) Presenter {
	return NewJSONPresenter(200, resp)
}
