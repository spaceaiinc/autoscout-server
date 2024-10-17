package routes

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
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
// 求職者情報の登録
func CreateJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param           entity.CreateOrUpdateJobSeekerParam
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.CreateJobSeeker(param, uint(agentStaffIDInt))
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

// 求職者情報の更新
func UpdateJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param           entity.CreateOrUpdateJobSeekerParam
			jobSeekerIDStr  = c.Param("job_seeker_id")
			agentStaffIDStr = c.Param("agent_staff_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateJobSeeker(param, uint(jobSeekerIDInt), uint(agentStaffIDInt))
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

// 求職者情報の更新
func UpdateActivityMemoByJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param          entity.ActivityMemoParam
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateActivityMemoByJobSeekerID(param, uint(jobSeekerIDInt))
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

// 求職者のマッチング求人を閲覧可能かを管理する値を更新
func UpdateCanViewMatchingJob(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateJobSeekerCanViewMatchingJobParam
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateCanViewMatchingJob(param)
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

// 求職者情報の削除
func DeleteJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.DeleteJobSeekerParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeeker(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者IDから求職者情報を取得
func GetJobSeekerByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerByID(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者IDから求職者情報を取得
func GetJobSeekerByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerByUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者uuidから求職者の応募書類情報を取得
func GetJobSeekerDocumentByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerDocumentByUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者IDから求職者情報を取得
func GetJobSeekerByTaskGroupUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			uuidStr = c.Param("task_group_uuid")
		)

		taskGroupUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerByTaskGroupUUID(taskGroupUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDとページ番号で50件取得
func GetJobSeekerListByIDList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr = c.QueryParam("agent_id")
		)

		// クエリパラム
		idListParam, err := parseIDListQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		agentIDInt, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerListByIDList(idListParam.IDList, uint(agentIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// クエリパラム（last_name, first_name, last_furigana, first_furigana, email, phone_number）に合致する求職者情報を取得
func GetDuplicateJobSeekerList(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr    = c.QueryParam("agent_id")
			lastName      = c.QueryParam("last_name")
			firstName     = c.QueryParam("first_name")
			lastFurigana  = c.QueryParam("last_furigana")
			firstFurigana = c.QueryParam("first_furigana")
			email         = c.QueryParam("email")
			phoneNumber   = c.QueryParam("phone_number")
		)

		agentID, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetDuplicateJobSeekerList(uint(agentID), lastName, firstName, lastFurigana, firstFurigana, email, phoneNumber)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者編集ページで使用されるセレクトボックスに必要なデータを取得（自社の担当者一覧、アライアンスエージェント一覧、流入経路）
func GetSelectListForCreateOrUpdateJobSeekerByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetSelectListForCreateOrUpdateJobSeekerByAgentID(firebaseToken, uint(agentID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 絞り込み検索API
//
// エージェントIDとクエリパラムで求職者一覧を絞り込み
func GetSearchJobSeekerListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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
		searchParam, err := parseSearchJobSeekerQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetSearchJobSeekerListByAgentID(uint(agentID), uint(pageNumberInt), searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

func GetSearchActiveJobSeekerListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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
		searchParam, err := parseSearchJobSeekerQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetSearchActiveJobSeekerListByAgentID(uint(agentID), uint(pageNumberInt), searchParam)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// エージェントIDとクエリパラムで求職者一覧を絞り込み
func GetSearchAllianceJobSeekerListByAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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
		searchParam, err := parseSearchJobSeekerQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetSearchAllianceJobSeekerListByAgentID(uint(agentID), uint(pageNumberInt), searchParam)
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
func GetSearchJobSeekerListByAgentIDAndType(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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
		searchParam, err := parseSearchJobSeekerQueryParams(c)
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

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetSearchJobSeekerListByAgentIDAndType(uint(agentIDInt), uint(pageNumberInt), searchParam, entity.JobSeekerType(searchTypeInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// シェア求人検索→自社求職者検索(絞り込み) API
//
func GetSearchPublicJobSeekerListByAgentIDAndPage(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr              = c.Param("agent_id")
			pageNumberStr           = c.QueryParam("page_number")
			jobInformationIDStrList = c.QueryParams()["job_information_id_list[]"]
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

		jobInformationIDList, err := parseQueryParamUINT(jobInformationIDStrList)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// クエリパラム
		searchParam, err := parseSearchJobSeekerQueryParams(c)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetSearchPublicJobSeekerListByAgentIDAndPage(uint(agentID), uint(pageNumberInt), jobInformationIDList, searchParam)
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
// 資料関連API
//
// 求職者の資料情報登録
func CreateJobSeekerDocument(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateJobSeekerDocumentParam
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.CreateJobSeekerDocument(param)
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

// 求職者の資料情報更新
func UpdateJobSeekerDocument(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateOrUpdateJobSeekerDocumentParam
		)

		jobSeekerIDStr := c.Param("job_seeker_id")

		jobSeekerIDInt, _ := strconv.Atoi(jobSeekerIDStr)

		param.JobSeekerID = uint(jobSeekerIDInt)

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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateJobSeekerDocument(param)
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

// 求職者資料の更新 + 現在のタスクが「書類選考/応募書類準備」の場合は「書類選考/エントリー依頼」に変える
func UpdateJobSeekerDocumentForTask(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param                 entity.CreateOrUpdateJobSeekerDocumentParam
			jobSeekerUUIDStr      = c.Param("job_seeker_uuid")
			jobInformationUUIDStr = c.Param("job_information_uuid")
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "job_seeker_uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		jobInformationUUID, err := uuid.Parse(jobInformationUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "job_information_uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateJobSeekerDocumentForTask(param, jobSeekerUUID, jobInformationUUID)
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

// 求職者の資料情報取得
func GetJobSeekerDocumentByJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerID, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerDocumentByJobSeekerID(uint(jobSeekerID))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// CSV操作　API
// csvファイルの読み込み
func ImportJobSeekerCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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

		//csvの求職者リストを、jobSeeker型のリストに変換
		jobSeekerList, err := parseJobSeekerCSV(r)
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.ImportJobSeekerCSV(jobSeekerList, uint(agentID))
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

// csvファイルのプレビュー *csvを読み込んでdbには登録せずに、求職者情報のリストを返す
func PreviewJobSeekerCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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

		//csvの求職者リストを、jobSeeker型のリストに変換
		jobSeekerList, err := parseJobSeekerCSV(r)
		if err != nil {
			wrapped := fmt.Errorf("ファイルの変換エラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// jobSeekerリストのまま返す
		p := presenter.NewMissedRecordsAndJobSeekerListJSONPresenter(
			responses.NewMissedRecordsAndJobSeekerList([]uint{}, jobSeekerList),
		)

		renderJSON(c, p)
		return nil
	}
}

// csvファイルの吐き出し
func ExportJobSeekerCSV(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.ExportJobSeekerCSV(uint(agentID))

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
// LINE関連　API
//
// LINEIDの更新
func UpdateJobSeekerLineID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateJobSeekerLineIDParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.UpdateJobSeekerLineID(param)
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

// 全ての求職者一覧
func GetAllJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
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

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetAllJobSeeker(uint(pageNumberInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// 面談前アンケート関連　API
//
func CreateInitialQuestionnaire(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateInitialQuestionnaireParam
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.CreateInitialQuestionnaire(param)
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

// 履歴書PDFの削除（resume_pdf_urlカラムを空文字で更新）
func DeleteJobSeekerResumePDFURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerResumePDFURL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 履歴書原本の削除（resume_origin_urlカラムを空文字で更新）
func DeleteJobSeekerResumeOriginURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerResumeOriginURL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 職務経歴書PDFの削除（cv_pdf_urlカラムを空文字で更新）
func DeleteJobSeekerCVPDFURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerCVPDFURL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 職務経歴書原本の削除（cv_origin_urlカラムを空文字で更新）
func DeleteJobSeekerCVOriginURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerCVOriginURL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 推薦状PDFの削除（recommendation_pdf_urlカラムを空文字で更新）
func DeleteJobSeekerRecommendationPDFURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerRecommendationPDFURL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 推薦状原本の削除（recommendation_origin_urlカラムを空文字で更新）
func DeleteJobSeekerRecommendationOriginURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerRecommendationOriginURL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 証明写真の削除（id_photo_urlカラムを空文字で更新）
func DeleteJobSeekerIDPhotoURL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerIDPhotoURL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// その他①の削除（other_document1_urlカラムを空文字で更新）
func DeleteJobSeekerOtherDocument1URL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerOtherDocument1URL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// その他②の削除（other_document2_urlカラムを空文字で更新）
func DeleteJobSeekerOtherDocument2URL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerOtherDocument2URL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// その他③の削除（other_document3_urlカラムを空文字で更新）
func DeleteJobSeekerOtherDocument3URL(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.DeleteJobSeekerOtherDocument3URL(uint(jobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
// ゲストページ用 API
//

// 求職者UUIDから求職者情報（ゲスト用）を取得
func GetJobSeekerForInitialStepByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerForInitialStepByUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者UUIDから求職者情報（ゲスト用）を取得
func GetGuestJobSeekerForByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetGuestJobSeekerForByUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者UUIDから求職者情報（ゲスト用）を取得
func GetJobSeekerDesiredForGuestByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerDesiredForGuestByUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者のUUIDからエージェントIDを取得 ＊ゲストログイン用
func GetJobSeekerAgentIDByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerAgentIDByUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者のUUIDと名前から一致するか確認　*面談設定フォームの署名時に使用
func CheckJobSeekerByUUIDAndName(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CheckJobSeekerByUUIDAndNameParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.CheckJobSeekerByUUIDAndName(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 求職者のuuidとメールアドレスが一致するか確認 ※ゲストページログイン用
func SendJobSeekerResetPasswordEmail(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendJobSeekerResetPasswordEmailParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.SendJobSeekerResetPasswordEmail(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// 求職者のお問い合わせ
func SendJobSeekerContact(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendJobSeekerContactParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.SendJobSeekerContact(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// 求職者のお問い合わせ
func UpdateInterviewDateByJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateJobSeekerInterviewDateFromGuestPageParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.UpdateInterviewDateByJobSeekerID(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// 求職者のパスワード更新 ※ゲストページログイン用
func UpdateJobSeekerPassword(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateJobSeekerPasswordParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateJobSeekerPassword(param)
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
// LP用 API
//

// LPから求職者登録
func CreateJobSeekerFromLP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateJobSeekerFromLPParam
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.CreateJobSeekerFromLP(param)
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

// LPから求職者の電話番号を更新
func UpdateJobSeekerPhoneFromLP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateJobSeekerPhoneFromLPParam
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateJobSeekerPhoneFromLP(param)
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

// LPから求職者の希望条件を更新
func UpdateJobSeekerDesiredFromLP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.UpdateJobSeekerDesiredFromLPParam
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

		h := di.InitializeJobSeekerHandler(firebase, tx, sendgrid, oneSignal, slack)
		p, err := h.UpdateJobSeekerDesiredFromLP(param)
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

// LPの登録状況を取得する
func GetJobSeekerLPRegisterStatusByUUID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerUUIDStr = c.Param("job_seeker_uuid")
		)

		jobSeekerUUID, err := uuid.Parse(jobSeekerUUIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", "uuidのフォーマットが不正です", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.GetJobSeekerLPRegisterStatusByUUID(jobSeekerUUID)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// LPからのお問い合わせ
func SendLPContact(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendContactFromLPParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.SendLPContact(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// パスワード再設定ページのリンクをメールで送信
func SendJobSeekerResetPasswordEmailForLP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendJobSeekerResetPasswordEmailFromLPParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.SendJobSeekerResetPasswordEmailForLP(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

func ResetPasswordForLP(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.ResetPasswordFromLPParam
		)

		if err := bindAndValidate(c, &param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return err
		}

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.ResetPasswordForLP(param)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}

// LPからパスワードトークンの有効性を確認
func CheckResetPasswordToken(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal, slack config.Slack) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			resetPasswordToken = c.Param("password_token")
		)

		h := di.InitializeJobSeekerHandler(firebase, db, sendgrid, oneSignal, slack)
		p, err := h.CheckResetPasswordToken(resetPasswordToken)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}
		renderJSON(c, p)
		return nil
	}
}
