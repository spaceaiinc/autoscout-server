package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentRobot struct {
	AgentRobot *entity.AgentRobot `json:"agent_robot"`
}

func NewAgentRobot(scoutService *entity.AgentRobot) AgentRobot {
	return AgentRobot{
		AgentRobot: scoutService,
	}
}

type AgentRobotList struct {
	AgentRobotList []*entity.AgentRobot `json:"agent_robot_list"`
}

func NewAgentRobotList(scoutServices []*entity.AgentRobot) AgentRobotList {
	return AgentRobotList{
		AgentRobotList: scoutServices,
	}
}
