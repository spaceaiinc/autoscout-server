package interactor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type JobSeekerInteractor interface {
	// 汎用系 API
	CreateJobSeeker(input CreateJobSeekerInput) (CreateJobSeekerOutput, error)
	UpdateJobSeeker(input UpdateJobSeekerInput) (UpdateJobSeekerOutput, error)
	DeleteJobSeeker(input DeleteJobSeekerInput) (DeleteJobSeekerOutput, error)
	GetJobSeekerByID(input GetJobSeekerByIDInput) (GetJobSeekerByIDOutput, error)                                  // 指定IDの求職者を取得する関数
	GetJobSeekerByUUID(input GetJobSeekerByUUIDInput) (GetJobSeekerByUUIDOutput, error)                            // 指定uuidの求職者を取得する関数
	GetJobSeekerDocumentByUUID(input GetJobSeekerDocumentByUUIDInput) (GetJobSeekerDocumentByUUIDOutput, error)    // 指定uuidの求職者の応募書類データを取得する関数
	GetJobSeekerByTaskGroupUUID(input GetJobSeekerByTaskGroupUUIDInput) (GetJobSeekerByTaskGroupUUIDOutput, error) // 指定IDの求職者を取得する関数
	GetJobSeekerListByIDList(input GetJobSeekerListByIDListInput) (GetJobSeekerListByIDListOutput, error)
	GetDuplicateJobSeekerList(input GetDuplicateJobSeekerListInput) (GetDuplicateJobSeekerListOutput, error)
	GetSelectListForCreateOrUpdateJobSeekerByAgentID(input GetSelectListForCreateOrUpdateJobSeekerByAgentIDInput) (GetSelectListForCreateOrUpdateJobSeekerByAgentIDOutput, error)
	UpdateActivityMemoByJobSeekerID(input UpdateActivityMemoByJobSeekerIDInput) (UpdateActivityMemoByJobSeekerIDOutput, error)
	UpdateCanViewMatchingJob(input UpdateCanViewMatchingJobInput) (UpdateCanViewMatchingJobOutput, error)

	// 絞り込み検索
	GetSearchJobSeekerListByAgentID(input GetSearchJobSeekerListByAgentIDInput) (GetSearchJobSeekerListByAgentIDOutput, error)
	GetSearchActiveJobSeekerListByAgentID(input GetSearchActiveJobSeekerListByAgentIDInput) (GetSearchActiveJobSeekerListByAgentIDOutput, error)
	GetSearchAllianceJobSeekerListByAgentID(input GetSearchAllianceJobSeekerListByAgentIDInput) (GetSearchAllianceJobSeekerListByAgentIDOutput, error)

	// 求人検索→求職者検索 API
	GetSearchJobSeekerListByAgentIDAndType(input GetSearchJobSeekerListByAgentIDAndTypeInput) (GetSearchJobSeekerListByAgentIDAndTypeOutput, error) // 求人検索→求職者検索(絞り込み) API

	// シェア求人検索→自社求職者検索 API
	GetSearchPublicJobSeekerListByAgentIDAndPage(input GetSearchPublicJobSeekerListByAgentIDAndPageInput) (GetSearchPublicJobSeekerListByAgentIDAndPageOutput, error) // シェア求人検索→自社求職者検索(絞り込み) API

	// 求職者資料関連 API
	CreateJobSeekerDocument(input CreateJobSeekerDocumentInput) (CreateJobSeekerDocumentOutput, error)
	UpdateJobSeekerDocument(input UpdateJobSeekerDocumentInput) (UpdateJobSeekerDocumentOutput, error)
	UpdateJobSeekerDocumentForTask(input UpdateJobSeekerDocumentForTaskInput) (UpdateJobSeekerDocumentForTaskOutput, error)
	GetJobSeekerDocumentByJobSeekerID(input GetJobSeekerDocumentByJobSeekerIDInput) (GetJobSeekerDocumentByJobSeekerIDOutput, error)

	// csv API
	ImportJobSeekerCSV(input ImportJobSeekerCSVInput) (ImportJobSeekerCSVOutput, error) // csvから求職者一覧を登録する関数
	ExportJobSeekerCSV(input ExportJobSeekerCSVInput) (ExportJobSeekerCSVOutput, error) // 求職者一覧をcsvで取得する関数

	// LINE関連 API
	UpdateJobSeekerLineID(input UpdateJobSeekerLineIDInput) (UpdateJobSeekerLineIDOutput, error) // LINE IDを更新する関数

	// 面談前アンケート関連 API
	CreateInitialQuestionnaire(input CreateInitialQuestionnaireInput) (CreateInitialQuestionnairdhutput, error) // 面談前アンケートを登録 (求職者の同意項目、業界、職種、勤務地、質問要望、ファイル（履歴書（原本）、職務経歴書（原本）））

	// Admin API
	GetAllJobSeeker(input GetAllJobSeekerInput) (GetAllJobSeekerOutput, error) // 求職者一覧（すべて）

	DeleteJobSeekerResumePDFURL(input DeleteJobSeekerResumePDFURLInput) (DeleteJobSeekerResumePDFURLOutput, error)
	DeleteJobSeekerResumeOriginURL(input DeleteJobSeekerResumeOriginURLInput) (DeleteJobSeekerResumeOriginURLOutput, error)
	DeleteJobSeekerCVPDFURL(input DeleteJobSeekerCVPDFURLInput) (DeleteJobSeekerCVPDFURLOutput, error)
	DeleteJobSeekerCVOriginURL(input DeleteJobSeekerCVOriginURLInput) (DeleteJobSeekerCVOriginURLOutput, error)
	DeleteJobSeekerRecommendationPDFURL(input DeleteJobSeekerRecommendationPDFURLInput) (DeleteJobSeekerRecommendationPDFURLOutput, error)
	DeleteJobSeekerRecommendationOriginURL(input DeleteJobSeekerRecommendationOriginURLInput) (DeleteJobSeekerRecommendationOriginURLOutput, error)
	DeleteJobSeekerIDPhotoURL(input DeleteJobSeekerIDPhotoURLInput) (DeleteJobSeekerIDPhotoURLOutput, error)
	DeleteJobSeekerOtherDocument1URL(input DeleteJobSeekerOtherDocument1URLInput) (DeleteJobSeekerOtherDocument1URLOutput, error)
	DeleteJobSeekerOtherDocument2URL(input DeleteJobSeekerOtherDocument2URLInput) (DeleteJobSeekerOtherDocument2URLOutput, error)
	DeleteJobSeekerOtherDocument3URL(input DeleteJobSeekerOtherDocument3URLInput) (DeleteJobSeekerOtherDocument3URLOutput, error)

	// ゲストページ用 API
	GetJobSeekerForInitialStepByUUID(input GetJobSeekerForInitialStepByUUIDInput) (GetJobSeekerForInitialStepByUUIDOutput, error)    // 指定uuidの求職者を取得する関数
	GetGuestJobSeekerForByUUID(input GetGuestJobSeekerForByUUIDInput) (GetGuestJobSeekerForByUUIDOutput, error)                      // 指定uuidの求職者を取得する関数
	GetJobSeekerDesiredForGuestByUUID(input GetJobSeekerDesiredForGuestByUUIDInput) (GetJobSeekerDesiredForGuestByUUIDOutput, error) // 指定uuidの求職者を取得する関数
	GetJobSeekerAgentIDByUUID(input GetJobSeekerAgentIDByUUIDInput) (GetJobSeekerAgentIDByUUIDOutput, error)
	CheckJobSeekerByUUIDAndName(input CheckJobSeekerByUUIDAndNameInput) (CheckJobSeekerByUUIDAndNameOutput, error)
	UpdateJobSeekerPassword(input UpdateJobSeekerPasswordInput) (UpdateJobSeekerPasswordOutput, error)
	SendJobSeekerResetPasswordEmail(input SendJobSeekerResetPasswordEmailInput) (SendJobSeekerResetPasswordEmailOutput, error)
	SendJobSeekerContact(input SendJobSeekerContactInput) (SendJobSeekerContactOutput, error)
	UpdateInterviewDateByJobSeekerID(input UpdateInterviewDateByJobSeekerIDInput) (UpdateInterviewDateByJobSeekerIDOutput, error)

	// LP
	CreateJobSeekerFromLP(input CreateJobSeekerFromLPInput) (CreateJobSeekerFromLPOutput, error)
	UpdateJobSeekerPhoneFromLP(input UpdateJobSeekerPhoneFromLPInput) (UpdateJobSeekerPhoneFromLPOutput, error)
	UpdateJobSeekerDesiredFromLP(input UpdateJobSeekerDesiredFromLPInput) (UpdateJobSeekerDesiredFromLPOutput, error)
	GetJobSeekerLPRegisterStatusByUUID(input GetJobSeekerLPRegisterStatusByUUIDInput) (GetJobSeekerLPRegisterStatusByUUIDOutput, error)
	SendLPContact(input SendLPContactInput) (SendLPContactOutput, error)
	SendJobSeekerResetPasswordEmailForLP(input SendJobSeekerResetPasswordEmailForLPInput) (SendJobSeekerResetPasswordEmailForLPOutput, error)
	ResetPasswordForLP(input ResetPasswordForLPInput) (ResetPasswordForLPOutput, error)
	CheckResetPasswordToken(input CheckResetPasswordTokenInput) (CheckResetPasswordTokenOutput, error)
}

type JobSeekerInteractorImpl struct {
	firebase                                           usecase.Firebase
	sendgrid                                           config.Sendgrid
	oneSignal                                          config.OneSignal
	slack                                              config.Slack
	jobSeekerRepository                                usecase.JobSeekerRepository
	jobSeekerStudentHistoryRepository                  usecase.JobSeekerStudentHistoryRepository
	jobSeekerWorkHistoryRepository                     usecase.JobSeekerWorkHistoryRepository
	jobSeekerExperienceIndustryRepository              usecase.JobSeekerExperienceIndustryRepository
	jobSeekerDepartmentHistoryRepository               usecase.JobSeekerDepartmentHistoryRepository
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
	jobSeekerHideToAgentRepository                     usecase.JobSeekerHideToAgentRepository
	jobSeekerExperienceOccupationRepository            usecase.JobSeekerExperienceOccupationRepository
	jobSeekerDesiredCompanyScaleRepository             usecase.JobSeekerDesiredCompanyScaleRepository
	jobSeekerExperienceJobRepository                   usecase.JobSeekerExperienceJobRepository
	jobSeekerLPLoginTokenRepository                    usecase.JobSeekerLPLoginTokenRepository
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
	jobInfoOccupationRepository                        usecase.JobInformationOccupationRepository
	jobInfoRequiredConditionRepository                 usecase.JobInformationRequiredConditionRepository
	jobInfoRequiredLanguageTypeRepository              usecase.JobInformationRequiredLanguageTypeRepository
	jobInfoRequiredExperienceDevelopmentTypeRepository usecase.JobInformationRequiredExperienceDevelopmentTypeRepository
	agentRepository                                    usecase.AgentRepository
	agentStaffRepository                               usecase.AgentStaffRepository
	agentAllianceRepository                            usecase.AgentAllianceRepository
	agentInflowChannelOptionRepository                 usecase.AgentInflowChannelOptionRepository
	enterpriseProfileRepository                        usecase.EnterpriseProfileRepository
	enterpriseIndustryRepository                       usecase.EnterpriseIndustryRepository
	enterpriseReferenceMaterialRepository              usecase.EnterpriseReferenceMaterialRepository
	chatGroupWithJobSeekerRepository                   usecase.ChatGroupWithJobSeekerRepository
	chatMessageWithJobSeekerRepository                 usecase.ChatMessageWithJobSeekerRepository
	initialQuestionnaireRepository                     usecase.InitialQuestionnaireRepository
	initialQuestionnaireDesiredIndustryRepository      usecase.InitialQuestionnaireDesiredIndustryRepository
	initialQuestionnaireDesiredOccupationRepository    usecase.InitialQuestionnaireDesiredOccupationRepository
	initialQuestionnaireDesiredWorkLocationRepository  usecase.InitialQuestionnaireDesiredWorkLocationRepository
	taskGroupRepository                                usecase.TaskGroupRepository
	taskRepository                                     usecase.TaskRepository
	interviewTaskRepository                            usecase.InterviewTaskRepository
	interviewTaskGroupRepository                       usecase.InterviewTaskGroupRepository
	sendingJobSeekerRepository                         usecase.SendingJobSeekerRepository
	sendingJobSeekerStudentHistoryRepository           usecase.SendingJobSeekerStudentHistoryRepository
	sendingJobSeekerWorkHistoryRepository              usecase.SendingJobSeekerWorkHistoryRepository
	sendingJobSeekerExperienceIndustryRepository       usecase.SendingJobSeekerExperienceIndustryRepository
	sendingJobSeekerDepartmentHistoryRepository        usecase.SendingJobSeekerDepartmentHistoryRepository
	sendingJobSeekerLicenseRepository                  usecase.SendingJobSeekerLicenseRepository
	sendingJobSeekerSelfPromotionRepository            usecase.SendingJobSeekerSelfPromotionRepository
	sendingJobSeekerDocumentRepository                 usecase.SendingJobSeekerDocumentRepository
	sendingJobSeekerDesiredIndustryRepository          usecase.SendingJobSeekerDesiredIndustryRepository
	sendingJobSeekerDesiredOccupationRepository        usecase.SendingJobSeekerDesiredOccupationRepository
	sendingJobSeekerDesiredWorkLocationRepository      usecase.SendingJobSeekerDesiredWorkLocationRepository
	sendingJobSeekerDesiredHolidayTypeRepository       usecase.SendingJobSeekerDesiredHolidayTypeRepository
	sendingJobSeekerDevelopmentSkillRepository         usecase.SendingJobSeekerDevelopmentSkillRepository
	sendingJobSeekerLanguageSkillRepository            usecase.SendingJobSeekerLanguageSkillRepository
	sendingJobSeekerPCToolRepository                   usecase.SendingJobSeekerPCToolRepository
	sendingJobSeekerExperienceOccupationRepository     usecase.SendingJobSeekerExperienceOccupationRepository
	sendingJobSeekerDesiredCompanyScaleRepository      usecase.SendingJobSeekerDesiredCompanyScaleRepository
	chatGroupWithSendingJobSeekerRepository            usecase.ChatGroupWithSendingJobSeekerRepository
	chatMessageWithSendingJobSeekerRepository          usecase.ChatMessageWithSendingJobSeekerRepository
}

