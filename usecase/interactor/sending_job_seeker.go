package interactor

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerInteractor interface {
	// 汎用系 API
	CreateSendingJobSeeker(input CreateSendingJobSeekerInput) (CreateSendingJobSeekerOutput, error)
	CreateSendingJobSeekerFromJobSeeker(input CreateSendingJobSeekerFromJobSeekerInput) (CreateSendingJobSeekerFromJobSeekerOutput, error) // CRM求職者から送客求職者を作成
	FirstUpdateSendingJobSeeker(input FirstUpdateSendingJobSeekerInput) (FirstUpdateSendingJobSeekerOutput, error)
	UpdateSendingJobSeeker(input UpdateSendingJobSeekerInput) (UpdateSendingJobSeekerOutput, error)
	UpdateSendingJobSeekerPhase(input UpdateSendingJobSeekerPhaseInput) (UpdateSendingJobSeekerPhaseOutput, error)
	UpdateSendingInterviewDateBySendingJobSeekerID(input UpdateSendingJobSeekerInterviewDateInput) (UpdateSendingJobSeekerInterviewDateOutput, error)
	UpdateSendingJobSeekerActivityMemo(input UpdateSendingJobSeekerActivityMemoInput) (UpdateSendingJobSeekerActivityMemoOutput, error)
	UpdateIsVewForWating(input UpdateIsVewForWatingInput) (UpdateIsVewForWatingOutput, error)
	UpdateIsVewForUnregister(input UpdateIsVewForUnregisterInput) (UpdateIsVewForUnregisterOutput, error)
	DeleteSendingJobSeeker(input DeleteSendingJobSeekerInput) (DeleteSendingJobSeekerOutput, error)

	GetSendingJobSeekerByID(input GetSendingJobSeekerByIDInput) (GetSendingJobSeekerByIDOutput, error)       // 指定IDの送客求職者を取得する関数
	GetSendingJobSeekerByUUID(input GetSendingJobSeekerByUUIDInput) (GetSendingJobSeekerByUUIDOutput, error) // 指定uuidの送客求職者を取得する関
	GetIsNotViewSendingJobSeekerCountByAgentStaffID(input GetIsNotViewSendingJobSeekerCountByAgentStaffIDInput) (GetIsNotViewSendingJobSeekerCountByAgentStaffIDOutput, error)

	// { value: number; label: string }[]の形式でリストを取得する
	GetSearchListForSendingJobSeekerManagementByAgentID(input GetSearchListForSendingJobSeekerManagementByAgentIDInput) (GetSearchListForSendingJobSeekerManagementByAgentIDOutput, error)

	// 面談前アンケート関連 API
	CreateSendingInitialQuestionnaire(input CreateSendingInitialQuestionnaireInput) (CreateSendingInitialQuestionnairdhutput, error) // 面談前アンケートを登録 (求職者の同意項目、業界、職種、勤務地、質問要望、ファイル（履歴書（原本）、職務経歴書（原本）））

	// 送客終了理由 API
	CreateSendingJobSeekerEndStatus(input CreateSendingJobSeekerEndStatusInput) (CreateSendingJobSeekerEndStatusOutput, error) // 終了理由を登録

	// 送客求職者資料関連 API
	UpdateSendingJobSeekerDocument(input UpdateSendingJobSeekerDocumentInput) (UpdateSendingJobSeekerDocumentOutput, error)
	GetSendingJobSeekerDocumentByUUID(input GetSendingJobSeekerDocumentByUUIDInput) (GetSendingJobSeekerDocumentByUUIDOutput, error) // 指定uuidの送客求職者の応募書類データを取得する関数
	GetSendingJobSeekerDocumentBySendingJobSeekerID(input GetSendingJobSeekerDocumentBySendingJobSeekerIDInput) (GetSendingJobSeekerDocumentBySendingJobSeekerIDOutput, error)
	DeleteSendingJobSeekerResumePDFURL(input DeleteSendingJobSeekerResumePDFURLInput) (DeleteSendingJobSeekerResumePDFURLOutput, error)
	DeleteSendingJobSeekerResumeOriginURL(input DeleteSendingJobSeekerResumeOriginURLInput) (DeleteSendingJobSeekerResumeOriginURLOutput, error)
	DeleteSendingJobSeekerCVPDFURL(input DeleteSendingJobSeekerCVPDFURLInput) (DeleteSendingJobSeekerCVPDFURLOutput, error)
	DeleteSendingJobSeekerCVOriginURL(input DeleteSendingJobSeekerCVOriginURLInput) (DeleteSendingJobSeekerCVOriginURLOutput, error)
	DeleteSendingJobSeekerRecommendationPDFURL(input DeleteSendingJobSeekerRecommendationPDFURLInput) (DeleteSendingJobSeekerRecommendationPDFURLOutput, error)
	DeleteSendingJobSeekerRecommendationOriginURL(input DeleteSendingJobSeekerRecommendationOriginURLInput) (DeleteSendingJobSeekerRecommendationOriginURLOutput, error)
	DeleteSendingJobSeekerIDPhotoURL(input DeleteSendingJobSeekerIDPhotoURLInput) (DeleteSendingJobSeekerIDPhotoURLOutput, error)
	DeleteSendingJobSeekerOtherDocument1URL(input DeleteSendingJobSeekerOtherDocument1URLInput) (DeleteSendingJobSeekerOtherDocument1URLOutput, error)
	DeleteSendingJobSeekerOtherDocument2URL(input DeleteSendingJobSeekerOtherDocument2URLInput) (DeleteSendingJobSeekerOtherDocument2URLOutput, error)
	DeleteSendingJobSeekerOtherDocument3URL(input DeleteSendingJobSeekerOtherDocument3URLInput) (DeleteSendingJobSeekerOtherDocument3URLOutput, error)

	// 子テーブルのセット
	setSendingJobSeekerChildTable(sendingJobSeeker *entity.SendingJobSeeker) (*entity.SendingJobSeeker, error)
}

type SendingJobSeekerInteractorImpl struct {
	firebase                                       usecase.Firebase
	sendgrid                                       config.Sendgrid
	oneSignal                                      config.OneSignal
	sendingCustomerRepository                      usecase.SendingCustomerRepository
	sendingJobSeekerRepository                     usecase.SendingJobSeekerRepository
	sendingJobSeekerStudentHistoryRepository       usecase.SendingJobSeekerStudentHistoryRepository
	sendingJobSeekerWorkHistoryRepository          usecase.SendingJobSeekerWorkHistoryRepository
	sendingJobSeekerExperienceIndustryRepository   usecase.SendingJobSeekerExperienceIndustryRepository
	sendingJobSeekerDepartmentHistoryRepository    usecase.SendingJobSeekerDepartmentHistoryRepository
	sendingJobSeekerLicenseRepository              usecase.SendingJobSeekerLicenseRepository
	sendingJobSeekerSelfPromotionRepository        usecase.SendingJobSeekerSelfPromotionRepository
	sendingJobSeekerDocumentRepository             usecase.SendingJobSeekerDocumentRepository
	sendingJobSeekerDesiredIndustryRepository      usecase.SendingJobSeekerDesiredIndustryRepository
	sendingJobSeekerDesiredOccupationRepository    usecase.SendingJobSeekerDesiredOccupationRepository
	sendingJobSeekerDesiredWorkLocationRepository  usecase.SendingJobSeekerDesiredWorkLocationRepository
	sendingJobSeekerDesiredHolidayTypeRepository   usecase.SendingJobSeekerDesiredHolidayTypeRepository
	sendingJobSeekerDevelopmentSkillRepository     usecase.SendingJobSeekerDevelopmentSkillRepository
	sendingJobSeekerLanguageSkillRepository        usecase.SendingJobSeekerLanguageSkillRepository
	sendingJobSeekerPCToolRepository               usecase.SendingJobSeekerPCToolRepository
	sendingJobSeekerExperienceOccupationRepository usecase.SendingJobSeekerExperienceOccupationRepository
	sendingJobSeekerDesiredCompanyScaleRepository  usecase.SendingJobSeekerDesiredCompanyScaleRepository
	agentRepository                                usecase.AgentRepository
	agentStaffRepository                           usecase.AgentStaffRepository
	chatGroupWithSendingJobSeekerRepository        usecase.ChatGroupWithSendingJobSeekerRepository
	jobSeekerRepository                            usecase.JobSeekerRepository
	sendingJobSeekerEndStatusRepository            usecase.SendingJobSeekerEndStatusRepository
	sendingJobSeekerIsViewRepository               usecase.SendingJobSeekerIsViewRepository
	sendingEnterpriseRepository                    usecase.SendingEnterpriseRepository
}

