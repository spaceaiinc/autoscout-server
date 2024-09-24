package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewSendingEnterpriseJSONPresenter(resp responses.SendingEnterprise) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingEnterpriseListJSONPresenter(resp responses.SendingEnterpriseList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingEnterprisePageListAndMaxPageAndIDListJSONPresenter(resp responses.SendingEnterpriseListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingEnterpriseReferenceMaterialJSONPresenter(resp responses.SendingEnterpriseReferenceMaterial) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSendingEnterpriseForSigninJSONPresenter(resp responses.SendingEnterpriseForSignin) Presenter {
	return NewJSONPresenter(200, resp)
}
