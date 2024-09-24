package interactor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

// null.NewInt(0, false)
var NullInt = null.NewInt(0, false)

// 第一引数のスライスに第二引数の値が含まれているかを確認する関数
func includeUINT(uintSlice []uint, targetV uint) bool {
	for _, v := range uintSlice {
		if v == targetV {
			return true
		}
	}
	return false
}

// 乱数生成.
func MakeRandomStr(strRange uint32) (string, error) {
	//　乱数の範囲
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, strRange)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// I が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil

}

// トークンの暗号化
func encrypt(token string) (cipherText string, err error) {
	tokenBytes := []byte(token)

	encryptionText := os.Getenv("ENCRYPTION_KEY")
	encryptionKey := []byte(encryptionText)

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		return "", err
	}

	encryptedToken := gcm.Seal(nonce, nonce, tokenBytes, nil)

	return base64.StdEncoding.EncodeToString(encryptedToken), nil
}

// トークンの復号化
func decryption(encryptedToken string) (decodeText string, err error) {
	encryptionText := os.Getenv("ENCRYPTION_KEY")
	encryptionKey := []byte(encryptionText)

	// 1. 復号化するためのAES鍵を生成する
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 2. 暗号化されたトークンをデコードする
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 3. 暗号化されたトークンを復号化する
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
		return "", err
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	decryptedTokenBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(decryptedTokenBytes), nil
}

func hashPassword(password string) (hashedPassword string, err error) {
	// パスワードをハッシュ化
	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println("Hashed Password:", string(hashedPasswordByte))

	return string(hashedPasswordByte), nil
}

