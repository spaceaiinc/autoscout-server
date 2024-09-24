package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type EvaluationPoint struct {
	EvaluationPoint *entity.EvaluationPoint `json:"evaluation_point"`
}

func NewEvaluationPoint(evaluationPoint *entity.EvaluationPoint) EvaluationPoint {
	return EvaluationPoint{
		EvaluationPoint: evaluationPoint,
	}
}

type EvaluationPointList struct {
	EvaluationPointList []*entity.EvaluationPoint `json:"evaluation_point_list"`
}

func NewEvaluationPointList(evaluationPoints []*entity.EvaluationPoint) EvaluationPointList {
	return EvaluationPointList{
		EvaluationPointList: evaluationPoints,
	}
}
