package entity

type LabelAndValue struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

func NewLabelAndValue(
	value uint,
	label string,
) *LabelAndValue {
	return &LabelAndValue{
		Value: value,
		Label: label,
	}
}
