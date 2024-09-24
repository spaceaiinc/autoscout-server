package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SendingBillingAddressHandler interface {
	// 汎用系 API
	UpdateSendingBillingAddress(param entity.UpdateSendingBillingAddressParam, billingAddressID uint) (presenter.Presenter, error)
	GetSendingBillingAddressByID(billingAddressID uint) (presenter.Presenter, error)
	GetSendingBillingAddressBySendingEnterpriseID(sendingEnterpriseID uint) (presenter.Presenter, error)
	GetSendingBillingAddressListBySendingEnterpriseIDList(sendingEnterpriseIDList []uint) (presenter.Presenter, error)
	GetAllSendingBillingAddress() (presenter.Presenter, error)

	// Admin API
}

type SendingBillingAddressHandlerImpl struct {
	billingAddressInteractor interactor.SendingBillingAddressInteractor
}

func NewSendingBillingAddressHandlerImpl(baI interactor.SendingBillingAddressInteractor) SendingBillingAddressHandler {
	return &SendingBillingAddressHandlerImpl{
		billingAddressInteractor: baI,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (h *SendingBillingAddressHandlerImpl) UpdateSendingBillingAddress(param entity.UpdateSendingBillingAddressParam, billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.UpdateSendingBillingAddress(interactor.UpdateSendingBillingAddressInput{
		UpdateParam:             param,
		SendingBillingAddressID: billingAddressID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingBillingAddressJSONPresenter(responses.NewSendingBillingAddress(output.SendingBillingAddress)), nil
}

// 請求先IDから請求先情報を取得する
func (h *SendingBillingAddressHandlerImpl) GetSendingBillingAddressByID(billingAddressID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetSendingBillingAddressByID(interactor.GetSendingBillingAddressByIDInput{
		SendingBillingAddressID: billingAddressID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingBillingAddressJSONPresenter(responses.NewSendingBillingAddress(output.SendingBillingAddress)), nil
}

// 求人企業IDから請求先一覧を取得する
func (h *SendingBillingAddressHandlerImpl) GetSendingBillingAddressBySendingEnterpriseID(enterpriseID uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetSendingBillingAddressBySendingEnterpriseID(interactor.GetSendingBillingAddressBySendingEnterpriseIDInput{
		SendingEnterpriseID: enterpriseID,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingBillingAddressJSONPresenter(responses.NewSendingBillingAddress(output.SendingBillingAddress)), nil
}

func (h *SendingBillingAddressHandlerImpl) GetSendingBillingAddressListBySendingEnterpriseIDList(enterpriseIDList []uint) (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetSendingBillingAddressListBySendingEnterpriseIDList(interactor.GetSendingBillingAddressListBySendingEnterpriseIDListInput{
		SendingEnterpriseIDList: enterpriseIDList,
	})

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingBillingAddressListJSONPresenter(responses.NewSendingBillingAddressList(output.SendingBillingAddressList)), nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API
// GetAllSendingBillingAddress() (presenter.Presenter, error)
func (h *SendingBillingAddressHandlerImpl) GetAllSendingBillingAddress() (presenter.Presenter, error) {
	output, err := h.billingAddressInteractor.GetAllSendingBillingAddress()

	if err != nil {
		return nil, err
	}

	return presenter.NewSendingBillingAddressListJSONPresenter(responses.NewSendingBillingAddressList(output.SendingBillingAddressList)), nil
}

/****************************************************************************************/
