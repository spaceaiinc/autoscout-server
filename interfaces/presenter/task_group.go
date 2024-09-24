package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewTaskGroupJSONPresenter(resp responses.TaskGroup) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewTaskGroupListJSONPresenter(resp responses.TaskGroupList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewTaskGroupPageListAndMaxPageJSONPresenter(resp responses.TaskGroupListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewSoundOutGroupListJSONPresenter(resp responses.SoundOutGroupList) Presenter {
	return NewJSONPresenter(200, resp)
}
