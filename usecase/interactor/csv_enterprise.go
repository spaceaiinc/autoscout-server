package interactor

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"gopkg.in/guregu/null.v4"
	// ブラウザ操作
)

/****************************************************************************************/
/// CSV操作 API
//
//企業・請求先のcsvファイルを読み込む
type ImportEnterpriseCSVInput struct {
	CreateParam   []*entity.EnterpriseAndBillingAddress
	MissedRecords []uint
	AgentID       uint
}

type ImportEnterpriseCSVOutput struct {
	MissedRecords []uint
	OK            bool
}

func (i *EnterpriseProfileInteractorImpl) ImportEnterpriseCSV(input ImportEnterpriseCSVInput) (ImportEnterpriseCSVOutput, error) {
	var (
		output ImportEnterpriseCSVOutput
		err    error
	)

	// routesで除外されたレコードを格納する
	output.MissedRecords = input.MissedRecords

	/*****被りを防ぐため既にDBに登録されている情報を取得****/
	// エージェントIDを指定して企業情報を取得
	enterpriseProfileListInDB, err := i.enterpriseProfileRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIndustryListInDB, err := i.enterpriseIndustryRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, profile := range enterpriseProfileListInDB {
		for _, industry := range enterpriseIndustryListInDB {
			if profile.ID == industry.EnterpriseID {
				profile.Industries = append(profile.Industries, industry.Industry)
			}
		}
	}

	/***************************/

	// 企業のかぶりチェック用
	type duplicatedEnterpriseChecker struct {
		ID          uint
		CompanyName string
		PostCode    string
	}

	var (
		duplicatedEnterpriseCheckerList []duplicatedEnterpriseChecker
		isDuplicatedEnterpriseInCSV     bool
		isDuplicatedEnterpriseInDB      bool
	)

	for index, enterprise := range input.CreateParam {

		fmt.Println("csv企業情報", enterprise.AgentStaffID, enterprise.CompanyName)

		/***** 企業情報の登録 *****/
		isDuplicatedEnterpriseInCSV = false

		// CSV内で企業名と郵便番号が一致する企業は保存しない
		// 一番初めは必ず保存
		if index > 0 {
			for _, duplicatedEnterpriseInCSV := range duplicatedEnterpriseCheckerList {

				if enterprise.CompanyName == duplicatedEnterpriseInCSV.CompanyName && enterprise.PostCode == duplicatedEnterpriseInCSV.PostCode {

					enterprise.EnterpriseID = duplicatedEnterpriseInCSV.ID
					fmt.Println("csv企業ID引き継ぎ:", duplicatedEnterpriseInCSV.ID)
					isDuplicatedEnterpriseInCSV = true
					break
				}

			}
		}

		if !isDuplicatedEnterpriseInCSV {

			// DBに既に入っている企業名と郵便番号が一致する企業は保存しない
			isDuplicatedEnterpriseInDB = false

			for _, enterpriseInDB := range enterpriseProfileListInDB {

				if enterprise.CompanyName == enterpriseInDB.CompanyName && enterprise.PostCode == enterpriseInDB.PostCode {

					enterprise.EnterpriseID = enterpriseInDB.ID
					isDuplicatedEnterpriseInDB = true
					break

				}
			}

			if !isDuplicatedEnterpriseInDB {

				enterpriseProfile := entity.NewEnterpriseProfile(
					enterprise.CompanyName,
					enterprise.AgentStaffID,
					enterprise.CorporateSiteURL,
					enterprise.Representative,
					enterprise.Establishment,
					enterprise.PostCode,
					enterprise.OfficeLocation,
					enterprise.EmployeeNumberSingle,
					enterprise.EmployeeNumberGroup,
					enterprise.Capital,
					enterprise.PublicOffering,
					enterprise.EarningsYear,
					enterprise.Earnings,
					enterprise.BusinessDetail,
				)

				err = i.enterpriseProfileRepository.Create(enterpriseProfile)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
					fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
					return output, err
				}

				// 空のレコードを作成
				referenceMaterial := entity.NewEnterpriseReferenceMaterial(
					enterpriseProfile.ID,
					"",
					"",
				)

				err = i.enterpriseReferenceMaterialRepository.Create(referenceMaterial)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
					fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
					return output, err
				}

				for _, industry := range enterprise.Industries {
					industry := entity.NewEnterpriseIndustry(
						enterpriseProfile.ID,
						industry,
					)

					err = i.enterpriseIndustryRepository.Create(industry)
					if err != nil {
						fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
						fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
						fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
						fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
						fmt.Println(err)
						output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
						return output, err
					}
				}

				// 次回、重複した際に、企業IDを使えるように保存する
				duplicatedEnterpriseCheckerList = append(
					duplicatedEnterpriseCheckerList,
					duplicatedEnterpriseChecker{
						enterpriseProfile.ID,
						enterprise.CompanyName,
						enterprise.PostCode,
					},
				)

				// 請求先で使用するために企業IDを保存する
				enterprise.EnterpriseID = enterpriseProfile.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedEnterpriseCheckerList = append(duplicatedEnterpriseCheckerList, duplicatedEnterpriseChecker{
				ID:          enterprise.EnterpriseID,
				CompanyName: enterprise.CompanyName,
				PostCode:    enterprise.PostCode,
			})

		}

		/***** 請求先情報の登録 *****/
		billingAddress := entity.NewBillingAddress(
			enterprise.EnterpriseID,
			enterprise.AgentStaffIDForBillingAddress,
			enterprise.ContractPhase,
			enterprise.ContractDate,
			enterprise.PaymentPolicy,
			enterprise.CompanyName,
			enterprise.BillingAddressAddress,
			enterprise.HowToRecommend,
			enterprise.BillingAddressTitle,
		)

		err = i.billingAddressRepository.Create(billingAddress)
		if err != nil {
			fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
			fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
			fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
			fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
			fmt.Println(err)
			output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
			return output, err
		}

		for _, hs := range enterprise.HRStaffs {
			hrStaff := entity.NewBillingAddressHRStaff(
				billingAddress.ID,
				hs.HRStaffName,
				hs.HRStaffEmail,
				hs.HRStaffPhoneNumber,
			)

			err = i.billingAddressHRStaffRepository.Create(hrStaff)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
				fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
				return output, err
			}
		}

		for _, rs := range enterprise.RAStaffs {
			raStaff := entity.NewBillingAddressRAStaff(
				billingAddress.ID,
				rs.BillingAddressStaffName,
				rs.BillingAddressStaffEmail,
				rs.BillingAddressStaffPhoneNumber,
			)

			err = i.billingAddressRAStaffRepository.Create(raStaff)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
				fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
				return output, err
			}
		}
	}

	fmt.Println("企業情報の登録が完了しました。")
	fmt.Println("----------------------------------------")
	fmt.Println("除外されたレコード※最後の3行はデフォルト:", output.MissedRecords)

	output.OK = true
	return output, nil
}

// 企業・請求先のcsvファイルを読み込む
type ImportJobInformationCSVInput struct {
	CreateParam   []*entity.JobInformation
	MissedRecords []uint
	AgentID       uint
}

type ImportJobInformationCSVOutput struct {
	MissedRecords []uint
	OK            bool
}

