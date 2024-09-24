package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type BillingAddressHandler interface {
	// 汎用系 API
	CreateBillingAddress(param entity.CreateBillingAddressParam, enterpriseID uint) (presenter.Presenter, error)
	UpdateBillingAddress(param entity.UpdateBillingAddressParam, billingAddressID uint) (presenter.Presenter, error)
	UpdateBillingAddressStaffIDByIDListAtOnce(agentStaffID uint, billingAddressIDList []uint) (presenter.Presenter, error)
	DeleteBillingAddress(billingAddressID uint) (presenter.Presenter, error)
	GetBillingAddressByID(billingAddressID uint) (presenter.Presenter, error)
	GetBillingAddressListByEnterpriseID(enterpriseID uint) (presenter.Presenter, error)
	GetBillingAddressListByPageAndAgentID(agentID, pageNumber uint) (presenter.Presenter, error)
	GetSearchBillingAddressListByPageAndAgentID(agentID, pageNumber uint, searchParam entity.SearchBillingAddress) (presenter.Presenter, error)
	GetAllBillingAddress() (presenter.Presenter, error)

	// Admin API
}

type BillingAddressHandlerImpl struct {
	billingAddressInteractor interactor.BillingAddressInteractor
}

func NewBillingAddressHandlerImpl(baI interactor.BillingAddressInteractor) BillingAddressHandler {
	return &BillingAddressHandlerImpl{
		billingAddressInteractor: baI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *BillingAddressHandlerImpl) CreateBillingAddress(param entity.CreateBillingAddressParam, enterpriseID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.CreateBillingAddress(interactor.CreateBillingAddressInput{
		CreateParam:  param,
		EnterpriseID: enterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewBillingAddressJSONPresenter(responses.NewBillingAddress(output.BillingAddress)), nil
}

func (h *BillingAddressHandlerImpl) UpdateBillingAddress(param entity.UpdateBillingAddressParam, billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.UpdateBillingAddress(interactor.UpdateBillingAddressInput{
		UpdateParam:      param,
		BillingAddressID: billingAddressID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewBillingAddressJSONPresenter(responses.NewBillingAddress(output.BillingAddress)), nil
}

// 複数の請求先のagent_staff_idカラムを一括更新
func (h *BillingAddressHandlerImpl) UpdateBillingAddressStaffIDByIDListAtOnce(agentStaffID uint, billingAddressIDList []uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.UpdateBillingAddressStaffIDByIDListAtOnce(interactor.UpdateBillingAddressAgentStaffInput{
		AgentStaffID:         agentStaffID,
		BillingAddressIDList: billingAddressIDList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *BillingAddressHandlerImpl) DeleteBillingAddress(billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.DeleteBillingAddress(interactor.DeleteBillingAddressInput{
		BillingAddressID: billingAddressID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

// 請求先IDから請求先情報を取得する
func (h *BillingAddressHandlerImpl) GetBillingAddressByID(billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetBillingAddressByID(interactor.GetBillingAddressByIDInput{
		BillingAddressID: billingAddressID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewBillingAddressJSONPresenter(responses.NewBillingAddress(output.BillingAddress)), nil
}

// 求人企業IDから請求先一覧を取得する
func (h *BillingAddressHandlerImpl) GetBillingAddressListByEnterpriseID(enterpriseID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetBillingAddressListByEnterpriseID(interactor.GetBillingAddressListByEnterpriseIDInput{
		EnterpriseID: enterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewBillingAddressListJSONPresenter(responses.NewBillingAddressList(output.BillingAddressList)), nil
}

// AgentIDから請求先一覧を取得する
func (h *BillingAddressHandlerImpl) GetBillingAddressListByPageAndAgentID(agentID, pageNumber uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetBillingAddressListByPageAndAgentID(interactor.GetBillingAddressListByPageAndAgentIDInput{
		AgentID:    agentID,
		PageNumber: pageNumber,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewBillingAddressPageListAndMaxPageAndIDListJSONPresenter(responses.NewBillingAddressListAndMaxPageAndIDList(output.BillingAddressList, output.MaxPageNumber, output.IDList)), nil
}

// エージェントIDとクエリパラムで請求一覧を絞り込み
func (h *BillingAddressHandlerImpl) GetSearchBillingAddressListByPageAndAgentID(agentID, pageNumber uint, searchParam entity.SearchBillingAddress) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetSearchBillingAddressListByPageAndAgentID(interactor.GetSearchBillingAddressListByPageAndAgentIDInput{
		AgentID:     agentID,
		PageNumber:  pageNumber,
		SearchParam: searchParam,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewBillingAddressPageListAndMaxPageAndIDListJSONPresenter(responses.NewBillingAddressListAndMaxPageAndIDList(output.BillingAddressList, output.MaxPageNumber, output.IDList)), nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API
// GetAllBillingAddress() (presenter.Presenter, error)
func (h *BillingAddressHandlerImpl) GetAllBillingAddress() (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetAllBillingAddress()

	if err != nil {
		return nil, err
	}

	return presenter.NewBillingAddressListJSONPresenter(responses.NewBillingAddressList(output.BillingAddressList)), nil
}

/****************************************************************************************/
