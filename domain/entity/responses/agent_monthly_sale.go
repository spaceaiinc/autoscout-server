package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentMonthlySale struct {
	AgentMonthlySale *entity.AgentMonthlySale `json:"agent_monthly_sales"`
}

func NewAgentMonthlySale(agentUser *entity.AgentMonthlySale) AgentMonthlySale {
	return AgentMonthlySale{
		AgentMonthlySale: agentUser,
	}
}

type AgentMonthlySaleList struct {
	AgentMonthlySaleList []*entity.AgentMonthlySale `json:"agent_monthly_sales_list"`
}

func NewAgentMonthlySaleList(agentUsers []*entity.AgentMonthlySale) AgentMonthlySaleList {
	return AgentMonthlySaleList{
		AgentMonthlySaleList: agentUsers,
	}
}
