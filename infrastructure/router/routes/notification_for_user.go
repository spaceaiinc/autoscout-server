package routes

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/di"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

/****************************************************************************************/
/// 汎用系 API
//
// Systemからのお知らせを作成
func CreateNotificationForUser(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.NotificationForUser
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeNotificationForUserHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.CreateNotificationForUser(param)
		if err != nil {
			tx.Rollback()
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		tx.Commit()
		renderJSON(c, p)
		return nil
	}
}

// ページと送信対象からお知らせを取得 @param: page_number, target_list
func GetNotificationForUserListByPageAndTargetList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			pageNumberStr = c.QueryParam("page_number")

			targetListStr = c.QueryParams()["target_list[]"]
			targetList    []entity.NotificationForUserTarget
		)

		pageNumber, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		if len(targetListStr) == 0 {
			wrapped := fmt.Errorf("%s:%w", "target_list is required", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// target_listを配列に変換
		for _, target := range targetListStr {
			targetInt, err := strconv.Atoi(target)
			if err != nil {
				wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
				renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
				return wrapped
			}
			targetList = append(targetList, entity.NotificationForUserTarget(targetInt))
		}

		h := di.InitializeNotificationForUserHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetNotificationForUserListByPageAndTargetList(uint(pageNumber), targetList)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
/// 既読判定系 API
//
// Systemからのお知らせを確認したユーザーを作成
func CreateUserNotificationView(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeNotificationForUserHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.CreateUserNotificationView(uint(agentStaffID))
		if err != nil {
			tx.Rollback()
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		tx.Commit()
		renderJSON(c, p)
		return nil
	}
}

// 担当者IDからお知らせの未読件数を取得
func GetUnwatchedNotificationForUserCountByAgentStaffID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeNotificationForUserHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetUnwatchedNotificationForUserCountByAgentStaffID(uint(agentStaffID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
