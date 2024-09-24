package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewEnterpriseProfileJSONPresenter(resp responses.EnterpriseProfile) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewEnterpriseProfileListJSONPresenter(resp responses.EnterpriseProfileList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewEnterpriseProfilePageListAndMaxPageAndIDListJSONPresenter(resp responses.EnterpriseProfileListAndMaxPageAndIDList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewEnterpriseReferenceMaterialJSONPresenter(resp responses.EnterpriseReferenceMaterial) Presenter {
	return NewJSONPresenter(200, resp)
}

// 企業と求人情報の一覧をJSONで返す　インポート用
func NewEnterpriseAndJobInformationListJSONPresenter(resp responses.EnterpriseAndJobInformationList) Presenter {
	return NewJSONPresenter(200, resp)
}
