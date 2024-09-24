package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type AdminHandler interface {
	Authorize(username, password string) (presenter.Presenter, error)
}

type AdminHandlerImpl struct {
	adminInteractor interactor.AdminInteractor
}

func NewAdminHandlerImpl(aI interactor.AdminInteractor) AdminHandler {
	return &AdminHandlerImpl{adminInteractor: aI}
}

func (h *AdminHandlerImpl) Authorize(username, password string) (presenter.Presenter, error) {
	var (
		input = interactor.AdminAuthorizeInput{Username: username, Password: password}
	)

	output, err := h.adminInteractor.Authorize(input)
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}
