package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentSaleManagement struct {
	AgentSaleManagement *entity.AgentSaleManagement `json:"agent_sale_management"`
}

func NewAgentSaleManagement(agentSaleManagement *entity.AgentSaleManagement) AgentSaleManagement {

	return AgentSaleManagement{
		AgentSaleManagement: agentSaleManagement,
	}
}

type AgentSaleManagementList struct {
	AgentSaleManagementList []*entity.AgentSaleManagement `json:"agent_sale_management_list"`
}

func NewAgentSaleManagementList(agentSaleManagements []*entity.AgentSaleManagement) AgentSaleManagementList {
	return AgentSaleManagementList{
		AgentSaleManagementList: agentSaleManagements,
	}
}

type AgentSaleManagementAndMonthlyList struct {
	AgentSaleManagement       *entity.AgentSaleManagement        `json:"agent_sale_management"`
	AgentStaffSaleManagements []*entity.AgentStaffSaleManagement `json:"agent_staff_sale_management_list"`
	AgentMonthlySales         []*entity.AgentMonthlySale         `json:"agent_monthly_sale_list"`
}

func NewAgentSaleManagementAndMonthlyList(
	agentSaleManagement *entity.AgentSaleManagement,
	agentStaffSaleManagements []*entity.AgentStaffSaleManagement,
	agentMonthlySales []*entity.AgentMonthlySale,
) AgentSaleManagementAndMonthlyList {
	return AgentSaleManagementAndMonthlyList{
		AgentSaleManagement:       agentSaleManagement,
		AgentStaffSaleManagements: agentStaffSaleManagements,
		AgentMonthlySales:         agentMonthlySales,
	}
}
