package interactor

import (
	"fmt"
	"strconv"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type ChatGroupWithJobSeekerInteractor interface {
	// 汎用系 API
	CreateChatGroupWithJobSeeker(input CreateChatGroupWithJobSeekerInput) (CreateChatGroupWithJobSeekerOutput, error)
	UpdateAgentLastWatched(input UpdateAgentLastWatchedInput) (UpdateAgentLastWatchedOutput, error)
	GetChatGroupWithJobSeekerByID(input GetChatGroupWithJobSeekerByIDInput) (GetChatGroupWithJobSeekerByIDOutput, error)
	GetChatGroupWithJobSeekerListByAgentID(input GetChatGroupWithJobSeekerListByAgentIDInput) (GetChatGroupWithJobSeekerListByAgentIDOutput, error)
	GetChatGroupWithJobSeekerSearchListAndPageByAgentID(input GetChatGroupWithJobSeekerSearchListAndPageByAgentIDInput) (GetChatGroupWithJobSeekerSearchListAndPageByAgentIDOutput, error)
	GetJobSeekerChatNotificationByAgentID(input GetJobSeekerChatNotificationByAgentIDInput) (GetJobSeekerChatNotificationByAgentIDOutput, error)
}

type ChatGroupWithJobSeekerInteractorImpl struct {
	firebase                         usecase.Firebase
	sendgrid                         config.Sendgrid
	chatGroupWithJobSeekerRepository usecase.ChatGroupWithJobSeekerRepository
}

// ChatGroupWithJobSeekerInteractorImpl is an implementation of ChatGroupWithJobSeekerInteractor
func NewChatGroupWithJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	cgR usecase.ChatGroupWithJobSeekerRepository,
) ChatGroupWithJobSeekerInteractor {
	return &ChatGroupWithJobSeekerInteractorImpl{
		firebase:                         fb,
		sendgrid:                         sg,
		chatGroupWithJobSeekerRepository: cgR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//エージェントと求職者のチャットグループの作成
type CreateChatGroupWithJobSeekerInput struct {
	CreateParam entity.CreateChatGroupWithJobSeekerParam
}

type CreateChatGroupWithJobSeekerOutput struct {
	ChatGroupWithJobSeeker *entity.ChatGroupWithJobSeeker
}

func (i *ChatGroupWithJobSeekerInteractorImpl) CreateChatGroupWithJobSeeker(input CreateChatGroupWithJobSeekerInput) (CreateChatGroupWithJobSeekerOutput, error) {
	var (
		output CreateChatGroupWithJobSeekerOutput
		err    error
	)

	// エージェントと求職者のチャットグループを作成
	chatGroup := entity.NewChatGroupWithJobSeeker(
		input.CreateParam.AgentID,
		input.CreateParam.JobSeekerID,
		false, // 初めはLINE
	)

	err = i.chatGroupWithJobSeekerRepository.Create(chatGroup)
	if err != nil {
		return output, err
	}

	output.ChatGroupWithJobSeeker = chatGroup

	return output, nil
}

// エージェントの最終閲覧時間を更新
type UpdateAgentLastWatchedInput struct {
	GroupID uint
}

type UpdateAgentLastWatchedOutput struct {
	OK bool
}

func (i *ChatGroupWithJobSeekerInteractorImpl) UpdateAgentLastWatched(input UpdateAgentLastWatchedInput) (UpdateAgentLastWatchedOutput, error) {
	var (
		output UpdateAgentLastWatchedOutput
	)

	err := i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true
	return output, nil
}

// IDを使ってエージェントと求職者のチャットグループ情報を取得する
type GetChatGroupWithJobSeekerByIDInput struct {
	ID uint
}

type GetChatGroupWithJobSeekerByIDOutput struct {
	ChatGroupWithJobSeeker *entity.ChatGroupWithJobSeeker
}

func (i *ChatGroupWithJobSeekerInteractorImpl) GetChatGroupWithJobSeekerByID(input GetChatGroupWithJobSeekerByIDInput) (GetChatGroupWithJobSeekerByIDOutput, error) {
	var (
		output GetChatGroupWithJobSeekerByIDOutput
		err    error
	)

	chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByID(input.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatGroupWithJobSeeker = chatGroup

	return output, nil
}

// エージェントIDからチャットグループ情報を取得する
type GetChatGroupWithJobSeekerListByAgentIDInput struct {
	AgentID uint
}

type GetChatGroupWithJobSeekerListByAgentIDOutput struct {
	ChatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker
}

func (i *ChatGroupWithJobSeekerInteractorImpl) GetChatGroupWithJobSeekerListByAgentID(input GetChatGroupWithJobSeekerListByAgentIDInput) (GetChatGroupWithJobSeekerListByAgentIDOutput, error) {
	var (
		output GetChatGroupWithJobSeekerListByAgentIDOutput
		err    error
	)

	chatGroupList, err := i.chatGroupWithJobSeekerRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.ChatGroupWithJobSeekerList = chatGroupList

	return output, nil
}

// エージェントIDからチャットグループ情報を取得する
type GetChatGroupWithJobSeekerSearchListAndPageByAgentIDInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchChatJobSeeker
}

type GetChatGroupWithJobSeekerSearchListAndPageByAgentIDOutput struct {
	ChatGroupWithJobSeekerList []*entity.ChatGroupWithJobSeeker
	MaxPageNumber              uint
}

func (i *ChatGroupWithJobSeekerInteractorImpl) GetChatGroupWithJobSeekerSearchListAndPageByAgentID(input GetChatGroupWithJobSeekerSearchListAndPageByAgentIDInput) (GetChatGroupWithJobSeekerSearchListAndPageByAgentIDOutput, error) {
	var (
		output GetChatGroupWithJobSeekerSearchListAndPageByAgentIDOutput
		err    error
		param  = input.SearchParam
	)

	/************ 1. チャットグループ情報を取得 **************/

	chatGroupList, err := i.chatGroupWithJobSeekerRepository.GetByAgentIDAndFreeWord(input.AgentID, param.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. 絞り込み処理 **************/

	var (
		chatListWithAgentStaffID []*entity.ChatGroupWithJobSeeker
		chatListWithPhase        []*entity.ChatGroupWithJobSeeker
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
	output.MaxPageNumber = getChatGroupJobSeekerListMaxPage(chatListWithPhase)

	// 指定ページの求職者100件を取得
	output.ChatGroupWithJobSeekerList = getChatGroupJobSeekerListWithPage(chatListWithPhase, input.PageNumber)

	return output, nil
}

// 未読情報を取得する
type GetJobSeekerChatNotificationByAgentIDInput struct {
	AgentID uint
}

type GetJobSeekerChatNotificationByAgentIDOutput struct {
	UnwatchedCount uint
}

func (i *ChatGroupWithJobSeekerInteractorImpl) GetJobSeekerChatNotificationByAgentID(input GetJobSeekerChatNotificationByAgentIDInput) (GetJobSeekerChatNotificationByAgentIDOutput, error) {
	var (
		output GetJobSeekerChatNotificationByAgentIDOutput
		err    error
	)

	count, err := i.chatGroupWithJobSeekerRepository.CountNotificationByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	output.UnwatchedCount = uint(count)

	return output, nil
}

/****************************************************************************************/
/// Admin API
//

/****************************************************************************************/
