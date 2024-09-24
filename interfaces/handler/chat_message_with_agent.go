package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ChatMessageWithAgentHandler interface {
	// 汎用系 API
	SendChatMessageWithAgent(param entity.SendChatMessageWithAgentParam) (presenter.Presenter, error)
	GetChatMessageWithAgentListByThreadID(threadID, agentStaffID uint) (presenter.Presenter, error)

	// 最終閲覧時間の更新
	UpdateChatMessageWithAgentWatchedAtByThreadID(threadID, agentStaffID uint) (presenter.Presenter, error)
}

type ChatMessageWithAgentHandlerImpl struct {
	chatMessageWithAgentInteractor interactor.ChatMessageWithAgentInteractor
}

func NewChatMessageWithAgentHandlerImpl(cmI interactor.ChatMessageWithAgentInteractor) ChatMessageWithAgentHandler {
	return &ChatMessageWithAgentHandlerImpl{
		chatMessageWithAgentInteractor: cmI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *ChatMessageWithAgentHandlerImpl) SendChatMessageWithAgent(param entity.SendChatMessageWithAgentParam) (presenter.Presenter, error) {
	output, err := h.chatMessageWithAgentInteractor.SendChatMessageWithAgent(interactor.SendChatMessageWithAgentInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithAgentJSONPresenter(responses.NewChatMessageWithAgent(output.ChatMessageWithAgent)), nil
}

// グループIDからチャットメッセージを取得する
func (h *ChatMessageWithAgentHandlerImpl) GetChatMessageWithAgentListByThreadID(threadID, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.chatMessageWithAgentInteractor.GetChatMessageWithAgentListByThreadID(interactor.GetChatMessageWithAgentListByThreadIDInput{
		ThreadID:     threadID,
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithAgentListJSONPresenter(responses.NewChatMessageWithAgentList(output.ChatMessageWithAgentList)), nil
}

// グループIDからチャットメッセージを取得する
func (h *ChatMessageWithAgentHandlerImpl) UpdateChatMessageWithAgentWatchedAtByThreadID(threadID, agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.chatMessageWithAgentInteractor.UpdateChatMessageWithAgentWatchedAtByThreadID(interactor.UpdateChatMessageWithAgentWatchedAtByThreadIDInput{
		ThreadID:     threadID,
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
