package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentStaffJSONPresenter(resp responses.AgentStaff) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentStaffListJSONPresenter(resp responses.AgentStaffList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentStaffListAndMaxPageJSONPresenter(resp responses.AgentStaffListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}
