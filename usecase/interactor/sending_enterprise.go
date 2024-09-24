package interactor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type SendingEnterpriseInteractor interface {
	// 汎用系 API
	CreateSendingEnterprise(input CreateSendingEnterpriseInput) (CreateSendingEnterpriseOutput, error)
	UpdateSendingEnterprise(input UpdateSendingEnterpriseInput) (UpdateSendingEnterpriseOutput, error)
	UpdateSendingEnterprisePassword(input UpdateSendingEnterprisePasswordInput) (UpdateSendingEnterprisePasswordOutput, error)
	DeleteSendingEnterprise(input DeleteSendingEnterpriseInput) (DeleteSendingEnterpriseOutput, error)
	GetSendingEnterpriseByID(input GetSendingEnterpriseByIDInput) (GetSendingEnterpriseByIDOutput, error)                                                    // 指定IDの送客先エージェント情報を削除する関数
	GetSendingEnterpriseAndBillingAddressByID(input GetSendingEnterpriseAndBillingAddressByIDInput) (GetSendingEnterpriseAndBillingAddressByIDOutput, error) // 指定IDの送客先エージェント情報を削除する関数
	GetSendingEnterpriseByUUID(input GetSendingEnterpriseByUUIDInput) (GetSendingEnterpriseByUUIDOutput, error)                                              // 指定UUIDの送客先エージェント情報を削除する関数

	//ページネーション
	GetAllSendingEnterpriseByPageAndFreeWord(input GetAllSendingEnterpriseByPageAndFreeWordInput) (GetAllSendingEnterpriseByPageAndFreeWordOutput, error) // 送客先エージェント一覧（すべて）

	// 固有ページ用
	GetSendingInformationForSendingMail(input GetSendingInformationForSendingMailInput) (GetSendingInformationForSendingMailOutput, error)
	SendSendingMail(input SendSendingMailInput) (SendSendingMailOutput, error)
	SendMailForRSVP(input SendMailForRSVPInput) (SendMailForRSVPOutput, error)
	GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(input GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDInput) (GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDOutput, error)

	SigninSendingEnterprise(input SigninSendingEnterpriseInput) (SigninSendingEnterpriseOutput, error)

	// 送客先エージェントの絞り込み検索
	// GetSearchSendingEnterpriseListByAgentID(input GetSearchSendingEnterpriseListByAgentIDInput) (GetSearchSendingEnterpriseListByAgentIDOutput, error)

	//送客先エージェント資料関連
	CreateSendingEnterpriseReferenceMaterial(input CreateSendingEnterpriseReferenceMaterialInput) (CreateSendingEnterpriseReferenceMaterialOutput, error)
	UpdateSendingEnterpriseReferenceMaterial(input UpdateSendingEnterpriseReferenceMaterialInput) (UpdateSendingEnterpriseReferenceMaterialOutput, error)
	DeleteSendingEnterpriseReferenceMaterial(input DeleteSendingEnterpriseReferenceMaterialInput) (DeleteSendingEnterpriseReferenceMaterialOutput, error)
	GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(input GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDInput) (GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDOutput, error)
}

type SendingEnterpriseInteractorImpl struct {
	firebase                                        usecase.Firebase
	sendgrid                                        config.Sendgrid
	sendingEnterpriseRepository                     usecase.SendingEnterpriseRepository
	sendingBillingAddressRepository                 usecase.SendingBillingAddressRepository
	sendingBillingAddressStaffRepository            usecase.SendingBillingAddressStaffRepository
	sendingJobInformationRepository                 usecase.SendingJobInformationRepository
	sendingEnterpriseReferenceMaterialRepository    usecase.SendingEnterpriseReferenceMaterialRepository
	sendingPhaseRepository                          usecase.SendingPhaseRepository
	sendingJobSeekerRepository                      usecase.SendingJobSeekerRepository
	sendingJobSeekerStudentHistoryRepository        usecase.SendingJobSeekerStudentHistoryRepository
	sendingJobSeekerWorkHistoryRepository           usecase.SendingJobSeekerWorkHistoryRepository
	sendingJobSeekerExperienceIndustryRepository    usecase.SendingJobSeekerExperienceIndustryRepository
	sendingJobSeekerDepartmentHistoryRepository     usecase.SendingJobSeekerDepartmentHistoryRepository
	sendingJobSeekerLicenseRepository               usecase.SendingJobSeekerLicenseRepository
	sendingJobSeekerSelfPromotionRepository         usecase.SendingJobSeekerSelfPromotionRepository
	sendingJobSeekerDocumentRepository              usecase.SendingJobSeekerDocumentRepository
	sendingJobSeekerDesiredIndustryRepository       usecase.SendingJobSeekerDesiredIndustryRepository
	sendingJobSeekerDesiredOccupationRepository     usecase.SendingJobSeekerDesiredOccupationRepository
	sendingJobSeekerDesiredWorkLocationRepository   usecase.SendingJobSeekerDesiredWorkLocationRepository
	sendingJobSeekerDesiredHolidayTypeRepository    usecase.SendingJobSeekerDesiredHolidayTypeRepository
	sendingJobSeekerDevelopmentSkillRepository      usecase.SendingJobSeekerDevelopmentSkillRepository
	sendingJobSeekerLanguageSkillRepository         usecase.SendingJobSeekerLanguageSkillRepository
	sendingJobSeekerPCToolRepository                usecase.SendingJobSeekerPCToolRepository
	sendingJobSeekerExperienceOccupationRepository  usecase.SendingJobSeekerExperienceOccupationRepository
	sendingJobSeekerDesiredCompanyScaleRepository   usecase.SendingJobSeekerDesiredCompanyScaleRepository
	sendingCustomerRepository                       usecase.SendingCustomerRepository
	sendingJobSeekerDesiredJobInformationRepository usecase.SendingJobSeekerDesiredJobInformationRepository
	sendingShareDocumentRepository                  usecase.SendingShareDocumentRepository
	sendingEnterpriseSpecialityRepository           usecase.SendingEnterpriseSpecialityRepository
}

