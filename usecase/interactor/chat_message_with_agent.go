package interactor

import (
	"fmt"
	"os"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatMessageWithAgentInteractor interface {
	// 汎用系 API
	SendChatMessageWithAgent(input SendChatMessageWithAgentInput) (SendChatMessageWithAgentOutput, error)
	GetChatMessageWithAgentListByThreadID(input GetChatMessageWithAgentListByThreadIDInput) (GetChatMessageWithAgentListByThreadIDOutput, error)

	// 最終閲覧時間の更新
	UpdateChatMessageWithAgentWatchedAtByThreadID(input UpdateChatMessageWithAgentWatchedAtByThreadIDInput) (UpdateChatMessageWithAgentWatchedAtByThreadIDOutput, error)
}

type ChatMessageWithAgentInteractorImpl struct {
	firebase                             usecase.Firebase
	sendgrid                             config.Sendgrid
	oneSignal                            config.OneSignal
	chatMessageWithAgentRepository       usecase.ChatMessageWithAgentRepository
	chatGroupWithAgentRepository         usecase.ChatGroupWithAgentRepository
	chatThreadWithAgentRepository        usecase.ChatThreadWithAgentRepository
	chatMessageToUserWithAgentRepository usecase.ChatMessageToUserWithAgentRepository
	agentStaffRepository                 usecase.AgentStaffRepository
}

// ChatMessageWithAgentInteractorImpl is an implementation of ChatMessageWithAgentInteractor
func NewChatMessageWithAgentInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	cmA usecase.ChatMessageWithAgentRepository,
	cgR usecase.ChatGroupWithAgentRepository,
	ctR usecase.ChatThreadWithAgentRepository,
	cmtuR usecase.ChatMessageToUserWithAgentRepository,
	asR usecase.AgentStaffRepository,
) ChatMessageWithAgentInteractor {
	return &ChatMessageWithAgentInteractorImpl{
		firebase:                             fb,
		sendgrid:                             sg,
		oneSignal:                            os,
		chatMessageWithAgentRepository:       cmA,
		chatGroupWithAgentRepository:         cgR,
		chatThreadWithAgentRepository:        ctR,
		chatMessageToUserWithAgentRepository: cmtuR,
		agentStaffRepository:                 asR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// エージェントがエージェントにLINE送信
type SendChatMessageWithAgentInput struct {
	SendParam entity.SendChatMessageWithAgentParam
}

type SendChatMessageWithAgentOutput struct {
	ChatMessageWithAgent *entity.ChatMessageWithAgent
}

func (i *ChatMessageWithAgentInteractorImpl) SendChatMessageWithAgent(input SendChatMessageWithAgentInput) (SendChatMessageWithAgentOutput, error) {
	var (
		output SendChatMessageWithAgentOutput
		err    error
	)

	// エージェントからエージェントへのチャットメッセージを作成
	chatMessage := entity.NewChatMessageWithAgent(
		input.SendParam.ThreadID,
		input.SendParam.AgentStaffID,
		input.SendParam.Message,
		input.SendParam.FileURL,
	)

	err = i.chatMessageWithAgentRepository.Create(chatMessage)
	if err != nil {
		return output, err
	}

	redirectURL := os.Getenv("BASE_DOMAIN") + "/chat/agent/"
	topic := "ChatAgent"

	for _, toAgentStaff := range input.SendParam.ToAgentStaffs {
		messageToUser := entity.ChatMessageToUserWithAgent{
			MessageID:    chatMessage.ID,
			AgentStaffID: toAgentStaff.ID,
		}

		err = i.chatMessageToUserWithAgentRepository.Create(&messageToUser)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		err = utility.WebPush(
			i.oneSignal.AppID,
			i.oneSignal.APIKey,
			toAgentStaff.FirebaseID,
			"チャット通知",
			"エージェントチャットに、あなた宛の新着メッセージがあります",
			topic,
			redirectURL,
		)
		if err != nil {
			fmt.Println("WebPushの通知でエラー")
			fmt.Println(err)
		}
	}

	// グループのLastSendAtを更新
	err = i.chatGroupWithAgentRepository.UpdateLastSendAtByThreadID(input.SendParam.ThreadID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// // エージェントチャットのプッシュ通知
	// if len(input.SendParam.ToAgentStaffIDs) > 0 {
	// 	redirectURL := os.Getenv("BASE_DOMAIN") + "/chat/agent/"
	// 	topic := "ChatAgent"

	// 	staffIDList := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(input.SendParam.ToAgentStaffIDs)), ", "), "[]")

	// 	staffList, err := i.agentStaffRepository.GetListByStaffIDList(staffIDList)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return output, err
	// 	}

	// 	var firebaseIDList []string

	// 	for _, staff := range staffList {
	// 		firebaseIDList = append(firebaseIDList, staff.FirebaseID)
	// 	}

	// 	firebaseIDList := strings.Join(firebaseIDList, ", ")

	// 	fmt.Println("ID")
	// 	fmt.Println(firebaseIDList)

	// 	err = utility.WebPush(
	// 		i.oneSignal.AppID,
	// 		i.oneSignal.APIKey,
	// 		firebaseIDList,
	// 		"チャット通知",
	// 		"エージェントチャットに、あなた宛の新着メッセージがあります",
	// 		topic,
	// 		redirectURL,
	// 	)
	// 	if err != nil {
	// 		fmt.Println("WebPushの通知でエラー")
	// 		fmt.Println(err)
	// 	}
	// }

	output.ChatMessageWithAgent = chatMessage

	return output, nil
}

// チャットグループIDからチャットメッセージ情報を取得する
type GetChatMessageWithAgentListByThreadIDInput struct {
	ThreadID     uint
	AgentStaffID uint
}

type GetChatMessageWithAgentListByThreadIDOutput struct {
	ChatMessageWithAgentList []*entity.ChatMessageWithAgent
}

func (i *ChatMessageWithAgentInteractorImpl) GetChatMessageWithAgentListByThreadID(input GetChatMessageWithAgentListByThreadIDInput) (GetChatMessageWithAgentListByThreadIDOutput, error) {
	var (
		output GetChatMessageWithAgentListByThreadIDOutput
		err    error
	)

	chatMessageList, err := i.chatMessageWithAgentRepository.GetByThreadID(input.ThreadID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ToUserのセット
	var (
		idListUint []uint
	)

	for _, message := range chatMessageList {
		idListUint = append(idListUint, message.ID)
	}

	toUsers, err := i.chatMessageToUserWithAgentRepository.GetByMessageIDList(idListUint)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, message := range chatMessageList {
		for _, toUser := range toUsers {
			if message.ID == toUser.MessageID {
				message.MessageToUsers = append(message.MessageToUsers, *toUser)

				// 自分宛のメッセージが既読かを判定
				if toUser.AgentStaffID == input.AgentStaffID && toUser.WatchedAt.Unix() <= toUser.SendAt.Unix() {
					message.NotWatched = true
				} else if toUser.AgentStaffID == input.AgentStaffID && toUser.WatchedAt.Unix() > toUser.SendAt.Unix() {
					message.NotWatched = false
				}
			}
		}
	}

	output.ChatMessageWithAgentList = chatMessageList

	return output, nil
}

// 最終閲覧時間を更新する
type UpdateChatMessageWithAgentWatchedAtByThreadIDInput struct {
	ThreadID     uint
	AgentStaffID uint
}

type UpdateChatMessageWithAgentWatchedAtByThreadIDOutput struct {
	OK bool
}

func (i *ChatMessageWithAgentInteractorImpl) UpdateChatMessageWithAgentWatchedAtByThreadID(input UpdateChatMessageWithAgentWatchedAtByThreadIDInput) (UpdateChatMessageWithAgentWatchedAtByThreadIDOutput, error) {
	var (
		output UpdateChatMessageWithAgentWatchedAtByThreadIDOutput
		err    error
	)

	err = i.chatMessageToUserWithAgentRepository.UpdateWatchedAtByThreadID(input.ThreadID, input.AgentStaffID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}
