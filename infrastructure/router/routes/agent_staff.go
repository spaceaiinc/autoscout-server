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
// 指定AgentIDの担当者一覧取得
func GetAgentStaffListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
			agentIDStr    = c.Param("agent_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetAgentStaffListByAgentID(firebaseToken, uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定AgentIDの担当者一覧取得
func GetAgentStaffListByAgentIDOrderByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr      = c.Param("agent_id")
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetAgentStaffListByAgentIDOrderByID(uint(agentID), uint(agentStaffID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定AgentStaffIDの担当者情報を更新
func UpdateAgentStaff(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param           entity.AgentStaffUpdateParam
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
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

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaff(uint(agentStaffID), param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定AgentStaffIDの担当メールアドレスを更新
func UpdateAgentStaffEmail(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param           entity.AgentStaffEmailUpdateParam
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
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

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaffEmail(param, uint(agentStaffID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定AgentStaffIDの担当パスワードを更新
func UpdateAgentStaffPassword(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param           entity.AgentStaffPasswordUpdateParam
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
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

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaffPassword(param, uint(agentStaffID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定AgentStaffIDの担当パスワードを更新
func UpdateAgentStaffUsageEnd(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var param entity.DeleteAgentStaffParam

		if err := bindAndValidate(c, &param); err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaffUsageEnd(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定AgentStaffIDの担当パスワードを更新
func UpdateAgentStaffUsageReStart(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaffUsageReStart(uint(agentStaffID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// メール通知（求職者）の更新 body: {agent_staff_id, notification_job_seeker}
func UpdateAgentStaffNotificationJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var param entity.UpdateAgentStaffNotificationJobSeekerParam

		if err := bindAndValidate(c, &param); err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaffNotificationJobSeeker(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// メール通知（未処理・未読）の更新 body: {agent_staff_id, notification_unwatched}
func UpdateAgentStaffNotificationUnwatched(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var param entity.UpdateAgentStaffNotificationUnwatchedParam

		if err := bindAndValidate(c, &param); err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaffNotificationUnwatched(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 管理権限の更新 body: {agent_staff_id, authority}
func UpdateAgentStaffAuthority(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var param entity.UpdateAgentStaffAuthorityParam

		if err := bindAndValidate(c, &param); err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentStaffAuthority(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 担当者削除　*firebaseのアイパスを削除&DBのis_deletedをtrueにする
func DeleteAgentStaff(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var param entity.DeleteAgentStaffParam

		if err := bindAndValidate(c, &param); err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.DeleteAgentStaff(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

func GetOhterAgentStaffListByAgentIDAndAllianceAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr         = c.QueryParam("agent_id")
			AllianceAgentIDStr = c.QueryParam("alliance_agent_id")
			AgentStaffIDStr    = c.QueryParam("agent_staff_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		allianceAgentID, err := strconv.Atoi(AllianceAgentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		agentStaffID, err := strconv.Atoi(AgentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetOhterAgentStaffListByAgentIDAndAllianceAgentID(uint(agentID), uint(allianceAgentID), uint(agentStaffID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

func GetAgentStaffListWithSaleNotCreated(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken   = GetFirebaseToken(c)
			agentIDStr      = c.Param("agent_id")
			managementIDStr = c.Param("management_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		managementID, err := strconv.Atoi(managementIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetAgentStaffListWithSaleNotCreated(firebaseToken, uint(agentID), uint(managementID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 利用可能の担当者を取得
func GetAgentStaffListByAgentIDAndUsageStatusAvailable(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
			agentIDStr    = c.Param("agent_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetAgentStaffListByAgentIDAndUsageStatusAvailable(firebaseToken, uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 未削除の担当者を取得
func GetAgentStaffListByAgentIDAndIsDeletedFalse(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
			agentIDStr    = c.Param("agent_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetAgentStaffListByAgentIDAndIsDeletedFalse(firebaseToken, uint(agentID))
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
// 担当者作成
func AgentStaffSignUp(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param      entity.AgentStaffSignUpForAdminParam
			agentIDStr = c.Param("agent_id")
		)

		agentID, err := strconv.Atoi(agentIDStr)
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

		h := di.InitializeAgentStaffHandler(firebase, tx, sendgrid)
		p, err := h.SignUpForAdmin(param, uint(agentID))
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

// 全担当者一覧取得
func GetAllAgentStaffList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetAllAgentStaffList()
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
/// Agent API
//
// ログインユーザー情報を取得
func GetAgentStaffMe(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
		)

		h := di.InitializeAgentStaffHandler(firebase, db, sendgrid)
		p, err := h.GetAgentStaffMe(firebaseToken)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

/****************************************************************************************/
