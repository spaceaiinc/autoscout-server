package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentStaff struct {
	AgentStaff *entity.AgentStaff `json:"agent_staff_user"`
}

func NewAgentStaff(agentStaffUser *entity.AgentStaff) AgentStaff {
	return AgentStaff{
		AgentStaff: agentStaffUser,
	}
}

type AgentStaffList struct {
	AgentStaffList []*entity.AgentStaff `json:"agent_staff_user_list"`
}

func NewAgentStaffList(agentStaffUsers []*entity.AgentStaff) AgentStaffList {
	return AgentStaffList{
		AgentStaffList: agentStaffUsers,
	}
}

type AgentStaffListAndMaxPage struct {
	AgentStaffList []*entity.AgentStaff `json:"agent_staff_user_list"`
	MaxPageNumber  uint                 `json:"max_page_number"`
}

func NewAgentStaffListAndMaxPage(AgentStaffs []*entity.AgentStaff, maxPageNumber uint) AgentStaffListAndMaxPage {
	return AgentStaffListAndMaxPage{
		AgentStaffList: AgentStaffs,
		MaxPageNumber:  maxPageNumber,
	}
}
