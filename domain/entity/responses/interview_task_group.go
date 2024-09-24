package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type InterviewTaskGroup struct {
	InterviewTaskGroup *entity.InterviewTaskGroup `json:"interview_task_group"`
}

func NewInterviewTaskGroup(taskGroup *entity.InterviewTaskGroup) InterviewTaskGroup {
	return InterviewTaskGroup{
		InterviewTaskGroup: taskGroup,
	}
}

type InterviewTaskGroupList struct {
	InterviewTaskGroupList []*entity.InterviewTaskGroup `json:"interview_task_group_list"`
}

func NewInterviewTaskGroupList(taskGroups []*entity.InterviewTaskGroup) InterviewTaskGroupList {
	return InterviewTaskGroupList{
		InterviewTaskGroupList: taskGroups,
	}
}
