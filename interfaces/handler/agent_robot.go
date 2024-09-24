package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AgentRobotHandler interface {
	// 汎用系 API
	UpdateAgentRobot(agentRobotID uint, param entity.CreateOrUpdateAgentRobotParam) (presenter.Presenter, error)
	GetAgentRobotListByAgentID(agentID uint) (presenter.Presenter, error)
	GetAgentRobotByID(id uint) (presenter.Presenter, error)

	// Admin API
	CreateAgentRobot(param entity.CreateOrUpdateAgentRobotParam) (presenter.Presenter, error)
	DeleteAgentRobot(id uint) (presenter.Presenter, error)
}

type AgentRobotHandlerImpl struct {
	agentRobotInteractor interactor.AgentRobotInteractor
}

func NewAgentRobotHandlerImpl(arI interactor.AgentRobotInteractor) AgentRobotHandler {
	return &AgentRobotHandlerImpl{
		agentRobotInteractor: arI,
	}
}

/****************************************************************************************/
// 汎用系 API
//
// エージェントロボットを更新
func (h *AgentRobotHandlerImpl) UpdateAgentRobot(agentRobotID uint, param entity.CreateOrUpdateAgentRobotParam) (presenter.Presenter, error) {
	output, err := h.agentRobotInteractor.UpdateAgentRobot(interactor.UpdateAgentRobotInput{
		AgentRobotID: agentRobotID,
		UpdateParam:  param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentRobotJSONPresenter(responses.NewAgentRobot(output.AgentRobot)), nil
}

// エージェントロボットを取得
func (h *AgentRobotHandlerImpl) GetAgentRobotByID(agentRobotID uint) (presenter.Presenter, error) {
	output, err := h.agentRobotInteractor.GetAgentRobotByID(interactor.GetAgentRobotByIDInput{
		AgentRobotID: agentRobotID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentRobotJSONPresenter(responses.NewAgentRobot(output.AgentRobot)), nil
}

// エージェントIDからエージェントロボットを取得
func (h *AgentRobotHandlerImpl) GetAgentRobotListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.agentRobotInteractor.GetAgentRobotListByAgentID(interactor.GetAgentRobotListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentRobotListJSONPresenter(responses.NewAgentRobotList(output.AgentRobotList)), nil
}

/****************************************************************************************/
/****************************************************************************************/
// Admin API
//
// エージェントロボットを作成
func (h *AgentRobotHandlerImpl) CreateAgentRobot(param entity.CreateOrUpdateAgentRobotParam) (presenter.Presenter, error) {
	output, err := h.agentRobotInteractor.CreateAgentRobot(interactor.CreateAgentRobotInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewAgentRobotJSONPresenter(responses.NewAgentRobot(output.AgentRobot)), nil
}

// エージェントロボットを削除
func (h *AgentRobotHandlerImpl) DeleteAgentRobot(agentRobotID uint) (presenter.Presenter, error) {
	output, err := h.agentRobotInteractor.DeleteAgentRobot(interactor.DeleteAgentRobotInput{
		AgentRobotID: agentRobotID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
