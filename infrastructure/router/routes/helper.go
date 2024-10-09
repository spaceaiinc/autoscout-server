package routes

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces/presenter"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/guregu/null.v4"
)

func GetFirebaseToken(c echo.Context) string {
	token := c.Request().Header.Get("FirebaseAuthorization")
	return token
}

func GetLineToken(c echo.Context) string {
	token := c.Request().Header.Get("LineAuthorization")
	return token
}

func bindAndValidate(c echo.Context, obj interface{}) (err error) {
	err = c.Bind(obj)
	if err != nil {
		return err
	}
	return validator.New().Struct(obj)
}

func renderJSON(c echo.Context, p presenter.Presenter) {
	err := c.JSON(p.StatusCode(), p.Data())
	if err != nil {
		fmt.Println(err)
	}
}

// ファイル送信用
func renderFile(c echo.Context, filePath string) error {
	err := c.File(filePath)
	if err != nil {
		return err
	}
	return nil
}

func ParseWebHook(c echo.Context) string {
	token := c.Request().Header.Get("FirebaseAuthorization")
	return token
}

/****************************************************************************************/
/// QueryParamの処理

// []string -> []null.Int
func parseQueryParams(strList []string) ([]null.Int, error) {
	var (
		nullIntList []null.Int
	)
	for _, str := range strList {
		intValue, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}

		nullIntList = append(nullIntList, null.NewInt(int64(intValue), true))
	}

	return nullIntList, nil
}

// string -> null.Int
func parseQueryParam(strValue string) (null.Int, error) {
	var (
		nullIntValue null.Int
	)
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return null.NewInt(0, false), err
	}

	nullIntValue = null.NewInt(int64(intValue), true)

	return nullIntValue, nil
}

func parseQueryParamUINT(strList []string) ([]uint, error) {
	var (
		uintList []uint
	)
	for _, str := range strList {
		intValue, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}

		uintList = append(uintList, uint(intValue))
	}

	return uintList, nil
}

func parseQueryParamBool(strValue string) (bool, error) {
	if strValue == "" {
		return false, nil
	}

	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		return false, err
	}

	return boolValue, nil
}

// 企業検索
func parseSearchEnterpriseQueryParams(c echo.Context) (entity.SearchEnterprise, error) {
	var (
		searchEnterpriseParamList entity.SearchEnterprise

		freeWordStr     = c.QueryParam("free_word")
		agentStaffIDStr = c.QueryParam("agent_staff_id")
		industryStrList = c.QueryParams()["industries[]"]
		// prefectureStrList   = c.QueryParams()["prefectures[]"]
		companyScaleStrList = c.QueryParams()["company_scale_types[]"]
	)

	//業界
	industries, err := parseQueryParams(industryStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchEnterpriseParamList, wrapped
	}

	// 所在地
	// prefectures, err := parseQueryParams(prefectureStrList)
	// if err != nil {
	// 	wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
	// 	renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
	// 	return searchEnterpriseParamList, wrapped
	// }

	// 企業規模
	companyScaleTypes, err := parseQueryParams(companyScaleStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchEnterpriseParamList, wrapped
	}

	searchEnterpriseParamList = *entity.NewSearchEnterprise(
		freeWordStr,
		agentStaffIDStr,
		industries,
		// prefectures,
		companyScaleTypes,
	)

	return searchEnterpriseParamList, nil
}

// 請求先検索
func parseSearchBillingAddressQueryParams(c echo.Context) (entity.SearchBillingAddress, error) {
	var (
		searchEnterpriseParamList entity.SearchBillingAddress
		freeWordStr               = c.QueryParam("free_word")
		agentStaffIDStr           = c.QueryParam("agent_staff_id")
	)

	searchEnterpriseParamList = *entity.NewSearchBillingAddress(
		freeWordStr,
		agentStaffIDStr,
	)

	return searchEnterpriseParamList, nil
}

