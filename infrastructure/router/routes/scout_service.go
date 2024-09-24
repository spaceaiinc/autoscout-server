package routes

import (
	"fmt"
	"log"
	"net/http"
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
// 汎用系 API
// スカウトサービスの登録
func CreateScoutService(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateScoutServiceParam
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

		h := di.InitializeScoutServiceHandler(firebase, tx, sendgrid, oneSignal, appVar, googleAPI, slack)
		p, err := h.CreateScoutService(param)
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

// スカウトサービスの更新
func UpdateScoutService(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateScoutServiceParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		scoutServiceIDStr := c.Param("scout_service_id")
		log.Println("scoutServiceIDStr", scoutServiceIDStr)
		scoutServiceIDInt, err := strconv.Atoi(scoutServiceIDStr)
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

		h := di.InitializeScoutServiceHandler(firebase, tx, sendgrid, oneSignal, appVar, googleAPI, slack)
		p, err := h.UpdateScoutService(param, uint(scoutServiceIDInt))
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

// スカウトサービスのパスワード更新
func UpdateScoutServicePassword(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateScoutServicePasswordParam
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

		h := di.InitializeScoutServiceHandler(firebase, tx, sendgrid, oneSignal, appVar, googleAPI, slack)
		p, err := h.UpdateScoutServicePassword(param)
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

// スカウトサービスの削除
func DeleteScoutService(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		scoutServiceIDStr := c.Param("scout_service_id")

		scoutServiceIDInt, err := strconv.Atoi(scoutServiceIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeScoutServiceHandler(firebase, db, sendgrid, oneSignal, appVar, googleAPI, slack)
		p, err := h.DeleteScoutService(uint(scoutServiceIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// IDからスカウトサービスを取得
func GetScoutServiceByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		scoutServiceIDStr := c.Param("scout_service_id")

		scoutServiceIDInt, err := strconv.Atoi(scoutServiceIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeScoutServiceHandler(firebase, db, sendgrid, oneSignal, appVar, googleAPI, slack)
		p, err := h.GetScoutServiceByID(uint(scoutServiceIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDからスカウトサービスを取得
func GetScoutServiceListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		agentIDStr := c.Param("agent_id")

		agentIDInt, err := strconv.Atoi(agentIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeScoutServiceHandler(firebase, db, sendgrid, oneSignal, appVar, googleAPI, slack)
		p, err := h.GetScoutServiceListByAgentID(uint(agentIDInt))
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
// Gmail API
func GmailWebHook(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, appVar config.App, googleAPI config.GoogleAPI, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			pubSubStruct = new(entity.PubsubStruct)
		)

		if err := c.Bind(pubSubStruct); err != nil {
			fmt.Println("Bindエラー: ", err)
			return c.JSON(http.StatusOK, map[string]bool{
				"success": false,
			})
		}

		h := di.InitializeScoutServiceHandler(firebase, db, sendgrid, oneSignal, appVar, googleAPI, slack)
		p, err := h.GmailWebHook(pubSubStruct)
		if err != nil {
			fmt.Println(err)

			// エラーでもstatus 200で返す
			hookError := c.JSON(http.StatusOK, map[string]bool{
				"success": false,
			})
			return hookError
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
