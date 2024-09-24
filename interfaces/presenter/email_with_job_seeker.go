package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewEmailWithJobSeekerJSONPresenter(resp responses.EmailWithJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewEmailWithJobSeekerListJSONPresenter(resp responses.EmailWithJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}
