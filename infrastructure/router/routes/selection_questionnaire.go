package routes

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
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

func CreateSelectionQuestionnaire(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param                entity.CreateOrUpdateSelectionQuestionnaireParam
			questionnaireUUIDStr = c.Param("questionnaire_uuid")
		)

		questionnaireUUID, err := uuid.Parse(questionnaireUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
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

		h := di.InitializeSelectionQuestionnaireHandler(firebase, tx, sendgrid)
		p, err := h.CreateSelectionQuestionnaire(param, questionnaireUUID)
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

// 選考後アンケートの更新
func UpdateSelectionQuestionnaire(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param              entity.CreateOrUpdateSelectionQuestionnaireParam
			questionnaireIDStr = c.Param("questionnaire_id")
		)

		questionnaireIDInt, err := strconv.Atoi(questionnaireIDStr)
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

		h := di.InitializeSelectionQuestionnaireHandler(firebase, tx, sendgrid)
		p, err := h.UpdateSelectionQuestionnaire(param, uint(questionnaireIDInt))
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

// 選考後アンケートの情報を取得
func GetSelectionQuestionnaireOrNullByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			questionnaireUUIDStr = c.Param("questionnaire_uuid")
		)

		questionnaireUUID, err := uuid.Parse(questionnaireUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSelectionQuestionnaireHandler(firebase, db, sendgrid)
		p, err := h.GetSelectionQuestionnaireOrNullByUUID(questionnaireUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GenerateSelectionQuestionnaireByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		h := di.InitializeSelectionQuestionnaireHandler(firebase, db, sendgrid)
		p, err := h.GenerateSelectionQuestionnaireByUUID()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetUnansweredQuestionnaireListByJobSeekerUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "求職者uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSelectionQuestionnaireHandler(firebase, db, sendgrid)
		p, err := h.GetUnansweredQuestionnaireListByJobSeekerUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者のアンケート情報を取得する（queryParam: job_seeker_uuid, job_information_uuid, selection_information_id）
func GetQuestionnaireForJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr      = c.QueryParam("job_seeker_uuid")
			jobInformationUUIDStr = c.QueryParam("job_information_uuid")
			selectionPhaseStr     = c.QueryParam("selection_phase")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "求職者uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		jobInformationUUID, err := uuid.Parse(jobInformationUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "求人uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		selectionPhaseInt, err := strconv.Atoi(selectionPhaseStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSelectionQuestionnaireHandler(firebase, db, sendgrid)
		p, err := h.GetQuestionnaireForJobSeeker(jobSeekerUUID, jobInformationUUID, uint(selectionPhaseInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
