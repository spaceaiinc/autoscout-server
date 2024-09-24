package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/di"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

func GetGoogleAuthCodeURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, googleAPI config.GoogleAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
		)

		h := di.InitializeGoogleAuthenticationHandler(firebase, db, sendgrid, googleAPI)
		p, err := h.GetGoogleAuthCodeURL(firebaseToken)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

func UpdateGoogleOauthToken(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, googleAPI config.GoogleAPI) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
			param         entity.GoogleAuthParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeGoogleAuthenticationHandler(firebase, db, sendgrid, googleAPI)
		p, err := h.UpdateGoogleOauthToken(firebaseToken, param.AuthCode)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}