func compareHashedPaasowd(hashedPassword, password string) (err error) {
	// ハッシュ化されたパスワードと元のパスワードを比較
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func timeStrFormat(date string) string {
	if date == "" {
		return ""
	}

	resDate, _ := time.Parse("2006-01-02T15:04", date)
	year, month, day := resDate.Date()
	h := resDate.Hour()
	m := resDate.Minute()

	hours := fmt.Sprintf("%d", h)
	minutes := fmt.Sprintf("%d", m)

	if h < 10 {
		hours = fmt.Sprintf("0%d", h)
	}
	if m < 10 {
		minutes = fmt.Sprintf("0%d", m)
	}

	dateString := fmt.Sprintf("%d年%d月%d日 %s:%s", year, month, day, hours, minutes)
	return dateString
}

type ExperienceOccupationDuplicatedChecker struct {
	departmentID              uint
	experienceOccupationMonth uint
}

// 経験職種の重複を削除
func removeDuplicateExperienceOccupation(experienceOccupationList []*ExperienceOccupationDuplicatedChecker) []*ExperienceOccupationDuplicatedChecker {
	var result []*ExperienceOccupationDuplicatedChecker
	m := make(map[uint]bool)

	for _, ele := range experienceOccupationList {
		if !m[ele.departmentID] {
			m[ele.departmentID] = true
			result = append(result, ele)
		}
	}
	return result
}

// ********* マッチング検索用 汎用関数　*********
// 誕生日から年齢を取得
func getAgeFromBirthday(birthdayStr string) uint {
	if birthdayStr == "" {
		return 0
	}

	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// 現在時刻
	now := time.Now()

	thisYear, thisMonth, day := now.Date()

	age := thisYear - birthday.Year()

	// 誕生日を迎えていない場合はageを「−1」する
	if thisMonth < birthday.Month() && day < birthday.Day() {
		age -= 1
	}

	return uint(age)
}

// "2020-01"　を　2020, 1に変換
func getIntFromDateStr(dateStr string) (year, month int64) {
	if dateStr == "" {
		return 0, 0
	}

	yearAndMonth, err := time.Parse("2006-01", dateStr)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	year = int64(yearAndMonth.Year())
	month = int64(yearAndMonth.Month())

	return year, month
}

// 入社年月と退社年月から勤続年数を取得して,月換算する
func getNumberOfMonths(startYearAndMonthStr, endYearAndMonthStr string) int64 {
	var (
		endYear    int64
		endMonth   int64
		startYear  int64
		startMonth int64
	)

	// 開始年月が空の場合は0を返す
	if startYearAndMonthStr == "" {
		return 0
	} else {
		startYear, startMonth = getIntFromDateStr(startYearAndMonthStr)
	}

	if endYearAndMonthStr == "" {
		// 退社年月が空の場合は現在時刻を取得
		endYearInt, endMonthTypedMonth, _ := time.Now().Date()
		endYear = int64(endYearInt)
		endMonth = int64(endMonthTypedMonth)
	} else {
		endYear, endMonth = getIntFromDateStr(endYearAndMonthStr)
	}

	// 月数を計算
	numberOfMonths := (endYear - startYear) * 12
	numberOfMonths += endMonth - startMonth

	return numberOfMonths
}

// マッチング検索用
// "10名未満",
// "10名以上100名未満",
// "100名以上500名未満",
// "500名以上1000名未満",
// "1000名以上",
// I=99 "不問",
func getCompanyScaleFromEmployeeNumber(employeeNumber int64) int64 {
	if employeeNumber >= 0 && employeeNumber < 10 {
		return 0
	} else if employeeNumber >= 10 && employeeNumber < 100 {
		return 1
	} else if employeeNumber >= 100 && employeeNumber < 500 {
		return 2
	} else if employeeNumber >= 500 && employeeNumber < 1000 {
		return 3
	} else if employeeNumber >= 1000 {
		return 4
	}
	return 99
}

//	必要マネジメント経験（1〜5名,6〜10名,11名〜30名,31名〜）
//
// "1〜5名",
// "6〜10名",
// "11名〜30名",
// "31名〜",
func getManagementFromNumber(inputNumber int64) int64 {
	if inputNumber >= 1 && inputNumber <= 5 {
		return 0
	} else if inputNumber >= 6 && inputNumber <= 10 {
		return 1
	} else if inputNumber >= 11 && inputNumber <= 30 {
		return 2
	} else if inputNumber >= 31 {
		return 3
	}
	return 99
}

// ************************** 文字列変換　汎用*********************************** //
// 指定した数字に対応する文字列を取得
// 有りor無し
func getStrAvailable(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.Available {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

func getStrAvailableFromBool(input bool) string {
	if input {
		return "あり"
	}
	return "なし"
}

// 不問or不可
func getStrNotAsk(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.NotAsk {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 不問or条件ありor不可
func getStrConditionOrNot(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.ConditionOrNot {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 利用可能or利用不可
func getStrUsageStatus(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.UsageStatus {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 雇用形態
func getStrEmploymentStatusForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.EmploymentStatus {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

func getStrEmploymentStatusForJobInfo(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.EmploymentStatusForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 業界
func getStrIndustry(inputNumber null.Int) string {
	if inputNumber == null.NewInt(9999, true) {
		return "まだ定まっていない"
	} else if inputNumber.Valid {
		for _, v := range entity.Industry {
			if v[uint(inputNumber.Int64)] != "" {
				return v[uint(inputNumber.Int64)]
			}
		}
	}

	return ""
}

// 職種
func getStrOccupation(inputNumber null.Int) string {
	if inputNumber == null.NewInt(9999, true) {
		return "まだ定まっていない"
	} else if inputNumber.Valid {
		for _, v := range entity.Occupation {
			if v[uint(inputNumber.Int64)] != "" {
				return v[uint(inputNumber.Int64)]
			}
		}
	}

	return ""
}

// 都道府県
func getStrPrefecture(inputNumber null.Int) string {
	if inputNumber == null.NewInt(9999, true) {
		return "まだ定まっていない"
	}
	if inputNumber.Valid {
		for i, v := range entity.Prefecture {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

func getStrPrefectureForJobInfo(inputNumber null.Int) string {
	if inputNumber == null.NewInt(9999, true) {
		return "全国各地"
	}
	if inputNumber.Valid {
		for i, v := range entity.Prefecture {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 性別
func getStrGenderForJobSeeker(inputNumber null.Int) string {
	if inputNumber == null.NewInt(99, true) {
		return "不問"
	} else if inputNumber.Valid {
		for i, v := range entity.GenderForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 性別
func getStrGenderForJobInfo(inputNumber null.Int) string {
	if inputNumber == null.NewInt(99, true) {
		return "不問"
	} else if inputNumber.Valid {
		for i, v := range entity.GenderForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 国籍
func getStrNationalityForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.NationalityForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

func getStrNationalityForJobInfo(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.NationalityForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 休日
func getStrHolidayForJobSeerker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.HolidayForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

func getStrHolidayForJobInfo(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.HolidayForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 転勤
func getStrTransfer(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.TransferForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 企業規模
func getStrCompanyScale(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.CompanyScale {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 株式公開
func getStrPublicOffering(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.PublicOffering {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 役職
func getStrPosition(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.Position {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 企業特徴
func getStrJobFeature(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.JobFeature {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 就業時期
func getStrJoinCompanyPeriod(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.JoinCompanyPeriod {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 学歴の開始ステータス
func getStrFirstStatusForStudentHistory(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.FirstStatusForStudentHistory {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 学歴の終了ステータス
func getStrLastStatusForStudentHistory(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.LastStatusForStudentHistory {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 職歴の開始ステータス
func getStrFirstStatusForWorkHistory(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.FirstStatusForWorkHistory {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 職歴の終了ステータス
func getStrLastStatusForWorkHistory(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.LastStatusForWorkHistory {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 言語
func getStrLanguage(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.LanguageType {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 言語レベル
func getStrLanguageLevel(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.LanguageLevel {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

func getStrEmploymentPeriod(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.EmploymentPeriod {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 受動喫煙
func getStrPassiveSmoking(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.PassiveSmoking {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 指定した数字に対応する文字列を取得
// 就業状態
func getStrStateOfEmployment(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.StateOfEmployment {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 指定した数字に対応する文字列を取得　募集状況
func getStrRecruitmentState(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.RecruitmentState {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 学校
func getStrSchoolCategory(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.SchoolCategory {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 学校
func getStrFinalEducationForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.FinalEducationForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 大学ランク
func getStrCollegeRank(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.CollegeRankForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 理系or文系
func getStrStudyCategoryForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.StudyCategoryForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 理系or文系
func getStrStudyCategoryForJobInfo(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.StudyCategoryForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 印象
func getStrAppearanceForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.AppearanceForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// コミュニケーション
func getStrCommunicationForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.CommunicationForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 思考
func getStrThinkingForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.ThinkingForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 印象
func getStrAppearanceForJobInfo(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.AppearanceForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// コミュニケーション
func getStrCommunicationForJobInfo(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.CommunicationForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 思考
func getStrThinkingForJobInfo(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.ThinkingForJobInfo {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 転職回数
func getStrJobChangeForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.JobChangeForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 最終学歴の状態
func getStrLastStatus(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.LastStatusForStudentHistory {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 職歴詳細
func getStrWorkHistoryStatus(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.WorkHistoryStatus {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 契約フェーズ
func getStrContractPhase(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.ContractPhase {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// ヨミ角度
func getStrAccuracy(inputNumber null.Int) string {
	if inputNumber.Valid {
		for index, value := range entity.Accuracy {
			if inputNumber == null.NewInt(int64(index), true) {
				return value
			}
		}
	}

	return ""
}

// 求職者のフェーズ（エントリー、面談調整中、面談予約、面談実施）
func getStrPhaseForJobSeeker(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.PhaseForJobSeeker {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 中途or 新卒　or 既卒orフリーター
func getStrUserStatus(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.UserStatus {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 開発スキルカテゴリーと指定カテゴリー内のスキルタイプ
func getStrDevelopmentSkillCategoryAndType(inputCategory, inputType null.Int) (outputCategory, outputType string) {

	if inputCategory.Valid {
		for categoryI, categoryV := range entity.DevelopmentCategory {
			if inputCategory == null.NewInt(int64(categoryI), true) {
				outputCategory = categoryV

				for typeI, typeV := range entity.DevelopmentTypeList[categoryI] {
					if inputType == null.NewInt(int64(typeI), true) {
						outputType = typeV
						return outputCategory, outputType
					}
				}
			}
		}
	}

	return outputCategory, outputType
}

// 　HelloWorkのコードに対応する文字列を取得する
// 資格
func getStrLicenseType(inputNumber null.Int) string {
	if inputNumber.Valid {
		return entity.LicenseType[uint(inputNumber.Int64)]
	}

	return ""
}

func getStrWordSkill(input null.Int) string {
	var (
		str string
	)

	if input.Valid {
		for entityI, entityV := range entity.WordSkill {
			if input.Int64 == int64(entityI) {
				str = entityV
			}
		}
	}

	return str
}

func getStrExcelSkill(input null.Int) string {
	var (
		str string
	)

	if input.Valid {
		for entityI, entityV := range entity.ExcelSkill {
			if input.Int64 == int64(entityI) {
				str = entityV
			}
		}
	}

	return str
}

func getStrPowerPointSkill(input null.Int) string {
	var (
		str string
	)

	if input.Valid {
		for entityI, entityV := range entity.PowerPointSkill {
			if input.Int64 == int64(entityI) {
				str = entityV
			}
		}
	}

	return str
}

/*******複数のリストを文字列に変換********/
//
// 業界リスト
func getStrIndustryList(inputNumberList []null.Int) string {
	var industryList []string
	for _, inputNumber := range inputNumberList {
		industryList = append(industryList, getStrIndustry(inputNumber))
	}

	return strings.Join(industryList, ",")
}

// 職種リスト
func getStrOccupationListForJobInfo(inputNumberList []entity.JobInformationOccupation) string {
	var list []string

	for _, inputNumber := range inputNumberList {
		list = append(list, getStrOccupation(inputNumber.Occupation))
	}

	return strings.Join(list, ",")
}

// 仕事の特徴
func getStrFeatureList(inputList []entity.JobInformationFeature) string {
	var featureList []string

	if len(inputList) < 1 {
		return ""
	}

	for _, feature := range inputList {
		if feature.Feature.Valid {
			for i, v := range entity.JobFeature {
				if feature.Feature.Int64 == int64(i) {
					featureList = append(featureList, v)
				}
			}
		}
	}

	return strings.Join(featureList, ",")
}

// 仕事の特徴
func getStrWorkCharmPointList(inputList []entity.JobInformationWorkCharmPoint) string {
	var featureList []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.Title != "" {
			workCharmPoint := inputV.Title + "：" + inputV.Contents
			featureList = append(featureList, workCharmPoint)
		}
	}

	return strings.Join(featureList, "\n")
}

// 求人の勤務地
func getStrPrefectureList(inputList []entity.JobInformationPrefecture) string {
	var list []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.Prefecture.Valid {
			for entityI, entityV := range entity.Prefecture {
				if inputV.Prefecture.Int64 == int64(entityI) {
					list = append(list, entityV)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求人のターゲット
func getStrTargetList(inputList []entity.JobInformationTarget) string {
	var list []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.Target.Valid {
			for entityI, entityV := range entity.UserStatus {
				if inputV.Target.Int64 == int64(entityI) {
					list = append(list, entityV)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求人の雇用形態
func getStrEmploymentStatusListForJobInfo(inputList []entity.JobInformationEmploymentStatus) string {
	var list []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.EmploymentStatus.Valid {
			for entityI, entityV := range entity.EmploymentStatusForJobInfo {
				if inputV.EmploymentStatus.Int64 == int64(entityI) {
					list = append(list, entityV)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の雇用形態
func getStrEmploymentStatusListForJobSeeker(inputList []entity.JobInformationEmploymentStatus) string {
	var list []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.EmploymentStatus.Valid {
			for entityI, entityV := range entity.EmploymentStatus {
				if inputV.EmploymentStatus.Int64 == int64(entityI) {
					list = append(list, entityV)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

func getStrRequiredSocialExperienceList(inputList []entity.JobInformationRequiredSocialExperience) string {
	var list []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.SocialExperienceType.Valid {
			for entityI, entityV := range entity.EmploymentStatusForJobInfo {
				if inputV.SocialExperienceType.Int64 == int64(entityI) {
					list = append(list, entityV)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求人の資格
func getStrLicenseListForJobInfo(inputList []entity.JobInformationRequiredLicense) string {
	var list []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.License.Valid && entity.LicenseType[uint(inputV.License.Int64)] != "" {
			list = append(list, entity.LicenseType[uint(inputV.License.Int64)])
		}
	}

	return strings.Join(list, ",")
}

// 求職者の資格
func getStrLicenseListForJobSeeker(inputList []entity.JobSeekerLicense) string {
	var list []string

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.LicenseType.Valid && entity.LicenseType[uint(inputV.LicenseType.Int64)] != "" {
			list = append(list, entity.LicenseType[uint(inputV.LicenseType.Int64)])
		}
	}

	return strings.Join(list, ",")
}

// 求人のPCスキル
func getStrPCToolListForJobInfo(inputList []entity.JobInformationRequiredPCTool) string {

	var (
		list []string
	)

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.Tool.Valid {
			for categoryI, categoryV := range entity.PCTool {
				if inputV.Tool.Int64 == int64(categoryI) {
					list = append(list, categoryV)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者のPCスキル
func getStrPCToolListForJobSeeker(inputList []entity.JobSeekerPCTool) string {

	var (
		list []string
	)

	if len(inputList) < 1 {
		return ""
	}

	for _, inputV := range inputList {
		if inputV.Tool.Valid {
			for categoryI, categoryV := range entity.PCTool {
				if inputV.Tool.Int64 == int64(categoryI) {
					list = append(list, categoryV)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求人の開発経験 言語
func getStrLanguageDevelopmentExperienceListForJobInfo(inputList []entity.JobInformationRequiredExperienceDevelopment) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		if inputV.DevelopmentCategory.Valid {
			for _, inputTypeV := range inputV.ExperienceDevelopmentTypes {
				for typeI, typeV := range entity.DevelopmentTypeList[0] {
					if inputTypeV.DevelopmentType == null.NewInt(int64(typeI), true) {
						develomentSkillStr := fmt.Sprint(typeV, getStrYearAndMonthWithOver(inputV.ExperienceYear, inputV.ExperienceMonth))
						list = append(list, develomentSkillStr)
					}
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の開発経験 言語
func getStrLanguageDevelopmentExperienceListForJobSeeker(inputList []entity.JobSeekerDevelopmentSkill) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		if inputV.DevelopmentCategory.Valid {
			for typeI, typeV := range entity.DevelopmentTypeList[0] {
				if inputV.DevelopmentType == null.NewInt(int64(typeI), true) {
					develomentSkillStr := fmt.Sprint(typeV, getStrYearAndMonthWithOver(inputV.ExperienceYear, inputV.ExperienceMonth))
					list = append(list, develomentSkillStr)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求人の開発経験 OS
func getStrOSDevelopmentExperienceListForJobInfo(inputList []entity.JobInformationRequiredExperienceDevelopment) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		if inputV.DevelopmentCategory.Valid {
			for _, inputTypeV := range inputV.ExperienceDevelopmentTypes {
				for typeI, typeV := range entity.DevelopmentTypeList[1] {
					if inputTypeV.DevelopmentType == null.NewInt(int64(typeI), true) {
						develomentSkillStr := fmt.Sprint(typeV, getStrYearAndMonthWithOver(inputV.ExperienceYear, inputV.ExperienceMonth))
						list = append(list, develomentSkillStr)
					}
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の開発経験 OS
func getStrOSDevelopmentExperienceListForJobSeeker(inputList []entity.JobSeekerDevelopmentSkill) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		if inputV.DevelopmentCategory.Valid {
			for typeI, typeV := range entity.DevelopmentTypeList[1] {
				if inputV.DevelopmentType == null.NewInt(int64(typeI), true) {
					develomentSkillStr := fmt.Sprint(typeV, getStrYearAndMonthWithOver(inputV.ExperienceYear, inputV.ExperienceMonth))
					list = append(list, develomentSkillStr)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求人の語学
func getStrRequiredLanguageForJobInfo(inputV entity.JobInformationRequiredLanguage) string {
	var (
		list []string
	)

	for _, inputTypeV := range inputV.LanguageTypes {
		for typeI, typeV := range entity.LanguageType {
			if inputTypeV.LanguageType.Int64 == int64(typeI) {
				languageStr := fmt.Sprint(
					"(", typeV, ")",
					"\n  レベル:", getStrLanguageLevel(inputV.LanguageLevel),
					"\n  TOEIC:", inputV.Toeic.Int64,
					"\n  TOEFL(ibt):", inputV.ToeflIBT.Int64,
					"\n  TOEFL(pbt):", inputV.ToeflPBT.Int64,
				)
				list = append(list, languageStr)
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の語学
func getStrLanguageListForJobSeeker(inputList []entity.JobSeekerLanguageSkill) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		if inputV.LanguageType.Valid {
			for typeI, typeV := range entity.LanguageType {
				if inputV.LanguageType.Int64 == int64(typeI) {
					languageStr := fmt.Sprint(
						"(", typeV, ")",
						"\n  レベル:", getStrLanguageLevel(inputV.LanguageLevel),
						"\n  TOEIC:", inputV.Toeic.Int64,
						"\n  TOEFL(ibt):", inputV.ToeflIBT.Int64,
						"\n  TOEFL(pbt):", inputV.ToeflPBT.Int64,
					)
					list = append(list, languageStr)
				}
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の自己PR
func getStrSelfPromotionList(inputList []entity.JobSeekerSelfPromotion) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		selfPrStr := fmt.Sprint(
			"・", inputV.Title,
			"\n  ", inputV.Contents,
		)
		list = append(list, selfPrStr)
	}

	return strings.Join(list, ",")
}

// 求職者の希望業界
func getStrDesiredIndustryList(inputList []entity.JobSeekerDesiredIndustry) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		for _, industryV := range entity.Industry {
			if industryV[uint(inputV.DesiredIndustry.Int64)] != "" {
				desiredIndustryStr := fmt.Sprint(
					industryV[uint(inputV.DesiredIndustry.Int64)], "(第", inputV.DesiredRank.Int64, "希望群),",
				)
				list = append(list, desiredIndustryStr)
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の希望職種
func getStrDesiredOccupationList(inputList []entity.JobSeekerDesiredOccupation) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		for _, occupationV := range entity.Occupation {
			if occupationV[uint(inputV.DesiredOccupation.Int64)] != "" {
				desiredOccupationStr := fmt.Sprint(
					occupationV[uint(inputV.DesiredOccupation.Int64)], "(第", inputV.DesiredRank.Int64, "希望群),",
				)
				list = append(list, desiredOccupationStr)
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の希望勤務地
func getStrDesiredWorkLocationList(inputList []entity.JobSeekerDesiredWorkLocation) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		for prefectureI, prefectureV := range entity.Prefecture {
			if prefectureI == int(inputV.DesiredWorkLocation.Int64) {
				desiredWorkingLocationStr := fmt.Sprint(
					prefectureV, "(第", inputV.DesiredRank.Int64, "希望群),",
				)
				list = append(list, desiredWorkingLocationStr)
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の希望休日タイプ
func getStrDesiredHolidayTypeList(inputList []entity.JobSeekerDesiredHolidayType) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		for holidayTypeI, holidayTypeV := range entity.HolidayForJobSeeker {
			if holidayTypeI == uint(inputV.HolidayType.Int64) {
				desiredHolidayTypeStr := fmt.Sprint(
					holidayTypeV, ",",
				)
				list = append(list, desiredHolidayTypeStr)
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の希望企業きぼ
func getStrDesiredCompanyScaleList(inputList []entity.JobSeekerDesiredCompanyScale) string {
	var (
		list []string
	)

	for _, inputV := range inputList {
		for companyScaleI, companyScaleV := range entity.CompanyScale {
			if companyScaleI == int(inputV.DesiredCompanyScale.Int64) {
				desiredCompanyScaleStr := fmt.Sprint(
					companyScaleV, ",",
				)
				list = append(list, desiredCompanyScaleStr)
			}
		}
	}

	return strings.Join(list, ",")
}

// 求職者の学歴情報
func getStrStudentHistoryList(inputList []entity.JobSeekerStudentHistory) string {
	var (
		list []string
	)

	for inputI, inputV := range inputList {
		// 学校カテゴリー
		schoolCategoryStr := getStrSchoolCategory(inputV.SchoolCategory)

		studentHistoryStr := fmt.Sprint(
			"<学歴", inputI+1, ">",
			"\n・カテゴリー:", schoolCategoryStr,
			"\n・学校名:", inputV.SchoolName,
			"\n・入学年:", inputV.EntranceYear,
			"\n・開始ステータス:", getStrFirstStatusForStudentHistory(inputV.FirstStatus),
			"\n・卒業年:", inputV.GraduationYear,
			"\n・終了ステータス:", getStrLastStatusForStudentHistory(inputV.LastStatus),
		)
		list = append(list, studentHistoryStr)
	}

	return strings.Join(list, "\n\n\n\n")
}

// 求職者の職歴情報
func getStrWorkHistoryList(inputList []entity.JobSeekerWorkHistory) string {
	var (
		list []string
	)

	for inputI, inputV := range inputList {

		// 会社業界
		var industryStr string
		for _, industryV := range inputV.ExperienceIndustries {
			for _, industryMasterV := range entity.Industry {
				if industryMasterV[uint(industryV.Industry.Int64)] != "" {
					industryStr += fmt.Sprint(
						industryMasterV[uint(industryV.Industry.Int64)],
					)
				}
			}
		}

		workHistoryStr := fmt.Sprint(
			"<職歴", inputI+1, ">",
			"\n・会社名:", inputV.CompanyName,
			"\n・会社業界:", industryStr,
			"\n・従業員数(単体):", inputV.EmployeeNumberSingle.Int64,
			"\n・従業員数(連結):", inputV.EmployeeNumberGroup.Int64,
			"\n・雇用形態:", getStrEmploymentStatusForJobSeeker(inputV.EmploymentStatus),
			"\n・入社年:", inputV.JoiningYear,
			"\n・開始ステータス:", getStrFirstStatusForWorkHistory(inputV.FirstStatus),
			"\n・退社年:", inputV.RetireYear,
			"\n・終了ステータス:", getStrLastStatusForWorkHistory(inputV.LastStatus),
		)

		// 部門歴
		for departmentI, departmentV := range inputV.DepartmentHistories {

			// 職種
			var occupationStr string
			for _, occupationV := range departmentV.ExperienceOccupations {
				for _, occupationMasterV := range entity.Occupation {
					if occupationMasterV[uint(occupationV.Occupation.Int64)] != "" {
						occupationStr += fmt.Sprint(
							occupationMasterV[uint(occupationV.Occupation.Int64)],
						)
					}
				}
			}

			fmt.Println(departmentI, departmentV)
			workHistoryStr += fmt.Sprint(
				"\n・部門歴", departmentI+1, ":",
				"\n・部門名:", departmentV.Department,
				"\n・職種:", occupationStr,
				"\n・マネジメント人数:", departmentV.ManagementNumber.Int64,
				"\n・マネジメント詳細:", departmentV.ManagementDetail,
				"\n・開始年:", departmentV.StartYear,
				"\n・終了年:", departmentV.EndYear,
			)
		}
		list = append(list, workHistoryStr)
	}

	return strings.Join(list, "\n\n\n\n")
}

// 求人の必要業界・職種経験
func getStrExperienceJob(inputV entity.JobInformationRequiredExperienceJob) string {
	var (
		output        string
		industryStr   string
		occupationStr string
	)

	for _, industryV := range entity.Industry {
		for _, experienceIndustry := range inputV.ExperienceIndustries {
			if industryV[uint(experienceIndustry.ExperienceIndustry.Int64)] != "" {
				industryStr += fmt.Sprint(
					industryV[uint(experienceIndustry.ExperienceIndustry.Int64)], ",",
				)
			}
		}
	}

	for _, occupationV := range entity.Occupation {
		for _, experienceOccupation := range inputV.ExperienceOccupations {
			if occupationV[uint(experienceOccupation.ExperienceOccupation.Int64)] != "" {
				occupationStr += fmt.Sprint(
					occupationV[uint(experienceOccupation.ExperienceOccupation.Int64)], ",",
				)
			}
		}
	}

	output = fmt.Sprint("\n業界(", industryStr, ")\n", "職種(", occupationStr, ")", "\n", "年数", getStrYearAndMonthWithOver(inputV.ExperienceYear, inputV.ExperienceMonth))

	return output
}

func getStrSelectionFlowPattern(inputList []entity.JobInformationSelectionFlowPattern) string {
	var (
		list []string
	)

	for inputI, inputV := range inputList {
		var pattern []string
		pattern = append(pattern, fmt.Sprint("<選考パターン", inputI+1, ": ", inputV.FlowTitle, ">"))

		for _, inputInfo := range inputV.SelectionInformations {

			seletionInfo := fmt.Sprint(
				"\n・選考ステップ:", getStrTaskPhase(inputInfo.SelectionType),
				"\n・選考ポイント:", inputInfo.SelectionPoint,
				"\n・合格率:", inputInfo.PassedExample,
				"\n・NG例:", inputInfo.FailExample,
				"\n・通過率:", inputInfo.PassingRate.Int64, "%",
				"\n・選考後アンケートの有無:", getStrAvailableFromBool(inputInfo.IsQuestionnairy),
				"\n\n",
			)
			pattern = append(pattern, seletionInfo)

		}

		list = append(list, strings.Join(pattern, "\n"))
	}

	return strings.Join(list, "\n\n\n\n")
}

func getStrHRStaffList(inputList []entity.BillingAddressHRStaff) string {
	var (
		list []string
	)

	for inputI, inputV := range inputList {
		hrStaffStr := fmt.Sprint(
			"<担当者", inputI+1, ">",
			"\n・担当者名:", inputV.HRStaffName,
			"\n・メールアドレス:", inputV.HRStaffEmail,
			"\n・電話番号:", inputV.HRStaffPhoneNumber,
		)
		list = append(list, hrStaffStr)
	}

	return strings.Join(list, "\n")
}

func getStrRAStaffList(inputList []entity.BillingAddressRAStaff) string {
	var (
		list []string
	)

	for inputI, inputV := range inputList {
		raStaffStr := fmt.Sprint(
			"<担当者", inputI+1, ">",
			"\n・担当者名:", inputV.BillingAddressStaffName,
			"\n・メールアドレス:", inputV.BillingAddressStaffEmail,
			"\n・電話番号:", inputV.BillingAddressStaffPhoneNumber,
		)

		list = append(list, raStaffStr)
	}

	return strings.Join(list, "\n")
}

// 求人の必須条件を文字列で取得
func getStrRequiredCondition(inputList []entity.JobInformationRequiredCondition) string {
	var (
		list []string
	)

	if len(inputList) == 0 {
		return ""
	}

	for inputI, inputV := range inputList {
		conditionStr := fmt.Sprint(
			"<条件", inputI+1, ">",
			"\n・業界職種経験:", getStrExperienceJob(inputV.RequiredExperienceJobs),
			"\n\n・マネジメント経験の有無:", getStrAvailable(inputV.RequiredManagement),
			"\n\n・資格:", getStrLicenseListForJobInfo(inputV.RequiredLicenses),
			"\n\n・業務ツール:", getStrPCToolListForJobInfo(inputV.RequiredPCTools),
			"\n\n・開発経験(言語):", getStrLanguageDevelopmentExperienceListForJobInfo(inputV.RequiredExperienceDevelopments),
			"\n\n・開発経験(OS):", getStrOSDevelopmentExperienceListForJobInfo(inputV.RequiredExperienceDevelopments),
			"\n\n・語学力:", getStrRequiredLanguageForJobInfo(inputV.RequiredLanguages),
		)

		list = append(list, conditionStr)
	}

	return strings.Join(list, "\n\n\n\n")
}

/*************************************/

// ********* 文字列合体 汎用関数　*********//
//
// 年と月を合体させた文字列を取得 求人票とマスクレジュメ用　 例：2019年12月
func getStrYearAndMonth(year, month null.Int) string {

	yearStr := strconv.Itoa(int(year.Int64))
	monthStr := strconv.Itoa(int(month.Int64))

	if year.Valid && month.Valid {
		return yearStr + "年" + monthStr + "月"

	} else if year.Valid && !month.Valid {
		return yearStr + "年"

	} else if !year.Valid && month.Valid {
		return monthStr + "月"
	}

	return ""
}

// 年と月を合体させた文字列を取得 　 例：（2019年12ヶ月以上）
func getStrYearAndMonthWithOver(year, month null.Int) string {

	yearStr := strconv.Itoa(int(year.Int64))
	monthStr := strconv.Itoa(int(month.Int64))

	if year.Valid && month.Valid {
		return "(" + yearStr + "年" + monthStr + "ヶ月以上)"

	} else if year.Valid && !month.Valid {
		return "(" + yearStr + "年以上)"

	} else if !year.Valid && month.Valid {
		return "(" + monthStr + "ヶ月以上)"
	}

	return ""
}

// 単体と連結の従業員数を取得　例：100人（単体）・200人（連結）
func getStrEmployeeSingleAndGroup(single, group null.Int) string {
	singleStr := strconv.Itoa(int(single.Int64))
	groupStr := strconv.Itoa(int(group.Int64))

	if single.Valid && group.Valid {
		return singleStr + "人（単体）・" + groupStr + "人（連結）"

	} else if single.Valid && !group.Valid {
		return singleStr + "人（単体）"
	} else if !single.Valid && group.Valid {
		return groupStr + "人（連結）"
	}

	return ""
}

// 給料の下限と上限を取得　例：100万円～200万円
func getStrSalaryRange(lower, upper null.Int) string {
	lowerStr := strconv.Itoa(int(lower.Int64))
	upperStr := strconv.Itoa(int(upper.Int64))

	if lower.Valid && upper.Valid {
		return lowerStr + "万円～" + upperStr + "万円"

	} else if lower.Valid && !upper.Valid {
		return lowerStr + "万円～"

	} else if !lower.Valid && upper.Valid {
		return "〜" + upperStr + "万円"
	}

	return ""
}

func getStrFullAddress(postCode, prefecture, detail string) string {
	return postCode + " " + prefecture + " " + detail
}

// 2020-01-01 → 2020年1月1日
func getDateKanji(date string) string {
	if date == "" {
		return ""
	}
	return date[0:4] + "年" + date[5:7] + "月" + date[8:10] + "日"
}

// nullの場合,0ではなく空文字を返す toeicなどのスコア用
func getStrNull(inputedNumber null.Int) string {
	if inputedNumber.Valid {
		return strconv.Itoa(int(inputedNumber.Int64))
	}

	return ""
}

// ********* PDF変換用HTML作成 汎用パーツ　*********
func createDevelopmentSkillHTML(list []*entity.JobSeekerDevelopmentSkill) (outputHTML string) {
	var (
		outputArr []string
		langArr   []string
		osArr     []string
	)

	if len(list) < 1 {
		return ""
	}

	for i, v := range list {
		if i == 0 {
			outputArr = append(outputArr, `
			</br>
			<div>
			<table cellpadding="6px">
				<tbody>
				`)
		}

		yearAndMonth := getStrYearAndMonth(v.ExperienceYear, v.ExperienceMonth)

		_, skillType := getStrDevelopmentSkillCategoryAndType(v.DevelopmentCategory, v.DevelopmentType)

		switch v.DevelopmentCategory.Int64 {
		case 0:
			langArr = append(langArr, fmt.Sprintf(`%v（%v）`, skillType, yearAndMonth))
		case 1:
			osArr = append(osArr, fmt.Sprintf(`%v（%v）`, skillType, yearAndMonth))
		}

	}

	body := fmt.Sprintf(`
				<tr class="single">
					<th>
						開発言語
					</th>
					<td colspan="3">
						%v
					</td>
				</tr>
				<tr class="single">
					<th>
						開発OS
					</th>
					<td colspan="3">
						%v
					</td>
				</tr>
		`,
		strings.Join(langArr, ","),
		strings.Join(osArr, ","),
	)

	outputArr = append(outputArr, body)

	outputArr = append(outputArr, `
			</tbody>
			</table>
			</div>
			`)

	outputHTML = strings.Join(outputArr, "")

	return outputHTML
}

func createPCToolHTML(list []*entity.JobSeekerPCTool) (outputHTML string) {
	var (
		outputArr     []string
		excelArr      []string
		accessArr     []string
		wordArr       []string
		powerPointArr []string
		otherArr      []string
	)

	if len(list) < 1 {
		return ""
	}

	// for _, v := range list {
	// skillCategory, skillType := getStrPCToolCategoryAndType(v.SkillCategory, v.SkillType)

	// 	if skillCategory == "Excelスキル" {
	// 		excelArr = append(excelArr, skillType)
	// 	} else if skillCategory == "Accessスキル" {
	// 		accessArr = append(accessArr, skillType)
	// 	} else if skillCategory == "PowerPointスキル" {
	// 		powerPointArr = append(powerPointArr, skillType)
	// 	} else if skillCategory == "Wordスキル" {
	// 		wordArr = append(wordArr, skillType)
	// 	} else if skillCategory == "その他PCスキル" {
	// 		otherArr = append(otherArr, skillType)
	// 	}
	// }

	lenExcel := len(excelArr)
	lenAccess := len(accessArr)
	lenWord := len(wordArr)
	lenPowerPoint := len(powerPointArr)
	lenOther := len(otherArr)

	if lenExcel > 0 || lenAccess > 0 || lenWord > 0 || lenPowerPoint > 0 || lenOther > 0 {
		outputArr = append(outputArr, `
			</br>
			<div>
			<table cellpadding="6px">
				<tbody>
				`)
	}

	if len(excelArr) > 0 {
		outputArr = append(outputArr, fmt.Sprintf(
			`
			<tr class="single">
			<th>
				%v
			</th>
			<td colspan="3">
				%v
			</td>
		</tr>
		`,
			"Excel",
			strings.Join(excelArr, ","),
		))
	}

	if len(accessArr) > 0 {
		outputArr = append(outputArr, fmt.Sprintf(
			`
			<tr class="single">
			<th>
				%v
			</th>
			<td colspan="3">
				%v
			</td>
		</tr>
		`,
			"Access",
			strings.Join(accessArr, ","),
		))
	}

	if len(wordArr) > 0 {
		outputArr = append(outputArr, fmt.Sprintf(
			`
			<tr class="single">
			<th>
				%v
			</th>
			<td colspan="3">
				%v
			</td>
		</tr>
		`,
			"Word",
			strings.Join(wordArr, ","),
		))
	}

	if len(powerPointArr) > 0 {
		outputArr = append(outputArr, fmt.Sprintf(
			`
			<tr class="single">
			<th>
				%v
			</th>
			<td colspan="3">
				%v
			</td>
		</tr>
		`,
			"PowerPoint",
			strings.Join(powerPointArr, ","),
		))
	}

	if len(otherArr) > 0 {
		outputArr = append(outputArr, fmt.Sprintf(
			`
			<tr class="single">
			<th>
				%v
			</th>
			<td colspan="3">
				%v
			</td>
		</tr>
		`,
			"その他",
			strings.Join(otherArr, ","),
		))
	}

	if lenExcel > 0 || lenAccess > 0 || lenWord > 0 || lenPowerPoint > 0 || lenOther > 0 {
		outputArr = append(outputArr, `
			</tbody>
			</table>
			</div>
			`)

		outputHTML = strings.Join(outputArr, "")
	}

	return outputHTML
}

func createLanguageSkillHTML(list []*entity.JobSeekerLanguageSkill) (outputHTML string) {
	var (
		outputArr []string
		body      string
	)

	if len(list) < 1 {
		return ""
	}

	for i, v := range list {
		if i == 0 {
			outputArr = append(outputArr, `
			</br>
			<div>
			<table cellpadding="6px">
				<tbody>
				`)
		}

		body = fmt.Sprintf(`
	<tr class="single">
		<th>
			語学
		</th>
		<td colspan="3">
			%v
		</td>
	</tr>
	<tr class="single">
		<th>
			言語レベル
		</th>
		<td colspan="3">
			%v
		</td>
	</tr>
	<tr class="single">
		<th>
			TOEIC
		</th>
		<td colspan="3">
			%v
		</td>
	</tr>
	<tr class="double">
		<th>
			TOEFL(P)
		</th>
		<td>
			%v
		</td>
		<th>
			TOEFL(i)
		</th>
		<td>
			%v
		</td>
	</tr>
		`,
			getStrLanguage(v.LanguageType),
			getStrLanguageLevel(v.LanguageLevel),
			// getStrTalkingSkill(v.TalkingSkill),
			// getStrReadingSkill(v.ReadingSkill),
			// getStrWritingSkill(v.WritingSkill),
			// getStrYearAndMonth(v.BusinessExperienceYear, v.BusinessExperienceMonth),
			getStrNull(v.Toeic),
			getStrNull(v.ToeflPBT),
			getStrNull(v.ToeflIBT),
		)

		outputArr = append(outputArr, body)
	}

	outputArr = append(outputArr, `
			</tbody>
			</table>
			</div>
			`)

	outputHTML = strings.Join(outputArr, "")

	return outputHTML
}

func createLicenseHTML(list []*entity.JobSeekerLicense) (outputHTML string) {
	var (
		outputArr []string
		body      string
	)

	if len(list) < 1 {
		return `
			</br>
			<div>
				<p style="text-align: left; margin-bottom: 0px;">
				■所持資格/スキル
				</p>
			</div>
			`
	}

	for i, v := range list {
		if i == 0 {
			outputArr = append(outputArr, `
			</br>
			<div>
			<p style="text-align: left; margin-bottom: 0px;">
			■所持資格/スキル
			</p>
			<table cellpadding="6px">
				<tbody>
				`)
		}

		body = fmt.Sprintf(`
		<tr class="single">
		<th style="background-color: white;">
			%v
		</th>
		<td colspan="3">
			%v
		</td>
		</tr>
		`,
			v.AcquisitionTime,
			getStrLicenseType(v.LicenseType),
		)

		outputArr = append(outputArr, body)
	}

	outputArr = append(outputArr, `
			</tbody>
			</table>
			</div>
			`)

	outputHTML = strings.Join(outputArr, "")

	return outputHTML
}

// 学歴　履歴書用
// func createStudentHistoryHTML(list []*entity.JobSeekerStudentHistory) (outputHTML string) {
// 	var (
// 		outputArr []string
// 		body      string
// 	)

// 	if len(list) < 1 {
// 		return ""
// 	}

// 	// schoolCategoryの小さい順に並び替え
// 	sort.Slice(list, func(i, j int) bool { return list[i].SchoolCategory.Int64 < list[j].SchoolCategory.Int64 })

// 	for i, v := range list {
// 		if i == 0 {
// 			outputArr = append(outputArr, `
// 			</br>
// 			<div>
// 				<p style="text-align: left; margin-bottom: 0px;">
// 					◆学歴
// 				</p>
// 			<table cellpadding="6px">
// 				<tbody>
// 				`)
// 		}

// 		body = fmt.Sprintf(`
// 		<tr>
// 		<td colspan="4">
// 			%v %v（%v入学・%v%v）
// 		</td>
// 	</tr>
// 		`,
// 			v.SchoolName, v.Subject, v.EntranceYear, v.GraduationYear, getStrLastStatus(v.LastStatus),
// 		)

// 		outputArr = append(outputArr, body)
// 	}

// 	outputArr = append(outputArr, `
// 			</tbody>
// 			</table>
// 			</div>
// 			`)

// 	outputHTML = strings.Join(outputArr, "")

// 	return outputHTML
// }

// 最終学歴のみを取り出す　マスクレジュメ用
func createFinalStudentHistoryHTML(list []*entity.JobSeekerStudentHistory) (outputHTML string) {

	if len(list) < 1 {
		return ""
	}

	// schoolCategoryが大きい順に並べ替え
	(sort.Slice(list, func(i, j int) bool { return list[i].SchoolCategory.Int64 > list[j].SchoolCategory.Int64 }))

	// 最終学歴の番号を使う
	finalHistoryNumber := 0

	// 卒業, 中退, 退学, 卒業見込み, 修了, 修了見込み,
	// LastStatusが中退(I:1)と退学(I:2)の場合は、一つ前の学歴を使う
	for _, history := range list {
		if history.LastStatus.Int64 == 1 || history.LastStatus.Int64 == 2 || !history.LastStatus.Valid {
			finalHistoryNumber++
		} else {
			break
		}
	}

	outputHTML = fmt.Sprintf(`
		</br>
		<div>
			<p style="text-align: left; margin-bottom: 0px;">
				■学歴
			</p>
		<table cellpadding="6px">
			<tbody>
			<tr class="single">
			<th>
				最終学歴
			</th>
			<td colspan="3">
				%v（%v%v）
			</td>
		</tr>
		<tr>
			<th>
				大学レベル
			</th>
			<td>
				%v
			</td>
			<th>
				学部・学科・コース
			</th>
			<td>
				%v
			</td>
		</tr>
	</tbody>
	</table>
	</div>
		`,
		getStrFinalEducationForJobSeeker(list[finalHistoryNumber].SchoolCategory),
		// list[finalHistoryNumber].EntranceYear,
		list[finalHistoryNumber].GraduationYear,
		getStrLastStatus(list[finalHistoryNumber].LastStatus),
		getStrCollegeRank(list[finalHistoryNumber].SchoolLevel),
		list[finalHistoryNumber].Subject,
	)

	return outputHTML
}

func createWorkHistoryHTML(list []*entity.JobSeekerWorkHistory) (outputHTML string) {
	var (
		outputArr []string
		body      string
	)

	if len(list) < 1 {
		return ""
	}

	for i, v := range list {
		if i == 0 {
			outputArr = append(outputArr, `
			</br>
			<div>
				<p style="text-align: left; margin-bottom: 0px;">
					■職務経歴
				</p>
				`)
		}

		var (
			industryTypeArr []string
			occupationArr   []string
		)
		for _, vSub := range v.ExperienceIndustries {
			industryTypeArr = append(industryTypeArr, getStrIndustry(vSub.Industry))
		}

		for _, vSub := range v.DepartmentHistories {

			// 複数の職種をまとめる
			var ocupationList []string
			fmt.Println("vSub.ExperienceOccupations:", len(vSub.ExperienceOccupations))
			for _, occupation := range vSub.ExperienceOccupations {
				fmt.Println("occupation: ", occupation.Occupation.Int64, getStrOccupation(occupation.Occupation))

				if occupation.Occupation.Valid {
					ocupationList = append(ocupationList, getStrOccupation(occupation.Occupation))
				}
			}

			occupationArr = append(occupationArr, fmt.Sprintf(`
			<tr class="single">
			<th>
				期間
			</th>
			<th colspan="3">
				業務内容
			</th>
			</tr>
			<tr class="single">
			<td rowspan="2">
				%v<br />〜<br />%v
			</td>
			<td colspan="3">
				%v %v %v
			</td>
		</tr>
		<tr class="single">
			<td colspan="3">
				%v
			</td>
		</tr>
		`,
				vSub.StartYear,
				vSub.EndYear,
				vSub.Department,
				strings.Join(ocupationList, ","), // 複数の職種をまとめる
				vSub.ManagementDetail,
				vSub.JobDescription,
			))

		}

		body = fmt.Sprintf(`
	<table cellpadding="6px">
	<tbody>
	<p style="text-align: left; padding-left: 5px; font-size: 12px;">
		%v〜%v
	</p>
	<tr class="double">
		<th>
			業界
		</th>
		<td>
			%v
		</td>
		<th>
			従業員数
		</th>
		<td>
			%v
		</td>
	</tr>
	<tr class="double">
		<th>
			株式公開
		</th>
		<td>
			%v
		</td>
		<th>
			雇用形態
		</th>
		<td>
			%v
		</td>
	</tr>
	%v
	</tbody>
	</table>
		`,
			v.JoiningYear,
			v.RetireYear,
			strings.Join(industryTypeArr, ","),
			getStrEmployeeSingleAndGroup(v.EmployeeNumberSingle, v.EmployeeNumberGroup),
			getStrPublicOffering(v.PublicOffering),
			getStrEmploymentStatusForJobInfo(v.EmploymentStatus),
			strings.Join(occupationArr, ""),
		)

		outputArr = append(outputArr, body)
	}

	outputArr = append(outputArr, `
			</div>
			</br>
			</br>
			`)

	outputHTML = strings.Join(outputArr, "")

	return outputHTML
}

// デプロイ情報の最大ページ数を返す（本番実装までは1ページあたり20件）
func getDeploymentMaxPage(deploymentList []*entity.DeploymentInformation) uint {
	var maxPage = len(deploymentList) / 20

	if 0 < (len(deploymentList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページのデプロイ情報一覧を返す（本番実装までは1ページあたり20件）
func getDeploymentListWithPage(deploymentList []*entity.DeploymentInformation, page uint) []*entity.DeploymentInformation {
	var (
		perPage uint = 20
		listLen uint = uint(len(deploymentList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return deploymentList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.DeploymentInformation{}
	}

	if (listLen - first) <= perPage {
		return deploymentList[first:]
	}
	return deploymentList[first:last]
}

// 経験年数（LPで入力）
func getStrExperienceYear(inputNumber null.Int) string {
	if inputNumber == null.NewInt(9999, true) {
		return "まだ定まっていない"
	}
	if inputNumber.Valid {
		for i, v := range entity.ExperienceYear {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// 従業員数（LPで入力）
// 0: "10人未満", 5で登録
// 1: "10〜49人", 10で登録
// 2: "50〜99人", 50で登録
// 3: "100〜299人", 100で登録
// 4: "300〜999人", 300で登録
// 5: "1000〜2999人", 1000で登録
// 6: "3000〜4999人", 3000で登録
// 7: "5000人以上", 5000で登録

func getStrEmployeeNumber(inputNumber null.Int) null.Int {
	if inputNumber.Valid {
		switch inputNumber.Int64 {
		case 0:
			return null.NewInt(5, true)
		case 1:
			return null.NewInt(10, true)
		case 2:
			return null.NewInt(50, true)
		case 3:
			return null.NewInt(100, true)
		case 4:
			return null.NewInt(300, true)
		case 5:
			return null.NewInt(1000, true)
		case 6:
			return null.NewInt(3000, true)
		case 7:
			return null.NewInt(5000, true)
		}
	}

	return NullInt
}
