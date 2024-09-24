package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingPhaseJSONPresenter(resp responses.SendingPhase) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingPhaseListJSONPresenter(resp responses.SendingPhaseList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingJobSeekerTableListAndMaxPageAndIDListAndListCountJSONPresenter(resp responses.SendingJobSeekerTableListAndMaxPageAndIDListAndListCount) Presenter {
	return NewJSONPresenter(200, resp)
}
