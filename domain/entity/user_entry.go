package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type UserEntry struct {
	ID          uint      `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	ServiceType null.Int  `db:"service_type" json:"service_type"`
	IsProcessed bool      `db:"is_processed" json:"is_processed"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func NewUserEntry(
	userID string,
	serviceType null.Int,
) *UserEntry {
	return &UserEntry{
		UserID:      userID,
		ServiceType: serviceType,
	}
}
