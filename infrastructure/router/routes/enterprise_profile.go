package routes

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
	"github.com/spaceaiinc/autoscout-server/infrastructure/di"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

/****************************************************************************************/
/// 汎用系 API
// 企業情報の登録
func CreateEnterpriseProfile(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateEnterpriseProfileParam
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

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.CreateEnterpriseProfile(param)
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
func UpdateEnterpriseProfile(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateEnterpriseProfileParam
		)

		enterpriseIDStr := c.Param("enterprise_id")

		enterpriseIDInt, _ := strconv.Atoi(enterpriseIDStr)

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

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.UpdateEnterpriseProfile(param, uint(enterpriseIDInt))
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
func DeleteEnterpriseProfile(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.DeleteEnterpriseProfileParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.DeleteEnterpriseProfile(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業の参考資料の削除(material_type = 1 or 2)
func DeleteEnterpriseReferenceMaterial(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			referenceMaterialIDStr = c.Param("reference_material_id")
			materialTypeStr        = c.Param("material_type")
		)
		referenceMaterialIDInt, err := strconv.Atoi(referenceMaterialIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		materialTypeInt, err := strconv.Atoi(materialTypeStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.DeleteEnterpriseReferenceMaterial(uint(referenceMaterialIDInt), uint(materialTypeInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDから企業情報を取得
func GetEnterpriseByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.GetEnterpriseByID(uint(enterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 担当者IDで企業情報一覧を取得
func GetEnterpriseListByAgentStaffID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.GetEnterpriseListByAgentStaffID(uint(agentStaffIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDで企業情報一覧を取得
func GetEnterpriseListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.GetEnterpriseListByAgentID(uint(agentIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDとページ番号で50件取得
func GetEnterpriseListByAgentIDAndPage(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr    = c.Param("agent_id")
			pageNumberStr = c.QueryParam("page_number")
		)

		agentIDInt, err := strconv.Atoi(agentIDStr)
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

		if pageNumberInt < 1 || agentIDInt < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number or agent_id is invalid", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.GetEnterpriseListByAgentIDAndPage(uint(agentIDInt), uint(pageNumberInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDとクエリパラムで企業一覧を絞り込み
func GetSearchEnterpriseListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.Param("agent_id")
			pageNumber = c.QueryParam("page_number")
		)

		agentID, err := strconv.Atoi(agentIDStr)
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

		if pageNumberInt < 1 || agentID < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number or agent_id is invalid", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// クエリパラム
		searchParam, err := parseSearchEnterpriseQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.GetSearchEnterpriseListByAgentID(uint(agentID), uint(pageNumberInt), searchParam)
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
func CreateEnterpriseReferenceMaterial(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateEnterpriseReferenceMaterialParam
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

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.CreateEnterpriseReferenceMaterial(param)
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
func UpdateEnterpriseReferenceMaterial(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateEnterpriseReferenceMaterialParam
		)

		enterpriseIDStr := c.Param("enterprise_id")

		enterpriseIDInt, _ := strconv.Atoi(enterpriseIDStr)

		param.EnterpriseID = uint(enterpriseIDInt)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		fmt.Println(param)

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.UpdateEnterpriseReferenceMaterial(param)
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

// 企業IDから企業資料を取得
func GetEnterpriseReferenceMaterialByEnterpriseID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.GetEnterpriseReferenceMaterialByEnterpriseID(uint(enterpriseIDInt))
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
/// 企業の追加情報 API
//
// 企業の追加情報を作成
func CreateEnterpriseActivity(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var param entity.CreateEnterpriseActivityParam

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			fmt.Println(err)
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			fmt.Println(err)
			return err
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.CreateEnterpriseActivity(param)
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

/****************************************************************************************/
/****************************************************************************************/
// CSV操作　API
//
// csvファイルの読み込み 企業・請求先
func ImportEnterpriseCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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
		enterpriseList, missedRecords, err := parseEnterpriseCSV(r)
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		log.Println("enterpriseList", enterpriseList)

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.ImportEnterpriseCSV(enterpriseList, missedRecords, uint(agentID))
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

// csvファイルの読み込み 企業・請求先　*プレビュー用
func PreviewEnterpriseCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {

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
		enterpriseList, missedRecords, err := parseEnterpriseCSV(r)
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// リストのまま返す
		p := presenter.NewMissedRecordsAndEnterpriseListJSONPresenter(
			responses.NewMissedRecordsAndEnterpriseList(missedRecords, enterpriseList),
		)
		renderJSON(c, p)
		return nil
	}
}

// csvファイルの読み込み 求人企業
func ImportJobInformationCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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
		jobInformationList, missedRecords, err := parseJobInformationCSV(r)
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.ImportJobInformationCSV(jobInformationList, missedRecords, uint(agentID))
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

// csvファイルの読み込み 求人企業 *プレビュー用
func PreviewJobInformationCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {

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
		jobInformationList, missedRecords, err := parseJobInformationCSV(r)
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// リストのまま返す
		p := presenter.NewMissedRecordsAndJobInformationListJSONPresenter(
			responses.NewMissedRecordsAndJobInformationList(missedRecords, jobInformationList),
		)

		renderJSON(c, p)
		return nil
	}
}

// csvファイルの読み込み 企業・請求先 サーカス
func ImportEnterpriseCSVForCircus(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		// csvファイルの拡張子チェック
		if !strings.HasSuffix(fileParam.Filename, ".csv") {
			wrapped := fmt.Errorf("ファイルの拡張子がcsvではありません: %s:%w", fileParam.Filename, entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		r := csv.NewReader(file)
		file.Close()

		r.TrimLeadingSpace = true // true の場合は、先頭の空白文字を無視する
		r.ReuseRecord = true      // true の場合は、Read で戻ってくるスライスを次回再利用する。パフォーマンスが上がる

		//csvの企業リストを、enterprise型のリストに変換
		enterpriseList, missedRecords, err := parseEnterpriseCSVForCircus(r, uint(agentStaffID))
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		log.Println("enterpriseList", enterpriseList)

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.ImportEnterpriseCSVForCircus(enterpriseList, missedRecords, uint(agentID))
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

// csvファイルの読み込み 企業・請求先 サーカス
func ImportEnterpriseCSVForAgentBank(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		// csvファイルの拡張子チェック
		if !strings.HasSuffix(fileParam.Filename, ".csv") {
			wrapped := fmt.Errorf("ファイルの拡張子がcsvではありません: %s:%w", fileParam.Filename, entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		r := csv.NewReader(file)
		file.Close()

		r.TrimLeadingSpace = true // true の場合は、先頭の空白文字を無視する
		r.ReuseRecord = true      // true の場合は、Read で戻ってくるスライスを次回再利用する。パフォーマンスが上がる

		//csvの企業リストを、enterprise型のリストに変換
		enterpriseList, missedRecords, err := parseEnterpriseCSVForAgentBank(r, uint(agentStaffID))
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		log.Println("enterpriseList", enterpriseList)

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, tx, sendgrid)
		p, err := h.ImportEnterpriseCSVForAgentBank(enterpriseList, missedRecords, uint(agentID))
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

// csvファイルの吐き出し
func ExportEnterpriseCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.ExportEnterpriseCSV(uint(agentID))

		if err != nil {
			fmt.Println(err)
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderFile(c, p)
		// ローカルファイルの削除
		os.Remove(p)

		return nil
	}
}

/****************************************************************************************/
/****************************************************************************************/
/// Admin API

// 全ての企業一覧
func GetInitialEnterprise(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			pageNumberStr = c.QueryParam("page_number")
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

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.GetInitialEnterprise(uint(pageNumberInt))
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
// 他サービスから求人取得 API
//
// // 他媒体の求人企業を一括インポートするテーブル
// func InitialEnterpriseImporter(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
// 	return func(c echo.Context) error {
// 		var (
// 			param entity.InitialEnterpriseImporter
// 		)

// 		if err := bindAndValidate(c, &param); err != nil {
// 			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
// 			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
// 			return wrapped
// 		}

// 		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
// 		p, err := h.InitialEnterpriseImporter(param)

// 		if err != nil {
// 			renderJSON(c, presenter.NewErrorJSONPresenter(err))
// 			return err
// 		}

// 		renderJSON(c, p)
// 		return nil
// 	}
// }

// JSONファイルから企業と求人のリストから作成
func ImportEnterpriseJSON(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentStaffIDStr = c.Param("agent_staff_id")
			// JSON変換用パラム
			converterParam entity.EnterpriseAndJobInformationListParam
		)

		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// ファイルの受け取り
		fileParam, err := c.FormFile("file")
		if err != nil {
			wrapped := fmt.Errorf("ファイルの受け取りエラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// ファイルのオープン
		file, err := fileParam.Open()
		if err != nil {
			wrapped := fmt.Errorf("ファイルが開けません: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// ファイルの拡張子チェック
		if !strings.HasSuffix(fileParam.Filename, ".json") {
			wrapped := fmt.Errorf("ファイルの拡張子がjsonではありません: %s:%w", fileParam.Filename, entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// ファイルの読み込み
		byteData, err := io.ReadAll(file)
		if err != nil {
			wrapped := fmt.Errorf("ファイルの読み込みエラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// jsonのパース
		if err := json.Unmarshal(byteData, &converterParam); err != nil {
			wrapped := fmt.Errorf("jsonのパースエラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeEnterpriseProfileHandler(firebase, db, sendgrid)
		p, err := h.ImportEnterpriseJSON(converterParam.EnterpriseAndJobInformationList, uint(agentStaffID))

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

/****************************************************************************************/
