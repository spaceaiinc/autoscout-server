package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SaleHandler interface {
	// 汎用系 API
	CreateSale(param entity.CreateOrUpdateSaleParam) (presenter.Presenter, error)
	UpdateSale(saleID uint, param entity.CreateOrUpdateSaleParam) (presenter.Presenter, error)
	GetSaleByID(saleID uint) (presenter.Presenter, error)
	GetSaleByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error)
	GetSaleListByIDList(idList []uint) (presenter.Presenter, error)
	GetAccuracySearchList(searchParam entity.SearchAccuracy) (presenter.Presenter, error)
}

type SaleHandlerImpl struct {
	saleInteractor interactor.SaleInteractor
}

func NewSaleHandlerImpl(itI interactor.SaleInteractor) SaleHandler {
	return &SaleHandlerImpl{
		saleInteractor: itI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *SaleHandlerImpl) CreateSale(param entity.CreateOrUpdateSaleParam) (presenter.Presenter, error) {
	output, err := h.saleInteractor.CreateSale(interactor.CreateSaleInput{
		CreateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SaleHandlerImpl) UpdateSale(saleID uint, param entity.CreateOrUpdateSaleParam) (presenter.Presenter, error) {
	output, err := h.saleInteractor.UpdateSale(interactor.UpdateSaleInput{
		SaleID:      saleID,
		UpdateParam: param,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SaleHandlerImpl) GetSaleByID(saleID uint) (presenter.Presenter, error) {
	output, err := h.saleInteractor.GetSaleByID(interactor.GetSaleByIDInput{
		SaleID: saleID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSaleJSONPresenter(responses.NewSale(output.Sale)), nil
}

func (h *SaleHandlerImpl) GetSaleByJobSeekerID(jobSeekerID uint) (presenter.Presenter, error) {
	output, err := h.saleInteractor.GetSaleByJobSeekerID(interactor.GetSaleByJobSeekerIDInput{
		JobSeekerID: jobSeekerID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSaleJSONPresenter(responses.NewSale(output.Sale)), nil
}

func (h *SaleHandlerImpl) GetSaleListByIDList(idList []uint) (presenter.Presenter, error) {
	output, err := h.saleInteractor.GetSaleListByIDList(interactor.GetSaleListByIDListInput{
		IDList: idList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSaleListJSONPresenter(responses.NewSaleList(output.SaleList)), nil
}

// 自社が絡むヨミ情報を全て取得
func (h *SaleHandlerImpl) GetAccuracySearchList(searchParam entity.SearchAccuracy) (presenter.Presenter, error) {
	output, err := h.saleInteractor.GetAccuracySearchList(interactor.GetAccuracySearchListInput{
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSaleListAndMaxPageAndIDListJSONPresenter(responses.NewSaleListAndMaxPageAndIDList(output.SaleList, output.MaxPageNumber, output.IDList)), nil
}
