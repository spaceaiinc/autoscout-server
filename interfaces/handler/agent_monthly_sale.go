package handler

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AgentMonthlySaleHandler interface {
	// 汎用系 API
	CreateAgentMonthlySale(param entity.CreateOrUpdateAgentMonthlyManagementParam) (presenter.Presenter, error)
	UpdateAgentMonthlySale(param entity.CreateOrUpdateAgentMonthlyManagementParam) (presenter.Presenter, error)
	GetAgentMonthlySaleList(agentID, managementID uint) (presenter.Presenter, error)
	GetDefaultPreviewAgentMonthlySaleListByAgentID(agentID uint) (presenter.Presenter, error)

	// sale_managementsテーブル関連
	GetSaleManagementListByAgentID(agentID uint) (presenter.Presenter, error)
	GetSaleManagementListAndStaffManagementByAgentID(agentID uint) (presenter.Presenter, error)
	GetSaleManagementByID(managementID uint) (presenter.Presenter, error)
	GetSaleManagementAndAgentMonthlyByID(managementID uint) (presenter.Presenter, error)

	GetSumOfStaffMonthlyByManagementID(managementID uint) (presenter.Presenter, error) // 個人売上の合算値を全体売り上げとして取得する
}
type AgentMonthlySaleHandlerImpl struct {
	agentMonthlySaleInteractor interactor.AgentMonthlySaleInteractor
}

func NewAgentMonthlySaleHandlerImpl(jslI interactor.AgentMonthlySaleInteractor) AgentMonthlySaleHandler {
	return &AgentMonthlySaleHandlerImpl{
		agentMonthlySaleInteractor: jslI,
	}
}

/****************************************************************************************/
// 汎用系 API
//
func (h *AgentMonthlySaleHandlerImpl) CreateAgentMonthlySale(param entity.CreateOrUpdateAgentMonthlyManagementParam) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.CreateAgentMonthlySale(interactor.CreateAgentMonthlySaleInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentMonthlySaleHandlerImpl) UpdateAgentMonthlySale(param entity.CreateOrUpdateAgentMonthlyManagementParam) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.UpdateAgentMonthlySale(interactor.UpdateAgentMonthlySaleInput{
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *AgentMonthlySaleHandlerImpl) GetAgentMonthlySaleList(agentID, managementID uint) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.GetAgentMonthlySaleList(interactor.GetAgentMonthlySaleListInput{
		AgentID:      agentID,
		ManagementID: managementID,
	})

	if err != nil {
		fmt.Println("handler error", err)
		return nil, err
	}

	// AgentSaleManagement *entity.AgentSaleManagement
	// AgentStaffSaleManagements []*entity.AgentStaffSaleManagement
	// AgentMonthlySales []*entity.AgentMonthlySale
	return presenter.NewAgentSaleManagementAndMonthlyListJSONPresenter(responses.NewAgentSaleManagementAndMonthlyList(
		output.AgentSaleManagement,
		output.AgentStaffSaleManagements,
		output.AgentMonthlySales,
	)), nil
}

func (h *AgentMonthlySaleHandlerImpl) GetDefaultPreviewAgentMonthlySaleListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.GetDefaultPreviewAgentMonthlySaleListByAgentID(interactor.GetDefaultPreviewAgentMonthlySaleListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentSaleManagementAndMonthlyListJSONPresenter(responses.NewAgentSaleManagementAndMonthlyList(
		output.AgentSaleManagement,
		output.AgentStaffSaleManagements,
		output.AgentMonthlySales,
	)), nil
}

/****************************************************************************************/
// sale_managementsテーブル API
//

func (h *AgentMonthlySaleHandlerImpl) GetSaleManagementListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.GetSaleManagementListByAgentID(interactor.GetSaleManagementListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentSaleManagementListJSONPresenter(responses.NewAgentSaleManagementList(output.AgentSaleManagementList)), nil
}

func (h *AgentMonthlySaleHandlerImpl) GetSaleManagementListAndStaffManagementByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.GetSaleManagementListAndStaffManagementByAgentID(interactor.GetSaleManagementListAndStaffManagementByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentSaleManagementListJSONPresenter(responses.NewAgentSaleManagementList(output.AgentSaleManagementList)), nil
}

func (h *AgentMonthlySaleHandlerImpl) GetSaleManagementByID(managementID uint) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.GetSaleManagementByID(interactor.GetSaleManagementByIDInput{
		ManagementID: managementID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentSaleManagementJSONPresenter(responses.NewAgentSaleManagement(output.AgentSaleManagement)), nil
}

func (h *AgentMonthlySaleHandlerImpl) GetSaleManagementAndAgentMonthlyByID(managementID uint) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.GetSaleManagementAndAgentMonthlyByID(interactor.GetSaleManagementAndAgentMonthlyByIDInput{
		ManagementID: managementID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentSaleManagementJSONPresenter(responses.NewAgentSaleManagement(output.AgentSaleManagement)), nil
}

func (h *AgentMonthlySaleHandlerImpl) GetSumOfStaffMonthlyByManagementID(managementID uint) (presenter.Presenter, error) {
	output, err := h.agentMonthlySaleInteractor.GetSumOfStaffMonthlyByManagementID(interactor.GetSumOfStaffMonthlyByManagementIDInput{
		ManagementID: managementID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentSaleManagementJSONPresenter(responses.NewAgentSaleManagement(output.AgentSaleManagement)), nil
}
