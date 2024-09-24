package interactor

import (
	"errors"
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatGroupWithAgentInteractor interface {
	// 汎用系 API
	CreateChatGroupWithAgent(inpdut CreateChatGroupWithAgentInput) (CreateChatGroupWithAgentOutput, error)
	CreateChatGroupWithAgentList(input CreateChatGroupWithAgentListInput) (CreateChatGroupWithAgentListOutput, error)
	UpdateAgentLastWatched(input UpdateAgentLastWatchedForAgentChatInput) (UpdateAgentLastWatchedForAgentChatOutput, error)
	GetChatGroupWithAgentByID(input GetChatGroupWithAgentByIDInput) (GetChatGroupWithAgentByIDOutput, error)
	GetChatGroupWithAgentListByAgentID(input GetChatGroupWithAgentListByAgentIDInput) (GetChatGroupWithAgentListByAgentIDOutput, error)
	GetChatGroupAndThreadWithAgentByChatGroupID(input GetChatGroupAndThreadWithAgentByChatGroupIDInput) (GetChatGroupAndThreadWithAgentByChatGroupIDOutput, error)
	GetChatGroupWithAgentByAgentIDAndOtherAgentID(input GetChatGroupWithAgentByAgentIDAndOtherAgentIDInput) (GetChatGroupWithAgentByAgentIDAndOtherAgentIDOutput, error)

	GetChatGroupWithAgentListByAgentIDAndAgentStaffID(input GetChatGroupWithAgentListByAgentIDAndAgentStaffIDInput) (GetChatGroupWithAgentListByAgentIDAndAgentStaffIDOutput, error)
	CheckAnyChatGroupWithoutGroup(input CheckAnyChatGroupWithoutGroupInput) (CheckAnyChatGroupWithoutGroupOutput, error)
}

type ChatGroupWithAgentInteractorImpl struct {
	firebase                             usecase.Firebase
	sendgrid                             config.Sendgrid
	oneSignal                            config.OneSignal
	chatGroupWithAgentRepository         usecase.ChatGroupWithAgentRepository
	chatThreadWithAgentRepository        usecase.ChatThreadWithAgentRepository
	chatMessageToUserWithAgentRepository usecase.ChatMessageToUserWithAgentRepository
	agentAllianceRepository              usecase.AgentAllianceRepository
}

// ChatGroupWithAgentInteractorImpl is an implementation of ChatGroupWithAgentInteractor
func NewChatGroupWithAgentInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	cgR usecase.ChatGroupWithAgentRepository,
	ctR usecase.ChatThreadWithAgentRepository,
	cmtuR usecase.ChatMessageToUserWithAgentRepository,
	aaR usecase.AgentAllianceRepository,
) ChatGroupWithAgentInteractor {
	return &ChatGroupWithAgentInteractorImpl{
		firebase:                             fb,
		sendgrid:                             sg,
		oneSignal:                            os,
		chatGroupWithAgentRepository:         cgR,
		chatThreadWithAgentRepository:        ctR,
		chatMessageToUserWithAgentRepository: cmtuR,
		agentAllianceRepository:              aaR,
	}
}

/****************************************************************************************/
// 汎用系API
//

// エージェントと求職者のチャットグループの作成
type CreateChatGroupWithAgentInput struct {
	CreateParam entity.CreateChatGroupWithAgentParam
}

type CreateChatGroupWithAgentOutput struct {
	ChatGroupWithAgent *entity.ChatGroupWithAgent
}

func (i *ChatGroupWithAgentInteractorImpl) CreateChatGroupWithAgent(input CreateChatGroupWithAgentInput) (CreateChatGroupWithAgentOutput, error) {
	var (
		output    CreateChatGroupWithAgentOutput
		err       error
		chatGroup *entity.ChatGroupWithAgent
		agent1ID  uint // 数字が小さい方
		agent2ID  uint // 数字が大きい方
	)

	if input.CreateParam.Agent1ID < input.CreateParam.Agent2ID {
		agent1ID = input.CreateParam.Agent1ID
		agent2ID = input.CreateParam.Agent2ID
	} else {
		agent1ID = input.CreateParam.Agent2ID
		agent2ID = input.CreateParam.Agent1ID
	}

	// グループが作成済みかを確認
	chatGroup, err = i.chatGroupWithAgentRepository.FindByAgentID(agent1ID, agent2ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// エージェントと求職者のチャットグループを作成
			chatGroup = entity.NewChatGroupWithAgent(
				agent1ID,
				agent2ID,
			)
			err = i.chatGroupWithAgentRepository.Create(chatGroup)
			if err != nil {
				if !errors.Is(err, entity.ErrDuplicateEntry) {
					// 「Duplicate」以外のエラーの場合は終了
					return output, err
				}
			}
		}
	}

	_, err = i.agentAllianceRepository.FindByAgentID(agent1ID, agent2ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// アライアンス申請のレコードがない場合は作成
			agentAlliance := entity.NewAgentAlliance(
				agent1ID,
				agent2ID,
				false,
				false,
				false, // agent1CancelRequest,
				false, // agent2CancelRequest,
			)
			err = i.agentAllianceRepository.Create(agentAlliance)
			if err != nil {
				if !errors.Is(err, entity.ErrDuplicateEntry) {
					// 「Duplicate」以外のエラーの場合は終了
					return output, err
				}
			}
		} else {
			fmt.Println(err)
			return output, err
		}
	}

	output.ChatGroupWithAgent = chatGroup

	return output, nil
}

