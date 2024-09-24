package routes

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

/***************** CSVファイルを読み込む（企業・請求先用） *********************/
func parseEnterpriseCSV(r *csv.Reader) ([]*entity.EnterpriseAndBillingAddress, []uint, error) {
	var (
		enterpriseList []*entity.EnterpriseAndBillingAddress
		hrStaff        entity.BillingAddressHRStaff
		raStaff        entity.BillingAddressRAStaff

		// 入らなかったレコードを記録
		recordCounter uint = 3
		missedRecords []uint
		breakCounter  uint
	)

	// カラムタイトルと項目説明のレコードを読み込み
	r.Read()
	r.Read()
	r.Read()

	for {
		var (
			enterprise entity.EnterpriseAndBillingAddress
		)

		record, err := r.Read()
		// 空行があった場合は終了
		if err == io.EOF || len(record) == 0 || breakCounter == 3 {
			break
		}

		// レコードの行数をカウント
		recordCounter++
		enterprise.RecordLine = recordCounter

		// 企業名が3回連続空の場合は終了
		if record[0] == "" {
			breakCounter++
			missedRecords = append(missedRecords, recordCounter)
			continue

		} else if breakCounter > 0 && record[0] != "" {
			breakCounter = 0
		}

		/**
		企業項目
		*/

		for i, column := range record {

			switch i {
			// 会社名
			case 0:
				enterprise.CompanyName = column

				// 企業担当者ID
			case 1:
				fmt.Println("raのID:", column)

				id, err := strconv.Atoi(column)
				if err != nil || column == "" {
					missedRecords = append(missedRecords, recordCounter)
					continue
				} else {
					enterprise.AgentStaffID = uint(id)
				}

				// 企業HP URL
			case 2:
				enterprise.CorporateSiteURL = column

			// 代表者名
			case 3:
				enterprise.Representative = column

				// 本社郵便番号 000-0000（半角ハイフンあり）
			case 4:
				fmt.Println("本社郵便番号:", column)
				enterprise.PostCode = column

				// 本社所在地　都道府県以下
			case 5:
				enterprise.OfficeLocation = column

				// 設立
			case 6:
				enterprise.Establishment = column

				// 従業員数（単体）
			case 7:
				val, err := strconv.Atoi(column)
				if err != nil {
					enterprise.EmployeeNumberSingle = null.NewInt(0, false)
				} else {
					enterprise.EmployeeNumberSingle = null.NewInt(int64(val), true)
				}

				// 従業員数（連結）
			case 8:
				val, err := strconv.Atoi(column)
				if err != nil {
					enterprise.EmployeeNumberGroup = null.NewInt(0, false)
				} else {
					enterprise.EmployeeNumberGroup = null.NewInt(int64(val), true)
				}

				// 資本金
			case 9:
				enterprise.Capital = column

				// 株式公開
			case 10:
				enterprise.PublicOffering = entity.GetIntPublicOffering(column)

				// 売上高（年度）
			case 11:
				val, err := strconv.Atoi(column)
				if err != nil {
					enterprise.EarningsYear = null.NewInt(0, false)
				} else {
					enterprise.EarningsYear = null.NewInt(int64(val), true)
				}

				// 売上（額）
			case 12:
				enterprise.Earnings = column

				// 事業内容
			case 13:
				enterprise.BusinessDetail = column

			// 業界
			case 14, 15, 16:

				// null.Int型に変換
				industry := entity.GetIntIndustry(column)

				// null.Intの配列作成
				if industry.Valid {
					enterprise.Industries = append(enterprise.Industries, industry)
				}

				/**
				請求先項目
				*/

				// 担当者ID
			case 17:
				// 未入力の場合、企業担当者を設定
				fmt.Println("担当者のID:", column)
				id, err := strconv.Atoi(column)
				if err != nil || column == "" {
					enterprise.AgentStaffIDForBillingAddress = enterprise.AgentStaffID
				} else {
					enterprise.AgentStaffIDForBillingAddress = uint(id)
				}

				// 請求先タイトル
			case 18:
				enterprise.BillingAddressTitle = column

				// 契約フェーズ
			case 19:
				enterprise.ContractPhase = entity.GetIntContractPhase(column)

				// 基本契約締結日
			case 20:
				enterprise.ContractDate = column

				// 支払い規定
			case 21:
				enterprise.PaymentPolicy = column

				// 請求先企業名
			case 22:
				if column != "" {
					column = enterprise.CompanyName
				}
				enterprise.BillingAddressCompanyName = column

				// 請求先住所
			case 23:
				enterprise.BillingAddressAddress = column

				// 人事担当者テーブル1
				// HR担当者名
			case 24, 27, 30:
				hrStaff = entity.BillingAddressHRStaff{}

				hrStaff.HRStaffName = column

			// HR担当者メールアドレス
			case 25, 28, 31:
				hrStaff.HRStaffEmail = column

			// HR担当者電話番号
			case 26, 29, 32:
				hrStaff.HRStaffPhoneNumber = column

				if hrStaff.HRStaffName != "" || hrStaff.HRStaffEmail != "" || hrStaff.HRStaffPhoneNumber != "" {
					enterprise.HRStaffs = append(enterprise.HRStaffs, hrStaff)
				}

				// 請求先担当者テーブル1
				// 請求先担当者名
			case 33, 36, 39:
				raStaff = entity.BillingAddressRAStaff{}

				raStaff.BillingAddressStaffName = column

				// 請求先担当者メールアドレス
			case 34, 37, 40:
				raStaff.BillingAddressStaffEmail = column

				// 請求先 電話番号
			case 35, 38, 41:
				raStaff.BillingAddressStaffPhoneNumber = column

				fmt.Println("raのstaff:", raStaff.BillingAddressStaffName)

				if raStaff.BillingAddressStaffName != "" || raStaff.BillingAddressStaffEmail != "" || raStaff.BillingAddressStaffPhoneNumber != "" {
					enterprise.RAStaffs = append(enterprise.RAStaffs, raStaff)
				}

				// 推薦方法
			case 42:
				enterprise.HowToRecommend = column

			}
		}

		enterpriseList = append(enterpriseList, &enterprise)
		fmt.Println("企業情報: ", enterprise, enterpriseList)
	}

	fmt.Println("除外されたレコード: ", missedRecords)

	return enterpriseList, missedRecords, nil
}

