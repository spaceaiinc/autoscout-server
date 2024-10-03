package interactor

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
	// ブラウザ操作
	// csvの文字化け変換
)

type ScoutServiceInteractor interface {
	// 汎用系 API
	CreateScoutService(input CreateScoutServiceInput) (CreateScoutServiceByIDOutput, error)
	UpdateScoutService(input UpdateScoutServiceInput) (UpdateScoutServiceByIDOutput, error)
	UpdateScoutServicePassword(input UpdateScoutServicePasswordInput) (UpdateScoutServicePasswordByIDOutput, error)
	DeleteScoutService(input DeleteScoutServiceInput) (DeleteScoutServiceByIDOutput, error)
	GetByID(input ScoutServiceGetByIDInput) (ScoutServiceGetByIDOutput, error)
	GetListByAgentID(input GetListByAgentIDInput) (GetListByAgentIDOutput, error)

	// Batch処理用 API
	BatchScout(input BatchScoutInput) (BatchScoutOutput, error)
	BatchEntry(input BatchEntryInput) (BatchEntryOutput, error)

	// エントリー求職者取得 API
	EntryOnRan(input EntryOnRanInput) (EntryOnRanOutput, error)
	EntryOnMynaviScouting(input EntryOnMynaviScoutingInput) (EntryOnMynaviScoutingOutput, error)
	EntryOnAmbi(input EntryOnAmbiInput) (EntryOnAmbiOutput, error)
	EntryOnMynaviAgentScout(input EntryOnMynaviAgentScoutInput) (EntryOnMynaviAgentScoutOutput, error)

	// スカウトメール送信 API
	ScoutOnRan(input ScoutOnRanInput) (ScoutOnRanOutput, error)
	ScoutOnMynaviScouting(input ScoutOnMynaviScoutingInput) (ScoutOnMynaviScoutingOutput, error)
	ScoutOnAmbi(input ScoutOnAmbiInput) (ScoutOnAmbiOutput, error)
	ScoutOnMynaviAgentScout(input ScoutOnMynaviAgentScoutInput) (ScoutOnMynaviAgentScoutOutput, error)

	// Gmail API
	GmailWebHook(input GmailWebHookInput) (GmailWebHookOutput, error)
}

type ScoutServiceInteractorImpl struct {
	firebase                                usecase.Firebase
	sendgrid                                config.Sendgrid
	oneSignal                               config.OneSignal
	app                                     config.App
	googleAPI                               config.GoogleAPI
	slack                                   config.Slack
	scoutServiceRepository                  usecase.ScoutServiceRepository
	scoutServiceTemplateRepository          usecase.ScoutServiceTemplateRepository
	scoutServiceGetEntryTimeRepository      usecase.ScoutServiceGetEntryTimeRepository
	agentRobotRepository                    usecase.AgentRobotRepository
	agentStaffRepository                    usecase.AgentStaffRepository
	agentRepository                         usecase.AgentRepository
	jobSeekerRepository                     usecase.JobSeekerRepository
	jobSeekerStudentHistoryRepository       usecase.JobSeekerStudentHistoryRepository
	jobSeekerWorkHistoryRepository          usecase.JobSeekerWorkHistoryRepository
	jobSeekerExperienceIndustryRepository   usecase.JobSeekerExperienceIndustryRepository
	jobSeekerDepartmentHistoryRepository    usecase.JobSeekerDepartmentHistoryRepository
	jobSeekerLicenseRepository              usecase.JobSeekerLicenseRepository
	jobSeekerSelfPromotionRepository        usecase.JobSeekerSelfPromotionRepository
	jobSeekerDocumentRepository             usecase.JobSeekerDocumentRepository
	jobSeekerDesiredIndustryRepository      usecase.JobSeekerDesiredIndustryRepository
	jobSeekerDesiredOccupationRepository    usecase.JobSeekerDesiredOccupationRepository
	jobSeekerDesiredWorkLocationRepository  usecase.JobSeekerDesiredWorkLocationRepository
	jobSeekerDesiredHolidayTypeRepository   usecase.JobSeekerDesiredHolidayTypeRepository
	jobSeekerDevelopmentSkillRepository     usecase.JobSeekerDevelopmentSkillRepository
	jobSeekerLanguageSkillRepository        usecase.JobSeekerLanguageSkillRepository
	jobSeekerPCToolRepository               usecase.JobSeekerPCToolRepository
	jobSeekerHideToAgentRepository          usecase.JobSeekerHideToAgentRepository
	jobSeekerExperienceOccupationRepository usecase.JobSeekerExperienceOccupationRepository
	jobSeekerDesiredCompanyScaleRepository  usecase.JobSeekerDesiredCompanyScaleRepository
	chatGroupWithJobSeekerRepository        usecase.ChatGroupWithJobSeekerRepository
	interviewTaskGroupRepository            usecase.InterviewTaskGroupRepository
	interviewTaskRepository                 usecase.InterviewTaskRepository
	interviewAdjustmentTemplateRepository   usecase.InterviewAdjustmentTemplateRepository
	googleAuthenticationRepository          usecase.GoogleAuthenticationRepository
	emailWithJobSeekerRepository            usecase.EmailWithJobSeekerRepository
	userEntryRepository                     usecase.UserEntryRepository
}

