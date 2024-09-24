package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewTaskJSONPresenter(resp responses.Task) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewTaskListJSONPresenter(resp responses.TaskList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewTaskListByTypeJSONPresenter(resp responses.TaskListByType) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewRATaskListByTypeJSONPresenter(resp responses.RATaskListByType) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewCATaskListByTypeJSONPresenter(resp responses.CATaskListByType) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewTaskPageListAndMaxPageJSONPresenter(resp responses.TaskListAndMaxPage) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewActiveTaskCountJSONPresenter(resp responses.ActiveTaskCount) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewTaskListAndJobSeekerListJSONPresenter(resp responses.TaskListAndJobSeekerList) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewJobSeekerTaskListJSONPresenter(resp responses.JobSeekerTaskList) Presenter {
	return NewJSONPresenter(200, resp)
}
