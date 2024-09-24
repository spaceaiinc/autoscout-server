package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewJobSeekerJSONPresenter(resp responses.JobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerListJSONPresenter(resp responses.JobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerListAndMaxPageAndIDListJSONPresenter(resp responses.JobSeekerListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerListAndMaxPageAndIDListAndListCountJSONPresenter(resp responses.JobSeekerListAndMaxPageAndIDListAndListCount) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerDocumentJSONPresenter(resp responses.JobSeekerDocument) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSelectListForCreateOrUpdateJobSeekerJSONPresenter(resp responses.SelectListForCreateOrUpdateJobSeeker) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerForGuestJSONPresenter(resp responses.JobSeekerForGuest) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerDesiredForGuestJSONPresenter(resp responses.JobSeekerDesiredForGuest) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerUUIDJSONPresenter(resp responses.JobSeekerUUID) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerAgentIDJSONPresenter(resp responses.JobSeekerAgentID) Presenter {
	return NewJSONPresenter(200, resp)
}

// LPの登録状況を返す（LPのみで使用）
func NewJobSeekerRegisterStatusJSONPresenter(resp responses.JobSeekerRegisterStatus) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerLoginFromLPJSONPresenter(resp responses.JobSeekerLoginFromLP) Presenter {
	return NewJSONPresenter(200, resp)
}
