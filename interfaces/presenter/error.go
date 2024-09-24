package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ErrorJSONPresenter struct {
	PresenterImpl
}

func NewErrorJSONPresenter(err error) Presenter {
	code, message := entity.ErrorInfo(err)
	return &ErrorJSONPresenter{
		PresenterImpl: PresenterImpl{
			statusCode: code,
			data: map[string]interface{}{
				"error":  message,
				"status": code,
			},
		},
	}
}
