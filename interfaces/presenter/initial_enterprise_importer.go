package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewInitialEnterpriseImporterJSONPresenter(resp responses.InitialEnterpriseImporter) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewInitialEnterpriseImporterListJSONPresenter(resp responses.InitialEnterpriseImporterList) Presenter {
	return NewJSONPresenter(200, resp)
}
