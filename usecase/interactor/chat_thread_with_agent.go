package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatThreadWithAgentInteractor interface {
	// 汎用系 API
	CreateChatThreadWithAgent(input CreateChatThreadWithAgentInput) (CreateChatThreadWithAgentOutput, error)
	// GetChatGroupAndThreadWithAgentByAgentID(input GetChatGroupAndThreadWithAgentByAgentIDInput) (GetChatGroupAndThreadWithAgentByAgentIDOutput, error)
	GetChatThreadAndMessageWithAgentByChatThreadID(input GetChatThreadAndMessageWithAgentByChatThreadIDInput) (GetChatThreadAndMessageWithAgentByChatThreadIDOutput, error)
}

type ChatThreadWithAgentInteractorImpl struct {
	firebase                             usecase.Firebase
	sendgrid                             config.Sendgrid
	oneSignal                            config.OneSignal
	chatGroupWithAgentRepository         usecase.ChatGroupWithAgentRepository
	chatThreadWithAgentRepository        usecase.ChatThreadWithAgentRepository
	chatMessageWithAgentRepository       usecase.ChatMessageWithAgentRepository
	chatMessageToUserWithAgentRepository usecase.ChatMessageToUserWithAgentRepository
	agentAllianceRepository              usecase.AgentAllianceRepository
}

// ChatThreadWithAgentInteractorImpl is an implementation of ChatThreadWithAgentInteractor
func NewChatThreadWithAgentInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	cgwaR usecase.ChatGroupWithAgentRepository,
	ctwaR usecase.ChatThreadWithAgentRepository,
	cmwaR usecase.ChatMessageWithAgentRepository,
	cmtuwaR usecase.ChatMessageToUserWithAgentRepository,
	aaR usecase.AgentAllianceRepository,
) ChatThreadWithAgentInteractor {
	return &ChatThreadWithAgentInteractorImpl{
		firebase:                             fb,
		sendgrid:                             sg,
		oneSignal:                            os,
		chatGroupWithAgentRepository:         cgwaR,
		chatThreadWithAgentRepository:        ctwaR,
		chatMessageWithAgentRepository:       cmwaR,
		chatMessageToUserWithAgentRepository: cmtuwaR,
		agentAllianceRepository:              aaR,
	}
}

/****************************************************************************************/
/// 汎用系API
//エージェントと求職者のチャットグループの作成
type CreateChatThreadWithAgentInput struct {
	CreateParam entity.CreateChatThreadWithAgentParam
}

type CreateChatThreadWithAgentOutput struct {
	ChatThreadWithAgent *entity.ChatThreadWithAgent
}

func (i *ChatThreadWithAgentInteractorImpl) CreateChatThreadWithAgent(input CreateChatThreadWithAgentInput) (CreateChatThreadWithAgentOutput, error) {
	var (
		output CreateChatThreadWithAgentOutput
		err    error
	)

	// エージェントと求職者のチャットグループを作成
	chatThread := entity.NewChatThreadWithAgent(
		input.CreateParam.GroupID,
		input.CreateParam.AgentStaffID,
		input.CreateParam.Title,
	)

	err = i.chatThreadWithAgentRepository.Create(chatThread)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatThreadWithAgent = chatThread

	return output, nil
}

// チャットグループIDからチャットメッセージ情報を取得する
type GetChatThreadAndMessageWithAgentByChatThreadIDInput struct {
	ThreadID uint
}

type GetChatThreadAndMessageWithAgentByChatThreadIDOutput struct {
	ChatThreadWithAgent *entity.ChatThreadWithAgent
}

func (i *ChatThreadWithAgentInteractorImpl) GetChatThreadAndMessageWithAgentByChatThreadID(input GetChatThreadAndMessageWithAgentByChatThreadIDInput) (GetChatThreadAndMessageWithAgentByChatThreadIDOutput, error) {
	var (
		output GetChatThreadAndMessageWithAgentByChatThreadIDOutput
		err    error
	)

	chatThread, err := i.chatThreadWithAgentRepository.FindByID(input.ThreadID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	chatMessageList, err := i.chatMessageWithAgentRepository.GetByThreadID(input.ThreadID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// chatMessageToUserWithAgentList, err := i.chatMessageToUserWithAgentRepository.GetByThreadID(input.ThreadID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// for _, chatMessageToUserWithAgent := range chatMessageToUserWithAgentList {
	// 	if chatMessage.ID == chatMessageToUserWithAgent.MessageID {
	// 		// 最終閲覧時間が0の場合は未読
	// 		if chatMessageToUserWithAgent.WatchedAt.Unix() <= chatMessageToUserWithAgent.SendAt.Unix() {
	// 			fmt.Println("未読:", chatMessageToUserWithAgent.WatchedAt.Unix(), "<=", chatMessageToUserWithAgent.SendAt.Unix())
	// 			chatMessageToUserWithAgent.IsWatched = false
	// 		} else {
	// 			fmt.Println("既読:", chatMessageToUserWithAgent.WatchedAt.Unix(), "<=", chatMessageToUserWithAgent.SendAt.Unix())
	// 			chatMessageToUserWithAgent.IsWatched = true
	// 		}

	// 		chatMessage.MessageToUsers = append(chatMessage.MessageToUsers, *chatMessageToUserWithAgent)
	// 	}
	// }

	for _, chatMessage := range chatMessageList {
		if chatThread.ID == chatMessage.ThreadID {
			chatThread.Messages = append(chatThread.Messages, *chatMessage)
		}
	}

	output.ChatThreadWithAgent = chatThread

	return output, nil
}
