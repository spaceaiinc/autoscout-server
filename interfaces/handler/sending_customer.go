package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingCustomerHandler interface {
	// 汎用系 API
	GetSendingCustomerByID(sendingCustomerID uint) (presenter.Presenter, error)
	GetSendingCustomerListByAgentID(agentID uint) (presenter.Presenter, error)
	GetSearchSendingCustomerListByPageAndTabAndAgentID(pageNumber, agentID, tabNumber uint, freeWord string) (presenter.Presenter, error)
}

type SendingCustomerHandlerImpl struct {
	sendingCustomerInteractor interactor.SendingCustomerInteractor
}

func NewSendingCustomerHandlerImpl(epI interactor.SendingCustomerInteractor) SendingCustomerHandler {
	return &SendingCustomerHandlerImpl{
		sendingCustomerInteractor: epI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *SendingCustomerHandlerImpl) GetSendingCustomerByID(sendingCustomerID uint) (presenter.Presenter, error) {
	output, err := h.sendingCustomerInteractor.GetSendingCustomerByID(interactor.GetSendingCustomerByIDInput{
		SendingCustomerID: sendingCustomerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingCustomerJSONPresenter(responses.NewSendingCustomer(output.SendingCustomer)), nil
}

func (h *SendingCustomerHandlerImpl) GetSendingCustomerListByAgentID(agentID uint) (presenter.Presenter, error) {
	output, err := h.sendingCustomerInteractor.GetSendingCustomerListByAgentID(interactor.GetSendingCustomerListByAgentIDInput{
		AgentID: agentID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingCustomerListJSONPresenter(responses.NewSendingCustomerList(output.SendingCustomerList)), nil
}

func (h *SendingCustomerHandlerImpl) GetSearchSendingCustomerListByPageAndTabAndAgentID(pageNumber, agentID, tabNumber uint, freeWord string) (presenter.Presenter, error) {
	output, err := h.sendingCustomerInteractor.GetSearchSendingCustomerListByPageAndTabAndAgentID(interactor.GetSearchSendingCustomerListByPageAndTabAndAgentIDInput{
		PageNumber: pageNumber,
		AgentID:    agentID,
		TabNumber:  tabNumber,
		FreeWord:   freeWord,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingCustomerListAndMaxPageAndCountJSONPresenter(responses.NewSendingCustomerListAndMaxPageAndCount(output.SendingCustomerList, output.MaxPageNumber, output.AllCount, output.SendingCount, output.CompleteCount, output.CloseCount)), nil
}
