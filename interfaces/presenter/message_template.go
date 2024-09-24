package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewMessageTemplateJSONPresenter(resp responses.MessageTemplate) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewMessageTemplateListJSONPresenter(resp responses.MessageTemplateList) Presenter {
	return NewJSONPresenter(200, resp)
}
