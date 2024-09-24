package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentStaffMonthlySaleJSONPresenter(resp responses.AgentStaffMonthlySale) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentStaffMonthlySaleListJSONPresenter(resp responses.AgentStaffMonthlySaleList) Presenter {
	return NewJSONPresenter(200, resp)
}
