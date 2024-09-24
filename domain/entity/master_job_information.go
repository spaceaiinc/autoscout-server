package entity

var RecruitmentState = []string{
	"Open",
	"Close",
	"下書き",
}

var OpenOrClose = []string{
	"Open",
	"Close",
}

var EmploymentPeriod = []string{
	"定め有り",
	"定め無し",
}

var MedicalHistory Type = []string{
	"不問",
	"不可",
}

var AppearanceForJobInfo Type = []string{
	"1:アピアランスによる採用判断はなし",
	"2:気になる点が見受けられるが、指摘で改善が見られれば可",
	"3:可もなく不可もないレベルで及第点",
	"4:比較的多くの人が好印象を持てるレベルを求める",
	"5:万人が心地よい印象を抱くレベルを求める",
}

var CommunicationForJobInfo Type = []string{
	"1:コミュニケーションスキルに難点が見受けられるレベルでも可",
	"2:一問一答が可能レベルを求める",
	"3:可もなく不可もないレベルで及第点",
	"4:比較的多くの人が好印象を持てるレベルを求める",
	"5:万人が心地よい印象を抱くレベルを求める",
}

var ThinkingForJobInfo Type = []string{
	"1:終始、論理性に懸念を抱くレベルでも可",
	"2:時折、論理性に懸念を抱くレベルでも可",
	"3:可もなく不可もなくで及第点",
	"4:会話から一定以上の論理性が窺えるレベルを求める",
	"5:万人に論理的な印象を与えるレベルを求める",
}

var JobChangeForJobInfo Type = []string{
	"0回のみ",
	"1回まで",
	"2回まで",
	"3回まで",
	"4回まで",
	"5回まで",
	// index=99 "不問",
}

var FinalEducationForJobInfo Type = []string{
	"中卒以上",
	"高卒以上",
	"専卒以上（短卒除く）",
	"短卒以上（専卒除く）",
	"専卒・短卒以上",
	"高専卒以上",
	"大卒以上",
	"院卒以上",
}

// 旧帝・早慶上,KKDR・MARCH,SKKR・NTKS,SSTT・DTATE,不問
var CollegeRankForJobInfo Type = []string{
	"旧帝・早慶上",
	"KKDR・GMARCH",
	"SKKR・NTKS",
	"SSTT・DTATE",
	// index=99 "不問",
}

var StudyCategoryForJobInfo Type = []string{
	"理系尚可",
	"理系のみ",
	// index=99 "不問",
}

var NationalityForJobInfo Type = []string{
	"日本国籍のみ",
	"外国国籍のみ",
	// index=99 "不問",
}

// 企業への連絡方法（マスクレジュメ送付時）　許可する,許可しない(RA経由),マスクレジュメ打診不要
var ContactManner Type = []string{
	"許可する",
	"許可しない(RA経由)",
	"マスクレジュメ打診不要",
}

// 受動喫煙対策　あり・屋内禁煙,あり・屋内禁煙・敷地内禁煙,あり・屋内禁煙・敷地内禁煙（屋外に喫煙場所設置）,あり・喫煙室設置・喫煙可能室設置,あり・喫煙室設置・喫煙専用室設置,あり・喫煙室設置・加熱式たばこ専用喫煙室設置,あり・喫煙室設置・喫煙目的室設置,あり・喫煙室設置・喫煙可の宿泊室あり,なし
var PassiveSmoking Type = []string{
	"あり・屋内禁煙",
	"あり・屋内禁煙・敷地内禁煙",
	"あり・屋内禁煙・敷地内禁煙（屋外に喫煙場所設置）",
	"あり・喫煙室設置・喫煙可能室設置",
	"あり・喫煙室設置・喫煙専用室設置",
	"あり・喫煙室設置・加熱式たばこ専用喫煙室設置",
	"あり・喫煙室設置・喫煙目的室設置",
	"あり・喫煙室設置・喫煙可の宿泊室あり",
	"なし",
}

// 求人募集性別{0: 男性のみ, 1: 女性のみ, 2男性尚可, 3: 女性尚可, 99: 不問}
var GenderForJobInfo Type = []string{
	"男性のみ",
	"女性のみ",
	"男性尚可",
	"女性尚可",
	// index99 "不問",
}

var EmploymentStatusForJobInfo Type = []string{
	"正社員",
	"契約社員",
	"派遣社員",
	"紹介予定派遣",
	"アルバイト・パート",
	"正社員（無期雇用派遣）",
	// "自営業",
	// "業務委託",
	// index = 99 不問
}

/** 内定率*/
var OfferRate Type = []string{
	"1~3%",
	"3~5%",
	"5~10%",
	"10~15%",
	"15~20%",
	"20%以上",
}

/** 書類通過率 */
var DocumentPassingRate Type = []string{
	"20~30%",
	"30~40%",
	"40~50%",
	"50~75%",
	"75~100%",
}

/** 直近の応募数 */
var NumberOfRecentApplications Type = []string{
	"1~5件",
	"5~10件",
	"10~15件",
	"15件以上",
}
