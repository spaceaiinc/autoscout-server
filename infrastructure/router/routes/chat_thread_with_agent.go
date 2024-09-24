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
func CreateChatThreadWithAgent(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateChatThreadWithAgentParam
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

		h := di.InitializeChatThreadWithAgentHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.CreateChatThreadWithAgent(param)
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

// エージェントIDからエージェントとエージェントのチャットメッセージメッセージの一覧取得

// エージェントIDからエージェントとエージェントのチャットメッセージメッセージの一覧取得
func GetChatThreadAndMessageWithAgentByChatThreadID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			threadIDStr = c.Param("thread_id")
		)

		threadIDInt, err := strconv.Atoi(threadIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeChatThreadWithAgentHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetChatThreadAndMessageWithAgentByChatThreadID(uint(threadIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