// JobSeekerInteractorImpl is an implementation of JobSeekerInteractor
func NewJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	sl config.Slack,
	jsR usecase.JobSeekerRepository,
	jsshR usecase.JobSeekerStudentHistoryRepository,
	jswhR usecase.JobSeekerWorkHistoryRepository,
	jseiR usecase.JobSeekerExperienceIndustryRepository,
	jseoR usecase.JobSeekerExperienceOccupationRepository,
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
	jshR usecase.JobSeekerHideToAgentRepository,
	jsdhR usecase.JobSeekerDepartmentHistoryRepository,
	jsdcsR usecase.JobSeekerDesiredCompanyScaleRepository,
	jsejR usecase.JobSeekerExperienceJobRepository,
	jlltR usecase.JobSeekerLPLoginTokenRepository,
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
	joR usecase.JobInformationOccupationRepository,
	jrcR usecase.JobInformationRequiredConditionRepository,
	jrltR usecase.JobInformationRequiredLanguageTypeRepository,
	jredtR usecase.JobInformationRequiredExperienceDevelopmentTypeRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	aaR usecase.AgentAllianceRepository,
	aicoR usecase.AgentInflowChannelOptionRepository,
	epR usecase.EnterpriseProfileRepository,
	eiR usecase.EnterpriseIndustryRepository,
	ermR usecase.EnterpriseReferenceMaterialRepository,
	cgR usecase.ChatGroupWithJobSeekerRepository,
	cmR usecase.ChatMessageWithJobSeekerRepository,
	iqR usecase.InitialQuestionnaireRepository,
	iqdiR usecase.InitialQuestionnaireDesiredIndustryRepository,
	iqdoR usecase.InitialQuestionnaireDesiredOccupationRepository,
	iqdwlR usecase.InitialQuestionnaireDesiredWorkLocationRepository,
	tgR usecase.TaskGroupRepository,
	tR usecase.TaskRepository,
	itR usecase.InterviewTaskRepository,
	itgR usecase.InterviewTaskGroupRepository,
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
	cgsR usecase.ChatGroupWithSendingJobSeekerRepository,
	cmsR usecase.ChatMessageWithSendingJobSeekerRepository,
) JobSeekerInteractor {
	return &JobSeekerInteractorImpl{
		firebase:                                           fb,
		sendgrid:                                           sg,
		oneSignal:                                          os,
		slack:                                              sl,
		jobSeekerRepository:                                jsR,
		jobSeekerStudentHistoryRepository:                  jsshR,
		jobSeekerWorkHistoryRepository:                     jswhR,
		jobSeekerExperienceIndustryRepository:              jseiR,
		jobSeekerExperienceOccupationRepository:            jseoR,
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
		jobSeekerHideToAgentRepository:                     jshR,
		jobSeekerDepartmentHistoryRepository:               jsdhR,
		jobSeekerDesiredCompanyScaleRepository:             jsdcsR,
		jobSeekerExperienceJobRepository:                   jsejR,
		jobSeekerLPLoginTokenRepository:                    jlltR,
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
		jobInfoOccupationRepository:                        joR,
		jobInfoRequiredConditionRepository:                 jrcR,
		jobInfoRequiredLanguageTypeRepository:              jrltR,
		jobInfoRequiredExperienceDevelopmentTypeRepository: jredtR,
		agentRepository:                                    aR,
		agentStaffRepository:                               asR,
		agentAllianceRepository:                            aaR,
		agentInflowChannelOptionRepository:                 aicoR,
		enterpriseProfileRepository:                        epR,
		enterpriseIndustryRepository:                       eiR,
		enterpriseReferenceMaterialRepository:              ermR,
		chatGroupWithJobSeekerRepository:                   cgR,
		chatMessageWithJobSeekerRepository:                 cmR,
		initialQuestionnaireRepository:                     iqR,
		initialQuestionnaireDesiredIndustryRepository:      iqdiR,
		initialQuestionnaireDesiredOccupationRepository:    iqdoR,
		initialQuestionnaireDesiredWorkLocationRepository:  iqdwlR,
		taskGroupRepository:                                tgR,
		taskRepository:                                     tR,
		interviewTaskRepository:                            itR,
		interviewTaskGroupRepository:                       itgR,
		sendingJobSeekerRepository:                         sjsR,
		sendingJobSeekerStudentHistoryRepository:           sjsshR,
		sendingJobSeekerWorkHistoryRepository:              sjswhR,
		sendingJobSeekerExperienceIndustryRepository:       sjseiR,
		sendingJobSeekerExperienceOccupationRepository:     sjseoR,
		sendingJobSeekerLicenseRepository:                  sjslR,
		sendingJobSeekerSelfPromotionRepository:            sjsspR,
		sendingJobSeekerDocumentRepository:                 sjsdR,
		sendingJobSeekerDesiredIndustryRepository:          sjsdiR,
		sendingJobSeekerDesiredOccupationRepository:        sjsdoR,
		sendingJobSeekerDesiredWorkLocationRepository:      sjsdwlR,
		sendingJobSeekerDesiredHolidayTypeRepository:       sjsdhtR,
		sendingJobSeekerDevelopmentSkillRepository:         sjsdsR,
		sendingJobSeekerLanguageSkillRepository:            sjslsR,
		sendingJobSeekerPCToolRepository:                   sjsptR,
		sendingJobSeekerDepartmentHistoryRepository:        sjsdhR,
		sendingJobSeekerDesiredCompanyScaleRepository:      sjsdcsR,
		chatGroupWithSendingJobSeekerRepository:            cgsR,
		chatMessageWithSendingJobSeekerRepository:          cmsR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//求職者の作成
type CreateJobSeekerInput struct {
	CreateParam  entity.CreateOrUpdateJobSeekerParam
	AgentStaffID uint
}

type CreateJobSeekerOutput struct {
	JobSeeker *entity.JobSeeker
}

func (i *JobSeekerInteractorImpl) CreateJobSeeker(input CreateJobSeekerInput) (CreateJobSeekerOutput, error) {
	var (
		output CreateJobSeekerOutput
		err    error
	)

	jobSeeker := entity.NewJobSeeker(
		input.CreateParam.AgentID,
		input.CreateParam.AgentStaffID,
		input.CreateParam.UserStatus,
		input.CreateParam.LastName,
		input.CreateParam.FirstName,
		input.CreateParam.LastFurigana,
		input.CreateParam.FirstFurigana,
		input.CreateParam.Gender,
		input.CreateParam.GenderRemarks,
		input.CreateParam.Birthday,
		input.CreateParam.Spouse,
		input.CreateParam.SupportObligation,
		input.CreateParam.Dependents,
		input.CreateParam.PhoneNumber,
		input.CreateParam.Email,
		input.CreateParam.EmergencyPhoneNumber,
		input.CreateParam.PostCode,
		input.CreateParam.Prefecture,
		input.CreateParam.Address,
		input.CreateParam.AddressFurigana,
		input.CreateParam.StateOfEmployment,
		input.CreateParam.JobSummary,
		input.CreateParam.HistorySupplement,
		input.CreateParam.ResearchContent,
		input.CreateParam.JoinCompanyPeriod,
		input.CreateParam.JobChange,
		input.CreateParam.AnnualIncome,
		input.CreateParam.DesiredAnnualIncome,
		input.CreateParam.Transfer,
		input.CreateParam.TransferRequirement,
		input.CreateParam.ShortResignation,
		input.CreateParam.ShortResignationRemarks,
		input.CreateParam.MedicalHistory,
		input.CreateParam.Nationality,
		input.CreateParam.Appearance,
		input.CreateParam.Communication,
		input.CreateParam.Thinking,
		input.CreateParam.RecommendationProfile,
		input.CreateParam.CandidProfile,
		input.CreateParam.SecretMemo,
		input.CreateParam.JobHuntingState,
		input.CreateParam.RecommendReason,
		input.CreateParam.Phase,
		input.CreateParam.InterviewDate,
		input.CreateParam.RegisterPhase,
		input.CreateParam.StudyCategory,
		input.CreateParam.WordSkill,
		input.CreateParam.ExcelSkill,
		input.CreateParam.PowerPointSkill,
		input.CreateParam.InflowChannelID,
		input.CreateParam.NationalityRemarks,
		input.CreateParam.MedicalHistoryRemarks,
		input.CreateParam.AcceptancePoints,
	)

	err = i.jobSeekerRepository.Create(jobSeeker)
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

	for _, sh := range input.CreateParam.StudentHistories {
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

	for _, wh := range input.CreateParam.WorkHistories {
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

		for _, dh := range wh.DepartmentHistories {
			departmentHistory := entity.NewJobSeekerDepartmentHistory(
				workHistory.ID,
				dh.Department,
				dh.ManagementNumber,
				dh.ManagementDetail,
				dh.JobDescription,
				dh.StartYear,
				dh.EndYear,
			)

			err = i.jobSeekerDepartmentHistoryRepository.Create(departmentHistory)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, occupations := range dh.ExperienceOccupations {
				occupation := entity.NewJobSeekerExperienceOccupation(
					departmentHistory.ID,
					occupations.Occupation,
				)

				err = i.jobSeekerExperienceOccupationRepository.Create(occupation)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

		}

	}

	for _, dcs := range input.CreateParam.DesiredCompanyScales {
		disiredCompanyScale := entity.NewJobSeekerDesiredCompanyScale(
			jobSeeker.ID,
			dcs.DesiredCompanyScale,
		)

		err = i.jobSeekerDesiredCompanyScaleRepository.Create(disiredCompanyScale)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, l := range input.CreateParam.Licenses {
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

	for _, sp := range input.CreateParam.SelfPromotions {
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

	for _, di := range input.CreateParam.DesiredIndustries {
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

	for _, do := range input.CreateParam.DesiredOccupations {
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

	for _, dwl := range input.CreateParam.DesiredWorkLocations {
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

	for _, dht := range input.CreateParam.DesiredHolidayTypes {
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

	for _, ds := range input.CreateParam.DevelopmentSkills {
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

	for _, ls := range input.CreateParam.LanguageSkills {
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

	for _, ps := range input.CreateParam.PCTools {
		PCTool := entity.NewJobSeekerPCTool(
			jobSeeker.ID,
			ps.Tool,
		)

		err = i.jobSeekerPCToolRepository.Create(PCTool)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// hideToAgent
	for _, hta := range input.CreateParam.HideToAgents {
		hideToAgent := entity.NewJobSeekerHideToAgent(
			jobSeeker.ID,
			hta.AgentID,
		)

		err = i.jobSeekerHideToAgentRepository.Create(hideToAgent)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
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

	// str := input.CreateParam.InterviewDate
	// layout := "2006-01-02T15:04"
	// t, _ := time.Parse(layout, str)
	// y, m, d := input.CreateParam.InterviewDate.Date()

	// 入寮したフェーズに応じて作成タスクを変更
	switch input.CreateParam.Phase {
	case EntryInterview: // エントリー
		phaseSub = null.NewInt(0, true) // 日程調整依頼
		// 当日（2020-02-01）
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

	case InvitationInterview: // 面談案内済み
		phaseSub = null.NewInt(0, true) // 面談調整中
		// 当日（2020-02-01）
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

	case ReservationInterview: // 面談予約完了
		phaseSub = null.NewInt(0, true) // 面談の前日確認を行う
		// input.CreateParam.InterviewDate（2022-12-15T16:29）の前日（2022-12-14）
		yesterday := input.CreateParam.InterviewDate.AddDate(0, 0, -1)
		date = yesterday.Format("2006-01-02")

	case WaitingInterview: // 面談実施待ち
		phaseSub = null.NewInt(1, true) // 面談実施待ち
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

	case PreparingAfterInterview: // 面談実施済み（準備中）
		phaseSub = null.NewInt(99, true) // 終了
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

		// 面談実施済みの場合は初回面談日時に記録する
		firstInterviewDate = input.CreateParam.InterviewDate

	case OperatingAfterInterview: // 面談実施済み（稼働中）
		phaseSub = null.NewInt(99, true) // 終了
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

		// 面談実施済みの場合は初回面談日時に記録する
		firstInterviewDate = input.CreateParam.InterviewDate

	case ReleasedAfterInterview: // 面談実施済み（リリース状態）
		phaseSub = null.NewInt(99, true) // 終了
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

		// 面談実施済みの場合は初回面談日時に記録する
		firstInterviewDate = input.CreateParam.InterviewDate

	case OfferedAfterInterview: // サービス終了/決定者
		phaseSub = null.NewInt(99, true) // 終了
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

	case ContinuingAfterInterview: // サービス終了/今後継続連絡
		phaseSub = null.NewInt(99, true) // 終了
		date = input.CreateParam.InterviewDate.Format("2006-01-02")

	case QuitedAfterInterview: // サービス終了/転職活動終了
		phaseSub = null.NewInt(99, true) // 終了
		date = input.CreateParam.InterviewDate.Format("2006-01-02")
	}

	// 面談調整タスクの作成
	interviewTaskGroup := entity.NewInterviewTaskGroup(
		input.CreateParam.AgentID,
		jobSeeker.ID,
		input.CreateParam.InterviewDate,
		firstInterviewDate, // 面談実施済みの場合は初回面談日時に記録する
	)

	err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	interviewTask := entity.NewInterviewTask(
		interviewTaskGroup.ID,
		null.NewInt(int64(input.AgentStaffID), true),
		input.CreateParam.AgentStaffID,
		input.CreateParam.Phase,
		phaseSub,
		"",
		date,
		null.NewInt(99, true),
		getStrPhaseForJobSeeker(jobSeeker.Phase), // SelectActionLabelは求職者のフェーズにする
	)

	err = i.interviewTaskRepository.Create(interviewTask)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeeker = jobSeeker
	output.JobSeeker.StudentHistories = input.CreateParam.StudentHistories
	output.JobSeeker.WorkHistories = input.CreateParam.WorkHistories
	output.JobSeeker.Licenses = input.CreateParam.Licenses
	output.JobSeeker.SelfPromotions = input.CreateParam.SelfPromotions
	output.JobSeeker.DesiredIndustries = input.CreateParam.DesiredIndustries
	output.JobSeeker.DesiredOccupations = input.CreateParam.DesiredOccupations
	output.JobSeeker.DesiredWorkLocations = input.CreateParam.DesiredWorkLocations
	output.JobSeeker.DesiredHolidayTypes = input.CreateParam.DesiredHolidayTypes
	output.JobSeeker.DevelopmentSkills = input.CreateParam.DevelopmentSkills
	output.JobSeeker.LanguageSkills = input.CreateParam.LanguageSkills
	output.JobSeeker.PCTools = input.CreateParam.PCTools
	output.JobSeeker.HideToAgents = input.CreateParam.HideToAgents

	return output, nil
}

// 求人企業の更新
type UpdateJobSeekerInput struct {
	UpdateParam  entity.CreateOrUpdateJobSeekerParam
	JobSeekerID  uint
	AgentStaffID uint
}

type UpdateJobSeekerOutput struct {
	JobSeeker *entity.JobSeeker
}

func (i *JobSeekerInteractorImpl) UpdateJobSeeker(input UpdateJobSeekerInput) (UpdateJobSeekerOutput, error) {
	var (
		output UpdateJobSeekerOutput
		err    error
	)

	jobSeeker := entity.NewJobSeeker(
		input.UpdateParam.AgentID, // 更新しない
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.UserStatus,
		input.UpdateParam.LastName,
		input.UpdateParam.FirstName,
		input.UpdateParam.LastFurigana,
		input.UpdateParam.FirstFurigana,
		input.UpdateParam.Gender,
		input.UpdateParam.GenderRemarks,
		input.UpdateParam.Birthday,
		input.UpdateParam.Spouse,
		input.UpdateParam.SupportObligation,
		input.UpdateParam.Dependents,
		input.UpdateParam.PhoneNumber,
		input.UpdateParam.Email,
		input.UpdateParam.EmergencyPhoneNumber,
		input.UpdateParam.PostCode,
		input.UpdateParam.Prefecture,
		input.UpdateParam.Address,
		input.UpdateParam.AddressFurigana,
		input.UpdateParam.StateOfEmployment,
		input.UpdateParam.JobSummary,
		input.UpdateParam.HistorySupplement,
		input.UpdateParam.ResearchContent,
		input.UpdateParam.JoinCompanyPeriod,
		input.UpdateParam.JobChange,
		input.UpdateParam.AnnualIncome,
		input.UpdateParam.DesiredAnnualIncome,
		input.UpdateParam.Transfer,
		input.UpdateParam.TransferRequirement,
		input.UpdateParam.ShortResignation,
		input.UpdateParam.ShortResignationRemarks,
		input.UpdateParam.MedicalHistory,
		input.UpdateParam.Nationality,
		input.UpdateParam.Appearance,
		input.UpdateParam.Communication,
		input.UpdateParam.Thinking,
		input.UpdateParam.RecommendationProfile,
		input.UpdateParam.CandidProfile,
		input.UpdateParam.SecretMemo,
		input.UpdateParam.JobHuntingState,
		input.UpdateParam.RecommendReason,
		input.UpdateParam.Phase,
		input.UpdateParam.InterviewDate,
		input.UpdateParam.RegisterPhase,
		input.UpdateParam.StudyCategory,
		input.UpdateParam.WordSkill,
		input.UpdateParam.ExcelSkill,
		input.UpdateParam.PowerPointSkill,
		input.UpdateParam.InflowChannelID,
		input.UpdateParam.NationalityRemarks,
		input.UpdateParam.MedicalHistoryRemarks,
		input.UpdateParam.AcceptancePoints,
	)

	err = i.jobSeekerRepository.Update(input.JobSeekerID, jobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.jobSeekerStudentHistoryRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sh := range input.UpdateParam.StudentHistories {
		studentHistory := entity.NewJobSeekerStudentHistory(
			input.JobSeekerID,
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

	err = i.jobSeekerWorkHistoryRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, wh := range input.UpdateParam.WorkHistories {
		workHistory := entity.NewJobSeekerWorkHistory(
			input.JobSeekerID,
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

		for _, dh := range wh.DepartmentHistories {
			departmentHistory := entity.NewJobSeekerDepartmentHistory(
				workHistory.ID,
				dh.Department,
				dh.ManagementNumber,
				dh.ManagementDetail,
				dh.JobDescription,
				dh.StartYear,
				dh.EndYear,
			)

			err = i.jobSeekerDepartmentHistoryRepository.Create(departmentHistory)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, occupations := range dh.ExperienceOccupations {
				occupation := entity.NewJobSeekerExperienceOccupation(
					departmentHistory.ID,
					occupations.Occupation,
				)

				err = i.jobSeekerExperienceOccupationRepository.Create(occupation)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

		}

	}

	err = i.jobSeekerDesiredCompanyScaleRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, dcs := range input.UpdateParam.DesiredCompanyScales {
		disiredCompanyScale := entity.NewJobSeekerDesiredCompanyScale(
			input.JobSeekerID,
			dcs.DesiredCompanyScale,
		)

		err = i.jobSeekerDesiredCompanyScaleRepository.Create(disiredCompanyScale)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobSeekerLicenseRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, l := range input.UpdateParam.Licenses {
		license := entity.NewJobSeekerLicense(
			input.JobSeekerID,
			l.LicenseType,
			l.AcquisitionTime,
		)

		err = i.jobSeekerLicenseRepository.Create(license)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobSeekerSelfPromotionRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sp := range input.UpdateParam.SelfPromotions {
		selfPromotion := entity.NewJobSeekerSelfPromotion(
			input.JobSeekerID,
			sp.Title,
			sp.Contents,
		)

		err = i.jobSeekerSelfPromotionRepository.Create(selfPromotion)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobSeekerDesiredIndustryRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, di := range input.UpdateParam.DesiredIndustries {
		desiredIndustry := entity.NewJobSeekerDesiredIndustry(
			input.JobSeekerID,
			di.DesiredIndustry,
			di.DesiredRank,
		)

		err = i.jobSeekerDesiredIndustryRepository.Create(desiredIndustry)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobSeekerDesiredOccupationRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, do := range input.UpdateParam.DesiredOccupations {
		desiredOccupation := entity.NewJobSeekerDesiredOccupation(
			input.JobSeekerID,
			do.DesiredOccupation,
			do.DesiredRank,
		)

		err = i.jobSeekerDesiredOccupationRepository.Create(desiredOccupation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobSeekerDesiredWorkLocationRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, dwl := range input.UpdateParam.DesiredWorkLocations {
		desiredWorkLocation := entity.NewJobSeekerDesiredWorkLocation(
			input.JobSeekerID,
			dwl.DesiredWorkLocation,
			dwl.DesiredRank,
		)

		err = i.jobSeekerDesiredWorkLocationRepository.Create(desiredWorkLocation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobSeekerDesiredHolidayTypeRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, dht := range input.UpdateParam.DesiredHolidayTypes {
		desiredHolidayType := entity.NewJobSeekerDesiredHolidayType(
			input.JobSeekerID,
			dht.HolidayType,
		)

		err = i.jobSeekerDesiredHolidayTypeRepository.Create(desiredHolidayType)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.jobSeekerDevelopmentSkillRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ds := range input.UpdateParam.DevelopmentSkills {
		developmentSkill := entity.NewJobSeekerDevelopmentSkill(
			input.JobSeekerID,
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

	err = i.jobSeekerLanguageSkillRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ls := range input.UpdateParam.LanguageSkills {
		languageSkill := entity.NewJobSeekerLanguageSkill(
			input.JobSeekerID,
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

	err = i.jobSeekerPCToolRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ps := range input.UpdateParam.PCTools {
		PCTool := entity.NewJobSeekerPCTool(
			input.JobSeekerID,
			ps.Tool,
		)

		err = i.jobSeekerPCToolRepository.Create(PCTool)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// hideToAgent
	err = i.jobSeekerHideToAgentRepository.DeleteByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hta := range input.UpdateParam.HideToAgents {
		hideToAgent := entity.NewJobSeekerHideToAgent(
			input.JobSeekerID,
			hta.AgentID,
		)

		err = i.jobSeekerHideToAgentRepository.Create(hideToAgent)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
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

	// 面談実施待ちより前のフェーズの場合は「面談調整タスク」を作成
	// if input.UpdateParam.Phase == EntryInterview ||
	// 	input.UpdateParam.Phase == InvitationInterview ||
	// 	input.UpdateParam.Phase == ReservationInterview ||
	// 	input.UpdateParam.Phase == WaitingInterview {

	var (
		phaseSub           null.Int
		date               string
		interviewTaskGroup *entity.InterviewTaskGroup
	)

	// 入寮したフェーズに応じて作成タスクを変更
	switch input.UpdateParam.Phase {
	case EntryInterview: // エントリー
		phaseSub = null.NewInt(0, true) // 日程調整依頼

		// 当日（2020-02-01）
		now := time.Now()
		date = now.Format("2006-01-02")

	case InvitationInterview: // 面談案内済み
		phaseSub = null.NewInt(0, true) // 面談調整中

		// 当日（2020-02-01）
		now := time.Now()
		date = now.Format("2006-01-02")

	case ReservationInterview: // 面談予約完了
		phaseSub = null.NewInt(0, true) // 面談の前日確認を行う

		// input.UpdateParam.InterviewDate（2022-12-15T16:29）の前日（2022-12-14）
		yesterday := input.UpdateParam.InterviewDate.AddDate(0, 0, -1)
		date = yesterday.Format("2006-01-02")

	case WaitingInterview: // 面談実施待ち
		phaseSub = null.NewInt(1, true) // 面談実施待ち

		// input.UpdateParam.InterviewDate（2022-12-15T16:29）当日（2022-12-14）
		date = input.UpdateParam.InterviewDate.Format("2006-01-02")

	case PreparingAfterInterview: // 面談実施済み（準備中）
		phaseSub = null.NewInt(99, true) // 終了
		now := time.Now()
		date = now.Format("2006-01-02")

	case OperatingAfterInterview: // 面談実施済み（稼働中）
		phaseSub = null.NewInt(99, true) // 終了
		now := time.Now()
		date = now.Format("2006-01-02")

	case ReleasedAfterInterview: // 面談実施済み（リリース状態）
		phaseSub = null.NewInt(99, true) // 終了
		now := time.Now()
		date = now.Format("2006-01-02")

	case OfferedAfterInterview: // サービス終了/決定者
		phaseSub = null.NewInt(99, true) // 終了
		now := time.Now()
		date = now.Format("2006-01-02")

	case ContinuingAfterInterview: // サービス終了/今後継続連絡
		phaseSub = null.NewInt(99, true) // 終了
		now := time.Now()
		date = now.Format("2006-01-02")

	case QuitedAfterInterview: // サービス終了/転職活動終了
		phaseSub = null.NewInt(99, true) // 終了
		now := time.Now()
		date = now.Format("2006-01-02")
	}

	// 面談調整タスクの作成
	interviewTaskGroup, err = i.interviewTaskGroupRepository.FindByAgentIDAndJobSeekerID(input.UpdateParam.AgentID, input.JobSeekerID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			interviewTaskGroup = entity.NewInterviewTaskGroup(
				input.UpdateParam.AgentID,
				input.JobSeekerID,
				input.UpdateParam.InterviewDate,
				utility.EarliestTime(), // 初期値
			)

			err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			interviewTask := entity.NewInterviewTask(
				interviewTaskGroup.ID,
				null.NewInt(int64(input.AgentStaffID), true),
				input.UpdateParam.AgentStaffID,
				input.UpdateParam.Phase,
				phaseSub,
				"",
				date,
				null.NewInt(99, true),
				getStrPhaseForJobSeeker(input.UpdateParam.Phase), // SelectActionLabelは求職者のフェーズにする
			)

			err = i.interviewTaskRepository.Create(interviewTask)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			fmt.Println(err)
			return output, err
		}
	} else {
		// 最新の面談タスクを取得
		latestInterviewTask, err := i.interviewTaskRepository.FindLatestByAgentIDAndJobSeekerID(input.UpdateParam.AgentID, input.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最新のタスクと編集後のタスクが違う場合はタスクを作成するする
		if latestInterviewTask.PhaseCategory != input.UpdateParam.Phase {
			interviewTask := entity.NewInterviewTask(
				interviewTaskGroup.ID,
				null.NewInt(int64(input.AgentStaffID), true),
				input.UpdateParam.AgentStaffID,
				input.UpdateParam.Phase,
				phaseSub,
				"",
				date,
				null.NewInt(99, true),
				getStrPhaseForJobSeeker(input.UpdateParam.Phase), // SelectActionLabelは求職者のフェーズにする
			)

			err = i.interviewTaskRepository.Create(interviewTask)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// すでにタスクグループが存在する場合は依頼時間を更新
			err = i.interviewTaskGroupRepository.UpdateLastRequestAt(interviewTaskGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			fmt.Println("同一タスクです")
			fmt.Println("---------")
		}
	}

	if input.UpdateParam.AgentStaffID.Valid {
		agentStaff, err := i.agentStaffRepository.FindByID(uint(input.UpdateParam.AgentStaffID.Int64))
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		jobSeeker.StaffName = agentStaff.StaffName
	}

	// 流入経路の名前を取得 *更新直後の画面で表示するため
	if input.UpdateParam.InflowChannelID.Valid {
		inflowChannel, err := i.agentInflowChannelOptionRepository.FindByID(uint(input.UpdateParam.InflowChannelID.Int64))
		if err != nil {
			fmt.Println(err)
			return output, err
		}
		jobSeeker.ChannelName = inflowChannel.ChannelName
	}

	jobSeeker.ID = input.JobSeekerID
	jobSeeker.UUID = input.UpdateParam.UUID
	jobSeeker.StudentHistories = input.UpdateParam.StudentHistories
	jobSeeker.WorkHistories = input.UpdateParam.WorkHistories
	jobSeeker.DesiredCompanyScales = input.UpdateParam.DesiredCompanyScales
	jobSeeker.Licenses = input.UpdateParam.Licenses
	jobSeeker.SelfPromotions = input.UpdateParam.SelfPromotions
	jobSeeker.DesiredIndustries = input.UpdateParam.DesiredIndustries
	jobSeeker.DesiredOccupations = input.UpdateParam.DesiredOccupations
	jobSeeker.DesiredWorkLocations = input.UpdateParam.DesiredWorkLocations
	jobSeeker.DesiredHolidayTypes = input.UpdateParam.DesiredHolidayTypes
	jobSeeker.DevelopmentSkills = input.UpdateParam.DevelopmentSkills
	jobSeeker.LanguageSkills = input.UpdateParam.LanguageSkills
	jobSeeker.PCTools = input.UpdateParam.PCTools
	jobSeeker.HideToAgents = input.UpdateParam.HideToAgents

	output.JobSeeker = jobSeeker

	return output, nil
}

// 求職者の削除
type DeleteJobSeekerInput struct {
	DeleteParam entity.DeleteJobSeekerParam
}

type DeleteJobSeekerOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeeker(input DeleteJobSeekerInput) (DeleteJobSeekerOutput, error) {
	var output DeleteJobSeekerOutput

	err := i.jobSeekerRepository.Delete(input.DeleteParam.ID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// 求職者IDを使って求職者情報を取得する
type GetJobSeekerByIDInput struct {
	JobSeekerID uint
}

type GetJobSeekerByIDOutput struct {
	JobSeeker *entity.JobSeeker
}

func (i *JobSeekerInteractorImpl) GetJobSeekerByID(input GetJobSeekerByIDInput) (GetJobSeekerByIDOutput, error) {
	var (
		output GetJobSeekerByIDOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者の子テーブル情報をセット
	jobSeeker, err = getJobSeekerChildTableData(jobSeeker, i)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeeker = jobSeeker

	return output, nil
}

// 求職者IDを使って求職者情報を取得する
type GetJobSeekerByUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetJobSeekerByUUIDOutput struct {
	JobSeeker *entity.JobSeeker
}

func (i *JobSeekerInteractorImpl) GetJobSeekerByUUID(input GetJobSeekerByUUIDInput) (GetJobSeekerByUUIDOutput, error) {
	var (
		output GetJobSeekerByUUIDOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeeker = jobSeeker

	return output, nil
}

// 求職者IDを使って求職者情報を取得する
type GetJobSeekerDocumentByUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetJobSeekerDocumentByUUIDOutput struct {
	Document *entity.JobSeekerDocument
}

func (i *JobSeekerInteractorImpl) GetJobSeekerDocumentByUUID(input GetJobSeekerDocumentByUUIDInput) (GetJobSeekerDocumentByUUIDOutput, error) {
	var (
		output GetJobSeekerDocumentByUUIDOutput
		err    error
	)

	document, err := i.jobSeekerDocumentRepository.FindByJobSeekerUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Document = document

	return output, nil
}

// タスクグループuuidを使って求職者情報を取得する
type GetJobSeekerByTaskGroupUUIDInput struct {
	TaskGroupUUID uuid.UUID
}

type GetJobSeekerByTaskGroupUUIDOutput struct {
	JobSeeker *entity.JobSeeker
}

func (i *JobSeekerInteractorImpl) GetJobSeekerByTaskGroupUUID(input GetJobSeekerByTaskGroupUUIDInput) (GetJobSeekerByTaskGroupUUIDOutput, error) {
	var (
		output GetJobSeekerByTaskGroupUUIDOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByTaskGroupUUID(input.TaskGroupUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者IDからスタッフ名を取得
	if jobSeeker.AgentStaffID.Valid {
		jobSeeker.StaffName, err = i.agentStaffRepository.FindStaffNameByJobSeekerID(jobSeeker.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selfPromotion, err := i.jobSeekerSelfPromotionRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	document, err := i.jobSeekerDocumentRepository.FindByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	pcSkill, err := i.jobSeekerPCToolRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgent, err := i.jobSeekerHideToAgentRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
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

	output.JobSeeker = jobSeeker

	return output, nil
}

type GetJobSeekerListByIDListInput struct {
	IDList  []uint
	AgentID uint
}

type GetJobSeekerListByIDListOutput struct {
	JobSeekerList []*entity.JobSeeker
}

func (i *JobSeekerInteractorImpl) GetJobSeekerListByIDList(input GetJobSeekerListByIDListInput) (GetJobSeekerListByIDListOutput, error) {
	var output GetJobSeekerListByIDListOutput

	jobSeekerList, err := i.jobSeekerRepository.GetByIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	studentHistory, err := i.jobSeekerStudentHistoryRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	license, err := i.jobSeekerLicenseRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selfPromotion, err := i.jobSeekerSelfPromotionRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	document, err := i.jobSeekerDocumentRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	pcSkill, err := i.jobSeekerPCToolRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgent, err := i.jobSeekerHideToAgentRepository.GetByJobSeekerIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, jobSeeker := range jobSeekerList {
		if jobSeeker.AgentID != input.AgentID {
			// 自社求職者でない場合はLINEを取得しない
			jobSeeker.LineID = ""
		}

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
									ID:                  eo.ID,
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
					JobSeekerID:       d.JobSeekerID,
					ResumeOriginURL:   d.ResumeOriginURL,
					ResumePDFURL:      d.ResumePDFURL,
					CVOriginURL:       d.CVOriginURL,
					CVPDFURL:          d.CVPDFURL,
					IDPhotoURL:        d.IDPhotoURL,
					OtherDocument1URL: d.OtherDocument1URL,
					OtherDocument2URL: d.OtherDocument2URL,
					OtherDocument3URL: d.OtherDocument3URL,
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

	output.JobSeekerList = jobSeekerList

	return output, nil
}

type GetDuplicateJobSeekerListInput struct {
	AgentID       uint
	LastName      string
	FirstName     string
	LastFurigana  string
	FirstFurigana string
	Email         string
	PhoneNumber   string
}

type GetDuplicateJobSeekerListOutput struct {
	JobSeekerList []*entity.JobSeeker
}

func (i *JobSeekerInteractorImpl) GetDuplicateJobSeekerList(input GetDuplicateJobSeekerListInput) (GetDuplicateJobSeekerListOutput, error) {
	var output GetDuplicateJobSeekerListOutput

	duplicateJobSeekerList, err := i.jobSeekerRepository.GetDuplicateByNameAndFuriganaAndEmailAndPhoneNumber(
		input.AgentID, input.LastName, input.FirstName, input.LastFurigana, input.FirstFurigana, input.Email, input.PhoneNumber,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerList = duplicateJobSeekerList

	return output, nil
}

type GetSelectListForCreateOrUpdateJobSeekerByAgentIDInput struct {
	Token   string
	AgentID uint
}

type GetSelectListForCreateOrUpdateJobSeekerByAgentIDOutput struct {
	AgentStaffList               []*entity.AgentStaff
	AgentInflowChannelOptionList []*entity.AgentInflowChannelOption
	AllianceAgentList            []*entity.Agent
}

func (i *JobSeekerInteractorImpl) GetSelectListForCreateOrUpdateJobSeekerByAgentID(input GetSelectListForCreateOrUpdateJobSeekerByAgentIDInput) (GetSelectListForCreateOrUpdateJobSeekerByAgentIDOutput, error) {
	var output GetSelectListForCreateOrUpdateJobSeekerByAgentIDOutput

	agentStaffList, err := i.agentStaffRepository.GetByAgentIDAndUsageStatusAvailable(input.AgentID)
	if err != nil {
		return output, err
	}

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	//  FirebaseIDが一致する担当者を一番上にソート
	for agentStaffI, agentStaff := range agentStaffList {
		if agentStaff.FirebaseID == firebaseID {
			agentStaffList = append(agentStaffList[:agentStaffI], agentStaffList[agentStaffI+1:]...)
			agentStaffList = append([]*entity.AgentStaff{agentStaff}, agentStaffList...)
			break
		}
	}

	allianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	allianceAgentList := getAllianceAgentList(input.AgentID, allianceList)

	// エージェントの流入経路マスタを取得
	agentInflowChannelOptionList, err := i.agentInflowChannelOptionRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentStaffList = agentStaffList
	output.AgentInflowChannelOptionList = agentInflowChannelOptionList
	output.AllianceAgentList = allianceAgentList

	return output, nil
}

// 求人企業の更新
type UpdateActivityMemoByJobSeekerIDInput struct {
	Param       entity.ActivityMemoParam
	JobSeekerID uint
}

type UpdateActivityMemoByJobSeekerIDOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) UpdateActivityMemoByJobSeekerID(input UpdateActivityMemoByJobSeekerIDInput) (UpdateActivityMemoByJobSeekerIDOutput, error) {
	var (
		output UpdateActivityMemoByJobSeekerIDOutput
		err    error
	)

	// 求職者テーブルの個人情報同意の有無を更新
	err = i.jobSeekerRepository.UpdateActivityMemo(input.JobSeekerID, input.Param.ActivityMemo)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// マッチング求人を閲覧可能かを管理する値
type UpdateCanViewMatchingJobInput struct {
	Param entity.UpdateJobSeekerCanViewMatchingJobParam
}

type UpdateCanViewMatchingJobOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) UpdateCanViewMatchingJob(input UpdateCanViewMatchingJobInput) (UpdateCanViewMatchingJobOutput, error) {
	var (
		output UpdateCanViewMatchingJobOutput
		err    error
		param  = input.Param
	)

	// マッチング求人を閲覧可能かを管理する値を更新
	err = i.jobSeekerRepository.UpdateCanViewMatchingJob(param.JobSeekerID, param.CanViewMatchingJob)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
// 求職者資料関連API
//
type CreateJobSeekerDocumentInput struct {
	CreateParam entity.CreateOrUpdateJobSeekerDocumentParam
}

type CreateJobSeekerDocumentOutput struct {
	JobSeekerDocument *entity.JobSeekerDocument
}

func (i *JobSeekerInteractorImpl) CreateJobSeekerDocument(input CreateJobSeekerDocumentInput) (CreateJobSeekerDocumentOutput, error) {
	var output CreateJobSeekerDocumentOutput

	jobSeekerDocument := entity.NewJobSeekerDocument(
		input.CreateParam.JobSeekerID,
		input.CreateParam.ResumeOriginURL,
		input.CreateParam.ResumePDFURL,
		input.CreateParam.CVOriginURL,
		input.CreateParam.CVPDFURL,
		input.CreateParam.RecommendationOriginURL,
		input.CreateParam.RecommendationPDFURL,
		input.CreateParam.IDPhotoURL,
		input.CreateParam.OtherDocument1URL,
		input.CreateParam.OtherDocument2URL,
		input.CreateParam.OtherDocument3URL,
	)

	err := i.jobSeekerDocumentRepository.Create(jobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerDocument = jobSeekerDocument

	return output, nil
}

type UpdateJobSeekerDocumentInput struct {
	UpdateParam entity.CreateOrUpdateJobSeekerDocumentParam
}

type UpdateJobSeekerDocumentOutput struct {
	JobSeekerDocument *entity.JobSeekerDocument
}

func (i *JobSeekerInteractorImpl) UpdateJobSeekerDocument(input UpdateJobSeekerDocumentInput) (UpdateJobSeekerDocumentOutput, error) {
	var output UpdateJobSeekerDocumentOutput

	jobSeekerDocument := entity.NewJobSeekerDocument(
		input.UpdateParam.JobSeekerID,
		input.UpdateParam.ResumeOriginURL,
		input.UpdateParam.ResumePDFURL,
		input.UpdateParam.CVOriginURL,
		input.UpdateParam.CVPDFURL,
		input.UpdateParam.RecommendationOriginURL,
		input.UpdateParam.RecommendationPDFURL,
		input.UpdateParam.IDPhotoURL,
		input.UpdateParam.OtherDocument1URL,
		input.UpdateParam.OtherDocument2URL,
		input.UpdateParam.OtherDocument3URL,
	)

	err := i.jobSeekerDocumentRepository.UpdateByJobSeekerID(input.UpdateParam.JobSeekerID, jobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerDocument = jobSeekerDocument

	return output, nil
}

type UpdateJobSeekerDocumentForTaskInput struct {
	UpdateParam        entity.CreateOrUpdateJobSeekerDocumentParam
	JobSeekerUUID      uuid.UUID
	JobInformationUUID uuid.UUID
}

type UpdateJobSeekerDocumentForTaskOutput struct {
	JobSeekerDocument *entity.JobSeekerDocument
}

func (i *JobSeekerInteractorImpl) UpdateJobSeekerDocumentForTask(input UpdateJobSeekerDocumentForTaskInput) (UpdateJobSeekerDocumentForTaskOutput, error) {
	var output UpdateJobSeekerDocumentForTaskOutput

	/************ 応募書類の更新 **************/

	jobSeekerDocument := entity.NewJobSeekerDocument(
		input.UpdateParam.JobSeekerID,
		input.UpdateParam.ResumeOriginURL,
		input.UpdateParam.ResumePDFURL,
		input.UpdateParam.CVOriginURL,
		input.UpdateParam.CVPDFURL,
		input.UpdateParam.RecommendationOriginURL,
		input.UpdateParam.RecommendationPDFURL,
		input.UpdateParam.IDPhotoURL,
		input.UpdateParam.OtherDocument1URL,
		input.UpdateParam.OtherDocument2URL,
		input.UpdateParam.OtherDocument3URL,
	)

	err := i.jobSeekerDocumentRepository.UpdateByJobSeekerID(input.UpdateParam.JobSeekerID, jobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ タスクを進める **************/

	latestTask, err := i.taskRepository.FindByJobSeekerUUIDAndJobInformationUUID(input.JobSeekerUUID, input.JobInformationUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 「書類選考/応募書類準備」
	IsPrepareDocument := latestTask.PhaseCategory == null.NewInt(int64(entity.DocumentSelection), true) && latestTask.PhaseSubCategory == null.NewInt(int64(entity.PrepareDocument), true)

	// 「書類選考/応募書類準備」の場合は「書類選考/エントリー依頼」
	if IsPrepareDocument {
		nextTask := entity.NewTask(
			latestTask.TaskGroupID,
			null.NewInt(int64(entity.DocumentSelection), true), // 書類選考
			null.NewInt(int64(entity.RequestEntry), true),      // エントリー依頼
			null.NewInt(int64(entity.CA), true),
			latestTask.ExecutedStaffID,
			"",
			latestTask.DeadlineDay,
			latestTask.DeadlineTime,
			"",
			"",
			"",
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(nextTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// CAからCAにタスク作成する場合
		err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// プッシュ通知に必要な変数定義
		var (
			redirectURL = os.Getenv("BASE_DOMAIN")
			contents    = "求職者がエントリーしました。"
			topic       = "TaskCA"
		)

		// 通知先のCA担当者の情報を取得
		caStaff, err := i.agentStaffRepository.FindByID(latestTask.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		err = utility.WebPush(
			i.oneSignal.AppID,
			i.oneSignal.APIKey,
			caStaff.FirebaseID,
			"タスク通知",
			contents,
			topic,
			redirectURL,
		)
		if err != nil {
			fmt.Println("WebPushの通知")
			fmt.Println(err)
		}
	}

	output.JobSeekerDocument = jobSeekerDocument

	return output, nil
}

type GetJobSeekerDocumentByJobSeekerIDInput struct {
	JobSeekerID uint
}

type GetJobSeekerDocumentByJobSeekerIDOutput struct {
	JobSeekerDocument *entity.JobSeekerDocument
}

func (i *JobSeekerInteractorImpl) GetJobSeekerDocumentByJobSeekerID(input GetJobSeekerDocumentByJobSeekerIDInput) (GetJobSeekerDocumentByJobSeekerIDOutput, error) {
	var output GetJobSeekerDocumentByJobSeekerIDOutput

	jobSeekerDocument, err := i.jobSeekerDocumentRepository.FindByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.JobSeekerDocument = jobSeekerDocument

	return output, nil
}

/****************************************************************************************/
// LINE関連API
//
// LINEIDを更新
type UpdateJobSeekerLineIDInput struct {
	Param entity.UpdateJobSeekerLineIDParam
}

type UpdateJobSeekerLineIDOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) UpdateJobSeekerLineID(input UpdateJobSeekerLineIDInput) (UpdateJobSeekerLineIDOutput, error) {
	var output UpdateJobSeekerLineIDOutput

	// AgentUUIDからLINE情報を取得
	agentLineChannel, err := i.agentRepository.FindLineChannelByAgentUUID(input.Param.AgentUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// アクセストークンを発行
	// docs:https://developers.line.biz/ja/reference/messaging-api/#issue-channel-access-token-v2-1
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Add("code", input.Param.Code)
	form.Add("redirect_uri", fmt.Sprintf("%s/init_job_seeker/complete/", os.Getenv("BASE_DOMAIN")))
	form.Add("client_id", agentLineChannel.LineLoginChannelID)
	form.Add("client_secret", agentLineChannel.LineLoginChannelSecret)

	body := strings.NewReader(form.Encode()) // リクエストのbodyを作成

	req, err := http.NewRequest(http.MethodPost, "https://api.line.me/oauth2/v2.1/token", body)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	defer res.Body.Close()

	// レスポンスを構造体に変換
	type lineAccessTokenJSON struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpresIn    int64  `json:"expires_in"`
		KeyID       string `json:"key_id"`
	}

	var tokenJSON lineAccessTokenJSON
	err = json.NewDecoder(res.Body).Decode(&tokenJSON)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 同じcodeで再度アクセストークンを取得すると400 Bad Requestが返ってくる
	if res.StatusCode == 400 && res.Status == "400 Bad Request" {
		return output, nil
	}

	if tokenJSON.AccessToken == "" {
		err = fmt.Errorf("%s:%w", "tokenJSON.AccessToken is empty", entity.ErrServerError)
		fmt.Println(err)
		return output, err
	}

	// AccessTokenを使って求職者のLINEIDを取得
	getProfileReq, err := http.NewRequest(http.MethodGet, "https://api.line.me/v2/profile", nil)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	getProfileReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenJSON.AccessToken))

	getProfileClient := &http.Client{}
	getProfileRes, err := getProfileClient.Do(getProfileReq)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	defer getProfileRes.Body.Close()

	type lineGetProfileJSON struct {
		UserID        string `json:"userId"`
		DisplayName   string `json:"displayName"`
		PictureURL    string `json:"pictureUrl"`
		StatusMessage string `json:"statusMessage"`
	}

	var profileJSON lineGetProfileJSON

	err = json.NewDecoder(getProfileRes.Body).Decode(&profileJSON)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if profileJSON.UserID == "" {
		err = fmt.Errorf("%s:%w", "profileJSON.UserID is empty", entity.ErrServerError)
		fmt.Println(err)
		return output, err
	}

	/**
	 * 1. 求職者テーブルに該当レコードがあるかを確認
	 * 2. 該当レコードがない場合は送客求職者テーブルに該当レコードがあるかを確認
	 * 3. どちらかのテーブルで該当のレコードがあれば処理成功
	 **/

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.Param.JobSeekerUUID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// 求職者テーブルに該当のuuidが存在しない場合
			sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByUUID(input.Param.JobSeekerUUID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// データベースにLINE_IDを保存する
			err = i.sendingJobSeekerRepository.UpdateLineID(input.Param.JobSeekerUUID, profileJSON.UserID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINEの利用状況を更新する（line_activeをtrueに変更）
			err = i.chatGroupWithSendingJobSeekerRepository.UpdateSendingJobSeekerLineActive(true, sendingJobSeeker.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

		} else {
			// 求職者テーブルでNotFound以外のエラーの場合はエラーを返す
			fmt.Println(err)
			return output, err
		}
	} else {
		// データベースにLINE_IDを保存する
		err := i.jobSeekerRepository.UpdateLineIDByUUID(input.Param.JobSeekerUUID, profileJSON.UserID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// LINEの利用状況を更新する（line_activeをtrueに変更）
		err = i.chatGroupWithJobSeekerRepository.UpdateJobSeekerLineActive(jobSeeker.ID, true)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
// 面談前アンケート関連 API
//

// 面談前アンケートを登録
// 1. 求職者の同意項目アップデート
// 2. アンケート情報の登録（業界、職種、勤務地、質問要望）
// 3. 求職者情報の更新（業界、職種、勤務地）
// 4. 求職者情報の更新（ファイル（履歴書（原本）、商務経歴書（原本）））
type CreateInitialQuestionnaireInput struct {
	Param entity.CreateInitialQuestionnaireParam
}

type CreateInitialQuestionnairdhutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) CreateInitialQuestionnaire(input CreateInitialQuestionnaireInput) (CreateInitialQuestionnairdhutput, error) {
	var (
		output CreateInitialQuestionnairdhutput
		err    error
	)

	// データベースに保存
	initialQuestionnaire := entity.NewInitialQuestionnaire(
		input.Param.JobSeekerID,
		input.Param.Question,
	)

	// 面談前アンケートを登録
	err = i.initialQuestionnaireRepository.Create(initialQuestionnaire)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者テーブルの個人情報同意の有無を更新
	err = i.jobSeekerRepository.UpdateAgreement(
		input.Param.JobSeekerID,
		true,
		// input.Param.Agreement,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 履歴書(原本)がアップロードされた場合は更新
	if input.Param.Documents.ResumeOriginURL != "" {
		err = i.jobSeekerDocumentRepository.UpdateResumeOriginURLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.ResumeOriginURL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 職務経歴書(原本)がアップロードされた場合は更新
	if input.Param.Documents.CVOriginURL != "" {
		err = i.jobSeekerDocumentRepository.UpdateCVOriginURLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.CVOriginURL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 履歴書(PDF)がアップロードされた場合は更新
	if input.Param.Documents.ResumePDFURL != "" {
		err = i.jobSeekerDocumentRepository.UpdateResumePDFURLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.ResumePDFURL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 職務経歴書(PDF)がアップロードされた場合は更新
	if input.Param.Documents.CVPDFURL != "" {
		err = i.jobSeekerDocumentRepository.UpdateCVPDFURLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.CVPDFURL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 証明写真がアップロードされた場合は更新
	if input.Param.Documents.IDPhotoURL != "" {
		err = i.jobSeekerDocumentRepository.UpdateIDPhotoURLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.IDPhotoURL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// その他①がアップロードされた場合は更新
	if input.Param.Documents.OtherDocument1URL != "" {
		err = i.jobSeekerDocumentRepository.UpdateOtherDocument1URLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.OtherDocument1URL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// その他②がアップロードされた場合は更新
	if input.Param.Documents.OtherDocument2URL != "" {
		err = i.jobSeekerDocumentRepository.UpdateOtherDocument2URLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.OtherDocument2URL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// その③がアップロードされた場合は更新
	if input.Param.Documents.OtherDocument3URL != "" {
		err = i.jobSeekerDocumentRepository.UpdateOtherDocument3URLByJobSeekerID(
			input.Param.JobSeekerID,
			input.Param.Documents.OtherDocument3URL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 求職者が存在しない場合は希望条件をそのまま作成
	if len(input.Param.DesiredIndustries) > 0 {
		// 求職者が存在する場合は求職者テーブルの希望条件は更新
		// 求職者IDから希望業界を削除
		err = i.jobSeekerDesiredIndustryRepository.DeleteByJobSeekerID(input.Param.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 希望業界
		for _, desiredIndustry := range input.Param.DesiredIndustries {
			// アンケートテーブルに保存
			err = i.initialQuestionnaireDesiredIndustryRepository.Create(
				entity.NewInitialQuestionnaireDesiredIndustry(
					initialQuestionnaire.ID,
					desiredIndustry.DesiredIndustry,
					desiredIndustry.DesiredRank,
				),
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// 求職者テーブルに保存
			err = i.jobSeekerDesiredIndustryRepository.Create(
				entity.NewJobSeekerDesiredIndustry(
					initialQuestionnaire.JobSeekerID,
					desiredIndustry.DesiredIndustry,
					desiredIndustry.DesiredRank,
				),
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// 求職者が存在しない場合は希望条件をそのまま作成
	if len(input.Param.DesiredOccupations) > 0 {
		// 求職者が存在する場合は求職者テーブルの希望条件は更新
		// 求職者IDから希望職種を削除
		err = i.jobSeekerDesiredOccupationRepository.DeleteByJobSeekerID(input.Param.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 希望職種
		for _, desiredOccupation := range input.Param.DesiredOccupations {
			// アンケートテーブルに保存
			err = i.initialQuestionnaireDesiredOccupationRepository.Create(
				entity.NewInitialQuestionnaireDesiredOccupation(
					initialQuestionnaire.ID,
					desiredOccupation.DesiredOccupation,
					desiredOccupation.DesiredRank,
				),
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// 求職者テーブルに保存
			err = i.jobSeekerDesiredOccupationRepository.Create(
				entity.NewJobSeekerDesiredOccupation(
					initialQuestionnaire.JobSeekerID,
					desiredOccupation.DesiredOccupation,
					desiredOccupation.DesiredRank,
				),
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// 求職者が存在しない場合は希望条件をそのまま作成
	if len(input.Param.DesiredWorkLocations) > 0 {
		// 求職者が存在する場合は求職者テーブルの希望条件は更新
		// 求職者IDから希望勤務地を削除
		err = i.jobSeekerDesiredWorkLocationRepository.DeleteByJobSeekerID(input.Param.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 希望勤務地
		for _, desiredWorkLocation := range input.Param.DesiredWorkLocations {
			// アンケートテーブルに保存
			err = i.initialQuestionnaireDesiredWorkLocationRepository.Create(
				entity.NewInitialQuestionnaireDesiredWorkLocation(
					initialQuestionnaire.ID,
					desiredWorkLocation.DesiredWorkLocation,
					desiredWorkLocation.DesiredRank,
				),
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// 求職者テーブルに保存
			err = i.jobSeekerDesiredWorkLocationRepository.Create(
				entity.NewJobSeekerDesiredWorkLocation(
					initialQuestionnaire.JobSeekerID,
					desiredWorkLocation.DesiredWorkLocation,
					desiredWorkLocation.DesiredRank,
				),
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
/// Admin API
//

// すべての求職者情報を取得
type GetAllJobSeekerInput struct {
	PageNumber uint
}

type GetAllJobSeekerOutput struct {
	JobSeekerList []*entity.JobSeeker
	MaxPageNumber uint
	IDList        []uint
}

func (i *JobSeekerInteractorImpl) GetAllJobSeeker(input GetAllJobSeekerInput) (GetAllJobSeekerOutput, error) {
	var (
		output GetAllJobSeekerOutput
		err    error
	)

	jobSeekerList, err := i.jobSeekerRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	studentHistory, err := i.jobSeekerStudentHistoryRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	workHistory, err := i.jobSeekerWorkHistoryRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceIndustry, err := i.jobSeekerExperienceIndustryRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory, err := i.jobSeekerDepartmentHistoryRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	experienceOccupation, err := i.jobSeekerExperienceOccupationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredCompanyScale, err := i.jobSeekerDesiredCompanyScaleRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	license, err := i.jobSeekerLicenseRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selfPromotion, err := i.jobSeekerSelfPromotionRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	document, err := i.jobSeekerDocumentRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredIndustry, err := i.jobSeekerDesiredIndustryRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupation, err := i.jobSeekerDesiredOccupationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocation, err := i.jobSeekerDesiredWorkLocationRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredHolidayType, err := i.jobSeekerDesiredHolidayTypeRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	developmentSkill, err := i.jobSeekerDevelopmentSkillRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	languageSkill, err := i.jobSeekerLanguageSkillRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	pcSkill, err := i.jobSeekerPCToolRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgent, err := i.jobSeekerHideToAgentRepository.All()
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
									ID:                  eo.ID,
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
					JobSeekerID:       d.JobSeekerID,
					ResumeOriginURL:   d.ResumeOriginURL,
					ResumePDFURL:      d.ResumePDFURL,
					CVOriginURL:       d.CVOriginURL,
					CVPDFURL:          d.CVPDFURL,
					IDPhotoURL:        d.IDPhotoURL,
					OtherDocument1URL: d.OtherDocument1URL,
					OtherDocument2URL: d.OtherDocument2URL,
					OtherDocument3URL: d.OtherDocument3URL,
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
			value := entity.JobSeekerHideToAgent{
				JobSeekerID: hta.JobSeekerID,
				AgentID:     hta.AgentID,
				AgentName:   hta.AgentName,
			}

			jobSeeker.HideToAgents = append(jobSeeker.HideToAgents, value)
		}
	}

	// IDListを返す
	for _, jobSeeker := range jobSeekerList {
		output.IDList = append(output.IDList, jobSeeker.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getJobSeekerListMaxPage(jobSeekerList)

	// 指定ページの求職者20件を取得（本番実装までは1ページあたり5件）
	output.JobSeekerList = getJobSeekerListWithPage(jobSeekerList, input.PageNumber)

	return output, nil
}

/****************************************************************************************/

func setJobSeekerChildTableByIDList(i *JobSeekerInteractorImpl, jobSeekerList []*entity.JobSeeker) ([]*entity.JobSeeker, error) {
	var output []*entity.JobSeeker

	idList := getJobSeekerIDList(jobSeekerList)

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

	selfPromotion, err := i.jobSeekerSelfPromotionRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	document, err := i.jobSeekerDocumentRepository.GetByJobSeekerIDList(idList)
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

	pcSkill, err := i.jobSeekerPCToolRepository.GetByJobSeekerIDList(idList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hideToAgent, err := i.jobSeekerHideToAgentRepository.GetByJobSeekerIDList(idList)
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
									ID:                  eo.ID,
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
					JobSeekerID:       d.JobSeekerID,
					ResumeOriginURL:   d.ResumeOriginURL,
					ResumePDFURL:      d.ResumePDFURL,
					CVOriginURL:       d.CVOriginURL,
					CVPDFURL:          d.CVPDFURL,
					IDPhotoURL:        d.IDPhotoURL,
					OtherDocument1URL: d.OtherDocument1URL,
					OtherDocument2URL: d.OtherDocument2URL,
					OtherDocument3URL: d.OtherDocument3URL,
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

	return jobSeekerList, nil
}

type DeleteJobSeekerResumePDFURLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerResumePDFURLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerResumePDFURL(input DeleteJobSeekerResumePDFURLInput) (DeleteJobSeekerResumePDFURLOutput, error) {
	var output DeleteJobSeekerResumePDFURLOutput

	err := i.jobSeekerDocumentRepository.UpdateResumePDFURLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerResumeOriginURLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerResumeOriginURLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerResumeOriginURL(input DeleteJobSeekerResumeOriginURLInput) (DeleteJobSeekerResumeOriginURLOutput, error) {
	var output DeleteJobSeekerResumeOriginURLOutput

	err := i.jobSeekerDocumentRepository.UpdateResumeOriginURLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerCVPDFURLInput struct {
	JobSeekerID uint
}
type DeleteJobSeekerCVPDFURLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerCVPDFURL(input DeleteJobSeekerCVPDFURLInput) (DeleteJobSeekerCVPDFURLOutput, error) {
	var output DeleteJobSeekerCVPDFURLOutput

	err := i.jobSeekerDocumentRepository.UpdateCVPDFURLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerCVOriginURLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerCVOriginURLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerCVOriginURL(input DeleteJobSeekerCVOriginURLInput) (DeleteJobSeekerCVOriginURLOutput, error) {
	var output DeleteJobSeekerCVOriginURLOutput

	err := i.jobSeekerDocumentRepository.UpdateCVOriginURLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerRecommendationPDFURLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerRecommendationPDFURLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerRecommendationPDFURL(input DeleteJobSeekerRecommendationPDFURLInput) (DeleteJobSeekerRecommendationPDFURLOutput, error) {
	var output DeleteJobSeekerRecommendationPDFURLOutput

	err := i.jobSeekerDocumentRepository.UpdateRecommendationPDFURLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerRecommendationOriginURLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerRecommendationOriginURLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerRecommendationOriginURL(input DeleteJobSeekerRecommendationOriginURLInput) (DeleteJobSeekerRecommendationOriginURLOutput, error) {
	var output DeleteJobSeekerRecommendationOriginURLOutput

	err := i.jobSeekerDocumentRepository.UpdateRecommendationOriginURLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerIDPhotoURLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerIDPhotoURLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerIDPhotoURL(input DeleteJobSeekerIDPhotoURLInput) (DeleteJobSeekerIDPhotoURLOutput, error) {
	var output DeleteJobSeekerIDPhotoURLOutput

	err := i.jobSeekerDocumentRepository.UpdateIDPhotoURLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerOtherDocument1URLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerOtherDocument1URLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerOtherDocument1URL(input DeleteJobSeekerOtherDocument1URLInput) (DeleteJobSeekerOtherDocument1URLOutput, error) {
	var output DeleteJobSeekerOtherDocument1URLOutput

	err := i.jobSeekerDocumentRepository.UpdateOtherDocument1URLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerOtherDocument2URLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerOtherDocument2URLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerOtherDocument2URL(input DeleteJobSeekerOtherDocument2URLInput) (DeleteJobSeekerOtherDocument2URLOutput, error) {
	var output DeleteJobSeekerOtherDocument2URLOutput

	err := i.jobSeekerDocumentRepository.UpdateOtherDocument2URLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteJobSeekerOtherDocument3URLInput struct {
	JobSeekerID uint
}

type DeleteJobSeekerOtherDocument3URLOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) DeleteJobSeekerOtherDocument3URL(input DeleteJobSeekerOtherDocument3URLInput) (DeleteJobSeekerOtherDocument3URLOutput, error) {
	var output DeleteJobSeekerOtherDocument3URLOutput

	err := i.jobSeekerDocumentRepository.UpdateOtherDocument3URLByJobSeekerID(input.JobSeekerID, "")
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
// ゲストページ用 API
//
// 求職者IDを使って求職者情報を取得する
type GetJobSeekerForInitialStepByUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetJobSeekerForInitialStepByUUIDOutput struct {
	JobSeekerForGuest *entity.JobSeekerForGuest
}

func (i *JobSeekerInteractorImpl) GetJobSeekerForInitialStepByUUID(input GetJobSeekerForInitialStepByUUIDInput) (GetJobSeekerForInitialStepByUUIDOutput, error) {
	var (
		output GetJobSeekerForInitialStepByUUIDOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(jobSeeker.AgentID, jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	isBlocked := false
	if !chatGroup.LineActive && jobSeeker.LineID != "" {
		isBlocked = true
	}

	output.JobSeekerForGuest = &entity.JobSeekerForGuest{
		ID:         jobSeeker.ID,
		Agreement:  jobSeeker.Agreement,
		LineActive: chatGroup.LineActive,
		IsBlocked:  isBlocked,
	}

	return output, nil
}

type GetGuestJobSeekerForByUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetGuestJobSeekerForByUUIDOutput struct {
	User *entity.GuestJobSeekerUser
}

func (i *JobSeekerInteractorImpl) GetGuestJobSeekerForByUUID(input GetGuestJobSeekerForByUUIDInput) (GetGuestJobSeekerForByUUIDOutput, error) {
	var (
		output GetGuestJobSeekerForByUUIDOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントアカウントのログインユーザー情報を作成
	guestJobSeeker := entity.NewGuestJobSeekerUser(
		jobSeeker.ID,
		input.JobSeekerUUID,
		jobSeeker.LastName,
		jobSeeker.FirstName,
		jobSeeker.Email,
		jobSeeker.AgentID,
		jobSeeker.Phase,
		jobSeeker.CanViewMatchingJob,
	)

	output.User = guestJobSeeker

	return output, nil
}

type GetJobSeekerDesiredForGuestByUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetJobSeekerDesiredForGuestByUUIDOutput struct {
	JobSeekerDesired *entity.JobSeekerDesiredForGuest
}

func (i *JobSeekerInteractorImpl) GetJobSeekerDesiredForGuestByUUID(input GetJobSeekerDesiredForGuestByUUIDInput) (GetJobSeekerDesiredForGuestByUUIDOutput, error) {
	var (
		output GetJobSeekerDesiredForGuestByUUIDOutput
		err    error

		desiredIndustryListNullInt     []null.Int
		desiredOccupationListNullInt   []null.Int
		desiredWorkLocationListNullInt []null.Int
	)

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

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

	return output, nil
}

// 求職者IDを使って求職者情報を取得する
type GetJobSeekerAgentIDByUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetJobSeekerAgentIDByUUIDOutput struct {
	AgentID uint
}

func (i *JobSeekerInteractorImpl) GetJobSeekerAgentIDByUUID(input GetJobSeekerAgentIDByUUIDInput) (GetJobSeekerAgentIDByUUIDOutput, error) {
	var (
		output GetJobSeekerAgentIDByUUIDOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentID = jobSeeker.AgentID

	return output, nil
}

type CheckJobSeekerByUUIDAndNameInput struct {
	Param entity.CheckJobSeekerByUUIDAndNameParam
}

type CheckJobSeekerByUUIDAndNameOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) CheckJobSeekerByUUIDAndName(input CheckJobSeekerByUUIDAndNameInput) (CheckJobSeekerByUUIDAndNameOutput, error) {
	var (
		output CheckJobSeekerByUUIDAndNameOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.Param.UUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if input.Param.Name == fmt.Sprintf("%s%s", jobSeeker.LastName, jobSeeker.FirstName) {
		output.OK = true
	}

	return output, nil
}

type UpdateJobSeekerPasswordInput struct {
	Param entity.UpdateJobSeekerPasswordParam
}

type UpdateJobSeekerPasswordOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) UpdateJobSeekerPassword(input UpdateJobSeekerPasswordInput) (UpdateJobSeekerPasswordOutput, error) {
	output := UpdateJobSeekerPasswordOutput{}

	// ResetPasswordTokenが一致するか確認
	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.Param.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if input.Param.ResetPasswordToken != jobSeeker.ResetPasswordToken {
		fmt.Println("ResetPasswordTokenが一致しません")
		return output, nil
	}

	// パスワードのハッシュ化
	hashedPassword, err := hashPassword(input.Param.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// パスワードの更新
	err = i.jobSeekerRepository.UpdatePassword(input.Param.JobSeekerUUID, hashedPassword)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// NOTE: ResetPasswordTokenを空で更新するのはトークンの有無で表示を切り替えるため
	// パスワードトークンを削除（空で更新する）
	err = i.jobSeekerRepository.UpdateResetPasswordToken(jobSeeker.UUID, "")
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type SendJobSeekerResetPasswordEmailInput struct {
	Param entity.SendJobSeekerResetPasswordEmailParam
}

type SendJobSeekerResetPasswordEmailOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) SendJobSeekerResetPasswordEmail(input SendJobSeekerResetPasswordEmailInput) (SendJobSeekerResetPasswordEmailOutput, error) {
	var (
		output SendJobSeekerResetPasswordEmailOutput
		err    error
	)

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.Param.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// メールアドレスが一致するか確認
	if input.Param.Email != jobSeeker.Email {
		err = errors.New("メールアドレスが一致しません")
		return output, err
	}

	// パスワードリセットトークンを生成
	resetPasswordToken := utility.CreateUUID()

	err = i.jobSeekerRepository.UpdateResetPasswordToken(jobSeeker.UUID, resetPasswordToken.String())
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	baseURL := os.Getenv("BASE_DOMAIN")
	resetPasswordURL := fmt.Sprintf("%s/guest_js/reset_password/?job_seeker=%s&reset_password_token=%s", baseURL, jobSeeker.UUID, resetPasswordToken)

	// メール送信
	err = utility.SendMailToSingleWithoutCC(
		i.sendgrid.APIKey,
		"パスワードリセットのお手続き",
		fmt.Sprintf("お客様\n\n %s アカウントのパスワードをリセットするには、次のリンクをクリックしてください。\n\n%s\n\nパスワードのリセットを依頼していない場合は、このメールを無視してください。\nよろしくお願いいたします。\n\nautoscout事務局",
			jobSeeker.Email,
			resetPasswordURL,
		),
		entity.EmailUser{
			Name:  "autoscout事務局",
			Email: "info@spaceai.jp",
		},
		entity.EmailUser{
			Name:  jobSeeker.LastName + jobSeeker.FirstName + "様",
			Email: jobSeeker.Email,
		},
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	output.OK = true

	return output, nil
}

type SendJobSeekerContactInput struct {
	Param entity.SendJobSeekerContactParam
}

type SendJobSeekerContactOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) SendJobSeekerContact(input SendJobSeekerContactInput) (SendJobSeekerContactOutput, error) {
	var (
		output SendJobSeekerContactOutput
		err    error
	)

	// メール送信
	content := fmt.Sprintf("*求職者マイページからのお問い合わせ*\n\n%s", input.Param.Content)

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// 参考: https://risaki-masa.com/how-to-get-api-token-in-slack/
	slack := utility.NewSlack(i.slack.ReachAccessToken)
	err = slack.SendContact(i.slack.ReachChanelID, content)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	output.OK = true

	return output, nil
}

// 送客求職者の面談日時の更新
type UpdateInterviewDateByJobSeekerIDInput struct {
	Param entity.UpdateJobSeekerInterviewDateFromGuestPageParam
}

type UpdateInterviewDateByJobSeekerIDOutput struct {
	User *entity.GuestJobSeekerUser
}

func (i *JobSeekerInteractorImpl) UpdateInterviewDateByJobSeekerID(input UpdateInterviewDateByJobSeekerIDInput) (UpdateInterviewDateByJobSeekerIDOutput, error) {
	var (
		output               UpdateInterviewDateByJobSeekerIDOutput
		param                = input.Param
		err                  error
		staffID              uint
		reservationInterview = null.NewInt(int64(entity.ReservationInterview), true) // 面談予約完了
	)

	jobSeeker, err := i.jobSeekerRepository.FindByID(param.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 面談日時を更新
	err = i.jobSeekerRepository.UpdateInterviewDate(param.JobSeekerID, param.InterviewDate)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// フェーズを更新する
	err = i.jobSeekerRepository.UpdatePhase(param.JobSeekerID, reservationInterview)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 当日（2020-02-01）
	now := time.Now()
	date := now.Format("2006-01-02")

	// 面談調整タスクの作成
	interviewTaskGroup, err := i.interviewTaskGroupRepository.FindByAgentIDAndJobSeekerID(jobSeeker.AgentID, param.JobSeekerID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			interviewTaskGroup = entity.NewInterviewTaskGroup(
				jobSeeker.AgentID,
				param.JobSeekerID,
				param.InterviewDate,
				param.InterviewDate, // 初期値
			)

			err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			interviewTask := entity.NewInterviewTask(
				interviewTaskGroup.ID,
				null.NewInt(int64(staffID), true),
				null.NewInt(int64(staffID), true),
				reservationInterview,
				null.NewInt(0, true),
				"",
				date,
				null.NewInt(99, true),
				getStrPhaseForJobSeeker(reservationInterview), // SelectActionLabelは求職者のフェーズにする
			)

			err = i.interviewTaskRepository.Create(interviewTask)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			fmt.Println(err)
			return output, err
		}
	} else {
		err = i.interviewTaskGroupRepository.UpdateInterviewDate(interviewTaskGroup.ID, param.InterviewDate)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最新の面談タスクを取得
		latestInterviewTask, err := i.interviewTaskRepository.FindLatestByAgentIDAndJobSeekerID(jobSeeker.AgentID, param.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 最新のタスクと編集後のタスクが違う場合はタスクを作成するする
		if latestInterviewTask.PhaseCategory != reservationInterview {
			interviewTask := entity.NewInterviewTask(
				interviewTaskGroup.ID,
				null.NewInt(int64(staffID), true),
				null.NewInt(int64(staffID), true),
				reservationInterview,
				null.NewInt(0, true),
				"",
				date,
				null.NewInt(99, true),
				getStrPhaseForJobSeeker(reservationInterview), // SelectActionLabelは求職者のフェーズにする
			)

			err = i.interviewTaskRepository.Create(interviewTask)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// すでにタスクグループが存在する場合は依頼時間を更新
			err = i.interviewTaskGroupRepository.UpdateLastRequestAt(interviewTaskGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			fmt.Println("同一タスクです")
			fmt.Println("---------")
		}
	}

	// Slackに送付するメッセージ
	content := fmt.Sprintf(
		"*マイページから面談予約が行われました。*\n\nID: %v\n\nお名前: %s%s\n\nメールアドレス: %s\n\n面談日時: %s\n\n担当: %s",
		jobSeeker.ID,
		jobSeeker.LastName,
		jobSeeker.FirstName,
		jobSeeker.Email,
		param.InterviewDate.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05"),
		param.StaffName,
	)

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// NOTE: https://risaki-masa.com/how-to-get-api-token-in-slack/
	slack := utility.NewSlack(i.slack.ReachAccessToken)
	err = slack.SendContact(i.slack.ReachChanelID, content)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントアカウントのログインユーザー情報を作成
	guestJobSeeker := entity.NewGuestJobSeekerUser(
		jobSeeker.ID,
		jobSeeker.UUID,
		jobSeeker.LastName,
		jobSeeker.FirstName,
		jobSeeker.Email,
		jobSeeker.AgentID,
		reservationInterview, // 面談予約完了
		jobSeeker.CanViewMatchingJob,
	)

	output.User = guestJobSeeker

	return output, nil
}

/****************************************************************************************/
// LP用 API
//
type CreateJobSeekerFromLPInput struct {
	Param entity.CreateJobSeekerFromLPParam
}

type CreateJobSeekerFromLPOutput struct {
	UUID uuid.UUID
}

func (i *JobSeekerInteractorImpl) CreateJobSeekerFromLP(input CreateJobSeekerFromLPInput) (CreateJobSeekerFromLPOutput, error) {
	var (
		output            CreateJobSeekerFromLPOutput
		param             = input.Param
		stateOfEmployment null.Int
		retireYear        string
		retireLastStatus  null.Int
		nationality       null.Int
		systemAgentID     uint = 1
		defaultCAStaffID  uint = 2 // 本番環境のIDをデフォルトCAとして設定
	)

	// すでに登録済みのアドレスかをチェック
	sameEmailCount, err := i.jobSeekerRepository.CountByEmail(param.Email)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	fmt.Println("param", param)

	if sameEmailCount > 0 {
		err = fmt.Errorf("%w:%s", entity.ErrDuplicateEntry, "すでに使用されているメールアドレスです")
		return output, err
	}

	// 離職有無から就業状況の値をセット
	if param.IsRetire {
		// 離職中
		stateOfEmployment = null.NewInt(1, true)
		retireYear = fmt.Sprintf("%s-%s", param.RetireYear, param.RetireMonth)
		retireLastStatus = null.NewInt(1, true) // 一身上の都合により退職
	} else {
		// 現職中
		stateOfEmployment = null.NewInt(0, true)
		retireLastStatus = null.NewInt(0, true) // 現在に至る
	}

	// FirstLanguage{0: 日本語, 1: 英語, 2: その他の言語}
	if param.FirstLanguage == null.NewInt(0, true) {
		// 日本国籍
		nationality = null.NewInt(0, true)
	} else {
		// 外国籍
		nationality = null.NewInt(1, true)
	}

	env := os.Getenv("APP_ENV")
	if env != "prd" {
		// 本番以外の場合はIDを1
		defaultCAStaffID = 1
	}

	jobSeeker := entity.NewJobSeeker(
		systemAgentID, // Systemで登録
		null.NewInt(int64(defaultCAStaffID), true),
		NullInt,
		param.LastName,
		param.FirstName,
		param.LastFurigana,
		param.FirstFurigana,
		param.Gender,
		"",
		fmt.Sprintf("%s-%s-%s", param.Birthyear, param.Birthmonth, param.Birthday), // 生年月日（1999-06-10）
		NullInt,                // 配偶者
		NullInt,                // 配偶者扶養義務
		NullInt,                // 扶養家族人数
		"",                     // 電話番号
		param.Email,            // メールアドレス
		"",                     // 緊急連絡先
		"",                     // 郵便番号
		param.Prefecture,       // 都道府県
		"",                     // 住所
		"",                     // 住所フリガナ
		stateOfEmployment,      // 就業状況
		param.PRPoint,          // 職務要約（JobSummary）
		param.JobSummary,       // 経歴補足（HistorySupplement）
		"",                     // 研究内容・学チカ
		NullInt,                // 入社可能時期
		param.CompanyNum,       // 転職回数（JobChange）
		param.Income,           // 直近の年収（AnnualIncome）
		NullInt,                // 希望年収
		NullInt,                // 転職可否
		"",                     // 転勤条件
		NullInt,                // 短期離職
		"",                     // 短期離職補足
		NullInt,                // 既往歴
		nationality,            // 国籍
		NullInt,                // アピアランス
		NullInt,                // コミュ力
		NullInt,                // 論理的思考力
		"",                     // 人物像（推薦状用）
		"",                     // 人物像（本音）
		"",                     // メモ
		NullInt,                // 転職・就活状況（JobHuntingState）
		"",                     // 推薦理由（RecommendReason）
		null.NewInt(0, true),   // フェーズ（エントリーで登録）
		utility.EarliestTime(), // 面談日時
		null.NewInt(1, true),   // 登録状況（仮登録で登録）
		NullInt,                // 専攻学科
		NullInt,                // ワードスキル
		NullInt,                // エクセルスキル
		NullInt,                // パワーポイントスキル
		NullInt,                // 流入経路
		"",                     // 国籍備考
		"",                     // 既往歴備考
		"",                     // 応募承諾のポイント
	)

	err = i.jobSeekerRepository.Create(jobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 応募可能求人の閲覧権限をTRUEで更新
	err = i.jobSeekerRepository.UpdateCanViewMatchingJob(jobSeeker.ID, true)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// パスワードのハッシュ化
	hashedPassword, err := hashPassword(param.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// パスワードの更新
	err = i.jobSeekerRepository.UpdatePassword(jobSeeker.UUID, hashedPassword)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 学歴
	studentHistory := entity.NewJobSeekerStudentHistory(
		jobSeeker.ID,
		param.SchoolCategory,
		param.SchoolName,
		NullInt, // 学校レベル
		param.Subject,
		"",      // 入学年月
		NullInt, // 入学ステータス
		fmt.Sprintf("%s-%s", param.GraduationYear, param.GraduationMonth), // 卒業年月
		null.NewInt(0, true), // 卒業ステータス（卒業）
	)

	err = i.jobSeekerStudentHistoryRepository.Create(studentHistory)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 職歴
	workHistory := entity.NewJobSeekerWorkHistory(
		jobSeeker.ID,
		param.CompanyName,
		getStrEmployeeNumber(param.EmployeeNumber), // 従業員数
		NullInt, // 従業員数
		NullInt, // 株式公開
		fmt.Sprintf("%s-%s", param.JoiningYear, param.JoiningMonth), // 入社年月
		param.EmployeeStatus, // 雇用形態
		"",                   // 退職理由（本音）
		"",                   // 退職理由（建前）
		retireYear,           // 退職年月
		null.NewInt(0, true), // 入社ステータス（入社）
		retireLastStatus,     // 退職ステータス
	)

	err = i.jobSeekerWorkHistoryRepository.Create(workHistory)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 経験業界
	for _, industry := range param.Industries {
		experienceIndustry := entity.NewJobSeekerExperienceIndustry(
			workHistory.ID,
			industry.Industry,
		)

		err = i.jobSeekerExperienceIndustryRepository.Create(experienceIndustry)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 経験職種

	// 職務内容に職種と経験年数をセットする
	var jobDescriptionArray []string
	for _, experienceOccupation := range param.ExperienceOccupations {
		occupationStr := getStrOccupation(experienceOccupation.Occupation)
		experienceYearStr := getStrExperienceYear(experienceOccupation.ExperienceYear)

		jobDescriptionArray = append(jobDescriptionArray, fmt.Sprintf("%s %s", occupationStr, experienceYearStr))
	}

	departmentHistory := entity.NewJobSeekerDepartmentHistory(
		workHistory.ID,
		"",                                      // 部署
		NullInt,                                 // マネジメント人数
		"",                                      // マネジメントの詳細
		strings.Join(jobDescriptionArray, "\n"), // 職務内容
		fmt.Sprintf("%s-%s", param.JoiningYear, param.JoiningMonth), // 開始年月
		retireYear, // 終了年月
	)

	err = i.jobSeekerDepartmentHistoryRepository.Create(departmentHistory)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, eo := range param.ExperienceOccupations {
		occupation := entity.NewJobSeekerExperienceOccupation(
			departmentHistory.ID,
			eo.Occupation,
		)

		err = i.jobSeekerExperienceOccupationRepository.Create(occupation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, aeo := range param.AllExperienceOccupations {
		allExperienceOccupation := entity.NewJobSeekerExperienceJob(
			jobSeeker.ID,
			aeo.Occupation,
			aeo.ExperienceYear,
		)

		err = i.jobSeekerExperienceJobRepository.Create(allExperienceOccupation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 資格
	if param.DriversLicense == null.NewInt(0, true) {
		// 「0: 持っている（AT車限定）」の場合
		// 4805: 普通自動車免許（AT）
		param.Licenses = append(param.Licenses, null.NewInt(4805, true))
	} else if param.DriversLicense == null.NewInt(1, true) {
		// 「1: 持っている（MT車）」の場合
		// 4803: 普通自動車免許（MT）
		param.Licenses = append(param.Licenses, null.NewInt(4803, true))
	}

	// 資格の重複チェック（資格選択で自動車免許を選択することを考慮して重複チェック）
	duplicateLicense := make(map[null.Int]bool)

	for _, l := range param.Licenses {
		if !duplicateLicense[l] {
			license := entity.NewJobSeekerLicense(
				jobSeeker.ID,
				l,
				"",
			)

			err = i.jobSeekerLicenseRepository.Create(license)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			duplicateLicense[l] = true
		}
	}

	for _, language := range param.Languages {
		languageSkill := entity.NewJobSeekerLanguageSkill(
			jobSeeker.ID,
			language.LanguageType,
			language.LanguageLevel,
			NullInt,
			"",
			NullInt,
			"",
			NullInt,
			"",
		)

		err = i.jobSeekerLanguageSkillRepository.Create(languageSkill)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// チャットグループ
	chatGroup := entity.NewChatGroupWithJobSeeker(
		jobSeeker.AgentID, // System
		jobSeeker.ID,
		false, // 初めはLINE連携してないから false
	)

	err = i.chatGroupWithJobSeekerRepository.Create(chatGroup)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 面談調整タスク
	phaseSub := null.NewInt(0, true) // 日程調整依頼
	// 当日（2020-02-01）
	date := jobSeeker.InterviewDate.Format("2006-01-02")

	interviewTaskGroup := entity.NewInterviewTaskGroup(
		jobSeeker.AgentID,
		jobSeeker.ID,
		jobSeeker.InterviewDate,
		jobSeeker.InterviewDate, // 面談実施済みの場合は初回面談日時に記録する
	)

	err = i.interviewTaskGroupRepository.Create(interviewTaskGroup)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	interviewTask := entity.NewInterviewTask(
		interviewTaskGroup.ID,
		NullInt, //　タスク実行者ID
		NullInt, // CAのID
		jobSeeker.Phase,
		phaseSub,
		"",
		date,
		null.NewInt(99, true),
		getStrPhaseForJobSeeker(jobSeeker.Phase), // SelectActionLabelは求職者のフェーズにする
	)

	err = i.interviewTaskRepository.Create(interviewTask)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 書類テーブルの作成
	document := entity.NewJobSeekerDocument(
		jobSeeker.ID, "", "", "", "", "", "", "", "", "", "",
	)

	err = i.jobSeekerDocumentRepository.Create(document)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.UUID = jobSeeker.UUID

	return output, nil
}

// LPから求職者の電話番号を更新
type UpdateJobSeekerPhoneFromLPInput struct {
	Param entity.UpdateJobSeekerPhoneFromLPParam
}

type UpdateJobSeekerPhoneFromLPOutput struct {
	UUID        uuid.UUID
	LogintToken uuid.UUID
}

func (i *JobSeekerInteractorImpl) UpdateJobSeekerPhoneFromLP(input UpdateJobSeekerPhoneFromLPInput) (UpdateJobSeekerPhoneFromLPOutput, error) {
	var (
		output        UpdateJobSeekerPhoneFromLPOutput
		param              = input.Param
		SystemAgentID uint = 1
	)

	// uuidで求職者情報を取得
	jobSeeker, err := i.jobSeekerRepository.FindByUUID(param.UUID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "再読み込みして再度お試しください。")
			return output, wrapped
		} else {
			// Not Found以外のエラーの場合はそのままサーバーエラー
			fmt.Println(err)
			return output, err
		}
	}

	// 桁数が10もしくは11かをチェック
	if len(param.PhoneNumber) < 10 || 11 < len(param.PhoneNumber) {
		wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "電話番号の形式が正しくありません。")
		return output, wrapped
	}

	// 全ての桁が同一かをチェック
	var (
		firstChar = param.PhoneNumber[0]
		isAllSame = true
	)

	for i := 1; i < len(param.PhoneNumber); i++ {
		if param.PhoneNumber[i] != firstChar {
			isAllSame = false
		}
	}

	if isAllSame {
		wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "電話番号の形式が正しくありません。")
		return output, wrapped
	}

	// 電話番号が登録済みの場合はエラー
	_, err = i.jobSeekerRepository.FindByPhoneNumberForLP(param.PhoneNumber, SystemAgentID)
	if err == nil {
		wrapped := fmt.Errorf("%w:%s", entity.ErrDuplicateEntry, "既に登録済みの電話番号です。")
		return output, wrapped
	}

	// 電話番号の更新
	err = i.jobSeekerRepository.UpdatePhoneNumberByUUID(param.UUID, param.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.UUID = jobSeeker.UUID

	return output, nil
}

// LPから求職者の希望条件を更新
type UpdateJobSeekerDesiredFromLPInput struct {
	Param entity.UpdateJobSeekerDesiredFromLPParam
}

type UpdateJobSeekerDesiredFromLPOutput struct {
	UUID        uuid.UUID
	LogintToken uuid.UUID
}

func (i *JobSeekerInteractorImpl) UpdateJobSeekerDesiredFromLP(input UpdateJobSeekerDesiredFromLPInput) (UpdateJobSeekerDesiredFromLPOutput, error) {
	var (
		output UpdateJobSeekerDesiredFromLPOutput
		param  = input.Param
	)

	// uuidで求職者情報を取得
	jobSeeker, err := i.jobSeekerRepository.FindByUUID(param.UUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ログイントークンを作成
	logintTokenFromLP := entity.NewJobSeekerLPLoginToken(
		jobSeeker.ID, utility.CreateUUID(),
	)

	err = i.jobSeekerLPLoginTokenRepository.Create(logintTokenFromLP)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 電話番号と希望年収の更新
	err = i.jobSeekerRepository.UpdateDesiredIncomeByUUID(param.UUID, param.DesiredIncome)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 希望勤務地の更新
	err = i.jobSeekerDesiredWorkLocationRepository.DeleteByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, dwl := range param.DesiredWorkLocations {
		desiredWorkLocation := entity.NewJobSeekerDesiredWorkLocation(
			jobSeeker.ID,
			dwl,
			null.NewInt(1, true), // 第一希望で登録する
		)

		err = i.jobSeekerDesiredWorkLocationRepository.Create(desiredWorkLocation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 希望職種の更新
	err = i.jobSeekerDesiredOccupationRepository.DeleteByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, do := range param.DesiredOccupations {
		desiredOccupation := entity.NewJobSeekerDesiredOccupation(
			jobSeeker.ID,
			do,
			null.NewInt(1, true), // 第一希望で登録する
		)

		err = i.jobSeekerDesiredOccupationRepository.Create(desiredOccupation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	var (
		baseURL = os.Getenv("BASE_DOMAIN")
		// lpBaseURL = os.Getenv("BASE_DOMAIN_LP")
		seekerName  = jobSeeker.LastName + jobSeeker.FirstName
		seekerEmail = jobSeeker.Email
		from        = entity.EmailUser{Name: "autoscout事務局", Email: "info@spaceai.jp"}
		to          = entity.EmailUser{Name: seekerName, Email: seekerEmail}
		signinPage  = fmt.Sprintf("%s/guest_js/signin/?job_seeker=%s&page=matching_job&lp=%s", baseURL, jobSeeker.UUID, logintTokenFromLP.LoginToken)
		contactPage = fmt.Sprintf("%s/contact", baseURL) // autoscout LPのお問い合わせ
	)

	contents := fmt.Sprintf(
		"%s様\n\nautoscoutをご利用いただき、ありがとうございます。\n\n本登録完了いたしました。\n引き続き、よろしくお願い申し上げます。\n\n\nマイページログイン\n%s\n\n\n○●----------------------------------------------------------●○\nautoscout事務局\nお問い合わせ: %s\n○●----------------------------------------------------------●○",
		seekerName, signinPage, contactPage,
	)

	// メール送信
	err = utility.SendMailToSingle(
		i.sendgrid.APIKey,
		"アカウント発行のお知らせ",
		contents,
		from,
		to,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.UUID = jobSeeker.UUID
	output.LogintToken = logintTokenFromLP.LoginToken

	return output, nil
}

type GetJobSeekerLPRegisterStatusByUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetJobSeekerLPRegisterStatusByUUIDOutput struct {
	JobSeekerRegisterStatus *entity.JobSeekerRegisterStatus
}

func (i *JobSeekerInteractorImpl) GetJobSeekerLPRegisterStatusByUUID(input GetJobSeekerLPRegisterStatusByUUIDInput) (GetJobSeekerLPRegisterStatusByUUIDOutput, error) {
	var output GetJobSeekerLPRegisterStatusByUUIDOutput

	// uuidで求職者情報を取得
	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredOccupations, err := i.jobSeekerDesiredOccupationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	desiredWorkLocations, err := i.jobSeekerDesiredWorkLocationRepository.GetByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobSeekerRegisterStatus := entity.JobSeekerRegisterStatus{
		ID:   jobSeeker.ID,
		UUID: jobSeeker.UUID,
	}

	// IsCompletedRegisterの判定（メールアドレスと自分で設定したパスワードがあるか）
	if jobSeeker.Email != "" && jobSeeker.Password != "" {
		jobSeekerRegisterStatus.IsCompletedRegister = true
	}

	// IsCompletedDesiredの判定（電話番号が空でないか）
	if jobSeeker.PhoneNumber != "" {
		jobSeekerRegisterStatus.IsCompletedPhoneNumber = true
	}

	// IsCompletedDesiredの判定（希望条件が空でないか）
	if jobSeeker.DesiredAnnualIncome.Valid && len(desiredOccupations) > 0 && len(desiredWorkLocations) > 0 {
		jobSeekerRegisterStatus.IsCompletedDesired = true
	}

	output.JobSeekerRegisterStatus = &jobSeekerRegisterStatus

	return output, nil
}

type SendLPContactInput struct {
	Param entity.SendContactFromLPParam
}

type SendLPContactOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) SendLPContact(input SendLPContactInput) (SendLPContactOutput, error) {
	var (
		output SendLPContactOutput
		err    error
		param  = input.Param
	)

	if param.CompanyName == "" {
		param.CompanyName = "入力なし"
	}

	// Slackに送付するメッセージ
	content := fmt.Sprintf(
		"*LPからのお問い合わせ*\n\nお名前: %s\n\n会社名: %s\n\nメールアドレス: %s\n\nお問い合わせ内容: \n%s",
		param.Name,
		param.CompanyName,
		param.Email,
		param.ContacMessage,
	)

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// NOTE: https://risaki-masa.com/how-to-get-api-token-in-slack/
	slack := utility.NewSlack(i.slack.ReachAccessToken)
	err = slack.SendContact(i.slack.ReachChanelID, content)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type SendJobSeekerResetPasswordEmailForLPInput struct {
	Param entity.SendJobSeekerResetPasswordEmailFromLPParam
}

type SendJobSeekerResetPasswordEmailForLPOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) SendJobSeekerResetPasswordEmailForLP(input SendJobSeekerResetPasswordEmailForLPInput) (SendJobSeekerResetPasswordEmailForLPOutput, error) {
	var (
		output        SendJobSeekerResetPasswordEmailForLPOutput
		err           error
		SystemAgentID uint = 1
	)

	// 求職者のメールアドレスが合致するか確認
	jobSeeker, err := i.jobSeekerRepository.FindByEmailForLP(input.Param.Email, SystemAgentID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "メールアドレスが存在しません。")
			return output, wrapped
		} else {
			// Not Found以外のエラーの場合はそのままサーバーエラー
			fmt.Println(err)
			return output, err
		}
	}

	// パスワードリセットトークンを生成
	resetPasswordToken := utility.CreateUUID()

	err = i.jobSeekerRepository.UpdateResetPasswordToken(jobSeeker.UUID, resetPasswordToken.String())
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	lpBaseURL := os.Getenv("LP_BASE_DOMAIN")
	resetPasswordURL := fmt.Sprintf("%s/reset_password/%s", lpBaseURL, resetPasswordToken)

	// メール送信
	err = utility.SendMailToSingleWithoutCC(
		i.sendgrid.APIKey,
		"パスワードリセットのお手続き",
		fmt.Sprintf("お客様\n\n %s アカウントのパスワードをリセットするには、次のリンクをクリックしてください。\n\n%s\n\nこのメールに覚えがない場合は上記のリンクはクリックせず、このメールを無視してください。\nよろしくお願いいたします。\n\nautoscout事務局",
			jobSeeker.Email,
			resetPasswordURL,
		),
		entity.EmailUser{
			Name:  "autoscout事務局",
			Email: "info@spaceai.jp",
		},
		entity.EmailUser{
			Name:  jobSeeker.LastName + jobSeeker.FirstName + "様",
			Email: jobSeeker.Email,
		},
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}
	output.OK = true

	return output, nil
}

type ResetPasswordForLPInput struct {
	Param entity.ResetPasswordFromLPParam
}

type ResetPasswordForLPOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) ResetPasswordForLP(input ResetPasswordForLPInput) (ResetPasswordForLPOutput, error) {
	var (
		output             = ResetPasswordForLPOutput{}
		SystemAgentID uint = 1
	)

	// ResetPasswordTokenが一致するか確認
	jobSeeker, err := i.jobSeekerRepository.FindByResetPasswordTokenForLP(input.Param.ResetPasswordToken, SystemAgentID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "不正なURLです。")
			return output, wrapped
		} else {
			// Not Found以外のエラーの場合はそのままサーバーエラー
			fmt.Println(err)
			return output, err
		}
	}

	// パスワードのハッシュ化
	hashedPassword, err := hashPassword(input.Param.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// パスワードの更新
	err = i.jobSeekerRepository.UpdatePassword(jobSeeker.UUID, hashedPassword)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// NOTE: ResetPasswordTokenを空で更新するのはトークンの有無で表示を切り替えるため
	// パスワードトークンを削除（空で更新する）
	err = i.jobSeekerRepository.UpdateResetPasswordToken(jobSeeker.UUID, "")
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type CheckResetPasswordTokenInput struct {
	ResetPasswordToken string
}

type CheckResetPasswordTokenOutput struct {
	OK bool
}

func (i *JobSeekerInteractorImpl) CheckResetPasswordToken(input CheckResetPasswordTokenInput) (CheckResetPasswordTokenOutput, error) {
	var (
		output             = CheckResetPasswordTokenOutput{}
		SystemAgentID uint = 1
	)

	// ResetPasswordTokenが一致するか確認
	_, err := i.jobSeekerRepository.FindByResetPasswordTokenForLP(input.ResetPasswordToken, SystemAgentID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "不正なURLです。")
			output.OK = false
			return output, wrapped
		} else {
			// Not Found以外のエラーの場合はそのままサーバーエラー
			fmt.Println(err)
			output.OK = false
			return output, err
		}
	}

	output.OK = true

	return output, nil
}