// エージェントと求職者のチャットグループの作成
type CreateChatGroupWithAgentListInput struct {
	CreateParam entity.MyAgentIDAndOtherAgentIDListParam
}

type CreateChatGroupWithAgentListOutput struct {
	ChatGroupWithAgentList []*entity.ChatGroupWithAgent
}

func (i *ChatGroupWithAgentInteractorImpl) CreateChatGroupWithAgentList(input CreateChatGroupWithAgentListInput) (CreateChatGroupWithAgentListOutput, error) {
	var (
		output        CreateChatGroupWithAgentListOutput
		chatGroup     *entity.ChatGroupWithAgent
		chatGroupList []*entity.ChatGroupWithAgent
		err           error
	)

	// エージェントと求職者のチャットグループを作成
	for _, unCreatedID := range input.CreateParam.OtherAgentIDList {
		var (
			agent1ID uint // 数字が小さい方
			agent2ID uint // 数字が大きい方
		)

		if input.CreateParam.MyAgentID < unCreatedID {
			agent1ID = input.CreateParam.MyAgentID
			agent2ID = unCreatedID
		} else {
			agent1ID = unCreatedID
			agent2ID = input.CreateParam.MyAgentID
		}

		chatGroup, err = i.chatGroupWithAgentRepository.FindByAgentID(agent1ID, agent2ID) // 先に申請されているかを確認
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				// エージェントと求職者のチャットグループを作成
				chatGroup = entity.NewChatGroupWithAgent(
					agent1ID,
					agent2ID,
				)

				err := i.chatGroupWithAgentRepository.Create(chatGroup)
				if err != nil {
					if !errors.Is(err, entity.ErrDuplicateEntry) {
						// 「Duplicate」以外の場合はエラー
						fmt.Println(err)
						return output, err
					}
				}
			}
		}
		// 32, 33, 34
		chatGroupList = append(chatGroupList, chatGroup)
	}

	output.ChatGroupWithAgentList = chatGroupList

	return output, nil
}

// エージェントの最終閲覧時間を更新
type UpdateAgentLastWatchedForAgentChatInput struct {
	GroupID uint
	AgentID uint
}

type UpdateAgentLastWatchedForAgentChatOutput struct {
	OK bool
}

