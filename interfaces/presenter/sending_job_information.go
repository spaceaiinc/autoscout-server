package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingJobInformationJSONPresenter(resp responses.SendingJobInformation) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobInformationListJSONPresenter(resp responses.SendingJobInformationList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobInformationListAndMaxPageAndIDListJSONPresenter(resp responses.SendingJobInformationListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobInformationListAndMaxPageAndIDListAndListCountJSONPresenter(resp responses.SendingJobInformationListAndMaxPageAndIDListAndListCount) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingEnterpriseListAndJobInformationListAndMaxPageJSONPresenter(resp responses.SendingEnterpriseListAndJobInformationListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobListingForSendingJSONPresenter(resp responses.JobListingForSending) Presenter {
	return NewJSONPresenter(200, resp)
}
