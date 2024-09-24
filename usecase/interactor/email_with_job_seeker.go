package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EmailWithJobSeekerInteractor interface {
	// 汎用系 API
	SendEmailWithJobSeeker(input SendEmailWithJobSeekerInput) (SendEmailWithJobSeekerOutput, error)
	GetEmailWithJobSeekerListByJobSeekerID(input GetEmailWithJobSeekerListByJobSeekerIDInput) (GetEmailWithJobSeekerListByJobSeekerIDOutput, error)
}

type EmailWithJobSeekerInteractorImpl struct {
	firebase                         usecase.Firebase
	sendgrid                         config.Sendgrid
	oneSignal                        config.OneSignal
	emailWithJobSeekerRepository     usecase.EmailWithJobSeekerRepository
	chatGroupWithJobSeekerRepository usecase.ChatGroupWithJobSeekerRepository
	agentStaffRepository             usecase.AgentStaffRepository
}

// EmailWithJobSeekerInteractorImpl is an implementation of EmailWithJobSeekerInteractor
func NewEmailWithJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	ewjsR usecase.EmailWithJobSeekerRepository,
	cgR usecase.ChatGroupWithJobSeekerRepository,
	asR usecase.AgentStaffRepository,
) EmailWithJobSeekerInteractor {
	return &EmailWithJobSeekerInteractorImpl{
		firebase:                         fb,
		sendgrid:                         sg,
		oneSignal:                        os,
		emailWithJobSeekerRepository:     ewjsR,
		chatGroupWithJobSeekerRepository: cgR,
		agentStaffRepository:             asR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// エージェントが求職者にメール送信&保存する
type SendEmailWithJobSeekerInput struct {
	SendParam entity.SendEmailWithJobSeekerParam
}

type SendEmailWithJobSeekerOutput struct {
	EmailWithJobSeeker *entity.EmailWithJobSeeker
}

func (i *EmailWithJobSeekerInteractorImpl) SendEmailWithJobSeeker(input SendEmailWithJobSeekerInput) (SendEmailWithJobSeekerOutput, error) {
	var (
		output SendEmailWithJobSeekerOutput
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
	err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントの最終閲覧時間を更新
	err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// メールを保存
	emailWithJobSeeker := entity.NewEmailWithJobSeeker(
		input.SendParam.JobSeekerID,
		input.SendParam.Subject,
		input.SendParam.Content,
		input.SendParam.FileName,
	)

	err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.EmailWithJobSeeker = emailWithJobSeeker

	return output, nil
}

// 求職者IDからメール一覧を取得する
type GetEmailWithJobSeekerListByJobSeekerIDInput struct {
	JobSeekerID uint
}

type GetEmailWithJobSeekerListByJobSeekerIDOutput struct {
	EmailWithJobSeekerList []*entity.EmailWithJobSeeker
}

func (i *EmailWithJobSeekerInteractorImpl) GetEmailWithJobSeekerListByJobSeekerID(input GetEmailWithJobSeekerListByJobSeekerIDInput) (GetEmailWithJobSeekerListByJobSeekerIDOutput, error) {
	var (
		output GetEmailWithJobSeekerListByJobSeekerIDOutput
		err    error
	)

	emailList, err := i.emailWithJobSeekerRepository.GetByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	fmt.Println("#######")
	fmt.Println(emailList)
	fmt.Println("#######")

	output.EmailWithJobSeekerList = emailList

	return output, nil
}

/****************************************************************************************/
