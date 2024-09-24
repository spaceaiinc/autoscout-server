package entity

import (
	"time"

	"github.com/google/uuid"
)

type AgentAlliance struct {
	ID                  uint      `db:"id" json:"id"`
	Agent1ID            uint      `db:"agent1_id" json:"agent1_id"`                         // エージェント1のid
	Agent2ID            uint      `db:"agent2_id" json:"agent2_id"`                         //エージェント2のid
	Agent1Request       bool      `db:"agent1_request" json:"agent1_request"`               // エージェント1の申請状況, 承認）
	Agent2Request       bool      `db:"agent2_request" json:"agent2_request"`               // エージェント2の申請状況, 承認）
	Agent1CancelRequest bool      `db:"agent1_cancel_request" json:"agent1_cancel_request"` // エージェント1の解除申請状況, 承認）
	Agent2CancelRequest bool      `db:"agent2_cancel_request" json:"agent2_cancel_request"` // エージェント2の解除申請状況, 承認）
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`

	AgentName      string `db:"agent_name" json:"agent_name"`           // エージェント名
	OfficeLocation string `db:"office_location" json:"office_location"` // 住所
	Representative string `db:"representative" json:"representative"`   // 代表者
	Establish      string `db:"establish" json:"establish"`             // 設立（年月）
}

func NewAgentAlliance(
	agent1ID uint,
	agent2ID uint,
	agent1Request bool,
	agent2Request bool,
	agent1CancelRequest bool,
	agent2CancelRequest bool,
) *AgentAlliance {
	return &AgentAlliance{
		Agent1ID:            agent1ID,
		Agent2ID:            agent2ID,
		Agent1Request:       agent1Request,
		Agent2Request:       agent2Request,
		Agent1CancelRequest: agent1CancelRequest,
		Agent2CancelRequest: agent2CancelRequest,
	}
}

type CreateAgentAllianceParam struct {
	MyAgentID     uint      `json:"my_agent_id"`  // 自社のID
	RequestUUID   uuid.UUID `json:"request_uuid"` // 申請先のuuid
	Agent1Request bool      `json:"agent1_request"`
	Agent2Request bool      `json:"agent2_request"`
}

type UpdateAgentAllianceRequestStateParam struct {
	AgentAllianceID uint `json:"agent_alliance_id" validate:"required"`
	AgentID         uint `json:"agent_id" validate:"required"`
}

type MyAgentIDAndOtherAgentIDParam struct {
	MyAgentID    uint `json:"my_agent_id"`    // 自社のID
	OtherAgentID uint `json:"other_agent_id"` // 他社のID
}

type MyAgentIDAndOtherAgentIDListParam struct {
	MyAgentID        uint   `json:"my_agent_id"`         // 自社のID
	OtherAgentIDList []uint `json:"other_agent_id_list"` // 他社のID
}

type AllianceAndChatGroupParam struct {
	MyAgentID                uint   `json:"my_agent_id"`                   // 自社のID
	UnAppliedAgentIDList     []uint `json:"un_applied_agent_id_list"`      // アライアンス未申請のIDリスト
	UnCreatedChatAgentIDList []uint `json:"un_created_chat_agent_id_list"` // チャットグループ未作成のIDリスト
}

type UpdateAgentAllianceCancelRequestParam struct {
	AgentAllianceID uint `json:"agent_alliance_id" validate:"required"`
	MyAgentID       uint `json:"my_agent_id" validate:"required"`
}
