package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentAlliance struct {
	AgentAlliance *entity.AgentAlliance `json:"agent_alliance"`
}

func NewAgentAlliance(agentAlliance *entity.AgentAlliance) AgentAlliance {
	return AgentAlliance{
		AgentAlliance: agentAlliance,
	}
}

type AgentAllianceList struct {
	AgentAllianceList []*entity.AgentAlliance `json:"agent_alliance_list"`
}

func NewAgentAllianceList(agentAlliances []*entity.AgentAlliance) AgentAllianceList {
	return AgentAllianceList{
		AgentAllianceList: agentAlliances,
	}
}

type AgentAllianceAndRemainingTaskList struct {
	AgentAlliance     *entity.AgentAlliance `json:"agent_alliance"`
	RemainingTaskList []*entity.Task        `json:"remaining_task_list"`
}

func NewAgentAllianceAndRemainingTaskList(
	agentAlliance *entity.AgentAlliance,
	remainingTaskList []*entity.Task,
) AgentAllianceAndRemainingTaskList {
	return AgentAllianceAndRemainingTaskList{
		AgentAlliance:     agentAlliance,
		RemainingTaskList: remainingTaskList,
	}
}
