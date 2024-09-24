package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewScoutServiceJSONPresenter(resp responses.ScoutService) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewScoutServiceListJSONPresenter(resp responses.ScoutServiceList) Presenter {
	return NewJSONPresenter(200, resp)
}