// 売り上げ検索
func parseSearchSaleQueryParams(c echo.Context) (entity.SearchSale, error) {
	var (
		searchSaleParamList entity.SearchSale

		raStaffIDStr    = c.QueryParam("ra_staff_id")
		caStaffIDStr    = c.QueryParam("ca_staff_id")
		accuracyStrList = c.QueryParams()["accuracy_types[]"]
	)

	//ヨミ
	accuracies, err := parseQueryParams(accuracyStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchSaleParamList, wrapped
	}

	searchSaleParamList = *entity.NewSearchSale(
		raStaffIDStr,
		caStaffIDStr,
		accuracies,
	)

	return searchSaleParamList, nil
}

// ヨミの絞り込み検索（設定のヨミ一覧ページに使用）
func parseSearchAccuracyQueryParams(c echo.Context) (entity.SearchAccuracy, error) {
	var (
		searchSaleParamList entity.SearchAccuracy

		agentIDStr                = c.QueryParam("agent_id")
		jobSeekerFreeWordStr      = c.QueryParam("job_seeker_free_word")
		jobInformationFreeWordStr = c.QueryParam("job_information_free_word")
		contractSignedMonthStr    = c.QueryParam("contract_signed_month")
		billingMonthStr           = c.QueryParam("billing_month")
		raStaffIDStr              = c.QueryParam("ra_staff_id")
		caStaffIDStr              = c.QueryParam("ca_staff_id")
		accuracyStrList           = c.QueryParams()["accuracies[]"]
		pageNumberStr             = c.QueryParam("page_number")
	)

	// エージェントIDを変換
	agentIDInt, err := strconv.Atoi(agentIDStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchSaleParamList, wrapped
	}

	// ページ番号を変換
	pageNumberInt, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchSaleParamList, wrapped
	}

	// ヨミ
	accuracies, err := parseQueryParamUINT(accuracyStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchSaleParamList, wrapped
	}

	searchSaleParamList = *entity.NewSearchSearchAccuracy(
		uint(agentIDInt),
		jobSeekerFreeWordStr,
		jobInformationFreeWordStr,
		contractSignedMonthStr,
		billingMonthStr,
		raStaffIDStr,
		caStaffIDStr,
		accuracies,
		uint(pageNumberInt),
	)

	return searchSaleParamList, nil
}

