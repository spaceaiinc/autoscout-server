package interactor

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type JobInformationInteractor interface {
	// 汎用系 API
	CreateJobInformation(input CreateJobInformationInput) (CreateJobInformationOutput, error)
	UpdateJobInformation(input UpdateJobInformationInput) (UpdateJobInformationOutput, error)
	DeleteJobInformation(input DeleteJobInformationInput) (DeleteJobInformationOutput, error)
	GetJobInformationByID(input GetJobInformationByIDInput) (GetJobInformationByIDOutput, error)
	GetJobInformationByUUID(input GetJobInformationByUUIDInput) (GetJobInformationByUUIDOutput, error)
	GetJobListingByJobInformationUUID(input GetJobListingByJobInformationUUIDInput) (GetJobListingByJobInformationUUIDOutput, error)
	GetJobListingForJobSeeker(input GetJobListingForJobSeekerInput) (GetJobListingForJobSeekerOutput, error)
	GetJobInformationListByBillingAddressID(input GetJobInformationListByBillingAddressIDInput) (GetJobInformationListByBillingAddressIDOutput, error)
	GetJobInformationListByEnterpriseID(input GetJobInformationListByEnterpriseIDInput) (GetJobInformationListByEnterpriseIDOutput, error)
	GetJobInformationListByAgentID(input GetJobInformationListByAgentIDInput) (GetJobInformationListByAgentIDOutput, error)
	GetSelectionFlowPatternListByJobInformationID(input GetSelectionFlowPatternListByJobInformationIDInput) (GetSelectionFlowPatternListByJobInformationIDOutput, error)
	GetOpenSelectionFlowPatternListByJobInformationID(input GetOpenSelectionFlowPatternListByJobInformationIDInput) (GetOpenSelectionFlowPatternListByJobInformationIDOutput, error)
	GetSelectionFlowPatternByID(input GetSelectionFlowPatternByIDInput) (GetSelectionFlowPatternByIDOutput, error)
	CreateSelectionFlowPattern(input CreateSelectionFlowPatternInput) (CreateSelectionFlowPatternOutput, error)
	UpdateSelectionFlowPattern(input UpdateSelectionFlowPatternInput) (UpdateSelectionFlowPatternOutput, error)
	DeltedSelectionFlowPattern(input DeltedSelectionFlowPatternInput) (DeltedSelectionFlowPatternOutput, error)
	GetJobInformationListByIDList(input GetJobInformationListByIDListInput) (GetJobInformationListByIDListOutput, error)
	GetJobListingListByJobSeekerUUID(input GetJobListingListByJobSeekerUUIDInput) (GetJobListingListByJobSeekerUUIDOutput, error)

	// 求職者検索→求人検索 API
	GetJobInformationListByAgentIDAndType(input GetJobInformationListByAgentIDAndTypeInput) (GetJobInformationListByAgentIDAndTypeOutput, error)                   // 求職者検索→求人検索
	GetSearchJobInformationListByAgentIDAndType(input GetSearchJobInformationListByAgentIDAndTypeInput) (GetSearchJobInformationListByAgentIDAndTypeOutput, error) // 求職者検索→求人検索(絞り込み)

	// シェア求職者検索→自社求人検索 API
	GetSearchPublicJobInformationListByAgentIDAndPage(input GetSearchPublicJobInformationListByAgentIDAndPageInput) (GetSearchPublicJobInformationListByAgentIDAndPageOutput, error) // シェア求職者検索→自社求人検索(絞り込み)

	// 求人の絞り込み検索
	GetSearchActiveJobInformationListByAgentID(input GetSearchActiveJobInformationListByAgentIDInput) (GetSearchActiveJobInformationListByAgentIDOutput, error)
	GetSearchJobInformationListByOtherAgentID(input GetSearchJobInformationListByOtherAgentIDInput) (GetSearchJobInformationListByOtherAgentIDOutput, error)

	// LP用 API
	GetSearchJobInformationCountByLPDiagnosis(input GetSearchJobInformationCountByLPDiagnosisInput) (GetSearchJobInformationCountByLPDiagnosisOutput, error)
	GetSearchJobListingListByJobSeekerUUID(input GetSearchJobListingListByJobSeekerUUIDInput) (GetSearchJobListingListByJobSeekerUUIDOutput, error)
	GetJobInformationListForDiagnosis() (GetJobInformationListForDiagnosisOutput, error)
	GetJobListingListAndJobSeekerDesiredForDiagnosis(input GetJobListingListAndJobSeekerDesiredForDiagnosisInput) (GetJobListingListAndJobSeekerDesiredForDiagnosisOutput, error)

	// 求職者のエントリー希望と興味あり求人の取得
	GetJobListingListByJobSeekerUUIDAndInterestedType(input GetJobListingListByJobSeekerUUIDAndInterestedTypeInput) (GetJobListingListByJobSeekerUUIDAndInterestedTypeOutput, error)

	//admin API
	GetAllJobInformation(input GetAllJobInformationInput) (GetAllJobInformationOutput, error)
}

type JobInformationInteractorImpl struct {
	firebase                                           usecase.Firebase
	sendgrid                                           config.Sendgrid
	jobInformationRepository                           usecase.JobInformationRepository
	jobInfoTargetRepository                            usecase.JobInformationTargetRepository
	jobInfoFeatureRepository                           usecase.JobInformationFeatureRepository
	jobInfoPrefectureRepository                        usecase.JobInformationPrefectureRepository
	jobInfoWorkCharmPointRepository                    usecase.JobInformationWorkCharmPointRepository
	jobInfoEmploymentStatusRepository                  usecase.JobInformationEmploymentStatusRepository
	jobInfoRequiredLicenseRepository                   usecase.JobInformationRequiredLicenseRepository
	jobInfoRequiredPCToolRepository                    usecase.JobInformationRequiredPCToolRepository
	jobInfoRequiredLanguageRepository                  usecase.JobInformationRequiredLanguageRepository
	jobInfoRequiredExperienceDevelopmentRepository     usecase.JobInformationRequiredExperienceDevelopmentRepository
	jobInfoRequiredExperienceJobRepository             usecase.JobInformationRequiredExperienceJobRepository
	jobInfoRequiredExperienceIndustryRepository        usecase.JobInformationRequiredExperienceIndustryRepository
	jobInfoRequiredExperienceOccupationRepository      usecase.JobInformationRequiredExperienceOccupationRepository
	jobInfoRequiredSocialExperienceRepository          usecase.JobInformationRequiredSocialExperienceRepository
	jobInfoSelectionFlowPatternRepository              usecase.JobInformationSelectionFlowPatternRepository
	jobInfoSelectionInformationRepository              usecase.JobInformationSelectionInformationRepository
	jobInfoHideToAgentRepository                       usecase.JobInformationHideToAgentRepository
	jobInfoOccupationRepository                        usecase.JobInformationOccupationRepository
	jobInfoRequiredConditionRepository                 usecase.JobInformationRequiredConditionRepository
	jobInfoRequiredExperienceDevelopmentTypeRepository usecase.JobInformationRequiredExperienceDevelopmentTypeRepository
	jobInfoRequiredLanguageTypeRepository              usecase.JobInformationRequiredLanguageTypeRepository
	jobSeekerRepository                                usecase.JobSeekerRepository
	jobSeekerStudentHistoryRepository                  usecase.JobSeekerStudentHistoryRepository
	jobSeekerWorkHistoryRepository                     usecase.JobSeekerWorkHistoryRepository
	jobSeekerExperienceIndustryRepository              usecase.JobSeekerExperienceIndustryRepository
	jobSeekerExperienceOccupationRepository            usecase.JobSeekerExperienceOccupationRepository
	jobSeekerExperienceJobRepository                   usecase.JobSeekerExperienceJobRepository
	jobSeekerLicenseRepository                         usecase.JobSeekerLicenseRepository
	jobSeekerSelfPromotionRepository                   usecase.JobSeekerSelfPromotionRepository
	jobSeekerDocumentRepository                        usecase.JobSeekerDocumentRepository
	jobSeekerDesiredIndustryRepository                 usecase.JobSeekerDesiredIndustryRepository
	jobSeekerDesiredOccupationRepository               usecase.JobSeekerDesiredOccupationRepository
	jobSeekerDesiredWorkLocationRepository             usecase.JobSeekerDesiredWorkLocationRepository
	jobSeekerDesiredHolidayTypeRepository              usecase.JobSeekerDesiredHolidayTypeRepository
	jobSeekerDevelopmentSkillRepository                usecase.JobSeekerDevelopmentSkillRepository
	jobSeekerLanguageSkillRepository                   usecase.JobSeekerLanguageSkillRepository
	jobSeekerPCToolRepository                          usecase.JobSeekerPCToolRepository
	jobSeekerDepartmentHistoryRepository               usecase.JobSeekerDepartmentHistoryRepository
	jobSeekerDesiredCompanyScaleRepository             usecase.JobSeekerDesiredCompanyScaleRepository
	agentRepository                                    usecase.AgentRepository
	agentStaffRepository                               usecase.AgentStaffRepository
	agentAllianceRepository                            usecase.AgentAllianceRepository
	enterpriseProfileRepository                        usecase.EnterpriseProfileRepository
	enterpriseIndustryRepository                       usecase.EnterpriseIndustryRepository
	enterpriseReferenceMaterialRepository              usecase.EnterpriseReferenceMaterialRepository
	taskGroupRepository                                usecase.TaskGroupRepository
	taskRepository                                     usecase.TaskRepository
	selectionQuestionnaireRepository                   usecase.SelectionQuestionnaireRepository
	evaluationPointRepository                          usecase.EvaluationPointRepository
	billingAddressRepository                           usecase.BillingAddressRepository
	jobSeekerInterestedJobListingRepository            usecase.JobSeekerInterestedJobListingRepository
}

