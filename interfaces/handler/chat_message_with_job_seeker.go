package handler

import (
	"net/http"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ChatMessageWithJobSeekerHandler interface {
	// 汎用系 API
	SendChatMessageWithJobSeekerLineMessage(param entity.SendChatMessageWithJobSeekerLineParam) (presenter.Presenter, error)
	SendChatMessageWithJobSeekerLineImage(param entity.SendChatMessageWithJobSeekerLineImageParam) (presenter.Presenter, error)
	GetChatMessageWithJobSeekerListByGroupID(groupID uint) (presenter.Presenter, error)

	// LINE WebHook
	LineWebHook(req *http.Request) (presenter.Presenter, error)

	// Admin API
}

type ChatMessageWithJobSeekerHandlerImpl struct {
	chatMessageWithJobSeekerInteractor interactor.ChatMessageWithJobSeekerInteractor
}

func NewChatMessageWithJobSeekerHandlerImpl(cmI interactor.ChatMessageWithJobSeekerInteractor) ChatMessageWithJobSeekerHandler {
	return &ChatMessageWithJobSeekerHandlerImpl{
		chatMessageWithJobSeekerInteractor: cmI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *ChatMessageWithJobSeekerHandlerImpl) SendChatMessageWithJobSeekerLineMessage(param entity.SendChatMessageWithJobSeekerLineParam) (presenter.Presenter, error) {
	output, err := h.chatMessageWithJobSeekerInteractor.SendChatMessageWithJobSeekerLineMessage(interactor.SendChatMessageWithJobSeekerLineMessageInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithJobSeekerJSONPresenter(responses.NewChatMessageWithJobSeeker(output.ChatMessageWithJobSeeker)), nil
}

func (h *ChatMessageWithJobSeekerHandlerImpl) SendChatMessageWithJobSeekerLineImage(param entity.SendChatMessageWithJobSeekerLineImageParam) (presenter.Presenter, error) {
	output, err := h.chatMessageWithJobSeekerInteractor.SendChatMessageWithJobSeekerLineImage(interactor.SendChatMessageWithJobSeekerLineImageInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithJobSeekerJSONPresenter(responses.NewChatMessageWithJobSeeker(output.ChatMessageWithJobSeeker)), nil
}

// グループIDからチャットメッセージを取得する
func (h *ChatMessageWithJobSeekerHandlerImpl) GetChatMessageWithJobSeekerListByGroupID(groupID uint) (presenter.Presenter, error) {
	output, err := h.chatMessageWithJobSeekerInteractor.GetChatMessageWithJobSeekerListByGroupID(interactor.GetChatMessageWithJobSeekerListByGroupIDInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatMessageWithJobSeekerListJSONPresenter(responses.NewChatMessageWithJobSeekerList(output.ChatMessageWithJobSeekerList)), nil
}

/****************************************************************************************/
/// LINE WebHook
//
func (h *ChatMessageWithJobSeekerHandlerImpl) LineWebHook(req *http.Request) (presenter.Presenter, error) {
	output, err := h.chatMessageWithJobSeekerInteractor.LineWebHook(interactor.LineWebHookInput{
		Request: req,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API

/****************************************************************************************/