// ScoutServiceInteractorImpl is an implementation of ScoutServiceInteractor
func NewScoutServiceInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	ap config.App,
	ga config.GoogleAPI,
	sl config.Slack,
	ssR usecase.ScoutServiceRepository,
	stR usecase.ScoutServiceTemplateRepository,
	ssgetR usecase.ScoutServiceGetEntryTimeRepository,
	arR usecase.AgentRobotRepository,
	asR usecase.AgentStaffRepository,
	aR usecase.AgentRepository,
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
	cgwjsR usecase.ChatGroupWithJobSeekerRepository,
	itgR usecase.InterviewTaskGroupRepository,
	itR usecase.InterviewTaskRepository,
	iatR usecase.InterviewAdjustmentTemplateRepository,
	gaR usecase.GoogleAuthenticationRepository,
	ewjsR usecase.EmailWithJobSeekerRepository,
	ueR usecase.UserEntryRepository,
) ScoutServiceInteractor {
	return &ScoutServiceInteractorImpl{
		firebase:                                fb,
		sendgrid:                                sg,
		oneSignal:                               os,
		app:                                     ap,
		googleAPI:                               ga,
		slack:                                   sl,
		scoutServiceRepository:                  ssR,
		scoutServiceTemplateRepository:          stR,
		scoutServiceGetEntryTimeRepository:      ssgetR,
		agentRobotRepository:                    arR,
		agentStaffRepository:                    asR,
		agentRepository:                         aR,
		jobSeekerRepository:                     jsR,
		jobSeekerStudentHistoryRepository:       jsshR,
		jobSeekerWorkHistoryRepository:          jswhR,
		jobSeekerExperienceIndustryRepository:   jseiR,
		jobSeekerExperienceOccupationRepository: jseoR,
		jobSeekerLicenseRepository:              jslR,
		jobSeekerSelfPromotionRepository:        jsspR,
		jobSeekerDocumentRepository:             jsdR,
		jobSeekerDesiredIndustryRepository:      jsdiR,
		jobSeekerDesiredOccupationRepository:    jsdoR,
		jobSeekerDesiredWorkLocationRepository:  jsdwlR,
		jobSeekerDesiredHolidayTypeRepository:   jsdhtR,
		jobSeekerDevelopmentSkillRepository:     jsdsR,
		jobSeekerLanguageSkillRepository:        jslsR,
		jobSeekerPCToolRepository:               jsptR,
		jobSeekerHideToAgentRepository:          jshR,
		jobSeekerDepartmentHistoryRepository:    jsdhR,
		jobSeekerDesiredCompanyScaleRepository:  jsdcsR,
		chatGroupWithJobSeekerRepository:        cgwjsR,
		interviewTaskGroupRepository:            itgR,
		interviewTaskRepository:                 itR,
		interviewAdjustmentTemplateRepository:   iatR,
		googleAuthenticationRepository:          gaR,
		emailWithJobSeekerRepository:            ewjsR,
		userEntryRepository:                     ueR,
	}
}

