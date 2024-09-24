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
// 請求先情報の登録
func CreateBillingAddress(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param           entity.CreateBillingAddressParam
			enterpriseIDStr = c.Param("enterprise_id")
		)

		enterpriseIDInt, err := strconv.Atoi(enterpriseIDStr)
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

		h := di.InitializeBillingAddressHandler(firebase, tx, sendgrid)
		p, err := h.CreateBillingAddress(param, uint(enterpriseIDInt))
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

// 請求先情報の更新
func UpdateBillingAddress(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateBillingAddressParam
		)

		billingAddressIDStr := c.Param("billing_address_id")

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

		h := di.InitializeBillingAddressHandler(firebase, tx, sendgrid)
		p, err := h.UpdateBillingAddress(param, uint(billingAddressIDInt))
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

// フェーズを更新
type UpdateBillingStaff struct {
	AgentStaffID         uint   `json:"agent_staff_id" validate:"required"`
	BillingAddressIDList []uint `json:"billing_id_list" validate:"required"`
}

// 複数の請求先のagent_staff_idカラムを一括更新
func UpdateBillingAddressStaffIDByIDListAtOnce(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param UpdateBillingStaff
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

		h := di.InitializeBillingAddressHandler(firebase, tx, sendgrid)
		p, err := h.UpdateBillingAddressStaffIDByIDListAtOnce(param.AgentStaffID, param.BillingAddressIDList)
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

// 請求先情報の削除
func DeleteBillingAddress(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			billingAddressIDStr = c.Param("billing_address_id")
		)

		billingAddressIDInt, err := strconv.Atoi(billingAddressIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.DeleteBillingAddress(uint(billingAddressIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 請求先IDから請求先情報を取得
func GetBillingAddressByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			billingAddressIDStr = c.Param("billing_address_id")
		)

		billingAddressIDInt, err := strconv.Atoi(billingAddressIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetBillingAddressByID(uint(billingAddressIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDから請求先情報を取得
func GetBillingAddressListByEnterpriseID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			enterpriseIDStr = c.Param("enterprise_id")
		)

		enterpriseIDInt, err := strconv.Atoi(enterpriseIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetBillingAddressListByEnterpriseID(uint(enterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// AgentIDから請求先情報を取得
func GetBillingAddressListByPageAndAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.Param("agent_id")
			pageNumber = c.QueryParam("page_number")
		)

		agentIDInt, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		pageNumberInt, err := strconv.Atoi(pageNumber)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		if pageNumberInt < 1 || agentIDInt < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number or agent_id is invalid", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetBillingAddressListByPageAndAgentID(uint(agentIDInt), uint(pageNumberInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// AgentIDから請求先情報を取得
func GetSearchBillingAddressListByPageAndAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.Param("agent_id")
			pageNumber = c.QueryParam("page_number")
		)

		agentIDInt, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		pageNumberInt, err := strconv.Atoi(pageNumber)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		if pageNumberInt < 1 || agentIDInt < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number or agent_id is invalid", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// クエリパラム
		searchParam, err := parseSearchBillingAddressQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetSearchBillingAddressListByPageAndAgentID(uint(agentIDInt), uint(pageNumberInt), searchParam)
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
//GetAllBillingAddress
func GetAllBillingAddress(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		h := di.InitializeBillingAddressHandler(firebase, db, sendgrid)
		p, err := h.GetAllBillingAddress()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
