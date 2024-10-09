package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewUserSessionJSONPresenter(resp responses.UserSession) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewGuestEnterpriseUserSessionJSONPresenter(resp responses.GuestEnterpriseUserSession) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewGuestJobSeekerUserSessionJSONPresenter(resp responses.GuestJobSeekerUserSession) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewGoogleAuthSessionJSONPresenter(resp responses.GoogleAuthSession) Presenter {
	return NewJSONPresenter(200, resp)
}
