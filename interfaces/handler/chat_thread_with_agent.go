package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ChatThreadWithAgentHandler interface {
	// 汎用系 API
	CreateChatThreadWithAgent(param entity.CreateChatThreadWithAgentParam) (presenter.Presenter, error)

	GetChatThreadAndMessageWithAgentByChatThreadID(threadID uint) (presenter.Presenter, error)
}

type ChatThreadWithAgentHandlerImpl struct {
	chatThreadWithAgentInteractor interactor.ChatThreadWithAgentInteractor
}

func NewChatThreadWithAgentHandlerImpl(cmI interactor.ChatThreadWithAgentInteractor) ChatThreadWithAgentHandler {
	return &ChatThreadWithAgentHandlerImpl{
		chatThreadWithAgentInteractor: cmI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *ChatThreadWithAgentHandlerImpl) CreateChatThreadWithAgent(param entity.CreateChatThreadWithAgentParam) (presenter.Presenter, error) {
	output, err := h.chatThreadWithAgentInteractor.CreateChatThreadWithAgent(interactor.CreateChatThreadWithAgentInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatThreadWithAgentJSONPresenter(responses.NewChatThreadWithAgent(output.ChatThreadWithAgent)), nil
}

// グループIDからチャットメッセージを取得する
func (h *ChatThreadWithAgentHandlerImpl) GetChatThreadAndMessageWithAgentByChatThreadID(threadID uint) (presenter.Presenter, error) {
	output, err := h.chatThreadWithAgentInteractor.GetChatThreadAndMessageWithAgentByChatThreadID(interactor.GetChatThreadAndMessageWithAgentByChatThreadIDInput{
		ThreadID: threadID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatThreadWithAgentJSONPresenter(responses.NewChatThreadWithAgent(output.ChatThreadWithAgent)), nil
}