func (i *EnterpriseProfileInteractorImpl) ImportJobInformationCSV(input ImportJobInformationCSVInput) (ImportJobInformationCSVOutput, error) {
	var (
		output ImportJobInformationCSVOutput
		err    error
	)

	// routesで除外されたレコードを格納する
	output.MissedRecords = input.MissedRecords

	// 企業のかぶりチェック用
	type duplicatedEnterpriseChecker struct {
		BillingAddressID    uint
		CompanyName         string
		PostCode            string
		BillingAddressTitle string
	}

	var (
		duplicatedEnterpriseCheckerList []duplicatedEnterpriseChecker
	)

	/*****被りを防ぐため既にDBに登録されている情報を取得****/
	// エージェントIDを指定して企業情報を取得
	enterpriseProfileListInDB, err := i.enterpriseProfileRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントIDから請求先一覧を取得
	billingAddressListInDB, err := i.billingAddressRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 企業名&&郵便番号&&請求先タイトルが一致する請求先を取得 *郵便番号は空でもOK
	for _, enterpriseInDB := range enterpriseProfileListInDB {
		for _, billingAddressInDB := range billingAddressListInDB {
			if enterpriseInDB.ID == billingAddressInDB.EnterpriseID {
				duplicatedEnterpriseCheckerList = append(duplicatedEnterpriseCheckerList, duplicatedEnterpriseChecker{
					BillingAddressID:    billingAddressInDB.ID,
					CompanyName:         enterpriseInDB.CompanyName,
					PostCode:            enterpriseInDB.PostCode,
					BillingAddressTitle: billingAddressInDB.Title,
				})
			}
		}
	}

	/****************************************************/

	for _, jobInformation := range input.CreateParam {

		/***** 企業情報の登録 *****/
		// DBに既に入っている企業名と郵便番号が一致する企業は保存しない

		for _, duplicatedEnterpriseChecker := range duplicatedEnterpriseCheckerList {
			fmt.Println("db企業情報", duplicatedEnterpriseChecker.CompanyName, duplicatedEnterpriseChecker.PostCode, duplicatedEnterpriseChecker.BillingAddressTitle)
			fmt.Println("csv企業情報", jobInformation.CompanyName, jobInformation.PostCode, jobInformation.BillingAddressTitle)

			// 企業名&&郵便番号&&請求先タイトルが一致する請求先を取得 *郵便番号は空でもOK
			if jobInformation.CompanyName == duplicatedEnterpriseChecker.CompanyName &&
				(jobInformation.PostCode == duplicatedEnterpriseChecker.PostCode || jobInformation.PostCode == "") &&
				jobInformation.BillingAddressTitle == duplicatedEnterpriseChecker.BillingAddressTitle {

				jobInformation.BillingAddressID = duplicatedEnterpriseChecker.BillingAddressID
				break
			}
		}

		// 企業情報が登録されていない場合はスキップ
		if jobInformation.BillingAddressID == 0 {
			fmt.Println("どの請求先ともヒットしなかったためスキップします。求人タイトル:", jobInformation.Title)
			output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
			continue
		}

		// デフォルト値をセット
		// 登録状況（0: 本登録, 1: 仮登録）
		jobInformation.RegisterPhase = null.NewInt(0, true)

		log.Println("求人情報の登録を開始します。求人タイトル:", jobInformation)

		/********* 求人情報 *********/

		err = i.jobInformationRepository.Create(jobInformation)
		if err != nil {
			fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
			fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
			fmt.Println("請求先ID", jobInformation.BillingAddressID)
			fmt.Println(err)
			output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
			return output, err
		}

		for _, target := range jobInformation.Targets {
			t := entity.NewJobInformationTarget(
				jobInformation.ID,
				target.Target,
			)

			err = i.jobInfoTargetRepository.Create(t)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
		}

		for _, feature := range jobInformation.Features {
			f := entity.NewJobInformationFeature(
				jobInformation.ID,
				feature.Feature,
			)

			err = i.jobInfoFeatureRepository.Create(f)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
		}

		for _, prefecture := range jobInformation.Prefectures {
			p := entity.NewJobInformationPrefecture(
				jobInformation.ID,
				prefecture.Prefecture,
			)

			err = i.jobInfoPrefectureRepository.Create(p)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
		}

		for _, workCharmPoint := range jobInformation.WorkCharmPoints {
			wcp := entity.NewJobInformationWorkCharmPoint(
				jobInformation.ID,
				workCharmPoint.Title,
				workCharmPoint.Contents,
			)

			err = i.jobInfoWorkCharmPointRepository.Create(wcp)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
		}

		for _, employmentStatus := range jobInformation.EmploymentStatuses {
			es := entity.NewJobInformationEmploymentStatus(
				jobInformation.ID,
				employmentStatus.EmploymentStatus,
			)

			err = i.jobInfoEmploymentStatusRepository.Create(es)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
		}

		// 必要条件
		for _, requiredCondition := range jobInformation.RequiredConditions {
			rc := entity.NewJobInformationRequiredCondition(
				jobInformation.ID,
				requiredCondition.IsCommon,
				requiredCondition.RequiredManagement,
			)

			err = i.jobInfoRequiredConditionRepository.Create(rc)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
			requiredCondition.ID = rc.ID

			// 必要資格　複数
			for _, requiredLicense := range requiredCondition.RequiredLicenses {
				rl := entity.NewJobInformationRequiredLicense(
					requiredCondition.ID,
					requiredLicense.License,
				)

				err = i.jobInfoRequiredLicenseRepository.Create(rl)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
					fmt.Println("請求先ID:", jobInformation.BillingAddressID)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
					return output, err
				}
			}

			// 必要PCツール　単数
			for _, requiredPCTool := range requiredCondition.RequiredPCTools {
				rpt := entity.NewJobInformationRequiredPCTool(
					requiredCondition.ID,
					requiredPCTool.Tool,
				)

				err = i.jobInfoRequiredPCToolRepository.Create(rpt)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
					fmt.Println("請求先ID:", jobInformation.BillingAddressID)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
					return output, err
				}
			}

			// 必要言語スキル　単数
			// いずれかの項目が入力されている場合は必要言語スキルを登録する
			if len(requiredCondition.RequiredLanguages.LanguageTypes) > 0 ||
				requiredCondition.RequiredLanguages.LanguageLevel.Valid ||
				requiredCondition.RequiredLanguages.Toeic.Valid ||
				requiredCondition.RequiredLanguages.ToeflIBT.Valid ||
				requiredCondition.RequiredLanguages.ToeflPBT.Valid {
				requiredLanguages := entity.NewJobInformationRequiredLanguage(
					requiredCondition.ID,
					requiredCondition.RequiredLanguages.LanguageLevel,
					requiredCondition.RequiredLanguages.Toeic,
					requiredCondition.RequiredLanguages.ToeflIBT,
					requiredCondition.RequiredLanguages.ToeflPBT,
				)

				err = i.jobInfoRequiredLanguageRepository.Create(requiredLanguages)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
					fmt.Println("請求先ID:", jobInformation.BillingAddressID)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
					return output, err
				}

				for _, languageType := range requiredCondition.RequiredLanguages.LanguageTypes {
					lt := entity.NewJobInformationRequiredLanguageType(
						requiredLanguages.ID,
						languageType.LanguageType,
					)

					err = i.jobInfoRequiredLanguageTypeRepository.Create(lt)
					if err != nil {
						fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
						fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
						fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
						fmt.Println("請求先ID:", jobInformation.BillingAddressID)
						fmt.Println(err)
						output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
						return output, err
					}
				}
			}

			// 開発スキル　　言語,OS 各1つずつ
			for _, requiredExperienceDevelopment := range requiredCondition.RequiredExperienceDevelopments {
				red := entity.NewJobInformationRequiredExperienceDevelopment(
					requiredCondition.ID,
					requiredExperienceDevelopment.DevelopmentCategory,
					requiredExperienceDevelopment.ExperienceYear,
					requiredExperienceDevelopment.ExperienceMonth,
				)

				err = i.jobInfoRequiredExperienceDevelopmentRepository.Create(red)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
					fmt.Println("請求先ID:", jobInformation.BillingAddressID)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
					return output, err
				}

				for _, experienceDevelopmentType := range requiredExperienceDevelopment.ExperienceDevelopmentTypes {
					dt := entity.NewJobInformationRequiredExperienceDevelopmentType(
						red.ID,
						experienceDevelopmentType.DevelopmentType,
					)

					err = i.jobInfoRequiredExperienceDevelopmentTypeRepository.Create(dt)
					if err != nil {
						fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
						fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
						fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
						fmt.Println("請求先ID:", jobInformation.BillingAddressID)
						fmt.Println(err)
						output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
						return output, err
					}
				}
			}

			// 必要経験　　単数
			if len(requiredCondition.RequiredExperienceJobs.ExperienceIndustries) > 0 ||
				len(requiredCondition.RequiredExperienceJobs.ExperienceOccupations) > 0 ||
				requiredCondition.RequiredExperienceJobs.ExperienceYear.Valid ||
				requiredCondition.RequiredExperienceJobs.ExperienceMonth.Valid {
				requiredExperienceJobs := entity.NewJobInformationRequiredExperienceJob(
					requiredCondition.ID,
					requiredCondition.RequiredExperienceJobs.ExperienceYear,
					requiredCondition.RequiredExperienceJobs.ExperienceMonth,
				)

				err = i.jobInfoRequiredExperienceJobRepository.Create(requiredExperienceJobs)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
					fmt.Println("請求先ID:", jobInformation.BillingAddressID)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
					return output, err
				}

				for _, experienceIndustry := range requiredCondition.RequiredExperienceJobs.ExperienceIndustries {
					ei := entity.NewJobInformationRequiredExperienceIndustry(
						requiredExperienceJobs.ID,
						experienceIndustry.ExperienceIndustry,
					)

					err = i.jobInfoRequiredExperienceIndustryRepository.Create(ei)
					if err != nil {
						fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
						fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
						fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
						fmt.Println("請求先ID:", jobInformation.BillingAddressID)
						fmt.Println(err)
						output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
						return output, err
					}
				}

				for _, experienceOccupation := range requiredCondition.RequiredExperienceJobs.ExperienceOccupations {
					eo := entity.NewJobInformationRequiredExperienceOccupation(
						requiredExperienceJobs.ID,
						experienceOccupation.ExperienceOccupation,
					)

					err = i.jobInfoRequiredExperienceOccupationRepository.Create(eo)
					if err != nil {
						fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
						fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
						fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
						fmt.Println("請求先ID:", jobInformation.BillingAddressID)
						fmt.Println(err)
						output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
						return output, err
					}
				}
			}
		}

		for _, requiredSocialExperience := range jobInformation.RequiredSocialExperiences {
			rse := entity.NewJobInformationRequiredSocialExperience(
				jobInformation.ID,
				requiredSocialExperience.SocialExperienceType,
			)

			err = i.jobInfoRequiredSocialExperienceRepository.Create(rse)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
		}

		for _, selectionFlowPattern := range jobInformation.SelectionFlowPatterns {
			sfp := entity.NewJobInformationSelectionFlowPattern(
				jobInformation.ID,
				selectionFlowPattern.PublicStatus,
				selectionFlowPattern.FlowTitle,
				selectionFlowPattern.FlowPattern,
				false,
			)

			err = i.jobInfoSelectionFlowPatternRepository.Create(sfp)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}

			for _, selectionInformation := range selectionFlowPattern.SelectionInformations {
				si := entity.NewJobInformationSelectionInformation(
					sfp.ID,
					selectionInformation.SelectionType,
					selectionInformation.SelectionPoint,
					selectionInformation.PassedExample,
					selectionInformation.FailExample,
					selectionInformation.PassingRate,
					selectionInformation.IsQuestionnairy,
				)

				err = i.jobInfoSelectionInformationRepository.Create(si)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
					fmt.Println("請求先ID", jobInformation.BillingAddressID)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
					return output, err
				}
			}
		}

		for _, occupation := range jobInformation.Occupations {
			oc := entity.NewJobInformationOccupation(
				jobInformation.ID,
				occupation.Occupation,
			)

			err = i.jobInfoOccupationRepository.Create(oc)
			if err != nil {
				fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
				fmt.Println("企業登録時にエラーの求人の行:", jobInformation.RecordLine)
				fmt.Println("企業登録時にエラーの企業ID:", jobInformation.EnterpriseID, "企業担当者ID:", jobInformation.AgentStaffID)
				fmt.Println("請求先ID:", jobInformation.BillingAddressID)
				fmt.Println(err)
				output.MissedRecords = append(output.MissedRecords, jobInformation.RecordLine)
				return output, err
			}
		}
	}

	fmt.Println("企業情報の登録が完了しました。")
	fmt.Println("----------------------------------------")
	fmt.Println("除外されたレコード※最後の3行はデフォルト:", output.MissedRecords)

	output.OK = true
	return output, nil
}