/***************** CSVファイルを読み込む（求人用） *********************/
func parseJobInformationCSV(r *csv.Reader) ([]*entity.JobInformation, []uint, error) {
	var (
		jobInformationList    []*entity.JobInformation
		target                entity.JobInformationTarget
		feature               entity.JobInformationFeature
		prefecture            entity.JobInformationPrefecture
		employmentStatus      entity.JobInformationEmploymentStatus
		workCharmPoint        entity.JobInformationWorkCharmPoint
		condition             entity.JobInformationRequiredCondition
		language              entity.JobInformationRequiredLanguage
		languageType          entity.JobInformationRequiredLanguageType
		pcTool                entity.JobInformationRequiredPCTool
		license               entity.JobInformationRequiredLicense
		developmentExperience entity.JobInformationRequiredExperienceDevelopment
		developmentType       entity.JobInformationRequiredExperienceDevelopmentType
		jobExperience         entity.JobInformationRequiredExperienceJob
		experienceIndustry    entity.JobInformationRequiredExperienceIndustry
		experienceOccupation  entity.JobInformationRequiredExperienceOccupation
		socialExperience      entity.JobInformationRequiredSocialExperience
		selection             entity.JobInformationSelectionFlowPattern
		selectionInfo         entity.JobInformationSelectionInformation
		occupation            entity.JobInformationOccupation

		//選考フロー数をカウント
		selectionFlowCounter int64
		isFirstSelection     bool
		isSecondSelection    bool
		isThirdSelection     bool
		isFourthSelection    bool
		isFifthSelection     bool

		// 入らなかったレコードを記録
		recordCounter uint = 3
		missedRecords []uint
		breakCounter  uint
	)

	// カラムタイトルと項目説明のレコードを読み込み
	r.Read()
	r.Read()
	r.Read()

	for {
		var (
			jobInformation entity.JobInformation
		)

		record, err := r.Read()
		// 空行があった場合は終了
		if err == io.EOF || len(record) == 0 || breakCounter == 3 {
			break
		}

		// レコードの行数をカウント
		recordCounter++
		jobInformation.RecordLine = recordCounter

		// 企業名が3回連続空の場合は終了
		if record[0] == "" {
			breakCounter++
			missedRecords = append(missedRecords, recordCounter)
			continue

		} else if breakCounter > 0 && record[0] != "" {
			breakCounter = 0
		}

		/**
		企業項目
		*/

		for i, column := range record {

			switch i {
			// 会社名
			case 0:
				jobInformation.CompanyName = column

				// 郵便番号
			case 1:
				jobInformation.PostCode = column

				// 請求先タイトル
			case 2:
				jobInformation.BillingAddressTitle = column

				/***************求人情報****************/

				//求人タイトル
			case 3:
				fmt.Println("求人タイトル:", column)
				jobInformation.Title = column

			// 募集対象　※複数選択
			case 4:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					target.Target = entity.GetIntUserStatus(str)

					if target.Target.Valid {
						jobInformation.Targets = append(jobInformation.Targets, target)
					}
				}

				//"雇用形態 ※複数選択
			case 5:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					employmentStatus.EmploymentStatus = entity.GetIntEmploymentStatusForJobInfo(str)
					if employmentStatus.EmploymentStatus.Valid {
						jobInformation.EmploymentStatuses = append(jobInformation.EmploymentStatuses, employmentStatus)
					}
				}

				// 雇用期間の定めの有無
			case 6:
				jobInformation.EmploymentPeriod = entity.GetIntAvailable(column)

				// 更新上限
			case 7:
				jobInformation.EmploymentPeriodDetail = column

				// 試用期間有無
			case 8:
				fmt.Println("試用期間有無:", column)
				jobInformation.TrialPeriod = entity.GetIntAvailable(column)

				// 試用期間詳細
			case 9:
				jobInformation.TrialPeriodDetail = column

				//募集状況
			case 10:
				jobInformation.RecruitmentState = entity.GetIntRecruitmentState(column)

			// 募集期限
			case 11:
				jobInformation.ExpirationDate = column

			// 職種 3つまで
			case 12, 13, 14:
				occupation.Occupation = entity.GetIntOccupation(column)
				if occupation.Occupation.Valid {
					jobInformation.Occupations = append(jobInformation.Occupations, occupation)
				}

				// '1,2,3'のようにカンマ区切りで複数の業界が入っている場合
				// strSlice := strings.Split(column, ",")
				// for _, str := range strSlice {
				// 	var strNullInt null.Int
				// 	str = strings.TrimSpace(str)
				// 	strInt, err := strconv.Atoi(str)
				// 	if err != nil {
				// 		continue
				// 	} else {
				// 		strNullInt = null.NewInt(int64(strInt), true)
				// 		println("nullInt:", fmt.Sprint(strNullInt))
				// 		jobInformation.Industries = append(jobInformation.Industries, strNullInt)
				// 	}
				// }

			// // 背景
			// case 15:
			// 	jobInformation.Background = entity.GetIntBackground(column)

			// 募集人数
			case 16:
				numberOfHires, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.NumberOfHires = null.NewInt(0, false)
				} else {
					jobInformation.NumberOfHires = null.NewInt(int64(numberOfHires), true)
				}

				// 求人の特徴 3つまで
			case 17, 18, 19:
				feature.Feature = entity.GetIntFeature(column)

				if feature.Feature.Valid {
					jobInformation.Features = append(jobInformation.Features, feature)
				}

			// 求人の魅力1 2つで1つ。3つまで
			case 20, 22, 24:
				workCharmPoint = entity.JobInformationWorkCharmPoint{}
				workCharmPoint.Title = column

			case 21, 23, 25:
				workCharmPoint.Contents = column

				if workCharmPoint.Title != "" || workCharmPoint.Contents != "" {
					jobInformation.WorkCharmPoints = append(jobInformation.WorkCharmPoints, workCharmPoint)
				}

			// 仕事内容
			case 26:
				jobInformation.WorkDetail = column

				//勤務地　都道府県 *"47都道府県 or 全国各地〜〜県まで入力複数ある場合は「,」で区切る"
			case 27:

				// 全国各地の場合は、47都道府県を設定
				// '1,2,3'のようにカンマ区切りで複数の業界が入っている場合
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {
					// var strNullInt null.Int
					str = strings.TrimSpace(str)

					strInt := entity.GetIntPrefecture(str)

					println("都道府県:", fmt.Sprint(strInt))

					prefecture.Prefecture = strInt

					if prefecture.Prefecture.Valid {
						jobInformation.Prefectures = append(jobInformation.Prefectures, prefecture)
					}
				}

			// 勤務地（雇入れ直後）
			case 28:
				jobInformation.WorkLocation = column

				// 転勤有無
			case 29:
				jobInformation.Transfer = entity.GetIntAvailable(column)

				//変更の範囲
			case 30:
				jobInformation.TransferDetail = column

				//年収（下限）
			case 31:
				underIncome, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.UnderIncome = null.NewInt(0, false)
				} else {
					jobInformation.UnderIncome = null.NewInt(int64(underIncome), true)
				}

				//年収（上限）
			case 32:
				overIncome, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.OverIncome = null.NewInt(0, false)
				} else {
					jobInformation.OverIncome = null.NewInt(int64(overIncome), true)
				}

				// 給与詳細・昇給賞与
			case 33:
				fmt.Println("給与詳細・昇給賞与:", column)
				jobInformation.Salary = column

			// 社会保険　雇用保険・労災保険・健康保険・厚生年金保険
			case 34:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)

					if str == "雇用保険" {
						jobInformation.EmploymentInsurance = true
					}
					if str == "労災保険" {
						jobInformation.AccidentInsurance = true
					}
					if str == "健康保険" {
						jobInformation.HealthInsurance = true
					}
					if str == "厚生年金保険" {
						jobInformation.PensionInsurance = true
					}
				}

				// 諸手当・福利厚生
			case 35:
				jobInformation.Insurance = column

			// 	// 固定残業代有無
			// case 36:
			// 	jobInformation.FixedOvertime = entity.GetIntAvailable(column)

			// 固定残業代超過分の支払い有無
			case 37:
				jobInformation.FixedOvertimePayment = entity.GetIntAvailable(column)

				// 固定残業代の詳細
			case 38:
				jobInformation.FixedOvertimeDetail = column

				// 勤務時間
			case 39:
				jobInformation.WorkTime = column

			// 	// 残業時間の有無
			// case 40:
			// 	jobInformation.Overtime = entity.GetIntAvailable(column)

			// 平均残業時間
			case 41:
				jobInformation.OvertimeAverage = column

				// 休日・休暇タイプ
			case 42:
				jobInformation.HolidayType = entity.GetIntHolidayForJobSeeker(column)

				// 休日・休暇詳細
			case 43:
				jobInformation.HolidayDetail = column

				// 受動喫煙対策の有無
			case 44:
				jobInformation.PassiveSmoking = entity.GetIntPassiveSmoking(column)

				//選考フロー
			case 45:
				jobInformation.SelectionFlow = column

				/*********業界・職種経験**************/

				// 応募資格（経験・スキルなど）
			case 46:
				fmt.Println("応募資格（経験・スキルなど）:", column)
				jobInformation.RequiredExperienceJobDetail = column

				//	応募資格（その他）
			case 47:
				// jobInformation.OtherRequired = column

				// 募集性別
			case 48:
				jobInformation.Gender = entity.GetIntGenderForJobInfo(column)

				//	国籍
			case 49:
				jobInformation.Nationality = entity.GetIntNationalityForJobInfo(column)

				//最終学歴
			case 50:
				jobInformation.FinalEducation = entity.GetIntFinalEducationForJobInfo(column)

				// 応募可能学科 文系/理系
			case 51:
				jobInformation.StudyCategory = entity.GetIntStudyCategoryForJobInfo(column)

				//"大学ランク（大卒以上のみ選択）
			case 52:
				jobInformation.SchoolLevel = entity.GetIntSchoolLevelForJobInfo(column)

				// 募集年齢（下限）
			case 53:
				ageUnder, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.AgeUnder = null.NewInt(0, false)
				} else {
					jobInformation.AgeUnder = null.NewInt(int64(ageUnder), true)
				}

				//募集年齢（上限）
			case 54:
				ageOver, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.AgeOver = null.NewInt(0, false)
				} else {
					jobInformation.AgeOver = null.NewInt(int64(ageOver), true)
				}

				//転職回数限度(●社まで)
			case 55:
				jobInformation.JobChange = entity.GetIntJobChange(column)

				//  短期離職（1年未満）
			case 56:
				jobInformation.ShortResignation = entity.GetIntConditionOrNot(column)

				//短期離職備考
			case 57:
				fmt.Println("短期離職備考:", column)
				jobInformation.ShortResignationRemarks = column

			// ※複数選択"	"社会人経験
			case 58:
				socialExperience = entity.JobInformationRequiredSocialExperience{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)

					socialExperience.SocialExperienceType = entity.GetIntSocialExperienceType(str)

					if socialExperience.SocialExperienceType.Valid {
						jobInformation.RequiredSocialExperiences = append(jobInformation.RequiredSocialExperiences, socialExperience)
					}
				}

				// 経験年数
			case 59:
				socialExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.SocialExperienceYear = null.NewInt(0, false)
				} else {
					jobInformation.SocialExperienceYear = null.NewInt(int64(socialExperienceYear), true)
				}

				//経験月数
			case 60:
				socialExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.SocialExperienceMonth = null.NewInt(0, false)
				} else {
					jobInformation.SocialExperienceMonth = null.NewInt(int64(socialExperienceMonth), true)
				}

			//"PCスキル（Excel）
			case 61:
				if column == "" {
					continue
				} else {
					// カテゴリーをexcellに指定
					jobInformation.ExcelSkill = entity.GetIntExcelSkill(column)
				}

			// 	"PCスキル（Word）
			case 62:
				if column == "" {
					continue
				} else {
					jobInformation.WordSkill = entity.GetIntWordSkill(column)

				}

			// "PCスキル（Power P）
			case 63:
				if column == "" {
					continue
				} else {
					jobInformation.PowerPointSkill = entity.GetIntPowerPointSkill(column)
				}

			//アピアランス
			case 64:
				jobInformation.Appearance = entity.GetIntAppearanceForJobInfo(column)

			//コミュニケーション
			case 65:
				jobInformation.Communication = entity.GetIntCommunicationForJobInfo(column)

			//論理的思考力
			case 66:
				fmt.Println("論理的思考力", column)
				jobInformation.Thinking = entity.GetIntThinkingForJobInfo(column)

			// 応募条件（エージェント向け情報）
			case 67:
				jobInformation.TargetDetail = column

				/***********   共通の必須条件（業界・職種・資格・語学力）  ***********/

				//必要業職種経験1
				condition = entity.JobInformationRequiredCondition{}
				condition.IsCommon = true
				jobExperience = entity.JobInformationRequiredExperienceJob{}

				//業界
			case 68:
				experienceIndustry = entity.JobInformationRequiredExperienceIndustry{}

				if column == "" {
					continue
				}

				industries := entity.GetIntIndustryFromBigCategory(column)
				fmt.Println("業界:", column, "->", industries)

				for _, industry := range industries {
					experienceIndustry.ExperienceIndustry = null.NewInt(int64(industry), true)
					if experienceIndustry.ExperienceIndustry.Valid {
						jobExperience.ExperienceIndustries = append(jobExperience.ExperienceIndustries, experienceIndustry)
					}
				}

				//職種 3つまで
			case 69, 70, 71:
				experienceOccupation = entity.JobInformationRequiredExperienceOccupation{}

				if column == "" {
					continue
				}

				experienceOccupation.ExperienceOccupation = entity.GetIntOccupation(column)
				fmt.Println("職種:", column, "->", experienceOccupation.ExperienceOccupation)

				if experienceOccupation.ExperienceOccupation.Valid {
					jobExperience.ExperienceOccupations = append(jobExperience.ExperienceOccupations, experienceOccupation)
				}

				//経験年数
			case 72:
				jobExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceYear = null.NewInt(int64(jobExperienceYear), true)
				}

				//経験年数
			case 73:
				jobExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceMonth = null.NewInt(int64(jobExperienceMonth), true)
				}

				// // 主要項目が未入力の場合は追加しない
				if len(jobExperience.ExperienceIndustries) > 0 ||
					len(jobExperience.ExperienceOccupations) > 0 ||
					jobExperience.ExperienceYear.Valid ||
					jobExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceJobs = jobExperience
				}

				fmt.Println("jobExperience:", jobExperience)
				fmt.Println("condition:", condition.RequiredExperienceJobs.ExperienceIndustries)
				fmt.Println("condition:", condition.RequiredExperienceJobs.ExperienceOccupations)

				// 必要マネジメント経験 "必要"という文字列で入ってくる。必要=0, true, 不要=0, false
			case 74:
				if column == "必要" {
					condition.RequiredManagement = null.NewInt(0, true)
				} else {
					condition.RequiredManagement = null.NewInt(0, false)
				}

			//PC業務ツール　3つまで
			case 75, 76, 77:
				pcTool = entity.JobInformationRequiredPCTool{}
				if column == "" {
					continue
				} else {
					pcTool.Tool = entity.GetIntPCTool(column)

					if pcTool.Tool.Valid {
						condition.RequiredPCTools = append(condition.RequiredPCTools, pcTool)
					}
				}

				developmentExperience = entity.JobInformationRequiredExperienceDevelopment{}

				/************* 開発経験 ****************/
				//	開発言語 3つまで
			case 78, 79, 80:
				developmentType = entity.JobInformationRequiredExperienceDevelopmentType{}

				langType := entity.GetIntDevelopmentType(0, column)
				developmentType.DevelopmentType = langType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 81:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 82:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

				//	開発OS  3つまで
			case 83, 84, 85:
				developmentType = entity.JobInformationRequiredExperienceDevelopmentType{}

				osType := entity.GetIntDevelopmentType(1, column)
				developmentType.DevelopmentType = osType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 86:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 87:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

			//免許・資格
			case 88:
				license = entity.JobInformationRequiredLicense{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					license.License = entity.GetIntLicenseType(str)

					if license.License.Valid {
						condition.RequiredLicenses = append(condition.RequiredLicenses, license)
					}
				}

				// ※複数選択"	語学
				// 語学タイプ 3つまで
				language = entity.JobInformationRequiredLanguage{}
			case 89, 90, 91:
				languageType = entity.JobInformationRequiredLanguageType{}

				languageType.LanguageType = entity.GetIntLanguageType(column)

				if languageType.LanguageType.Valid {
					language.LanguageTypes = append(language.LanguageTypes, languageType)
				}

			case 92:
				language.LanguageLevel = entity.GetIntLanguageLevel(column)

			//"TOEIC ※英語のみ"
			case 93:
				toeic, err := strconv.Atoi(column)
				if err != nil {
					language.Toeic = null.NewInt(0, false)
				} else {
					language.Toeic = null.NewInt(int64(toeic), true)
				}

				//	"TOEFL(i) ※英語のみ"
			case 94:
				toeflI, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflIBT = null.NewInt(0, false)
				} else {
					language.ToeflIBT = null.NewInt(int64(toeflI), true)
				}

				// "TOEFL(P) ※英語のみ"
			case 95:
				toeflP, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflPBT = null.NewInt(0, false)
				} else {
					language.ToeflPBT = null.NewInt(int64(toeflP), true)
				}

				// 全項目未入力時はスキップ
				if len(language.LanguageTypes) > 0 ||
					language.LanguageLevel.Valid ||
					language.Toeic.Valid ||
					language.ToeflIBT.Valid ||
					language.ToeflPBT.Valid {

					condition.RequiredLanguages = language
				}

				// 共通項目
				if len(condition.RequiredExperienceDevelopments) > 0 ||
					len(condition.RequiredPCTools) > 0 ||
					condition.RequiredManagement.Valid ||
					condition.RequiredExperienceJobs.ExperienceYear.Valid ||
					condition.RequiredExperienceJobs.ExperienceMonth.Valid ||
					len(condition.RequiredExperienceJobs.ExperienceIndustries) > 0 ||
					len(condition.RequiredExperienceJobs.ExperienceOccupations) > 0 ||
					len(condition.RequiredLicenses) > 0 ||
					len(condition.RequiredLanguages.LanguageTypes) > 0 {
					jobInformation.RequiredConditions = append(jobInformation.RequiredConditions, condition)
				}

				/***********   パターン別必須条件（業界・職種・資格・語学力）  ***********/

				//必要業職種経験1
				condition = entity.JobInformationRequiredCondition{}
				condition.IsCommon = false
				jobExperience = entity.JobInformationRequiredExperienceJob{}

				//業界
			case 96:
				experienceIndustry = entity.JobInformationRequiredExperienceIndustry{}

				if column == "" {
					continue
				}

				industries := entity.GetIntIndustryFromBigCategory(column)

				for _, industry := range industries {
					experienceIndustry.ExperienceIndustry = null.NewInt(int64(industry), true)
					if experienceIndustry.ExperienceIndustry.Valid {
						jobExperience.ExperienceIndustries = append(jobExperience.ExperienceIndustries, experienceIndustry)
					}
				}

				//職種 3つまで
			case 97, 98, 99:
				experienceOccupation = entity.JobInformationRequiredExperienceOccupation{}

				if column == "" {
					continue
				}

				experienceOccupation.ExperienceOccupation = entity.GetIntOccupation(column)

				if experienceOccupation.ExperienceOccupation.Valid {
					jobExperience.ExperienceOccupations = append(jobExperience.ExperienceOccupations, experienceOccupation)
				}

				//経験年数
			case 100:
				jobExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceYear = null.NewInt(int64(jobExperienceYear), true)
				}

				//経験年数
			case 101:
				jobExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceMonth = null.NewInt(int64(jobExperienceMonth), true)
				}

				// // 主要項目が未入力の場合は追加しない
				if len(jobExperience.ExperienceIndustries) > 0 ||
					len(jobExperience.ExperienceOccupations) > 0 ||
					jobExperience.ExperienceYear.Valid ||
					jobExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceJobs = jobExperience
				}

				// 必要マネジメント経験 "必要"という文字列で入ってくる。必要=0, true, 不要=0, false
			case 102:
				if column == "必要" {
					condition.RequiredManagement = null.NewInt(0, true)
				} else {
					condition.RequiredManagement = null.NewInt(0, false)
				}

			//PC業務ツール　3つまで
			case 103, 104, 105:
				pcTool = entity.JobInformationRequiredPCTool{}
				if column == "" {
					continue
				} else {
					pcTool.Tool = entity.GetIntPCTool(column)

					if pcTool.Tool.Valid {
						condition.RequiredPCTools = append(condition.RequiredPCTools, pcTool)
					}
				}

				developmentExperience = entity.JobInformationRequiredExperienceDevelopment{}

				/************* 開発経験 ****************/
				//	開発言語 3つまで
			case 106, 107, 108:
				developmentType = entity.JobInformationRequiredExperienceDevelopmentType{}

				langType := entity.GetIntDevelopmentType(0, column)
				developmentType.DevelopmentType = langType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 109:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 110:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

				//	開発OS  3つまで
			case 111, 112, 113:
				developmentType = entity.JobInformationRequiredExperienceDevelopmentType{}

				osType := entity.GetIntDevelopmentType(1, column)
				developmentType.DevelopmentType = osType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 114:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 115:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

			//免許・資格
			case 116:
				license = entity.JobInformationRequiredLicense{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					license.License = entity.GetIntLicenseType(str)

					if license.License.Valid {
						condition.RequiredLicenses = append(condition.RequiredLicenses, license)
					}
				}

				// ※複数選択"	語学
				// 語学タイプ 3つまで
				language = entity.JobInformationRequiredLanguage{}
			case 117, 118, 119:
				languageType = entity.JobInformationRequiredLanguageType{}

				languageType.LanguageType = entity.GetIntLanguageType(column)

				if languageType.LanguageType.Valid {
					language.LanguageTypes = append(language.LanguageTypes, languageType)
				}

			case 120:
				language.LanguageLevel = entity.GetIntLanguageLevel(column)

			//"TOEIC ※英語のみ"
			case 121:
				toeic, err := strconv.Atoi(column)
				if err != nil {
					language.Toeic = null.NewInt(0, false)
				} else {
					language.Toeic = null.NewInt(int64(toeic), true)
				}

				//	"TOEFL(i) ※英語のみ"
			case 122:
				toeflI, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflIBT = null.NewInt(0, false)
				} else {
					language.ToeflIBT = null.NewInt(int64(toeflI), true)
				}

				// "TOEFL(P) ※英語のみ"
			case 123:
				toeflP, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflPBT = null.NewInt(0, false)
				} else {
					language.ToeflPBT = null.NewInt(int64(toeflP), true)
				}

				// 全項目未入力時はスキップ
				if len(language.LanguageTypes) > 0 ||
					language.LanguageLevel.Valid ||
					language.Toeic.Valid ||
					language.ToeflIBT.Valid ||
					language.ToeflPBT.Valid {

					condition.RequiredLanguages = language
				}

				// 共通項目
				if len(condition.RequiredExperienceDevelopments) > 0 ||
					len(condition.RequiredPCTools) > 0 ||
					condition.RequiredManagement.Valid ||
					condition.RequiredExperienceJobs.ExperienceYear.Valid ||
					condition.RequiredExperienceJobs.ExperienceMonth.Valid ||
					len(condition.RequiredExperienceJobs.ExperienceIndustries) > 0 ||
					len(condition.RequiredExperienceJobs.ExperienceOccupations) > 0 ||
					len(condition.RequiredLicenses) > 0 ||
					len(condition.RequiredLanguages.LanguageTypes) > 0 {
					jobInformation.RequiredConditions = append(jobInformation.RequiredConditions, condition)
				}

				/***********   パターン別必須条件（業界・職種・資格・語学力）  ***********/

				//必要業職種経験1
				condition = entity.JobInformationRequiredCondition{}
				condition.IsCommon = false
				jobExperience = entity.JobInformationRequiredExperienceJob{}

				//業界
			case 124:
				experienceIndustry = entity.JobInformationRequiredExperienceIndustry{}

				if column == "" {
					continue
				}

				industries := entity.GetIntIndustryFromBigCategory(column)

				for _, industry := range industries {
					experienceIndustry.ExperienceIndustry = null.NewInt(int64(industry), true)
					if experienceIndustry.ExperienceIndustry.Valid {
						jobExperience.ExperienceIndustries = append(jobExperience.ExperienceIndustries, experienceIndustry)
					}
				}

				//職種 3つまで
			case 125, 126, 127:
				experienceOccupation = entity.JobInformationRequiredExperienceOccupation{}

				if column == "" {
					continue
				}

				experienceOccupation.ExperienceOccupation = entity.GetIntOccupation(column)

				if experienceOccupation.ExperienceOccupation.Valid {
					jobExperience.ExperienceOccupations = append(jobExperience.ExperienceOccupations, experienceOccupation)
				}

				//経験年数
			case 128:
				jobExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceYear = null.NewInt(int64(jobExperienceYear), true)
				}

				//経験年数
			case 129:
				jobExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceMonth = null.NewInt(int64(jobExperienceMonth), true)
				}

				// // 主要項目が未入力の場合は追加しない
				if len(jobExperience.ExperienceIndustries) > 0 ||
					len(jobExperience.ExperienceOccupations) > 0 ||
					jobExperience.ExperienceYear.Valid ||
					jobExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceJobs = jobExperience
				}

			// 必要マネジメント経験 "必要"という文字列で入ってくる。必要=0, true, 不要=0, false
			case 130:
				if column == "必要" {
					condition.RequiredManagement = null.NewInt(0, true)
				} else {
					condition.RequiredManagement = null.NewInt(0, false)
				}

			//PC業務ツール　3つまで
			case 131, 132, 133:
				pcTool = entity.JobInformationRequiredPCTool{}
				if column == "" {
					continue
				} else {
					pcTool.Tool = entity.GetIntPCTool(column)

					if pcTool.Tool.Valid {
						condition.RequiredPCTools = append(condition.RequiredPCTools, pcTool)
					}
				}

				developmentExperience = entity.JobInformationRequiredExperienceDevelopment{}

				/************* 開発経験 ****************/
				//	開発言語 3つまで
			case 134, 135, 136:
				developmentType = entity.JobInformationRequiredExperienceDevelopmentType{}

				langType := entity.GetIntDevelopmentType(0, column)
				developmentType.DevelopmentType = langType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 137:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 138:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

				//	開発OS  3つまで
			case 139, 140, 141:
				developmentType = entity.JobInformationRequiredExperienceDevelopmentType{}

				osType := entity.GetIntDevelopmentType(1, column)
				developmentType.DevelopmentType = osType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 142:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 143:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

			//免許・資格
			case 144:
				license = entity.JobInformationRequiredLicense{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					license.License = entity.GetIntLicenseType(str)

					if license.License.Valid {
						condition.RequiredLicenses = append(condition.RequiredLicenses, license)
					}
				}

				// ※複数選択"	語学
				// 語学タイプ 3つまで
				language = entity.JobInformationRequiredLanguage{}
			case 145, 146, 147:
				languageType = entity.JobInformationRequiredLanguageType{}

				languageType.LanguageType = entity.GetIntLanguageType(column)

				if languageType.LanguageType.Valid {
					language.LanguageTypes = append(language.LanguageTypes, languageType)
				}

			case 148:
				language.LanguageLevel = entity.GetIntLanguageLevel(column)

			//"TOEIC ※英語のみ"
			case 149:
				toeic, err := strconv.Atoi(column)
				if err != nil {
					language.Toeic = null.NewInt(0, false)
				} else {
					language.Toeic = null.NewInt(int64(toeic), true)
				}

				//	"TOEFL(i) ※英語のみ"
			case 150:
				toeflI, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflIBT = null.NewInt(0, false)
				} else {
					language.ToeflIBT = null.NewInt(int64(toeflI), true)
				}

				// "TOEFL(P) ※英語のみ"
			case 151:
				toeflP, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflPBT = null.NewInt(0, false)
				} else {
					language.ToeflPBT = null.NewInt(int64(toeflP), true)
				}

				// 全項目未入力時はスキップ
				if len(language.LanguageTypes) > 0 ||
					language.LanguageLevel.Valid ||
					language.Toeic.Valid ||
					language.ToeflIBT.Valid ||
					language.ToeflPBT.Valid {

					condition.RequiredLanguages = language
				}

				// 共通項目
				if len(condition.RequiredExperienceDevelopments) > 0 ||
					len(condition.RequiredPCTools) > 0 ||
					condition.RequiredManagement.Valid ||
					condition.RequiredExperienceJobs.ExperienceYear.Valid ||
					condition.RequiredExperienceJobs.ExperienceMonth.Valid ||
					len(condition.RequiredExperienceJobs.ExperienceIndustries) > 0 ||
					len(condition.RequiredExperienceJobs.ExperienceOccupations) > 0 ||
					len(condition.RequiredLicenses) > 0 ||
					len(condition.RequiredLanguages.LanguageTypes) > 0 {
					jobInformation.RequiredConditions = append(jobInformation.RequiredConditions, condition)
				}

				/*********** 備考・手数料 ***********/

			// 社内限定メモ
			case 152:
				jobInformation.SecretMemo = column

				// 推薦時必要な書類
			case 153:
				jobInformation.RequiredDocumentsDetail = column

				//成功報酬手数料（固定）
			case 154:
				commison, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.Commission = null.NewInt(0, false)
				} else {
					jobInformation.Commission = null.NewInt(int64(commison), true)
				}

				//成功報酬手数料（料率）
			case 155:
				fmt.Println("料率:", column)
				rate, err := strconv.Atoi(column)
				if err != nil {
					jobInformation.CommissionRate = null.NewInt(0, false)
				} else {
					jobInformation.CommissionRate = null.NewInt(int64(rate), true)
				}

				// 手数料補足
			case 156:
				jobInformation.CommissionDetail = column

				//返金規定
			case 157:
				jobInformation.RefundPolicy = column

				/******選考情報******/

				//選考フローのタイトル
			case 158:
				selection = entity.JobInformationSelectionFlowPattern{}
				selectionFlowCounter = 0
				isFirstSelection = false
				isSecondSelection = false
				isThirdSelection = false
				isFourthSelection = false
				isFifthSelection = false

				selection.FlowTitle = column

				fmt.Println("選考フロー1のタイトル:", column)

				//選考フローの公開設定
			case 159:
				// 未入力orエラーの場合は非公開にする
				openOrClose := entity.GetIntOpenOrClose(column)
				if !openOrClose.Valid {
					openOrClose = null.NewInt(1, true)
				}
				selection.PublicStatus = openOrClose

				/***書類選考***/
				//選考ポイント
			case 160:
				selectionInfo = entity.JobInformationSelectionInformation{}
				selectionInfo.SelectionPoint = column

				//合格事例
			case 161:
				selectionInfo.PassedExample = column

				//不合格事例
			case 162:
				selectionInfo.FailExample = column

				//通過率
			case 163:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 書類選考は、全項目が未入力の場合でも追加する
				selectionInfo.SelectionType = null.NewInt(1, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				/***一次選考***/
			case 164:
				selectionInfo = entity.JobInformationSelectionInformation{}
				// 選考の有無が「あり」の場合のみ、選考情報を追加する
				if column == "あり" {
					selectionFlowCounter += 1
					isFirstSelection = true
				}

				// アンケートの有無
			case 165:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 166:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 167:
				selectionInfo.PassedExample = column

				//不合格事例
			case 168:
				selectionInfo.FailExample = column

				//通過率
			case 169:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isFirstSelection {
					selectionInfo.SelectionType = null.NewInt(2, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***二次選考***/
				//選考ポイント
			case 170:
				selectionInfo = entity.JobInformationSelectionInformation{}
				if column == "あり" {
					selectionFlowCounter += 1
					isSecondSelection = true
				}

				// アンケートの有無
			case 171:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 172:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 173:
				selectionInfo.PassedExample = column

				//不合格事例
			case 174:
				selectionInfo.FailExample = column

				//通過率
			case 175:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isSecondSelection {
					selectionInfo.SelectionType = null.NewInt(3, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***三次選考***/
				//選考ポイント
			case 176:
				selectionInfo = entity.JobInformationSelectionInformation{}
				if column == "あり" {
					selectionFlowCounter += 1
					isThirdSelection = true
				}

				// アンケートの有無
			case 177:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

			case 178:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 179:
				selectionInfo.PassedExample = column

				//不合格事例
			case 180:
				selectionInfo.FailExample = column

				//通過率
			case 181:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isThirdSelection {
					selectionInfo.SelectionType = null.NewInt(4, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***四次選考***/
				//選考ポイント
			case 182:
				selectionInfo = entity.JobInformationSelectionInformation{}
				if column == "あり" {
					selectionFlowCounter += 1
					isFourthSelection = true
				}

				// アンケートの有無
			case 183:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 184:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 185:
				selectionInfo.PassedExample = column

				//不合格事例
			case 186:
				selectionInfo.FailExample = column

				//通過率
			case 187:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isFourthSelection {
					selectionInfo.SelectionType = null.NewInt(5, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***五次選考***/
			case 188:
				selectionInfo = entity.JobInformationSelectionInformation{}

				if column == "あり" {
					selectionFlowCounter += 1
					isFifthSelection = true
				}

				// アンケートの有無
			case 189:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 190:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 191:
				selectionInfo.PassedExample = column

				//不合格事例
			case 192:
				selectionInfo.FailExample = column

				//通過率
			case 193:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isFifthSelection {
					selectionInfo.SelectionType = null.NewInt(6, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***最終選考***/
				//選考ポイント

				// アンケートの有無
			case 194:
				selectionInfo = entity.JobInformationSelectionInformation{}

				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

			case 195:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 196:
				selectionInfo.PassedExample = column

				//不合格事例
			case 197:
				selectionInfo.FailExample = column

				//通過率
			case 198:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 最終選考は、全項目が未入力の場合でも追加する
				selectionInfo.SelectionType = null.NewInt(7, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				//内定(内定保留フェーズ)
			case 199:
				selectionInfo = entity.JobInformationSelectionInformation{}

				//内定承諾率
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 内定(内定保留)
				selectionInfo.IsQuestionnairy = false
				selectionInfo.SelectionPoint = ""
				selectionInfo.PassedExample = ""
				selectionInfo.FailExample = ""
				selectionInfo.SelectionType = null.NewInt(8, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				// 内定承諾
				selectionInfo = entity.JobInformationSelectionInformation{}
				selectionInfo.IsQuestionnairy = false
				selectionInfo.SelectionPoint = ""
				selectionInfo.PassedExample = ""
				selectionInfo.FailExample = ""
				selectionInfo.SelectionType = null.NewInt(9, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				// 選考パターンを追加
				// 書類選考と最終選考以外の選考でカウントされた数を追加
				selection.FlowPattern = null.NewInt(selectionFlowCounter, true)
				fmt.Println("選考パターン1: ", selection.FlowPattern)
				jobInformation.SelectionFlowPatterns = append(jobInformation.SelectionFlowPatterns, selection)

				/******選考情報 2つ目******/
				//選考フローのタイトル
			case 200:
				selection = entity.JobInformationSelectionFlowPattern{}
				selectionFlowCounter = 0
				isFirstSelection = false
				isSecondSelection = false
				isThirdSelection = false
				isFourthSelection = false
				isFifthSelection = false

				selection.FlowTitle = column

				fmt.Println("選考フロー2のタイトル:", column)

				//選考フローの公開設定
			case 201:
				// 未入力orエラーの場合は非公開にする
				openOrClose := entity.GetIntOpenOrClose(column)
				if !openOrClose.Valid {
					openOrClose = null.NewInt(1, true)
				}
				selection.PublicStatus = openOrClose

				/***書類選考***/
				//選考ポイント
			case 202:
				selectionInfo = entity.JobInformationSelectionInformation{}

				selectionInfo.SelectionPoint = column

				//合格事例
			case 203:
				selectionInfo.PassedExample = column

				//不合格事例
			case 204:
				selectionInfo.FailExample = column

				//通過率
			case 205:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 書類選考は、全項目が未入力の場合でも追加する
				selectionInfo.SelectionType = null.NewInt(1, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				/***一次選考***/
			case 206:
				selectionInfo = entity.JobInformationSelectionInformation{}

				// 選考の有無が「あり」の場合のみ、選考情報を追加する
				if column == "あり" {
					selectionFlowCounter += 1
					isFirstSelection = true
				}

				// アンケートの有無
			case 207:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 208:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 209:
				selectionInfo.PassedExample = column

				//不合格事例
			case 210:
				selectionInfo.FailExample = column

				//通過率
			case 211:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isFirstSelection {
					selectionInfo.SelectionType = null.NewInt(2, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***二次選考***/
				//選考ポイント
			case 212:
				selectionInfo = entity.JobInformationSelectionInformation{}

				if column == "あり" {
					selectionFlowCounter += 1
					isSecondSelection = true
				}

				// アンケートの有無
			case 213:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

			case 214:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 215:
				selectionInfo.PassedExample = column

				//不合格事例
			case 216:
				selectionInfo.FailExample = column

				//通過率
			case 217:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isSecondSelection {
					selectionInfo.SelectionType = null.NewInt(3, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***三次選考***/
				//選考ポイント
			case 218:
				selectionInfo = entity.JobInformationSelectionInformation{}

				if column == "あり" {
					selectionFlowCounter += 1
					isThirdSelection = true
				}

				// アンケートの有無
			case 219:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 220:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 221:
				selectionInfo.PassedExample = column

				//不合格事例
			case 222:
				selectionInfo.FailExample = column

				//通過率
			case 223:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isThirdSelection {
					selectionInfo.SelectionType = null.NewInt(4, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***四次選考***/
				//選考ポイント
			case 224:
				selectionInfo = entity.JobInformationSelectionInformation{}

				if column == "あり" {
					selectionFlowCounter += 1
					isFourthSelection = true
				}

				// アンケートの有無
			case 225:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				// 選考ポイント
			case 226:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 227:
				selectionInfo.PassedExample = column

				//不合格事例
			case 228:
				selectionInfo.FailExample = column

				//通過率
			case 229:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isFourthSelection {
					selectionInfo.SelectionType = null.NewInt(5, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***五次選考***/
			case 230:
				selectionInfo = entity.JobInformationSelectionInformation{}

				if column == "あり" {
					selectionFlowCounter += 1
					isFifthSelection = true
				}

				// アンケートの有無
			case 231:
				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 232:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 233:
				selectionInfo.PassedExample = column

				//不合格事例
			case 234:
				selectionInfo.FailExample = column

				//通過率
			case 235:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 主要項目が未入力の場合は追加しない
				if isFifthSelection {
					selectionInfo.SelectionType = null.NewInt(6, true)
					selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)
				}

				/***最終選考***/

				// アンケートの有無
			case 236:
				selectionInfo = entity.JobInformationSelectionInformation{}

				if column == "あり" {
					selectionInfo.IsQuestionnairy = true
				} else {
					selectionInfo.IsQuestionnairy = false
				}

				//選考ポイント
			case 237:
				selectionInfo.SelectionPoint = column

				//合格事例
			case 238:
				selectionInfo.PassedExample = column

				//不合格事例
			case 239:
				selectionInfo.FailExample = column

				//通過率
			case 240:
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 最終選考は、全項目が未入力の場合でも追加する
				selectionInfo.SelectionType = null.NewInt(7, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				//内定(内定保留)
			case 241:
				selectionInfo = entity.JobInformationSelectionInformation{}

				//内定承諾率
				rate, err := strconv.Atoi(column)
				if err != nil {
					selectionInfo.PassingRate = null.NewInt(0, false)
				} else {
					selectionInfo.PassingRate = null.NewInt(int64(rate), true)
				}

				// 内定(内定保留)
				selectionInfo.IsQuestionnairy = false
				selectionInfo.SelectionPoint = ""
				selectionInfo.PassedExample = ""
				selectionInfo.FailExample = ""
				selectionInfo.SelectionType = null.NewInt(8, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				// 内定承諾
				selectionInfo = entity.JobInformationSelectionInformation{}
				selectionInfo.IsQuestionnairy = false
				selectionInfo.SelectionPoint = ""
				selectionInfo.PassedExample = ""
				selectionInfo.FailExample = ""
				selectionInfo.SelectionType = null.NewInt(9, true)
				selection.SelectionInformations = append(selection.SelectionInformations, selectionInfo)

				// 選考パターンを追加
				// 書類選考と最終選考以外の選考でカウントされた数を追加
				selection.FlowPattern = null.NewInt(selectionFlowCounter, true)
				fmt.Println("選考パターン2: ", selection.FlowPattern)

				// 選考パターンが未入力の場合は追加しない
				if selection.FlowTitle != "" {
					jobInformation.SelectionFlowPatterns = append(jobInformation.SelectionFlowPatterns, selection)
				}

			}
		}

		jobInformationList = append(jobInformationList, &jobInformation)
	}

	for _, jobInformation := range jobInformationList {
		fmt.Println("選考パターン3: ", jobInformation.SelectionFlowPatterns)

		for _, condition := range jobInformation.RequiredConditions {
			fmt.Println("licnese: ", condition.RequiredLicenses)
			fmt.Println("language: ", condition.RequiredLanguages)
			fmt.Println("experienceJob: ", condition.RequiredExperienceJobs)
			fmt.Println("experienceIndustry: ", condition.RequiredExperienceJobs.ExperienceIndustries)
			fmt.Println("experienceDevelopment: ", condition.RequiredExperienceDevelopments)
		}
	}

	fmt.Println("除外されたレコードは: ", missedRecords)

	return jobInformationList, missedRecords, nil
}

/***************** CSVファイルを読み込む（サーカス求人用） *********************/
func parseEnterpriseCSVForCircus(r *csv.Reader, agentStaffID uint) ([]*entity.EnterpriseAndJobInformation, []uint, error) {
	var (
		enterpriseList []*entity.EnterpriseAndJobInformation

		// 入らなかったレコードを記録
		recordCounter uint = 3
		missedRecords []uint
		breakCounter  uint
		now           string = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
	)

	// カラムタイトルと項目説明のレコードを読み込み
	r.Read()

	for {

		// 企業情報を初期化
		enterprise := entity.EnterpriseAndJobInformation{
			AgentStaffID:  agentStaffID,
			RegisterPhase: null.NewInt(1, true),
			SecretMemo:    "・取り込み媒体: サーカスエージェント\n\n・取り込み日時: " + now,
		}

		record, err := r.Read()
		// 空行があった場合は終了
		if err == io.EOF || len(record) == 0 || breakCounter == 3 {
			break
		}

		// レコードの行数をカウント
		recordCounter++
		enterprise.RecordLine = recordCounter

		// 企業名が3回連続空の場合は終了
		if record[0] == "" {
			breakCounter++
			missedRecords = append(missedRecords, recordCounter)

		} else if breakCounter > 0 && record[0] != "" {
			breakCounter = 0
		}

		/**
		企業項目
		*/

		for i, column := range record {

			switch i {
			// 0: 企業ID
			case 0:
				// 被り判定ようにサーカス内での企業IDを取得
				i, _ := strconv.Atoi(column)
				enterprise.CircusEnterpriseID = uint(i)

				enterprise.SecretMemo += "\n\n・企業ID(サーカス内):" + column
			// 1: 求人ID
			case 1:
				enterprise.SecretMemo += "\n\n・求人ID(サーカス内):" + column

			// 2: 企業名
			case 2:
				enterprise.CompanyName = column
				enterprise.BillingAddressCompanyName = column
				enterprise.BillingAddressTitle = column + "請求先1"

				// 3: 会社HP
			case 3:
				enterprise.CorporateSiteURL = column

				// 4: 資本金
			case 4:
				enterprise.Capital = column

				// 5: 住所（都道府県）
			case 5:
				enterprise.OfficeLocation = column

				// 6: 住所
			case 6:
				enterprise.OfficeLocation += "\n" + column

				// 7: 従業員数
			case 7:
				enterprise.BusinessDetail += "\n\n・従業員数:" + column

				// 8: 上場区分
			case 8:
				enterprise.BusinessDetail += "\n\n・上場区分:" + column
				enterprise.PublicOffering = entity.GetIntPublicOfferingForCircus(column)

				// 9: 設立年
			case 9:
				// enterprise.Establishment = column + "年"
				enterprise.BusinessDetail += "\n\n・設立年:" + column + "年"

				// 10: 平均年齢
			case 10:
				enterprise.SecretMemo += "\n\n・平均年齢:" + column + "歳"

				// 11: 男女比率
			case 11:
				enterprise.SecretMemo += "\n\n・男女比率:" + column

				// 12: 求人名
			case 12:
				enterprise.Title = column

				// 13: 職種メイン
			case 13:
				enterprise.SecretMemo += "\n\n・職種メイン:" + column
				// occupation.Occupation

				// 14: 職種サブ
			case 14:
				enterprise.SecretMemo += "\n\n・職種サブ:" + column

			// 15: 業種メイン
			case 15:
				enterprise.SecretMemo += "\n\n・業種メイン:" + column

			// 16: 業種サブ
			case 16:
				enterprise.SecretMemo += "\n\n・業種サブ:" + column

			// 17: 最終学歴
			case 17:
				enterprise.SecretMemo += "\n\n・最終学歴:" + column
				enterprise.FinalEducation = entity.GetIntFinalEducationForCircus(column)

			// 18: PRポイント
			case 18:
				enterprise.WorkDetail += "\n\n・PRポイント:" + column

			// 19: 仕事内容・この仕事のミッション
			case 19:
				enterprise.WorkDetail += "\n\n・仕事内容・この仕事のミッション:" + column

			// 20: 応募時必須条件
			case 20:
				// enterprise.OtherRequired += column

			// 21: 職種経験
			case 21:
				enterprise.RequiredExperienceJobDetail += "・職種経験:" + column

			// 22: 業種経験
			case 22:
				enterprise.RequiredExperienceJobDetail += "\n\n・業種経験:" + column

			// 23: 勤務都道府県
			case 23:
				enterprise.WorkLocation = column
				prefectureInt := entity.GetIntPrefectureForCircus(column)
				if prefectureInt.Valid {
					prefecture := entity.JobInformationPrefecture{
						Prefecture: prefectureInt,
					}
					enterprise.Prefectures = append(enterprise.Prefectures, prefecture)
				}

				// 24: 勤務市区町村
			case 24:
				enterprise.WorkLocation += column

				// 25: 勤務地（雇入れ直後）
			case 25:
				enterprise.WorkLocation += column

				// 26: 休日
			case 26:
				enterprise.HolidayDetail += "\n\n・休日:" + column
				enterprise.HolidayType = entity.GetIntHolidayForCircus(column)

				// 27: 年間休日
			case 27:
				enterprise.HolidayDetail += "\n\n・年間休日:" + column + "日"

				// 28: その他の休日休暇
			case 28:
				enterprise.HolidayDetail += "\n\n・その他の休日休暇:" + column

				// 29: 休日休暇に関する補足情報
			case 29:
				enterprise.HolidayDetail += "\n\n・休日休暇に関する補足情報:" + column

				// 30: 福利厚生・諸手当
			case 30:
				enterprise.Insurance += "\n\n・福利厚生・諸手当:" + column

				// 31: その他の福利厚生・諸手当
			case 31:
				enterprise.Insurance += "\n\n・その他の福利厚生・諸手当:" + column

				// 32: 福利厚生・諸手当に関する補足情報
			case 32:
				enterprise.Insurance += "\n\n・福利厚生・諸手当に関する補足情報:" + column

				// 33: 始業時間
			case 33:
				// 0.357638888888889 -> 9:00
				columnFloat, err := strconv.ParseFloat(column, 64)
				if err != nil {
					enterprise.WorkTime += "\n\n・始業時間:" + column
					break
				}
				enterprise.WorkTime += "\n\n・始業時間:" + strconv.FormatFloat(columnFloat*24, 'f', 0, 64) + "時"

				// 34: 終業時間
			case 34:
				// 0.770833333333333 -> 18:30
				columnFloat, err := strconv.ParseFloat(column, 64)
				if err != nil {
					enterprise.WorkTime += "\n\n・終業時間:" + column
					break
				}
				enterprise.WorkTime += "\n\n・終業時間:" + strconv.FormatFloat(columnFloat*24, 'f', 0, 64) + "時"

				// 35: 夜間勤務
			case 35:
				enterprise.WorkTime += "\n\n・夜間勤務:" + column

				// 36: 月間平均残業時間
			case 36:
				enterprise.OvertimeAverage += column

				// 37: 勤務地・勤務時間について
			case 37:
				enterprise.WorkTime += "\n\n・勤務地・勤務時間について:" + column

				// 38: 年齢下限
			case 38:
				i, err := strconv.Atoi(column)
				if err != nil {
					enterprise.AgeUnder = null.NewInt(0, false)
				} else {
					enterprise.AgeUnder = null.NewInt(int64(i), true)
				}

				// 39: 年齢上限
			case 39:
				i, err := strconv.Atoi(column)
				if err != nil {
					enterprise.AgeOver = null.NewInt(0, false)
				} else {
					enterprise.AgeOver = null.NewInt(int64(i), true)
				}

				// 40: 性別
			case 40:
				enterprise.SecretMemo += "\n\n・性別:" + column
				enterprise.Gender = entity.GetIntGenderForCircus(column)

				// 41: 国籍
			case 41:
				enterprise.SecretMemo += "\n\n・国籍:" + column
				enterprise.Nationality = entity.GetIntNationalityForCircus(column)

				// 42: 想定年収下限
			case 42:
				enterprise.Salary = column + "〜"
				i, err := strconv.Atoi(column)
				if err != nil {
					enterprise.UnderIncome = null.NewInt(0, false)
				} else {
					enterprise.UnderIncome = null.NewInt(int64(i), true)
				}

				// 43: 想定年収上限
			case 43:
				enterprise.Salary += column + "万円"
				i, err := strconv.Atoi(column)
				if err != nil {
					enterprise.OverIncome = null.NewInt(0, false)
				} else {
					enterprise.OverIncome = null.NewInt(int64(i), true)
				}

				// 44: 年収例
			case 44:
				enterprise.Salary += "\n\n・年収例:" + column

				// 45: 給与条件について
			case 45:
				enterprise.Salary += "\n\n・給与条件について:" + column

				// 46: 給与情報に関する補足情報
			case 46:
				enterprise.Salary += "\n\n・給与情報に関する補足情報:" + column

				// 47: 賞与回数
			case 47:
				enterprise.Salary += "\n\n・賞与回数:" + column

				// 48: 昨年度賞与実績
			case 48:
				enterprise.Salary += "\n\n・昨年度賞与実績:" + column

				// 49: インセンティブ
			case 49:
				enterprise.Salary += "\n\n・インセンティブ:" + column

				// 50: 選考フロー
			case 50:
				enterprise.SelectionFlow = column

				// 51: 選考に関する補足情報
			case 51:
				enterprise.SelectionFlow += "\n\n・選考に関する補足情報:" + column

			// 	// 52: 募集背景
			// case 52:
			// 	enterprise.SecretMemo += "\n\n・募集背景:" + column
			// 	enterprise.Background = entity.GetIntBackgroundForCircus(column)

			// 53: 雇用形態
			case 53:
				enterprise.SecretMemo += "\n\n・雇用形態:" + column
				status := entity.GetIntEmploymentStatusForCircus(column)
				if status.Valid {
					employmentStatus := entity.JobInformationEmploymentStatus{
						EmploymentStatus: status,
					}
					enterprise.EmploymentStatuses = append(enterprise.EmploymentStatuses, employmentStatus)
				}

				// 54: 休暇制度
			case 54:
				enterprise.HolidayDetail += "\n\n・休暇制度:" + column

				// 55: 部署の人数
			case 55:
				enterprise.SecretMemo += "\n\n・部署の人数:" + column

				// 56: 部署の人員構成
			case 56:
				enterprise.SecretMemo += "\n\n・部署の人員構成:" + column

				// 57: 受動喫煙対策について
			case 57:
				enterprise.SecretMemo += "\n\n・受動喫煙対策について:" + column

				// 58: 求人募集状況
			case 58:
				enterprise.SecretMemo += "\n\n・求人募集状況:" + column
				enterprise.RecruitmentState = entity.GetIntOpenOrCloseForCircus(column)

			}
		}

		log.Println("企業情報を読み込みました: ", enterprise.CompanyName)

		enterpriseList = append(enterpriseList, &enterprise)
	}

	for _, enterprise := range enterpriseList {
		log.Println("企業情報を読み込みました after loop: ", enterprise.CompanyName)
	}

	fmt.Println("除外されたレコードは: ", missedRecords)

	return enterpriseList, missedRecords, nil
}

/***************** CSVファイルを読み込む（エージェントバンク求人用） *********************/
func parseEnterpriseCSVForAgentBank(r *csv.Reader, agentStaffID uint) ([]*entity.EnterpriseAndJobInformation, []uint, error) {
	var (
		enterpriseList []*entity.EnterpriseAndJobInformation

		// 入らなかったレコードを記録
		recordCounter uint = 3
		missedRecords []uint
		breakCounter  uint
	)

	// カラムタイトルと項目説明のレコードを読み込み
	r.Read()

	for {

		// 企業情報を初期化
		enterprise := entity.EnterpriseAndJobInformation{
			AgentStaffID:  agentStaffID,
			RegisterPhase: null.NewInt(1, true),
			SecretMemo:    "・取り込み媒体: エージェントバンク\n\n・取り込み日時: " + time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05"),
		}

		record, err := r.Read()
		// 空行があった場合は終了
		if err == io.EOF || len(record) == 0 || breakCounter == 3 {
			break
		}

		// レコードの行数をカウント
		recordCounter++
		enterprise.RecordLine = recordCounter

		// 企業名が3回連続空の場合は終了
		if record[0] == "" {
			breakCounter++
			missedRecords = append(missedRecords, recordCounter)

		} else if breakCounter > 0 && record[0] != "" {
			breakCounter = 0
		}

		/**
		企業項目
		*/

		for i, column := range record {

			switch i {
			// 0: agentbank上のid
			case 0:
				enterprise.SecretMemo += "\n\n・agentbank上のid: " + column

			// 1: ステータス
			case 1:
				enterprise.RecruitmentState = entity.GetIntRecruitmentState(column)

			// 2: 職種名
			case 2:
				enterprise.SecretMemo += "\n\n・職種名: " + column

			// 3: 求人キャッチコピー
			case 3:
				enterprise.Title = column

			// 4: 会社名
			case 4:
				enterprise.CompanyName = column

			// 5: 勤務地
			case 5:
				enterprise.WorkLocation = column

			// 6: 応募用メールアドレス
			case 6:
				enterprise.SecretMemo += "\n\n・応募用メールアドレス: " + column

			// 7: 電話番号（半角）
			case 7:
				enterprise.SecretMemo += "\n\n・電話番号（半角）: " + column

			// 8: 履歴書の有無
			case 8:
				enterprise.SecretMemo += "\n\n・履歴書の有無: " + column

			// 9: 直接訪問先の住所
			case 9:
				enterprise.SecretMemo += "\n\n・直接訪問先の住所: " + column

			// 10: 直接訪問先の追加説明
			case 10:
				enterprise.SecretMemo += "\n\n・直接訪問先の追加説明: " + column

			// 11: 雇用形態
			case 11:
				enterprise.SecretMemo += "\n\n・雇用形態: " + column
				employmentStatus := entity.GetIntEmploymentStatus(column)
				if employmentStatus.Valid {
					enterprise.EmploymentStatuses = append(enterprise.EmploymentStatuses, entity.JobInformationEmploymentStatus{
						EmploymentStatus: employmentStatus,
					})
				}

			// 12: 雇用形態
			case 12:
				enterprise.SecretMemo += "\n\n・雇用形態2: " + column
				employmentStatus := entity.GetIntEmploymentStatus(column)
				if employmentStatus.Valid {
					enterprise.EmploymentStatuses = append(enterprise.EmploymentStatuses, entity.JobInformationEmploymentStatus{
						EmploymentStatus: employmentStatus,
					})
				}

			// 13: 仕事内容（仕事内容）
			case 13:
				enterprise.WorkDetail = "\n\n・仕事内容（仕事内容）: " + column

			// 14: 仕事内容（アピールポイント）
			case 14:
				enterprise.WorkDetail += "\n\n・仕事内容（アピールポイント）: " + column

			// 15: 仕事内容（求める人材）
			case 15:
				enterprise.WorkDetail += "\n\n・仕事内容（求める人材）: " + column

			// 16: 仕事内容（勤務時間・曜日）
			case 16:
				enterprise.WorkTime = column

			// 17: 仕事内容（休日・休暇）
			case 17:
				enterprise.HolidayDetail = column

				// 18:  受動喫煙対策： 喫煙スペースあり
			case 18:
				enterprise.SecretMemo += "\n\n・受動喫煙対策： 喫煙スペースあり: " + column

			// 19: 仕事内容（アクセス）
			case 19:
				enterprise.WorkLocation += "\n\n・仕事内容（アクセス）: " + column

			// 20: 仕事内容（待遇・福利厚生）
			case 20:
				enterprise.Insurance = column

			// 21: 仕事内容（その他）
			case 21:
				if column != "" {
					enterprise.WorkDetail += "\n\n・仕事内容（その他）: " + column
				}

			// 22: 給与（下限）
			case 22:
				income, err := strconv.Atoi(column)
				if err != nil {
					enterprise.UnderIncome = null.NewInt(0, false)
				} else {
					enterprise.UnderIncome = null.NewInt(int64(income), true)
				}

			// 23: 給与（上限）
			case 23:
				income, err := strconv.Atoi(column)
				if err != nil {
					enterprise.OverIncome = null.NewInt(0, false)
				} else {
					enterprise.OverIncome = null.NewInt(int64(income), true)
				}

			// 24: 給与種別
			case 24:
				enterprise.SecretMemo += "\n\n・給与種別: " + column

			// 25: 職種カテゴリー
			case 25:
				enterprise.SecretMemo += "\n\n・職種カテゴリー: " + column

			// 26: 掲載画像
			case 26:
				enterprise.SecretMemo += "\n\n・掲載画像: " + column

			// 27: タグ
			case 27:
				if column != "" {
					enterprise.SecretMemo += "\n\n・タグ: " + column
				}

			// 28: タグ
			case 28:
				if column != "" {
					enterprise.SecretMemo += "\n\n・タグ2: " + column
				}

			// 29: タグ
			case 29:
				if column != "" {
					enterprise.SecretMemo += "\n\n・タグ3: " + column
				}

			// 30: タグ
			case 30:
				if column != "" {
					enterprise.SecretMemo += "\n\n・タグ4: " + column
				}

			// 31: タグ
			case 31:
				if column != "" {
					enterprise.SecretMemo += "\n\n・タグ5: " + column
				}

			// 32: 質問項目（職務経験）
			case 32:
				enterprise.SecretMemo += "\n\n・質問項目（職務経験）: " + column

			// 33: 質問項目（学歴）
			case 33:
				enterprise.SecretMemo += "\n\n・質問項目（学歴）: " + column

			// 34: 質問項目（勤務地）
			case 34:
				enterprise.SecretMemo += "\n\n・質問項目（勤務地）: " + column

			// 35: 質問項目（資格・免許）
			case 35:
				enterprise.SecretMemo += "\n\n・質問項目（資格・免許）: " + column

			// 36: 質問項目（言語）
			case 36:
				enterprise.SecretMemo += "\n\n・質問項目（言語）: " + column

			// 37: 質問項目（言語）
			case 37:
				enterprise.SecretMemo += "\n\n・質問項目（言語）2: " + column

			// 38: 質問項目（生年月日）
			case 38:
				enterprise.SecretMemo += "\n\n・質問項目（生年月日）: " + column

			// 39: 質問項目（性別）
			case 39:
				enterprise.SecretMemo += "\n\n・質問項目（性別）: " + column

			// 40: 質問項目（氏名のふりがな）
			case 40:
				enterprise.SecretMemo += "\n\n・質問項目（氏名のふりがな）: " + column

			// 41: 全ての要件を満たした候補者のみの通知を受け取る
			case 41:
				enterprise.SecretMemo += "\n\n・全ての要件を満たした候補者のみの通知を受け取る: " + column

			}
		}

		log.Println("企業情報を読み込みました: ", enterprise.CompanyName)

		enterpriseList = append(enterpriseList, &enterprise)
	}

	for _, enterprise := range enterpriseList {
		log.Println("企業情報を読み込みました after loop: ", enterprise.CompanyName)
	}

	fmt.Println("除外されたレコードは: ", missedRecords)

	return enterpriseList, missedRecords, nil
}

/***************** CSVファイルを読み込む（送客先求人用） *********************/
func parseSendingJobInformationCSV(r *csv.Reader) ([]*entity.SendingJobInformation, []uint, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
		target                    entity.SendingJobInformationTarget
		feature                   entity.SendingJobInformationFeature
		prefecture                entity.SendingJobInformationPrefecture
		employmentStatus          entity.SendingJobInformationEmploymentStatus
		workCharmPoint            entity.SendingJobInformationWorkCharmPoint
		condition                 entity.SendingJobInformationRequiredCondition
		language                  entity.SendingJobInformationRequiredLanguage
		languageType              entity.SendingJobInformationRequiredLanguageType
		pcTool                    entity.SendingJobInformationRequiredPCTool
		license                   entity.SendingJobInformationRequiredLicense
		developmentExperience     entity.SendingJobInformationRequiredExperienceDevelopment
		developmentType           entity.SendingJobInformationRequiredExperienceDevelopmentType
		jobExperience             entity.SendingJobInformationRequiredExperienceJob
		experienceIndustry        entity.SendingJobInformationRequiredExperienceIndustry
		experienceOccupation      entity.SendingJobInformationRequiredExperienceOccupation
		socialExperience          entity.SendingJobInformationRequiredSocialExperience
		occupation                entity.SendingJobInformationOccupation

		// 入らなかったレコードを記録
		recordCounter uint = 3
		missedRecords []uint
		breakCounter  uint
	)

	// カラムタイトルと項目説明のレコードを読み込み
	r.Read()
	r.Read()
	r.Read()

	for {
		var (
			sendingJobInformation entity.SendingJobInformation
		)

		// 初期値
		sendingJobInformation.RegisterPhase = null.NewInt(0, true) // 本登録

		record, err := r.Read()
		// 空行があった場合は終了
		if err == io.EOF || len(record) == 0 || breakCounter == 3 {
			break
		}

		// レコードの行数をカウント
		recordCounter++
		sendingJobInformation.RecordLine = recordCounter

		// 企業名が3回連続空の場合は終了
		if record[0] == "" {
			breakCounter++
			missedRecords = append(missedRecords, recordCounter)
			continue

		} else if breakCounter > 0 && record[0] != "" {
			breakCounter = 0
		}

		/**
		企業項目
		*/

		for i, column := range record {

			switch i {
			// 会社名
			case 0:
				sendingJobInformation.CompanyName = column

				// 企業HP URL
			case 1:
				sendingJobInformation.CorporateSiteURL = column

				// 本社郵便番号 000-0000（半角ハイフンあり）
			case 2:
				fmt.Println("本社郵便番号:", column)
				sendingJobInformation.PostCode = column

				// 本社所在地　都道府県以下
			case 3:
				sendingJobInformation.OfficeLocation = column

				// 設立
			case 4:
				sendingJobInformation.Establishment = column

				// 従業員数（単体）
			case 5:
				val, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.EmployeeNumberSingle = null.NewInt(0, false)
				} else {
					sendingJobInformation.EmployeeNumberSingle = null.NewInt(int64(val), true)
				}

				// 従業員数（連結）
			case 6:
				val, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.EmployeeNumberGroup = null.NewInt(0, false)
				} else {
					sendingJobInformation.EmployeeNumberGroup = null.NewInt(int64(val), true)
				}

				// 株式公開
			case 7:
				sendingJobInformation.PublicOffering = entity.GetIntPublicOffering(column)

				// 売上高（年度）
			case 8:
				val, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.EarningsYear = null.NewInt(0, false)
				} else {
					sendingJobInformation.EarningsYear = null.NewInt(int64(val), true)
				}

				// 売上（額）
			case 9:
				sendingJobInformation.Earnings = column

				// 事業内容
			case 10:
				sendingJobInformation.BusinessDetail = column

			// 業界
			case 11, 12, 13:

				// null.Int型に変換
				industryInt := entity.GetIntIndustry(column)

				// null.Intの配列作成
				if industryInt.Valid {
					industry := entity.SendingJobInformationIndustry{
						Industry: industryInt,
					}
					sendingJobInformation.Industries = append(sendingJobInformation.Industries, industry)
				}

				/***************求人情報****************/

				//求人タイトル
			case 14:
				fmt.Println("求人タイトル:", column)
				sendingJobInformation.Title = column

			// 募集対象
			case 15:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					target.Target = entity.GetIntUserStatus(str)

					if target.Target.Valid {
						sendingJobInformation.Targets = append(sendingJobInformation.Targets, target)
					}
				}

				//"雇用形態
			case 16:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					employmentStatus.EmploymentStatus = entity.GetIntEmploymentStatusForJobInfo(str)
					if employmentStatus.EmploymentStatus.Valid {
						sendingJobInformation.EmploymentStatuses = append(sendingJobInformation.EmploymentStatuses, employmentStatus)
					}
				}

				// 雇用期間の定めの有無
			case 17:
				sendingJobInformation.EmploymentPeriod = entity.GetIntAvailable(column)

				// 更新上限
			case 18:
				sendingJobInformation.EmploymentPeriodDetail = column

				// 試用期間有無
			case 19:
				fmt.Println("試用期間有無:", column)
				sendingJobInformation.TrialPeriod = entity.GetIntAvailable(column)

				// 試用期間詳細
			case 20:
				sendingJobInformation.TrialPeriodDetail = column

				//募集状況
			case 21:
				sendingJobInformation.RecruitmentState = entity.GetIntRecruitmentState(column)

			// 募集期限
			case 22:
				sendingJobInformation.ExpirationDate = column

			// 職種 3つまで
			case 23, 24, 25:
				occupation.Occupation = entity.GetIntOccupation(column)
				if occupation.Occupation.Valid {
					sendingJobInformation.Occupations = append(sendingJobInformation.Occupations, occupation)
				}

			// 背景
			case 26:
				// sendingJobInformation.Background = entity.GetIntBackground(column)
				sendingJobInformation.Background = null.NewInt(0, false)

			// 募集人数
			case 27:
				numberOfHires, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.NumberOfHires = null.NewInt(0, false)
				} else {
					sendingJobInformation.NumberOfHires = null.NewInt(int64(numberOfHires), true)
				}

			// 求人の特徴 3つまで
			case 28, 29, 30:
				feature.Feature = entity.GetIntFeature(column)

				if feature.Feature.Valid {
					sendingJobInformation.Features = append(sendingJobInformation.Features, feature)
				}

			// 求人の魅力1 2つで1つ。3つまで
			case 31, 33, 35:
				workCharmPoint = entity.SendingJobInformationWorkCharmPoint{}
				workCharmPoint.Title = column

			case 32, 34, 36:
				workCharmPoint.Contents = column

				if workCharmPoint.Title != "" || workCharmPoint.Contents != "" {
					sendingJobInformation.WorkCharmPoints = append(sendingJobInformation.WorkCharmPoints, workCharmPoint)
				}

			// 仕事内容
			case 37:
				sendingJobInformation.WorkDetail = column

				//勤務地　都道府県 *"47都道府県 or 全国各地〜〜県まで入力複数ある場合は「,」で区切る"
			case 38:

				// 全国各地の場合は、47都道府県を設定
				// '1,2,3'のようにカンマ区切りで複数の業界が入っている場合
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {
					str = strings.TrimSpace(str)

					strInt := entity.GetIntPrefecture(str)

					println("都道府県:", fmt.Sprint(strInt))

					prefecture.Prefecture = strInt

					if prefecture.Prefecture.Valid {
						sendingJobInformation.Prefectures = append(sendingJobInformation.Prefectures, prefecture)
					}
				}

			// 勤務地（雇入れ直後）
			case 39:
				sendingJobInformation.WorkLocation = column

				// 転勤有無
			case 40:
				sendingJobInformation.Transfer = entity.GetIntAvailable(column)

				//変更の範囲
			case 41:
				sendingJobInformation.TransferDetail = column

				//年収（下限）
			case 42:
				underIncome, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.UnderIncome = null.NewInt(0, false)
				} else {
					sendingJobInformation.UnderIncome = null.NewInt(int64(underIncome), true)
				}

				//年収（上限）
			case 43:
				overIncome, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.OverIncome = null.NewInt(0, false)
				} else {
					sendingJobInformation.OverIncome = null.NewInt(int64(overIncome), true)
				}

				// 給与詳細・昇給賞与
			case 44:
				fmt.Println("給与詳細・昇給賞与:", column)
				sendingJobInformation.Salary = column

			// 社会保険　雇用保険・労災保険・健康保険・厚生年金保険
			case 45:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)

					if str == "雇用保険" {
						sendingJobInformation.EmploymentInsurance = true
					}
					if str == "労災保険" {
						sendingJobInformation.AccidentInsurance = true
					}
					if str == "健康保険" {
						sendingJobInformation.HealthInsurance = true
					}
					if str == "厚生年金保険" {
						sendingJobInformation.PensionInsurance = true
					}
				}

				// 諸手当・福利厚生
			case 46:
				sendingJobInformation.Insurance = column

				// 固定残業代有無
			case 47:
				sendingJobInformation.FixedOvertime = entity.GetIntAvailable(column)

				// 固定残業代超過分の支払い有無
			case 48:
				sendingJobInformation.FixedOvertimePayment = entity.GetIntAvailable(column)

				// 固定残業代の詳細
			case 49:
				sendingJobInformation.FixedOvertimeDetail = column

				// 勤務時間
			case 50:
				sendingJobInformation.WorkTime = column

				// 残業時間の有無
			case 51:
				sendingJobInformation.Overtime = entity.GetIntAvailable(column)

				// 平均残業時間
			case 52:
				sendingJobInformation.OvertimeAverage = column

				// 休日・休暇タイプ
			case 53:
				sendingJobInformation.HolidayType = entity.GetIntHolidayForJobSeeker(column)

				// 休日・休暇詳細
			case 54:
				sendingJobInformation.HolidayDetail = column

				// 受動喫煙対策の有無
			case 55:
				sendingJobInformation.PassiveSmoking = entity.GetIntPassiveSmoking(column)

				//選考フロー
			case 56:
				sendingJobInformation.SelectionFlow = column

				/*********業界・職種経験**************/

				// 応募資格（経験・スキルなど）
			case 57:
				fmt.Println("応募資格（経験・スキルなど）:", column)
				sendingJobInformation.RequiredExperienceJobDetail = column

				//	応募資格（その他）
			case 58:
				// sendingJobInformation.OtherRequired = column

				// 募集性別
			case 59:
				sendingJobInformation.Gender = entity.GetIntGenderForJobInfo(column)

				//	国籍
			case 60:
				sendingJobInformation.Nationality = entity.GetIntNationalityForJobInfo(column)

				//最終学歴
			case 61:
				sendingJobInformation.FinalEducation = entity.GetIntFinalEducationForJobInfo(column)

				// 応募可能学科 文系/理系
			case 62:
				sendingJobInformation.StudyCategory = entity.GetIntStudyCategoryForJobInfo(column)

				//"大学ランク（大卒以上のみ選択）
			case 63:
				sendingJobInformation.SchoolLevel = entity.GetIntSchoolLevelForJobInfo(column)

				// 募集年齢（下限）
			case 64:
				ageUnder, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.AgeUnder = null.NewInt(0, false)
				} else {
					sendingJobInformation.AgeUnder = null.NewInt(int64(ageUnder), true)
				}

				//募集年齢（上限）
			case 65:
				ageOver, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.AgeOver = null.NewInt(0, false)
				} else {
					sendingJobInformation.AgeOver = null.NewInt(int64(ageOver), true)
				}

				//転職回数限度(●社まで)
			case 66:
				sendingJobInformation.JobChange = entity.GetIntJobChange(column)

				//  短期離職（1年未満）
			case 67:
				sendingJobInformation.ShortResignation = entity.GetIntConditionOrNot(column)

				//短期離職備考
			case 68:
				fmt.Println("短期離職備考:", column)
				sendingJobInformation.ShortResignationRemarks = column

			// ※複数選択"	"社会人経験
			case 69:
				socialExperience = entity.SendingJobInformationRequiredSocialExperience{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)

					socialExperience.SocialExperienceType = entity.GetIntSocialExperienceType(str)

					if socialExperience.SocialExperienceType.Valid {
						sendingJobInformation.RequiredSocialExperiences = append(sendingJobInformation.RequiredSocialExperiences, socialExperience)
					}
				}

				// 経験年数
			case 70:
				socialExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.SocialExperienceYear = null.NewInt(0, false)
				} else {
					sendingJobInformation.SocialExperienceYear = null.NewInt(int64(socialExperienceYear), true)
				}

				//経験月数
			case 71:
				socialExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					sendingJobInformation.SocialExperienceMonth = null.NewInt(0, false)
				} else {
					sendingJobInformation.SocialExperienceMonth = null.NewInt(int64(socialExperienceMonth), true)
				}

			//"PCスキル（Excel）
			case 72:
				if column == "" {
					continue
				} else {
					// カテゴリーをexcellに指定
					sendingJobInformation.ExcelSkill = entity.GetIntExcelSkill(column)
				}

			// 	"PCスキル（Word）
			case 73:
				if column == "" {
					continue
				} else {
					sendingJobInformation.WordSkill = entity.GetIntWordSkill(column)

				}

			// "PCスキル（Power P）
			case 74:
				if column == "" {
					continue
				} else {
					sendingJobInformation.PowerPointSkill = entity.GetIntPowerPointSkill(column)
				}

			//アピアランス
			case 75:
				sendingJobInformation.Appearance = entity.GetIntAppearanceForJobInfo(column)

			//コミュニケーション
			case 76:
				sendingJobInformation.Communication = entity.GetIntCommunicationForJobInfo(column)

			//論理的思考力
			case 77:
				fmt.Println("論理的思考力", column)
				sendingJobInformation.Thinking = entity.GetIntThinkingForJobInfo(column)

			// 応募条件（エージェント向け情報）
			case 78:
				sendingJobInformation.TargetDetail = column

				/***********   共通の必須条件（業界・職種・資格・語学力）  ***********/

				//必要業職種経験1
				condition = entity.SendingJobInformationRequiredCondition{}
				condition.IsCommon = true
				jobExperience = entity.SendingJobInformationRequiredExperienceJob{}

				//業界
			case 79:
				experienceIndustry = entity.SendingJobInformationRequiredExperienceIndustry{}

				if column == "" {
					continue
				}

				industries := entity.GetIntIndustryFromBigCategory(column)
				fmt.Println("業界:", column, "->", industries)

				for _, industry := range industries {
					experienceIndustry.ExperienceIndustry = null.NewInt(int64(industry), true)
					if experienceIndustry.ExperienceIndustry.Valid {
						jobExperience.ExperienceIndustries = append(jobExperience.ExperienceIndustries, experienceIndustry)
					}
				}

				//職種 3つまで
			case 80, 81, 82:
				experienceOccupation = entity.SendingJobInformationRequiredExperienceOccupation{}

				if column == "" {
					continue
				}

				experienceOccupation.ExperienceOccupation = entity.GetIntOccupation(column)
				fmt.Println("職種:", column, "->", experienceOccupation.ExperienceOccupation)

				if experienceOccupation.ExperienceOccupation.Valid {
					jobExperience.ExperienceOccupations = append(jobExperience.ExperienceOccupations, experienceOccupation)
				}

				//経験年数
			case 83:
				jobExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceYear = null.NewInt(int64(jobExperienceYear), true)
				}

				//経験年数
			case 84:
				jobExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceMonth = null.NewInt(int64(jobExperienceMonth), true)
				}

				// // 主要項目が未入力の場合は追加しない
				if len(jobExperience.ExperienceIndustries) > 0 ||
					len(jobExperience.ExperienceOccupations) > 0 ||
					jobExperience.ExperienceYear.Valid ||
					jobExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceJobs = jobExperience
				}

				fmt.Println("jobExperience:", jobExperience)
				fmt.Println("condition:", condition.RequiredExperienceJobs.ExperienceIndustries)
				fmt.Println("condition:", condition.RequiredExperienceJobs.ExperienceOccupations)

				// 必要マネジメント経験 "必要"という文字列で入ってくる。必要=0, true, 不要=0, false
			case 85:
				if column == "必要" {
					condition.RequiredManagement = null.NewInt(0, true)
				} else {
					condition.RequiredManagement = null.NewInt(0, false)
				}

			//PC業務ツール　3つまで
			case 86, 87, 88:
				pcTool = entity.SendingJobInformationRequiredPCTool{}
				if column == "" {
					continue
				} else {
					pcTool.Tool = entity.GetIntPCTool(column)

					if pcTool.Tool.Valid {
						condition.RequiredPCTools = append(condition.RequiredPCTools, pcTool)
					}
				}

				developmentExperience = entity.SendingJobInformationRequiredExperienceDevelopment{}

				/************* 開発経験 ****************/
				//	開発言語 3つまで
			case 89, 90, 91:
				developmentType = entity.SendingJobInformationRequiredExperienceDevelopmentType{}

				langType := entity.GetIntDevelopmentType(0, column)
				developmentType.DevelopmentType = langType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 92:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 93:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

				//	開発OS  3つまで
			case 94, 95, 96:
				developmentType = entity.SendingJobInformationRequiredExperienceDevelopmentType{}

				osType := entity.GetIntDevelopmentType(1, column)
				developmentType.DevelopmentType = osType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 97:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 98:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

			//免許・資格
			case 99:
				license = entity.SendingJobInformationRequiredLicense{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					license.License = entity.GetIntLicenseType(str)

					if license.License.Valid {
						condition.RequiredLicenses = append(condition.RequiredLicenses, license)
					}
				}

				// ※複数選択"	語学
				// 語学タイプ 3つまで
				language = entity.SendingJobInformationRequiredLanguage{}
			case 100, 101, 102:
				languageType = entity.SendingJobInformationRequiredLanguageType{}

				languageType.LanguageType = entity.GetIntLanguageType(column)

				if languageType.LanguageType.Valid {
					language.LanguageTypes = append(language.LanguageTypes, languageType)
				}

			case 103:
				language.LanguageLevel = entity.GetIntLanguageLevel(column)

			//"TOEIC ※英語のみ"
			case 104:
				toeic, err := strconv.Atoi(column)
				if err != nil {
					language.Toeic = null.NewInt(0, false)
				} else {
					language.Toeic = null.NewInt(int64(toeic), true)
				}

				//	"TOEFL(i) ※英語のみ"
			case 105:
				toeflI, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflIBT = null.NewInt(0, false)
				} else {
					language.ToeflIBT = null.NewInt(int64(toeflI), true)
				}

				// "TOEFL(P) ※英語のみ"
			case 106:
				toeflP, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflPBT = null.NewInt(0, false)
				} else {
					language.ToeflPBT = null.NewInt(int64(toeflP), true)
				}

				// 全項目未入力時はスキップ
				if len(language.LanguageTypes) > 0 ||
					language.LanguageLevel.Valid ||
					language.Toeic.Valid ||
					language.ToeflIBT.Valid ||
					language.ToeflPBT.Valid {

					condition.RequiredLanguages = language
				}

				// 共通項目
				if len(condition.RequiredExperienceDevelopments) > 0 ||
					len(condition.RequiredPCTools) > 0 ||
					condition.RequiredManagement.Valid ||
					condition.RequiredExperienceJobs.ExperienceYear.Valid ||
					condition.RequiredExperienceJobs.ExperienceMonth.Valid ||
					len(condition.RequiredExperienceJobs.ExperienceIndustries) > 0 ||
					len(condition.RequiredExperienceJobs.ExperienceOccupations) > 0 ||
					len(condition.RequiredLicenses) > 0 ||
					len(condition.RequiredLanguages.LanguageTypes) > 0 {
					sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, condition)
				}

				/***********   パターン別必須条件（業界・職種・資格・語学力）  ***********/

				//必要業職種経験1
				condition = entity.SendingJobInformationRequiredCondition{}
				condition.IsCommon = false
				jobExperience = entity.SendingJobInformationRequiredExperienceJob{}

				//業界
			case 107:
				experienceIndustry = entity.SendingJobInformationRequiredExperienceIndustry{}

				if column == "" {
					continue
				}

				industries := entity.GetIntIndustryFromBigCategory(column)

				for _, industry := range industries {
					experienceIndustry.ExperienceIndustry = null.NewInt(int64(industry), true)
					if experienceIndustry.ExperienceIndustry.Valid {
						jobExperience.ExperienceIndustries = append(jobExperience.ExperienceIndustries, experienceIndustry)
					}
				}

				//職種 3つまで
			case 108, 109, 110:
				experienceOccupation = entity.SendingJobInformationRequiredExperienceOccupation{}

				if column == "" {
					continue
				}

				experienceOccupation.ExperienceOccupation = entity.GetIntOccupation(column)

				if experienceOccupation.ExperienceOccupation.Valid {
					jobExperience.ExperienceOccupations = append(jobExperience.ExperienceOccupations, experienceOccupation)
				}

				//経験年数
			case 111:
				jobExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceYear = null.NewInt(int64(jobExperienceYear), true)
				}

				//経験年数
			case 112:
				jobExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceMonth = null.NewInt(int64(jobExperienceMonth), true)
				}

				// // 主要項目が未入力の場合は追加しない
				if len(jobExperience.ExperienceIndustries) > 0 ||
					len(jobExperience.ExperienceOccupations) > 0 ||
					jobExperience.ExperienceYear.Valid ||
					jobExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceJobs = jobExperience
				}

				// 必要マネジメント経験 "必要"という文字列で入ってくる。必要=0, true, 不要=0, false
			case 113:
				if column == "必要" {
					condition.RequiredManagement = null.NewInt(0, true)
				} else {
					condition.RequiredManagement = null.NewInt(0, false)
				}

			//PC業務ツール　3つまで
			case 114, 115, 116:
				pcTool = entity.SendingJobInformationRequiredPCTool{}
				if column == "" {
					continue
				} else {
					pcTool.Tool = entity.GetIntPCTool(column)

					if pcTool.Tool.Valid {
						condition.RequiredPCTools = append(condition.RequiredPCTools, pcTool)
					}
				}

				developmentExperience = entity.SendingJobInformationRequiredExperienceDevelopment{}

				/************* 開発経験 ****************/
				//	開発言語 3つまで
			case 117, 118, 119:
				developmentType = entity.SendingJobInformationRequiredExperienceDevelopmentType{}

				langType := entity.GetIntDevelopmentType(0, column)
				developmentType.DevelopmentType = langType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 120:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 121:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

				//	開発OS  3つまで
			case 122, 123, 124:
				developmentType = entity.SendingJobInformationRequiredExperienceDevelopmentType{}

				osType := entity.GetIntDevelopmentType(1, column)
				developmentType.DevelopmentType = osType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 125:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 126:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

			//免許・資格
			case 127:
				license = entity.SendingJobInformationRequiredLicense{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					license.License = entity.GetIntLicenseType(str)

					if license.License.Valid {
						condition.RequiredLicenses = append(condition.RequiredLicenses, license)
					}
				}

				// ※複数選択"	語学
				// 語学タイプ 3つまで
				language = entity.SendingJobInformationRequiredLanguage{}
			case 128, 129, 130:
				languageType = entity.SendingJobInformationRequiredLanguageType{}

				languageType.LanguageType = entity.GetIntLanguageType(column)

				if languageType.LanguageType.Valid {
					language.LanguageTypes = append(language.LanguageTypes, languageType)
				}

			case 131:
				language.LanguageLevel = entity.GetIntLanguageLevel(column)

			//"TOEIC ※英語のみ"
			case 132:
				toeic, err := strconv.Atoi(column)
				if err != nil {
					language.Toeic = null.NewInt(0, false)
				} else {
					language.Toeic = null.NewInt(int64(toeic), true)
				}

				//	"TOEFL(i) ※英語のみ"
			case 133:
				toeflI, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflIBT = null.NewInt(0, false)
				} else {
					language.ToeflIBT = null.NewInt(int64(toeflI), true)
				}

				// "TOEFL(P) ※英語のみ"
			case 134:
				toeflP, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflPBT = null.NewInt(0, false)
				} else {
					language.ToeflPBT = null.NewInt(int64(toeflP), true)
				}

				// 全項目未入力時はスキップ
				if len(language.LanguageTypes) > 0 ||
					language.LanguageLevel.Valid ||
					language.Toeic.Valid ||
					language.ToeflIBT.Valid ||
					language.ToeflPBT.Valid {

					condition.RequiredLanguages = language
				}

				// 共通項目
				if len(condition.RequiredExperienceDevelopments) > 0 ||
					len(condition.RequiredPCTools) > 0 ||
					condition.RequiredManagement.Valid ||
					condition.RequiredExperienceJobs.ExperienceYear.Valid ||
					condition.RequiredExperienceJobs.ExperienceMonth.Valid ||
					len(condition.RequiredExperienceJobs.ExperienceIndustries) > 0 ||
					len(condition.RequiredExperienceJobs.ExperienceOccupations) > 0 ||
					len(condition.RequiredLicenses) > 0 ||
					len(condition.RequiredLanguages.LanguageTypes) > 0 {
					sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, condition)
				}

				/***********   パターン別必須条件（業界・職種・資格・語学力）  ***********/

				//必要業職種経験1
				condition = entity.SendingJobInformationRequiredCondition{}
				condition.IsCommon = false
				jobExperience = entity.SendingJobInformationRequiredExperienceJob{}

				//業界
			case 135:
				experienceIndustry = entity.SendingJobInformationRequiredExperienceIndustry{}

				if column == "" {
					continue
				}

				industries := entity.GetIntIndustryFromBigCategory(column)

				for _, industry := range industries {
					experienceIndustry.ExperienceIndustry = null.NewInt(int64(industry), true)
					if experienceIndustry.ExperienceIndustry.Valid {
						jobExperience.ExperienceIndustries = append(jobExperience.ExperienceIndustries, experienceIndustry)
					}
				}

				//職種 3つまで
			case 136, 137, 138:
				experienceOccupation = entity.SendingJobInformationRequiredExperienceOccupation{}

				if column == "" {
					continue
				}

				experienceOccupation.ExperienceOccupation = entity.GetIntOccupation(column)

				if experienceOccupation.ExperienceOccupation.Valid {
					jobExperience.ExperienceOccupations = append(jobExperience.ExperienceOccupations, experienceOccupation)
				}

				//経験年数
			case 139:
				jobExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceYear = null.NewInt(int64(jobExperienceYear), true)
				}

				//経験年数
			case 140:
				jobExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					jobExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					jobExperience.ExperienceMonth = null.NewInt(int64(jobExperienceMonth), true)
				}

				// // 主要項目が未入力の場合は追加しない
				if len(jobExperience.ExperienceIndustries) > 0 ||
					len(jobExperience.ExperienceOccupations) > 0 ||
					jobExperience.ExperienceYear.Valid ||
					jobExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceJobs = jobExperience
				}

			// 必要マネジメント経験 "必要"という文字列で入ってくる。必要=0, true, 不要=0, false
			case 141:
				if column == "必要" {
					condition.RequiredManagement = null.NewInt(0, true)
				} else {
					condition.RequiredManagement = null.NewInt(0, false)
				}

			//PC業務ツール　3つまで
			case 142, 143, 144:
				pcTool = entity.SendingJobInformationRequiredPCTool{}
				if column == "" {
					continue
				} else {
					pcTool.Tool = entity.GetIntPCTool(column)

					if pcTool.Tool.Valid {
						condition.RequiredPCTools = append(condition.RequiredPCTools, pcTool)
					}
				}

				developmentExperience = entity.SendingJobInformationRequiredExperienceDevelopment{}

				/************* 開発経験 ****************/
				//	開発言語 3つまで
			case 145, 146, 147:
				developmentType = entity.SendingJobInformationRequiredExperienceDevelopmentType{}

				langType := entity.GetIntDevelopmentType(0, column)
				developmentType.DevelopmentType = langType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 148:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 149:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

				//	開発OS  3つまで
			case 150, 151, 152:
				developmentType = entity.SendingJobInformationRequiredExperienceDevelopmentType{}

				osType := entity.GetIntDevelopmentType(1, column)
				developmentType.DevelopmentType = osType

				if developmentType.DevelopmentType.Valid {
					developmentExperience.ExperienceDevelopmentTypes = append(developmentExperience.ExperienceDevelopmentTypes, developmentType)
				}

				//経験年数（年）
			case 153:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 154:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentExperience.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentExperience.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if len(developmentExperience.ExperienceDevelopmentTypes) > 0 ||
					developmentExperience.ExperienceYear.Valid ||
					developmentExperience.ExperienceMonth.Valid {
					condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, developmentExperience)
				}

			//免許・資格
			case 155:
				license = entity.SendingJobInformationRequiredLicense{}

				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {

					str = strings.TrimSpace(str)
					license.License = entity.GetIntLicenseType(str)

					if license.License.Valid {
						condition.RequiredLicenses = append(condition.RequiredLicenses, license)
					}
				}

				// ※複数選択"	語学
				// 語学タイプ 3つまで
				language = entity.SendingJobInformationRequiredLanguage{}
			case 156, 157, 158:
				languageType = entity.SendingJobInformationRequiredLanguageType{}

				languageType.LanguageType = entity.GetIntLanguageType(column)

				if languageType.LanguageType.Valid {
					language.LanguageTypes = append(language.LanguageTypes, languageType)
				}

			case 159:
				language.LanguageLevel = entity.GetIntLanguageLevel(column)

			//"TOEIC ※英語のみ"
			case 160:
				toeic, err := strconv.Atoi(column)
				if err != nil {
					language.Toeic = null.NewInt(0, false)
				} else {
					language.Toeic = null.NewInt(int64(toeic), true)
				}

				//	"TOEFL(i) ※英語のみ"
			case 161:
				toeflI, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflIBT = null.NewInt(0, false)
				} else {
					language.ToeflIBT = null.NewInt(int64(toeflI), true)
				}

				// "TOEFL(P) ※英語のみ"
			case 162:
				toeflP, err := strconv.Atoi(column)
				if err != nil {
					language.ToeflPBT = null.NewInt(0, false)
				} else {
					language.ToeflPBT = null.NewInt(int64(toeflP), true)
				}

				// 全項目未入力時はスキップ
				if len(language.LanguageTypes) > 0 ||
					language.LanguageLevel.Valid ||
					language.Toeic.Valid ||
					language.ToeflIBT.Valid ||
					language.ToeflPBT.Valid {

					condition.RequiredLanguages = language
				}

				// 共通項目
				if len(condition.RequiredExperienceDevelopments) > 0 ||
					len(condition.RequiredPCTools) > 0 ||
					condition.RequiredManagement.Valid ||
					condition.RequiredExperienceJobs.ExperienceYear.Valid ||
					condition.RequiredExperienceJobs.ExperienceMonth.Valid ||
					len(condition.RequiredExperienceJobs.ExperienceIndustries) > 0 ||
					len(condition.RequiredExperienceJobs.ExperienceOccupations) > 0 ||
					len(condition.RequiredLicenses) > 0 ||
					len(condition.RequiredLanguages.LanguageTypes) > 0 {
					sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, condition)
				}
			}
		}

		sendingJobInformationList = append(sendingJobInformationList, &sendingJobInformation)
	}

	for _, sendingJobInformation := range sendingJobInformationList {
		for _, condition := range sendingJobInformation.RequiredConditions {
			fmt.Println("licnese: ", condition.RequiredLicenses)
			fmt.Println("language: ", condition.RequiredLanguages)
			fmt.Println("experienceJob: ", condition.RequiredExperienceJobs)
			fmt.Println("experienceIndustry: ", condition.RequiredExperienceJobs.ExperienceIndustries)
			fmt.Println("experienceDevelopment: ", condition.RequiredExperienceDevelopments)
		}
	}

	fmt.Println("除外されたレコードは: ", missedRecords)

	return sendingJobInformationList, missedRecords, nil
}

