package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewAgentSaleManagementJSONPresenter(resp responses.AgentSaleManagement) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentSaleManagementListJSONPresenter(resp responses.AgentSaleManagementList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewAgentSaleManagementAndMonthlyListJSONPresenter(resp responses.AgentSaleManagementAndMonthlyList) Presenter {
	return NewJSONPresenter(200, resp)
}
