package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentJSONPresenter(resp responses.Agent) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentListJSONPresenter(resp responses.AgentList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentListAndMaxPageAndIDListJSONPresenter(resp responses.AgentListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentLineJSONPresenter(resp responses.AgentLine) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentLineLoginChannelIDJSONPresenter(resp responses.AgentLineLoginChannelID) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentBotInformationJSONPresenter(resp responses.AgentBotInformation) Presenter {
	return NewJSONPresenter(200, resp)
}
