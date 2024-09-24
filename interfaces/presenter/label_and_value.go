package presenter

import "github.com/spaceaiinc/autoscout-server/domain/entity/responses"

// 送客管理ページの検索に使用する値を返すのに使用
func NewLabelListForSendingJobSeekerManagementJSONPresenter(resp responses.LabelListForSendingJobSeekerManagement) Presenter {
	return NewJSONPresenter(200, resp)
}
