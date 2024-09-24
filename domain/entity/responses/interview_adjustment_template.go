package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type InterviewAdjustmentTemplate struct {
	InterviewAdjustmentTemplate *entity.InterviewAdjustmentTemplate `json:"interview_adjustment_template"`
}

func NewInterviewAdjustmentTemplate(task *entity.InterviewAdjustmentTemplate) InterviewAdjustmentTemplate {
	return InterviewAdjustmentTemplate{
		InterviewAdjustmentTemplate: task,
	}
}

type InterviewAdjustmentTemplateList struct {
	InterviewAdjustmentTemplateList []*entity.InterviewAdjustmentTemplate `json:"interview_adjustment_template_list"`
}

func NewInterviewAdjustmentTemplateList(tasks []*entity.InterviewAdjustmentTemplate) InterviewAdjustmentTemplateList {
	return InterviewAdjustmentTemplateList{
		InterviewAdjustmentTemplateList: tasks,
	}
}
