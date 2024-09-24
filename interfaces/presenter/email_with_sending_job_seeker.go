package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewEmailWithSendingJobSeekerJSONPresenter(resp responses.EmailWithSendingJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewEmailWithSendingJobSeekerListJSONPresenter(resp responses.EmailWithSendingJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}
