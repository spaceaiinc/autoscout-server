package entity

import "time"

// お知らせを確認したユーザーをを管理するテーブル（ユーザーが閲覧したらレコードを追加）
type UserNotificationView struct {
	ID              uint      `db:"id" json:"id"`
	NotificationID  uint      `db:"notification_id" json:"notification_id"` // notification_for_usersテーブルのID
	AgentStaffID    uint      `db:"agent_staff_id" json:"agent_staff_id"`   // agent_staffsテーブルのID
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

func NewUserNotificationView(
	notificationID uint,
	agentStaffID uint,
) *UserNotificationView {
	return &UserNotificationView{
		NotificationID:  notificationID,
		AgentStaffID:    agentStaffID,
	}
}