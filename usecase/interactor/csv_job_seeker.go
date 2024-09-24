package interactor

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"gopkg.in/guregu/null.v4"
)

/****************************************************************************************/
/// CSV操作 API
//
//求職者のcsvファイルを読み込む
type ImportJobSeekerCSVInput struct {
	CreateParam   []*entity.JobSeeker
	MissedRecords []uint
	AgentID       uint
}

type ImportJobSeekerCSVOutput struct {
	MissedRecords []uint
	OK            bool
}

func (i *JobSeekerInteractorImpl) ImportJobSeekerCSV(input ImportJobSeekerCSVInput) (ImportJobSeekerCSVOutput, error) {
	var (
		output ImportJobSeekerCSVOutput
		err    error
	)

	// routesで除外されたレコードを格納する
	output.MissedRecords = input.MissedRecords

	inflowChannelOptionList, err := i.agentInflowChannelOptionRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, jobSeeker := range input.CreateParam {
		jobSeeker.AgentID = input.AgentID              // エージェントID
		jobSeeker.RegisterPhase = null.NewInt(0, true) // 求職者の登録状況（0: 本登録, 1: 仮登録）

		// 流入マスタIDがエージェント内に存在するか確認
		isMatched := false
		for _, inflowChannelOption := range inflowChannelOptionList {
			isMatched = false
			if uint(jobSeeker.InflowChannelID.Int64) == inflowChannelOption.ID {
				isMatched = true
				break
			}
		}
		// 流入マスタIDがエージェント内に存在しない場合はnullにする
		if !isMatched {
			jobSeeker.InflowChannelID = null.NewInt(0, false)
		}

		err = i.jobSeekerRepository.Create(jobSeeker)
		if err != nil {
			fmt.Println(err)
			// output.MissedRecords = append(output.MissedRecords, jobSeeker.RecordLine)
			return output, err
		}

		// created_atを更新する処理 *csvインポートでエントリー日が入力されている場合に使用
		fmt.Println("created_atを更新する処理", jobSeeker.CreatedAt)
		err = i.jobSeekerRepository.UpdateCreatedAt(jobSeeker.ID, jobSeeker.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 開発環境の場合はユーザー名などを統一する
		if os.Getenv("APP_ENV") != "prd" {
			err = i.jobSeekerRepository.UpdateForDev(jobSeeker.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 空の書類テーブルの情報を登録
		err := i.jobSeekerDocumentRepository.Create(entity.NewJobSeekerDocument(
			jobSeeker.ID,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者作成時にメッセージグループ作成
		// エージェントと求職者のチャットグループを作成
		chatGroup := entity.NewChatGroupWithJobSeeker(
			jobSeeker.AgentID,
			jobSeeker.ID,
			false, // 初めはLINE連携してないから false
		)

		err = i.chatGroupWithJobSeekerRepository.Create(chatGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		var (
			EntryInterview           = null.NewInt(int64(entity.EntryInterview), true)
			InvitationInterview      = null.NewInt(int64(entity.InvitationInterview), true)
			ReservationInterview     = null.NewInt(int64(entity.ReservationInterview), true)
			WaitingInterview         = null.NewInt(int64(entity.WaitingInterview), true)
			PreparingAfterInterview  = null.NewInt(int64(entity.PreparingAfterInterview), true)
			OperatingAfterInterview  = null.NewInt(int64(entity.OperatingAfterInterview), true)
			ReleasedAfterInterview   = null.NewInt(int64(entity.ReleasedAfterInterview), true)
			OfferedAfterInterview    = null.NewInt(int64(entity.OfferedAfterInterview), true)
			ContinuingAfterInterview = null.NewInt(int64(entity.ContinuingAfterInterview), true)
			QuitedAfterInterview     = null.NewInt(int64(entity.QuitedAfterInterview), true)
		)

		var (
			phaseSub           null.Int
			date               string
			firstInterviewDate time.Time = utility.EarliestTime()
		)

		// 入寮したフェーズに応じて作成タスクを変更
		switch jobSeeker.Phase {
		case EntryInterview: // エントリー
			phaseSub = null.NewInt(0, true) //日程調整依頼
			// 当日（2020-02-01）
			date = jobSeeker.InterviewDate.Format("2006-01-02")

		case InvitationInterview: // 面談案内済み
			phaseSub = null.NewInt(0, true) //面談調整中
			// 当日（2020-02-01）
			date = jobSeeker.InterviewDate.Format("2006-01-02")

		case ReservationInterview: // 面談予約完了
			phaseSub = null.NewInt(0, true) //面談の前日確認を行う
			// jobSeeker.InterviewDate（2022-12-15T16:29）の前日（2022-12-14）
			yesterday := jobSeeker.InterviewDate.AddDate(0, 0, -1)
			date = yesterday.Format("2006-01-02")

		case WaitingInterview: // 面談実施待ち
			phaseSub = null.NewInt(1, true) //面談実施待ち
			date = jobSeeker.InterviewDate.Format("2006-01-02")

		case PreparingAfterInterview: // 面談実施済み（準備中）
			phaseSub = null.NewInt(99, true) //終了
			date = jobSeeker.InterviewDate.Format("2006-01-02")

			// 面談実施済みの場合は初回面談日時に記録する
			firstInterviewDate = jobSeeker.InterviewDate

		case OperatingAfterInterview: // 面談実施済み（稼働中）
			phaseSub = null.NewInt(99, true) //終了
			date = jobSeeker.InterviewDate.Format("2006-01-02")

			// 面談実施済みの場合は初回面談日時に記録する
			firstInterviewDate = jobSeeker.InterviewDate

		case ReleasedAfterInterview: // 面談実施済み（リリース状態）
			phaseSub = null.NewInt(99, true) //終了
			date = jobSeeker.InterviewDate.Format("2006-01-02")

			// 面談実施済みの場合は初回面談日時に記録する
			firstInterviewDate = jobSeeker.InterviewDate

		case OfferedAfterInterview: // サービス終了/決定者
			phaseSub = null.NewInt(99, true) //終了
			date = jobSeeker.InterviewDate.Format("2006-01-02")

		case ContinuingAfterInterview: // サービス終了/今後継続連絡
			phaseSub = null.NewInt(99, true) //終了
			date = jobSeeker.InterviewDate.Format("2006-01-02")

		case QuitedAfterInterview: // サービス終了/転職活動終了
			phaseSub = null.NewInt(99, true) //終了
			date = jobSeeker.InterviewDate.Format("2006-01-02")
		}

		// 面談調整タスクの作成
		interviewTaskGroup := entity.NewInterviewTaskGroup(
			input.AgentID,
			jobSeeker.ID,
			jobSeeker.InterviewDate,
			firstInterviewDate, // 面談実施済みの場合は初回面談日時に記録する
		)

		err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		interviewTask := entity.NewInterviewTask(
			interviewTaskGroup.ID,
			jobSeeker.AgentStaffID,
			jobSeeker.AgentStaffID,
			jobSeeker.Phase,
			phaseSub,
			"",
			date,
			null.NewInt(99, true),
			getStrPhaseForJobSeeker(jobSeeker.Phase),
		)

		err = i.interviewTaskRepository.Create(interviewTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, sh := range jobSeeker.StudentHistories {
			studentHistory := entity.NewJobSeekerStudentHistory(
				jobSeeker.ID,
				sh.SchoolCategory,
				sh.SchoolName,
				sh.SchoolLevel,
				sh.Subject,
				sh.EntranceYear,
				sh.FirstStatus,
				sh.GraduationYear,
				sh.LastStatus,
			)

			err = i.jobSeekerStudentHistoryRepository.Create(studentHistory)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		fmt.Println("jobSeeker.WorkHistories")
		fmt.Println(jobSeeker.LastName + jobSeeker.FirstName)
		fmt.Println(jobSeeker.WorkHistories)
		for _, wh := range jobSeeker.WorkHistories {
			workHistory := entity.NewJobSeekerWorkHistory(
				jobSeeker.ID,
				wh.CompanyName,
				wh.EmployeeNumberSingle,
				wh.EmployeeNumberGroup,
				wh.PublicOffering,
				wh.JoiningYear,
				wh.EmploymentStatus,
				wh.RetireReasonOfTruth,
				wh.RetireReasonOfPublic,
				wh.RetireYear,
				wh.FirstStatus,
				wh.LastStatus,
			)

			err = i.jobSeekerWorkHistoryRepository.Create(workHistory)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, ei := range wh.ExperienceIndustries {
				experienceIndustry := entity.NewJobSeekerExperienceIndustry(
					workHistory.ID,
					ei.Industry,
				)

				err = i.jobSeekerExperienceIndustryRepository.Create(experienceIndustry)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

			fmt.Println("wh.DepartmentHistories")
			fmt.Println(wh.DepartmentHistories)

			for _, dh := range wh.DepartmentHistories {

				departmentHistpry := entity.NewJobSeekerDepartmentHistory(
					workHistory.ID,
					dh.Department,
					dh.ManagementNumber,
					dh.ManagementDetail,
					dh.JobDescription,
					dh.StartYear,
					dh.EndYear,
				)

				err = i.jobSeekerDepartmentHistoryRepository.Create(departmentHistpry)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				for _, eo := range dh.ExperienceOccupations {
					experienceOccupation := entity.NewJobSeekerExperienceOccupation(
						departmentHistpry.ID,
						eo.Occupation,
					)

					err = i.jobSeekerExperienceOccupationRepository.Create(experienceOccupation)
					if err != nil {
						fmt.Println(err)
						return output, err
					}
				}
			}

		}

		for _, l := range jobSeeker.Licenses {
			license := entity.NewJobSeekerLicense(
				jobSeeker.ID,
				l.LicenseType,
				l.AcquisitionTime,
			)

			err = i.jobSeekerLicenseRepository.Create(license)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, sp := range jobSeeker.SelfPromotions {
			selfPromotion := entity.NewJobSeekerSelfPromotion(
				jobSeeker.ID,
				sp.Title,
				sp.Contents,
			)

			err = i.jobSeekerSelfPromotionRepository.Create(selfPromotion)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, di := range jobSeeker.DesiredIndustries {
			desiredIndustry := entity.NewJobSeekerDesiredIndustry(
				jobSeeker.ID,
				di.DesiredIndustry,
				di.DesiredRank,
			)

			err = i.jobSeekerDesiredIndustryRepository.Create(desiredIndustry)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, do := range jobSeeker.DesiredOccupations {
			desiredOccupation := entity.NewJobSeekerDesiredOccupation(
				jobSeeker.ID,
				do.DesiredOccupation,
				do.DesiredRank,
			)

			err = i.jobSeekerDesiredOccupationRepository.Create(desiredOccupation)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, dwl := range jobSeeker.DesiredWorkLocations {
			desiredWorkLocation := entity.NewJobSeekerDesiredWorkLocation(
				jobSeeker.ID,
				dwl.DesiredWorkLocation,
				dwl.DesiredRank,
			)

			err = i.jobSeekerDesiredWorkLocationRepository.Create(desiredWorkLocation)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, dht := range jobSeeker.DesiredHolidayTypes {
			desiredHolidayType := entity.NewJobSeekerDesiredHolidayType(
				jobSeeker.ID,
				dht.HolidayType,
			)

			err = i.jobSeekerDesiredHolidayTypeRepository.Create(desiredHolidayType)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, dcs := range jobSeeker.DesiredCompanyScales {
			desiredCompanyScale := entity.NewJobSeekerDesiredCompanyScale(
				jobSeeker.ID,
				dcs.DesiredCompanyScale,
			)

			err = i.jobSeekerDesiredCompanyScaleRepository.Create(desiredCompanyScale)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, ds := range jobSeeker.DevelopmentSkills {
			developmentSkill := entity.NewJobSeekerDevelopmentSkill(
				jobSeeker.ID,
				ds.DevelopmentCategory,
				ds.DevelopmentType,
				ds.ExperienceYear,
				ds.ExperienceMonth,
			)

			err = i.jobSeekerDevelopmentSkillRepository.Create(developmentSkill)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, ls := range jobSeeker.LanguageSkills {
			languageSkill := entity.NewJobSeekerLanguageSkill(
				jobSeeker.ID,
				ls.LanguageType,
				ls.LanguageLevel,
				ls.Toeic,
				ls.ToeicExaminationYear,
				ls.ToeflIBT,
				ls.ToeflIBTExaminationYear,
				ls.ToeflPBT,
				ls.ToeflPBTExaminationYear,
			)

			err = i.jobSeekerLanguageSkillRepository.Create(languageSkill)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, pt := range jobSeeker.PCTools {
			PCTool := entity.NewJobSeekerPCTool(
				jobSeeker.ID,
				pt.Tool,
			)

			err = i.jobSeekerPCToolRepository.Create(PCTool)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true
	return output, nil
}

// 求職者のcsvファイルを出力
type ExportJobSeekerCSVInput struct {
	AgentID uint
}

type ExportJobSeekerCSVOutput struct {
	FilePath *entity.FilePath
}

func (i *JobSeekerInteractorImpl) ExportJobSeekerCSV(input ExportJobSeekerCSVInput) (ExportJobSeekerCSVOutput, error) {
	var (
		output  ExportJobSeekerCSVOutput
		err     error
		records [][]string
		record  []string
	)

	jobSeekerList, err := i.jobSeekerRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selfPromotion, err := i.jobSeekerSelfPromotionRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	document, err := i.jobSeekerDocumentRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	pcSkill, err := i.jobSeekerPCToolRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgent, err := i.jobSeekerHideToAgentRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, jobSeeker := range jobSeekerList {
		for _, sh := range studentHistory {
			if jobSeeker.ID == sh.JobSeekerID {
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
		}

		for _, wh := range workHistory {
			if jobSeeker.ID == wh.JobSeekerID {
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
					if wh.ID == ei.WorkHistoryID {
						valueEI := entity.JobSeekerExperienceIndustry{
							ID:            ei.ID,
							WorkHistoryID: ei.WorkHistoryID,
							Industry:      ei.Industry,
						}

						value.ExperienceIndustries = append(value.ExperienceIndustries, valueEI)
					}
				}

				for _, dh := range departmentHistory {
					if wh.ID == dh.WorkHistoryID {
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
							if dh.ID == eo.DepartmentHistoryID {
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
		}

		for _, dcs := range desiredCompanyScale {
			if jobSeeker.ID == dcs.JobSeekerID {
				value := entity.JobSeekerDesiredCompanyScale{
					JobSeekerID:         dcs.JobSeekerID,
					DesiredCompanyScale: dcs.DesiredCompanyScale,
				}

				jobSeeker.DesiredCompanyScales = append(jobSeeker.DesiredCompanyScales, value)
			}
		}

		for _, l := range license {
			if jobSeeker.ID == l.JobSeekerID {
				value := entity.JobSeekerLicense{
					JobSeekerID:     l.JobSeekerID,
					LicenseType:     l.LicenseType,
					AcquisitionTime: l.AcquisitionTime,
				}

				jobSeeker.Licenses = append(jobSeeker.Licenses, value)
			}
		}

		for _, sp := range selfPromotion {
			if jobSeeker.ID == sp.JobSeekerID {
				value := entity.JobSeekerSelfPromotion{
					JobSeekerID: sp.JobSeekerID,
					Title:       sp.Title,
					Contents:    sp.Contents,
				}

				jobSeeker.SelfPromotions = append(jobSeeker.SelfPromotions, value)
			}
		}

		for _, d := range document {
			if jobSeeker.ID == d.JobSeekerID {
				value := entity.JobSeekerDocument{
					JobSeekerID:     d.JobSeekerID,
					ResumeOriginURL: d.ResumeOriginURL,
					ResumePDFURL:    d.ResumePDFURL,
					CVOriginURL:     d.CVOriginURL,
					CVPDFURL:        d.CVPDFURL,
					IDPhotoURL:      d.IDPhotoURL,
				}

				jobSeeker.Documents = value
			}
		}

		for _, di := range desiredIndustry {
			if jobSeeker.ID == di.JobSeekerID {
				value := entity.JobSeekerDesiredIndustry{
					JobSeekerID:     di.JobSeekerID,
					DesiredIndustry: di.DesiredIndustry,
					DesiredRank:     di.DesiredRank,
				}

				jobSeeker.DesiredIndustries = append(jobSeeker.DesiredIndustries, value)
			}
		}

		for _, do := range desiredOccupation {
			if jobSeeker.ID == do.JobSeekerID {
				value := entity.JobSeekerDesiredOccupation{
					JobSeekerID:       do.JobSeekerID,
					DesiredOccupation: do.DesiredOccupation,
					DesiredRank:       do.DesiredRank,
				}

				jobSeeker.DesiredOccupations = append(jobSeeker.DesiredOccupations, value)
			}
		}

		for _, dwl := range desiredWorkLocation {
			if jobSeeker.ID == dwl.JobSeekerID {
				value := entity.JobSeekerDesiredWorkLocation{
					JobSeekerID:         dwl.JobSeekerID,
					DesiredWorkLocation: dwl.DesiredWorkLocation,
					DesiredRank:         dwl.DesiredRank,
				}

				jobSeeker.DesiredWorkLocations = append(jobSeeker.DesiredWorkLocations, value)
			}
		}

		for _, dht := range desiredHolidayType {
			if jobSeeker.ID == dht.JobSeekerID {
				value := entity.JobSeekerDesiredHolidayType{
					JobSeekerID: dht.JobSeekerID,
					HolidayType: dht.HolidayType,
				}

				jobSeeker.DesiredHolidayTypes = append(jobSeeker.DesiredHolidayTypes, value)
			}
		}

		for _, ds := range developmentSkill {
			if jobSeeker.ID == ds.JobSeekerID {

				value := entity.JobSeekerDevelopmentSkill{
					JobSeekerID:         ds.JobSeekerID,
					DevelopmentCategory: ds.DevelopmentCategory,
					DevelopmentType:     ds.DevelopmentType,
					ExperienceYear:      ds.ExperienceYear,
					ExperienceMonth:     ds.ExperienceMonth,
				}

				jobSeeker.DevelopmentSkills = append(jobSeeker.DevelopmentSkills, value)
			}
		}

		for _, ls := range languageSkill {
			if jobSeeker.ID == ls.JobSeekerID {
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
		}

		for _, ps := range pcSkill {
			if jobSeeker.ID == ps.JobSeekerID {
				value := entity.JobSeekerPCTool{
					JobSeekerID: ps.JobSeekerID,
					Tool:        ps.Tool,
				}

				jobSeeker.PCTools = append(jobSeeker.PCTools, value)
			}
		}

		for _, hta := range hideToAgent {
			if jobSeeker.ID == hta.JobSeekerID {
				value := entity.JobSeekerHideToAgent{
					JobSeekerID: hta.JobSeekerID,
					AgentID:     hta.AgentID,
					AgentName:   hta.AgentName,
				}

				jobSeeker.HideToAgents = append(jobSeeker.HideToAgents, value)
			}
		}
	}

	// csvの一行目を作成
	records = append(
		records,
		[]string{
			"CA担当者",
			"フェーズ",
			"流入経路",
			"エントリー日時",
			"面談日時",
			"ステータス",
			"苗字",
			"名前",
			"苗字(カナ)",
			"名前(カナ)",
			"性別",
			"性別備考",
			"国籍",
			"国籍備考",
			"既往歴",
			"既往歴備考",
			"生年月日",
			"配偶者有無",
			"配偶者の扶養義務",
			"扶養家族（配偶者を除く）",
			"メールアドレス",
			"電話番号",
			"緊急連絡先（電話番号）",
			"住所(都道府県)",
			"住所",
			"住所(フリガナ)",
			"直近の年収",

			// 学歴
			"文系・理系",
			"学歴詳細",

			// 職歴
			"就業状況",
			"転職回数",
			"短期離職（1年未満）",
			"短期離職備考",
			"職歴詳細",

			// PCスキル
			"Excel",
			"Word",
			"PowerPoint",

			// 活かせるスキル・資格
			"資格",
			"業務ツール",
			"開発経験(言語)",
			"開発経験(OS)",
			"語学",

			// 希望条件
			"希望休日タイプ",
			"希望企業規模",
			"希望業界",
			"希望職種",
			"希望勤務地",
			"転勤可否",
			"転勤備考",
			"希望年収",
			"入社可能時期",

			// お人柄
			"アピアランス",
			// "アピアランス(本音)",
			// "アピアランス(推薦状用)",
			"コミュニケーション",
			// "コミュニケーション(本音)",
			// "コミュニケーション(推薦状用)",
			"論理的思考力",
			// "論理的思考力(本音)",
			// "論理的思考力(推薦状用)",
			"自己PR",
			"研究内容・学チカ",
			"応募承諾のポイント",

			// メモ
			"社内限定メモ",
			// "他社エージェント向けメモ",
		},
	)

	for _, jobSeeker := range jobSeekerList {
		// 求人ごとにデータを格納
		record =
			[]string{
				jobSeeker.StaffName,
				getStrPhaseForJobSeeker(jobSeeker.Phase),
				jobSeeker.ChannelName,
				jobSeeker.CreatedAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05"),
				jobSeeker.InterviewDate.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05"),
				getStrUserStatus(jobSeeker.UserStatus),
				jobSeeker.LastName,
				jobSeeker.FirstName,
				jobSeeker.LastFurigana,
				jobSeeker.FirstFurigana,
				getStrGenderForJobSeeker(jobSeeker.Gender),
				jobSeeker.GenderRemarks,
				getStrNationalityForJobSeeker(jobSeeker.Nationality),
				jobSeeker.NationalityRemarks,
				getStrAvailable(jobSeeker.MedicalHistory),
				jobSeeker.MedicalHistoryRemarks,
				jobSeeker.Birthday,
				getStrAvailable(jobSeeker.Spouse),
				getStrAvailable(jobSeeker.SupportObligation),
				fmt.Sprint(jobSeeker.Dependents.Int64),
				jobSeeker.Email,
				jobSeeker.PhoneNumber,
				jobSeeker.EmergencyPhoneNumber,
				getStrPrefecture(jobSeeker.Prefecture),
				fmt.Sprint(jobSeeker.PostCode, "\n", getStrPrefecture(jobSeeker.Prefecture), jobSeeker.Address),
				jobSeeker.AddressFurigana,
				fmt.Sprint(jobSeeker.AnnualIncome.Int64) + "万円",

				// 学歴
				getStrStudyCategoryForJobSeeker(jobSeeker.StudyCategory),
				getStrStudentHistoryList(jobSeeker.StudentHistories),

				// 職歴
				getStrStateOfEmployment(jobSeeker.StateOfEmployment),
				getStrJobChangeForJobSeeker(jobSeeker.JobChange),
				getStrAvailable(jobSeeker.ShortResignation),
				jobSeeker.ShortResignationRemarks,
				getStrWorkHistoryList(jobSeeker.WorkHistories),

				// PCスキル
				getStrExcelSkill(jobSeeker.ExcelSkill),
				getStrWordSkill(jobSeeker.WordSkill),
				getStrPowerPointSkill(jobSeeker.PowerPointSkill),

				// 活かせるスキル・資格
				getStrLicenseListForJobSeeker(jobSeeker.Licenses),
				getStrPCToolListForJobSeeker(jobSeeker.PCTools),
				getStrLanguageDevelopmentExperienceListForJobSeeker(jobSeeker.DevelopmentSkills),
				getStrOSDevelopmentExperienceListForJobSeeker(jobSeeker.DevelopmentSkills),
				getStrLanguageListForJobSeeker(jobSeeker.LanguageSkills),

				// 希望条件
				getStrDesiredHolidayTypeList(jobSeeker.DesiredHolidayTypes),
				getStrDesiredCompanyScaleList(jobSeeker.DesiredCompanyScales),
				getStrDesiredIndustryList(jobSeeker.DesiredIndustries),
				getStrDesiredOccupationList(jobSeeker.DesiredOccupations),
				getStrDesiredWorkLocationList(jobSeeker.DesiredWorkLocations),
				getStrTransfer(jobSeeker.Transfer),
				jobSeeker.TransferRequirement,
				fmt.Sprint(jobSeeker.DesiredAnnualIncome.Int64) + "万円",
				fmt.Sprint(getStrJoinCompanyPeriod(jobSeeker.JoinCompanyPeriod)),

				getStrAppearanceForJobSeeker(jobSeeker.Appearance),
				// jobSeeker.AppearanceDetailOfTruth,
				// jobSeeker.AppearanceDetail,
				getStrCommunicationForJobSeeker(jobSeeker.Communication),
				// jobSeeker.CommunicationDetailOfTruth,
				// jobSeeker.CommunicationDetail,
				getStrThinkingForJobSeeker(jobSeeker.Thinking),
				// jobSeeker.ThinkingDetailOfTruth,
				// jobSeeker.ThinkingDetail,
				getStrSelfPromotionList(jobSeeker.SelfPromotions),
				jobSeeker.ResearchContent,
				jobSeeker.AcceptancePoints,

				// メモ
				jobSeeker.SecretMemo,
				// jobSeeker.PublicMemo,
			}

			// レコードを追加
		records = append(records, record)
	}

	//CSVファイルを作成
	filePath := ("./job-seeker-" + fmt.Sprint(utility.CreateUUID()) + ".csv")

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
