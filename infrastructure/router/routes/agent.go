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
// エージェント情報の取得
func GetAgentByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAgentByAgentID(uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェント情報の取得
func GetAgentByAgentUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentUUIDStr = c.Param("agent_uuid")
		)

		agentUUID, err := uuid.Parse(agentUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAgentByAgentUUID(agentUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェント情報の更新
func UpdateAgent(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.Param("agent_id")
			param      entity.CreateOrUpdateAgentParam
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

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgent(uint(agentID), param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func UpdateAgentAgreementFileURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.AgentAgreementFileURLParam
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

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentAgreementFileURL(param)
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

// CRM機能の有無を更新
func UpdateAgentForAdmin(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.AgentForAdminParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentForAdmin(param)
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
// エージェント情報の登録
func AgentSignUp(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateAgentParam
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

		h := di.InitializeAgentHandler(firebase, tx, sendgrid)
		p, err := h.AgentSignUp(param)
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

// エージェント情報の登録
func AgentAndAgentStaffSignUp(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateAgentAndAgentStaffParam
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

		h := di.InitializeAgentHandler(firebase, tx, sendgrid)
		p, err := h.AgentAndAgentStaffSignUp(param)
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

// 全てのエージェント一覧
func GetAllAgentList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAllAgentList()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定エージェントIDのアライアンスエージェントリスト
func GetAllianceAgentListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAllianceAgentListByAgentID(uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定エージェントIDのアライアンスエージェントリスト（セレクト用）
func GetAllianceAgentListByAgentIDForSelect(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAllianceAgentListByAgentIDForSelect(uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定エージェントIDのアライアンスエージェントリスト（セレクト用）
func GetAgreementFileURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
		)

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAgreementFileURL(firebaseToken)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil

	}
}

// 指定エージェントIDのアライアンスエージェントリスト（セレクト用）
func GetAgentClaimList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			pageNumberStr = c.QueryParam("page_number")
			firebaseToken = GetFirebaseToken(c)
		)
		pageNumberInt, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		if pageNumberInt < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number or agent_id is invalid", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAgentClaimList(firebaseToken, uint(pageNumberInt))
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

// エージェントLINEのアカウント情報を登録
// func CreateAgentLineInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
// 	return func(c echo.Context) error {
// 		var (
// 			param entity.CreateOrUpdateAgentLineInformation
// 		)

// 		if err := bindAndValidate(c, &param); err != nil {
// 			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
// 			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
// 			return wrapped
// 		}

// 		tx, err := db.Begin()
// 		if err != nil {
// 			renderJSON(c, presenter.NewErrorJSONPresenter(err))
// 			return err
// 		}

// 		h := di.InitializeAgentHandler(firebase, tx, sendgrid)
// 		p, err := h.CreateAgentLineInformation(param)
// 		if err != nil {
// 			tx.Rollback()
// 			renderJSON(c, presenter.NewErrorJSONPresenter(err))
// 			return err
// 		}

// 		tx.Commit()
// 		renderJSON(c, p)
// 		return nil
// 	}
// }

/****************************************************************************************/

/****************************************************************************************/
/// LINE 関連
//
func UpdateAgentLineChannel(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.AgentLineChannelParam
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

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.UpdateAgentLineChannel(param)
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

func GetAgentLineChannelByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAgentLineChannelByAgentID(uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetAgentLineLoginChannelIDByAgentUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentUUIDStr = c.Param("agent_uuid")
		)

		agentUUID, err := uuid.Parse(agentUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeAgentHandler(firebase, db, sendgrid)
		p, err := h.GetAgentLineLoginChannelIDByAgentUUID(agentUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
