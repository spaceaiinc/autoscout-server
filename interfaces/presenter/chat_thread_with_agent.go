package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewChatThreadWithAgentJSONPresenter(resp responses.ChatThreadWithAgent) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatThreadWithAgentListJSONPresenter(resp responses.ChatThreadWithAgentList) Presenter {
	return NewJSONPresenter(200, resp)
}
