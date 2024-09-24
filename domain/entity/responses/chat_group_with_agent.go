package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ChatGroupWithAgent struct {
	ChatGroupWithAgent *entity.ChatGroupWithAgent `json:"chat_group_with_agent"`
}

func NewChatGroupWithAgent(chatGroupWithAgent *entity.ChatGroupWithAgent) ChatGroupWithAgent {
	return ChatGroupWithAgent{
		ChatGroupWithAgent: chatGroupWithAgent,
	}
}

type ChatGroupWithAgentList struct {
	ChatGroupWithAgentList []*entity.ChatGroupWithAgent `json:"chat_group_with_agent_list"`
}

func NewChatGroupWithAgentList(chatGroupWithAgents []*entity.ChatGroupWithAgent) ChatGroupWithAgentList {
	return ChatGroupWithAgentList{
		ChatGroupWithAgentList: chatGroupWithAgents,
	}
}
