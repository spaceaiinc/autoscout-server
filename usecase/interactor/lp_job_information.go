package interactor

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

/***********************************************************************************************************************/
// LP用　API
type GetSearchJobInformationCountByLPDiagnosisInput struct {
	SearchParam entity.DiagnosisParam
}

type GetSearchJobInformationCountByLPDiagnosisOutput struct {
	//	応募資格に合致した求人: 条件で絞った時の求人数
	Count uint

	// 未経験職で応募資格に合致した求人: 上記の求人の中で求人の特徴に「職種未経験OK」or「業界・職種未経験OK」が選択されていて、かつ「求職者が経験していない職種」の求人数
	GuaranteedInterviewCount uint

	// 年収期待値: 求人の年収下限と年収上限の中央値の平均値
	ExpectedIncome uint
}

func (i *JobInformationInteractorImpl) GetSearchJobInformationCountByLPDiagnosis(input GetSearchJobInformationCountByLPDiagnosisInput) (GetSearchJobInformationCountByLPDiagnosisOutput, error) {
	var (
		output GetSearchJobInformationCountByLPDiagnosisOutput
		err    error
		// jobInformationList   = input.SearchParam.JobInformationList
		// jobInformationIDList []uint
		// motoyuiAgentID       uint = 1

		// 年収の合計
		sumIncome uint

		// features []*entity.JobInformationFeature
	)

	// 絞り込み検索処理
	jobInformationList, err := searchJobInformationByLPJobSeekerForDiagnosisParam(input.SearchParam.JobInformationList, input.SearchParam)
	if err != nil {
		return output, err
	}

	for _, jobInformation := range jobInformationList {
		/**
		絞り込み通過した場合の処理
		*/
		// 応募資格に合致した求人: 条件で絞った時の求人数 ※1つでも未経験職種が含まれていたらカウント
		output.Count++

		// 未経験職で応募資格に合致した求人: 上記の求人の中で求人の特徴に「職種未経験OK」or「業界・職種未経験OK」が選択されていて、かつ「求職者が経験していない職種」の求人数
		if jobInformation.IsGuaranteedInterview {
			output.GuaranteedInterviewCount++
		}

		// 年収期待値: 求人の年収下限と年収上限の中央値の平均値
		var medianIncome int64
		if jobInformation.UnderIncome.Valid && jobInformation.OverIncome.Valid {
			medianIncome = (jobInformation.OverIncome.Int64 + jobInformation.UnderIncome.Int64) / 2
		} else if jobInformation.OverIncome.Valid {
			medianIncome = jobInformation.OverIncome.Int64
		} else if jobInformation.UnderIncome.Valid {
			medianIncome = jobInformation.UnderIncome.Int64
		}
		// 年収合計に追加
		sumIncome += uint(medianIncome)
	}

	// 元結求人のアクティブ求人を取得
	// jobInformationList, err = i.jobInformationRepository.GetActiveAllByAgentIDAndDiagnosisParamWithoutExternal(motoyuiAgentID, input.SearchParam)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// for _, jobInformation := range jobInformationList {
	// 	jobInformationIDList = append(jobInformationIDList, jobInformation.ID)
	// }

	// // 求人リストから、IDが合致する子テーブルの情報を取得
	// // 検索で使用しないかつリスト表示しない情報は、取得しない
	// requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByJobInformationIDList(jobInformationIDList)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByJobInformationIDListAndLicenceTypeList(jobInformationIDList, input.SearchParam.Licenses)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByJobInformationIDList(jobInformationIDList)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByJobInformationIDListAndLanguageTypeList(jobInformationIDList, input.SearchParam.Languages)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByJobInformationIDList(jobInformationIDList)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByJobInformationIDListAndIndustryList(jobInformationIDList, input.SearchParam.Industries)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByJobInformationIDListAndOccupationList(jobInformationIDList, input.SearchParam.ExperienceOccupations)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	// if input.SearchParam.FirstLanguage == null.NewInt(1, true) || input.SearchParam.FirstLanguage == null.NewInt(2, true) {
	// 	features, err = i.jobInfoFeatureRepository.GetByJobInformationIDList(jobInformationIDList)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return output, err
	// 	}
	// }

	// jobInfoLoop:
	// 	for _, jobInformation := range jobInformationList {

	// 		/**
	// 		国籍(第一言語)
	// 		*/
	// 		isMatchNationality := false

	// 		// 求人側で未入力の場合
	// 		if !jobInformation.Nationality.Valid ||

	// 			// 求人側で不問の場合
	// 			jobInformation.Nationality == null.NewInt(99, true) ||

	// 			// 求人側の募集国籍が日本国籍のみで求職者が日本国籍
	// 			(jobInformation.Nationality == null.NewInt(0, true) && input.SearchParam.FirstLanguage == null.NewInt(0, true)) ||

	// 			// 求人側の募集国籍が外国籍のみで求職者が外国籍
	// 			(jobInformation.Nationality == null.NewInt(1, true) && (input.SearchParam.FirstLanguage == null.NewInt(1, true) || input.SearchParam.FirstLanguage == null.NewInt(2, true))) {
	// 			isMatchNationality = true
	// 		}

	// 		// 求人の特徴に外国籍可がある場合も外国籍OK
	// 		if input.SearchParam.FirstLanguage == null.NewInt(1, true) || input.SearchParam.FirstLanguage == null.NewInt(2, true) {
	// 			for _, f := range features {
	// 				// 求人の特徴に外国籍化の求人があるかをチェック
	// 				if jobInformation.ID == f.JobInformationID && f.Feature == null.NewInt(21, true) {
	// 					isMatchNationality = true
	// 					break
	// 				}
	// 			}
	// 		}

	// 		if !isMatchNationality {
	// 			continue jobInfoLoop
	// 		}

	// 		/**
	// 		必要条件
	// 		*/
	// 		// for _, condition := range requiredConditions {
	// 		// 	if jobInformation.ID == condition.JobInformationID {
	// 		// 		for _, rl := range requiredLicenses {
	// 		// 			if condition.ID == rl.ConditionID {
	// 		// 				condition.RequiredLicenses = append(condition.RequiredLicenses, *rl)
	// 		// 			}
	// 		// 		}

	// 		// 		for _, rl := range requiredLanguages {
	// 		// 			if condition.ID == rl.ConditionID {
	// 		// 				// 言語タイプ
	// 		// 				for _, languageType := range requiredLanguageTypes {
	// 		// 					if rl.ID == languageType.LanguageID {
	// 		// 						rl.LanguageTypes = append(rl.LanguageTypes, *languageType)
	// 		// 					}
	// 		// 				}

	// 		// 				condition.RequiredLanguages = *rl
	// 		// 			}
	// 		// 		}

	// 		// 		for _, rej := range requiredExperienceJobs {
	// 		// 			if condition.ID == rej.ConditionID {

	// 		// 				// 業界
	// 		// 				for _, industry := range requiredExperienceIndustries {
	// 		// 					if rej.ID == industry.ExperienceJobID {
	// 		// 						rej.ExperienceIndustries = append(rej.ExperienceIndustries, *industry)
	// 		// 					}
	// 		// 				}

	// 		// 				// 職種
	// 		// 				for _, occupation := range requiredExperienceOccupations {
	// 		// 					if rej.ID == occupation.ExperienceJobID {
	// 		// 						rej.ExperienceOccupations = append(rej.ExperienceOccupations, *occupation)
	// 		// 					}
	// 		// 				}

	// 		// 				condition.RequiredExperienceJobs = *rej
	// 		// 			}
	// 		// 		}

	// 		// 		jobInformation.RequiredConditions = append(jobInformation.RequiredConditions, *condition)
	// 		// 	}
	// 		// }

	// 		/**
	// 		必要条件の絞り込み
	// 		- 必要条件の入力がない求人は一括OK
	// 		- 共通条件マッチとパターン条件のうち1がマッチすればOK
	// 		*/
	// 		if len(jobInformation.RequiredConditions) > 0 {
	// 			var (
	// 				// 共通要件のマッチ
	// 				isMatchCommonCondition bool
	// 				// パターン条件のマッチ
	// 				isMatchPatternCondition bool
	// 			)
	// 			for _, rc := range jobInformation.RequiredConditions {
	// 				// LP側で運転免許選択したら資格の運転免許を選択した状態にする
	// 				var (
	// 					isRequiredIndustry   bool // 業界の条件があるか
	// 					isRequiredOccupation bool // 職種の条件があるか
	// 					isRequiredLanguage   bool // 言語の条件があるか
	// 					isRequiredLicense    bool // 資格の条件があるか

	// 					isMatchIndustry   bool // 業界が合致している
	// 					isMatchOccupation bool // 職種が合致している
	// 					isMatchLanguage   bool // 言語が合致している
	// 					isMatchLicense    bool // 資格が合致している
	// 				)

	// 				if len(rc.RequiredExperienceJobs.ExperienceIndustries) > 0 {
	// 					isRequiredIndustry = true
	// 				}
	// 				if len(rc.RequiredExperienceJobs.ExperienceOccupations) > 0 {
	// 					isRequiredOccupation = true
	// 				}
	// 				if len(rc.RequiredLanguages.LanguageTypes) > 0 {
	// 					isRequiredLanguage = true
	// 				}
	// 				if len(rc.RequiredLicenses) > 0 {
	// 					isRequiredLicense = true
	// 				}

	// 				// 必要条件がない場合は
	// 				if !isRequiredIndustry && !isRequiredOccupation && !isRequiredLanguage && !isRequiredLicense {
	// 					isMatchPatternCondition = true
	// 					continue
	// 				}

	// 				/************************************************************/
	// 				// 業界
	// 				//
	// 				var (
	// 					ej                            = rc.RequiredExperienceJobs // 必要経験業職種
	// 					requiredExperienceMonth int64 = 0                         // 求人の必要経験月数
	// 				)

	// 				// 必要月数を計算
	// 				if ej.ExperienceYear.Valid && ej.ExperienceMonth.Valid {
	// 					// 年数と月数の入力がある場合
	// 					requiredExperienceMonth = (ej.ExperienceYear.Int64 * 6) + ej.ExperienceMonth.Int64
	// 				} else if ej.ExperienceYear.Valid && !ej.ExperienceMonth.Valid {
	// 					// 年数のみ入力がある場合
	// 					requiredExperienceMonth = ej.ExperienceYear.Int64 * 6
	// 				} else if !ej.ExperienceYear.Valid && ej.ExperienceMonth.Valid {
	// 					// 月数のみ入力がある場合
	// 					requiredExperienceMonth = ej.ExperienceMonth.Int64
	// 				}

	// 			industriesLoop:
	// 				for _, requiredIndustry := range ej.ExperienceIndustries {
	// 					for _, seekerIndustry := range input.SearchParam.Industries {
	// 						// 業界が合致していて、かつ「求人の業職種経験の必要経験年数」が未入力、または「求職者の業界経験年数」が「求人の業職種経験の必要経験年数」以上の場合
	// 						if seekerIndustry.Industry == requiredIndustry.ExperienceIndustry && (requiredExperienceMonth == 0 || requiredExperienceMonth <= seekerIndustry.ExperienceMonth.Int64) {
	// 							isMatchIndustry = true
	// 							break industriesLoop
	// 						}
	// 					}
	// 				}

	// 			occupationLoop:
	// 				for _, requiredOccupation := range ej.ExperienceOccupations {
	// 					for _, seekerOccupation := range input.SearchParam.AllExperienceOccupations {
	// 						// 求職者の職種の経験月数
	// 						var seekerExperienceMonth int64 = 0

	// 						// seekerOccupation.ExperienceYear: {0: 1年未満, 1: 1年以上, 2: 2年以上, ..., 10: 10年以上}
	// 						if seekerOccupation.ExperienceYear == null.NewInt(0, true) {
	// 							seekerExperienceMonth = 6 // 1年未満は半年としてカウント
	// 						} else if seekerOccupation.ExperienceYear.Valid && seekerOccupation.ExperienceYear.Int64 > 0 && seekerOccupation.ExperienceYear.Int64 <= 10 {
	// 							seekerExperienceMonth = seekerOccupation.ExperienceYear.Int64 * 12
	// 						}

	// 						// 職種が合致していて、かつ「求人の業職種経験の必要経験年数」が未入力、または「求職者の職種経験年数」が「求人の業職種経験の必要経験年数」以上の場合
	// 						if seekerOccupation.Occupation == requiredOccupation.ExperienceOccupation && (requiredExperienceMonth == 0 || requiredExperienceMonth <= seekerExperienceMonth) {
	// 							isMatchOccupation = true
	// 							break occupationLoop
	// 						}
	// 					}
	// 				}

	// 				/************************************************************/
	// 				// 言語
	// 				//
	// 			languageLoop:
	// 				for _, languageType := range rc.RequiredLanguages.LanguageTypes {
	// 					for _, seekerLanguage := range input.SearchParam.Languages {

	// 						// 言語タイプが合致していて、かつ必要なレベルも合致している場合
	// 						if seekerLanguage.LanguageType == languageType.LanguageType && seekerLanguage.LanguageLevel == rc.RequiredLanguages.LanguageLevel {
	// 							isMatchLanguage = true
	// 							break languageLoop
	// 						}
	// 					}
	// 				}

	// 				/************************************************************/
	// 				// 資格
	// 				//
	// 			licenseLoop:
	// 				for _, license := range rc.RequiredLicenses {
	// 					for _, seekerLicense := range input.SearchParam.Licenses {
	// 						// 言語タイプが合致していて、かつ必要なレベルも合致している場合
	// 						if seekerLicense == license.License {
	// 							isMatchLicense = true
	// 							break licenseLoop
	// 						}
	// 					}
	// 				}

	// 				var (
	// 					// 1つ
	// 					conditions1 = isRequiredIndustry && isMatchIndustry
	// 					conditions2 = isRequiredOccupation && isMatchOccupation
	// 					conditions3 = isRequiredLanguage && isMatchLanguage
	// 					conditions4 = isRequiredLicense && isMatchLicense

	// 					// 2つ
	// 					conditions5  = isRequiredIndustry && isRequiredOccupation && isMatchIndustry && isMatchOccupation
	// 					conditions6  = isRequiredIndustry && isRequiredLanguage && isMatchIndustry && isMatchLanguage
	// 					conditions7  = isRequiredIndustry && isRequiredLicense && isMatchIndustry && isMatchLicense
	// 					conditions8  = isRequiredOccupation && isRequiredLanguage && isMatchOccupation && isMatchLanguage
	// 					conditions9  = isRequiredOccupation && isRequiredLicense && isMatchOccupation && isMatchLicense
	// 					conditions10 = isRequiredLanguage && isRequiredLicense && isMatchLanguage && isMatchLicense

	// 					// 3つ
	// 					conditions11 = isRequiredIndustry && isRequiredOccupation && isRequiredLanguage && isMatchIndustry && isMatchOccupation && isMatchLanguage
	// 					conditions12 = isRequiredIndustry && isRequiredOccupation && isRequiredLicense && isMatchIndustry && isMatchOccupation && isMatchLicense
	// 					conditions13 = isRequiredIndustry && isRequiredLanguage && isRequiredLicense && isMatchIndustry && isMatchLanguage && isMatchLicense
	// 					conditions14 = isRequiredOccupation && isRequiredLanguage && isRequiredLicense && isMatchOccupation && isMatchLanguage && isMatchLicense

	// 					// 4つ
	// 					conditions15 = isRequiredIndustry && isRequiredOccupation && isRequiredLanguage && isRequiredLicense && isMatchIndustry && isMatchOccupation && isMatchLanguage && isMatchLicense
	// 				)

	// 				// 必要条件全てクリアしている場合
	// 				if conditions1 || conditions2 || conditions3 || conditions4 || conditions5 ||
	// 					conditions6 || conditions7 || conditions8 || conditions9 || conditions10 ||
	// 					conditions11 || conditions12 || conditions13 || conditions14 || conditions15 {
	// 					// 必要条件の場合
	// 					if rc.IsCommon {
	// 						isMatchCommonCondition = true
	// 						// 共通条件のみの場合は共通のみマッチしていればOK
	// 						if len(jobInformation.RequiredConditions) == 1 {
	// 							isMatchPatternCondition = true
	// 						}

	// 						// パターン条件の場合
	// 					} else {
	// 						isMatchPatternCondition = true
	// 						// 共通も含めてマッチしていればOK
	// 						if isMatchCommonCondition {
	// 							break
	// 						}
	// 					}
	// 					continue

	// 					// 必要条件を満たしておらず、かつ比較対象が共通条件の場合は求人のループに戻る
	// 				} else if rc.IsCommon {
	// 					continue jobInfoLoop
	// 				}
	// 			}
	// 			// いずれかの条件にマッチしなかった場合はスキップ
	// 			if !isMatchPatternCondition || !isMatchCommonCondition {
	// 				continue jobInfoLoop
	// 			}
	// 		}

	// 		/**
	// 		絞り込み通過した場合の処理
	// 		*/
	// 		// 応募資格に合致した求人: 条件で絞った時の求人数 ※1つでも未経験職種が含まれていたらカウント
	// 		output.Count++

	// 		// 未経験職で応募資格に合致した求人: 上記の求人の中で求人の特徴に「職種未経験OK」or「業界・職種未経験OK」が選択されていて、かつ「求職者が経験していない職種」の求人数
	// 		if jobInformation.IsGuaranteedInterview {
	// 			output.GuaranteedInterviewCount++
	// 		}

	// 		// 年収期待値: 求人の年収下限と年収上限の中央値の平均値
	// 		var medianIncome int64
	// 		if jobInformation.UnderIncome.Valid && jobInformation.OverIncome.Valid {
	// 			medianIncome = (jobInformation.OverIncome.Int64 + jobInformation.UnderIncome.Int64) / 2
	// 		} else if jobInformation.OverIncome.Valid {
	// 			medianIncome = jobInformation.OverIncome.Int64
	// 		} else if jobInformation.UnderIncome.Valid {
	// 			medianIncome = jobInformation.UnderIncome.Int64
	// 		}
	// 		// 年収合計に追加
	// 		sumIncome += uint(medianIncome)
	// 	}

	fmt.Println("------------------------------")
	fmt.Println("応募資格に合致した求人:", output.Count)
	fmt.Println("------------------------------")

	fmt.Println("------------------------------")
	fmt.Println("面談確約した求人:", output.GuaranteedInterviewCount)
	fmt.Println("------------------------------")

	if output.Count != 0 {
		output.ExpectedIncome = sumIncome / output.Count
	}
	fmt.Println("------------------------------")
	fmt.Println("年収期待値:", output.ExpectedIncome)
	fmt.Println("------------------------------")

	return output, nil
}

