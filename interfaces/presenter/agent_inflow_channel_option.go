package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentInflowChannelOptionJSONPresenter(resp responses.AgentInflowChannelOption) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentInflowChannelOptionListJSONPresenter(resp responses.AgentInflowChannelOptionList) Presenter {
	return NewJSONPresenter(200, resp)
}
