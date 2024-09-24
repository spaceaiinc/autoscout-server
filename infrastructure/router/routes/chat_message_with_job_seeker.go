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
// エージェントと求職者のチャットメッセージの登録
func SendChatMessageWithJobSeekerLineMessage(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendChatMessageWithJobSeekerLineParam
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

		h := di.InitializeChatMessageWithJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.SendChatMessageWithJobSeekerLineMessage(param)
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

// エージェントと求職者のチャットメッセージの登録
func SendChatMessageWithJobSeekerLineImage(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		// var (
		// 	param entity.SendChatMessageWithJobSeekerLineParam
		// )

		// フォームパラムの取得
		formParams, err := c.FormParams()
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		groupIDStr := formParams.Get("group_id")

		// groupID 型変換
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// ファイル
		fileParam, err := c.FormFile("file")
		if err != nil {
			wrapped := fmt.Errorf("ファイルの受け取りエラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		param := entity.SendChatMessageWithJobSeekerLineImageParam{
			GroupID: uint(groupID),
			File:    fileParam,
		}

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

		h := di.InitializeChatMessageWithJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.SendChatMessageWithJobSeekerLineImage(param)
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

// エージェントIDからエージェントと求職者のチャットメッセージメッセージの一覧取得
func GetChatMessageWithJobSeekerListByGroupID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			groupIDStr = c.Param("group_id")
		)

		groupIDInt, err := strconv.Atoi(groupIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeChatMessageWithJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetChatMessageWithJobSeekerListByGroupID(uint(groupIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
/// LINE WebHook
func LineWebHook(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := c.Request()

		h := di.InitializeChatMessageWithJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.LineWebHook(req)
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
