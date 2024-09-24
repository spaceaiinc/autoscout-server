package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewChatMessageWithAgentJSONPresenter(resp responses.ChatMessageWithAgent) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewChatMessageWithAgentListJSONPresenter(resp responses.ChatMessageWithAgentList) Presenter {
	return NewJSONPresenter(200, resp)
}
