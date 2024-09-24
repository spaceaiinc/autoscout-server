package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type DashboardHandler interface {
	// デフォルト表示
	GetDashboardForDefaultPreview(staffID uint) (presenter.Presenter, error)
	GetSaleListForDefaultPreview(staffID, pageNumber uint) (presenter.Presenter, error)
	GetReleaseJobSeekerListForDefaultPreview(staffID, pageNumber uint) (presenter.Presenter, error)

	// 検索
	GetDashboardForSearch(parm entity.SearchDashboard) (presenter.Presenter, error)
	GetSaleListForSearch(pageNumber uint, parm entity.SearchDashboard) (presenter.Presenter, error)
	GetReleaseJobSeekerListForSearch(pageNumber uint, parm entity.SearchDashboard) (presenter.Presenter, error)

	// csv API
	ExportAccuracyCSV(agentID uint) (string, error)
}

type DashboardHandlerImpl struct {
	dashboardInteractor interactor.DashboardInteractor
}

func NewDashboardHandlerImpl(dI interactor.DashboardInteractor) DashboardHandler {
	return &DashboardHandlerImpl{
		dashboardInteractor: dI,
	}
}

/****************************************************************************************/
// デフォルト表示
//
func (h *DashboardHandlerImpl) GetDashboardForDefaultPreview(staffID uint) (presenter.Presenter, error) {
	output, err := h.dashboardInteractor.GetDashboardForDefaultPreview(interactor.GetDashboardForDefaultPreviewInput{
		AgentStaffID: staffID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewDashboardJSONPresenter(responses.NewDashboard(output.Dashboard)), nil
}

func (h *DashboardHandlerImpl) GetSaleListForDefaultPreview(staffID, pageNumber uint) (presenter.Presenter, error) {
	output, err := h.dashboardInteractor.GetSaleListForDefaultPreview(interactor.GetSaleListForDefaultPreviewInput{
		AgentStaffID: staffID,
		PageNumber:   pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSaleListAndMaxPageAndIDListJSONPresenter(responses.NewSaleListAndMaxPageAndIDList(output.SaleLlist, output.MaxPageNumber, output.IDList)), nil
}

func (h *DashboardHandlerImpl) GetReleaseJobSeekerListForDefaultPreview(staffID, pageNumber uint) (presenter.Presenter, error) {
	output, err := h.dashboardInteractor.GetReleaseJobSeekerListForDefaultPreview(interactor.GetReleaseJobSeekerListForDefaultPreviewInput{
		AgentStaffID: staffID,
		PageNumber:   pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDList(output.JobSeekerList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
// 検索
//
func (h *DashboardHandlerImpl) GetDashboardForSearch(param entity.SearchDashboard) (presenter.Presenter, error) {
	output, err := h.dashboardInteractor.GetDashboardForSearch(interactor.GetDashboardForSearchInput{
		SearchParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewDashboardJSONPresenter(responses.NewDashboard(output.Dashboard)), nil
}

func (h *DashboardHandlerImpl) GetSaleListForSearch(pageNumber uint, param entity.SearchDashboard) (presenter.Presenter, error) {
	output, err := h.dashboardInteractor.GetSaleListForSearch(interactor.GetSaleListForSearchInput{
		PageNumber:  pageNumber,
		SearchParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSaleListAndMaxPageAndIDListJSONPresenter(responses.NewSaleListAndMaxPageAndIDList(output.SaleLlist, output.MaxPageNumber, output.IDList)), nil
}

func (h *DashboardHandlerImpl) GetReleaseJobSeekerListForSearch(pageNumber uint, param entity.SearchDashboard) (presenter.Presenter, error) {
	output, err := h.dashboardInteractor.GetReleaseJobSeekerListForSearch(interactor.GetReleaseJobSeekerListForSearchInput{
		PageNumber:  pageNumber,
		SearchParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerListAndMaxPageAndIDListJSONPresenter(responses.NewJobSeekerListAndMaxPageAndIDList(output.JobSeekerList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/
/// CSV API
//
//csvファイルを出力する
func (h *DashboardHandlerImpl) ExportAccuracyCSV(agentID uint) (string, error) {
	output, err := h.dashboardInteractor.ExportAccuracyCSV(interactor.ExportAccuracyCSVInput{
		AgentID: agentID,
	})

	if err != nil {
		return "", err
	}

	return output.FilePath.PathName, nil
}
