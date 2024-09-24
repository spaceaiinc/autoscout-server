package handler

import (
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingPhaseHandler interface {
	// 汎用系 API
	UpdateSendingPhaseSendingDate(id uint, sendingDate time.Time) (presenter.Presenter, error)

	GetSendingPhaseByID(id uint) (presenter.Presenter, error)
	GetSendingPhaseListBySendingJobSeekerID(sendingJobSeekerID uint) (presenter.Presenter, error) // 指定求職者IDの送客進捗情報を取得する関数

	// ページネーション
	GetSearchSendingPhaseListByPageAndTabAndAgentID(agentID, pageNumber, tabNumber uint, freeWord string, staffIDList, senderIDList, sendAgentIDList []uint) (presenter.Presenter, error)

	// 送客応諾後の終了理由
	CreateSendingPhaseEndStatus(param entity.CreateSendingPhaseEndStatusParam) (presenter.Presenter, error)
}

type SendingPhaseHandlerImpl struct {
	sendingPhaseInteractor interactor.SendingPhaseInteractor
}

func NewSendingPhaseHandlerImpl(epI interactor.SendingPhaseInteractor) SendingPhaseHandler {
	return &SendingPhaseHandlerImpl{
		sendingPhaseInteractor: epI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//

func (h *SendingPhaseHandlerImpl) UpdateSendingPhaseSendingDate(id uint, sendingDate time.Time) (presenter.Presenter, error) {
	output, err := h.sendingPhaseInteractor.UpdateSendingPhaseSendingDate(interactor.UpdateSendingPhaseSendingDateInput{
		SendingPhaseID: id,
		SendingDate:    sendingDate,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 送客進捗IDから送客進捗情報を取得する
func (h *SendingPhaseHandlerImpl) GetSendingPhaseByID(id uint) (presenter.Presenter, error) {
	output, err := h.sendingPhaseInteractor.GetSendingPhaseByID(interactor.GetSendingPhaseByIDInput{
		SendingPhaseID: id,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingPhaseJSONPresenter(responses.NewSendingPhase(output.SendingPhase)), nil
}

// 指定求職者IDの送客進捗情報を取得する関数
func (h *SendingPhaseHandlerImpl) GetSendingPhaseListBySendingJobSeekerID(sendingJobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.sendingPhaseInteractor.GetSendingPhaseListBySendingJobSeekerID(interactor.GetSendingPhaseListBySendingJobSeekerIDInput{
		SendingJobSeekerID: sendingJobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingPhaseListJSONPresenter(responses.NewSendingPhaseList(output.SendingPhaseList)), nil
}

/****************************************************************************************/
/// ページネーション
//
func (h *SendingPhaseHandlerImpl) GetSearchSendingPhaseListByPageAndTabAndAgentID(agentID, pageNumber, tabNumber uint, freeWord string, staffIDList, senderIDList, sendAgentIDList []uint) (presenter.Presenter, error) {
	output, err := h.sendingPhaseInteractor.GetSearchSendingPhaseListByPageAndTabAndAgentID(interactor.GetSearchSendingPhaseListByPageAndTabAndAgentIDInput{
		AgentID:         agentID,
		PageNumber:      pageNumber,
		TabNumber:       tabNumber,
		FreeWord:        freeWord,
		StaffIDList:     staffIDList,
		SenderIDList:    senderIDList,
		SendAgentIDList: sendAgentIDList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingJobSeekerTableListAndMaxPageAndIDListAndListCountJSONPresenter(responses.NewSendingJobSeekerTableListAndMaxPageAndIDListAndListCount(
		output.SendingJobSeekerTableList,
		output.MaxPageNumber,
		output.CompleteSendingCount,
		output.AcceptSendingCount,
		output.PreparingAfterInterviewCount,
		output.WaitingForInterviewCount,
		output.UnregisterDetailCount,
		output.UnregisterScheduleCount,
		output.CloseSendingCount,
	)), nil
}

/****************************************************************************************/
/// 送客応諾後の終了理由
//
func (h *SendingPhaseHandlerImpl) CreateSendingPhaseEndStatus(param entity.CreateSendingPhaseEndStatusParam) (presenter.Presenter, error) {
	output, err := h.sendingPhaseInteractor.CreateSendingPhaseEndStatus(interactor.CreateSendingPhaseEndStatusInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