func (i *ChatGroupWithAgentInteractorImpl) UpdateAgentLastWatched(input UpdateAgentLastWatchedForAgentChatInput) (UpdateAgentLastWatchedForAgentChatOutput, error) {
	var (
		output UpdateAgentLastWatchedForAgentChatOutput
	)

	err := i.chatGroupWithAgentRepository.UpdateAgentLastWatchedAt(input.GroupID, input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// IDを使ってエージェントとエージェントのチャットグループ情報を取得する
type GetChatGroupWithAgentByIDInput struct {
	AgentID uint
	ID      uint
}

type GetChatGroupWithAgentByIDOutput struct {
	ChatGroupWithAgent *entity.ChatGroupWithAgent
}

func (i *ChatGroupWithAgentInteractorImpl) GetChatGroupWithAgentByID(input GetChatGroupWithAgentByIDInput) (GetChatGroupWithAgentByIDOutput, error) {
	var (
		output GetChatGroupWithAgentByIDOutput
		err    error
	)

	chatGroup, err := i.chatGroupWithAgentRepository.FindByIDAndAgentID(input.AgentID, input.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatGroupWithAgent = chatGroup

	return output, nil
}

// エージェントIDからチャットグループ情報を取得する
type GetChatGroupWithAgentListByAgentIDInput struct {
	AgentID uint
}

type GetChatGroupWithAgentListByAgentIDOutput struct {
	ChatGroupWithAgentList []*entity.ChatGroupWithAgent
}

func (i *ChatGroupWithAgentInteractorImpl) GetChatGroupWithAgentListByAgentID(input GetChatGroupWithAgentListByAgentIDInput) (GetChatGroupWithAgentListByAgentIDOutput, error) {
	var (
		output GetChatGroupWithAgentListByAgentIDOutput
		err    error
	)

	chatGroupList, err := i.chatGroupWithAgentRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatGroupWithAgentList = chatGroupList

	return output, nil
}

// IDを使ってエージェントとエージェントのチャットグループ情報を取得する
type GetChatGroupAndThreadWithAgentByChatGroupIDInput struct {
	GroupID      uint
	AgentStaffID uint
}

type GetChatGroupAndThreadWithAgentByChatGroupIDOutput struct {
	ChatGroupWithAgent *entity.ChatGroupWithAgent
}

func (i *ChatGroupWithAgentInteractorImpl) GetChatGroupAndThreadWithAgentByChatGroupID(input GetChatGroupAndThreadWithAgentByChatGroupIDInput) (GetChatGroupAndThreadWithAgentByChatGroupIDOutput, error) {
	var (
		output GetChatGroupAndThreadWithAgentByChatGroupIDOutput
		err    error
	)

	chatGroup, err := i.chatGroupWithAgentRepository.FindByID(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	chatThread, err := i.chatThreadWithAgentRepository.GetByGroupID(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	chatMessageToUser, err := i.chatMessageToUserWithAgentRepository.GetByGroupIDAndUnwatched(input.GroupID, input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// スレッドごとに未読件数をカウントする
	for _, thread := range chatThread {
		if chatGroup.ID == thread.GroupID {

			for _, messageToUser := range chatMessageToUser {
				if messageToUser.ThreadID == thread.ID {
					thread.UnwatchedCount++
				}
			}

			chatGroup.Threads = append(chatGroup.Threads, *thread)
		}
	}

	output.ChatGroupWithAgent = chatGroup

	return output, nil
}

// 自社のエージェントIDと他社のエージェントIDを使ってグループ情報を取得
type GetChatGroupWithAgentByAgentIDAndOtherAgentIDInput struct {
	AgentID      uint
	OtherAgentID uint
}

type GetChatGroupWithAgentByAgentIDAndOtherAgentIDOutput struct {
	ChatGroupWithAgent *entity.ChatGroupWithAgent
}

func (i *ChatGroupWithAgentInteractorImpl) GetChatGroupWithAgentByAgentIDAndOtherAgentID(input GetChatGroupWithAgentByAgentIDAndOtherAgentIDInput) (GetChatGroupWithAgentByAgentIDAndOtherAgentIDOutput, error) {
	var (
		output GetChatGroupWithAgentByAgentIDAndOtherAgentIDOutput
		err    error
	)

	chatGroup, err := i.chatGroupWithAgentRepository.FindByAgentID(input.AgentID, input.OtherAgentID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			fmt.Println(err)
			return output, nil
		} else {
			fmt.Println(err)
			return output, err
		}
	}

	output.ChatGroupWithAgent = chatGroup

	return output, nil
}

// エージェントIDと担当者からチャットグループ情報を取得する 未読件数も取得する
type GetChatGroupWithAgentListByAgentIDAndAgentStaffIDInput struct {
	AgentID      uint
	AgentStaffID uint
}

type GetChatGroupWithAgentListByAgentIDAndAgentStaffIDOutput struct {
	ChatGroupWithAgentList []*entity.ChatGroupWithAgent
}

func (i *ChatGroupWithAgentInteractorImpl) GetChatGroupWithAgentListByAgentIDAndAgentStaffID(input GetChatGroupWithAgentListByAgentIDAndAgentStaffIDInput) (GetChatGroupWithAgentListByAgentIDAndAgentStaffIDOutput, error) {
	var (
		output GetChatGroupWithAgentListByAgentIDAndAgentStaffIDOutput
		err    error
	)

	chatGroupList, err := i.chatGroupWithAgentRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	messageToUsers, err := i.chatMessageToUserWithAgentRepository.GetByAgentStaffIDAndUnwatched(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 未読件数をグループごとにカウントする
	for _, chatGroup := range chatGroupList {
		for _, message := range messageToUsers {
			if chatGroup.ID == message.GroupID {
				chatGroup.UnwatchedCount++
			}
		}
	}
	output.ChatGroupWithAgentList = chatGroupList

	return output, nil
}

// チャットグループ未作成が一つでも存在することを確認する。
type CheckAnyChatGroupWithoutGroupInput struct {
	MyAgentID        uint
	OtherAgentIDList []uint
}

type CheckAnyChatGroupWithoutGroupOutput struct {
	IDList []uint
}

func (i *ChatGroupWithAgentInteractorImpl) CheckAnyChatGroupWithoutGroup(input CheckAnyChatGroupWithoutGroupInput) (CheckAnyChatGroupWithoutGroupOutput, error) {
	var (
		output CheckAnyChatGroupWithoutGroupOutput
		err    error
	)

	chatGroupList, err := i.chatGroupWithAgentRepository.GetByMyAgentIDAndOtherIDList(input.MyAgentID, input.OtherAgentIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// appliedAllianceListの値をマップに格納する
	existsOtherAgentID := make(map[uint]bool)
	for _, chatGroup := range chatGroupList {
		if input.MyAgentID == chatGroup.Agent1ID {
			existsOtherAgentID[chatGroup.Agent2ID] = true
		} else {
			existsOtherAgentID[chatGroup.Agent1ID] = true
		}
	}

	// input.OtherAgentIDListの各値がマップに存在しない場合、結果に追加する
	var unAppliedAgentIDList []uint
	for _, value := range input.OtherAgentIDList {
		if !existsOtherAgentID[value] {
			unAppliedAgentIDList = append(unAppliedAgentIDList, value)
		}
	}

	output.IDList = unAppliedAgentIDList

	return output, nil
}
