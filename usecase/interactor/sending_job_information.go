package interactor

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationInteractor interface {
	// 汎用系 API
	CreateSendingJobInformation(input CreateSendingJobInformationInput) (CreateSendingJobInformationOutput, error)
	UpdateSendingJobInformation(input UpdateSendingJobInformationInput) (UpdateSendingJobInformationOutput, error)
	DeleteSendingJobInformation(input DeleteSendingJobInformationInput) (DeleteSendingJobInformationOutput, error)

	GetSendingJobInformationByID(input GetSendingJobInformationByIDInput) (GetSendingJobInformationByIDOutput, error)
	GetSendingJobInformationByUUID(input GetSendingJobInformationByUUIDInput) (GetSendingJobInformationByUUIDOutput, error)
	GetJobListingBySendingJobInformationUUID(input GetJobListingBySendingJobInformationUUIDInput) (GetJobListingBySendingJobInformationUUIDOutput, error)
	GetSendingJobInformationListBySendingEnterpriseID(input GetSendingJobInformationListBySendingEnterpriseIDInput) (GetSendingJobInformationListBySendingEnterpriseIDOutput, error)

	// ページネーション系 API
	GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetInput) (GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetOutput, error)
	GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetInput) (GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetOutput, error)

	// 求人の絞り込み検索
	// GetSearchAllSendingJobInformationByPage(input GetSearchAllSendingJobInformationByPageInput) (GetSearchAllSendingJobInformationByPageOutput, error)

	// CSV API
	ImportSendingJobInformationCSV(input ImportSendingJobInformationCSVInput) (ImportSendingJobInformationCSVOutput, error)
}

type SendingJobInformationInteractorImpl struct {
	firebase                                                  usecase.Firebase
	sendgrid                                                  config.Sendgrid
	sendingJobInformationRepository                           usecase.SendingJobInformationRepository
	sendingEnterpriseRepository                               usecase.SendingEnterpriseRepository
	sendingJobInfoTargetRepository                            usecase.SendingJobInformationTargetRepository
	sendingJobInfoFeatureRepository                           usecase.SendingJobInformationFeatureRepository
	sendingJobInfoPrefectureRepository                        usecase.SendingJobInformationPrefectureRepository
	sendingJobInfoWorkCharmPointRepository                    usecase.SendingJobInformationWorkCharmPointRepository
	sendingJobInfoEmploymentStatusRepository                  usecase.SendingJobInformationEmploymentStatusRepository
	sendingJobInfoRequiredLicenseRepository                   usecase.SendingJobInformationRequiredLicenseRepository
	sendingJobInfoRequiredPCToolRepository                    usecase.SendingJobInformationRequiredPCToolRepository
	sendingJobInfoRequiredLanguageRepository                  usecase.SendingJobInformationRequiredLanguageRepository
	sendingJobInfoRequiredExperienceDevelopmentRepository     usecase.SendingJobInformationRequiredExperienceDevelopmentRepository
	sendingJobInfoRequiredExperienceJobRepository             usecase.SendingJobInformationRequiredExperienceJobRepository
	sendingJobInfoRequiredExperienceIndustryRepository        usecase.SendingJobInformationRequiredExperienceIndustryRepository
	sendingJobInfoRequiredExperienceOccupationRepository      usecase.SendingJobInformationRequiredExperienceOccupationRepository
	sendingJobInfoRequiredSocialExperienceRepository          usecase.SendingJobInformationRequiredSocialExperienceRepository
	sendingJobInfoOccupationRepository                        usecase.SendingJobInformationOccupationRepository
	sendingJobInfoRequiredConditionRepository                 usecase.SendingJobInformationRequiredConditionRepository
	sendingJobInfoRequiredExperienceDevelopmentTypeRepository usecase.SendingJobInformationRequiredExperienceDevelopmentTypeRepository
	sendingJobInfoRequiredLanguageTypeRepository              usecase.SendingJobInformationRequiredLanguageTypeRepository
	sendingJobInfoIndustryRepository                          usecase.SendingJobInformationIndustryRepository
}

// SendingJobInformationInteractorImpl is an implementation of SendingJobInformationInteractor
func NewSendingJobInformationInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	sjR usecase.SendingJobInformationRepository,
	seR usecase.SendingEnterpriseRepository,
	sjtR usecase.SendingJobInformationTargetRepository,
	sjfR usecase.SendingJobInformationFeatureRepository,
	sjpR usecase.SendingJobInformationPrefectureRepository,
	sjwcpR usecase.SendingJobInformationWorkCharmPointRepository,
	sjesR usecase.SendingJobInformationEmploymentStatusRepository,
	sjrlR usecase.SendingJobInformationRequiredLicenseRepository,
	sjrptR usecase.SendingJobInformationRequiredPCToolRepository,
	sjrlgR usecase.SendingJobInformationRequiredLanguageRepository,
	sjredR usecase.SendingJobInformationRequiredExperienceDevelopmentRepository,
	sjrejR usecase.SendingJobInformationRequiredExperienceJobRepository,
	sjreiR usecase.SendingJobInformationRequiredExperienceIndustryRepository,
	sjreoR usecase.SendingJobInformationRequiredExperienceOccupationRepository,
	sjrseR usecase.SendingJobInformationRequiredSocialExperienceRepository,
	sjoR usecase.SendingJobInformationOccupationRepository,
	sjrcR usecase.SendingJobInformationRequiredConditionRepository,
	sjiredtR usecase.SendingJobInformationRequiredExperienceDevelopmentTypeRepository,
	sjirltR usecase.SendingJobInformationRequiredLanguageTypeRepository,
	sjiR usecase.SendingJobInformationIndustryRepository,
) SendingJobInformationInteractor {
	return &SendingJobInformationInteractorImpl{
		firebase:                                                  fb,
		sendgrid:                                                  sg,
		sendingJobInformationRepository:                           sjR,
		sendingEnterpriseRepository:                               seR,
		sendingJobInfoTargetRepository:                            sjtR,
		sendingJobInfoFeatureRepository:                           sjfR,
		sendingJobInfoPrefectureRepository:                        sjpR,
		sendingJobInfoWorkCharmPointRepository:                    sjwcpR,
		sendingJobInfoEmploymentStatusRepository:                  sjesR,
		sendingJobInfoRequiredLicenseRepository:                   sjrlR,
		sendingJobInfoRequiredPCToolRepository:                    sjrptR,
		sendingJobInfoRequiredLanguageRepository:                  sjrlgR,
		sendingJobInfoRequiredExperienceDevelopmentRepository:     sjredR,
		sendingJobInfoRequiredExperienceJobRepository:             sjrejR,
		sendingJobInfoRequiredExperienceIndustryRepository:        sjreiR,
		sendingJobInfoRequiredExperienceOccupationRepository:      sjreoR,
		sendingJobInfoRequiredSocialExperienceRepository:          sjrseR,
		sendingJobInfoOccupationRepository:                        sjoR,
		sendingJobInfoRequiredConditionRepository:                 sjrcR,
		sendingJobInfoRequiredExperienceDevelopmentTypeRepository: sjiredtR,
		sendingJobInfoRequiredLanguageTypeRepository:              sjirltR,
		sendingJobInfoIndustryRepository:                          sjiR,
	}
}

/****************************************************************************************/
/// Admin API
//

type CreateSendingJobInformationInput struct {
	CreateParam entity.CreateSendingJobInformationParam
}

type CreateSendingJobInformationOutput struct {
	SendingJobInformation *entity.SendingJobInformation
}