type GetJobInformationListForDiagnosisOutput struct {
	JobInformationList []*entity.JobInformationForDiagnosis
}

func (i *JobInformationInteractorImpl) GetJobInformationListForDiagnosis() (GetJobInformationListForDiagnosisOutput, error) {
	var (
		output         GetJobInformationListForDiagnosisOutput
		motoyuiAgentID uint = 1
	)

	// 並行処理用のチャネルとエラーチャネル
	type result struct {
		data interface{}
		err  error
	}

	resultChan := make(chan result, 8)

	// 各リポジトリからデータを取得するゴルーチンを起動
	jobInformationList, err := i.jobInformationRepository.GetActiveAllByAgentIDWithoutExternal(motoyuiAgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobInformationIDList := make([]uint, 0, len(jobInformationList))
	for _, jobInformation := range jobInformationList {
		jobInformationIDList = append(jobInformationIDList, jobInformation.ID)
	}

	go func() {
		requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: requiredConditions, err: err}
	}()
	go func() {
		requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: requiredLicenses, err: err}
	}()
	go func() {
		requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: requiredLanguages, err: err}
	}()
	go func() {
		requiredLanguageType, err := i.jobInfoRequiredLanguageTypeRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: requiredLanguageType, err: err}
	}()
	go func() {
		requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: requiredExperienceJobs, err: err}
	}()
	go func() {
		requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: requiredExperienceIndustries, err: err}
	}()
	go func() {
		requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: requiredExperienceOccupations, err: err}
	}()
	go func() {
		features, err := i.jobInfoFeatureRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: features, err: err}
	}()
	go func() {
		prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationIDList(jobInformationIDList)

		resultChan <- result{data: prefectures, err: err}
	}()

	// 結果の受け取りとエラーチェック
	var (
		requiredConditionsResult            []*entity.JobInformationRequiredCondition
		requiredLicensesResult              []*entity.JobInformationRequiredLicense
		requiredLanguagesResult             []*entity.JobInformationRequiredLanguage
		requiredLanguageTypeResult          []*entity.JobInformationRequiredLanguageType
		requiredExperienceJobsResult        []*entity.JobInformationRequiredExperienceJob
		requiredExperienceIndustriesResult  []*entity.JobInformationRequiredExperienceIndustry
		requiredExperienceOccupationsResult []*entity.JobInformationRequiredExperienceOccupation
		featuresResult                      []*entity.JobInformationFeature
		prefecturesResult                   []*entity.JobInformationPrefecture
	)

	for j := 0; j < 9; j++ {
		res := <-resultChan
		if res.err != nil {
			fmt.Println(res.err)
			return output, res.err
		}

		switch v := res.data.(type) {
		case []*entity.JobInformationRequiredCondition:
			requiredConditionsResult = v
			fmt.Println("requiredConditionsResult", len(v))
		case []*entity.JobInformationRequiredLicense:
			requiredLicensesResult = v
			fmt.Println("requiredLicensesResult", len(v))
		case []*entity.JobInformationRequiredLanguage:
			requiredLanguagesResult = v
			fmt.Println("requiredLanguagesResult", len(v))
		case []*entity.JobInformationRequiredLanguageType:
			requiredLanguageTypeResult = v
			fmt.Println("requiredLanguageResult", len(v))
		case []*entity.JobInformationRequiredExperienceJob:
			requiredExperienceJobsResult = v
			fmt.Println("requiredExperienceJobsResult", len(v))
		case []*entity.JobInformationRequiredExperienceIndustry:
			requiredExperienceIndustriesResult = v
			fmt.Println("requiredExperienceIndustriesResult", len(v))
		case []*entity.JobInformationRequiredExperienceOccupation:
			requiredExperienceOccupationsResult = v
			fmt.Println("requiredExperienceOccupationsResult", len(v))
		case []*entity.JobInformationFeature:
			featuresResult = v
			fmt.Println("featuresResult", len(v))
		case []*entity.JobInformationPrefecture:
			prefecturesResult = v
			fmt.Println("prefecturesResult", len(v))
		}
	}

	fmt.Println("ループするよ")
	fmt.Println("1. jobInformationList", len(jobInformationList))
	fmt.Println("2. requiredConditions", len(requiredConditionsResult))
	fmt.Println("3. requiredLicenses", len(requiredLicensesResult))
	fmt.Println("4. requiredLanguages", len(requiredLanguagesResult))
	fmt.Println("5. requiredLanguageType", len(requiredLanguageTypeResult))
	fmt.Println("6. requiredExperienceJobs", len(requiredExperienceJobsResult))
	fmt.Println("7. requiredExperienceIndustries", len(requiredExperienceIndustriesResult))
	fmt.Println("8. requiredExperienceOccupations", len(requiredExperienceOccupationsResult))
	fmt.Println("9. features", len(featuresResult))
	fmt.Println("10. prefectures", len(prefecturesResult))
	for _, jobInformation := range jobInformationList {
		for _, condition := range requiredConditionsResult {
			if jobInformation.ID == condition.JobInformationID {
				for _, rl := range requiredLicensesResult {
					if condition.ID == rl.ConditionID {
						condition.RequiredLicenses = append(condition.RequiredLicenses, *rl)
					}
				}

				for _, rl := range requiredLanguagesResult {
					if condition.ID == rl.ConditionID {
						// 言語タイプ
						for _, languageType := range requiredLanguageTypeResult {
							if rl.ID == languageType.LanguageID {
								rl.LanguageTypes = append(rl.LanguageTypes, *languageType)
							}
						}
						condition.RequiredLanguages = *rl
					}
				}

				for _, rej := range requiredExperienceJobsResult {
					if condition.ID == rej.ConditionID {

						// 業界
						for _, industry := range requiredExperienceIndustriesResult {
							if rej.ID == industry.ExperienceJobID {
								rej.ExperienceIndustries = append(rej.ExperienceIndustries, *industry)
							}
						}

						// 職種
						for _, occupation := range requiredExperienceOccupationsResult {
							if rej.ID == occupation.ExperienceJobID {
								rej.ExperienceOccupations = append(rej.ExperienceOccupations, *occupation)
							}
						}
						condition.RequiredExperienceJobs = *rej
					}
				}

				jobInformation.RequiredConditions = append(jobInformation.RequiredConditions, *condition)
			}
		}

		for _, f := range featuresResult {
			if jobInformation.ID == f.JobInformationID {
				value := entity.JobInformationFeature{
					JobInformationID: f.JobInformationID,
					Feature:          f.Feature,
				}
				jobInformation.Features = append(jobInformation.Features, value)
			}
		}

		for _, p := range prefecturesResult {
			if jobInformation.ID == p.JobInformationID {
				value := entity.JobInformationPrefecture{
					JobInformationID: p.JobInformationID,
					Prefecture:       p.Prefecture,
				}
				jobInformation.Prefectures = append(jobInformation.Prefectures, value)
			}
		}

		jobInformationForDiagnosis := entity.NewJobInformationForDiagnosis(
			jobInformation.BillingAddressID,
			jobInformation.UnderIncome,
			jobInformation.OverIncome,
			jobInformation.Gender,
			jobInformation.Nationality,
			jobInformation.FinalEducation,
			jobInformation.AgeUnder,
			jobInformation.AgeOver,
			jobInformation.JobChange,
			jobInformation.DriverLicence,
			jobInformation.IsGuaranteedInterview,
		)
		jobInformationForDiagnosis.ID = jobInformation.ID
		jobInformationForDiagnosis.UUID = jobInformation.UUID
		jobInformationForDiagnosis.RequiredConditions = jobInformation.RequiredConditions
		jobInformationForDiagnosis.Features = jobInformation.Features
		jobInformationForDiagnosis.Prefectures = jobInformation.Prefectures
		output.JobInformationList = append(output.JobInformationList, jobInformationForDiagnosis)
	}

	return output, nil
}

