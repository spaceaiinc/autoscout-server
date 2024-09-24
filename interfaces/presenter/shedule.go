package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewScheduleJSONPresenter(resp responses.Schedule) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewScheduleListJSONPresenter(resp responses.ScheduleList) Presenter {
	return NewJSONPresenter(200, resp)
}
