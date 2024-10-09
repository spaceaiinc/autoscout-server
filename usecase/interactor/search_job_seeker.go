package interactor

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"gopkg.in/guregu/null.v4"
)

// エージェントIDから企業情報一覧を取得する
type GetSearchJobSeekerListByAgentIDInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchJobSeeker
}

type GetSearchJobSeekerListByAgentIDOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
}

func (i *JobSeekerInteractorImpl) GetSearchJobSeekerListByAgentID(input GetSearchJobSeekerListByAgentIDInput) (GetSearchJobSeekerListByAgentIDOutput, error) {
	var (
		output        GetSearchJobSeekerListByAgentIDOutput
		err           error
		jobSeekerList []*entity.JobSeeker
	)
	/**
	GetJobInformationListByAgentIDAndFreeWordは
	フリーワードの有無で処理を分岐

	フリーワードは社名のみ
	*/

	jobSeekerList, err = i.jobSeekerRepository.GetByAgentIDAndFreeWord(input.AgentID, input.SearchParam.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者の子テーブル情報をセット
	jobSeekerList, err = setJobSeekerChildTableByIDList(i, jobSeekerList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 絞り込み検索処理
	jobSeekerListWithThinking, err := searchJobSeekerList(jobSeekerList, input.SearchParam)
	if err != nil {
		return output, err
	}

	// IDListを返す
	for _, jobSeeker := range jobSeekerListWithThinking {
		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerListWithThinking)

	// 指定ページの求職者20件を取得
	output.JobSeekerList = getJobSeekerListWithPage(jobSeekerListWithThinking, input.PageNumber)

	return output, nil

}

type GetSearchActiveJobSeekerListByAgentIDInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchJobSeeker
}

type GetSearchActiveJobSeekerListByAgentIDOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
}

func (i *JobSeekerInteractorImpl) GetSearchActiveJobSeekerListByAgentID(input GetSearchActiveJobSeekerListByAgentIDInput) (GetSearchActiveJobSeekerListByAgentIDOutput, error) {
	var (
		output        GetSearchActiveJobSeekerListByAgentIDOutput
		err           error
		jobSeekerList []*entity.JobSeeker
	)
	/**
	GetActiveListByAgentIDAndFreeWordは
	フリーワードの有無で処理を分岐

	フリーワードは社名のみ
	*/

	//
	jobSeekerList, err = i.jobSeekerRepository.GetActiveOwnAndFreeWord(input.AgentID, input.SearchParam.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者の子テーブル情報をセット
	jobSeekerList, err = setJobSeekerChildTableByIDList(i, jobSeekerList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 絞り込み検索処理
	jobSeekerListWithThinking, err := searchJobSeekerList(jobSeekerList, input.SearchParam)
	if err != nil {
		return output, err
	}

	// IDListを返す
	for _, jobSeeker := range jobSeekerListWithThinking {
		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerListWithThinking)

	// 指定ページの求職者20件を取得
	output.JobSeekerList = getJobSeekerListWithPage(jobSeekerListWithThinking, input.PageNumber)

	return output, nil

}

// エージェントIDから企業情報一覧を取得する
type GetSearchAllianceJobSeekerListByAgentIDInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchJobSeeker
}

type GetSearchAllianceJobSeekerListByAgentIDOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
}

func (i *JobSeekerInteractorImpl) GetSearchAllianceJobSeekerListByAgentID(input GetSearchAllianceJobSeekerListByAgentIDInput) (GetSearchAllianceJobSeekerListByAgentIDOutput, error) {
	var (
		output        GetSearchAllianceJobSeekerListByAgentIDOutput
		err           error
		jobSeekerList []*entity.JobSeeker
	)
	/**
	他社エージェントの求職者を取得する
	*/

	//アライアンスを締結しているエージェントを取得
	agentAllianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// アライアンス締結済みの他社求職者（非公開除外前）
	listBeforeHideToAgent, err := i.jobSeekerRepository.GetByOtherAgentIDAndFreeWord(input.AgentID, input.SearchParam.FreeWord, agentAllianceList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 他のエージェントを含んでいる場合、hideToAgentをチェック
	// 他エージェントの求人IDリストから非公開エージェントを取得
	hideToAgent, err := i.jobSeekerHideToAgentRepository.GetHideByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if len(hideToAgent) > 0 {
		// 他エージェントの求人IDリストから非公開エージェントを除外
		jobSeekerList = checkJobSeekerByHideToAgent(listBeforeHideToAgent, hideToAgent)
	} else {
		jobSeekerList = listBeforeHideToAgent
	}

	if len(jobSeekerList) < 1 {
		fmt.Println("求職者が見つかりませんでした")
		return output, nil
	}

	idList := getJobSeekerIDList(jobSeekerList)

	// 求職者IDリストから関連テーブル情報を取得
	// selfpromotionとdocumentは使わないので取得しない
	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	pcTool, err := i.jobSeekerPCToolRepository.GetByJobSeekerIDList(idList)
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
					LastStatus:     sh.LastStatus,
					GraduationYear: sh.GraduationYear,
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
				}

				for _, ei := range experienceIndustry {
					if wh.ID == ei.WorkHistoryID {
						valueEI := entity.JobSeekerExperienceIndustry{
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

		for _, ps := range pcTool {
			if jobSeeker.ID == ps.JobSeekerID {
				value := entity.JobSeekerPCTool{
					JobSeekerID: ps.JobSeekerID,
					Tool:        ps.Tool,
				}

				jobSeeker.PCTools = append(jobSeeker.PCTools, value)
			}
		}

	}

	// 絞り込み検索処理
	jobSeekerListWithThinking, err := searchJobSeekerList(jobSeekerList, input.SearchParam)
	if err != nil {
		return output, err
	}

	// IDListを返す
	for _, jobSeeker := range jobSeekerListWithThinking {
		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerListWithThinking)

	// 指定ページの求職者20件を取得
	jobSeekerList20 := getJobSeekerListWithPage(jobSeekerListWithThinking, input.PageNumber)

	for _, jobSeeker := range jobSeekerList20 {
		// SecretMemo, Phaseを空にする
		jobSeeker.SecretMemo = ""
		jobSeeker.Phase = null.NewInt(0, false)
		jobSeeker.InflowChannelID = null.NewInt(0, false)
		jobSeeker.ChannelName = ""
		output.JobSeekerList = append(output.JobSeekerList, jobSeeker)
	}

	return output, nil

}

/****************************************************************************************/
// 求人検索→求職者検索(絞り込み) API
//

// エージェントIDとページ番号から20件取得
type GetSearchJobSeekerListByAgentIDAndTypeInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchJobSeeker
	Type        entity.JobSeekerType
}

type GetSearchJobSeekerListByAgentIDAndTypeOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
	AllCount      uint
	OwnCount      uint
	AllianceCount uint
}

func (i *JobSeekerInteractorImpl) GetSearchJobSeekerListByAgentIDAndType(input GetSearchJobSeekerListByAgentIDAndTypeInput) (GetSearchJobSeekerListByAgentIDAndTypeOutput, error) {
	var (
		output                GetSearchJobSeekerListByAgentIDAndTypeOutput
		err                   error
		jobSeekerList         []*entity.JobSeeker
		AllJobSeekerList      []*entity.JobSeeker
		OwnJobSeekerList      []*entity.JobSeeker
		AllianceJobSeekerList []*entity.JobSeeker
	)

	// 1. 求職者全て取得
	// 2. 同一求職者・非公開先の除外と絞り込み
	// 3. Typeに応じて返す求職者リストを変更する
	// 4. その他必要な処理

	/************ 1. 求職者全て取得 **************/

	jobSeekerListBeforeDuplicate, err := i.jobSeekerRepository.GetActiveAllAndFreeWord(input.AgentID, input.SearchParam.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 特別仕様: 本番環境のみ「2: 株式会社テスト」と「3: 株式会社Space AI（非公開求人管理用）」を除外して他社エージェントに非表示にする
	jobSeekerListBeforeDuplicate = excludeTestJobSeeker(jobSeekerListBeforeDuplicate, input.AgentID)

	/************ 2. 同一求職者・非公開先の除外と絞り込み **************/

	// 求職者の子テーブル情報をセット
	jobSeekerListBeforeDuplicate, err = setJobSeekerChildTableByIDList(i, jobSeekerListBeforeDuplicate)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 絞り込み検索処理
	jobSeekerListBeforeDuplicate, err = searchJobSeekerList(jobSeekerListBeforeDuplicate, input.SearchParam)
	if err != nil {
		return output, err
	}

	//　同一求人の除外
	jobSeekerList = excludeDuplicateJobSeeker(jobSeekerListBeforeDuplicate, input.AgentID)

	// 指定のAgentIDを非公開先にしている求人の非公開設定情報を取得
	hideToAgent, err := i.jobSeekerHideToAgentRepository.GetHideByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if len(hideToAgent) > 0 {
		// 非公開情報と合致する求人を除外して新しく求人リストを取得
		jobSeekerList = checkJobSeekerByHideToAgent(jobSeekerList, hideToAgent)
	}

	/************ 3. Typeに応じて返す求職者リストを変更する **************/

	// 指定AgentIDのアライアンス情報を取得
	var allianceIDList []uint

	agentAllianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// アライアンスエージェントのIDリストを作成
	for _, alliance := range agentAllianceList {
		if alliance.Agent1ID != input.AgentID {
			allianceIDList = append(allianceIDList, alliance.Agent1ID)
		} else {
			allianceIDList = append(allianceIDList, alliance.Agent2ID)
		}
	}

	// 「すべて・自社・シェア・お助け」のそれぞれの件数を取得する
	for _, jobSeeker := range jobSeekerList {
		// 自社のカウント
		if jobSeeker.AgentID == input.AgentID {
			output.OwnCount = output.OwnCount + 1
			OwnJobSeekerList = append(OwnJobSeekerList, jobSeeker)

			output.AllCount = output.AllCount + 1
			AllJobSeekerList = append(AllJobSeekerList, jobSeeker)
		} else if jobSeeker.AgentID != input.AgentID && includeUINT(allianceIDList, jobSeeker.AgentID) {
			// SecretMemo, Phaseを空にする
			jobSeeker.SecretMemo = ""
			jobSeeker.Phase = null.NewInt(0, false)
			jobSeeker.InflowChannelID = null.NewInt(0, false)
			jobSeeker.ChannelName = ""

			output.AllianceCount = output.AllianceCount + 1
			AllianceJobSeekerList = append(AllianceJobSeekerList, jobSeeker)

			output.AllCount = output.AllCount + 1
			AllJobSeekerList = append(AllJobSeekerList, jobSeeker)
		}
	}

	// 検索の種類に応じて返すリストを変更
	if input.Type == entity.TypeAllJobSeeker {
		jobSeekerList = AllJobSeekerList
	} else if input.Type == entity.TypeOwnJobSeeker {
		jobSeekerList = OwnJobSeekerList
	} else if input.Type == entity.TypeAllianceJobSeeker {
		jobSeekerList = AllianceJobSeekerList
	} else {
		err = fmt.Errorf("%v:%w", err, entity.ErrRequestError)
		return output, err
	}

	/************ 4. その他必要な処理 **************/

	// IDListを返す
	for _, jobSeeker := range jobSeekerList {
		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerList)

	// 指定ページの求職者20件を取得
	output.JobSeekerList = getJobSeekerListWithPage(jobSeekerList, input.PageNumber)

	return output, nil
}

/****************************************************************************************/
// シェア求人検索→自社求職者検索(絞り込み) API
//
type GetSearchPublicJobSeekerListByAgentIDAndPageInput struct {
	AgentID              uint
	PageNumber           uint
	JobInformationIDList []uint
	SearchParam          entity.SearchJobSeeker
}

type GetSearchPublicJobSeekerListByAgentIDAndPageOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
}

