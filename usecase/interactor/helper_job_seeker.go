package interactor

import (
	"fmt"
	"os"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

// 求人の最大ページ数を返す（本番実装までは1ページあたり5件）
func getJobSeekerListMaxPage(jobSeekerList []*entity.JobSeeker) uint {
	var maxPage = len(jobSeekerList) / 20

	if 0 < (len(jobSeekerList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページの求人一覧を返す（本番実装までは1ページあたり5件）
func getJobSeekerListWithPage(jobSeekerList []*entity.JobSeeker, page uint) []*entity.JobSeeker {
	var (
		perPage uint = 20
		listLen uint = uint(len(jobSeekerList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return jobSeekerList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.JobSeeker{}
	}

	if (listLen - first) <= perPage {
		return jobSeekerList[first:]
	}
	return jobSeekerList[first:last]
}

// 非公開エージェントに含まれているかどうかをチェックする
func checkJobSeekerByHideToAgent(jobSeekerList []*entity.JobSeeker, hideToAgentList []*entity.JobSeekerHideToAgent) []*entity.JobSeeker {
	var (
		outputList         []*entity.JobSeeker
		hideToAgentChecker = make(map[uint]bool)
	)

	// 非公開エージェントIDリストが空の場合は求人リストをそのまま返す
	if len(hideToAgentList) < 1 {
		return jobSeekerList
	}

	// hideToAgentで弾かれる求人IDを取得
	for _, hta := range hideToAgentList {
		hideToAgentChecker[hta.JobSeekerID] = true

	}

	// 検索元エージェントと一致する非公開エージェント情報を取得
	// 求人一覧から非公開エージェントに含まれているものを除外
	for _, jobSeeker := range jobSeekerList {
		if !hideToAgentChecker[jobSeeker.ID] {
			outputList = append(outputList, jobSeeker)
		}
	}

	return outputList
}

// 求人リストからIDリストを取得する
func getJobSeekerIDList(jobSeekerList []*entity.JobSeeker) []uint {
	var idListUint []uint

	if len(jobSeekerList) < 1 {
		return idListUint
	}

	for _, jobSeeker := range jobSeekerList {
		idListUint = append(idListUint, jobSeeker.ID)
	}

	return idListUint
}

// フリガナ + 電話番号が一致する場合は除外する
func excludeDuplicateJobSeeker(
	jobSeekerListBeforeDuplicate []*entity.JobSeeker,
	agentID uint,
) []*entity.JobSeeker {
	var (
		jobSeekerList                      []*entity.JobSeeker
		furiganaAndPhoneNumberCheckerOwn   = make(map[string]uint) // 自社のフリガナの重複
		furiganaAndPhoneNumberCheckerOther = make(map[string]uint) // 他社のフリガナの重複
	)

	// 自社と他社エージェントで同一社名の求人がある場合は他社の同一社名の求人を除外する
	if len(jobSeekerListBeforeDuplicate) > 0 {
		for _, jobSeeker := range jobSeekerListBeforeDuplicate {
			// この形式にする→「タナカタロウ_080-0000-0000」
			fp := jobSeeker.LastFurigana + jobSeeker.FirstFurigana + "_" + jobSeeker.PhoneNumber

			if jobSeeker.AgentID == agentID {
				furiganaAndPhoneNumberCheckerOwn[fp] = furiganaAndPhoneNumberCheckerOwn[fp] + 1
			} else {
				furiganaAndPhoneNumberCheckerOther[fp] = furiganaAndPhoneNumberCheckerOther[fp] + 1
			}
		}

		for _, jobSeeker := range jobSeekerListBeforeDuplicate {
			// この形式にする→「タナカタロウ_080-0000-0000」
			fp := jobSeeker.LastFurigana + jobSeeker.FirstFurigana + "_" + jobSeeker.PhoneNumber

			if jobSeeker.AgentID == agentID {
				// 自社求人は全て取得
				jobSeekerList = append(jobSeekerList, jobSeeker)
			} else if jobSeeker.AgentID != agentID && (furiganaAndPhoneNumberCheckerOwn[fp] == 0 && furiganaAndPhoneNumberCheckerOther[fp] > 0) {
				// シェア求人は社名被りないもののみ取得
				jobSeekerList = append(jobSeekerList, jobSeeker)
			}
		}
	}

	// 他社同士の重複チェック
	if len(jobSeekerList) > 0 {
		for _, jobSeeker := range jobSeekerList {
			// この形式にする→「タナカタロウ_080-0000-0000」
			fp := jobSeeker.LastFurigana + jobSeeker.FirstFurigana + "_" + jobSeeker.PhoneNumber

			if furiganaAndPhoneNumberCheckerOther[fp] > 1 {
				jobSeeker.IsDuplicate = true
			}
		}
	}

	return jobSeekerList
}

// 特別仕様: 本番環境のみ「2: 株式会社テスト」と「3: 株式会社Space AI（非公開求人管理用）」を除外して他社エージェントに非表示にするための関数
func excludeTestJobSeeker(
	jobSeekerList []*entity.JobSeeker,
	agentID uint,
) []*entity.JobSeeker {
	env := os.Getenv("APP_ENV")

	/**
	 * 条件
	 *
	 * ユーザーが「1: 株式会社Space AI」と「2: 株式会社テスト」と「3: 株式会社Space AI（非公開求人管理用）」以外で
	 * 求職者の担当エージェントが「2: 株式会社テスト」と「3: 株式会社Space AI（非公開求人管理用）」の場合
	 * スライスから除外する
	**/

	if env == "prd" {
		for i := 0; i < len(jobSeekerList); i++ {
			// ユーザーが「1: 株式会社Space AI」と「2: 株式会社テスト」と「3: 株式会社Space AI（非公開求人管理用）」以外の場合
			if agentID != 1 && agentID != 2 && agentID != 3 {
				// 求職者の担当エージェントが「2: 株式会社テスト」と「3: 株式会社Space AI（非公開求人管理用）」の場合
				if jobSeekerList[i].AgentID == 2 || jobSeekerList[i].AgentID == 3 {
					// 除外する
					jobSeekerList = append(jobSeekerList[:i], jobSeekerList[i+1:]...)
					i-- // スライスの要素が前にシフトされたため、現在のインデックスを調整する
				}
			}
		}
	}

	return jobSeekerList
}

// 特別仕様: 本番環境のみ「2: 株式会社テスト」と「3: 株式会社Space AI（非公開求人管理用）」を除外して他社エージェントに非表示にするための関数
func getJobSeekerChildTableData(
	jobSeeker *entity.JobSeeker,
	i *JobSeekerInteractorImpl,
) (*entity.JobSeeker, error) {
	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	selfPromotion, err := i.jobSeekerSelfPromotionRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	document, err := i.jobSeekerDocumentRepository.FindByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	pcSkill, err := i.jobSeekerPCToolRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	hideToAgent, err := i.jobSeekerHideToAgentRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return jobSeeker, err
	}

	for _, sh := range studentHistory {
		value := entity.JobSeekerStudentHistory{
			JobSeekerID:    sh.JobSeekerID,
			SchoolCategory: sh.SchoolCategory,
			SchoolName:     sh.SchoolName,
			SchoolLevel:    sh.SchoolLevel,
			Subject:        sh.Subject,
			EntranceYear:   sh.EntranceYear,
			FirstStatus:    sh.FirstStatus,
			GraduationYear: sh.GraduationYear,
			LastStatus:     sh.LastStatus,
		}

		jobSeeker.StudentHistories = append(jobSeeker.StudentHistories, value)
	}

	for _, wh := range workHistory {
		value := entity.JobSeekerWorkHistory{
			ID:                   wh.ID,
			JobSeekerID:          wh.JobSeekerID,
			CompanyName:          wh.CompanyName,
			EmployeeNumberSingle: wh.EmployeeNumberSingle,
			EmployeeNumberGroup:  wh.EmployeeNumberGroup,
			PublicOffering:       wh.PublicOffering,
			JoiningYear:          wh.JoiningYear,
			EmploymentStatus:     wh.EmploymentStatus,
			RetireReasonOfPublic: wh.RetireReasonOfPublic,
			RetireReasonOfTruth:  wh.RetireReasonOfTruth,
			RetireYear:           wh.RetireYear,
			FirstStatus:          wh.FirstStatus,
			LastStatus:           wh.LastStatus,
		}

		for _, ei := range experienceIndustry {
			if ei.WorkHistoryID == wh.ID {
				valueEI := entity.JobSeekerExperienceIndustry{
					ID:            ei.ID,
					WorkHistoryID: ei.WorkHistoryID,
					Industry:      ei.Industry,
				}

				value.ExperienceIndustries = append(value.ExperienceIndustries, valueEI)
			}
		}

		for _, dh := range departmentHistory {
			if dh.WorkHistoryID == wh.ID {
				valuedh := entity.JobSeekerDepartmentHistory{
					ID:               dh.ID,
					WorkHistoryID:    dh.WorkHistoryID,
					Department:       dh.Department,
					ManagementNumber: dh.ManagementNumber,
					ManagementDetail: dh.ManagementDetail,
					JobDescription:   dh.JobDescription,
					StartYear:        dh.StartYear,
					EndYear:          dh.EndYear,
				}

				for _, eo := range experienceOccupation {
					if valuedh.ID == eo.DepartmentHistoryID {
						valueEO := entity.JobSeekerExperienceOccupation{
							DepartmentHistoryID: eo.DepartmentHistoryID,
							Occupation:          eo.Occupation,
						}

						valuedh.ExperienceOccupations = append(valuedh.ExperienceOccupations, valueEO)
					}
				}

				value.DepartmentHistories = append(value.DepartmentHistories, valuedh)
			}
		}
		jobSeeker.WorkHistories = append(jobSeeker.WorkHistories, value)
	}

	for _, dcs := range desiredCompanyScale {
		value := entity.JobSeekerDesiredCompanyScale{
			JobSeekerID:         dcs.JobSeekerID,
			DesiredCompanyScale: dcs.DesiredCompanyScale,
		}

		jobSeeker.DesiredCompanyScales = append(jobSeeker.DesiredCompanyScales, value)
	}

	for _, l := range license {
		value := entity.JobSeekerLicense{
			JobSeekerID:     l.JobSeekerID,
			LicenseType:     l.LicenseType,
			AcquisitionTime: l.AcquisitionTime,
		}

		jobSeeker.Licenses = append(jobSeeker.Licenses, value)
	}

	for _, sp := range selfPromotion {
		value := entity.JobSeekerSelfPromotion{
			JobSeekerID: sp.JobSeekerID,
			Title:       sp.Title,
			Contents:    sp.Contents,
		}

		jobSeeker.SelfPromotions = append(jobSeeker.SelfPromotions, value)
	}

	valueDocument := entity.JobSeekerDocument{
		JobSeekerID:             document.JobSeekerID,
		ResumeOriginURL:         document.ResumeOriginURL,
		ResumePDFURL:            document.ResumePDFURL,
		CVOriginURL:             document.CVOriginURL,
		CVPDFURL:                document.CVPDFURL,
		RecommendationOriginURL: document.RecommendationOriginURL,
		RecommendationPDFURL:    document.RecommendationPDFURL,
		IDPhotoURL:              document.IDPhotoURL,
		OtherDocument1URL:       document.OtherDocument1URL,
		OtherDocument2URL:       document.OtherDocument2URL,
		OtherDocument3URL:       document.OtherDocument3URL,
	}

	jobSeeker.Documents = valueDocument

	for _, di := range desiredIndustry {
		value := entity.JobSeekerDesiredIndustry{
			JobSeekerID:     di.JobSeekerID,
			DesiredIndustry: di.DesiredIndustry,
			DesiredRank:     di.DesiredRank,
		}

		jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, value)
	}

	for _, do := range desiredOccupation {
		value := entity.JobSeekerDesiredOccupation{
			JobSeekerID:       do.JobSeekerID,
			DesiredOccupation: do.DesiredOccupation,
			DesiredRank:       do.DesiredRank,
		}

		jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, value)
	}

	for _, dwl := range desiredWorkLocation {
		value := entity.JobSeekerDesiredWorkLocation{
			JobSeekerID:         dwl.JobSeekerID,
			DesiredWorkLocation: dwl.DesiredWorkLocation,
			DesiredRank:         dwl.DesiredRank,
		}

		jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, value)
	}

	for _, dht := range desiredHolidayType {
		value := entity.JobSeekerDesiredHolidayType{
			JobSeekerID: dht.JobSeekerID,
			HolidayType: dht.HolidayType,
		}

		jobSeeker.DesiredHolidayTypes = append(jobSeeker.DesiredHolidayTypes, value)
	}

	for _, ds := range developmentSkill {
		value := entity.JobSeekerDevelopmentSkill{
			JobSeekerID:         ds.JobSeekerID,
			DevelopmentCategory: ds.DevelopmentCategory,
			DevelopmentType:     ds.DevelopmentType,
			ExperienceYear:      ds.ExperienceYear,
			ExperienceMonth:     ds.ExperienceMonth,
		}

		jobSeeker.DevelopmentSkills = append(jobSeeker.DevelopmentSkills, value)
	}

	for _, ls := range languageSkill {
		value := entity.JobSeekerLanguageSkill{
			JobSeekerID:             ls.JobSeekerID,
			LanguageType:            ls.LanguageType,
			LanguageLevel:           ls.LanguageLevel,
			Toeic:                   ls.Toeic,
			ToeicExaminationYear:    ls.ToeicExaminationYear,
			ToeflIBT:                ls.ToeflIBT,
			ToeflIBTExaminationYear: ls.ToeflIBTExaminationYear,
			ToeflPBT:                ls.ToeflPBT,
			ToeflPBTExaminationYear: ls.ToeflPBTExaminationYear,
		}

		jobSeeker.LanguageSkills = append(jobSeeker.LanguageSkills, value)
	}

	for _, ps := range pcSkill {
		value := entity.JobSeekerPCTool{
			JobSeekerID: ps.JobSeekerID,
			Tool:        ps.Tool,
		}

		jobSeeker.PCTools = append(jobSeeker.PCTools, value)
	}

	for _, hta := range hideToAgent {
		value := entity.JobSeekerHideToAgent{
			JobSeekerID: hta.JobSeekerID,
			AgentID:     hta.AgentID,
			AgentName:   hta.AgentName,
		}

		jobSeeker.HideToAgents = append(jobSeeker.HideToAgents, value)
	}

	return jobSeeker, nil
}
