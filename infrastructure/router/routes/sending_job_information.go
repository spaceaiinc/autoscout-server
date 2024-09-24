package routes

import (
	"encoding/csv"
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
// 求人の作成
func CreateSendingJobInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateSendingJobInformationParam
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

		h := di.InitializeSendingJobInformationHandler(firebase, tx, sendgrid)
		p, err := h.CreateSendingJobInformation(param)
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

// 求人の更新
func UpdateSendingJobInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param               entity.UpdateSendingJobInformationParam
			jobInformationIDStr = c.Param("sending_job_information_id")
		)

		jobInformationIDInt, err := strconv.Atoi(jobInformationIDStr)

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

		h := di.InitializeSendingJobInformationHandler(firebase, tx, sendgrid)
		p, err := h.UpdateSendingJobInformation(param, uint(jobInformationIDInt))
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

// 求人情報の削除
func DeleteSendingJobInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobInformationIDStr = c.Param("sending_job_information_id")
		)

		jobInformationIDInt, err := strconv.Atoi(jobInformationIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobInformationHandler(firebase, db, sendgrid)
		p, err := h.DeleteSendingJobInformation(uint(jobInformationIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求人IDから求人情報を取得
func GetSendingJobInformationByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobInformationIDStr = c.Param("sending_job_information_id")
		)

		jobInformationIDInt, err := strconv.Atoi(jobInformationIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSendingJobInformationByID(uint(jobInformationIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求人のuuidから求人情報を取得
func GetSendingJobInformationByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			uuidStr = c.Param("sending_job_information_uuid")
		)

		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeSendingJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSendingJobInformationByUUID(uuid)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求人のuuidから求人情報を取得
func GetJobListingBySendingJobInformationUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			uuidStr = c.Param("sending_job_information_uuid")
		)

		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeSendingJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobListingBySendingJobInformationUUID(uuid)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDで求人情報一覧を取得
func GetSendingJobInformationListBySendingEnterpriseID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeSendingJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSendingJobInformationListBySendingEnterpriseID(uint(sendingEnterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// まだ送客していない送客先のリストと指定IDの送客先が保有する求人の一覧と最大ページ数を取得するapi
func GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr  = c.QueryParam("sending_job_seeker_id")
			sendingEnterpriseIDStr = c.QueryParam("sending_enterprise_id")
			pageNumberStr          = c.QueryParam("page_number")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		sendingEnterpriseIDInt, err := strconv.Atoi(sendingEnterpriseIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

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

		h := di.InitializeSendingJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(uint(sendingJobSeekerIDInt), uint(sendingEnterpriseIDInt), uint(pageNumberInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 絞り込み まだ送客していない送客先のリストと指定IDの送客先が保有する求人の一覧と最大ページ数を取得するapi
func GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingJobSeekerIDStr  = c.QueryParam("sending_job_seeker_id")
			sendingEnterpriseIDStr = c.QueryParam("sending_enterprise_id")
			pageNumberStr          = c.QueryParam("page_number")
		)

		sendingJobSeekerIDInt, err := strconv.Atoi(sendingJobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		sendingEnterpriseIDInt, err := strconv.Atoi(sendingEnterpriseIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

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

		// クエリパラム
		searchParam, err := parseSearchJobInformationQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(uint(sendingJobSeekerIDInt), uint(sendingEnterpriseIDInt), uint(pageNumberInt), searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// csvファイルの読み込み
func ImportSendingJobInformationCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingBillingAddressIDStr = c.Param("sending_billing_address_id")
		)

		sendingBillingAddressIDInt, err := strconv.Atoi(sendingBillingAddressIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		fileParam, err := c.FormFile("file")
		if err != nil {
			wrapped := fmt.Errorf("ファイルの受け取りエラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		file, err := fileParam.Open()
		if err != nil {
			wrapped := fmt.Errorf("ファイルが開けません: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		r := csv.NewReader(file)
		file.Close()

		r.TrimLeadingSpace = true // true の場合は、先頭の空白文字を無視する
		r.ReuseRecord = true      // true の場合は、Read で戻ってくるスライスを次回再利用する。パフォーマンスが上がる

		//csvの企業リストを、enterprise型のリストに変換
		sendingJobInformationList, missedRecords, err := parseSendingJobInformationCSV(r)
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// 請求先IDをセット
		for _, sendingJobInformation := range sendingJobInformationList {
			sendingJobInformation.SendingBillingAddressID = uint(sendingBillingAddressIDInt)
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeSendingJobInformationHandler(firebase, tx, sendgrid)
		p, err := h.ImportSendingJobInformationCSV(sendingJobInformationList, missedRecords)
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
