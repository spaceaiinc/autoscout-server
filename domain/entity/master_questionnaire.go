package entity

// 志望度
var MyRanking Type = []string{
	0: "第1志望",
	1: "第1志望群",
	2: "第2志望",
	3: "第3志望",
	4: "第4志望以下",
}

// 他社の志望度
var OtherMyRanking Type = []string{
	0: "第1志望",
	1: "第1志望群",
	2: "第2志望群以下",
}

// 内定が出た場合の意向
var IntentionToJobOffer Type = []string{
	0: "他社選考を辞退して内定承諾",
	1: "内定を保留したい（●月●日まで）",
	2: "内定を辞退する",
	3: "その他",
}

// 選考継続希望有無
var ContinueSelectionType Type = []string{
	0: "選考継続を希望する",
  1: "選考を辞退する",
}

// 選考後アンケートの志望企業のフェーズ
var SelectionPhase Type = []string{
	0: "1次選考調整中",
	1: "1次選考実施",
	2: "2次選考調整中",
	3: "2次選考実施",
	4: "3次選考調整中",
	5: "3次選考実施",
	6: "最終選考調整中",
	7: "最終選考実施",
	8: "内定",
}


var Questionnaire []Type = []Type{
	MyRanking,
	IntentionToJobOffer,
	ContinueSelectionType,
	SelectionPhase,
}

