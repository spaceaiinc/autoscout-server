package interactor

import (
	"errors"
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type NotificationForUserInteractor interface {
	// 汎用系
	// Systemからのお知らせを作成
	CreateNotificationForUser(input CreateNotificationForUserInput) (CreateNotificationForUserOutput, error)
	// 全てのお知らせを取得
	GetNotificationForUserListByPageAndTargetList(input GetNotificationForUserListByPageAndTargetListInput) (GetNotificationForUserListByPageAndTargetListOutput, error)

	// 既読判定系
	// Systemからのお知らせを確認したユーザーを作成
	CreateUserNotificationView(input CreateUserNotificationViewInput) (CreateUserNotificationViewOutput, error)
	// 担当者IDからお知らせの未読件数を取得
	GetUnwatchedNotificationForUserCountByAgentStaffID(input GetUnwatchedNotificationForUserCountByAgentStaffIDInput) (GetUnwatchedNotificationForUserCountByAgentStaffIDOutput, error)
}

type NotificationForUserInteractorImpl struct {
	firebase                       usecase.Firebase
	sendgrid                       config.Sendgrid
	oneSignal                      config.OneSignal
	notificationForUserRepository  usecase.NotificationForUserRepository
	userNotificationViewRepository usecase.UserNotificationViewRepository
	agentRepository                usecase.AgentRepository
}

// NotificationForUserInteractorImpl is an implementation of NotificationForUserInteractor
func NewNotificationForUserInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	nfuR usecase.NotificationForUserRepository,
	unvR usecase.UserNotificationViewRepository,
	aR usecase.AgentRepository,
) NotificationForUserInteractor {
	return &NotificationForUserInteractorImpl{
		firebase:                       fb,
		sendgrid:                       sg,
		oneSignal:                      os,
		notificationForUserRepository:  nfuR,
		userNotificationViewRepository: unvR,
		agentRepository:                aR,
	}
}

/****************************************************************************************/
/// 汎用系
//
// Systemからのお知らせを作成
type CreateNotificationForUserInput struct {
	CreateParam entity.NotificationForUser
}

type CreateNotificationForUserOutput struct {
	NotificationForUser *entity.NotificationForUser
}

func (i *NotificationForUserInteractorImpl) CreateNotificationForUser(input CreateNotificationForUserInput) (CreateNotificationForUserOutput, error) {
	var (
		output CreateNotificationForUserOutput
		err    error
	)

	// お知らせの作成
	notificationForUser := entity.NewNotificationForUser(
		input.CreateParam.Title,
		input.CreateParam.Body,
		input.CreateParam.Target,
	)

	err = i.notificationForUserRepository.Create(notificationForUser)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.NotificationForUser = notificationForUser

	return output, nil
}

// 全てのお知らせを取得
type GetNotificationForUserListByPageAndTargetListInput struct {
	PageNumber uint
	TargetList []entity.NotificationForUserTarget
}

type GetNotificationForUserListByPageAndTargetListOutput struct {
	NotificationForUserList []*entity.NotificationForUser
	MaxPageNumber           uint
}

