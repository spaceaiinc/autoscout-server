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
// 企業情報の登録
func CreateSendingEnterprise(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateSendingEnterpriseParam
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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.CreateSendingEnterprise(param)
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

// 企業情報の更新
func UpdateSendingEnterprise(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateSendingEnterpriseParam
		)

		sendingEnterpriseIDStr := c.Param("sending_enterprise_id")

		sendingEnterpriseIDInt, _ := strconv.Atoi(sendingEnterpriseIDStr)

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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.UpdateSendingEnterprise(param, uint(sendingEnterpriseIDInt))
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

func UpdateSendingEnterprisePassword(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateSendingEnterprisePasswordParam
		)

		sendingEnterpriseIDStr := c.Param("sending_enterprise_id")
		sendingEnterpriseIDInt, err := strconv.Atoi(sendingEnterpriseIDStr)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.UpdateSendingEnterprisePassword(uint(sendingEnterpriseIDInt), param.Password)
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

func SigninSendingEnterprise(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SigninSendingEnterprisePasswordParam
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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.SigninSendingEnterprise(param)
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

// 企業情報の削除
func DeleteSendingEnterprise(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.DeleteSendingEnterprise(uint(sendingEnterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDから企業情報を取得
func GetSendingEnterpriseByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingEnterpriseByID(uint(sendingEnterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDから企業情報を取得
func GetSendingEnterpriseAndBillingAddressByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingEnterpriseAndBillingAddressByID(uint(sendingEnterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDから企業情報を取得
func GetSendingEnterpriseByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingEnterpriseUUIDStr = c.Param("sending_enterprise_uuid")
		)

		sendingEnterpriseUUID, err := uuid.Parse(sendingEnterpriseUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingEnterpriseByUUID(sendingEnterpriseUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 資料関連API
//
// 企業資料の作成
func CreateSendingEnterpriseReferenceMaterial(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam
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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.CreateSendingEnterpriseReferenceMaterial(param)
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

// 企業資料の更新
func UpdateSendingEnterpriseReferenceMaterial(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam
		)

		sendingEnterpriseIDStr := c.Param("sending_enterprise_id")

		sendingEnterpriseIDInt, _ := strconv.Atoi(sendingEnterpriseIDStr)

		param.SendingEnterpriseID = uint(sendingEnterpriseIDInt)

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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.UpdateSendingEnterpriseReferenceMaterial(param)
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

// 送客先エージェントの参考資料の削除 sending_enterprise_id, file_type: '送客先エージェント画像' | '参考資料1' | '参考資料2'
func DeleteSendingEnterpriseReferenceMaterial(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingEnterpriseIDStr = c.Param("sending_enterprise_id")
			fileTypeStr            = c.Param("file_type")
		)

		sendingEnterpriseIDInt, err := strconv.Atoi(sendingEnterpriseIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		isMatched := false
		for _, fileType := range entity.SendingEnterpriseFileType {
			if fileType == fileTypeStr {
				isMatched = true
				break
			}
		}
		if !isMatched {
			wrapped := fmt.Errorf("%s:%w", "file_typeが正しくありません", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.DeleteSendingEnterpriseReferenceMaterial(uint(sendingEnterpriseIDInt), fileTypeStr)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDから企業資料を取得
func GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(uint(sendingEnterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// ページの全ての企業一覧
func GetAllSendingEnterpriseByPageAndFreeWord(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			pageNumberStr = c.QueryParam("page_number")
			freeWord      = c.QueryParam("free_word")
		)

		pageNumberInt, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		if pageNumberInt < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number is invalid", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.GetAllSendingEnterpriseByPageAndFreeWord((uint(pageNumberInt)), freeWord)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 送客メール送信時に必要な情報をまとめて取得するapi
func GetSendingInformationForSendingMail(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobInformationIDStrList = c.QueryParams()["sending_job_information_id_list[]"]
		)

		sendingJobInformationIDList, err := parseQueryParamUINT(sendingJobInformationIDStrList)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingInformationForSendingMail(sendingJobInformationIDList)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func SendSendingMail(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendSendingMailParam
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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.SendSendingMail(param)
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

func SendMailForRSVP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendSendingMailForRSVPParam
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

		h := di.InitializeSendingEnterpriseHandler(firebase, tx, sendgrid)
		p, err := h.SendMailForRSVP(param)
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

func GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.Param("agent_id")
		)

		agentIDInt, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingEnterpriseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(uint(agentIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