// 求人検索
func parseSearchJobInformationQueryParams(c echo.Context) (entity.SearchJobInformation, error) {
	var (
		searchJobInformationParamList entity.SearchJobInformation

		freeWordStr       = c.QueryParam("free_word")
		agentStaffIDStr   = c.QueryParam("agent_staff_id")
		industryStrList   = c.QueryParams()["industries[]"]
		occupationStrList = c.QueryParams()["occupations[]"]
		employmentStrList = c.QueryParams()["employments[]"]
		prefectureStrList = c.QueryParams()["prefectures[]"]
		underIncomeStr    = c.QueryParam("under_income")
		overIncomeStr     = c.QueryParam("over_income")

		genderStrList           = c.QueryParams()["gender_types[]"]
		ageStr                  = c.QueryParam("age")
		finalEducationStrList   = c.QueryParams()["final_education_types[]"]
		schoolLevelStrList      = c.QueryParams()["school_level_types[]"]
		studyCategoryStrList    = c.QueryParams()["study_category_types[]"]
		nationalityStrList      = c.QueryParams()["nationality_types[]"]
		medicalHistoryStrList   = c.QueryParams()["medical_history_types[]"]
		jobChangeStrList        = c.QueryParams()["job_change_types[]"]
		shortResignationStrList = c.QueryParams()["short_resignation_types[]"]
		driverLicenceStrList    = c.QueryParams()["driver_licence_types[]"]
		appearanceStrList       = c.QueryParams()["appearance_types[]"]
		communicationStrList    = c.QueryParams()["communication_types[]"]
		thinkingStrList         = c.QueryParams()["thinking_types[]"]

		requiredExperienceIndustryStrList   = c.QueryParams()["required_experience_industries[]"]
		requiredExperienceOccupationStrList = c.QueryParams()["required_experience_occupations[]"]
		requiredSocialExperienceTypeStrList = c.QueryParams()["required_social_experience_types[]"]
		requiredSocialExperienceYearStr     = c.QueryParam("required_social_experience_year")
		requiredSocialExperienceMonthStr    = c.QueryParam("required_social_experience_month")
		requiredManagementStr               = c.QueryParam("required_management")
		requiredLicenseStrList              = c.QueryParams()["required_licenses[]"]
		requiredLanguageStrList             = c.QueryParams()["required_languages[]"]
		requiredLanguageLevelStrList        = c.QueryParams()["required_language_levels[]"]
		requiredExcelSkillStrList           = c.QueryParams()["required_excel_skills[]"]
		requiredWordSkillStrList            = c.QueryParams()["required_word_skills[]"]
		requiredPowerPointSkillStrList      = c.QueryParams()["required_power_point_skills[]"]
		requiredAnotherPCSkillStrList       = c.QueryParams()["required_another_pc_skills[]"]
		requiredDevelopmentLanguageStrList  = c.QueryParams()["required_development_languages[]"]
		requiredDevelopmentOSStrList        = c.QueryParams()["required_development_os[]"]

		transferStrList     = c.QueryParams()["transfer_types[]"]
		holidayTypeStrList  = c.QueryParams()["holiday_types[]"]
		companyScaleStrList = c.QueryParams()["company_scale_types[]"]
		featureStrList      = c.QueryParams()["features[]"]

		offerRateStrList                  = c.QueryParams()["offer_rate_types[]"]
		documentPassingRateStrList        = c.QueryParams()["document_passing_rate_types[]"]
		numberOfRecentApplicationsStrList = c.QueryParams()["number_of_recent_applications_types[]"]
		isGuaranteedInterviewStr          = c.QueryParam("is_guaranteed_interview")
	)

	//業界
	industries, err := parseQueryParams(industryStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 職種
	occupations, err := parseQueryParams(occupationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 雇用形態
	employments, err := parseQueryParams(employmentStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 所在地
	prefectures, err := parseQueryParams(prefectureStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 性別
	genderTypes, err := parseQueryParams(genderStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 最終学歴
	finalEducationTypes, err := parseQueryParams(finalEducationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 大学レベル
	schoolLevelTypes, err := parseQueryParams(schoolLevelStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 文系・理系
	studyCategoryTypes, err := parseQueryParams(studyCategoryStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 国籍
	nationalityTypes, err := parseQueryParams(nationalityStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 既往歴
	medicalHistoryTypes, err := parseQueryParams(medicalHistoryStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 転職回数
	jobChangeTypes, err := parseQueryParams(jobChangeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 短期離職
	shortResignationTypes, err := parseQueryParams(shortResignationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 普通自動車免許
	driverLicenceTypes, err := parseQueryParams(driverLicenceStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// アピアランス
	appearanceTypes, err := parseQueryParams(appearanceStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// コミュ力
	communicationTypes, err := parseQueryParams(communicationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 論理的思考力
	thinkingTypes, err := parseQueryParams(thinkingStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 経験業界
	experienceIndustries, err := parseQueryParams(requiredExperienceIndustryStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 経験職種
	experienceOccupations, err := parseQueryParams(requiredExperienceOccupationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 社会人経験
	socialExperienceTypes, err := parseQueryParams(requiredSocialExperienceTypeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 資格
	licenses, err := parseQueryParams(requiredLicenseStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 言語スキル
	languages, err := parseQueryParams(requiredLanguageStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 言語レベル
	languageLevels, err := parseQueryParams(requiredLanguageLevelStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// excelスキル
	// excelスキル
	excelSkills, err := parseQueryParams(requiredExcelSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// wordスキル
	wordSkills, err := parseQueryParams(requiredWordSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// powerpointスキル
	powerPointSkills, err := parseQueryParams(requiredPowerPointSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 言語スキル
	anotherPCSkills, err := parseQueryParams(requiredAnotherPCSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 開発言語
	developmentLanguages, err := parseQueryParams(requiredDevelopmentLanguageStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 開発OS
	developmentOS, err := parseQueryParams(requiredDevelopmentOSStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 転勤
	transferTypes, err := parseQueryParams(transferStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 休日タイプ
	holidayTypes, err := parseQueryParams(holidayTypeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 企業規模
	companyScaleTypes, err := parseQueryParams(companyScaleStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 求人の特徴
	features, err := parseQueryParams(featureStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 内定率
	offerRateTypes, err := parseQueryParams(offerRateStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 書類通過率
	documentPassingRateTypes, err := parseQueryParams(documentPassingRateStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 直近応募数
	numberOfRecentApplicationsTypes, err := parseQueryParams(numberOfRecentApplicationsStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	// 面接確約
	isGuaranteedInterview, err := parseQueryParamBool(isGuaranteedInterviewStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobInformationParamList, wrapped
	}

	searchJobInformationParamList = *entity.NewSearchJobInformation(
		freeWordStr,
		agentStaffIDStr,
		industries,
		occupations,
		employments,
		prefectures,
		underIncomeStr,
		overIncomeStr,
		genderTypes,
		ageStr,
		finalEducationTypes,
		schoolLevelTypes,
		studyCategoryTypes,
		nationalityTypes,
		medicalHistoryTypes,
		jobChangeTypes,
		shortResignationTypes,
		driverLicenceTypes,
		appearanceTypes,
		communicationTypes,
		thinkingTypes,
		experienceIndustries,
		experienceOccupations,
		socialExperienceTypes,
		requiredSocialExperienceYearStr,
		requiredSocialExperienceMonthStr,
		requiredManagementStr,
		licenses,
		languages,
		languageLevels,
		excelSkills,
		wordSkills,
		powerPointSkills,
		anotherPCSkills,
		developmentLanguages,
		developmentOS,
		transferTypes,
		holidayTypes,
		companyScaleTypes,
		features,
		offerRateTypes,
		documentPassingRateTypes,
		numberOfRecentApplicationsTypes,
		isGuaranteedInterview,
	)

	return searchJobInformationParamList, nil
}

// 求職者検索
func parseSearchJobSeekerQueryParams(c echo.Context) (entity.SearchJobSeeker, error) {
	var (
		searchJobSeekerParamList entity.SearchJobSeeker

		freeWordStr                = c.QueryParam("free_word")
		agentStaffIDStr            = c.QueryParam("agent_staff_id")
		phaseStrList               = c.QueryParams()["phase_types[]"]
		genderStrList              = c.QueryParams()["gender_types[]"]
		ageUnderStr                = c.QueryParam("age_under")
		ageOverStr                 = c.QueryParam("age_over")
		desiredIndustryStrList     = c.QueryParams()["desired_industries[]"]
		desiredOccupationStrList   = c.QueryParams()["desired_occupations[]"]
		desiredEorkLocationStrList = c.QueryParams()["desired_work_locations[]"]
		finalEducationStrList      = c.QueryParams()["final_education_types[]"]
		studyCategoryTypeStrList   = c.QueryParams()["study_category_types[]"]
		schoolLevelStrList         = c.QueryParams()["school_level_types[]"]
		nationalityStrList         = c.QueryParams()["nationality_types[]"]
		jobChangeStrList           = c.QueryParams()["job_change_types[]"]
		shortResignationStrList    = c.QueryParams()["short_resignation_types[]"]

		underIncomeStr                 = c.QueryParam("under_income")
		overIncomeStr                  = c.QueryParam("over_income")
		desiredTransferTypeStrList     = c.QueryParams()["desired_transfer_types[]"]
		desiredHolidayTypeStrList      = c.QueryParams()["desired_holiday_types[]"]
		desiredCompanyScaleTypeStrList = c.QueryParams()["desired_company_scale_types[]"]

		experienceIndustryStrList   = c.QueryParams()["experience_industries[]"]
		experienceOccupationStrList = c.QueryParams()["experience_occupations[]"]
		socialExperienceStrList     = c.QueryParams()["social_experiences[]"]
		managementStr               = c.QueryParam("management")
		licenseStrList              = c.QueryParams()["licenses[]"]
		languageStrList             = c.QueryParams()["languages[]"]
		excelSkillStrList           = c.QueryParams()["excel_skills[]"]
		wordSkillStrList            = c.QueryParams()["word_skills[]"]
		powerPointSkillStrList      = c.QueryParams()["power_point_skills[]"]
		anotherPCSkillStrList       = c.QueryParams()["another_pc_skills[]"]
		developmentLanguageStrList  = c.QueryParams()["development_languages[]"]
		developmentOSStrList        = c.QueryParams()["development_os[]"]

		appearanceStrList    = c.QueryParams()["appearance_types[]"]
		communicationStrList = c.QueryParams()["communication_types[]"]
		thinkingStrList      = c.QueryParams()["thinking_types[]"]
	)

	// フェーズ
	phaseTypes, err := parseQueryParams(phaseStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 性別
	genderTypes, err := parseQueryParams(genderStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	//業界
	desiredIndustries, err := parseQueryParams(desiredIndustryStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 職種
	desiredOccupations, err := parseQueryParams(desiredOccupationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 勤務地
	desiredWorkLocations, err := parseQueryParams(desiredEorkLocationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 最終学歴
	finalEducationTypes, err := parseQueryParams(finalEducationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 文系・理系
	studyCategoryTypes, err := parseQueryParams(studyCategoryTypeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 大学レベル
	schoolLevelTypes, err := parseQueryParams(schoolLevelStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 国籍
	nationalityTypes, err := parseQueryParams(nationalityStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 転職回数
	jobChangeTypes, err := parseQueryParams(jobChangeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 短期離職
	shortResignationTypes, err := parseQueryParams(shortResignationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 転勤
	desiredTransferTypes, err := parseQueryParams(desiredTransferTypeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 休日タイプ
	desiredHolidayTypes, err := parseQueryParams(desiredHolidayTypeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 企業規模
	desiredCompanyScaleTypes, err := parseQueryParams(desiredCompanyScaleTypeStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 経験業界
	experienceIndustries, err := parseQueryParams(experienceIndustryStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 経験職種
	experienceOccupations, err := parseQueryParams(experienceOccupationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 社会人経験
	socialExperiences, err := parseQueryParams(socialExperienceStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 資格
	licenses, err := parseQueryParams(licenseStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 言語スキル
	languages, err := parseQueryParams(languageStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// excelスキル
	excelSkills, err := parseQueryParams(excelSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// wordスキル
	wordSkills, err := parseQueryParams(wordSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// powerpointスキル
	powerPointSkills, err := parseQueryParams(powerPointSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 言語スキル
	anotherPCSkills, err := parseQueryParams(anotherPCSkillStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 開発言語
	developmentLanguages, err := parseQueryParams(developmentLanguageStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 開発OS
	developmentOS, err := parseQueryParams(developmentOSStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// アピアランス
	appearanceTypes, err := parseQueryParams(appearanceStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// コミュ力
	communicationTypes, err := parseQueryParams(communicationStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	// 論理的思考力
	thinkingTypes, err := parseQueryParams(thinkingStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchJobSeekerParamList, wrapped
	}

	searchJobSeekerParamList = *entity.NewSearchJobSeeker(
		freeWordStr,
		agentStaffIDStr,
		phaseTypes,
		genderTypes,
		ageUnderStr,
		ageOverStr,
		desiredIndustries,
		desiredOccupations,
		desiredWorkLocations,
		finalEducationTypes,
		studyCategoryTypes,
		schoolLevelTypes,
		nationalityTypes,
		jobChangeTypes,
		shortResignationTypes,
		underIncomeStr,
		overIncomeStr,
		desiredTransferTypes,
		desiredHolidayTypes,
		desiredCompanyScaleTypes,
		experienceIndustries,
		experienceOccupations,
		socialExperiences,
		managementStr,
		licenses,
		languages,
		excelSkills,
		wordSkills,
		powerPointSkills,
		anotherPCSkills,
		developmentLanguages,
		developmentOS,
		appearanceTypes,
		communicationTypes,
		thinkingTypes,
	)
	return searchJobSeekerParamList, nil
}

// 求職者検索
func parseSearchChatJobSeekerQueryParams(c echo.Context) (entity.SearchChatJobSeeker, error) {
	var (
		searchChatJobSeekerParamList entity.SearchChatJobSeeker

		freeWordStr     = c.QueryParam("free_word")
		agentStaffIDStr = c.QueryParam("agent_staff_id")
		phaseStrList    = c.QueryParams()["phase_types[]"]
	)

	// フェーズ
	phaseTypes, err := parseQueryParams(phaseStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchChatJobSeekerParamList, wrapped
	}

	searchChatJobSeekerParamList = *entity.NewSearchChatJobSeeker(
		freeWordStr,
		agentStaffIDStr,
		phaseTypes,
	)
	return searchChatJobSeekerParamList, nil
}

// 送客のチャット求職者検索
func parseSearchChatSendingJobSeekerQueryParams(c echo.Context) (entity.SearchChatSendingJobSeeker, error) {
	var (
		searchChatSendingJobSeekerParamList entity.SearchChatSendingJobSeeker

		freeWordStr     = c.QueryParam("free_word")
		agentStaffIDStr = c.QueryParam("agent_staff_id")
		phaseStrList    = c.QueryParams()["phase_types[]"]
	)

	// フェーズ
	phaseTypes, err := parseQueryParams(phaseStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchChatSendingJobSeekerParamList, wrapped
	}

	searchChatSendingJobSeekerParamList = *entity.NewSearchChatSendingJobSeeker(
		freeWordStr,
		agentStaffIDStr,
		phaseTypes,
	)
	return searchChatSendingJobSeekerParamList, nil
}

type IDList struct {
	IDList []uint
}

func NewIDList(
	idList []uint,
) *IDList {
	return &IDList{
		IDList: idList,
	}
}

// IDのリスト
func parseIDListQueryParams(c echo.Context) (IDList, error) {
	var (
		idParamList IDList

		idStrList = c.QueryParams()["id_list[]"]
	)

	// IDのリスト
	idList, err := parseQueryParamUINT(idStrList)
	if err != nil {
		return idParamList, err
	}

	idParamList = *NewIDList(idList)

	return idParamList, nil
}

// タスクグループ検索
func parseSearchTaskQueryParams(c echo.Context) (entity.SearchTask, error) {
	var (
		searchSearchTaskParamList entity.SearchTask

		jobSeekerFreeWordStr  = c.QueryParam("job_seeker_freeWord")
		enterpriseFreeWordStr = c.QueryParam("enterprise_freeWord")
		raStaffIDStr          = c.QueryParam("ra_staff_id")
		caStaffIDStr          = c.QueryParam("ca_staff_id")
		taskPhaseStrList      = c.QueryParams()["task_phase_types[]"]
	)

	// タスクフェーズ
	taskPhaseTypes, err := parseQueryParams(taskPhaseStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchSearchTaskParamList, wrapped
	}

	searchTaskGroupParamList := *entity.NewSearchTask(
		jobSeekerFreeWordStr,
		enterpriseFreeWordStr,
		raStaffIDStr,
		caStaffIDStr,
		taskPhaseTypes,
	)

	return searchTaskGroupParamList, nil
}

/****************************************************************************************/

func parseDesiredIndustriesSearchParam(c echo.Context) ([]null.Int, error) {
	var (
		it0Str = c.QueryParam("it_0")
		it1Str = c.QueryParam("it_1")
		it2Str = c.QueryParam("it_2")
		it3Str = c.QueryParam("it_3")
		it4Str = c.QueryParam("it_4")
		it5Str = c.QueryParam("it_5")

		manufacturer0Str = c.QueryParam("manufacturer_0")
		manufacturer1Str = c.QueryParam("manufacturer_1")
		manufacturer2Str = c.QueryParam("manufacturer_2")
		manufacturer3Str = c.QueryParam("manufacturer_3")
		manufacturer4Str = c.QueryParam("manufacturer_4")
		manufacturer5Str = c.QueryParam("manufacturer_5")
		manufacturer6Str = c.QueryParam("manufacturer_6")
		manufacturer7Str = c.QueryParam("manufacturer_7")
		manufacturer8Str = c.QueryParam("manufacturer_8")
		manufacturer9Str = c.QueryParam("manufacturer_9")

		trading0Str = c.QueryParam("trading_0")
		trading1Str = c.QueryParam("trading_1")
		trading2Str = c.QueryParam("trading_2")
		trading3Str = c.QueryParam("trading_3")
		trading4Str = c.QueryParam("trading_4")
		trading5Str = c.QueryParam("trading_5")
		trading6Str = c.QueryParam("trading_6")
		trading7Str = c.QueryParam("trading_7")
		trading8Str = c.QueryParam("trading_8")
		trading9Str = c.QueryParam("trading_9")

		service0Str  = c.QueryParam("service_0")
		service1Str  = c.QueryParam("service_1")
		service2Str  = c.QueryParam("service_2")
		service3Str  = c.QueryParam("service_3")
		service4Str  = c.QueryParam("service_4")
		service5Str  = c.QueryParam("service_5")
		service6Str  = c.QueryParam("service_6")
		service7Str  = c.QueryParam("service_7")
		service8Str  = c.QueryParam("service_8")
		service9Str  = c.QueryParam("service_9")
		service10Str = c.QueryParam("service_10")
		service11Str = c.QueryParam("service_11")

		advertisement0Str = c.QueryParam("advertisement_0")

		consulting0Str = c.QueryParam("consulting_0")

		finance0Str = c.QueryParam("finance_0")
		finance1Str = c.QueryParam("finance_1")
		finance2Str = c.QueryParam("finance_2")
		finance3Str = c.QueryParam("finance_3")

		medical0Str = c.QueryParam("medical_0")
		medical1Str = c.QueryParam("medical_1")

		distribution0Str = c.QueryParam("distribution_0")
		distribution1Str = c.QueryParam("distribution_1")

		other0Str = c.QueryParam("other_0")
		other1Str = c.QueryParam("other_1")
		other2Str = c.QueryParam("other_2")
		other3Str = c.QueryParam("other_3")

		queryList []string

		industries []null.Int
	)

	queryList = append(
		queryList,
		it0Str, it1Str, it2Str, it3Str, it4Str, it5Str,
		manufacturer0Str, manufacturer1Str, manufacturer2Str, manufacturer3Str, manufacturer4Str, manufacturer5Str, manufacturer6Str, manufacturer7Str, manufacturer8Str, manufacturer9Str,
		trading0Str, trading1Str, trading2Str, trading3Str, trading4Str, trading5Str, trading6Str, trading7Str, trading8Str, trading9Str,
		service0Str, service1Str, service2Str, service3Str, service4Str, service5Str, service6Str, service7Str, service8Str, service9Str, service10Str, service11Str,
		advertisement0Str,
		consulting0Str,
		finance0Str, finance1Str, finance2Str, finance3Str,
		medical0Str, medical1Str,
		distribution0Str, distribution1Str,
		other0Str, other1Str, other2Str, other3Str,
	)

	for _, q := range queryList {
		if q == "" {
			industries = append(industries, null.NewInt(0, false))
		}

		if q != "" {
			e, err := strconv.Atoi(q)
			if err != nil {
				wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
				renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
				return nil, wrapped
			}
			industries = append(industries, null.NewInt(int64(e), true))
		}
	}

	return industries, nil
}

// ダッシュボード
func parseSearchDashboardQueryParams(c echo.Context) (entity.SearchDashboard, error) {
	var (
		searchDashboardParam entity.SearchDashboard

		agentIDStr      = c.QueryParam("agent_id")
		agentStaffIDStr = c.QueryParam("agent_staff_id")
		managementIDStr = c.QueryParam("management_id")
		periodStr       = c.QueryParam("period")
		targetStr       = c.QueryParam("target")
		rangeStr        = c.QueryParam("range")
		accuracyStrList = c.QueryParams()["accuracy_search_list[]"]
	)

	agentIDInt, err := strconv.Atoi(agentIDStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchDashboardParam, wrapped
	}

	agentStaffIDInt, err := strconv.Atoi(agentStaffIDStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchDashboardParam, wrapped
	}

	managementIDInt, err := strconv.Atoi(managementIDStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchDashboardParam, wrapped
	}

	periodInt, err := strconv.Atoi(periodStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchDashboardParam, wrapped
	}

	targetInt, err := strconv.Atoi(targetStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchDashboardParam, wrapped
	}

	rangeInt, err := strconv.Atoi(rangeStr)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchDashboardParam, wrapped
	}

	// ヨミ
	accuracyTypes, err := parseQueryParams(accuracyStrList)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
		renderJSON(c, presenter.NewErrorJSONPresenter(wrapped))
		return searchDashboardParam, wrapped
	}

	searchDashboardParam = *entity.NewSearchDashboard(
		uint(agentIDInt),
		uint(agentStaffIDInt),
		uint(managementIDInt),
		uint(periodInt),
		uint(targetInt),
		uint(rangeInt),
		accuracyTypes,
	)

	return searchDashboardParam, nil
}
