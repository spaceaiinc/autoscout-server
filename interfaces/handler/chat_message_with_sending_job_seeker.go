package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ChatMessageWithSendingJobSeekerHandler interface {
	// 汎用系 API
	SendChatMessageWithSendingJobSeekerLineMessage(param entity.SendChatMessageWithSendingJobSeekerLineParam) (presenter.Presenter, error)
	SendChatMessageWithSendingJobSeekerLineImage(param entity.SendChatMessageWithSendingJobSeekerLineImageParam) (presenter.Presenter, error)
	GetChatMessageWithSendingJobSeekerListByGroupID(groupID uint) (presenter.Presenter, error)

	// Admin API
}

type ChatMessageWithSendingJobSeekerHandlerImpl struct {
	chatMessageWithSendingJobSeekerInteractor interactor.ChatMessageWithSendingJobSeekerInteractor
}

func NewChatMessageWithSendingJobSeekerHandlerImpl(cmI interactor.ChatMessageWithSendingJobSeekerInteractor) ChatMessageWithSendingJobSeekerHandler {
	return &ChatMessageWithSendingJobSeekerHandlerImpl{
		chatMessageWithSendingJobSeekerInteractor: cmI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *ChatMessageWithSendingJobSeekerHandlerImpl) SendChatMessageWithSendingJobSeekerLineMessage(param entity.SendChatMessageWithSendingJobSeekerLineParam) (presenter.Presenter, error) {
	output, err := h.chatMessageWithSendingJobSeekerInteractor.SendChatMessageWithSendingJobSeekerLineMessage(interactor.SendChatMessageWithSendingJobSeekerLineMessageInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithSendingJobSeekerJSONPresenter(responses.NewChatMessageWithSendingJobSeeker(output.ChatMessageWithSendingJobSeeker)), nil
}

func (h *ChatMessageWithSendingJobSeekerHandlerImpl) SendChatMessageWithSendingJobSeekerLineImage(param entity.SendChatMessageWithSendingJobSeekerLineImageParam) (presenter.Presenter, error) {
	output, err := h.chatMessageWithSendingJobSeekerInteractor.SendChatMessageWithSendingJobSeekerLineImage(interactor.SendChatMessageWithSendingJobSeekerLineImageInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithSendingJobSeekerJSONPresenter(responses.NewChatMessageWithSendingJobSeeker(output.ChatMessageWithSendingJobSeeker)), nil
}

// グループIDからチャットメッセージを取得する
func (h *ChatMessageWithSendingJobSeekerHandlerImpl) GetChatMessageWithSendingJobSeekerListByGroupID(groupID uint) (presenter.Presenter, error) {
	output, err := h.chatMessageWithSendingJobSeekerInteractor.GetChatMessageWithSendingJobSeekerListByGroupID(interactor.GetChatMessageWithSendingJobSeekerListByGroupIDInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithSendingJobSeekerListJSONPresenter(responses.NewChatMessageWithSendingJobSeekerList(output.ChatMessageWithSendingJobSeekerList)), nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API

/****************************************************************************************/
