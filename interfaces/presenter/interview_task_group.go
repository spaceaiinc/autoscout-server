package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewInterviewTaskGroupJSONPresenter(resp responses.InterviewTaskGroup) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewInterviewTaskGroupListJSONPresenter(resp responses.InterviewTaskGroupList) Presenter {
	return NewJSONPresenter(200, resp)
}
