package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewChatGroupWithJobSeekerJSONPresenter(resp responses.ChatGroupWithJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatGroupWithJobSeekerListJSONPresenter(resp responses.ChatGroupWithJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatGroupWithJobSeekerListAndMaxPageJSONPresenter(resp responses.ChatGroupWithJobSeekerListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatGroupWithJobSeekerUnwatchedCountJSONPresenter(resp responses.ChatGroupWithJobSeekerUnWatchedCount) Presenter {
	return NewJSONPresenter(200, resp)
}
