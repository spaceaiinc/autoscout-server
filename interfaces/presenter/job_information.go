package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewJobInformationJSONPresenter(resp responses.JobInformation) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobInformationListJSONPresenter(resp responses.JobInformationList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobInformationListAndMaxPageAndIDListJSONPresenter(resp responses.JobInformationListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobInformationListAndMaxPageAndIDListAndListCountJSONPresenter(resp responses.JobInformationListAndMaxPageAndIDListAndListCount) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobInformationCountAndIncomeJSONPresenter(resp responses.JobInformationCountAndIncome) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobListingListAndMaxPageJSONPresenter(resp responses.JobListingListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobListingListAndCountJSONPresenter(resp responses.JobListingListAndCount) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobListingJSONPresenter(resp responses.JobListing) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobListingListJSONPresenter(resp responses.JobListingList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobListingListForJobSeekerJSONPresenter(resp responses.JobListingListForJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobListingListAndJobSeekerDesiredJSONPresenter(resp responses.JobListingListAndJobSeekerDesired) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobInformationListForDiagnosisJSONPresenter(resp responses.JobInformationListForDiagnosis) Presenter {
	return NewJSONPresenter(200, resp)
}