// JobInformationInteractorImpl is an implementation of JobInformationInteractor
func NewJobInformationInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	jR usecase.JobInformationRepository,
	jtR usecase.JobInformationTargetRepository,
	jfR usecase.JobInformationFeatureRepository,
	jpR usecase.JobInformationPrefectureRepository,
	jwcpR usecase.JobInformationWorkCharmPointRepository,
	jesR usecase.JobInformationEmploymentStatusRepository,
	jrlR usecase.JobInformationRequiredLicenseRepository,
	jrptR usecase.JobInformationRequiredPCToolRepository,
	jrlgR usecase.JobInformationRequiredLanguageRepository,
	jredR usecase.JobInformationRequiredExperienceDevelopmentRepository,
	jrejR usecase.JobInformationRequiredExperienceJobRepository,
	jreiR usecase.JobInformationRequiredExperienceIndustryRepository,
	jreoR usecase.JobInformationRequiredExperienceOccupationRepository,
	jrseR usecase.JobInformationRequiredSocialExperienceRepository,
	jsfpR usecase.JobInformationSelectionFlowPatternRepository,
	jsiR usecase.JobInformationSelectionInformationRepository,
	jhtaR usecase.JobInformationHideToAgentRepository,
	joR usecase.JobInformationOccupationRepository,
	jrcR usecase.JobInformationRequiredConditionRepository,
	jiredtR usecase.JobInformationRequiredExperienceDevelopmentTypeRepository,
	jirltR usecase.JobInformationRequiredLanguageTypeRepository,
	jsR usecase.JobSeekerRepository,
	jsshR usecase.JobSeekerStudentHistoryRepository,
	jswhR usecase.JobSeekerWorkHistoryRepository,
	jseiR usecase.JobSeekerExperienceIndustryRepository,
	jseoR usecase.JobSeekerExperienceOccupationRepository,
	jsejR usecase.JobSeekerExperienceJobRepository,
	jslR usecase.JobSeekerLicenseRepository,
	jsspR usecase.JobSeekerSelfPromotionRepository,
	jsdR usecase.JobSeekerDocumentRepository,
	jsdiR usecase.JobSeekerDesiredIndustryRepository,
	jsdoR usecase.JobSeekerDesiredOccupationRepository,
	jsdwlR usecase.JobSeekerDesiredWorkLocationRepository,
	jsdhtR usecase.JobSeekerDesiredHolidayTypeRepository,
	jsdsR usecase.JobSeekerDevelopmentSkillRepository,
	jslsR usecase.JobSeekerLanguageSkillRepository,
	jsptR usecase.JobSeekerPCToolRepository,
	jsdhR usecase.JobSeekerDepartmentHistoryRepository,
	jsdcsR usecase.JobSeekerDesiredCompanyScaleRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	aaR usecase.AgentAllianceRepository,
	epR usecase.EnterpriseProfileRepository,
	eiR usecase.EnterpriseIndustryRepository,
	ermR usecase.EnterpriseReferenceMaterialRepository,
	tgR usecase.TaskGroupRepository,
	tR usecase.TaskRepository,
	sqR usecase.SelectionQuestionnaireRepository,
	evpR usecase.EvaluationPointRepository,
	baR usecase.BillingAddressRepository,
	jsijlR usecase.JobSeekerInterestedJobListingRepository,
) JobInformationInteractor {
	return &JobInformationInteractorImpl{
		firebase:                                           fb,
		sendgrid:                                           sg,
		jobInformationRepository:                           jR,
		jobInfoTargetRepository:                            jtR,
		jobInfoFeatureRepository:                           jfR,
		jobInfoPrefectureRepository:                        jpR,
		jobInfoWorkCharmPointRepository:                    jwcpR,
		jobInfoEmploymentStatusRepository:                  jesR,
		jobInfoRequiredLicenseRepository:                   jrlR,
		jobInfoRequiredPCToolRepository:                    jrptR,
		jobInfoRequiredLanguageRepository:                  jrlgR,
		jobInfoRequiredExperienceDevelopmentRepository:     jredR,
		jobInfoRequiredExperienceJobRepository:             jrejR,
		jobInfoRequiredExperienceIndustryRepository:        jreiR,
		jobInfoRequiredExperienceOccupationRepository:      jreoR,
		jobInfoRequiredSocialExperienceRepository:          jrseR,
		jobInfoSelectionFlowPatternRepository:              jsfpR,
		jobInfoSelectionInformationRepository:              jsiR,
		jobInfoHideToAgentRepository:                       jhtaR,
		jobInfoOccupationRepository:                        joR,
		jobInfoRequiredConditionRepository:                 jrcR,
		jobInfoRequiredExperienceDevelopmentTypeRepository: jiredtR,
		jobInfoRequiredLanguageTypeRepository:              jirltR,
		jobSeekerRepository:                                jsR,
		jobSeekerStudentHistoryRepository:                  jsshR,
		jobSeekerWorkHistoryRepository:                     jswhR,
		jobSeekerExperienceIndustryRepository:              jseiR,
		jobSeekerExperienceOccupationRepository:            jseoR,
		jobSeekerExperienceJobRepository:                   jsejR,
		jobSeekerLicenseRepository:                         jslR,
		jobSeekerSelfPromotionRepository:                   jsspR,
		jobSeekerDocumentRepository:                        jsdR,
		jobSeekerDesiredIndustryRepository:                 jsdiR,
		jobSeekerDesiredOccupationRepository:               jsdoR,
		jobSeekerDesiredWorkLocationRepository:             jsdwlR,
		jobSeekerDesiredHolidayTypeRepository:              jsdhtR,
		jobSeekerDevelopmentSkillRepository:                jsdsR,
		jobSeekerLanguageSkillRepository:                   jslsR,
		jobSeekerPCToolRepository:                          jsptR,
		jobSeekerDepartmentHistoryRepository:               jsdhR,
		jobSeekerDesiredCompanyScaleRepository:             jsdcsR,
		agentRepository:                                    aR,
		agentStaffRepository:                               asR,
		agentAllianceRepository:                            aaR,
		enterpriseProfileRepository:                        epR,
		enterpriseIndustryRepository:                       eiR,
		enterpriseReferenceMaterialRepository:              ermR,
		taskGroupRepository:                                tgR,
		taskRepository:                                     tR,
		selectionQuestionnaireRepository:                   sqR,
		evaluationPointRepository:                          evpR,
		billingAddressRepository:                           baR,
		jobSeekerInterestedJobListingRepository:            jsijlR,
	}
}

/****************************************************************************************/
/// Admin API
//

type CreateJobInformationInput struct {
	CreateParam      entity.CreateJobInformationParam
	BillingAddressID uint
}

type CreateJobInformationOutput struct {
	JobInformation *entity.JobInformation
}

