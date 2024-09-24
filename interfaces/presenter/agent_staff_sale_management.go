package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentStaffSaleManagementJSONPresenter(resp responses.AgentStaffSaleManagement) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentStaffSaleManagementListJSONPresenter(resp responses.AgentStaffSaleManagementList) Presenter {
	return NewJSONPresenter(200, resp)
}
