package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentAllianceJSONPresenter(resp responses.AgentAlliance) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentAllianceListJSONPresenter(resp responses.AgentAllianceList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentAllianceAndRemainingTaskListJSONPresenter(resp responses.AgentAllianceAndRemainingTaskList) Presenter {
	return NewJSONPresenter(200, resp)
}
