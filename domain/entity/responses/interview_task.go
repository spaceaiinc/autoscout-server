package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type InterviewTask struct {
	InterviewTask *entity.InterviewTask `json:"interview_task"`
}

func NewInterviewTask(task *entity.InterviewTask) InterviewTask {
	return InterviewTask{
		InterviewTask: task,
	}
}

type InterviewTaskList struct {
	InterviewTaskList []*entity.InterviewTask `json:"interview_task_list"`
}

func NewInterviewTaskList(tasks []*entity.InterviewTask) InterviewTaskList {
	return InterviewTaskList{
		InterviewTaskList: tasks,
	}
}
