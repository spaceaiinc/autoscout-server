package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type Dashboard struct {
	Dashboard *entity.Dashboard `json:"dashboard"`
}

func NewDashboard(
	dashboard *entity.Dashboard,
) Dashboard {
	return Dashboard{
		Dashboard: dashboard,
	}
}
