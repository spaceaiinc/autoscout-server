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
// 求人の作成
func CreateJobInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param               entity.CreateJobInformationParam
			billingAddressIDStr = c.Param("billing_address_id")
		)

		billingAddressIDInt, err := strconv.Atoi(billingAddressIDStr)
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

		h := di.InitializeJobInformationHandler(firebase, tx, sendgrid)
		p, err := h.CreateJobInformation(param, uint(billingAddressIDInt))
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
func UpdateJobInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param               entity.UpdateJobInformationParam
			jobInformationIDStr = c.Param("job_information_id")
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

		h := di.InitializeJobInformationHandler(firebase, tx, sendgrid)
		p, err := h.UpdateJobInformation(param, uint(jobInformationIDInt))
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
func DeleteJobInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobInformationIDStr = c.Param("job_information_id")
		)

		jobInformationIDInt, err := strconv.Atoi(jobInformationIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.DeleteJobInformation(uint(jobInformationIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求人IDから求人情報を取得
func GetJobInformationByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobInformationIDStr = c.Param("job_information_id")
		)

		jobInformationIDInt, err := strconv.Atoi(jobInformationIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationByID(uint(jobInformationIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求人のuuidから求人情報を取得
func GetJobInformationByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			uuidStr = c.Param("job_information_uuid")
		)

		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationByUUID(uuid)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求人のuuidを使って求人票の情報を取得する関数
func GetJobListingByJobInformationUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			uuidStr = c.Param("job_information_uuid")
		)

		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobListingByJobInformationUUID(uuid)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者が確認する求人票情報を取得（求人票 + タスクに紐づいた選考情報）
func GetJobListingForJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobInformationUUIDStr = c.QueryParam("job_information_uuid")
			jobSeekerUUIDStr      = c.QueryParam("job_seeker_uuid")
		)

		jobInformationUUID, err := uuid.Parse(jobInformationUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "job_information_uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "job_seeker_uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobListingForJobSeeker(jobInformationUUID, jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDで求人情報一覧を取得
func GetJobInformationListByBillingAddressID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationListByBillingAddressID(uint(billingAddressIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 企業IDで求人情報一覧を取得
func GetJobInformationListByEnterpriseID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationListByEnterpriseID(uint(enterpriseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// agentIDを使って求人情報一覧を取得する関数
func GetJobInformationListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationListByAgentID(uint(agentIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDとページ番号で50件取得
func GetJobInformationListByIDList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {

		// クエリパラム
		idListParam, err := parseIDListQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationListByIDList(idListParam.IDList)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 求職者検索→求人検索 API
//

// 全ての求人
func GetJobInformationListByAgentIDAndType(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr    = c.Param("agent_id")
			pageNumberStr = c.QueryParam("page_number")
			searchTypeStr = c.QueryParam("type")
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

		searchTypeInt, err := strconv.Atoi(searchTypeStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationListByAgentIDAndType(uint(agentIDInt), uint(pageNumberInt), entity.JobInformationType(searchTypeInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 全ての求人
func GetSearchJobInformationListByAgentIDAndType(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr    = c.Param("agent_id")
			pageNumberStr = c.QueryParam("page_number")
			searchTypeStr = c.QueryParam("type")
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

		// クエリパラム
		searchParam, err := parseSearchJobInformationQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		searchTypeInt, err := strconv.Atoi(searchTypeStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSearchJobInformationListByAgentIDAndType(uint(agentIDInt), uint(pageNumberInt), searchParam, entity.JobInformationType(searchTypeInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// シェア求職者検索→自社求人検索(絞り込み) API
//
func GetSearchPublicJobInformationListByAgentIDAndPage(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr         = c.Param("agent_id")
			pageNumberStr      = c.QueryParam("page_number")
			jobSeekerIDStrList = c.QueryParams()["job_seeker_id_list[]"]
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

		jobSeekerIDList, err := parseQueryParamUINT(jobSeekerIDStrList)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
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

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSearchPublicJobInformationListByAgentIDAndPage(uint(agentIDInt), uint(pageNumberInt), jobSeekerIDList, searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/

// 求人IDを使って選考フローパターンの一覧を取得する
func GetSelectionFlowPatternListByJobInformationID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobInformationIDStr = c.Param("job_information_id")
		)

		jobInformationIDInt, err := strconv.Atoi(jobInformationIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSelectionFlowPatternListByJobInformationID(uint(jobInformationIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求人IDを使って選考フローパターンの一覧を取得する
func GetOpenSelectionFlowPatternListByJobInformationID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobInformationIDStr = c.Param("job_information_id")
		)

		jobInformationIDInt, err := strconv.Atoi(jobInformationIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetOpenSelectionFlowPatternListByJobInformationID(uint(jobInformationIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 絞り込み検索 API
//
func GetSearchActiveJobInformationListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr    = c.Param("agent_id")
			pageNumberStr = c.QueryParam("page_number")
		)

		agentID, err := strconv.Atoi(agentIDStr)
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

		if pageNumberInt < 1 || agentID < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number or agent_id is invalid", entity.ErrRequestError)
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

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSearchActiveJobInformationListByAgentID(uint(agentID), uint(pageNumberInt), searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 選考フロー API
//
func GetSearchJobInformationListByOtherAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr    = c.Param("agent_id")
			pageNumberStr = c.QueryParam("page_number")
		)

		agentID, err := strconv.Atoi(agentIDStr)
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

		if pageNumberInt < 1 || agentID < 1 {
			wrapped := fmt.Errorf("%s:%w", "page_number or agent_id is invalid", entity.ErrRequestError)
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

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSearchJobInformationListByOtherAgentID(uint(agentID), uint(pageNumberInt), searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// IDを使って選考フローパターンを取得する
func GetSelectionFlowPatternByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			selectionFlowIDStr = c.Param("selection_flow_id")
		)

		selectionFlowIDInt, err := strconv.Atoi(selectionFlowIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSelectionFlowPatternByID(uint(selectionFlowIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 選考フローパターンを更新
func CreateSelectionFlowPattern(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateAndUpdateSelectionFlowPatternParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.CreateSelectionFlowPattern(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 選考フローパターンを更新
func UpdateSelectionFlowPattern(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param              entity.CreateAndUpdateSelectionFlowPatternParam
			selectionFlowIDStr = c.Param("selection_flow_id")
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		selectionFlowIDInt, err := strconv.Atoi(selectionFlowIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.UpdateSelectionFlowPattern(param, uint(selectionFlowIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 選考フローパターンを削除
func DeltedSelectionFlowPattern(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			selectionFlowIDStr = c.Param("selection_flow_id")
		)

		selectionFlowIDInt, err := strconv.Atoi(selectionFlowIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.DeltedSelectionFlowPattern(uint(selectionFlowIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者uuidで打診されている求人一覧を取得する
func GetJobListingListByJobSeekerUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobListingListByJobSeekerUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// LP用 API
// LPの診断ページから合致求人数と年収期待値を取得
func GetSearchJobInformationCountByLPDiagnosis(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.DiagnosisParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSearchJobInformationCountByLPDiagnosis(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// LPの診断ページから合致求人数と年収期待値を取得
func GetSearchJobListingListByJobSeekerUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SearchMatchingJobListParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetSearchJobListingListByJobSeekerUUID(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// LPから診断に使用する求人を取得
func GetJobInformationListForDiagnosis(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobInformationListForDiagnosis()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// LPから診断に使用する求人を取得
func GetJobListingListAndJobSeekerDesiredForDiagnosis(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return nil
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobListingListAndJobSeekerDesiredForDiagnosis(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者のエントリー希望と興味あり求人の取得
func GetJobListingListByJobSeekerUUIDAndInterestedType(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.InterestedTypeJobListParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetJobListingListByJobSeekerUUIDAndInterestedType(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
/// Admin API

// 全ての求人一覧
func GetAllJobInformation(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeJobInformationHandler(firebase, db, sendgrid)
		p, err := h.GetAllJobInformation(uint(pageNumberInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}
