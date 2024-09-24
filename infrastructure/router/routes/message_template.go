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
// メッセージテンプレートの登録
func CreateMessageTemplate(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateMessageTemplateParam
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

		h := di.InitializeMessageTemplateHandler(firebase, tx, sendgrid)
		p, err := h.CreateMessageTemplate(param)
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

// メッセージテンプレートの更新
func UpdateMessageTemplate(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateMessageTemplateParam
		)

		templateIDStr := c.Param("template_id")

		templateIDInt, _ := strconv.Atoi(templateIDStr)

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

		h := di.InitializeMessageTemplateHandler(firebase, tx, sendgrid)
		p, err := h.UpdateMessageTemplate(param, uint(templateIDInt))
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

// メッセージテンプレートの削除
func DeleteMessageTemplate(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		templateIDStr := c.Param("template_id")

		templateIDInt, err := strconv.Atoi(templateIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeMessageTemplateHandler(firebase, db, sendgrid)
		p, err := h.DeleteMessageTemplate(uint(templateIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// IDからメッセージテンプレートを取得
func GetMessageTemplateByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		templateIDStr := c.Param("template_id")

		templateIDInt, err := strconv.Atoi(templateIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeMessageTemplateHandler(firebase, db, sendgrid)
		p, err := h.GetMessageTemplateByID(uint(templateIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 担当者IDからメッセージテンプレートを取得
func GetMessageTemplateListByAgentStaffID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeMessageTemplateHandler(firebase, db, sendgrid)
		p, err := h.GetMessageTemplateListByAgentStaffID(uint(agentStaffIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 担当者IDと送信シーンからメッセージテンプレートを取得
func GetMessageTemplateListByAgentStaffIDAndSendScene(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
			sendSceneStr    = c.Param("send_scene")
		)

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		sendSceneInt, err := strconv.Atoi(sendSceneStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeMessageTemplateHandler(firebase, db, sendgrid)
		p, err := h.GetMessageTemplateListByAgentStaffIDAndSendScene(uint(agentStaffIDInt), uint(sendSceneInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDからメッセージテンプレートを取得
func GetMessageTemplateListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		agentIDStr := c.Param("agent_id")

		agentIDInt, err := strconv.Atoi(agentIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeMessageTemplateHandler(firebase, db, sendgrid)
		p, err := h.GetMessageTemplateListByAgentID(uint(agentIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API

/****************************************************************************************/
