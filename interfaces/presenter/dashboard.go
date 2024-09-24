package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewDashboardJSONPresenter(resp responses.Dashboard) Presenter {
	return NewJSONPresenter(200, resp)
}
