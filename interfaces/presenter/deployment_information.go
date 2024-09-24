package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewDeploymentInformationJSONPresenter(resp responses.DeploymentInformation) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewDeploymentInformationListJSONPresenter(resp responses.DeploymentInformationList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewDeploymentInformationListAndMaxPageJSONPresenter(resp responses.DeploymentInformationListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}
