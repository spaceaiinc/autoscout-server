package entity

var StateOfEmployment Type = []string{
	"就業中",
	"離職中",
}

// 求職者のフェーズ（エントリー、面談調整中、面談予約、面談実施）
var PhaseForJobSeeker Type = []string{
	"エントリー",
	"面談案内済み（面談調整中）",
	"面談予約完了",
	"面談実施待ち",
	"面談実施済み（準備中）",
	"面談実施済み（稼働中）",
	"面談実施済み（リリース状態）",
	"サービス終了 / 転職活動終了",
}

var AppearanceForJobSeeker Type = []string{
	"1:大きな懸念が見受けられる",
	"2:気になる点が見受けられる",
	"3:可もなく不可もなく",
	"4:好印象を持てる",
	"5:印象値のみで内定が出る",
}

var CommunicationForJobSeeker Type = []string{
	"1:相手に不快感を与えるレベル",
	"2:一問一答が可能",
	"3:印象を左右しないレベル〈可もなく不可もなく〉",
	"4:比較的多くの人が好印象を抱くレベル",
	"5:万人が心地よい印象を抱くレベル",
}

var ThinkingForJobSeeker Type = []string{
	"1:終始、論理性に懸念を抱くレベル",
	"2:時折、論理性に懸念を抱くレベル",
	"3:可もなく不可もなく",
	"4:一定以上の論理性が窺えるレベル",
	"5:万人に論理的な印象を与えるレベル",
}

var SchoolCategory Type = []string{
	"中学",
	"高校",
	"専門",
	"短大",
	"高専",
	"大学",
	"大学院",
}

var FinalEducationForJobSeeker Type = []string{
	"中学卒",
	"高校卒",
	"専門学校卒",
	"短期大学卒",
	"高等専門学校卒",
	"大学卒",
	"大学院卒",
}

var CollegeRankForJobSeeker Type = []string{
	"旧帝",
	"早慶上",
	"KKDR・GMARCH",
	"SKKR・NTKS",
	"SSTT・DTATE",
	"その他",
}

var StudyCategoryForJobSeeker Type = []string{
	"理系",
	"文系",
}

var JobChangeForJobSeeker Type = []string{
	"0回",
	"1回",
	"2回",
	"3回",
	"4回",
	"5回以上",
}

var TransferForJobSeeker Type = []string{
	"転勤可",
	"転勤可（条件あり）",
	"転勤不可",
}

// 入学, 編集学, 転入学
var FirstStatusForStudentHistory Type = []string{
	"入学",
	"編入学",
	"転入学",
}

var LastStatusForStudentHistory Type = []string{
	"卒業",
	"中退",
	"退学",
	"卒業見込み",
	"修了",    //院卒用
	"修了見込み", //院卒用
}

// 入社,入行,入局,入庁,入省,入職,出向,転籍,帰任,拝命
var FirstStatusForWorkHistory Type = []string{
	"入社",
	"入行",
	"入局",
	"入庁",
	"入省",
	"入職",
	"出向",
	"転籍",
	"帰任",
	"拝命",
	// ""雇用形態の変更" = 99
}

// 一身上の都合により退職,派遣期間満了につき退職,契約期間満了につき退職,会社都合により退職
var LastStatusForWorkHistory Type = []string{
	"現在に至る",
	"一身上の都合により退職",
	"派遣期間満了につき退職",
	"契約期間満了につき退職",
	"会社都合により退職",
	// ""雇用形態の変更" = 99
}

var WorkHistoryStatus Type = []string{
	"入社",
	"退職",
}

var JobHuntingState Type = []string{
	"活動中",
	"活動終了",
}

var GenderForJobSeeker Type = []string{
	"男性",
	"女性",
	// index99 "不問",
}

var NationalityForJobSeeker Type = []string{
	"日本国籍",
	"外国国籍",
	// index99 "不問",
}
