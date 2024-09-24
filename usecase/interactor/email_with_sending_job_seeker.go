package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EmailWithSendingJobSeekerInteractor interface {
	// 汎用系 API
	SendEmailWithSendingJobSeeker(input SendEmailWithSendingJobSeekerInput) (SendEmailWithSendingJobSeekerOutput, error)
	GetEmailWithSendingJobSeekerListBySendingJobSeekerID(input GetEmailWithSendingJobSeekerListBySendingJobSeekerIDInput) (GetEmailWithSendingJobSeekerListBySendingJobSeekerIDOutput, error)
}

type EmailWithSendingJobSeekerInteractorImpl struct {
	firebase                                usecase.Firebase
	sendgrid                                config.Sendgrid
	oneSignal                               config.OneSignal
	emailWithSendingJobSeekerRepository     usecase.EmailWithSendingJobSeekerRepository
	chatGroupWithSendingJobSeekerRepository usecase.ChatGroupWithSendingJobSeekerRepository
	agentStaffRepository                    usecase.AgentStaffRepository
}

// EmailWithSendingJobSeekerInteractorImpl is an implementation of EmailWithSendingJobSeekerInteractor
func NewEmailWithSendingJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	ewjsR usecase.EmailWithSendingJobSeekerRepository,
	cgR usecase.ChatGroupWithSendingJobSeekerRepository,
	asR usecase.AgentStaffRepository,
) EmailWithSendingJobSeekerInteractor {
	return &EmailWithSendingJobSeekerInteractorImpl{
		firebase:                                fb,
		sendgrid:                                sg,
		oneSignal:                               os,
		emailWithSendingJobSeekerRepository:     ewjsR,
		chatGroupWithSendingJobSeekerRepository: cgR,
		agentStaffRepository:                    asR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// エージェントが送客求職者にメール送信&保存する
type SendEmailWithSendingJobSeekerInput struct {
	SendParam entity.SendEmailWithSendingJobSeekerParam
}

type SendEmailWithSendingJobSeekerOutput struct {
	EmailWithSendingJobSeeker *entity.EmailWithSendingJobSeeker
}

func (i *EmailWithSendingJobSeekerInteractorImpl) SendEmailWithSendingJobSeeker(input SendEmailWithSendingJobSeekerInput) (SendEmailWithSendingJobSeekerOutput, error) {
	var (
		output SendEmailWithSendingJobSeekerOutput
		err    error
	)

	// 担当者情報を取得
	agentStaff, err := i.agentStaffRepository.FindByID(input.SendParam.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	from := entity.EmailUser{
		Name:  agentStaff.StaffName,
		Email: agentStaff.Email,
	}

	// メール送信
	err = utility.SendMailToSingle(
		i.sendgrid.APIKey,
		input.SendParam.Subject,
		input.SendParam.Content,
		from,
		input.SendParam.To,
		input.SendParam.Files,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントの最終送信時間を更新
	err = i.chatGroupWithSendingJobSeekerRepository.UpdateAgentLastSendAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントの最終閲覧時間を更新
	err = i.chatGroupWithSendingJobSeekerRepository.UpdateAgentLastWatchedAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// メールを保存
	emailWithSendingJobSeeker := entity.NewEmailWithSendingJobSeeker(
		input.SendParam.SendingJobSeekerID,
		input.SendParam.Subject,
		input.SendParam.Content,
		input.SendParam.FileName,
	)

	err = i.emailWithSendingJobSeekerRepository.Create(emailWithSendingJobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.EmailWithSendingJobSeeker = emailWithSendingJobSeeker

	return output, nil
}

// 送客求職者IDからメール一覧を取得する
type GetEmailWithSendingJobSeekerListBySendingJobSeekerIDInput struct {
	SendingJobSeekerID uint
}

type GetEmailWithSendingJobSeekerListBySendingJobSeekerIDOutput struct {
	EmailWithSendingJobSeekerList []*entity.EmailWithSendingJobSeeker
}

func (i *EmailWithSendingJobSeekerInteractorImpl) GetEmailWithSendingJobSeekerListBySendingJobSeekerID(input GetEmailWithSendingJobSeekerListBySendingJobSeekerIDInput) (GetEmailWithSendingJobSeekerListBySendingJobSeekerIDOutput, error) {
	var (
		output GetEmailWithSendingJobSeekerListBySendingJobSeekerIDOutput
		err    error
	)

	emailList, err := i.emailWithSendingJobSeekerRepository.GetListBySendingJobSeekerID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	fmt.Println("#######")
	fmt.Println(emailList)
	fmt.Println("#######")

	output.EmailWithSendingJobSeekerList = emailList

	return output, nil
}

/****************************************************************************************/
