package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentRobotJSONPresenter(resp responses.AgentRobot) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentRobotListJSONPresenter(resp responses.AgentRobotList) Presenter {
	return NewJSONPresenter(200, resp)
}
