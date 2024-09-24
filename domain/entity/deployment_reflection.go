package entity

import (
	"time"
)

type DeploymentReflection struct {
	ID           uint      `db:"id" json:"id"` // 重複しないID
	DeploymentID uint      `db:"deployment_id" json:"deployment_id"`
	AgentStaffID uint      `db:"agent_staff_id" json:"agent_staff_id"`
	IsReflected  bool      `db:"is_reflected" json:"is_reflected"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"` // 作成日時
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"` // 更新日時
}

func NewDeploymentReflection(
	deploymentID uint,
	agentStaffID uint,
	isReflected bool,
) *DeploymentReflection {
	return &DeploymentReflection{
		DeploymentID: deploymentID,
		AgentStaffID: agentStaffID,
		IsReflected:  isReflected,
	}
}
