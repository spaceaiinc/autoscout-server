package handler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ScoutServiceHandler interface {
	// 汎用系 API
	CreateScoutService(param entity.CreateOrUpdateScoutServiceParam) (presenter.Presenter, error)
	UpdateScoutService(param entity.CreateOrUpdateScoutServiceParam, scoutServiceID uint) (presenter.Presenter, error)
	UpdateScoutServicePassword(param entity.UpdateScoutServicePasswordParam) (presenter.Presenter, error)
	DeleteScoutService(id uint) (presenter.Presenter, error)
	GetByID(id uint) (presenter.Presenter, error)
	GetListByAgentID() func(c echo.Context) error

	// Batch処理 API
	BatchScout(now time.Time, agentRobotID uint) (presenter.Presenter, error)
	BatchEntry(now time.Time) (presenter.Presenter, error)

	// Gmail API
	GmailWebHook(pubSubStruct *entity.PubsubStruct) (presenter.Presenter, error)
}

type ScoutServiceHandlerImpl struct {
	scoutServiceInteractor interactor.ScoutServiceInteractor
}

func NewScoutServiceHandlerImpl(ssI interactor.ScoutServiceInteractor) ScoutServiceHandler {
	return &ScoutServiceHandlerImpl{
		scoutServiceInteractor: ssI,
	}
}

/****************************************************************************************/
// 汎用系 API
//
// スカウトサービスを作成
func (h *ScoutServiceHandlerImpl) CreateScoutService(param entity.CreateOrUpdateScoutServiceParam) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.CreateScoutService(interactor.CreateScoutServiceInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewScoutServiceJSONPresenter(responses.NewScoutService(output.ScoutService)), nil
}

// スカウトサービスを更新
func (h *ScoutServiceHandlerImpl) UpdateScoutService(param entity.CreateOrUpdateScoutServiceParam, scoutServiceID uint) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.UpdateScoutService(interactor.UpdateScoutServiceInput{
		UpdateParam:    param,
		ScoutServiceID: scoutServiceID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewScoutServiceJSONPresenter(responses.NewScoutService(output.ScoutService)), nil
}

// スカウトサービスのパスワード更新
func (h *ScoutServiceHandlerImpl) UpdateScoutServicePassword(param entity.UpdateScoutServicePasswordParam) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.UpdateScoutServicePassword(interactor.UpdateScoutServicePasswordInput{
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// スカウトサービスを削除
func (h *ScoutServiceHandlerImpl) DeleteScoutService(scoutServiceID uint) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.DeleteScoutService(interactor.DeleteScoutServiceInput{
		ScoutServiceID: scoutServiceID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// スカウトサービスを取得
func (h *ScoutServiceHandlerImpl) GetByID(scoutServiceID uint) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.GetByID(interactor.ScoutServiceGetByIDInput{
		ScoutServiceID: scoutServiceID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewScoutServiceJSONPresenter(responses.NewScoutService(output.ScoutService)), nil
}

// エージェントIDからスカウトサービスを取得
func (h *ScoutServiceHandlerImpl) GetListByAgentID() func(c echo.Context) error {
	return func(c echo.Context) error {
		scoutServiceIDStr := c.Param("scout_service_id")

		scoutServiceIDInt, err := strconv.Atoi(scoutServiceIDStr)
		if err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
			return wrapped
		}

		output, err := h.scoutServiceInteractor.GetListByAgentID(interactor.GetListByAgentIDInput{
			AgentID: uint(scoutServiceIDInt),
		})
		if err != nil {
			renderJSON(c, presenter.NewErrorJSONPresenter(err))
			return err
		}

		renderJSON(c, presenter.NewScoutServiceListJSONPresenter(responses.NewScoutServiceList(output.ScoutServiceList)))
		return nil
	}
}

/****************************************************************************************/
// Batch処理 API
//

// 各スカウトサービスからスカウト送信する
func (h *ScoutServiceHandlerImpl) BatchScout(now time.Time, agentRobotID uint) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.BatchScout(interactor.BatchScoutInput{
		Now:          now,
		AgentRobotID: agentRobotID,
	})

	if err != nil {
		return nil, err
	}

	if !output.OK {
		return nil, errors.New("エントリー取得またはスカウト送信でpanicエラーが発生しました")
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *ScoutServiceHandlerImpl) BatchEntry(now time.Time) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.BatchEntry(interactor.BatchEntryInput{
		Now: now,
	})
	if err != nil {
		return nil, err
	}

	if !output.OK {
		return nil, errors.New("エントリー取得またはスカウト送信でpanicエラーが発生しました")
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/

/****************************************************************************************/
// Gmail API

func (h *ScoutServiceHandlerImpl) GmailWebHook(pubSubStruct *entity.PubsubStruct) (presenter.Presenter, error) {
	output, err := h.scoutServiceInteractor.GmailWebHook(interactor.GmailWebHookInput{
		PubSubStruct: pubSubStruct,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