type GetSearchJobListingListByJobSeekerUUIDInput struct {
	SearchParam entity.SearchMatchingJobListParam
}

type GetSearchJobListingListByJobSeekerUUIDOutput struct {
	// 応募資格に合致した求人リスト
	JobListingList []*entity.JobListing

	//	応募資格に合致した求人: 条件で絞った時の求人数
	Count uint

	// 面接確約の求人数
	CountOfGuaranteedInterview uint

	// 年収期待値: 年収の最小値
	ExpectedUnderIncome uint

	// 年収期待値: 年収の最大値
	ExpectedOverIncome uint

	MaxPageNumber uint
}

func (i *JobInformationInteractorImpl) GetSearchJobListingListByJobSeekerUUID(input GetSearchJobListingListByJobSeekerUUIDInput) (GetSearchJobListingListByJobSeekerUUIDOutput, error) {
	var (
		output GetSearchJobListingListByJobSeekerUUIDOutput

		jobListingListWithDesiredIndustries   []*entity.JobListing
		jobListingListWithDesiredOccupations  []*entity.JobListing
		jobListingListWithPrefectures         []*entity.JobListing
		jobListingListWithFeatures            []*entity.JobListing
		jobListingListWithIsIncome            []*entity.JobListing
		jobListingListWithGuaranteedInterview []*entity.JobListing
	)

	/**
	希望条件の絞り込み
	*/
	// 希望業種
	if len(input.SearchParam.DesiredIndustries) == 0 {
		jobListingListWithDesiredIndustries = input.SearchParam.JobListingList
	} else {
		for _, jobListing := range input.SearchParam.JobListingList {
			isMatched := false
			for _, industry := range input.SearchParam.DesiredIndustries {
				for _, jobIndustry := range jobListing.Industries {
					if industry == jobIndustry.Industry {
						isMatched = true
						break
					}
				}
			}
			if !isMatched {
				continue
			}

			if isMatched {
				jobListingListWithDesiredIndustries = append(jobListingListWithDesiredIndustries, jobListing)
			}
		}
	}

	// 希望職種
	if len(input.SearchParam.DesiredOccupations) == 0 {
		jobListingListWithDesiredOccupations = jobListingListWithDesiredIndustries
	} else {
		for _, jobListing := range jobListingListWithDesiredIndustries {
			isMatched := false
			for _, occupation := range input.SearchParam.DesiredOccupations {
				for _, jobOccupation := range jobListing.Occupations {
					if occupation == jobOccupation.Occupation {
						isMatched = true
						break
					}
				}
			}
			if !isMatched {
				continue
			}

			if isMatched {
				jobListingListWithDesiredOccupations = append(jobListingListWithDesiredOccupations, jobListing)
			}
		}
	}

	// 希望勤務地
	if len(input.SearchParam.Prefectures) == 0 {
		jobListingListWithPrefectures = jobListingListWithDesiredOccupations
	} else {
		for _, jobListing := range jobListingListWithDesiredOccupations {
			isMatched := false
			for _, prefecture := range input.SearchParam.Prefectures {
				for _, jobPrefecture := range jobListing.Prefectures {
					if prefecture == jobPrefecture.Prefecture {
						isMatched = true
						break
					}
				}
			}
			if !isMatched {
				continue
			}

			if isMatched {
				jobListingListWithPrefectures = append(jobListingListWithPrefectures, jobListing)
			}
		}
	}

	// 特徴
	if len(input.SearchParam.Features) == 0 {
		jobListingListWithFeatures = jobListingListWithPrefectures
	} else {
		for _, jobListing := range jobListingListWithPrefectures {
			isMatched := false
			for _, feature := range input.SearchParam.Features {
				for _, jobFeature := range jobListing.Features {
					if feature == jobFeature.Feature {
						isMatched = true
						break
					}
				}
			}
			if !isMatched {
				continue
			}

			if isMatched {
				jobListingListWithFeatures = append(jobListingListWithFeatures, jobListing)
			}
		}
	}

	// 年収
	if input.SearchParam.Income.Valid && input.SearchParam.Income.Int64 != 0 {
		for _, jobListing := range jobListingListWithFeatures {
			if jobListing.UnderIncome.Int64 <= input.SearchParam.Income.Int64 &&
				(jobListing.OverIncome.Int64 >= input.SearchParam.Income.Int64 || !jobListing.OverIncome.Valid) {
				jobListingListWithIsIncome = append(jobListingListWithIsIncome, jobListing)
			}
		}
	} else {
		jobListingListWithIsIncome = jobListingListWithFeatures
	}

	// 面接確約
	if input.SearchParam.IsGuaranteedInterview {
		for _, jobListing := range jobListingListWithIsIncome {
			if jobListing.IsGuaranteedInterview {
				jobListingListWithGuaranteedInterview = append(jobListingListWithGuaranteedInterview, jobListing)
			}
		}
	} else {
		jobListingListWithGuaranteedInterview = jobListingListWithIsIncome
	}

	//	応募資格に合致した求人: 条件で絞った時の求人数 ※1つでも未経験職種が含まれていたらカウント
	output.Count = uint(len(jobListingListWithGuaranteedInterview))

	var (
		maxOverIncome  = int64(math.MinInt64) // 年収上限の最大値
		minUnderIncome = int64(math.MaxInt64) // 年収下限の最小値
	)
	for _, jobListing := range jobListingListWithGuaranteedInterview {

		// 面接確約の求人数
		if jobListing.IsGuaranteedInterview {
			output.CountOfGuaranteedInterview++
		}

		// 年収上限の最大値を更新
		if jobListing.OverIncome.Valid && jobListing.OverIncome.Int64 > maxOverIncome {
			maxOverIncome = jobListing.OverIncome.Int64
		}
		// 年収下限の最小値を更新
		if jobListing.UnderIncome.Valid && jobListing.UnderIncome.Int64 < minUnderIncome {
			minUnderIncome = jobListing.UnderIncome.Int64
		}
	}

	output.ExpectedUnderIncome = uint(minUnderIncome)
	output.ExpectedOverIncome = uint(maxOverIncome)

	// if validLen != 0 {
	// 	output.ExpectedIncome = sumIncome / validLen
	// }
	fmt.Println("------------------------------")
	fmt.Println("応募資格に合致した求人数:", output.Count)
	fmt.Println("------------------------------")
	fmt.Println("面接確約の求人数:", output.CountOfGuaranteedInterview)
	fmt.Println("------------------------------")
	fmt.Println("年収期待値:", output.ExpectedUnderIncome, " ~ ", output.ExpectedOverIncome)
	fmt.Println("------------------------------")

	/**
	ページング
	*/
	output.MaxPageNumber = getJobListingListMaxPage(jobListingListWithGuaranteedInterview)

	if !input.SearchParam.PageNumber.Valid {
		input.SearchParam.PageNumber = null.NewInt(1, true)
	}

	output.JobListingList = getJobListingListWithPage(jobListingListWithGuaranteedInterview, uint(input.SearchParam.PageNumber.Int64))

	return output, nil
}

