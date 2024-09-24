package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewNotificationForUserJSONPresenter(resp responses.NotificationForUser) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewNotificationForUserListJSONPresenter(resp responses.NotificationForUserList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewNotificationForUserListAndMaxPageNumberJSONPresenter(resp responses.NotificationForUserListAndMaxPageNumber) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewUnwatchedNotificationForUserCountJSONPresenter(resp responses.UnwatchedNotificationForUserCount) Presenter {
	return NewJSONPresenter(200, resp)
}
