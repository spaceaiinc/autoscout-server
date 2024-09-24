package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingJobSeekerDesiredJobInformationHandler interface {
	// 汎用系 API
	CreateMultiSendingJobSeekerDesiredJobInformation(param entity.CreateMultiSendingJobSeekerDesiredJobInformationParam) (presenter.Presenter, error)

	GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerID(sendingJobSeekerID uint) (presenter.Presenter, error)
}

type SendingJobSeekerDesiredJobInformationHandlerImpl struct {
	sendingJobSeekerDesiredJobInformationInteractor interactor.SendingJobSeekerDesiredJobInformationInteractor
}

func NewSendingJobSeekerDesiredJobInformationHandlerImpl(epI interactor.SendingJobSeekerDesiredJobInformationInteractor) SendingJobSeekerDesiredJobInformationHandler {
	return &SendingJobSeekerDesiredJobInformationHandlerImpl{
		sendingJobSeekerDesiredJobInformationInteractor: epI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *SendingJobSeekerDesiredJobInformationHandlerImpl) CreateMultiSendingJobSeekerDesiredJobInformation(param entity.CreateMultiSendingJobSeekerDesiredJobInformationParam) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerDesiredJobInformationInteractor.CreateMultiSendingJobSeekerDesiredJobInformation(interactor.CreateMultiSendingJobSeekerDesiredJobInformationInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerDesiredJobInformationListJSONPresenter(responses.NewSendingJobSeekerDesiredJobInformationList(output.SendingJobSeekerDesiredJobInformationList)), nil
}

// すべての送客進捗情報を取得する
func (h *SendingJobSeekerDesiredJobInformationHandlerImpl) GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerID(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingJobSeekerDesiredJobInformationInteractor.GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerID(interactor.GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerDesiredJobInformationListJSONPresenter(responses.NewSendingJobSeekerDesiredJobInformationList(output.SendingJobSeekerDesiredJobInformationList)), nil
}
