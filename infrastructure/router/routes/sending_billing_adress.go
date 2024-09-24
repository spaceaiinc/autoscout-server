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
// 請求先情報の更新
func UpdateSendingBillingAddress(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateSendingBillingAddressParam
		)

		billingAddressIDStr := c.Param("sending_billing_address_id")

		billingAddressIDInt, _ := strconv.Atoi(billingAddressIDStr)

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

		h := di.InitializeSendingBillingAddressHandler(firebase, tx, sendgrid)
		p, err := h.UpdateSendingBillingAddress(param, uint(billingAddressIDInt))
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

// 請求先IDから請求先情報を取得
func GetSendingBillingAddressByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			billingAddressIDStr = c.Param("sending_billing_address_id")
		)

		billingAddressIDInt, err := strconv.Atoi(billingAddressIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetSendingBillingAddressByID(uint(billingAddressIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDから請求先情報を取得
func GetSendingBillingAddressBySendingEnterpriseID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingEnterpriseIDStr = c.Param("sending_enterprise_id")
		)

		sendingEnterpriseIDInt, err := strconv.Atoi(sendingEnterpriseIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetSendingBillingAddressBySendingEnterpriseID(uint(sendingEnterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetSendingBillingAddressListBySendingEnterpriseIDList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingEnterpriseIDStrList = c.QueryParams()["id_list[]"]
		)
		sendingEnterpriseIDList, err := parseQueryParamUINT(sendingEnterpriseIDStrList)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		fmt.Println("sendingEnterpriseIDList", len(sendingEnterpriseIDList))

		h := di.InitializeSendingBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetSendingBillingAddressListBySendingEnterpriseIDList(sendingEnterpriseIDList)
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
/// Admin API
//
//GetAllSendingBillingAddress
func GetAllSendingBillingAddress(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		h := di.InitializeSendingBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetAllSendingBillingAddress()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