/***************** CSVファイルを読み込む（求職者用）） *********************/
func parseJobSeekerCSV(r *csv.Reader) ([]*entity.JobSeeker, error) {
	var (
		jobSeekerList []*entity.JobSeeker
	)

	// カラムタイトルと項目説明のレコードを読み込み
	r.Read()
	r.Read()
	r.Read()

	for {
		var (
			jobSeeker        entity.JobSeeker
			license          entity.JobSeekerLicense
			pcTool           entity.JobSeekerPCTool
			developmentSkill entity.JobSeekerDevelopmentSkill
			studentHistory   entity.JobSeekerStudentHistory

			languageSkill1 entity.JobSeekerLanguageSkill
			languageSkill2 entity.JobSeekerLanguageSkill
			languageSkill3 entity.JobSeekerLanguageSkill

			workHistory1 entity.JobSeekerWorkHistory
			workHistory2 entity.JobSeekerWorkHistory
			workHistory3 entity.JobSeekerWorkHistory
			workHistory4 entity.JobSeekerWorkHistory
			workHistory5 entity.JobSeekerWorkHistory

			experienceIndustry entity.JobSeekerExperienceIndustry

			departmentHistory1  entity.JobSeekerDepartmentHistory
			departmentHistory2  entity.JobSeekerDepartmentHistory
			departmentHistory3  entity.JobSeekerDepartmentHistory
			departmentHistory4  entity.JobSeekerDepartmentHistory
			departmentHistory5  entity.JobSeekerDepartmentHistory
			departmentHistory6  entity.JobSeekerDepartmentHistory
			departmentHistory7  entity.JobSeekerDepartmentHistory
			departmentHistory8  entity.JobSeekerDepartmentHistory
			departmentHistory9  entity.JobSeekerDepartmentHistory
			departmentHistory10 entity.JobSeekerDepartmentHistory
			departmentHistory11 entity.JobSeekerDepartmentHistory
			departmentHistory12 entity.JobSeekerDepartmentHistory
			departmentHistory13 entity.JobSeekerDepartmentHistory
			departmentHistory14 entity.JobSeekerDepartmentHistory
			departmentHistory15 entity.JobSeekerDepartmentHistory

			experienceOccupation1  entity.JobSeekerExperienceOccupation
			experienceOccupation2  entity.JobSeekerExperienceOccupation
			experienceOccupation3  entity.JobSeekerExperienceOccupation
			experienceOccupation4  entity.JobSeekerExperienceOccupation
			experienceOccupation5  entity.JobSeekerExperienceOccupation
			experienceOccupation6  entity.JobSeekerExperienceOccupation
			experienceOccupation7  entity.JobSeekerExperienceOccupation
			experienceOccupation8  entity.JobSeekerExperienceOccupation
			experienceOccupation9  entity.JobSeekerExperienceOccupation
			experienceOccupation10 entity.JobSeekerExperienceOccupation
			experienceOccupation11 entity.JobSeekerExperienceOccupation
			experienceOccupation12 entity.JobSeekerExperienceOccupation
			experienceOccupation13 entity.JobSeekerExperienceOccupation
			experienceOccupation14 entity.JobSeekerExperienceOccupation
			experienceOccupation15 entity.JobSeekerExperienceOccupation

			selfPromotion       entity.JobSeekerSelfPromotion
			desiredIndustry     entity.JobSeekerDesiredIndustry
			desiredOccupation   entity.JobSeekerDesiredOccupation
			desiredWorkLocation entity.JobSeekerDesiredWorkLocation
			desiredHolidayType  entity.JobSeekerDesiredHolidayType
			desiredCompanyScale entity.JobSeekerDesiredCompanyScale
		)

		// workHistory = entity.JobSeekerWorkHistory{}

		record, err := r.Read()
		// 空行があった場合は終了
		if err == io.EOF || len(record) == 0 || record[6] == "" {
			break
		}

		/**
		求職者項目
		*/

		for i, column := range record {
			switch i {

			// 営業担当者ID
			case 0:
				id, err := strconv.Atoi(column)
				if err != nil {
					jobSeeker.AgentStaffID = null.NewInt(0, false)
				} else {
					jobSeeker.AgentStaffID = null.NewInt(int64(id), true)
				}

				// 求職者のフェーズ
			case 1:
				jobSeeker.Phase = entity.GetIntPhaseForJobSeeker(column)

			// 流入経路ID
			case 2:
				id, err := strconv.Atoi(column)
				if err != nil {
					jobSeeker.InflowChannelID = null.NewInt(0, false)
				} else {
					jobSeeker.InflowChannelID = null.NewInt(int64(id), true)
				}

			// エントリー日時
			case 3:
				entryDate, err := time.Parse("2006-01-02 15:04", column)
				if err != nil {
					jobSeeker.CreatedAt = time.Now().UTC()
				} else {
					// エントリー日時はUTCに変換
					jobSeeker.CreatedAt = entryDate.Add(-9 * time.Hour)
				}

				// 面談日時
			case 4:
				interviewDate, err := time.Parse("2006-01-02 15:04", column)
				if err != nil {
					jobSeeker.InterviewDate = time.Now().UTC()
				} else {
					// UTCに変換
					jobSeeker.InterviewDate = interviewDate.Add(-9 * time.Hour)
				}

				// ステータス
			case 5:
				jobSeeker.UserStatus = entity.GetIntUserStatus(column)

			// 苗字
			case 6:
				jobSeeker.LastName = column

				// 名前
			case 7:
				jobSeeker.FirstName = column

			case 8:
				jobSeeker.LastFurigana = column

			case 9:
				jobSeeker.FirstFurigana = column

				// 性別
			case 10:
				jobSeeker.Gender = entity.GetIntGenderForJobSeeker(column)

				// 性別備考
			case 11:
				jobSeeker.GenderRemarks = column

			//国籍
			case 12:
				jobSeeker.Nationality = entity.GetIntNationalityForJobSeeker(column)

				//国籍備考
			case 13:
				jobSeeker.NationalityRemarks = column

				//既往歴
			case 14:
				jobSeeker.MedicalHistory = entity.GetIntAvailable(column)

				//既往歴備考
			case 15:
				jobSeeker.MedicalHistoryRemarks = column

				// 生年月日
			case 16:
				jobSeeker.Birthday = column

				// 配偶者の有無
			case 17:
				jobSeeker.Spouse = entity.GetIntAvailable(column)

				// 配偶者の扶養義務
			case 18:
				jobSeeker.SupportObligation = entity.GetIntAvailable(column)

				// 扶養者数
			case 19:
				val, err := strconv.Atoi(column)
				if err != nil {
					jobSeeker.Dependents = null.NewInt(0, false)
				} else {
					jobSeeker.Dependents = null.NewInt(int64(val), true)
				}

				// 電話番号
			case 20:
				jobSeeker.PhoneNumber = column

				// メールアドレス
			case 21:
				jobSeeker.Email = column

				// 緊急連絡先
			case 22:
				jobSeeker.EmergencyPhoneNumber = column

				// 郵便番号
			case 23:
				jobSeeker.PostCode = column

			case 24:
				//住所(都道府県)
				jobSeeker.Prefecture = entity.GetIntPrefecture(column)

				// 住所（市町村以下）
			case 25:
				jobSeeker.Address = column

				// 住所（フリガナ）
			case 26:
				jobSeeker.AddressFurigana = column

				//直近の年収
			case 27:
				val, err := strconv.Atoi(column)
				if err != nil {
					jobSeeker.AnnualIncome = null.NewInt(0, false)
				} else {
					jobSeeker.AnnualIncome = null.NewInt(int64(val), true)
				}

				/************学歴情報************/

			// 文系/理系
			case 28:
				jobSeeker.StudyCategory = entity.GetIntStudyCategoryForJobSeeker(column)

				/************学歴1************/
				//学校フェーズ
			case 29:
				studentHistory.SchoolCategory = entity.GetIntSchoolCategoryForJobSeeker(column)

				//学校名
			case 30:
				studentHistory.SchoolName = column

				//大学レベル ※大卒なら
			case 31:
				studentHistory.SchoolLevel = entity.GetIntSchoolLevelForJobSeeker(column)

				// 学部・学科・コース
			case 32:
				studentHistory.Subject = column

				//入学年月
			case 33:
				studentHistory.EntranceYear = column

			// 開始ステータス
			case 34:
				studentHistory.FirstStatus = entity.GetIntFirstStatusForStudentHistory(column)

				//卒業年月
			case 35:
				studentHistory.GraduationYear = column

				//終了ステータス
			case 36:
				studentHistory.LastStatus = entity.GetIntLastStatusForStudentHistory(column)

				// 学歴情報を追加
				// 学校フェーズが空白でない場合のみ追加
				if studentHistory.SchoolCategory.Valid {
					jobSeeker.StudentHistories = append(jobSeeker.StudentHistories, studentHistory)
				}

				/************学歴2************/
				//学校フェーズ
			case 37:
				studentHistory.SchoolCategory = entity.GetIntSchoolCategoryForJobSeeker(column)

				//学校名
			case 38:
				studentHistory.SchoolName = column

				//大学レベル ※大卒なら
			case 39:
				studentHistory.SchoolLevel = entity.GetIntSchoolLevelForJobSeeker(column)

				// 学部・学科・コース
			case 40:
				studentHistory.Subject = column

				//入学年月
			case 41:
				studentHistory.EntranceYear = column

			// 開始ステータス
			case 42:
				studentHistory.FirstStatus = entity.GetIntFirstStatusForStudentHistory(column)

				//卒業年月
			case 43:
				studentHistory.GraduationYear = column

				//終了ステータス
			case 44:
				studentHistory.LastStatus = entity.GetIntLastStatusForStudentHistory(column)

				// 学歴情報を追加
				// 学校フェーズが空白でない場合のみ追加
				if studentHistory.SchoolCategory.Valid {
					jobSeeker.StudentHistories = append(jobSeeker.StudentHistories, studentHistory)
				}

				/************学歴3************/
				//学校フェーズ
			case 45:
				studentHistory.SchoolCategory = entity.GetIntSchoolCategoryForJobSeeker(column)

				//学校名
			case 46:
				studentHistory.SchoolName = column

				//大学レベル ※大卒なら
			case 47:
				studentHistory.SchoolLevel = entity.GetIntSchoolLevelForJobSeeker(column)

				// 学部・学科・コース
			case 48:
				studentHistory.Subject = column

				//入学年月
			case 49:
				studentHistory.EntranceYear = column

			// 開始ステータス
			case 50:
				studentHistory.FirstStatus = entity.GetIntFirstStatusForStudentHistory(column)

				//卒業年月
			case 51:
				studentHistory.GraduationYear = column

				//終了ステータス
			case 52:
				studentHistory.LastStatus = entity.GetIntLastStatusForStudentHistory(column)

				// 学歴情報を追加
				// 学校フェーズが空白でない場合のみ追加
				if studentHistory.SchoolCategory.Valid {
					jobSeeker.StudentHistories = append(jobSeeker.StudentHistories, studentHistory)
				}

				/***************職歴情報***************/

				//就業状況
			case 53:
				jobSeeker.StateOfEmployment = entity.GetIntStateOfEmployment(column)

			//転職回数
			case 54:
				jobChange, err := strconv.Atoi(column)
				if err != nil {
					jobSeeker.JobChange = null.NewInt(0, false)
				} else {
					// 5回以上の場合は5
					if jobChange > 4 {
						jobChange = 5
					}
					jobSeeker.JobChange = null.NewInt(int64(jobChange), true)
				}

				//	短期離職（1年未満）	の有無
			case 55:
				jobSeeker.ShortResignation = entity.GetIntAvailable(column)

				//短期離職備考
			case 56:
				jobSeeker.ShortResignationRemarks = column

			//職務経歴要約
			case 57:
				jobSeeker.JobSummary = column

				// 経歴補足
			case 58:
				jobSeeker.HistorySupplement = column

				/***************職歴1***************/
				//会社名
			case 59:
				workHistory1.CompanyName = column

				//業界 ※≒経験業界 3つ
			case 60, 61, 62:
				experienceIndustry.Industry = entity.GetIntIndustry(column)

				if experienceIndustry.Industry.Valid {
					workHistory1.ExperienceIndustries = append(workHistory1.ExperienceIndustries, experienceIndustry)
				}

				// 従業員数（単体）
			case 63:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory1.EmployeeNumberSingle = null.NewInt(0, false)
				} else {
					workHistory1.EmployeeNumberSingle = null.NewInt(int64(val), true)
				}

				//従業員数（連結）
			case 64:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory1.EmployeeNumberGroup = null.NewInt(0, false)
				} else {
					workHistory1.EmployeeNumberGroup = null.NewInt(int64(val), true)
				}

				//株式公開
			case 65:
				workHistory1.PublicOffering = entity.GetIntPublicOffering(column)

				// 雇用形態
			case 66:
				workHistory1.EmploymentStatus = entity.GetIntEmploymentStatus(column)

			// 入社年 *2020-12 を取得
			case 67:
				workHistory1.JoiningYear = column

			// 開始ステータス
			case 68:
				workHistory1.FirstStatus = entity.GetIntFirstStatusForWorkHistory(column)

				//退社年* 2020-12 を取得
			case 69:
				workHistory1.RetireYear = column

			//終了ステータス
			case 70:
				workHistory1.LastStatus = entity.GetIntLastStatusForWorkHistory(column)

				// 退職理由（本音）
			case 71:
				workHistory1.RetireReasonOfTruth = column

				//退職理由（建前）
			case 72:
				workHistory1.RetireReasonOfPublic = column

				// 職種・役職歴1
				//職種 3つ
			case 73, 74, 75:
				experienceOccupation1.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation1.Occupation.Valid {
					departmentHistory1.ExperienceOccupations = append(departmentHistory1.ExperienceOccupations, experienceOccupation1)
				}

				//部門名
			case 76:
				departmentHistory1.Department = column

				//マネジメント人数
			case 77:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory1.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory1.ManagementNumber = null.NewInt(int64(v), true)
				}

				// マネジメント経験詳細
			case 78:
				departmentHistory1.ManagementDetail = column

				//職務内容
			case 79:
				departmentHistory1.JobDescription = column

				//開始年月* 2020-12 を取得
			case 80:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory1.StartYear = column

				// 終了年月* 2020-12 を取得
			case 81:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory1.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory1.Department != "" ||
					departmentHistory1.ManagementNumber.Valid || departmentHistory1.ManagementDetail != "" ||
					departmentHistory1.StartYear != "" || departmentHistory1.EndYear != "" ||
					departmentHistory1.JobDescription != "" {
					workHistory1.DepartmentHistories = append(workHistory1.DepartmentHistories, departmentHistory1)
				}

			// 職種・役職歴2
			//職種 3つ
			case 82, 83, 84:
				experienceOccupation2.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation2.Occupation.Valid {
					departmentHistory2.ExperienceOccupations = append(departmentHistory2.ExperienceOccupations, experienceOccupation2)
				}

				//部門名
			case 85:
				departmentHistory2.Department = column

				//マネジメント人数
			case 86:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory2.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory2.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 87:
				// マネジメント経験詳細
				departmentHistory2.ManagementDetail = column

			//職務内容
			case 88:
				departmentHistory2.JobDescription = column

			// 開始年月 * 2020-12 を取得
			case 89:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory2.StartYear = column

			// 終了年月 * 2020-12 を取得
			case 90:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory2.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory2.Department != "" ||
					departmentHistory2.ManagementNumber.Valid || departmentHistory2.ManagementDetail != "" ||
					departmentHistory2.StartYear != "" || departmentHistory2.EndYear != "" ||
					departmentHistory2.JobDescription != "" {
					workHistory1.DepartmentHistories = append(workHistory1.DepartmentHistories, departmentHistory2)
				}

				// 職種・役職歴3
				//職種 3つ
			case 91, 92, 93:
				experienceOccupation3.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation3.Occupation.Valid {
					departmentHistory3.ExperienceOccupations = append(departmentHistory3.ExperienceOccupations, experienceOccupation3)
				}

				//部門名
			case 94:
				departmentHistory3.Department = column

				//マネジメント人数
			case 95:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory3.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory3.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 96:
				// マネジメント経験詳細
				departmentHistory3.ManagementDetail = column

			//職務内容
			case 97:
				departmentHistory3.JobDescription = column

			// //開始年月
			// 2020-12 を取得
			case 98:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory3.StartYear = column

			// 	//終了年月
			// 	2020-12 を取得
			case 99:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory3.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory3.Department != "" ||
					departmentHistory3.ManagementNumber.Valid || departmentHistory3.ManagementDetail != "" ||
					departmentHistory3.StartYear != "" || departmentHistory3.EndYear != "" ||
					departmentHistory3.JobDescription != "" {
					workHistory1.DepartmentHistories = append(workHistory1.DepartmentHistories, departmentHistory3)
				}

				if workHistory1.CompanyName != "" {
					jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, workHistory1)
				}

				/***************職歴2***************/
				//会社名
			case 100:
				workHistory2.CompanyName = column

				//業界 ※≒経験業界
			case 101, 102, 103:
				experienceIndustry.Industry = entity.GetIntIndustry(column)

				if experienceIndustry.Industry.Valid {
					workHistory2.ExperienceIndustries = append(workHistory2.ExperienceIndustries, experienceIndustry)
				}

				// 従業員数（単体）
			case 104:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory2.EmployeeNumberSingle = null.NewInt(0, false)
				} else {
					workHistory2.EmployeeNumberSingle = null.NewInt(int64(val), true)
				}

				//従業員数（連結）
			case 105:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory2.EmployeeNumberGroup = null.NewInt(0, false)
				} else {
					workHistory2.EmployeeNumberGroup = null.NewInt(int64(val), true)
				}

				//株式公開
			case 106:
				workHistory2.PublicOffering = entity.GetIntPublicOffering(column)

				// 雇用形態
			case 107:
				workHistory2.EmploymentStatus = entity.GetIntEmploymentStatus(column)

			// 入社年
			case 108:
				workHistory2.JoiningYear = column

			// 開始ステータス
			case 109:
				workHistory2.FirstStatus = entity.GetIntFirstStatusForWorkHistory(column)

				//退社年
			case 110:
				workHistory2.RetireYear = column

			//終了ステータス
			case 111:
				workHistory2.LastStatus = entity.GetIntLastStatusForWorkHistory(column)

				// 退職理由（本音）
			case 112:
				workHistory2.RetireReasonOfTruth = column

				//退職理由（建前）
			case 113:
				workHistory2.RetireReasonOfPublic = column

				// 職種・役職歴1
				//職種 3つ
			case 114, 115, 116:
				experienceOccupation4.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation4.Occupation.Valid {
					departmentHistory4.ExperienceOccupations = append(departmentHistory4.ExperienceOccupations, experienceOccupation4)
				}

				//部門名
			case 117:
				departmentHistory4.Department = column

				//マネジメント人数
			case 118:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory4.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory4.ManagementNumber = null.NewInt(int64(v), true)
				}

				// マネジメント経験詳細
			case 119:
				departmentHistory4.ManagementDetail = column

				//職務内容
			case 120:
				departmentHistory4.JobDescription = column

				// 開始年月
			case 121:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory4.StartYear = column

				// 終了年月
			case 122:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory4.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory4.Department != "" ||
					departmentHistory4.ManagementNumber.Valid || departmentHistory4.ManagementDetail != "" ||
					departmentHistory4.StartYear != "" || departmentHistory4.EndYear != "" ||
					departmentHistory4.JobDescription != "" {
					workHistory2.DepartmentHistories = append(workHistory2.DepartmentHistories, departmentHistory4)
				}

			// 職種・役職歴2
			//職種
			case 123, 124, 125:
				experienceOccupation5.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation5.Occupation.Valid {
					departmentHistory5.ExperienceOccupations = append(departmentHistory5.ExperienceOccupations, experienceOccupation5)
				}

				//部門名
			case 126:
				departmentHistory5.Department = column

				//マネジメント人数
			case 127:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory5.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory5.ManagementNumber = null.NewInt(int64(v), true)
				}

				// マネジメント経験詳細
			case 128:
				departmentHistory5.ManagementDetail = column

			//職務内容
			case 129:
				departmentHistory5.JobDescription = column

			// 開始年月
			case 130:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory5.StartYear = column

			// 終了年月
			case 131:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory5.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory5.Department != "" ||
					departmentHistory5.ManagementNumber.Valid || departmentHistory5.ManagementDetail != "" ||
					departmentHistory5.StartYear != "" || departmentHistory5.EndYear != "" ||
					departmentHistory5.JobDescription != "" {
					workHistory2.DepartmentHistories = append(workHistory2.DepartmentHistories, departmentHistory5)
				}

				// 職種・役職歴3
				//職種
			case 132, 133, 134:
				experienceOccupation6.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation6.Occupation.Valid {
					departmentHistory6.ExperienceOccupations = append(departmentHistory6.ExperienceOccupations, experienceOccupation6)
				}

				//部門名
			case 135:
				departmentHistory6.Department = column

				//マネジメント人数
			case 136:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory6.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory6.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 137:
				// マネジメント経験詳細
				departmentHistory6.ManagementDetail = column

			//職務内容
			case 138:
				departmentHistory6.JobDescription = column

			// 開始年月
			case 139:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory6.StartYear = column

			// 終了年月
			case 140:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory6.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory6.Department != "" ||
					departmentHistory6.ManagementNumber.Valid || departmentHistory6.ManagementDetail != "" ||
					departmentHistory6.StartYear != "" || departmentHistory6.EndYear != "" ||
					departmentHistory6.JobDescription != "" {
					workHistory2.DepartmentHistories = append(workHistory2.DepartmentHistories, departmentHistory6)
				}

				if workHistory2.CompanyName != "" {
					jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, workHistory2)
				}
				/***************職歴3***************/
				//会社名
			case 141:
				workHistory3.CompanyName = column

				//業界 ※≒経験業界
			case 142, 143, 144:
				experienceIndustry.Industry = entity.GetIntIndustry(column)

				if experienceIndustry.Industry.Valid {
					workHistory3.ExperienceIndustries = append(workHistory3.ExperienceIndustries, experienceIndustry)
				}

				// 従業員数（単体）
			case 145:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory3.EmployeeNumberSingle = null.NewInt(0, false)
				} else {
					workHistory3.EmployeeNumberSingle = null.NewInt(int64(val), true)
				}

				//従業員数（連結）
			case 146:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory3.EmployeeNumberGroup = null.NewInt(0, false)
				} else {
					workHistory3.EmployeeNumberGroup = null.NewInt(int64(val), true)
				}

				//株式公開
			case 147:
				workHistory3.PublicOffering = entity.GetIntPublicOffering(column)

				// 雇用形態
			case 148:
				workHistory3.EmploymentStatus = entity.GetIntEmploymentStatus(column)

			// 入社年
			case 149:
				workHistory3.JoiningYear = column

			// 開始ステータス
			case 150:
				workHistory3.FirstStatus = entity.GetIntFirstStatusForWorkHistory(column)

				//退社年
			case 151:
				workHistory3.RetireYear = column

			//終了ステータス
			case 152:
				workHistory3.LastStatus = entity.GetIntLastStatusForWorkHistory(column)

				// 退職理由（本音）
			case 153:
				workHistory3.RetireReasonOfTruth = column

				//退職理由（建前）
			case 154:
				workHistory3.RetireReasonOfPublic = column

				// 職種・役職歴1
				//職種
			case 155, 156, 157:
				experienceOccupation7.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation7.Occupation.Valid {
					departmentHistory7.ExperienceOccupations = append(departmentHistory7.ExperienceOccupations, experienceOccupation7)
				}

				//部門名
			case 158:
				departmentHistory7.Department = column

				//マネジメント人数
			case 159:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory7.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory7.ManagementNumber = null.NewInt(int64(v), true)
				}

				// マネジメント経験詳細
			case 160:
				departmentHistory7.ManagementDetail = column

				//職務内容
			case 161:
				departmentHistory7.JobDescription = column

				// 開始年月
			case 162:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory7.StartYear = column

				// 終了年月
			case 163:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory7.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory7.Department != "" ||
					departmentHistory7.ManagementNumber.Valid || departmentHistory7.ManagementDetail != "" ||
					departmentHistory7.StartYear != "" || departmentHistory7.EndYear != "" ||
					departmentHistory7.JobDescription != "" {
					workHistory3.DepartmentHistories = append(workHistory3.DepartmentHistories, departmentHistory7)
				}

			// 職種・役職歴2
			//職種
			case 164, 165, 166:
				experienceOccupation8.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation8.Occupation.Valid {
					departmentHistory8.ExperienceOccupations = append(departmentHistory8.ExperienceOccupations, experienceOccupation8)
				}

				//部門名
			case 167:
				departmentHistory8.Department = column

				//マネジメント人数
			case 168:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory8.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory8.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 169:
				// マネジメント経験詳細
				departmentHistory8.ManagementDetail = column

			//職務内容
			case 170:
				departmentHistory8.JobDescription = column

			// 開始年月
			case 171:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory8.StartYear = column

			// 終了年月
			case 172:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory8.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory8.Department != "" ||
					departmentHistory8.ManagementNumber.Valid || departmentHistory8.ManagementDetail != "" ||
					departmentHistory8.StartYear != "" || departmentHistory8.EndYear != "" ||
					departmentHistory8.JobDescription != "" {
					workHistory3.DepartmentHistories = append(workHistory3.DepartmentHistories, departmentHistory8)
				}

				// 職種・役職歴3
				//職種
			case 173, 174, 175:
				experienceOccupation9.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation9.Occupation.Valid {
					departmentHistory9.ExperienceOccupations = append(departmentHistory9.ExperienceOccupations, experienceOccupation9)
				}

				//部門名
			case 176:
				departmentHistory9.Department = column

				//マネジメント人数
			case 177:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory9.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory9.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 178:
				// マネジメント経験詳細
				departmentHistory9.ManagementDetail = column

			//職務内容
			case 179:
				departmentHistory9.JobDescription = column

			//開始年月
			case 180:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory9.StartYear = column

			// 	終了年月
			case 181:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory9.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory9.Department != "" ||
					departmentHistory9.ManagementNumber.Valid || departmentHistory9.ManagementDetail != "" ||
					departmentHistory9.StartYear != "" || departmentHistory9.EndYear != "" ||
					departmentHistory9.JobDescription != "" {
					workHistory3.DepartmentHistories = append(workHistory3.DepartmentHistories, departmentHistory9)
				}

				if workHistory3.CompanyName != "" {
					jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, workHistory3)
				}
				/***************職歴4***************/
				//会社名
			case 182:
				workHistory4.CompanyName = column

				//業界 ※≒経験業界
			case 183, 184, 185:
				experienceIndustry.Industry = entity.GetIntIndustry(column)

				if experienceIndustry.Industry.Valid {
					workHistory4.ExperienceIndustries = append(workHistory4.ExperienceIndustries, experienceIndustry)
				}

				// 従業員数（単体）
			case 186:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory4.EmployeeNumberSingle = null.NewInt(0, false)
				} else {
					workHistory4.EmployeeNumberSingle = null.NewInt(int64(val), true)
				}

				//従業員数（連結）
			case 187:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory4.EmployeeNumberGroup = null.NewInt(0, false)
				} else {
					workHistory4.EmployeeNumberGroup = null.NewInt(int64(val), true)
				}

				//株式公開
			case 188:
				workHistory4.PublicOffering = entity.GetIntPublicOffering(column)

				// 雇用形態
			case 189:
				workHistory4.EmploymentStatus = entity.GetIntEmploymentStatus(column)

			// 入社年
			case 190:
				workHistory4.JoiningYear = column

			// 開始ステータス
			case 191:
				workHistory4.FirstStatus = entity.GetIntFirstStatusForWorkHistory(column)

				//退社年
			case 192:
				workHistory4.RetireYear = column

			//終了ステータス
			case 193:
				workHistory4.LastStatus = entity.GetIntLastStatusForWorkHistory(column)

				// 退職理由（本音）
			case 194:
				workHistory4.RetireReasonOfTruth = column

				//退職理由（建前）
			case 195:
				workHistory4.RetireReasonOfPublic = column

				// 職種・役職歴1
				//職種
			case 196, 197, 198:
				experienceOccupation10.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation10.Occupation.Valid {
					departmentHistory10.ExperienceOccupations = append(departmentHistory10.ExperienceOccupations, experienceOccupation10)
				}

				//部門名
			case 199:
				departmentHistory10.Department = column

				//マネジメント人数
			case 200:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory10.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory10.ManagementNumber = null.NewInt(int64(v), true)
				}

				// マネジメント経験詳細
			case 201:
				departmentHistory10.ManagementDetail = column

				//職務内容
			case 202:
				departmentHistory10.JobDescription = column

				// 開始年月
			case 203:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory10.StartYear = column

				// 終了年月
			case 204:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory10.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory10.Department != "" ||
					departmentHistory10.ManagementNumber.Valid || departmentHistory10.ManagementDetail != "" ||
					departmentHistory10.StartYear != "" || departmentHistory10.EndYear != "" ||
					departmentHistory10.JobDescription != "" {
					workHistory4.DepartmentHistories = append(workHistory4.DepartmentHistories, departmentHistory10)
				}

			// 職種・役職歴2
			//職種
			case 205, 206, 207:
				experienceOccupation11.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation11.Occupation.Valid {
					departmentHistory11.ExperienceOccupations = append(departmentHistory11.ExperienceOccupations, experienceOccupation11)
				}

				//部門名
			case 208:
				departmentHistory11.Department = column

				//マネジメント人数
			case 209:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory11.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory11.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 210:
				// マネジメント経験詳細
				departmentHistory11.ManagementDetail = column

			//職務内容
			case 211:
				departmentHistory11.JobDescription = column

			// 開始年月
			case 212:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory11.StartYear = column

			// 終了年月
			case 213:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory11.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory11.Department != "" ||
					departmentHistory11.ManagementNumber.Valid || departmentHistory11.ManagementDetail != "" ||
					departmentHistory11.StartYear != "" || departmentHistory11.EndYear != "" ||
					departmentHistory11.JobDescription != "" {
					workHistory4.DepartmentHistories = append(workHistory4.DepartmentHistories, departmentHistory11)
				}

				// 職種・役職歴3
				//職種
			case 214, 215, 216:
				experienceOccupation12.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation12.Occupation.Valid {
					departmentHistory12.ExperienceOccupations = append(departmentHistory12.ExperienceOccupations, experienceOccupation12)
				}

				//部門名
			case 217:
				departmentHistory12.Department = column

				//マネジメント人数
			case 218:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory12.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory12.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 219:
				// マネジメント経験詳細
				departmentHistory12.ManagementDetail = column

			//職務内容
			case 220:
				departmentHistory12.JobDescription = column

			// 開始年月
			case 221:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory12.StartYear = column

			// 終了年月
			case 222:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory12.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory12.Department != "" ||
					departmentHistory12.ManagementNumber.Valid || departmentHistory12.ManagementDetail != "" ||
					departmentHistory12.StartYear != "" || departmentHistory12.EndYear != "" ||
					departmentHistory12.JobDescription != "" {
					workHistory4.DepartmentHistories = append(workHistory4.DepartmentHistories, departmentHistory12)
				}

				if workHistory4.CompanyName != "" {
					jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, workHistory4)
				}
				/***************職歴5***************/
				//会社名
			case 223:
				workHistory5.CompanyName = column

				//業界 ※≒経験業界
			case 224, 225, 226:
				experienceIndustry.Industry = entity.GetIntIndustry(column)

				if experienceIndustry.Industry.Valid {
					workHistory5.ExperienceIndustries = append(workHistory5.ExperienceIndustries, experienceIndustry)
				}

				// 従業員数（単体）
			case 227:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory5.EmployeeNumberSingle = null.NewInt(0, false)
				} else {
					workHistory5.EmployeeNumberSingle = null.NewInt(int64(val), true)
				}

				//従業員数（連結）
			case 228:
				val, err := strconv.Atoi(column)
				if err != nil {
					workHistory5.EmployeeNumberGroup = null.NewInt(0, false)
				} else {
					workHistory5.EmployeeNumberGroup = null.NewInt(int64(val), true)
				}

				//株式公開
			case 229:
				workHistory5.PublicOffering = entity.GetIntPublicOffering(column)

				// 雇用形態
			case 230:
				workHistory5.EmploymentStatus = entity.GetIntEmploymentStatus(column)

			// 入社年
			case 231:
				workHistory5.JoiningYear = column

			// 開始ステータス
			case 232:
				workHistory5.FirstStatus = entity.GetIntFirstStatusForWorkHistory(column)

				//退社年
			case 233:
				workHistory5.RetireYear = column

			//終了ステータス
			case 234:
				workHistory5.LastStatus = entity.GetIntLastStatusForWorkHistory(column)

				// 退職理由（本音）
			case 235:
				workHistory5.RetireReasonOfTruth = column

				//退職理由（建前）
			case 236:
				workHistory5.RetireReasonOfPublic = column

				// 職種・役職歴1
				//職種
			case 237, 238, 239:
				experienceOccupation13.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation13.Occupation.Valid {
					departmentHistory13.ExperienceOccupations = append(departmentHistory13.ExperienceOccupations, experienceOccupation13)
				}

				//部門名
			case 240:
				departmentHistory13.Department = column

				//マネジメント人数
			case 241:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory13.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory13.ManagementNumber = null.NewInt(int64(v), true)
				}

				// マネジメント経験詳細
			case 242:
				departmentHistory13.ManagementDetail = column

				//職務内容
			case 243:
				departmentHistory13.JobDescription = column

				// 開始年月
			case 244:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory13.StartYear = column

				// 終了年月
			case 245:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory13.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory13.Department != "" ||
					departmentHistory13.ManagementNumber.Valid || departmentHistory13.ManagementDetail != "" ||
					departmentHistory13.StartYear != "" || departmentHistory13.EndYear != "" ||
					departmentHistory13.JobDescription != "" {
					workHistory5.DepartmentHistories = append(workHistory5.DepartmentHistories, departmentHistory13)
				}

			// 職種・役職歴2
			//職種
			case 246, 247, 248:
				experienceOccupation14.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation14.Occupation.Valid {
					departmentHistory14.ExperienceOccupations = append(departmentHistory14.ExperienceOccupations, experienceOccupation14)
				}

				//部門名
			case 249:
				departmentHistory14.Department = column

				//マネジメント人数
			case 250:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory14.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory14.ManagementNumber = null.NewInt(int64(v), true)
				}

			case 251:
				// マネジメント経験詳細
				departmentHistory14.ManagementDetail = column

			//職務内容
			case 252:
				departmentHistory14.JobDescription = column

			// 開始年月
			case 253:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory14.StartYear = column

			// 終了年月
			case 254:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory14.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory14.Department != "" ||
					departmentHistory14.ManagementNumber.Valid || departmentHistory14.ManagementDetail != "" ||
					departmentHistory14.StartYear != "" || departmentHistory14.EndYear != "" ||
					departmentHistory14.JobDescription != "" {
					workHistory5.DepartmentHistories = append(workHistory5.DepartmentHistories, departmentHistory14)
				}

				// 職種・役職歴3
				//職種
			case 255, 256, 257:
				experienceOccupation15.Occupation = entity.GetIntOccupation(column)
				if experienceOccupation15.Occupation.Valid {
					departmentHistory15.ExperienceOccupations = append(departmentHistory15.ExperienceOccupations, experienceOccupation15)
				}

				//部門名
			case 258:
				departmentHistory15.Department = column

				//マネジメント人数
			case 259:
				v, err := strconv.Atoi(column)
				if err != nil {
					departmentHistory15.ManagementNumber = null.NewInt(0, false)
				} else {
					departmentHistory15.ManagementNumber = null.NewInt(int64(v), true)
				}

				// マネジメント経験詳細
			case 260:
				departmentHistory15.ManagementDetail = column

			//職務内容
			case 261:
				departmentHistory15.JobDescription = column

			// 開始年月
			case 262:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory15.StartYear = column

			// 終了年月
			case 263:
				// 長さが超えている場合は、切り取る
				if len(column) > 7 {
					column = column[:7]
				}
				departmentHistory15.EndYear = column

				// どれかひとつでも入力されていれば、部門歴に追加
				if departmentHistory15.Department != "" ||
					departmentHistory15.ManagementNumber.Valid || departmentHistory15.ManagementDetail != "" ||
					departmentHistory15.StartYear != "" || departmentHistory15.EndYear != "" ||
					departmentHistory15.JobDescription != "" {
					workHistory5.DepartmentHistories = append(workHistory5.DepartmentHistories, departmentHistory15)
				}

				if workHistory5.CompanyName != "" {
					jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, workHistory5)
				}

				/************* PCスキル ****************/
				// PCスキル（Excel）
			case 264:
				if column == "" {
					continue

				} else {
					// カテゴリーをexcellに指定
					jobSeeker.ExcelSkill = entity.GetIntExcelSkill(column)
				}

				// PCスキル（Word）
			case 265:
				if column == "" {
					continue
				} else {
					jobSeeker.WordSkill = entity.GetIntWordSkill(column)

				}

				// PCスキル（Power P）
			case 266:
				if column == "" {
					continue
				} else {
					jobSeeker.PowerPointSkill = entity.GetIntPowerPointSkill(column)
				}

			/******* 免許・資格 ********/
			// 運転免許の有無
			// case
			// 	if column == "あり" {
			// 		license.LicenseType = null.NewInt(4803, true) // 	{4803, "普通自動車免許"},
			// 	} else {
			// 		license.LicenseType = null.NewInt(0, false)
			// 	}

			// case 346:
			// 	license.AcquisitionTime = column

			// 	if license.LicenseType.Valid {
			// 		jobSeeker.Licenses = append(jobSeeker.Licenses, license)
			// 	}

			//免許・資格
			case 267, 269, 271:
				license.LicenseType = entity.GetIntLicenseType(column)

			case 268, 270, 272:
				license.AcquisitionTime = column

				if license.LicenseType.Valid {
					jobSeeker.Licenses = append(jobSeeker.Licenses, license)
				}

			/************* 業務ツール経験 ****************/
			// ※複数選択"	"PCスキル（その他）
			case 273, 274, 275:
				if column == "" {
					continue
				} else {
					pcTool.Tool = entity.GetIntPCTool(column)

					if pcTool.Tool.Valid {
						jobSeeker.PCTools = append(jobSeeker.PCTools, pcTool)
					}

				}

			/************* 開発経験 ****************/
			//	開発言語
			case 276, 279, 282:
				developmentType := entity.GetIntDevelopmentType(0, column)
				developmentSkill.DevelopmentCategory = null.NewInt(0, true)
				developmentSkill.DevelopmentType = developmentType

				//経験年数（年）
			case 277, 280, 283:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentSkill.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentSkill.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 278, 281, 284:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentSkill.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentSkill.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if developmentSkill.DevelopmentType.Valid {
					jobSeeker.DevelopmentSkills = append(jobSeeker.DevelopmentSkills, developmentSkill)
				}

				// 開発OS
			case 285, 288, 291:
				developmentType := entity.GetIntDevelopmentType(1, column)
				developmentSkill.DevelopmentCategory = null.NewInt(1, true)
				developmentSkill.DevelopmentType = developmentType

				//経験年数（年）
			case 286, 289, 292:
				developmentExperienceYear, err := strconv.Atoi(column)
				if err != nil {
					developmentSkill.ExperienceYear = null.NewInt(0, false)
				} else {
					developmentSkill.ExperienceYear = null.NewInt(int64(developmentExperienceYear), true)
				}

				//経験年数（月）
			case 287, 290, 293:
				developmentExperienceMonth, err := strconv.Atoi(column)
				if err != nil {
					developmentSkill.ExperienceMonth = null.NewInt(0, false)
				} else {
					developmentSkill.ExperienceMonth = null.NewInt(int64(developmentExperienceMonth), true)
				}

				// 主要項目が未入力の場合は追加しない
				if developmentSkill.DevelopmentType.Valid {
					jobSeeker.DevelopmentSkills = append(jobSeeker.DevelopmentSkills, developmentSkill)
				}

				/*********** 語学スキル ***************/
				// ※複数選択"
				//英語力
			// case 294:
			// 	languageSkill1.LanguageType = entity.GetIntLanguageType(column)

			case 294:
				languageSkill1.LanguageLevel = entity.GetIntLanguageLevel(column)

			//"TOEIC ※英語のみ"
			case 295:
				toeic, err := strconv.Atoi(column)
				if err != nil {
					languageSkill1.Toeic = null.NewInt(0, false)
				} else {
					languageSkill1.Toeic = null.NewInt(int64(toeic), true)
				}

				// TOEIC 受験年月 (yyyyy-mm)
			case 296:
				languageSkill1.ToeicExaminationYear = column

				//	"TOEFL(i) ※英語のみ"
			case 297:
				toeflI, err := strconv.Atoi(column)
				if err != nil {
					languageSkill1.ToeflIBT = null.NewInt(0, false)
				} else {
					languageSkill1.ToeflIBT = null.NewInt(int64(toeflI), true)
				}

				// TOEFL(i) 受験年月 (yyyyy-mm)
			case 298:
				languageSkill1.ToeflIBTExaminationYear = column

				// "TOEFL(P) ※英語のみ"
			case 299:
				toeflP, err := strconv.Atoi(column)
				if err != nil {
					languageSkill1.ToeflPBT = null.NewInt(0, false)
				} else {
					languageSkill1.ToeflPBT = null.NewInt(int64(toeflP), true)
				}

				// TOEFL(P) 受験年月 (yyyyy-mm)
			case 300:
				languageSkill1.ToeflPBTExaminationYear = column

				// 全項目未入力時はスキップ
				if languageSkill1.LanguageLevel.Valid ||
					languageSkill1.Toeic.Valid || languageSkill1.ToeicExaminationYear != "" ||
					languageSkill1.ToeflIBT.Valid || languageSkill1.ToeflIBTExaminationYear != "" ||
					languageSkill1.ToeflPBT.Valid || languageSkill1.ToeflPBTExaminationYear != "" {
					// 英語を指定
					languageSkill1.LanguageType = null.NewInt(0, true)
					jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, languageSkill1)
				}

				//語学2
			case 301:
				languageSkill2.LanguageType = entity.GetIntLanguageType(column)

			case 302:
				languageSkill2.LanguageLevel = entity.GetIntLanguageLevel(column)

				// 全項目未入力時はスキップ
				if languageSkill2.LanguageType.Valid {
					jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, languageSkill2)
				}

				//語学3
			case 303:
				languageSkill3.LanguageType = entity.GetIntLanguageType(column)

			case 304:
				languageSkill3.LanguageLevel = entity.GetIntLanguageLevel(column)

				// 全項目未入力時はスキップ
				if languageSkill3.LanguageType.Valid {
					jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, languageSkill3)
				}

				/********** 希望 **********/
			//希望休日タイプ 3つ
			case 305, 306, 307:
				desiredHolidayType.HolidayType = entity.GetIntHolidayForJobSeeker(column)

				if desiredHolidayType.HolidayType.Valid {
					jobSeeker.DesiredHolidayTypes = append(jobSeeker.DesiredHolidayTypes, desiredHolidayType)
				}

			//希望企業規模 3つ
			case 308, 309, 310:
				desiredCompanyScale.DesiredCompanyScale = entity.GetIntCompanyScale(column)

				if desiredCompanyScale.DesiredCompanyScale.Valid {
					jobSeeker.DesiredCompanyScales = append(jobSeeker.DesiredCompanyScales, desiredCompanyScale)
				}

				// 希望業界
				// 第一希望群
			case 311, 312, 313:
				desiredIndustry.DesiredIndustry = entity.GetIntIndustry(column)
				desiredIndustry.DesiredRank = null.NewInt(1, true)

				if desiredIndustry.DesiredIndustry.Valid {
					jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, desiredIndustry)
				}

				// 第二希望群
			case 314, 315, 316:
				desiredIndustry.DesiredIndustry = entity.GetIntIndustry(column)
				desiredIndustry.DesiredRank = null.NewInt(2, true)

				if desiredIndustry.DesiredIndustry.Valid {
					jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, desiredIndustry)
				}

				// 第三希望群
			case 317, 318, 319:
				desiredIndustry.DesiredIndustry = entity.GetIntIndustry(column)
				desiredIndustry.DesiredRank = null.NewInt(3, true)

				if desiredIndustry.DesiredIndustry.Valid {
					jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, desiredIndustry)
				}

				// 希望職種
				// 第一希望群
			case 320, 321, 322:
				desiredOccupation.DesiredOccupation = entity.GetIntOccupation(column)
				desiredOccupation.DesiredRank = null.NewInt(1, true)

				if desiredOccupation.DesiredOccupation.Valid {
					jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, desiredOccupation)
				}

				// 第二希望群
			case 323, 324, 325:
				desiredOccupation.DesiredOccupation = entity.GetIntOccupation(column)
				desiredOccupation.DesiredRank = null.NewInt(2, true)

				if desiredOccupation.DesiredOccupation.Valid {
					jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, desiredOccupation)
				}

				// 第三希望群
			case 326, 327, 328:
				desiredOccupation.DesiredOccupation = entity.GetIntOccupation(column)
				desiredOccupation.DesiredRank = null.NewInt(3, true)

				if desiredOccupation.DesiredOccupation.Valid {
					jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, desiredOccupation)
				}

				// 希望勤務地
				// 第一希望 *複数ある場合は「,」で区切る
			case 329:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {
					str = strings.TrimSpace(str)
					desiredWorkLocation = entity.JobSeekerDesiredWorkLocation{
						DesiredWorkLocation: entity.GetIntPrefecture(str),
						DesiredRank:         null.NewInt(1, true),
					}

					if desiredWorkLocation.DesiredWorkLocation.Valid {
						jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, desiredWorkLocation)
					}
				}

				// 第二希望 *複数ある場合は「,」で区切る
			case 330:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {
					str = strings.TrimSpace(str)
					desiredWorkLocation = entity.JobSeekerDesiredWorkLocation{
						DesiredWorkLocation: entity.GetIntPrefecture(str),
						DesiredRank:         null.NewInt(2, true),
					}

					if desiredWorkLocation.DesiredWorkLocation.Valid {
						jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, desiredWorkLocation)
					}
				}

				// 第三希望 *複数ある場合は「,」で区切る
			case 331:
				strSlice := strings.Split(column, ",")
				for _, str := range strSlice {
					str = strings.TrimSpace(str)
					desiredWorkLocation = entity.JobSeekerDesiredWorkLocation{
						DesiredWorkLocation: entity.GetIntPrefecture(str),
						DesiredRank:         null.NewInt(3, true),
					}

					if desiredWorkLocation.DesiredWorkLocation.Valid {
						jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, desiredWorkLocation)
					}
				}

			//転勤可否
			case 332:
				jobSeeker.Transfer = entity.GetIntTransferForJobSeeker(column)

				//転勤備考
			case 333:
				jobSeeker.TransferRequirement = column

				//希望年収
			case 334:
				val, err := strconv.Atoi(column)
				if err != nil {
					jobSeeker.DesiredAnnualIncome = null.NewInt(0, false)
				} else {
					jobSeeker.DesiredAnnualIncome = null.NewInt(int64(val), true)
				}

				//入社可能時期
			case 335:
				jobSeeker.JoinCompanyPeriod = entity.GetIntJoinCompanyPeriod(column)

			//アピアランス
			case 336:
				jobSeeker.Appearance = entity.GetIntAppearanceForJobSeeker(column)

				// アピアランス(本音)
			// case 337:
			// jobSeeker.AppearanceDetailOfTruth = column

			// アピアランス(推薦状用)
			// case 338:
			// jobSeeker.AppearanceDetail = column

			//コミニュケーション
			case 339:
				jobSeeker.Communication = entity.GetIntCommunicationForJobSeeker(column)

				// コミニュケーション(本音)
			// case 340:
			// jobSeeker.CommunicationDetailOfTruth = column

			// コミニュケーション(推薦状用)
			// case 341:
			// jobSeeker.CommunicationDetail = column

			//論理的思考力
			case 342:
				jobSeeker.Thinking = entity.GetIntThinkingForJobSeeker(column)

				// 論理的思考力(本音)
			// case 343:
			// jobSeeker.ThinkingDetailOfTruth = column

			// 論理的思考力(推薦状用)
			// case 344:
			// jobSeeker.ThinkingDetail = column

			/********* 自己PR ***********/

			case 345, 347:
				selfPromotion.Title = column

			case 346, 348:
				selfPromotion.Contents = column

				// どっちか書いておけばいい
				if selfPromotion.Title != "" || selfPromotion.Contents != "" {
					jobSeeker.SelfPromotions = append(jobSeeker.SelfPromotions, selfPromotion)
				}

				// 研究・学チカ
			case 349:
				jobSeeker.ResearchContent = column

				//社内限定メモ
			case 350:
				jobSeeker.SecretMemo = column

				//他社エージェント向けメモ
				// case 351:
				// jobSeeker.PublicMemo = column

				/**************************************/

			}
		}

		jobSeekerList = append(jobSeekerList, &jobSeeker)
	}

	return jobSeekerList, nil
}
