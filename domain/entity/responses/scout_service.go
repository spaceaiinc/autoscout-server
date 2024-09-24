package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type ScoutService struct {
	ScoutService *entity.ScoutService `json:"scout_service"`
}

func NewScoutService(scoutService *entity.ScoutService) ScoutService {
	return ScoutService{
		ScoutService: scoutService,
	}
}

type ScoutServiceList struct {
	ScoutServiceList []*entity.ScoutService `json:"scout_service_list"`
}

func NewScoutServiceList(scoutServices []*entity.ScoutService) ScoutServiceList {
	return ScoutServiceList{
		ScoutServiceList: scoutServices,
	}
}
