package handler

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity/responses"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"github.com/spaceaiinc/autoscout-server/usecase/interactor"
)

type SessionHandler interface {
	// Gest API
	SignIn(token string) (presenter.Presenter, error)
	SignOut(token string) (presenter.Presenter, error)
	GetSignInUser(token string) (presenter.Presenter, error)

	// 求人企業が求人票修正するためのログイン
	SignInForGestEnterprise(password string, jobInformationUUID uuid.UUID) (presenter.Presenter, error)
	SignInForGestEnterpriseByTaskGroupUUID(password string, taskGroupUUID uuid.UUID) (presenter.Presenter, error)

	// ゲスト求職者ためのログイン
	SignInForGestJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error)
	SignInForGestJobSeekerFromLP(jobSeekerUUID, loginToken uuid.UUID) (presenter.Presenter, error)
	SignInForGestSendingJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error)

	// LPのログインフォームのログイン処理
	LoginGestJobSeekerForLP(email, password string) (presenter.Presenter, error)
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
/// Gest API
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

func (h *SessionHandlerImpl) SignInForGestEnterprise(password string, jobInformationUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGestEnterprise(interactor.SessionSignInForGestEnterpriseInput{
		Password:           password,
		JobInformationUUID: jobInformationUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGestEnterpriseUserSessionJSONPresenter(responses.NewGestEnterpriseUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignInForGestEnterpriseByTaskGroupUUID(password string, taskGroupUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGestEnterpriseByTaskGroupUUID(interactor.SessionSignInForGestEnterpriseByTaskGroupUUIDInput{
		Password:      password,
		TaskGroupUUID: taskGroupUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGestEnterpriseUserSessionJSONPresenter(responses.NewGestEnterpriseUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignInForGestJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGestJobSeeker(interactor.SessionSignInForGestJobSeekerInput{
		Password:      password,
		JobSeekerUUID: jobSeekerUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGestJobSeekerUserSessionJSONPresenter(responses.NewGestJobSeekerUserSession(output.User)), nil
}

// マイページログイン（LPからログイントークンを使ってログイン）
func (h *SessionHandlerImpl) SignInForGestJobSeekerFromLP(jobSeekerUUID, loginToken uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGestJobSeekerFromLP(interactor.SessionSignInForGestJobSeekerFromLPInput{
		JobSeekerUUID: jobSeekerUUID,
		LoginToken:    loginToken,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGestJobSeekerUserSessionJSONPresenter(responses.NewGestJobSeekerUserSession(output.User)), nil
}

func (h *SessionHandlerImpl) SignInForGestSendingJobSeeker(password string, jobSeekerUUID uuid.UUID) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.SignInForGestSendingJobSeeker(interactor.SessionSignInForGestSendingJobSeekerInput{
		Password:      password,
		JobSeekerUUID: jobSeekerUUID,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewGestJobSeekerUserSessionJSONPresenter(responses.NewGestJobSeekerUserSession(output.User)), nil
}

/****************************************************************************************/
// LP API

// LPのログインフォームのログイン処理
func (h *SessionHandlerImpl) LoginGestJobSeekerForLP(email, password string) (presenter.Presenter, error) {
	output, err := h.sessionInteractor.LoginGestJobSeekerForLP(interactor.LoginGestJobSeekerForLPInput{
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
