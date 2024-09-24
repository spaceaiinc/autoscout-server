package routes

import (
	"fmt"
	"strconv"
	"time"

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
//
// 送客進捗の送客予定日時の更新
type updateSendingPhaseSendingDateParam struct {
	SendingPhaseID uint      `json:"sending_phase_id" validate:"required"`
	SendingDate    time.Time `json:"sending_date" validate:"required"`
}

func UpdateSendingPhaseSendingDate(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param updateSendingPhaseSendingDateParam
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

		h := di.InitializeSendingPhaseHandler(firebase, tx, sendgrid)
		p, err := h.UpdateSendingPhaseSendingDate(param.SendingPhaseID, param.SendingDate)
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

// 送客進捗IDから送客進捗情報を取得
func GetSendingPhaseByID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			sendingPhaseIDStr = c.Param("sending_phase_id")
		)

		sendingPhaseIDInt, err := strconv.Atoi(sendingPhaseIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeSendingPhaseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingPhaseByID(uint(sendingPhaseIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

// 指定求職者IDの送客進捗情報を取得する関数
func GetSendingPhaseListBySendingJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
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

		h := di.InitializeSendingPhaseHandler(firebase, db, sendgrid)
		p, err := h.GetSendingPhaseListBySendingJobSeekerID(uint(sendingJobSeekerIDInt))
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
/// ページネーション
//

// 送客の進捗一覧ページに表示するリストを取得するapi（query: page_number, tab_number, staff_id_list, sender_id_list, send_agent_id_list）
func GetSearchSendingPhaseListByPageAndTabAndAgentID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			agentIDStr         = c.Param("agent_id")
			pageNumberStr      = c.QueryParam("page_number")
			tabNumberStr       = c.QueryParam("tab_number")
			freeWordStr        = c.QueryParam("free_word")
			staffIDListStr     = c.QueryParams()["staff_id_list[]"]
			senderIDListStr    = c.QueryParams()["sender_id_list[]"]
			sendAgentIDListStr = c.QueryParams()["send_agent_id_list[]"]

			staffIDList     []uint
			senderIDList    []uint
			sendAgentIDList []uint
		)

		// エージェントのID
		agentIDInt, err := strconv.Atoi(agentIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// ページ番号
		pageNumberInt, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// タブの種類
		tabNumberInt, err := strconv.Atoi(tabNumberStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// 担当者(絞り込み)
		for _, staffStrr := range staffIDListStr {
			staffInt, err := strconv.Atoi(staffStrr)
			if err != nil {
				wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
				renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
				return wrapped
			}
			staffIDList = append(staffIDList, uint(staffInt))
		}

		// 送客元(絞り込み)
		for _, senderStrr := range senderIDListStr {
			senderInt, err := strconv.Atoi(senderStrr)
			if err != nil {
				wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
				renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
				return wrapped
			}
			senderIDList = append(senderIDList, uint(senderInt))
		}

		// 送客先(絞り込み)
		for _, sendAgentStrr := range sendAgentIDListStr {
			sendAgentInt, err := strconv.Atoi(sendAgentStrr)
			if err != nil {
				wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
				renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
				return wrapped
			}
			sendAgentIDList = append(sendAgentIDList, uint(sendAgentInt))
		}

		h := di.InitializeSendingPhaseHandler(firebase, db, sendgrid)
		p, err := h.GetSearchSendingPhaseListByPageAndTabAndAgentID(
			uint(agentIDInt),
			uint(pageNumberInt),
			uint(tabNumberInt),
			freeWordStr,
			staffIDList,
			senderIDList,
			sendAgentIDList,
		)
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, p)
		return nil
	}
}

/****************************************************************************************/
/// 送客応諾後の終了理由
//
func CreateSendingPhaseEndStatus(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.CreateSendingPhaseEndStatusParam
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

		h := di.InitializeSendingPhaseHandler(firebase, tx, sendgrid)
		p, err := h.CreateSendingPhaseEndStatus(param)
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
