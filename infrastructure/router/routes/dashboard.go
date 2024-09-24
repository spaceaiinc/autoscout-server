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
// デフォルト表示関連
//

// ダッシュボードのデフォルト表示に必要な情報を取得（当月が含まれる売上期間の受注ベースの個人売上情報を取得する）
func GetDashboardForDefaultPreview(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
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

		h := di.InitializeDashboardHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetDashboardForDefaultPreview(uint(agentStaffIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetSaleListForDefaultPreview(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
			pageNumberStr   = c.QueryParam("page_number")
		)

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		pageNumberInt, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeDashboardHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSaleListForDefaultPreview(uint(agentStaffIDInt), uint(pageNumberInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetReleaseJobSeekerListForDefaultPreview(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
			pageNumberStr   = c.QueryParam("page_number")
		)

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		pageNumberInt, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeDashboardHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetReleaseJobSeekerListForDefaultPreview(uint(agentStaffIDInt), uint(pageNumberInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 検索表示関連
//

// ダッシュボードの条件検索
func GetDashboardForSearch(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {

		// クエリパラム
		searchParam, err := parseSearchDashboardQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeDashboardHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetDashboardForSearch(searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetSaleListForSearch(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			pageNumberStr = c.QueryParam("page_number")
		)

		pageNumberInt, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// クエリパラム
		searchParam, err := parseSearchDashboardQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeDashboardHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSaleListForSearch(uint(pageNumberInt), searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetReleaseJobSeekerListForSearch(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			pageNumberStr = c.QueryParam("page_number")
		)

		pageNumberInt, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// クエリパラム
		searchParam, err := parseSearchDashboardQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeDashboardHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetReleaseJobSeekerListForSearch(uint(pageNumberInt), searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// CSV API
//
// csvファイルの吐き出し
func ExportAccuracyCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.Param("agent_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeDashboardHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.ExportAccuracyCSV(uint(agentID))

		if err != nil {
			fmt.Println(err)
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderFile(c, p)
		// ローカルファイルの削除
		// os.Remove(p)

		return nil
	}
}

/****************************************************************************************/