// SendingJobSeekerInteractorImpl is an implementation of SendingJobSeekerInteractor
func NewSendingJobSeekerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	scR usecase.SendingCustomerRepository,
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
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	cgR usecase.ChatGroupWithSendingJobSeekerRepository,
	jsR usecase.JobSeekerRepository,
	sjsesR usecase.SendingJobSeekerEndStatusRepository,
	sjsivR usecase.SendingJobSeekerIsViewRepository,
	seR usecase.SendingEnterpriseRepository,
) SendingJobSeekerInteractor {
	return &SendingJobSeekerInteractorImpl{
		firebase:                                       fb,
		sendgrid:                                       sg,
		oneSignal:                                      os,
		sendingCustomerRepository:                      scR,
		sendingJobSeekerRepository:                     sjsR,
		sendingJobSeekerStudentHistoryRepository:       sjsshR,
		sendingJobSeekerWorkHistoryRepository:          sjswhR,
		sendingJobSeekerExperienceIndustryRepository:   sjseiR,
		sendingJobSeekerExperienceOccupationRepository: sjseoR,
		sendingJobSeekerLicenseRepository:              sjslR,
		sendingJobSeekerSelfPromotionRepository:        sjsspR,
		sendingJobSeekerDocumentRepository:             sjsdR,
		sendingJobSeekerDesiredIndustryRepository:      sjsdiR,
		sendingJobSeekerDesiredOccupationRepository:    sjsdoR,
		sendingJobSeekerDesiredWorkLocationRepository:  sjsdwlR,
		sendingJobSeekerDesiredHolidayTypeRepository:   sjsdhtR,
		sendingJobSeekerDevelopmentSkillRepository:     sjsdsR,
		sendingJobSeekerLanguageSkillRepository:        sjslsR,
		sendingJobSeekerPCToolRepository:               sjsptR,
		sendingJobSeekerDepartmentHistoryRepository:    sjsdhR,
		sendingJobSeekerDesiredCompanyScaleRepository:  sjsdcsR,
		agentRepository:                                aR,
		agentStaffRepository:                           asR,
		chatGroupWithSendingJobSeekerRepository:        cgR,
		jobSeekerRepository:                            jsR,
		sendingJobSeekerEndStatusRepository:            sjsesR,
		sendingJobSeekerIsViewRepository:               sjsivR,
		sendingEnterpriseRepository:                    seR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// 最初の作成 *作成時は名前と連絡先のみで、それ以外の項目は更新時に作成する
type CreateSendingJobSeekerInput struct {
	CreateParam  entity.CreateSendingJobSeekerParam
	AgentStaffID uint
}

type CreateSendingJobSeekerOutput struct {
	SendingCustomer *entity.SendingCustomer
}

func (i *SendingJobSeekerInteractorImpl) CreateSendingJobSeeker(input CreateSendingJobSeekerInput) (CreateSendingJobSeekerOutput, error) {
	var (
		output CreateSendingJobSeekerOutput
		err    error
	)

	/************ 1. 重複チェック **************/

	// ①求職者の重複はできない。名前、電話番号が一致している場合に、重複で送客できないようにする。
	duplicatedJobSeeker, err := i.jobSeekerRepository.FindByNameAndPhoneNumberByMotoyuiAgent(
		input.CreateParam.FirstName, input.CreateParam.LastName,
		input.CreateParam.FirstFurigana, input.CreateParam.LastFurigana,
		input.CreateParam.PhoneNumber,
	)
	if err != nil {
		if !errors.Is(err, entity.ErrNotFound) {
			// Not Found以外のエラーの場合はエラーでストップ
			fmt.Println(err)
			return output, err
		}
	}

	// Not Foundの場合は次の重複チェックへ
	duplicatedSendingJobSeeker, err := i.sendingJobSeekerRepository.FindByNameAndPhoneNumber(
		input.CreateParam.FirstName, input.CreateParam.LastName,
		input.CreateParam.FirstFurigana, input.CreateParam.LastFurigana,
		input.CreateParam.PhoneNumber,
	)
	if err != nil {
		if !errors.Is(err, entity.ErrNotFound) {
			// Not Found以外のエラーの場合は意図指定なためエラーでストップ
			fmt.Println(err)
			return output, err
		}
	}

	// 例外）但し送客登録日（日程を登録するをクリックした日）が、元結CRMの面談日から6ヶ月を経過している かつ フェーズが「転職活動終了」の求職者は送客できる
	if duplicatedJobSeeker != nil &&
		duplicatedJobSeeker.CreatedAt.AddDate(0, 6, 0).After(time.Now()) &&
		duplicatedJobSeeker.Phase.Int64 != int64(entity.QuitedAfterInterview) {
		fmt.Println("既にMotoyuiに登録済みのため送客できません")
		return output, errors.New("既にMotoyuiに登録済みのため送客できません")
	}

	// ※上の条件を満たしていても、送客管理に入っていれば対象外。
	if duplicatedSendingJobSeeker != nil {
		fmt.Println("既に送客管理に登録済みのため送客できません")
		return output, errors.New("既に送客管理に登録済みのため送客できません")
	}

	/************ 2. sending_customersにレコード作成 **************/

	sendingCustomer := entity.NewSendingCustomer(
		input.CreateParam.AgentID,
		input.CreateParam.LastName,
		input.CreateParam.FirstName,
		input.CreateParam.LastFurigana,
		input.CreateParam.FirstFurigana,
		input.CreateParam.PhoneNumber,
		input.CreateParam.Email,
		"",                     // input.CreateParam.ResumeOriginURL,
		"",                     // input.CreateParam.ResumePDFURL,
		"",                     // input.CreateParam.CVOriginURL,
		"",                     // input.CreateParam.CVPDFURL,
		utility.EarliestTime(), // input.CreateParam.InterviewDate *初期値
		"",                     // input.CreateParam.InterviewInformation,
		"",                     // input.CreateParam.Remarks,
		NullInt,                // input.CreateParam.Gender,
		NullInt,                // input.CreateParam.Nationality,
		"",                     // input.CreateParam.NationalityRemarks,
		"",                     // input.CreateParam.Birthday,
		"",                     // input.CreateParam.PostCode,
		NullInt,                // input.CreateParam.Prefecture,
		"",                     // input.CreateParam.Address,
		"",                     // input.CreateParam.AddressFurigana,
		"",                     // input.CreateParam.SchoolName,
		"",                     // input.CreateParam.Subject,
		"",                     // input.CreateParam.EntranceYear,
		"",                     // input.CreateParam.GraduationYear,
		NullInt,                // input.CreateParam.StateOfEmployment,
		NullInt,                // input.CreateParam.JobChange,
		"",                     // input.CreateParam.JobSummary,
		"",                     // input.CreateParam.CompanyName,
		"",                     // input.CreateParam.JoiningYear,
		"",                     // input.CreateParam.RetireYear,
		NullInt,                // input.CreateParam.FirstStatus,
		NullInt,                // input.CreateParam.LastStatus,
		"",                     // input.CreateParam.JobDescription,
		"",                     // input.CreateParam.HistorySupplement,
		null.NewInt(int64(entity.UnregisterSchedule), true), // Phase = 日程未登録
	)

	err = i.sendingCustomerRepository.Create(sendingCustomer)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 開発環境の場合はユーザー名などを統一する
	if os.Getenv("APP_ENV") != "prd" {
		err = i.sendingCustomerRepository.UpdateForDev(sendingCustomer)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 3. sending_job_seekersテーブルにレコード作成 **************/

	sendingJobSeeker := entity.NewSendingJobSeeker(
		1,       // 通常の送客フローの場合はアンドイーズに送客されるため「1」で登録する
		NullInt, // agentStaffID
		"",      // LineID
		NullInt, // UserStatus
		input.CreateParam.LastName,
		input.CreateParam.FirstName,
		input.CreateParam.LastFurigana,
		input.CreateParam.FirstFurigana,
		NullInt, // input.CreateParam.Gender,
		"",      // GenderRemarks
		"",      // input.CreateParam.Birthday,
		NullInt, // Spouse,
		NullInt, // SupportObligation,
		NullInt, // Dependents,
		input.CreateParam.PhoneNumber,
		input.CreateParam.Email,
		"",      // EmergencyPhoneNumber,
		"",      // input.CreateParam.PostCode,
		NullInt, // input.CreateParam.Prefecture,
		"",      // input.CreateParam.Address,
		"",      // input.CreateParam.AddressFurigana,
		NullInt, // input.CreateParam.StateOfEmployment,
		"",      // input.CreateParam.JobSummary,
		"",      // input.CreateParam.HistorySupplement,
		"",      // ResearchContent,
		NullInt, // JoinCompanyPeriod,
		NullInt, // input.CreateParam.JobChange,
		NullInt, // AnnualIncome,
		NullInt, // DesiredAnnualIncome,
		NullInt, // Transfer,
		"",      // TransferRequirement,
		NullInt, // ShortResignation,
		"",      // ShortResignationRemarks,
		NullInt, // MedicalHistory,
		NullInt, // input.CreateParam.Nationality,
		NullInt, // Appearance,
		"",      // AppearanceDetailOfTruth,
		"",      // AppearanceDetail,
		NullInt, // Communication,
		"",      // CommunicationDetailOfTruth,
		"",      // CommunicationDetail,
		NullInt, // Thinking,
		"",      // ThinkingDetailOfTruth,
		"",      // ThinkingDetail,
		"",      // SecretMemo,
		NullInt, // JobHuntingState,
		"",      // RecommendReason,
		null.NewInt(int64(entity.UnregisterSchedule), true), // Phase = 日程未登録
		utility.EarliestTime(),                              // input.CreateParam.InterviewDate *初期値
		false,                                               // Agreement, 同意なし
		NullInt,                                             // StudyCategory,
		NullInt,                                             // WordSkill,
		NullInt,                                             // ExcelSkill,
		NullInt,                                             // PowerPointSkill,
		"",                                                  // PublicMemo,
		"",                                                  // input.CreateParam.NationalityRemarks,
		"",                                                  // MedicalHistoryRemarks,
		"",                                                  // AcceptancePoints,
		sendingCustomer.ID,                                  // SendingCustomerID,
	)

	err = i.sendingJobSeekerRepository.Create(sendingJobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 開発環境の場合はユーザー名などを統一する
	if os.Getenv("APP_ENV") != "prd" {
		err = i.sendingJobSeekerRepository.UpdateForDev(sendingJobSeeker)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 4. sending_job_seekersテーブルの子テーブル情報を登録 **************/

	// 書類
	sendingJobSeekerDocument := entity.NewSendingJobSeekerDocument(
		sendingJobSeeker.ID,
		"", // ResumeOriginURL,
		"", // ResumePDFURL,
		"", // CVOriginURL,
		"", // CVPDFURL,
		"", // RecommendationOriginURL,
		"", // RecommendationPDFURL,
		"", // IDPhotoURL,
		"", // OtherDocument1URL,
		"", // OtherDocument2URL,
		"", // OtherDocument3URL,
	)

	err = i.sendingJobSeekerDocumentRepository.Create(sendingJobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者作成時にメッセージグループ作成
	// エージェントと求職者のチャットグループを作成
	chatGroup := entity.NewChatGroupWithSendingJobSeeker(
		sendingJobSeeker.ID,
		false, // 初めはLINE連携してないから false
	)

	err = i.chatGroupWithSendingJobSeekerRepository.Create(chatGroup)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客求職者の進捗管理上での閲覧状況を管理する情報を登録
	isView := entity.NewSendingJobSeekerIsView(
		sendingJobSeeker.ID,
		true,
		true,
	)

	err = i.sendingJobSeekerIsViewRepository.Create(isView)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingCustomer = sendingCustomer

	return output, nil
}

// CRM求職者から送客求職者を作成
type CreateSendingJobSeekerFromJobSeekerInput struct {
	CreateParam entity.JobSeeker
}

type CreateSendingJobSeekerFromJobSeekerOutput struct {
	SendingJobSeeker *entity.SendingJobSeeker
}

func (i *SendingJobSeekerInteractorImpl) CreateSendingJobSeekerFromJobSeeker(input CreateSendingJobSeekerFromJobSeekerInput) (CreateSendingJobSeekerFromJobSeekerOutput, error) {
	var (
		output CreateSendingJobSeekerFromJobSeekerOutput
		err    error
	)

	// Not Foundの場合は次の重複チェックへ
	duplicatedSendingJobSeeker, err := i.sendingJobSeekerRepository.FindByNameAndPhoneNumber(
		input.CreateParam.FirstName, input.CreateParam.LastName,
		input.CreateParam.FirstFurigana, input.CreateParam.LastFurigana,
		input.CreateParam.PhoneNumber,
	)
	if err != nil {
		if !errors.Is(err, entity.ErrNotFound) {
			// Not Found以外のエラーの場合は意図指定なためエラーでストップ
			fmt.Println(err)
			return output, err
		}
	}

	// ※上の条件を満たしていても、送客管理に入っていれば対象外。
	if duplicatedSendingJobSeeker != nil {
		fmt.Println("既に送客管理に登録済みのため送客できません")
		return output, errors.New("既に送客管理に登録済みのため送客できません")
	}

	// interviewDate, err := time.Parse("2006-01-02T15:04", input.CreateParam.InterviewDate)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return output, err
	// }

	/************ sending_customersに作成 **************/

	sendingCustomer := entity.NewSendingCustomer(
		input.CreateParam.AgentID,
		input.CreateParam.LastName,
		input.CreateParam.FirstName,
		input.CreateParam.LastFurigana,
		input.CreateParam.FirstFurigana,
		input.CreateParam.PhoneNumber,
		input.CreateParam.Email,
		input.CreateParam.Documents.ResumeOriginURL,
		input.CreateParam.Documents.ResumePDFURL,
		input.CreateParam.Documents.CVOriginURL,
		input.CreateParam.Documents.CVPDFURL,
		input.CreateParam.InterviewDate,
		"", // input.CreateParam.InterviewInformation,
		"", // input.CreateParam.Remarks,
		input.CreateParam.Gender,
		input.CreateParam.Nationality,
		input.CreateParam.NationalityRemarks,
		input.CreateParam.Birthday,
		input.CreateParam.PostCode,
		input.CreateParam.Prefecture,
		input.CreateParam.Address,
		input.CreateParam.AddressFurigana,
		"", // input.CreateParam.SchoolName,
		"", // input.CreateParam.Subject,
		"", // input.CreateParam.EntranceYear,
		"", // input.CreateParam.GraduationYear,
		input.CreateParam.StateOfEmployment,
		input.CreateParam.JobChange,
		input.CreateParam.JobSummary,
		"",      // input.CreateParam.CompanyName,
		"",      // input.CreateParam.JoiningYear,
		"",      // input.CreateParam.RetireYear,
		NullInt, // input.CreateParam.FirstStatus,
		NullInt, // input.CreateParam.LastStatus,
		"",      // input.CreateParam.JobDescription,
		input.CreateParam.HistorySupplement,
		null.NewInt(int64(entity.WaitingForInterview), true), // Phase = 面談実施待ち
	)

	err = i.sendingCustomerRepository.Create(sendingCustomer)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ job_seekersに作成 **************/

	sendingJobSeeker := entity.NewSendingJobSeeker(
		input.CreateParam.AgentID,
		input.CreateParam.AgentStaffID,
		input.CreateParam.LineID,
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
		"", // input.CreateParam.AppearanceDetailOfTruth,
		"", // input.CreateParam.AppearanceDetail,
		input.CreateParam.Communication,
		"", // input.CreateParam.CommunicationDetailOfTruth,
		"", // input.CreateParam.CommunicationDetail,
		input.CreateParam.Thinking,
		"", // input.CreateParam.ThinkingDetailOfTruth,
		"", // input.CreateParam.ThinkingDetail,
		input.CreateParam.SecretMemo,
		input.CreateParam.JobHuntingState,
		input.CreateParam.RecommendReason,
		null.NewInt(int64(entity.WaitingForInterview), true), // Phase = 面談実施待ち
		input.CreateParam.InterviewDate,
		input.CreateParam.Agreement,
		input.CreateParam.StudyCategory,
		input.CreateParam.WordSkill,
		input.CreateParam.ExcelSkill,
		input.CreateParam.PowerPointSkill,
		"", // input.CreateParam.PublicMemo,
		input.CreateParam.NationalityRemarks,
		input.CreateParam.MedicalHistoryRemarks,
		input.CreateParam.AcceptancePoints,
		sendingCustomer.ID, // SendingCustomerID,
	)

	err = i.sendingJobSeekerRepository.Create(sendingJobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客求職者の進捗管理上での閲覧状況を管理する情報を登録
	isView := entity.NewSendingJobSeekerIsView(
		sendingJobSeeker.ID,
		true,
		true,
	)

	err = i.sendingJobSeekerIsViewRepository.Create(isView)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 書類
	sendingJobSeekerDocument := entity.NewSendingJobSeekerDocument(
		sendingJobSeeker.ID,
		input.CreateParam.Documents.ResumeOriginURL,
		input.CreateParam.Documents.ResumePDFURL,
		input.CreateParam.Documents.CVOriginURL,
		input.CreateParam.Documents.CVPDFURL,
		input.CreateParam.Documents.RecommendationOriginURL,
		input.CreateParam.Documents.RecommendationPDFURL,
		input.CreateParam.Documents.IDPhotoURL,
		input.CreateParam.Documents.OtherDocument1URL,
		input.CreateParam.Documents.OtherDocument2URL,
		input.CreateParam.Documents.OtherDocument3URL,
	)

	err = i.sendingJobSeekerDocumentRepository.Create(sendingJobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sh := range input.CreateParam.StudentHistories {
		studentHistory := entity.NewSendingJobSeekerStudentHistory(
			sendingJobSeeker.ID,
			sh.SchoolCategory,
			sh.SchoolName,
			sh.SchoolLevel,
			sh.Subject,
			sh.EntranceYear,
			sh.FirstStatus,
			sh.GraduationYear,
			sh.LastStatus,
		)

		err = i.sendingJobSeekerStudentHistoryRepository.Create(studentHistory)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, wh := range input.CreateParam.WorkHistories {
		workHistory := entity.NewSendingJobSeekerWorkHistory(
			sendingJobSeeker.ID,
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

		err = i.sendingJobSeekerWorkHistoryRepository.Create(workHistory)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, ei := range wh.ExperienceIndustries {
			experienceIndustry := entity.NewSendingJobSeekerExperienceIndustry(
				workHistory.ID,
				ei.Industry,
			)

			err = i.sendingJobSeekerExperienceIndustryRepository.Create(experienceIndustry)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, dh := range wh.DepartmentHistories {
			departmentHistory := entity.NewSendingJobSeekerDepartmentHistory(
				workHistory.ID,
				dh.Department,
				dh.ManagementNumber,
				dh.ManagementDetail,
				dh.JobDescription,
				dh.StartYear,
				dh.EndYear,
			)

			err = i.sendingJobSeekerDepartmentHistoryRepository.Create(departmentHistory)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, occupations := range dh.ExperienceOccupations {
				occupation := entity.NewSendingJobSeekerExperienceOccupation(
					departmentHistory.ID,
					occupations.Occupation,
				)

				err = i.sendingJobSeekerExperienceOccupationRepository.Create(occupation)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

		}
	}

	for _, dcs := range input.CreateParam.DesiredCompanyScales {
		disiredCompanyScale := entity.NewSendingJobSeekerDesiredCompanyScale(
			sendingJobSeeker.ID,
			dcs.DesiredCompanyScale,
		)

		err = i.sendingJobSeekerDesiredCompanyScaleRepository.Create(disiredCompanyScale)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, l := range input.CreateParam.Licenses {
		license := entity.NewSendingJobSeekerLicense(
			sendingJobSeeker.ID,
			l.LicenseType,
			l.AcquisitionTime,
		)

		err = i.sendingJobSeekerLicenseRepository.Create(license)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, sp := range input.CreateParam.SelfPromotions {
		selfPromotion := entity.NewSendingJobSeekerSelfPromotion(
			sendingJobSeeker.ID,
			sp.Title,
			sp.Contents,
		)

		err = i.sendingJobSeekerSelfPromotionRepository.Create(selfPromotion)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, di := range input.CreateParam.DesiredIndustries {
		desiredIndustry := entity.NewSendingJobSeekerDesiredIndustry(
			sendingJobSeeker.ID,
			di.DesiredIndustry,
			di.DesiredRank,
		)

		err = i.sendingJobSeekerDesiredIndustryRepository.Create(desiredIndustry)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, do := range input.CreateParam.DesiredOccupations {
		desiredOccupation := entity.NewSendingJobSeekerDesiredOccupation(
			sendingJobSeeker.ID,
			do.DesiredOccupation,
			do.DesiredRank,
		)

		err = i.sendingJobSeekerDesiredOccupationRepository.Create(desiredOccupation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, dwl := range input.CreateParam.DesiredWorkLocations {
		desiredWorkLocation := entity.NewSendingJobSeekerDesiredWorkLocation(
			sendingJobSeeker.ID,
			dwl.DesiredWorkLocation,
			dwl.DesiredRank,
		)

		err = i.sendingJobSeekerDesiredWorkLocationRepository.Create(desiredWorkLocation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, dht := range input.CreateParam.DesiredHolidayTypes {
		desiredHolidayType := entity.NewSendingJobSeekerDesiredHolidayType(
			sendingJobSeeker.ID,
			dht.HolidayType,
		)

		err = i.sendingJobSeekerDesiredHolidayTypeRepository.Create(desiredHolidayType)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, ds := range input.CreateParam.DevelopmentSkills {
		developmentSkill := entity.NewSendingJobSeekerDevelopmentSkill(
			sendingJobSeeker.ID,
			ds.DevelopmentCategory,
			ds.DevelopmentType,
			ds.ExperienceYear,
			ds.ExperienceMonth,
		)

		err = i.sendingJobSeekerDevelopmentSkillRepository.Create(developmentSkill)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, ls := range input.CreateParam.LanguageSkills {
		languageSkill := entity.NewSendingJobSeekerLanguageSkill(
			sendingJobSeeker.ID,
			ls.LanguageType,
			ls.LanguageLevel,
			ls.Toeic,
			ls.ToeicExaminationYear,
			ls.ToeflIBT,
			ls.ToeflIBTExaminationYear,
			ls.ToeflPBT,
			ls.ToeflPBTExaminationYear,
		)

		err = i.sendingJobSeekerLanguageSkillRepository.Create(languageSkill)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, ps := range input.CreateParam.PCTools {
		PCTool := entity.NewSendingJobSeekerPCTool(
			sendingJobSeeker.ID,
			ps.Tool,
		)

		err = i.sendingJobSeekerPCToolRepository.Create(PCTool)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 求職者作成時にメッセージグループ作成
	// エージェントと求職者のチャットグループを作成
	chatGroup := entity.NewChatGroupWithSendingJobSeeker(
		sendingJobSeeker.ID,
		false, // 初めはLINE連携してないから false
	)

	err = i.chatGroupWithSendingJobSeekerRepository.Create(chatGroup)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingJobSeeker = sendingJobSeeker

	return output, nil
}

// 最初の更新 *作成時は名前と連絡先のみで、それ以外の項目は更新時に作成する
type FirstUpdateSendingJobSeekerInput struct {
	UpdateParam       entity.FirstUpdateSendingJobSeekerParam
	SendingCustomerID uint
}

type FirstUpdateSendingJobSeekerOutput struct {
	SendingJobSeeker *entity.SendingJobSeeker
}

func (i *SendingJobSeekerInteractorImpl) FirstUpdateSendingJobSeeker(input FirstUpdateSendingJobSeekerInput) (FirstUpdateSendingJobSeekerOutput, error) {
	var (
		output FirstUpdateSendingJobSeekerOutput
		err    error
	)

	// 送客IDから登録済みのjob_seekersを取得
	prevSendingJobSeeker, err := i.sendingJobSeekerRepository.FindBySendingCustomerID(input.SendingCustomerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ sending_customersに作成 **************/
	sendingCustomer := entity.NewSendingCustomer(
		input.UpdateParam.AgentID,
		input.UpdateParam.LastName,
		input.UpdateParam.FirstName,
		input.UpdateParam.LastFurigana,
		input.UpdateParam.FirstFurigana,
		input.UpdateParam.PhoneNumber,
		input.UpdateParam.Email,
		input.UpdateParam.ResumeOriginURL,
		input.UpdateParam.ResumePDFURL,
		input.UpdateParam.CVOriginURL,
		input.UpdateParam.CVPDFURL,
		input.UpdateParam.InterviewDate,
		input.UpdateParam.InterviewInformation,
		input.UpdateParam.Remarks,
		input.UpdateParam.Gender,
		input.UpdateParam.Nationality,
		input.UpdateParam.NationalityRemarks,
		input.UpdateParam.Birthday,
		input.UpdateParam.PostCode,
		input.UpdateParam.Prefecture,
		input.UpdateParam.Address,
		input.UpdateParam.AddressFurigana,
		input.UpdateParam.SchoolName,
		input.UpdateParam.Subject,
		input.UpdateParam.EntranceYear,
		input.UpdateParam.GraduationYear,
		input.UpdateParam.StateOfEmployment,
		input.UpdateParam.JobChange,
		input.UpdateParam.JobSummary,
		input.UpdateParam.CompanyName,
		input.UpdateParam.JoiningYear,
		input.UpdateParam.RetireYear,
		input.UpdateParam.FirstStatus,
		input.UpdateParam.LastStatus,
		input.UpdateParam.JobDescription,
		input.UpdateParam.HistorySupplement,
		null.NewInt(int64(entity.WaitingForInterview), true), // Phase = 面談実施待ち
	)

	err = i.sendingCustomerRepository.Update(prevSendingJobSeeker.SendingCustomerID, sendingCustomer)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ job_seekersを更新 **************/

	sendingJobSeeker := entity.NewSendingJobSeeker(
		input.UpdateParam.AgentID,
		prevSendingJobSeeker.AgentStaffID, // agentStaffIDはDBのデータをそのまま使う
		"",                                // LineID
		NullInt,                           // UserStatus
		input.UpdateParam.LastName,
		input.UpdateParam.FirstName,
		input.UpdateParam.LastFurigana,
		input.UpdateParam.FirstFurigana,
		input.UpdateParam.Gender,
		"", // GenderRemarks
		input.UpdateParam.Birthday,
		NullInt, // Spouse,
		NullInt, // SupportObligation,
		NullInt, // Dependents,
		input.UpdateParam.PhoneNumber,
		input.UpdateParam.Email,
		"", // EmergencyPhoneNumber,
		input.UpdateParam.PostCode,
		input.UpdateParam.Prefecture,
		input.UpdateParam.Address,
		input.UpdateParam.AddressFurigana,
		input.UpdateParam.StateOfEmployment,
		input.UpdateParam.JobSummary,
		input.UpdateParam.HistorySupplement,
		"",      // ResearchContent,
		NullInt, // JoinCompanyPeriod,
		input.UpdateParam.JobChange,
		NullInt, // AnnualIncome,
		NullInt, // DesiredAnnualIncome,
		NullInt, // Transfer,
		"",      // TransferRequirement,
		NullInt, // ShortResignation,
		"",      // ShortResignationRemarks,
		NullInt, // MedicalHistory,
		input.UpdateParam.Nationality,
		NullInt, // Appearance,
		"",      // AppearanceDetailOfTruth,
		"",      // AppearanceDetail,
		NullInt, // Communication,
		"",      // CommunicationDetailOfTruth,
		"",      // CommunicationDetail,
		NullInt, // Thinking,
		"",      // ThinkingDetailOfTruth,
		"",      // ThinkingDetail,
		"",      // SecretMemo,
		NullInt, // JobHuntingState,
		"",      // RecommendReason,
		null.NewInt(int64(entity.WaitingForInterview), true), // Phase = 面談実施待ち
		input.UpdateParam.InterviewDate,
		true,    // Agreement,
		NullInt, // StudyCategory,
		NullInt, // WordSkill,
		NullInt, // ExcelSkill,
		NullInt, // PowerPointSkill,
		"",      // PublicMemo,
		input.UpdateParam.NationalityRemarks,
		"",                                     // MedicalHistoryRemarks,
		input.UpdateParam.InterviewInformation, // AcceptancePoints,
		sendingCustomer.ID,                     // SendingCustomerID,
	)

	err = i.sendingJobSeekerRepository.Update(sendingJobSeeker, prevSendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 書類
	sendingJobSeekerDocument := entity.NewSendingJobSeekerDocument(
		prevSendingJobSeeker.ID,
		"", // ResumeOriginURL,
		input.UpdateParam.ResumePDFURL,
		"", // CVOriginURL,
		input.UpdateParam.CVPDFURL,
		"", // RecommendationOriginURL,
		"", // RecommendationPDFURL,
		"", // IDPhotoURL,
		"", // OtherDocument1URL,
		"", // OtherDocument2URL,
		"", // OtherDocument3URL,
	)

	err = i.sendingJobSeekerDocumentRepository.Update(sendingJobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 学歴
	studentHistory := entity.NewSendingJobSeekerStudentHistory(
		prevSendingJobSeeker.ID,
		NullInt, // SchoolCategory,
		input.UpdateParam.SchoolName,
		NullInt, // SchoolLevel,
		input.UpdateParam.Subject,
		input.UpdateParam.EntranceYear,
		input.UpdateParam.FirstStatus,
		input.UpdateParam.GraduationYear,
		input.UpdateParam.LastStatus,
	)

	err = i.sendingJobSeekerStudentHistoryRepository.Create(studentHistory)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 職歴
	workHistory := entity.NewSendingJobSeekerWorkHistory(
		prevSendingJobSeeker.ID,
		input.UpdateParam.CompanyName,
		NullInt, // EmployeeNumberSingle,
		NullInt, // EmployeeNumberGroup,
		NullInt, // PublicOffering,
		input.UpdateParam.JoiningYear,
		NullInt, // EmploymentStatus,
		"",      // RetireReasonOfTruth,
		"",      // RetireReasonOfPublic,
		input.UpdateParam.RetireYear,
		input.UpdateParam.FirstStatus,
		input.UpdateParam.LastStatus,
	)

	err = i.sendingJobSeekerWorkHistoryRepository.Create(workHistory)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	departmentHistory := entity.NewSendingJobSeekerDepartmentHistory(
		workHistory.ID,
		"",      // Department,
		NullInt, // ManagementNumber,
		"",      // ManagementDetail,
		input.UpdateParam.JobDescription,
		input.UpdateParam.JoiningYear, // StartYear,
		input.UpdateParam.RetireYear,  // EndYear,
	)

	err = i.sendingJobSeekerDepartmentHistoryRepository.Create(departmentHistory)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingJobSeeker.ID = prevSendingJobSeeker.ID
	output.SendingJobSeeker = sendingJobSeeker

	return output, nil
}

// 送客求職者の更新
type UpdateSendingJobSeekerInput struct {
	UpdateParam        entity.UpdateSendingJobSeekerParam
	SendingJobSeekerID uint
}

type UpdateSendingJobSeekerOutput struct {
	SendingJobSeeker *entity.SendingJobSeeker
}

func (i *SendingJobSeekerInteractorImpl) UpdateSendingJobSeeker(input UpdateSendingJobSeekerInput) (UpdateSendingJobSeekerOutput, error) {
	var (
		output UpdateSendingJobSeekerOutput
		err    error
	)

	sendingJobSeeker := entity.NewSendingJobSeeker(
		input.UpdateParam.AgentID, // 更新しない
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.LineID, // 更新しない
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
		input.UpdateParam.AppearanceDetailOfTruth,
		input.UpdateParam.AppearanceDetail,
		input.UpdateParam.Communication,
		input.UpdateParam.CommunicationDetailOfTruth,
		input.UpdateParam.CommunicationDetail,
		input.UpdateParam.Thinking,
		input.UpdateParam.ThinkingDetailOfTruth,
		input.UpdateParam.ThinkingDetail,
		input.UpdateParam.SecretMemo,
		input.UpdateParam.JobHuntingState,
		input.UpdateParam.RecommendReason,
		input.UpdateParam.Phase,
		input.UpdateParam.InterviewDate,
		input.UpdateParam.Agreement,
		input.UpdateParam.StudyCategory,
		input.UpdateParam.WordSkill,
		input.UpdateParam.ExcelSkill,
		input.UpdateParam.PowerPointSkill,
		input.UpdateParam.PublicMemo,
		input.UpdateParam.NationalityRemarks,
		input.UpdateParam.MedicalHistoryRemarks,
		input.UpdateParam.AcceptancePoints,
		input.UpdateParam.SendingCustomerID,
	)

	err = i.sendingJobSeekerRepository.Update(sendingJobSeeker, input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.sendingJobSeekerStudentHistoryRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sh := range input.UpdateParam.StudentHistories {
		studentHistory := entity.NewSendingJobSeekerStudentHistory(
			input.SendingJobSeekerID,
			sh.SchoolCategory,
			sh.SchoolName,
			sh.SchoolLevel,
			sh.Subject,
			sh.EntranceYear,
			sh.FirstStatus,
			sh.GraduationYear,
			sh.LastStatus,
		)

		err = i.sendingJobSeekerStudentHistoryRepository.Create(studentHistory)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerWorkHistoryRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, wh := range input.UpdateParam.WorkHistories {
		workHistory := entity.NewSendingJobSeekerWorkHistory(
			input.SendingJobSeekerID,
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

		err = i.sendingJobSeekerWorkHistoryRepository.Create(workHistory)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, ei := range wh.ExperienceIndustries {
			experienceIndustry := entity.NewSendingJobSeekerExperienceIndustry(
				workHistory.ID,
				ei.Industry,
			)

			err = i.sendingJobSeekerExperienceIndustryRepository.Create(experienceIndustry)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		for _, dh := range wh.DepartmentHistories {
			departmentHistory := entity.NewSendingJobSeekerDepartmentHistory(
				workHistory.ID,
				dh.Department,
				dh.ManagementNumber,
				dh.ManagementDetail,
				dh.JobDescription,
				dh.StartYear,
				dh.EndYear,
			)

			err = i.sendingJobSeekerDepartmentHistoryRepository.Create(departmentHistory)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			for _, occupations := range dh.ExperienceOccupations {
				occupation := entity.NewSendingJobSeekerExperienceOccupation(
					departmentHistory.ID,
					occupations.Occupation,
				)

				err = i.sendingJobSeekerExperienceOccupationRepository.Create(occupation)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}

		}

	}

	err = i.sendingJobSeekerDesiredCompanyScaleRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, dcs := range input.UpdateParam.DesiredCompanyScales {
		disiredCompanyScale := entity.NewSendingJobSeekerDesiredCompanyScale(
			input.SendingJobSeekerID,
			dcs.DesiredCompanyScale,
		)

		err = i.sendingJobSeekerDesiredCompanyScaleRepository.Create(disiredCompanyScale)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerLicenseRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, l := range input.UpdateParam.Licenses {
		license := entity.NewSendingJobSeekerLicense(
			input.SendingJobSeekerID,
			l.LicenseType,
			l.AcquisitionTime,
		)

		err = i.sendingJobSeekerLicenseRepository.Create(license)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerSelfPromotionRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sp := range input.UpdateParam.SelfPromotions {
		selfPromotion := entity.NewSendingJobSeekerSelfPromotion(
			input.SendingJobSeekerID,
			sp.Title,
			sp.Contents,
		)

		err = i.sendingJobSeekerSelfPromotionRepository.Create(selfPromotion)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerDesiredIndustryRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, di := range input.UpdateParam.DesiredIndustries {
		desiredIndustry := entity.NewSendingJobSeekerDesiredIndustry(
			input.SendingJobSeekerID,
			di.DesiredIndustry,
			di.DesiredRank,
		)

		err = i.sendingJobSeekerDesiredIndustryRepository.Create(desiredIndustry)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerDesiredOccupationRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, do := range input.UpdateParam.DesiredOccupations {
		desiredOccupation := entity.NewSendingJobSeekerDesiredOccupation(
			input.SendingJobSeekerID,
			do.DesiredOccupation,
			do.DesiredRank,
		)

		err = i.sendingJobSeekerDesiredOccupationRepository.Create(desiredOccupation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerDesiredWorkLocationRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, dwl := range input.UpdateParam.DesiredWorkLocations {
		desiredWorkLocation := entity.NewSendingJobSeekerDesiredWorkLocation(
			input.SendingJobSeekerID,
			dwl.DesiredWorkLocation,
			dwl.DesiredRank,
		)

		err = i.sendingJobSeekerDesiredWorkLocationRepository.Create(desiredWorkLocation)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerDesiredHolidayTypeRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, dht := range input.UpdateParam.DesiredHolidayTypes {
		desiredHolidayType := entity.NewSendingJobSeekerDesiredHolidayType(
			input.SendingJobSeekerID,
			dht.HolidayType,
		)

		err = i.sendingJobSeekerDesiredHolidayTypeRepository.Create(desiredHolidayType)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerDevelopmentSkillRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ds := range input.UpdateParam.DevelopmentSkills {
		developmentSkill := entity.NewSendingJobSeekerDevelopmentSkill(
			input.SendingJobSeekerID,
			ds.DevelopmentCategory,
			ds.DevelopmentType,
			ds.ExperienceYear,
			ds.ExperienceMonth,
		)

		err = i.sendingJobSeekerDevelopmentSkillRepository.Create(developmentSkill)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerLanguageSkillRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ls := range input.UpdateParam.LanguageSkills {
		languageSkill := entity.NewSendingJobSeekerLanguageSkill(
			input.SendingJobSeekerID,
			ls.LanguageType,
			ls.LanguageLevel,
			ls.Toeic,
			ls.ToeicExaminationYear,
			ls.ToeflIBT,
			ls.ToeflIBTExaminationYear,
			ls.ToeflPBT,
			ls.ToeflPBTExaminationYear,
		)

		err = i.sendingJobSeekerLanguageSkillRepository.Create(languageSkill)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.sendingJobSeekerPCToolRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ps := range input.UpdateParam.PCTools {
		PCTool := entity.NewSendingJobSeekerPCTool(
			input.SendingJobSeekerID,
			ps.Tool,
		)

		err = i.sendingJobSeekerPCToolRepository.Create(PCTool)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	sendingJobSeeker.ID = input.SendingJobSeekerID
	sendingJobSeeker.UUID = input.UpdateParam.UUID
	sendingJobSeeker.StudentHistories = input.UpdateParam.StudentHistories
	sendingJobSeeker.WorkHistories = input.UpdateParam.WorkHistories
	sendingJobSeeker.DesiredCompanyScales = input.UpdateParam.DesiredCompanyScales
	sendingJobSeeker.Licenses = input.UpdateParam.Licenses
	sendingJobSeeker.SelfPromotions = input.UpdateParam.SelfPromotions
	sendingJobSeeker.DesiredIndustries = input.UpdateParam.DesiredIndustries
	sendingJobSeeker.DesiredOccupations = input.UpdateParam.DesiredOccupations
	sendingJobSeeker.DesiredWorkLocations = input.UpdateParam.DesiredWorkLocations
	sendingJobSeeker.DesiredHolidayTypes = input.UpdateParam.DesiredHolidayTypes
	sendingJobSeeker.DevelopmentSkills = input.UpdateParam.DevelopmentSkills
	sendingJobSeeker.LanguageSkills = input.UpdateParam.LanguageSkills
	sendingJobSeeker.PCTools = input.UpdateParam.PCTools

	output.SendingJobSeeker = sendingJobSeeker

	return output, nil
}

// 送客求職者のフェーズの更新
type UpdateSendingJobSeekerPhaseInput struct {
	Phase              uint
	SendingJobSeekerID uint
}

type UpdateSendingJobSeekerPhaseOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) UpdateSendingJobSeekerPhase(input UpdateSendingJobSeekerPhaseInput) (UpdateSendingJobSeekerPhaseOutput, error) {
	var (
		output UpdateSendingJobSeekerPhaseOutput
		err    error
	)

	err = i.sendingJobSeekerRepository.UpdatePhase(input.SendingJobSeekerID, null.NewInt(int64(input.Phase), true))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.sendingCustomerRepository.UpdatePhase(input.SendingJobSeekerID, null.NewInt(int64(input.Phase), true))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 送客求職者の面談日時の更新
type UpdateSendingJobSeekerInterviewDateInput struct {
	InterviewDate      time.Time
	SendingJobSeekerID uint
}

type UpdateSendingJobSeekerInterviewDateOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) UpdateSendingInterviewDateBySendingJobSeekerID(input UpdateSendingJobSeekerInterviewDateInput) (UpdateSendingJobSeekerInterviewDateOutput, error) {
	var (
		output UpdateSendingJobSeekerInterviewDateOutput
		err    error
	)

	// 管理側の送客求職者情報を取得
	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 管理者側の更新
	err = i.sendingJobSeekerRepository.UpdateInterviewDateByCustomerID(sendingJobSeeker.SendingCustomerID, input.InterviewDate)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客元の情報を更新
	err = i.sendingCustomerRepository.UpdateInterviewDate(sendingJobSeeker.SendingCustomerID, input.InterviewDate)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// アクティビティーメモを更新
type UpdateSendingJobSeekerActivityMemoInput struct {
	SendingJobSeekerID uint
	ActivityMemo       string
}

type UpdateSendingJobSeekerActivityMemoOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) UpdateSendingJobSeekerActivityMemo(input UpdateSendingJobSeekerActivityMemoInput) (UpdateSendingJobSeekerActivityMemoOutput, error) {
	var (
		output UpdateSendingJobSeekerActivityMemoOutput
	)

	err := i.sendingJobSeekerRepository.UpdateActivityMemo(input.SendingJobSeekerID, input.ActivityMemo)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 面談実施待ちの未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
type UpdateIsVewForWatingInput struct {
	SendingJobSeekerID uint
}

type UpdateIsVewForWatingOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) UpdateIsVewForWating(input UpdateIsVewForWatingInput) (UpdateIsVewForWatingOutput, error) {
	var (
		output UpdateIsVewForWatingOutput
	)

	// 閲覧済み（false）に更新
	err := i.sendingJobSeekerIsViewRepository.UpdateIsNotWaitingViewedBySendingJobSeekerID(input.SendingJobSeekerID, false)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 未登録の未閲覧の更新（送客進捗管理上で未閲覧クリックした時に実行）
type UpdateIsVewForUnregisterInput struct {
	SendingJobSeekerID uint
}

type UpdateIsVewForUnregisterOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) UpdateIsVewForUnregister(input UpdateIsVewForUnregisterInput) (UpdateIsVewForUnregisterOutput, error) {
	var (
		output UpdateIsVewForUnregisterOutput
	)

	// 閲覧済み（false）に更新
	err := i.sendingJobSeekerIsViewRepository.UpdateIsNotUnregisterViewedBySendingJobSeekerID(input.SendingJobSeekerID, false)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 送客求職者の削除
type DeleteSendingJobSeekerInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeeker(input DeleteSendingJobSeekerInput) (DeleteSendingJobSeekerOutput, error) {
	var (
		output DeleteSendingJobSeekerOutput
	)

	err := i.sendingJobSeekerRepository.Delete(input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// 送客求職者IDを使って送客求職者情報を取得する
type GetSendingJobSeekerByIDInput struct {
	SendingJobSeekerID uint
}

type GetSendingJobSeekerByIDOutput struct {
	SendingJobSeeker *entity.SendingJobSeeker
}

func (i *SendingJobSeekerInteractorImpl) GetSendingJobSeekerByID(input GetSendingJobSeekerByIDInput) (GetSendingJobSeekerByIDOutput, error) {
	var (
		output GetSendingJobSeekerByIDOutput
		err    error
	)

	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客求職者の子テーブル情報をセット
	sendingJobSeeker, err = i.setSendingJobSeekerChildTable(sendingJobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingJobSeeker = sendingJobSeeker

	return output, nil
}

// 送客求職者IDを使って送客求職者情報を取得する
type GetSendingJobSeekerByUUIDInput struct {
	SendingJobSeekerUUID uuid.UUID
}

type GetSendingJobSeekerByUUIDOutput struct {
	SendingJobSeeker *entity.SendingJobSeeker
}

func (i *SendingJobSeekerInteractorImpl) GetSendingJobSeekerByUUID(input GetSendingJobSeekerByUUIDInput) (GetSendingJobSeekerByUUIDOutput, error) {
	var (
		output GetSendingJobSeekerByUUIDOutput
		err    error
	)

	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByUUID(input.SendingJobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客求職者の子テーブル情報をセット
	sendingJobSeeker, err = i.setSendingJobSeekerChildTable(sendingJobSeeker)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingJobSeeker = sendingJobSeeker

	return output, nil
}

// 送客求職者IDを使って送客求職者情報を取得する
type GetIsNotViewSendingJobSeekerCountByAgentStaffIDInput struct {
	AgentStaffID uint
}

type GetIsNotViewSendingJobSeekerCountByAgentStaffIDOutput struct {
	IsViewCount *entity.IsViewCount
}

func (i *SendingJobSeekerInteractorImpl) GetIsNotViewSendingJobSeekerCountByAgentStaffID(input GetIsNotViewSendingJobSeekerCountByAgentStaffIDInput) (GetIsNotViewSendingJobSeekerCountByAgentStaffIDOutput, error) {
	var (
		output GetIsNotViewSendingJobSeekerCountByAgentStaffIDOutput
		err    error
	)

	waitingCount, err := i.sendingJobSeekerIsViewRepository.GetIsNotWaitingViewCountByAgentStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	unregisterCount, err := i.sendingJobSeekerIsViewRepository.GetIsNotUnregisterViewCountByAgentStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.IsViewCount = &entity.IsViewCount{
		WaitingCount:    waitingCount,
		UnregisterCount: unregisterCount,
	}

	return output, nil
}

// 送客求職者IDを使って送客求職者情報を取得する
type GetSendingJobSeekerDocumentByUUIDInput struct {
	SendingJobSeekerUUID uuid.UUID
}

type GetSendingJobSeekerDocumentByUUIDOutput struct {
	Document *entity.SendingJobSeekerDocument
}

func (i *SendingJobSeekerInteractorImpl) GetSendingJobSeekerDocumentByUUID(input GetSendingJobSeekerDocumentByUUIDInput) (GetSendingJobSeekerDocumentByUUIDOutput, error) {
	var (
		output GetSendingJobSeekerDocumentByUUIDOutput
		err    error
	)

	document, err := i.sendingJobSeekerDocumentRepository.FindBySendingJobSeekerUUID(input.SendingJobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Document = document

	return output, nil
}

/****************************************************************************************/
// { value: number; label: string }[]の形式でリストを取得するapi
//

// 送客側の進捗一覧画面の絞り込み検索に使用する担当者, 送客元, 送客先のリストを取得する
type GetSearchListForSendingJobSeekerManagementByAgentIDInput struct {
	AgentID uint
}

type GetSearchListForSendingJobSeekerManagementByAgentIDOutput struct {
	StaffList     []*entity.LabelAndValue // 担当者のリスト
	SendAgentList []*entity.LabelAndValue // 送客先
	SenderList    []*entity.LabelAndValue // 送客元
}

func (i *SendingJobSeekerInteractorImpl) GetSearchListForSendingJobSeekerManagementByAgentID(input GetSearchListForSendingJobSeekerManagementByAgentIDInput) (GetSearchListForSendingJobSeekerManagementByAgentIDOutput, error) {
	var (
		output GetSearchListForSendingJobSeekerManagementByAgentIDOutput
		err    error

		// アンドイーズアカウントのチェック
		IsAndes = input.AgentID == 1

		staffList     []*entity.AgentStaff
		sendAgentList []*entity.SendingEnterprise
		senderList    []*entity.Agent
	)

	/************ 1. 担当者, 送客先, 送客元のリスト **************/

	// 担当者のリスト
	staffList, err = i.agentStaffRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客先のリスト
	sendAgentList, err = i.sendingEnterpriseRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客元のリスト（送客利用アカウントのエージェントを取得）
	if IsAndes {
		// アンドイーズの場合は送客利用しているすべてのエージェントを取得
		senderList, err = i.agentRepository.GetSendingActive()
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	} else {
		// アンドイーズ以外のagentsテーブルのsendig_typeが1になっているエージェントは場合は自社のみ
		agent, err := i.agentRepository.FindByID(input.AgentID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		senderList = append(senderList, agent)
	}

	/************ 2. 取得したリストをLabelAndValueに変換 **************/

	// 担当者
	for _, staff := range staffList {
		staffLabel := entity.NewLabelAndValue(staff.ID, staff.StaffName)
		output.StaffList = append(output.StaffList, staffLabel)
	}

	// 送客先
	for _, sendAgent := range sendAgentList {
		sendAgentLabel := entity.NewLabelAndValue(sendAgent.ID, sendAgent.CompanyName)
		output.SendAgentList = append(output.SendAgentList, sendAgentLabel)
	}

	// 送客元
	for _, sender := range senderList {
		senderLabel := entity.NewLabelAndValue(sender.ID, sender.AgentName)
		output.SenderList = append(output.SenderList, senderLabel)
	}

	return output, nil
}

/****************************************************************************************/
// 送客求職者資料関連API
//
type UpdateSendingJobSeekerDocumentInput struct {
	UpdateParam entity.CreateOrUpdateSendingJobSeekerDocumentParam
}

type UpdateSendingJobSeekerDocumentOutput struct {
	SendingJobSeekerDocument *entity.SendingJobSeekerDocument
}

func (i *SendingJobSeekerInteractorImpl) UpdateSendingJobSeekerDocument(input UpdateSendingJobSeekerDocumentInput) (UpdateSendingJobSeekerDocumentOutput, error) {
	var output UpdateSendingJobSeekerDocumentOutput

	sendingJobSeekerDocument := entity.NewSendingJobSeekerDocument(
		uint(input.UpdateParam.SendingJobSeekerID.Int64),
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

	err := i.sendingJobSeekerDocumentRepository.Update(sendingJobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByID(uint(input.UpdateParam.SendingJobSeekerID.Int64))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// フェーズが詳細未登録の場合(送客画面のため)、customersの資料も更新する
	if sendingJobSeeker.Phase == null.NewInt(int64(entity.WaitingForInterview), true) {
		err := i.sendingCustomerRepository.UpdateDocument(
			sendingJobSeeker.SendingCustomerID,
			input.UpdateParam.ResumePDFURL,
			input.UpdateParam.ResumeOriginURL,
			input.UpdateParam.CVPDFURL,
			input.UpdateParam.CVOriginURL,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.SendingJobSeekerDocument = sendingJobSeekerDocument

	return output, nil
}

type GetSendingJobSeekerDocumentBySendingJobSeekerIDInput struct {
	SendingJobSeekerID uint
}

type GetSendingJobSeekerDocumentBySendingJobSeekerIDOutput struct {
	SendingJobSeekerDocument *entity.SendingJobSeekerDocument
}

func (i *SendingJobSeekerInteractorImpl) GetSendingJobSeekerDocumentBySendingJobSeekerID(input GetSendingJobSeekerDocumentBySendingJobSeekerIDInput) (GetSendingJobSeekerDocumentBySendingJobSeekerIDOutput, error) {
	var output GetSendingJobSeekerDocumentBySendingJobSeekerIDOutput

	sendingJobSeekerDocument, err := i.sendingJobSeekerDocumentRepository.FindBySendingJobSeekerID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingJobSeekerDocument = sendingJobSeekerDocument

	return output, nil
}

type DeleteSendingJobSeekerResumePDFURLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerResumePDFURLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerResumePDFURL(input DeleteSendingJobSeekerResumePDFURLInput) (DeleteSendingJobSeekerResumePDFURLOutput, error) {
	var (
		output DeleteSendingJobSeekerResumePDFURLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateResumePDFURL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerResumeOriginURLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerResumeOriginURLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerResumeOriginURL(input DeleteSendingJobSeekerResumeOriginURLInput) (DeleteSendingJobSeekerResumeOriginURLOutput, error) {
	var (
		output DeleteSendingJobSeekerResumeOriginURLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateResumeOriginURL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerCVPDFURLInput struct {
	SendingJobSeekerID uint
}
type DeleteSendingJobSeekerCVPDFURLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerCVPDFURL(input DeleteSendingJobSeekerCVPDFURLInput) (DeleteSendingJobSeekerCVPDFURLOutput, error) {
	var (
		output DeleteSendingJobSeekerCVPDFURLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateCVPDFURL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerCVOriginURLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerCVOriginURLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerCVOriginURL(input DeleteSendingJobSeekerCVOriginURLInput) (DeleteSendingJobSeekerCVOriginURLOutput, error) {
	var (
		output DeleteSendingJobSeekerCVOriginURLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateCVOriginURL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerRecommendationPDFURLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerRecommendationPDFURLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerRecommendationPDFURL(input DeleteSendingJobSeekerRecommendationPDFURLInput) (DeleteSendingJobSeekerRecommendationPDFURLOutput, error) {
	var (
		output DeleteSendingJobSeekerRecommendationPDFURLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateRecommendationPDFURL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerRecommendationOriginURLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerRecommendationOriginURLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerRecommendationOriginURL(input DeleteSendingJobSeekerRecommendationOriginURLInput) (DeleteSendingJobSeekerRecommendationOriginURLOutput, error) {
	var (
		output DeleteSendingJobSeekerRecommendationOriginURLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateRecommendationOriginURL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerIDPhotoURLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerIDPhotoURLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerIDPhotoURL(input DeleteSendingJobSeekerIDPhotoURLInput) (DeleteSendingJobSeekerIDPhotoURLOutput, error) {
	var (
		output DeleteSendingJobSeekerIDPhotoURLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateIDPhotoURL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerOtherDocument1URLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerOtherDocument1URLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerOtherDocument1URL(input DeleteSendingJobSeekerOtherDocument1URLInput) (DeleteSendingJobSeekerOtherDocument1URLOutput, error) {
	var (
		output DeleteSendingJobSeekerOtherDocument1URLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateOtherDocument1URL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerOtherDocument2URLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerOtherDocument2URLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerOtherDocument2URL(input DeleteSendingJobSeekerOtherDocument2URLInput) (DeleteSendingJobSeekerOtherDocument2URLOutput, error) {
	var (
		output DeleteSendingJobSeekerOtherDocument2URLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateOtherDocument2URL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type DeleteSendingJobSeekerOtherDocument3URLInput struct {
	SendingJobSeekerID uint
}

type DeleteSendingJobSeekerOtherDocument3URLOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) DeleteSendingJobSeekerOtherDocument3URL(input DeleteSendingJobSeekerOtherDocument3URLInput) (DeleteSendingJobSeekerOtherDocument3URLOutput, error) {
	var (
		output DeleteSendingJobSeekerOtherDocument3URLOutput
	)

	err := i.sendingJobSeekerDocumentRepository.UpdateOtherDocument3URL("", input.SendingJobSeekerID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
// 面談前アンケート関連 API
//
type CreateSendingInitialQuestionnaireInput struct {
	Param entity.CreateSendingInitialQuestionnaireParam
}

type CreateSendingInitialQuestionnairdhutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) CreateSendingInitialQuestionnaire(input CreateSendingInitialQuestionnaireInput) (CreateSendingInitialQuestionnairdhutput, error) {
	var (
		output CreateSendingInitialQuestionnairdhutput
		err    error
		param  = input.Param
	)

	/************ 1. 求職者の同意項目アップデート **************/

	// アンケートの質問・要望を更新
	err = i.sendingJobSeekerRepository.UpdateQuestion(param.SendingJobSeekerID, param.Question)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者テーブルの個人情報同意の有無を更新
	err = i.sendingJobSeekerRepository.UpdateAgreement(
		param.SendingJobSeekerID,
		true, // input.Param.Agreement,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. 応募書類の更新 **************/

	sendingJobSeekerDocument := entity.NewSendingJobSeekerDocument(
		param.SendingJobSeekerID,
		param.Documents.ResumeOriginURL,         // ResumeOriginURL
		param.Documents.ResumePDFURL,            // ResumePDFURL
		param.Documents.CVOriginURL,             // CVOriginURL
		param.Documents.CVPDFURL,                // CVPDFURL
		param.Documents.RecommendationOriginURL, // RecommendationOriginURL
		param.Documents.RecommendationPDFURL,    // RecommendationPDFURL
		param.Documents.IDPhotoURL,              // IDPhotoURL
		param.Documents.OtherDocument1URL,       // OtherDocument1URL
		param.Documents.OtherDocument2URL,       // OtherDocument2URL
		param.Documents.OtherDocument3URL,       // OtherDocument3URL
	)

	// 応募書類の更新
	err = i.sendingJobSeekerDocumentRepository.Update(sendingJobSeekerDocument)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 3. 希望項目の更新 **************/

	// 希望業界の入力がある場合は更新
	if len(param.DesiredIndustries) > 0 {
		err = i.sendingJobSeekerDesiredIndustryRepository.Delete(param.SendingJobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		err = i.sendingJobSeekerDesiredIndustryRepository.CreateMulti(param.SendingJobSeekerID, param.DesiredIndustries)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 希望職種の入力がある場合は更新
	if len(param.DesiredOccupations) > 0 {
		err = i.sendingJobSeekerDesiredOccupationRepository.Delete(param.SendingJobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		err = i.sendingJobSeekerDesiredOccupationRepository.CreateMulti(param.SendingJobSeekerID, param.DesiredOccupations)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 希望勤務地の入力がある場合は更新
	if len(param.DesiredWorkLocations) > 0 {
		err = i.sendingJobSeekerDesiredWorkLocationRepository.Delete(param.SendingJobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		err = i.sendingJobSeekerDesiredWorkLocationRepository.CreateMulti(param.SendingJobSeekerID, param.DesiredWorkLocations)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
/// 送客理由 API
//
type CreateSendingJobSeekerEndStatusInput struct {
	CreateParam entity.CreateSendingJobSeekerEndStatusParam
}

type CreateSendingJobSeekerEndStatusOutput struct {
	OK bool
}

func (i *SendingJobSeekerInteractorImpl) CreateSendingJobSeekerEndStatus(input CreateSendingJobSeekerEndStatusInput) (CreateSendingJobSeekerEndStatusOutput, error) {
	var (
		output CreateSendingJobSeekerEndStatusOutput
		err    error
	)

	/********* 対象求職者のフェーズを終了に更新 *********/

	// 管理側のフェーズを終了に更新
	err = i.sendingJobSeekerRepository.UpdatePhase(input.CreateParam.SendingJobSeekerID, null.NewInt(int64(entity.CloseSending), true))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送客元のフェーズを終了に更新
	err = i.sendingCustomerRepository.UpdatePhaseBySendingJobSeekerID(input.CreateParam.SendingJobSeekerID, null.NewInt(int64(entity.CloseSending), true))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/********* 終了理由の登録 *********/

	sendingJobSeekerEndStatus := entity.NewSendingJobSeekerEndStatus(
		input.CreateParam.SendingJobSeekerID,
		input.CreateParam.EndReason,
		input.CreateParam.EndStatus,
	)

	err = i.sendingJobSeekerEndStatusRepository.Create(sendingJobSeekerEndStatus)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
/****************************************************************************************/
// 送客求職者の子テーブル情報をセットする関数
//
func (i *SendingJobSeekerInteractorImpl) setSendingJobSeekerChildTable(sendingJobSeeker *entity.SendingJobSeeker) (*entity.SendingJobSeeker, error) {
	studentHistory, err := i.sendingJobSeekerStudentHistoryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	workHistory, err := i.sendingJobSeekerWorkHistoryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	experienceIndustry, err := i.sendingJobSeekerExperienceIndustryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	departmentHistory, err := i.sendingJobSeekerDepartmentHistoryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	experienceOccupation, err := i.sendingJobSeekerExperienceOccupationRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	desiredCompanyScale, err := i.sendingJobSeekerDesiredCompanyScaleRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	license, err := i.sendingJobSeekerLicenseRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	selfPromotion, err := i.sendingJobSeekerSelfPromotionRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	document, err := i.sendingJobSeekerDocumentRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	desiredIndustry, err := i.sendingJobSeekerDesiredIndustryRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	desiredOccupation, err := i.sendingJobSeekerDesiredOccupationRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	desiredWorkLocation, err := i.sendingJobSeekerDesiredWorkLocationRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	desiredHolidayType, err := i.sendingJobSeekerDesiredHolidayTypeRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	developmentSkill, err := i.sendingJobSeekerDevelopmentSkillRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	languageSkill, err := i.sendingJobSeekerLanguageSkillRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
	}

	pcSkill, err := i.sendingJobSeekerPCToolRepository.FindBySendingJobSeekerID(sendingJobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return sendingJobSeeker, err
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

	return sendingJobSeeker, nil
}
