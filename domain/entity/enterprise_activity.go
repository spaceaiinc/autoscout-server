package entity

import (
	"time"
)

type EnterpriseActivity struct {
	ID           uint      `db:"id" json:"id"`
	EnterpriseID uint      `db:"enterprise_id" json:"enterprise_id"`
	Content      string    `db:"content" json:"content"`
	AddedAt      time.Time `db:"added_at" json:"added_at"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func NewEnterpriseActivity(
	enterpriseID uint,
	content string,
	addedAt time.Time,
) *EnterpriseActivity {
	return &EnterpriseActivity{
		EnterpriseID: enterpriseID,
		Content:      content,
		AddedAt:      addedAt,
	}
}

type CreateEnterpriseActivityParam struct {
	EnterpriseID uint      `db:"enterprise_id" json:"enterprise_id" validate:"required"`
	Content      string    `db:"content" json:"content" validate:"required"`
	AddedAt      time.Time `db:"added_at" json:"added_at"`
}
