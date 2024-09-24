package interactor

import (
	"fmt"
	"strconv"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatGroupWithSendingJobSeekerInteractor interface {
	// 汎用系 API
	CreateChatGroupWithSendingJobSeeker(input CreateChatGroupWithSendingJobSeekerInput) (CreateChatGroupWithSendingJobSeekerOutput, error)
	UpdateSendingJobSeekerAgentLastWatched(input UpdateSendingJobSeekerAgentLastWatchedInput) (UpdateSendingJobSeekerAgentLastWatchedOutput, error)
	GetChatGroupWithSendingJobSeekerByID(input GetChatGroupWithSendingJobSeekerByIDInput) (GetChatGroupWithSendingJobSeekerByIDOutput, error)
	GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentID(input GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDInput) (GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDOutput, error)
	GetSendingJobSeekerChatNotificationByAgentID(input GetSendingJobSeekerChatNotificationByAgentIDInput) (GetSendingJobSeekerChatNotificationByAgentIDOutput, error)
}

type ChatGroupWithSendingJobSeekerInteractorImpl struct {
	firebase                                usecase.Firebase
	sendgrid                                config.Sendgrid
	chatGroupWithSendingJobSeekerRepository usecase.ChatGroupWithSendingJobSeekerRepository
}

// ChatGroupWithSendingJobSeekerInteractorImpl is an implementation of ChatGroupWithSendingJobSeekerInteractor
func NewChatGroupWithSendingJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	cgR usecase.ChatGroupWithSendingJobSeekerRepository,
) ChatGroupWithSendingJobSeekerInteractor {
	return &ChatGroupWithSendingJobSeekerInteractorImpl{
		firebase:                                fb,
		sendgrid:                                sg,
		chatGroupWithSendingJobSeekerRepository: cgR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//エージェントと求職者のチャットグループの作成
type CreateChatGroupWithSendingJobSeekerInput struct {
	CreateParam entity.CreateChatGroupWithSendingJobSeekerParam
}

type CreateChatGroupWithSendingJobSeekerOutput struct {
	ChatGroupWithSendingJobSeeker *entity.ChatGroupWithSendingJobSeeker
}

func (i *ChatGroupWithSendingJobSeekerInteractorImpl) CreateChatGroupWithSendingJobSeeker(input CreateChatGroupWithSendingJobSeekerInput) (CreateChatGroupWithSendingJobSeekerOutput, error) {
	var (
		output CreateChatGroupWithSendingJobSeekerOutput
		err    error
	)

	// エージェントと求職者のチャットグループを作成
	chatGroup := entity.NewChatGroupWithSendingJobSeeker(
		input.CreateParam.SendingJobSeekerID,
		false, // 初めはLINE
	)

	err = i.chatGroupWithSendingJobSeekerRepository.Create(chatGroup)
	if err != nil {
		return output, err
	}

	output.ChatGroupWithSendingJobSeeker = chatGroup

	return output, nil
}

// エージェントの最終閲覧時間を更新
type UpdateSendingJobSeekerAgentLastWatchedInput struct {
	GroupID uint
}

type UpdateSendingJobSeekerAgentLastWatchedOutput struct {
	OK bool
}

func (i *ChatGroupWithSendingJobSeekerInteractorImpl) UpdateSendingJobSeekerAgentLastWatched(input UpdateSendingJobSeekerAgentLastWatchedInput) (UpdateSendingJobSeekerAgentLastWatchedOutput, error) {
	var (
		output UpdateSendingJobSeekerAgentLastWatchedOutput
	)

	err := i.chatGroupWithSendingJobSeekerRepository.UpdateAgentLastWatchedAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true
	return output, nil
}

// IDを使ってエージェントと求職者のチャットグループ情報を取得する
type GetChatGroupWithSendingJobSeekerByIDInput struct {
	ID uint
}

type GetChatGroupWithSendingJobSeekerByIDOutput struct {
	ChatGroupWithSendingJobSeeker *entity.ChatGroupWithSendingJobSeeker
}

func (i *ChatGroupWithSendingJobSeekerInteractorImpl) GetChatGroupWithSendingJobSeekerByID(input GetChatGroupWithSendingJobSeekerByIDInput) (GetChatGroupWithSendingJobSeekerByIDOutput, error) {
	var (
		output GetChatGroupWithSendingJobSeekerByIDOutput
		err    error
	)

	chatGroup, err := i.chatGroupWithSendingJobSeekerRepository.FindByID(input.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatGroupWithSendingJobSeeker = chatGroup

	return output, nil
}

// 絞り込んでチャットグループ情報を取得する
type GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchChatSendingJobSeeker
}

type GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDOutput struct {
	ChatGroupWithSendingJobSeekerList []*entity.ChatGroupWithSendingJobSeeker
	MaxPageNumber                     uint
}

func (i *ChatGroupWithSendingJobSeekerInteractorImpl) GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentID(input GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDInput) (GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDOutput, error) {
	var (
		output GetChatGroupWithSendingJobSeekerSearchListAndPageByAgentIDOutput
		err    error
		param  = input.SearchParam
	)

	chatGroupList, err := i.chatGroupWithSendingJobSeekerRepository.GetListByAgentIDAndFreeWord(input.AgentID, param.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. 絞り込み処理 **************/

	var (
		chatListWithAgentStaffID []*entity.ChatGroupWithSendingJobSeeker
		chatListWithPhase        []*entity.ChatGroupWithSendingJobSeeker
	)

	// 営業担当者IDがある場合
	agentStaffID, err := strconv.Atoi(param.AgentStaffID)
	if !(err != nil || agentStaffID == 0) {
		for _, chatGroup := range chatGroupList {
			if chatGroup.CAStaffID.Int64 == int64(agentStaffID) {
				chatListWithAgentStaffID = append(chatListWithAgentStaffID, chatGroup)
			}
		}
	}

	// 営業担当者IDが無い場合
	if err != nil || agentStaffID == 0 {
		chatListWithAgentStaffID = chatGroupList
	}

	fmt.Println("CA担当者: ", len(chatListWithAgentStaffID))
	// Note: フェーズ
	// フェーズがある場合
	if !(len(param.PhaseTypes) == 0) {
	phaseLoop:
		for _, chatGroup := range chatListWithAgentStaffID {
			for _, phase := range input.SearchParam.PhaseTypes {
				if !phase.Valid {
					continue
				}

				if phase == chatGroup.Phase {
					// 求職者のフェーズが一つでも合致していればcontinueで次の求職者へ

					// lineActiveがfalseかつLineIDが入力済みの場合
					if !chatGroup.LineActive && chatGroup.LineID != "" {
						chatGroup.IsBlocked = true
					}
					chatListWithPhase = append(chatListWithPhase, chatGroup)
					continue phaseLoop
				}
			}
		}
	}

	// フェーズが入っていない場合
	// if len(param.PhaseTypes) == 0 {
	// chatListWithPhase = chatListWithAgentStaffID
	// }

	fmt.Println("フェーズ: ", len(chatListWithPhase))

	/************ 3. ページ関連の処理 **************/
	// ページの最大数を取得
	output.MaxPageNumber = getChatGroupSendingJobSeekerListMaxPage(chatListWithPhase)

	// 指定ページの求職者200件を取得
	output.ChatGroupWithSendingJobSeekerList = getChatGroupSendingJobSeekerListWithPage(chatListWithPhase, input.PageNumber)

	return output, nil
}

// 絞り込んでチャットグループ情報を取得する
type GetSendingJobSeekerChatNotificationByAgentIDInput struct {
	AgentID uint
}

type GetSendingJobSeekerChatNotificationByAgentIDOutput struct {
	UnwatchedCount uint
}

func (i *ChatGroupWithSendingJobSeekerInteractorImpl) GetSendingJobSeekerChatNotificationByAgentID(input GetSendingJobSeekerChatNotificationByAgentIDInput) (GetSendingJobSeekerChatNotificationByAgentIDOutput, error) {
	var (
		output GetSendingJobSeekerChatNotificationByAgentIDOutput
		err    error
	)

	count, err := i.chatGroupWithSendingJobSeekerRepository.GetNotificationCountByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.UnwatchedCount = uint(count)

	return output, nil
}

/****************************************************************************************/
