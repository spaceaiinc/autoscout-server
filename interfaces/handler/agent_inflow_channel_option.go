package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AgentInflowChannelOptionHandler interface {
	// 汎用系 API
	CreateAgentInflowChannelOption(param entity.AgentInflowChannelOption) (presenter.Presenter, error)
	UpdateAgentInflowChannelOption(agentInflowChannelOptionID uint, param entity.AgentInflowChannelOption) (presenter.Presenter, error)
	GetAgentInflowChannelOptionByID(id uint) (presenter.Presenter, error)
	GetAgentInflowChannelOptionListByAgentID(agentID uint) (presenter.Presenter, error)
}

type AgentInflowChannelOptionHandlerImpl struct {
	agentInflowChannelOptionInteractor interactor.AgentInflowChannelOptionInteractor
}

func NewAgentInflowChannelOptionHandlerImpl(arI interactor.AgentInflowChannelOptionInteractor) AgentInflowChannelOptionHandler {
	return &AgentInflowChannelOptionHandlerImpl{
		agentInflowChannelOptionInteractor: arI,
	}
}

/****************************************************************************************/
// 汎用系 API
//
// エージェントの流入経路マスタを作成
func (h *AgentInflowChannelOptionHandlerImpl) CreateAgentInflowChannelOption(param entity.AgentInflowChannelOption) (presenter.Presenter, error) {
	output, err := h.agentInflowChannelOptionInteractor.CreateAgentInflowChannelOption(interactor.CreateAgentInflowChannelOptionInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentInflowChannelOptionJSONPresenter(responses.NewAgentInflowChannelOption(output.AgentInflowChannelOption)), nil
}

// エージェントの流入経路マスタを更新
func (h *AgentInflowChannelOptionHandlerImpl) UpdateAgentInflowChannelOption(agentInflowChannelOptionID uint, param entity.AgentInflowChannelOption) (presenter.Presenter, error) {
	output, err := h.agentInflowChannelOptionInteractor.UpdateAgentInflowChannelOption(interactor.UpdateAgentInflowChannelOptionInput{
		AgentInflowChannelOptionID: agentInflowChannelOptionID,
		UpdateParam:                param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentInflowChannelOptionJSONPresenter(responses.NewAgentInflowChannelOption(output.AgentInflowChannelOption)), nil
}

// エージェントの流入経路マスタを取得
func (h *AgentInflowChannelOptionHandlerImpl) GetAgentInflowChannelOptionByID(agentInflowChannelOptionID uint) (presenter.Presenter, error) {
	output, err := h.agentInflowChannelOptionInteractor.GetAgentInflowChannelOptionByID(interactor.GetAgentInflowChannelOptionByIDInput{
		AgentInflowChannelOptionID: agentInflowChannelOptionID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentInflowChannelOptionJSONPresenter(responses.NewAgentInflowChannelOption(output.AgentInflowChannelOption)), nil
}

// エージェントIDからエージェントの流入経路マスタを取得
func (h *AgentInflowChannelOptionHandlerImpl) GetAgentInflowChannelOptionListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.agentInflowChannelOptionInteractor.GetAgentInflowChannelOptionListByAgentID(interactor.GetAgentInflowChannelOptionListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentInflowChannelOptionListJSONPresenter(responses.NewAgentInflowChannelOptionList(output.AgentInflowChannelOptionList)), nil
}
