package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type Task struct {
	Task *entity.Task `json:"task"`
}

func NewTask(task *entity.Task) Task {
	return Task{
		Task: task,
	}
}

type TaskList struct {
	TaskList []*entity.Task `json:"task_list"`
}

func NewTaskList(tasks []*entity.Task) TaskList {
	return TaskList{
		TaskList: tasks,
	}
}

type TaskListByType struct {
	TaskListByType *entity.TaskListByType `json:"task_list_by_type"`
}

func NewTaskListByType(taskListByType *entity.TaskListByType) TaskListByType {
	return TaskListByType{
		TaskListByType: taskListByType,
	}
}

type RATaskListByType struct {
	RATaskListByType *entity.RATaskListByType `json:"ra_task_list_by_type"`
}

func NewRATaskListByType(raTaskListByType *entity.RATaskListByType) RATaskListByType {
	return RATaskListByType{
		RATaskListByType: raTaskListByType,
	}
}

type CATaskListByType struct {
	CATaskListByType *entity.CATaskListByType `json:"ca_task_list_by_type"`
}

func NewCATaskListByType(caTaskListByType *entity.CATaskListByType) CATaskListByType {
	return CATaskListByType{
		CATaskListByType: caTaskListByType,
	}
}

type TaskListAndMaxPage struct {
	MaxPageNumber uint           `json:"max_page_number"`
	TaskList      []*entity.Task `json:"task_list"`
}

func NewTaskListAndMaxPage(enterpriseProfiles []*entity.Task, maxPageNumber uint) TaskListAndMaxPage {
	return TaskListAndMaxPage{
		MaxPageNumber: maxPageNumber,
		TaskList:      enterpriseProfiles,
	}
}

type ActiveTaskCount struct {
	ActiveTaskCount uint `json:"active_task_count"`
}

func NewActiveTaskCount(activeTaskCount uint) ActiveTaskCount {
	return ActiveTaskCount{
		ActiveTaskCount: activeTaskCount,
	}
}

type TaskListAndJobSeekerList struct {
	TaskList      []*entity.Task      `json:"task_list"`
	JobSeekerList []*entity.JobSeeker `json:"job_seeker_list"`
}

func NewTaskListAndJobSeekerList(tasks []*entity.Task, jobSeekers []*entity.JobSeeker) TaskListAndJobSeekerList {
	return TaskListAndJobSeekerList{
		TaskList:      tasks,
		JobSeekerList: jobSeekers,
	}
}

type JobSeekerTaskList struct {
	TaskList []*entity.JobSeekerTask `json:"task_list"`
}

func NewJobSeekerTaskList(tasks []*entity.JobSeekerTask) JobSeekerTaskList {
	return JobSeekerTaskList{
		TaskList: tasks,
	}
}
