package entity

import (
	"strconv"
	"strings"

	"gopkg.in/guregu/null.v4"
)

// 有or無のチェック
func GetIntAvailable(label string) null.Int {
	var (
		availableNumber = null.NewInt(0, false)
	)

	// 空文字の場合は無しとする
	if label == "" {
		return null.NewInt(1, true)
	} else if label == "あり" {
		return null.NewInt(0, true)
	}

	return availableNumber
}

func GetIntContractPhase(label string) null.Int {
	var (
		output = null.NewInt(0, false)
	)

	for masterI, masterL := range ContractPhase {
		if masterL == label {
			output = null.NewInt(int64(masterI), true)
			break
		}
	}

	return output
}

// 求職者のエントリーフェーズ
func GetIntPhaseForJobSeeker(label string) null.Int {
	var (
		output = null.NewInt(0, false)
	)

	for masterI, masterL := range PhaseForJobSeeker {
		if masterL == label {
			output = null.NewInt(int64(masterI), true)
			break
		}
	}

	return output
}

// 転職状況
func GetIntJobHuntingState(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "" {
		return null.NewInt(0, true)
	}

	for masterI, masterL := range JobHuntingState {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// 就業状況
func GetIntStateOfEmployment(label string) null.Int {
	var (
		output = null.NewInt(0, false)
	)

	for masterI, masterL := range StateOfEmployment {
		if masterL == label {
			output = null.NewInt(int64(masterI), true)
			break
		}
	}

	return output
}

// 転勤可否　求職者用
func GetIntTransferForJobSeeker(label string) null.Int {
	var (
		output = null.NewInt(0, false)
	)

	for masterI, masterL := range TransferForJobSeeker {
		if masterL == label {
			output = null.NewInt(int64(masterI), true)
			break
		}
	}

	return output
}

// 転職回数　求職者用
func GetIntJobChangeForJobSeeker(label string) null.Int {
	var (
		output = null.NewInt(0, false)
	)

	for masterI, masterL := range JobChangeForJobSeeker {
		if masterL == label {
			output = null.NewInt(int64(masterI), true)
			break
		}
	}

	return output
}

// 入社可能時期
func GetIntJoinCompanyPeriod(label string) null.Int {
	var (
		output = null.NewInt(0, false)
	)

	for masterI, masterL := range JoinCompanyPeriod {
		if masterL == label {
			output = null.NewInt(int64(masterI), true)
			break
		}
	}

	return output
}

func GetIntIndustry(label string) null.Int {

	if label == "まだ定まっていない" || label == "不問" {
		return null.NewInt(9999, true)
	}

	for _, industry := range Industry {
		for key, value := range industry {
			if value == label {
				return null.NewInt(int64(key), true)
			}
		}
	}

	return null.NewInt(0, false)
}

// 業界大分類Stringから業界小分類Intを取得
func GetIntIndustryFromBigCategory(label string) []uint {
	var (
		list []uint
	)

	if label == "まだ定まっていない" || label == "不問" {
		list = append(list, 9999)
		return list
	}

	for key := range Industry[label] {
		list = append(list, key)
	}

	return list
}

func GetIntOccupation(label string) null.Int {

	if label == "まだ定まっていない" || label == "不問" {
		return null.NewInt(9999, true)
	}

	for _, occupation := range Occupation {
		for key, value := range occupation {
			if value == label {
				return null.NewInt(int64(key), true)
			}
		}
	}

	return null.NewInt(0, false)
}

// 職種大分類Stringから職種小分類Intを取得
func GetIntOccupationFromBigCategory(label string) []uint {
	var (
		list []uint
	)

	if label == "まだ定まっていない" || label == "不問" {
		list = append(list, 9999)
		return list
	}

	for key := range Occupation[label] {
		list = append(list, key)
	}

	return list
}

func GetIntPublicOffering(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range PublicOffering {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntRecruitmentState(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "" {
		return null.NewInt(1, true)
	}

	for masterI, masterL := range RecruitmentState {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntOpenOrClose(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range OpenOrClose {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntPrefecture(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "全国各地" || label == "全国" || label == "不問" {
		return null.NewInt(9999, true)
	}

	for masterI, masterL := range Prefecture {
		if masterL == label {
			return null.NewInt(int64(masterI), true)
		}
	}

	return number
}

func GetIntUserStatus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range UserStatus {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntEmploymentStatus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range EmploymentStatus {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntEmploymentStatusForJobInfo(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range EmploymentStatusForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntFinalEducationForJobInfo(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range FinalEducationForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntFinalEducationForJobSeeker(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range FinalEducationForJobSeeker {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntSchoolLevelForJobInfo(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range CollegeRankForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntSchoolLevelForJobSeeker(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range CollegeRankForJobSeeker {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntStudyCategoryForJobInfo(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "" || label == "不問" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range StudyCategoryForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntStudyCategoryForJobSeeker(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range StudyCategoryForJobSeeker {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntSchoolCategoryForJobSeeker(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range SchoolCategory {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntFirstStatusForStudentHistory(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range FirstStatusForStudentHistory {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntLastStatusForStudentHistory(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range LastStatusForStudentHistory {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntFirstStatusForWorkHistory(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "雇用形態の変更" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range FirstStatusForWorkHistory {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntLastStatusForWorkHistory(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "雇用形態の変更" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range LastStatusForWorkHistory {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntJobChange(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range JobChangeForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntSocialExperienceType(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range EmploymentStatus {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntLanguageType(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range LanguageType {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntLanguageLevel(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range LanguageLevel {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// func GetIntManagement(label string) null.Int {
// 	var (
// 		number = null.NewInt(0, false)
// 	)

// 	for index, value := range Management {
// 		if value == label {
// 			number = null.NewInt(int64(masterI), true)
// 			break
// 		}
// 	}

// 	return number
// }

func GetIntLicenseType(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for masterI, masterL := range LicenseType {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// excelSkill
func GetIntExcelSkill(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range ExcelSkill {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// wordSkill
func GetIntWordSkill(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range WordSkill {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// powerPointSkill
func GetIntPowerPointSkill(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range PowerPointSkill {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// pcTool
func GetIntPCTool(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(1000, true)
	} else if label == "その他会計ソフト" {
		return null.NewInt(99, true)
	}

	for masterI, masterL := range PCTool {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntDevelopmentType(categoryInt int, label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	for categoryIndex := range DevelopmentCategory {
		if categoryIndex == categoryInt {
			for typeIndex, typeValue := range DevelopmentTypeList[categoryIndex] {
				if typeValue == label {
					number = null.NewInt(int64(typeIndex), true)
					break
				}
			}
		}
	}

	return number
}

func GetIntFeature(label string) null.Int {
	var (
		featureNumber = null.NewInt(0, false)
	)

	for index, feature := range JobFeature {
		if feature == label {
			featureNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return featureNumber
}

func GetIntHolidayForJobSeeker(label string) null.Int {
	var (
		holidayNumber = null.NewInt(0, false)
	)

	for index, holiday := range HolidayForJobSeeker {
		if holiday == label {
			holidayNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return holidayNumber
}

func GetIntHolidayForJobInfo(label string) null.Int {
	var (
		holidayNumber = null.NewInt(0, false)
	)

	for index, holiday := range HolidayForJobInfo {
		if holiday == label {
			holidayNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return holidayNumber
}

func GetIntCompanyScale(label string) null.Int {
	var (
		holidayNumber = null.NewInt(0, false)
	)

	if label == "不問" {
		return null.NewInt(99, true)
	}

	for index, holiday := range CompanyScale {
		if holiday == label {
			holidayNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return holidayNumber
}

func GetIntPassiveSmoking(label string) null.Int {
	var (
		passiveSmokingNumber = null.NewInt(0, false)
	)

	for index, passiveSmoking := range PassiveSmoking {
		if passiveSmoking == label {
			passiveSmokingNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return passiveSmokingNumber
}

func GetIntGenderForJobSeeker(label string) null.Int {
	var (
		genderNumber = null.NewInt(0, false)
	)

	// if label == "不問" || label == "" {
	// 	return null.NewInt(99, true)
	// }

	for index, gender := range GenderForJobSeeker {
		if gender == label {
			genderNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return genderNumber
}

func GetIntConditionOrNot(label string) null.Int {
	var (
		genderNumber = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	}

	for index, gender := range ConditionOrNot {
		if gender == label {
			genderNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return genderNumber
}

func GetIntGenderForJobInfo(label string) null.Int {
	var (
		genderNumber = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	}

	for index, gender := range GenderForJobInfo {
		if gender == label {
			genderNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return genderNumber
}

func GetIntMedicalHistory(label string) null.Int {
	var (
		medicalHistoryNumber = null.NewInt(0, false)
	)

	for index, medicalHistory := range MedicalHistory {
		if medicalHistory == label {
			medicalHistoryNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return medicalHistoryNumber
}

func GetIntNationalityForJobSeeker(label string) null.Int {
	var (
		nationalityNumber = null.NewInt(0, false)
	)

	if label == "不問" {
		nationalityNumber = null.NewInt(99, true)
		return nationalityNumber
	}

	for index, nationality := range NationalityForJobSeeker {
		if nationality == label {
			nationalityNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return nationalityNumber
}

func GetIntNationalityForJobInfo(label string) null.Int {
	var (
		nationalityNumber = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	}

	for index, nationality := range NationalityForJobInfo {
		if nationality == label {
			nationalityNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return nationalityNumber
}

func GetIntAppearanceForJobSeeker(label string) null.Int {
	var (
		appearanceNumber = null.NewInt(0, false)
	)

	for index, appearance := range AppearanceForJobSeeker {
		if appearance == label {
			appearanceNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return appearanceNumber
}

func GetIntCommunicationForJobSeeker(label string) null.Int {
	var (
		communicationNumber = null.NewInt(0, false)
	)

	for index, communication := range CommunicationForJobSeeker {
		if communication == label {
			communicationNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return communicationNumber
}

func GetIntThinkingForJobSeeker(label string) null.Int {
	var (
		thinkingNumber = null.NewInt(0, false)
	)

	for index, thinking := range ThinkingForJobSeeker {
		if thinking == label {
			thinkingNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return thinkingNumber
}

func GetIntAppearanceForJobInfo(label string) null.Int {
	var (
		appearanceNumber = null.NewInt(0, false)
	)

	for index, appearance := range AppearanceForJobInfo {
		if appearance == label {
			appearanceNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return appearanceNumber
}

func GetIntCommunicationForJobInfo(label string) null.Int {
	var (
		communicationNumber = null.NewInt(0, false)
	)

	for index, communication := range CommunicationForJobInfo {
		if communication == label {
			communicationNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return communicationNumber
}

func GetIntThinkingForJobInfo(label string) null.Int {
	var (
		thinkingNumber = null.NewInt(0, false)
	)

	for index, thinking := range ThinkingForJobInfo {
		if thinking == label {
			thinkingNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return thinkingNumber
}

/***
 * AMBIマスタ変換
 */
// AMBIの職種マスタをautoscoutの職種マスタに変換
func ConvertAmbiOccupationInt(occupationNullInt null.Int) (int64, int64) {
	if !occupationNullInt.Valid {
		return 0, 0
	}
	if occupationNullInt.Int64 == 9999 || occupationNullInt.Int64 == 1000 {
		return 9999, 0
	}

	occupation := occupationNullInt.Int64

	if occupation == 0 || occupation == 1 || occupation == 2 || occupation == 3 || occupation == 4 {
		return 100, 0
	} else if occupation == 5 {
		return 802, 0
	} else if occupation == 6 || occupation == 7 {
		return 801, 0
	} else if occupation == 8 || occupation == 9 {
		return 802, 0
	} else if occupation == 10 || occupation == 11 || occupation == 12 || occupation == 13 || occupation == 14 {
		return 800, 0
	} else if occupation == 15 {
		return 803, 0
	} else if occupation == 16 {
		return 100, 0
	} else if occupation == 17 {
		return 1002, 0
	} else if occupation == 18 {
		return 1000, 0
	} else if occupation == 19 {
		return 1000, 0
	} else if occupation == 20 {
		return 503, 0
	} else if occupation == 21 {
		return 503, 0
	} else if occupation == 22 {
		return 600, 0
	} else if occupation == 23 {
		return 503, 0
	} else if occupation == 24 {
		return 600, 0
	} else if occupation == 25 {
		return 600, 503
	} else if occupation == 26 {
		return 201, 0
	} else if occupation == 27 {
		return 200, 0
	} else if occupation == 28 {
		return 200, 0
	} else if occupation == 29 {
		return 200, 0
	} else if occupation == 30 {
		return 200, 0
	} else if occupation == 31 {
		return 201, 200
	} else if occupation == 32 {
		return 201, 200

	} else if occupation == 33 {
		return 501, 0
	} else if occupation == 34 {
		return 501, 0
	} else if occupation == 35 {
		return 700, 0
	} else if occupation == 36 {
		return 501, 0
	} else if occupation == 37 {
		return 100, 0
	} else if occupation == 38 {
		return 502, 0
	} else if occupation == 39 {
		return 502, 0
	} else if occupation == 40 {
		return 502, 0

	} else if occupation == 41 {
		return 300, 0
	} else if occupation == 42 {
		return 301, 0
	} else if occupation == 43 {
		return 200, 0
	} else if occupation == 44 {
		return 1202, 0
	} else if occupation == 45 {
		return 901, 0
	} else if occupation == 46 {
		return 900, 0
	} else if occupation == 47 {
		return 200, 0

	} else if occupation >= 48 && occupation <= 61 {
		return 1203, 0
	} else if occupation == 62 || occupation == 63 || occupation == 64 || occupation == 65 {
		return 200, 201
	} else if occupation == 66 || occupation == 67 {
		return 404, 0
	} else if occupation == 68 {
		return 200, 201
	} else if occupation == 69 {
		return 100, 0
	} else if occupation == 70 || occupation == 71 {
		return 200, 0
	} else if occupation >= 72 && occupation <= 77 {
		return 1101, 0
	} else if occupation == 78 || occupation == 79 {
		return 1102, 0
	} else if occupation == 80 {
		return 200, 0
	} else if occupation == 81 {
		return 1104, 0
	} else if occupation == 82 {
		return 100, 0
	} else if occupation == 83 {
		return 1201, 0
	} else if occupation == 84 {
		return 1100, 0
	} else if occupation == 85 {
		return 1103, 0
	} else if occupation == 86 {
		return 1101, 1102

	} else if occupation == 87 {
		return 1105, 1107
	} else if occupation == 88 {
		return 1105, 1107
	} else if occupation == 89 {
		return 1105, 1107
	} else if occupation == 90 {
		return 1105, 1107
	} else if occupation == 91 {
		return 1105, 1107
	} else if occupation == 92 {
		return 1108, 0
	} else if occupation == 93 {
		return 1109, 0
	} else if occupation == 94 {
		return 200, 0
	} else if occupation == 95 {
		return 1105, 1107
	} else if occupation == 96 {
		return 1105, 1107
	} else if occupation == 97 {
		return 1105, 1107

	} else if occupation == 98 {
		return 1105, 1107
	} else if occupation == 99 {
		return 1105, 1107
	} else if occupation == 100 {
		return 1108, 0
	} else if occupation == 101 {
		return 1109, 0
	} else if occupation == 102 {
		return 200, 0
	} else if occupation == 103 {
		return 1105, 1107
	} else if occupation == 104 {
		return 1105, 1107
	} else if occupation == 105 {
		return 1105, 1107
	} else if occupation == 106 {
		return 1105, 1106
	} else if occupation == 107 {
		return 1105, 1106
	} else if occupation == 108 {
		return 1108, 0
	} else if occupation == 109 {
		return 1109, 0
	} else if occupation == 110 {
		return 200, 0
	} else if occupation == 111 {
		return 1105, 1107
	} else if occupation == 112 {
		return 1105, 1107
	} else if occupation == 113 {
		return 1105, 1107

	} else if occupation == 114 {
		return 1110, 0
	} else if occupation == 115 {
		return 1110, 0
	} else if occupation == 116 {
		return 1110, 0
	} else if occupation == 117 {
		return 1110, 0
	} else if occupation == 118 {
		return 1110, 0
	} else if occupation == 119 {
		return 1110, 0
	} else if occupation == 120 {
		return 1110, 0
	} else if occupation == 121 {
		return 1110, 0
	} else if occupation == 122 {
		return 1110, 0

	} else if occupation == 123 {
		return 1106, 1107
	} else if occupation == 124 {
		return 1106, 1107
	} else if occupation == 125 {
		return 1106, 1107
	} else if occupation == 126 {
		return 1106, 1107
	} else if occupation == 127 {
		return 1106, 1107
	} else if occupation == 128 {
		return 1106, 0
	} else if occupation == 129 {
		return 1108, 0
	} else if occupation == 130 {
		return 1109, 0
	} else if occupation == 131 {
		return 200, 0
	} else if occupation == 132 {
		return 1106, 1107
	} else if occupation == 133 {
		return 1106, 1107
	} else if occupation == 134 {
		return 1106, 1107
	} else if occupation == 135 {
		return 1106, 1107
	} else if occupation == 136 {
		return 408, 0
	} else if occupation == 137 {
		return 408, 0
	} else if occupation == 138 {
		return 406, 0
	} else if occupation == 139 {
		return 1204, 0
	} else if occupation == 140 {
		return 501, 0
	} else if occupation == 141 {
		return 1301, 0
	} else if occupation == 142 {
		return 501, 0
	} else if occupation == 143 {
		return 406, 0
	} else if occupation == 144 {
		return 1000, 0
	} else if occupation == 145 {
		return 0, 0
	} else if occupation == 146 {
		return 700, 0
	} else if occupation == 147 {
		return 701, 0
	} else if occupation == 148 {
		return 702, 0
	} else if occupation == 149 {
		return 502, 0
	} else if occupation == 150 {
		return 701, 0
	} else if occupation == 151 {
		return 703, 0

	} else if occupation == 152 {
		return 1400, 0
	} else if occupation == 153 {
		return 1400, 0
	} else if occupation == 154 {
		return 1400, 0
	}

	return 0, 0
}

var AmbiOccupation = map[string]map[int]string{
	"経営・経営企画・事業企画系": {
		0: "経営者・COO・経営幹部・カントリーヘッド",
		1: "経営企画・事業企画",
		2: "M＆A",
		3: "新規事業",
		4: "その他・経営・経営企画・事業企画系",
	},
	"管理部門系": {
		5:  "総務",
		6:  "人事（採用・労務・教育など）",
		7:  "人事制度・企画",
		8:  "法務・コンプライアンス",
		9:  "特許・知的財産関連",
		10: "CFO",
		11: "経理",
		12: "財務・コントローラー",
		13: "内部監査",
		14: "会計・税務",
		15: "広報・IR",
		16: "管理部長",
		17: "秘書・セクレタリー・アシスタント",
		18: "一般事務・営業事務",
		19: "その他・管理部門系",
	},
	"SCM・ロジスティクス・物流・購買・貿易系": {
		20: "購買・調達",
		21: "SCM",
		22: "物流企画・ロジスティクス",
		23: "貿易・通関",
		24: "センター・倉庫管理・運行・配車管理",
		25: "その他・SCM・ロジスティクス・物流・貿易系",
	},
	"営業系": {
		26: "営業（個人営業）",
		27: "営業（法人営業）",
		28: "海外営業",
		29: "営業マネージャー・管理職",
		30: "MR（医療情報担当者）・MS（医療品卸販売担当者）",
		31: "人材コンサルタント・コーディネーター",
		32: "その他・営業系",
	},
	"マーケティング・販売企画・商品開発系": {
		33: "商品企画・開発",
		34: "マーケティング・販売企画",
		35: "マーケティングプランナー・Webプランナー",
		36: "営業企画",
		37: "ブランド・プロダクトマネージャー",
		38: "Web・デジタルマーケティング",
		39: "マーケティンリサーチ・分析",
		40: "その他・マーケティング系",
	},

	"コンサルタント系（戦略、財務、組織、その他専門）": {
		41: "戦略コンサルタント",
		42: "財務・会計コンサルタント",
		43: "組織・人事コンサルタント",
		44: "調査員・リサーチャー",
		45: "弁護士・弁理士",
		46: "会計士・税理士",
		47: "その他・コンサルタント系",
	},
	"金融系専門職": {
		48: "法人営業（金融）",
		49: "個人営業（金融）・FP",
		50: "代理店営業・ホールセラー",
		51: "投資研究・アナリスト・エコノミスト",
		52: "ファンドマネージャー・ディーラー・トレーダー",
		53: "インベストメントバンキング・M＆A",
		54: "コーポレートファイナンス",
		55: "リスク管理・与信管理・債務管理",
		56: "コンプライアンス・監査",
		57: "金融事務・決済・計理・主計",
		58: "アンダーライター・損害調査",
		59: "金融商品企画・ストラクチャード",
		60: "アクチュアリー",
		61: "その他・金融系",
	},

	"不動産系専門職": {
		62: "不動産企画・仕入・開発",
		63: "アセットマネジメント・ヘッジファンド・PE投資",
		64: "プロパティマネジメント",
		65: "不動産鑑定評価（デューデリジェンス）",
		66: "ファシリティマネジメント・設備管理",
		67: "フロント・マンション管理",
		68: "その他・不動産系専門職",
	},

	"技術系（IT・WEB・通信系）": {
		69: "CTO・CIO",
		70: "ITコンサルタント",
		71: "ビジネスアナリスト・アーキテクト",
		72: "PM（WEB・オープン）",
		73: "PM（汎用系）",
		74: "PM（パッケージ・ミドルウェア系）",
		75: "SE（WEB・オープン）",
		76: "SE（汎用系）",
		77: "SE（パッケージ・ミドルウェア系）",
		78: "サーバー・ネットワークエンジニア",
		79: "DBエンジニア",
		80: "プリセールス・セールスエンジニア",
		81: "社内SE・システム管理",
		82: "プロダクトマネージャー",
		83: "データサイエンティスト",
		84: "製品開発・研究",
		85: "テクニカルサポート",
		86: "その他・技術系（IT・WEB・通信系）",
	},

	"技術系（電気・電子・半導体）": {
		87: "設計・開発エンジニア（電気・電子回路）",
		88: "設計・開発エンジニア（半導体）",
		89: "設計・開発エンジニア（その他・電気・電子・半導体）",
		90: "アプリケーション開発エンジニア（制御・組み込み系）",
		91: "PM（制御・組み込み系）",
		92: "生産技術・製造技術・エンジニアリング（電気・電子）",
		93: "生産管理・品質管理・品質保証・工場長（電気・電子）",
		94: "セールスエンジニア（電気・電子）",
		95: "サポートエンジニア（電気・電子）",
		96: "PM（電気・電子）",
		97: "その他・技術系（電気・電子・半導体）",
	},

	"技術系（機械・メカトロ・自動車）": {
		98:  "設計・開発エンジニア（自動車・輸送機器）",
		99:  "設計・開発エンジニア（その他・機械・メカトロ・自動車）",
		100: "生産技術・製造技術・エンジニアリング（機械・自動車）",
		101: "生産管理・品質管理・品質保証・工場長（機械・自動車）",
		102: "セールスエンジニア（機械・自動車）",
		103: "サポートエンジニア（機械・自動車）",
		104: "PM（機械・自動車）",
		105: "その他・技術系（その他・機械・メカトロ・自動車）",
	},

	"技術系（化学・素材・食品・衣料）": {
		106: "研究・開発（化学・素材・食品・衣料）",
		107: "研究・開発（その他・化学・素材・食品・衣料）",
		108: "生産技術・製造技術・エンジニアリング（化学・素材・食品・衣料）",
		109: "生産管理・品質管理・品質保証・工場長（化学・素材・食品・衣料）",
		110: "セールスエンジニア（化学・素材・食品・衣料）",
		111: "サポートエンジニア（化学・素材・食品・衣料）",
		112: "PM（化学・素材・食品・衣料）",
		113: "その他・技術系（その他・化学・素材・食品・衣料）",
	},

	"技術系（建設・設備・土木・プラント）": {
		114: "設計（建築）",
		115: "設計（設備）",
		116: "設計（土木）",
		117: "施工管理（建築）",
		118: "施工管理（設備）",
		119: "施工管理（土木）",
		120: "プラントエンジニアリング",
		121: "建築・土木技術開発・建設コンサルタント",
		122: "その他・技術系（建築・設備・土木・プラント）",
	},

	"技術系（メディカル）": {
		123: "研究・開発（医薬品）",
		124: "研究・開発（医療用具・医療機器）",
		125: "研究・開発（その他・メディカル）",
		126: "臨床開発・治験",
		127: "薬事",
		128: "学術",
		129: "生産技術・製造技術・エンジニアリング（メディカル）",
		130: "生産管理・品質管理・品質保証・工場長（メディカル）",
		131: "セールスエンジニア（メディカル）",
		132: "サポートエンジニア（メディカル）",
		133: "PM（メディカル）",
		134: "CRA・CRC",
		135: "その他・技術系（メディカル）",
	},

	"接客・販売・サービス・流通系": {
		136: "SV",
		137: "バイヤー、マーチャンダイザー、VMD",
		138: "支配人・ホテルフロント",
		139: "施設長・事務長・その他介護福祉系職",
		140: "店舗開発・FC開発",
		141: "講師・教師・インストラクター",
		142: "コールセンター運営・管理",
		143: "料理長・シェフ・調理師・メニュー開発",
		144: "通訳・翻訳",
		145: "その他・サービス・流通系",
	},

	"クリエイティブ系": {
		146: "プランナー",
		147: "プロデューサー・ディレクター",
		148: "デザイナー",
		149: "Webサイト運営・コンテンツ企画",
		150: "編集・コピーライター",
		151: "その他・クリエイティブ",
	},

	"総合職": {
		152: "総合職（文系）",
		153: "総合職（理系）",
	},

	"その他（オープンポジション）": {
		154: "オープンポジション",
	},

	"まだ定まっていない": {
		1000: "まだ定まっていない",
	},
}

func ConvertAmbiIndustryInt(industryNullInt null.Int) int64 {
	if !industryNullInt.Valid {
		return 9999
	}

	if industryNullInt == null.NewInt(0, true) {
		return 100
	} else if industryNullInt == null.NewInt(1, true) {
		return 101
	} else if industryNullInt == null.NewInt(2, true) {
		return 704
	} else if industryNullInt == null.NewInt(3, true) {
		return 706
	} else if industryNullInt == null.NewInt(4, true) {
		return 707
	} else if industryNullInt == null.NewInt(5, true) {
		return 100
	} else if industryNullInt == null.NewInt(6, true) {
		return 203
	} else if industryNullInt == null.NewInt(7, true) {
		return 203
	} else if industryNullInt == null.NewInt(8, true) {
		return 202
	} else if industryNullInt == null.NewInt(9, true) {
		return 202
	} else if industryNullInt == null.NewInt(10, true) {
		return 200
	} else if industryNullInt == null.NewInt(11, true) {
		return 206
	} else if industryNullInt == null.NewInt(12, true) {
		return 204
	} else if industryNullInt == null.NewInt(13, true) {
		return 207
	} else if industryNullInt == null.NewInt(14, true) {
		return 209
	} else if industryNullInt == null.NewInt(15, true) {
		return 1703
	} else if industryNullInt == null.NewInt(16, true) {
		return 300
	} else if industryNullInt == null.NewInt(17, true) {
		return 400
	} else if industryNullInt == null.NewInt(18, true) {
		return 401
	} else if industryNullInt == null.NewInt(19, true) {
		return 408
	} else if industryNullInt == null.NewInt(20, true) {
		return 404
	} else if industryNullInt == null.NewInt(21, true) {
		return 406
	} else if industryNullInt == null.NewInt(22, true) {
		return 400
	} else if industryNullInt == null.NewInt(23, true) {
		return 407
	} else if industryNullInt == null.NewInt(24, true) {
		return 409
	} else if industryNullInt == null.NewInt(25, true) {
		return 1703
	} else if industryNullInt == null.NewInt(26, true) {
		return 600
	} else if industryNullInt == null.NewInt(27, true) {
		return 600
	} else if industryNullInt == null.NewInt(28, true) {
		return 600
	} else if industryNullInt == null.NewInt(29, true) {
		return 600
	} else if industryNullInt == null.NewInt(30, true) {
		return 600
	} else if industryNullInt == null.NewInt(31, true) {
		return 601
	} else if industryNullInt == null.NewInt(32, true) {
		return 1302
	} else if industryNullInt == null.NewInt(33, true) {
		return 1303
	} else if industryNullInt == null.NewInt(34, true) {
		return 500
	} else if industryNullInt == null.NewInt(35, true) {
		return 900
	} else if industryNullInt == null.NewInt(36, true) {
		return 1300
	} else if industryNullInt == null.NewInt(37, true) {
		return 1703
	} else if industryNullInt == null.NewInt(38, true) {
		return 700
	} else if industryNullInt == null.NewInt(39, true) {
		return 800
	} else if industryNullInt == null.NewInt(40, true) {
		return 1100
	} else if industryNullInt == null.NewInt(41, true) {
		return 1101
	} else if industryNullInt == null.NewInt(42, true) {
		return 1100
	} else if industryNullInt == null.NewInt(43, true) {
		return 1100
	} else if industryNullInt == null.NewInt(44, true) {
		return 1400
	} else if industryNullInt == null.NewInt(45, true) {
		return 1401
	} else if industryNullInt == null.NewInt(46, true) {
		return 1200
	} else if industryNullInt == null.NewInt(47, true) {
		return 1201
	} else if industryNullInt == null.NewInt(48, true) {
		return 1500
	} else if industryNullInt == null.NewInt(49, true) {
		return 1501
	} else if industryNullInt == null.NewInt(50, true) {
		return 1600
	} else if industryNullInt == null.NewInt(51, true) {
		return 1700
	} else if industryNullInt == null.NewInt(52, true) {
		return 1702
	} else if industryNullInt == null.NewInt(53, true) {
		return 1703
	}
	return 9999
}

var AmbiIndustry map[string]map[uint]string = map[string]map[uint]string{
	"IT・インターネット・ゲーム": {
		0: "IT",
		1: "通信キャリア",
		2: "インターネット広告・メディア",
		3: "WEB制作・WEBデザイン",
		4: "ゲーム",
		5: "IT・インターネット・ゲーム（その他）",
	},
	"メーカー": {
		6:  "メーカー（コンピューター・通信系）",
		7:  "メーカー（電気・電子・半導体）",
		8:  "メーカー（自動車・輸送機器）",
		9:  "メーカー（機械）",
		10: "メーカー（化学・素材）",
		11: "メーカー（食品）",
		12: "メーカー（医薬品・医療機器）",
		13: "メーカー（ファッション・アパレル）",
		14: "メーカー（日用品・化粧品）",
		15: "メーカー（その他）",
	},
	"商社": {
		16: "商社（コンピューター・通信）",
		17: "商社（電気・電子・半導体）",
		18: "商社（自動車・輸送機器）",
		19: "商社（機械）",
		20: "商社（化学・素材）",
		21: "商社（食品）",
		22: "商社（医薬品・医療機器）",
		23: "商社（ファッション・アパレル）",
		24: "商社（日用品・化粧品）",
		25: "商社（その他）",
	},
	"流通・小売・サービス": {
		26: "流通・小売（百貨店・スーパー・コンビニ）",
		27: "流通・小売（ファッション・アパレル）",
		28: "流通・小売（医薬品・化粧品）",
		29: "流通・小売（食品）",
		30: "流通・小売（家電）",
		31: "通信販売",
		32: "フード・レストラン",
		33: "レジャー・アミューズメント",
		34: "人材ビジネス",
		35: "コールセンター",
		36: "ホテル・観光",
		37: "流通・小売・サービス（その他）",
	},
	"広告・出版・マスコミ": {
		38: "放送・広告・印刷・出版",
	},
	"コンサルティング": {
		39: "コンサルティングファーム・シンクタンク",
	},
	"金融": {
		40: "金融（銀行）",
		41: "金融（保険）",
		42: "金融（証券）",
		43: "金融（その他）",
	},
	"建築・不動産": {
		44: "不動産",
		45: "建築・土木",
	},
	"メディカル": {
		46: "医療",
		47: "福祉・介護",
	},
	"物流・運輸": {
		48: "物流・倉庫",
		49: "陸運・海運・航空・鉄道",
	},
	"インフラ・教育・官公庁・その他": {
		50: "電気・ガス・水道",
		51: "教育・学校",
		52: "団体・連合会・官公庁",
		53: "その他の業種",
	},
	"まだ定まっていない": {
		9999: "まだ定まっていない",
	},
}

/***
 * マイナビマスタ変換 *途中
 */
func ConvertMynaviOccupation(occupationStr string) (null.Int, null.Int) {
	if occupationStr == "" {
		return null.NewInt(1400, true), null.NewInt(0, false)
	}

	switch occupationStr {
	case "営業・企画営業（法人向け）":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "営業・企画営業（個人向け）":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "営業マネジャー・営業管理職":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "代理店営業・パートナーセールス":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "内勤営業・カウンターセールス":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "ルートセールス・渉外・外商":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "海外営業":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "メディカル営業（MR・MS・その他）":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "キャリアカウンセラー":
		return null.NewInt(201, true), null.NewInt(0, false)
	case "人材コーディネーター":
		return null.NewInt(201, true), null.NewInt(0, false)
	case "コールセンター管理・運営":
		return null.NewInt(201, true), null.NewInt(0, false)
	case "カスタマーサポート・サポートデスク":
		return null.NewInt(201, true), null.NewInt(0, false)
	case "テクニカルサポート":
		return null.NewInt(201, true), null.NewInt(0, false)
	case "オペレーター・アポインター":
		return null.NewInt(201, true), null.NewInt(0, false)
	case "マーケティングリサーチ・分析":
	case "販促企画・営業企画":
	case "商品企画・商品開発":
	case "広告宣伝":
	case "マーチャンダイザー":
	case "仕入れ・バイヤー":
	case "店舗開発・FC開発":
	case "事業企画・事業プロデュース":
	case "海外事業企画":
	case "経営企画":
	case "CIO・CTO":
	case "経営幹部（CEO・CFO・COO等）":
	case "FCオーナー・代理店研修生":
	case "購買・資材調達":
	case "貿易業務・国際業務":
	case "物流企画・物流管理":
	case "商品管理・在庫管理":
	case "経理・財務":
	case "会計・税務":
	case "総務":
	case "人事・労務・採用":
	case "法務・コンプライアンス":
	case "知的財産・特許":
	case "広報":
	case "IR":
	case "内部監査":
	case "情報セキュリティ":
	case "一般事務・庶務":
	case "営業事務・営業アシスタント":
	case "受付":
	case "秘書":
	case "スーパーバイザー・エリアマネジャー":
	case "教育・研修トレーナー":
	case "店長・店長候補（小売・流通系）":
	case "販売・販売アドバイザー・売り場担当":
	case "美容部員":
	case "店長・店長候補（フード・アミューズメント系）":
	case "ホール・フロアスタッフ":
	case "調理・調理補助":
	case "エステティシャン":
	case "理容師・美容師":
	case "アロマセラピスト・ネイリスト":
	case "トリマー":
	case "ブライダルコーディネーター・ウェディングプランナー":
	case "葬祭ディレクター・プランナー":
	case "旅行コーディネーター・添乗":
	case "カウンタースタッフ・予約手配・オペレーター":
	case "ホテル・宿泊施設サービス":
	case "タクシードライバー・ハイヤードライバー":
	case "バス運転手・バス乗務員":
	case "パイロット・航空管制官等空輸職":
	case "フライトアテンダント（CA・FA）・グランドスタッフ":
	case "鉄道乗務員・船舶乗務員":
	case "管理薬剤師・薬剤師・登録販売者":
	case "医師":
	case "看護師・准看護師・保健師・助産師":
	case "歯科医師":
	case "歯科技工士・歯科衛生士":
	case "PT（理学療法士）・OT（作業療法士）・ST（言語聴覚士）・ORT（視能訓練士）":
	case "マッサージ師・柔道整復師・鍼師・灸師":
	case "各種検査技師":
	case "臨床心理士・カウンセラー・セラピスト":
	case "医療事務・医療秘書":
	case "獣医":
	case "介護・福祉事業責任者・施設長":
	case "ケアマネジャー（介護支援専門員）":
	case "サービス提供責任者":
	case "介護職・ヘルパー":
	case "生活相談員・生活支援員":
	case "介護系事務職":
	case "管理栄養士・栄養士・フードコーディネーター":
	case "保育士・幼稚園教諭":
	case "児童相談員":
	case "スクール運営・マネジメント":
	case "教師":
	case "講師":
	case "スポーツインストラクター・トレーナー":
	case "インストラクター(OA・その他)":
	case "教務事務":
	case "通訳":
	case "翻訳":
	case "コンサルタント（経営戦略）":
	case "コンサルタント（財務・会計）":
	case "コンサルタント（業務プロセス）":
	case "コンサルタント（組織・人事）":
	case "コンサルタント（生産・物流）":
	case "コンサルタント（営業・マーケティング）":
	case "ISOコンサルタント・ISO審査員":
	case "公開業務(IPO)":
	case "M&A":
	case "研究調査員・リサーチャー":
	case "公認会計士":
	case "税理士":
	case "弁護士":
	case "弁理士・特許技術者":
	case "司法書士・行政書士":
	case "社会保険労務士":
	case "士業補助者":
	case "金融営業（法人）":
	case "金融営業（個人）・リテール・FP":
	case "金融営業（代理店）・パートナーセールス":
	case "投資銀行業務（インベストバンキング）":
	case "運用業務・ファンドマネジャー":
	case "トレーダー・ディーラー":
	case "アナリスト・エコノミスト":
	case "ストラテジックファイナンス":
	case "金融系専門職（ミドル・バック・その他）":
	case "金融商品開発・アクチュアリー":
	case "投資理論・クオンツ":
	case "金融システム企画":
	case "リスク・与信・債権管理":
	case "金融事務":
	case "生損保系専門職（査定・損害調査等）":
	case "不動産営業":
	case "アセットマネジャー":
	case "不動産鑑定・デューデリジェンス":
	case "プロパティマネジャー":
	case "ファシリティマネジャー":
	case "用地仕入":
	case "不動産事業企画":
	case "不動産管理":
	case "フロント（マンション管理・ビル管理）":
	case "アカウントエグゼクティブ（AE）・アカウントプランナー":
	case "メディアプランナー（MP)":
	case "クリエイティブディレクター":
	case "コピーライター":
	case "アートディレクター":
	case "グラフィックデザイナー・CGデザイナー・イラストレーター(広告系)":
	case "フォトグラファー":
	case "ディレクター・プロデューサー・進行管理（編集・制作系）":
	case "編集・校正":
	case "記者・ライター":
	case "テクニカルライター":
	case "DTPオペレーター":
	case "ファッションデザイナー・服飾雑貨デザイナー・テキスタイルデザイナー":
	case "パタンナー・縫製":
	case "生産管理・品質管理（アパレル・ファッション）":
	case "スタイリスト・ヘアメイク":
	case "プロダクトデザイナー（工業デザイン）":
	case "生産管理・品質管理（工業プロダクト）":
	case "インテリアデザイナー・インテリアコーディネーター":
	case "空間・ディスプレイ・店舗デザイナー":
	case "プロデューサー・ディレクター・プランナー・演出":
	case "脚本家・放送作家":
	case "アシスタントプロデューサー・アシスタントディレクター・進行":
	case "アナウンサー・イベントコンパニオン・モデル・俳優":
	case "芸能マネジャー":
	case "グラフィックデザイナー・CGデザイナー・イラストレーター（映像系）":
	case "制作関連技術者（カメラ・照明・音響）":
	case "SEOコンサルタント・SEMコンサルタント":
	case "インターネットサービス企画":
	case "WEBプロデューサー・ディレクター":
	case "情報アーキテクト・UI/UXデザイナー":
	case "システムディレクター・テクニカルディレクター":
	case "アクセス解析・統計解析":
	case "WEBコンテンツ企画・制作":
	case "WEBデザイナー":
	case "フロントエンドエンジニア・コーダー":
	case "プログラマー(WEBサイト・インターネットサービス系)":
	case "ディレクター・プロデューサー（ゲーム・アミューズメント系)":
	case "ゲームプランナー・ゲーム企画":
	case "シナリオライター":
	case "グラフィックデザイナー・CGデザイナー・イラストレーター(ゲーム・アミューズメント系)":
	case "プログラマー(ゲーム・アミューズメント系)":
	case "サウンドクリエイター":
	case "WEBショップ・ECサイト運営":
	case "システムアナリスト":
	case "ITアーキテクト":
	case "システムコンサルタント（業務系）":
	case "システムコンサルタント（DB・ミドルウェア）":
	case "システムコンサルタント（ネットワーク・通信）":
	case "パッケージ導入コンサルタント（ERP・SCM・CRM等）":
	case "セキュリティコンサルタント":
	case "プリセールス・セールスエンジニア":
	case "プロジェクトマネジャー・リーダー(WEB・オープン・モバイル系）":
	case "システムエンジニア（アプリ設計／WEB・オープン・モバイル系）":
	case "システムエンジニア（DB・ミドルウェア設計／WEB・オープン・モバイル系）":
	case "プログラマー（WEB・オープン・モバイル系）":
	case "プロジェクトマネジャー・リーダー（汎用機系）":
	case "システムエンジニア（アプリ設計／汎用機系）":
	case "システムエンジニア（DB・ミドルウェア設計／汎用機系）":
	case "プログラマー（汎用機系）":
	case "プロジェクトマネジャー(制御系)":
	case "システムエンジニア(通信制御ソフト開発／制御系)":
	case "システムエンジニア(マイコン・計測・画像等／制御系)":
	case "プログラマー(制御系)":
	case "プロダクトマネジャー（パッケージソフト・ミドルウェア）":
	case "システムエンジニア（パッケージソフト・ミドルウェア）":
	case "プログラマー（パッケージソフト・ミドルウェア）":
	case "ローカライズ （パッケージソフト・ミドルウェア）":
	case "ネットワーク設計・構築（LAN・WAN・インターネット）":
	case "サーバ設計・構築（LAN・WAN・インターネット）":
	case "通信設備計画策定":
	case "通信設備設計・構築（有線系）":
	case "通信設備設計・構築（無線系）":
	case "通信設備設置・テスト":
	case "サーバ・マシン運用・監視":
	case "ネットワーク運用・監視":
	case "テクニカルサポート（ソフトウェア・ネットワーク）":
	case "導入・運用トレーナー":
	case "社内情報化戦略・推進":
	case "社内システム開発・運用":
	case "研究開発（ソフトウェア・ネットワーク）":
	case "特許技術者（ソフトウェア・ネットワーク）":
	case "品質管理（ソフトウェア・ネットワーク）":
	case "基礎研究（電気・電子・機械・半導体・材料系）":
	case "応用研究（電気・電子・機械・半導体・材料系）":
	case "特許技術者（電気・電子・機械・半導体・材料系）":
	case "システム設計・アーキテクチャー":
	case "デジタル回路設計":
	case "アナログ回路設計":
	case "高周波回路設計":
	case "混載回路設計":
	case "システムLSI設計":
	case "デジタルIC設計":
	case "アナログIC設計":
	case "高周波IC設計":
	case "混載IC設計":
	case "パワーIC設計":
	case "光学設計・その他光学関連職":
	case "制御設計（家電・コンピューター・通信機器系）":
	case "制御設計（精密・医療用機器系）":
	case "制御設計（自動車・輸送用機器系）":
	case "制御設計（工作機械・ロボット系）":
	case "制御設計（その他）":
	case "機械・機構設計（家電・コンピューター・通信機器系）":
	case "機械・機構設計（精密・医療用機器系）":
	case "機械・機構設計（自動車・輸送用機器系）":
	case "機械・機構設計（工作機械・ロボット・機械系）":
	case "機械・機構設計（その他）":
	case "金型設計":
	case "解析":
	case "生産・製造・プロセス技術（家電・コンピューター・通信機器系）":
	case "生産・製造・プロセス技術（精密・医療用機器系）":
	case "生産・製造・プロセス技術（自動車・輸送用機器系）":
	case "生産・製造・プロセス技術（工作機械・ロボット系）":
	case "生産・製造・プロセス技術（半導体・電子部品系）":
	case "生産・製造・プロセス技術（機械部品・金型・治工具系）":
	case "品質保証（電気・電子・機械・半導体・材料系）":
	case "品質管理（電気・電子・機械・半導体・材料系）":
	case "生産管理・製造管理（電気・電子・機械・半導体・材料系）":
	case "セールスエンジニア・FAE（家電・コンピューター・通信機器系）":
	case "セールスエンジニア・FAE（精密・医療用機器系）":
	case "セールスエンジニア・FAE（自動車・輸送用機器系）":
	case "セールスエンジニア・FAE（工作機械・ロボット・機械系）":
	case "サービスエンジニア・サポートエンジニア（家電・コンピューター・通信機器系）":
	case "サービスエンジニア・サポートエンジニア（精密・医療用機器系）":
	case "サービスエンジニア・サポートエンジニア（自動車・輸送用機器系）":
	case "サービスエンジニア・サポートエンジニア（工作機械・ロボット・機械系）":
	case "CAD・CAMオペレーター（電気・電子・機械・半導体系）":
	case "評価・検査（家電・コンピューター・通信機器系）":
	case "評価・検査（精密・医療用機器系）":
	case "評価・検査（自動車・輸送用機器系）":
	case "評価・検査（工作機械・ロボット・機械系）":
	case "建設コンサルタント":
	case "測量":
	case "建築設計":
	case "土木設計":
	case "プラント設計":
	case "電気設備設計":
	case "空調設備設計":
	case "その他設計・設備設計":
	case "CADオペレーター（建築）・製図":
	case "積算":
	case "構造解析（建築・土木）":
	case "建築施工管理・工事監理者":
	case "内装施工管理・工事監理者（リフォーム・住宅・商業施設等）":
	case "土木施工管理・工事監理者":
	case "造園施工管理・工事監理者":
	case "プラント施工管理・工事監理者":
	case "電気設備施工管理・工事監理者":
	case "空調設備施工管理・工事監理者":
	case "管工事施工管理・工事監理者":
	case "設備保全・メンテナンス":
	case "環境保全・管理・調査・分析":
	case "製品・研究開発（建築・土木・プラント・設備）":
	case "生産技術・生産管理（建築・土木・プラント・設備）":
	case "品質管理・保証（建築・土木・プラント・設備）":
	case "特許技術・調査":
	case "基礎・応用研究・技術開発":
	case "生産技術・製造技術":
	case "生産管理":
	case "設備管理":
	case "品質管理・保証(化学・素材・バイオ系)":
	case "フィールドエンジニア":
	case "セールスエンジニア(化学・素材・バイオ系)":
	case "基礎研究":
	case "商品開発":
	case "生産技術・生産管理・製造技術（食品・化粧品系）":
	case "品質管理・保証(食品・化粧品系)":
	case "申請関連":
	case "研究(基礎・シーズ探索・スクリーニング)":
	case "研究(ゲノム・バイオ)":
	case "前臨床研究":
	case "臨床開発モニター(CRA)":
	case "治験コーディネーター（CRC）":
	case "臨床開発":
	case "薬事申請（医薬品・医療機器・医薬部外品・化粧品系）":
	case "生産技術・生産管理・製造技術（医薬品・医療機器系）":
	case "品質管理・保証(医薬品・医療機器系)":
	case "データマネジメント・生物統計":
	case "学術":
	case "整備・メカニック（自動車・二輪車）":
	case "工場生産・製造（輸送用機器・家電・電子機器系）":
	case "工場生産・製造（食品・化粧品・医薬品系）":
	case "工場生産・製造（アパレル・ファッション・その他製品）":
	case "土木・建築・解体工事（とび工・鉄筋工等）":
	case "外装・内装工事（塗装工・防水工等）":
	case "設備工事（電気・通信）":
	case "警備・守衛・清掃":
	case "マンション・ビル管理者":
	case "設備管理・保守（ガス・空調・上下水・消防等）":
	case "ビル施設管理":
	case "配送・宅配・セールスドライバー":
	case "運送ドライバー（中・長距離）":
	case "新聞配達・集金":
	case "倉庫作業・管理":
	case "農林水産関連（農業）":
	case "農林水産関連（林業）":
	case "農林水産関連（水産業）":
	case "農林水産関連（畜産業・その他）":
	case "飼育員（ブリーダー・調教師等）":
	case "公務員（事務系）":
	case "公務員（技術系）":
	case "警察官":
	case "消防士":
	case "自衛隊":
	case "団体職員":
	case "学校法人職員":
	case "公共施設職員（図書館・美術館等）":
	}

	return null.NewInt(1400, true), null.NewInt(0, false)
}

func ConvertMynaviIndustry(industryStr string) (null.Int, null.Int) {
	if industryStr == "" {
		return null.NewInt(1703, true), null.NewInt(0, false)
	}

	switch industryStr {
	case "ソフトウェア・情報処理":
		return null.NewInt(0, true), null.NewInt(0, false)
	case "インターネット関連":
	case "通信関連":
	case "ゲーム関連":
	case "総合電機":
	case "コンピューター機器":
	case "家電・AV機器":
	case "その他電気・電子関連":
	case "ゲーム・アミューズメント製品":
	case "精密機器":
	case "通信機器":
	case "半導体・電子・電気機器":
	case "医療用機器・医療関連":
	case "重電・産業用電気機器":
	case "輸送用機器（自動車含む）":
	case "プラント・エンジニアリング":
	case "鉱業・金属製品・鉄鋼":
	case "ガラス・化学・石油":
	case "繊維":
	case "紙・パルプ":
	case "セメント":
	case "ゴム":
	case "窯業・セラミック":
	case "非鉄金属":
	case "住宅・建材・エクステリア":
	case "インテリア・住宅関連":
	case "食品":
	case "化粧品・医薬品":
	case "スポーツ・レジャー用品（メーカー）":
	case "その他メーカー":
	case "文具・事務機器関連":
	case "繊維・アパレル":
	case "宝飾品・貴金属":
	case "日用品・雑貨":
	case "玩具":
	case "総合商社":
	case "専門商社":
	case "人材派遣・人材紹介":
	case "エステティック・美容・理容":
	case "セキュリティ":
	case "ビル管理・メンテナンス":
	case "アウトソーシング":
	case "医療・福祉・介護サービス":
	case "サービス（その他）":
	case "教育":
	case "冠婚葬祭":
	case "フィットネスクラブ":
	case "旅行・観光":
	case "ホテル・旅館":
	case "レジャーサービス・アミューズメント":
	case "通信販売・ネット販売":
	case "流通・チェーンストア":
	case "ホームセンター":
	case "百貨店":
	case "専門店（総合）":
	case "専門店（食品関連）":
	case "専門店（自動車関連）":
	case "専門店（その他小売）":
	case "専門店（カメラ・OA関連）":
	case "専門店（電気機器関連）":
	case "専門店（書籍・音楽関連）":
	case "コンビニエンスストア":
	case "ドラッグストア・調剤薬局":
	case "専門店（メガネ・貴金属）":
	case "専門店（ファッション・服飾関連）":
	case "専門店（インテリア関連）":
	case "専門店（スポーツ用品）":
	case "フードビジネス（総合）":
	case "フードビジネス（洋食）":
	case "フードビジネス（ファストフード）":
	case "フードビジネス（アジア系）":
	case "フードビジネス（和食）":
	case "放送・映像・音響":
	case "新聞・出版・印刷":
	case "ディスプレイ・空間デザイン・イベント":
	case "広告":
	case "アート・芸能関連":
	case "外資系銀行":
	case "信用組合・信用金庫・労働金庫":
	case "銀行":
	case "信託銀行":
	case "証券・投資銀行":
	case "商品取引":
	case "投資信託委託・投資顧問":
	case "生命保険・損害保険":
	case "リース・レンタル":
	case "クレジット・信販":
	case "その他金融":
	case "政府系・系統金融機関":
	case "共済":
	case "外資系金融":
	case "金融総合グループ":
	case "ベンチャーキャピタル":
	case "事業者金融・消費者金融":
	case "シンクタンク・マーケティング・調査":
	case "専門コンサルタント":
	case "個人事務所（士業）":
	case "設備工事":
	case "設計":
	case "建設コンサルタント":
	case "建設・土木":
	case "不動産":
	case "リフォーム・内装工事":
	case "物流・倉庫":
	case "海運・鉄道・空輸・陸運":
	case "環境関連設備":
	case "環境・リサイクル":
	case "電力・ガス・エネルギー":
	case "警察・消防・自衛隊":
	case "官公庁":
	case "公益・特殊・独立行政法人":
	case "農林・水産":
	case "生活協同組合":
	case "農業協同組合（JA金融機関を含む）":
	}

	return null.NewInt(1703, true), null.NewInt(0, false)
}

/******************************************************************/
// サーカスエージェントマスタ変換(https://docs.google.com/spreadsheets/d/1y5L5Wj3dSb3R7rYLDyKavxOheyaNS4Ev9qAPB57jEDo/edit#gid=1957537079)
//
func GetIntPublicOfferingForCircus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	switch label {
	case "未上場":
		return null.NewInt(3, true)
	case "東証一部":
		return null.NewInt(0, true)
	case "東証二部":
		return null.NewInt(1, true)
	case "マザーズ":
		return null.NewInt(2, true)
	case "ジャスダック":
		return null.NewInt(1, true)
	case "名証一部":
		return null.NewInt(5, true)
	case "名証二部":
		return null.NewInt(5, true)
	case "セントレックス":
		return null.NewInt(5, true)
	case "本則市場":
		return null.NewInt(5, true)
	case "アンビシャス":
		return null.NewInt(5, true)
	case "Q-Board":
		return null.NewInt(5, true)
	case "その他株式市場上場":
		return null.NewInt(5, true)
	case "東証プライム":
		return null.NewInt(0, true)
	case "東証スタンダード":
		return null.NewInt(1, true)
	case "東証グロース":
		return null.NewInt(2, true)
	}

	for masterI, masterL := range PublicOffering {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntGenderForCircus(label string) null.Int {
	var (
		genderNumber = null.NewInt(0, false)
	)

	if label == "不問" || label == "" || label == "性別不問" {
		return null.NewInt(99, true)
	} else if label == "男性限定" {
		return null.NewInt(0, true)
	} else if label == "女性限定" {
		return null.NewInt(1, true)
	} else if label == "男性であれば尚良し" {
		return null.NewInt(2, true)
	} else if label == "女性であれば尚良し" {
		return null.NewInt(3, true)
	}

	for index, gender := range GenderForJobInfo {
		if gender == label {
			genderNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return genderNumber
}
func GetIntNationalityForCircus(label string) null.Int {
	var (
		nationalityNumber = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	} else if label == "外国籍NG" {
		return null.NewInt(0, true)
	} else if label == "外国籍OK" {
		return null.NewInt(99, true)
	}

	for index, nationality := range NationalityForJobInfo {
		if nationality == label {
			nationalityNumber = null.NewInt(int64(index), true)
			break
		}
	}

	return nationalityNumber
}

func GetIntEmploymentStatusForCircus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "[正社員]" {
		return null.NewInt(0, true)
	} else if label == "[契約社員]" {
		return null.NewInt(1, true)
	} else if label == "[派遣社員]" {
		return null.NewInt(2, true)
	} else if label == "[紹介予定派遣]" {
		return null.NewInt(3, true)
	} else if label == "[アルバイト・パート]" {
		return null.NewInt(4, true)
	}

	for masterI, masterL := range EmploymentStatusForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntOpenOrCloseForCircus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "募集中" {
		return null.NewInt(0, true)
	}

	for masterI, masterL := range OpenOrClose {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

/*
{ value: 0, label: '中卒以上' },
{ value: 1, label: '高卒以上' },
{ value: 2, label: '専卒以上（短大除く）' },
{ value: 3, label: '短卒以上（専卒除く）' },
{ value: 4, label: '専卒・短卒以上' },
{ value: 5, label: '高専卒以上' },
{ value: 6, label: '大卒以上' },
{ value: 7, label: '院卒以上' },
*/

// 高卒以上, 専門卒以上, 大卒以上
func GetIntFinalEducationForCircus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" || label == "" || label == "学歴不問" {
		return null.NewInt(0, false)
	} else if label == "専門卒以上" || label == "短大卒以上" {
		return null.NewInt(4, true)
	}

	for masterI, masterL := range FinalEducationForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntHolidayForCircus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(99, true)
	}

	switch label {
	case "土日祝休み":
		return null.NewInt(0, true)
	case "土日休み":
		return null.NewInt(0, true)
	case "週休2日（土日以外）":
		return null.NewInt(1, true)
	case "シフト制":
		return null.NewInt(3, true)
	case "その他":
		return null.NewInt(99, true)
	}

	for masterI, masterL := range HolidayForJobSeeker {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

func GetIntPrefectureForCircus(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "全国各地" || label == "全国" || label == "不問" {
		return null.NewInt(9999, true)
	}

	// 県や府が抜けていても、ヒットするようにする
	for index, value := range Prefecture {
		if strings.Contains(value, label) {
			return null.NewInt(int64(index), true)
		}
	}

	return number
}

// 必要条件
func GetIntRequiredConditionForCircus(label string) null.Int {
	// 最終学歴
	outputNumber := GetIntFinalEducationForCircus(label)
	if outputNumber.Valid {
		return outputNumber
	}
	// 休日
	outputNumber = GetIntHolidayForCircus(label)
	if outputNumber.Valid {
		return outputNumber
	}
	// 国籍
	outputNumber = GetIntNationalityForCircus(label)
	if outputNumber.Valid {
		return outputNumber
	}
	// 性別
	outputNumber = GetIntGenderForCircus(label)
	if outputNumber.Valid {
		return outputNumber
	}

	return null.NewInt(0, false)
}

func GetIntIndustryForCircus(label string) (null.Int, null.Int) {
	if label == "その他" {
		return null.NewInt(1703, true), null.NewInt(0, false)
	}

	// 例外
	switch label {
	case "その他":
		return null.NewInt(1703, true), null.NewInt(0, false)

		// IT
	case "ソフトウェア・情報処理":
		return null.NewInt(100, true), null.NewInt(0, false)
	case "web・インターネット・ゲーム":
		return null.NewInt(101, true), null.NewInt(0, false)
	case "通信":
		return null.NewInt(100, true), null.NewInt(0, false)

		// メーカー
	case "電気機器メーカー":
		return null.NewInt(203, true), null.NewInt(0, false)
	case "精密機器・医療機器メーカー":
		return null.NewInt(203, true), null.NewInt(0, false)
	case "半導体、電気・電子部品メーカー":
		return null.NewInt(203, true), null.NewInt(0, false)
	case "コンピュータ・通信機器メーカー":
		return null.NewInt(203, true), null.NewInt(204, true)
	case "自動車・輸送用機器メーカー":
		return null.NewInt(202, true), null.NewInt(0, false)
	case "化学、素材メーカー":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "医薬品、バイオ、食品メーカー":
		return null.NewInt(204, true), null.NewInt(206, true)
	case "その他メーカー":
		return null.NewInt(1703, true), null.NewInt(0, false)

		// 金融・保険
	case "銀行":
		return null.NewInt(1100, true), null.NewInt(0, false)
	case "生命保険・損害保険":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "証券・投資関連":
		return null.NewInt(1100, true), null.NewInt(0, false)
	case "金融・保険（その他）":
		return null.NewInt(1100, true), null.NewInt(0, false)

	case "総合商社":
		return null.NewInt(400, true), null.NewInt(0, false)
	case "専門商社":
		return null.NewInt(1703, true), null.NewInt(0, false)
	case "商社（その他）":
		return null.NewInt(1703, true), null.NewInt(0, false)

		// 小売
	case "ファッション・アパレル":
		return null.NewInt(600, true), null.NewInt(601, true)
	case "美容・コスメ":
		return null.NewInt(600, true), null.NewInt(601, true)
	case "飲食・レストラン":
		return null.NewInt(1302, true), null.NewInt(0, false)
	case "アミューズメント・ホテル":
		return null.NewInt(1300, true), null.NewInt(0, false)
	case "教育":
		return null.NewInt(502, true), null.NewInt(0, false)
	case "医療・福祉":
		return null.NewInt(1200, true), null.NewInt(0, false)
	case "冠婚葬祭":
		return null.NewInt(1305, true), null.NewInt(0, false)
	case "その他サービス":
		return null.NewInt(1703, true), null.NewInt(0, false)

	case "人材ビジネス":
		return null.NewInt(500, true), null.NewInt(0, false)
	// case "その他":

	case "物流・倉庫":
		return null.NewInt(1500, true), null.NewInt(0, false)
	case "運送":
		return null.NewInt(1501, true), null.NewInt(0, false)
	// case "その他":

	case "コンサルティングファーム":
		return null.NewInt(800, true), null.NewInt(0, false)
	case "監査法人・税理士法人":
		return null.NewInt(800, true), null.NewInt(1000, true)
	case "シンクタンク":
		return null.NewInt(800, true), null.NewInt(0, false)
	case "専門コンサルティング":
		return null.NewInt(800, true), null.NewInt(0, false)
	// case "その他":

	case "広告":
		return null.NewInt(704, true), null.NewInt(0, false)
	case "PR・イベント":
		return null.NewInt(702, true), null.NewInt(0, false)
	case "マスコミ":
		return null.NewInt(700, true), null.NewInt(0, false)
	case "出版・印刷":
		return null.NewInt(703, true), null.NewInt(0, false)
	case "映像・音量":
		return null.NewInt(701, true), null.NewInt(0, false)
	// case "その他":

	case "不動産":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "建築":
		return null.NewInt(1401, true), null.NewInt(0, false)
	case "土木":
		return null.NewInt(1401, true), null.NewInt(0, false)
	// case "その他":

	case "エネルギー（電気・ガス）":
		return null.NewInt(1600, true), null.NewInt(0, false)
	case "官公庁・団体":
		return null.NewInt(1702, true), null.NewInt(0, false)
	case "医療機関":
		return null.NewInt(1200, true), null.NewInt(0, false)

		// 大分類:エネルギー関連・公共関連すべて
	// case "その他":

	case "就業経験なし":
		// return null.NewInt(1703, true), null.NewInt(0, false)

	}

	return null.NewInt(0, false), null.NewInt(0, false)
}

func GetIntOccupationForCircus(label string) (null.Int, null.Int) {
	if label == "その他" {
		return null.NewInt(1400, true), null.NewInt(0, false)
	}

	// 例外
	switch label {
	case "法人営業":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "個人営業":
		return null.NewInt(201, true), null.NewInt(0, false)
	case "ルートセールス":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "内勤営業・カウンターセールス":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "キャリアカウンセラー・人材コーディネーター":
		return null.NewInt(200, true), null.NewInt(201, true)
	case "MR":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "技術営業":
		return null.NewInt(200, true), null.NewInt(0, false)
	case "営業企画":
		return null.NewInt(501, true), null.NewInt(0, false)
	case "その他営業":
		return null.NewInt(200, true), null.NewInt(201, true)
	case "マーケティング":
		return null.NewInt(502, true), null.NewInt(0, false)
	case "商品企画・商品開発":
		return null.NewInt(501, true), null.NewInt(0, false)
	case "バイヤー・マーチャンダイザー・物流":
		return null.NewInt(503, true), null.NewInt(0, false)
	case "リサーチャー・データサイエンティスト":
		return null.NewInt(1201, true), null.NewInt(1202, false)
	case "その他企画・マーケティング関連職":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "経営者・CxO":
		return null.NewInt(100, true), null.NewInt(0, false)
	case "経営企画・事業企画":
		return null.NewInt(500, true), null.NewInt(0, false)
	case "経理":
		return null.NewInt(800, true), null.NewInt(0, false)
	case "財務・会計・監査":
		return null.NewInt(800, true), null.NewInt(0, false)
	case "人事（採用・教育・評価）":
		return null.NewInt(801, true), null.NewInt(0, false)
	case "人事（労務管理）":
		return null.NewInt(801, true), null.NewInt(0, false)
	case "総務":
		return null.NewInt(801, true), null.NewInt(0, false)
	case "法務・特許":
		return null.NewInt(802, true), null.NewInt(0, false)
	case "広報・IR":
		return null.NewInt(803, true), null.NewInt(0, false)
	case "その他管理部門関連職":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "一般事務":
		return null.NewInt(1000, true), null.NewInt(0, false)
	case "営業事務（営業アシスタント）":
		return null.NewInt(1000, true), null.NewInt(0, false)
	case "専門事務（医療・金融・不動産・貿易）":
		return null.NewInt(1000, true), null.NewInt(0, false)
	case "秘書":
		return null.NewInt(1002, true), null.NewInt(0, false)
	case "受付・コンシェルジュ":
		return null.NewInt(1001, true), null.NewInt(0, false)
	case "コールセンターSV":
		return null.NewInt(402, true), null.NewInt(403, true)
	case "カスタマーサポート・ヘルプデスク":
		return null.NewInt(402, true), null.NewInt(0, false)
	case "テレフォンオペレーター・アポインター":
		return null.NewInt(403, true), null.NewInt(0, false)
	case "販売・サービススタッフ":
		return null.NewInt(408, true), null.NewInt(401, true)
	case "スーパーバイザー、エリアマネージャー、店長":
		return null.NewInt(408, true), null.NewInt(0, false)
	case "飲食（ホール・調理スタッフ）":
		return null.NewInt(400, true), null.NewInt(401, true)
	case "美容関連接客・販売":
		return null.NewInt(407, true), null.NewInt(0, false)
	case "アパレル関連接客販売":
		return null.NewInt(408, true), null.NewInt(0, false)
	case "旅行・ホテル・ブライダル関連接客":
		return null.NewInt(406, true), null.NewInt(0, false)
	case "アミューズメント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "介護・福祉系接客":
		return null.NewInt(1204, true), null.NewInt(0, false)
	case "警備員・清掃":
		return null.NewInt(404, true), null.NewInt(0, false)
	case "その他接客・販売・サービス関連職":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "インストラクター・講師":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "翻訳・通訳":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "ドライバー":
		return null.NewInt(405, true), null.NewInt(0, false)
	case "web・インターネット関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "Webデザイナー":
		return null.NewInt(702, true), null.NewInt(0, false)
	case "Webディレクター・Webプロデューサー":
		return null.NewInt(701, true), null.NewInt(0, false)
	case "ゲームクリエイター・ゲームプランナー":
		return null.NewInt(703, true), null.NewInt(0, false)
	case "編集・ライター":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "制作進行管理":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "ファッションデザイナー・パタンナー":
		return null.NewInt(702, true), null.NewInt(0, false)
	case "グラフィックデザイナー・イラストレーター":
		return null.NewInt(702, true), null.NewInt(0, false)
	case "映像、音響、イベント関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "店舗・空間デザイナー":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "その他クリエイティブ関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "システムコンサルタント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "SAP・ERP導入コンサルタント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "セキュリティエンジニア・コンサルタント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "プリセールス・セールスエンジニア":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "PG（オープンweb系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "SE（オープンweb系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "PL・PM（オープンweb系）":
		return null.NewInt(1101, true), null.NewInt(1100, true)
	case "アプリ開発（オープンweb系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "PG（汎用機系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "SE（汎用機系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "PL・PM（汎用機系）":
		return null.NewInt(1101, true), null.NewInt(1100, true)
	case "PG（制御・組込系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "SE（制御・組込系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "PL・PM（制御・組込系）":
		return null.NewInt(1101, true), null.NewInt(1100, true)
	case "アプリ開発（制御・組込系）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "PG（パッケージソフト開発）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "SE（パッケージソフト開発）":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "PL・PM（パッケージソフト開発）":
		return null.NewInt(1101, true), null.NewInt(1100, true)
	case "サーバー設計・構築":
		return null.NewInt(1102, true), null.NewInt(1100, true)
	case "ネットワーク・構築":
		return null.NewInt(1102, true), null.NewInt(1100, true)
	case "データベース設計・構築":
		return null.NewInt(1102, true), null.NewInt(1100, true)
	case "インフラ保守・運用":
		return null.NewInt(1102, true), null.NewInt(1100, true)
	case "テクニカルサポート":
		return null.NewInt(1103, true), null.NewInt(0, false)
	case "オペレーター":
		return null.NewInt(1103, true), null.NewInt(0, false)
	case "システム導入・運用トレーナー":
		return null.NewInt(1103, true), null.NewInt(0, false)
	case "社内SE":
		return null.NewInt(1104, true), null.NewInt(0, false)
	case "テスター・テストエンジニア":
		return null.NewInt(1103, true), null.NewInt(0, false)
	case "フロントエンド・マークアップエンジニア":
		return null.NewInt(1101, true), null.NewInt(0, false)
	case "その他エンジニア関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
		// 技術系
	case "設計開発（電気・電子・半導体）":
		return null.NewInt(1105, true), null.NewInt(0, false)
	case "生産管理・品質管理（電気・電子・半導体）":
		return null.NewInt(1109, true), null.NewInt(0, false)
	case "セールスエンジニア・サポートエンジニア（電気・電子・半導体）":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "PM（電気・電子・半導体）":
		return null.NewInt(1107, true), null.NewInt(1105, true)
	case "研究開発（電気・電子・半導体）":
		return null.NewInt(1107, true), null.NewInt(1106, true)
	case "生産技術（電気・電子・半導体）":
		return null.NewInt(1108, true), null.NewInt(0, false)
	case "その他電気・電子・半導体関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "設計開発（機械・メカトロ・自動車）":
		return null.NewInt(1105, true), null.NewInt(0, false)
	case "生産管理・品質管理（機械・メカトロ・自動車）":
		return null.NewInt(1108, true), null.NewInt(0, false)
	case "セールスエンジニア・サポートエンジニア（機械・メカトロ・自動車）":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "PM（機械・メカトロ・自動車）":
		return null.NewInt(1107, true), null.NewInt(1105, true)
	case "研究開発（機械・メカトロ・自動車）":
		return null.NewInt(1107, true), null.NewInt(1106, true)
	case "生産技術（機械・メカトロ・自動車）":
		return null.NewInt(1108, true), null.NewInt(0, false)
	case "その他機械・メカトロ・自動車関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "研究開発（素材・化学・食品）":
		return null.NewInt(1107, true), null.NewInt(1106, true)
	case "生産管理・品質管理（素材・化学・食品）":
		return null.NewInt(1108, true), null.NewInt(0, false)
	case "セールスエンジニア・サポートエンジニア（素材・化学・食品）":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "PM（素材・化学・食品）":
		return null.NewInt(1107, true), null.NewInt(1105, true)
	case "生産技術（素材・化学・食品）":
		return null.NewInt(1108, true), null.NewInt(0, false)
	case "その他素材・化学・食品関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "設計（建築・設備）":
		return null.NewInt(1110, true), null.NewInt(0, false)
	case "設計（土木）":
		return null.NewInt(1110, true), null.NewInt(0, false)
	case "施工管理（建築・設備）":
		return null.NewInt(1110, true), null.NewInt(0, false)
	case "施工管理（土木）":
		return null.NewInt(1110, true), null.NewInt(0, false)
	case "プラント":
		return null.NewInt(1110, true), null.NewInt(0, false)
	case "その他建築・土木関連":
		return null.NewInt(1110, true), null.NewInt(0, false)
	case "戦略・経営コンサルタント":
		return null.NewInt(300, true), null.NewInt(0, false)
	case "財務・会計コンサルタント":
		return null.NewInt(301, true), null.NewInt(0, false)
	case "組織・人事コンサルタント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "生産・物流コンサルタント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "営業・マーケティングコンサルタント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "その他ビジネスコンサルタント系関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "金融商品開発・アクチュアリー":
		return null.NewInt(1203, true), null.NewInt(0, false)
	case "ディーラー・トレーダー":
		return null.NewInt(1203, true), null.NewInt(0, false)
	case "A・ストラクチャードファイナンス":
		return null.NewInt(1203, true), null.NewInt(0, false)
	case "金融アナリスト・リサーチ":
		return null.NewInt(1203, true), null.NewInt(0, false)
	case "融資審査・査定":
		return null.NewInt(1203, true), null.NewInt(0, false)
	case "その他金融系関連":
		return null.NewInt(1203, true), null.NewInt(0, false)
	case "不動産企画・開発":
		return null.NewInt(200, true), null.NewInt(201, true)
	case "不動産仕入れ":
		return null.NewInt(200, true), null.NewInt(201, true)
	case "不動産デューデリジェンス（査定）":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "プロパティマネジメント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "アセットマネジメント":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "ビル・マンション管理":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "その他不動産系関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "研究開発（医療系）":
		return null.NewInt(1106, true), null.NewInt(1107, true)
	case "生産管理・品質管理（医療系）":
		return null.NewInt(1109, true), null.NewInt(0, false)
	case "セールスエンジニア・サポートエンジニア（医療系）":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "PM（医療系）":
		return null.NewInt(1106, true), null.NewInt(1107, true)
	case "生産技術（医療系）":
		return null.NewInt(1109, true), null.NewInt(0, false)
	case "薬事":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "医療学術":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "その他医療系関連":
		return null.NewInt(1400, true), null.NewInt(0, false)
	case "医師・看護師・薬剤師":
		return null.NewInt(902, true), null.NewInt(903, true)
	case "士業（弁護士・税理士・会計士等）":
		return null.NewInt(900, true), null.NewInt(901, true)
	case "士業（歯科衛生士・栄養士・臨床心理士等）":
		return null.NewInt(902, true), null.NewInt(903, true)
	case "就業経験なし":
		return null.NewInt(1400, true), null.NewInt(0, false)
	}

	return null.NewInt(1400, true), null.NewInt(0, false)
}

/******************************************************************/
// agentBankマスタ変換
//
func GetIntFeatureForAgentBank(label string) null.Int {
	var (
		featureNumber = null.NewInt(0, false)
	)

	// 正社員経験なしOK 3
	// 業界未経験可 0
	// 職種未経験可 1
	// 業界・職種未経験可 2
	// 外国籍可 ×
	// 新卒限定 ×
	// 上場企業 7
	// 車通勤OK ×
	// 転勤なし 6
	// 服装自由 ×
	// 社員寮あり 18
	// 年間休日120日以上 19
	// 土日祝休み ×
	// 残業時間20時間以内 20
	// 時短勤務可 ×

	switch label {
	/** エージェント限定情報 */
	case "正社員経験なしOK":
		featureNumber = null.NewInt(3, true)
	case "業界未経験可":
		featureNumber = null.NewInt(0, true)
	case "職種未経験可":
		featureNumber = null.NewInt(1, true)
	case "業界・職種未経験可":
		featureNumber = null.NewInt(2, true)
	case "上場企業":
		featureNumber = null.NewInt(7, true)
	case "転勤なし":
		featureNumber = null.NewInt(6, true)
	case "社員寮あり":
		featureNumber = null.NewInt(18, true)
	case "年間休日120日以上":
		featureNumber = null.NewInt(19, true)
	case "残業20時間以内":
		featureNumber = null.NewInt(20, true)
	}

	return featureNumber
}

func GetIntEmploymentStatusForAgentBank(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "正社員" {
		return null.NewInt(0, true)
	} else if label == "契約社員" {
		return null.NewInt(1, true)
	} else if label == "パート・アルバイト" {
		return null.NewInt(4, true)
	} else if label == "無期雇用派遣" {
		return null.NewInt(5, true)
	}

	for masterI, masterL := range EmploymentStatusForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// 不問, 高校卒業以上, 高専卒業以上, 短大・専門卒業以上, 大学卒業以上. MARCH以上, 早慶・国公立以上, 大学院卒以上
func GetIntFinalEducationForAgentBank(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)

	if label == "不問" || label == "" {
		return null.NewInt(0, false)
	} else if label == "高校卒業以上" {
		return null.NewInt(1, true)
	} else if label == "高専卒業以上" {
		return null.NewInt(5, true)
	} else if label == "短大・専門卒業以上" {
		return null.NewInt(4, true)
	} else if label == "大学卒業以上" || label == "MARCH以上" || label == "早慶・国公立以上" {
		return null.NewInt(6, true)
	} else if label == "大学院卒以上" {
		return null.NewInt(7, true)
	}
	for masterI, masterL := range FinalEducationForJobInfo {
		if masterL == label {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

/*
agentbank: 1社まで可, 2社まで可, 3社まで可, 4社まで可, 5社まで可, 6社まで可, 不問
autoscout:

	0: '0回のみ',
	1: '1回まで',
	2: '2回まで',
	3: '3回まで',
	4: '4回まで',
	5: '5回まで',
	99: '不問',
*/
func GetIntJobChangeForAgentBank(label string) null.Int {
	if label == "不問" {
		return null.NewInt(99, true)
	}

	if strings.Contains(label, "社まで可") {
		labenInt, _ := strconv.Atoi(string([]rune(label)[0]))
		masterValue := labenInt - 1
		return null.NewInt(int64(masterValue), true)
	} else {
		return null.NewInt(0, false)
	}
}

func GetIntIndustryForAgentBank(label string) []null.Int {
	var industries = []null.Int{}

	switch label {
	case "ソフトウェア 情報処理":
		industries = append(industries, null.NewInt(100, true))
	case "インターネット関連 ゲーム":
		industries = append(industries, null.NewInt(100, true), null.NewInt(707, true))
	case "通信":
		industries = append(industries, null.NewInt(100, true), null.NewInt(101, true))
	case "メーカー（自動車・輸送機器）":
		industries = append(industries, null.NewInt(202, true))
	case "メーカー（電子・電子部品・半導体）":
		industries = append(industries, null.NewInt(203, true))
	case "メーカー（機械関連）":
		industries = append(industries, null.NewInt(202, true))
	case "メーカー（医療機器）":
		industries = append(industries, null.NewInt(204, true))
	case "メーカー（食料品）":
		industries = append(industries, null.NewInt(206, true))
	case "メーカー（医薬品・化粧品・生活用品）":
		industries = append(industries, null.NewInt(204, true), null.NewInt(209, true))
	case "メーカー（化学・繊維・素材）":
		industries = append(industries, null.NewInt(200, true))
	case "メーカー（ファッション・アパレル）":
		industries = append(industries, null.NewInt(207, true))
	case "メーカー（その他）":
		industries = append(industries, null.NewInt(210, true))
	case "総合商社":
		industries = append(industries, null.NewInt(300, true))
	case "専門商社":
		industries = append(industries, null.NewInt(410, true))
	case "商社（その他）":
		industries = append(industries, null.NewInt(410, true))
	case "物流 運輸 倉庫":
		industries = append(industries, null.NewInt(1500, true), null.NewInt(1501, true))
	case "小売（ファッション アパレル）":
		industries = append(industries, null.NewInt(600, true))
	case "小売（スーパー CVS）":
		industries = append(industries, null.NewInt(600, true))
	case "小売（専門店）":
		industries = append(industries, null.NewInt(600, true))
	case "小売（その他）":
		industries = append(industries, null.NewInt(600, true))
	case "銀行":
		industries = append(industries, null.NewInt(1100, true))
	case "生命保険 損害保険":
		industries = append(industries, null.NewInt(1101, true))
	case "証券 投資関連":
		industries = append(industries, null.NewInt(1102, true))
	case "金融 保険（その他）":
		industries = append(industries, null.NewInt(1100, true))
	case "建築・土木・設計":
		industries = append(industries, null.NewInt(1401, true))
	case "不動産":
		industries = append(industries, null.NewInt(1400, true))
	case "不動産・建設系（その他）":
		industries = append(industries, null.NewInt(1400, true))
	case "人材サービス":
		industries = append(industries, null.NewInt(500, true))
	case "教育関連":
		industries = append(industries, null.NewInt(502, true))
	case "フード・レストラン":
		industries = append(industries, null.NewInt(1302, true))
	case "レジャー・ホテル・旅行":
		industries = append(industries, null.NewInt(1300, true))
	case "エンターテイメント・スポーツ":
		industries = append(industries, null.NewInt(1303, true))
	case "医療":
		industries = append(industries, null.NewInt(1200, true))
	case "福祉・保育":
		industries = append(industries, null.NewInt(1201, true))
	case "冠婚葬祭":
		industries = append(industries, null.NewInt(1305, true))
	case "サービス（その他）":
		industries = append(industries, null.NewInt(1307, true))
	case "シンクタンク・コンサルティングファーム":
		industries = append(industries, null.NewInt(800, true))
	case "監査法人・税理士事務所・法律事務所":
		industries = append(industries, null.NewInt(1000, true))
	case "専門コンサル（その他）":
		industries = append(industries, null.NewInt(800, true))
	case "広告":
		industries = append(industries, null.NewInt(703, true), null.NewInt(704, true))
	case "イベント PR":
		industries = append(industries, null.NewInt(702, true))
	case "放送 映像 音響":
		industries = append(industries, null.NewInt(700, true))
	case "印刷 出版":
		industries = append(industries, null.NewInt(700, true))
	case "マスコミ（その他）":
		industries = append(industries, null.NewInt(700, true))
	case "鉄道・航空":
		industries = append(industries, null.NewInt(1501, true))
	case "団体・連合会・官公庁":
		industries = append(industries, null.NewInt(1702, true))
	case "電気・ガス・水道・エネルギー・環境関連":
		industries = append(industries, null.NewInt(1600, true))
	case "その他の業":
		industries = append(industries, null.NewInt(1703, true))
	}

	return industries
}

func GetIntOccupationForAgentBank(label string) []null.Int {
	var occupations = []null.Int{}

	switch label {
	case "法人営業":
		occupations = append(occupations, null.NewInt(200, true))
	case "個人営業":
		occupations = append(occupations, null.NewInt(201, true))
	case "ルートセールス・代理店営業":
		occupations = append(occupations, null.NewInt(200, true))
	case "内勤営業・カウンターセールス":
		occupations = append(occupations, null.NewInt(204, true))
	case "海外営業":
		occupations = append(occupations, null.NewInt(203, true))
	case "カスタマーサポート・コールセンター運営":
		occupations = append(occupations, null.NewInt(402, true))
	case "キャリアカウンセラー・人材コーディネーター":
		occupations = append(occupations, null.NewInt(201, true))
	case "財務・経理":
		occupations = append(occupations, null.NewInt(800, true))
	case "人事":
		occupations = append(occupations, null.NewInt(801, true))
	case "労務":
		occupations = append(occupations, null.NewInt(801, true))
	case "総務・事務":
		occupations = append(occupations, null.NewInt(801, true))
	case "法務":
		occupations = append(occupations, null.NewInt(802, true))
	case "広報・IR":
		occupations = append(occupations, null.NewInt(803, true))
	case "物流・貿易":
		occupations = append(occupations, null.NewInt(600, true))
	case "一般事務・営業事務":
		occupations = append(occupations, null.NewInt(1000, true))
	case "秘書":
		occupations = append(occupations, null.NewInt(1002, true))
	case "商品企画・商品開発":
		occupations = append(occupations, null.NewInt(501, true))
	case "ブランドマネージャー・プロダクトマネージャー":
		occupations = append(occupations, null.NewInt(505, true))
	case "広告・宣伝":
		occupations = append(occupations, null.NewInt(803, true))
	case "販売促進・販促企画":
		occupations = append(occupations, null.NewInt(501, true))
	case "営業企画":
		occupations = append(occupations, null.NewInt(501, true))
	case "イベント企画・運営":
		occupations = append(occupations, null.NewInt(501, true))
	case "Web・SNSマーケティング":
		occupations = append(occupations, null.NewInt(502, true))
	case "データアナリスト":
		occupations = append(occupations, null.NewInt(1200, true))
	case "市場調査・分析":
		occupations = append(occupations, null.NewInt(1202, true))
	case "経営企画・事業統括":
		occupations = append(occupations, null.NewInt(500, true))
	case "管理職・エグゼクティブ":
		occupations = append(occupations, null.NewInt(100, true))
	case "MD・バイヤー・店舗開発":
		occupations = append(occupations, null.NewInt(408, true))
	case "Webディレクター・Webプロデューサー":
		occupations = append(occupations, null.NewInt(701, true))
	case "テクニカルディレクター・プロジェクトマネージャー":
		occupations = append(occupations, null.NewInt(504, true), null.NewInt(701, true))
	case "クリエイティブディレクター":
		occupations = append(occupations, null.NewInt(701, true), null.NewInt(703, true))
	case "制作・進行管理（その他）":
		occupations = append(occupations, null.NewInt(705, true))
	case "Webデザイナー":
		occupations = append(occupations, null.NewInt(702, true))
	case "UI・UXデザイナー":
		occupations = append(occupations, null.NewInt(702, true))
	case "ゲームデザイナー・イラストレーター":
		occupations = append(occupations, null.NewInt(702, true))
	case "CGデザイナー":
		occupations = append(occupations, null.NewInt(702, true))
	case "Web・モバイル・ソーシャル・ゲーム制作／開発":
		occupations = append(occupations, null.NewInt(702, true), null.NewInt(1101, true))
	case "編集・ライター":
		occupations = append(occupations, null.NewInt(704, true))
	case "映像・動画関連":
		occupations = append(occupations, null.NewInt(703, true))
	case "ファッション・インテリア・空間・プロダクトデザイン":
		occupations = append(occupations, null.NewInt(702, true))
	case "業務系アプリケーションエンジニア・プログラマ":
		occupations = append(occupations, null.NewInt(1101, true))
	case "Webサービス系エンジニア・プログラマ":
		occupations = append(occupations, null.NewInt(1101, true))
	case "制御系ソフトウェア開発（通信・ネットワーク・IoT関連）":
		occupations = append(occupations, null.NewInt(1101, true))
	case "インフラエンジニア":
		occupations = append(occupations, null.NewInt(1102, true))
	case "ITヘルプデスク・カスタマーサポート":
		occupations = append(occupations, null.NewInt(1103, true), null.NewInt(402, true))
	case "IT・システムコンサルタント":
		occupations = append(occupations, null.NewInt(302, true))
	case "社内情報システム（社内SE）":
		occupations = append(occupations, null.NewInt(1104, true))
	case "機械・機構設計・金型設計":
		occupations = append(occupations, null.NewInt(1105, true))
	case "回路・システム設計":
		occupations = append(occupations, null.NewInt(1105, true))
	case "サービスエンジニア・サポートエンジニア":
		occupations = append(occupations, null.NewInt(1103, true))
	case "素材・半導体素材・化成品関連":
		occupations = append(occupations, null.NewInt(1112, true))
	case "化粧品・食品・香料関連":
		occupations = append(occupations, null.NewInt(1112, true))
	case "医薬品関連":
		occupations = append(occupations, null.NewInt(1112, true))
	case "医療用具関連":
		occupations = append(occupations, null.NewInt(1112, true))
	case "研究開発・技術開発・構造解析・特許":
		occupations = append(occupations, null.NewInt(1107, true))
	case "施工管理・設備・環境保全":
		occupations = append(occupations, null.NewInt(1110, true))
	case "プランニング・測量・設計・積算":
		occupations = append(occupations, null.NewInt(1105, true))
	case "技能工（整備・工場生産・製造・工事）":
		occupations = append(occupations, null.NewInt(1111, true), null.NewInt(1108, true))
	case "生産・品質管理":
		occupations = append(occupations, null.NewInt(1109, true))
	case "運輸・配送・倉庫関連":
		occupations = append(occupations, null.NewInt(405, true))
	case "交通（鉄道・バス・タクシー）関連":
		occupations = append(occupations, null.NewInt(405, true))
	case "店長・SV（スーパーバイザー）":
		occupations = append(occupations, null.NewInt(408, true))
	case "ホールスタッフ":
		occupations = append(occupations, null.NewInt(406, true))
	case "料理長":
		occupations = append(occupations, null.NewInt(406, true))
	case "調理":
		occupations = append(occupations, null.NewInt(406, true))
	case "警備・施設管理関連職":
		occupations = append(occupations, null.NewInt(404, true))
	case "販売・サービススタッフ":
		occupations = append(occupations, null.NewInt(401, true))
	case "宿泊施設・ホテル":
		occupations = append(occupations, null.NewInt(406, true))
	case "ビジネスコンサルタント・シンクタンク":
		occupations = append(occupations, null.NewInt(300, true), null.NewInt(303, true))
	case "士業・専門コンサルタント":
		occupations = append(occupations, null.NewInt(901, true), null.NewInt(303, true))
	case "金融系専門職":
		occupations = append(occupations, null.NewInt(1203, true))
	case "不動産・プロパティマネジメント系専門職":
		occupations = append(occupations, null.NewInt(1208, true))
	case "医療・看護":
		occupations = append(occupations, null.NewInt(902, true))
	case "薬事":
		occupations = append(occupations, null.NewInt(903, true))
	case "臨床開発":
		occupations = append(occupations, null.NewInt(1206, true))
	case "福祉・介護":
		occupations = append(occupations, null.NewInt(1204, true))
	case "教育・保育":
		occupations = append(occupations, null.NewInt(1207, true))
	case "インストラクター・通訳・翻訳":
		occupations = append(occupations, null.NewInt(1209, true))
	case "その他":
		occupations = append(occupations, null.NewInt(1400, true))
	}

	return occupations
}

// 内定率
func GetIntOfferRateForAgentBank(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)
	label = strings.ReplaceAll(label, " ", "")

	for masterI, masterL := range OfferRate {
		if strings.Contains(label, masterL) {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// 書類通過率
func GetIntDocumentPassingRateForAgentBank(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)
	label = strings.ReplaceAll(label, " ", "")

	for masterI, masterL := range DocumentPassingRate {
		if strings.Contains(label, masterL) {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}

// 直近の応募数
func GetIntNumberOfRecentApplicationsForAgentBank(label string) null.Int {
	var (
		number = null.NewInt(0, false)
	)
	label = strings.ReplaceAll(label, " ", "")

	for masterI, masterL := range NumberOfRecentApplications {
		if strings.Contains(label, masterL) {
			number = null.NewInt(int64(masterI), true)
			break
		}
	}

	return number
}
