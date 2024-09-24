package routes

import (
	"fmt"
	"strconv"
	"time"

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
// 汎用系 API
//
// 作成 *作成時は名前と連絡先のみで、それ以外の項目は更新時に作成する
func CreateSendingJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateSendingJobSeekerParam
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.CreateSendingJobSeeker(param)
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

// CRM求職者から送客求職者を作成
func CreateSendingJobSeekerFromJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.JobSeeker
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.CreateSendingJobSeekerFromJobSeeker(param)
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

// 面談前アンケート
func CreateSendingInitialQuestionnaire(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateSendingInitialQuestionnaireParam
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.CreateSendingInitialQuestionnaire(param)
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

// 作成後最初の更新 *作成時は名前と連絡先のみで、それ以外の項目は更新時に作成する
func FirstUpdateSendingJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingCustomerIDStr = c.Param("sending_customer_id")
			param                entity.FirstUpdateSendingJobSeekerParam
		)

		sendingCustomerID, err := strconv.Atoi(sendingCustomerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		if sendingCustomerID == 0 {
			wrapped := fmt.Errorf("%s:%w", "sending_custmomer_id is required", entity.ErrRequestError)
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.FirstUpdateSendingJobSeeker(uint(sendingCustomerID), param)
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

// 更新
func UpdateSendingJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
			param                 entity.UpdateSendingJobSeekerParam
		)

		sendingJobSeekerID, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		if sendingJobSeekerID == 0 {
			wrapped := fmt.Errorf("%s:%w", "sending_job_seeker_id is required", entity.ErrRequestError)
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.UpdateSendingJobSeeker(uint(sendingJobSeekerID), param)
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
type UpdateSendingJobSeekerPhaseParam struct {
	Phase              uint `json:"phase" validate:"required"`
	SendingJobSeekerID uint `json:"sending_job_seeker_id" validate:"required"`
}

func UpdateSendingJobSeekerPhase(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param UpdateSendingJobSeekerPhaseParam
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.UpdateSendingJobSeekerPhase(param.SendingJobSeekerID, param.Phase)
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

// 面談日時を更新
type UpdateSendingJobSeekerInterviewDateParam struct {
	InterviewDate      time.Time `json:"interview_date" validate:"required"`
	SendingJobSeekerID uint      `json:"sending_job_seeker_id" validate:"required"`
}

func UpdateSendingInterviewDateBySendingJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param UpdateSendingJobSeekerInterviewDateParam
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.UpdateSendingInterviewDateBySendingJobSeekerID(param.SendingJobSeekerID, param.InterviewDate)
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

// アクティビティーメモの更新
// 面談日時を更新
type UpdateSendingJobSeekerActivityMemoParam struct {
	SendingJobSeekerID uint   `json:"sending_job_seeker_id" validate:"required"`
	ActivityMemo       string `json:"activity_memo" validate:"required"`
}

func UpdateSendingJobSeekerActivityMemo(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param UpdateSendingJobSeekerActivityMemoParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			tx.Rollback()
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.UpdateSendingJobSeekerActivityMemo(param.SendingJobSeekerID, param.ActivityMemo)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		tx.Commit()
		renderJSON(c, p)
		return nil
	}
}

// 面談実施待ちの未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
func UpdateIsVewForWating(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerID, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			tx.Rollback()
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.UpdateIsVewForWating(uint(sendingJobSeekerID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		tx.Commit()
		renderJSON(c, p)
		return nil
	}
}

// 未登録の未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
func UpdateIsVewForUnregister(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerID, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			tx.Rollback()
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.UpdateIsVewForUnregister(uint(sendingJobSeekerID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		tx.Commit()
		renderJSON(c, p)
		return nil
	}
}

// 指定IDの送客求職者を削除
func DeleteSendingJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "sending_job_seeker_id is required", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeeker(uint(sendingJobSeekerIDInt))
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

// 指定IDの送客求職者を取得
func GetSendingJobSeekerByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerID, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSendingJobSeekerByID(uint(sendingJobSeekerID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 指定uuidの送客求職者を取得
func GetSendingJobSeekerByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerUUIDStr = c.Param("sending_job_seeker_uuid")
		)

		sendingJobSeekerUUID, err := uuid.Parse(sendingJobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSendingJobSeekerByUUID(sendingJobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetIsNotViewSendingJobSeekerCountByAgentStaffID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agebtStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffIDInt, err := strconv.Atoi(agebtStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "agent_staff_id is required", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetIsNotViewSendingJobSeekerCountByAgentStaffID(uint(agentStaffIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetSearchListForSendingJobSeekerManagementByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.Param("agent_id")
		)

		agentIDInt, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "agent_id is required", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSearchListForSendingJobSeekerManagementByAgentID(uint(agentIDInt))
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
// 書類関連API
//
// 送客求職者の書類情報更新
func UpdateSendingJobSeekerDocument(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateSendingJobSeekerDocumentParam
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

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.UpdateSendingJobSeekerDocument(param)
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

// 送客求職者の書類情報取得
func GetSendingJobSeekerDocumentBySendingJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerID, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetSendingJobSeekerDocumentBySendingJobSeekerID(uint(sendingJobSeekerID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 履歴書PDFの削除（resume_pdf_urlカラムを空文字で更新）
func DeleteSendingJobSeekerResumePDFURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerResumePDFURL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 履歴書原本の削除（resume_origin_urlカラムを空文字で更新）
func DeleteSendingJobSeekerResumeOriginURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerResumeOriginURL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 職務経歴書PDFの削除（cv_pdf_urlカラムを空文字で更新）
func DeleteSendingJobSeekerCVPDFURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerCVPDFURL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 職務経歴書原本の削除（cv_origin_urlカラムを空文字で更新）
func DeleteSendingJobSeekerCVOriginURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerCVOriginURL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 推薦状PDFの削除（recommendation_pdf_urlカラムを空文字で更新）
func DeleteSendingJobSeekerRecommendationPDFURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerRecommendationPDFURL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 推薦状原本の削除（recommendation_origin_urlカラムを空文字で更新）
func DeleteSendingJobSeekerRecommendationOriginURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerRecommendationOriginURL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 証明写真の削除（id_photo_urlカラムを空文字で更新）
func DeleteSendingJobSeekerIDPhotoURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerIDPhotoURL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// その他①の削除（other_document1_urlカラムを空文字で更新）
func DeleteSendingJobSeekerOtherDocument1URL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerOtherDocument1URL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// その他②の削除（other_document2_urlカラムを空文字で更新）
func DeleteSendingJobSeekerOtherDocument2URL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerOtherDocument2URL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// その他③の削除（other_document3_urlカラムを空文字で更新）
func DeleteSendingJobSeekerOtherDocument3URL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.DeleteSendingJobSeekerOtherDocument3URL(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 送客の終了理由
//
// 送客の終了理由を作成
func CreateSendingJobSeekerEndStatus(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) (err error) {
		var (
			param entity.CreateSendingJobSeekerEndStatusParam
		)

		if err = bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		h := di.InitializeSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.CreateSendingJobSeekerEndStatus(param)
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
