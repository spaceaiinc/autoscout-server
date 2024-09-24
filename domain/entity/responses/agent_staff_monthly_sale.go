package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentStaffMonthlySale struct {
	AgentStaffMonthlySale *entity.AgentStaffMonthlySale `json:"agent_staff_sale_target_management"`
}

func NewAgentStaffMonthlySale(agentStaffUser *entity.AgentStaffMonthlySale) AgentStaffMonthlySale {
	return AgentStaffMonthlySale{
		AgentStaffMonthlySale: agentStaffUser,
	}
}

type AgentStaffMonthlySaleList struct {
	AgentStaffMonthlySaleList []*entity.AgentStaffMonthlySale `json:"agent_staff_sale_target_management_list"`
}

func NewAgentStaffMonthlySaleList(agentStaffUsers []*entity.AgentStaffMonthlySale) AgentStaffMonthlySaleList {
	return AgentStaffMonthlySaleList{
		AgentStaffMonthlySaleList: agentStaffUsers,
	}
}
