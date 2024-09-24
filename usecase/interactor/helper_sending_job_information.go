package interactor

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

// 求人の最大ページ数を返す（本番実装までは1ページあたり20件）
func getSendingJobInformationListMaxPage(sendingJobInformationList []*entity.SendingJobInformation) uint {
	var maxPage = len(sendingJobInformationList) / 20
	if 0 < (len(sendingJobInformationList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの求人一覧を返す（本番実装までは1ページあたり20件）
func getSendingJobInformationListWithPage(sendingJobInformationList []*entity.SendingJobInformation, page uint) []*entity.SendingJobInformation {
	var (
		perPage uint = 20
		listLen uint = uint(len(sendingJobInformationList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return sendingJobInformationList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.SendingJobInformation{}
	}

	if (listLen - first) <= perPage {
		return sendingJobInformationList[first:]
	}
	return sendingJobInformationList[first:last]
}

// 求人リストからIDリストを取得する
func getSendingJobInformationIDList(sendingJobInformationList []*entity.SendingJobInformation) []uint {
	var idListUint []uint

	if len(sendingJobInformationList) == 0 {
		return idListUint
	}

	for _, sendingJobInformation := range sendingJobInformationList {
		idListUint = append(idListUint, sendingJobInformation.ID)
	}

	return idListUint
}

// 絞り込み検索処理 本体
func searchSendingJobInformationList(jobInformationList []*entity.SendingJobInformation, searchParam entity.SearchJobInformation) ([]*entity.SendingJobInformation, error) {

	// 絞り込み項目の結果を代入するための変数を用意
	var (
		jobInformationListWithAgentStaffID []*entity.SendingJobInformation
		/**
		アライアンスの中間テーブル作成後に
		アライアンスの絞り込みを追加
		*/
		jobInformationListWithIndustry                     []*entity.SendingJobInformation // or
		jobInformationListWithOccupation                   []*entity.SendingJobInformation // or
		jobInformationListWithEmployment                   []*entity.SendingJobInformation // or
		jobInformationListWithPrefecture                   []*entity.SendingJobInformation // or
		jobInformationListWithUnderIncome                  []*entity.SendingJobInformation // or
		jobInformationListWithOverIncome                   []*entity.SendingJobInformation // or
		jobInformationListWithGender                       []*entity.SendingJobInformation
		jobInformationListWithAge                          []*entity.SendingJobInformation
		jobInformationListWithFinalEducation               []*entity.SendingJobInformation
		jobInformationListWithSchoolLevel                  []*entity.SendingJobInformation
		jobInformationListWithStudyCategory                []*entity.SendingJobInformation
		jobInformationListWithNationality                  []*entity.SendingJobInformation
		jobInformationListWithJobChange                    []*entity.SendingJobInformation
		jobInformationListWithShortResignation             []*entity.SendingJobInformation
		jobInformationListWithAppearance                   []*entity.SendingJobInformation
		jobInformationListWithCommunication                []*entity.SendingJobInformation
		jobInformationListWithThinking                     []*entity.SendingJobInformation
		jobInformationListWithRequiredExperienceIndustry   []*entity.SendingJobInformation // and
		jobInformationListWithRequiredExperienceOccupation []*entity.SendingJobInformation // and
		jobInformationListWithRequiredSocialExperience     []*entity.SendingJobInformation
		jobInformationListWithRequiredManagement           []*entity.SendingJobInformation
		jobInformationListWithRequiredLicense              []*entity.SendingJobInformation // and
		jobInformationListWithRequiredLanguage             []*entity.SendingJobInformation // and
		jobInformationListWithRequiredExcelSkill           []*entity.SendingJobInformation
		jobInformationListWithRequiredWordSkill            []*entity.SendingJobInformation
		jobInformationListWithRequiredPowerPointSkill      []*entity.SendingJobInformation
		jobInformationListWithRequiredAnotherPCSkill       []*entity.SendingJobInformation // and
		jobInformationListWithRequiredDevelopmentLanguage  []*entity.SendingJobInformation // and
		jobInformationListWithRequiredDevelopmentOS        []*entity.SendingJobInformation // and
		jobInformationListWithTransfer                     []*entity.SendingJobInformation
		jobInformationListWithHolidayType                  []*entity.SendingJobInformation
		jobInformationListWithCompanyScale                 []*entity.SendingJobInformation
		jobInformationListWithFeature                      []*entity.SendingJobInformation
	)

	// 営業担当者IDがある場合
	agentStaffID, err := strconv.Atoi(searchParam.AgentStaffID)
	if !(err != nil || agentStaffID == 0) {
		for _, jobInformation := range jobInformationList {
			if jobInformation.AgentStaffID == uint(agentStaffID) {
				jobInformationListWithAgentStaffID = append(jobInformationListWithAgentStaffID, jobInformation)
			}
		}
	}

	// 営業担当者IDが無い場合
	if err != nil || agentStaffID == 0 {
		jobInformationListWithAgentStaffID = jobInformationList
	}

	fmt.Println("RA担当者: ", len(jobInformationListWithAgentStaffID))

	// NOTE: 業界
	// 業界のいずれかが入っている場合
	if !(len(searchParam.Industries) == 0) {
	industryLoop:
		for _, jobInformation := range jobInformationListWithAgentStaffID {
			for _, industryParam := range searchParam.Industries {
				if !industryParam.Valid {
					continue
				}
				for _, jobInfoIndustry := range jobInformation.Industries {
					if !jobInfoIndustry.Industry.Valid {
						continue
					}
					if industryParam == jobInfoIndustry.Industry {
						jobInformationListWithIndustry = append(jobInformationListWithIndustry, jobInformation)
						continue industryLoop
					}
				}
			}
		}
	}

	// 業界のいずれかも入っていない場合
	if len(searchParam.Industries) == 0 {
		jobInformationListWithIndustry = jobInformationListWithAgentStaffID
	}

	fmt.Println("業界: ", len(jobInformationListWithIndustry))

	// NOTE: 職種
	// 職種のいずれかが入っている場合
	if !(len(searchParam.Occupations) == 0) {
	occupationLoop:
		for _, jobInformation := range jobInformationListWithIndustry {
			for _, jobInfoOccupation := range jobInformation.Occupations {
				for _, occupation := range searchParam.Occupations {
					if !occupation.Valid {
						continue
					}
					if occupation == jobInfoOccupation.Occupation {
						// 求人の職種が一つでも合致していればcontinueで次の求人へ
						jobInformationListWithOccupation = append(jobInformationListWithOccupation, jobInformation)
						continue occupationLoop
					}
				}
			}
		}
	}

	// 職種のいずれかも入っていない場合
	if len(searchParam.Occupations) == 0 {
		jobInformationListWithOccupation = jobInformationListWithIndustry
	}

	fmt.Println("職種: ", len(jobInformationListWithOccupation))

	// NOTE: 雇用形態
	// 雇用形態のいずれかが入っている場合
	if !(len(searchParam.Employments) == 0) {
	employmentLoop:
		for _, jobInformation := range jobInformationListWithOccupation {
			for _, jobInfoEmployment := range jobInformation.EmploymentStatuses {
				for _, employment := range searchParam.Employments {
					if !employment.Valid {
						continue
					}
					if employment == jobInfoEmployment.EmploymentStatus {
						// 求人の職種が一つでも合致していればcontinueで次の求人へ
						jobInformationListWithEmployment = append(jobInformationListWithEmployment, jobInformation)
						continue employmentLoop
					}
				}
			}
		}
	}

	// 職種のいずれかも入っていない場合
	if len(searchParam.Employments) == 0 {
		jobInformationListWithEmployment = jobInformationListWithOccupation
	}

	fmt.Println("雇用形態: ", len(jobInformationListWithEmployment))

	// NOTE: 勤務地
	// 勤務地のいずれかが入っている場合
	if !(len(searchParam.Prefectures) == 0) {
	prefectureLoop:
		for _, jobInformation := range jobInformationListWithEmployment {
			for _, prefectureParam := range searchParam.Prefectures {
				if !prefectureParam.Valid {
					continue
				}
				for _, jobInfoPrefecture := range jobInformation.Prefectures {
					// 合致する or 勤務地が全国各地の場合
					if prefectureParam == jobInfoPrefecture.Prefecture || jobInfoPrefecture.Prefecture == null.NewInt(9999, true) {
						jobInformationListWithPrefecture = append(jobInformationListWithPrefecture, jobInformation)
						continue prefectureLoop
					}
				}
			}
		}
	}

	// 勤務地のいずれかも入っていない場合
	if len(searchParam.Prefectures) == 0 {
		jobInformationListWithPrefecture = jobInformationListWithEmployment
	}

	fmt.Println("勤務地: ", len(jobInformationListWithPrefecture))

	// Note: 年収下限
	// 年収下限がある場合
	underIncome, err := strconv.Atoi(searchParam.UnderIncome)
	if !(err != nil) {
		for _, jobInformation := range jobInformationListWithPrefecture {
			if jobInformation.UnderIncome == null.NewInt(0, false) {
				continue
			}
			if underIncome <= int(jobInformation.UnderIncome.Int64) {
				jobInformationListWithUnderIncome = append(jobInformationListWithUnderIncome, jobInformation)
			}
		}
	}

	// 年収下限が無い場合
	if err != nil {
		jobInformationListWithUnderIncome = jobInformationListWithPrefecture
	}

	fmt.Println("年収下限: ", len(jobInformationListWithUnderIncome))

	// Note: 年収上限
	// 年収上限がある場合
	overIncome, err := strconv.Atoi(searchParam.OverIncome)
	if !(err != nil) {
		for _, jobInformation := range jobInformationListWithUnderIncome {
			if jobInformation.OverIncome == null.NewInt(0, false) {
				continue
			}

			if overIncome >= int(jobInformation.OverIncome.Int64) {
				jobInformationListWithOverIncome = append(jobInformationListWithOverIncome, jobInformation)
			}
		}
	}

	// 年収上限が無い場合
	if err != nil {
		jobInformationListWithOverIncome = jobInformationListWithUnderIncome
	}

	fmt.Println("年収上限: ", len(jobInformationListWithOverIncome))

	// Note: 性別
	// 性別がある場合
	if !(len(searchParam.GenderTypes) == 0) {
	genderLoop:
		for _, jobInformation := range jobInformationListWithOverIncome {
			for _, gender := range searchParam.GenderTypes {
				if !gender.Valid {
					continue
				}
				if gender == jobInformation.Gender {
					// 求人の希望性別が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithGender = append(jobInformationListWithGender, jobInformation)
					continue genderLoop
				}
			}
		}
	}

	// 性別のいずれかも入っていない場合
	if len(searchParam.GenderTypes) == 0 {
		jobInformationListWithGender = jobInformationListWithOverIncome
	}

	fmt.Println("性別: ", len(jobInformationListWithGender))

	// Note: 年齢
	// 年齢がある場合
	age, err := strconv.Atoi(searchParam.Age)
	if !(err != nil) {
		for _, jobInformation := range jobInformationListWithGender {
			// 募集年齢: 不問（年齢上限、年齢下限が未入力）
			if !jobInformation.AgeUnder.Valid && !jobInformation.AgeOver.Valid {
				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
				continue
			}

			// 年齢下限と上限が入力されており、検索パラムの値が下限以上で上限以下の場合
			if jobInformation.AgeUnder.Valid && age >= int(jobInformation.AgeUnder.Int64) && jobInformation.AgeOver.Valid && age <= int(jobInformation.AgeOver.Int64) {
				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
				continue
			}

			// 年齢下限が入力されている + 検索パラムの値が上限以下 + 年齢上限が未入力
			if jobInformation.AgeUnder.Valid && age >= int(jobInformation.AgeUnder.Int64) && !jobInformation.AgeOver.Valid {
				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
				continue
			}

			// 年齢上限が入力されている + 検索パラムの値が上限以下 + 年齢下限が未入力
			if jobInformation.AgeOver.Valid && age <= int(jobInformation.AgeOver.Int64) && !jobInformation.AgeUnder.Valid {
				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
				continue
			}
		}
	}

	// 年齢下限が無い場合
	if err != nil {
		jobInformationListWithAge = jobInformationListWithGender
	}

	fmt.Println("年齢: ", len(jobInformationListWithAge))

	// // Note: 年齢上限
	// // 年齢上限がある場合
	// ageOver, err := strconv.Atoi(searchParam.AgeOver)
	// if !(err != nil) {
	// 	for _, jobInformation := range jobInformationListWithAgeUnder {
	// 		if jobInformation.AgeOver == null.NewInt(0, false) {
	// 			continue
	// 		}

	// 		if ageOver >= int(jobInformation.AgeOver.Int64) {
	// 			jobInformationListWithAgeOver = append(jobInformationListWithAgeOver, jobInformation)
	// 		}
	// 	}
	// }

	// // 年齢上限が無い場合
	// if len(searchParam.AgeOver) == 0 {
	// 	jobInformationListWithAgeOver = jobInformationListWithAgeUnder
	// }

	// fmt.Println("年齢上限: ", len(jobInformationListWithAgeOver))

	// Note: 最終学歴
	// 最終学歴がある場合
	if !(len(searchParam.FinalEducationTypes) == 0) {
	finalEducationLoop:
		for _, jobInformation := range jobInformationListWithAge {
			for _, finalEducation := range searchParam.FinalEducationTypes {
				if !finalEducation.Valid {
					continue
				}
				if finalEducation == jobInformation.FinalEducation {
					// 求人の希望最終学歴が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithFinalEducation = append(jobInformationListWithFinalEducation, jobInformation)
					continue finalEducationLoop
				}
			}
		}
	}

	//  最終学歴のいずれかも入っていない場合
	if len(searchParam.FinalEducationTypes) == 0 {
		jobInformationListWithFinalEducation = jobInformationListWithAge
	}

	fmt.Println("最終学歴: ", len(jobInformationListWithFinalEducation))

	// Note: 大学レベル
	// 大学レベルがある場合
	if !(len(searchParam.SchoolLevelTypes) == 0) {
	schoolLevelLoop:
		for _, jobInformation := range jobInformationListWithFinalEducation {
			for _, schoolLevel := range searchParam.SchoolLevelTypes {
				if !schoolLevel.Valid {
					continue
				}
				if schoolLevel == jobInformation.SchoolLevel {
					// 求人の希望大学レベルが一つでも合致していればcontinueで次の求人へ
					jobInformationListWithSchoolLevel = append(jobInformationListWithSchoolLevel, jobInformation)
					continue schoolLevelLoop
				}
			}
		}
	}

	// 大学レベルが無い場合
	if len(searchParam.SchoolLevelTypes) == 0 {
		jobInformationListWithSchoolLevel = jobInformationListWithFinalEducation
	}

	fmt.Println("大学レベル: ", len(jobInformationListWithSchoolLevel))

	// Note: 文系・理系
	// 文系・理系がある場合
	if !(len(searchParam.StudyCategoryTypes) == 0) {
	schoolCategoryLoop:
		for _, jobInformation := range jobInformationListWithSchoolLevel {
			for _, studyCategory := range searchParam.StudyCategoryTypes {
				if !studyCategory.Valid {
					continue
				}
				if studyCategory == jobInformation.StudyCategory {
					// 求人の文系・理系が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithStudyCategory = append(jobInformationListWithStudyCategory, jobInformation)
					continue schoolCategoryLoop
				}
			}
		}
	}

	// 文系・理系が無い場合
	if len(searchParam.StudyCategoryTypes) == 0 {
		jobInformationListWithStudyCategory = jobInformationListWithSchoolLevel
	}

	fmt.Println("文系・理系: ", len(jobInformationListWithStudyCategory))

	// Note: 国籍
	// 国籍がある場合
	if !(len(searchParam.NationalityTypes) == 0) {
	nationalityloop:
		for _, jobInformation := range jobInformationListWithStudyCategory {
			for _, nationality := range searchParam.NationalityTypes {
				if !nationality.Valid {
					continue
				}
				if nationality == jobInformation.Nationality {
					// 求人の希望国籍が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithNationality = append(jobInformationListWithNationality, jobInformation)
					continue nationalityloop
				}
			}
		}
	}

	// 国籍が無い場合
	if len(searchParam.NationalityTypes) == 0 {
		jobInformationListWithNationality = jobInformationListWithStudyCategory
	}

	fmt.Println("国籍: ", len(jobInformationListWithNationality))

	// Note: 転職回数
	// 転職回数がある場合
	if !(len(searchParam.JobChangeTypes) == 0) {
	jobChangeLoop:
		for _, jobInformation := range jobInformationListWithNationality {
			for _, jobChange := range searchParam.JobChangeTypes {
				if !jobChange.Valid {
					continue
				}
				if jobChange == jobInformation.JobChange {
					// 求人の転職回数が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithJobChange = append(jobInformationListWithJobChange, jobInformation)
					continue jobChangeLoop
				}
			}
		}
	}

	// 転職回数が無い場合
	if len(searchParam.JobChangeTypes) == 0 {
		jobInformationListWithJobChange = jobInformationListWithNationality
	}

	fmt.Println("転職回数: ", len(jobInformationListWithJobChange))

	// Note: 短期離職
	// 短期離職がある場合
	if !(len(searchParam.ShortResignationTypes) == 0) {
	shortResignationloop:
		for _, jobInformation := range jobInformationListWithJobChange {
			for _, shortResignation := range searchParam.ShortResignationTypes {
				if !shortResignation.Valid {
					continue
				}
				if shortResignation == jobInformation.ShortResignation {
					// 求人の短期離職が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithShortResignation = append(jobInformationListWithShortResignation, jobInformation)
					continue shortResignationloop
				}
			}
		}
	}

	// 短期離職が無い場合
	if len(searchParam.ShortResignationTypes) == 0 {
		jobInformationListWithShortResignation = jobInformationListWithJobChange
	}

	fmt.Println("短期離職: ", len(jobInformationListWithShortResignation))

	// Note: アピアランス
	// アピアランスがある場合
	if !(len(searchParam.AppearanceTypes) == 0) {
	appearanceLoop:
		for _, jobInformation := range jobInformationListWithShortResignation {
			for _, appearance := range searchParam.AppearanceTypes {
				if !appearance.Valid {
					continue
				}
				if appearance == jobInformation.Appearance {
					// 求人のアピアランスが一つでも合致していればcontinueで次の求人へ
					jobInformationListWithAppearance = append(jobInformationListWithAppearance, jobInformation)
					continue appearanceLoop
				}
			}
		}
	}

	// アピアランスが無い場合
	if len(searchParam.AppearanceTypes) == 0 {
		jobInformationListWithAppearance = jobInformationListWithShortResignation
	}

	fmt.Println("アピアランス: ", len(jobInformationListWithAppearance))

	// Note: コミュ力
	// コミュ力がある場合
	if !(len(searchParam.CommunicationTypes) == 0) {
	communicationLoop:
		for _, jobInformation := range jobInformationListWithAppearance {
			for _, communication := range searchParam.CommunicationTypes {
				if !communication.Valid {
					continue
				}
				if communication == jobInformation.Communication {
					// 求人のコミュ力が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithCommunication = append(jobInformationListWithCommunication, jobInformation)
					continue communicationLoop
				}
			}
		}
	}

	// コミュ力が無い場合
	if len(searchParam.CommunicationTypes) == 0 {
		jobInformationListWithCommunication = jobInformationListWithAppearance
	}

	fmt.Println("コミュ力: ", len(jobInformationListWithCommunication))

	// Note: 論理的思考力
	// 論理的思考力がある場合
	if !(len(searchParam.ThinkingTypes) == 0) {
	thinkingLoop:
		for _, jobInformation := range jobInformationListWithCommunication {
			for _, thinking := range searchParam.ThinkingTypes {
				if !thinking.Valid {
					continue
				}
				if thinking == jobInformation.Thinking {
					// 求人の論理的思考力が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithThinking = append(jobInformationListWithThinking, jobInformation)
					continue thinkingLoop
				}
			}
		}
	}

	// 論理的思考力が無い場合
	if len(searchParam.ThinkingTypes) == 0 {
		jobInformationListWithThinking = jobInformationListWithCommunication
	}

	fmt.Println("論理的思考力: ", len(jobInformationListWithThinking))

	// NOTE: 必要経験業界（And検索）
	// 必要経験業界のいずれかが入っている場合
	if !(len(searchParam.RequiredExperienceIndustries) == 0) {

		// 検索パラムのmap
		var experienceIndustrieParams = make(map[null.Int]bool)
		for _, industry := range searchParam.RequiredExperienceIndustries {
			experienceIndustrieParams[industry] = false
		}

	experienceIndustryLoop:
		for _, jobInformation := range jobInformationListWithThinking {
			if len(jobInformation.RequiredConditions) == 0 {
				continue
			}

			/**
			合致した時にexperienceIndustryLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 経験業界のリスト作成
			var experienceIndustries []entity.SendingJobInformationRequiredExperienceIndustry
			for _, infoCondition := range jobInformation.RequiredConditions {
				experienceIndustries = append(experienceIndustries, infoCondition.RequiredExperienceJobs.ExperienceIndustries...)
			}

			// 求人の経験業界のmap
			var experienceIndustryDatas = make(map[null.Int]bool)

			// パラムのループ
			for _, industry := range searchParam.RequiredExperienceIndustries {
				if !industry.Valid {
					continue
				}

				for _, experienceIndustry := range experienceIndustries {
					if industry == experienceIndustry.ExperienceIndustry {
						experienceIndustrieParams[industry] = true
						experienceIndustryDatas[experienceIndustry.ExperienceIndustry] = true
					}
				}

				if reflect.DeepEqual(experienceIndustrieParams, experienceIndustryDatas) {
					jobInformationListWithRequiredExperienceIndustry = append(jobInformationListWithRequiredExperienceIndustry, jobInformation)
					continue experienceIndustryLoop
				}
			}
		}
	}

	// 必要経験業界のいずれかも入っていない場合
	if len(searchParam.RequiredExperienceIndustries) == 0 {
		jobInformationListWithRequiredExperienceIndustry = jobInformationListWithThinking
	}

	fmt.Println("必要経験業界: ", len(jobInformationListWithRequiredExperienceIndustry))

	// NOTE: 必要経験職種（And検索）
	// 必要経験職種のいずれかが入っている場合
	if !(len(searchParam.RequiredExperienceOccupations) == 0) {

		// 検索パラムのmap
		var experienceOccupationParams = make(map[null.Int]bool)
		for _, occupation := range searchParam.RequiredExperienceOccupations {
			experienceOccupationParams[occupation] = false
		}

	experienceOccupationLoop:
		for _, jobInformation := range jobInformationListWithRequiredExperienceIndustry {
			/**
			合致した時にexperienceOccupationLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 経験職種のリスト作成
			var experienceOccupations []entity.SendingJobInformationRequiredExperienceOccupation
			for _, infoCondition := range jobInformation.RequiredConditions {
				experienceOccupations = append(experienceOccupations, infoCondition.RequiredExperienceJobs.ExperienceOccupations...)
			}

			// 求人の経験職種のmap
			var experienceOccupationDatas = make(map[null.Int]bool)

			// パラムのループ
			for _, occupation := range searchParam.RequiredExperienceOccupations {
				if !occupation.Valid {
					continue
				}

				for _, experienceOccupation := range experienceOccupations {
					if occupation == experienceOccupation.ExperienceOccupation {
						experienceOccupationParams[occupation] = true
						experienceOccupationDatas[experienceOccupation.ExperienceOccupation] = true
					}
				}

				if reflect.DeepEqual(experienceOccupationParams, experienceOccupationDatas) {
					jobInformationListWithRequiredExperienceOccupation = append(jobInformationListWithRequiredExperienceOccupation, jobInformation)
					continue experienceOccupationLoop
				}
			}
		}
	}

	// 必要経験職種のいずれかも入っていない場合
	if len(searchParam.RequiredExperienceOccupations) == 0 {
		jobInformationListWithRequiredExperienceOccupation = jobInformationListWithRequiredExperienceIndustry
	}

	fmt.Println("必要経験職種: ", len(jobInformationListWithRequiredExperienceOccupation))

	// NOTE: 必要社会人経験
	// 必要社会人経験のいずれかが入っている場合
	if !(len(searchParam.RequiredSocialExperienceTypes) == 0) {
		experienceYear, errYear := strconv.Atoi(searchParam.RequiredSocialExperienceYear)
		experienceMonth, errMonth := strconv.Atoi(searchParam.RequiredSocialExperienceMonth)

	socialExperienceLoop:
		for _, jobInformation := range jobInformationListWithRequiredExperienceOccupation {
			/**
			合致した時にtrueに変える
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			for _, socialExperience := range searchParam.RequiredSocialExperienceTypes {
				if !socialExperience.Valid {
					continue
				}
				if len(jobInformation.RequiredSocialExperiences) == 0 {
					continue
				}

				for _, requiredSocialExperience := range jobInformation.RequiredSocialExperiences {
					// 不問の場合は年数と月数は関係ない
					if requiredSocialExperience.SocialExperienceType.Int64 == 99 && socialExperience.Int64 == 99 {
						jobInformationListWithRequiredSocialExperience = append(jobInformationListWithRequiredSocialExperience, jobInformation)
						continue socialExperienceLoop
					}

					var (
						year  int64 = 0
						month int64 = 0
					)

					// パラムの年数と月数がどちらも正しく取れていない場合は
					if errYear != nil && errMonth != nil {
						continue
					}

					// 年数が正しく取れていれば代入
					if errYear == nil {
						year = int64(experienceYear)
					}

					// 月数が正しく取れていれば代入
					if errMonth == nil {
						month = int64(experienceMonth)
					}

					// 社会人種別が同じで経験年数がパラムの年数以下の場合
					if socialExperience == requiredSocialExperience.SocialExperienceType &&
						((jobInformation.SocialExperienceYear.Int64*12)+jobInformation.SocialExperienceMonth.Int64) <= ((year*12)+month) {
						/**
						求人の必要社会人経験が合致したタイミングでalreadyMatchをtrueにして
						ダブりが発生しないように制御
						*/
						jobInformationListWithRequiredSocialExperience = append(jobInformationListWithRequiredSocialExperience, jobInformation)
						continue socialExperienceLoop
					}
				}
			}
		}

	}

	// 必要社会人経験のいずれかも入っていない場合
	if len(searchParam.RequiredSocialExperienceTypes) == 0 {
		jobInformationListWithRequiredSocialExperience = jobInformationListWithRequiredExperienceOccupation
	}

	fmt.Println("必要社会人経験: ", len(jobInformationListWithRequiredSocialExperience))

	// Note: マネジメント
	// マネジメントがある場合
	management, err := strconv.Atoi(searchParam.RequiredManagement)
	if !(err != nil) {
	managementLoop:
		for _, jobInformation := range jobInformationListWithRequiredSocialExperience {
			if len(jobInformation.RequiredConditions) == 0 {
				continue
			}
			/**
			合致した時にmanagementLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/
			// パラムのループ
			for _, requiredCondition := range jobInformation.RequiredConditions {
				if requiredCondition.RequiredManagement == null.NewInt(int64(management), true) {
					jobInformationListWithRequiredManagement = append(jobInformationListWithRequiredManagement, jobInformation)
					continue managementLoop
				}
			}
		}
	}

	// マネジメントが無い場合
	if err != nil {
		jobInformationListWithRequiredManagement = jobInformationListWithRequiredSocialExperience
	}

	fmt.Println("マネジメント: ", len(jobInformationListWithRequiredManagement))

	// NOTE: 必要資格（And検索）
	// 必要資格のいずれかが入っている場合
	if !(len(searchParam.RequiredLicenses) == 0) {
		jobInformationListWithRequiredLicense = licenceSearchSendingJobInformation(jobInformationListWithRequiredManagement, searchParam.RequiredLicenses)
	}

	// 必要資格のいずれかも入っていない場合
	if len(searchParam.RequiredLicenses) == 0 {
		jobInformationListWithRequiredLicense = jobInformationListWithRequiredManagement
	}

	fmt.Println("必要資格: ", len(jobInformationListWithRequiredLicense))

	// NOTE: 必要語学（And検索）
	// 必要語学のいずれかが入っている場合
	if !(len(searchParam.RequiredLanguages) == 0) {

		// 検索パラムのmap
		var languageParams = make(map[null.Int]bool)
		for _, languageParam := range searchParam.RequiredLanguages {
			languageParams[languageParam] = false
		}

	languageLoop:
		for _, jobInformation := range jobInformationListWithRequiredLicense {
			if len(jobInformation.RequiredConditions) == 0 {
				continue
			}
			/**
			合致した時にlanguageLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 必要語学のリスト作成
			var languages []entity.SendingJobInformationRequiredLanguageType
			for _, infoCondition := range jobInformation.RequiredConditions {
				languages = append(languages, infoCondition.RequiredLanguages.LanguageTypes...)
			}

			// 求人の語学のmap
			var languageDatas = make(map[null.Int]bool)

			// パラムのループ
			for _, languageParam := range searchParam.RequiredLanguages {
				if !languageParam.Valid {
					continue
				}

				for _, language := range languages {
					if languageParam == language.LanguageType {
						languageParams[languageParam] = true
						languageDatas[language.LanguageType] = true
					}
				}

				if reflect.DeepEqual(languageParams, languageDatas) {
					jobInformationListWithRequiredLanguage = append(jobInformationListWithRequiredLanguage, jobInformation)
					continue languageLoop
				}
			}
		}
	}

	// 必要語学のいずれかも入っていない場合
	if len(searchParam.RequiredLanguages) == 0 {
		jobInformationListWithRequiredLanguage = jobInformationListWithRequiredLicense
	}

	fmt.Println("必要語学: ", len(jobInformationListWithRequiredLanguage))

	// NOTE: excelスキル
	// excelスキルのいずれかが入っている場合
	if !(len(searchParam.RequiredExcelSkills) == 0) {

	excelSkillLoop:
		for _, jobInformation := range jobInformationListWithRequiredLanguage {
			if !jobInformation.ExcelSkill.Valid {
				continue
			}

			for _, excelSkill := range searchParam.RequiredExcelSkills {
				if !excelSkill.Valid {
					continue
				}

				if excelSkill == jobInformation.ExcelSkill {
					/**
					求人の必要excelスキルが合致したタイミングで次の求人情報に移動して
					ダブりが発生しないように制御
					*/
					jobInformationListWithRequiredExcelSkill = append(jobInformationListWithRequiredExcelSkill, jobInformation)
					continue excelSkillLoop
				}
			}
		}
	}

	// 必要excelスキルのいずれかも入っていない場合
	if len(searchParam.RequiredExcelSkills) == 0 {
		jobInformationListWithRequiredExcelSkill = jobInformationListWithRequiredLanguage
	}

	fmt.Println("必要excelスキル: ", len(jobInformationListWithRequiredExcelSkill))

	// NOTE: wordスキル
	// wordスキルのいずれかが入っている場合
	if !(len(searchParam.RequiredWordSkills) == 0) {

	wordSkillLoop:
		for _, jobInformation := range jobInformationListWithRequiredExcelSkill {
			if !jobInformation.WordSkill.Valid {
				continue
			}

			for _, wordSkill := range searchParam.RequiredWordSkills {
				if !wordSkill.Valid {
					continue
				}

				if wordSkill == jobInformation.WordSkill {
					/**
					求人の必要wordスキルが合致したタイミングで次の求人へ移る
					ダブりが発生しないように制御
					*/
					jobInformationListWithRequiredWordSkill = append(jobInformationListWithRequiredWordSkill, jobInformation)
					continue wordSkillLoop
				}
			}
		}
	}

	// 必要wordスキルのいずれかも入っていない場合
	if len(searchParam.RequiredWordSkills) == 0 {
		jobInformationListWithRequiredWordSkill = jobInformationListWithRequiredExcelSkill
	}

	fmt.Println("必要wordスキル: ", len(jobInformationListWithRequiredWordSkill))

	// NOTE: powerpointスキル
	// powerpointスキルのいずれかが入っている場合
	if !(len(searchParam.RequiredPowerPointSkills) == 0) {

	powerPointSkillLoop:
		for _, jobInformation := range jobInformationListWithRequiredWordSkill {
			if !jobInformation.PowerPointSkill.Valid {
				continue
			}

			for _, powerPointSkill := range searchParam.RequiredPowerPointSkills {
				if !powerPointSkill.Valid {
					continue
				}

				if powerPointSkill == jobInformation.PowerPointSkill {
					/**
					求人の必要powerpointスキルが合致したタイミングで次の求人へ移る
					ダブりが発生しないように制御
					*/

					jobInformationListWithRequiredPowerPointSkill = append(jobInformationListWithRequiredPowerPointSkill, jobInformation)
					continue powerPointSkillLoop
				}
			}
		}
	}

	// 必要powerpointスキルのいずれかも入っていない場合
	if len(searchParam.RequiredPowerPointSkills) == 0 {
		jobInformationListWithRequiredPowerPointSkill = jobInformationListWithRequiredWordSkill
	}

	fmt.Println("必要powerpointスキル: ", len(jobInformationListWithRequiredPowerPointSkill))

	// NOTE: 必要PC業務ツール（And検索）
	// 必要PCスキル（その他）のいずれかが入っている場合
	if !(len(searchParam.RequiredAnotherPCSkills) == 0) {
		// 検索パラムのmap
		var pcToolParams = make(map[null.Int]bool)
		for _, pcToolParam := range searchParam.RequiredAnotherPCSkills {
			pcToolParams[pcToolParam] = false
		}

	pcToolLoop:
		for _, jobInformation := range jobInformationListWithRequiredPowerPointSkill {
			if len(jobInformation.RequiredConditions) == 0 {
				continue
			}

			/**
			合致した時にpcToolLoopを抜けるためのラベル
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 必要資格のリスト作成
			var pcTools []entity.SendingJobInformationRequiredPCTool
			for _, infoCondition := range jobInformation.RequiredConditions {
				pcTools = append(pcTools, infoCondition.RequiredPCTools...)
			}

			// 求人の業務ツール経験のmap
			var pcToolDatas = make(map[null.Int]bool)

			// パラムのループ
			for _, pcToolParam := range searchParam.RequiredAnotherPCSkills {
				if !pcToolParam.Valid {
					continue
				}

				for _, pcTool := range pcTools {
					if pcToolParam == pcTool.Tool {
						pcToolParams[pcToolParam] = true
						pcToolDatas[pcTool.Tool] = true
					}
				}

				if reflect.DeepEqual(pcToolParams, pcToolDatas) {
					jobInformationListWithRequiredAnotherPCSkill = append(jobInformationListWithRequiredAnotherPCSkill, jobInformation)
					continue pcToolLoop
				}
			}
		}
	}

	// 必要PC業務ツールのいずれかも入っていない場合
	if len(searchParam.RequiredAnotherPCSkills) == 0 {
		jobInformationListWithRequiredAnotherPCSkill = jobInformationListWithRequiredPowerPointSkill
	}

	fmt.Println("必要PC業務ツール: ", len(jobInformationListWithRequiredAnotherPCSkill))

	// NOTE: 開発言語
	// 開発言語のいずれかが入っている場合
	if !(len(searchParam.RequiredDevelopmentLanguages) == 0) {

		// 検索パラムのmap
		var developmentLanguageParams = make(map[null.Int]bool)
		for _, developmentLanguageParam := range searchParam.RequiredDevelopmentLanguages {
			developmentLanguageParams[developmentLanguageParam] = false
		}

	developmentLanguageLoop:
		for _, jobInformation := range jobInformationListWithRequiredAnotherPCSkill {
			if len(jobInformation.RequiredConditions) == 0 {
				continue
			}

			/**
			合致した時にdevelopmentLanguageLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 必要開発言語のリスト作成
			var developments []entity.SendingJobInformationRequiredExperienceDevelopment
			for _, infoCondition := range jobInformation.RequiredConditions {
				developments = append(developments, infoCondition.RequiredExperienceDevelopments...)
			}

			// 求人の必要開発言語のmap
			var developmentLanguageDatas = make(map[null.Int]bool)

			// パラムのループ
			for _, developmentLanguageParam := range searchParam.RequiredDevelopmentLanguages {
				if !developmentLanguageParam.Valid {
					continue
				}

				for _, development := range developments {
					for _, developmentType := range development.ExperienceDevelopmentTypes {
						if development.DevelopmentCategory == null.NewInt(0, true) && developmentLanguageParam == developmentType.DevelopmentType {
							developmentLanguageParams[developmentLanguageParam] = true
							developmentLanguageDatas[developmentType.DevelopmentType] = true
						}
					}
				}

				if reflect.DeepEqual(developmentLanguageParams, developmentLanguageDatas) {
					jobInformationListWithRequiredDevelopmentLanguage = append(jobInformationListWithRequiredDevelopmentLanguage, jobInformation)
					continue developmentLanguageLoop
				}
			}
		}
	}

	// 開発言語のいずれかも入っていない場合
	if len(searchParam.RequiredDevelopmentLanguages) == 0 {
		jobInformationListWithRequiredDevelopmentLanguage = jobInformationListWithRequiredAnotherPCSkill
	}

	fmt.Println("開発言語: ", len(jobInformationListWithRequiredDevelopmentLanguage))

	// NOTE: 開発OS
	// 開発OSのいずれかが入っている場合
	if !(len(searchParam.RequiredDevelopmentOS) == 0) {

		// 検索パラムのmap
		var developmentOSParams = make(map[null.Int]bool)
		for _, developmentOSParam := range searchParam.RequiredDevelopmentOS {
			developmentOSParams[developmentOSParam] = false
		}

	developmentOSLoop:
		for _, jobInformation := range jobInformationListWithRequiredDevelopmentLanguage {
			if len(jobInformation.RequiredConditions) == 0 {
				continue
			}

			/**
			合致した時にdevelopmentOSLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 必要開発OSのリスト作成
			var developments []entity.SendingJobInformationRequiredExperienceDevelopment
			for _, infoCondition := range jobInformation.RequiredConditions {
				developments = append(developments, infoCondition.RequiredExperienceDevelopments...)
			}

			// 求人の必要開発OSのmap
			var developmentOSDatas = make(map[null.Int]bool)

			for _, developmentOSParam := range searchParam.RequiredDevelopmentOS {
				if !developmentOSParam.Valid {
					continue
				}

				for _, development := range developments {
					for _, developmentType := range development.ExperienceDevelopmentTypes {
						if development.DevelopmentCategory == null.NewInt(1, true) && developmentOSParam == developmentType.DevelopmentType {
							developmentOSParams[developmentOSParam] = true
							developmentOSDatas[developmentType.DevelopmentType] = true
						}
					}
				}

				if reflect.DeepEqual(developmentOSParams, developmentOSDatas) {
					jobInformationListWithRequiredDevelopmentOS = append(jobInformationListWithRequiredDevelopmentOS, jobInformation)
					continue developmentOSLoop
				}
			}
		}
	}

	// 開発OSのいずれかも入っていない場合
	if len(searchParam.RequiredDevelopmentOS) == 0 {
		jobInformationListWithRequiredDevelopmentOS = jobInformationListWithRequiredDevelopmentLanguage
	}

	fmt.Println("開発OS: ", len(jobInformationListWithRequiredDevelopmentOS))

	// Note: 転勤有無
	// 転勤有無がある場合
	if !(len(searchParam.TransferTypes) == 0) {
	transferLoop:
		for _, jobInformation := range jobInformationListWithRequiredDevelopmentOS {
			for _, transfer := range searchParam.TransferTypes {
				if !transfer.Valid {
					continue
				}
				if transfer == jobInformation.Transfer {
					// 求人の転勤有無が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithTransfer = append(jobInformationListWithTransfer, jobInformation)
					continue transferLoop
				}
			}
		}
	}

	// 転勤有無が無い場合
	if len(searchParam.TransferTypes) == 0 {
		jobInformationListWithTransfer = jobInformationListWithRequiredDevelopmentOS
	}

	fmt.Println("転勤有無: ", len(jobInformationListWithTransfer))

	// Note: 休日タイプ
	// 休日タイプがある場合
	if !(len(searchParam.HolidayTypes) == 0) {
	holidayLoop:
		for _, jobInformation := range jobInformationListWithTransfer {
			for _, holidayType := range searchParam.HolidayTypes {
				if !holidayType.Valid {
					continue
				}
				if holidayType == jobInformation.HolidayType {
					// 求人の休日休暇が一つでも合致していればcontinueで次の求人へ
					jobInformationListWithHolidayType = append(jobInformationListWithHolidayType, jobInformation)
					continue holidayLoop
				}
			}
		}
	}

	// 休日タイプが無い場合
	if len(searchParam.HolidayTypes) == 0 {
		jobInformationListWithHolidayType = jobInformationListWithTransfer
	}

	fmt.Println("休日タイプ: ", len(jobInformationListWithHolidayType))

	// Note: 企業規模の絞り込みは企業の従業員数（単体）と比較する
	// 企業規模がある場合
	if !(len(searchParam.CompanyScaleTypes) == 0) {
	companyScaleLoop:
		for _, jobInformation := range jobInformationListWithHolidayType {
			if jobInformation.EmployeeNumberSingle == null.NewInt(0, false) {
				continue
			}
			for _, companyScale := range searchParam.CompanyScaleTypes {
				if !companyScale.Valid {
					continue
				}

				if companyScale == null.NewInt(0, true) {
					// 10名未満の場合
					if jobInformation.EmployeeNumberSingle.Int64 < 10 {
						jobInformationListWithCompanyScale = append(jobInformationListWithCompanyScale, jobInformation)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(1, true) {
					// 10名以上100名未満の場合
					if jobInformation.EmployeeNumberSingle.Int64 >= 10 && jobInformation.EmployeeNumberSingle.Int64 < 100 {
						jobInformationListWithCompanyScale = append(jobInformationListWithCompanyScale, jobInformation)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(2, true) {
					// 100名以上200名未満の場合
					if jobInformation.EmployeeNumberSingle.Int64 >= 100 && jobInformation.EmployeeNumberSingle.Int64 < 200 {
						jobInformationListWithCompanyScale = append(jobInformationListWithCompanyScale, jobInformation)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(3, true) {
					// 200名以上1000名未満の場合
					if jobInformation.EmployeeNumberSingle.Int64 >= 200 && jobInformation.EmployeeNumberSingle.Int64 < 1000 {
						jobInformationListWithCompanyScale = append(jobInformationListWithCompanyScale, jobInformation)
						continue companyScaleLoop
					}
				} else if companyScale == null.NewInt(4, true) {
					// 1000名以上の場合
					if jobInformation.EmployeeNumberSingle.Int64 >= 1000 {
						jobInformationListWithCompanyScale = append(jobInformationListWithCompanyScale, jobInformation)
						continue companyScaleLoop
					}
				}
			}
		}
	}

	// 企業規模が無い場合
	if len(searchParam.CompanyScaleTypes) == 0 {
		jobInformationListWithCompanyScale = jobInformationListWithHolidayType
	}

	// 求人の特徴
	if !(len(searchParam.Features) == 0) {
	featureLoop:
		for _, jobInformation := range jobInformationListWithCompanyScale {
			for _, jobFeature := range jobInformation.Features {
				for _, feature := range searchParam.Features {
					if !feature.Valid {
						continue
					}

					// 検索パラムが「0: 業界未経験OK」or「1: 職種未経験OK」で求人が「2: 業界・職種未経験OK」
					var inexperienced = (feature == null.NewInt(0, true) || feature == null.NewInt(1, true)) && jobFeature.Feature == null.NewInt(2, true)

					// 【ヒッt条件】
					// 求人の特徴が合致している or 検索パラムが「0: 業界未経験OK」or「1: 職種未経験OK」で求人が「2: 業界・職種未経験OK」
					if feature == jobFeature.Feature || inexperienced {
						jobInformationListWithFeature = append(jobInformationListWithFeature, jobInformation)
						continue featureLoop
					}
				}
			}
		}
	}

	// 求人の特徴が無い場合
	if len(searchParam.Features) == 0 {
		jobInformationListWithFeature = jobInformationListWithCompanyScale
	}

	return jobInformationListWithFeature, nil
}

func licenceSearchSendingJobInformation(jobInformationList []*entity.SendingJobInformation, licenseParams []null.Int) []*entity.SendingJobInformation {
	var (
		infoListAtLicense []*entity.SendingJobInformation
	)
	// ライセンスのマッチング
licenseLoop:
	for _, info := range jobInformationList {
		for _, licenseParam := range licenseParams {
			for _, condition := range info.RequiredConditions {
				for _, infoLicense := range condition.RequiredLicenses {

					if !licenseParam.Valid {
						infoListAtLicense = append(infoListAtLicense, info)
						continue licenseLoop
						// 完全一致 ||
						// 求職者の保有資格：普通自動車免許（MT） && 求人の必要資格：普通自動車免許（MT）＋普通自動車免許（AT）がヒット
					} else if (infoLicense.License == licenseParam) ||
						(infoLicense.License.Int64 == 4803 && licenseParam.Int64 == 4805) ||
						(infoLicense.License.Int64 == 1205 && (licenseParam.Int64 == 1203 || licenseParam.Int64 == 1204)) ||
						(infoLicense.License.Int64 == 1204 && licenseParam.Int64 == 1205) ||
						(infoLicense.License.Int64 == 1206 && licenseParam.Int64 == 1207) ||
						(infoLicense.License.Int64 == 1211 && licenseParam.Int64 == 1212) ||
						(infoLicense.License.Int64 == 1218 && licenseParam.Int64 == 1219) ||
						(infoLicense.License.Int64 == 1224 && (licenseParam.Int64 == 1223 || licenseParam.Int64 == 1222)) ||
						(infoLicense.License.Int64 == 1223 && licenseParam.Int64 == 1224) ||
						(infoLicense.License.Int64 == 1238 && licenseParam.Int64 == 1239) ||
						(infoLicense.License.Int64 == 1301 && licenseParam.Int64 == 1302) ||
						(infoLicense.License.Int64 == 1305 && licenseParam.Int64 == 1306) ||
						(infoLicense.License.Int64 == 1315 && licenseParam.Int64 == 1316) ||
						(infoLicense.License.Int64 == 1322 && (licenseParam.Int64 == 1321 || licenseParam.Int64 == 1320)) ||
						(infoLicense.License.Int64 == 1321 && licenseParam.Int64 == 1322) ||
						(infoLicense.License.Int64 == 1326 && (licenseParam.Int64 == 1324 || licenseParam.Int64 == 1325 || licenseParam.Int64 == 1323)) ||
						(infoLicense.License.Int64 == 1325 && (licenseParam.Int64 == 1324 || licenseParam.Int64 == 1323)) ||
						(infoLicense.License.Int64 == 1324 && licenseParam.Int64 == 1323) ||
						(infoLicense.License.Int64 == 1401 && licenseParam.Int64 == 1402) ||
						(infoLicense.License.Int64 == 1404 && licenseParam.Int64 == 1405) ||
						(infoLicense.License.Int64 == 1406 && licenseParam.Int64 == 1407) ||
						(infoLicense.License.Int64 == 1408 && licenseParam.Int64 == 1409) ||
						(infoLicense.License.Int64 == 1547 && (licenseParam.Int64 == 1546 || licenseParam.Int64 == 1545)) ||
						(infoLicense.License.Int64 == 1546 && licenseParam.Int64 == 1547) ||
						(infoLicense.License.Int64 == 1550 && (licenseParam.Int64 == 1549 || licenseParam.Int64 == 1548)) ||
						(infoLicense.License.Int64 == 1549 && licenseParam.Int64 == 1550) ||
						(infoLicense.License.Int64 == 1551 && licenseParam.Int64 == 1552) ||
						(infoLicense.License.Int64 == 1553 && licenseParam.Int64 == 1554) ||
						(infoLicense.License.Int64 == 1563 && licenseParam.Int64 == 1564) ||
						(infoLicense.License.Int64 == 1605 && licenseParam.Int64 == 1606) ||
						(infoLicense.License.Int64 == 1610 && licenseParam.Int64 == 1611) ||
						(infoLicense.License.Int64 == 2202 && licenseParam.Int64 == 2203) ||
						(infoLicense.License.Int64 == 2306 && (licenseParam.Int64 == 2305 || licenseParam.Int64 == 2304)) ||
						(infoLicense.License.Int64 == 2305 && licenseParam.Int64 == 2306) ||
						(infoLicense.License.Int64 == 2314 && licenseParam.Int64 == 2313) ||
						(infoLicense.License.Int64 == 2409 && (licenseParam.Int64 == 2408 || licenseParam.Int64 == 2407)) ||
						(infoLicense.License.Int64 == 2408 && licenseParam.Int64 == 2409) ||
						(infoLicense.License.Int64 == 2517 && (licenseParam.Int64 == 2516 || licenseParam.Int64 == 2515)) ||
						(infoLicense.License.Int64 == 2516 && licenseParam.Int64 == 2517) ||
						(infoLicense.License.Int64 == 2701 && licenseParam.Int64 == 2702) ||
						(infoLicense.License.Int64 == 2703 && licenseParam.Int64 == 2704) ||
						(infoLicense.License.Int64 == 2902 && licenseParam.Int64 == 2901) ||
						(infoLicense.License.Int64 == 2901 && (licenseParam.Int64 == 2902 || licenseParam.Int64 == 2903)) ||
						(infoLicense.License.Int64 == 2905 && licenseParam.Int64 == 2904) ||
						(infoLicense.License.Int64 == 2904 && (licenseParam.Int64 == 2905 || licenseParam.Int64 == 2906)) ||
						(infoLicense.License.Int64 == 2908 && licenseParam.Int64 == 2907) ||
						(infoLicense.License.Int64 == 2907 && (licenseParam.Int64 == 2908 || licenseParam.Int64 == 2909)) ||
						(infoLicense.License.Int64 == 2913 && licenseParam.Int64 == 2912) ||
						(infoLicense.License.Int64 == 2916 && licenseParam.Int64 == 2917) ||
						(infoLicense.License.Int64 == 3001 && licenseParam.Int64 == 3002) ||
						(infoLicense.License.Int64 == 3009 && licenseParam.Int64 == 3010) ||
						(infoLicense.License.Int64 == 3014 && (licenseParam.Int64 == 3013 || licenseParam.Int64 == 3012)) ||
						(infoLicense.License.Int64 == 3013 && licenseParam.Int64 == 3014) ||
						(infoLicense.License.Int64 == 3017 && (licenseParam.Int64 == 3016 || licenseParam.Int64 == 3015)) ||
						(infoLicense.License.Int64 == 3016 && licenseParam.Int64 == 3017) ||
						(infoLicense.License.Int64 == 3107 && licenseParam.Int64 == 3108) ||
						(infoLicense.License.Int64 == 3201 && licenseParam.Int64 == 3202) ||
						(infoLicense.License.Int64 == 3309 && (licenseParam.Int64 == 3308 || licenseParam.Int64 == 3307)) ||
						(infoLicense.License.Int64 == 3308 && licenseParam.Int64 == 3309) ||
						(infoLicense.License.Int64 == 3316 && (licenseParam.Int64 == 3314 || licenseParam.Int64 == 3315 || licenseParam.Int64 == 3313)) ||
						(infoLicense.License.Int64 == 3315 && (licenseParam.Int64 == 3314 || licenseParam.Int64 == 3313)) ||
						(infoLicense.License.Int64 == 3314 && licenseParam.Int64 == 3313) ||
						(infoLicense.License.Int64 == 3325 && licenseParam.Int64 == 3326) ||
						(infoLicense.License.Int64 == 3327 && licenseParam.Int64 == 3328) ||
						(infoLicense.License.Int64 == 3329 && licenseParam.Int64 == 3330) ||
						(infoLicense.License.Int64 == 3331 && licenseParam.Int64 == 3332) ||
						(infoLicense.License.Int64 == 3337 && (licenseParam.Int64 == 3336 || licenseParam.Int64 == 3335)) ||
						(infoLicense.License.Int64 == 3336 && licenseParam.Int64 == 3337) ||
						(infoLicense.License.Int64 == 3338 && licenseParam.Int64 == 3339) ||
						(infoLicense.License.Int64 == 3403 && (licenseParam.Int64 == 3402 || licenseParam.Int64 == 3401)) ||
						(infoLicense.License.Int64 == 3402 && licenseParam.Int64 == 3403) ||
						(infoLicense.License.Int64 == 3407 && (licenseParam.Int64 == 3406 || licenseParam.Int64 == 3405)) ||
						(infoLicense.License.Int64 == 3406 && licenseParam.Int64 == 3407) ||
						(infoLicense.License.Int64 == 3515 && (licenseParam.Int64 == 3514 || licenseParam.Int64 == 3513)) ||
						(infoLicense.License.Int64 == 3514 && licenseParam.Int64 == 3515) ||
						(infoLicense.License.Int64 == 3523 && (licenseParam.Int64 == 3522 || licenseParam.Int64 == 3521)) ||
						(infoLicense.License.Int64 == 3522 && licenseParam.Int64 == 3523) ||
						(infoLicense.License.Int64 == 3527 && (licenseParam.Int64 == 3525 || licenseParam.Int64 == 3526 || licenseParam.Int64 == 3524)) ||
						(infoLicense.License.Int64 == 3526 && (licenseParam.Int64 == 3524 || licenseParam.Int64 == 3525)) ||
						(infoLicense.License.Int64 == 3525 && licenseParam.Int64 == 3524) ||
						(infoLicense.License.Int64 == 3617 && (licenseParam.Int64 == 3616 || licenseParam.Int64 == 3615)) ||
						(infoLicense.License.Int64 == 3616 && licenseParam.Int64 == 3617) ||
						(infoLicense.License.Int64 == 3621 && (licenseParam.Int64 == 3620 || licenseParam.Int64 == 3619)) ||
						(infoLicense.License.Int64 == 3620 && licenseParam.Int64 == 3621) ||
						(infoLicense.License.Int64 == 3625 && (licenseParam.Int64 == 3624 || licenseParam.Int64 == 3623)) ||
						(infoLicense.License.Int64 == 3624 && licenseParam.Int64 == 3625) ||
						(infoLicense.License.Int64 == 3629 && (licenseParam.Int64 == 3627 || licenseParam.Int64 == 3628 || licenseParam.Int64 == 3626)) ||
						(infoLicense.License.Int64 == 3628 && (licenseParam.Int64 == 3627 || licenseParam.Int64 == 3626)) ||
						(infoLicense.License.Int64 == 3627 && licenseParam.Int64 == 3626) ||
						(infoLicense.License.Int64 == 3632 && (licenseParam.Int64 == 3631 || licenseParam.Int64 == 3630)) ||
						(infoLicense.License.Int64 == 3631 && licenseParam.Int64 == 3632) ||
						(infoLicense.License.Int64 == 3636 && (licenseParam.Int64 == 3635 || licenseParam.Int64 == 3634)) ||
						(infoLicense.License.Int64 == 3635 && licenseParam.Int64 == 3636) ||
						(infoLicense.License.Int64 == 3708 && licenseParam.Int64 == 3709) ||
						(infoLicense.License.Int64 == 3722 && (licenseParam.Int64 == 3721 || licenseParam.Int64 == 3720)) ||
						(infoLicense.License.Int64 == 3721 && licenseParam.Int64 == 3722) ||
						(infoLicense.License.Int64 == 3730 && (licenseParam.Int64 == 3729 || licenseParam.Int64 == 3728)) ||
						(infoLicense.License.Int64 == 3729 && licenseParam.Int64 == 3730) ||
						(infoLicense.License.Int64 == 3801 && licenseParam.Int64 == 3802) ||
						(infoLicense.License.Int64 == 3814 && (licenseParam.Int64 == 3813 || licenseParam.Int64 == 3812)) ||
						(infoLicense.License.Int64 == 3813 && licenseParam.Int64 == 3814) ||
						(infoLicense.License.Int64 == 3815 && licenseParam.Int64 == 3816) ||
						(infoLicense.License.Int64 == 3817 && licenseParam.Int64 == 3818) ||
						(infoLicense.License.Int64 == 3843 && (licenseParam.Int64 == 3842 || licenseParam.Int64 == 3841)) ||
						(infoLicense.License.Int64 == 3842 && licenseParam.Int64 == 3843) ||
						(infoLicense.License.Int64 == 3846 && (licenseParam.Int64 == 3845 || licenseParam.Int64 == 3844)) ||
						(infoLicense.License.Int64 == 3845 && licenseParam.Int64 == 3846) ||
						(infoLicense.License.Int64 == 3851 && (licenseParam.Int64 == 3850 || licenseParam.Int64 == 3849)) ||
						(infoLicense.License.Int64 == 3850 && licenseParam.Int64 == 3851) ||
						(infoLicense.License.Int64 == 3854 && (licenseParam.Int64 == 3853 || licenseParam.Int64 == 3852)) ||
						(infoLicense.License.Int64 == 3853 && licenseParam.Int64 == 3854) ||
						(infoLicense.License.Int64 == 3858 && (licenseParam.Int64 == 3856 || licenseParam.Int64 == 3857 || licenseParam.Int64 == 3855)) ||
						(infoLicense.License.Int64 == 3857 && (licenseParam.Int64 == 3856 || licenseParam.Int64 == 3855)) ||
						(infoLicense.License.Int64 == 3856 && licenseParam.Int64 == 3855) ||
						(infoLicense.License.Int64 == 3862 && (licenseParam.Int64 == 3860 || licenseParam.Int64 == 3861 || licenseParam.Int64 == 3859)) ||
						(infoLicense.License.Int64 == 3861 && (licenseParam.Int64 == 3860 || licenseParam.Int64 == 3859)) ||
						(infoLicense.License.Int64 == 3860 && licenseParam.Int64 == 3859) ||
						(infoLicense.License.Int64 == 3866 && (licenseParam.Int64 == 3864 || licenseParam.Int64 == 3865 || licenseParam.Int64 == 3863)) ||
						(infoLicense.License.Int64 == 3865 && (licenseParam.Int64 == 3864 || licenseParam.Int64 == 3863)) ||
						(infoLicense.License.Int64 == 3864 && licenseParam.Int64 == 3863) ||
						(infoLicense.License.Int64 == 3875 && (licenseParam.Int64 == 3874 || licenseParam.Int64 == 3873)) ||
						(infoLicense.License.Int64 == 3874 && licenseParam.Int64 == 3875) ||
						(infoLicense.License.Int64 == 3884 && (licenseParam.Int64 == 3883 || licenseParam.Int64 == 3882)) ||
						(infoLicense.License.Int64 == 3883 && licenseParam.Int64 == 3884) ||
						(infoLicense.License.Int64 == 3903 && (licenseParam.Int64 == 3902 || licenseParam.Int64 == 3901)) ||
						(infoLicense.License.Int64 == 3902 && licenseParam.Int64 == 3903) ||
						(infoLicense.License.Int64 == 3906 && (licenseParam.Int64 == 3905 || licenseParam.Int64 == 3904)) ||
						(infoLicense.License.Int64 == 3905 && licenseParam.Int64 == 3906) ||
						(infoLicense.License.Int64 == 4103 && (licenseParam.Int64 == 4102 || licenseParam.Int64 == 4101)) ||
						(infoLicense.License.Int64 == 4102 && licenseParam.Int64 == 4103) ||
						(infoLicense.License.Int64 == 4104 && licenseParam.Int64 == 4105) ||
						(infoLicense.License.Int64 == 4113 && (licenseParam.Int64 == 4112 || licenseParam.Int64 == 4111)) ||
						(infoLicense.License.Int64 == 4112 && licenseParam.Int64 == 4113) ||
						(infoLicense.License.Int64 == 4311 && (licenseParam.Int64 == 4310 || licenseParam.Int64 == 4309)) ||
						(infoLicense.License.Int64 == 4310 && licenseParam.Int64 == 4311) ||
						(infoLicense.License.Int64 == 4319 && (licenseParam.Int64 == 4318 || licenseParam.Int64 == 4317)) ||
						(infoLicense.License.Int64 == 4318 && licenseParam.Int64 == 4319) ||
						(infoLicense.License.Int64 == 4320 && licenseParam.Int64 == 4321) ||
						(infoLicense.License.Int64 == 4323 && licenseParam.Int64 == 4324) ||
						(infoLicense.License.Int64 == 4325 && licenseParam.Int64 == 4326) ||
						(infoLicense.License.Int64 == 4330 && (licenseParam.Int64 == 4329 || licenseParam.Int64 == 4328)) ||
						(infoLicense.License.Int64 == 4329 && licenseParam.Int64 == 4330) ||
						(infoLicense.License.Int64 == 4334 && licenseParam.Int64 == 4335) ||
						(infoLicense.License.Int64 == 4338 && (licenseParam.Int64 == 4337 || licenseParam.Int64 == 4336)) ||
						(infoLicense.License.Int64 == 4337 && licenseParam.Int64 == 4338) ||
						(infoLicense.License.Int64 == 4406 && (licenseParam.Int64 == 4405 || licenseParam.Int64 == 4404)) ||
						(infoLicense.License.Int64 == 4405 && licenseParam.Int64 == 4406) ||
						(infoLicense.License.Int64 == 4412 && (licenseParam.Int64 == 4411 || licenseParam.Int64 == 4410)) ||
						(infoLicense.License.Int64 == 4411 && licenseParam.Int64 == 4412) ||
						(infoLicense.License.Int64 == 4501 && licenseParam.Int64 == 4502) ||
						(infoLicense.License.Int64 == 4503 && licenseParam.Int64 == 4504) ||
						(infoLicense.License.Int64 == 4505 && licenseParam.Int64 == 4506) ||
						(infoLicense.License.Int64 == 4507 && licenseParam.Int64 == 4508) ||
						(infoLicense.License.Int64 == 4510 && licenseParam.Int64 == 4511) ||
						(infoLicense.License.Int64 == 4512 && licenseParam.Int64 == 4513) ||
						(infoLicense.License.Int64 == 4514 && licenseParam.Int64 == 4515) ||
						(infoLicense.License.Int64 == 4516 && licenseParam.Int64 == 4517) ||
						(infoLicense.License.Int64 == 4518 && licenseParam.Int64 == 4519) ||
						(infoLicense.License.Int64 == 4520 && licenseParam.Int64 == 4521) ||
						(infoLicense.License.Int64 == 4603 && (licenseParam.Int64 == 4602 || licenseParam.Int64 == 4601)) ||
						(infoLicense.License.Int64 == 4602 && licenseParam.Int64 == 4603) ||
						(infoLicense.License.Int64 == 4801 && licenseParam.Int64 == 4823) ||
						(infoLicense.License.Int64 == 4803 && (licenseParam.Int64 == 4801 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4823 || licenseParam.Int64 == 4803 || licenseParam.Int64 == 4804)) ||
						(infoLicense.License.Int64 == 4804 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4821 || licenseParam.Int64 == 4823 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4806 || licenseParam.Int64 == 4807 || licenseParam.Int64 == 4808 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4803 || licenseParam.Int64 == 4805 || licenseParam.Int64 == 4809 || licenseParam.Int64 == 4810)) ||
						(infoLicense.License.Int64 == 4806 && licenseParam.Int64 == 4807) ||
						(infoLicense.License.Int64 == 4808 && licenseParam.Int64 == 4898) ||
						(infoLicense.License.Int64 == 4809 && (licenseParam.Int64 == 4807 || licenseParam.Int64 == 4806)) ||
						(infoLicense.License.Int64 == 4810 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4821 || licenseParam.Int64 == 4823 || licenseParam.Int64 == 4804 || licenseParam.Int64 == 4805 || licenseParam.Int64 == 4810 || licenseParam.Int64 == 4807 || licenseParam.Int64 == 4806 || licenseParam.Int64 == 4898 || licenseParam.Int64 == 4808 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4811)) ||
						(infoLicense.License.Int64 == 4812 && licenseParam.Int64 == 4811) ||
						(infoLicense.License.Int64 == 4817 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818)) ||
						(infoLicense.License.Int64 == 4818 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801)) ||
						(infoLicense.License.Int64 == 4819 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819)) ||
						(infoLicense.License.Int64 == 4820 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820)) ||
						(infoLicense.License.Int64 == 4821 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4823)) ||
						(infoLicense.License.Int64 == 4822 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4822)) ||
						(infoLicense.License.Int64 == 4823 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4821 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4822)) ||
						(infoLicense.License.Int64 == 4907 && licenseParam.Int64 == 4908) ||
						(infoLicense.License.Int64 == 5003 && (licenseParam.Int64 == 5002 || licenseParam.Int64 == 5001)) ||
						(infoLicense.License.Int64 == 5002 && licenseParam.Int64 == 5003) ||
						(infoLicense.License.Int64 == 5004 && licenseParam.Int64 == 5005) ||
						(infoLicense.License.Int64 == 5006 && licenseParam.Int64 == 5007) ||
						(infoLicense.License.Int64 == 5010 && (licenseParam.Int64 == 5009 || licenseParam.Int64 == 5008)) ||
						(infoLicense.License.Int64 == 5009 && licenseParam.Int64 == 5010) ||
						(infoLicense.License.Int64 == 5103 && licenseParam.Int64 == 5104) ||
						(infoLicense.License.Int64 == 5107 && (licenseParam.Int64 == 5106 || licenseParam.Int64 == 5105)) ||
						(infoLicense.License.Int64 == 5106 && licenseParam.Int64 == 5107) ||
						(infoLicense.License.Int64 == 5108 && licenseParam.Int64 == 5109) ||
						(infoLicense.License.Int64 == 5113 && (licenseParam.Int64 == 5111 || licenseParam.Int64 == 5112 || licenseParam.Int64 == 5110)) ||
						(infoLicense.License.Int64 == 5112 && (licenseParam.Int64 == 5110 || licenseParam.Int64 == 5111)) ||
						(infoLicense.License.Int64 == 5111 && licenseParam.Int64 == 5110) ||
						(infoLicense.License.Int64 == 5114 && licenseParam.Int64 == 5115) ||
						(infoLicense.License.Int64 == 5116 && licenseParam.Int64 == 5117) ||
						(infoLicense.License.Int64 == 5121 && (licenseParam.Int64 == 5119 || licenseParam.Int64 == 5120 || licenseParam.Int64 == 5118)) ||
						(infoLicense.License.Int64 == 5120 && (licenseParam.Int64 == 5119 || licenseParam.Int64 == 5118)) ||
						(infoLicense.License.Int64 == 5119 && licenseParam.Int64 == 5118) ||
						(infoLicense.License.Int64 == 5204 && licenseParam.Int64 == 5205) ||
						(infoLicense.License.Int64 == 5208 && (licenseParam.Int64 == 5207 || licenseParam.Int64 == 5206)) ||
						(infoLicense.License.Int64 == 5207 && licenseParam.Int64 == 5208) ||
						(infoLicense.License.Int64 == 5303 && licenseParam.Int64 == 5304) ||
						(infoLicense.License.Int64 == 5307 && licenseParam.Int64 == 5308) ||
						(infoLicense.License.Int64 == 5404 && (licenseParam.Int64 == 5402 || licenseParam.Int64 == 5403 || licenseParam.Int64 == 5401)) ||
						(infoLicense.License.Int64 == 5403 && (licenseParam.Int64 == 5401 || licenseParam.Int64 == 5402)) ||
						(infoLicense.License.Int64 == 5402 && licenseParam.Int64 == 5401) ||
						(infoLicense.License.Int64 == 5407 && (licenseParam.Int64 == 5406 || licenseParam.Int64 == 5405)) ||
						(infoLicense.License.Int64 == 5406 && licenseParam.Int64 == 5407) ||
						(infoLicense.License.Int64 == 5408 && licenseParam.Int64 == 5409) ||
						(infoLicense.License.Int64 == 5412 && (licenseParam.Int64 == 5411 || licenseParam.Int64 == 5410)) ||
						(infoLicense.License.Int64 == 5411 && licenseParam.Int64 == 5412) ||
						(infoLicense.License.Int64 == 5413 && licenseParam.Int64 == 5414) ||
						(infoLicense.License.Int64 == 5419 && (licenseParam.Int64 == 5417 || licenseParam.Int64 == 5418 || licenseParam.Int64 == 5416)) ||
						(infoLicense.License.Int64 == 5418 && (licenseParam.Int64 == 5417 || licenseParam.Int64 == 5416)) ||
						(infoLicense.License.Int64 == 5417 && licenseParam.Int64 == 5416) ||
						(infoLicense.License.Int64 == 5420 && licenseParam.Int64 == 5421) ||
						(infoLicense.License.Int64 == 5425 && (licenseParam.Int64 == 5424 || licenseParam.Int64 == 5423)) ||
						(infoLicense.License.Int64 == 5424 && licenseParam.Int64 == 5425) ||
						(infoLicense.License.Int64 == 5426 && licenseParam.Int64 == 5427) ||
						(infoLicense.License.Int64 == 5428 && licenseParam.Int64 == 5429) ||
						(infoLicense.License.Int64 == 5433 && (licenseParam.Int64 == 5432 || licenseParam.Int64 == 5431)) ||
						(infoLicense.License.Int64 == 5432 && licenseParam.Int64 == 5433) ||
						(infoLicense.License.Int64 == 5440 && (licenseParam.Int64 == 5439 || licenseParam.Int64 == 5438)) ||
						(infoLicense.License.Int64 == 5439 && licenseParam.Int64 == 5440) ||
						(infoLicense.License.Int64 == 5441 && licenseParam.Int64 == 5442) ||
						(infoLicense.License.Int64 == 5445 && (licenseParam.Int64 == 5444 || licenseParam.Int64 == 5443)) ||
						(infoLicense.License.Int64 == 5444 && licenseParam.Int64 == 5445) ||
						(infoLicense.License.Int64 == 5501 && licenseParam.Int64 == 5502) ||
						(infoLicense.License.Int64 == 5504 && licenseParam.Int64 == 5505) ||
						(infoLicense.License.Int64 == 5513 && (licenseParam.Int64 == 5511 || licenseParam.Int64 == 5512 || licenseParam.Int64 == 5510)) ||
						(infoLicense.License.Int64 == 5512 && (licenseParam.Int64 == 5511 || licenseParam.Int64 == 5510)) ||
						(infoLicense.License.Int64 == 5511 && licenseParam.Int64 == 5510) ||
						(infoLicense.License.Int64 == 5515 && licenseParam.Int64 == 5516) ||
						(infoLicense.License.Int64 == 5517 && licenseParam.Int64 == 5518) ||
						(infoLicense.License.Int64 == 5519 && licenseParam.Int64 == 5520) ||
						(infoLicense.License.Int64 == 5603 && (licenseParam.Int64 == 5602 || licenseParam.Int64 == 5601)) ||
						(infoLicense.License.Int64 == 5602 && licenseParam.Int64 == 5603) ||
						(infoLicense.License.Int64 == 5606 && (licenseParam.Int64 == 5605 || licenseParam.Int64 == 5604)) ||
						(infoLicense.License.Int64 == 5605 && licenseParam.Int64 == 5606) ||
						(infoLicense.License.Int64 == 5611 && (licenseParam.Int64 == 5609 || licenseParam.Int64 == 5610 || licenseParam.Int64 == 5608)) ||
						(infoLicense.License.Int64 == 5610 && (licenseParam.Int64 == 5608 || licenseParam.Int64 == 5609)) ||
						(infoLicense.License.Int64 == 5609 && licenseParam.Int64 == 5608) ||
						(infoLicense.License.Int64 == 5614 && (licenseParam.Int64 == 5613 || licenseParam.Int64 == 5612)) ||
						(infoLicense.License.Int64 == 5613 && licenseParam.Int64 == 5614) ||
						(infoLicense.License.Int64 == 5617 && (licenseParam.Int64 == 5616 || licenseParam.Int64 == 5615)) ||
						(infoLicense.License.Int64 == 5616 && licenseParam.Int64 == 5617) ||
						(infoLicense.License.Int64 == 5618 && licenseParam.Int64 == 5619) ||
						(infoLicense.License.Int64 == 5622 && (licenseParam.Int64 == 5621 || licenseParam.Int64 == 5620)) ||
						(infoLicense.License.Int64 == 5621 && licenseParam.Int64 == 5622) ||
						(infoLicense.License.Int64 == 5623 && licenseParam.Int64 == 5624) ||
						(infoLicense.License.Int64 == 5703 && (licenseParam.Int64 == 5702 || licenseParam.Int64 == 5701)) ||
						(infoLicense.License.Int64 == 5702 && licenseParam.Int64 == 5703) ||
						(infoLicense.License.Int64 == 5707 && (licenseParam.Int64 == 5705 || licenseParam.Int64 == 5706 || licenseParam.Int64 == 5704)) ||
						(infoLicense.License.Int64 == 5706 && (licenseParam.Int64 == 5705 || licenseParam.Int64 == 5704)) ||
						(infoLicense.License.Int64 == 5705 && licenseParam.Int64 == 5704) ||
						(infoLicense.License.Int64 == 5714 && (licenseParam.Int64 == 5713 || licenseParam.Int64 == 5712)) ||
						(infoLicense.License.Int64 == 5713 && licenseParam.Int64 == 5714) ||
						(infoLicense.License.Int64 == 5717 && (licenseParam.Int64 == 5716 || licenseParam.Int64 == 5715)) ||
						(infoLicense.License.Int64 == 5716 && licenseParam.Int64 == 5717) ||
						(infoLicense.License.Int64 == 5722 && (licenseParam.Int64 == 5720 || licenseParam.Int64 == 5721 || licenseParam.Int64 == 5719)) ||
						(infoLicense.License.Int64 == 5721 && (licenseParam.Int64 == 5720 || licenseParam.Int64 == 5719)) ||
						(infoLicense.License.Int64 == 5720 && licenseParam.Int64 == 5719) ||
						(infoLicense.License.Int64 == 5804 && licenseParam.Int64 == 5805) ||
						(infoLicense.License.Int64 == 5808 && licenseParam.Int64 == 5809) ||
						(infoLicense.License.Int64 == 5810 && licenseParam.Int64 == 5811) ||
						(infoLicense.License.Int64 == 5813 && licenseParam.Int64 == 5815) ||
						(infoLicense.License.Int64 == 5814 && licenseParam.Int64 == 5816) ||
						(infoLicense.License.Int64 == 5896 && licenseParam.Int64 == 5817) ||
						(infoLicense.License.Int64 == 5908 && (licenseParam.Int64 == 5907 || licenseParam.Int64 == 5906)) ||
						(infoLicense.License.Int64 == 5907 && licenseParam.Int64 == 5908) ||
						(infoLicense.License.Int64 == 6003 && (licenseParam.Int64 == 6002 || licenseParam.Int64 == 6001)) ||
						(infoLicense.License.Int64 == 6002 && licenseParam.Int64 == 6003) ||
						(infoLicense.License.Int64 == 6005 && licenseParam.Int64 == 6006) ||
						(infoLicense.License.Int64 == 6008 && licenseParam.Int64 == 6009) ||
						(infoLicense.License.Int64 == 6010 && licenseParam.Int64 == 6011) ||
						(infoLicense.License.Int64 == 6012 && licenseParam.Int64 == 6013) ||
						(infoLicense.License.Int64 == 6014 && licenseParam.Int64 == 6015) ||
						(infoLicense.License.Int64 == 6101 && licenseParam.Int64 == 6102) ||
						(infoLicense.License.Int64 == 6103 && licenseParam.Int64 == 6104) ||
						(infoLicense.License.Int64 == 6107 && (licenseParam.Int64 == 6106 || licenseParam.Int64 == 6105)) ||
						(infoLicense.License.Int64 == 6106 && licenseParam.Int64 == 6107) ||
						(infoLicense.License.Int64 == 6110 && (licenseParam.Int64 == 6109 || licenseParam.Int64 == 6108)) ||
						(infoLicense.License.Int64 == 6109 && licenseParam.Int64 == 6110) ||
						(infoLicense.License.Int64 == 6113 && (licenseParam.Int64 == 6112 || licenseParam.Int64 == 6111)) ||
						(infoLicense.License.Int64 == 6112 && licenseParam.Int64 == 6113) ||
						(infoLicense.License.Int64 == 6114 && licenseParam.Int64 == 6115) ||
						(infoLicense.License.Int64 == 6116 && licenseParam.Int64 == 6117) ||
						(infoLicense.License.Int64 == 6118 && licenseParam.Int64 == 6119) ||
						(infoLicense.License.Int64 == 6202 && licenseParam.Int64 == 6202) ||
						(infoLicense.License.Int64 == 6206 && licenseParam.Int64 == 6207) ||
						(infoLicense.License.Int64 == 6208 && licenseParam.Int64 == 6209) ||
						(infoLicense.License.Int64 == 6212 && licenseParam.Int64 == 6213) ||
						(infoLicense.License.Int64 == 6214 && licenseParam.Int64 == 6215) ||
						(infoLicense.License.Int64 == 6301 && licenseParam.Int64 == 6302) ||
						(infoLicense.License.Int64 == 6303 && licenseParam.Int64 == 6304) ||
						(infoLicense.License.Int64 == 6305 && licenseParam.Int64 == 6306) ||
						(infoLicense.License.Int64 == 6404 && (licenseParam.Int64 == 6402 || licenseParam.Int64 == 6403 || licenseParam.Int64 == 6401)) ||
						(infoLicense.License.Int64 == 6403 && (licenseParam.Int64 == 6402 || licenseParam.Int64 == 6401)) ||
						(infoLicense.License.Int64 == 6402 && licenseParam.Int64 == 6401) ||
						(infoLicense.License.Int64 == 6405 && licenseParam.Int64 == 6406) ||
						(infoLicense.License.Int64 == 6503 && licenseParam.Int64 == 6504) ||
						(infoLicense.License.Int64 == 6507 && licenseParam.Int64 == 6508) ||
						(infoLicense.License.Int64 == 6509 && licenseParam.Int64 == 6510) ||
						(infoLicense.License.Int64 == 6601 && licenseParam.Int64 == 6602) ||
						(infoLicense.License.Int64 == 6603 && licenseParam.Int64 == 6604) ||
						(infoLicense.License.Int64 == 6608 && (licenseParam.Int64 == 6607 || licenseParam.Int64 == 6606)) ||
						(infoLicense.License.Int64 == 6607 && licenseParam.Int64 == 6608) ||
						(infoLicense.License.Int64 == 6611 && (licenseParam.Int64 == 6610 || licenseParam.Int64 == 6609)) ||
						(infoLicense.License.Int64 == 6610 && licenseParam.Int64 == 6611) ||
						(infoLicense.License.Int64 == 6612 && licenseParam.Int64 == 6613) ||
						(infoLicense.License.Int64 == 6620 && (licenseParam.Int64 == 6618 || licenseParam.Int64 == 6619 || licenseParam.Int64 == 6617)) ||
						(infoLicense.License.Int64 == 6619 && (licenseParam.Int64 == 6619 || licenseParam.Int64 == 6617)) ||
						(infoLicense.License.Int64 == 6618 && licenseParam.Int64 == 6617) ||
						(infoLicense.License.Int64 == 6624 && (licenseParam.Int64 == 6622 || licenseParam.Int64 == 6623 || licenseParam.Int64 == 6621)) ||
						(infoLicense.License.Int64 == 6623 && (licenseParam.Int64 == 6622 || licenseParam.Int64 == 6621)) ||
						(infoLicense.License.Int64 == 6622 && licenseParam.Int64 == 6621) ||
						(infoLicense.License.Int64 == 6625 && (licenseParam.Int64 == 6626 || licenseParam.Int64 == 6627)) ||
						(infoLicense.License.Int64 == 6626 && licenseParam.Int64 == 6627) ||
						(infoLicense.License.Int64 == 6630 && licenseParam.Int64 == 6631) ||
						(infoLicense.License.Int64 == 6634 && (licenseParam.Int64 == 6633 || licenseParam.Int64 == 6632)) ||
						(infoLicense.License.Int64 == 6633 && licenseParam.Int64 == 6634) ||
						(infoLicense.License.Int64 == 6635 && licenseParam.Int64 == 6636) ||
						(infoLicense.License.Int64 == 6703 && (licenseParam.Int64 == 6702 || licenseParam.Int64 == 6701)) ||
						(infoLicense.License.Int64 == 6702 && licenseParam.Int64 == 6703) ||
						(infoLicense.License.Int64 == 6803 && (licenseParam.Int64 == 6802 || licenseParam.Int64 == 6801)) ||
						(infoLicense.License.Int64 == 6802 && licenseParam.Int64 == 6803) ||
						(infoLicense.License.Int64 == 6806 && (licenseParam.Int64 == 6805 || licenseParam.Int64 == 6804)) ||
						(infoLicense.License.Int64 == 6805 && licenseParam.Int64 == 6806) ||
						(infoLicense.License.Int64 == 6808 && licenseParam.Int64 == 6809) ||
						(infoLicense.License.Int64 == 6814 && (licenseParam.Int64 == 6813 || licenseParam.Int64 == 6812)) ||
						(infoLicense.License.Int64 == 6813 && licenseParam.Int64 == 6814) ||
						(infoLicense.License.Int64 == 6908 && licenseParam.Int64 == 6909) ||
						(infoLicense.License.Int64 == 7001 && licenseParam.Int64 == 7002) ||
						(infoLicense.License.Int64 == 7005 && licenseParam.Int64 == 7006) ||
						(infoLicense.License.Int64 == 7009 && (licenseParam.Int64 == 7008 || licenseParam.Int64 == 7007)) ||
						(infoLicense.License.Int64 == 7008 && licenseParam.Int64 == 7009) ||
						(infoLicense.License.Int64 == 7010 && licenseParam.Int64 == 7011) ||
						(infoLicense.License.Int64 == 7013 && licenseParam.Int64 == 7014) ||
						(infoLicense.License.Int64 == 7015 && licenseParam.Int64 == 7016) ||
						(infoLicense.License.Int64 == 7018 && licenseParam.Int64 == 7019) ||
						(infoLicense.License.Int64 == 7020 && licenseParam.Int64 == 7021) ||
						(infoLicense.License.Int64 == 7023 && licenseParam.Int64 == 7024) ||
						(infoLicense.License.Int64 == 7103 && licenseParam.Int64 == 7104) ||
						(infoLicense.License.Int64 == 7105 && licenseParam.Int64 == 7106) ||
						(infoLicense.License.Int64 == 7109 && (licenseParam.Int64 == 7108 || licenseParam.Int64 == 7107)) ||
						(infoLicense.License.Int64 == 7108 && licenseParam.Int64 == 7109) ||
						(infoLicense.License.Int64 == 7110 && licenseParam.Int64 == 7111) ||
						(infoLicense.License.Int64 == 7113 && licenseParam.Int64 == 7114) ||
						(infoLicense.License.Int64 == 7117 && (licenseParam.Int64 == 7116 || licenseParam.Int64 == 7115)) ||
						(infoLicense.License.Int64 == 7116 && licenseParam.Int64 == 7117) ||
						(infoLicense.License.Int64 == 7118 && licenseParam.Int64 == 7119) ||
						(infoLicense.License.Int64 == 7120 && licenseParam.Int64 == 7121) ||
						(infoLicense.License.Int64 == 7122 && licenseParam.Int64 == 7123) ||
						(infoLicense.License.Int64 == 7124 && licenseParam.Int64 == 7125) ||
						(infoLicense.License.Int64 == 7127 && licenseParam.Int64 == 7128) ||
						(infoLicense.License.Int64 == 7201 && licenseParam.Int64 == 7202) ||
						(infoLicense.License.Int64 == 7302 && licenseParam.Int64 == 7303) ||
						(infoLicense.License.Int64 == 7305 && licenseParam.Int64 == 7306) ||
						(infoLicense.License.Int64 == 7406 && licenseParam.Int64 == 7407) {
						infoListAtLicense = append(infoListAtLicense, info)
						continue licenseLoop
					}
				}
			}
		}
	}

	fmt.Println("資格: ", infoListAtLicense)

	return infoListAtLicense
}