func (i *JobInformationInteractorImpl) CreateJobInformation(input CreateJobInformationInput) (CreateJobInformationOutput, error) {
	var (
		output CreateJobInformationOutput
		err    error
	)

	jobInformation := entity.NewJobInformation(
		input.BillingAddressID,
		input.CreateParam.Title,
		input.CreateParam.RecruitmentState,
		input.CreateParam.ExpirationDate,
		input.CreateParam.WorkDetail,
		input.CreateParam.NumberOfHires,
		input.CreateParam.WorkLocation,
		input.CreateParam.Transfer,
		input.CreateParam.TransferDetail,
		input.CreateParam.UnderIncome,
		input.CreateParam.OverIncome,
		input.CreateParam.Salary,
		input.CreateParam.Insurance,
		input.CreateParam.WorkTime,
		input.CreateParam.OvertimeAverage,
		input.CreateParam.FixedOvertimePayment,
		input.CreateParam.FixedOvertimeDetail,
		input.CreateParam.TrialPeriod,
		input.CreateParam.TrialPeriodDetail,
		input.CreateParam.EmploymentPeriod,
		input.CreateParam.EmploymentPeriodDetail,
		input.CreateParam.HolidayType,
		input.CreateParam.HolidayDetail,
		input.CreateParam.PassiveSmoking,
		input.CreateParam.SelectionFlow,
		input.CreateParam.Gender,
		input.CreateParam.Nationality,
		input.CreateParam.FinalEducation,
		input.CreateParam.SchoolLevel,
		input.CreateParam.MedicalHistory,
		input.CreateParam.AgeUnder,
		input.CreateParam.AgeOver,
		input.CreateParam.JobChange,
		input.CreateParam.ShortResignation,
		input.CreateParam.ShortResignationRemarks,
		input.CreateParam.SocialExperienceYear,
		input.CreateParam.SocialExperienceMonth,
		input.CreateParam.Appearance,
		input.CreateParam.Communication,
		input.CreateParam.Thinking,
		input.CreateParam.TargetDetail,
		input.CreateParam.Commission,
		input.CreateParam.CommissionRate,
		input.CreateParam.CommissionDetail,
		input.CreateParam.RefundPolicy,
		input.CreateParam.RequiredExperienceJobDetail,
		input.CreateParam.SecretMemo,
		input.CreateParam.RequiredDocumentsDetail,
		input.CreateParam.EmploymentInsurance,
		input.CreateParam.AccidentInsurance,
		input.CreateParam.HealthInsurance,
		input.CreateParam.PensionInsurance,
		input.CreateParam.RegisterPhase,
		input.CreateParam.StudyCategory,
		input.CreateParam.DriverLicence,
		input.CreateParam.WordSkill,
		input.CreateParam.ExcelSkill,
		input.CreateParam.PowerPointSkill,
		input.CreateParam.IsExternal,
		input.CreateParam.WorkDetailAfterHiring,
		input.CreateParam.WorkDetailScopeOfChange,
		input.CreateParam.OfferRate,
		input.CreateParam.DocumentPassingRate,
		input.CreateParam.NumberOfRecentApplications,
		input.CreateParam.IsGuaranteedInterview,
	)

	err = i.jobInformationRepository.Create(jobInformation)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, target := range input.CreateParam.Targets {
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

	for _, feature := range input.CreateParam.Features {
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

	for _, prefecture := range input.CreateParam.Prefectures {
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

	for _, workCharmPoint := range input.CreateParam.WorkCharmPoints {
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

	for _, employmentStatus := range input.CreateParam.EmploymentStatuses {
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

	// 必要条件
	for _, requiredCondition := range input.CreateParam.RequiredConditions {
		// 開発経験が2つ未満の場合は作成しない
		if len(requiredCondition.RequiredExperienceDevelopments) < 2 {
			continue
		}
		// 全て未入力の場合は作成しない
		if !requiredCondition.RequiredManagement.Valid &&
			!requiredCondition.RequiredExperienceDevelopments[0].DevelopmentCategory.Valid &&
			!requiredCondition.RequiredExperienceDevelopments[1].DevelopmentCategory.Valid &&
			len(requiredCondition.RequiredExperienceJobs.ExperienceIndustries) == 0 &&
			len(requiredCondition.RequiredExperienceJobs.ExperienceOccupations) == 0 &&
			len(requiredCondition.RequiredLanguages.LanguageTypes) == 0 &&
			len(requiredCondition.RequiredPCTools) == 0 &&
			len(requiredCondition.RequiredLicenses) == 0 {
			continue
		}

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

	for _, requiredSocialExperience := range input.CreateParam.RequiredSocialExperiences {
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

	for _, selectionFlowPattern := range input.CreateParam.SelectionFlowPatterns {
		sfp := entity.NewJobInformationSelectionFlowPattern(
			jobInformation.ID,
			selectionFlowPattern.PublicStatus,
			selectionFlowPattern.FlowTitle,
			selectionFlowPattern.FlowPattern,
			false,
		)

		sfp.SelectionInformations = selectionFlowPattern.SelectionInformations

		err = i.jobInfoSelectionFlowPatternRepository.Create(sfp)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, selectionInformation := range sfp.SelectionInformations {
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
				fmt.Println(err)
				return output, err
			}
		}
	}

	// 非公開エージェント
	for _, hideToAgent := range input.CreateParam.HideToAgents {
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

	// 募集職種
	for _, occupation := range input.CreateParam.Occupations {
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

	output.JobInformation = jobInformation
	output.JobInformation.Targets = input.CreateParam.Targets
	output.JobInformation.Prefectures = input.CreateParam.Prefectures
	output.JobInformation.Features = input.CreateParam.Features
	output.JobInformation.WorkCharmPoints = input.CreateParam.WorkCharmPoints
	output.JobInformation.EmploymentStatuses = input.CreateParam.EmploymentStatuses
	output.JobInformation.RequiredConditions = input.CreateParam.RequiredConditions
	output.JobInformation.RequiredSocialExperiences = input.CreateParam.RequiredSocialExperiences
	output.JobInformation.SelectionFlowPatterns = input.CreateParam.SelectionFlowPatterns
	output.JobInformation.HideToAgents = input.CreateParam.HideToAgents
	output.JobInformation.Occupations = input.CreateParam.Occupations

	return output, nil
}

// 求人の更新
type UpdateJobInformationInput struct {
	UpdateParam      entity.UpdateJobInformationParam
	JobInformationID uint
}

type UpdateJobInformationOutput struct {
	JobInformation *entity.JobInformation
}

func (i *JobInformationInteractorImpl) UpdateJobInformation(input UpdateJobInformationInput) (UpdateJobInformationOutput, error) {
	var (
		output UpdateJobInformationOutput
		err    error
	)

	jobInformation := entity.NewJobInformation(
		input.UpdateParam.BillingAddressID,
		input.UpdateParam.Title,
		input.UpdateParam.RecruitmentState,
		input.UpdateParam.ExpirationDate,
		input.UpdateParam.WorkDetail,
		input.UpdateParam.NumberOfHires,
		input.UpdateParam.WorkLocation,
		input.UpdateParam.Transfer,
		input.UpdateParam.TransferDetail,
		input.UpdateParam.UnderIncome,
		input.UpdateParam.OverIncome,
		input.UpdateParam.Salary,
		input.UpdateParam.Insurance,
		input.UpdateParam.WorkTime,
		input.UpdateParam.OvertimeAverage,
		input.UpdateParam.FixedOvertimePayment,
		input.UpdateParam.FixedOvertimeDetail,
		input.UpdateParam.TrialPeriod,
		input.UpdateParam.TrialPeriodDetail,
		input.UpdateParam.EmploymentPeriod,
		input.UpdateParam.EmploymentPeriodDetail,
		input.UpdateParam.HolidayType,
		input.UpdateParam.HolidayDetail,
		input.UpdateParam.PassiveSmoking,
		input.UpdateParam.SelectionFlow,
		input.UpdateParam.Gender,
		input.UpdateParam.Nationality,
		input.UpdateParam.FinalEducation,
		input.UpdateParam.SchoolLevel,
		input.UpdateParam.MedicalHistory,
		input.UpdateParam.AgeUnder,
		input.UpdateParam.AgeOver,
		input.UpdateParam.JobChange,
		input.UpdateParam.ShortResignation,
		input.UpdateParam.ShortResignationRemarks,
		input.UpdateParam.SocialExperienceYear,
		input.UpdateParam.SocialExperienceMonth,
		input.UpdateParam.Appearance,
		input.UpdateParam.Communication,
		input.UpdateParam.Thinking,
		input.UpdateParam.TargetDetail,
		input.UpdateParam.Commission,
		input.UpdateParam.CommissionRate,
		input.UpdateParam.CommissionDetail,
		input.UpdateParam.RefundPolicy,
		input.UpdateParam.RequiredExperienceJobDetail,
		input.UpdateParam.SecretMemo,
		input.UpdateParam.RequiredDocumentsDetail,
		input.UpdateParam.EmploymentInsurance,
		input.UpdateParam.AccidentInsurance,
		input.UpdateParam.HealthInsurance,
		input.UpdateParam.PensionInsurance,
		input.UpdateParam.RegisterPhase,
		input.UpdateParam.StudyCategory,
		input.UpdateParam.DriverLicence,
		input.UpdateParam.WordSkill,
		input.UpdateParam.ExcelSkill,
		input.UpdateParam.PowerPointSkill,
		input.UpdateParam.IsExternal,
		input.UpdateParam.WorkDetailAfterHiring,
		input.UpdateParam.WorkDetailScopeOfChange,
		input.UpdateParam.OfferRate,
		input.UpdateParam.DocumentPassingRate,
		input.UpdateParam.NumberOfRecentApplications,
		input.UpdateParam.IsGuaranteedInterview,
	)

	err = i.jobInformationRepository.Update(input.JobInformationID, jobInformation)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.jobInfoTargetRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, target := range input.UpdateParam.Targets {
		t := entity.NewJobInformationTarget(
			input.JobInformationID,
			target.Target,
		)

		err = i.jobInfoTargetRepository.Create(t)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobInfoFeatureRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, feature := range input.UpdateParam.Features {
		f := entity.NewJobInformationFeature(
			input.JobInformationID,
			feature.Feature,
		)

		err = i.jobInfoFeatureRepository.Create(f)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobInfoPrefectureRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, prefecture := range input.UpdateParam.Prefectures {
		p := entity.NewJobInformationPrefecture(
			input.JobInformationID,
			prefecture.Prefecture,
		)

		err = i.jobInfoPrefectureRepository.Create(p)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobInfoWorkCharmPointRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, workCharmPoint := range input.UpdateParam.WorkCharmPoints {
		wcp := entity.NewJobInformationWorkCharmPoint(
			input.JobInformationID,
			workCharmPoint.Title,
			workCharmPoint.Contents,
		)

		err = i.jobInfoWorkCharmPointRepository.Create(wcp)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobInfoEmploymentStatusRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, employmentStatus := range input.UpdateParam.EmploymentStatuses {
		es := entity.NewJobInformationEmploymentStatus(
			input.JobInformationID,
			employmentStatus.EmploymentStatus,
		)

		err = i.jobInfoEmploymentStatusRepository.Create(es)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 必要条件以下　削除
	err = i.jobInfoRequiredConditionRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, requiredCondition := range input.UpdateParam.RequiredConditions {
		// 開発経験が2つ未満の場合は作成しない
		if len(requiredCondition.RequiredExperienceDevelopments) < 2 {
			continue
		}
		// 全て未入力の場合は作成しない
		if !requiredCondition.RequiredManagement.Valid &&
			!requiredCondition.RequiredExperienceDevelopments[0].DevelopmentCategory.Valid &&
			!requiredCondition.RequiredExperienceDevelopments[1].DevelopmentCategory.Valid &&
			len(requiredCondition.RequiredExperienceJobs.ExperienceIndustries) == 0 &&
			len(requiredCondition.RequiredExperienceJobs.ExperienceOccupations) == 0 &&
			len(requiredCondition.RequiredLanguages.LanguageTypes) == 0 &&
			len(requiredCondition.RequiredPCTools) == 0 &&
			len(requiredCondition.RequiredLicenses) == 0 {
			continue
		}

		rc := entity.NewJobInformationRequiredCondition(
			input.JobInformationID,
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
	// 必要条件以上

	err = i.jobInfoRequiredSocialExperienceRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, requiredSocialExperience := range input.UpdateParam.RequiredSocialExperiences {
		rse := entity.NewJobInformationRequiredSocialExperience(
			input.JobInformationID,
			requiredSocialExperience.SocialExperienceType,
		)

		err = i.jobInfoRequiredSocialExperienceRepository.Create(rse)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 非公開エージェント hideToAgent
	err = i.jobInfoHideToAgentRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hideToAgent := range input.UpdateParam.HideToAgents {
		hta := entity.NewJobInformationHideToAgent(
			input.JobInformationID,
			hideToAgent.AgentID,
		)

		err = i.jobInfoHideToAgentRepository.Create(hta)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 募集職種
	err = i.jobInfoOccupationRepository.DeleteByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	for _, occupation := range input.UpdateParam.Occupations {
		oc := entity.NewJobInformationOccupation(
			input.JobInformationID,
			occupation.Occupation,
		)

		err = i.jobInfoOccupationRepository.Create(oc)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.JobInformation = jobInformation
	output.JobInformation.Targets = input.UpdateParam.Targets
	output.JobInformation.Features = input.UpdateParam.Features
	output.JobInformation.Prefectures = input.UpdateParam.Prefectures
	output.JobInformation.WorkCharmPoints = input.UpdateParam.WorkCharmPoints
	output.JobInformation.EmploymentStatuses = input.UpdateParam.EmploymentStatuses
	output.JobInformation.RequiredConditions = input.UpdateParam.RequiredConditions
	output.JobInformation.RequiredSocialExperiences = input.UpdateParam.RequiredSocialExperiences
	output.JobInformation.HideToAgents = input.UpdateParam.HideToAgents
	output.JobInformation.Occupations = input.UpdateParam.Occupations

	// 求人更新時に紐付けた請求先情報を取得
	billing, err := i.billingAddressRepository.FindByID(input.UpdateParam.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 請求先が変更されたことでRAとCA担当が同じになったタスクがないかを確認
	tasGroupkList, err := i.taskGroupRepository.GetNotDoubleSidedSameRAAndCAByStaffID(billing.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// RAとCA担当が同じになったタスクがある場合はそのタスクの「is_double_sided」をtrueに更新する
	if len(tasGroupkList) > 0 {
		var groupIDList []uint

		for _, tg := range tasGroupkList {
			groupIDList = append(groupIDList, tg.ID)
		}

		err := i.taskGroupRepository.UpdateListIsDoubleSided(groupIDList, true)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	return output, nil
}

type DeleteJobInformationInput struct {
	JobInformationID uint
}

type DeleteJobInformationOutput struct {
	OK bool
}

func (i *JobInformationInteractorImpl) DeleteJobInformation(input DeleteJobInformationInput) (DeleteJobInformationOutput, error) {
	var (
		output DeleteJobInformationOutput
	)

	err := i.jobInformationRepository.Delete(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 求人IDを使って求人情報を取得する
type GetJobInformationByIDInput struct {
	JobInformationID uint
}

type GetJobInformationByIDOutput struct {
	JobInformation *entity.JobInformation
}

func (i *JobInformationInteractorImpl) GetJobInformationByID(input GetJobInformationByIDInput) (GetJobInformationByIDOutput, error) {
	var (
		output GetJobInformationByIDOutput
		err    error
	)

	jobInformation, err := i.jobInformationRepository.FindByID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	targets, err := i.jobInfoTargetRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.jobInfoFeatureRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.jobInfoRequiredPCToolRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformations, err := i.jobInfoSelectionInformationRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// hideToAgent
	hideToAgents, err := i.jobInfoHideToAgentRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.enterpriseIndustryRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, t := range targets {
		value := entity.JobInformationTarget{
			JobInformationID: jobInformation.ID,
			Target:           t.Target,
		}
		jobInformation.Targets = append(jobInformation.Targets, value)
	}

	for _, f := range features {
		value := entity.JobInformationFeature{
			JobInformationID: f.JobInformationID,
			Feature:          f.Feature,
		}
		jobInformation.Features = append(jobInformation.Features, value)
	}

	for _, p := range prefectures {
		value := entity.JobInformationPrefecture{
			JobInformationID: p.JobInformationID,
			Prefecture:       p.Prefecture,
		}
		jobInformation.Prefectures = append(jobInformation.Prefectures, value)
	}

	for _, wcp := range workCharmPoints {
		value := entity.JobInformationWorkCharmPoint{
			JobInformationID: wcp.JobInformationID,
			Title:            wcp.Title,
			Contents:         wcp.Contents,
		}
		jobInformation.WorkCharmPoints = append(jobInformation.WorkCharmPoints, value)
	}

	for _, es := range employmentStatuses {
		value := entity.JobInformationEmploymentStatus{
			JobInformationID: es.JobInformationID,
			EmploymentStatus: es.EmploymentStatus,
		}
		jobInformation.EmploymentStatuses = append(jobInformation.EmploymentStatuses, value)
	}

	// 必要条件
	for _, condition := range requiredConditions {
		// 必要資格
		for _, license := range requiredLicenses {
			if condition.ID == license.ConditionID {
				condition.RequiredLicenses = append(condition.RequiredLicenses, *license)
			}
		}

		// 必要PCツール
		for _, pcTool := range requiredPCTools {
			if condition.ID == pcTool.ConditionID {
				condition.RequiredPCTools = append(condition.RequiredPCTools, *pcTool)
			}
		}

		// 必要言語
		for _, language := range requiredLanguages {
			if condition.ID == language.ConditionID {

				// 必要言語タイプ
				for _, languageType := range requiredLanguageTypes {
					if language.ID == languageType.LanguageID {
						language.LanguageTypes = append(language.LanguageTypes, *languageType)
					}
				}
				condition.RequiredLanguages = *language
			}
		}

		// 必要経験（開発）
		for _, experienceDevelopment := range requiredExperienceDevelopments {
			if condition.ID == experienceDevelopment.ConditionID {

				// 必要言語タイプ
				for _, experienceDevelopmentType := range requiredExperienceDevelopmentTypes {
					if experienceDevelopment.ID == experienceDevelopmentType.ExperienceDevelopmentID {
						experienceDevelopment.ExperienceDevelopmentTypes = append(experienceDevelopment.ExperienceDevelopmentTypes, *experienceDevelopmentType)
					}
				}
				condition.RequiredExperienceDevelopments = append(condition.RequiredExperienceDevelopments, *experienceDevelopment)
			}
		}

		// 必要業職種経験
		for _, experienceJob := range requiredExperienceJobs {
			if condition.ID == uint(experienceJob.ConditionID) {

				// 必要経験（業界）
				for _, industry := range requiredExperienceIndustries {
					if experienceJob.ID == industry.ExperienceJobID {
						experienceJob.ExperienceIndustries = append(experienceJob.ExperienceIndustries, *industry)
					}
				}

				// 必要経験（職種）
				for _, occupation := range requiredExperienceOccupations {
					if experienceJob.ID == occupation.ExperienceJobID {
						experienceJob.ExperienceOccupations = append(experienceJob.ExperienceOccupations, *occupation)
					}
				}

				condition.RequiredExperienceJobs = *experienceJob
			}
		}

		jobInformation.RequiredConditions = append(jobInformation.RequiredConditions, *condition)
	}

	for _, rse := range requiredSocialExperiences {
		value := entity.JobInformationRequiredSocialExperience{
			JobInformationID:     rse.JobInformationID,
			SocialExperienceType: rse.SocialExperienceType,
		}
		jobInformation.RequiredSocialExperiences = append(jobInformation.RequiredSocialExperiences, value)
	}

	for _, sfp := range selectionFlowPatterns {
		value := entity.JobInformationSelectionFlowPattern{
			JobInformationID: sfp.JobInformationID,
			PublicStatus:     sfp.PublicStatus,
			FlowTitle:        sfp.FlowTitle,
			FlowPattern:      sfp.FlowPattern,
		}

		for _, si := range selectionInformations {
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
		jobInformation.SelectionFlowPatterns = append(jobInformation.SelectionFlowPatterns, value)
	}

	for _, hta := range hideToAgents {
		value := entity.JobInformationHideToAgent{
			JobInformationID: hta.JobInformationID,
			AgentID:          hta.AgentID,
			AgentName:        hta.AgentName,
		}
		jobInformation.HideToAgents = append(jobInformation.HideToAgents, value)
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

	output.JobInformation = jobInformation

	return output, nil
}

// 求人のuuidを使って求人情報を取得する
type GetJobInformationByUUIDInput struct {
	UUID uuid.UUID
}

type GetJobInformationByUUIDOutput struct {
	JobInformation *entity.JobInformation
}

func (i *JobInformationInteractorImpl) GetJobInformationByUUID(input GetJobInformationByUUIDInput) (GetJobInformationByUUIDOutput, error) {
	var (
		output GetJobInformationByUUIDOutput
		err    error
	)

	jobInformation, err := i.jobInformationRepository.FindByUUID(input.UUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	targets, err := i.jobInfoTargetRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.jobInfoFeatureRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.jobInfoRequiredPCToolRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformations, err := i.jobInfoSelectionInformationRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// hideToAgent
	hideToAgents, err := i.jobInfoHideToAgentRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIndustries, err := i.enterpriseIndustryRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, t := range targets {
		value := entity.JobInformationTarget{
			JobInformationID: jobInformation.ID,
			Target:           t.Target,
		}
		jobInformation.Targets = append(jobInformation.Targets, value)
	}

	for _, f := range features {
		value := entity.JobInformationFeature{
			JobInformationID: f.JobInformationID,
			Feature:          f.Feature,
		}
		jobInformation.Features = append(jobInformation.Features, value)
	}

	for _, p := range prefectures {
		value := entity.JobInformationPrefecture{
			JobInformationID: p.JobInformationID,
			Prefecture:       p.Prefecture,
		}
		jobInformation.Prefectures = append(jobInformation.Prefectures, value)
	}

	for _, wcp := range workCharmPoints {
		value := entity.JobInformationWorkCharmPoint{
			JobInformationID: wcp.JobInformationID,
			Title:            wcp.Title,
			Contents:         wcp.Contents,
		}
		jobInformation.WorkCharmPoints = append(jobInformation.WorkCharmPoints, value)
	}

	for _, es := range employmentStatuses {
		value := entity.JobInformationEmploymentStatus{
			JobInformationID: es.JobInformationID,
			EmploymentStatus: es.EmploymentStatus,
		}
		jobInformation.EmploymentStatuses = append(jobInformation.EmploymentStatuses, value)
	}

	for _, condition := range requiredConditions {
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

	for _, rse := range requiredSocialExperiences {
		value := entity.JobInformationRequiredSocialExperience{
			JobInformationID:     rse.JobInformationID,
			SocialExperienceType: rse.SocialExperienceType,
		}
		jobInformation.RequiredSocialExperiences = append(jobInformation.RequiredSocialExperiences, value)
	}

	for _, sfp := range selectionFlowPatterns {
		value := entity.JobInformationSelectionFlowPattern{
			JobInformationID: sfp.JobInformationID,
			PublicStatus:     sfp.PublicStatus,
			FlowTitle:        sfp.FlowTitle,
			FlowPattern:      sfp.FlowPattern,
		}

		for _, si := range selectionInformations {
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
		jobInformation.SelectionFlowPatterns = append(jobInformation.SelectionFlowPatterns, value)
	}

	for _, hta := range hideToAgents {
		value := entity.JobInformationHideToAgent{
			JobInformationID: hta.JobInformationID,
			AgentID:          hta.AgentID,
			AgentName:        hta.AgentName,
		}
		jobInformation.HideToAgents = append(jobInformation.HideToAgents, value)
	}

	for _, oc := range occupations {
		value := entity.JobInformationOccupation{
			JobInformationID: oc.JobInformationID,
			Occupation:       oc.Occupation,
		}
		jobInformation.Occupations = append(jobInformation.Occupations, value)
	}

	for _, industry := range enterpriseIndustries {
		if jobInformation.EnterpriseID == industry.EnterpriseID {
			valueIn := entity.EnterpriseIndustry{
				Industry: industry.Industry,
			}
			jobInformation.Industries = append(jobInformation.Industries, valueIn)
		}
	}

	output.JobInformation = jobInformation

	return output, nil
}

// 求人のuuidを使って求人情報を取得する
type GetJobListingByJobInformationUUIDInput struct {
	UUID uuid.UUID
}

type GetJobListingByJobInformationUUIDOutput struct {
	JobListing *entity.JobListing
}

func (i *JobInformationInteractorImpl) GetJobListingByJobInformationUUID(input GetJobListingByJobInformationUUIDInput) (GetJobListingByJobInformationUUIDOutput, error) {
	var (
		output GetJobListingByJobInformationUUIDOutput
		err    error
	)

	jobInformation, err := i.jobInformationRepository.FindByUUID(input.UUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIndustries, err := i.enterpriseIndustryRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, p := range prefectures {
		value := entity.JobInformationPrefecture{
			JobInformationID: p.JobInformationID,
			Prefecture:       p.Prefecture,
		}
		jobInformation.Prefectures = append(jobInformation.Prefectures, value)
	}

	for _, wcp := range workCharmPoints {
		value := entity.JobInformationWorkCharmPoint{
			JobInformationID: wcp.JobInformationID,
			Title:            wcp.Title,
			Contents:         wcp.Contents,
		}
		jobInformation.WorkCharmPoints = append(jobInformation.WorkCharmPoints, value)
	}

	for _, es := range employmentStatuses {
		value := entity.JobInformationEmploymentStatus{
			JobInformationID: es.JobInformationID,
			EmploymentStatus: es.EmploymentStatus,
		}
		jobInformation.EmploymentStatuses = append(jobInformation.EmploymentStatuses, value)
	}

	for _, industry := range enterpriseIndustries {
		if jobInformation.EnterpriseID == industry.EnterpriseID {
			valueIn := entity.EnterpriseIndustry{
				Industry: industry.Industry,
			}
			jobInformation.Industries = append(jobInformation.Industries, valueIn)
		}
	}

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

	jobListing.JobInformationUUID = jobInformation.UUID
	jobListing.Industries = jobInformation.Industries
	jobListing.Prefectures = jobInformation.Prefectures
	jobListing.EmploymentStatuses = jobInformation.EmploymentStatuses
	jobListing.WorkCharmPoints = jobInformation.WorkCharmPoints

	output.JobListing = jobListing

	return output, nil
}

// 求職者が確認する求人票情報を取得（求人票 + タスクに紐づいた選考情報）
type GetJobListingForJobSeekerInput struct {
	JobInformationUUID uuid.UUID
	JobSeekerUUID      uuid.UUID
}

type GetJobListingForJobSeekerOutput struct {
	JobListing *entity.JobListing
}

func (i *JobInformationInteractorImpl) GetJobListingForJobSeeker(input GetJobListingForJobSeekerInput) (GetJobListingForJobSeekerOutput, error) {
	var (
		output GetJobListingForJobSeekerOutput
		err    error
	)

	jobInformation, err := i.jobInformationRepository.FindByUUID(input.JobInformationUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIndustries, err := i.enterpriseIndustryRepository.GetByJobInformationID(jobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	latestTask, err := i.taskRepository.FindLatestWithRelatedByJobSeekerUUIDAndJobInformationUUID(input.JobSeekerUUID, input.JobInformationUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, p := range prefectures {
		value := entity.JobInformationPrefecture{
			JobInformationID: p.JobInformationID,
			Prefecture:       p.Prefecture,
		}
		jobInformation.Prefectures = append(jobInformation.Prefectures, value)
	}

	for _, wcp := range workCharmPoints {
		value := entity.JobInformationWorkCharmPoint{
			JobInformationID: wcp.JobInformationID,
			Title:            wcp.Title,
			Contents:         wcp.Contents,
		}
		jobInformation.WorkCharmPoints = append(jobInformation.WorkCharmPoints, value)
	}

	for _, es := range employmentStatuses {
		value := entity.JobInformationEmploymentStatus{
			JobInformationID: es.JobInformationID,
			EmploymentStatus: es.EmploymentStatus,
		}
		jobInformation.EmploymentStatuses = append(jobInformation.EmploymentStatuses, value)
	}

	for _, industry := range enterpriseIndustries {
		if jobInformation.EnterpriseID == industry.EnterpriseID {
			valueIn := entity.EnterpriseIndustry{
				Industry: industry.Industry,
			}
			jobInformation.Industries = append(jobInformation.Industries, valueIn)
		}
	}

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

	jobListing.Industries = jobInformation.Industries
	jobListing.Prefectures = jobInformation.Prefectures
	jobListing.EmploymentStatuses = jobInformation.EmploymentStatuses
	jobListing.WorkCharmPoints = jobInformation.WorkCharmPoints
	// 最終タスク情報をセット *不要な情報は削除
	jobListing.LatestTask = *entity.NewTaskForJobListing(
		latestTask.TaskGroupID,
		latestTask.JobInformationID,
		latestTask.JobSeekerID,
		latestTask.RAStaffID,
		latestTask.CAStaffID,
		latestTask.PhaseCategory,
		latestTask.PhaseSubCategory,
		latestTask.StaffType,
		latestTask.ExecutedStaffID,
		latestTask.ExternalJobListingURL,
	)

	output.JobListing = jobListing

	return output, nil
}

// 請求先IDから求人情報一覧を取得する
type GetJobInformationListByBillingAddressIDInput struct {
	BillingAddressID uint
}

type GetJobInformationListByBillingAddressIDOutput struct {
	JobInformationList []*entity.JobInformation
}

func (i *JobInformationInteractorImpl) GetJobInformationListByBillingAddressID(input GetJobInformationListByBillingAddressIDInput) (GetJobInformationListByBillingAddressIDOutput, error) {
	var (
		output GetJobInformationListByBillingAddressIDOutput
		err    error
	)

	jobInformationList, err := i.jobInformationRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	targets, err := i.jobInfoTargetRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.jobInfoFeatureRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.jobInfoRequiredPCToolRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.enterpriseIndustryRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformations, err := i.jobInfoSelectionInformationRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgents, err := i.jobInfoHideToAgentRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

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

		for _, hta := range hideToAgents {
			if jobInformation.ID == hta.JobInformationID {
				value := entity.JobInformationHideToAgent{
					JobInformationID: hta.JobInformationID,
					AgentID:          hta.AgentID,
					AgentName:        hta.AgentName,
				}
				jobInformation.HideToAgents = append(jobInformation.HideToAgents, value)
			}
		}
	}

	output.JobInformationList = jobInformationList

	return output, nil
}

// 企業IDから求人情報一覧を取得する
type GetJobInformationListByEnterpriseIDInput struct {
	EnterpriseID uint
}

type GetJobInformationListByEnterpriseIDOutput struct {
	JobInformationList []*entity.JobInformation
}

func (i *JobInformationInteractorImpl) GetJobInformationListByEnterpriseID(input GetJobInformationListByEnterpriseIDInput) (GetJobInformationListByEnterpriseIDOutput, error) {
	var (
		output GetJobInformationListByEnterpriseIDOutput
		err    error
	)

	jobInformationList, err := i.jobInformationRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	targets, err := i.jobInfoTargetRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.jobInfoFeatureRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.jobInfoRequiredPCToolRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.enterpriseIndustryRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformations, err := i.jobInfoSelectionInformationRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	//hideToAgent
	hideToAgents, err := i.jobInfoHideToAgentRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

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

		for _, hta := range hideToAgents {
			if jobInformation.ID == hta.JobInformationID {
				value := entity.JobInformationHideToAgent{
					JobInformationID: hta.JobInformationID,
					AgentID:          hta.AgentID,
					AgentName:        hta.AgentName,
				}
				jobInformation.HideToAgents = append(jobInformation.HideToAgents, value)
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

	}

	output.JobInformationList = jobInformationList

	return output, nil
}

// エージェントIDから求人情報一覧を取得する
type GetJobInformationListByAgentIDInput struct {
	AgentID uint
}

type GetJobInformationListByAgentIDOutput struct {
	JobInformationList []*entity.JobInformation
}

func (i *JobInformationInteractorImpl) GetJobInformationListByAgentID(input GetJobInformationListByAgentIDInput) (GetJobInformationListByAgentIDOutput, error) {
	var (
		output               GetJobInformationListByAgentIDOutput
		err                  error
		jobInformationIDList []uint
	)

	jobInformationList, err := i.jobInformationRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, jobInformation := range jobInformationList {
		jobInformationIDList = append(jobInformationIDList, jobInformation.ID)
	}

	targets, err := i.jobInfoTargetRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.jobInfoFeatureRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

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

	requiredPCTools, err := i.jobInfoRequiredPCToolRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.GetByJobInformationIDList(jobInformationIDList)
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

	requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.enterpriseIndustryRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformations, err := i.jobInfoSelectionInformationRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgents, err := i.jobInfoHideToAgentRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByJobInformationIDList(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

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

		for _, hta := range hideToAgents {
			if jobInformation.ID == hta.JobInformationID {
				value := entity.JobInformationHideToAgent{
					JobInformationID: hta.JobInformationID,
					AgentID:          hta.AgentID,
					AgentName:        hta.AgentName,
				}
				jobInformation.HideToAgents = append(jobInformation.HideToAgents, value)
			}
		}

		// occupation
		for _, oc := range occupations {
			if jobInformation.ID == oc.JobInformationID {
				value := entity.JobInformationOccupation{
					JobInformationID: oc.JobInformationID,
					Occupation:       oc.Occupation,
				}
				jobInformation.Occupations = append(jobInformation.Occupations, value)
			}
		}

	}

	output.JobInformationList = jobInformationList

	return output, nil
}

// 求人IDを使って選考フローパターンの一覧を取得する
type GetSelectionFlowPatternListByJobInformationIDInput struct {
	JobInformationID uint
}

type GetSelectionFlowPatternListByJobInformationIDOutput struct {
	SelectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
}

func (i *JobInformationInteractorImpl) GetSelectionFlowPatternListByJobInformationID(input GetSelectionFlowPatternListByJobInformationIDInput) (GetSelectionFlowPatternListByJobInformationIDOutput, error) {
	var (
		output GetSelectionFlowPatternListByJobInformationIDOutput
		err    error
	)

	selectionFlowPatternList, err := i.jobInfoSelectionFlowPatternRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformationList, err := i.jobInfoSelectionInformationRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sfp := range selectionFlowPatternList {
		for _, si := range selectionInformationList {
			if sfp.ID == si.SelectionFlowID {
				valueSi := *si
				sfp.SelectionInformations = append(sfp.SelectionInformations, valueSi)
			}
		}
	}

	output.SelectionFlowPatternList = selectionFlowPatternList

	return output, nil
}

// 求人IDを使って選考フローパターンの一覧を取得する
type GetOpenSelectionFlowPatternListByJobInformationIDInput struct {
	JobInformationID uint
}

type GetOpenSelectionFlowPatternListByJobInformationIDOutput struct {
	SelectionFlowPatternList []*entity.JobInformationSelectionFlowPattern
}

func (i *JobInformationInteractorImpl) GetOpenSelectionFlowPatternListByJobInformationID(input GetOpenSelectionFlowPatternListByJobInformationIDInput) (GetOpenSelectionFlowPatternListByJobInformationIDOutput, error) {
	var (
		output GetOpenSelectionFlowPatternListByJobInformationIDOutput
		err    error
	)

	selectionFlowPatternList, err := i.jobInfoSelectionFlowPatternRepository.GetOpenByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformationList, err := i.jobInfoSelectionInformationRepository.GetByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sfp := range selectionFlowPatternList {
		for _, si := range selectionInformationList {
			if sfp.ID == si.SelectionFlowID {
				valueSi := *si
				sfp.SelectionInformations = append(sfp.SelectionInformations, valueSi)
			}
		}
	}

	output.SelectionFlowPatternList = selectionFlowPatternList

	return output, nil
}

// 求人IDを使って選考フローパターンの一覧を取得する
type GetSelectionFlowPatternByIDInput struct {
	SelectionFlowID uint
}

type GetSelectionFlowPatternByIDOutput struct {
	SelectionFlowPattern *entity.JobInformationSelectionFlowPattern
}

func (i *JobInformationInteractorImpl) GetSelectionFlowPatternByID(input GetSelectionFlowPatternByIDInput) (GetSelectionFlowPatternByIDOutput, error) {
	var (
		output GetSelectionFlowPatternByIDOutput
		err    error
	)

	selectionFlowPattern, err := i.jobInfoSelectionFlowPatternRepository.FindByID(input.SelectionFlowID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformationList, err := i.jobInfoSelectionInformationRepository.GetBySelectionFlowID(input.SelectionFlowID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// アンケート
	selectionQuestionnaireList, err := i.selectionQuestionnaireRepository.GetBySelectionFlowID(input.SelectionFlowID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 評価点
	evaluationPointList, err := i.evaluationPointRepository.GetBySelectionFlowID(input.SelectionFlowID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, si := range selectionInformationList {
		if selectionFlowPattern.ID == si.SelectionFlowID {
			for _, sq := range selectionQuestionnaireList {
				if si.ID == sq.SelectionInformationID {
					valueSq := *sq
					si.SelectionQuestionnaires = append(si.SelectionQuestionnaires, valueSq)

					// trueならカウントする
					if sq.IsSelfIntroduction {
						si.IsSelfIntroductionCount++
					}
					if sq.IsSelfPR {
						si.IsSelfPRCount++
					}
					if sq.IsRetireReason {
						si.IsRetireReasonCount++
					}
					if sq.IsJobChangeAxis {
						si.IsJobChangeAxisCount++
					}
					if sq.IsApplyingReason {
						si.IsApplyingReasonCount++
					}
					if sq.IsCareerVision {
						si.IsCareerVisionCount++
					}
					if sq.IsReverseQuestion {
						si.IsReverseQuestionCount++
					}
				}
			}

			// 選考と求職者の重複をチェック
			checkReInterviewSelectionInfoID := make(map[null.Int]bool)
			checkJobSeekerID := make(map[uint]bool)

			for _, ep := range evaluationPointList {
				if si.ID == uint(ep.SelectionInformationID.Int64) {
					fmt.Println(ep.LastName+ep.FirstName, ep.ID, ep.SelectionInformationID)
					fmt.Println("IsReInterview", ep.IsReInterview)

					if ep.IsPassed {
						if !checkReInterviewSelectionInfoID[ep.SelectionInformationID] || !checkJobSeekerID[ep.JobSeekerID] {
							// 求職者をチェック
							checkJobSeekerID[ep.JobSeekerID] = true

							if !ep.IsReInterview {
								// 再面談以外なら選考をチェック
								checkReInterviewSelectionInfoID[ep.SelectionInformationID] = true
							}
							si.PassedExamples = append(si.PassedExamples, *ep)
						}
					} else {
						if !checkReInterviewSelectionInfoID[ep.SelectionInformationID] || !checkJobSeekerID[ep.JobSeekerID] {
							// 求職者をチェック
							checkJobSeekerID[ep.JobSeekerID] = true

							if !ep.IsReInterview {
								// 再面談以外なら選考をチェック
								checkReInterviewSelectionInfoID[ep.SelectionInformationID] = true
							}
							si.FailureExamples = append(si.FailureExamples, *ep)
						}
					}
				}
			}

			valueSi := *si
			selectionFlowPattern.SelectionInformations = append(selectionFlowPattern.SelectionInformations, valueSi)
		}
	}
	output.SelectionFlowPattern = selectionFlowPattern

	return output, nil
}

// 選考フローパターンの作成
type CreateSelectionFlowPatternInput struct {
	CreateParam entity.CreateAndUpdateSelectionFlowPatternParam
}

type CreateSelectionFlowPatternOutput struct {
	SelectionFlowPattern *entity.JobInformationSelectionFlowPattern
}

func (i *JobInformationInteractorImpl) CreateSelectionFlowPattern(input CreateSelectionFlowPatternInput) (CreateSelectionFlowPatternOutput, error) {
	var (
		output CreateSelectionFlowPatternOutput
		err    error
	)

	// 選考フローパターンの作成
	selectionFlowPattern := entity.NewJobInformationSelectionFlowPattern(
		input.CreateParam.JobInformationID,
		input.CreateParam.PublicStatus,
		input.CreateParam.FlowTitle,
		input.CreateParam.FlowPattern,
		false,
	)

	err = i.jobInfoSelectionFlowPatternRepository.Create(selectionFlowPattern)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 選考情報を作成
	for _, si := range input.CreateParam.SelectionInformations {
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

	output.SelectionFlowPattern = selectionFlowPattern

	return output, nil
}

// 選考フローパターンの更新
type UpdateSelectionFlowPatternInput struct {
	UpdateParam     entity.CreateAndUpdateSelectionFlowPatternParam
	SelectionFlowID uint
}

type UpdateSelectionFlowPatternOutput struct {
	SelectionFlowPattern *entity.JobInformationSelectionFlowPattern
}

func (i *JobInformationInteractorImpl) UpdateSelectionFlowPattern(input UpdateSelectionFlowPatternInput) (UpdateSelectionFlowPatternOutput, error) {
	var (
		output UpdateSelectionFlowPatternOutput
		err    error
	)

	selectionFlowPattern := entity.NewJobInformationSelectionFlowPattern(
		input.UpdateParam.JobInformationID,
		input.UpdateParam.PublicStatus,
		input.UpdateParam.FlowTitle,
		input.UpdateParam.FlowPattern,
		false,
	)
	selectionFlowPattern.ID = input.SelectionFlowID
	selectionFlowPattern.SelectionInformations = input.UpdateParam.SelectionInformations

	// 選考フローパターンを更新
	err = i.jobInfoSelectionFlowPatternRepository.Update(selectionFlowPattern.ID, selectionFlowPattern)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/**
	 * 選考情報のテーブルは縦持ちだが、
	 * 選考情報のIDを選考後アンケートのテーブルが参照してるため、
	 * 「Delete→Create」ではなく「Update」でデータを更新する。
	 * ※ただし「job_information_selection_informationsテーブルにレコードがない場合は、Create
	**/
	selectionInformations, err := i.jobInfoSelectionInformationRepository.GetBySelectionFlowID(selectionFlowPattern.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if len(selectionInformations) == 0 {
		for _, si := range selectionFlowPattern.SelectionInformations {
			// SelectionFlowIDをセット
			si.SelectionFlowID = selectionFlowPattern.ID

			// レコードを作成
			err = i.jobInfoSelectionInformationRepository.Create(&si)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		for _, si := range selectionFlowPattern.SelectionInformations {
			// len(selectionInformations) != 0の場合はsi.IDを使用してアップデートする
			// レコードを更新
			err = i.jobInfoSelectionInformationRepository.Update(si.ID, &si)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.SelectionFlowPattern = selectionFlowPattern

	return output, nil
}

// 選考フローパターンの更新
type DeltedSelectionFlowPatternInput struct {
	SelectionFlowID uint
}

type DeltedSelectionFlowPatternOutput struct {
	OK bool
}

func (i *JobInformationInteractorImpl) DeltedSelectionFlowPattern(input DeltedSelectionFlowPatternInput) (DeltedSelectionFlowPatternOutput, error) {
	var (
		output DeltedSelectionFlowPatternOutput
		err    error
	)

	// 選考情報を削除
	err = i.jobInfoSelectionFlowPatternRepository.Delete(input.SelectionFlowID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 企業IDから求人情報一覧を取得する
type GetJobInformationListByIDListInput struct {
	IDList []uint
}

type GetJobInformationListByIDListOutput struct {
	JobInformationList []*entity.JobInformation
}

func (i *JobInformationInteractorImpl) GetJobInformationListByIDList(input GetJobInformationListByIDListInput) (GetJobInformationListByIDListOutput, error) {
	var (
		output GetJobInformationListByIDListOutput
		err    error
	)

	jobInformationList, err := i.jobInformationRepository.GetByIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	enterpriseIDList := getEnterpriseIDList(jobInformationList)

	targets, err := i.jobInfoTargetRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.jobInfoFeatureRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.jobInfoRequiredConditionRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.jobInfoRequiredLicenseRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.jobInfoRequiredPCToolRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.enterpriseIndustryRepository.GetByEnterpriseIDList(enterpriseIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformations, err := i.jobInfoSelectionInformationRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	//hideToAgent
	hideToAgents, err := i.jobInfoHideToAgentRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.GetByJobInformationIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

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

		for _, hta := range hideToAgents {
			if jobInformation.ID == hta.JobInformationID {
				value := entity.JobInformationHideToAgent{
					JobInformationID: hta.JobInformationID,
					AgentID:          hta.AgentID,
					AgentName:        hta.AgentName,
				}
				jobInformation.HideToAgents = append(jobInformation.HideToAgents, value)
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
	}

	output.JobInformationList = jobInformationList

	return output, nil
}

// 求人票情報を求職者のuuidを使用して取得
// 求職者のマイページで使用するapi
// 求人票情報に合わせて求職者の選考情報（フェーズ）も合わせて取得
type GetJobListingListByJobSeekerUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetJobListingListByJobSeekerUUIDOutput struct {
	NotYetEntryJobListingList    []*entity.JobListing // 未エントリー
	AcceptJobOfferJobListingList []*entity.JobListing // 内定承諾
	HoldJobOfferJobListingList   []*entity.JobListing // 内定
	SelectionJobListingList      []*entity.JobListing // 選考中
	EndJobListingList            []*entity.JobListing // 終了
}

func (i *JobInformationInteractorImpl) GetJobListingListByJobSeekerUUID(input GetJobListingListByJobSeekerUUIDInput) (GetJobListingListByJobSeekerUUIDOutput, error) {
	var (
		output GetJobListingListByJobSeekerUUIDOutput
		err    error
	)

	jobInformationList, err := i.jobInformationRepository.GetAlreadySoundOutByJobSeekerUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 打診された求人がない場合はreturnする
	if len(jobInformationList) == 0 {
		return output, nil
	}

	jobInfoIDList := getJobInformationIDList(jobInformationList)

	latestTaskList, err := i.taskRepository.GetLatestByJobSeekerUUIDAndJobInformationIDList(input.JobSeekerUUID, jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, jobInformation := range jobInformationList {
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

		for _, lt := range latestTaskList {

			if jobInformation.ID == lt.JobInformationID {
				// 最終タスク情報をセット *不要な情報は削除
				jobListing.LatestTask = *entity.NewTaskForJobListing(
					lt.TaskGroupID,
					lt.JobInformationID,
					lt.JobSeekerID,
					lt.RAStaffID,
					lt.CAStaffID,
					lt.PhaseCategory,
					lt.PhaseSubCategory,
					lt.StaffType,
					lt.ExecutedStaffID,
					lt.ExternalJobListingURL,
				)

				// 外部求人の場合は外部の求人タイトルと会社名を表示できるようにセット
				if jobListing.IsExternal {
					jobListing.Title = lt.Title
					jobListing.CompanyName = lt.CompanyName
				}
			}
		}

		jobListing.JobInformationUUID = jobInformation.UUID

		// 外部求人で求人票の設定がされていない場合は返さない既存の外部求人を求職者のマイページに表示しないための対策
		if jobListing.IsExternal && jobListing.LatestTask.ExternalJobListingURL == "" {
			fmt.Printf("外部求人の求人票が設定されていません\n求職者: %s\n求人: %s", input.JobSeekerUUID, jobListing.JobInformationUUID)
			continue
		}

		// まだ応募意思確認をしていないしていない（応募意思確認中より前のタスクの場合）
		var IsNotYetConfirmApplicationIntention = jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.Entry), true) &&
			(jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.SoundOutMask), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.InHouseNG), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.CollectResultOfMask), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.EnterpriseNG), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.SoundOutJobInformation), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.RequestShareJobSeeker), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.WithoutSoundOut), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.RequestShareJobInformation), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.ConfirmPossibility), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.Unlikely), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.CloseEntryPhaseForUnlikely), true) ||
				jobListing.LatestTask.PhaseSubCategory == null.NewInt(int64(entity.CloseEntryPhaseForWithoutSoundOut), true))

		if IsNotYetConfirmApplicationIntention {
			continue
		}

		if jobListing.LatestTask.PhaseSubCategory.Int64 >= 90 && jobListing.LatestTask.PhaseSubCategory.Int64 <= 99 {
			// 終了
			output.EndJobListingList = append(output.EndJobListingList, jobListing)
		} else if jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.AcceptJobOffer), true) {
			// 内定承諾
			output.AcceptJobOfferJobListingList = append(output.AcceptJobOfferJobListingList, jobListing)
		} else if jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.HoldJobOffer), true) {
			// 内定保留
			output.HoldJobOfferJobListingList = append(output.HoldJobOfferJobListingList, jobListing)
		} else if jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.Entry), true) {
			// 未エントリー
			output.NotYetEntryJobListingList = append(output.NotYetEntryJobListingList, jobListing)
		} else if jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.DocumentSelection), true) ||
			jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.FirstSelection), true) ||
			jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.SecondSelection), true) ||
			jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.ThirdSelection), true) ||
			jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.FourthSelection), true) ||
			jobListing.LatestTask.PhaseCategory == null.NewInt(int64(entity.FinalSelection), true) {
			// 選考中
			output.SelectionJobListingList = append(output.SelectionJobListingList, jobListing)
		}
	}

	return output, nil
}

/***********************************************************************************************************************/
// 求職者検索→求人検索 API
//

// 全ての求人(自社求人 + シェア求人 + お助け求人)
type GetJobInformationListByAgentIDAndTypeInput struct {
	AgentID    uint
	PageNumber uint
	Type       entity.JobInformationType
}

type GetJobInformationListByAgentIDAndTypeOutput struct {
	JobInformationList []*entity.JobInformation
	MaxPageNumber      uint
	IDList             []uint
	AllCount           uint
	OwnCount           uint
	AllianceCount      uint
	// HelpCount          uint
}

func (i *JobInformationInteractorImpl) GetJobInformationListByAgentIDAndType(input GetJobInformationListByAgentIDAndTypeInput) (GetJobInformationListByAgentIDAndTypeOutput, error) {
	var (
		output                     GetJobInformationListByAgentIDAndTypeOutput
		jobInformationList         []*entity.JobInformation
		allJobInformationList      []*entity.JobInformation
		ownJobInformationList      []*entity.JobInformation
		allianceJobInformationList []*entity.JobInformation
		// helpJobInformationList     []*entity.JobInformation
	)
	// 1. 求人全て取得
	// 2. 同一求人と非公開先の除外
	// 3. Typeに応じて返す求人リストを変更する
	// 4. その他必要な処理

	/************ 1. 求人全て取得 **************/

	// 募集状況（0: Open）求人の登録状況（0: 本登録）の全ての求人
	jobInformationListBeforeDuplicate, err := i.jobInformationRepository.GetActiveAllByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 特別仕様: 本番環境のみ「2: 株式会社テスト」と「3: 株式会社Motoyui（非公開求人管理用）」を除外して他社エージェントに非表示にする
	jobInformationListBeforeDuplicate = excludeTestJobInformation(jobInformationListBeforeDuplicate, input.AgentID)

	/************ 2. 同一求人と非公開先の除外 **************/

	//　同一求人の除外
	jobInformationList = excludeDuplicateJobInformation(jobInformationListBeforeDuplicate, input.AgentID)

	// 指定のAgentIDを非公開先にしている求人の非公開設定情報を取得
	hideToAgent, err := i.jobInfoHideToAgentRepository.GetHideByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if len(hideToAgent) > 0 {
		// 非公開情報と合致する求人を除外して新しく求人リストを取得
		jobInformationList = checkJobInformationByHideToAgent(jobInformationList, hideToAgent)
	}

	/************ 3. Typeに応じて返す求人リストを変更する **************/

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

	// シェア求人で「雇用形態」が「2: 派遣社員 or 3: 紹介予定派遣」のみの場合は除外する
	// 仮に表示してこの求人が稼働した場合、二重派遣で法的にアウトなるための措置
	jobInformationList, err = filterAllianceJobInformationExcludedTemporaryWorker(i, jobInformationList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 「すべて・自社・シェア・お助け」のそれぞれの件数を取得する
	for _, jobInformation := range jobInformationList {
		// 自社のカウント
		if jobInformation.AgentID == input.AgentID {
			output.OwnCount = output.OwnCount + 1
			ownJobInformationList = append(ownJobInformationList, jobInformation)

			output.AllCount = output.AllCount + 1
			allJobInformationList = append(allJobInformationList, jobInformation)
		} else if jobInformation.AgentID != input.AgentID && includeUINT(allianceIDList, jobInformation.AgentID) {
			// SecretMemoを空にする
			jobInformation.SecretMemo = ""

			output.AllianceCount = output.AllianceCount + 1
			allianceJobInformationList = append(allianceJobInformationList, jobInformation)

			output.AllCount = output.AllCount + 1
			allJobInformationList = append(allJobInformationList, jobInformation)
		}
	}

	// 検索の種類に応じて返すリストを変更
	if input.Type == entity.TypeAllJobInformation {
		jobInformationList = allJobInformationList
	} else if input.Type == entity.TypeOwnJobInformation {
		jobInformationList = ownJobInformationList
	} else if input.Type == entity.TypeAllianceJobInformation {
		jobInformationList = allianceJobInformationList
	} else {
		err = fmt.Errorf("%v:%w", err, entity.ErrRequestError)
		return output, err
	}

	/************ 4. その他必要な処理 **************/

	// IDListを返す
	for _, jobInformation := range jobInformationList {
		output.IDList = append(output.IDList, jobInformation.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobInformationListMaxPage(jobInformationList)

	// 指定ページの求人20件を取得（本番実装までは1ページあたり5件）
	jobInformationList20 := getJobInformationListWithPage(jobInformationList, input.PageNumber)

	jobInformationList20, err = setJobInformationChildTableByIDList(i, jobInformationList20)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobInformationList = jobInformationList20

	return output, nil
}

/***********************************************************************************************************************/
// 管理系API
//すべての求人情報一覧を取得する

type GetAllJobInformationInput struct {
	PageNumber uint
}

type GetAllJobInformationOutput struct {
	JobInformationList []*entity.JobInformation
	MaxPageNumber      uint
	IDList             []uint
}

func (i *JobInformationInteractorImpl) GetAllJobInformation(input GetAllJobInformationInput) (GetAllJobInformationOutput, error) {
	var (
		output GetAllJobInformationOutput
		err    error
	)

	jobInformationList, err := i.jobInformationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	targets, err := i.jobInfoTargetRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.jobInfoFeatureRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.jobInfoPrefectureRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.jobInfoWorkCharmPointRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.jobInfoEmploymentStatusRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.jobInfoRequiredConditionRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.jobInfoRequiredLicenseRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.jobInfoRequiredPCToolRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.jobInfoRequiredLanguageRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.jobInfoRequiredLanguageTypeRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.jobInfoRequiredExperienceDevelopmentRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.jobInfoRequiredExperienceDevelopmentTypeRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.jobInfoRequiredExperienceJobRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.jobInfoRequiredExperienceIndustryRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.jobInfoRequiredExperienceOccupationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.jobInfoRequiredSocialExperienceRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatterns, err := i.jobInfoSelectionFlowPatternRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionInformations, err := i.jobInfoSelectionInformationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgents, err := i.jobInfoHideToAgentRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.jobInfoOccupationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

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

		for _, hta := range hideToAgents {
			if jobInformation.ID == hta.JobInformationID {
				value := entity.JobInformationHideToAgent{
					JobInformationID: hta.JobInformationID,
					AgentID:          hta.AgentID,
					AgentName:        hta.AgentName,
				}
				jobInformation.HideToAgents = append(jobInformation.HideToAgents, value)
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
	}

	// IDListを返す
	for _, jobInformation := range jobInformationList {
		output.IDList = append(output.IDList, jobInformation.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobInformationListMaxPage(jobInformationList)

	// 指定ページの求人20件を取得（本番実装までは1ページあたり5件）
	jobInformationList20 := getJobInformationListWithPage(jobInformationList, input.PageNumber)

	if len(jobInformationList20) > 0 {

		idList, enterpriseIDList := getJobInformationIDListAndEnterpriseIDList(jobInformationList20)

		// 求人リストから、IDが合致する子テーブルの情報を取得
		// リスト表示しない情報は、取得しない

		prefecture, err := i.jobInfoPrefectureRepository.GetByJobInformationIDList(idList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		industries, err := i.enterpriseIndustryRepository.GetByEnterpriseIDList(enterpriseIDList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, jobInformation := range jobInformationList20 {
			for _, p := range prefecture {
				if jobInformation.ID == p.JobInformationID {
					value := entity.JobInformationPrefecture{
						JobInformationID: p.JobInformationID,
						Prefecture:       p.Prefecture,
					}

					jobInformation.Prefectures = append(jobInformation.Prefectures, value)
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
		}
	}

	output.JobInformationList = jobInformationList20

	return output, nil
}
