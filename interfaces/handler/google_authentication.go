package handler

import (
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type GoogleAuthenticationHandler interface {
	GetGoogleAuthCodeURL(firebaseToken string) (presenter.Presenter, error)
	UpdateGoogleOauthToken(firebaseToken, code string) (presenter.Presenter, error)
}

type GoogleAuthenticationHandlerImpl struct {
	googleAuthenticationInteractor interactor.GoogleAuthenticationInteractor
}

func NewGoogleAuthenticationHandlerImpl(gai interactor.GoogleAuthenticationInteractor) GoogleAuthenticationHandler {
	return &GoogleAuthenticationHandlerImpl{
		googleAuthenticationInteractor: gai,
	}
}

/****************************************************************************************/
//  API
//
func (h *GoogleAuthenticationHandlerImpl) GetGoogleAuthCodeURL(token string) (presenter.Presenter, error) {
	output, err := h.googleAuthenticationInteractor.GetGoogleAuthCodeURL(interactor.GetGoogleAuthCodeURLInput{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGoogleAuthSessionJSONPresenter(responses.NewGoogleAuthSession(output.AuthURL)), nil
}

func (h *GoogleAuthenticationHandlerImpl) UpdateGoogleOauthToken(token, code string) (presenter.Presenter, error) {
	output, err := h.googleAuthenticationInteractor.UpdateGoogleOauthToken(interactor.UpdateGoogleOauthTokenInput{
		Token: token,
		Code:  code,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

/****************************************************************************************/
