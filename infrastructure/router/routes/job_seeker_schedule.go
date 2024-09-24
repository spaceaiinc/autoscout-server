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
// 求職者情報の登録
func CreateJobSeekerSchedule(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateJobSeekerScheduleParam
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

		h := di.InitializeJobSeekerScheduleHandler(firebase, tx, sendgrid)
		p, err := h.CreateJobSeekerSchedule(param)
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

func UpdateJobSeekerSchedule(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param                  entity.CreateOrUpdateJobSeekerScheduleParam
			jobSeekerScheduleIDStr = c.Param("schedule_id")
		)

		jobSeekerScheduleIDInt, err := strconv.Atoi(jobSeekerScheduleIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
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

		h := di.InitializeJobSeekerScheduleHandler(firebase, tx, sendgrid)
		p, err := h.UpdateJobSeekerSchedule(uint(jobSeekerScheduleIDInt), param)
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

// CA担当者に新規追加した日程を追加
func ShareScheduleToCAStaff(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param          entity.ShareScheduleParam
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
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

		h := di.InitializeJobSeekerScheduleHandler(firebase, tx, sendgrid)
		p, err := h.ShareScheduleToCAStaff(uint(jobSeekerIDInt), param)
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

// 求職者IDでスケジュール情報を取得
func GetJobSeekerScheduleListByJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		jobSeekerIDStr := c.Param("job_seeker_id")

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerScheduleHandler(firebase, db, sendgrid)
		p, err := h.GetJobSeekerScheduleListByJobSeekerID(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者IDとスケジュールタイプでスケジュール情報を取得
func GetJobSeekerScheduleTypeListByJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr  = c.Param("job_seeker_id")
			scheduleTypeStr = c.QueryParam("schedule_type")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		scheduleTypeInt, err := strconv.Atoi(scheduleTypeStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerScheduleHandler(firebase, db, sendgrid)
		p, err := h.GetJobSeekerScheduleTypeListByJobSeekerID(uint(jobSeekerIDInt), uint(scheduleTypeInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func DeleteJobSeekerSchedule(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		jobSeekerScheduleIDStr := c.Param("schedule_id")

		jobSeekerScheduleIDInt, err := strconv.Atoi(jobSeekerScheduleIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerScheduleHandler(firebase, db, sendgrid)
		p, err := h.DeleteJobSeekerSchedule(uint(jobSeekerScheduleIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
