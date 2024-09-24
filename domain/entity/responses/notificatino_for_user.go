package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type NotificationForUser struct {
	NotificationForUser *entity.NotificationForUser `json:"notification_for_user"`
}

func NewNotificationForUser(notificationForUser *entity.NotificationForUser) NotificationForUser {
	return NotificationForUser{
		NotificationForUser: notificationForUser,
	}
}

type NotificationForUserList struct {
	NotificationForUserList []*entity.NotificationForUser `json:"notification_for_user_list"`
}

func NewNotificationForUserList(notificationForUsers []*entity.NotificationForUser) NotificationForUserList {
	return NotificationForUserList{
		NotificationForUserList: notificationForUsers,
	}
}

type NotificationForUserListAndMaxPageNumber struct {
	NotificationForUserList []*entity.NotificationForUser `json:"notification_for_user_list"`
	MaxPageNumber           uint                          `json:"max_page_number"`
}

func NewNotificationForUserListAndMaxPageNumber(
	notificationForUsers []*entity.NotificationForUser,
	maxPageNumber uint,
) NotificationForUserListAndMaxPageNumber {
	return NotificationForUserListAndMaxPageNumber{
		NotificationForUserList: notificationForUsers,
		MaxPageNumber:           maxPageNumber,
	}
}

type UnwatchedNotificationForUserCount struct {
	UnwatchedNotificationForUserCount uint `json:"unwatched_notification_for_user_count"`
}

func NewUnwatchedNotificationForUserCount(unwatchedNotificationForUserCount uint) UnwatchedNotificationForUserCount {
	return UnwatchedNotificationForUserCount{
		UnwatchedNotificationForUserCount: unwatchedNotificationForUserCount,
	}
}
