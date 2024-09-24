package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type ScheduleHandler interface {
	// 汎用系 API
	GetScheduleListWithInPeriod(staffID uint, startDate, endDate string) (presenter.Presenter, error)
	GetScheduleListWithInPeriodByStaffIDList(staffIDList []uint, startDate, endDate string) (presenter.Presenter, error)
}

type ScheduleHandlerImpl struct {
	scheduleInteractor interactor.ScheduleInteractor
}

func NewScheduleHandlerImpl(sI interactor.ScheduleInteractor) ScheduleHandler {
	return &ScheduleHandlerImpl{
		scheduleInteractor: sI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *ScheduleHandlerImpl) GetScheduleListWithInPeriod(staffID uint, startDate, endDate string) (presenter.Presenter, error) {
	output, err := h.scheduleInteractor.GetScheduleListWithInPeriod(interactor.GetScheduleListWithInPeriodInput{
		AgentStaffID: staffID,
		StartDate:    startDate,
		EndDate:      endDate,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewScheduleListJSONPresenter(responses.NewScheduleList(output.ScheduleList)), nil
}

func (h *ScheduleHandlerImpl) GetScheduleListWithInPeriodByStaffIDList(staffIDList []uint, startDate, endDate string) (presenter.Presenter, error) {
	output, err := h.scheduleInteractor.GetScheduleListWithInPeriodByStaffIDList(interactor.GetScheduleListWithInPeriodByStaffIDListInput{
		StaffIDList: staffIDList,
		StartDate:   startDate,
		EndDate:     endDate,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewScheduleListJSONPresenter(responses.NewScheduleList(output.ScheduleList)), nil
}
