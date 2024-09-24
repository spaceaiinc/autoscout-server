package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewJobSeekerScheduleJSONPresenter(resp responses.JobSeekerSchedule) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerScheduleListJSONPresenter(resp responses.JobSeekerScheduleList) Presenter {
	return NewJSONPresenter(200, resp)
}