// 絞り込み検索処理 本体
func searchJobInformationByLPJobSeeker(jobInformationList []*entity.JobInformation, searchParam entity.DiagnosisParam) ([]*entity.JobInformation, error) {
	fmt.Println("ヒット件数: ", len(jobInformationList))

	/**
	性別
	年齢
	国籍
	最終学歴
	経験社数（転職回数）
	必要条件（業界・職種・語学・資格）
	*/

	// 絞り込み項目の結果を代入するための変数を用意
	var (
		jobInformationListWithGender             []*entity.JobInformation
		jobInformationListWithAge                []*entity.JobInformation
		jobInformationListWithNationality        []*entity.JobInformation
		jobInformationListWithFinalEducation     []*entity.JobInformation
		jobInformationListWithCompanyNum         []*entity.JobInformation
		jobInformationListWithRequiredConditions []*entity.JobInformation
	)

	/************ 1. 性別の絞り込み **************/

	for _, jobInformation := range jobInformationList {
		// 求人側で未入力の場合
		if !jobInformation.Gender.Valid ||

			// 求人側で特定の性別のみでない場合（2: 男性尚可 or 3: 女性尚可　or 99: 不問）
			(jobInformation.Gender != null.NewInt(0, true) && jobInformation.Gender != null.NewInt(1, true)) ||

			// 求人側と求職者の性別が一致する場合
			(searchParam.Gender == jobInformation.Gender) {
			jobInformationListWithGender = append(jobInformationListWithGender, jobInformation)
		}
	}

	fmt.Println("性別: ", len(jobInformationListWithGender))

	/************ 2. 年齢の絞り込み **************/

	if searchParam.Birthyear != "" && searchParam.Birthmonth != "" && searchParam.Birthday != "" {
		// 現在の日付を取得
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().In(jst)

		// フォーマットに従って time.Time 型に変換
		dateStr := fmt.Sprintf("%s-%s-%s", searchParam.Birthyear, searchParam.Birthmonth, searchParam.Birthday)
		birthDay, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return nil, err
		}

		// 年齢を計算
		age := now.Year() - birthDay.Year()

		// まだ誕生日を迎えていない場合は、年齢から1を引く
		if now.Month() < birthDay.Month() || (now.Month() == birthDay.Month() && now.Day() < birthDay.Day()) {
			age--
		}

		for _, jobInformation := range jobInformationListWithGender {
			// 募集年齢: 不問（年齢上限、年齢下限が未入力）
			if (!jobInformation.AgeUnder.Valid && !jobInformation.AgeOver.Valid) ||

				// 年齢下限と上限が入力されており、検索パラムの値が下限以上で上限以下の場合
				(jobInformation.AgeUnder.Valid && age >= int(jobInformation.AgeUnder.Int64) && jobInformation.AgeOver.Valid && age <= int(jobInformation.AgeOver.Int64)) ||

				// 年齢下限が入力されている + 検索パラムの値が上限以下 + 年齢上限が未入力
				(jobInformation.AgeUnder.Valid && age >= int(jobInformation.AgeUnder.Int64) && !jobInformation.AgeOver.Valid) ||

				// 年齢上限が入力されている + 検索パラムの値が上限以下 + 年齢下限が未入力
				(jobInformation.AgeOver.Valid && age <= int(jobInformation.AgeOver.Int64) && !jobInformation.AgeUnder.Valid) {

				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
				continue
			}
		}
	} else {
		// 生年月日の入力がない場合は年齢の条件が不問の求人のみ合致
		for _, jobInformation := range jobInformationListWithGender {
			// 募集年齢: 不問（年齢上限、年齢下限が未入力）
			if !jobInformation.AgeUnder.Valid && !jobInformation.AgeOver.Valid {
				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
			}
		}
	}

	fmt.Println("年齢: ", len(jobInformationListWithAge))

	/************ 3. 国籍の絞り込み **************/

	for _, jobInformation := range jobInformationListWithAge {
		var isPossibleForeign bool

		// 求人の特徴に外国籍化の求人があるかをチェック
		for _, feature := range jobInformation.Features {
			if feature.Feature == null.NewInt(21, true) {
				isPossibleForeign = true
				break
			}
		}

		// 求人側で未入力の場合
		if !jobInformation.Nationality.Valid ||

			// 求人側で不問の場合
			jobInformation.Nationality == null.NewInt(99, true) ||

			// 求人側の募集国籍が日本国籍のみで求職者が日本国籍
			(jobInformation.Nationality == null.NewInt(0, true) && searchParam.FirstLanguage == null.NewInt(0, true)) ||

			// 求人側の募集国籍が外国籍のみで求職者が外国籍
			(jobInformation.Nationality == null.NewInt(1, true) && (searchParam.FirstLanguage == null.NewInt(1, true) || searchParam.FirstLanguage == null.NewInt(2, true))) ||

			// 求人側が外国籍可で求職者が外国籍
			(isPossibleForeign && (searchParam.FirstLanguage == null.NewInt(1, true) || searchParam.FirstLanguage == null.NewInt(2, true))) {
			jobInformationListWithNationality = append(jobInformationListWithNationality, jobInformation)
		}
	}

	/************ 4. 最終学歴の絞り込み **************/

	if searchParam.SchoolCategory.Valid {
		for _, jobInformation := range jobInformationListWithNationality {
			// 求人側で未入力の場合: !jobInformation.FinalEducation.Valid
			// 求人側で不問の場合: jobInformation.FinalEducation == null.NewInt(99, true)
			// 求人側で中卒以上の場合: jobInformation.FinalEducation == null.NewInt(0, true)
			// 求人側で高卒以上で求職者が高卒以上の場合: jobInformation.FinalEducation == null.NewInt(1, true) && 1 <= searchParam.SchoolCategory.Int64
			// 求人側で専卒以上（短大除く）で求職者が高卒以上の場合: jobInformation.FinalEducation == null.NewInt(1, true) && 1 <= searchParam.SchoolCategory.Int64

			var (
				// 高卒以上
				IsHighSchoolOrMore = jobInformation.FinalEducation == null.NewInt(1, true) && 1 <= searchParam.SchoolCategory.Int64

				// 専卒以上
				IsVocationalSchoolOrMore = jobInformation.FinalEducation == null.NewInt(2, true) && (searchParam.SchoolCategory == null.NewInt(2, true) || searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 短大卒以上
				IsJuniorCollegeOrMore = jobInformation.FinalEducation == null.NewInt(3, true) && (searchParam.SchoolCategory == null.NewInt(3, true) || searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 専卒と短大卒以上
				IsVocationalSchoolAndJuniorCollegeOrMore = jobInformation.FinalEducation == null.NewInt(4, true) && (searchParam.SchoolCategory == null.NewInt(2, true) || searchParam.SchoolCategory == null.NewInt(3, true) || searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 高専卒以上
				IsTechnicalCollegeOrMore = jobInformation.FinalEducation == null.NewInt(5, true) && (searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 大卒以上
				IsUniversityOrMore = jobInformation.FinalEducation == null.NewInt(6, true) && (searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 院卒以上
				IsGraduateSchoolOrMore = jobInformation.FinalEducation == null.NewInt(7, true) && (searchParam.SchoolCategory == null.NewInt(6, true))
			)

			if !jobInformation.FinalEducation.Valid ||
				jobInformation.FinalEducation == null.NewInt(99, true) ||
				jobInformation.FinalEducation == null.NewInt(0, true) ||
				IsHighSchoolOrMore ||
				IsVocationalSchoolOrMore ||
				IsJuniorCollegeOrMore ||
				IsVocationalSchoolAndJuniorCollegeOrMore ||
				IsTechnicalCollegeOrMore ||
				IsUniversityOrMore ||
				IsGraduateSchoolOrMore {
				jobInformationListWithFinalEducation = append(jobInformationListWithFinalEducation, jobInformation)
			}
		}
	} else {
		// 最終学歴未入力の場合も「求人側で設定なし」or「求人側で不問」or「求人側で中卒以上」の場合は合致
		for _, jobInformation := range jobInformationListWithNationality {
			if !jobInformation.FinalEducation.Valid ||
				jobInformation.FinalEducation == null.NewInt(99, true) ||
				jobInformation.FinalEducation == null.NewInt(0, true) {
				jobInformationListWithFinalEducation = append(jobInformationListWithFinalEducation, jobInformation)
			}
		}
	}

	/************ 5. 経験社数の絞り込み **************/

	/*

		jobSeeker:
		{ value: 0, label: '1社' }, // 転職: 0回
		{ value: 1, label: '2社' }, // 転職: 1回
		{ value: 2, label: '3社' }, // 転職: 2回
		{ value: 3, label: '4社' }, // 転職: 3回
		{ value: 4, label: '5社' }, // 転職: 4回
		{ value: 5, label: '6社' }, // 転職: 5回
		{ value: 6, label: '７社以上' }, // 転職: ６回以上

		jobInfo:
			0: '0回のみ',
			1: '1回まで',
			2: '2回まで',
			3: '3回まで',
			4: '4回まで',
			5: '5回まで',
			99: '不問',
	*/
	if searchParam.CompanyNum.Valid {
		for _, jobInformation := range jobInformationListWithFinalEducation {
			// 求人側不問 or 求人側未入力 or 求職者が6社以上で、求人の経験社数が該当している場合

			// 求人側未入力
			if !jobInformation.JobChange.Valid ||

				// 求人側不問
				jobInformation.JobChange.Int64 == 99 ||

				// 求職者が6社以下で、求職者の転職回数（経験社数）が求人の指定する転職回数以下の場合
				(searchParam.CompanyNum.Int64 != 6 && searchParam.CompanyNum.Int64 <= jobInformation.JobChange.Int64) {
				jobInformationListWithCompanyNum = append(jobInformationListWithCompanyNum, jobInformation)
				continue
			}
		}
	} else {
		for _, jobInformation := range jobInformationListWithFinalEducation {
			// 求人側未入力
			if !jobInformation.JobChange.Valid ||

				// 求人側不問
				jobInformation.JobChange.Int64 == 99 {
				jobInformationListWithCompanyNum = append(jobInformationListWithCompanyNum, jobInformation)
			}
		}
	}

	fmt.Println("経験社数: ", len(jobInformationListWithCompanyNum))

	/************ 6. 必要条件の絞り込み **************/

jobInformationLoop:
	for _, jobInformation := range jobInformationListWithCompanyNum {
		// 必要経験業界がない（不問）場合求人をリストに追加
		if len(jobInformation.RequiredConditions) == 0 {
			jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
			continue
		}

		/**
		合致した時にexperienceIndustryLoopを抜ける
		絞り込みと比較する値がどちらもスライスの場合は必要
		*/

		// 必須要件の有無がTrueの項目のisMatchがtrueの場合に求人リストにセットする
		// LP側で運転免許選択したら資格の運転免許を選択した状態にする
		var (
			isMatchCommonCondition  bool // 共通条件が合致しているか
			isMatchPatternCondition bool // パターン条件が合致しているか
		)
		for _, rc := range jobInformation.RequiredConditions {
			var (
				isRequiredIndustry   bool // 業界の条件があるか
				isRequiredOccupation bool // 職種の条件があるか
				isRequiredLanguage   bool // 言語の条件があるか
				isRequiredLicense    bool // 資格の条件があるか

				isMatchIndustry   bool // 業界が合致している
				isMatchOccupation bool // 職種が合致している
				isMatchLanguage   bool // 言語が合致している
				isMatchLicense    bool // 資格が合致している
			)

			if len(rc.RequiredExperienceJobs.ExperienceIndustries) > 0 {
				isRequiredIndustry = true
			}
			if len(rc.RequiredExperienceJobs.ExperienceOccupations) > 0 {
				isRequiredOccupation = true
			}
			if len(rc.RequiredLanguages.LanguageTypes) > 0 {
				isRequiredLanguage = true
			}
			if len(rc.RequiredLicenses) > 0 {
				isRequiredLicense = true
			}

			// 必要条件がない場合は
			if !isRequiredIndustry && !isRequiredOccupation && !isRequiredLanguage && !isRequiredLicense {
				jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
				break
			}

			/************************************************************/
			// 業界
			//
			ej := rc.RequiredExperienceJobs       // 必要経験業職種
			var requiredExperienceMonth int64 = 0 // 求人の必要経験月数

			// 必要月数を計算
			if ej.ExperienceYear.Valid && ej.ExperienceMonth.Valid {
				// 年数と月数の入力がある場合
				requiredExperienceMonth = (ej.ExperienceYear.Int64 * 6) + ej.ExperienceMonth.Int64
			} else if ej.ExperienceYear.Valid && !ej.ExperienceMonth.Valid {
				// 年数のみ入力がある場合
				requiredExperienceMonth = ej.ExperienceYear.Int64 * 6
			} else if !ej.ExperienceYear.Valid && ej.ExperienceMonth.Valid {
				// 月数のみ入力がある場合
				requiredExperienceMonth = ej.ExperienceMonth.Int64
			}

		industriesLoop:
			for _, requiredIndustry := range ej.ExperienceIndustries {
				for _, seekerIndustry := range searchParam.Industries {
					// 業界が合致していて、かつ「求人の業職種経験の必要経験年数」が未入力、または「求職者の業界経験年数」が「求人の業職種経験の必要経験年数」以上の場合
					if seekerIndustry.Industry == requiredIndustry.ExperienceIndustry && (requiredExperienceMonth == 0 || requiredExperienceMonth <= seekerIndustry.ExperienceMonth.Int64) {
						isMatchIndustry = true
						break industriesLoop
					}
				}
			}

		occupationLoop:
			for _, requiredOccupation := range ej.ExperienceOccupations {
				for _, seekerOccupation := range searchParam.AllExperienceOccupations {
					// 求職者の職種の経験月数
					var seekerExperienceMonth int64 = 0

					// seekerOccupation.ExperienceYear: {0: 1年未満, 1: 1年以上, 2: 2年以上, ..., 10: 10年以上}
					if seekerOccupation.ExperienceYear == null.NewInt(0, true) {
						seekerExperienceMonth = 6 // 1年未満は半年としてカウント
					} else if seekerOccupation.ExperienceYear.Valid && seekerOccupation.ExperienceYear.Int64 > 0 && seekerOccupation.ExperienceYear.Int64 <= 10 {
						seekerExperienceMonth = seekerOccupation.ExperienceYear.Int64 * 12
					}

					// 職種が合致していて、かつ「求人の業職種経験の必要経験年数」が未入力、または「求職者の職種経験年数」が「求人の業職種経験の必要経験年数」以上の場合
					if seekerOccupation.Occupation == requiredOccupation.ExperienceOccupation && (requiredExperienceMonth == 0 || requiredExperienceMonth <= seekerExperienceMonth) {
						isMatchOccupation = true
						break occupationLoop
					}
				}
			}

			/************************************************************/
			// 言語
			//
		languageLoop:
			for _, languageType := range rc.RequiredLanguages.LanguageTypes {
				for _, seekerLanguage := range searchParam.Languages {
					// 言語タイプが合致していて、かつ必要なレベルも合致している場合
					if seekerLanguage.LanguageType == languageType.LanguageType && seekerLanguage.LanguageLevel == rc.RequiredLanguages.LanguageLevel {
						isMatchLanguage = true
						break languageLoop
					}
				}
			}

			/************************************************************/
			// 資格
			//
		licenseLoop:
			for _, license := range rc.RequiredLicenses {
				for _, seekerLicense := range searchParam.Licenses {
					// 言語タイプが合致していて、かつ必要なレベルも合致している場合
					if seekerLicense == license.License {
						isMatchLicense = true
						break licenseLoop
					}
				}
			}

			var (
				// 1つ
				conditions1 = isRequiredIndustry && isMatchIndustry
				conditions2 = isRequiredOccupation && isMatchOccupation
				conditions3 = isRequiredLanguage && isMatchLanguage
				conditions4 = isRequiredLicense && isMatchLicense

				// 2つ
				conditions5  = isRequiredIndustry && isRequiredOccupation && isMatchIndustry && isMatchOccupation
				conditions6  = isRequiredIndustry && isRequiredLanguage && isMatchIndustry && isMatchLanguage
				conditions7  = isRequiredIndustry && isRequiredLicense && isMatchIndustry && isMatchLicense
				conditions8  = isRequiredOccupation && isRequiredLanguage && isMatchOccupation && isMatchLanguage
				conditions9  = isRequiredOccupation && isRequiredLicense && isMatchOccupation && isMatchLicense
				conditions10 = isRequiredLanguage && isRequiredLicense && isMatchLanguage && isMatchLicense

				// 3つ
				conditions11 = isRequiredIndustry && isRequiredOccupation && isRequiredLanguage && isMatchIndustry && isMatchOccupation && isMatchLanguage
				conditions12 = isRequiredIndustry && isRequiredOccupation && isRequiredLicense && isMatchIndustry && isMatchOccupation && isMatchLicense
				conditions13 = isRequiredIndustry && isRequiredLanguage && isRequiredLicense && isMatchIndustry && isMatchLanguage && isMatchLicense
				conditions14 = isRequiredOccupation && isRequiredLanguage && isRequiredLicense && isMatchOccupation && isMatchLanguage && isMatchLicense

				// 4つ
				conditions15 = isRequiredIndustry && isRequiredOccupation && isRequiredLanguage && isRequiredLicense && isMatchIndustry && isMatchOccupation && isMatchLanguage && isMatchLicense
			)

			if conditions1 || conditions2 || conditions3 || conditions4 || conditions5 ||
				conditions6 || conditions7 || conditions8 || conditions9 || conditions10 ||
				conditions11 || conditions12 || conditions13 || conditions14 || conditions15 {
				// 必要条件の場合
				if rc.IsCommon {
					fmt.Println("OK isMatchCommonCondition: ", jobInformation.ID)
					isMatchCommonCondition = true
					// 共通条件のみの場合は共通のみマッチしていればOK
					if len(jobInformation.RequiredConditions) == 1 {
						jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
						continue jobInformationLoop
					}

					// パターン条件の場合
				} else {
					fmt.Println("OK isMatchPatternCondition: ", jobInformation.ID)
					isMatchPatternCondition = true
					// 共通も含めてマッチしていればOK
					if isMatchCommonCondition {
						jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
						continue jobInformationLoop
					}
				}
				continue

				// 必要条件を満たしておらず、かつ比較対象が共通条件の場合は求人のループに戻る
			} else if rc.IsCommon {
				continue jobInformationLoop
			}
		}
		// 必要条件全てクリアしている場合は求人リストにセット
		if isMatchPatternCondition && isMatchCommonCondition {
			jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
			continue jobInformationLoop
		}
	}

	fmt.Println("必要要件: ", len(jobInformationListWithRequiredConditions))

	return jobInformationListWithRequiredConditions, nil
}

// 絞り込み検索処理 本体
func searchJobInformationByLPJobSeekerForDiagnosisParam(jobInformationList []*entity.JobInformationForDiagnosis, searchParam entity.DiagnosisParam) ([]*entity.JobInformationForDiagnosis, error) {
	fmt.Println("ヒット件数: ", len(jobInformationList))

	/**
	性別
	年齢
	国籍
	最終学歴
	経験社数（転職回数）
	必要条件（業界・職種・語学・資格）
	*/

	// 絞り込み項目の結果を代入するための変数を用意
	var (
		jobInformationListWithGender             []*entity.JobInformationForDiagnosis
		jobInformationListWithAge                []*entity.JobInformationForDiagnosis
		jobInformationListWithNationality        []*entity.JobInformationForDiagnosis
		jobInformationListWithFinalEducation     []*entity.JobInformationForDiagnosis
		jobInformationListWithCompanyNum         []*entity.JobInformationForDiagnosis
		jobInformationListWithRequiredConditions []*entity.JobInformationForDiagnosis
	)

	/************ 1. 性別の絞り込み **************/

	for _, jobInformation := range jobInformationList {
		// 求人側で未入力の場合
		if !jobInformation.Gender.Valid ||

			// 求人側で特定の性別のみでない場合（2: 男性尚可 or 3: 女性尚可　or 99: 不問）
			(jobInformation.Gender != null.NewInt(0, true) && jobInformation.Gender != null.NewInt(1, true)) ||

			// 求人側と求職者の性別が一致する場合
			(searchParam.Gender == jobInformation.Gender) {
			jobInformationListWithGender = append(jobInformationListWithGender, jobInformation)
		}
	}

	fmt.Println("性別: ", len(jobInformationListWithGender))

	/************ 2. 年齢の絞り込み **************/

	if searchParam.Birthyear != "" && searchParam.Birthmonth != "" && searchParam.Birthday != "" {
		// 現在の日付を取得
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().In(jst)

		// フォーマットに従って time.Time 型に変換
		dateStr := fmt.Sprintf("%s-%s-%s", searchParam.Birthyear, searchParam.Birthmonth, searchParam.Birthday)
		birthDay, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return nil, err
		}

		// 年齢を計算
		age := now.Year() - birthDay.Year()

		// まだ誕生日を迎えていない場合は、年齢から1を引く
		if now.Month() < birthDay.Month() || (now.Month() == birthDay.Month() && now.Day() < birthDay.Day()) {
			age--
		}

		for _, jobInformation := range jobInformationListWithGender {
			// 募集年齢: 不問（年齢上限、年齢下限が未入力）
			if (!jobInformation.AgeUnder.Valid && !jobInformation.AgeOver.Valid) ||

				// 年齢下限と上限が入力されており、検索パラムの値が下限以上で上限以下の場合
				(jobInformation.AgeUnder.Valid && age >= int(jobInformation.AgeUnder.Int64) && jobInformation.AgeOver.Valid && age <= int(jobInformation.AgeOver.Int64)) ||

				// 年齢下限が入力されている + 検索パラムの値が上限以下 + 年齢上限が未入力
				(jobInformation.AgeUnder.Valid && age >= int(jobInformation.AgeUnder.Int64) && !jobInformation.AgeOver.Valid) ||

				// 年齢上限が入力されている + 検索パラムの値が上限以下 + 年齢下限が未入力
				(jobInformation.AgeOver.Valid && age <= int(jobInformation.AgeOver.Int64) && !jobInformation.AgeUnder.Valid) {

				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
				continue
			}
		}
	} else {
		// 生年月日の入力がない場合は年齢の条件が不問の求人のみ合致
		for _, jobInformation := range jobInformationListWithGender {
			// 募集年齢: 不問（年齢上限、年齢下限が未入力）
			if !jobInformation.AgeUnder.Valid && !jobInformation.AgeOver.Valid {
				jobInformationListWithAge = append(jobInformationListWithAge, jobInformation)
			}
		}
	}

	fmt.Println("年齢: ", len(jobInformationListWithAge))

	/************ 3. 国籍の絞り込み **************/

	for _, jobInformation := range jobInformationListWithAge {
		var isPossibleForeign bool

		// 求人の特徴に外国籍化の求人があるかをチェック
		for _, feature := range jobInformation.Features {
			if feature.Feature == null.NewInt(21, true) {
				isPossibleForeign = true
				break
			}
		}

		// 求人側で未入力の場合
		if !jobInformation.Nationality.Valid ||

			// 求人側で不問の場合
			jobInformation.Nationality == null.NewInt(99, true) ||

			// 求人側の募集国籍が日本国籍のみで求職者が日本国籍
			(jobInformation.Nationality == null.NewInt(0, true) && searchParam.FirstLanguage == null.NewInt(0, true)) ||

			// 求人側の募集国籍が外国籍のみで求職者が外国籍
			(jobInformation.Nationality == null.NewInt(1, true) && (searchParam.FirstLanguage == null.NewInt(1, true) || searchParam.FirstLanguage == null.NewInt(2, true))) ||

			// 求人側が外国籍可で求職者が外国籍
			(isPossibleForeign && (searchParam.FirstLanguage == null.NewInt(1, true) || searchParam.FirstLanguage == null.NewInt(2, true))) {
			jobInformationListWithNationality = append(jobInformationListWithNationality, jobInformation)
		}
	}

	/************ 4. 最終学歴の絞り込み **************/

	if searchParam.SchoolCategory.Valid {
		for _, jobInformation := range jobInformationListWithNationality {
			// 求人側で未入力の場合: !jobInformation.FinalEducation.Valid
			// 求人側で不問の場合: jobInformation.FinalEducation == null.NewInt(99, true)
			// 求人側で中卒以上の場合: jobInformation.FinalEducation == null.NewInt(0, true)
			// 求人側で高卒以上で求職者が高卒以上の場合: jobInformation.FinalEducation == null.NewInt(1, true) && 1 <= searchParam.SchoolCategory.Int64
			// 求人側で専卒以上（短大除く）で求職者が高卒以上の場合: jobInformation.FinalEducation == null.NewInt(1, true) && 1 <= searchParam.SchoolCategory.Int64

			var (
				// 高卒以上
				IsHighSchoolOrMore = jobInformation.FinalEducation == null.NewInt(1, true) && 1 <= searchParam.SchoolCategory.Int64

				// 専卒以上
				IsVocationalSchoolOrMore = jobInformation.FinalEducation == null.NewInt(2, true) && (searchParam.SchoolCategory == null.NewInt(2, true) || searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 短大卒以上
				IsJuniorCollegeOrMore = jobInformation.FinalEducation == null.NewInt(3, true) && (searchParam.SchoolCategory == null.NewInt(3, true) || searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 専卒と短大卒以上
				IsVocationalSchoolAndJuniorCollegeOrMore = jobInformation.FinalEducation == null.NewInt(4, true) && (searchParam.SchoolCategory == null.NewInt(2, true) || searchParam.SchoolCategory == null.NewInt(3, true) || searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 高専卒以上
				IsTechnicalCollegeOrMore = jobInformation.FinalEducation == null.NewInt(5, true) && (searchParam.SchoolCategory == null.NewInt(4, true) || searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 大卒以上
				IsUniversityOrMore = jobInformation.FinalEducation == null.NewInt(6, true) && (searchParam.SchoolCategory == null.NewInt(5, true) || searchParam.SchoolCategory == null.NewInt(6, true))

				// 院卒以上
				IsGraduateSchoolOrMore = jobInformation.FinalEducation == null.NewInt(7, true) && (searchParam.SchoolCategory == null.NewInt(6, true))
			)

			if !jobInformation.FinalEducation.Valid ||
				jobInformation.FinalEducation == null.NewInt(99, true) ||
				jobInformation.FinalEducation == null.NewInt(0, true) ||
				IsHighSchoolOrMore ||
				IsVocationalSchoolOrMore ||
				IsJuniorCollegeOrMore ||
				IsVocationalSchoolAndJuniorCollegeOrMore ||
				IsTechnicalCollegeOrMore ||
				IsUniversityOrMore ||
				IsGraduateSchoolOrMore {
				jobInformationListWithFinalEducation = append(jobInformationListWithFinalEducation, jobInformation)
			}
		}
	} else {
		// 最終学歴未入力の場合も「求人側で設定なし」or「求人側で不問」or「求人側で中卒以上」の場合は合致
		for _, jobInformation := range jobInformationListWithNationality {
			if !jobInformation.FinalEducation.Valid ||
				jobInformation.FinalEducation == null.NewInt(99, true) ||
				jobInformation.FinalEducation == null.NewInt(0, true) {
				jobInformationListWithFinalEducation = append(jobInformationListWithFinalEducation, jobInformation)
			}
		}
	}

	/************ 5. 経験社数の絞り込み **************/

	/*

		jobSeeker:
		{ value: 0, label: '1社' }, // 転職: 0回
		{ value: 1, label: '2社' }, // 転職: 1回
		{ value: 2, label: '3社' }, // 転職: 2回
		{ value: 3, label: '4社' }, // 転職: 3回
		{ value: 4, label: '5社' }, // 転職: 4回
		{ value: 5, label: '6社' }, // 転職: 5回
		{ value: 6, label: '７社以上' }, // 転職: ６回以上

		jobInfo:
			0: '0回のみ',
			1: '1回まで',
			2: '2回まで',
			3: '3回まで',
			4: '4回まで',
			5: '5回まで',
			99: '不問',
	*/
	if searchParam.CompanyNum.Valid {
		for _, jobInformation := range jobInformationListWithFinalEducation {
			// 求人側不問 or 求人側未入力 or 求職者が6社以上で、求人の経験社数が該当している場合

			// 求人側未入力
			if !jobInformation.JobChange.Valid ||

				// 求人側不問
				jobInformation.JobChange.Int64 == 99 ||

				// 求職者が6社以下で、求職者の転職回数（経験社数）が求人の指定する転職回数以下の場合
				(searchParam.CompanyNum.Int64 != 6 && searchParam.CompanyNum.Int64 <= jobInformation.JobChange.Int64) {
				jobInformationListWithCompanyNum = append(jobInformationListWithCompanyNum, jobInformation)
				continue
			}
		}
	} else {
		for _, jobInformation := range jobInformationListWithFinalEducation {
			// 求人側未入力
			if !jobInformation.JobChange.Valid ||

				// 求人側不問
				jobInformation.JobChange.Int64 == 99 {
				jobInformationListWithCompanyNum = append(jobInformationListWithCompanyNum, jobInformation)
			}
		}
	}

	fmt.Println("経験社数: ", len(jobInformationListWithCompanyNum))

	/************ 6. 必要条件の絞り込み **************/

jobInformationLoop:
	for _, jobInformation := range jobInformationListWithCompanyNum {
		// 必要経験業界がない（不問）場合求人をリストに追加
		if len(jobInformation.RequiredConditions) == 0 {
			jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
			continue
		}

		/**
		合致した時にexperienceIndustryLoopを抜ける
		絞り込みと比較する値がどちらもスライスの場合は必要
		*/

		// 必須要件の有無がTrueの項目のisMatchがtrueの場合に求人リストにセットする
		// LP側で運転免許選択したら資格の運転免許を選択した状態にする
		var (
			isMatchCommonCondition  bool // 共通条件が合致しているか
			isMatchPatternCondition bool // パターン条件が合致しているか
		)
		for _, rc := range jobInformation.RequiredConditions {
			var (
				isRequiredIndustry   bool // 業界の条件があるか
				isRequiredOccupation bool // 職種の条件があるか
				isRequiredLanguage   bool // 言語の条件があるか
				isRequiredLicense    bool // 資格の条件があるか

				isMatchIndustry   bool // 業界が合致している
				isMatchOccupation bool // 職種が合致している
				isMatchLanguage   bool // 言語が合致している
				isMatchLicense    bool // 資格が合致している
			)

			if len(rc.RequiredExperienceJobs.ExperienceIndustries) > 0 {
				isRequiredIndustry = true
			}
			if len(rc.RequiredExperienceJobs.ExperienceOccupations) > 0 {
				isRequiredOccupation = true
			}
			if len(rc.RequiredLanguages.LanguageTypes) > 0 {
				isRequiredLanguage = true
			}
			if len(rc.RequiredLicenses) > 0 {
				isRequiredLicense = true
			}

			// 必要条件がない場合は
			if !isRequiredIndustry && !isRequiredOccupation && !isRequiredLanguage && !isRequiredLicense {
				jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
				break
			}

			/************************************************************/
			// 業界
			//
			ej := rc.RequiredExperienceJobs       // 必要経験業職種
			var requiredExperienceMonth int64 = 0 // 求人の必要経験月数

			// 必要月数を計算
			if ej.ExperienceYear.Valid && ej.ExperienceMonth.Valid {
				// 年数と月数の入力がある場合
				requiredExperienceMonth = (ej.ExperienceYear.Int64 * 6) + ej.ExperienceMonth.Int64
			} else if ej.ExperienceYear.Valid && !ej.ExperienceMonth.Valid {
				// 年数のみ入力がある場合
				requiredExperienceMonth = ej.ExperienceYear.Int64 * 6
			} else if !ej.ExperienceYear.Valid && ej.ExperienceMonth.Valid {
				// 月数のみ入力がある場合
				requiredExperienceMonth = ej.ExperienceMonth.Int64
			}

		industriesLoop:
			for _, requiredIndustry := range ej.ExperienceIndustries {
				for _, seekerIndustry := range searchParam.Industries {
					// 業界が合致していて、かつ「求人の業職種経験の必要経験年数」が未入力、または「求職者の業界経験年数」が「求人の業職種経験の必要経験年数」以上の場合
					if seekerIndustry.Industry == requiredIndustry.ExperienceIndustry && (requiredExperienceMonth == 0 || requiredExperienceMonth <= seekerIndustry.ExperienceMonth.Int64) {
						isMatchIndustry = true
						break industriesLoop
					}
				}
			}

		occupationLoop:
			for _, requiredOccupation := range ej.ExperienceOccupations {
				for _, seekerOccupation := range searchParam.AllExperienceOccupations {
					// 求職者の職種の経験月数
					var seekerExperienceMonth int64 = 0

					// seekerOccupation.ExperienceYear: {0: 1年未満, 1: 1年以上, 2: 2年以上, ..., 10: 10年以上}
					if seekerOccupation.ExperienceYear == null.NewInt(0, true) {
						seekerExperienceMonth = 6 // 1年未満は半年としてカウント
					} else if seekerOccupation.ExperienceYear.Valid && seekerOccupation.ExperienceYear.Int64 > 0 && seekerOccupation.ExperienceYear.Int64 <= 10 {
						seekerExperienceMonth = seekerOccupation.ExperienceYear.Int64 * 12
					}

					// 職種が合致していて、かつ「求人の業職種経験の必要経験年数」が未入力、または「求職者の職種経験年数」が「求人の業職種経験の必要経験年数」以上の場合
					if seekerOccupation.Occupation == requiredOccupation.ExperienceOccupation && (requiredExperienceMonth == 0 || requiredExperienceMonth <= seekerExperienceMonth) {
						isMatchOccupation = true
						break occupationLoop
					}
				}
			}

			/************************************************************/
			// 言語
			//
		languageLoop:
			for _, languageType := range rc.RequiredLanguages.LanguageTypes {
				for _, seekerLanguage := range searchParam.Languages {
					// 言語タイプが合致していて、かつ必要なレベルも合致している場合
					if seekerLanguage.LanguageType == languageType.LanguageType && seekerLanguage.LanguageLevel == rc.RequiredLanguages.LanguageLevel {
						isMatchLanguage = true
						break languageLoop
					}
				}
			}

			/************************************************************/
			// 資格
			//
		licenseLoop:
			for _, license := range rc.RequiredLicenses {
				for _, seekerLicense := range searchParam.Licenses {
					// 言語タイプが合致していて、かつ必要なレベルも合致している場合
					if seekerLicense == license.License {
						isMatchLicense = true
						break licenseLoop
					}
				}
			}

			var (
				// 1つ
				conditions1 = isRequiredIndustry && isMatchIndustry
				conditions2 = isRequiredOccupation && isMatchOccupation
				conditions3 = isRequiredLanguage && isMatchLanguage
				conditions4 = isRequiredLicense && isMatchLicense

				// 2つ
				conditions5  = isRequiredIndustry && isRequiredOccupation && isMatchIndustry && isMatchOccupation
				conditions6  = isRequiredIndustry && isRequiredLanguage && isMatchIndustry && isMatchLanguage
				conditions7  = isRequiredIndustry && isRequiredLicense && isMatchIndustry && isMatchLicense
				conditions8  = isRequiredOccupation && isRequiredLanguage && isMatchOccupation && isMatchLanguage
				conditions9  = isRequiredOccupation && isRequiredLicense && isMatchOccupation && isMatchLicense
				conditions10 = isRequiredLanguage && isRequiredLicense && isMatchLanguage && isMatchLicense

				// 3つ
				conditions11 = isRequiredIndustry && isRequiredOccupation && isRequiredLanguage && isMatchIndustry && isMatchOccupation && isMatchLanguage
				conditions12 = isRequiredIndustry && isRequiredOccupation && isRequiredLicense && isMatchIndustry && isMatchOccupation && isMatchLicense
				conditions13 = isRequiredIndustry && isRequiredLanguage && isRequiredLicense && isMatchIndustry && isMatchLanguage && isMatchLicense
				conditions14 = isRequiredOccupation && isRequiredLanguage && isRequiredLicense && isMatchOccupation && isMatchLanguage && isMatchLicense

				// 4つ
				conditions15 = isRequiredIndustry && isRequiredOccupation && isRequiredLanguage && isRequiredLicense && isMatchIndustry && isMatchOccupation && isMatchLanguage && isMatchLicense
			)

			if conditions1 || conditions2 || conditions3 || conditions4 || conditions5 ||
				conditions6 || conditions7 || conditions8 || conditions9 || conditions10 ||
				conditions11 || conditions12 || conditions13 || conditions14 || conditions15 {
				// 必要条件の場合
				if rc.IsCommon {
					fmt.Println("OK isMatchCommonCondition: ", jobInformation.ID)
					isMatchCommonCondition = true
					// 共通条件のみの場合は共通のみマッチしていればOK
					if len(jobInformation.RequiredConditions) == 1 {
						jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
						continue jobInformationLoop
					}

					// パターン条件の場合
				} else {
					fmt.Println("OK isMatchPatternCondition: ", jobInformation.ID)
					isMatchPatternCondition = true
					// 共通も含めてマッチしていればOK
					if isMatchCommonCondition {
						jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
						continue jobInformationLoop
					}
				}
				continue

				// 必要条件を満たしておらず、かつ比較対象が共通条件の場合は求人のループに戻る
			} else if rc.IsCommon {
				continue jobInformationLoop
			}
		}
		// 必要条件全てクリアしている場合は求人リストにセット
		if isMatchPatternCondition && isMatchCommonCondition {
			jobInformationListWithRequiredConditions = append(jobInformationListWithRequiredConditions, jobInformation)
			continue jobInformationLoop
		}
	}

	fmt.Println("必要要件: ", len(jobInformationListWithRequiredConditions))

	return jobInformationListWithRequiredConditions, nil
}

type GetJobListingListAndJobSeekerDesiredForDiagnosisInput struct {
	// 応募資格に合致した求人リスト
	JobSeekerUUID uuid.UUID
}

type GetJobListingListAndJobSeekerDesiredForDiagnosisOutput struct {
	// 応募資格に合致した求人リスト
	JobListingList []*entity.JobListing

	JobSeekerDesired *entity.JobSeekerDesiredForGuest
}

func (i *JobInformationInteractorImpl) GetJobListingListAndJobSeekerDesiredForDiagnosis(input GetJobListingListAndJobSeekerDesiredForDiagnosisInput) (GetJobListingListAndJobSeekerDesiredForDiagnosisOutput, error) {
	var (
		output                         GetJobListingListAndJobSeekerDesiredForDiagnosisOutput
		err                            error
		jobInformationList             []*entity.JobInformation
		jobInformationListAfterMapping []*entity.JobInformation
		motoyuiAgentID                 uint = 1

		// 絞り込み
		searchParamLicenses              []null.Int
		searchParamLanguages             []entity.Language
		searchParamExperienceIndustries  []entity.ExperienceIndustry
		searchParamExperienceOccupations []entity.ExperienceOccupation
		birthYear                        string
		birthMonth                       string
		birthDay                         string
		schoolCategory                   null.Int

		// 希望条件
		desiredIndustryListNullInt     []null.Int
		desiredOccupationListNullInt   []null.Int
		desiredWorkLocationListNullInt []null.Int
	)

	/************ 1. 絞り込みに使用する求職者情報を取得 **************/
	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. 資格情報のセット **************/

	jobSeekerLicenses, err := i.jobSeekerLicenseRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var (
		driverLicense null.Int
		licenseMT     = null.NewInt(4803, true) // 資格マスタ
		licenseAT     = null.NewInt(4805, true) // 資格マスタ
		hasLicenseMT  = null.NewInt(1, true)    // LPマスタ（持っている（MT車））
		hasLicenseAT  = null.NewInt(0, true)    // LPマスタ（持っている（AT車限定））
	)

	for _, jobSeekerLicense := range jobSeekerLicenses {
		if jobSeekerLicense.LicenseType.Valid {
			searchParamLicenses = append(searchParamLicenses, jobSeekerLicense.LicenseType)

			// 普通自動車運転免許の場合は、運転免許フラグをtrueにする
			if jobSeekerLicense.LicenseType == licenseAT {
				driverLicense = hasLicenseAT
			}
			if jobSeekerLicense.LicenseType == licenseMT {
				driverLicense = hasLicenseMT
			}
		}
	}

	/************ 3. 語学情報のセット **************/

	jobSeekerLanguages, err := i.jobSeekerLanguageSkillRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	for _, jobSeekerLanguage := range jobSeekerLanguages {
		searchParamLanguages = append(searchParamLanguages, entity.Language{
			LanguageType:  jobSeekerLanguage.LanguageType,
			LanguageLevel: jobSeekerLanguage.LanguageLevel,
		})
	}

	/************ 4. 経験業界のセット **************/

	workHistories, err := i.jobSeekerWorkHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustries, err := i.jobSeekerExperienceIndustryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 経験年数を合計するために業界ごとに何ヶ月経験があるかを計算
	industryMonth := make(map[null.Int]int64)

	for _, workHistory := range workHistories {
		var experienceMonth int64 = 0

		if workHistory.JoiningYear != "" {
			// フォーマットに従って time.Time 型に変換
			// workHistoryのJoiningYearとRetireYearは「YYYY-MM」形式
			joinDateStr := fmt.Sprintf("%s-%s", workHistory.JoiningYear, "01")
			joinDate, joingError := time.Parse("2006-01-02", joinDateStr)

			var retireDateStr string
			if workHistory.RetireYear == "" {
				now := time.Now()
				monthStr := fmt.Sprintf("%02d", int(now.Month()))
				retireDateStr = fmt.Sprintf("%v-%v-%s", now.Year(), monthStr, "01")
			} else {
				retireDateStr = fmt.Sprintf("%s-%s", workHistory.RetireYear, "01")
			}
			retireDate, retireError := time.Parse("2006-01-02", retireDateStr)

			diffYear := int64(retireDate.Year() - joinDate.Year())
			diffMonth := int64(retireDate.Month() - joinDate.Month())

			if joingError == nil && retireError == nil {
				// 求職者の経験月数
				experienceMonth = (diffYear * 12) + diffMonth
			}
		}

		// 経験業界
		for _, industry := range experienceIndustries {
			if workHistory.ID == industry.WorkHistoryID {
				industryMonth[industry.Industry] += experienceMonth
			}
		}
	}

	// 計算した経験月数をLPのマスタに変換
	for industry, experienceMonth := range industryMonth {
		// 職種の重複がない場合は新しく職種をセット
		searchParamExperienceIndustries = append(
			searchParamExperienceIndustries,
			entity.ExperienceIndustry{
				Industry:        industry,
				ExperienceMonth: null.NewInt(experienceMonth, true),
			},
		)
	}

	/************ 5. 経験職種のセット **************/

	// 職種絞り込みは専用テーブルの値を使用する。 職歴に登録されている場合(面談済み)、そちらを使用して絞り込む。
	departmentHistories, err := i.jobSeekerDepartmentHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupations, err := i.jobSeekerExperienceOccupationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if len(experienceOccupations) != 0 {
		// 経験職種が入力されている（面談済み）
		// 経験年数を合計するために職種ごとに何ヶ月経験があるかを計算
		occupationMonth := make(map[null.Int]int64)

		for _, departmentHistory := range departmentHistories {
			for _, experienceOccupation := range experienceOccupations {
				if departmentHistory.ID == experienceOccupation.DepartmentHistoryID {

					// 在職中（現在に至る）
					endYear := departmentHistory.EndYear

					// departmentHistory.EndYearがからの場合は現職中の判断する
					// TODO: WorkHistoryのlastStatusが「現在に至る」の場合に現在時刻を入れるように修正
					if endYear == "" {
						now := time.Now()
						endYear = now.Format("2006-01")
					}

					experienceYear := getNumberOfMonths(departmentHistory.StartYear, endYear)
					occupationMonth[experienceOccupation.Occupation] += experienceYear
				}
			}
		}

		// 計算した経験月数をLPのマスタに変換
		for occupation, experienceMonth := range occupationMonth {
			// LPの経験年数のマスタに変換
			var (
				experienceYearForLP null.Int
				under1Year          = experienceMonth < 12   // 0: 1年未満
				over1Year           = experienceMonth < 24   // 1: 1年以上
				over2Year           = experienceMonth < 36   // 2: 2年以上
				over3Year           = experienceMonth < 48   // 3: 3年以上
				over4Year           = experienceMonth < 60   // 4: 4年以上
				over5Year           = experienceMonth < 72   // 5: 5年以上
				over6Year           = experienceMonth < 84   // 6: 6年以上
				over7Year           = experienceMonth < 96   // 7: 7年以上
				over8Year           = experienceMonth < 108  // 8: 8年以上
				over9Year           = experienceMonth < 120  // 9: 9年以上
				over10Year          = 120 <= experienceMonth // 10: 10年以上
			)

			if under1Year {
				experienceYearForLP = null.NewInt(0, true)
			} else if over1Year {
				experienceYearForLP = null.NewInt(1, true)
			} else if over2Year {
				experienceYearForLP = null.NewInt(2, true)
			} else if over3Year {
				experienceYearForLP = null.NewInt(3, true)
			} else if over4Year {
				experienceYearForLP = null.NewInt(4, true)
			} else if over5Year {
				experienceYearForLP = null.NewInt(5, true)
			} else if over6Year {
				experienceYearForLP = null.NewInt(6, true)
			} else if over7Year {
				experienceYearForLP = null.NewInt(7, true)
			} else if over8Year {
				experienceYearForLP = null.NewInt(8, true)
			} else if over9Year {
				experienceYearForLP = null.NewInt(9, true)
			} else if over10Year {
				experienceYearForLP = null.NewInt(10, true)
			}

			// 職種の重複がない場合は新しく職種をセット
			searchParamExperienceOccupations = append(
				searchParamExperienceOccupations,
				entity.ExperienceOccupation{
					Occupation:     occupation,
					ExperienceYear: experienceYearForLP,
				},
			)
		}

	} else {
		// 診断で入力された職歴のみ（面談実施前）
		jobSeekerExperienceJobs, err := i.jobSeekerExperienceJobRepository.GetByJobSeekerID(jobSeeker.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, occupation := range jobSeekerExperienceJobs {
			searchParamExperienceOccupations = append(
				searchParamExperienceOccupations,
				entity.ExperienceOccupation{
					Occupation:     occupation.Occupation,
					ExperienceYear: occupation.ExperienceYear,
				},
			)
		}
	}

	/************ 6. 生年月日のセット **************/

	if jobSeeker.Birthday != "" {
		// "-"で分割
		parts := strings.Split(jobSeeker.Birthday, "-")
		if len(parts) == 3 {
			// 各要素を文字列として変数に格納
			birthYear = parts[0]
			birthMonth = parts[1]
			birthDay = parts[2]
		}
	}

	/************ 7. 学歴のセット **************/

	studentHistories, err := i.jobSeekerStudentHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var schoolCategoryList []null.Int // 学校の種類のリスト

	if len(studentHistories) > 0 {
		// 指定された条件に合致する学校の種類をリストに格納
		for _, history := range jobSeeker.StudentHistories {
			// {0: 卒業, 1: 中退, 2: 退学, 3: 卒業見込み, 4: 修了, 5: 修了見込み}
			// 1と2以外は最終学歴としてカウント
			if !history.SchoolCategory.Valid && history.LastStatus != null.NewInt(1, true) && history.LastStatus != null.NewInt(2, true) {
				schoolCategoryList = append(schoolCategoryList, history.SchoolCategory)
			}
		}

		// 学校の種類のリストが空でない場合
		if len(schoolCategoryList) > 0 {
			// 最大の学校の種類を求める
			schoolCategory = schoolCategoryList[0]
			for _, category := range schoolCategoryList[1:] {
				if category.Int64 > schoolCategory.Int64 {
					schoolCategory = category
				}
			}
		}
	}

	/************ 8. 条件検索用のパラム変換 **************/
	searchParam := entity.DiagnosisParam{
		Gender:                jobSeeker.Gender,
		Birthyear:             birthYear,
		Birthmonth:            birthMonth,
		Birthday:              birthDay,
		FirstLanguage:         jobSeeker.Nationality, // 国籍はそのままFirstLanguageにそのまま入れる
		SchoolCategory:        schoolCategory,
		CompanyNum:            jobSeeker.JobChange, // 転職回数（経験社数）はそのまま入れる
		Industries:            searchParamExperienceIndustries,
		ExperienceOccupations: searchParamExperienceOccupations,
		Languages:             searchParamLanguages,
		DriversLicense:        driverLicense,
		Licenses:              searchParamLicenses,
	}

	/************ 9. 希望条件の設定 **************/
	desiredIndustryList, err := i.jobSeekerDesiredIndustryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupationList, err := i.jobSeekerDesiredOccupationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocationList, err := i.jobSeekerDesiredWorkLocationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// マッチ求人検索項目をnullint[]に変換
	for _, di := range desiredIndustryList {
		desiredIndustryListNullInt = append(desiredIndustryListNullInt, di.DesiredIndustry)
	}
	for _, do := range desiredOccupationList {
		desiredOccupationListNullInt = append(desiredOccupationListNullInt, do.DesiredOccupation)
	}
	for _, dwl := range desiredWorkLocationList {
		desiredWorkLocationListNullInt = append(desiredWorkLocationListNullInt, dwl.DesiredWorkLocation)
	}

	output.JobSeekerDesired = &entity.JobSeekerDesiredForGuest{
		ID:                   jobSeeker.ID,
		Phase:                jobSeeker.Phase,
		DesiredAnnualIncome:  jobSeeker.DesiredAnnualIncome,
		DesiredIndustries:    desiredIndustryListNullInt,
		DesiredOccupations:   desiredOccupationListNullInt,
		DesiredWorkLocations: desiredWorkLocationListNullInt,
	}

	/************ 10. 求人リストの取得 **************/
	// 元結求人&外部求人除く
	jobInformationList, err = i.jobInformationRepository.GetActiveAllByAgentIDWithoutExternal(motoyuiAgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobInformationIDList := make([]uint, 0, len(jobInformationList))
	for _, jobInformation := range jobInformationList {
		jobInformationIDList = append(jobInformationIDList, jobInformation.ID)
	}

	// 求人リストから、IDが合致する子テーブルの情報を取得
	// 検索で使用しないかつリスト表示しない情報は、取得しない
	requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguage, err := i.jobInfoRequiredLanguageTypeRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者ページ表示用項目のみ取得 0:業界未経験OK, 1:職種未経験OK, 2:業界・職種未経験OK, 6:転勤なし
	features, err := i.jobInfoFeatureRepository.GetByJobInformationIDListForGuestJobSeeker(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.enterpriseIndustryRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	taskGroups, err := i.taskGroupRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

jobInfoLoop:
	for _, jobInformation := range jobInformationList {
		// すでにタスクがある求人を除外
		for _, taskGroup := range taskGroups {
			if jobInformation.ID == taskGroup.JobInformationID {
				continue jobInfoLoop
			}
		}

		for _, condition := range requiredConditions {
			if jobInformation.ID == condition.JobInformationID {
				for _, rl := range requiredLicenses {
					if condition.ID == rl.ConditionID {
						condition.RequiredLicenses = append(condition.RequiredLicenses, *rl)
					}
				}

				for _, rl := range requiredLanguages {
					if condition.ID == rl.ConditionID {

						// 言語タイプ
						for _, languageType := range requiredLanguage {
							if rl.ID == languageType.LanguageID {
								rl.LanguageTypes = append(rl.LanguageTypes, *languageType)
							}
						}

						condition.RequiredLanguages = *rl
					}
				}

				for _, rej := range requiredExperienceJobs {
					if condition.ID == rej.ConditionID {

						// 業界
						for _, industry := range requiredExperienceIndustries {
							if rej.ID == industry.ExperienceJobID {
								rej.ExperienceIndustries = append(rej.ExperienceIndustries, *industry)
							}
						}

						// 職種
						for _, occupation := range requiredExperienceOccupations {
							if rej.ID == occupation.ExperienceJobID {
								rej.ExperienceOccupations = append(rej.ExperienceOccupations, *occupation)
							}
						}

						condition.RequiredExperienceJobs = *rej
					}
				}

				jobInformation.RequiredConditions = append(jobInformation.RequiredConditions, *condition)
			}
		}

		for _, f := range features {
			if jobInformation.ID == f.JobInformationID {
				value := entity.JobInformationFeature{
					JobInformationID: f.JobInformationID,
					Feature:          f.Feature,
				}
				jobInformation.Features = append(jobInformation.Features, value)
			}
		}

		for _, p := range prefectures {
			if jobInformation.ID == p.JobInformationID {
				value := entity.JobInformationPrefecture{
					JobInformationID: p.JobInformationID,
					Prefecture:       p.Prefecture,
				}
				jobInformation.Prefectures = append(jobInformation.Prefectures, value)
			}
		}

		for _, oc := range occupations {
			if jobInformation.ID == oc.JobInformationID {
				value := entity.JobInformationOccupation{
					JobInformationID: oc.JobInformationID,
					Occupation:       oc.Occupation,
				}
				jobInformation.Occupations = append(jobInformation.Occupations, value)
			}
		}

		for _, ind := range industries {
			if jobInformation.EnterpriseID == ind.EnterpriseID {
				value := entity.EnterpriseIndustry{
					EnterpriseID: ind.EnterpriseID,
					Industry:     ind.Industry,
				}
				jobInformation.Industries = append(jobInformation.Industries, value)
			}
		}

		jobInformationListAfterMapping = append(jobInformationListAfterMapping, jobInformation)
	}

	/************ 11. 必要条件の絞り込み **************/
	jobInformationListAfterSearch, err := searchJobInformationByLPJobSeeker(jobInformationListAfterMapping, searchParam)
	if err != nil {
		return output, err
	}

	for _, jobInformation := range jobInformationListAfterSearch {
		jobListing := entity.NewJobListing(
			jobInformation.ID,
			jobInformation.AgentStaffID,
			jobInformation.CompanyName,
			jobInformation.CorporateSiteURL,
			jobInformation.PostCode,
			jobInformation.OfficeLocation,
			jobInformation.EmployeeNumberSingle,
			jobInformation.EmployeeNumberGroup,
			jobInformation.Establishment,
			jobInformation.PublicOffering,
			jobInformation.Earnings,
			jobInformation.EarningsYear,
			jobInformation.BusinessDetail,
			jobInformation.Title,
			jobInformation.WorkDetail,
			jobInformation.WorkLocation,
			jobInformation.Transfer,
			jobInformation.TransferDetail,
			jobInformation.UnderIncome,
			jobInformation.OverIncome,
			jobInformation.Salary,
			jobInformation.Insurance,
			jobInformation.WorkTime,
			jobInformation.OvertimeAverage,
			jobInformation.FixedOvertimePayment,
			jobInformation.FixedOvertimeDetail,
			jobInformation.TrialPeriod,
			jobInformation.TrialPeriodDetail,
			jobInformation.EmploymentPeriod,
			jobInformation.EmploymentPeriodDetail,
			jobInformation.HolidayDetail,
			jobInformation.PassiveSmoking,
			jobInformation.SelectionFlow,
			jobInformation.EmploymentInsurance,
			jobInformation.AccidentInsurance,
			jobInformation.HealthInsurance,
			jobInformation.PensionInsurance,
			jobInformation.IsExternal,
			jobInformation.WorkDetailAfterHiring,
			jobInformation.WorkDetailScopeOfChange,
		)

		jobListing.IsGuaranteedInterview = jobInformation.IsGuaranteedInterview
		jobListing.JobInformationUUID = jobInformation.UUID
		jobListing.Prefectures = jobInformation.Prefectures
		jobListing.Occupations = jobInformation.Occupations
		jobListing.Industries = jobInformation.Industries
		jobListing.Features = jobInformation.Features

		output.JobListingList = append(output.JobListingList, jobListing)
	}

	return output, nil
}

type GetJobListingListByJobSeekerUUIDAndInterestedTypeInput struct {
	Param entity.InterestedTypeJobListParam
}

type GetJobListingListByJobSeekerUUIDAndInterestedTypeOutput struct {
	JobListingList []*entity.JobListing
	MaxPageNumber  uint
}

func (i *JobInformationInteractorImpl) GetJobListingListByJobSeekerUUIDAndInterestedType(input GetJobListingListByJobSeekerUUIDAndInterestedTypeInput) (GetJobListingListByJobSeekerUUIDAndInterestedTypeOutput, error) {
	var (
		output             GetJobListingListByJobSeekerUUIDAndInterestedTypeOutput
		inputParam         = input.Param
		err                error
		jobInformationList []*entity.JobInformation
		jobListingList     []*entity.JobListing
	)

	/************ エラーハンドリング **************/

	if !inputParam.InterestedType.Valid {
		wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "再読み込みして再度お試しください。")
		return output, wrapped
	}

	if !inputParam.PageNumber.Valid {
		wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "指定のページが存在しません。\n再読み込みして再度お試しください。")
		return output, wrapped
	}

	/************ InteresterJobListingを取得 **************/

	interestedJobListingList, err := i.jobSeekerInterestedJobListingRepository.GetByJobSeekerUUIDAndInterestedType(inputParam.JobSeekerUUID, entity.InterestedType(inputParam.InterestedType.Int64))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobInformationIDList := make([]uint, 0, len(interestedJobListingList))
	for _, interestedJobListing := range interestedJobListingList {
		jobInformationIDList = append(jobInformationIDList, interestedJobListing.JobInformationID)
	}

	/************ 求人リストの取得 **************/

	// 元結求人&外部求人除く
	jobInformationList, err = i.jobInformationRepository.GetByIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者ページ表示用項目のみ取得 0:業界未経験OK, 1:職種未経験OK, 2:業界・職種未経験OK, 6:転勤なし
	features, err := i.jobInfoFeatureRepository.GetByJobInformationIDListForGuestJobSeeker(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.enterpriseIndustryRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// jobInfoLoop:
	for _, jobInformation := range jobInformationList {

		for _, f := range features {
			if jobInformation.ID == f.JobInformationID {
				value := entity.JobInformationFeature{
					JobInformationID: f.JobInformationID,
					Feature:          f.Feature,
				}
				jobInformation.Features = append(jobInformation.Features, value)
			}
		}

		for _, p := range prefectures {
			if jobInformation.ID == p.JobInformationID {
				value := entity.JobInformationPrefecture{
					JobInformationID: p.JobInformationID,
					Prefecture:       p.Prefecture,
				}
				jobInformation.Prefectures = append(jobInformation.Prefectures, value)
			}
		}

		for _, oc := range occupations {
			if jobInformation.ID == oc.JobInformationID {
				value := entity.JobInformationOccupation{
					JobInformationID: oc.JobInformationID,
					Occupation:       oc.Occupation,
				}
				jobInformation.Occupations = append(jobInformation.Occupations, value)
			}
		}

		for _, ind := range industries {
			if jobInformation.EnterpriseID == ind.EnterpriseID {
				value := entity.EnterpriseIndustry{
					EnterpriseID: ind.EnterpriseID,
					Industry:     ind.Industry,
				}
				jobInformation.Industries = append(jobInformation.Industries, value)
			}
		}

		/************ JobListingに変換 **************/

		jobListing := entity.NewJobListing(
			jobInformation.ID,
			jobInformation.AgentStaffID,
			jobInformation.CompanyName,
			jobInformation.CorporateSiteURL,
			jobInformation.PostCode,
			jobInformation.OfficeLocation,
			jobInformation.EmployeeNumberSingle,
			jobInformation.EmployeeNumberGroup,
			jobInformation.Establishment,
			jobInformation.PublicOffering,
			jobInformation.Earnings,
			jobInformation.EarningsYear,
			jobInformation.BusinessDetail,
			jobInformation.Title,
			jobInformation.WorkDetail,
			jobInformation.WorkLocation,
			jobInformation.Transfer,
			jobInformation.TransferDetail,
			jobInformation.UnderIncome,
			jobInformation.OverIncome,
			jobInformation.Salary,
			jobInformation.Insurance,
			jobInformation.WorkTime,
			jobInformation.OvertimeAverage,
			jobInformation.FixedOvertimePayment,
			jobInformation.FixedOvertimeDetail,
			jobInformation.TrialPeriod,
			jobInformation.TrialPeriodDetail,
			jobInformation.EmploymentPeriod,
			jobInformation.EmploymentPeriodDetail,
			jobInformation.HolidayDetail,
			jobInformation.PassiveSmoking,
			jobInformation.SelectionFlow,
			jobInformation.EmploymentInsurance,
			jobInformation.AccidentInsurance,
			jobInformation.HealthInsurance,
			jobInformation.PensionInsurance,
			jobInformation.IsExternal,
			jobInformation.WorkDetailAfterHiring,
			jobInformation.WorkDetailScopeOfChange,
		)

		jobListing.IsGuaranteedInterview = jobInformation.IsGuaranteedInterview
		jobListing.JobInformationUUID = jobInformation.UUID
		jobListing.Prefectures = jobInformation.Prefectures
		jobListing.Occupations = jobInformation.Occupations
		jobListing.Industries = jobInformation.Industries
		jobListing.Features = jobInformation.Features

		jobListingList = append(jobListingList, jobListing)
	}

	/**
	ページング
	*/
	output.MaxPageNumber = getJobListingListMaxPage(jobListingList)
	output.JobListingList = getJobListingListWithPage(jobListingList, uint(inputParam.PageNumber.Int64))

	return output, nil
}
