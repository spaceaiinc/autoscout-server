package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type EnterpriseIndustry struct {
	ID           uint      `db:"id" json:"id"`
	EnterpriseID uint      `db:"enterprise_id" json:"enterprise_id"`
	Industry     null.Int  `db:"industry" json:"industry"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}

func NewEnterpriseIndustry(
	enterpriseID uint,
	industry null.Int,
) *EnterpriseIndustry {
	return &EnterpriseIndustry{
		EnterpriseID: enterpriseID,
		Industry:     industry,
	}
}