func (i *SendingJobInformationInteractorImpl) CreateSendingJobInformation(input CreateSendingJobInformationInput) (CreateSendingJobInformationOutput, error) {
	var (
		output CreateSendingJobInformationOutput
		err    error
	)

	sendingJobInformation := entity.NewSendingJobInformation(
		uint(input.CreateParam.SendingBillingAddressID.Int64),
		input.CreateParam.CompanyName,
		input.CreateParam.Title,
		input.CreateParam.RecruitmentState,
		input.CreateParam.ExpirationDate,
		input.CreateParam.Background,
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
		input.CreateParam.Overtime,
		input.CreateParam.OvertimeAverage,
		input.CreateParam.FixedOvertime,
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
		input.CreateParam.OtherRequired,
		input.CreateParam.Appearance,
		input.CreateParam.Communication,
		input.CreateParam.Thinking,
		input.CreateParam.TargetDetail,
		input.CreateParam.RequiredExperienceJobDetail,
		input.CreateParam.EmploymentInsurance,
		input.CreateParam.AccidentInsurance,
		input.CreateParam.HealthInsurance,
		input.CreateParam.PensionInsurance,
		input.CreateParam.RegisterPhase,
		input.CreateParam.StudyCategory,
		input.CreateParam.WordSkill,
		input.CreateParam.ExcelSkill,
		input.CreateParam.PowerPointSkill,
		input.CreateParam.CorporateSiteURL,
		input.CreateParam.PostCode,
		input.CreateParam.OfficeLocation,
		input.CreateParam.Establishment,
		input.CreateParam.EmployeeNumberSingle,
		input.CreateParam.EmployeeNumberGroup,
		input.CreateParam.PublicOffering,
		input.CreateParam.EarningsYear,
		input.CreateParam.Earnings,
		input.CreateParam.BusinessDetail,
	)

	err = i.sendingJobInformationRepository.Create(sendingJobInformation)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, target := range input.CreateParam.Targets {
		t := entity.NewSendingJobInformationTarget(
			sendingJobInformation.ID,
			target.Target,
		)

		err = i.sendingJobInfoTargetRepository.Create(t)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, feature := range input.CreateParam.Features {
		f := entity.NewSendingJobInformationFeature(
			sendingJobInformation.ID,
			feature.Feature,
		)

		err = i.sendingJobInfoFeatureRepository.Create(f)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, prefecture := range input.CreateParam.Prefectures {
		p := entity.NewSendingJobInformationPrefecture(
			sendingJobInformation.ID,
			prefecture.Prefecture,
		)

		err = i.sendingJobInfoPrefectureRepository.Create(p)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, workCharmPoint := range input.CreateParam.WorkCharmPoints {
		wcp := entity.NewSendingJobInformationWorkCharmPoint(
			sendingJobInformation.ID,
			workCharmPoint.Title,
			workCharmPoint.Contents,
		)

		err = i.sendingJobInfoWorkCharmPointRepository.Create(wcp)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, employmentStatus := range input.CreateParam.EmploymentStatuses {
		es := entity.NewSendingJobInformationEmploymentStatus(
			sendingJobInformation.ID,
			employmentStatus.EmploymentStatus,
		)

		err = i.sendingJobInfoEmploymentStatusRepository.Create(es)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 必要条件
	for _, requiredCondition := range input.CreateParam.RequiredConditions {
		rc := entity.NewSendingJobInformationRequiredCondition(
			sendingJobInformation.ID,
			requiredCondition.IsCommon,
			requiredCondition.RequiredManagement,
		)

		err = i.sendingJobInfoRequiredConditionRepository.Create(rc)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		requiredCondition.ID = rc.ID

		// 必要資格　複数
		for _, requiredLicense := range requiredCondition.RequiredLicenses {
			rl := entity.NewSendingJobInformationRequiredLicense(
				requiredCondition.ID,
				requiredLicense.License,
			)

			err = i.sendingJobInfoRequiredLicenseRepository.Create(rl)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 必要PCツール　単数
		for _, requiredPCTool := range requiredCondition.RequiredPCTools {
			rpt := entity.NewSendingJobInformationRequiredPCTool(
				requiredCondition.ID,
				requiredPCTool.Tool,
			)

			err = i.sendingJobInfoRequiredPCToolRepository.Create(rpt)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 必要言語スキル　単数
		requiredLanguages := entity.NewSendingJobInformationRequiredLanguage(
			requiredCondition.ID,
			requiredCondition.RequiredLanguages.LanguageLevel,
			requiredCondition.RequiredLanguages.Toeic,
			requiredCondition.RequiredLanguages.ToeflIBT,
			requiredCondition.RequiredLanguages.ToeflPBT,
		)

		err = i.sendingJobInfoRequiredLanguageRepository.Create(requiredLanguages)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, languageType := range requiredCondition.RequiredLanguages.LanguageTypes {
			lt := entity.NewSendingJobInformationRequiredLanguageType(
				requiredLanguages.ID,
				languageType.LanguageType,
			)

			err = i.sendingJobInfoRequiredLanguageTypeRepository.Create(lt)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 開発スキル　　言語,OS 各1つずつ
		for _, requiredExperienceDevelopment := range requiredCondition.RequiredExperienceDevelopments {
			red := entity.NewSendingJobInformationRequiredExperienceDevelopment(
				requiredCondition.ID,
				requiredExperienceDevelopment.DevelopmentCategory,
				requiredExperienceDevelopment.ExperienceYear,
				requiredExperienceDevelopment.ExperienceMonth,
			)

			err = i.sendingJobInfoRequiredExperienceDevelopmentRepository.Create(red)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, experienceDevelopmentType := range requiredExperienceDevelopment.ExperienceDevelopmentTypes {
				dt := entity.NewSendingJobInformationRequiredExperienceDevelopmentType(
					red.ID,
					experienceDevelopmentType.DevelopmentType,
				)

				err = i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.Create(dt)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}

		// 必要経験　　単数
		requiredExperienceJobs := entity.NewSendingJobInformationRequiredExperienceJob(
			requiredCondition.ID,
			requiredCondition.RequiredExperienceJobs.ExperienceYear,
			requiredCondition.RequiredExperienceJobs.ExperienceMonth,
		)

		err = i.sendingJobInfoRequiredExperienceJobRepository.Create(requiredExperienceJobs)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, experienceIndustry := range requiredCondition.RequiredExperienceJobs.ExperienceIndustries {
			ei := entity.NewSendingJobInformationRequiredExperienceIndustry(
				requiredExperienceJobs.ID,
				experienceIndustry.ExperienceIndustry,
			)

			err = i.sendingJobInfoRequiredExperienceIndustryRepository.Create(ei)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, experienceOccupation := range requiredCondition.RequiredExperienceJobs.ExperienceOccupations {
			eo := entity.NewSendingJobInformationRequiredExperienceOccupation(
				requiredExperienceJobs.ID,
				experienceOccupation.ExperienceOccupation,
			)

			err = i.sendingJobInfoRequiredExperienceOccupationRepository.Create(eo)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	for _, requiredSocialExperience := range input.CreateParam.RequiredSocialExperiences {
		rse := entity.NewSendingJobInformationRequiredSocialExperience(
			sendingJobInformation.ID,
			requiredSocialExperience.SocialExperienceType,
		)

		err = i.sendingJobInfoRequiredSocialExperienceRepository.Create(rse)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 募集職種
	for _, occupation := range input.CreateParam.Occupations {
		oc := entity.NewSendingJobInformationOccupation(
			sendingJobInformation.ID,
			occupation.Occupation,
		)

		err = i.sendingJobInfoOccupationRepository.Create(oc)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 業界
	for _, industry := range input.CreateParam.Industries {
		ind := entity.NewSendingJobInformationIndustry(
			sendingJobInformation.ID,
			industry.Industry,
		)

		err = i.sendingJobInfoIndustryRepository.Create(ind)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.SendingJobInformation = sendingJobInformation
	output.SendingJobInformation.Targets = input.CreateParam.Targets
	output.SendingJobInformation.Prefectures = input.CreateParam.Prefectures
	output.SendingJobInformation.Features = input.CreateParam.Features
	output.SendingJobInformation.WorkCharmPoints = input.CreateParam.WorkCharmPoints
	output.SendingJobInformation.EmploymentStatuses = input.CreateParam.EmploymentStatuses
	output.SendingJobInformation.RequiredConditions = input.CreateParam.RequiredConditions
	output.SendingJobInformation.RequiredSocialExperiences = input.CreateParam.RequiredSocialExperiences
	output.SendingJobInformation.Occupations = input.CreateParam.Occupations
	output.SendingJobInformation.Industries = input.CreateParam.Industries

	return output, nil
}

// 求人の更新
type UpdateSendingJobInformationInput struct {
	UpdateParam             entity.UpdateSendingJobInformationParam
	SendingJobInformationID uint
}

type UpdateSendingJobInformationOutput struct {
	SendingJobInformation *entity.SendingJobInformation
}

func (i *SendingJobInformationInteractorImpl) UpdateSendingJobInformation(input UpdateSendingJobInformationInput) (UpdateSendingJobInformationOutput, error) {
	var (
		output UpdateSendingJobInformationOutput
		err    error
	)

	sendingJobInformation := entity.NewSendingJobInformation(
		input.UpdateParam.SendingBillingAddressID,
		input.UpdateParam.CompanyName,
		input.UpdateParam.Title,
		input.UpdateParam.RecruitmentState,
		input.UpdateParam.ExpirationDate,
		input.UpdateParam.Background,
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
		input.UpdateParam.Overtime,
		input.UpdateParam.OvertimeAverage,
		input.UpdateParam.FixedOvertime,
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
		input.UpdateParam.OtherRequired,
		input.UpdateParam.Appearance,
		input.UpdateParam.Communication,
		input.UpdateParam.Thinking,
		input.UpdateParam.TargetDetail,
		input.UpdateParam.RequiredExperienceJobDetail,
		input.UpdateParam.EmploymentInsurance,
		input.UpdateParam.AccidentInsurance,
		input.UpdateParam.HealthInsurance,
		input.UpdateParam.PensionInsurance,
		input.UpdateParam.RegisterPhase,
		input.UpdateParam.StudyCategory,
		input.UpdateParam.WordSkill,
		input.UpdateParam.ExcelSkill,
		input.UpdateParam.PowerPointSkill,
		input.UpdateParam.CorporateSiteURL,
		input.UpdateParam.PostCode,
		input.UpdateParam.OfficeLocation,
		input.UpdateParam.Establishment,
		input.UpdateParam.EmployeeNumberSingle,
		input.UpdateParam.EmployeeNumberGroup,
		input.UpdateParam.PublicOffering,
		input.UpdateParam.EarningsYear,
		input.UpdateParam.Earnings,
		input.UpdateParam.BusinessDetail,
	)

	err = i.sendingJobInformationRepository.Update(sendingJobInformation, input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.sendingJobInfoTargetRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, target := range input.UpdateParam.Targets {
		t := entity.NewSendingJobInformationTarget(
			input.SendingJobInformationID,
			target.Target,
		)

		err = i.sendingJobInfoTargetRepository.Create(t)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobInfoFeatureRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, feature := range input.UpdateParam.Features {
		f := entity.NewSendingJobInformationFeature(
			input.SendingJobInformationID,
			feature.Feature,
		)

		err = i.sendingJobInfoFeatureRepository.Create(f)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobInfoPrefectureRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, prefecture := range input.UpdateParam.Prefectures {
		p := entity.NewSendingJobInformationPrefecture(
			input.SendingJobInformationID,
			prefecture.Prefecture,
		)

		err = i.sendingJobInfoPrefectureRepository.Create(p)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobInfoWorkCharmPointRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, workCharmPoint := range input.UpdateParam.WorkCharmPoints {
		wcp := entity.NewSendingJobInformationWorkCharmPoint(
			input.SendingJobInformationID,
			workCharmPoint.Title,
			workCharmPoint.Contents,
		)

		err = i.sendingJobInfoWorkCharmPointRepository.Create(wcp)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobInfoEmploymentStatusRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, employmentStatus := range input.UpdateParam.EmploymentStatuses {
		es := entity.NewSendingJobInformationEmploymentStatus(
			input.SendingJobInformationID,
			employmentStatus.EmploymentStatus,
		)

		err = i.sendingJobInfoEmploymentStatusRepository.Create(es)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 必要条件以下　削除
	err = i.sendingJobInfoRequiredConditionRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, requiredCondition := range input.UpdateParam.RequiredConditions {
		rc := entity.NewSendingJobInformationRequiredCondition(
			input.SendingJobInformationID,
			requiredCondition.IsCommon,
			requiredCondition.RequiredManagement,
		)

		err = i.sendingJobInfoRequiredConditionRepository.Create(rc)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		requiredCondition.ID = rc.ID

		// 必要資格　複数
		for _, requiredLicense := range requiredCondition.RequiredLicenses {
			rl := entity.NewSendingJobInformationRequiredLicense(
				requiredCondition.ID,
				requiredLicense.License,
			)

			err = i.sendingJobInfoRequiredLicenseRepository.Create(rl)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 必要PCツール　単数
		for _, requiredPCTool := range requiredCondition.RequiredPCTools {
			rpt := entity.NewSendingJobInformationRequiredPCTool(
				requiredCondition.ID,
				requiredPCTool.Tool,
			)

			err = i.sendingJobInfoRequiredPCToolRepository.Create(rpt)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 必要言語スキル　単数
		requiredLanguages := entity.NewSendingJobInformationRequiredLanguage(
			requiredCondition.ID,
			requiredCondition.RequiredLanguages.LanguageLevel,
			requiredCondition.RequiredLanguages.Toeic,
			requiredCondition.RequiredLanguages.ToeflIBT,
			requiredCondition.RequiredLanguages.ToeflPBT,
		)

		err = i.sendingJobInfoRequiredLanguageRepository.Create(requiredLanguages)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, languageType := range requiredCondition.RequiredLanguages.LanguageTypes {
			lt := entity.NewSendingJobInformationRequiredLanguageType(
				requiredLanguages.ID,
				languageType.LanguageType,
			)

			err = i.sendingJobInfoRequiredLanguageTypeRepository.Create(lt)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// 開発スキル　　言語,OS 各1つずつ
		for _, requiredExperienceDevelopment := range requiredCondition.RequiredExperienceDevelopments {
			red := entity.NewSendingJobInformationRequiredExperienceDevelopment(
				requiredCondition.ID,
				requiredExperienceDevelopment.DevelopmentCategory,
				requiredExperienceDevelopment.ExperienceYear,
				requiredExperienceDevelopment.ExperienceMonth,
			)

			err = i.sendingJobInfoRequiredExperienceDevelopmentRepository.Create(red)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, experienceDevelopmentType := range requiredExperienceDevelopment.ExperienceDevelopmentTypes {
				dt := entity.NewSendingJobInformationRequiredExperienceDevelopmentType(
					red.ID,
					experienceDevelopmentType.DevelopmentType,
				)

				err = i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.Create(dt)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}

		// 必要経験　　単数
		requiredExperienceJobs := entity.NewSendingJobInformationRequiredExperienceJob(
			requiredCondition.ID,
			requiredCondition.RequiredExperienceJobs.ExperienceYear,
			requiredCondition.RequiredExperienceJobs.ExperienceMonth,
		)

		err = i.sendingJobInfoRequiredExperienceJobRepository.Create(requiredExperienceJobs)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, experienceIndustry := range requiredCondition.RequiredExperienceJobs.ExperienceIndustries {
			ei := entity.NewSendingJobInformationRequiredExperienceIndustry(
				requiredExperienceJobs.ID,
				experienceIndustry.ExperienceIndustry,
			)

			err = i.sendingJobInfoRequiredExperienceIndustryRepository.Create(ei)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, experienceOccupation := range requiredCondition.RequiredExperienceJobs.ExperienceOccupations {
			eo := entity.NewSendingJobInformationRequiredExperienceOccupation(
				requiredExperienceJobs.ID,
				experienceOccupation.ExperienceOccupation,
			)

			err = i.sendingJobInfoRequiredExperienceOccupationRepository.Create(eo)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}
	// 必要条件以上

	err = i.sendingJobInfoRequiredSocialExperienceRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, requiredSocialExperience := range input.UpdateParam.RequiredSocialExperiences {
		rse := entity.NewSendingJobInformationRequiredSocialExperience(
			input.SendingJobInformationID,
			requiredSocialExperience.SocialExperienceType,
		)

		err = i.sendingJobInfoRequiredSocialExperienceRepository.Create(rse)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 募集職種
	err = i.sendingJobInfoOccupationRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	for _, occupation := range input.UpdateParam.Occupations {
		oc := entity.NewSendingJobInformationOccupation(
			input.SendingJobInformationID,
			occupation.Occupation,
		)

		err = i.sendingJobInfoOccupationRepository.Create(oc)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 業界
	err = i.sendingJobInfoIndustryRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, industry := range input.UpdateParam.Industries {
		ind := entity.NewSendingJobInformationIndustry(
			input.SendingJobInformationID,
			industry.Industry,
		)

		err = i.sendingJobInfoIndustryRepository.Create(ind)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.SendingJobInformation = sendingJobInformation
	output.SendingJobInformation.Targets = input.UpdateParam.Targets
	output.SendingJobInformation.Features = input.UpdateParam.Features
	output.SendingJobInformation.Prefectures = input.UpdateParam.Prefectures
	output.SendingJobInformation.WorkCharmPoints = input.UpdateParam.WorkCharmPoints
	output.SendingJobInformation.EmploymentStatuses = input.UpdateParam.EmploymentStatuses
	output.SendingJobInformation.RequiredConditions = input.UpdateParam.RequiredConditions
	output.SendingJobInformation.RequiredSocialExperiences = input.UpdateParam.RequiredSocialExperiences
	output.SendingJobInformation.Occupations = input.UpdateParam.Occupations
	output.SendingJobInformation.Industries = input.UpdateParam.Industries

	return output, nil
}

type DeleteSendingJobInformationInput struct {
	SendingJobInformationID uint
}

type DeleteSendingJobInformationOutput struct {
	OK bool
}

func (i *SendingJobInformationInteractorImpl) DeleteSendingJobInformation(input DeleteSendingJobInformationInput) (DeleteSendingJobInformationOutput, error) {
	var (
		output DeleteSendingJobInformationOutput
	)

	err := i.sendingJobInformationRepository.Delete(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 求人IDを使って求人情報を取得する
type GetSendingJobInformationByIDInput struct {
	SendingJobInformationID uint
}

type GetSendingJobInformationByIDOutput struct {
	SendingJobInformation *entity.SendingJobInformation
}

func (i *SendingJobInformationInteractorImpl) GetSendingJobInformationByID(input GetSendingJobInformationByIDInput) (GetSendingJobInformationByIDOutput, error) {
	var (
		output GetSendingJobInformationByIDOutput
		err    error
	)

	sendingJobInformation, err := i.sendingJobInformationRepository.FindByID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	targets, err := i.sendingJobInfoTargetRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.sendingJobInfoFeatureRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.sendingJobInfoPrefectureRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.sendingJobInfoWorkCharmPointRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.sendingJobInfoEmploymentStatusRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.sendingJobInfoRequiredConditionRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.sendingJobInfoRequiredLicenseRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.sendingJobInfoRequiredPCToolRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.sendingJobInfoRequiredLanguageRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.sendingJobInfoRequiredLanguageTypeRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.sendingJobInfoRequiredExperienceDevelopmentRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.sendingJobInfoRequiredExperienceJobRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.sendingJobInfoRequiredExperienceIndustryRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.sendingJobInfoRequiredExperienceOccupationRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.sendingJobInfoRequiredSocialExperienceRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.sendingJobInfoOccupationRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.sendingJobInfoIndustryRepository.GetListBySendingJobInformationID(input.SendingJobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, t := range targets {
		value := entity.SendingJobInformationTarget{
			SendingJobInformationID: sendingJobInformation.ID,
			Target:                  t.Target,
		}
		sendingJobInformation.Targets = append(sendingJobInformation.Targets, value)
	}

	for _, f := range features {
		value := entity.SendingJobInformationFeature{
			SendingJobInformationID: f.SendingJobInformationID,
			Feature:                 f.Feature,
		}
		sendingJobInformation.Features = append(sendingJobInformation.Features, value)
	}

	for _, p := range prefectures {
		value := entity.SendingJobInformationPrefecture{
			SendingJobInformationID: p.SendingJobInformationID,
			Prefecture:              p.Prefecture,
		}
		sendingJobInformation.Prefectures = append(sendingJobInformation.Prefectures, value)
	}

	for _, wcp := range workCharmPoints {
		value := entity.SendingJobInformationWorkCharmPoint{
			SendingJobInformationID: wcp.SendingJobInformationID,
			Title:                   wcp.Title,
			Contents:                wcp.Contents,
		}
		sendingJobInformation.WorkCharmPoints = append(sendingJobInformation.WorkCharmPoints, value)
	}

	for _, es := range employmentStatuses {
		value := entity.SendingJobInformationEmploymentStatus{
			SendingJobInformationID: es.SendingJobInformationID,
			EmploymentStatus:        es.EmploymentStatus,
		}
		sendingJobInformation.EmploymentStatuses = append(sendingJobInformation.EmploymentStatuses, value)
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

		sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, *condition)
	}

	for _, rse := range requiredSocialExperiences {
		value := entity.SendingJobInformationRequiredSocialExperience{
			SendingJobInformationID: rse.SendingJobInformationID,
			SocialExperienceType:    rse.SocialExperienceType,
		}
		sendingJobInformation.RequiredSocialExperiences = append(sendingJobInformation.RequiredSocialExperiences, value)
	}

	for _, oc := range occupations {
		value := entity.SendingJobInformationOccupation{
			SendingJobInformationID: oc.SendingJobInformationID,
			Occupation:              oc.Occupation,
		}
		sendingJobInformation.Occupations = append(sendingJobInformation.Occupations, value)
	}

	for _, ind := range industries {
		value := entity.SendingJobInformationIndustry{
			SendingJobInformationID: ind.SendingJobInformationID,
			Industry:                ind.Industry,
		}
		sendingJobInformation.Industries = append(sendingJobInformation.Industries, value)
	}

	output.SendingJobInformation = sendingJobInformation

	return output, nil
}

// 求人のuuidを使って求人情報を取得する
type GetSendingJobInformationByUUIDInput struct {
	UUID uuid.UUID
}

type GetSendingJobInformationByUUIDOutput struct {
	SendingJobInformation *entity.SendingJobInformation
}

func (i *SendingJobInformationInteractorImpl) GetSendingJobInformationByUUID(input GetSendingJobInformationByUUIDInput) (GetSendingJobInformationByUUIDOutput, error) {
	var (
		output GetSendingJobInformationByUUIDOutput
		err    error
	)

	sendingJobInformation, err := i.sendingJobInformationRepository.FindByUUID(input.UUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	targets, err := i.sendingJobInfoTargetRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.sendingJobInfoFeatureRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.sendingJobInfoPrefectureRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.sendingJobInfoWorkCharmPointRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.sendingJobInfoEmploymentStatusRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.sendingJobInfoRequiredConditionRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.sendingJobInfoRequiredLicenseRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.sendingJobInfoRequiredPCToolRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.sendingJobInfoRequiredLanguageRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.sendingJobInfoRequiredLanguageTypeRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.sendingJobInfoRequiredExperienceDevelopmentRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.sendingJobInfoRequiredExperienceJobRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.sendingJobInfoRequiredExperienceIndustryRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.sendingJobInfoRequiredExperienceOccupationRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.sendingJobInfoRequiredSocialExperienceRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.sendingJobInfoOccupationRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.sendingJobInfoIndustryRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, t := range targets {
		value := entity.SendingJobInformationTarget{
			SendingJobInformationID: sendingJobInformation.ID,
			Target:                  t.Target,
		}
		sendingJobInformation.Targets = append(sendingJobInformation.Targets, value)
	}

	for _, f := range features {
		value := entity.SendingJobInformationFeature{
			SendingJobInformationID: f.SendingJobInformationID,
			Feature:                 f.Feature,
		}
		sendingJobInformation.Features = append(sendingJobInformation.Features, value)
	}

	for _, p := range prefectures {
		value := entity.SendingJobInformationPrefecture{
			SendingJobInformationID: p.SendingJobInformationID,
			Prefecture:              p.Prefecture,
		}
		sendingJobInformation.Prefectures = append(sendingJobInformation.Prefectures, value)
	}

	for _, wcp := range workCharmPoints {
		value := entity.SendingJobInformationWorkCharmPoint{
			SendingJobInformationID: wcp.SendingJobInformationID,
			Title:                   wcp.Title,
			Contents:                wcp.Contents,
		}
		sendingJobInformation.WorkCharmPoints = append(sendingJobInformation.WorkCharmPoints, value)
	}

	for _, es := range employmentStatuses {
		value := entity.SendingJobInformationEmploymentStatus{
			SendingJobInformationID: es.SendingJobInformationID,
			EmploymentStatus:        es.EmploymentStatus,
		}
		sendingJobInformation.EmploymentStatuses = append(sendingJobInformation.EmploymentStatuses, value)
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

		sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, *condition)
	}

	for _, rse := range requiredSocialExperiences {
		value := entity.SendingJobInformationRequiredSocialExperience{
			SendingJobInformationID: rse.SendingJobInformationID,
			SocialExperienceType:    rse.SocialExperienceType,
		}
		sendingJobInformation.RequiredSocialExperiences = append(sendingJobInformation.RequiredSocialExperiences, value)
	}

	for _, oc := range occupations {
		value := entity.SendingJobInformationOccupation{
			SendingJobInformationID: oc.SendingJobInformationID,
			Occupation:              oc.Occupation,
		}
		sendingJobInformation.Occupations = append(sendingJobInformation.Occupations, value)
	}

	for _, ind := range industries {
		value := entity.SendingJobInformationIndustry{
			SendingJobInformationID: ind.SendingJobInformationID,
			Industry:                ind.Industry,
		}
		sendingJobInformation.Industries = append(sendingJobInformation.Industries, value)
	}

	output.SendingJobInformation = sendingJobInformation

	return output, nil
}

// 求人のuuidを使って求人情報を取得する
type GetJobListingBySendingJobInformationUUIDInput struct {
	UUID uuid.UUID
}

type GetJobListingBySendingJobInformationUUIDOutput struct {
	JobListing *entity.JobListingForSending
}

func (i *SendingJobInformationInteractorImpl) GetJobListingBySendingJobInformationUUID(input GetJobListingBySendingJobInformationUUIDInput) (GetJobListingBySendingJobInformationUUIDOutput, error) {
	var (
		output GetJobListingBySendingJobInformationUUIDOutput
		err    error
	)

	sendingJobInformation, err := i.sendingJobInformationRepository.FindByUUID(input.UUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.sendingJobInfoPrefectureRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.sendingJobInfoWorkCharmPointRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.sendingJobInfoEmploymentStatusRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.sendingJobInfoIndustryRepository.GetListBySendingJobInformationID(sendingJobInformation.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, p := range prefectures {
		value := entity.SendingJobInformationPrefecture{
			SendingJobInformationID: p.SendingJobInformationID,
			Prefecture:              p.Prefecture,
		}
		sendingJobInformation.Prefectures = append(sendingJobInformation.Prefectures, value)
	}

	for _, wcp := range workCharmPoints {
		value := entity.SendingJobInformationWorkCharmPoint{
			SendingJobInformationID: wcp.SendingJobInformationID,
			Title:                   wcp.Title,
			Contents:                wcp.Contents,
		}
		sendingJobInformation.WorkCharmPoints = append(sendingJobInformation.WorkCharmPoints, value)
	}

	for _, es := range employmentStatuses {
		value := entity.SendingJobInformationEmploymentStatus{
			SendingJobInformationID: es.SendingJobInformationID,
			EmploymentStatus:        es.EmploymentStatus,
		}
		sendingJobInformation.EmploymentStatuses = append(sendingJobInformation.EmploymentStatuses, value)
	}

	for _, ind := range industries {
		value := entity.SendingJobInformationIndustry{
			SendingJobInformationID: ind.SendingJobInformationID,
			Industry:                ind.Industry,
		}
		sendingJobInformation.Industries = append(sendingJobInformation.Industries, value)
	}

	jobListing := entity.NewJobListingForSending(
		sendingJobInformation.ID,
		sendingJobInformation.AgentStaffID,
		sendingJobInformation.CompanyName,
		sendingJobInformation.CorporateSiteURL,
		sendingJobInformation.PostCode,
		sendingJobInformation.OfficeLocation,
		sendingJobInformation.EmployeeNumberSingle,
		sendingJobInformation.EmployeeNumberGroup,
		sendingJobInformation.Establishment,
		sendingJobInformation.PublicOffering,
		sendingJobInformation.Earnings,
		sendingJobInformation.EarningsYear,
		sendingJobInformation.BusinessDetail,
		sendingJobInformation.Title,
		sendingJobInformation.Background,
		sendingJobInformation.WorkDetail,
		sendingJobInformation.WorkLocation,
		sendingJobInformation.Transfer,
		sendingJobInformation.TransferDetail,
		sendingJobInformation.UnderIncome,
		sendingJobInformation.OverIncome,
		sendingJobInformation.Salary,
		sendingJobInformation.Insurance,
		sendingJobInformation.WorkTime,
		sendingJobInformation.Overtime,
		sendingJobInformation.OvertimeAverage,
		sendingJobInformation.FixedOvertime,
		sendingJobInformation.FixedOvertimePayment,
		sendingJobInformation.FixedOvertimeDetail,
		sendingJobInformation.TrialPeriod,
		sendingJobInformation.TrialPeriodDetail,
		sendingJobInformation.EmploymentPeriod,
		sendingJobInformation.EmploymentPeriodDetail,
		sendingJobInformation.HolidayDetail,
		sendingJobInformation.PassiveSmoking,
		sendingJobInformation.SelectionFlow,
		sendingJobInformation.EmploymentInsurance,
		sendingJobInformation.AccidentInsurance,
		sendingJobInformation.HealthInsurance,
		sendingJobInformation.PensionInsurance,
	)

	jobListing.Industries = sendingJobInformation.Industries
	jobListing.Prefectures = sendingJobInformation.Prefectures
	jobListing.EmploymentStatuses = sendingJobInformation.EmploymentStatuses
	jobListing.WorkCharmPoints = sendingJobInformation.WorkCharmPoints

	output.JobListing = jobListing

	return output, nil
}

// 企業IDから求人情報一覧を取得する
type GetSendingJobInformationListBySendingEnterpriseIDInput struct {
	SendingEnterpriseID uint
}

type GetSendingJobInformationListBySendingEnterpriseIDOutput struct {
	SendingJobInformationList []*entity.SendingJobInformation
}

func (i *SendingJobInformationInteractorImpl) GetSendingJobInformationListBySendingEnterpriseID(input GetSendingJobInformationListBySendingEnterpriseIDInput) (GetSendingJobInformationListBySendingEnterpriseIDOutput, error) {
	var (
		output GetSendingJobInformationListBySendingEnterpriseIDOutput
		err    error
	)

	sendingJobInformationList, err := i.sendingJobInformationRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	targets, err := i.sendingJobInfoTargetRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.sendingJobInfoFeatureRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.sendingJobInfoPrefectureRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.sendingJobInfoWorkCharmPointRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.sendingJobInfoEmploymentStatusRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.sendingJobInfoRequiredLicenseRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.sendingJobInfoRequiredConditionRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.sendingJobInfoRequiredPCToolRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.sendingJobInfoRequiredLanguageRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.sendingJobInfoRequiredLanguageTypeRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.sendingJobInfoRequiredExperienceDevelopmentRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.sendingJobInfoRequiredExperienceJobRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.sendingJobInfoRequiredExperienceIndustryRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.sendingJobInfoRequiredExperienceOccupationRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.sendingJobInfoRequiredSocialExperienceRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.sendingJobInfoOccupationRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	industries, err := i.sendingJobInfoIndustryRepository.GetListBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sendingJobInformation := range sendingJobInformationList {
		for _, t := range targets {
			if sendingJobInformation.ID == t.SendingJobInformationID {
				value := entity.SendingJobInformationTarget{
					SendingJobInformationID: t.SendingJobInformationID,
					Target:                  t.Target,
				}

				sendingJobInformation.Targets = append(sendingJobInformation.Targets, value)
			}
		}

		for _, f := range features {
			if sendingJobInformation.ID == f.SendingJobInformationID {
				value := entity.SendingJobInformationFeature{
					SendingJobInformationID: f.SendingJobInformationID,
					Feature:                 f.Feature,
				}

				sendingJobInformation.Features = append(sendingJobInformation.Features, value)
			}
		}

		for _, p := range prefectures {
			if sendingJobInformation.ID == p.SendingJobInformationID {
				value := entity.SendingJobInformationPrefecture{
					SendingJobInformationID: p.SendingJobInformationID,
					Prefecture:              p.Prefecture,
				}

				sendingJobInformation.Prefectures = append(sendingJobInformation.Prefectures, value)
			}
		}

		for _, wcp := range workCharmPoints {
			if sendingJobInformation.ID == wcp.SendingJobInformationID {
				value := entity.SendingJobInformationWorkCharmPoint{
					SendingJobInformationID: wcp.SendingJobInformationID,
					Title:                   wcp.Title,
					Contents:                wcp.Contents,
				}

				sendingJobInformation.WorkCharmPoints = append(sendingJobInformation.WorkCharmPoints, value)
			}
		}

		for _, es := range employmentStatuses {
			if sendingJobInformation.ID == es.SendingJobInformationID {
				value := entity.SendingJobInformationEmploymentStatus{
					SendingJobInformationID: es.SendingJobInformationID,
					EmploymentStatus:        es.EmploymentStatus,
				}

				sendingJobInformation.EmploymentStatuses = append(sendingJobInformation.EmploymentStatuses, value)
			}
		}

		for _, condition := range requiredConditions {
			if sendingJobInformation.ID == condition.SendingJobInformationID {
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

				sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, *condition)
			}
		}

		for _, rse := range requiredSocialExperiences {
			if sendingJobInformation.ID == rse.SendingJobInformationID {
				value := entity.SendingJobInformationRequiredSocialExperience{
					SendingJobInformationID: rse.SendingJobInformationID,
					SocialExperienceType:    rse.SocialExperienceType,
				}
				sendingJobInformation.RequiredSocialExperiences = append(sendingJobInformation.RequiredSocialExperiences, value)
			}
		}

		for _, oc := range occupations {
			if sendingJobInformation.ID == oc.SendingJobInformationID {
				value := entity.SendingJobInformationOccupation{
					SendingJobInformationID: oc.SendingJobInformationID,
					Occupation:              oc.Occupation,
				}
				sendingJobInformation.Occupations = append(sendingJobInformation.Occupations, value)
			}
		}

		for _, ind := range industries {
			value := entity.SendingJobInformationIndustry{
				SendingJobInformationID: ind.SendingJobInformationID,
				Industry:                ind.Industry,
			}
			sendingJobInformation.Industries = append(sendingJobInformation.Industries, value)
		}

	}

	output.SendingJobInformationList = sendingJobInformationList

	return output, nil
}

// まだ送客していない送客先とその送客先が保有する求人数の一覧を取得するapi
type GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetInput struct {
	SendingJobSeekerID  uint
	SendingEnterpriseID uint
	PageNumber          uint
}

type GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetOutput struct {
	SendingEnterpriseList     []*entity.SendingEnterprise
	SendingJobInformationList []*entity.SendingJobInformation
	MaxPageNumber             uint
}

func (i *SendingJobInformationInteractorImpl) GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYet(input GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetInput) (GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetOutput, error) {
	var (
		output GetSendingEnterpriseListAndJobInfoPageListByHaveNotSentYetOutput
		err    error

		sendingEnterpriseID uint
		jobInfoIDList       []uint
	)

	/************ 企業情報を取得 **************/

	// 企業情報と求人数を取得
	sendingEnterpriseList, err := i.sendingEnterpriseRepository.GetSendingEnterpriseAndJobInfoCountByHaveNotSentYetBySendingJobSeekerID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 企業情報が取得できない場合は処理終了
	if len(sendingEnterpriseList) == 0 {
		return output, nil
	}

	output.SendingEnterpriseList = sendingEnterpriseList

	/************ 求人情報を取得 **************/

	// 初回レンダリング時はinput.SendingEnterpriseID==0になるためその場合の処理を記述
	if input.SendingEnterpriseID == 0 {
		sendingEnterpriseID = sendingEnterpriseList[0].ID
	} else {
		sendingEnterpriseID = input.SendingEnterpriseID
	}

	sendingJobInformationList, err := i.sendingJobInformationRepository.GetListBySendingEnterpriseID(sendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, jobInfo := range sendingJobInformationList {
		jobInfoIDList = append(jobInfoIDList, jobInfo.ID)
	}

	/************ 求人の子テーブル情報を取得 **************/

	targets, err := i.sendingJobInfoTargetRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.sendingJobInfoFeatureRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.sendingJobInfoPrefectureRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.sendingJobInfoWorkCharmPointRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.sendingJobInfoEmploymentStatusRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.sendingJobInfoRequiredConditionRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.sendingJobInfoRequiredLicenseRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.sendingJobInfoRequiredPCToolRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.sendingJobInfoRequiredLanguageRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.sendingJobInfoRequiredLanguageTypeRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.sendingJobInfoRequiredExperienceDevelopmentRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.sendingJobInfoRequiredExperienceJobRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.sendingJobInfoRequiredExperienceIndustryRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.sendingJobInfoRequiredExperienceOccupationRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.sendingJobInfoRequiredSocialExperienceRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.sendingJobInfoOccupationRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sendingJobInformation := range sendingJobInformationList {
		for _, t := range targets {
			if sendingJobInformation.ID == t.SendingJobInformationID {
				value := entity.SendingJobInformationTarget{
					SendingJobInformationID: sendingJobInformation.ID,
					Target:                  t.Target,
				}
				sendingJobInformation.Targets = append(sendingJobInformation.Targets, value)
			}
		}

		for _, f := range features {
			if sendingJobInformation.ID == f.SendingJobInformationID {
				value := entity.SendingJobInformationFeature{
					SendingJobInformationID: f.SendingJobInformationID,
					Feature:                 f.Feature,
				}
				sendingJobInformation.Features = append(sendingJobInformation.Features, value)
			}
		}

		for _, p := range prefectures {
			if sendingJobInformation.ID == p.SendingJobInformationID {
				value := entity.SendingJobInformationPrefecture{
					SendingJobInformationID: p.SendingJobInformationID,
					Prefecture:              p.Prefecture,
				}
				sendingJobInformation.Prefectures = append(sendingJobInformation.Prefectures, value)
			}
		}

		for _, wcp := range workCharmPoints {
			if sendingJobInformation.ID == wcp.SendingJobInformationID {
				value := entity.SendingJobInformationWorkCharmPoint{
					SendingJobInformationID: wcp.SendingJobInformationID,
					Title:                   wcp.Title,
					Contents:                wcp.Contents,
				}
				sendingJobInformation.WorkCharmPoints = append(sendingJobInformation.WorkCharmPoints, value)
			}
		}

		for _, es := range employmentStatuses {
			if sendingJobInformation.ID == es.SendingJobInformationID {
				value := entity.SendingJobInformationEmploymentStatus{
					SendingJobInformationID: es.SendingJobInformationID,
					EmploymentStatus:        es.EmploymentStatus,
				}
				sendingJobInformation.EmploymentStatuses = append(sendingJobInformation.EmploymentStatuses, value)
			}
		}

		// 必要条件
		for _, condition := range requiredConditions {
			if sendingJobInformation.ID == condition.SendingJobInformationID {
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

				sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, *condition)
			}
		}

		for _, rse := range requiredSocialExperiences {
			if sendingJobInformation.ID == rse.SendingJobInformationID {
				value := entity.SendingJobInformationRequiredSocialExperience{
					SendingJobInformationID: rse.SendingJobInformationID,
					SocialExperienceType:    rse.SocialExperienceType,
				}
				sendingJobInformation.RequiredSocialExperiences = append(sendingJobInformation.RequiredSocialExperiences, value)
			}
		}

		for _, oc := range occupations {
			if sendingJobInformation.ID == oc.SendingJobInformationID {
				value := entity.SendingJobInformationOccupation{
					SendingJobInformationID: oc.SendingJobInformationID,
					Occupation:              oc.Occupation,
				}
				sendingJobInformation.Occupations = append(sendingJobInformation.Occupations, value)
			}
		}
	}

	/************ ページ関連の処理 **************/

	// ページの最大数を取得
	output.MaxPageNumber = getSendingJobInformationListMaxPage(sendingJobInformationList)

	// 指定ページの求人20件を取得（本番実装までは1ページあたり5件）
	sendingJobInformationList20 := getSendingJobInformationListWithPage(sendingJobInformationList, input.PageNumber)

	output.SendingJobInformationList = sendingJobInformationList20

	return output, nil
}

// 絞り込み まだ送客していない送客先とその送客先が保有する求人数の一覧を取得するapi
type GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetInput struct {
	SendingJobSeekerID  uint
	SendingEnterpriseID uint
	PageNumber          uint
	SearchParam         entity.SearchJobInformation
}

type GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetOutput struct {
	SendingEnterpriseList     []*entity.SendingEnterprise
	SendingJobInformationList []*entity.SendingJobInformation
	MaxPageNumber             uint
}

func (i *SendingJobInformationInteractorImpl) GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYet(input GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetInput) (GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetOutput, error) {
	var (
		output                          GetSendingEnterpriseListAndJobInfoPageSearchListByHaveNotSentYetOutput
		err                             error
		sendingJobInformationList       []*entity.SendingJobInformation
		sendingEnterpriseID             uint
		enterpriseIDList                []uint
		jobInfoIDList                   []uint
		selectSendingJobInformationList []*entity.SendingJobInformation
	)

	/************ 企業情報を取得 **************/

	sendingEnterpriseList, err := i.sendingEnterpriseRepository.GetSendingEnterpriseByHaveNotSentYetBySendingJobSeekerID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 企業情報が取得できない場合は処理終了
	if len(sendingEnterpriseList) == 0 {
		return output, nil
	}

	output.SendingEnterpriseList = sendingEnterpriseList

	// 企業IDのリストを作成
	for _, enterprise := range sendingEnterpriseList {
		enterpriseIDList = append(enterpriseIDList, enterprise.ID)
	}

	/************ 求人情報を取得 **************/

	sendingJobInformationList, err = i.sendingJobInformationRepository.GetListBySendingEnterpriseIDListAndFreeWord(enterpriseIDList, input.SearchParam.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sendingJobInformation := range sendingJobInformationList {
		jobInfoIDList = append(jobInfoIDList, sendingJobInformation.ID)
	}

	/************ 求人の子テーブル情報を取得 **************/

	targets, err := i.sendingJobInfoTargetRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	features, err := i.sendingJobInfoFeatureRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prefectures, err := i.sendingJobInfoPrefectureRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workCharmPoints, err := i.sendingJobInfoWorkCharmPointRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	employmentStatuses, err := i.sendingJobInfoEmploymentStatusRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredConditions, err := i.sendingJobInfoRequiredConditionRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLicenses, err := i.sendingJobInfoRequiredLicenseRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredPCTools, err := i.sendingJobInfoRequiredPCToolRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguages, err := i.sendingJobInfoRequiredLanguageRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredLanguageTypes, err := i.sendingJobInfoRequiredLanguageTypeRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopments, err := i.sendingJobInfoRequiredExperienceDevelopmentRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceDevelopmentTypes, err := i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceJobs, err := i.sendingJobInfoRequiredExperienceJobRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceIndustries, err := i.sendingJobInfoRequiredExperienceIndustryRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredExperienceOccupations, err := i.sendingJobInfoRequiredExperienceOccupationRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	requiredSocialExperiences, err := i.sendingJobInfoRequiredSocialExperienceRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	occupations, err := i.sendingJobInfoOccupationRepository.GetListByIDList(jobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sendingJobInformation := range sendingJobInformationList {
		for _, t := range targets {
			if sendingJobInformation.ID == t.SendingJobInformationID {
				value := entity.SendingJobInformationTarget{
					SendingJobInformationID: sendingJobInformation.ID,
					Target:                  t.Target,
				}
				sendingJobInformation.Targets = append(sendingJobInformation.Targets, value)
			}
		}

		for _, f := range features {
			if sendingJobInformation.ID == f.SendingJobInformationID {
				value := entity.SendingJobInformationFeature{
					SendingJobInformationID: f.SendingJobInformationID,
					Feature:                 f.Feature,
				}
				sendingJobInformation.Features = append(sendingJobInformation.Features, value)
			}
		}

		for _, p := range prefectures {
			if sendingJobInformation.ID == p.SendingJobInformationID {
				value := entity.SendingJobInformationPrefecture{
					SendingJobInformationID: p.SendingJobInformationID,
					Prefecture:              p.Prefecture,
				}
				sendingJobInformation.Prefectures = append(sendingJobInformation.Prefectures, value)
			}
		}

		for _, wcp := range workCharmPoints {
			if sendingJobInformation.ID == wcp.SendingJobInformationID {
				value := entity.SendingJobInformationWorkCharmPoint{
					SendingJobInformationID: wcp.SendingJobInformationID,
					Title:                   wcp.Title,
					Contents:                wcp.Contents,
				}
				sendingJobInformation.WorkCharmPoints = append(sendingJobInformation.WorkCharmPoints, value)
			}
		}

		for _, es := range employmentStatuses {
			if sendingJobInformation.ID == es.SendingJobInformationID {
				value := entity.SendingJobInformationEmploymentStatus{
					SendingJobInformationID: es.SendingJobInformationID,
					EmploymentStatus:        es.EmploymentStatus,
				}
				sendingJobInformation.EmploymentStatuses = append(sendingJobInformation.EmploymentStatuses, value)
			}
		}

		// 必要条件
		for _, condition := range requiredConditions {
			if sendingJobInformation.ID == condition.SendingJobInformationID {
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

				sendingJobInformation.RequiredConditions = append(sendingJobInformation.RequiredConditions, *condition)
			}
		}

		for _, rse := range requiredSocialExperiences {
			if sendingJobInformation.ID == rse.SendingJobInformationID {
				value := entity.SendingJobInformationRequiredSocialExperience{
					SendingJobInformationID: rse.SendingJobInformationID,
					SocialExperienceType:    rse.SocialExperienceType,
				}
				sendingJobInformation.RequiredSocialExperiences = append(sendingJobInformation.RequiredSocialExperiences, value)
			}
		}

		for _, oc := range occupations {
			if sendingJobInformation.ID == oc.SendingJobInformationID {
				value := entity.SendingJobInformationOccupation{
					SendingJobInformationID: oc.SendingJobInformationID,
					Occupation:              oc.Occupation,
				}
				sendingJobInformation.Occupations = append(sendingJobInformation.Occupations, value)
			}
		}
	}

	/************ 絞り込み検索 **************/

	// 絞り込み検索処理
	searchSendingJobInformationList, err := searchSendingJobInformationList(sendingJobInformationList, input.SearchParam)
	if err != nil {
		return output, err
	}

	/************ 絞り込みの結果を使って求人数をセット **************/

	// 初回レンダリング時はinput.SendingEnterpriseID==0になるためその場合の処理を記述
	if input.SendingEnterpriseID == 0 {
		sendingEnterpriseID = sendingEnterpriseList[0].ID
	} else {
		sendingEnterpriseID = input.SendingEnterpriseID
	}

	for _, sendingEnterprise := range sendingEnterpriseList {
		for _, sendingJobInformation := range searchSendingJobInformationList {
			if sendingEnterprise.ID == sendingJobInformation.SendingEnterpriseID {
				sendingEnterprise.JobInformationCount = sendingEnterprise.JobInformationCount + 1

				// 選択中の送客先の求人のみを「selectSendingJobInformationList」にセット
				if sendingJobInformation.SendingEnterpriseID == sendingEnterpriseID {
					selectSendingJobInformationList = append(selectSendingJobInformationList, sendingJobInformation)
				}
			}
		}
	}

	output.SendingEnterpriseList = sendingEnterpriseList

	/************ ページ関連の処理 **************/

	// ページの最大数を取得
	output.MaxPageNumber = getSendingJobInformationListMaxPage(selectSendingJobInformationList)

	// 指定ページの求人20件を取得（本番実装までは1ページあたり5件）
	sendingJobInformationList20 := getSendingJobInformationListWithPage(selectSendingJobInformationList, input.PageNumber)

	output.SendingJobInformationList = sendingJobInformationList20

	return output, nil
}

/****************************************************************************************/
/// CSV API
//
// CSVインポート
type ImportSendingJobInformationCSVInput struct {
	CreateParamList []*entity.SendingJobInformation
	MissedRecords   []uint
}

type ImportSendingJobInformationCSVOutput struct {
	MissedRecords []uint
	OK            bool
}

func (i *SendingJobInformationInteractorImpl) ImportSendingJobInformationCSV(input ImportSendingJobInformationCSVInput) (ImportSendingJobInformationCSVOutput, error) {
	var (
		output ImportSendingJobInformationCSVOutput
	)

	// CSVインポート
	for _, sendingJobInformation := range input.CreateParamList {
		err := i.sendingJobInformationRepository.Create(sendingJobInformation)
		if err != nil {
			output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
			fmt.Println(err)
			return output, err
		}

		for _, target := range sendingJobInformation.Targets {
			t := entity.NewSendingJobInformationTarget(
				sendingJobInformation.ID,
				target.Target,
			)

			err = i.sendingJobInfoTargetRepository.Create(t)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
		}

		for _, feature := range sendingJobInformation.Features {
			f := entity.NewSendingJobInformationFeature(
				sendingJobInformation.ID,
				feature.Feature,
			)

			err = i.sendingJobInfoFeatureRepository.Create(f)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
		}

		for _, prefecture := range sendingJobInformation.Prefectures {
			p := entity.NewSendingJobInformationPrefecture(
				sendingJobInformation.ID,
				prefecture.Prefecture,
			)

			err = i.sendingJobInfoPrefectureRepository.Create(p)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
		}

		for _, workCharmPoint := range sendingJobInformation.WorkCharmPoints {
			wcp := entity.NewSendingJobInformationWorkCharmPoint(
				sendingJobInformation.ID,
				workCharmPoint.Title,
				workCharmPoint.Contents,
			)

			err = i.sendingJobInfoWorkCharmPointRepository.Create(wcp)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
		}

		for _, employmentStatus := range sendingJobInformation.EmploymentStatuses {
			es := entity.NewSendingJobInformationEmploymentStatus(
				sendingJobInformation.ID,
				employmentStatus.EmploymentStatus,
			)

			err = i.sendingJobInfoEmploymentStatusRepository.Create(es)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
		}

		// 必要条件
		for _, requiredCondition := range sendingJobInformation.RequiredConditions {
			rc := entity.NewSendingJobInformationRequiredCondition(
				sendingJobInformation.ID,
				requiredCondition.IsCommon,
				requiredCondition.RequiredManagement,
			)

			err = i.sendingJobInfoRequiredConditionRepository.Create(rc)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
			requiredCondition.ID = rc.ID

			// 必要資格　複数
			for _, requiredLicense := range requiredCondition.RequiredLicenses {
				rl := entity.NewSendingJobInformationRequiredLicense(
					requiredCondition.ID,
					requiredLicense.License,
				)

				err = i.sendingJobInfoRequiredLicenseRepository.Create(rl)
				if err != nil {
					output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
					fmt.Println(err)
					return output, err
				}
			}

			// 必要PCツール　単数
			for _, requiredPCTool := range requiredCondition.RequiredPCTools {
				rpt := entity.NewSendingJobInformationRequiredPCTool(
					requiredCondition.ID,
					requiredPCTool.Tool,
				)

				err = i.sendingJobInfoRequiredPCToolRepository.Create(rpt)
				if err != nil {
					output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
					fmt.Println(err)
					return output, err
				}
			}

			// 必要言語スキル　単数
			requiredLanguages := entity.NewSendingJobInformationRequiredLanguage(
				requiredCondition.ID,
				requiredCondition.RequiredLanguages.LanguageLevel,
				requiredCondition.RequiredLanguages.Toeic,
				requiredCondition.RequiredLanguages.ToeflIBT,
				requiredCondition.RequiredLanguages.ToeflPBT,
			)

			err = i.sendingJobInfoRequiredLanguageRepository.Create(requiredLanguages)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}

			for _, languageType := range requiredCondition.RequiredLanguages.LanguageTypes {
				lt := entity.NewSendingJobInformationRequiredLanguageType(
					requiredLanguages.ID,
					languageType.LanguageType,
				)

				err = i.sendingJobInfoRequiredLanguageTypeRepository.Create(lt)
				if err != nil {
					output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
					fmt.Println(err)
					return output, err
				}
			}

			// 開発スキル　　言語,OS 各1つずつ
			for _, requiredExperienceDevelopment := range requiredCondition.RequiredExperienceDevelopments {
				red := entity.NewSendingJobInformationRequiredExperienceDevelopment(
					requiredCondition.ID,
					requiredExperienceDevelopment.DevelopmentCategory,
					requiredExperienceDevelopment.ExperienceYear,
					requiredExperienceDevelopment.ExperienceMonth,
				)

				err = i.sendingJobInfoRequiredExperienceDevelopmentRepository.Create(red)
				if err != nil {
					output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
					fmt.Println(err)
					return output, err
				}

				for _, experienceDevelopmentType := range requiredExperienceDevelopment.ExperienceDevelopmentTypes {
					dt := entity.NewSendingJobInformationRequiredExperienceDevelopmentType(
						red.ID,
						experienceDevelopmentType.DevelopmentType,
					)

					err = i.sendingJobInfoRequiredExperienceDevelopmentTypeRepository.Create(dt)
					if err != nil {
						output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
						fmt.Println(err)
						return output, err
					}
				}
			}

			// 必要経験　　単数
			requiredExperienceJobs := entity.NewSendingJobInformationRequiredExperienceJob(
				requiredCondition.ID,
				requiredCondition.RequiredExperienceJobs.ExperienceYear,
				requiredCondition.RequiredExperienceJobs.ExperienceMonth,
			)

			err = i.sendingJobInfoRequiredExperienceJobRepository.Create(requiredExperienceJobs)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}

			for _, experienceIndustry := range requiredCondition.RequiredExperienceJobs.ExperienceIndustries {
				ei := entity.NewSendingJobInformationRequiredExperienceIndustry(
					requiredExperienceJobs.ID,
					experienceIndustry.ExperienceIndustry,
				)

				err = i.sendingJobInfoRequiredExperienceIndustryRepository.Create(ei)
				if err != nil {
					output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
					fmt.Println(err)
					return output, err
				}
			}

			for _, experienceOccupation := range requiredCondition.RequiredExperienceJobs.ExperienceOccupations {
				eo := entity.NewSendingJobInformationRequiredExperienceOccupation(
					requiredExperienceJobs.ID,
					experienceOccupation.ExperienceOccupation,
				)

				err = i.sendingJobInfoRequiredExperienceOccupationRepository.Create(eo)
				if err != nil {
					output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
					fmt.Println(err)
					return output, err
				}
			}
		}

		for _, requiredSocialExperience := range sendingJobInformation.RequiredSocialExperiences {
			rse := entity.NewSendingJobInformationRequiredSocialExperience(
				sendingJobInformation.ID,
				requiredSocialExperience.SocialExperienceType,
			)

			err = i.sendingJobInfoRequiredSocialExperienceRepository.Create(rse)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
		}

		// 募集職種
		for _, occupation := range sendingJobInformation.Occupations {
			oc := entity.NewSendingJobInformationOccupation(
				sendingJobInformation.ID,
				occupation.Occupation,
			)

			err = i.sendingJobInfoOccupationRepository.Create(oc)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
				return output, err
			}
		}

		// 業界
		for _, industry := range sendingJobInformation.Industries {
			ind := entity.NewSendingJobInformationIndustry(
				sendingJobInformation.ID,
				industry.Industry,
			)

			err = i.sendingJobInfoIndustryRepository.Create(ind)
			if err != nil {
				output.MissedRecords = append(output.MissedRecords, sendingJobInformation.RecordLine)
				fmt.Println(err)
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
