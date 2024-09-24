package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AgentStaffMonthlySaleHandler interface {
	// 汎用系 API
	CreateAgentStaffMonthlySale(param entity.CreateOrUpdateStaffMonthlyManagementParam) (presenter.Presenter, error)
	UpdateAgentStaffMonthlySale(param entity.CreateOrUpdateStaffMonthlyManagementParam) (presenter.Presenter, error)
	GetStaffMonthlySaleList(agentStaffID, managementID uint) (presenter.Presenter, error)
	GetStaffSaleManagementAndAgentMonthlyByID(agentStaffID, managementID uint) (presenter.Presenter, error)
}

type AgentStaffMonthlySaleHandlerImpl struct {
	agentStaffMonthlySaleInteractor interactor.AgentStaffMonthlySaleInteractor
}

func NewAgentStaffMonthlySaleHandlerImpl(asmI interactor.AgentStaffMonthlySaleInteractor) AgentStaffMonthlySaleHandler {
	return &AgentStaffMonthlySaleHandlerImpl{
		agentStaffMonthlySaleInteractor: asmI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *AgentStaffMonthlySaleHandlerImpl) CreateAgentStaffMonthlySale(param entity.CreateOrUpdateStaffMonthlyManagementParam) (presenter.Presenter, error) {
	output, err := h.agentStaffMonthlySaleInteractor.CreateAgentStaffMonthlySale(interactor.CreateAgentStaffMonthlySaleInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentStaffMonthlySaleHandlerImpl) UpdateAgentStaffMonthlySale(param entity.CreateOrUpdateStaffMonthlyManagementParam) (presenter.Presenter, error) {
	output, err := h.agentStaffMonthlySaleInteractor.UpdateAgentStaffMonthlySale(interactor.UpdateAgentStaffMonthlySaleInput{
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentStaffMonthlySaleHandlerImpl) GetStaffMonthlySaleList(agentStaffID, managementID uint) (presenter.Presenter, error) {
	output, err := h.agentStaffMonthlySaleInteractor.GetStaffMonthlySaleList(interactor.GetStaffMonthlySaleListInput{
		AgentStaffID: agentStaffID,
		ManagementID: managementID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffSaleManagementJSONPresenter(responses.NewAgentStaffSaleManagement(output.StaffSaleManagement)), nil
}

func (h *AgentStaffMonthlySaleHandlerImpl) GetStaffSaleManagementAndAgentMonthlyByID(agentStaffID, managementID uint) (presenter.Presenter, error) {
	output, err := h.agentStaffMonthlySaleInteractor.GetStaffSaleManagementAndAgentMonthlyByID(interactor.GetStaffSaleManagementAndAgentMonthlyByIDInput{
		AgentStaffID: agentStaffID,
		ManagementID: managementID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentStaffSaleManagementJSONPresenter(responses.NewAgentStaffSaleManagement(output.StaffSaleManagement)), nil
}
