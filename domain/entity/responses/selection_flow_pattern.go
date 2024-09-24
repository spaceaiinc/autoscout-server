package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SelectionFlowPattern struct {
	SelectionFlowPattern *entity.JobInformationSelectionFlowPattern `json:"selection_flow_pattern"`
}

func NewSelectionFlowPattern(selectionFlowPattern *entity.JobInformationSelectionFlowPattern) SelectionFlowPattern {
	return SelectionFlowPattern{
		SelectionFlowPattern: selectionFlowPattern,
	}
}

type SelectionFlowPatternList struct {
	SelectionFlowPatternList []*entity.JobInformationSelectionFlowPattern `json:"selection_flow_pattern_list"`
}

func NewSelectionFlowPatternList(selectionFlowPatterns []*entity.JobInformationSelectionFlowPattern) SelectionFlowPatternList {
	return SelectionFlowPatternList{
		SelectionFlowPatternList: selectionFlowPatterns,
	}
}