// SendingEnterpriseInteractorImpl is an implementation of SendingEnterpriseInteractor
func NewSendingEnterpriseInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	seR usecase.SendingEnterpriseRepository,
	sbR usecase.SendingBillingAddressRepository,
	sjiR usecase.SendingJobInformationRepository,
	sbsR usecase.SendingBillingAddressStaffRepository,
	sarmR usecase.SendingEnterpriseReferenceMaterialRepository,
	spR usecase.SendingPhaseRepository,
	sjsR usecase.SendingJobSeekerRepository,
	sjsshR usecase.SendingJobSeekerStudentHistoryRepository,
	sjswhR usecase.SendingJobSeekerWorkHistoryRepository,
	sjseiR usecase.SendingJobSeekerExperienceIndustryRepository,
	sjseoR usecase.SendingJobSeekerExperienceOccupationRepository,
	sjslR usecase.SendingJobSeekerLicenseRepository,
	sjsspR usecase.SendingJobSeekerSelfPromotionRepository,
	sjsdR usecase.SendingJobSeekerDocumentRepository,
	sjsdiR usecase.SendingJobSeekerDesiredIndustryRepository,
	sjsdoR usecase.SendingJobSeekerDesiredOccupationRepository,
	sjsdwlR usecase.SendingJobSeekerDesiredWorkLocationRepository,
	sjsdhtR usecase.SendingJobSeekerDesiredHolidayTypeRepository,
	sjsdsR usecase.SendingJobSeekerDevelopmentSkillRepository,
	sjslsR usecase.SendingJobSeekerLanguageSkillRepository,
	sjsptR usecase.SendingJobSeekerPCToolRepository,
	sjsdhR usecase.SendingJobSeekerDepartmentHistoryRepository,
	sjsdcsR usecase.SendingJobSeekerDesiredCompanyScaleRepository,
	scR usecase.SendingCustomerRepository,
	sjsdjiR usecase.SendingJobSeekerDesiredJobInformationRepository,
	ssdR usecase.SendingShareDocumentRepository,
	sesR usecase.SendingEnterpriseSpecialityRepository,
) SendingEnterpriseInteractor {
	return &SendingEnterpriseInteractorImpl{
		firebase:                                        fb,
		sendgrid:                                        sg,
		sendingEnterpriseRepository:                     seR,
		sendingBillingAddressRepository:                 sbR,
		sendingBillingAddressStaffRepository:            sbsR,
		sendingJobInformationRepository:                 sjiR,
		sendingEnterpriseReferenceMaterialRepository:    sarmR,
		sendingPhaseRepository:                          spR,
		sendingJobSeekerRepository:                      sjsR,
		sendingJobSeekerStudentHistoryRepository:        sjsshR,
		sendingJobSeekerWorkHistoryRepository:           sjswhR,
		sendingJobSeekerExperienceIndustryRepository:    sjseiR,
		sendingJobSeekerExperienceOccupationRepository:  sjseoR,
		sendingJobSeekerLicenseRepository:               sjslR,
		sendingJobSeekerSelfPromotionRepository:         sjsspR,
		sendingJobSeekerDocumentRepository:              sjsdR,
		sendingJobSeekerDesiredIndustryRepository:       sjsdiR,
		sendingJobSeekerDesiredOccupationRepository:     sjsdoR,
		sendingJobSeekerDesiredWorkLocationRepository:   sjsdwlR,
		sendingJobSeekerDesiredHolidayTypeRepository:    sjsdhtR,
		sendingJobSeekerDevelopmentSkillRepository:      sjsdsR,
		sendingJobSeekerLanguageSkillRepository:         sjslsR,
		sendingJobSeekerPCToolRepository:                sjsptR,
		sendingJobSeekerDepartmentHistoryRepository:     sjsdhR,
		sendingJobSeekerDesiredCompanyScaleRepository:   sjsdcsR,
		sendingCustomerRepository:                       scR,
		sendingJobSeekerDesiredJobInformationRepository: sjsdjiR,
		sendingShareDocumentRepository:                  ssdR,
		sendingEnterpriseSpecialityRepository:           sesR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//求人送客先エージェントの作成
type CreateSendingEnterpriseInput struct {
	CreateParam entity.CreateOrUpdateSendingEnterpriseParam
}

type CreateSendingEnterpriseOutput struct {
	SendingEnterprise *entity.SendingEnterprise
}

func (i *SendingEnterpriseInteractorImpl) CreateSendingEnterprise(input CreateSendingEnterpriseInput) (CreateSendingEnterpriseOutput, error) {
	var (
		output CreateSendingEnterpriseOutput
		err    error
	)

	sendingEnterprise := entity.NewSendingEnterprise(
		input.CreateParam.CompanyName,
		input.CreateParam.AgentStaffID,
		input.CreateParam.CorporateSiteURL,
		input.CreateParam.Representative,
		input.CreateParam.Establishment,
		input.CreateParam.PostCode,
		input.CreateParam.OfficeLocation,
		input.CreateParam.EmployeeNumberSingle,
		input.CreateParam.PublicOffering,
		input.CreateParam.SendingTarget,
		input.CreateParam.SendingRequiredInfo,
		input.CreateParam.Remarks,
		input.CreateParam.Password,
	)

	err = i.sendingEnterpriseRepository.Create(sendingEnterprise)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// uint, uint, null.Int, string, string, string, string, string, string, null.Int)go
	// 空の請求先も作成
	sendingBillingAddress := entity.NewSendingBillingAddress(
		sendingEnterprise.ID,
		input.CreateParam.AgentStaffID, // input.CreateParam.AgentStaffID,
		null.NewInt(0, false),          // input.CreateParam.ContractPhase,
		"",                             // input.CreateParam.ContractDate,
		"",                             // input.CreateParam.PaymentPolicy,
		input.CreateParam.CompanyName,  // input.CreateParam.AgentName,
		"",                             // input.CreateParam.Address,
		input.CreateParam.CompanyName,  // input.CreateParam.Title,
		"",                             // input.CreateParam.ScheduleAdjustmentURL,
		null.NewInt(0, false),          // input.CreateParam.Commission,
	)

	err = i.sendingBillingAddressRepository.Create(sendingBillingAddress)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントの特徴
	sendingEnterpriseSpeciality := entity.NewSendingEnterpriseSpeciality(
		sendingEnterprise.ID,
		"", // input.CreateParam.Speciality.ImageURL,
		input.CreateParam.Speciality.JobInformationCount,
		input.CreateParam.Speciality.SpecializedOccupation,
		input.CreateParam.Speciality.SpecializedIndustry,
		input.CreateParam.Speciality.SpecializedArea,
		input.CreateParam.Speciality.SpecializedCompanyType,
		input.CreateParam.Speciality.SpecializedJobSeekerType,
		input.CreateParam.Speciality.ConsultingStrengths,
		input.CreateParam.Speciality.SupportContent,
		input.CreateParam.Speciality.PRPoint,
	)

	err = i.sendingEnterpriseSpecialityRepository.Create(sendingEnterpriseSpeciality)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterprise.Speciality = *sendingEnterpriseSpeciality

	output.SendingEnterprise = sendingEnterprise

	return output, nil
}

// 求人送客先エージェントの更新
type UpdateSendingEnterpriseInput struct {
	UpdateParam         entity.CreateOrUpdateSendingEnterpriseParam
	SendingEnterpriseID uint
}

type UpdateSendingEnterpriseOutput struct {
	SendingEnterprise *entity.SendingEnterprise
}

func (i *SendingEnterpriseInteractorImpl) UpdateSendingEnterprise(input UpdateSendingEnterpriseInput) (UpdateSendingEnterpriseOutput, error) {
	var (
		output UpdateSendingEnterpriseOutput
		err    error
	)

	sendingEnterprise := entity.NewSendingEnterprise(
		input.UpdateParam.CompanyName,
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.CorporateSiteURL,
		input.UpdateParam.Representative,
		input.UpdateParam.Establishment,
		input.UpdateParam.PostCode,
		input.UpdateParam.OfficeLocation,
		input.UpdateParam.EmployeeNumberSingle,
		input.UpdateParam.PublicOffering,
		input.UpdateParam.SendingTarget,
		input.UpdateParam.SendingRequiredInfo,
		input.UpdateParam.Remarks,
		"", // passwordは単体でのみ更新する
	)

	err = i.sendingEnterpriseRepository.Update(sendingEnterprise, input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントの特徴
	sendingEnterpriseSpeciality := entity.NewSendingEnterpriseSpeciality(
		input.SendingEnterpriseID,
		input.UpdateParam.Speciality.ImageURL,
		input.UpdateParam.Speciality.JobInformationCount,
		input.UpdateParam.Speciality.SpecializedOccupation,
		input.UpdateParam.Speciality.SpecializedIndustry,
		input.UpdateParam.Speciality.SpecializedArea,
		input.UpdateParam.Speciality.SpecializedCompanyType,
		input.UpdateParam.Speciality.SpecializedJobSeekerType,
		input.UpdateParam.Speciality.ConsultingStrengths,
		input.UpdateParam.Speciality.SupportContent,
		input.UpdateParam.Speciality.PRPoint,
	)

	err = i.sendingEnterpriseSpecialityRepository.Update(sendingEnterpriseSpeciality)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterprise.Speciality = *sendingEnterpriseSpeciality

	output.SendingEnterprise = sendingEnterprise

	return output, nil
}

// 求人送客先エージェントの更新
type UpdateSendingEnterprisePasswordInput struct {
	SendingEnterpriseID uint
	Password            string
}

type UpdateSendingEnterprisePasswordOutput struct {
	OK bool
}

func (i *SendingEnterpriseInteractorImpl) UpdateSendingEnterprisePassword(input UpdateSendingEnterprisePasswordInput) (UpdateSendingEnterprisePasswordOutput, error) {
	var (
		output UpdateSendingEnterprisePasswordOutput
		err    error
	)

	/************ 1. テーブルの更新 **************/

	err = i.sendingEnterpriseRepository.UpdatePassword(input.SendingEnterpriseID, input.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 求人送客先エージェントの更新
type SigninSendingEnterpriseInput struct {
	SigninParam entity.SigninSendingEnterprisePasswordParam
}

type SigninSendingEnterpriseOutput struct {
	SendingEnterprise    *entity.SendingEnterprise
	SendingJobSeeker     *entity.SendingJobSeeker
	SendingShareDocument *entity.SendingShareDocument
}

func (i *SendingEnterpriseInteractorImpl) SigninSendingEnterprise(input SigninSendingEnterpriseInput) (SigninSendingEnterpriseOutput, error) {
	var (
		output SigninSendingEnterpriseOutput
		err    error
	)

	/************ 1. uuidの形式をチェック **************/

	enterpriseUUID, err := uuid.Parse(input.SigninParam.SendingEnterpriseUUID)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", "送客先のuuidのフォーマットが不正です", entity.ErrRequestError)
		return output, wrapped
	}

	jobSeekerUUID, err := uuid.Parse(input.SigninParam.SendingJobSeekerUUID)
	if err != nil {
		wrapped := fmt.Errorf("%s:%w", "送客求職者のuuidのフォーマットが不正です", entity.ErrRequestError)
		return output, wrapped
	}

	/************ 2. 送客先エージェントの情報を取得 **************/

	sendingEnterprise, err := i.sendingEnterpriseRepository.FindByUUID(enterpriseUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 3. パスワードが正しいかを確認 **************/

	if sendingEnterprise.Password != input.SigninParam.Password {
		err = fmt.Errorf("パスワードが正しくありません。")
		return output, err
	}

	/************ 4. 送客求職者情報を取得 **************/

	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByUUID(jobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	studentHistory, err := i.sendingJobSeekerStudentHistoryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workHistory, err := i.sendingJobSeekerWorkHistoryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustry, err := i.sendingJobSeekerExperienceIndustryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory, err := i.sendingJobSeekerDepartmentHistoryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupation, err := i.sendingJobSeekerExperienceOccupationRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredCompanyScale, err := i.sendingJobSeekerDesiredCompanyScaleRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	license, err := i.sendingJobSeekerLicenseRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selfPromotion, err := i.sendingJobSeekerSelfPromotionRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	document, err := i.sendingJobSeekerDocumentRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredIndustry, err := i.sendingJobSeekerDesiredIndustryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupation, err := i.sendingJobSeekerDesiredOccupationRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocation, err := i.sendingJobSeekerDesiredWorkLocationRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredHolidayType, err := i.sendingJobSeekerDesiredHolidayTypeRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	developmentSkill, err := i.sendingJobSeekerDevelopmentSkillRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	languageSkill, err := i.sendingJobSeekerLanguageSkillRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	pcSkill, err := i.sendingJobSeekerPCToolRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sh := range studentHistory {
		value := entity.SendingJobSeekerStudentHistory{
			SendingJobSeekerID: sh.SendingJobSeekerID,
			SchoolCategory:     sh.SchoolCategory,
			SchoolName:         sh.SchoolName,
			SchoolLevel:        sh.SchoolLevel,
			Subject:            sh.Subject,
			EntranceYear:       sh.EntranceYear,
			FirstStatus:        sh.FirstStatus,
			GraduationYear:     sh.GraduationYear,
			LastStatus:         sh.LastStatus,
		}

		sendingJobSeeker.StudentHistories = append(sendingJobSeeker.StudentHistories, value)
	}

	for _, wh := range workHistory {
		value := entity.SendingJobSeekerWorkHistory{
			ID:                   wh.ID,
			SendingJobSeekerID:   wh.SendingJobSeekerID,
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
				valueEI := entity.SendingJobSeekerExperienceIndustry{
					ID:            ei.ID,
					WorkHistoryID: ei.WorkHistoryID,
					Industry:      ei.Industry,
				}

				value.ExperienceIndustries = append(value.ExperienceIndustries, valueEI)
			}
		}

		for _, dh := range departmentHistory {
			if dh.WorkHistoryID == wh.ID {
				valuedh := entity.SendingJobSeekerDepartmentHistory{
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
						valueEO := entity.SendingJobSeekerExperienceOccupation{
							DepartmentHistoryID: eo.DepartmentHistoryID,
							Occupation:          eo.Occupation,
						}

						valuedh.ExperienceOccupations = append(valuedh.ExperienceOccupations, valueEO)
					}
				}

				value.DepartmentHistories = append(value.DepartmentHistories, valuedh)
			}
		}
		sendingJobSeeker.WorkHistories = append(sendingJobSeeker.WorkHistories, value)
	}

	for _, dcs := range desiredCompanyScale {
		value := entity.SendingJobSeekerDesiredCompanyScale{
			SendingJobSeekerID:  dcs.SendingJobSeekerID,
			DesiredCompanyScale: dcs.DesiredCompanyScale,
		}

		sendingJobSeeker.DesiredCompanyScales = append(sendingJobSeeker.DesiredCompanyScales, value)
	}

	for _, l := range license {
		value := entity.SendingJobSeekerLicense{
			SendingJobSeekerID: l.SendingJobSeekerID,
			LicenseType:        l.LicenseType,
			AcquisitionTime:    l.AcquisitionTime,
		}

		sendingJobSeeker.Licenses = append(sendingJobSeeker.Licenses, value)
	}

	for _, sp := range selfPromotion {
		value := entity.SendingJobSeekerSelfPromotion{
			SendingJobSeekerID: sp.SendingJobSeekerID,
			Title:              sp.Title,
			Contents:           sp.Contents,
		}

		sendingJobSeeker.SelfPromotions = append(sendingJobSeeker.SelfPromotions, value)
	}

	valueDocument := entity.SendingJobSeekerDocument{
		SendingJobSeekerID:      document.SendingJobSeekerID,
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

	sendingJobSeeker.Documents = valueDocument

	for _, di := range desiredIndustry {
		value := entity.SendingJobSeekerDesiredIndustry{
			SendingJobSeekerID: di.SendingJobSeekerID,
			DesiredIndustry:    di.DesiredIndustry,
			DesiredRank:        di.DesiredRank,
		}

		sendingJobSeeker.DesiredIndustries = append(sendingJobSeeker.DesiredIndustries, value)
	}

	for _, do := range desiredOccupation {
		value := entity.SendingJobSeekerDesiredOccupation{
			SendingJobSeekerID: do.SendingJobSeekerID,
			DesiredOccupation:  do.DesiredOccupation,
			DesiredRank:        do.DesiredRank,
		}

		sendingJobSeeker.DesiredOccupations = append(sendingJobSeeker.DesiredOccupations, value)
	}

	for _, dwl := range desiredWorkLocation {
		value := entity.SendingJobSeekerDesiredWorkLocation{
			SendingJobSeekerID:  dwl.SendingJobSeekerID,
			DesiredWorkLocation: dwl.DesiredWorkLocation,
			DesiredRank:         dwl.DesiredRank,
		}

		sendingJobSeeker.DesiredWorkLocations = append(sendingJobSeeker.DesiredWorkLocations, value)
	}

	for _, dht := range desiredHolidayType {
		value := entity.SendingJobSeekerDesiredHolidayType{
			SendingJobSeekerID: dht.SendingJobSeekerID,
			HolidayType:        dht.HolidayType,
		}

		sendingJobSeeker.DesiredHolidayTypes = append(sendingJobSeeker.DesiredHolidayTypes, value)
	}

	for _, ds := range developmentSkill {
		value := entity.SendingJobSeekerDevelopmentSkill{
			SendingJobSeekerID:  ds.SendingJobSeekerID,
			DevelopmentCategory: ds.DevelopmentCategory,
			DevelopmentType:     ds.DevelopmentType,
			ExperienceYear:      ds.ExperienceYear,
			ExperienceMonth:     ds.ExperienceMonth,
		}

		sendingJobSeeker.DevelopmentSkills = append(sendingJobSeeker.DevelopmentSkills, value)
	}

	for _, ls := range languageSkill {
		value := entity.SendingJobSeekerLanguageSkill{
			SendingJobSeekerID:      ls.SendingJobSeekerID,
			LanguageType:            ls.LanguageType,
			LanguageLevel:           ls.LanguageLevel,
			Toeic:                   ls.Toeic,
			ToeicExaminationYear:    ls.ToeicExaminationYear,
			ToeflIBT:                ls.ToeflIBT,
			ToeflIBTExaminationYear: ls.ToeflIBTExaminationYear,
			ToeflPBT:                ls.ToeflPBT,
			ToeflPBTExaminationYear: ls.ToeflPBTExaminationYear,
		}

		sendingJobSeeker.LanguageSkills = append(sendingJobSeeker.LanguageSkills, value)
	}

	for _, ps := range pcSkill {
		value := entity.SendingJobSeekerPCTool{
			SendingJobSeekerID: ps.SendingJobSeekerID,
			Tool:               ps.Tool,
		}

		sendingJobSeeker.PCTools = append(sendingJobSeeker.PCTools, value)
	}

	/************ 5. シェアが許可された書類情報を取得 **************/

	shareDocuments, err := i.sendingShareDocumentRepository.FindBySendingJobSeekerIDAndSendingEnterpriseID(sendingJobSeeker.ID, sendingEnterprise.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 6. レスポンス **************/

	output.SendingEnterprise = sendingEnterprise
	output.SendingJobSeeker = sendingJobSeeker
	output.SendingShareDocument = shareDocuments

	return output, nil
}

// 求人送客先エージェントの削除
type DeleteSendingEnterpriseInput struct {
	SendingEnterpriseID uint
}

type DeleteSendingEnterpriseOutput struct {
	OK bool
}

func (i *SendingEnterpriseInteractorImpl) DeleteSendingEnterprise(input DeleteSendingEnterpriseInput) (DeleteSendingEnterpriseOutput, error) {
	var (
		output DeleteSendingEnterpriseOutput
	)

	err := i.sendingEnterpriseRepository.Delete(input.SendingEnterpriseID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// 送客先エージェントIDを使って送客先エージェント情報を取得する
type GetSendingEnterpriseByIDInput struct {
	SendingEnterpriseID uint
}

type GetSendingEnterpriseByIDOutput struct {
	SendingEnterprise *entity.SendingEnterprise
}

func (i *SendingEnterpriseInteractorImpl) GetSendingEnterpriseByID(input GetSendingEnterpriseByIDInput) (GetSendingEnterpriseByIDOutput, error) {
	var (
		output GetSendingEnterpriseByIDOutput
		err    error
	)

	sendingEnterprise, err := i.sendingEnterpriseRepository.FindByID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterpriseReferenceMaterial, err := i.sendingEnterpriseReferenceMaterialRepository.FindBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterpriseSpeciality, err := i.sendingEnterpriseSpecialityRepository.FindBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterprise.ReferenceMaterialID = sendingEnterpriseReferenceMaterial.ID

	sendingEnterprise.ReferenceMaterial = *sendingEnterpriseReferenceMaterial
	sendingEnterprise.Speciality = *sendingEnterpriseSpeciality

	output.SendingEnterprise = sendingEnterprise

	return output, nil
}

// 送客先エージェントIDを使って送客先エージェント情報を取得する
type GetSendingEnterpriseAndBillingAddressByIDInput struct {
	SendingEnterpriseID uint
}

type GetSendingEnterpriseAndBillingAddressByIDOutput struct {
	SendingEnterprise *entity.SendingEnterprise
}

func (i *SendingEnterpriseInteractorImpl) GetSendingEnterpriseAndBillingAddressByID(input GetSendingEnterpriseAndBillingAddressByIDInput) (GetSendingEnterpriseAndBillingAddressByIDOutput, error) {
	var (
		output GetSendingEnterpriseAndBillingAddressByIDOutput
		err    error
	)

	/************ 1. 送客先エージェントの情報を取得 **************/

	sendingEnterprise, err := i.sendingEnterpriseRepository.FindByID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. 送客先エージェントの子テーブル情報を取得 **************/

	sendingEnterpriseReferenceMaterial, err := i.sendingEnterpriseReferenceMaterialRepository.FindBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterpriseSpeciality, err := i.sendingEnterpriseSpecialityRepository.FindBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingBillingAddress, err := i.sendingBillingAddressRepository.FindBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingBillingStaffList, err := i.sendingBillingAddressStaffRepository.FindBySendingBillingAddressID(sendingBillingAddress.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 3. レスポンスの情報を整理 **************/

	for _, staff := range sendingBillingStaffList {
		sendingBillingAddress.Staffs = append(sendingBillingAddress.Staffs, *staff)
	}

	sendingEnterprise.SendingBillingAddress = *sendingBillingAddress
	sendingEnterprise.ReferenceMaterialID = sendingEnterpriseReferenceMaterial.ID
	sendingEnterprise.ReferenceMaterial = *sendingEnterpriseReferenceMaterial
	sendingEnterprise.Speciality = *sendingEnterpriseSpeciality
	output.SendingEnterprise = sendingEnterprise

	return output, nil
}

// 送客先エージェントIDを使って送客先エージェント情報を取得する
type GetSendingEnterpriseByUUIDInput struct {
	SendingEnterpriseUUID uuid.UUID
}

type GetSendingEnterpriseByUUIDOutput struct {
	SendingEnterprise *entity.SendingEnterprise
}

func (i *SendingEnterpriseInteractorImpl) GetSendingEnterpriseByUUID(input GetSendingEnterpriseByUUIDInput) (GetSendingEnterpriseByUUIDOutput, error) {
	var (
		output GetSendingEnterpriseByUUIDOutput
		err    error
	)

	sendingEnterprise, err := i.sendingEnterpriseRepository.FindByUUID(input.SendingEnterpriseUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterpriseReferenceMaterial, err := i.sendingEnterpriseReferenceMaterialRepository.FindBySendingEnterpriseID(sendingEnterprise.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterpriseSpeciality, err := i.sendingEnterpriseSpecialityRepository.FindBySendingEnterpriseID(sendingEnterprise.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterprise.ReferenceMaterialID = sendingEnterpriseReferenceMaterial.ID

	sendingEnterprise.ReferenceMaterial = *sendingEnterpriseReferenceMaterial
	sendingEnterprise.Speciality = *sendingEnterpriseSpeciality

	output.SendingEnterprise = sendingEnterprise

	return output, nil
}

/****************************************************************************************/
/// ページネーション
//
// ページごとにすべての送客先エージェント情報を取得
type GetAllSendingEnterpriseByPageAndFreeWordInput struct {
	PageNumber uint
	FreeWord   string
}

type GetAllSendingEnterpriseByPageAndFreeWordOutput struct {
	SendingEnterpriseList []*entity.SendingEnterprise
	MaxPageNumber         uint
	IDList                []uint
}

func (i *SendingEnterpriseInteractorImpl) GetAllSendingEnterpriseByPageAndFreeWord(input GetAllSendingEnterpriseByPageAndFreeWordInput) (GetAllSendingEnterpriseByPageAndFreeWordOutput, error) {
	var (
		output GetAllSendingEnterpriseByPageAndFreeWordOutput
		err    error
	)

	sendingEnterpriseList, err := i.sendingEnterpriseRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingEnterpriseReferenceMaterialList, err := i.sendingEnterpriseReferenceMaterialRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// freewordの絞り込み 企業名or担当者名
	if input.FreeWord != "" {
		filteredList := []*entity.SendingEnterprise{}
		for i := range sendingEnterpriseList {
			if strings.Contains(sendingEnterpriseList[i].CompanyName, input.FreeWord) {
				filteredList = append(filteredList, sendingEnterpriseList[i])
			} else if strings.Contains(sendingEnterpriseList[i].StaffName, input.FreeWord) {
				filteredList = append(filteredList, sendingEnterpriseList[i])
			}
		}

		sendingEnterpriseList = filteredList
	}

	// ページの最大数を取得
	output.MaxPageNumber = getSendingEnterpriseListMaxPage(sendingEnterpriseList)

	// IDListを返す
	for _, enterprise := range sendingEnterpriseList {
		output.IDList = append(output.IDList, enterprise.ID)
	}

	// 指定ページの送客先エージェント50件を取得（本番実装までは1ページあたり5件）
	sendingEnterpriseList50 := getSendingEnterpriseListWithPage(sendingEnterpriseList, input.PageNumber)

	for i := range sendingEnterpriseList50 {
		for _, file := range sendingEnterpriseReferenceMaterialList {
			if sendingEnterpriseList50[i].ID == file.SendingEnterpriseID {
				sendingEnterpriseList50[i].ReferenceMaterial.Reference1PDFURL = file.Reference1PDFURL
				sendingEnterpriseList50[i].ReferenceMaterial.Reference2PDFURL = file.Reference2PDFURL
			}
		}
	}

	output.SendingEnterpriseList = sendingEnterpriseList50

	return output, nil
}

/****************************************************************************************/
/// 送客先エージェントの絞り込み検索
//
// // エージェントIDから送客先エージェント名一覧を取得する
// type GetSearchSendingEnterpriseListByAgentIDInput struct {
// 	AgentID     uint
// 	PageNumber  uint
// 	SearchParam entity.SearchSendingEnterprise
// }

// type GetSearchSendingEnterpriseListByAgentIDOutput struct {
// 	SendingEnterpriseList []*entity.SendingEnterprise
// 	MaxPageNumber    uint
// 	IDList           []uint
// }

// func (i *SendingEnterpriseInteractorImpl) GetSearchSendingEnterpriseListByAgentID(input GetSearchSendingEnterpriseListByAgentIDInput) (GetSearchSendingEnterpriseListByAgentIDOutput, error) {
// 	var (
// 		output GetSearchSendingEnterpriseListByAgentIDOutput
// 		err    error
// 	)

// 	/**
// 	GetSendingEnterpriseListByAgentIDAndFreeWordは
// 	フリーワードの有無で処理を分岐

// 	フリーワードは社名のみ
// 	*/
// 	sendingEnterpriseList, err := i.sendingEnterpriseRepository.GetListByAgentIDAndFreeWord(input.AgentID, input.SearchParam.FreeWord)
// 	if err != nil {
// 		fmt.Println(err)
// 		return output, err
// 	}

// 	sendingEnterpriseReferenceMaterialList, err := i.sendingEnterpriseReferenceMaterialRepository.GetListByAgentID(input.AgentID)
// 	if err != nil {
// 		fmt.Println(err)
// 		return output, err
// 	}

// 	// 縦持ちテーブルの処理
// 	for _, sendingEnterprise := range sendingEnterpriseList {
// 		for _, industry := range enterpriseIndustryList {
// 			if sendingEnterprise.ID == industry.SendingEnterpriseID {
// 				sendingEnterprise.Industries = append(sendingEnterprise.Industries, industry.Industry)
// 			}
// 		}

// 		for _, file := range sendingEnterpriseReferenceMaterialList {
// 			if sendingEnterprise.ID == file.SendingEnterpriseID {
// 				sendingEnterprise.ReferenceMaterial.Reference1PDFURL = file.Reference1PDFURL
// 				sendingEnterprise.ReferenceMaterial.Reference2PDFURL = file.Reference2PDFURL
// 			}
// 		}
// 	}

// 	// 絞り込み項目の結果を代入するための変数を用意
// 	var (
// 		sendingEnterpriseListWithAgentStaffID []*entity.SendingEnterprise
// 		sendingEnterpriseListWithIndustry     []*entity.SendingEnterprise
// 		// sendingEnterpriseListWithPrefecture   []*entity.SendingEnterprise
// 		sendingEnterpriseListWithCompanyScale []*entity.SendingEnterprise
// 	)

// 	// 営業担当者IDがある場合
// 	agentStaffID, err := strconv.Atoi(input.SearchParam.AgentStaffID)
// 	if !(err != nil || agentStaffID == 0) {
// 		for _, enterprise := range sendingEnterpriseList {
// 			if enterprise.AgentStaffID != uint(agentStaffID) {
// 				continue
// 			}
// 			sendingEnterpriseListWithAgentStaffID = append(sendingEnterpriseListWithAgentStaffID, enterprise)
// 		}
// 	}

// 	// 営業担当者IDが無い場合
// 	if err != nil || agentStaffID == 0 {
// 		sendingEnterpriseListWithAgentStaffID = sendingEnterpriseList
// 	}

// 	fmt.Println("RA担当者: ", sendingEnterpriseListWithAgentStaffID)

// 	// Note: 送客先エージェント規模の絞り込みは送客先エージェントの従業員数（単体）と比較する
// 	// 送客先エージェント規模がある場合
// 	if !(len(input.SearchParam.CompanyScaleTypes) == 0) {
// 	companyScaleLoop:
// 		for _, enterprise := range sendingEnterpriseListWithIndustry {
// 			if enterprise.EmployeeNumberSingle == null.NewInt(0, false) {
// 				continue
// 			}
// 			for _, companyScale := range input.SearchParam.CompanyScaleTypes {
// 				if !companyScale.Valid {
// 					continue
// 				}

// 				if companyScale == null.NewInt(0, true) {
// 					// 10名未満の場合
// 					if enterprise.EmployeeNumberSingle.Int64 < 10 {
// 						sendingEnterpriseListWithCompanyScale = append(sendingEnterpriseListWithCompanyScale, enterprise)
// 						continue companyScaleLoop
// 					}
// 				} else if companyScale == null.NewInt(1, true) {
// 					// 10名以上100名未満の場合
// 					if enterprise.EmployeeNumberSingle.Int64 >= 10 && enterprise.EmployeeNumberSingle.Int64 < 100 {
// 						sendingEnterpriseListWithCompanyScale = append(sendingEnterpriseListWithCompanyScale, enterprise)
// 						continue companyScaleLoop
// 					}
// 				} else if companyScale == null.NewInt(2, true) {
// 					// 100名以上200名未満の場合
// 					if enterprise.EmployeeNumberSingle.Int64 >= 100 && enterprise.EmployeeNumberSingle.Int64 < 200 {
// 						sendingEnterpriseListWithCompanyScale = append(sendingEnterpriseListWithCompanyScale, enterprise)
// 						continue companyScaleLoop
// 					}
// 				} else if companyScale == null.NewInt(3, true) {
// 					// 200名以上1000名未満の場合
// 					if enterprise.EmployeeNumberSingle.Int64 >= 200 && enterprise.EmployeeNumberSingle.Int64 < 1000 {
// 						sendingEnterpriseListWithCompanyScale = append(sendingEnterpriseListWithCompanyScale, enterprise)
// 						continue companyScaleLoop
// 					}
// 				} else if companyScale == null.NewInt(4, true) {
// 					// 1000名以上の場合
// 					if enterprise.EmployeeNumberSingle.Int64 >= 1000 {
// 						sendingEnterpriseListWithCompanyScale = append(sendingEnterpriseListWithCompanyScale, enterprise)
// 						continue companyScaleLoop
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// 送客先エージェント規模が無い場合
// 	if len(input.SearchParam.CompanyScaleTypes) == 0 {
// 		sendingEnterpriseListWithCompanyScale = sendingEnterpriseListWithIndustry
// 	}

// 	fmt.Println("送客先エージェント規模: ", sendingEnterpriseListWithCompanyScale)

// 	// IDListを返す
// 	for _, enterprise := range sendingEnterpriseListWithCompanyScale {
// 		output.IDList = append(output.IDList, enterprise.ID)
// 	}

// 	// ページの最大数を取得
// 	output.MaxPageNumber = getSendingEnterpriseListMaxPage(sendingEnterpriseListWithCompanyScale)

// 	// 指定ページの送客先エージェント20件を取得（本番実装までは1ページあたり5件）
// 	output.SendingEnterpriseList = getSendingEnterpriseListWithPage(sendingEnterpriseListWithCompanyScale, input.PageNumber)

// 	return output, nil

// }

/****************************************************************************************/
// 送客先エージェント資料関連
//
// 求人送客先エージェント資料の作成
type CreateSendingEnterpriseReferenceMaterialInput struct {
	CreateParam entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam
}

type CreateSendingEnterpriseReferenceMaterialOutput struct {
	SendingEnterpriseReferenceMaterial *entity.SendingEnterpriseReferenceMaterial
}

func (i *SendingEnterpriseInteractorImpl) CreateSendingEnterpriseReferenceMaterial(input CreateSendingEnterpriseReferenceMaterialInput) (CreateSendingEnterpriseReferenceMaterialOutput, error) {
	var (
		output CreateSendingEnterpriseReferenceMaterialOutput
		err    error
	)

	sendingEnterpriseReferenceMaterial := entity.NewSendingEnterpriseReferenceMaterial(
		input.CreateParam.SendingEnterpriseID,
		input.CreateParam.Reference1PDFURL,
		input.CreateParam.Reference2PDFURL,
	)

	err = i.sendingEnterpriseReferenceMaterialRepository.Create(sendingEnterpriseReferenceMaterial)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// specialitiesのimage_urlも更新
	err = i.sendingEnterpriseSpecialityRepository.UpdateImageURLBySendingEnterpriseID(
		input.CreateParam.SendingEnterpriseID,
		input.CreateParam.ImageURL,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingEnterpriseReferenceMaterial = sendingEnterpriseReferenceMaterial

	return output, nil

}

// 求人送客先エージェント資料の更新
type UpdateSendingEnterpriseReferenceMaterialInput struct {
	UpdateParam entity.CreateOrUpdateSendingEnterpriseReferenceMaterialParam
}

type UpdateSendingEnterpriseReferenceMaterialOutput struct {
	SendingEnterpriseReferenceMaterial *entity.SendingEnterpriseReferenceMaterial
}

func (i *SendingEnterpriseInteractorImpl) UpdateSendingEnterpriseReferenceMaterial(input UpdateSendingEnterpriseReferenceMaterialInput) (UpdateSendingEnterpriseReferenceMaterialOutput, error) {
	var (
		output UpdateSendingEnterpriseReferenceMaterialOutput
		err    error
	)

	sendingEnterpriseReferenceMaterial := entity.NewSendingEnterpriseReferenceMaterial(
		input.UpdateParam.SendingEnterpriseID,
		input.UpdateParam.Reference1PDFURL,
		input.UpdateParam.Reference2PDFURL,
	)

	err = i.sendingEnterpriseReferenceMaterialRepository.Update(sendingEnterpriseReferenceMaterial)

	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// specialitiesのimage_urlも更新
	err = i.sendingEnterpriseSpecialityRepository.UpdateImageURLBySendingEnterpriseID(
		input.UpdateParam.SendingEnterpriseID,
		input.UpdateParam.ImageURL,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingEnterpriseReferenceMaterial = sendingEnterpriseReferenceMaterial

	return output, nil
}

// 求人送客先エージェント資料の削除
type DeleteSendingEnterpriseReferenceMaterialInput struct {
	SendingEnterpriseID uint
	FileType            string
}

type DeleteSendingEnterpriseReferenceMaterialOutput struct {
	OK bool
}

func (i *SendingEnterpriseInteractorImpl) DeleteSendingEnterpriseReferenceMaterial(input DeleteSendingEnterpriseReferenceMaterialInput) (DeleteSendingEnterpriseReferenceMaterialOutput, error) {
	var (
		output DeleteSendingEnterpriseReferenceMaterialOutput
	)

	if input.FileType == "参考資料1" || input.FileType == "参考資料2" {
		err := i.sendingEnterpriseReferenceMaterialRepository.UpdateReferenceURLBySendingEnterpriseIDAndMaterialType(input.SendingEnterpriseID, input.FileType)
		if err != nil {
			return output, err
		}
	} else if input.FileType == "送客先エージェント画像" {
		err := i.sendingEnterpriseSpecialityRepository.UpdateImageURLBySendingEnterpriseID(input.SendingEnterpriseID, "")
		if err != nil {
			return output, err
		}
	}

	output.OK = true

	return output, nil
}

// 求人送客先エージェント資料の取得
type GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDInput struct {
	SendingEnterpriseID uint
}

type GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDOutput struct {
	SendingEnterpriseReferenceMaterial *entity.SendingEnterpriseReferenceMaterial
}

func (i *SendingEnterpriseInteractorImpl) GetSendingEnterpriseReferenceMaterialBySendingEnterpriseID(input GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDInput) (GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDOutput, error) {
	var (
		output GetSendingEnterpriseReferenceMaterialBySendingEnterpriseIDOutput
		err    error
	)

	sendingEnterpriseReferenceMaterial, err := i.sendingEnterpriseReferenceMaterialRepository.FindBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingEnterpriseReferenceMaterial = sendingEnterpriseReferenceMaterial

	return output, nil

}

// 送客メール送信時に必要な情報をまとめて取得するapi
type GetSendingInformationForSendingMailInput struct {
	SendingJobInformationIDList []uint
}

type GetSendingInformationForSendingMailOutput struct {
	SendingEnterpriseList []*entity.SendingEnterprise
}

func (i *SendingEnterpriseInteractorImpl) GetSendingInformationForSendingMail(input GetSendingInformationForSendingMailInput) (GetSendingInformationForSendingMailOutput, error) {
	var (
		output                   GetSendingInformationForSendingMailOutput
		err                      error
		encounteredForBilling    = map[uint]bool{}
		billingAddressIDList     []uint
		encounteredForEnterprise = map[uint]bool{}
		enterpriseIDList         []uint
	)

	/************ 求人情報の取得 **************/

	if len(input.SendingJobInformationIDList) == 0 {
		return output, nil
	}

	sendingJobInformationList, err := i.sendingJobInformationRepository.GetListByIDList(input.SendingJobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 請求先情報の取得 **************/

	for _, sendingJobInformation := range sendingJobInformationList {
		if !encounteredForBilling[sendingJobInformation.SendingBillingAddressID] {
			encounteredForBilling[sendingJobInformation.SendingBillingAddressID] = true
			billingAddressIDList = append(billingAddressIDList, sendingJobInformation.SendingBillingAddressID)
		}
	}

	if len(billingAddressIDList) == 0 {
		return output, nil
	}

	sendingBillingAddressList, err := i.sendingBillingAddressRepository.GetByIDList(billingAddressIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 請求先の担当者情報の取得
	sendingBillingAddressStaffList, err := i.sendingBillingAddressStaffRepository.GetByBillingAdressIDList(billingAddressIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 企業情報の取得 **************/

	for _, sendingBillingAddress := range sendingBillingAddressList {
		if !encounteredForEnterprise[sendingBillingAddress.SendingEnterpriseID] {
			encounteredForEnterprise[sendingBillingAddress.SendingEnterpriseID] = true
			enterpriseIDList = append(enterpriseIDList, sendingBillingAddress.SendingEnterpriseID)
		}
	}

	if len(enterpriseIDList) == 0 {
		return output, nil
	}

	sendingEnterpriseList, err := i.sendingEnterpriseRepository.GetByIDList(enterpriseIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, billingAddress := range sendingBillingAddressList {
		for _, jobInformation := range sendingJobInformationList {
			// 請求先に求人情報をセット
			if billingAddress.ID == jobInformation.SendingBillingAddressID {
				billingAddress.SendingJobInformationList = append(billingAddress.SendingJobInformationList, *jobInformation)
			}
		}

		// 請求先に担当者情報をセット
		for _, billingStaff := range sendingBillingAddressStaffList {
			if billingAddress.ID == billingStaff.SendingBillingAddressID {
				billingAddress.Staffs = append(billingAddress.Staffs, *billingStaff)
			}
		}
	}

	for _, enterprise := range sendingEnterpriseList {
		for _, billingAddress := range sendingBillingAddressList {
			// 企業に請求先の情報をセット
			if enterprise.ID == billingAddress.SendingEnterpriseID {
				enterprise.SendingBillingAddress = *billingAddress
			}
		}
	}

	output.SendingEnterpriseList = sendingEnterpriseList

	return output, nil
}

type SendSendingMailInput struct {
	SendParam entity.SendSendingMailParam
}

type SendSendingMailOutput struct {
	OK bool
}

func (i *SendingEnterpriseInteractorImpl) SendSendingMail(input SendSendingMailInput) (SendSendingMailOutput, error) {
	var (
		output SendSendingMailOutput
		err    error
	)

	from := entity.EmailUser{
		Name:  "アンドイーズ株式会社",
		Email: "info@spaceai.jp",
	}

	/************ 1. sending_job_seekerテーブルのphaseカラムを更新 **************/

	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByID(input.SendParam.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// sending_customerテーブルのphaseカラムがentity.AcceptSending(3)より進んでいない場合は更新
	if entity.SendingJobSeekerPhase(sendingJobSeeker.Phase.Int64) < entity.AcceptSending {
		err = i.sendingJobSeekerRepository.UpdatePhase(input.SendParam.SendingJobSeekerID, null.NewInt(int64(entity.AcceptSending), true))
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 2. sending_customerテーブルのphaseカラムを更新 **************/

	sendingCustomer, err := i.sendingCustomerRepository.FindByID(sendingJobSeeker.SendingCustomerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if entity.SendingJobSeekerPhase(sendingCustomer.Phase.Int64) < entity.AcceptSending {
		err = i.sendingCustomerRepository.UpdatePhase(sendingCustomer.ID, null.NewInt(int64(entity.AcceptSending), true))
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 3. sending_phasesテーブルにレコード作成・送客メール送信 **************/

	for _, sendingInfo := range input.SendParam.SendingList {

		sendingShareDocument := entity.NewSendingShareDocument(
			input.SendParam.SendingJobSeekerID,
			sendingInfo.SendingEnterpriseID,
			sendingInfo.IsShareUploadResume,
			sendingInfo.IsShareUploadCV,
			sendingInfo.IsShareUploadRecommendation,
			sendingInfo.IsShareGeneratedResume,
			sendingInfo.IsShareGeneratedCV,
			sendingInfo.IsShareGeneratedRecommendation,
		)

		err = i.sendingShareDocumentRepository.Create(sendingShareDocument)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		sendingPhase, err := i.sendingPhaseRepository.FindBySendingJobSeekerIDAndSendingEnterpriseID(input.SendParam.SendingJobSeekerID, sendingInfo.SendingEnterpriseID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				// sending_phasesテーブルにレコードが無い場合はレコード作成
				newSendingPhase := entity.NewSendingPhase(
					input.SendParam.SendingJobSeekerID,
					sendingInfo.SendingEnterpriseID,
					null.NewInt(int64(entity.AcceptSending), true), // 送客応諾
					sendingInfo.SendingDate,
					false,
				)

				err = i.sendingPhaseRepository.Create(newSendingPhase)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else {
				// エラーの場合
				fmt.Println(err)
				return output, err
			}
		} else {
			// すでにsending_phasesテーブルにレコードがある場合は更新
			err = i.sendingPhaseRepository.UpdatePhase(sendingPhase.ID, uint(entity.AcceptSending))
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		// メール送信
		err = utility.SendMailToMultiple(
			i.sendgrid.APIKey,
			sendingInfo.Mail.Subject,
			sendingInfo.Mail.Content,
			from,
			sendingInfo.Mail.Tos,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 4. 興味あり求人の更新 **************/

	// 興味あり求人のレコードを削除
	err = i.sendingJobSeekerDesiredJobInformationRepository.DeleteBySendingJobSeekerID(input.SendParam.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 興味あり求人のレコードを作成
	err = i.sendingJobSeekerDesiredJobInformationRepository.CreateMulti(input.SendParam.SendingJobSeekerID, input.SendParam.InterestingJobInfoIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type SendMailForRSVPInput struct {
	SendParam entity.SendSendingMailForRSVPParam
}

type SendMailForRSVPOutput struct {
	OK bool
}

// RSVP（参加確認）のメール送信
func (i *SendingEnterpriseInteractorImpl) SendMailForRSVP(input SendMailForRSVPInput) (SendMailForRSVPOutput, error) {
	var (
		output SendMailForRSVPOutput
		err    error
	)

	from := entity.EmailUser{
		Name:  "アンドイーズ株式会社",
		Email: "info@spaceai.jp",
	}

	/************ 1. 参加確認メールの送信 **************/

	for _, sendingInfo := range input.SendParam.SendingList {
		// メール送信
		err = utility.SendMailToMultiple(
			i.sendgrid.APIKey,
			sendingInfo.Mail.Subject,
			sendingInfo.Mail.Content,
			from,
			sendingInfo.Mail.Tos,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 2. 参加確認有無の更新 **************/

	for _, phaseID := range input.SendParam.SendingPhaseIDList {
		err = i.sendingPhaseRepository.UpdateIsAttended(phaseID, true)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true

	return output, nil
}

type GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDInput struct {
	AgentID uint
}

type GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDOutput struct {
	SendingEnterpriseList []*entity.SendingEnterprise
}

// 参加確認に使用
func (i *SendingEnterpriseInteractorImpl) GetSendingEnterpriseAndAcceptJobSeekerListByAgentID(input GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDInput) (GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDOutput, error) {
	var (
		output                 GetSendingEnterpriseAndAcceptJobSeekerListByAgentIDOutput
		err                    error
		responseEnterpriseList []*entity.SendingEnterprise
	)

	/************ 1. sending_enterpriseテーブルから全てのレコードを取得 **************/

	sendingEnterpriseList, err := i.sendingEnterpriseRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var sendingEnterpriseIDList []uint

	for _, sendingEnterprise := range sendingEnterpriseList {
		sendingEnterpriseIDList = append(sendingEnterpriseIDList, sendingEnterprise.ID)
	}

	/************ 2. sending_phasesテーブルからレコードを取得 **************/

	sendingPhaseList, err := i.sendingPhaseRepository.GetListByAgentIDAndPhaseAndEnterpriseIDList(input.AgentID, uint(entity.AcceptSending), sendingEnterpriseIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sendingEnterprise := range sendingEnterpriseList {
		for _, sendingPhase := range sendingPhaseList {
			if sendingEnterprise.ID == sendingPhase.SendingEnterpriseID && !sendingPhase.IsAttended {
				sendingEnterprise.AcceptSendingPhaseList = append(sendingEnterprise.AcceptSendingPhaseList, *sendingPhase)
			}
		}

		if len(sendingEnterprise.AcceptSendingPhaseList) > 0 {
			responseEnterpriseList = append(responseEnterpriseList, sendingEnterprise)
		}
	}

	output.SendingEnterpriseList = responseEnterpriseList

	return output, nil
}
