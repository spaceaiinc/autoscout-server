package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type LabelListForSendingJobSeekerManagement struct {
	StaffList     []*entity.LabelAndValue `json:"staff_list"`      // 担当者
	SendAgentList []*entity.LabelAndValue `json:"send_agent_list"` // 送客先
	SenderList    []*entity.LabelAndValue `json:"sender_list"`     // 送客元
}

func NewLabelListForSendingJobSeekerManagement(staffList []*entity.LabelAndValue, sendAgentList []*entity.LabelAndValue, senderList []*entity.LabelAndValue) LabelListForSendingJobSeekerManagement {
	return LabelListForSendingJobSeekerManagement{
		StaffList:     staffList,
		SendAgentList: sendAgentList,
		SenderList:    senderList,
	}
}
