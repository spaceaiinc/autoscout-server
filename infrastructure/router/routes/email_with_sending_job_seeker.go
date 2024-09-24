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
//
// 送客求職者へのメール送信
func SendEmailWithSendingJobSeeker(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SendEmailWithSendingJobSeekerParam
		)

		// フォームパラムの取得
		formParams, err := c.FormParams()
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		groupIDStr := formParams.Get("group_id")
		agentStaffIDStr := formParams.Get("agent_staff_id")
		jobSeekerIDStr := formParams.Get("sending_job_seeker_id")
		subject := formParams.Get("subject")
		content := formParams.Get("content")
		toEmail := formParams.Get("to_email")

		if toEmail == "" {
			wrapped := fmt.Errorf("送信先のメールアドレスが空です:%w", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		toName := formParams.Get("to_name")
		if toName == "" {
			wrapped := fmt.Errorf("送信先の名前が空です:%w", entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// groupID 型変換
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// agentStaffID 型変換
		agentStaffID, err := strconv.Atoi(agentStaffIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// jobSeekerID 型変換
		jobSeekerID, err := strconv.Atoi(jobSeekerIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		// ファイル
		multipartForm, err := c.MultipartForm()
		if err != nil {
			wrapped := fmt.Errorf("ファイルの受け取りエラー: %s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}
		fmt.Println("multipart form", multipartForm)

		files := multipartForm.File["files"]
		fmt.Println("multimart form file", files)

		// 複数ある場合はカンマ区切り
		fileName := ""
		for _, file := range files {
			fileName += file.Filename + "|--|--T--|--|"
		}

		// 送信先
		to := entity.EmailUser{
			Name:  toName,
			Email: toEmail,
		}

		param = entity.SendEmailWithSendingJobSeekerParam{
			GroupID:            uint(groupID),
			AgentStaffID:       uint(agentStaffID),
			SendingJobSeekerID: uint(jobSeekerID),
			Subject:            subject,
			Content:            content,
			To:                 to,
			Files:              files,
			FileName:           fileName,
		}
		fmt.Println(param)

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

		h := di.InitializeEmailWithSendingJobSeekerHandler(firebase, tx, sendgrid, oneSignal)
		p, err := h.SendEmailWithSendingJobSeeker(param)
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

// エージェントIDからエージェントと求職者のチャットメッセージメッセージの一覧取得
func GetEmailWithSendingJobSeekerListBySendingJobSeekerID(db *database.DB, firebase usecase.Firebase, sendgrid config.Sendgrid, oneSignal config.OneSignal) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			jobSeekerIDStr = c.Param("sending_job_seeker_id")
		)

		jobSeekerIDInt, err := strconv.Atoi(jobSeekerIDStr)

		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		h := di.InitializeEmailWithSendingJobSeekerHandler(firebase, db, sendgrid, oneSignal)
		p, err := h.GetEmailWithSendingJobSeekerListBySendingJobSeekerID(uint(jobSeekerIDInt))
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

/****************************************************************************************/
