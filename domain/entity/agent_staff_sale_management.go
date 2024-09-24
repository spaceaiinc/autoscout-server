package entity

import (
	"time"
)

type AgentStaffSaleManagement struct {
	ID           uint      `json:"id" db:"id"`
	ManagementID uint      `json:"management_id" db:"management_id"`
	AgentID      uint      `json:"agent_id" db:"agent_id"`
	AgentStaffID uint      `json:"agent_staff_id" db:"agent_staff_id"`
	StaffName    string    `db:"staff_name" json:"staff_name"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
	FiscalYear   string    `json:"fiscal_year" db:"fiscal_year"`

	IsOpen            bool                    `json:"is_open" db:"is_open"`
	StaffMonthlySales []AgentStaffMonthlySale `json:"staff_monthly_sales" db:"staff_monthly_sales"`
}

func NewAgentStaffSaleManagement(
	managementID uint,
	agentStaffID uint,
) *AgentStaffSaleManagement {
	return &AgentStaffSaleManagement{
		ManagementID: managementID,
		AgentStaffID: agentStaffID,
	}
}

type CreateOrUpdateStaffMonthlyManagementParam struct {
	ManagementID      uint                    `json:"management_id" validate:"required"`
	AgentStaffID      uint                    `json:"agent_staff_id" validate:"required"`
	StaffMonthlySales []AgentStaffMonthlySale `json:"staff_monthly_sales" validate:"required"`
}
