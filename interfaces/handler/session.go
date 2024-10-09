package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SessionHandler interface {
	// Guest API
	SignIn(token string) (presenter.Presenter, error)
	SignOut(token string) (presenter.Presenter, error)
	GetSignInUser(token string) (presenter.Presenter, error)

	// 求人企業が求人票修正するためのログイン
	SignInForGuestEnterprise(password string, jobInformationUUID uuid.UUID) (presenter.Presenter, error)
	SignInForGuestEnterpriseByTaskGroupUUID(password string, taskGroupUUID uuid.UUID) (presenter.Presenter, error)

	// ゲスト求職者ためのログイン
	SignInForGuestJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	SignInForGuestJobSeekerFromLP(jobSeekerUUID, loginToken uuid.UUID) (presenter.Presenter, error)
	SignInForGuestSendingJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error)

	// LPのログインフォームのログイン処理
	LoginGuestJobSeekerForLP(email, password string) (presenter.Presenter, error)
}

type SessionHandlerImpl struct {
	sessionInteractor interactor.SessionInteractor
}

func NewSessionHandlerImpl(si interactor.SessionInteractor) SessionHandler {
	return &SessionHandlerImpl{
		sessionInteractor: si,
	}
}

/****************************************************************************************/
/// Guest API
func (h *SessionHandlerImpl) SignIn(token string) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignIn(interactor.SessionSignInInput{Token: token})
	if err != nil {
		return nil, err
	}

	return presenter.NewUserSessionJSONPresenter(responses.NewUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignOut(token string) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignOut(interactor.SessionSignOutInput{Token: token})
	if err != nil {
		return nil, err
	}

	return presenter.NewOKJSONPresenter(responses.NewOK(output.OK)), nil
}

func (h *SessionHandlerImpl) GetSignInUser(token string) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.GetSignInUser(interactor.GetSignInUserInput{Token: token})
	if err != nil {
		return nil, err
	}

	return presenter.NewUserSessionJSONPresenter(responses.NewUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignInForGuestEnterprise(password string, jobInformationUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGuestEnterprise(interactor.SessionSignInForGuestEnterpriseInput{
		Password:           password,
		JobInformationUUID: jobInformationUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGuestEnterpriseUserSessionJSONPresenter(responses.NewGuestEnterpriseUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignInForGuestEnterpriseByTaskGroupUUID(password string, taskGroupUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGuestEnterpriseByTaskGroupUUID(interactor.SessionSignInForGuestEnterpriseByTaskGroupUUIDInput{
		Password:      password,
		TaskGroupUUID: taskGroupUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGuestEnterpriseUserSessionJSONPresenter(responses.NewGuestEnterpriseUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignInForGuestJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGuestJobSeeker(interactor.SessionSignInForGuestJobSeekerInput{
		Password:      password,
		JobSeekerUUID: jobSeekerUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGuestJobSeekerUserSessionJSONPresenter(responses.NewGuestJobSeekerUserSession(output.User)), nil
}

// マイページログイン（LPからログイントークンを使ってログイン）
func (h *SessionHandlerImpl) SignInForGuestJobSeekerFromLP(jobSeekerUUID, loginToken uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGuestJobSeekerFromLP(interactor.SessionSignInForGuestJobSeekerFromLPInput{
		JobSeekerUUID: jobSeekerUUID,
		LoginToken:    loginToken,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGuestJobSeekerUserSessionJSONPresenter(responses.NewGuestJobSeekerUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignInForGuestSendingJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGuestSendingJobSeeker(interactor.SessionSignInForGuestSendingJobSeekerInput{
		Password:      password,
		JobSeekerUUID: jobSeekerUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGuestJobSeekerUserSessionJSONPresenter(responses.NewGuestJobSeekerUserSession(output.User)), nil
}

/****************************************************************************************/
// LP API

// LPのログインフォームのログイン処理
func (h *SessionHandlerImpl) LoginGuestJobSeekerForLP(email, password string) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.LoginGuestJobSeekerForLP(interactor.LoginGuestJobSeekerForLPInput{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewJobSeekerLoginFromLPJSONPresenter(responses.NewJobSeekerLoginFromLP(output.UUID, output.LoginToken)), nil
}

/****************************************************************************************/
/// Agent API

/****************************************************************************************/
