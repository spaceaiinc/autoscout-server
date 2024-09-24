package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ChatThreadWithAgent struct {
	ChatThreadWithAgent *entity.ChatThreadWithAgent `json:"chat_thread_with_agent"`
}

func NewChatThreadWithAgent(chatThreadWithAgent *entity.ChatThreadWithAgent) ChatThreadWithAgent {
	return ChatThreadWithAgent{
		ChatThreadWithAgent: chatThreadWithAgent,
	}
}

type ChatThreadWithAgentList struct {
	ChatThreadWithAgentList []*entity.ChatThreadWithAgent `json:"chat_thread_with_agent_list"`
}

func NewChatThreadWithAgentList(chatThreadWithAgents []*entity.ChatThreadWithAgent) ChatThreadWithAgentList {
	return ChatThreadWithAgentList{
		ChatThreadWithAgentList: chatThreadWithAgents,
	}
}
