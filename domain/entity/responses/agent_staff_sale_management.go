package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentStaffSaleManagement struct {
	AgentStaffSaleManagement *entity.AgentStaffSaleManagement `json:"staff_sale_management"`
}

func NewAgentStaffSaleManagement(agentUser *entity.AgentStaffSaleManagement) AgentStaffSaleManagement {
	return AgentStaffSaleManagement{
		AgentStaffSaleManagement: agentUser,
	}
}

type AgentStaffSaleManagementList struct {
	AgentStaffSaleManagementList []*entity.AgentStaffSaleManagement `json:"staff_sale_management_list"`
}

func NewAgentStaffSaleManagementList(agentUsers []*entity.AgentStaffSaleManagement) AgentStaffSaleManagementList {
	return AgentStaffSaleManagementList{
		AgentStaffSaleManagementList: agentUsers,
	}
}
