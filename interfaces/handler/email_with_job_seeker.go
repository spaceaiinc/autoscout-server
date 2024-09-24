package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type EmailWithJobSeekerHandler interface {
	// 汎用系 API
	SendEmailWithJobSeeker(param entity.SendEmailWithJobSeekerParam) (presenter.Presenter, error)
	GetEmailWithJobSeekerListByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error)

	// Admin API
}

type EmailWithJobSeekerHandlerImpl struct {
	emailWithJobSeekerInteractor interactor.EmailWithJobSeekerInteractor
}

func NewEmailWithJobSeekerHandlerImpl(cmI interactor.EmailWithJobSeekerInteractor) EmailWithJobSeekerHandler {
	return &EmailWithJobSeekerHandlerImpl{
		emailWithJobSeekerInteractor: cmI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
// メール送信
func (h *EmailWithJobSeekerHandlerImpl) SendEmailWithJobSeeker(param entity.SendEmailWithJobSeekerParam) (presenter.Presenter, error) {
	output, err := h.emailWithJobSeekerInteractor.SendEmailWithJobSeeker(interactor.SendEmailWithJobSeekerInput{
		SendParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEmailWithJobSeekerJSONPresenter(responses.NewEmailWithJobSeeker(output.EmailWithJobSeeker)), nil
}

// グループIDからチャットメッセージを取得する
func (h *EmailWithJobSeekerHandlerImpl) GetEmailWithJobSeekerListByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.emailWithJobSeekerInteractor.GetEmailWithJobSeekerListByJobSeekerID(interactor.GetEmailWithJobSeekerListByJobSeekerIDInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewEmailWithJobSeekerListJSONPresenter(responses.NewEmailWithJobSeekerList(output.EmailWithJobSeekerList)), nil
}

/****************************************************************************************/
/// Admin API

/****************************************************************************************/