/**************************************************************************************/
// 求人企業のcsvファイルを読み込む (サーカス用)
//
type ImportEnterpriseCSVForCircusInput struct {
	CreateParam   []*entity.EnterpriseAndJobInformation
	MissedRecords []uint
	AgentID       uint
}

type ImportEnterpriseCSVForCircusOutput struct {
	MissedRecords []uint
	OK            bool
}

func (i *EnterpriseProfileInteractorImpl) ImportEnterpriseCSVForCircus(input ImportEnterpriseCSVForCircusInput) (ImportEnterpriseCSVForCircusOutput, error) {
	var (
		output ImportEnterpriseCSVForCircusOutput
		err    error
	)

	// routesで除外されたレコードを格納する
	output.MissedRecords = input.MissedRecords

	/*****被りを防ぐため既にDBに登録されている情報を取得****/
	// エージェントIDを指定して企業情報を取得
	enterpriseProfileListInDB, err := i.enterpriseProfileRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントIDから請求先一覧を取得
	billingAddressListInDB, err := i.billingAddressRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/***************************/

	// 企業のかぶりチェック用
	type duplicatedEnterpriseChecker struct {
		ID                 uint
		CompanyName        string
		CircusEnterpriseID uint
	}

	var (
		duplicatedEnterpriseCheckerList     []duplicatedEnterpriseChecker
		isDuplicatedEnterpriseInCSV         bool
		isDuplicatedEnterpriseInDB          bool
		duplicatedBillingAddressCheckerList []entity.BillingAddress
		isDuplicatedBillingAddressInCSV     bool
		isDuplicatedBillingAddressInDB      bool
	)

	for index, enterprise := range input.CreateParam {

		fmt.Println("csv企業情報", enterprise.AgentStaffID, enterprise.CompanyName)
		if enterprise.CompanyName == "" {
			fmt.Println("企業名が空欄のためスキップします。")
			continue
		}

		/***** 企業情報の登録 *****/
		isDuplicatedEnterpriseInCSV = false

		// CSV内で企業名とサーカス内の企業IDが一致する企業は保存しない
		// 一番初めは必ず保存
		if index > 0 {
			for _, duplicatedEnterpriseInCSV := range duplicatedEnterpriseCheckerList {

				if enterprise.CompanyName == duplicatedEnterpriseInCSV.CompanyName &&
					(enterprise.CircusEnterpriseID == duplicatedEnterpriseInCSV.CircusEnterpriseID || enterprise.CircusEnterpriseID == 0) {
					enterprise.EnterpriseID = duplicatedEnterpriseInCSV.ID
					fmt.Println("csv企業ID引き継ぎ:", duplicatedEnterpriseInCSV.ID)
					isDuplicatedEnterpriseInCSV = true
					break
				}

			}
		}

		if !isDuplicatedEnterpriseInCSV {

			// DBに既に入っている企業名と郵便番号が一致する企業は保存しない
			isDuplicatedEnterpriseInDB = false

			for _, enterpriseInDB := range enterpriseProfileListInDB {

				if enterprise.CompanyName == enterpriseInDB.CompanyName {

					enterprise.EnterpriseID = enterpriseInDB.ID
					isDuplicatedEnterpriseInDB = true
					break

				}
			}

			if !isDuplicatedEnterpriseInDB {

				enterpriseProfile := entity.NewEnterpriseProfile(
					enterprise.CompanyName,
					enterprise.AgentStaffID,
					enterprise.CorporateSiteURL,
					enterprise.Representative,
					enterprise.Establishment,
					enterprise.PostCode,
					enterprise.OfficeLocation,
					enterprise.EmployeeNumberSingle,
					enterprise.EmployeeNumberGroup,
					enterprise.Capital,
					enterprise.PublicOffering,
					enterprise.EarningsYear,
					enterprise.Earnings,
					enterprise.BusinessDetail,
				)

				err = i.enterpriseProfileRepository.Create(enterpriseProfile)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				// 空のレコードを作成
				referenceMaterial := entity.NewEnterpriseReferenceMaterial(
					enterpriseProfile.ID,
					"",
					"",
				)

				err = i.enterpriseReferenceMaterialRepository.Create(referenceMaterial)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
					fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
					return output, err
				}

				for _, industry := range enterprise.Industries {
					industry := entity.NewEnterpriseIndustry(
						enterpriseProfile.ID,
						industry,
					)

					err = i.enterpriseIndustryRepository.Create(industry)
					if err != nil {
						fmt.Println(err)
						return output, err
					}
				}

				// 請求先で使用するために企業IDを保存する
				enterprise.EnterpriseID = enterpriseProfile.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedEnterpriseCheckerList = append(
				duplicatedEnterpriseCheckerList,
				duplicatedEnterpriseChecker{
					ID:                 enterprise.EnterpriseID,
					CompanyName:        enterprise.CompanyName,
					CircusEnterpriseID: enterprise.CircusEnterpriseID,
				},
			)

		}

		/***** 請求先情報の登録 *****/
		// CSV内で請求先が一致する企業は保存しない
		isDuplicatedBillingAddressInCSV = false

		// 一番初めは必ず保存する
		if index > 0 {
			for _, duplicatedbBillingAddressInCSV := range duplicatedBillingAddressCheckerList {
				// すべての請求先担当者情報が一致した場合は飛ばす
				if enterprise.CompanyName == duplicatedbBillingAddressInCSV.CompanyName {
					enterprise.BillingAddressID = duplicatedbBillingAddressInCSV.ID
					isDuplicatedBillingAddressInCSV = true
					break
				}
			}
		}

		if !isDuplicatedBillingAddressInCSV {

			// DBに既に入っている企業名と郵便番号が一致する企業は保存しない
			isDuplicatedBillingAddressInDB = false

			for _, billingAddressInDB := range billingAddressListInDB {
				if enterprise.CompanyName == billingAddressInDB.CompanyName {
					enterprise.BillingAddressID = billingAddressInDB.ID
					isDuplicatedBillingAddressInDB = true
					break
				}
			}

			if !isDuplicatedBillingAddressInDB {

				billingAddress := entity.NewBillingAddress(
					enterprise.EnterpriseID,
					enterprise.AgentStaffID,
					enterprise.ContractPhase,
					enterprise.ContractDate,
					enterprise.PaymentPolicy,
					enterprise.CompanyName,
					enterprise.BillingAddressAddress,
					enterprise.HowToRecommend,
					enterprise.BillingAddressTitle,
				)

				err = i.billingAddressRepository.Create(billingAddress)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				// 求人登録時に使用するために請求先IDを保存する
				enterprise.BillingAddressID = billingAddress.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedBillingAddressCheckerList = append(duplicatedBillingAddressCheckerList, entity.BillingAddress{
				ID:             enterprise.BillingAddressID,
				CompanyName:    enterprise.CompanyName,
				ContractDate:   enterprise.ContractDate,
				PaymentPolicy:  enterprise.PaymentPolicy,
				HowToRecommend: enterprise.HowToRecommend,
				AgentStaffID:   enterprise.AgentStaffIDForBillingAddress,
				HRStaffs:       enterprise.HRStaffs,
				RAStaffs:       enterprise.RAStaffs,
			})
		}

		/********* 求人情報 *********/
		jobInformation := entity.NewJobInformation(
			enterprise.BillingAddressID,
			enterprise.Title,
			enterprise.RecruitmentState,
			enterprise.ExpirationDate,
			enterprise.WorkDetail,
			enterprise.NumberOfHires,
			enterprise.WorkLocation,
			enterprise.Transfer,
			enterprise.TransferDetail,
			enterprise.UnderIncome,
			enterprise.OverIncome,
			enterprise.Salary,
			enterprise.Insurance,
			enterprise.WorkTime,
			enterprise.OvertimeAverage,
			enterprise.FixedOvertimePayment,
			enterprise.FixedOvertimeDetail,
			enterprise.TrialPeriod,
			enterprise.TrialPeriodDetail,
			enterprise.EmploymentPeriod,
			enterprise.EmploymentPeriodDetail,
			enterprise.HolidayType,
			enterprise.HolidayDetail,
			enterprise.PassiveSmoking,
			enterprise.SelectionFlow,
			enterprise.Gender,
			enterprise.Nationality,
			enterprise.FinalEducation,
			enterprise.SchoolLevel,
			enterprise.MedicalHistory,
			enterprise.AgeUnder,
			enterprise.AgeOver,
			enterprise.JobChange,
			enterprise.ShortResignation,
			enterprise.ShortResignationRemarks,
			enterprise.SocialExperienceYear,
			enterprise.SocialExperienceMonth,
			enterprise.Appearance,
			enterprise.Communication,
			enterprise.Thinking,
			enterprise.TargetDetail,
			enterprise.Commission,
			enterprise.CommissionRate,
			enterprise.CommissionDetail,
			enterprise.RefundPolicy,
			enterprise.RequiredExperienceJobDetail,
			enterprise.SecretMemo,
			enterprise.RequiredDocumentsDetail,
			enterprise.EmploymentInsurance,
			enterprise.AccidentInsurance,
			enterprise.HealthInsurance,
			enterprise.PensionInsurance,
			enterprise.RegisterPhase,
			enterprise.StudyCategory,
			enterprise.DriverLicence,
			enterprise.WordSkill,
			enterprise.ExcelSkill,
			enterprise.PowerPointSkill,
			false, // import時はautoscout求人のためfalse
			enterprise.WorkDetailAfterHiring,
			enterprise.WorkDetailScopeOfChange,
			enterprise.OfferRate,
			enterprise.DocumentPassingRate,
			enterprise.NumberOfRecentApplications,
			enterprise.IsGuaranteedInterview,
		)

		err = i.jobInformationRepository.Create(jobInformation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, prefecture := range enterprise.Prefectures {
			p := entity.NewJobInformationPrefecture(
				jobInformation.ID,
				prefecture.Prefecture,
			)

			err = i.jobInfoPrefectureRepository.Create(p)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, employmentStatus := range enterprise.EmploymentStatuses {
			es := entity.NewJobInformationEmploymentStatus(
				jobInformation.ID,
				employmentStatus.EmploymentStatus,
			)

			err = i.jobInfoEmploymentStatusRepository.Create(es)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 募集職種
		for _, occupation := range enterprise.Occupations {
			oc := entity.NewJobInformationOccupation(
				jobInformation.ID,
				occupation.Occupation,
			)

			err = i.jobInfoOccupationRepository.Create(oc)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

	}

	fmt.Println("企業情報の登録が完了しました。")
	fmt.Println("----------------------------------------")
	fmt.Println("除外されたレコード:", output.MissedRecords)

	output.OK = true
	return output, nil
}

/**************************************************************************************/
// 求人企業のcsvファイルを読み込む (サーカス用)
//
type ImportEnterpriseCSVForAgentBankInput struct {
	CreateParam   []*entity.EnterpriseAndJobInformation
	MissedRecords []uint
	AgentID       uint
}

type ImportEnterpriseCSVForAgentBankOutput struct {
	MissedRecords []uint
	OK            bool
}

func (i *EnterpriseProfileInteractorImpl) ImportEnterpriseCSVForAgentBank(input ImportEnterpriseCSVForAgentBankInput) (ImportEnterpriseCSVForAgentBankOutput, error) {
	var (
		output ImportEnterpriseCSVForAgentBankOutput
		err    error
	)

	// routesで除外されたレコードを格納する
	output.MissedRecords = input.MissedRecords

	/*****被りを防ぐため既にDBに登録されている情報を取得****/
	// エージェントIDを指定して企業情報を取得
	enterpriseProfileListInDB, err := i.enterpriseProfileRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントIDから請求先一覧を取得
	billingAddressListInDB, err := i.billingAddressRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/***************************/

	// 企業のかぶりチェック用
	type duplicatedEnterpriseChecker struct {
		ID                 uint
		CompanyName        string
		CircusEnterpriseID uint
	}

	var (
		duplicatedEnterpriseCheckerList     []duplicatedEnterpriseChecker
		isDuplicatedEnterpriseInCSV         bool
		isDuplicatedEnterpriseInDB          bool
		duplicatedBillingAddressCheckerList []entity.BillingAddress
		isDuplicatedBillingAddressInCSV     bool
		isDuplicatedBillingAddressInDB      bool
	)

	for index, enterprise := range input.CreateParam {

		fmt.Println("csv企業情報", enterprise.AgentStaffID, enterprise.CompanyName)
		if enterprise.CompanyName == "" {
			fmt.Println("企業名が空欄のためスキップします。")
			continue
		}

		/***** 企業情報の登録 *****/
		isDuplicatedEnterpriseInCSV = false

		// CSV内で企業名とサーカス内の企業IDが一致する企業は保存しない
		// 一番初めは必ず保存
		if index > 0 {
			for _, duplicatedEnterpriseInCSV := range duplicatedEnterpriseCheckerList {
				if enterprise.CompanyName == duplicatedEnterpriseInCSV.CompanyName &&
					(enterprise.CircusEnterpriseID == duplicatedEnterpriseInCSV.CircusEnterpriseID || enterprise.CircusEnterpriseID == 0) {
					enterprise.EnterpriseID = duplicatedEnterpriseInCSV.ID
					fmt.Println("csv企業ID引き継ぎ:", duplicatedEnterpriseInCSV.ID)
					isDuplicatedEnterpriseInCSV = true
					break
				}
			}
		}

		if !isDuplicatedEnterpriseInCSV {

			// DBに既に入っている企業名と郵便番号が一致する企業は保存しない
			isDuplicatedEnterpriseInDB = false

			for _, enterpriseInDB := range enterpriseProfileListInDB {

				if enterprise.CompanyName == enterpriseInDB.CompanyName {

					enterprise.EnterpriseID = enterpriseInDB.ID
					isDuplicatedEnterpriseInDB = true
					break

				}
			}

			if !isDuplicatedEnterpriseInDB {

				enterpriseProfile := entity.NewEnterpriseProfile(
					enterprise.CompanyName,
					enterprise.AgentStaffID,
					enterprise.CorporateSiteURL,
					enterprise.Representative,
					enterprise.Establishment,
					enterprise.PostCode,
					enterprise.OfficeLocation,
					enterprise.EmployeeNumberSingle,
					enterprise.EmployeeNumberGroup,
					enterprise.Capital,
					enterprise.PublicOffering,
					enterprise.EarningsYear,
					enterprise.Earnings,
					enterprise.BusinessDetail,
				)

				err = i.enterpriseProfileRepository.Create(enterpriseProfile)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				// 空のレコードを作成
				referenceMaterial := entity.NewEnterpriseReferenceMaterial(
					enterpriseProfile.ID,
					"",
					"",
				)

				err = i.enterpriseReferenceMaterialRepository.Create(referenceMaterial)
				if err != nil {
					fmt.Println("除外された求人の行※最後の3行はデフォルト:", output.MissedRecords)
					fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
					fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
					fmt.Println(err)
					output.MissedRecords = append(output.MissedRecords, enterprise.RecordLine)
					return output, err
				}

				for _, industry := range enterprise.Industries {
					industry := entity.NewEnterpriseIndustry(
						enterpriseProfile.ID,
						industry,
					)

					err = i.enterpriseIndustryRepository.Create(industry)
					if err != nil {
						fmt.Println(err)
						return output, err
					}
				}

				// 請求先で使用するために企業IDを保存する
				enterprise.EnterpriseID = enterpriseProfile.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedEnterpriseCheckerList = append(
				duplicatedEnterpriseCheckerList,
				duplicatedEnterpriseChecker{
					ID:                 enterprise.EnterpriseID,
					CompanyName:        enterprise.CompanyName,
					CircusEnterpriseID: enterprise.CircusEnterpriseID,
				},
			)

		}

		/***** 請求先情報の登録 *****/
		// CSV内で請求先が一致する企業は保存しない
		isDuplicatedBillingAddressInCSV = false

		// 一番初めは必ず保存する
		if index > 0 {
			for _, duplicatedbBillingAddressInCSV := range duplicatedBillingAddressCheckerList {
				// すべての請求先担当者情報が一致した場合は飛ばす
				if enterprise.CompanyName == duplicatedbBillingAddressInCSV.CompanyName {
					enterprise.BillingAddressID = duplicatedbBillingAddressInCSV.ID
					isDuplicatedBillingAddressInCSV = true
					break
				}
			}
		}

		if !isDuplicatedBillingAddressInCSV {

			// DBに既に入っている企業名と郵便番号が一致する企業は保存しない
			isDuplicatedBillingAddressInDB = false

			for _, billingAddressInDB := range billingAddressListInDB {
				if enterprise.CompanyName == billingAddressInDB.CompanyName {
					enterprise.BillingAddressID = billingAddressInDB.ID
					isDuplicatedBillingAddressInDB = true
					break
				}
			}

			if !isDuplicatedBillingAddressInDB {

				billingAddress := entity.NewBillingAddress(
					enterprise.EnterpriseID,
					enterprise.AgentStaffID,
					enterprise.ContractPhase,
					enterprise.ContractDate,
					enterprise.PaymentPolicy,
					enterprise.CompanyName,
					enterprise.BillingAddressAddress,
					enterprise.HowToRecommend,
					enterprise.BillingAddressTitle,
				)

				err = i.billingAddressRepository.Create(billingAddress)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				// 求人登録時に使用するために請求先IDを保存する
				enterprise.BillingAddressID = billingAddress.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedBillingAddressCheckerList = append(duplicatedBillingAddressCheckerList, entity.BillingAddress{
				ID:             enterprise.BillingAddressID,
				CompanyName:    enterprise.CompanyName,
				ContractDate:   enterprise.ContractDate,
				PaymentPolicy:  enterprise.PaymentPolicy,
				HowToRecommend: enterprise.HowToRecommend,
				AgentStaffID:   enterprise.AgentStaffIDForBillingAddress,
				HRStaffs:       enterprise.HRStaffs,
				RAStaffs:       enterprise.RAStaffs,
			})
		}

		/********* 求人情報 *********/
		jobInformation := entity.NewJobInformation(
			enterprise.BillingAddressID,
			enterprise.Title,
			enterprise.RecruitmentState,
			enterprise.ExpirationDate,
			enterprise.WorkDetail,
			enterprise.NumberOfHires,
			enterprise.WorkLocation,
			enterprise.Transfer,
			enterprise.TransferDetail,
			enterprise.UnderIncome,
			enterprise.OverIncome,
			enterprise.Salary,
			enterprise.Insurance,
			enterprise.WorkTime,
			enterprise.OvertimeAverage,
			enterprise.FixedOvertimePayment,
			enterprise.FixedOvertimeDetail,
			enterprise.TrialPeriod,
			enterprise.TrialPeriodDetail,
			enterprise.EmploymentPeriod,
			enterprise.EmploymentPeriodDetail,
			enterprise.HolidayType,
			enterprise.HolidayDetail,
			enterprise.PassiveSmoking,
			enterprise.SelectionFlow,
			enterprise.Gender,
			enterprise.Nationality,
			enterprise.FinalEducation,
			enterprise.SchoolLevel,
			enterprise.MedicalHistory,
			enterprise.AgeUnder,
			enterprise.AgeOver,
			enterprise.JobChange,
			enterprise.ShortResignation,
			enterprise.ShortResignationRemarks,
			enterprise.SocialExperienceYear,
			enterprise.SocialExperienceMonth,
			enterprise.Appearance,
			enterprise.Communication,
			enterprise.Thinking,
			enterprise.TargetDetail,
			enterprise.Commission,
			enterprise.CommissionRate,
			enterprise.CommissionDetail,
			enterprise.RefundPolicy,
			enterprise.RequiredExperienceJobDetail,
			enterprise.SecretMemo,
			enterprise.RequiredDocumentsDetail,
			enterprise.EmploymentInsurance,
			enterprise.AccidentInsurance,
			enterprise.HealthInsurance,
			enterprise.PensionInsurance,
			enterprise.RegisterPhase,
			enterprise.StudyCategory,
			enterprise.DriverLicence,
			enterprise.WordSkill,
			enterprise.ExcelSkill,
			enterprise.PowerPointSkill,
			false, // import時はautoscout求人のためfalse
			enterprise.WorkDetailAfterHiring,
			enterprise.WorkDetailScopeOfChange,
			enterprise.OfferRate,
			enterprise.DocumentPassingRate,
			enterprise.NumberOfRecentApplications,
			enterprise.IsGuaranteedInterview,
		)

		err = i.jobInformationRepository.Create(jobInformation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, prefecture := range enterprise.Prefectures {
			p := entity.NewJobInformationPrefecture(
				jobInformation.ID,
				prefecture.Prefecture,
			)

			err = i.jobInfoPrefectureRepository.Create(p)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, employmentStatus := range enterprise.EmploymentStatuses {
			es := entity.NewJobInformationEmploymentStatus(
				jobInformation.ID,
				employmentStatus.EmploymentStatus,
			)

			err = i.jobInfoEmploymentStatusRepository.Create(es)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 募集職種
		for _, occupation := range enterprise.Occupations {
			oc := entity.NewJobInformationOccupation(
				jobInformation.ID,
				occupation.Occupation,
			)

			err = i.jobInfoOccupationRepository.Create(oc)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

	}

	fmt.Println("企業情報の登録が完了しました。")
	fmt.Println("----------------------------------------")
	fmt.Println("除外されたレコード:", output.MissedRecords)

	output.OK = true
	return output, nil
}

/**************************************************************************************/
// 企業と請求先、求人のリストを読み込む
//
type ImportEnterpriseJSONInput struct {
	CreateParam  []entity.EnterpriseAndJobInformation
	AgentStaffID uint
}

type ImportEnterpriseJSONOutput struct {
	OK bool
}

func (i *EnterpriseProfileInteractorImpl) ImportEnterpriseJSON(input ImportEnterpriseJSONInput) (ImportEnterpriseJSONOutput, error) {
	var (
		output ImportEnterpriseJSONOutput
		err    error
	)

	/*****被りを防ぐため既にDBに登録されている情報を取得****/
	// エージェント担当者を取得
	agentStaff, err := i.agentStaffRepository.FindByID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントIDを指定して企業情報を取得
	enterpriseProfileListInDB, err := i.enterpriseProfileRepository.GetByAgentID(agentStaff.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントIDから請求先一覧を取得
	billingAddressListInDB, err := i.billingAddressRepository.GetByAgentID(agentStaff.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobInformationListInDB, err := i.jobInformationRepository.GetByAgentIDAndExternalType(agentStaff.AgentID, entity.JobInformatinoExternalTypeAgentBank)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/***************************/

	// 企業のかぶりチェック用
	type duplicatedEnterpriseChecker struct {
		ID                 uint
		CompanyName        string
		CircusEnterpriseID uint
	}

	var (
		duplicatedEnterpriseCheckerList     []duplicatedEnterpriseChecker
		isDuplicatedEnterpriseInCSV         bool
		isDuplicatedEnterpriseInDB          bool
		duplicatedBillingAddressCheckerList []entity.BillingAddress
		isDuplicatedBillingAddressInCSV     bool
		isDuplicatedBillingAddressInDB      bool
	)

enterpriseLoop:
	for index, enterprise := range input.CreateParam {
		// 初期値設定
		enterprise.AgentStaffID = input.AgentStaffID
		enterprise.AgentStaffIDForBillingAddress = input.AgentStaffID

		fmt.Println("csv企業情報", enterprise.AgentStaffID, enterprise.CompanyName)
		if enterprise.CompanyName == "" {
			fmt.Println("企業名が空欄のためスキップします。")
			continue
		}

		/***** 企業情報の登録 *****/
		isDuplicatedEnterpriseInCSV = false

		// CSV内で企業名が一致する企業は保存しない
		// 一番初めは必ず保存
		if index > 0 {
			for _, duplicatedEnterpriseInCSV := range duplicatedEnterpriseCheckerList {

				if enterprise.CompanyName == duplicatedEnterpriseInCSV.CompanyName {
					enterprise.EnterpriseID = duplicatedEnterpriseInCSV.ID
					fmt.Println("csv企業ID引き継ぎ:", duplicatedEnterpriseInCSV.ID)
					isDuplicatedEnterpriseInCSV = true
					break
				}

			}
		}

		if !isDuplicatedEnterpriseInCSV {

			// DBに既に入っている企業名が一致する企業は保存しない
			isDuplicatedEnterpriseInDB = false

			for _, enterpriseInDB := range enterpriseProfileListInDB {

				if enterprise.CompanyName == enterpriseInDB.CompanyName {

					enterprise.EnterpriseID = enterpriseInDB.ID
					isDuplicatedEnterpriseInDB = true
					break

				}
			}

			if !isDuplicatedEnterpriseInDB {

				enterpriseProfile := entity.NewEnterpriseProfile(
					enterprise.CompanyName,
					enterprise.AgentStaffID,
					enterprise.CorporateSiteURL,
					enterprise.Representative,
					enterprise.Establishment,
					enterprise.PostCode,
					enterprise.OfficeLocation,
					enterprise.EmployeeNumberSingle,
					enterprise.EmployeeNumberGroup,
					enterprise.Capital,
					enterprise.PublicOffering,
					enterprise.EarningsYear,
					enterprise.Earnings,
					enterprise.BusinessDetail,
				)

				err = i.enterpriseProfileRepository.Create(enterpriseProfile)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				// 空のレコードを作成
				referenceMaterial := entity.NewEnterpriseReferenceMaterial(
					enterpriseProfile.ID,
					"",
					"",
				)

				err = i.enterpriseReferenceMaterialRepository.Create(referenceMaterial)
				if err != nil {
					fmt.Println("企業登録時にエラーの求人の行:", enterprise.RecordLine)
					fmt.Println("企業登録時にエラーの企業ID:", enterprise.EnterpriseID, "企業担当者ID:", enterprise.AgentStaffID)
					fmt.Println("請求先ID", enterprise.BillingAddressID, "請求先担当者ID", enterprise.AgentStaffIDForBillingAddress)
					fmt.Println(err)
					return output, err
				}

				for _, industry := range enterprise.Industries {
					industry := entity.NewEnterpriseIndustry(
						enterpriseProfile.ID,
						industry,
					)

					err = i.enterpriseIndustryRepository.Create(industry)
					if err != nil {
						fmt.Println(err)
						return output, err
					}
				}

				// 請求先で使用するために企業IDを保存する
				enterprise.EnterpriseID = enterpriseProfile.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedEnterpriseCheckerList = append(
				duplicatedEnterpriseCheckerList,
				duplicatedEnterpriseChecker{
					ID:                 enterprise.EnterpriseID,
					CompanyName:        enterprise.CompanyName,
					CircusEnterpriseID: enterprise.CircusEnterpriseID,
				},
			)

		}

		/***** 請求先情報の登録 *****/
		// CSV内の重複判定
		isDuplicatedBillingAddressInCSV = false

		// 一番初めは必ず保存する
		if index > 0 {
			for _, duplicatedbBillingAddressInCSV := range duplicatedBillingAddressCheckerList {
				// 請求先社名と企業名が一致した場合は飛ばす
				if enterprise.BillingAddressCompanyName == duplicatedbBillingAddressInCSV.CompanyName &&
					enterprise.CompanyName == duplicatedbBillingAddressInCSV.EnterpriseCompanyName {
					enterprise.BillingAddressID = duplicatedbBillingAddressInCSV.ID
					isDuplicatedBillingAddressInCSV = true
					break
				}
			}
		}

		// CSV内で被らなかった場合、DBとの重複判定
		if !isDuplicatedBillingAddressInCSV {

			// DBとの重複判定
			isDuplicatedBillingAddressInDB = false

			for _, billingAddressInDB := range billingAddressListInDB {
				// 請求先社名と企業名が一致した場合は飛ばす
				if enterprise.BillingAddressCompanyName == billingAddressInDB.CompanyName &&
					enterprise.CompanyName == billingAddressInDB.EnterpriseCompanyName {
					enterprise.BillingAddressID = billingAddressInDB.ID
					isDuplicatedBillingAddressInDB = true
					break
				}
			}

			if !isDuplicatedBillingAddressInDB {

				billingAddress := entity.NewBillingAddress(
					enterprise.EnterpriseID,
					enterprise.AgentStaffIDForBillingAddress,
					enterprise.ContractPhase,
					enterprise.ContractDate,
					enterprise.PaymentPolicy,
					enterprise.BillingAddressCompanyName,
					enterprise.BillingAddressAddress,
					enterprise.HowToRecommend,
					enterprise.BillingAddressTitle,
				)

				err = i.billingAddressRepository.Create(billingAddress)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				// 求人登録時に使用するために請求先IDを保存する
				enterprise.BillingAddressID = billingAddress.ID

			}

			// csvで被らなかった場合、IDを取得してから被りチェックリストにもアペンドする
			duplicatedBillingAddressCheckerList = append(duplicatedBillingAddressCheckerList, entity.BillingAddress{
				ID:                    enterprise.BillingAddressID,
				CompanyName:           enterprise.BillingAddressCompanyName,
				ContractDate:          enterprise.ContractDate,
				PaymentPolicy:         enterprise.PaymentPolicy,
				HowToRecommend:        enterprise.HowToRecommend,
				AgentStaffID:          enterprise.AgentStaffIDForBillingAddress,
				HRStaffs:              enterprise.HRStaffs,
				RAStaffs:              enterprise.RAStaffs,
				EnterpriseCompanyName: enterprise.CompanyName,
			})
		}

		/********* 求人情報 *********/
		// EcternalIDとExternalTypeが一致する求人は、既に登録されているためIsGuaranteedInterviewをtrueにする
		for _, jobInformationInDB := range jobInformationListInDB {
			if jobInformationInDB.ExternalID == enterprise.ExternalID &&
				jobInformationInDB.ExternalType == enterprise.ExternalType {
				fmt.Println("すでに登録されている求人のため、面接確約のみを更新します。", jobInformationInDB.ID, jobInformationInDB.Title, jobInformationInDB.ExternalID)
				err = i.jobInformationRepository.UpdateIsGuaranteedInterview(jobInformationInDB.ID, true)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
				continue enterpriseLoop
			}
		}

		jobInformation := entity.NewJobInformation(
			enterprise.BillingAddressID,
			enterprise.Title,
			enterprise.RecruitmentState,
			enterprise.ExpirationDate,
			enterprise.WorkDetail,
			enterprise.NumberOfHires,
			enterprise.WorkLocation,
			enterprise.Transfer,
			enterprise.TransferDetail,
			enterprise.UnderIncome,
			enterprise.OverIncome,
			enterprise.Salary,
			enterprise.Insurance,
			enterprise.WorkTime,
			enterprise.OvertimeAverage,
			enterprise.FixedOvertimePayment,
			enterprise.FixedOvertimeDetail,
			enterprise.TrialPeriod,
			enterprise.TrialPeriodDetail,
			enterprise.EmploymentPeriod,
			enterprise.EmploymentPeriodDetail,
			enterprise.HolidayType,
			enterprise.HolidayDetail,
			enterprise.PassiveSmoking,
			enterprise.SelectionFlow,
			enterprise.Gender,
			enterprise.Nationality,
			enterprise.FinalEducation,
			enterprise.SchoolLevel,
			enterprise.MedicalHistory,
			enterprise.AgeUnder,
			enterprise.AgeOver,
			enterprise.JobChange,
			enterprise.ShortResignation,
			enterprise.ShortResignationRemarks,
			enterprise.SocialExperienceYear,
			enterprise.SocialExperienceMonth,
			enterprise.Appearance,
			enterprise.Communication,
			enterprise.Thinking,
			enterprise.TargetDetail,
			enterprise.Commission,
			enterprise.CommissionRate,
			enterprise.CommissionDetail,
			enterprise.RefundPolicy,
			enterprise.RequiredExperienceJobDetail,
			enterprise.SecretMemo,
			enterprise.RequiredDocumentsDetail,
			enterprise.EmploymentInsurance,
			enterprise.AccidentInsurance,
			enterprise.HealthInsurance,
			enterprise.PensionInsurance,
			enterprise.RegisterPhase,
			enterprise.StudyCategory,
			enterprise.DriverLicence,
			enterprise.WordSkill,
			enterprise.ExcelSkill,
			enterprise.PowerPointSkill,
			false, // IsExternal import時はautoscout求人のためfalse
			enterprise.WorkDetailAfterHiring,
			enterprise.WorkDetailScopeOfChange,
			enterprise.OfferRate,
			enterprise.DocumentPassingRate,
			enterprise.NumberOfRecentApplications,
			// 面接確約がある場合のみ取得する今回の場合は全てtrue
			true, // enterprise.IsGuaranteedInterview,
		)

		err = i.jobInformationRepository.Create(jobInformation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 他媒体求人IDの登録
		externalJobID := entity.NewJobInformationExternalID(
			jobInformation.ID,
			enterprise.ExternalType,
			enterprise.ExternalID,
		)

		err = i.jobInformationExternalIDRepository.Create(externalJobID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, prefecture := range enterprise.Prefectures {
			p := entity.NewJobInformationPrefecture(
				jobInformation.ID,
				prefecture.Prefecture,
			)

			err = i.jobInfoPrefectureRepository.Create(p)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, employmentStatus := range enterprise.EmploymentStatuses {
			es := entity.NewJobInformationEmploymentStatus(
				jobInformation.ID,
				employmentStatus.EmploymentStatus,
			)

			err = i.jobInfoEmploymentStatusRepository.Create(es)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 募集職種
		for _, occupation := range enterprise.Occupations {
			oc := entity.NewJobInformationOccupation(
				jobInformation.ID,
				occupation.Occupation,
			)

			err = i.jobInfoOccupationRepository.Create(oc)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, sf := range enterprise.SelectionFlowPatterns {

			// 選考フローパターンの作成
			selectionFlowPattern := entity.NewJobInformationSelectionFlowPattern(
				jobInformation.ID,
				sf.PublicStatus,
				sf.FlowTitle,
				sf.FlowPattern,
				false,
			)

			err = i.jobInfoSelectionFlowPatternRepository.Create(selectionFlowPattern)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// 選考情報を作成
			for _, si := range sf.SelectionInformations {
				selectionInformation := entity.NewJobInformationSelectionInformation(
					selectionFlowPattern.ID,
					si.SelectionType,
					si.SelectionPoint,
					si.PassedExample,
					si.FailExample,
					si.PassingRate,
					si.IsQuestionnairy,
				)

				err = i.jobInfoSelectionInformationRepository.Create(selectionInformation)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}

		for _, target := range enterprise.Targets {
			t := entity.NewJobInformationTarget(
				jobInformation.ID,
				target.Target,
			)

			err = i.jobInfoTargetRepository.Create(t)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, feature := range enterprise.Features {
			f := entity.NewJobInformationFeature(
				jobInformation.ID,
				feature.Feature,
			)

			err = i.jobInfoFeatureRepository.Create(f)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, workCharmPoint := range enterprise.WorkCharmPoints {
			wcp := entity.NewJobInformationWorkCharmPoint(
				jobInformation.ID,
				workCharmPoint.Title,
				workCharmPoint.Contents,
			)

			err = i.jobInfoWorkCharmPointRepository.Create(wcp)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 必要条件
		for _, requiredCondition := range enterprise.RequiredConditions {
			rc := entity.NewJobInformationRequiredCondition(
				jobInformation.ID,
				requiredCondition.IsCommon,
				requiredCondition.RequiredManagement,
			)

			err = i.jobInfoRequiredConditionRepository.Create(rc)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			requiredCondition.ID = rc.ID

			// 必要資格　複数
			for _, requiredLicense := range requiredCondition.RequiredLicenses {
				rl := entity.NewJobInformationRequiredLicense(
					requiredCondition.ID,
					requiredLicense.License,
				)

				err = i.jobInfoRequiredLicenseRepository.Create(rl)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

			// 必要PCツール　単数
			for _, requiredPCTool := range requiredCondition.RequiredPCTools {
				rpt := entity.NewJobInformationRequiredPCTool(
					requiredCondition.ID,
					requiredPCTool.Tool,
				)

				err = i.jobInfoRequiredPCToolRepository.Create(rpt)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

			// 必要言語スキル　単数
			requiredLanguages := entity.NewJobInformationRequiredLanguage(
				requiredCondition.ID,
				requiredCondition.RequiredLanguages.LanguageLevel,
				requiredCondition.RequiredLanguages.Toeic,
				requiredCondition.RequiredLanguages.ToeflIBT,
				requiredCondition.RequiredLanguages.ToeflPBT,
			)

			err = i.jobInfoRequiredLanguageRepository.Create(requiredLanguages)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, languageType := range requiredCondition.RequiredLanguages.LanguageTypes {
				lt := entity.NewJobInformationRequiredLanguageType(
					requiredLanguages.ID,
					languageType.LanguageType,
				)

				err = i.jobInfoRequiredLanguageTypeRepository.Create(lt)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

			// 開発スキル　　言語,OS 各1つずつ
			for _, requiredExperienceDevelopment := range requiredCondition.RequiredExperienceDevelopments {
				red := entity.NewJobInformationRequiredExperienceDevelopment(
					requiredCondition.ID,
					requiredExperienceDevelopment.DevelopmentCategory,
					requiredExperienceDevelopment.ExperienceYear,
					requiredExperienceDevelopment.ExperienceMonth,
				)

				err = i.jobInfoRequiredExperienceDevelopmentRepository.Create(red)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				for _, experienceDevelopmentType := range requiredExperienceDevelopment.ExperienceDevelopmentTypes {
					dt := entity.NewJobInformationRequiredExperienceDevelopmentType(
						red.ID,
						experienceDevelopmentType.DevelopmentType,
					)

					err = i.jobInfoRequiredExperienceDevelopmentTypeRepository.Create(dt)
					if err != nil {
						fmt.Println(err)
						return output, err
					}
				}
			}

			// 必要経験　　単数
			requiredExperienceJobs := entity.NewJobInformationRequiredExperienceJob(
				requiredCondition.ID,
				requiredCondition.RequiredExperienceJobs.ExperienceYear,
				requiredCondition.RequiredExperienceJobs.ExperienceMonth,
			)

			err = i.jobInfoRequiredExperienceJobRepository.Create(requiredExperienceJobs)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, experienceIndustry := range requiredCondition.RequiredExperienceJobs.ExperienceIndustries {
				ei := entity.NewJobInformationRequiredExperienceIndustry(
					requiredExperienceJobs.ID,
					experienceIndustry.ExperienceIndustry,
				)

				err = i.jobInfoRequiredExperienceIndustryRepository.Create(ei)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

			for _, experienceOccupation := range requiredCondition.RequiredExperienceJobs.ExperienceOccupations {
				eo := entity.NewJobInformationRequiredExperienceOccupation(
					requiredExperienceJobs.ID,
					experienceOccupation.ExperienceOccupation,
				)

				err = i.jobInfoRequiredExperienceOccupationRepository.Create(eo)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}

		for _, requiredSocialExperience := range enterprise.RequiredSocialExperiences {
			rse := entity.NewJobInformationRequiredSocialExperience(
				jobInformation.ID,
				requiredSocialExperience.SocialExperienceType,
			)

			err = i.jobInfoRequiredSocialExperienceRepository.Create(rse)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 非公開エージェント
		for _, hideToAgent := range enterprise.HideToAgents {
			hta := entity.NewJobInformationHideToAgent(
				jobInformation.ID,
				hideToAgent.AgentID,
			)

			err = i.jobInfoHideToAgentRepository.Create(hta)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	fmt.Println("企業情報の登録が完了しました。")
	fmt.Println("----------------------------------------")

	output.OK = true
	return output, nil
}

/*********************************************************************************/
// CSVファイルの出力
//
// 求人企業のcsvファイルを出力
type ExportEnterpriseCSVInput struct {
	AgentID uint
}

type ExportEnterpriseCSVOutput struct {
	FilePath *entity.FilePath
}

func (i *EnterpriseProfileInteractorImpl) ExportEnterpriseCSV(input ExportEnterpriseCSVInput) (ExportEnterpriseCSVOutput, error) {
	var (
		output  ExportEnterpriseCSVOutput
		err     error
		records [][]string
		record  []string
	)

	// csvの一行目を作成
	records = append(
		records,
		[]string{
			//企業情報
			"企業名",
			"企業担当者名",
			"企業HP URL",
			"代表者名",
			"本社所在地", // 郵便番号・住所
			"設立(年月)",
			"従業員数", // 単体・連結
			"資本金",
			"株式公開",
			"売上高", // 年度・額
			"事業内容",
			"業界",

			//請求先情報
			"RA担当者名",
			"契約フェーズ",
			"契約締結日",
			"支払い規定",
			"人事担当者",  // 名前・メールアドレス・電話番号
			"請求先担当者", // 名前・メールアドレス・電話番号
			"推薦方法",

			//求人情報
			//募集概要
			"求人タイトル",
			"募集対象",
			"雇用形態",
			"雇用期間の定め有無",
			"更新上限",
			"試用期間有無",
			"試用期間詳細",
			"募集状況",
			"募集期限",
			"募集職種",
			"採用人数",
			"求人の特徴",
			"求人の魅力",
			"仕事内容（雇入れ直後）",
			"仕事内容（変更の範囲）",
			"仕事内容",
			"勤務地（区分）",
			"勤務地（雇入れ直後）",
			"転勤有無",
			"変更の範囲",
			"年収", //下限・上限
			"給与詳細・昇給賞与",
			"雇用保険",
			"労災保険",
			"健康保険",
			"厚生年金保険",
			"諸手当・福利厚生",
			"固定残業代超過分の支払い有無",
			"固定残業代の詳細",
			"勤務時間",
			"平均残業時間",
			"休日・休暇タイプ",
			"休日・休暇詳細",
			"受動喫煙対策有無",
			"選考フロー",
			"選考フロー詳細",
			"応募資格(経験・スキルなど)",

			// 人物要件
			"性別",
			"国籍",
			"最終学歴",
			"大学ランク（大卒以上のみ選択）",
			"専攻大分類(理系/文系)",
			"募集年齢", //下限・上限
			"転職回数限度 (●社まで)",
			"短期離職（1年未満）",
			"短期離職備考",
			"必要社会人経験",
			"必要社会人経験年数",
			"必要Excelレベル",
			"必要Wordレベル",
			"必要PowerPointレベル",
			"アピアランス",
			"コミュニケーション",
			"論理的思考力",
			"応募条件（エージェント向け情報）",

			// 必須条件（業界・職種・資格・語学力）
			"共通の必須条件（業界・職種・資格・語学力）",
			"パターン別の必須条件（業界・職種・資格・語学力）",

			// 備考・手数料
			"社内限定メモ",
			"推薦時に必要な情報・書類",
			"成功報酬手数料(固定額)",
			"成功報酬手数料(料率)",
			"返金規定",
		},
	)

	//企業情報
	enterpriseProfileList, err := i.enterpriseProfileRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIndustryList, err := i.enterpriseIndustryRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, profile := range enterpriseProfileList {
		for _, industry := range enterpriseIndustryList {
			if profile.ID == industry.EnterpriseID {
				profile.Industries = append(profile.Industries, industry.Industry)
			}
		}

		billingAddressList, err := i.billingAddressRepository.GetByEnterpriseID(profile.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		hrStaff, err := i.billingAddressHRStaffRepository.GetByEnterpriseID(profile.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		raStaff, err := i.billingAddressRAStaffRepository.GetByEnterpriseID(profile.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, billingAddress := range billingAddressList {
			for _, hs := range hrStaff {
				if billingAddress.ID == hs.BillingAddressID {
					value := entity.BillingAddressHRStaff{
						BillingAddressID:   hs.BillingAddressID,
						HRStaffName:        hs.HRStaffName,
						HRStaffEmail:       hs.HRStaffEmail,
						HRStaffPhoneNumber: hs.HRStaffPhoneNumber,
					}

					billingAddress.HRStaffs = append(billingAddress.HRStaffs, value)
				}
			}

			for _, rs := range raStaff {
				if billingAddress.ID == rs.BillingAddressID {
					value := entity.BillingAddressRAStaff{
						BillingAddressID:               rs.BillingAddressID,
						BillingAddressStaffName:        rs.BillingAddressStaffName,
						BillingAddressStaffEmail:       rs.BillingAddressStaffEmail,
						BillingAddressStaffPhoneNumber: rs.BillingAddressStaffPhoneNumber,
					}

					billingAddress.RAStaffs = append(billingAddress.RAStaffs, value)
				}
			}

			jobInformationList, err := i.jobInformationRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			targets, err := i.jobInfoTargetRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			features, err := i.jobInfoFeatureRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			prefectures, err := i.jobInfoPrefectureRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredPCTools, err := i.jobInfoRequiredPCToolRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			occupations, err := i.jobInfoOccupationRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			industries, err := i.enterpriseIndustryRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			selectionInformations, err := i.jobInfoSelectionInformationRepository.GetByBillingAddressID(billingAddress.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// hideToAgents, err := i.jobInfoHideToAgentRepository.GetByBillingAddressID(billingAddress.ID)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return output, err
			// }

			for _, jobInformation := range jobInformationList {
				for _, t := range targets {
					if jobInformation.ID == t.JobInformationID {
						value := entity.JobInformationTarget{
							JobInformationID: t.JobInformationID,
							Target:           t.Target,
						}

						jobInformation.Targets = append(jobInformation.Targets, value)
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

				for _, wcp := range workCharmPoints {
					if jobInformation.ID == wcp.JobInformationID {
						value := entity.JobInformationWorkCharmPoint{
							JobInformationID: wcp.JobInformationID,
							Title:            wcp.Title,
							Contents:         wcp.Contents,
						}

						jobInformation.WorkCharmPoints = append(jobInformation.WorkCharmPoints, value)
					}
				}

				for _, es := range employmentStatuses {
					if jobInformation.ID == es.JobInformationID {
						value := entity.JobInformationEmploymentStatus{
							JobInformationID: es.JobInformationID,
							EmploymentStatus: es.EmploymentStatus,
						}

						jobInformation.EmploymentStatuses = append(jobInformation.EmploymentStatuses, value)
					}
				}

				for _, condition := range requiredConditions {
					if jobInformation.ID == condition.JobInformationID {
						for _, rl := range requiredLicenses {
							if condition.ID == rl.ConditionID {
								condition.RequiredLicenses = append(condition.RequiredLicenses, *rl)
							}
						}

						for _, rps := range requiredPCTools {
							if condition.ID == rps.ConditionID {
								condition.RequiredPCTools = append(condition.RequiredPCTools, *rps)
							}
						}

						for _, rl := range requiredLanguages {
							if condition.ID == rl.ConditionID {

								// 言語タイプ
								for _, languageType := range requiredLanguageTypes {
									if rl.ID == languageType.LanguageID {
										rl.LanguageTypes = append(rl.LanguageTypes, *languageType)
									}
								}

								condition.RequiredLanguages = *rl
							}
						}

						for _, red := range requiredExperienceDevelopments {
							if condition.ID == red.ConditionID {

								// 開発タイプ
								for _, developmentType := range requiredExperienceDevelopmentTypes {
									if red.ID == developmentType.ExperienceDevelopmentID {
										red.ExperienceDevelopmentTypes = append(red.ExperienceDevelopmentTypes, *developmentType)
									}
								}
								condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, *red)
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

				for _, rse := range requiredSocialExperiences {
					if jobInformation.ID == rse.JobInformationID {
						value := entity.JobInformationRequiredSocialExperience{
							JobInformationID:     rse.JobInformationID,
							SocialExperienceType: rse.SocialExperienceType,
						}
						jobInformation.RequiredSocialExperiences = append(jobInformation.RequiredSocialExperiences, value)
					}
				}

				for _, oc := range occupations {
					value := entity.JobInformationOccupation{
						JobInformationID: oc.JobInformationID,
						Occupation:       oc.Occupation,
					}
					jobInformation.Occupations = append(jobInformation.Occupations, value)
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

				for _, sfp := range selectionFlowPatterns {
					if jobInformation.ID == sfp.JobInformationID {
						value := entity.JobInformationSelectionFlowPattern{
							JobInformationID: sfp.JobInformationID,
							PublicStatus:     sfp.PublicStatus,
							FlowTitle:        sfp.FlowTitle,
							FlowPattern:      sfp.FlowPattern,
						}

						for _, si := range selectionInformations {
							if sfp.ID == si.SelectionFlowID {
								valueSi := entity.JobInformationSelectionInformation{
									SelectionFlowID: si.SelectionFlowID,
									SelectionType:   si.SelectionType,
									SelectionPoint:  si.SelectionPoint,
									PassedExample:   si.PassedExample,
									FailExample:     si.FailExample,
									PassingRate:     si.PassingRate,
								}
								value.SelectionInformations = append(value.SelectionInformations, valueSi)
							}
						}
						jobInformation.SelectionFlowPatterns = append(jobInformation.SelectionFlowPatterns, value)
					}
				}

				// requiredConditionをcommonConditionとpatternConditionに分ける
				for _, condition := range jobInformation.RequiredConditions {
					if condition.IsCommon {
						jobInformation.CommonCondition = condition
					} else {
						jobInformation.PatternConditions = append(jobInformation.PatternConditions, condition)
					}
				}

				// 文字列に変換用にスライスに変換
				commonConditionList := []entity.JobInformationRequiredCondition{
					jobInformation.CommonCondition,
				}

				// 求人ごとにデータを格納
				// 縦持ちカラムは1セルにまとめる
				record =
					[]string{
						profile.CompanyName,
						profile.StaffName,
						profile.CorporateSiteURL,
						profile.Representative,
						profile.PostCode + "\n" + " " + profile.OfficeLocation,
						profile.Establishment,
						fmt.Sprint("単体:", profile.EmployeeNumberSingle.Int64, "\n連結:", profile.EmployeeNumberGroup.Int64),
						profile.Capital,
						getStrPublicOffering(profile.PublicOffering),
						fmt.Sprint(profile.Earnings, "(", profile.EarningsYear.Int64, "年度)"),
						profile.BusinessDetail,
						getStrIndustryList(profile.Industries),

						// 請求先情報
						billingAddress.StaffName,
						getStrContractPhase(billingAddress.ContractPhase),
						billingAddress.ContractDate,
						billingAddress.PaymentPolicy,
						getStrHRStaffList(billingAddress.HRStaffs),
						getStrRAStaffList(billingAddress.RAStaffs),
						billingAddress.HowToRecommend,

						// 求人情報
						// 募集概要
						jobInformation.Title,
						getStrTargetList(jobInformation.Targets),
						getStrEmploymentStatusListForJobInfo(jobInformation.EmploymentStatuses),
						getStrAvailable(jobInformation.EmploymentPeriod),
						jobInformation.EmploymentPeriodDetail,
						getStrAvailable(jobInformation.TrialPeriod),
						jobInformation.TrialPeriodDetail,
						getStrRecruitmentState(jobInformation.RecruitmentState),
						jobInformation.ExpirationDate,
						getStrOccupationListForJobInfo(jobInformation.Occupations),
						fmt.Sprint(jobInformation.NumberOfHires.Int64),
						getStrFeatureList(jobInformation.Features),
						getStrWorkCharmPointList(jobInformation.WorkCharmPoints),
						jobInformation.WorkDetailAfterHiring,
						jobInformation.WorkDetailScopeOfChange,
						jobInformation.WorkDetail,
						getStrPrefectureList(jobInformation.Prefectures),
						jobInformation.WorkLocation,
						getStrAvailable(jobInformation.Transfer),
						jobInformation.TransferDetail,
						fmt.Sprint(jobInformation.UnderIncome.Int64, "万円〜", jobInformation.OverIncome.Int64, "万円"),
						jobInformation.Salary,
						getStrAvailableFromBool(jobInformation.EmploymentInsurance),
						getStrAvailableFromBool(jobInformation.AccidentInsurance),
						getStrAvailableFromBool(jobInformation.HealthInsurance),
						getStrAvailableFromBool(jobInformation.PensionInsurance),
						jobInformation.Insurance,
						getStrAvailable(jobInformation.FixedOvertimePayment),
						jobInformation.FixedOvertimeDetail,
						jobInformation.WorkTime,
						jobInformation.OvertimeAverage,
						getStrHolidayForJobInfo(jobInformation.HolidayType),
						jobInformation.HolidayDetail,
						getStrAvailable(jobInformation.PassiveSmoking),
						jobInformation.SelectionFlow,
						getStrSelectionFlowPattern(jobInformation.SelectionFlowPatterns),
						jobInformation.RequiredExperienceJobDetail,

						// 人物条件
						getStrGenderForJobInfo(jobInformation.Gender),
						getStrNationalityForJobInfo(jobInformation.Nationality),
						getStrFinalEducationForJobSeeker(jobInformation.FinalEducation),
						getStrCollegeRank(jobInformation.SchoolLevel),
						getStrStudyCategoryForJobInfo(jobInformation.StudyCategory),
						fmt.Sprint(jobInformation.AgeUnder.Int64, "歳〜", jobInformation.AgeUnder.Int64, "歳"),
						getStrJobChangeForJobSeeker(jobInformation.JobChange),
						getStrAvailable(jobInformation.ShortResignation),
						jobInformation.ShortResignationRemarks,
						getStrRequiredSocialExperienceList(jobInformation.RequiredSocialExperiences),
						fmt.Sprint(jobInformation.SocialExperienceYear.Int64, "年", jobInformation.SocialExperienceMonth.Int64, "ヶ月"),
						getStrExcelSkill(jobInformation.ExcelSkill),
						getStrWordSkill(jobInformation.WordSkill),
						getStrPowerPointSkill(jobInformation.PowerPointSkill),
						getStrAppearanceForJobInfo(jobInformation.Appearance),
						getStrCommunicationForJobInfo(jobInformation.Communication),
						getStrThinkingForJobInfo(jobInformation.Thinking),
						jobInformation.TargetDetail,

						// 必要条件
						// 必須条件
						getStrRequiredCondition(commonConditionList),
						// パターン別条件
						getStrRequiredCondition(jobInformation.PatternConditions),

						// 備考・手数料
						jobInformation.SecretMemo,
						jobInformation.RequiredDocumentsDetail,
						fmt.Sprint(jobInformation.Commission.Int64, "万円"),
						fmt.Sprint(jobInformation.CommissionRate.Int64, "%"),
						jobInformation.RefundPolicy,
					}

					// レコードを追加
				records = append(records, record)
			}
		}
	}

	//CSVファイルを作成
	filePath := ("./enterprise" + fmt.Sprint(utility.CreateUUID()) + ".csv")

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	cw := csv.NewWriter(file)

	defer file.Close()
	defer cw.Flush()

	cw.WriteAll(records)

	if err := cw.Error(); err != nil {
		fmt.Println("error writing csv:", err)
		return output, err
	}

	output.FilePath = entity.NewFilePath(filePath)

	return output, nil
}
