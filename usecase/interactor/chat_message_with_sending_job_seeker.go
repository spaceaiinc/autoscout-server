package interactor

import (
	"fmt"

	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type ChatMessageWithSendingJobSeekerInteractor interface {
	// 汎用系 API
	SendChatMessageWithSendingJobSeekerLineMessage(input SendChatMessageWithSendingJobSeekerLineMessageInput) (SendChatMessageWithSendingJobSeekerLineMessageOutput, error)
	SendChatMessageWithSendingJobSeekerLineImage(input SendChatMessageWithSendingJobSeekerLineImageInput) (SendChatMessageWithSendingJobSeekerLineImageOutput, error)
	GetChatMessageWithSendingJobSeekerListByGroupID(input GetChatMessageWithSendingJobSeekerListByGroupIDInput) (GetChatMessageWithSendingJobSeekerListByGroupIDOutput, error)
}

type ChatMessageWithSendingJobSeekerInteractorImpl struct {
	firebase                                  usecase.Firebase
	sendgrid                                  config.Sendgrid
	oneSignal                                 config.OneSignal
	chatMessageWithSendingJobSeekerRepository usecase.ChatMessageWithSendingJobSeekerRepository
	chatGroupWithSendingJobSeekerRepository   usecase.ChatGroupWithSendingJobSeekerRepository
	jobSeekerRepository                       usecase.SendingJobSeekerRepository
	agentRepository                           usecase.AgentRepository
	agentStaffRepository                      usecase.AgentStaffRepository
}

// ChatMessageWithSendingJobSeekerInteractorImpl is an implementation of ChatMessageWithSendingJobSeekerInteractor
func NewChatMessageWithSendingJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	cmR usecase.ChatMessageWithSendingJobSeekerRepository,
	cgR usecase.ChatGroupWithSendingJobSeekerRepository,
	jsR usecase.SendingJobSeekerRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
) ChatMessageWithSendingJobSeekerInteractor {
	return &ChatMessageWithSendingJobSeekerInteractorImpl{
		firebase:  fb,
		sendgrid:  sg,
		oneSignal: os,
		chatMessageWithSendingJobSeekerRepository: cmR,
		chatGroupWithSendingJobSeekerRepository:   cgR,
		jobSeekerRepository:                       jsR,
		agentRepository:                           aR,
		agentStaffRepository:                      asR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// エージェントが求職者にLINE送信
type SendChatMessageWithSendingJobSeekerLineMessageInput struct {
	SendParam entity.SendChatMessageWithSendingJobSeekerLineParam
}

type SendChatMessageWithSendingJobSeekerLineMessageOutput struct {
	ChatMessageWithSendingJobSeeker *entity.ChatMessageWithSendingJobSeeker
}

func (i *ChatMessageWithSendingJobSeekerInteractorImpl) SendChatMessageWithSendingJobSeekerLineMessage(input SendChatMessageWithSendingJobSeekerLineMessageInput) (SendChatMessageWithSendingJobSeekerLineMessageOutput, error) {
	var (
		output SendChatMessageWithSendingJobSeekerLineMessageOutput
		err    error
	)

	// 送信元がエージェントかつ、テキストメッセージの場合、LINEに送信
	agent, err := i.agentRepository.FindByID(1) // SystemのID
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// シークレットとトークンを設定
	bot, err := linebot.New(agent.LineMessagingChannelSecret, agent.LineMessagingChannelAccessToken)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	chatGroup, err := i.chatGroupWithSendingJobSeekerRepository.FindByID(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobSeeker, err := i.jobSeekerRepository.FindByID(chatGroup.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// LINE メッセージを送信する
	call, err := bot.PushMessage(
		jobSeeker.LineID,
		linebot.NewTextMessage(input.SendParam.Message),
	).Do()
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	fmt.Println("メッセージの送信に成功しました:", call)

	// エージェントから求職者へのチャットメッセージを作成
	chatMessage := entity.NewChatMessageWithSendingJobSeeker(
		input.SendParam.GroupID,
		null.NewInt(0, true), // エージェント
		input.SendParam.Message,
		"",
		"",
		"",
		"",
		null.NewInt(0, false),
		"",
		null.NewInt(0, true),
		time.Now().In(time.UTC),
	)

	err = i.chatMessageWithSendingJobSeekerRepository.Create(chatMessage)
	if err != nil {
		return output, err
	}

	// エージェントの最終送信時間を更新
	err = i.chatGroupWithSendingJobSeekerRepository.UpdateAgentLastSendAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatMessageWithSendingJobSeeker = chatMessage

	return output, nil
}

type SendChatMessageWithSendingJobSeekerLineImageInput struct {
	SendParam entity.SendChatMessageWithSendingJobSeekerLineImageParam
}

type SendChatMessageWithSendingJobSeekerLineImageOutput struct {
	ChatMessageWithSendingJobSeeker *entity.ChatMessageWithSendingJobSeeker
}

func (i *ChatMessageWithSendingJobSeekerInteractorImpl) SendChatMessageWithSendingJobSeekerLineImage(input SendChatMessageWithSendingJobSeekerLineImageInput) (SendChatMessageWithSendingJobSeekerLineImageOutput, error) {
	var (
		output SendChatMessageWithSendingJobSeekerLineImageOutput
		err    error
	)

	// エージェントから求職者への画像送信（LINE）
	agent, err := i.agentRepository.FindByID(1) // SystemのID
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	url, err := i.firebase.UploadToStorageForAgentLine(input.SendParam.File, agent.UUID.String())
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	chatGroup, err := i.chatGroupWithSendingJobSeekerRepository.FindByID(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobSeeker, err := i.jobSeekerRepository.FindByID(chatGroup.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// シークレットとトークンを設定
	bot, err := linebot.New(agent.LineMessagingChannelSecret, agent.LineMessagingChannelAccessToken)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// LINE メッセージを送信する
	call, err := bot.PushMessage(
		jobSeeker.LineID,
		linebot.NewImageMessage(url, url),
	).Do()
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	fmt.Println("メッセージの送信に成功しました:", call)

	// エージェントから求職者へのチャットメッセージを作成
	chatMessage := entity.NewChatMessageWithSendingJobSeeker(
		input.SendParam.GroupID,
		null.NewInt(0, true), // エージェント
		"",
		"",
		"",
		url,
		url,
		null.NewInt(0, false),
		"",
		null.NewInt(2, true),
		time.Now().In(time.UTC),
	)

	err = i.chatMessageWithSendingJobSeekerRepository.Create(chatMessage)
	if err != nil {
		return output, err
	}

	// エージェントの最終送信時間を更新
	err = i.chatGroupWithSendingJobSeekerRepository.UpdateAgentLastSendAt(input.SendParam.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatMessageWithSendingJobSeeker = chatMessage

	return output, nil
}

// チャットグループIDからチャットメッセージ情報を取得する
type GetChatMessageWithSendingJobSeekerListByGroupIDInput struct {
	GroupID uint
}

type GetChatMessageWithSendingJobSeekerListByGroupIDOutput struct {
	ChatMessageWithSendingJobSeekerList []*entity.ChatMessageWithSendingJobSeeker
}

func (i *ChatMessageWithSendingJobSeekerInteractorImpl) GetChatMessageWithSendingJobSeekerListByGroupID(input GetChatMessageWithSendingJobSeekerListByGroupIDInput) (GetChatMessageWithSendingJobSeekerListByGroupIDOutput, error) {
	var (
		output GetChatMessageWithSendingJobSeekerListByGroupIDOutput
		err    error
	)

	chatMessageList, err := i.chatMessageWithSendingJobSeekerRepository.GetListByGroupID(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	fmt.Println("#######")
	fmt.Println(chatMessageList)
	fmt.Println("#######")

	output.ChatMessageWithSendingJobSeekerList = chatMessageList

	return output, nil
}
