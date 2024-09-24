package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ChatGroupWithSendingJobSeekerHandler interface {
	// 汎用系 API
	CreateChatGroupWithSendingJobSeeker(param entity.CreateChatGroupWithSendingJobSeekerParam) (presenter.Presenter, error)
	UpdateSendingJobSeekerAgentLastWatched(groupID uint) (presenter.Presenter, error)
	GetChatGroupWithSendingJobSeekerByID(id uint) (presenter.Presenter, error)
	GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentID(agentID, pageNumber uint, searchParam entity.SearchChatSendingJobSeeker) (presenter.Presenter, error)
	GetSendingJobSeekerChatNotificationByAgentID(agentID uint) (presenter.Presenter, error)
}

type ChatGroupWithSendingJobSeekerHandlerImpl struct {
	chatGroupWithSendingJobSeekerInteractor interactor.ChatGroupWithSendingJobSeekerInteractor
}

func NewChatGroupWithSendingJobSeekerHandlerImpl(cgI interactor.ChatGroupWithSendingJobSeekerInteractor) ChatGroupWithSendingJobSeekerHandler {
	return &ChatGroupWithSendingJobSeekerHandlerImpl{
		chatGroupWithSendingJobSeekerInteractor: cgI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *ChatGroupWithSendingJobSeekerHandlerImpl) CreateChatGroupWithSendingJobSeeker(param entity.CreateChatGroupWithSendingJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.chatGroupWithSendingJobSeekerInteractor.CreateChatGroupWithSendingJobSeeker(interactor.CreateChatGroupWithSendingJobSeekerInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithSendingJobSeekerJSONPresenter(responses.NewChatGroupWithSendingJobSeeker(output.ChatGroupWithSendingJobSeeker)), nil
}

func (h *ChatGroupWithSendingJobSeekerHandlerImpl) UpdateSendingJobSeekerAgentLastWatched(groupID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithSendingJobSeekerInteractor.UpdateSendingJobSeekerAgentLastWatched(interactor.UpdateSendingJobSeekerAgentLastWatchedInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// IDからチャットグループを取得する
func (h *ChatGroupWithSendingJobSeekerHandlerImpl) GetChatGroupWithSendingJobSeekerByID(id uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithSendingJobSeekerInteractor.GetChatGroupWithSendingJobSeekerByID(interactor.GetChatGroupWithSendingJobSeekerByIDInput{
		ID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithSendingJobSeekerJSONPresenter(responses.NewChatGroupWithSendingJobSeeker(output.ChatGroupWithSendingJobSeeker)), nil
}

func (h *ChatGroupWithSendingJobSeekerHandlerImpl) GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentID(agentID, pageNumber uint, searchParam entity.SearchChatSendingJobSeeker) (presenter.Presenter, error) {
	output, err := h.chatGroupWithSendingJobSeekerInteractor.GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentID(interactor.GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithSendingJobSeekerListAndMaxPageJSONPresenter(responses.NewChatGroupWithSendingJobSeekerListAndMaxPage(output.ChatGroupWithSendingJobSeekerList, output.MaxPageNumber)), nil
}

func (h *ChatGroupWithSendingJobSeekerHandlerImpl) GetSendingJobSeekerChatNotificationByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithSendingJobSeekerInteractor.GetSendingJobSeekerChatNotificationByAgentID(interactor.GetSendingJobSeekerChatNotificationByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithSendingJobSeekerUnwatchedCountJSONPresenter(responses.NewChatGroupWithSendingJobSeekerUnWatchedCount(output.UnwatchedCount)), nil
}

/****************************************************************************************/
