package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingJobSeekerJSONPresenter(resp responses.SendingJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobSeekerListJSONPresenter(resp responses.SendingJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobSeekerListAndMaxPageAndIDListJSONPresenter(resp responses.SendingJobSeekerListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobSeekerDocumentJSONPresenter(resp responses.SendingJobSeekerDocument) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewIsNotViewCountJSONPresenter(resp responses.IsNotViewCount) Presenter {
	return NewJSONPresenter(200, resp)
}
