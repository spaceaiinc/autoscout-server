package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AgentAllianceHandler interface {
	// 汎用系 API
	RequestAlliance(param entity.CreateAgentAllianceParam) (presenter.Presenter, error)
	UpdateOrCreateAgentAllianceList(param entity.MyAgentIDAndOtherAgentIDListParam) (presenter.Presenter, error)
	UpdateAgentAllianceCancelRequest(param entity.UpdateAgentAllianceCancelRequestParam) (presenter.Presenter, error)
	GetAgentAllianceByID(agentAllianceID uint) (presenter.Presenter, error)
	GetAgentAllianceListByAgentID(agentID uint) (presenter.Presenter, error)
	CheckAnyAllianceWithoutApplication(myAgentID uint, otherAgentIDList []uint) (presenter.Presenter, error)
}

type AgentAllianceHandlerImpl struct {
	AgentAllianceInteractor interactor.AgentAllianceInteractor
}

func NewAgentAllianceHandlerImpl(acuI interactor.AgentAllianceInteractor) AgentAllianceHandler {
	return &AgentAllianceHandlerImpl{
		AgentAllianceInteractor: acuI,
	}
}

/****************************************************************************************/
// 汎用系 API
//

func (h *AgentAllianceHandlerImpl) RequestAlliance(param entity.CreateAgentAllianceParam) (presenter.Presenter, error) {
	output, err := h.AgentAllianceInteractor.RequestAlliance(interactor.RequestAllianceInput{
		Param: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentAllianceJSONPresenter(responses.NewAgentAlliance(output.Alliance)), nil
}

func (h *AgentAllianceHandlerImpl) UpdateOrCreateAgentAllianceList(param entity.MyAgentIDAndOtherAgentIDListParam) (presenter.Presenter, error) {
	output, err := h.AgentAllianceInteractor.UpdateOrCreateAgentAllianceList(interactor.UpdateOrCreateAgentAllianceListInput{
		Param: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentAllianceListJSONPresenter(responses.NewAgentAllianceList(output.AgentAllianceList)), nil
}

func (h *AgentAllianceHandlerImpl) UpdateAgentAllianceCancelRequest(param entity.UpdateAgentAllianceCancelRequestParam) (presenter.Presenter, error) {
	output, err := h.AgentAllianceInteractor.UpdateAgentAllianceCancelRequest(interactor.UpdateAgentAllianceCancelRequestInput{
		Param: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentAllianceAndRemainingTaskListJSONPresenter(
		responses.NewAgentAllianceAndRemainingTaskList(output.AgentAlliance, output.RemainingTaskList),
	), nil
}

func (h *AgentAllianceHandlerImpl) GetAgentAllianceByID(agentAllianceID uint) (presenter.Presenter, error) {
	output, err := h.AgentAllianceInteractor.GetAgentAllianceByID(interactor.GetAgentAllianceByIDInput{
		AgentAllianceID: agentAllianceID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentAllianceJSONPresenter(responses.NewAgentAlliance(output.AgentAlliance)), nil
}

func (h *AgentAllianceHandlerImpl) GetAgentAllianceListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.AgentAllianceInteractor.GetAgentAllianceListByAgentID(interactor.GetAgentAllianceListByAgentIDInput{
		AgentID: agentID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewAgentAllianceListJSONPresenter(responses.NewAgentAllianceList(output.AgentAllianceList)), nil
}

func (h *AgentAllianceHandlerImpl) CheckAnyAllianceWithoutApplication(myAgentID uint, otherAgentIDList []uint) (presenter.Presenter, error) {
	output, err := h.AgentAllianceInteractor.CheckAnyAllianceWithoutApplication(interactor.CheckAnyAllianceWithoutApplicationInput{
		MyAgentID:        myAgentID,
		OtherAgentIDList: otherAgentIDList,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewIDListJSONPresenter(responses.NewIDList(output.IDList)), nil
}
