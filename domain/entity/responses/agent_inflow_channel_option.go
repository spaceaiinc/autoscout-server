package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type AgentInflowChannelOption struct {
	AgentInflowChannelOption *entity.AgentInflowChannelOption `json:"agent_inflow_channel_option"`
}

func NewAgentInflowChannelOption(agentInflowChannelOption *entity.AgentInflowChannelOption) AgentInflowChannelOption {
	return AgentInflowChannelOption{
		AgentInflowChannelOption: agentInflowChannelOption,
	}
}

type AgentInflowChannelOptionList struct {
	AgentInflowChannelOptionList []*entity.AgentInflowChannelOption `json:"agent_inflow_channel_option_list"`
}

func NewAgentInflowChannelOptionList(agentInflowChannelOptions []*entity.AgentInflowChannelOption) AgentInflowChannelOptionList {
	return AgentInflowChannelOptionList{
		AgentInflowChannelOptionList: agentInflowChannelOptions,
	}
}
