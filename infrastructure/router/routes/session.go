package routes

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/di"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SignInParam struct {
	Token string `json:"token" validate:"required"`
}

type SignInPasswordParam struct {
	Password string `json:"password" validate:"required"`
}

type SignInLoginTokenParam struct {
	LoginToken string `json:"login_token" validate:"required"`
}

type LoginForLPParam struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CheckSignInParam struct {
	FirebaseID string `json:"firebase_id" validate:"required"`
}

// SignIn
func SignIn(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param = new(SignInParam)
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.SignIn(param.Token)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// Signout
func SignOut(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param = new(SignInParam)
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.SignOut(param.Token)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

func GetSignInUser(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
		)

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.GetSignInUser(firebaseToken)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// SignIn
func SignInForGestEnterprise(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param   = new(SignInPasswordParam)
			uuidStr = c.Param("job_information_uuid")
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		jobInformationUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.SignInForGestEnterprise(param.Password, jobInformationUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// SignIn
func SignInForGestEnterpriseByTaskGroupUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param   = new(SignInPasswordParam)
			uuidStr = c.Param("task_group_uuid")
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		taskGroupUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.SignInForGestEnterpriseByTaskGroupUUID(param.Password, taskGroupUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// SignIn
func SignInForGestJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param   = new(SignInPasswordParam)
			uuidStr = c.Param("job_seeker_uuid")
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		jobSeekerUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			// uuidはクエリに含まれているものを使用するため、「URLが正しくありません。ご確認の上もう一度お試しください。」を返す
			wrapped := fmt.Errorf("URLが正しくありません。ご確認の上もう一度お試しください。")
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.SignInForGestJobSeeker(param.Password, jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// マイページログイン（LPからログイントークンを使ってログイン）
func SignInForGestJobSeekerFromLP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param   = new(SignInLoginTokenParam)
			uuidStr = c.Param("job_seeker_uuid")
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		loginToken, err := uuid.Parse(param.LoginToken)
		if err != nil {
			// uuidはクエリに含まれているものを使用するため、「URLが正しくありません。ご確認の上もう一度お試しください。」を返す
			wrapped := fmt.Errorf("URLが正しくありません。ご確認の上もう一度お試しください。")
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		jobSeekerUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			// uuidはクエリに含まれているものを使用するため、「URLが正しくありません。ご確認の上もう一度お試しください。」を返す
			wrapped := fmt.Errorf("URLが正しくありません。ご確認の上もう一度お試しください。")
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.SignInForGestJobSeekerFromLP(jobSeekerUUID, loginToken)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// SignIn
func SignInForGestSendingJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param   = new(SignInPasswordParam)
			uuidStr = c.Param("sending_job_seeker_uuid")
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		jobSeekerUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.SignInForGestSendingJobSeeker(param.Password, jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// LPのログインフォームのログイン処理
func LoginGestJobSeekerForLP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param = new(LoginForLPParam)
		)

		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeSessionHandler(firebase, db, sendgrid)
		p, err := h.LoginGestJobSeekerForLP(param.Email, param.Password)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}
