package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

func NewInterviewAdjustmentTemplateJSONPresenter(resp responses.InterviewAdjustmentTemplate) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewInterviewAdjustmentTemplateListJSONPresenter(resp responses.InterviewAdjustmentTemplateList) Presenter {
	return NewJSONPresenter(200, resp)
}