/****************************************************************************************/
// 汎用系API
//
// スカウトサービスを作成する
type CreateScoutServiceInput struct {
	CreateParam entity.CreateOrUpdateScoutServiceParam
}

type CreateScoutServiceByIDOutput struct {
	ScoutService *entity.ScoutService
}

func (i *ScoutServiceInteractorImpl) CreateScoutService(input CreateScoutServiceInput) (CreateScoutServiceByIDOutput, error) {
	var (
		output       CreateScoutServiceByIDOutput
		scoutService *entity.ScoutService
		err          error
	)

	// パスワードの暗号化
	input.CreateParam.Password, err = encrypt(input.CreateParam.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	scoutService = entity.NewScoutService(
		input.CreateParam.AgentRobotID,
		input.CreateParam.AgentStaffID,
		input.CreateParam.LoginID,
		input.CreateParam.Password,
		input.CreateParam.ServiceType,
		input.CreateParam.IsActive,
		input.CreateParam.Memo,
		input.CreateParam.TemplateTitleForEmployed,
		input.CreateParam.TemplateTitleForUnemployed,
		input.CreateParam.InterviewAdjustmentTemplateID,
		input.CreateParam.InflowChannelID,
	)

	// スカウトサービスを作成
	err = i.scoutServiceRepository.Create(scoutService)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// エントリー取得時間を作成
	for _, getEntryTime := range input.CreateParam.GetEntryTimes {
		newGetEntryTime := entity.NewScoutServiceGetEntryTime(
			scoutService.ID,
			getEntryTime.StartHour,
			getEntryTime.StartMinute,
		)

		err = i.scoutServiceGetEntryTimeRepository.Create(newGetEntryTime)
		if err != nil {
			log.Println(err)
			return output, err
		}
	}

	// スカウトテンプレートを作成
	for _, scoutServiceTemplate := range input.CreateParam.Templates {
		newScoutServiceTemplate := entity.NewScoutServiceTemplate(
			scoutService.ID,
			scoutServiceTemplate.StartHour,
			scoutServiceTemplate.StartMinute,
			scoutServiceTemplate.RunOnMonday,
			scoutServiceTemplate.RunOnTuesday,
			scoutServiceTemplate.RunOnWednesday,
			scoutServiceTemplate.RunOnThursday,
			scoutServiceTemplate.RunOnFriday,
			scoutServiceTemplate.RunOnSaturday,
			scoutServiceTemplate.RunOnSunday,
			scoutServiceTemplate.ScoutCount,
			scoutServiceTemplate.SearchTitle,
			scoutServiceTemplate.MessageTitle,
			scoutServiceTemplate.JobInformationTitle,
			scoutServiceTemplate.JobInformationID,
			scoutServiceTemplate.AgeLimit,
			scoutServiceTemplate.ScoutType,
			scoutServiceTemplate.AutoRemind,
			scoutServiceTemplate.ReplyLimit,
		)

		// テンプレートを作成
		err = i.scoutServiceTemplateRepository.Create(newScoutServiceTemplate)
		if err != nil {
			return output, err
		}
	}

	// 作成後、担当者名とロボット名を取得するために再度取得
	scoutService, err = i.scoutServiceRepository.FindByID(scoutService.ID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	scoutService.Templates = input.CreateParam.Templates
	scoutService.GetEntryTimes = input.CreateParam.GetEntryTimes
	scoutService.Password = ""
	output.ScoutService = scoutService

	return output, nil
}

// スカウトサービスを更新する
type UpdateScoutServiceInput struct {
	ScoutServiceID uint
	UpdateParam    entity.CreateOrUpdateScoutServiceParam
}

type UpdateScoutServiceByIDOutput struct {
	ScoutService *entity.ScoutService
}

func (i *ScoutServiceInteractorImpl) UpdateScoutService(input UpdateScoutServiceInput) (UpdateScoutServiceByIDOutput, error) {
	var (
		output       UpdateScoutServiceByIDOutput
		scoutService *entity.ScoutService
		err          error
	)

	// パスワードの暗号化
	input.UpdateParam.Password, err = encrypt(input.UpdateParam.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	scoutService = entity.NewScoutService(
		input.UpdateParam.AgentRobotID,
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.LoginID,
		input.UpdateParam.Password,
		input.UpdateParam.ServiceType,
		input.UpdateParam.IsActive,
		input.UpdateParam.Memo,
		input.UpdateParam.TemplateTitleForEmployed,
		input.UpdateParam.TemplateTitleForUnemployed,
		input.UpdateParam.InterviewAdjustmentTemplateID,
		input.UpdateParam.InflowChannelID,
	)

	// スカウトサービスを更新
	err = i.scoutServiceRepository.Update(input.ScoutServiceID, scoutService)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// エントリー取得時間を更新
	err = i.scoutServiceGetEntryTimeRepository.DeleteByScoutServiceID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	for _, getEntryTime := range input.UpdateParam.GetEntryTimes {
		newGetEntryTime := entity.NewScoutServiceGetEntryTime(
			input.ScoutServiceID,
			getEntryTime.StartHour,
			getEntryTime.StartMinute,
		)

		err = i.scoutServiceGetEntryTimeRepository.Create(newGetEntryTime)
		if err != nil {
			log.Println(err)
			return output, err
		}
	}

	// スカウトテンプレートを更新
	err = i.scoutServiceTemplateRepository.DeleteByScoutServiceID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	for _, scoutServiceTemplate := range input.UpdateParam.Templates {
		newScoutServiceTemplate := entity.NewScoutServiceTemplate(
			input.ScoutServiceID,
			scoutServiceTemplate.StartHour,
			scoutServiceTemplate.StartMinute,
			scoutServiceTemplate.RunOnMonday,
			scoutServiceTemplate.RunOnTuesday,
			scoutServiceTemplate.RunOnWednesday,
			scoutServiceTemplate.RunOnThursday,
			scoutServiceTemplate.RunOnFriday,
			scoutServiceTemplate.RunOnSaturday,
			scoutServiceTemplate.RunOnSunday,
			scoutServiceTemplate.ScoutCount,
			scoutServiceTemplate.SearchTitle,
			scoutServiceTemplate.MessageTitle,
			scoutServiceTemplate.JobInformationTitle,
			scoutServiceTemplate.JobInformationID,
			scoutServiceTemplate.AgeLimit,
			scoutServiceTemplate.ScoutType,
			scoutServiceTemplate.AutoRemind,
			scoutServiceTemplate.ReplyLimit,
		)

		// テンプレートを更新
		err = i.scoutServiceTemplateRepository.Create(newScoutServiceTemplate)
		if err != nil {
			log.Println(err)
			return output, err
		}
	}

	// 更新後、担当者名とロボット名を取得するために再度取得
	scoutService, err = i.scoutServiceRepository.FindByID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	scoutService.Templates = input.UpdateParam.Templates
	scoutService.GetEntryTimes = input.UpdateParam.GetEntryTimes
	scoutService.Password = ""
	output.ScoutService = scoutService

	return output, nil
}

type UpdateScoutServicePasswordInput struct {
	UpdateParam entity.UpdateScoutServicePasswordParam
}

type UpdateScoutServicePasswordByIDOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) UpdateScoutServicePassword(input UpdateScoutServicePasswordInput) (UpdateScoutServicePasswordByIDOutput, error) {
	var (
		output UpdateScoutServicePasswordByIDOutput
		err    error
	)

	// パスワードの暗号化
	input.UpdateParam.Password, err = encrypt(input.UpdateParam.Password)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// スカウトサービスのパスワードを更新
	err = i.scoutServiceRepository.UpdatePassword(input.UpdateParam)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// スカウトサービスを削除する
type DeleteScoutServiceInput struct {
	ScoutServiceID uint
}

type DeleteScoutServiceByIDOutput struct {
	OK bool
}

func (i *ScoutServiceInteractorImpl) DeleteScoutService(input DeleteScoutServiceInput) (DeleteScoutServiceByIDOutput, error) {
	var (
		output DeleteScoutServiceByIDOutput
		err    error
	)

	// スカウトサービスを削除
	err = i.scoutServiceRepository.Delete(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// エントリー取得時間を削除
	err = i.scoutServiceGetEntryTimeRepository.DeleteByScoutServiceID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	err = i.scoutServiceTemplateRepository.DeleteByScoutServiceID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// スカウトサービスを取得する
type ScoutServiceGetByIDInput struct {
	ScoutServiceID uint
}

type ScoutServiceGetByIDOutput struct {
	ScoutService *entity.ScoutService
}

func (i *ScoutServiceInteractorImpl) GetByID(input ScoutServiceGetByIDInput) (ScoutServiceGetByIDOutput, error) {
	var (
		output       ScoutServiceGetByIDOutput
		scoutService *entity.ScoutService
		err          error
	)

	// スカウトサービスを取得
	scoutService, err = i.scoutServiceRepository.FindByID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// エントリー取得時間を取得
	scoutServiceGetEntryTimeList, err := i.scoutServiceGetEntryTimeRepository.GetByScoutServiceID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// テンプレートを取得
	scoutServiceTemplateList, err := i.scoutServiceTemplateRepository.GetByScoutServiceID(input.ScoutServiceID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// エントリー取得時間をマッピング
	for _, scoutServiceGetEntryTime := range scoutServiceGetEntryTimeList {
		scoutService.GetEntryTimes = append(scoutService.GetEntryTimes, *scoutServiceGetEntryTime)
	}

	// テンプレートをマッピング
	for _, scoutServiceTemplate := range scoutServiceTemplateList {
		scoutService.Templates = append(scoutService.Templates, *scoutServiceTemplate)
	}

	scoutService.Password = ""
	output.ScoutService = scoutService

	return output, nil
}

// エージェントIDからスカウトサービスを取得する
type GetListByAgentIDInput struct {
	AgentID uint
}

type GetListByAgentIDOutput struct {
	ScoutServiceList []*entity.ScoutService
}

func (i *ScoutServiceInteractorImpl) GetListByAgentID(input GetListByAgentIDInput) (GetListByAgentIDOutput, error) {
	var (
		output           GetListByAgentIDOutput
		scoutServiceList []*entity.ScoutService
		err              error
	)

	// スカウトサービスを取得
	scoutServiceList, err = i.scoutServiceRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	getEntryTimeList, err := i.scoutServiceGetEntryTimeRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// 取得時間でソート
	sort.Slice(getEntryTimeList, func(i, j int) bool {
		return getEntryTimeList[i].StartHour.Int64 < getEntryTimeList[j].StartHour.Int64
	})

	// スカウトテンプレートを取得
	scoutServiceTemplateList, err := i.scoutServiceTemplateRepository.GetByAgentID(input.AgentID)
	if err != nil {
		log.Println(err)
		return output, err
	}

	// スカウトテンプレートをマッピング
	for _, scoutService := range scoutServiceList {
		scoutService.Password = ""
		for _, getEntryTime := range getEntryTimeList {
			if scoutService.ID == getEntryTime.ScoutServiceID {
				scoutService.GetEntryTimes = append(scoutService.GetEntryTimes, *getEntryTime)
			}
		}

		for _, scoutServiceTemplate := range scoutServiceTemplateList {
			if scoutService.ID == scoutServiceTemplate.ScoutServiceID {
				scoutService.Templates = append(scoutService.Templates, *scoutServiceTemplate)
			}
		}
	}

	output.ScoutServiceList = scoutServiceList

	return output, nil
}

/**********************************************************************************************************/
// 担当エージェントスタッフにメールを送信する
//
func (i *ScoutServiceInteractorImpl) sendErrorMail(messageText string) error {
	var (
	// err error
	)

	log.Println(messageText)

	if messageText == "" {
		return errors.New("メッセージが空です")
	}

	adminEmail := "hidenariyuda@gmail.com"

	err := utility.SendEmail([]string{adminEmail}, "RPAスカウトエラー通知", messageText)
	if err != nil {
		log.Println(err)
		return err
	}

	// if i.app.BatchType != "scout" {
	// 	return nil
	// }

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// 参考: https://risaki-masa.com/how-to-get-api-token-in-slack/
	// slack := utility.NewSlack(i.slack.ReachAccessToken)
	// err = slack.SendContact("C07DMLNMRLG", messageText)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// // メール送信
	// err = utility.SendMailList(
	// 	i.sendgrid.APIKey,
	// 	"RPA処理エラー通知",
	// 	messageText,
	// 	entity.EmailUser{
	// 		Name:  "autoscout事務局",
	// 		Email: "sender@spaceai.jp",
	// 	},
	// 	// 開発確認用
	// 	[]entity.EmailUser{
	// 	},
	// )
	// if err != nil {
	// 	log.Println("開発確認用メール送信失敗")
	// 	log.Println(err)
	// 	return err
	// }

	return nil
}

// スカウト成功メール
func (i *ScoutServiceInteractorImpl) sendScoutSuccessMail(agentStaffID uint, scoutService *entity.ScoutService, scoutServiceTemplateList []*entity.ScoutServiceTemplate) error {
	// IDListから更新後のデータを取得
	scoutServiceTemplateIDList := make([]uint, 0)
	for _, scoutServiceTemplate := range scoutServiceTemplateList {
		scoutServiceTemplateIDList = append(scoutServiceTemplateIDList, scoutServiceTemplate.ID)
	}

	updatedScoutServiceTemplateList, err := i.scoutServiceTemplateRepository.GetByIDList(scoutServiceTemplateIDList)
	if err != nil {
		log.Println(err)
		return err
	}

	// 現在時刻
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))

	totalCount := 0
	templateDetail := ""
	for _, updatedScoutServiceTemplate := range updatedScoutServiceTemplateList {
		lastSendAtStr := updatedScoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05")
		// 12時間以内の場合は成功とする
		if updatedScoutServiceTemplate.LastSendAt.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Before(now.Add(-12 * time.Hour)) {
			// AMBIの数合わせの場合は表示しない
			if scoutService.ServiceType.Int64 == entity.ScoutServiceTypeAmbi &&
				strings.Contains(updatedScoutServiceTemplate.SearchTitle, "数合わせ用") &&
				updatedScoutServiceTemplate.LastSendCount.Int64 == 0 {
				continue
			}
			updatedScoutServiceTemplate.LastSendCount = null.NewInt(0, false)
			lastSendAtStr = ""
		}

		searchTitle := updatedScoutServiceTemplate.SearchTitle
		// マイナビスカウティングの場合はメッセージタイトルを表示する
		if scoutService.ServiceType.Int64 == entity.ScoutServiceTypeMynaviScouting {
			searchTitle = updatedScoutServiceTemplate.JobInformationTitle
		}

		totalCount += int(updatedScoutServiceTemplate.LastSendCount.Int64)

		templateDetail += fmt.Sprintf(
			"・検索タイトル: %s,\n 送信時間: %v,\n 送信数(件): %v\n\n",
			searchTitle,
			lastSendAtStr,
			updatedScoutServiceTemplate.LastSendCount.Int64,
		)
	}
	messageText := fmt.Sprintf(
		"%sのスカウト送信が完了しました。\n\n・合計送信数(件): %v\n\n・詳細\n%s",
		entity.ScoutServiceTypeLabel[scoutService.ServiceType.Int64],
		totalCount,
		templateDetail,
	)

	log.Println(messageText)

	if i.app.BatchType != "scout" {
		return nil
	}

	agentStaff, err := i.agentStaffRepository.FindByID(agentStaffID)
	if err != nil {
		log.Println(err)
		return err
	}

	err = utility.SendEmail([]string{agentStaff.Email}, "RPAスカウト完了通知", messageText)
	if err != nil {
		log.Println(err)
		return err
	}

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// 参考: https://risaki-masa.com/how-to-get-api-token-in-slack/
	// slack := utility.NewSlack(i.slack.ReachAccessToken)
	// err = slack.SendContact("C07DMLNMRLG", messageText)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// err = utility.SendMailList(
	// 	i.sendgrid.APIKey,
	// 	"RPAスカウト通知",
	// 	messageText,
	// 	entity.EmailUser{
	// 		Name:  "autoscout事務局",
	// 		Email: "sender@spaceai.jp",
	// 	},
	// 	[]entity.EmailUser{
	// 		{
	// 			Name:  "autoscout事務局",
	// 			Email: "info@spaceai.jp",
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	log.Println("メール送信失敗")
	// 	log.Println(err)
	// 	return err
	// } else {
	// 	log.Println("メール送信成功")
	// }

	return nil
}
