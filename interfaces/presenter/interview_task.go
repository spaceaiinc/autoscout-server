package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewInterviewTaskJSONPresenter(resp responses.InterviewTask) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewInterviewTaskListJSONPresenter(resp responses.InterviewTaskList) Presenter {
	return NewJSONPresenter(200, resp)
}
