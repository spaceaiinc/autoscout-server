package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewUserSessionJSONPresenter(resp responses.UserSession) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewGestEnterpriseUserSessionJSONPresenter(resp responses.GestEnterpriseUserSession) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewGestJobSeekerUserSessionJSONPresenter(resp responses.GestJobSeekerUserSession) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewGoogleAuthSessionJSONPresenter(resp responses.GoogleAuthSession) Presenter {
	return NewJSONPresenter(200, resp)
}
