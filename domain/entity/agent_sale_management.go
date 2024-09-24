package entity

import (
	"time"
)

type AgentSaleManagement struct {
	ID         uint      `json:"id" db:"id"`
	AgentID    uint      `json:"agent_id" db:"agent_id"`
	FiscalYear string    `json:"fiscal_year" db:"fiscal_year"`
	IsOpen     bool      `json:"is_open" db:"is_open"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
	UpdatedAt  time.Time `db:"updated_at" json:"-"`

	AgentMonthlySales         []AgentMonthlySale         `json:"agent_monthly_sales" db:"agent_monthly_sales"`
	AgentStaffSaleManagements []AgentStaffSaleManagement `json:"staff_sale_managements" db:"staff_sale_managements"`
}

func NewAgentSaleManagement(
	agentID uint,
	fiscalYear string,
	isOpen bool,
) *AgentSaleManagement {
	return &AgentSaleManagement{
		AgentID:    agentID,
		FiscalYear: fiscalYear,
		IsOpen:     isOpen,
	}
}
