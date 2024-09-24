package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ChatMessageWithAgent struct {
	ChatMessageWithAgent *entity.ChatMessageWithAgent `json:"chat_message_with_agent"`
}

func NewChatMessageWithAgent(chatMessageWithAgent *entity.ChatMessageWithAgent) ChatMessageWithAgent {
	return ChatMessageWithAgent{
		ChatMessageWithAgent: chatMessageWithAgent,
	}
}

type ChatMessageWithAgentList struct {
	ChatMessageWithAgentList []*entity.ChatMessageWithAgent `json:"chat_message_with_agent_list"`
}

func NewChatMessageWithAgentList(chatMessageWithAgents []*entity.ChatMessageWithAgent) ChatMessageWithAgentList {
	return ChatMessageWithAgentList{
		ChatMessageWithAgentList: chatMessageWithAgents,
	}
}
