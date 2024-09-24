package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type NotificationForUserHandler interface {
	// 汎用系 API
	// Motoyuiからのお知らせを作成
	CreateNotificationForUser(param entity.NotificationForUser) (presenter.Presenter, error)
	// 全てのMotoyuiからのお知らせを取得
	GetNotificationForUserListByPageAndTargetList(pageNumber uint, targetList []entity.NotificationForUserTarget) (presenter.Presenter, error)

	// 既読判定系
	// Motoyuiからのお知らせを確認したユーザーを作成
	CreateUserNotificationView(agentStaffID uint) (presenter.Presenter, error)
	// 担当者IDからお知らせの未読件数を取得
	GetUnwatchedNotificationForUserCountByAgentStaffID(agentStaffID uint) (presenter.Presenter, error)
}

type NotificationForUserHandlerImpl struct {
	notificationForUserInteractor interactor.NotificationForUserInteractor
}

func NewNotificationForUserHandlerImpl(nfuI interactor.NotificationForUserInteractor) NotificationForUserHandler {
	return &NotificationForUserHandlerImpl{
		notificationForUserInteractor: nfuI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
// Motoyuiからのお知らせを作成
func (h *NotificationForUserHandlerImpl) CreateNotificationForUser(param entity.NotificationForUser) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.CreateNotificationForUser(interactor.CreateNotificationForUserInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewNotificationForUserJSONPresenter(responses.NewNotificationForUser(output.NotificationForUser)), nil
}

// 送信対象からMotoyuiからのお知らせを取得
func (h *NotificationForUserHandlerImpl) GetNotificationForUserListByPageAndTargetList(pageNumber uint, targetList []entity.NotificationForUserTarget) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.GetNotificationForUserListByPageAndTargetList(interactor.GetNotificationForUserListByPageAndTargetListInput{
		PageNumber: pageNumber,
		TargetList: targetList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewNotificationForUserListAndMaxPageNumberJSONPresenter(responses.NewNotificationForUserListAndMaxPageNumber(output.NotificationForUserList, output.MaxPageNumber)), nil
}

/****************************************************************************************/
/// 既読判定系 API
//
// Motoyuiからのお知らせを確認したユーザーを作成
func (h *NotificationForUserHandlerImpl) CreateUserNotificationView(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.CreateUserNotificationView(interactor.CreateUserNotificationViewInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 担当者IDからお知らせの未読件数を取得
func (h *NotificationForUserHandlerImpl) GetUnwatchedNotificationForUserCountByAgentStaffID(agentStaffID uint) (presenter.Presenter, error) {
	output, err := h.notificationForUserInteractor.GetUnwatchedNotificationForUserCountByAgentStaffID(interactor.GetUnwatchedNotificationForUserCountByAgentStaffIDInput{
		AgentStaffID: agentStaffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewUnwatchedNotificationForUserCountJSONPresenter(responses.NewUnwatchedNotificationForUserCount(output.UnwatchedNotificationForUserCount)), nil
}
