package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// Motoyuiからのお知らせを管理するテーブル
type NotificationForUser struct {
	ID        uint      `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`   // お知らせのタイトル
	Body      string    `db:"body" json:"body"`     // お知らせの本文
	Target    null.Int  `db:"target" json:"target"` // 送信対象（0: 全てのユーザー, 1: CRMのみ, 2: 送客のみ）
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// お知らせを確認したユーザー
	UserNotificationViews []UserNotificationView `db:"-" json:"user_notification_views"`
}

func NewNotificationForUser(
	title string,
	body string,
	target null.Int,
) *NotificationForUser {
	return &NotificationForUser{
		Title:  title,
		Body:   body,
		Target: target,
	}
}

// 送信対象タイプ(0: 全てのユーザー, 1: CRMのみ, 2: 送客のみ)
type NotificationForUserTarget int64

const (
	NotificationForUserTargetAll NotificationForUserTarget = iota
	NotificationForUserTargetCRM
	NotificationForUserTargetSending
)
