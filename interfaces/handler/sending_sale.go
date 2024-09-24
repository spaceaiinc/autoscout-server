package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingSaleHandler interface {
	// 汎用系 API
	CreateSendingSale(param entity.CreateSendingSaleParam) (presenter.Presenter, error)
	UpdateSendingSale(sendingSaleID uint, param entity.UpdateSendingSaleParam) (presenter.Presenter, error)
	GetSendingSaleByID(sendingSaleID uint) (presenter.Presenter, error)
	GetSendingSaleByJobSeekerIDAndEnterpriseID(sendingJobSeekerID, sendingEnterpriseID uint) (presenter.Presenter, error)
	GetSendingSaleListBySenderAgentIDForMonthly(senderAgentID uint, startMonth, endMonth string) (presenter.Presenter, error)
	GetSendingSaleListByAgentIDForMonthly(agentID uint, startMonth, endMonth string) (presenter.Presenter, error)

	// Admin API
	GetAllSendingSaleForMonthly(startMonth, endMonth string) (presenter.Presenter, error)
}

type SendingSaleHandlerImpl struct {
	sendingSaleInteractor interactor.SendingSaleInteractor
}

func NewSendingSaleHandlerImpl(itI interactor.SendingSaleInteractor) SendingSaleHandler {
	return &SendingSaleHandlerImpl{
		sendingSaleInteractor: itI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *SendingSaleHandlerImpl) CreateSendingSale(param entity.CreateSendingSaleParam) (presenter.Presenter, error) {
	output, err := h.sendingSaleInteractor.CreateSendingSale(interactor.CreateSendingSaleInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingSaleHandlerImpl) UpdateSendingSale(sendingSaleID uint, param entity.UpdateSendingSaleParam) (presenter.Presenter, error) {
	output, err := h.sendingSaleInteractor.UpdateSendingSale(interactor.UpdateSendingSaleInput{
		SendingSaleID: sendingSaleID,
		UpdateParam:   param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SendingSaleHandlerImpl) GetSendingSaleByID(sendingSaleID uint) (presenter.Presenter, error) {
	output, err := h.sendingSaleInteractor.GetSendingSaleByID(interactor.GetSendingSaleByIDInput{
		SendingSaleID: sendingSaleID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingSaleJSONPresenter(responses.NewSendingSale(output.SendingSale)), nil
}

func (h *SendingSaleHandlerImpl) GetSendingSaleByJobSeekerIDAndEnterpriseID(sendingJobSeekerID, sendingEnterpriseID uint) (presenter.Presenter, error) {
	output, err := h.sendingSaleInteractor.GetSendingSaleByJobSeekerIDAndEnterpriseID(interactor.GetSendingSaleByJobSeekerIDAndEnterpriseIDInput{
		SendingJobSeekerID:  sendingJobSeekerID,
		SendingEnterpriseID: sendingEnterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingSaleJSONPresenter(responses.NewSendingSale(output.SendingSale)), nil
}

func (h *SendingSaleHandlerImpl) GetSendingSaleListBySenderAgentIDForMonthly(senderAgentID uint, startMonth, endMonth string) (presenter.Presenter, error) {
	output, err := h.sendingSaleInteractor.GetSendingSaleListBySenderAgentIDForMonthly(interactor.GetSendingSaleListBySenderAgentIDForMonthlyInput{
		SenderAgentID: senderAgentID,
		StartMonth:    startMonth,
		EndMonth:      endMonth,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingSaleListJSONPresenter(responses.NewSendingSaleList(output.SendingSaleList)), nil
}

func (h *SendingSaleHandlerImpl) GetSendingSaleListByAgentIDForMonthly(agentID uint, startMonth, endMonth string) (presenter.Presenter, error) {
	output, err := h.sendingSaleInteractor.GetSendingSaleListByAgentIDForMonthly(interactor.GetSendingSaleListByAgentIDForMonthlyInput{
		AgentID:    agentID,
		StartMonth: startMonth,
		EndMonth:   endMonth,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingSaleListJSONPresenter(responses.NewSendingSaleList(output.SendingSaleList)), nil
}

/****************************************************************************************/
/// Admin API
//
func (h *SendingSaleHandlerImpl) GetAllSendingSaleForMonthly(startMonth, endMonth string) (presenter.Presenter, error) {
	output, err := h.sendingSaleInteractor.GetAllSendingSaleForMonthly(interactor.GetAllSendingSaleForMonthlyInput{
		StartMonth: startMonth,
		EndMonth:   endMonth,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingSaleListJSONPresenter(responses.NewSendingSaleList(output.SendingSaleList)), nil
}
