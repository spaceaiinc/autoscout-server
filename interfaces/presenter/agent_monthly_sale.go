package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentMonthlySaleJSONPresenter(resp responses.AgentMonthlySale) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentMonthlySaleListJSONPresenter(resp responses.AgentMonthlySaleList) Presenter {
	return NewJSONPresenter(200, resp)
}
