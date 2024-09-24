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
// カレンダに表示するスケジュール一覧を取得
//
func GetScheduleListWithInPeriod(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
			startDateStr    = c.QueryParam("start_date")
			endDateStr      = c.QueryParam("end_date")
		)

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeScheduleHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetScheduleListWithInPeriod(uint(agentStaffIDInt), startDateStr, endDateStr)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetScheduleListWithInPeriodByStaffIDList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			staffIDStrList = c.QueryParams()["staff_id_list[]"]
			startDateStr   = c.QueryParam("start_date")
			endDateStr     = c.QueryParam("end_date")
		)

		staffIDList, err := parseQueryParamUINT(staffIDStrList)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeScheduleHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetScheduleListWithInPeriodByStaffIDList(staffIDList, startDateStr, endDateStr)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
