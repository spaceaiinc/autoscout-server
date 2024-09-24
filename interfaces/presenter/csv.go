package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewMissedRecordsAndOKJSONPresenter(resp responses.MissedRecordsAndOK) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewMissedRecordsAndEnterpriseListJSONPresenter(resp responses.MissedRecordsAndEnterpriseList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewMissedRecordsAndJobInformationListJSONPresenter(resp responses.MissedRecordsAndJobInformationList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewMissedRecordsAndJobSeekerListJSONPresenter(resp responses.MissedRecordsAndJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}
