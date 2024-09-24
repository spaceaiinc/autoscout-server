package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type EmailWithSendingJobSeekerHandler interface {
	// 汎用系 API
	SendEmailWithSendingJobSeeker(param entity.SendEmailWithSendingJobSeekerParam) (presenter.Presenter, error)
	GetEmailWithSendingJobSeekerListBySendingJobSeekerID(jobSeekerID uint) (presenter.Presenter, error)

	// Admin API
}

type EmailWithSendingJobSeekerHandlerImpl struct {
	emailWithSendingJobSeekerInteractor interactor.EmailWithSendingJobSeekerInteractor
}

func NewEmailWithSendingJobSeekerHandlerImpl(cmI interactor.EmailWithSendingJobSeekerInteractor) EmailWithSendingJobSeekerHandler {
	return &EmailWithSendingJobSeekerHandlerImpl{
		emailWithSendingJobSeekerInteractor: cmI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
// メール送信
func (h *EmailWithSendingJobSeekerHandlerImpl) SendEmailWithSendingJobSeeker(param entity.SendEmailWithSendingJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.emailWithSendingJobSeekerInteractor.SendEmailWithSendingJobSeeker(interactor.SendEmailWithSendingJobSeekerInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEmailWithSendingJobSeekerJSONPresenter(responses.NewEmailWithSendingJobSeeker(output.EmailWithSendingJobSeeker)), nil
}

// グループIDからチャットメッセージを取得する
func (h *EmailWithSendingJobSeekerHandlerImpl) GetEmailWithSendingJobSeekerListBySendingJobSeekerID(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.emailWithSendingJobSeekerInteractor.GetEmailWithSendingJobSeekerListBySendingJobSeekerID(interactor.GetEmailWithSendingJobSeekerListBySendingJobSeekerIDInput{
		SendingJobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEmailWithSendingJobSeekerListJSONPresenter(responses.NewEmailWithSendingJobSeekerList(output.EmailWithSendingJobSeekerList)), nil
}

/****************************************************************************************/
/// Admin API

/****************************************************************************************/
