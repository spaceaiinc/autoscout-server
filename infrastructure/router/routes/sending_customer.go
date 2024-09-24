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
// 汎用系 API
//
// 指定IDの送客求職者を取得
func GetSendingCustomerByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingCustomerIDStr = c.Param("sending_customer_id")
		)

		sendingCustomerID, err := strconv.Atoi(sendingCustomerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingCustomerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSendingCustomerByID(uint(sendingCustomerID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 指定エージェントIDの送客求職者リストを取得
func GetSendingCustomerListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.QueryParam("agent_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingCustomerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSendingCustomerListByAgentID(uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetSearchSendingCustomerListByPageAndTabAndAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr    = c.QueryParam("agent_id")
			pageNumberStr = c.QueryParam("page_number")
			tabNumberStr  = c.QueryParam("tab_number")
			freeWordStr   = c.QueryParam("free_word")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		pageNumber, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tabNumber, err := strconv.Atoi(tabNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingCustomerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSearchSendingCustomerListByPageAndTabAndAgentID(uint(pageNumber), uint(agentID), uint(tabNumber), freeWordStr)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
