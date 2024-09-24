package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type Agent struct {
	Agent *entity.Agent `json:"agent"`
}

func NewAgent(agentAgent *entity.Agent) Agent {
	return Agent{
		Agent: agentAgent,
	}
}

type AgentList struct {
	AgentList []*entity.Agent `json:"agent_list"`
}

func NewAgentList(agentAgents []*entity.Agent) AgentList {
	return AgentList{
		AgentList: agentAgents,
	}
}

type AgentLine struct {
	AgentLine *entity.AgentLineChannelParam `json:"agent_line"`
}

func NewAgentLine(agentLine *entity.AgentLineChannelParam) AgentLine {
	return AgentLine{
		AgentLine: agentLine,
	}
}

// LoginChannelID
type AgentLineLoginChannelID struct {
	LineLoginChannelID string `json:"line_login_channel_id"`
}

func NewAgentLineLoginChannelID(lineLoginChannelID string) AgentLineLoginChannelID {
	return AgentLineLoginChannelID{
		LineLoginChannelID: lineLoginChannelID,
	}
}

type AgentBotInformation struct {
	AgentBot *entity.Bot `json:"agent_bot"`
}

func NewAgentBotInformation(agentBot *entity.Bot) AgentBotInformation {
	return AgentBotInformation{
		AgentBot: agentBot,
	}
}

type AgentListAndMaxPageAndIDList struct {
	MaxPageNumber uint            `json:"max_page_number"`
	IDList        []uint          `json:"id_list"`
	AgentList     []*entity.Agent `json:"agent_list"`
}

func NewAgentListAndMaxPageAndIDList(agents []*entity.Agent, maxPageNumber uint, idList []uint) AgentListAndMaxPageAndIDList {
	return AgentListAndMaxPageAndIDList{
		MaxPageNumber: maxPageNumber,
		IDList:        idList,
		AgentList:     agents,
	}
}