func (i *NotificationForUserInteractorImpl) GetNotificationForUserListByPageAndTargetList(input GetNotificationForUserListByPageAndTargetListInput) (GetNotificationForUserListByPageAndTargetListOutput, error) {
	var (
		output                  GetNotificationForUserListByPageAndTargetListOutput
		notificationForUserList []*entity.NotificationForUser
		notificationIDList      []uint
	)

	notificationForUserList, err := i.notificationForUserRepository.GetByTargetList(input.TargetList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ページの最大数を取得 10件ずつ
	const per = 10
	output.MaxPageNumber = uint(len(notificationForUserList) / per)
	if 0 < (len(notificationForUserList) % per) {
		output.MaxPageNumber = uint(output.MaxPageNumber + 1)
	}

	// ページに応じてお知らせを取得
	start := (input.PageNumber - 1) * uint(per)
	end := input.PageNumber * uint(per)

	if len(notificationForUserList) < int(end) {
		end = uint(len(notificationForUserList))
	}

	notificationForUserList = notificationForUserList[start:end]

	for _, notificationForUser := range notificationForUserList {
		notificationIDList = append(notificationIDList, notificationForUser.ID)
	}

	userNotificationViewList, err := i.userNotificationViewRepository.GetByNotificationIDList(notificationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, notificationForUser := range notificationForUserList {
		for _, userNotificationView := range userNotificationViewList {
			if notificationForUser.ID == userNotificationView.NotificationID {
				notificationForUser.UserNotificationViews = append(notificationForUser.UserNotificationViews, *userNotificationView)
			}
		}
	}

	output.NotificationForUserList = notificationForUserList

	return output, nil
}

/****************************************************************************************/
/// 既読判定系
//
// Systemからのお知らせを確認したユーザーを作成
type CreateUserNotificationViewInput struct {
	AgentStaffID uint
}

type CreateUserNotificationViewOutput struct {
	OK bool
}

func (i *NotificationForUserInteractorImpl) CreateUserNotificationView(input CreateUserNotificationViewInput) (CreateUserNotificationViewOutput, error) {
	var (
		output CreateUserNotificationViewOutput
		err    error
	)

	// 最新の既読を取得
	latestView, err := i.userNotificationViewRepository.FindLatestByAgentStaffID(input.AgentStaffID)
	if errors.Is(err, entity.ErrNotFound) {
		// お知らせがない場合はそのまま処理を続ける
		err = nil
	} else if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 最新のお知らせを取得
	latestNotification, err := i.notificationForUserRepository.FindLatest()
	if errors.Is(err, entity.ErrNotFound) {
		// お知らせがない場合は何もしない
		return output, nil
	} else if err != nil {
		fmt.Println(err)
		return output, err
	} else {

		// 既読がない場合 or 最新の既読と最新のお知らせが異なる場合 は最新のお知らせを既読にする
		if latestView == nil || latestView.NotificationID != latestNotification.ID {
			userNotificationView := entity.NewUserNotificationView(
				latestNotification.ID,
				input.AgentStaffID,
			)

			err = i.userNotificationViewRepository.Create(userNotificationView)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true

	return output, nil
}

// 担当者IDからお知らせの未読件数を取得
type GetUnwatchedNotificationForUserCountByAgentStaffIDInput struct {
	AgentStaffID uint
}

type GetUnwatchedNotificationForUserCountByAgentStaffIDOutput struct {
	UnwatchedNotificationForUserCount uint
}

func (i *NotificationForUserInteractorImpl) GetUnwatchedNotificationForUserCountByAgentStaffID(input GetUnwatchedNotificationForUserCountByAgentStaffIDInput) (GetUnwatchedNotificationForUserCountByAgentStaffIDOutput, error) {
	var (
		output                           GetUnwatchedNotificationForUserCountByAgentStaffIDOutput
		unwatchedNotificationForUserList []*entity.NotificationForUser
		targetList                       []entity.NotificationForUserTarget = []entity.NotificationForUserTarget{entity.NotificationForUserTargetAll}
	)

	agent, err := i.agentRepository.FindByAgentStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントの利用タイプから該当するお知らせを取得
	if agent.IsCRMActive {
		targetList = append(targetList, entity.NotificationForUserTargetCRM)
	}
	if agent.IsSendingActive {
		targetList = append(targetList, entity.NotificationForUserTargetSending)
	}

	// 既読の中で最新のお知らせを取得
	latestView, err := i.userNotificationViewRepository.FindLatestByAgentStaffID(input.AgentStaffID)
	// 既読がない場合は最新のお知らせを取得
	if errors.Is(err, entity.ErrNotFound) {
		unwatchedNotificationForUserList, err = i.notificationForUserRepository.GetByTargetList(targetList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// NotFound以外のエラーの場合
	} else if err != nil {
		fmt.Println(err)
		return output, err

		// レコードがある場合は既読以降のお知らせを取得
	} else {
		unwatchedNotificationForUserList, err = i.notificationForUserRepository.GetAfterIDAndTargetList(latestView.NotificationID, targetList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.UnwatchedNotificationForUserCount = uint(len(unwatchedNotificationForUserList))

	return output, nil
}