func (i *JobSeekerInteractorImpl) GetSearchPublicJobSeekerListByAgentIDAndPage(input GetSearchPublicJobSeekerListByAgentIDAndPageInput) (GetSearchPublicJobSeekerListByAgentIDAndPageOutput, error) {
	var (
		output                         GetSearchPublicJobSeekerListByAgentIDAndPageOutput
		err                            error
		jobSeekerList                  []*entity.JobSeeker
		jobSeekerlistBeforeHideToAgent []*entity.JobSeeker
		agentIDList                    []uint
	)

	// 1. 自社求職者全て取得
	// 2. 非公開先の除外と絞り込み
	// 3. その他必要な処理

	/************ 1. 求職者全て取得 **************/

	jobSeekerlistBeforeHideToAgent, err = i.jobSeekerRepository.GetActiveOwnAndFreeWord(input.AgentID, input.SearchParam.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. 非公開先の除外と絞り込み **************/

	// 求職者の子テーブル情報をセット
	jobSeekerlistBeforeHideToAgent, err = setJobSeekerChildTableByIDList(i, jobSeekerlistBeforeHideToAgent)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 絞り込み検索処理
	jobSeekerlistBeforeHideToAgent, err = searchJobSeekerList(jobSeekerlistBeforeHideToAgent, input.SearchParam)
	if err != nil {
		return output, err
	}

	jobInformationList, err := i.jobInformationRepository.GetByIDList(input.JobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求人に紐づいたエージェントIDのリストを作成
	for _, jobInformation := range jobInformationList {
		agentIDList = append(agentIDList, jobInformation.AgentID)
	}

	if len(agentIDList) > 0 {
		// 指定エージェントを非公開先に設定している求職者の情報取得
		hideToAgent, err := i.jobSeekerHideToAgentRepository.GetByAgentIDList(agentIDList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 他エージェントの求人IDリストから非公開エージェントを除外
		if len(hideToAgent) > 0 {
			jobSeekerList = checkJobSeekerByHideToAgent(jobSeekerlistBeforeHideToAgent, hideToAgent)
		} else {
			jobSeekerList = jobSeekerlistBeforeHideToAgent
		}
	}

	/************ 3. その他必要な処理 **************/

	// IDListを返す
	for _, jobSeeker := range jobSeekerList {
		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerList)

	// 指定ページの求職者20件を取得
	output.JobSeekerList = getJobSeekerListWithPage(jobSeekerList, input.PageNumber)

	return output, nil
}

/****************************************************************************************/

// 絞り込み検索処理 本体
func searchJobSeekerList(jobSeekerList []*entity.JobSeeker, searchParam entity.SearchJobSeeker) ([]*entity.JobSeeker, error) {

	type searchJobSeekerListInput struct {
		SearchParam entity.SearchJobSeeker
	}

	input := searchJobSeekerListInput{
		SearchParam: searchParam,
	}

	// 絞り込み項目の結果を代入するための変数を用意
	var (
		jobSeekerListWithAgentStaffID []*entity.JobSeeker
		/**
		アライアンスの中間テーブル作成後に
		アライアンスの絞り込みを追加
		*/
		jobSeekerListWithPhase                []*entity.JobSeeker
		jobSeekerListWithGender               []*entity.JobSeeker
		jobSeekerListWithAgeUnder             []*entity.JobSeeker
		jobSeekerListWithAgeOver              []*entity.JobSeeker
		jobSeekerListWithDesiredIndustry      []*entity.JobSeeker
		jobSeekerListWithDesiredOccupation    []*entity.JobSeeker
		jobSeekerListWithDesiredWorkLocation  []*entity.JobSeeker
		jobSeekerListWithFinalEducation       []*entity.JobSeeker
		jobSeekerListWithStudyCategoryTypes   []*entity.JobSeeker
		jobSeekerListWithSchoolLevel          []*entity.JobSeeker
		jobSeekerListWithNationality          []*entity.JobSeeker
		jobSeekerListWithJobChange            []*entity.JobSeeker
		jobSeekerListWithShortResignation     []*entity.JobSeeker
		jobSeekerListWithUnderIncome          []*entity.JobSeeker
		jobSeekerListWithOverIncome           []*entity.JobSeeker
		jobSeekerListWithDesiredTransfer      []*entity.JobSeeker
		jobSeekerListWithDesiredHolidayType   []*entity.JobSeeker
		jobSeekerListWithDesiredCompanyScale  []*entity.JobSeeker
		jobSeekerListWithExperienceIndustry   []*entity.JobSeeker // And検索
		jobSeekerListWithExperienceOccupation []*entity.JobSeeker // And検索
		jobSeekerListWithSocialExperience     []*entity.JobSeeker
		jobSeekerListWithManagement           []*entity.JobSeeker // And検索
		jobSeekerListWithLicense              []*entity.JobSeeker // And検索
		jobSeekerListWithLanguage             []*entity.JobSeeker // And検索
		jobSeekerListWithExcelSkill           []*entity.JobSeeker
		jobSeekerListWithWordSkill            []*entity.JobSeeker
		jobSeekerListWithPowerPointSkill      []*entity.JobSeeker
		jobSeekerListWithAnotherPCSkill       []*entity.JobSeeker // And検索
		jobSeekerListWithDevelopmentLanguage  []*entity.JobSeeker // And検索
		jobSeekerListWithDevelopmentOS        []*entity.JobSeeker // And検索
		jobSeekerListWithAppearance           []*entity.JobSeeker
		jobSeekerListWithCommunication        []*entity.JobSeeker
		jobSeekerListWithThinking             []*entity.JobSeeker

		/**
		年齢の絞り込みの際に使用する
		*/
		now                      = time.Now()
		thisYear, thisMonth, day = now.Date()
	)

	// 営業担当者IDがある場合
	agentStaffID, err := strconv.Atoi(input.SearchParam.AgentStaffID)
	if !(err != nil || agentStaffID == 0) {
		for _, jobInformation := range jobSeekerList {
			if jobInformation.AgentStaffID.Int64 == int64(agentStaffID) {
				jobSeekerListWithAgentStaffID = append(jobSeekerListWithAgentStaffID, jobInformation)
			}
		}
	}

	// 営業担当者IDが無い場合
	if err != nil || agentStaffID == 0 {
		jobSeekerListWithAgentStaffID = jobSeekerList
	}

	fmt.Println("CA担当者: ", len(jobSeekerListWithAgentStaffID))

	// Note: フェーズ
	// フェーズがある場合
	if !(len(input.SearchParam.PhaseTypes) == 0) {
	phaseLoop:
		for _, jobSeeker := range jobSeekerListWithAgentStaffID {
			for _, phase := range input.SearchParam.PhaseTypes {
				if !phase.Valid {
					continue
				}

				if phase == jobSeeker.Phase {
					// 求職者のフェーズが一つでも合致していればcontinueで次の求職者へ
					jobSeekerListWithPhase = append(jobSeekerListWithPhase, jobSeeker)
					continue phaseLoop
				}
			}
		}
	}

	// 性別のいずれかも入っていない場合
	if len(input.SearchParam.PhaseTypes) == 0 {
		jobSeekerListWithPhase = jobSeekerListWithAgentStaffID
	}

	fmt.Println("フェーズ: ", len(jobSeekerListWithPhase))

	// Note: 性別
	// 性別がある場合
	if !(len(input.SearchParam.GenderTypes) == 0) {
	genderLoop:
		for _, jobSeeker := range jobSeekerListWithPhase {
			for _, gender := range input.SearchParam.GenderTypes {
				if !gender.Valid {
					continue
				}
				if gender == jobSeeker.Gender {
					// 求職者の希望性別が一つでも合致していればcontinueで次の求職者へ
					jobSeekerListWithGender = append(jobSeekerListWithGender, jobSeeker)
					continue genderLoop
				}
			}
		}
	}

	// 性別のいずれかも入っていない場合
	if len(input.SearchParam.GenderTypes) == 0 {
		jobSeekerListWithGender = jobSeekerListWithPhase
	}

	// Note: 年齢下限
	// 年齢下限がある場合
	ageUnder, err := strconv.Atoi(input.SearchParam.AgeUnder)
	if !(err != nil) {
		for _, jobSeeker := range jobSeekerListWithGender {
			if jobSeeker.Birthday == "" {
				continue
			}

			// 求職者の誕生日をパース
			birthday, err := time.Parse("2006-01-02", jobSeeker.Birthday)
			if err != nil {
				continue
			}

			// 求職者の年齢
			jobSeekerAge := thisYear - birthday.Year()

			// 誕生日を迎えていない場合はageを「−1」する
			if thisMonth <= birthday.Month() && day <= birthday.Day() {
				jobSeekerAge -= 1
			}

			if ageUnder <= jobSeekerAge {
				jobSeekerListWithAgeUnder = append(jobSeekerListWithAgeUnder, jobSeeker)
			}
		}
	}

	// 年齢下限が無い場合
	if len(input.SearchParam.AgeUnder) == 0 {
		jobSeekerListWithAgeUnder = jobSeekerListWithGender
	}

	fmt.Println("年齢下限: ", len(jobSeekerListWithAgeUnder))

	// Note: 年齢上限
	// 年齢上限がある場合
	ageOver, err := strconv.Atoi(input.SearchParam.AgeOver)
	if !(err != nil) {
		for _, jobSeeker := range jobSeekerListWithAgeUnder {
			if jobSeeker.Birthday == "" {
				continue
			}

			// 求職者の誕生日をパース
			birthday, err := time.Parse("2006-01-02", jobSeeker.Birthday)
			if err != nil {
				continue
			}

			// 求職者の年齢
			jobSeekerAge := thisYear - birthday.Year()

			// 誕生日を迎えていない場合はageを「−1」する
			if thisMonth <= birthday.Month() && day <= birthday.Day() {
				jobSeekerAge -= 1
			}

			if ageOver >= jobSeekerAge {
				jobSeekerListWithAgeOver = append(jobSeekerListWithAgeOver, jobSeeker)
			}
		}
	}

	// 年齢上限が無い場合
	if len(input.SearchParam.AgeOver) == 0 {
		jobSeekerListWithAgeOver = jobSeekerListWithAgeUnder
	}

	fmt.Println("年齢上限: ", len(jobSeekerListWithAgeOver))

	// NOTE: 業界（And検索）
	// 業界のいずれかが入っている場合
	if !(len(input.SearchParam.DesiredIndustries) == 0) {
	industryLoop:
		for _, jobSeeker := range jobSeekerListWithAgeOver {
			if len(jobSeeker.DesiredIndustries) == 0 {
				continue
			}

			/**
			合致した時にindustryLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/
			for _, desiredIndustryParam := range input.SearchParam.DesiredIndustries {
				if !desiredIndustryParam.Valid {
					continue
				}

				for _, jobSeekerIndustry := range jobSeeker.DesiredIndustries {
					if desiredIndustryParam == jobSeekerIndustry.DesiredIndustry || jobSeekerIndustry.DesiredIndustry == null.NewInt(9999, true) {
						jobSeekerListWithDesiredIndustry = append(jobSeekerListWithDesiredIndustry, jobSeeker)
						continue industryLoop
					}
				}
			}
		}
	}

	// 業界のいずれかも入っていない場合
	if len(input.SearchParam.DesiredIndustries) == 0 {
		jobSeekerListWithDesiredIndustry = jobSeekerListWithAgeOver
	}

	fmt.Println("業界: ", len(jobSeekerListWithDesiredIndustry))

	// NOTE: 職種（And検索）
	// 職種のいずれかが入っている場合
	if !(len(input.SearchParam.DesiredOccupations) == 0) {
	occuptionLoop:
		for _, jobSeeker := range jobSeekerListWithDesiredIndustry {
			if len(jobSeeker.DesiredOccupations) == 0 {
				continue
			}

			/**
			合致した時にoccuptionLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/
			for _, desiredOccupationParam := range input.SearchParam.DesiredOccupations {
				if !desiredOccupationParam.Valid {
					continue
				}

				for _, jobSeekerOccupation := range jobSeeker.DesiredOccupations {
					if desiredOccupationParam == jobSeekerOccupation.DesiredOccupation || jobSeekerOccupation.DesiredOccupation == null.NewInt(9999, true) {
						jobSeekerListWithDesiredOccupation = append(jobSeekerListWithDesiredOccupation, jobSeeker)
						continue occuptionLoop
					}
				}
			}
		}
	}

	// 職種のいずれかも入っていない場合
	if len(input.SearchParam.DesiredOccupations) == 0 {
		jobSeekerListWithDesiredOccupation = jobSeekerListWithDesiredIndustry
	}

	fmt.Println("職種: ", len(jobSeekerListWithDesiredOccupation))

	// NOTE: 勤務地
	// 勤務地のいずれかが入っている場合
	if !(len(input.SearchParam.DesiredWorkLocations) == 0) {

	locationLoop:
		for _, jobSeeker := range jobSeekerListWithDesiredOccupation {
			if len(jobSeeker.DesiredWorkLocations) == 0 {
				continue
			}
			/**
			合致した時にlocationLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			for _, desiredWorkLocationParam := range input.SearchParam.DesiredWorkLocations {
				if !desiredWorkLocationParam.Valid {
					continue
				}

				for _, jobSeekerWorkLocation := range jobSeeker.DesiredWorkLocations {
					if desiredWorkLocationParam == jobSeekerWorkLocation.DesiredWorkLocation || jobSeekerWorkLocation.DesiredWorkLocation == null.NewInt(9999, true) {
						jobSeekerListWithDesiredWorkLocation = append(jobSeekerListWithDesiredWorkLocation, jobSeeker)
						continue locationLoop
					}
				}
			}
		}
	}

	// 勤務地のいずれかも入っていない場合
	if len(input.SearchParam.DesiredWorkLocations) == 0 {
		jobSeekerListWithDesiredWorkLocation = jobSeekerListWithDesiredOccupation
	}

	fmt.Println("勤務地: ", len(jobSeekerListWithDesiredWorkLocation))

	// Note: 最終学歴
	// 最終学歴がある場合
	if !(len(input.SearchParam.FinalEducationTypes) == 0) {
	finalEducationLoop:
		for _, jobSeeker := range jobSeekerListWithDesiredWorkLocation {
			studentHistoryLength := len(jobSeeker.StudentHistories)
			if studentHistoryLength == 0 {
				continue
			}

			// seeker.StudentHistoriesの中から一番高い学歴を取得
			(sort.Slice(
				jobSeeker.StudentHistories,
				func(i, j int) bool {
					return jobSeeker.StudentHistories[i].SchoolCategory.Int64 > jobSeeker.StudentHistories[j].SchoolCategory.Int64
				}))

			// 最終学歴の番号を使う
			finalHistoryNumber := 0

			// 卒業, 中退, 退学, 卒業見込み, 修了
			// LastStatusが中退(index:1)と退学(index:2)の場合は、一つ前の学歴を使う
			for _, history := range jobSeeker.StudentHistories {
				if history.LastStatus.Int64 == 1 || history.LastStatus.Int64 == 2 || !history.LastStatus.Valid {
					finalHistoryNumber++
				} else {
					break
				}
			}
			// 全て中退or退学の場合は、次に移る
			if finalHistoryNumber >= studentHistoryLength {
				continue
			}
			for _, finalEducation := range input.SearchParam.FinalEducationTypes {
				if !finalEducation.Valid {
					continue
				}
				// 絞り込みパラムと合致している場合はリストに追加
				if jobSeeker.StudentHistories[finalHistoryNumber].SchoolCategory == finalEducation {
					jobSeekerListWithFinalEducation = append(jobSeekerListWithFinalEducation, jobSeeker)
					continue finalEducationLoop
				}
			}
		}
	}

	//  最終学歴のいずれかも入っていない場合
	if len(input.SearchParam.FinalEducationTypes) == 0 {
		jobSeekerListWithFinalEducation = jobSeekerListWithDesiredWorkLocation
	}

	fmt.Println("最終学歴: ", len(jobSeekerListWithFinalEducation))

	// Note: 文系・理系
	// ある場合
	if !(len(input.SearchParam.StudyCategoryTypes) == 0) {
	studyCategoryLoop:
		for _, jobSeeker := range jobSeekerListWithFinalEducation {
			if jobSeeker.StudyCategory == null.NewInt(0, false) {
				continue
			}
			for _, studyCategoryType := range input.SearchParam.StudyCategoryTypes {
				if !studyCategoryType.Valid {
					continue
				}
				// 絞り込みパラムと合致している場合はリストに追加
				if jobSeeker.StudyCategory == studyCategoryType {
					jobSeekerListWithStudyCategoryTypes = append(jobSeekerListWithStudyCategoryTypes, jobSeeker)
					continue studyCategoryLoop
				}
			}
		}
	}

	// いずれも入っていない場合
	if len(input.SearchParam.StudyCategoryTypes) == 0 {
		jobSeekerListWithStudyCategoryTypes = jobSeekerListWithFinalEducation
	}

	fmt.Println("文系・理系: ", len(jobSeekerListWithStudyCategoryTypes))

	// Note: 大学レベル
	// 大学レベルがある場合
	if !(len(input.SearchParam.SchoolLevelTypes) == 0) {
	schoolLevelLoop:
		for _, jobSeeker := range jobSeekerListWithStudyCategoryTypes {
			for _, schoolLevel := range input.SearchParam.SchoolLevelTypes {
				if !schoolLevel.Valid {
					continue
				}

				for _, history := range jobSeeker.StudentHistories {
					if !history.SchoolCategory.Valid ||
						(history.SchoolCategory.Int64 != 5 && history.SchoolCategory.Int64 != 6) ||
						history.LastStatus.Int64 == 1 ||
						history.LastStatus.Int64 == 2 {
						continue
					}

					// 絞り込みパラムと合致している場合はリストに追加
					if history.SchoolLevel == schoolLevel {
						jobSeekerListWithSchoolLevel = append(jobSeekerListWithSchoolLevel, jobSeeker)
						continue schoolLevelLoop
					}
				}
			}
		}
	}

	//  大学レベルのいずれかも入っていない場合
	if len(input.SearchParam.SchoolLevelTypes) == 0 {
		jobSeekerListWithSchoolLevel = jobSeekerListWithStudyCategoryTypes
	}

	fmt.Println("大学レベル: ", len(jobSeekerListWithSchoolLevel))

	// Note: 国籍
	// 国籍がある場合
	if !(len(input.SearchParam.NationalityTypes) == 0) {
	nationalityLoop:
		for _, jobSeeker := range jobSeekerListWithSchoolLevel {
			for _, nationality := range input.SearchParam.NationalityTypes {
				if !nationality.Valid {
					continue
				}
				if nationality == jobSeeker.Nationality {
					// 求人の希望国籍が一つでも合致していればcontinueで次の求職者へ
					jobSeekerListWithNationality = append(jobSeekerListWithNationality, jobSeeker)
					continue nationalityLoop
				}
			}
		}
	}

	// 国籍が無い場合
	if len(input.SearchParam.NationalityTypes) == 0 {
		jobSeekerListWithNationality = jobSeekerListWithSchoolLevel
	}

	fmt.Println("国籍: ", len(jobSeekerListWithNationality))

	// Note: 転職回数
	// 転職回数がある場合
	if !(len(input.SearchParam.JobChangeTypes) == 0) {
	jobChangeLoop:
		for _, jobSeeker := range jobSeekerListWithNationality {
			if jobSeeker.JobChange == null.NewInt(0, false) {
				continue
			}
			for _, jobChange := range input.SearchParam.JobChangeTypes {
				if !jobChange.Valid {
					continue
				}

				if jobChange == jobSeeker.JobChange {
					// jobChangeが合致していた場合、continueで次の求職者へ
					jobSeekerListWithJobChange = append(jobSeekerListWithJobChange, jobSeeker)
					continue jobChangeLoop
				}
			}
		}
	}

	// 転職回数が無い場合
	if len(input.SearchParam.JobChangeTypes) == 0 {
		jobSeekerListWithJobChange = jobSeekerListWithNationality
	}

	fmt.Println("転職回数: ", len(jobSeekerListWithJobChange))

	// Note: 短期離職
	// 短期離職がある場合
	if !(len(input.SearchParam.ShortResignationTypes) == 0) {
	shortResignationLoop:
		for _, jobSeeker := range jobSeekerListWithJobChange {
			for _, shortResignation := range input.SearchParam.ShortResignationTypes {
				if !shortResignation.Valid {
					continue
				}
				if shortResignation == jobSeeker.ShortResignation {
					// 求人の短期離職が一つでも合致していればcontinueで次の求職者へ
					jobSeekerListWithShortResignation = append(jobSeekerListWithShortResignation, jobSeeker)
					continue shortResignationLoop
				}
			}
		}
	}

	// 短期離職が無い場合
	if len(input.SearchParam.ShortResignationTypes) == 0 {
		jobSeekerListWithShortResignation = jobSeekerListWithJobChange
	}

	fmt.Println("短期離職: ", len(jobSeekerListWithShortResignation))

	// Note: 年収下限
	// 年収下限がある場合
	underIncome, err := strconv.Atoi(input.SearchParam.UnderIncome)
	if !(err != nil) {
		for _, jobSeeker := range jobSeekerListWithShortResignation {
			if jobSeeker.DesiredAnnualIncome == null.NewInt(0, false) {
				continue
			}
			if underIncome <= int(jobSeeker.DesiredAnnualIncome.Int64) {
				jobSeekerListWithUnderIncome = append(jobSeekerListWithUnderIncome, jobSeeker)
			}
		}
	}

	// 年収下限が無い場合
	if err != nil {
		jobSeekerListWithUnderIncome = jobSeekerListWithShortResignation
	}

	fmt.Println("年収下限: ", len(jobSeekerListWithUnderIncome))

	// Note: 年収上限
	// 年収上限がある場合
	overIncome, err := strconv.Atoi(input.SearchParam.OverIncome)
	if !(err != nil) {
		for _, jobSeeker := range jobSeekerListWithUnderIncome {
			if jobSeeker.DesiredAnnualIncome == null.NewInt(0, false) {
				continue
			}

			if overIncome >= int(jobSeeker.DesiredAnnualIncome.Int64) {
				jobSeekerListWithOverIncome = append(jobSeekerListWithOverIncome, jobSeeker)
			}
		}
	}

	// 年収上限が無い場合
	if err != nil {
		jobSeekerListWithOverIncome = jobSeekerListWithUnderIncome
	}

	fmt.Println("年収上限: ", len(jobSeekerListWithOverIncome))

	// Note: 転勤有無
	// 転勤有無がある場合
	if !(len(input.SearchParam.DesiredTransferTypes) == 0) {
	transferLoop:
		for _, jobSeeker := range jobSeekerListWithOverIncome {
			for _, transfer := range input.SearchParam.DesiredTransferTypes {
				if !transfer.Valid {
					continue
				}
				if transfer == jobSeeker.Transfer {
					// 転勤有無が一つでも合致していればcontinueでループ抜ける
					jobSeekerListWithDesiredTransfer = append(jobSeekerListWithDesiredTransfer, jobSeeker)
					continue transferLoop
				}

				// if !transfer.Valid || transfer == null.NewInt(1, true)  {
				// 	// 転勤有無が一つでも合致していればcontinueでループ抜ける
				// 	jobSeekerListWithDesiredTransfer = append(jobSeekerListWithDesiredTransfer, jobSeeker)
				// 	continue transferLoop
				// }

				// // パラムが「転勤あり」で、求職者が「転勤可能」もしくは「転勤可能(条件あり)」
				// if (transfer == null.NewInt(0, true) && (jobSeeker.Transfer == null.NewInt(0, true) || jobSeeker.Transfer == null.NewInt(1, true))) {
				// 	// 転勤有無が一つでも合致していればcontinueでループ抜ける
				// 	jobSeekerListWithDesiredTransfer = append(jobSeekerListWithDesiredTransfer, jobSeeker)
				// 	continue transferLoop
				// }
			}
		}
	}

	// 転勤有無が無い場合
	if len(input.SearchParam.DesiredTransferTypes) == 0 {
		jobSeekerListWithDesiredTransfer = jobSeekerListWithOverIncome
	}

	fmt.Println("転勤有無: ", len(jobSeekerListWithDesiredTransfer))

	// Note: 希望休日タイプ
	// 希望休日タイプがある場合
	if !(len(input.SearchParam.DesiredHolidayTypes) == 0) {
	holidayLoop:
		for _, jobSeeker := range jobSeekerListWithDesiredTransfer {
			if len(jobSeeker.DesiredHolidayTypes) == 0 {
				continue
			}

			/**
			合致した時にtrueに変える
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			for _, holidayType := range input.SearchParam.DesiredHolidayTypes {
				if !holidayType.Valid {
					continue
				}

				for _, jobSeekerHolidayType := range jobSeeker.DesiredHolidayTypes {
					if holidayType == jobSeekerHolidayType.HolidayType || jobSeekerHolidayType.HolidayType == null.NewInt(99, true) {
						// 休日休暇が一つでも合致していればcontinueで次の求職者へ
						jobSeekerListWithDesiredHolidayType = append(jobSeekerListWithDesiredHolidayType, jobSeeker)
						continue holidayLoop
					}
				}
			}
		}
	}

	// 希望休日タイプが無い場合
	if len(input.SearchParam.DesiredHolidayTypes) == 0 {
		jobSeekerListWithDesiredHolidayType = jobSeekerListWithDesiredTransfer
	}

	fmt.Println("休日タイプ: ", len(jobSeekerListWithDesiredHolidayType))

	// Note: 希望企業規模の絞り込みは企業の従業員数（単体）と比較する
	// 希望企業規模がある場合
	if !(len(input.SearchParam.DesiredCompanyScaleTypes) == 0) {
	companyScaleLoop:
		for _, jobSeeker := range jobSeekerListWithDesiredHolidayType {
			if len(jobSeeker.DesiredCompanyScales) == 0 {
				continue
			}

			for _, companyScale := range input.SearchParam.DesiredCompanyScaleTypes {
				if !companyScale.Valid {
					continue
				}
				for _, seekerCompanyScale := range jobSeeker.DesiredCompanyScales {
					if companyScale == seekerCompanyScale.DesiredCompanyScale || seekerCompanyScale.DesiredCompanyScale == null.NewInt(99, true) {
						jobSeekerListWithDesiredCompanyScale = append(jobSeekerListWithDesiredCompanyScale, jobSeeker)
						continue companyScaleLoop
					}
				}
			}
		}
	}

	// 希望企業規模が無い場合
	if len(input.SearchParam.DesiredCompanyScaleTypes) == 0 {
		jobSeekerListWithDesiredCompanyScale = jobSeekerListWithDesiredHolidayType
	}

	fmt.Println("希望企業規模: ", len(jobSeekerListWithDesiredCompanyScale))

	// NOTE: 経験業界（And検索）
	// 経験業界のいずれかが入っている場合
	if !(len(input.SearchParam.ExperienceIndustries) == 0) {

		// 検索パラムのmap
		var experienceIndustrieParams = make(map[null.Int]bool)
		for _, industry := range input.SearchParam.ExperienceIndustries {
			experienceIndustrieParams[industry] = false
		}

	experienceIndustryLoop:
		for _, jobSeeker := range jobSeekerListWithDesiredCompanyScale {
			if len(jobSeeker.WorkHistories) == 0 {
				continue
			}

			/**
			合致した時にexperienceIndustryLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 経験業界のリスト作成
			var experienceIndustries []entity.JobSeekerExperienceIndustry
			for _, workHistory := range jobSeeker.WorkHistories {
				experienceIndustries = append(experienceIndustries, workHistory.ExperienceIndustries...)
			}

			// 求職者の経験業界のmap
			var experienceIndustryDatas = make(map[null.Int]bool)

			for _, industryParam := range input.SearchParam.ExperienceIndustries {
				if !industryParam.Valid {
					continue
				}

				for _, experienceIndustry := range experienceIndustries {
					if industryParam == experienceIndustry.Industry {
						experienceIndustrieParams[industryParam] = true
						experienceIndustryDatas[experienceIndustry.Industry] = true
					}
				}

				if reflect.DeepEqual(experienceIndustrieParams, experienceIndustryDatas) {
					jobSeekerListWithExperienceIndustry = append(jobSeekerListWithExperienceIndustry, jobSeeker)
					continue experienceIndustryLoop
				}
			}
		}
	}

	// 経験業界のいずれかも入っていない場合
	if len(input.SearchParam.ExperienceIndustries) == 0 {
		jobSeekerListWithExperienceIndustry = jobSeekerListWithDesiredCompanyScale
	}

	fmt.Println("経験業界: ", len(jobSeekerListWithExperienceIndustry))

	// NOTE: 経験職種（And検索）
	// 経験職種のいずれかが入っている場合
	if !(len(input.SearchParam.ExperienceOccupations) == 0) {

		// 検索パラムのmap
		var experienceOccupationParams = make(map[null.Int]bool)
		for _, occupation := range input.SearchParam.ExperienceOccupations {
			experienceOccupationParams[occupation] = false
		}

	experienceOccupationLoop:
		for _, jobSeeker := range jobSeekerListWithExperienceIndustry {
			if len(jobSeeker.WorkHistories) == 0 {
				continue
			}

			/**
			合致した時にexperienceOccupationLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 経験職種のリスト作成
			var experienceOccupations []entity.JobSeekerExperienceOccupation
			for _, workHistory := range jobSeeker.WorkHistories {
				if len(workHistory.DepartmentHistories) == 0 {
					continue
				}
				for _, department := range workHistory.DepartmentHistories {
					experienceOccupations = append(experienceOccupations, department.ExperienceOccupations...)
				}
			}

			// 求職者の経験職種のmap
			var experienceOccupationDatas = make(map[null.Int]bool)

			for _, occupationParam := range input.SearchParam.ExperienceOccupations {
				if !occupationParam.Valid {
					continue
				}

				for _, experienceOccupation := range experienceOccupations {
					if occupationParam == experienceOccupation.Occupation {
						experienceOccupationParams[occupationParam] = true
						experienceOccupationDatas[experienceOccupation.Occupation] = true
					}
				}

				if reflect.DeepEqual(experienceOccupationParams, experienceOccupationDatas) {
					jobSeekerListWithExperienceOccupation = append(jobSeekerListWithExperienceOccupation, jobSeeker)
					continue experienceOccupationLoop
				}
			}
		}
	}

	// 経験職種のいずれかも入っていない場合
	if len(input.SearchParam.ExperienceOccupations) == 0 {
		jobSeekerListWithExperienceOccupation = jobSeekerListWithExperienceIndustry
	}

	fmt.Println("経験職種: ", len(jobSeekerListWithExperienceOccupation))

	// NOTE: 社会人経験
	// 社会人経験のいずれかが入っている場合
	if !(len(input.SearchParam.SocialExperiences) == 0) {

	socialExperienceLoop:
		for _, jobSeeker := range jobSeekerListWithExperienceOccupation {
			if len(jobSeeker.WorkHistories) == 0 {
				continue
			}

			/**
			合致した時にsocialExperienceLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/
			for _, socialExperience := range input.SearchParam.SocialExperiences {
				if !socialExperience.Valid {
					continue
				}

				for _, workHistory := range jobSeeker.WorkHistories {
					if socialExperience == workHistory.EmploymentStatus {
						jobSeekerListWithSocialExperience = append(jobSeekerListWithSocialExperience, jobSeeker)
						continue socialExperienceLoop
					}
				}
			}
		}
	}

	// 社会人経験のいずれかも入っていない場合
	if len(input.SearchParam.SocialExperiences) == 0 {
		jobSeekerListWithSocialExperience = jobSeekerListWithExperienceOccupation
	}

	fmt.Println("社会人経験: ", len(jobSeekerListWithSocialExperience))

	// NOTE: マネジメント
	// マネジメントのいずれかが入っている場合
	management, err := strconv.Atoi(input.SearchParam.Management)
	if !(err != nil) {
	managementLoop:
		for _, jobSeeker := range jobSeekerListWithSocialExperience {
			if len(jobSeeker.WorkHistories) == 0 {
				continue
			}
			/**
			合致した時にmanagementLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/
			for _, workHistory := range jobSeeker.WorkHistories {
				if len(workHistory.DepartmentHistories) == 0 {
					continue
				}
				for _, departmentHistory := range workHistory.DepartmentHistories {
					if 0 < departmentHistory.ManagementNumber.Int64 && management == 0 {
						jobSeekerListWithManagement = append(jobSeekerListWithManagement, jobSeeker)
						continue managementLoop
					}
				}
			}
		}
	}

	// マネジメントのいずれかも入っていない場合
	if err != nil {
		jobSeekerListWithManagement = jobSeekerListWithSocialExperience
	}

	fmt.Println("マネジメント: ", len(jobSeekerListWithManagement))

	// NOTE: 資格（And検索）
	// 資格のいずれかが入っている場合
	if !(len(input.SearchParam.Licenses) == 0) {
		// 資格の絞り込み
		jobSeekerListWithLicense = licenceSearchJobSeeker(jobSeekerList, input.SearchParam.Licenses)
	}

	// 資格のいずれかも入っていない場合
	if len(input.SearchParam.Licenses) == 0 {
		jobSeekerListWithLicense = jobSeekerListWithManagement
	}

	fmt.Println("資格: ", len(jobSeekerListWithLicense))

	// NOTE: 語学（And検索）
	// 語学のいずれかが入っている場合
	if !(len(input.SearchParam.Languages) == 0) {

		// 検索パラムのmap
		var languageParams = make(map[null.Int]bool)
		for _, languageParam := range input.SearchParam.Languages {
			languageParams[languageParam] = false
		}

	languageLoop:
		for _, jobSeeker := range jobSeekerListWithLicense {
			if len(jobSeeker.LanguageSkills) == 0 {
				continue
			}

			/**
			合致した時にlanguageLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 求人の語学のmap
			var languageDatas = make(map[null.Int]bool)

			for _, languageParam := range input.SearchParam.Languages {
				if !languageParam.Valid {
					continue
				}

				for _, jobSeekerLanguage := range jobSeeker.LanguageSkills {
					if languageParam == jobSeekerLanguage.LanguageType {
						languageParams[languageParam] = true
						languageDatas[jobSeekerLanguage.LanguageType] = true
					}
				}
				if reflect.DeepEqual(languageParams, languageDatas) {
					jobSeekerListWithLanguage = append(jobSeekerListWithLanguage, jobSeeker)
					continue languageLoop
				}
			}
		}
	}

	// 語学のいずれかも入っていない場合
	if len(input.SearchParam.Languages) == 0 {
		jobSeekerListWithLanguage = jobSeekerListWithLicense
	}

	fmt.Println("語学: ", len(jobSeekerListWithLanguage))

	// NOTE: 必要PCスキル（excel）
	// 必要PCスキル（excel）のいずれかが入っている場合
	if !(len(input.SearchParam.ExcelSkills) == 0) {
	excelSkillLoop:
		for _, jobSeeker := range jobSeekerListWithLanguage {
			if !jobSeeker.ExcelSkill.Valid {
				continue
			}

			for _, excelSkill := range input.SearchParam.ExcelSkills {
				if !excelSkill.Valid {
					continue
				}

				if excelSkill == jobSeeker.ExcelSkill {
					/**
					必要PCスキル（excel）が合致したタイミングで次の求職者へ移る
					ダブりが発生しないように制御
					*/
					jobSeekerListWithExcelSkill = append(jobSeekerListWithExcelSkill, jobSeeker)
					continue excelSkillLoop
				}
			}
		}
	}

	// 必要PCスキル（excel）のいずれかも入っていない場合
	if len(input.SearchParam.ExcelSkills) == 0 {
		jobSeekerListWithExcelSkill = jobSeekerListWithLanguage
	}

	fmt.Println("必要PCスキル（excel）: ", len(jobSeekerListWithExcelSkill))

	// NOTE: 必要PCスキル（word）
	// 必要PCスキル（word）のいずれかが入っている場合
	if !(len(input.SearchParam.WordSkills) == 0) {
	wordSkillLoop:
		for _, jobSeeker := range jobSeekerListWithExcelSkill {
			if !jobSeeker.WordSkill.Valid {
				continue
			}

			for _, wordSkill := range input.SearchParam.WordSkills {
				if !wordSkill.Valid {
					continue
				}

				if wordSkill == jobSeeker.WordSkill {
					/**
					必要PCスキル（word）が合致したタイミングで次の求人へ移る
					ダブりが発生しないように制御
					*/
					jobSeekerListWithWordSkill = append(jobSeekerListWithWordSkill, jobSeeker)
					continue wordSkillLoop
				}
			}
		}
	}

	// 必要PCスキル（word）のいずれかも入っていない場合
	if len(input.SearchParam.WordSkills) == 0 {
		jobSeekerListWithWordSkill = jobSeekerListWithExcelSkill
	}

	fmt.Println("必要PCスキル（word）: ", len(jobSeekerListWithWordSkill))

	// NOTE: 必要PCスキル（powerpoint）
	// 必要PCスキル（powerpoint）のいずれかが入っている場合
	if !(len(input.SearchParam.PowerPointSkills) == 0) {
	powerPointSkillLoop:
		for _, jobSeeker := range jobSeekerListWithWordSkill {
			if !jobSeeker.PowerPointSkill.Valid {
				continue
			}

			for _, powerPointSkill := range input.SearchParam.PowerPointSkills {
				if !powerPointSkill.Valid {
					continue
				}

				if powerPointSkill == jobSeeker.PowerPointSkill {
					/**
					必要PCスキル（powerpoint）が合致したタイミングで次の求人へ移る
					ダブりが発生しないように制御
					*/
					jobSeekerListWithPowerPointSkill = append(jobSeekerListWithPowerPointSkill, jobSeeker)
					continue powerPointSkillLoop
				}
			}
		}
	}

	// 必要PCスキル（powerpoint）のいずれかも入っていない場合
	if len(input.SearchParam.PowerPointSkills) == 0 {
		jobSeekerListWithPowerPointSkill = jobSeekerListWithWordSkill
	}

	fmt.Println("必要PCスキル（powerpoint）: ", len(jobSeekerListWithPowerPointSkill))

	// NOTE: 業務ツール（And検索）
	// 業務ツールのいずれかが入っている場合
	if !(len(input.SearchParam.AnotherPCSkills) == 0) {
		// 検索パラムのmap
		var pcToolParams = make(map[null.Int]bool)
		for _, pcToolParam := range input.SearchParam.AnotherPCSkills {
			pcToolParams[pcToolParam] = false
		}

	pcToolLoop:
		for _, jobSeeker := range jobSeekerListWithPowerPointSkill {
			if len(jobSeeker.PCTools) == 0 {
				continue
			}

			/**
			合致した時にpcToolLoopを抜けるためのラベル
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 求人の業務ツール経験のmap
			var pcToolDatas = make(map[null.Int]bool)

			for _, pcToolParam := range input.SearchParam.AnotherPCSkills {
				if !pcToolParam.Valid {
					continue
				}
				for _, jobSeekerPCTool := range jobSeeker.PCTools {
					if !jobSeekerPCTool.Tool.Valid {
						continue
					}

					if pcToolParam == jobSeekerPCTool.Tool {
						pcToolParams[pcToolParam] = true
						pcToolDatas[jobSeekerPCTool.Tool] = true
					}
				}

				if reflect.DeepEqual(pcToolParams, pcToolDatas) {
					jobSeekerListWithAnotherPCSkill = append(jobSeekerListWithAnotherPCSkill, jobSeeker)
					continue pcToolLoop
				}
			}
		}
	}

	// 必要PCスキル（その他）のいずれかも入っていない場合
	if len(input.SearchParam.AnotherPCSkills) == 0 {
		jobSeekerListWithAnotherPCSkill = jobSeekerListWithPowerPointSkill
	}

	fmt.Println("業務ツール: ", len(jobSeekerListWithAnotherPCSkill))

	// NOTE: 開発言語
	// 開発言語のいずれかが入っている場合
	if !(len(input.SearchParam.DevelopmentLanguages) == 0) {

		// 検索パラムのmap
		var developmentLanguageParams = make(map[null.Int]bool)
		for _, developmentLanguageParam := range input.SearchParam.DevelopmentLanguages {
			developmentLanguageParams[developmentLanguageParam] = false
		}

	developmentLanguageLoop:
		for _, jobSeeker := range jobSeekerListWithAnotherPCSkill {
			if len(jobSeeker.DevelopmentSkills) == 0 {
				continue
			}

			/**
			合致した時にdevelopmentLanguageLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 求人の必要開発言語のmap
			var developmentLanguageDatas = make(map[null.Int]bool)

			for _, devLanguageParam := range input.SearchParam.DevelopmentLanguages {
				if !devLanguageParam.Valid {
					continue
				}

				for _, jobSeekerDevelopment := range jobSeeker.DevelopmentSkills {
					if !jobSeekerDevelopment.DevelopmentCategory.Valid {
						continue
					}

					if jobSeekerDevelopment.DevelopmentCategory.Int64 == 0 && devLanguageParam == jobSeekerDevelopment.DevelopmentType {
						developmentLanguageParams[devLanguageParam] = true
						developmentLanguageDatas[jobSeekerDevelopment.DevelopmentType] = true
					}
				}

				if reflect.DeepEqual(developmentLanguageParams, developmentLanguageDatas) {
					jobSeekerListWithDevelopmentLanguage = append(jobSeekerListWithDevelopmentLanguage, jobSeeker)
					continue developmentLanguageLoop
				}
			}
		}
	}

	// 開発言語のいずれかも入っていない場合
	if len(input.SearchParam.DevelopmentLanguages) == 0 {
		jobSeekerListWithDevelopmentLanguage = jobSeekerListWithAnotherPCSkill
	}

	fmt.Println("開発言語: ", len(jobSeekerListWithDevelopmentLanguage))

	// NOTE: 開発OS
	// 開発OSのいずれかが入っている場合
	if !(len(input.SearchParam.DevelopmentOS) == 0) {

		// 検索パラムのmap
		var developmentOSParams = make(map[null.Int]bool)
		for _, developmentOSParam := range input.SearchParam.DevelopmentOS {
			developmentOSParams[developmentOSParam] = false
		}

	developmentOSLoop:
		for _, jobSeeker := range jobSeekerListWithDevelopmentLanguage {
			if len(jobSeeker.DevelopmentSkills) == 0 {
				continue
			}

			/**
			合致した時にdevelopmentOSLoopを抜ける
			絞り込みと比較する値がどちらもスライスの場合は必要
			*/

			// 求人の必要開発OSのmap
			var developmentOSDatas = make(map[null.Int]bool)

			for _, devOSParam := range input.SearchParam.DevelopmentOS {
				if !devOSParam.Valid {
					continue
				}

				for _, jobSeekerDevelopment := range jobSeeker.DevelopmentSkills {
					if !jobSeekerDevelopment.DevelopmentCategory.Valid {
						continue
					}

					if jobSeekerDevelopment.DevelopmentCategory.Int64 == 1 && devOSParam == jobSeekerDevelopment.DevelopmentType {
						developmentOSParams[devOSParam] = true
						developmentOSDatas[jobSeekerDevelopment.DevelopmentType] = true
					}
				}
				if reflect.DeepEqual(developmentOSParams, developmentOSDatas) {
					jobSeekerListWithDevelopmentOS = append(jobSeekerListWithDevelopmentOS, jobSeeker)
					continue developmentOSLoop
				}
			}
		}
	}

	// 開発OSのいずれかも入っていない場合
	if len(input.SearchParam.DevelopmentOS) == 0 {
		jobSeekerListWithDevelopmentOS = jobSeekerListWithDevelopmentLanguage
	}

	fmt.Println("開発OS: ", len(jobSeekerListWithDevelopmentOS))

	// Note: アピアランス
	// アピアランスがある場合
	if !(len(input.SearchParam.AppearanceTypes) == 0) {
	appearanceLoop:
		for _, jobSeeker := range jobSeekerListWithDevelopmentOS {
			for _, appearance := range input.SearchParam.AppearanceTypes {
				if !appearance.Valid {
					continue
				}
				if appearance == jobSeeker.Appearance {
					// アピアランスが一つでも合致していればcontinueで次の求職者へ
					jobSeekerListWithAppearance = append(jobSeekerListWithAppearance, jobSeeker)
					continue appearanceLoop
				}
			}
		}
	}

	// アピアランスが無い場合
	if len(input.SearchParam.AppearanceTypes) == 0 {
		jobSeekerListWithAppearance = jobSeekerListWithDevelopmentOS
	}

	fmt.Println("アピアランス: ", len(jobSeekerListWithAppearance))

	// Note: コミュ力
	// コミュ力がある場合
	if !(len(input.SearchParam.CommunicationTypes) == 0) {
	communicationLoop:
		for _, jobSeeker := range jobSeekerListWithAppearance {
			for _, communication := range input.SearchParam.CommunicationTypes {
				if !communication.Valid {
					continue
				}
				if communication == jobSeeker.Communication {
					// コミュ力が一つでも合致していればcontinueで次の求職者へ
					jobSeekerListWithCommunication = append(jobSeekerListWithCommunication, jobSeeker)
					continue communicationLoop
				}
			}
		}
	}

	// コミュ力が無い場合
	if len(input.SearchParam.CommunicationTypes) == 0 {
		jobSeekerListWithCommunication = jobSeekerListWithAppearance
	}

	fmt.Println("コミュ力: ", len(jobSeekerListWithCommunication))

	// Note: 論理的思考力
	// 論理的思考力がある場合
	if !(len(input.SearchParam.ThinkingTypes) == 0) {
	thinkingLoop:
		for _, jobInformation := range jobSeekerListWithCommunication {
			for _, thinking := range input.SearchParam.ThinkingTypes {
				if !thinking.Valid {
					continue
				}
				if thinking == jobInformation.Thinking {
					// 論理的思考力が一つでも合致していればcontinueで次の求職者へ
					jobSeekerListWithThinking = append(jobSeekerListWithThinking, jobInformation)
					continue thinkingLoop
				}
			}
		}
	}

	// 論理的思考力が無い場合
	if len(input.SearchParam.ThinkingTypes) == 0 {
		jobSeekerListWithThinking = jobSeekerListWithCommunication
	}

	fmt.Println("論理的思考力: ", len(jobSeekerListWithThinking))

	return jobSeekerListWithThinking, nil
}

func licenceSearchJobSeeker(jobSeekerList []*entity.JobSeeker, licenseParams []null.Int) []*entity.JobSeeker {

	var (
		seekerListAtLicense []*entity.JobSeeker
	)
	// ライセンスのマッチング
licenseLoop:
	for _, seeker := range jobSeekerList {
		for _, licenseParam := range licenseParams {
			for _, seekerLicense := range seeker.Licenses {
				if !licenseParam.Valid {
					seekerListAtLicense = append(seekerListAtLicense, seeker)
					continue licenseLoop
					// 完全一致 ||
					// 求職者の保有資格：普通自動車免許（MT） && 求人の必要資格：普通自動車免許（MT）＋普通自動車免許（AT）がヒット
				} else if (seekerLicense.LicenseType == licenseParam) ||
					(seekerLicense.LicenseType.Int64 == 4803 && licenseParam.Int64 == 4805) ||
					(seekerLicense.LicenseType.Int64 == 1205 && (licenseParam.Int64 == 1203 || licenseParam.Int64 == 1204)) ||
					(seekerLicense.LicenseType.Int64 == 1204 && licenseParam.Int64 == 1205) ||
					(seekerLicense.LicenseType.Int64 == 1206 && licenseParam.Int64 == 1207) ||
					(seekerLicense.LicenseType.Int64 == 1211 && licenseParam.Int64 == 1212) ||
					(seekerLicense.LicenseType.Int64 == 1218 && licenseParam.Int64 == 1219) ||
					(seekerLicense.LicenseType.Int64 == 1224 && (licenseParam.Int64 == 1223 || licenseParam.Int64 == 1222)) ||
					(seekerLicense.LicenseType.Int64 == 1223 && licenseParam.Int64 == 1224) ||
					(seekerLicense.LicenseType.Int64 == 1238 && licenseParam.Int64 == 1239) ||
					(seekerLicense.LicenseType.Int64 == 1301 && licenseParam.Int64 == 1302) ||
					(seekerLicense.LicenseType.Int64 == 1305 && licenseParam.Int64 == 1306) ||
					(seekerLicense.LicenseType.Int64 == 1315 && licenseParam.Int64 == 1316) ||
					(seekerLicense.LicenseType.Int64 == 1322 && (licenseParam.Int64 == 1321 || licenseParam.Int64 == 1320)) ||
					(seekerLicense.LicenseType.Int64 == 1321 && licenseParam.Int64 == 1322) ||
					(seekerLicense.LicenseType.Int64 == 1326 && (licenseParam.Int64 == 1324 || licenseParam.Int64 == 1325 || licenseParam.Int64 == 1323)) ||
					(seekerLicense.LicenseType.Int64 == 1325 && (licenseParam.Int64 == 1324 || licenseParam.Int64 == 1323)) ||
					(seekerLicense.LicenseType.Int64 == 1324 && licenseParam.Int64 == 1323) ||
					(seekerLicense.LicenseType.Int64 == 1401 && licenseParam.Int64 == 1402) ||
					(seekerLicense.LicenseType.Int64 == 1404 && licenseParam.Int64 == 1405) ||
					(seekerLicense.LicenseType.Int64 == 1406 && licenseParam.Int64 == 1407) ||
					(seekerLicense.LicenseType.Int64 == 1408 && licenseParam.Int64 == 1409) ||
					(seekerLicense.LicenseType.Int64 == 1547 && (licenseParam.Int64 == 1546 || licenseParam.Int64 == 1545)) ||
					(seekerLicense.LicenseType.Int64 == 1546 && licenseParam.Int64 == 1547) ||
					(seekerLicense.LicenseType.Int64 == 1550 && (licenseParam.Int64 == 1549 || licenseParam.Int64 == 1548)) ||
					(seekerLicense.LicenseType.Int64 == 1549 && licenseParam.Int64 == 1550) ||
					(seekerLicense.LicenseType.Int64 == 1551 && licenseParam.Int64 == 1552) ||
					(seekerLicense.LicenseType.Int64 == 1553 && licenseParam.Int64 == 1554) ||
					(seekerLicense.LicenseType.Int64 == 1563 && licenseParam.Int64 == 1564) ||
					(seekerLicense.LicenseType.Int64 == 1605 && licenseParam.Int64 == 1606) ||
					(seekerLicense.LicenseType.Int64 == 1610 && licenseParam.Int64 == 1611) ||
					(seekerLicense.LicenseType.Int64 == 2202 && licenseParam.Int64 == 2203) ||
					(seekerLicense.LicenseType.Int64 == 2306 && (licenseParam.Int64 == 2305 || licenseParam.Int64 == 2304)) ||
					(seekerLicense.LicenseType.Int64 == 2305 && licenseParam.Int64 == 2306) ||
					(seekerLicense.LicenseType.Int64 == 2314 && licenseParam.Int64 == 2313) ||
					(seekerLicense.LicenseType.Int64 == 2409 && (licenseParam.Int64 == 2408 || licenseParam.Int64 == 2407)) ||
					(seekerLicense.LicenseType.Int64 == 2408 && licenseParam.Int64 == 2409) ||
					(seekerLicense.LicenseType.Int64 == 2517 && (licenseParam.Int64 == 2516 || licenseParam.Int64 == 2515)) ||
					(seekerLicense.LicenseType.Int64 == 2516 && licenseParam.Int64 == 2517) ||
					(seekerLicense.LicenseType.Int64 == 2701 && licenseParam.Int64 == 2702) ||
					(seekerLicense.LicenseType.Int64 == 2703 && licenseParam.Int64 == 2704) ||
					(seekerLicense.LicenseType.Int64 == 2902 && licenseParam.Int64 == 2901) ||
					(seekerLicense.LicenseType.Int64 == 2901 && (licenseParam.Int64 == 2902 || licenseParam.Int64 == 2903)) ||
					(seekerLicense.LicenseType.Int64 == 2905 && licenseParam.Int64 == 2904) ||
					(seekerLicense.LicenseType.Int64 == 2904 && (licenseParam.Int64 == 2905 || licenseParam.Int64 == 2906)) ||
					(seekerLicense.LicenseType.Int64 == 2908 && licenseParam.Int64 == 2907) ||
					(seekerLicense.LicenseType.Int64 == 2907 && (licenseParam.Int64 == 2908 || licenseParam.Int64 == 2909)) ||
					(seekerLicense.LicenseType.Int64 == 2913 && licenseParam.Int64 == 2912) ||
					(seekerLicense.LicenseType.Int64 == 2916 && licenseParam.Int64 == 2917) ||
					(seekerLicense.LicenseType.Int64 == 3001 && licenseParam.Int64 == 3002) ||
					(seekerLicense.LicenseType.Int64 == 3009 && licenseParam.Int64 == 3010) ||
					(seekerLicense.LicenseType.Int64 == 3014 && (licenseParam.Int64 == 3013 || licenseParam.Int64 == 3012)) ||
					(seekerLicense.LicenseType.Int64 == 3013 && licenseParam.Int64 == 3014) ||
					(seekerLicense.LicenseType.Int64 == 3017 && (licenseParam.Int64 == 3016 || licenseParam.Int64 == 3015)) ||
					(seekerLicense.LicenseType.Int64 == 3016 && licenseParam.Int64 == 3017) ||
					(seekerLicense.LicenseType.Int64 == 3107 && licenseParam.Int64 == 3108) ||
					(seekerLicense.LicenseType.Int64 == 3201 && licenseParam.Int64 == 3202) ||
					(seekerLicense.LicenseType.Int64 == 3309 && (licenseParam.Int64 == 3308 || licenseParam.Int64 == 3307)) ||
					(seekerLicense.LicenseType.Int64 == 3308 && licenseParam.Int64 == 3309) ||
					(seekerLicense.LicenseType.Int64 == 3316 && (licenseParam.Int64 == 3314 || licenseParam.Int64 == 3315 || licenseParam.Int64 == 3313)) ||
					(seekerLicense.LicenseType.Int64 == 3315 && (licenseParam.Int64 == 3314 || licenseParam.Int64 == 3313)) ||
					(seekerLicense.LicenseType.Int64 == 3314 && licenseParam.Int64 == 3313) ||
					(seekerLicense.LicenseType.Int64 == 3325 && licenseParam.Int64 == 3326) ||
					(seekerLicense.LicenseType.Int64 == 3327 && licenseParam.Int64 == 3328) ||
					(seekerLicense.LicenseType.Int64 == 3329 && licenseParam.Int64 == 3330) ||
					(seekerLicense.LicenseType.Int64 == 3331 && licenseParam.Int64 == 3332) ||
					(seekerLicense.LicenseType.Int64 == 3337 && (licenseParam.Int64 == 3336 || licenseParam.Int64 == 3335)) ||
					(seekerLicense.LicenseType.Int64 == 3336 && licenseParam.Int64 == 3337) ||
					(seekerLicense.LicenseType.Int64 == 3338 && licenseParam.Int64 == 3339) ||
					(seekerLicense.LicenseType.Int64 == 3403 && (licenseParam.Int64 == 3402 || licenseParam.Int64 == 3401)) ||
					(seekerLicense.LicenseType.Int64 == 3402 && licenseParam.Int64 == 3403) ||
					(seekerLicense.LicenseType.Int64 == 3407 && (licenseParam.Int64 == 3406 || licenseParam.Int64 == 3405)) ||
					(seekerLicense.LicenseType.Int64 == 3406 && licenseParam.Int64 == 3407) ||
					(seekerLicense.LicenseType.Int64 == 3515 && (licenseParam.Int64 == 3514 || licenseParam.Int64 == 3513)) ||
					(seekerLicense.LicenseType.Int64 == 3514 && licenseParam.Int64 == 3515) ||
					(seekerLicense.LicenseType.Int64 == 3523 && (licenseParam.Int64 == 3522 || licenseParam.Int64 == 3521)) ||
					(seekerLicense.LicenseType.Int64 == 3522 && licenseParam.Int64 == 3523) ||
					(seekerLicense.LicenseType.Int64 == 3527 && (licenseParam.Int64 == 3525 || licenseParam.Int64 == 3526 || licenseParam.Int64 == 3524)) ||
					(seekerLicense.LicenseType.Int64 == 3526 && (licenseParam.Int64 == 3524 || licenseParam.Int64 == 3525)) ||
					(seekerLicense.LicenseType.Int64 == 3525 && licenseParam.Int64 == 3524) ||
					(seekerLicense.LicenseType.Int64 == 3617 && (licenseParam.Int64 == 3616 || licenseParam.Int64 == 3615)) ||
					(seekerLicense.LicenseType.Int64 == 3616 && licenseParam.Int64 == 3617) ||
					(seekerLicense.LicenseType.Int64 == 3621 && (licenseParam.Int64 == 3620 || licenseParam.Int64 == 3619)) ||
					(seekerLicense.LicenseType.Int64 == 3620 && licenseParam.Int64 == 3621) ||
					(seekerLicense.LicenseType.Int64 == 3625 && (licenseParam.Int64 == 3624 || licenseParam.Int64 == 3623)) ||
					(seekerLicense.LicenseType.Int64 == 3624 && licenseParam.Int64 == 3625) ||
					(seekerLicense.LicenseType.Int64 == 3629 && (licenseParam.Int64 == 3627 || licenseParam.Int64 == 3628 || licenseParam.Int64 == 3626)) ||
					(seekerLicense.LicenseType.Int64 == 3628 && (licenseParam.Int64 == 3627 || licenseParam.Int64 == 3626)) ||
					(seekerLicense.LicenseType.Int64 == 3627 && licenseParam.Int64 == 3626) ||
					(seekerLicense.LicenseType.Int64 == 3632 && (licenseParam.Int64 == 3631 || licenseParam.Int64 == 3630)) ||
					(seekerLicense.LicenseType.Int64 == 3631 && licenseParam.Int64 == 3632) ||
					(seekerLicense.LicenseType.Int64 == 3636 && (licenseParam.Int64 == 3635 || licenseParam.Int64 == 3634)) ||
					(seekerLicense.LicenseType.Int64 == 3635 && licenseParam.Int64 == 3636) ||
					(seekerLicense.LicenseType.Int64 == 3708 && licenseParam.Int64 == 3709) ||
					(seekerLicense.LicenseType.Int64 == 3722 && (licenseParam.Int64 == 3721 || licenseParam.Int64 == 3720)) ||
					(seekerLicense.LicenseType.Int64 == 3721 && licenseParam.Int64 == 3722) ||
					(seekerLicense.LicenseType.Int64 == 3730 && (licenseParam.Int64 == 3729 || licenseParam.Int64 == 3728)) ||
					(seekerLicense.LicenseType.Int64 == 3729 && licenseParam.Int64 == 3730) ||
					(seekerLicense.LicenseType.Int64 == 3801 && licenseParam.Int64 == 3802) ||
					(seekerLicense.LicenseType.Int64 == 3814 && (licenseParam.Int64 == 3813 || licenseParam.Int64 == 3812)) ||
					(seekerLicense.LicenseType.Int64 == 3813 && licenseParam.Int64 == 3814) ||
					(seekerLicense.LicenseType.Int64 == 3815 && licenseParam.Int64 == 3816) ||
					(seekerLicense.LicenseType.Int64 == 3817 && licenseParam.Int64 == 3818) ||
					(seekerLicense.LicenseType.Int64 == 3843 && (licenseParam.Int64 == 3842 || licenseParam.Int64 == 3841)) ||
					(seekerLicense.LicenseType.Int64 == 3842 && licenseParam.Int64 == 3843) ||
					(seekerLicense.LicenseType.Int64 == 3846 && (licenseParam.Int64 == 3845 || licenseParam.Int64 == 3844)) ||
					(seekerLicense.LicenseType.Int64 == 3845 && licenseParam.Int64 == 3846) ||
					(seekerLicense.LicenseType.Int64 == 3851 && (licenseParam.Int64 == 3850 || licenseParam.Int64 == 3849)) ||
					(seekerLicense.LicenseType.Int64 == 3850 && licenseParam.Int64 == 3851) ||
					(seekerLicense.LicenseType.Int64 == 3854 && (licenseParam.Int64 == 3853 || licenseParam.Int64 == 3852)) ||
					(seekerLicense.LicenseType.Int64 == 3853 && licenseParam.Int64 == 3854) ||
					(seekerLicense.LicenseType.Int64 == 3858 && (licenseParam.Int64 == 3856 || licenseParam.Int64 == 3857 || licenseParam.Int64 == 3855)) ||
					(seekerLicense.LicenseType.Int64 == 3857 && (licenseParam.Int64 == 3856 || licenseParam.Int64 == 3855)) ||
					(seekerLicense.LicenseType.Int64 == 3856 && licenseParam.Int64 == 3855) ||
					(seekerLicense.LicenseType.Int64 == 3862 && (licenseParam.Int64 == 3860 || licenseParam.Int64 == 3861 || licenseParam.Int64 == 3859)) ||
					(seekerLicense.LicenseType.Int64 == 3861 && (licenseParam.Int64 == 3860 || licenseParam.Int64 == 3859)) ||
					(seekerLicense.LicenseType.Int64 == 3860 && licenseParam.Int64 == 3859) ||
					(seekerLicense.LicenseType.Int64 == 3866 && (licenseParam.Int64 == 3864 || licenseParam.Int64 == 3865 || licenseParam.Int64 == 3863)) ||
					(seekerLicense.LicenseType.Int64 == 3865 && (licenseParam.Int64 == 3864 || licenseParam.Int64 == 3863)) ||
					(seekerLicense.LicenseType.Int64 == 3864 && licenseParam.Int64 == 3863) ||
					(seekerLicense.LicenseType.Int64 == 3875 && (licenseParam.Int64 == 3874 || licenseParam.Int64 == 3873)) ||
					(seekerLicense.LicenseType.Int64 == 3874 && licenseParam.Int64 == 3875) ||
					(seekerLicense.LicenseType.Int64 == 3884 && (licenseParam.Int64 == 3883 || licenseParam.Int64 == 3882)) ||
					(seekerLicense.LicenseType.Int64 == 3883 && licenseParam.Int64 == 3884) ||
					(seekerLicense.LicenseType.Int64 == 3903 && (licenseParam.Int64 == 3902 || licenseParam.Int64 == 3901)) ||
					(seekerLicense.LicenseType.Int64 == 3902 && licenseParam.Int64 == 3903) ||
					(seekerLicense.LicenseType.Int64 == 3906 && (licenseParam.Int64 == 3905 || licenseParam.Int64 == 3904)) ||
					(seekerLicense.LicenseType.Int64 == 3905 && licenseParam.Int64 == 3906) ||
					(seekerLicense.LicenseType.Int64 == 4103 && (licenseParam.Int64 == 4102 || licenseParam.Int64 == 4101)) ||
					(seekerLicense.LicenseType.Int64 == 4102 && licenseParam.Int64 == 4103) ||
					(seekerLicense.LicenseType.Int64 == 4104 && licenseParam.Int64 == 4105) ||
					(seekerLicense.LicenseType.Int64 == 4113 && (licenseParam.Int64 == 4112 || licenseParam.Int64 == 4111)) ||
					(seekerLicense.LicenseType.Int64 == 4112 && licenseParam.Int64 == 4113) ||
					(seekerLicense.LicenseType.Int64 == 4311 && (licenseParam.Int64 == 4310 || licenseParam.Int64 == 4309)) ||
					(seekerLicense.LicenseType.Int64 == 4310 && licenseParam.Int64 == 4311) ||
					(seekerLicense.LicenseType.Int64 == 4319 && (licenseParam.Int64 == 4318 || licenseParam.Int64 == 4317)) ||
					(seekerLicense.LicenseType.Int64 == 4318 && licenseParam.Int64 == 4319) ||
					(seekerLicense.LicenseType.Int64 == 4320 && licenseParam.Int64 == 4321) ||
					(seekerLicense.LicenseType.Int64 == 4323 && licenseParam.Int64 == 4324) ||
					(seekerLicense.LicenseType.Int64 == 4325 && licenseParam.Int64 == 4326) ||
					(seekerLicense.LicenseType.Int64 == 4330 && (licenseParam.Int64 == 4329 || licenseParam.Int64 == 4328)) ||
					(seekerLicense.LicenseType.Int64 == 4329 && licenseParam.Int64 == 4330) ||
					(seekerLicense.LicenseType.Int64 == 4334 && licenseParam.Int64 == 4335) ||
					(seekerLicense.LicenseType.Int64 == 4338 && (licenseParam.Int64 == 4337 || licenseParam.Int64 == 4336)) ||
					(seekerLicense.LicenseType.Int64 == 4337 && licenseParam.Int64 == 4338) ||
					(seekerLicense.LicenseType.Int64 == 4406 && (licenseParam.Int64 == 4405 || licenseParam.Int64 == 4404)) ||
					(seekerLicense.LicenseType.Int64 == 4405 && licenseParam.Int64 == 4406) ||
					(seekerLicense.LicenseType.Int64 == 4412 && (licenseParam.Int64 == 4411 || licenseParam.Int64 == 4410)) ||
					(seekerLicense.LicenseType.Int64 == 4411 && licenseParam.Int64 == 4412) ||
					(seekerLicense.LicenseType.Int64 == 4501 && licenseParam.Int64 == 4502) ||
					(seekerLicense.LicenseType.Int64 == 4503 && licenseParam.Int64 == 4504) ||
					(seekerLicense.LicenseType.Int64 == 4505 && licenseParam.Int64 == 4506) ||
					(seekerLicense.LicenseType.Int64 == 4507 && licenseParam.Int64 == 4508) ||
					(seekerLicense.LicenseType.Int64 == 4510 && licenseParam.Int64 == 4511) ||
					(seekerLicense.LicenseType.Int64 == 4512 && licenseParam.Int64 == 4513) ||
					(seekerLicense.LicenseType.Int64 == 4514 && licenseParam.Int64 == 4515) ||
					(seekerLicense.LicenseType.Int64 == 4516 && licenseParam.Int64 == 4517) ||
					(seekerLicense.LicenseType.Int64 == 4518 && licenseParam.Int64 == 4519) ||
					(seekerLicense.LicenseType.Int64 == 4520 && licenseParam.Int64 == 4521) ||
					(seekerLicense.LicenseType.Int64 == 4603 && (licenseParam.Int64 == 4602 || licenseParam.Int64 == 4601)) ||
					(seekerLicense.LicenseType.Int64 == 4602 && licenseParam.Int64 == 4603) ||
					(seekerLicense.LicenseType.Int64 == 4801 && licenseParam.Int64 == 4823) ||
					(seekerLicense.LicenseType.Int64 == 4803 && (licenseParam.Int64 == 4801 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4823 || licenseParam.Int64 == 4803 || licenseParam.Int64 == 4804)) ||
					(seekerLicense.LicenseType.Int64 == 4804 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4821 || licenseParam.Int64 == 4823 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4806 || licenseParam.Int64 == 4807 || licenseParam.Int64 == 4808 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4803 || licenseParam.Int64 == 4805 || licenseParam.Int64 == 4809 || licenseParam.Int64 == 4810)) ||
					(seekerLicense.LicenseType.Int64 == 4806 && licenseParam.Int64 == 4807) ||
					(seekerLicense.LicenseType.Int64 == 4808 && licenseParam.Int64 == 4898) ||
					(seekerLicense.LicenseType.Int64 == 4809 && (licenseParam.Int64 == 4807 || licenseParam.Int64 == 4806)) ||
					(seekerLicense.LicenseType.Int64 == 4810 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4821 || licenseParam.Int64 == 4823 || licenseParam.Int64 == 4804 || licenseParam.Int64 == 4805 || licenseParam.Int64 == 4810 || licenseParam.Int64 == 4807 || licenseParam.Int64 == 4806 || licenseParam.Int64 == 4898 || licenseParam.Int64 == 4808 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4811)) ||
					(seekerLicense.LicenseType.Int64 == 4812 && licenseParam.Int64 == 4811) ||
					(seekerLicense.LicenseType.Int64 == 4817 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818)) ||
					(seekerLicense.LicenseType.Int64 == 4818 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801)) ||
					(seekerLicense.LicenseType.Int64 == 4819 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819)) ||
					(seekerLicense.LicenseType.Int64 == 4820 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820)) ||
					(seekerLicense.LicenseType.Int64 == 4821 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4822 || licenseParam.Int64 == 4823)) ||
					(seekerLicense.LicenseType.Int64 == 4822 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4822)) ||
					(seekerLicense.LicenseType.Int64 == 4823 && (licenseParam.Int64 == 4802 || licenseParam.Int64 == 4801 || licenseParam.Int64 == 4818 || licenseParam.Int64 == 4817 || licenseParam.Int64 == 4820 || licenseParam.Int64 == 4821 || licenseParam.Int64 == 4819 || licenseParam.Int64 == 4822)) ||
					(seekerLicense.LicenseType.Int64 == 4907 && licenseParam.Int64 == 4908) ||
					(seekerLicense.LicenseType.Int64 == 5003 && (licenseParam.Int64 == 5002 || licenseParam.Int64 == 5001)) ||
					(seekerLicense.LicenseType.Int64 == 5002 && licenseParam.Int64 == 5003) ||
					(seekerLicense.LicenseType.Int64 == 5004 && licenseParam.Int64 == 5005) ||
					(seekerLicense.LicenseType.Int64 == 5006 && licenseParam.Int64 == 5007) ||
					(seekerLicense.LicenseType.Int64 == 5010 && (licenseParam.Int64 == 5009 || licenseParam.Int64 == 5008)) ||
					(seekerLicense.LicenseType.Int64 == 5009 && licenseParam.Int64 == 5010) ||
					(seekerLicense.LicenseType.Int64 == 5103 && licenseParam.Int64 == 5104) ||
					(seekerLicense.LicenseType.Int64 == 5107 && (licenseParam.Int64 == 5106 || licenseParam.Int64 == 5105)) ||
					(seekerLicense.LicenseType.Int64 == 5106 && licenseParam.Int64 == 5107) ||
					(seekerLicense.LicenseType.Int64 == 5108 && licenseParam.Int64 == 5109) ||
					(seekerLicense.LicenseType.Int64 == 5113 && (licenseParam.Int64 == 5111 || licenseParam.Int64 == 5112 || licenseParam.Int64 == 5110)) ||
					(seekerLicense.LicenseType.Int64 == 5112 && (licenseParam.Int64 == 5110 || licenseParam.Int64 == 5111)) ||
					(seekerLicense.LicenseType.Int64 == 5111 && licenseParam.Int64 == 5110) ||
					(seekerLicense.LicenseType.Int64 == 5114 && licenseParam.Int64 == 5115) ||
					(seekerLicense.LicenseType.Int64 == 5116 && licenseParam.Int64 == 5117) ||
					(seekerLicense.LicenseType.Int64 == 5121 && (licenseParam.Int64 == 5119 || licenseParam.Int64 == 5120 || licenseParam.Int64 == 5118)) ||
					(seekerLicense.LicenseType.Int64 == 5120 && (licenseParam.Int64 == 5119 || licenseParam.Int64 == 5118)) ||
					(seekerLicense.LicenseType.Int64 == 5119 && licenseParam.Int64 == 5118) ||
					(seekerLicense.LicenseType.Int64 == 5204 && licenseParam.Int64 == 5205) ||
					(seekerLicense.LicenseType.Int64 == 5208 && (licenseParam.Int64 == 5207 || licenseParam.Int64 == 5206)) ||
					(seekerLicense.LicenseType.Int64 == 5207 && licenseParam.Int64 == 5208) ||
					(seekerLicense.LicenseType.Int64 == 5303 && licenseParam.Int64 == 5304) ||
					(seekerLicense.LicenseType.Int64 == 5307 && licenseParam.Int64 == 5308) ||
					(seekerLicense.LicenseType.Int64 == 5404 && (licenseParam.Int64 == 5402 || licenseParam.Int64 == 5403 || licenseParam.Int64 == 5401)) ||
					(seekerLicense.LicenseType.Int64 == 5403 && (licenseParam.Int64 == 5401 || licenseParam.Int64 == 5402)) ||
					(seekerLicense.LicenseType.Int64 == 5402 && licenseParam.Int64 == 5401) ||
					(seekerLicense.LicenseType.Int64 == 5407 && (licenseParam.Int64 == 5406 || licenseParam.Int64 == 5405)) ||
					(seekerLicense.LicenseType.Int64 == 5406 && licenseParam.Int64 == 5407) ||
					(seekerLicense.LicenseType.Int64 == 5408 && licenseParam.Int64 == 5409) ||
					(seekerLicense.LicenseType.Int64 == 5412 && (licenseParam.Int64 == 5411 || licenseParam.Int64 == 5410)) ||
					(seekerLicense.LicenseType.Int64 == 5411 && licenseParam.Int64 == 5412) ||
					(seekerLicense.LicenseType.Int64 == 5413 && licenseParam.Int64 == 5414) ||
					(seekerLicense.LicenseType.Int64 == 5419 && (licenseParam.Int64 == 5417 || licenseParam.Int64 == 5418 || licenseParam.Int64 == 5416)) ||
					(seekerLicense.LicenseType.Int64 == 5418 && (licenseParam.Int64 == 5417 || licenseParam.Int64 == 5416)) ||
					(seekerLicense.LicenseType.Int64 == 5417 && licenseParam.Int64 == 5416) ||
					(seekerLicense.LicenseType.Int64 == 5420 && licenseParam.Int64 == 5421) ||
					(seekerLicense.LicenseType.Int64 == 5425 && (licenseParam.Int64 == 5424 || licenseParam.Int64 == 5423)) ||
					(seekerLicense.LicenseType.Int64 == 5424 && licenseParam.Int64 == 5425) ||
					(seekerLicense.LicenseType.Int64 == 5426 && licenseParam.Int64 == 5427) ||
					(seekerLicense.LicenseType.Int64 == 5428 && licenseParam.Int64 == 5429) ||
					(seekerLicense.LicenseType.Int64 == 5433 && (licenseParam.Int64 == 5432 || licenseParam.Int64 == 5431)) ||
					(seekerLicense.LicenseType.Int64 == 5432 && licenseParam.Int64 == 5433) ||
					(seekerLicense.LicenseType.Int64 == 5440 && (licenseParam.Int64 == 5439 || licenseParam.Int64 == 5438)) ||
					(seekerLicense.LicenseType.Int64 == 5439 && licenseParam.Int64 == 5440) ||
					(seekerLicense.LicenseType.Int64 == 5441 && licenseParam.Int64 == 5442) ||
					(seekerLicense.LicenseType.Int64 == 5445 && (licenseParam.Int64 == 5444 || licenseParam.Int64 == 5443)) ||
					(seekerLicense.LicenseType.Int64 == 5444 && licenseParam.Int64 == 5445) ||
					(seekerLicense.LicenseType.Int64 == 5501 && licenseParam.Int64 == 5502) ||
					(seekerLicense.LicenseType.Int64 == 5504 && licenseParam.Int64 == 5505) ||
					(seekerLicense.LicenseType.Int64 == 5513 && (licenseParam.Int64 == 5511 || licenseParam.Int64 == 5512 || licenseParam.Int64 == 5510)) ||
					(seekerLicense.LicenseType.Int64 == 5512 && (licenseParam.Int64 == 5511 || licenseParam.Int64 == 5510)) ||
					(seekerLicense.LicenseType.Int64 == 5511 && licenseParam.Int64 == 5510) ||
					(seekerLicense.LicenseType.Int64 == 5515 && licenseParam.Int64 == 5516) ||
					(seekerLicense.LicenseType.Int64 == 5517 && licenseParam.Int64 == 5518) ||
					(seekerLicense.LicenseType.Int64 == 5519 && licenseParam.Int64 == 5520) ||
					(seekerLicense.LicenseType.Int64 == 5603 && (licenseParam.Int64 == 5602 || licenseParam.Int64 == 5601)) ||
					(seekerLicense.LicenseType.Int64 == 5602 && licenseParam.Int64 == 5603) ||
					(seekerLicense.LicenseType.Int64 == 5606 && (licenseParam.Int64 == 5605 || licenseParam.Int64 == 5604)) ||
					(seekerLicense.LicenseType.Int64 == 5605 && licenseParam.Int64 == 5606) ||
					(seekerLicense.LicenseType.Int64 == 5611 && (licenseParam.Int64 == 5609 || licenseParam.Int64 == 5610 || licenseParam.Int64 == 5608)) ||
					(seekerLicense.LicenseType.Int64 == 5610 && (licenseParam.Int64 == 5608 || licenseParam.Int64 == 5609)) ||
					(seekerLicense.LicenseType.Int64 == 5609 && licenseParam.Int64 == 5608) ||
					(seekerLicense.LicenseType.Int64 == 5614 && (licenseParam.Int64 == 5613 || licenseParam.Int64 == 5612)) ||
					(seekerLicense.LicenseType.Int64 == 5613 && licenseParam.Int64 == 5614) ||
					(seekerLicense.LicenseType.Int64 == 5617 && (licenseParam.Int64 == 5616 || licenseParam.Int64 == 5615)) ||
					(seekerLicense.LicenseType.Int64 == 5616 && licenseParam.Int64 == 5617) ||
					(seekerLicense.LicenseType.Int64 == 5618 && licenseParam.Int64 == 5619) ||
					(seekerLicense.LicenseType.Int64 == 5622 && (licenseParam.Int64 == 5621 || licenseParam.Int64 == 5620)) ||
					(seekerLicense.LicenseType.Int64 == 5621 && licenseParam.Int64 == 5622) ||
					(seekerLicense.LicenseType.Int64 == 5623 && licenseParam.Int64 == 5624) ||
					(seekerLicense.LicenseType.Int64 == 5703 && (licenseParam.Int64 == 5702 || licenseParam.Int64 == 5701)) ||
					(seekerLicense.LicenseType.Int64 == 5702 && licenseParam.Int64 == 5703) ||
					(seekerLicense.LicenseType.Int64 == 5707 && (licenseParam.Int64 == 5705 || licenseParam.Int64 == 5706 || licenseParam.Int64 == 5704)) ||
					(seekerLicense.LicenseType.Int64 == 5706 && (licenseParam.Int64 == 5705 || licenseParam.Int64 == 5704)) ||
					(seekerLicense.LicenseType.Int64 == 5705 && licenseParam.Int64 == 5704) ||
					(seekerLicense.LicenseType.Int64 == 5714 && (licenseParam.Int64 == 5713 || licenseParam.Int64 == 5712)) ||
					(seekerLicense.LicenseType.Int64 == 5713 && licenseParam.Int64 == 5714) ||
					(seekerLicense.LicenseType.Int64 == 5717 && (licenseParam.Int64 == 5716 || licenseParam.Int64 == 5715)) ||
					(seekerLicense.LicenseType.Int64 == 5716 && licenseParam.Int64 == 5717) ||
					(seekerLicense.LicenseType.Int64 == 5722 && (licenseParam.Int64 == 5720 || licenseParam.Int64 == 5721 || licenseParam.Int64 == 5719)) ||
					(seekerLicense.LicenseType.Int64 == 5721 && (licenseParam.Int64 == 5720 || licenseParam.Int64 == 5719)) ||
					(seekerLicense.LicenseType.Int64 == 5720 && licenseParam.Int64 == 5719) ||
					(seekerLicense.LicenseType.Int64 == 5804 && licenseParam.Int64 == 5805) ||
					(seekerLicense.LicenseType.Int64 == 5808 && licenseParam.Int64 == 5809) ||
					(seekerLicense.LicenseType.Int64 == 5810 && licenseParam.Int64 == 5811) ||
					(seekerLicense.LicenseType.Int64 == 5813 && licenseParam.Int64 == 5815) ||
					(seekerLicense.LicenseType.Int64 == 5814 && licenseParam.Int64 == 5816) ||
					(seekerLicense.LicenseType.Int64 == 5896 && licenseParam.Int64 == 5817) ||
					(seekerLicense.LicenseType.Int64 == 5908 && (licenseParam.Int64 == 5907 || licenseParam.Int64 == 5906)) ||
					(seekerLicense.LicenseType.Int64 == 5907 && licenseParam.Int64 == 5908) ||
					(seekerLicense.LicenseType.Int64 == 6003 && (licenseParam.Int64 == 6002 || licenseParam.Int64 == 6001)) ||
					(seekerLicense.LicenseType.Int64 == 6002 && licenseParam.Int64 == 6003) ||
					(seekerLicense.LicenseType.Int64 == 6005 && licenseParam.Int64 == 6006) ||
					(seekerLicense.LicenseType.Int64 == 6008 && licenseParam.Int64 == 6009) ||
					(seekerLicense.LicenseType.Int64 == 6010 && licenseParam.Int64 == 6011) ||
					(seekerLicense.LicenseType.Int64 == 6012 && licenseParam.Int64 == 6013) ||
					(seekerLicense.LicenseType.Int64 == 6014 && licenseParam.Int64 == 6015) ||
					(seekerLicense.LicenseType.Int64 == 6101 && licenseParam.Int64 == 6102) ||
					(seekerLicense.LicenseType.Int64 == 6103 && licenseParam.Int64 == 6104) ||
					(seekerLicense.LicenseType.Int64 == 6107 && (licenseParam.Int64 == 6106 || licenseParam.Int64 == 6105)) ||
					(seekerLicense.LicenseType.Int64 == 6106 && licenseParam.Int64 == 6107) ||
					(seekerLicense.LicenseType.Int64 == 6110 && (licenseParam.Int64 == 6109 || licenseParam.Int64 == 6108)) ||
					(seekerLicense.LicenseType.Int64 == 6109 && licenseParam.Int64 == 6110) ||
					(seekerLicense.LicenseType.Int64 == 6113 && (licenseParam.Int64 == 6112 || licenseParam.Int64 == 6111)) ||
					(seekerLicense.LicenseType.Int64 == 6112 && licenseParam.Int64 == 6113) ||
					(seekerLicense.LicenseType.Int64 == 6114 && licenseParam.Int64 == 6115) ||
					(seekerLicense.LicenseType.Int64 == 6116 && licenseParam.Int64 == 6117) ||
					(seekerLicense.LicenseType.Int64 == 6118 && licenseParam.Int64 == 6119) ||
					(seekerLicense.LicenseType.Int64 == 6202 && licenseParam.Int64 == 6202) ||
					(seekerLicense.LicenseType.Int64 == 6206 && licenseParam.Int64 == 6207) ||
					(seekerLicense.LicenseType.Int64 == 6208 && licenseParam.Int64 == 6209) ||
					(seekerLicense.LicenseType.Int64 == 6212 && licenseParam.Int64 == 6213) ||
					(seekerLicense.LicenseType.Int64 == 6214 && licenseParam.Int64 == 6215) ||
					(seekerLicense.LicenseType.Int64 == 6301 && licenseParam.Int64 == 6302) ||
					(seekerLicense.LicenseType.Int64 == 6303 && licenseParam.Int64 == 6304) ||
					(seekerLicense.LicenseType.Int64 == 6305 && licenseParam.Int64 == 6306) ||
					(seekerLicense.LicenseType.Int64 == 6404 && (licenseParam.Int64 == 6402 || licenseParam.Int64 == 6403 || licenseParam.Int64 == 6401)) ||
					(seekerLicense.LicenseType.Int64 == 6403 && (licenseParam.Int64 == 6402 || licenseParam.Int64 == 6401)) ||
					(seekerLicense.LicenseType.Int64 == 6402 && licenseParam.Int64 == 6401) ||
					(seekerLicense.LicenseType.Int64 == 6405 && licenseParam.Int64 == 6406) ||
					(seekerLicense.LicenseType.Int64 == 6503 && licenseParam.Int64 == 6504) ||
					(seekerLicense.LicenseType.Int64 == 6507 && licenseParam.Int64 == 6508) ||
					(seekerLicense.LicenseType.Int64 == 6509 && licenseParam.Int64 == 6510) ||
					(seekerLicense.LicenseType.Int64 == 6601 && licenseParam.Int64 == 6602) ||
					(seekerLicense.LicenseType.Int64 == 6603 && licenseParam.Int64 == 6604) ||
					(seekerLicense.LicenseType.Int64 == 6608 && (licenseParam.Int64 == 6607 || licenseParam.Int64 == 6606)) ||
					(seekerLicense.LicenseType.Int64 == 6607 && licenseParam.Int64 == 6608) ||
					(seekerLicense.LicenseType.Int64 == 6611 && (licenseParam.Int64 == 6610 || licenseParam.Int64 == 6609)) ||
					(seekerLicense.LicenseType.Int64 == 6610 && licenseParam.Int64 == 6611) ||
					(seekerLicense.LicenseType.Int64 == 6612 && licenseParam.Int64 == 6613) ||
					(seekerLicense.LicenseType.Int64 == 6620 && (licenseParam.Int64 == 6618 || licenseParam.Int64 == 6619 || licenseParam.Int64 == 6617)) ||
					(seekerLicense.LicenseType.Int64 == 6619 && (licenseParam.Int64 == 6619 || licenseParam.Int64 == 6617)) ||
					(seekerLicense.LicenseType.Int64 == 6618 && licenseParam.Int64 == 6617) ||
					(seekerLicense.LicenseType.Int64 == 6624 && (licenseParam.Int64 == 6622 || licenseParam.Int64 == 6623 || licenseParam.Int64 == 6621)) ||
					(seekerLicense.LicenseType.Int64 == 6623 && (licenseParam.Int64 == 6622 || licenseParam.Int64 == 6621)) ||
					(seekerLicense.LicenseType.Int64 == 6622 && licenseParam.Int64 == 6621) ||
					(seekerLicense.LicenseType.Int64 == 6625 && (licenseParam.Int64 == 6626 || licenseParam.Int64 == 6627)) ||
					(seekerLicense.LicenseType.Int64 == 6626 && licenseParam.Int64 == 6627) ||
					(seekerLicense.LicenseType.Int64 == 6630 && licenseParam.Int64 == 6631) ||
					(seekerLicense.LicenseType.Int64 == 6634 && (licenseParam.Int64 == 6633 || licenseParam.Int64 == 6632)) ||
					(seekerLicense.LicenseType.Int64 == 6633 && licenseParam.Int64 == 6634) ||
					(seekerLicense.LicenseType.Int64 == 6635 && licenseParam.Int64 == 6636) ||
					(seekerLicense.LicenseType.Int64 == 6703 && (licenseParam.Int64 == 6702 || licenseParam.Int64 == 6701)) ||
					(seekerLicense.LicenseType.Int64 == 6702 && licenseParam.Int64 == 6703) ||
					(seekerLicense.LicenseType.Int64 == 6803 && (licenseParam.Int64 == 6802 || licenseParam.Int64 == 6801)) ||
					(seekerLicense.LicenseType.Int64 == 6802 && licenseParam.Int64 == 6803) ||
					(seekerLicense.LicenseType.Int64 == 6806 && (licenseParam.Int64 == 6805 || licenseParam.Int64 == 6804)) ||
					(seekerLicense.LicenseType.Int64 == 6805 && licenseParam.Int64 == 6806) ||
					(seekerLicense.LicenseType.Int64 == 6808 && licenseParam.Int64 == 6809) ||
					(seekerLicense.LicenseType.Int64 == 6814 && (licenseParam.Int64 == 6813 || licenseParam.Int64 == 6812)) ||
					(seekerLicense.LicenseType.Int64 == 6813 && licenseParam.Int64 == 6814) ||
					(seekerLicense.LicenseType.Int64 == 6908 && licenseParam.Int64 == 6909) ||
					(seekerLicense.LicenseType.Int64 == 7001 && licenseParam.Int64 == 7002) ||
					(seekerLicense.LicenseType.Int64 == 7005 && licenseParam.Int64 == 7006) ||
					(seekerLicense.LicenseType.Int64 == 7009 && (licenseParam.Int64 == 7008 || licenseParam.Int64 == 7007)) ||
					(seekerLicense.LicenseType.Int64 == 7008 && licenseParam.Int64 == 7009) ||
					(seekerLicense.LicenseType.Int64 == 7010 && licenseParam.Int64 == 7011) ||
					(seekerLicense.LicenseType.Int64 == 7013 && licenseParam.Int64 == 7014) ||
					(seekerLicense.LicenseType.Int64 == 7015 && licenseParam.Int64 == 7016) ||
					(seekerLicense.LicenseType.Int64 == 7018 && licenseParam.Int64 == 7019) ||
					(seekerLicense.LicenseType.Int64 == 7020 && licenseParam.Int64 == 7021) ||
					(seekerLicense.LicenseType.Int64 == 7023 && licenseParam.Int64 == 7024) ||
					(seekerLicense.LicenseType.Int64 == 7103 && licenseParam.Int64 == 7104) ||
					(seekerLicense.LicenseType.Int64 == 7105 && licenseParam.Int64 == 7106) ||
					(seekerLicense.LicenseType.Int64 == 7109 && (licenseParam.Int64 == 7108 || licenseParam.Int64 == 7107)) ||
					(seekerLicense.LicenseType.Int64 == 7108 && licenseParam.Int64 == 7109) ||
					(seekerLicense.LicenseType.Int64 == 7110 && licenseParam.Int64 == 7111) ||
					(seekerLicense.LicenseType.Int64 == 7113 && licenseParam.Int64 == 7114) ||
					(seekerLicense.LicenseType.Int64 == 7117 && (licenseParam.Int64 == 7116 || licenseParam.Int64 == 7115)) ||
					(seekerLicense.LicenseType.Int64 == 7116 && licenseParam.Int64 == 7117) ||
					(seekerLicense.LicenseType.Int64 == 7118 && licenseParam.Int64 == 7119) ||
					(seekerLicense.LicenseType.Int64 == 7120 && licenseParam.Int64 == 7121) ||
					(seekerLicense.LicenseType.Int64 == 7122 && licenseParam.Int64 == 7123) ||
					(seekerLicense.LicenseType.Int64 == 7124 && licenseParam.Int64 == 7125) ||
					(seekerLicense.LicenseType.Int64 == 7127 && licenseParam.Int64 == 7128) ||
					(seekerLicense.LicenseType.Int64 == 7201 && licenseParam.Int64 == 7202) ||
					(seekerLicense.LicenseType.Int64 == 7302 && licenseParam.Int64 == 7303) ||
					(seekerLicense.LicenseType.Int64 == 7305 && licenseParam.Int64 == 7306) ||
					(seekerLicense.LicenseType.Int64 == 7406 && licenseParam.Int64 == 7407) {
					seekerListAtLicense = append(seekerListAtLicense, seeker)
					continue licenseLoop
				}
			}

		}
	}

	fmt.Println("資格: ", seekerListAtLicense)

	return seekerListAtLicense
}
