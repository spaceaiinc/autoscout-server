package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ChatGroupWithAgentHandler interface {
	// 汎用系 API
	CreateChatGroupWithAgent(param entity.CreateChatGroupWithAgentParam) (presenter.Presenter, error)
	CreateChatGroupWithAgentList(param entity.MyAgentIDAndOtherAgentIDListParam) (presenter.Presenter, error)
	UpdateAgentLastWatched(groupID uint, agentID uint) (presenter.Presenter, error)
	GetChatGroupWithAgentByID(agent_id uint, id uint) (presenter.Presenter, error)
	GetChatGroupWithAgentListByAgentID(agentID uint) (presenter.Presenter, error)
	GetChatGroupAndThreadWithAgentByChatGroupID(groupID, agentStaffID uint) (presenter.Presenter, error)
	GetChatGroupWithAgentListByAgentIDAndAgentStaffID(agentID, agentStaffID uint) (presenter.Presenter, error)
	GetChatGroupWithAgentByAgentIDAndOtherAgentID(agentID, otherAgentID uint) (presenter.Presenter, error)
	CheckAnyChatGroupWithoutGroup(myAgentID uint, otherAgentIDList []uint) (presenter.Presenter, error)
}

type ChatGroupWithAgentHandlerImpl struct {
	chatGroupWithAgentInteractor interactor.ChatGroupWithAgentInteractor
}

func NewChatGroupWithAgentHandlerImpl(cgI interactor.ChatGroupWithAgentInteractor) ChatGroupWithAgentHandler {
	return &ChatGroupWithAgentHandlerImpl{
		chatGroupWithAgentInteractor: cgI,
	}
}

/****************************************************************************************/
// 汎用系 API
//

func (h *ChatGroupWithAgentHandlerImpl) CreateChatGroupWithAgent(param entity.CreateChatGroupWithAgentParam) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.CreateChatGroupWithAgent(interactor.CreateChatGroupWithAgentInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithAgentJSONPresenter(responses.NewChatGroupWithAgent(output.ChatGroupWithAgent)), nil
}

func (h *ChatGroupWithAgentHandlerImpl) CreateChatGroupWithAgentList(param entity.MyAgentIDAndOtherAgentIDListParam) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.CreateChatGroupWithAgentList(interactor.CreateChatGroupWithAgentListInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithAgentListJSONPresenter(responses.NewChatGroupWithAgentList(output.ChatGroupWithAgentList)), nil
}

func (h *ChatGroupWithAgentHandlerImpl) UpdateAgentLastWatched(groupID uint, agentID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.UpdateAgentLastWatched(interactor.UpdateAgentLastWatchedForAgentChatInput{
		GroupID: groupID,
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// IDからチャットグループを取得する
func (h *ChatGroupWithAgentHandlerImpl) GetChatGroupWithAgentByID(agent_id uint, id uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.GetChatGroupWithAgentByID(interactor.GetChatGroupWithAgentByIDInput{
		AgentID: agent_id,
		ID:      id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithAgentJSONPresenter(responses.NewChatGroupWithAgent(output.ChatGroupWithAgent)), nil
}

// グループIDからチャットグループを取得する
func (h *ChatGroupWithAgentHandlerImpl) GetChatGroupWithAgentListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.GetChatGroupWithAgentListByAgentID(interactor.GetChatGroupWithAgentListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithAgentListJSONPresenter(responses.NewChatGroupWithAgentList(output.ChatGroupWithAgentList)), nil
}

// グループIDからチャットメッセージを取得する
func (h *ChatGroupWithAgentHandlerImpl) GetChatGroupAndThreadWithAgentByChatGroupID(groupID, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.GetChatGroupAndThreadWithAgentByChatGroupID(interactor.GetChatGroupAndThreadWithAgentByChatGroupIDInput{
		GroupID:      groupID,
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithAgentJSONPresenter(responses.NewChatGroupWithAgent(output.ChatGroupWithAgent)), nil
}

// エージェントIDとエージェントスタッフIDからチャットグループを取得する
func (h *ChatGroupWithAgentHandlerImpl) GetChatGroupWithAgentByAgentIDAndOtherAgentID(agentID, otherAgentID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.GetChatGroupWithAgentByAgentIDAndOtherAgentID(interactor.GetChatGroupWithAgentByAgentIDAndOtherAgentIDInput{
		AgentID:      agentID,
		OtherAgentID: otherAgentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithAgentJSONPresenter(responses.NewChatGroupWithAgent(output.ChatGroupWithAgent)), nil
}

// エージェントIDとエージェントスタッフIDからチャットグループを取得する
func (h *ChatGroupWithAgentHandlerImpl) GetChatGroupWithAgentListByAgentIDAndAgentStaffID(agentID, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.GetChatGroupWithAgentListByAgentIDAndAgentStaffID(interactor.GetChatGroupWithAgentListByAgentIDAndAgentStaffIDInput{
		AgentID:      agentID,
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithAgentListJSONPresenter(responses.NewChatGroupWithAgentList(output.ChatGroupWithAgentList)), nil
}

func (h *ChatGroupWithAgentHandlerImpl) CheckAnyChatGroupWithoutGroup(myAgentID uint, otherAgentIDList []uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithAgentInteractor.CheckAnyChatGroupWithoutGroup(interactor.CheckAnyChatGroupWithoutGroupInput{
		MyAgentID:        myAgentID,
		OtherAgentIDList: otherAgentIDList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewIDListJSONPresenter(responses.NewIDList(output.IDList)), nil
}
