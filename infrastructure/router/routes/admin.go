package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/di"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
)

type AdminAuthorizeParam struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AdminAuthorize is
func AdminAuthorize(db *database.DB, appConfig config.App) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param = new(AdminAuthorizeParam)
		)
		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAdminHandler(db, appConfig)
		p, err := h.Authorize(param.Username, param.Password)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
