package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingJobSeekerDesiredJobInformationJSONPresenter(resp responses.SendingJobSeekerDesiredJobInformation) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobSeekerDesiredJobInformationListJSONPresenter(resp responses.SendingJobSeekerDesiredJobInformationList) Presenter {
	return NewJSONPresenter(200, resp)
}
