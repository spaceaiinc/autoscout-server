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
// エージェントと送客求職者のチャットメッセージの登録
func SendChatMessageWithSendingJobSeekerLineMessage(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendChatMessageWithSendingJobSeekerLineParam
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

		h := di.InitializeChatMessageWithSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.SendChatMessageWithSendingJobSeekerLineMessage(param)
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

// エージェントと送客求職者のチャットメッセージの登録
func SendChatMessageWithSendingJobSeekerLineImage(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		// var (
		// 	param entity.SendChatMessageWithSendingJobSeekerLineParam
		// )

		// フォームパラむの取得
		formParams, err := c.FormParams()
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		groupIDStr := formParams.Get("group_id")
		// lineIDStr := formParams.Get("line_id")

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

		// file, err := fileParam.Open()
		// if err != nil {
		// 	wrapped := fmt.Errorf("ファイルが開けません: %s:%w", err.Error(), entity.ErrRequestError)
		// 	renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		// 	return wrapped
		// }

		param := entity.SendChatMessageWithSendingJobSeekerLineImageParam{
			GroupID: uint(groupID),
			// LineID:  lineIDStr,
			File: fileParam,
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

		h := di.InitializeChatMessageWithSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.SendChatMessageWithSendingJobSeekerLineImage(param)
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

// エージェントIDからエージェントと送客求職者のチャットメッセージメッセージの一覧取得
func GetChatMessageWithSendingJobSeekerListByGroupID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
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

		h := di.InitializeChatMessageWithSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetChatMessageWithSendingJobSeekerListByGroupID(uint(groupIDInt))
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
