package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDesiredJobInformationInteractor interface {
	// 汎用系 API
	CreateMultiSendingJobSeekerDesiredJobInformation(input CreateMultiSendingJobSeekerDesiredJobInformationInput) (CreateMultiSendingJobSeekerDesiredJobInformationOutput, error)

	GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerID(input GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDInput) (GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDOutput, error)
}

type SendingJobSeekerDesiredJobInformationInteractorImpl struct {
	firebase                                        usecase.Firebase
	sendgrid                                        config.Sendgrid
	sendingJobSeekerDesiredJobInformationRepository usecase.SendingJobSeekerDesiredJobInformationRepository
}

// SendingJobSeekerDesiredJobInformationInteractorImpl is an implementation of SendingJobSeekerDesiredJobInformationInteractor
func NewSendingJobSeekerDesiredJobInformationInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	sjsdjiR usecase.SendingJobSeekerDesiredJobInformationRepository,
) SendingJobSeekerDesiredJobInformationInteractor {
	return &SendingJobSeekerDesiredJobInformationInteractorImpl{
		firebase: fb,
		sendgrid: sg,
		sendingJobSeekerDesiredJobInformationRepository: sjsdjiR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// 興味あり求人の作成
type CreateMultiSendingJobSeekerDesiredJobInformationInput struct {
	CreateParam entity.CreateMultiSendingJobSeekerDesiredJobInformationParam
}

type CreateMultiSendingJobSeekerDesiredJobInformationOutput struct {
	SendingJobSeekerDesiredJobInformationList []*entity.SendingJobSeekerDesiredJobInformation
}

func (i *SendingJobSeekerDesiredJobInformationInteractorImpl) CreateMultiSendingJobSeekerDesiredJobInformation(input CreateMultiSendingJobSeekerDesiredJobInformationInput) (CreateMultiSendingJobSeekerDesiredJobInformationOutput, error) {
	var (
		output CreateMultiSendingJobSeekerDesiredJobInformationOutput
		err    error
	)

	err = i.sendingJobSeekerDesiredJobInformationRepository.DeleteBySendingJobSeekerID(input.CreateParam.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, desiredJobInformationID := range input.CreateParam.SendingJobInformationIDList {

		sendingJobSeekerDesiredJobInformation := entity.NewSendingJobSeekerDesiredJobInformation(
			input.CreateParam.SendingJobSeekerID,
			desiredJobInformationID,
		)

		err = i.sendingJobSeekerDesiredJobInformationRepository.Create(sendingJobSeekerDesiredJobInformation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		output.SendingJobSeekerDesiredJobInformationList = append(output.SendingJobSeekerDesiredJobInformationList, sendingJobSeekerDesiredJobInformation)
	}

	return output, nil
}

// 送客求職者IDを使って興味あり求人情報を取得する
type GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDInput struct {
	SendingJobSeekerID uint
}

type GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDOutput struct {
	SendingJobSeekerDesiredJobInformationList []*entity.SendingJobSeekerDesiredJobInformation
}

func (i *SendingJobSeekerDesiredJobInformationInteractorImpl) GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerID(input GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDInput) (GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDOutput, error) {
	var (
		output GetSendingJobSeekerDesiredJobInformationListBySendingJobSeekerIDOutput
		err    error
	)

	desiredJobInformationList, err := i.sendingJobSeekerDesiredJobInformationRepository.GetListBySendingJobSeekerID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingJobSeekerDesiredJobInformationList = desiredJobInformationList

	return output, nil
}
