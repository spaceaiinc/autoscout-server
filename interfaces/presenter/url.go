package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewURLJSONPresenter(resp responses.URL) Presenter {
	return NewJSONPresenter(200, resp)
}
