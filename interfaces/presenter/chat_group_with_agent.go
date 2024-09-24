package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewChatGroupWithAgentJSONPresenter(resp responses.ChatGroupWithAgent) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatGroupWithAgentListJSONPresenter(resp responses.ChatGroupWithAgentList) Presenter {
	return NewJSONPresenter(200, resp)
}
