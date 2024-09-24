package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ChatGroupWithJobSeekerHandler interface {
	// 汎用系 API
	CreateChatGroupWithJobSeeker(param entity.CreateChatGroupWithJobSeekerParam) (presenter.Presenter, error)
	UpdateAgentLastWatched(groupID uint) (presenter.Presenter, error)
	GetChatGroupWithJobSeekerByID(id uint) (presenter.Presenter, error)
	GetChatGroupWithJobSeekerListByAgentID(agentID uint) (presenter.Presenter, error)
	GetChatGroupWithJobSeekerSearchListAndPageByAgentID(agentID, pageNumber uint, searchParam entity.SearchChatJobSeeker) (presenter.Presenter, error)
	GetJobSeekerChatNotificationByAgentID(agentID uint) (presenter.Presenter, error)

	// Admin API
}

type ChatGroupWithJobSeekerHandlerImpl struct {
	chatGroupWithJobSeekerInteractor interactor.ChatGroupWithJobSeekerInteractor
}

func NewChatGroupWithJobSeekerHandlerImpl(cgI interactor.ChatGroupWithJobSeekerInteractor) ChatGroupWithJobSeekerHandler {
	return &ChatGroupWithJobSeekerHandlerImpl{
		chatGroupWithJobSeekerInteractor: cgI,
	}
}

/****************************************************************************************/
// 汎用系 API
//
func (h *ChatGroupWithJobSeekerHandlerImpl) CreateChatGroupWithJobSeeker(param entity.CreateChatGroupWithJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.chatGroupWithJobSeekerInteractor.CreateChatGroupWithJobSeeker(interactor.CreateChatGroupWithJobSeekerInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithJobSeekerJSONPresenter(responses.NewChatGroupWithJobSeeker(output.ChatGroupWithJobSeeker)), nil
}

func (h *ChatGroupWithJobSeekerHandlerImpl) UpdateAgentLastWatched(groupID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithJobSeekerInteractor.UpdateAgentLastWatched(interactor.UpdateAgentLastWatchedInput{
		GroupID: groupID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// IDからチャットグループを取得する
func (h *ChatGroupWithJobSeekerHandlerImpl) GetChatGroupWithJobSeekerByID(id uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithJobSeekerInteractor.GetChatGroupWithJobSeekerByID(interactor.GetChatGroupWithJobSeekerByIDInput{
		ID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithJobSeekerJSONPresenter(responses.NewChatGroupWithJobSeeker(output.ChatGroupWithJobSeeker)), nil
}

// グループIDからチャットグループを取得する
func (h *ChatGroupWithJobSeekerHandlerImpl) GetChatGroupWithJobSeekerListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithJobSeekerInteractor.GetChatGroupWithJobSeekerListByAgentID(interactor.GetChatGroupWithJobSeekerListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithJobSeekerListJSONPresenter(responses.NewChatGroupWithJobSeekerList(output.ChatGroupWithJobSeekerList)), nil
}

// グループIDからチャットグループを取得する
func (h *ChatGroupWithJobSeekerHandlerImpl) GetChatGroupWithJobSeekerSearchListAndPageByAgentID(agentID, pageNumber uint, searchParam entity.SearchChatJobSeeker) (presenter.Presenter, error) {
	output, err := h.chatGroupWithJobSeekerInteractor.GetChatGroupWithJobSeekerSearchListAndPageByAgentID(interactor.GetChatGroupWithJobSeekerSearchListAndPageByAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithJobSeekerListAndMaxPageJSONPresenter(responses.NewChatGroupWithJobSeekerListAndMaxPage(output.ChatGroupWithJobSeekerList, output.MaxPageNumber)), nil
}

func (h *ChatGroupWithJobSeekerHandlerImpl) GetJobSeekerChatNotificationByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.chatGroupWithJobSeekerInteractor.GetJobSeekerChatNotificationByAgentID(interactor.GetJobSeekerChatNotificationByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewChatGroupWithJobSeekerUnwatchedCountJSONPresenter(responses.NewChatGroupWithJobSeekerUnWatchedCount(output.UnwatchedCount)), nil
}

/****************************************************************************************/
/****************************************************************************************/
// Admin API

/****************************************************************************************/
