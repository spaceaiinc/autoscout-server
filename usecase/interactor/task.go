package interactor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type TaskInteractor interface {
	// 旧タスク開始
	// Memo: 旧版の為、新しいapiの動作確認後に削除（が対応）
	SoundOutForJobInformation(input SoundOutForJobInformationInput) (SoundOutForJobInformationOutput, error)
	SoundOutForSendJobListing(input SoundOutForSendJobListingInput) (SoundOutForSendJobListingOutput, error)
	SoundOutForMaskResume(input SoundOutForMaskResumeInput) (SoundOutForMaskResumeOutput, error)
	SoundOutForRequestShareJobSeeker(input SoundOutForRequestShareJobSeekerInput) (SoundOutForRequestShareJobSeekerOutput, error)
	SoundOutForRequestShareJobInformation(input SoundOutForRequestShareJobInformationInput) (SoundOutForRequestShareJobInformationOutput, error)

	// 新タスク開始
	SoundOutGroupForJobInformation(input SoundOutGroupForJobInformationInput) (SoundOutGroupForJobInformationOutput, error)
	SoundOutGroupForSendJobListing(input SoundOutGroupForSendJobListingInput) (SoundOutGroupForSendJobListingOutput, error)

	// 汎用系
	GetLatestTaskListByAgentStaffID(input GetLatestTaskListByAgentStaffIDInput) (GetLatestTaskListByAgentStaffIDOutput, error)
	GetLatestSameTaskListByJobSeekerID(input GetLatestSameTaskListByJobSeekerIDInput) (GetLatestSameTaskListByJobSeekerIDOutput, error)
	GetLatestTaskByJobSeekerIDAndJobInformationID(input GetLatestTaskByJobSeekerIDAndJobInformationIDInput) (GetLatestTaskByJobSeekerIDAndJobInformationIDOutput, error)

	// 複数のタスクを順番に処理するapi
	CreateTaskInBatchProcessing(input CreateTaskInBatchProcessingInput) (CreateTaskInBatchProcessingOutput, error)

	//
	GetSoundOutGroupList(input GetSoundOutGroupListInput) (GetSoundOutGroupListOutput, error)

	// フェーズごとのタスク
	CreateNextTaskAfterEntryPhase(input CreateNextTaskAfterEntryPhaseInput) (CreateNextTaskAfterEntryPhaseOutput, error)
	CreateNextTaskAfterDocumentSelectionPhase(input CreateNextTaskAfterDocumentSelectionPhaseInput) (CreateNextTaskAfterDocumentSelectionPhaseOutput, error)
	CreateNextTaskAfterSelectionPhase(input CreateNextTaskAfterSelectionPhaseInput) (CreateNextTaskAfterSelectionPhaseOutput, error)
	CreateNextTaskAfterDeclinePhase(input CreateNextTaskAfterDeclinePhaseInput) (CreateNextTaskAfterDeclinePhaseOutput, error)                      // 内定辞退
	CreateNextTaskAfterHoldJobOfferPhase(input CreateNextTaskAfterHoldJobOfferPhaseInput) (CreateNextTaskAfterHoldJobOfferPhaseOutput, error)       // 内定保留
	CreateNextTaskAfterAcceptJobOfferPhase(input CreateNextTaskAfterAcceptJobOfferPhaseInput) (CreateNextTaskAfterAcceptJobOfferPhaseOutput, error) // 内定承諾
	CreateNextSameTaskList(input CreateNextSameTaskListInput) (CreateNextSameTaskListOutput, error)
	CreateEntryTaskFromMatchingJob(input CreateEntryTaskFromMatchingJobInput) (CreateEntryTaskFromMatchingJobOutput, error) // マイページのマッチ求人からエントリー
	UpdateTaskGroupDocument(input UpdateTaskGroupDocumentInput) (UpdateTaskGroupDocumentOutput, error)

	UpdateRALastWatched(input UpdateRALastWatchedInput) (UpdateRALastWatchedOutput, error)
	UpdateRALastRequest(input UpdateRALastRequestInput) (UpdateRALastRequestOutput, error)
	UpdateCALastWatched(input UpdateCALastWatchedInput) (UpdateCALastWatchedOutput, error)
	UpdateCALastRequest(input UpdateCALastRequestInput) (UpdateCALastRequestOutput, error)
	UpdateLastRequest(input UpdateLastRequestInput) (UpdateLastRequestOutput, error)
	UpdateExternalJob(input UpdateExternalJobInput) (UpdateExternalJobOutput, error)

	// 指定IDのタスク情報を取得する関数
	GetTaskByID(input GetTaskByIDInput) (GetTaskByIDOutput, error)
	GetTaskListByAgentIDAndPage(input GetTaskListByAgentIDAndPageInput) (GetTaskListByAgentIDAndPageOutput, error)
	GetSearchTaskListByAgentIDAndPage(input GetSearchTaskListByAgentIDAndPageInput) (GetSearchTaskListByAgentIDAndPageOutput, error)
	GetTaskListAfterEntryByJobSeekerID(input GetTaskListAfterEntryByJobSeekerIDInput) (GetTaskListAfterEntryByJobSeekerIDOutput, error)
	GetTaskGroupByID(input GetTaskGroupByIDInput) (GetTaskGroupByIDOutput, error)
	GetActiveTaskCountByBillingAddressID(input GetActiveTaskCountByBillingAddressIDInput) (GetActiveTaskCountByBillingAddressIDOutput, error)
	GetActiveTaskCountByJobInformationID(input GetActiveTaskCountByJobInformationIDInput) (GetActiveTaskCountByJobInformationIDOutput, error)
	GetActiveTaskCountBySelectionID(input GetActiveTaskCountBySelectionIDInput) (GetActiveTaskCountBySelectionIDOutput, error)

	DeleteTask(input DeleteTaskInput) (output DeleteTaskOutput, err error) // GoogleCalendarからタスクを取得する関数
	// Admin API

	// Batch API
	BatchNotifyUnwatched(input BatchNotifyUnwatchedInput) (BatchNotifyUnwatchedOutput, error)

	GetJobSeekerTaskListByAgentStaffID(input GetJobSeekerTaskListByAgentStaffIDInput) (GetJobSeekerTaskListByAgentStaffIDOutput, error)
}

type TaskInteractorImpl struct {
	firebase                                  usecase.Firebase
	sendgrid                                  config.Sendgrid
	oneSignal                                 config.OneSignal
	taskRepository                            usecase.TaskRepository
	taskGroupRepository                       usecase.TaskGroupRepository
	taskGroupDocumentRepository               usecase.TaskGroupDocumentRepository
	agentRepository                           usecase.AgentRepository
	agentStaffRepository                      usecase.AgentStaffRepository
	jobInformationRepository                  usecase.JobInformationRepository
	jobSeekerRepository                       usecase.JobSeekerRepository
	jobSeekerStudentHistoryRepository         usecase.JobSeekerStudentHistoryRepository
	jobSeekerWorkHistoryRepository            usecase.JobSeekerWorkHistoryRepository
	jobSeekerExperienceIndustryRepository     usecase.JobSeekerExperienceIndustryRepository
	jobSeekerDepartmentHistoryRepository      usecase.JobSeekerDepartmentHistoryRepository
	jobSeekerLicenseRepository                usecase.JobSeekerLicenseRepository
	jobSeekerSelfPromotionRepository          usecase.JobSeekerSelfPromotionRepository
	jobSeekerDocumentRepository               usecase.JobSeekerDocumentRepository
	jobSeekerDesiredIndustryRepository        usecase.JobSeekerDesiredIndustryRepository
	jobSeekerDesiredOccupationRepository      usecase.JobSeekerDesiredOccupationRepository
	jobSeekerDesiredWorkLocationRepository    usecase.JobSeekerDesiredWorkLocationRepository
	jobSeekerDesiredHolidayTypeRepository     usecase.JobSeekerDesiredHolidayTypeRepository
	jobSeekerDevelopmentSkillRepository       usecase.JobSeekerDevelopmentSkillRepository
	jobSeekerLanguageSkillRepository          usecase.JobSeekerLanguageSkillRepository
	jobSeekerPCToolRepository                 usecase.JobSeekerPCToolRepository
	jobSeekerHideToAgentRepository            usecase.JobSeekerHideToAgentRepository
	jobSeekerExperienceOccupationRepository   usecase.JobSeekerExperienceOccupationRepository
	jobSeekerDesiredCompanyScaleRepository    usecase.JobSeekerDesiredCompanyScaleRepository
	chatGroupWithJobSeekerRepository          usecase.ChatGroupWithJobSeekerRepository
	evaluationPointRepository                 usecase.EvaluationPointRepository
	chatMessageWithJobSeekerRepository        usecase.ChatMessageWithJobSeekerRepository
	selectionQuestionnaireRepository          usecase.SelectionQuestionnaireRepository
	selectionQuestionnaireMyRankingRepository usecase.SelectionQuestionnaireMyRankingRepository
	jobInfoSelectionFlowPatternRepository     usecase.JobInformationSelectionFlowPatternRepository
	jobInfoSelectionInformationRepository     usecase.JobInformationSelectionInformationRepository
	saleRepository                            usecase.SaleRepository
	taskIsRecommendDocumentRepository         usecase.TaskIsRecommendDocumentRepository
	jobSeekerScheduleRepository               usecase.JobSeekerScheduleRepository
	jobSeekerRescheduleRepository             usecase.JobSeekerRescheduleRepository
	chatGroupWithAgentRepository              usecase.ChatGroupWithAgentRepository
	chatThreadWithAgentRepository             usecase.ChatThreadWithAgentRepository
	chatMessageWithAgentRepository            usecase.ChatMessageWithAgentRepository
	chatMessageToUserWithAgentRepository      usecase.ChatMessageToUserWithAgentRepository
	emailWithJobSeekerRepository              usecase.EmailWithJobSeekerRepository
	jobSeekerInterestedJobListingRepository   usecase.JobSeekerInterestedJobListingRepository
}

// TaskInteractorImpl is an implementation of TaskInteractor
func NewTaskInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	tR usecase.TaskRepository,
	tgR usecase.TaskGroupRepository,
	tgdR usecase.TaskGroupDocumentRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	jIR usecase.JobInformationRepository,
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
	cgR usecase.ChatGroupWithJobSeekerRepository,
	epR usecase.EvaluationPointRepository,
	cmjR usecase.ChatMessageWithJobSeekerRepository,
	sqR usecase.SelectionQuestionnaireRepository,
	smR usecase.SelectionQuestionnaireMyRankingRepository,
	jsfpR usecase.JobInformationSelectionFlowPatternRepository,
	jsiR usecase.JobInformationSelectionInformationRepository,
	saR usecase.SaleRepository,
	irdR usecase.TaskIsRecommendDocumentRepository,
	jssR usecase.JobSeekerScheduleRepository,
	jsrR usecase.JobSeekerRescheduleRepository,
	cgwaR usecase.ChatGroupWithAgentRepository,
	ctwaR usecase.ChatThreadWithAgentRepository,
	cmwaA usecase.ChatMessageWithAgentRepository,
	cmtuwaR usecase.ChatMessageToUserWithAgentRepository,
	ewjsR usecase.EmailWithJobSeekerRepository,
	jsijlR usecase.JobSeekerInterestedJobListingRepository,
) TaskInteractor {
	return &TaskInteractorImpl{
		firebase:                                  fb,
		sendgrid:                                  sg,
		oneSignal:                                 os,
		taskRepository:                            tR,
		taskGroupRepository:                       tgR,
		taskGroupDocumentRepository:               tgdR,
		agentRepository:                           aR,
		agentStaffRepository:                      asR,
		jobInformationRepository:                  jIR,
		jobSeekerRepository:                       jsR,
		jobSeekerStudentHistoryRepository:         jsshR,
		jobSeekerWorkHistoryRepository:            jswhR,
		jobSeekerExperienceIndustryRepository:     jseiR,
		jobSeekerExperienceOccupationRepository:   jseoR,
		jobSeekerLicenseRepository:                jslR,
		jobSeekerSelfPromotionRepository:          jsspR,
		jobSeekerDocumentRepository:               jsdR,
		jobSeekerDesiredIndustryRepository:        jsdiR,
		jobSeekerDesiredOccupationRepository:      jsdoR,
		jobSeekerDesiredWorkLocationRepository:    jsdwlR,
		jobSeekerDesiredHolidayTypeRepository:     jsdhtR,
		jobSeekerDevelopmentSkillRepository:       jsdsR,
		jobSeekerLanguageSkillRepository:          jslsR,
		jobSeekerPCToolRepository:                 jsptR,
		jobSeekerHideToAgentRepository:            jshR,
		jobSeekerDepartmentHistoryRepository:      jsdhR,
		jobSeekerDesiredCompanyScaleRepository:    jsdcsR,
		chatGroupWithJobSeekerRepository:          cgR,
		evaluationPointRepository:                 epR,
		chatMessageWithJobSeekerRepository:        cmjR,
		selectionQuestionnaireRepository:          sqR,
		selectionQuestionnaireMyRankingRepository: smR,
		jobInfoSelectionFlowPatternRepository:     jsfpR,
		jobInfoSelectionInformationRepository:     jsiR,
		saleRepository:                            saR,
		taskIsRecommendDocumentRepository:         irdR,
		jobSeekerScheduleRepository:               jssR,
		jobSeekerRescheduleRepository:             jsrR,
		chatGroupWithAgentRepository:              cgwaR,
		chatThreadWithAgentRepository:             ctwaR,
		chatMessageWithAgentRepository:            cmwaA,
		chatMessageToUserWithAgentRepository:      cmtuwaR,
		emailWithJobSeekerRepository:              ewjsR,
		jobSeekerInterestedJobListingRepository:   jsijlR,
	}
}

/****************************************************************************************/
// 旧タスク開始の関数
// Memo: 旧版の為、新しいapiの動作確認後に削除（が対応）
//

type SoundOutForJobInformationInput struct {
	SoundOutParam entity.SoundOutForJobInformationParam
}

type SoundOutForJobInformationOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) SoundOutForJobInformation(input SoundOutForJobInformationInput) (SoundOutForJobInformationOutput, error) {
	var (
		output      SoundOutForJobInformationOutput
		err         error
		idList      []uint
		encountered = map[uint]bool{}
	)

	for _, group := range input.SoundOutParam.GroupList {
		fmt.Println("RA: ", group.RAStaffID)
		fmt.Println("CA: ", group.CAStaffID)

		var IsDoubleSided = false

		if group.RAStaffID == group.CAStaffID {
			IsDoubleSided = true
		}

		taskGroup := entity.NewTaskGroup(
			group.JobSeekerID,
			group.JobInformationID,
			IsDoubleSided,
			"",
			"",
			"",
			group.RAAgentID,
			group.CAAgentID,
		)

		// タスクグループを作成
		err = i.taskGroupRepository.Create(taskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 「エントリー（求人打診依頼）」のタスクを作成
		task := entity.NewTask(
			taskGroup.ID,
			null.NewInt(int64(entity.Entry), true), // エントリー
			null.NewInt(int64(entity.SoundOutJobInformation), true), // 求人打診依頼
			null.NewInt(int64(entity.CA), true),                     // 担当者タイプ
			group.CAStaffID,                                         // CAタスクの場合はCA担当者のIDを入れる
			input.SoundOutParam.Remarks,
			input.SoundOutParam.DeadlineDay,
			input.SoundOutParam.DeadlineTime,
			input.SoundOutParam.TalkAboutInInterview,
			input.SoundOutParam.ScheduleCollectionCondition,
			input.SoundOutParam.ExamGuideContent,
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(task)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// RAからCAにタスク作成する場合
		err := i.taskGroupRepository.UpdateRALastRequestAt(taskGroup.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task.Title = group.Title
		task.CompanyName = group.CompanyName
		task.LastName = group.LastName
		task.FirstName = group.FirstName
		task.LastFurigana = group.LastFurigana
		task.FirstFurigana = group.FirstFurigana
		task.RAAgentID = group.RAAgentID
		task.CAAgentID = group.CAAgentID

		output.TaskList = append(output.TaskList, task)

		// push通知を送るためのCA担当者IDのリストを作成（重複除外）
		if !encountered[group.CAStaffID] {
			encountered[group.CAStaffID] = true
			idList = append(idList, group.CAStaffID)
		}
	}

	for _, caStaffID := range idList {
		// プッシュ通知
		redirectURL := os.Getenv("BASE_DOMAIN")
		contents := "新着CAタスクの依頼があります。"
		topic := "TaskCA"

		if caStaffID != 0 {
			caStaff, err := i.agentStaffRepository.FindByID(caStaffID)
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
				fmt.Println("WebPushの通知でエラー")
				fmt.Println(err)
			}
		}
	}

	return output, nil
}

type SoundOutForSendJobListingInput struct {
	SoundOutParam entity.SoundOutForJobInformationForSendMessageParam
}

type SoundOutForSendJobListingOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) SoundOutForSendJobListing(input SoundOutForSendJobListingInput) (SoundOutForSendJobListingOutput, error) {
	var (
		output SoundOutForSendJobListingOutput
		err    error

		urlTextList          = map[uint]string{}
		sendGroupEncountered = map[uint]bool{}
		caIDList             []uint
		caIDEncountered      = map[uint]bool{}
	)

	for _, group := range input.SoundOutParam.GroupList {
		fmt.Println("RA: ", group.RAStaffID)
		fmt.Println("CA: ", group.CAStaffID)

		var IsDoubleSided = false

		if group.RAStaffID == group.CAStaffID {
			IsDoubleSided = true
		}

		taskGroup := entity.NewTaskGroup(
			group.JobSeekerID,
			group.JobInformationID,
			IsDoubleSided,
			"",
			"",
			"",
			group.RAAgentID,
			group.CAAgentID,
		)

		// タスクグループを作成
		err = i.taskGroupRepository.Create(taskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// タスクを定義
		var task *entity.Task

		// CA担当が自分の場合
		if group.CAStaffID == input.SoundOutParam.AgentStaffID {
			// 「エントリー（応募意思確認中）」のタスクを作成
			task = entity.NewTask(
				taskGroup.ID,
				null.NewInt(int64(entity.Entry), true), // エントリー
				null.NewInt(int64(entity.ConfirmApplicationIntention), true), // 応募意思確認中
				null.NewInt(int64(entity.CA), true),                          // 担当者タイプ
				group.CAStaffID,                                              // CAタスクの場合はCA担当者のIDを入れる
				input.SoundOutParam.Remarks,
				input.SoundOutParam.DeadlineDay,
				input.SoundOutParam.DeadlineTime,
				input.SoundOutParam.TalkAboutInInterview,
				input.SoundOutParam.ScheduleCollectionCondition,
				input.SoundOutParam.ExamGuideContent,
				false,
			)

			// タスク作成
			err = i.taskRepository.Create(task)
			if err != nil {
				return output, err
			}

			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(taskGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メッセージに含める求人票URLをまとめる
			baseURL := os.Getenv("BASE_DOMAIN")

			// 求人票送付に使用するURLを生成
			listingURL := fmt.Sprintf(
				"\n■%s\n%s/jobs/sound_out/?job_seeker=%s&job_information=%s\n",
				group.CompanyName, baseURL, group.JobSeekerUUID, group.JobInformationUUID,
			)

			// 求職者に複数の求人を一括で送るために送付URLをまとめる
			urlTextList[group.JobSeekerID] = urlTextList[group.JobSeekerID] + listingURL

			if !sendGroupEncountered[group.JobSeekerID] {
				sendGroupEncountered[group.JobSeekerID] = true
			}

		} else {
			// 「エントリー（求人打診依頼）」のタスクを作成
			task := entity.NewTask(
				taskGroup.ID,
				null.NewInt(int64(entity.Entry), true), // エントリー
				null.NewInt(int64(entity.SoundOutJobInformation), true), // 求人打診依頼
				null.NewInt(int64(entity.CA), true),                     // 担当者タイプ
				group.CAStaffID,                                         // CAタスクの場合はCA担当者のIDを入れる
				"",
				input.SoundOutParam.DeadlineDay,
				input.SoundOutParam.DeadlineTime,
				input.SoundOutParam.TalkAboutInInterview,
				input.SoundOutParam.ScheduleCollectionCondition,
				input.SoundOutParam.ExamGuideContent,
				false,
			)

			// タスク作成
			err = i.taskRepository.Create(task)
			if err != nil {
				return output, err
			}
		}

		task.Title = group.Title
		task.CompanyName = group.CompanyName
		task.LastName = group.LastName
		task.FirstName = group.FirstName
		task.LastFurigana = group.LastFurigana
		task.FirstFurigana = group.FirstFurigana
		task.RAAgentID = group.RAAgentID
		task.CAAgentID = group.CAAgentID

		output.TaskList = append(output.TaskList, task)

		// push通知を送るためのCA担当者IDのリストを作成（重複除外）
		if !caIDEncountered[group.CAStaffID] {
			caIDEncountered[group.CAStaffID] = true
			caIDList = append(caIDList, group.CAStaffID)
		}
	}

	for _, message := range input.SoundOutParam.Messages {
		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(uint(input.SoundOutParam.AgentStaffID))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(input.SoundOutParam.AgentID, message.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		if chatGroup.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			jobSeeker, err := i.jobSeekerRepository.FindByID(message.JobSeekerID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			fmt.Println("メッセージの送信に成功しました:", call)

			// エージェントから求職者へのチャットメッセージを作成
			chatMessage := entity.NewChatMessageWithJobSeeker(
				chatGroup.ID,
				null.NewInt(0, true), // エージェント
				message.Content,
				"",
				"",
				"",
				"",
				null.NewInt(0, false),
				"",
				null.NewInt(0, true),
				time.Now().In(time.UTC),
			)

			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				return output, err
			}
		} else {
			// メール送信
			err = utility.SendMailToSingle(
				i.sendgrid.APIKey,
				message.Subject,
				message.Content,
				entity.EmailUser{
					Name:  staff.StaffName,
					Email: staff.Email,
				},
				entity.EmailUser{
					Name:  message.Name,
					Email: message.Email,
				},
				nil,
			)

			fmt.Println("メール送信")

			if err != nil {
				fmt.Println("メール送信失敗")
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				message.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// プッシュ通知
	for _, caStaffID := range caIDList {
		redirectURL := os.Getenv("BASE_DOMAIN")
		contents := "新着CAタスクの依頼があります。"
		topic := "TaskCA"

		if caStaffID != 0 {
			caStaff, err := i.agentStaffRepository.FindByID(caStaffID)
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
				fmt.Println("WebPushの通知でエラー")
				fmt.Println(err)
			}
		}
	}

	return output, nil
}

type SoundOutForMaskResumeInput struct {
	SoundOutParam entity.SoundOutForJobInformationParam
}

type SoundOutForMaskResumeOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) SoundOutForMaskResume(input SoundOutForMaskResumeInput) (SoundOutForMaskResumeOutput, error) {
	var (
		output      SoundOutForMaskResumeOutput
		err         error
		idList      []uint
		encountered = map[uint]bool{}
	)

	for _, group := range input.SoundOutParam.GroupList {
		fmt.Println("RA: ", group.RAStaffID)
		fmt.Println("CA: ", group.CAStaffID)

		var IsDoubleSided = false

		if group.RAStaffID == group.CAStaffID {
			IsDoubleSided = true
		}

		taskGroup := entity.NewTaskGroup(
			group.JobSeekerID,
			group.JobInformationID,
			IsDoubleSided,
			"",
			"",
			"",
			group.RAAgentID,
			group.CAAgentID,
		)

		// タスクグループを作成
		err = i.taskGroupRepository.Create(taskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task := entity.NewTask(
			taskGroup.ID,
			null.NewInt(int64(entity.Entry), true),        // エントリー
			null.NewInt(int64(entity.SoundOutMask), true), // マスクレジュメ打診の依頼
			null.NewInt(int64(entity.RA), true),           // 担当者タイプ
			group.RAStaffID,                               // RAタスクの場合はRA担当者のIDを入れる
			input.SoundOutParam.Remarks,
			input.SoundOutParam.DeadlineDay,
			input.SoundOutParam.DeadlineTime,
			input.SoundOutParam.TalkAboutInInterview,
			input.SoundOutParam.ScheduleCollectionCondition,
			input.SoundOutParam.ExamGuideContent,
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(task)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// CAからRAにタスク作成する場合
		err := i.taskGroupRepository.UpdateCALastRequestAt(taskGroup.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task.Title = group.Title
		task.CompanyName = group.CompanyName
		task.LastName = group.LastName
		task.FirstName = group.FirstName
		task.LastFurigana = group.LastFurigana
		task.FirstFurigana = group.FirstFurigana
		task.RAAgentID = group.RAAgentID
		task.CAAgentID = group.CAAgentID

		output.TaskList = append(output.TaskList, task)

		// push通知を送るためのRA担当者IDのリストを作成（重複除外）
		if !encountered[group.RAStaffID] {
			encountered[group.RAStaffID] = true
			idList = append(idList, group.RAStaffID)
		}
	}

	for _, raStaffID := range idList {
		// プッシュ通知
		redirectURL := os.Getenv("BASE_DOMAIN")
		contents := "新着RAタスクの依頼があります。"
		topic := "TaskRA"

		if raStaffID != 0 {
			raStaff, err := i.agentStaffRepository.FindByID(raStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			err = utility.WebPush(
				i.oneSignal.AppID,
				i.oneSignal.APIKey,
				raStaff.FirebaseID,
				"タスク通知",
				contents,
				topic,
				redirectURL,
			)
			if err != nil {
				fmt.Println("WebPushの通知でエラー")
				fmt.Println(err)
			}
		}
	}

	return output, nil
}

// シェア依頼処理（求人検索ページから）
type SoundOutForRequestShareJobSeekerInput struct {
	SoundOutParam entity.SoundOutForJobInformationParam
}

type SoundOutForRequestShareJobSeekerOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) SoundOutForRequestShareJobSeeker(input SoundOutForRequestShareJobSeekerInput) (SoundOutForRequestShareJobSeekerOutput, error) {
	var (
		output      SoundOutForRequestShareJobSeekerOutput
		err         error
		idList      []uint
		encountered = map[uint]bool{}
	)

	// 「エントリー/求職者シェア依頼」のタスクを作成
	for _, group := range input.SoundOutParam.GroupList {
		var IsDoubleSided = false

		if group.RAStaffID == group.CAStaffID {
			IsDoubleSided = true
		}

		taskGroup := entity.NewTaskGroup(
			group.JobSeekerID,
			group.JobInformationID,
			IsDoubleSided,
			"",
			"",
			"",
			group.RAAgentID,
			group.CAAgentID,
		)

		// タスクグループを作成
		err = i.taskGroupRepository.Create(taskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task := entity.NewTask(
			taskGroup.ID,
			null.NewInt(int64(entity.Entry), true), // エントリー
			null.NewInt(int64(entity.RequestShareJobSeeker), true), // 求職者シェア依頼
			null.NewInt(int64(entity.CA), true),                    // 担当者タイプ
			group.CAStaffID,                                        // CAタスクの場合はCA担当者のIDを入れる
			input.SoundOutParam.Remarks,
			input.SoundOutParam.DeadlineDay,
			input.SoundOutParam.DeadlineTime,
			input.SoundOutParam.TalkAboutInInterview,
			input.SoundOutParam.ScheduleCollectionCondition,
			input.SoundOutParam.ExamGuideContent,
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(task)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// RAからCAにタスク作成する場合
		err := i.taskGroupRepository.UpdateRALastRequestAt(taskGroup.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task.Title = group.Title
		task.CompanyName = group.CompanyName
		task.LastName = group.LastName
		task.FirstName = group.FirstName
		task.LastFurigana = group.LastFurigana
		task.FirstFurigana = group.FirstFurigana
		task.RAAgentID = group.RAAgentID
		task.CAAgentID = group.CAAgentID

		output.TaskList = append(output.TaskList, task)

		// push通知を送るためのRA担当者IDのリストを作成（重複除外）
		if !encountered[group.RAStaffID] {
			encountered[group.CAStaffID] = true
			idList = append(idList, group.CAStaffID)
		}
	}

	for _, caStaffID := range idList {
		// プッシュ通知
		redirectURL := os.Getenv("BASE_DOMAIN")
		contents := "新着CAタスクの依頼があります。"
		topic := "TaskCA"

		if caStaffID != 0 {
			raStaff, err := i.agentStaffRepository.FindByID(caStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			err = utility.WebPush(
				i.oneSignal.AppID,
				i.oneSignal.APIKey,
				raStaff.FirebaseID,
				"タスク通知",
				contents,
				topic,
				redirectURL,
			)
			if err != nil {
				fmt.Println("WebPushの通知でエラー")
				fmt.Println(err)
			}
		}
	}

	return output, nil
}

// シェア依頼処理（求職者検索ページから）
type SoundOutForRequestShareJobInformationInput struct {
	SoundOutParam entity.SoundOutForJobInformationParam
}

type SoundOutForRequestShareJobInformationOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) SoundOutForRequestShareJobInformation(input SoundOutForRequestShareJobInformationInput) (SoundOutForRequestShareJobInformationOutput, error) {
	var (
		output      SoundOutForRequestShareJobInformationOutput
		err         error
		idList      []uint
		encountered = map[uint]bool{}
	)

	// 「エントリー/求人シェア依頼」のタスクを作成
	for _, group := range input.SoundOutParam.GroupList {
		fmt.Println("RA: ", group.RAStaffID)
		fmt.Println("CA: ", group.CAStaffID)

		var IsDoubleSided = false

		if group.RAStaffID == group.CAStaffID {
			IsDoubleSided = true
		}

		taskGroup := entity.NewTaskGroup(
			group.JobSeekerID,
			group.JobInformationID,
			IsDoubleSided,
			"",
			"",
			"",
			group.RAAgentID,
			group.CAAgentID,
		)

		// タスクグループを作成
		err = i.taskGroupRepository.Create(taskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task := entity.NewTask(
			taskGroup.ID,
			null.NewInt(int64(entity.Entry), true), // エントリー
			null.NewInt(int64(entity.RequestShareJobInformation), true), // 求人シェア依頼
			null.NewInt(int64(entity.RA), true),                         // 担当者タイプ
			group.RAStaffID,                                             // RAタスクの場合はRA担当者のIDを入れる
			input.SoundOutParam.Remarks,
			input.SoundOutParam.DeadlineDay,
			input.SoundOutParam.DeadlineTime,
			input.SoundOutParam.TalkAboutInInterview,
			input.SoundOutParam.ScheduleCollectionCondition,
			input.SoundOutParam.ExamGuideContent,
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(task)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// CAからRAにタスク作成する場合
		err := i.taskGroupRepository.UpdateCALastRequestAt(taskGroup.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task.Title = group.Title
		task.CompanyName = group.CompanyName
		task.LastName = group.LastName
		task.FirstName = group.FirstName
		task.LastFurigana = group.LastFurigana
		task.FirstFurigana = group.FirstFurigana
		task.RAAgentID = group.RAAgentID
		task.CAAgentID = group.CAAgentID

		output.TaskList = append(output.TaskList, task)

		// push通知を送るためのRA担当者IDのリストを作成（重複除外）
		if !encountered[group.RAStaffID] {
			encountered[group.RAStaffID] = true
			idList = append(idList, group.RAStaffID)
		}
	}

	for _, raStaffID := range idList {
		// プッシュ通知
		redirectURL := os.Getenv("BASE_DOMAIN")
		contents := "新着RAタスクの依頼があります。"
		topic := "TaskRA"

		if raStaffID != 0 {
			raStaff, err := i.agentStaffRepository.FindByID(raStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			err = utility.WebPush(
				i.oneSignal.AppID,
				i.oneSignal.APIKey,
				raStaff.FirebaseID,
				"タスク通知",
				contents,
				topic,
				redirectURL,
			)
			if err != nil {
				fmt.Println("WebPushの通知でエラー")
				fmt.Println(err)
			}
		}
	}

	return output, nil
}

/****************************************************************************************/
// 新タスク開始の関数
//

type SoundOutGroupForJobInformationInput struct {
	SoundOutParam entity.SoundOutForJobInformationParam
}

type SoundOutGroupForJobInformationOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) SoundOutGroupForJobInformation(input SoundOutGroupForJobInformationInput) (SoundOutGroupForJobInformationOutput, error) {
	var (
		output      SoundOutGroupForJobInformationOutput
		err         error
		idList      []uint
		encountered = map[uint]bool{}
	)

	for _, group := range input.SoundOutParam.GroupList {
		fmt.Println("RA: ", group.RAStaffID)
		fmt.Println("CA: ", group.CAStaffID)

		var IsDoubleSided = false

		if group.RAStaffID == group.CAStaffID {
			IsDoubleSided = true
		}

		// 外部求人関連の値を格納する変数を定義
		var (
			externalJobInformationTitle string
			externalCompanyName         string
			externalJobListingURL       string
		)

		// 外部求人の場合は変数に格納
		if group.IsExternal {
			externalJobInformationTitle = group.Title
			externalCompanyName = group.CompanyName
			externalJobListingURL = group.JobListingURL
		}

		taskGroup := entity.NewTaskGroup(
			group.JobSeekerID,
			group.JobInformationID,
			IsDoubleSided,
			externalJobInformationTitle,
			externalCompanyName,
			externalJobListingURL,
			group.RAAgentID,
			group.CAAgentID,
		)

		// タスクグループを作成
		err = i.taskGroupRepository.Create(taskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 「エントリー（求人打診依頼）」のタスクを作成
		task := entity.NewTask(
			taskGroup.ID,
			null.NewInt(int64(entity.Entry), true), // エントリー
			null.NewInt(int64(entity.SoundOutJobInformation), true), // 求人打診依頼
			null.NewInt(int64(entity.CA), true),                     // 担当者タイプ
			group.CAStaffID,                                         // CAタスクの場合はCA担当者のIDを入れる
			input.SoundOutParam.Remarks,
			input.SoundOutParam.DeadlineDay,
			input.SoundOutParam.DeadlineTime,
			input.SoundOutParam.TalkAboutInInterview,
			input.SoundOutParam.ScheduleCollectionCondition,
			input.SoundOutParam.ExamGuideContent,
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(task)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// RAからCAにタスク作成する場合
		err := i.taskGroupRepository.UpdateRALastRequestAt(taskGroup.ID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		task.Title = group.Title
		task.CompanyName = group.CompanyName
		task.LastName = group.LastName
		task.FirstName = group.FirstName
		task.LastFurigana = group.LastFurigana
		task.FirstFurigana = group.FirstFurigana
		task.RAAgentID = group.RAAgentID
		task.CAAgentID = group.CAAgentID

		output.TaskList = append(output.TaskList, task)

		// push通知を送るためのCA担当者IDのリストを作成（重複除外）
		if !encountered[group.CAStaffID] {
			encountered[group.CAStaffID] = true
			idList = append(idList, group.CAStaffID)
		}
	}

	for _, caStaffID := range idList {
		// プッシュ通知
		redirectURL := os.Getenv("BASE_DOMAIN")
		contents := "新着CAタスクの依頼があります。"
		topic := "TaskCA"

		if caStaffID != 0 {
			caStaff, err := i.agentStaffRepository.FindByID(caStaffID)
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
				fmt.Println("WebPushの通知でエラー")
				fmt.Println(err)
			}
		}
	}

	return output, nil
}

type SoundOutGroupForSendJobListingInput struct {
	SoundOutParam entity.SoundOutForJobInformationForSendMessageParam
}

type SoundOutGroupForSendJobListingOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) SoundOutGroupForSendJobListing(input SoundOutGroupForSendJobListingInput) (SoundOutGroupForSendJobListingOutput, error) {
	var (
		output SoundOutGroupForSendJobListingOutput
		err    error

		urlTextList          = map[uint]string{}
		sendGroupEncountered = map[uint]bool{}
		caIDList             []uint
		caIDEncountered      = map[uint]bool{}
	)

	for _, group := range input.SoundOutParam.GroupList {
		fmt.Println("RA: ", group.RAStaffID)
		fmt.Println("CA: ", group.CAStaffID)

		var IsDoubleSided = false

		if group.RAStaffID == group.CAStaffID {
			IsDoubleSided = true
		}

		// 外部求人関連の値を格納する変数を定義
		var (
			externalJobInformationTitle string
			externalCompanyName         string
			externalJobListingURL       string
		)

		// 外部求人の場合は変数に格納
		if group.IsExternal {
			externalJobInformationTitle = group.Title
			externalCompanyName = group.CompanyName
			externalJobListingURL = group.JobListingURL
		}

		taskGroup := entity.NewTaskGroup(
			group.JobSeekerID,
			group.JobInformationID,
			IsDoubleSided,
			externalJobInformationTitle,
			externalCompanyName,
			externalJobListingURL,
			group.RAAgentID,
			group.CAAgentID,
		)

		// タスクグループを作成
		err = i.taskGroupRepository.Create(taskGroup)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// タスクを定義
		var task *entity.Task

		// CA担当が自分の場合
		if group.CAStaffID == input.SoundOutParam.AgentStaffID {
			// 「エントリー（応募意思確認中）」のタスクを作成
			task = entity.NewTask(
				taskGroup.ID,
				null.NewInt(int64(entity.Entry), true), // エントリー
				null.NewInt(int64(entity.ConfirmApplicationIntention), true), // 応募意思確認中
				null.NewInt(int64(entity.CA), true),                          // 担当者タイプ
				group.CAStaffID,                                              // CAタスクの場合はCA担当者のIDを入れる
				input.SoundOutParam.Remarks,
				input.SoundOutParam.DeadlineDay,
				input.SoundOutParam.DeadlineTime,
				input.SoundOutParam.TalkAboutInInterview,
				input.SoundOutParam.ScheduleCollectionCondition,
				input.SoundOutParam.ExamGuideContent,
				false,
			)

			// タスク作成
			err = i.taskRepository.Create(task)
			if err != nil {
				return output, err
			}

			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(taskGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メッセージに含める求人票URLをまとめる
			baseURL := os.Getenv("BASE_DOMAIN")

			// 求人票送付に使用するURLを生成
			listingURL := fmt.Sprintf(
				"\n■%s\n%s/jobs/sound_out/?job_seeker=%s&job_information=%s\n",
				group.CompanyName, baseURL, group.JobSeekerUUID, group.JobInformationUUID,
			)

			// 求職者に複数の求人を一括で送るために送付URLをまとめる
			urlTextList[group.JobSeekerID] = urlTextList[group.JobSeekerID] + listingURL

			if !sendGroupEncountered[group.JobSeekerID] {
				sendGroupEncountered[group.JobSeekerID] = true
			}

		} else {
			// 「エントリー（求人打診依頼）」のタスクを作成
			task := entity.NewTask(
				taskGroup.ID,
				null.NewInt(int64(entity.Entry), true), // エントリー
				null.NewInt(int64(entity.SoundOutJobInformation), true), // 求人打診依頼
				null.NewInt(int64(entity.CA), true),                     // 担当者タイプ
				group.CAStaffID,                                         // CAタスクの場合はCA担当者のIDを入れる
				"",
				input.SoundOutParam.DeadlineDay,
				input.SoundOutParam.DeadlineTime,
				input.SoundOutParam.TalkAboutInInterview,
				input.SoundOutParam.ScheduleCollectionCondition,
				input.SoundOutParam.ExamGuideContent,
				false,
			)

			// タスク作成
			err = i.taskRepository.Create(task)
			if err != nil {
				return output, err
			}
		}

		task.Title = group.Title
		task.CompanyName = group.CompanyName
		task.LastName = group.LastName
		task.FirstName = group.FirstName
		task.LastFurigana = group.LastFurigana
		task.FirstFurigana = group.FirstFurigana
		task.RAAgentID = group.RAAgentID
		task.CAAgentID = group.CAAgentID

		output.TaskList = append(output.TaskList, task)

		// push通知を送るためのCA担当者IDのリストを作成（重複除外）
		if !caIDEncountered[group.CAStaffID] {
			caIDEncountered[group.CAStaffID] = true
			caIDList = append(caIDList, group.CAStaffID)
		}
	}

	for _, message := range input.SoundOutParam.Messages {
		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(uint(input.SoundOutParam.AgentStaffID))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(input.SoundOutParam.AgentID, message.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		if chatGroup.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			jobSeeker, err := i.jobSeekerRepository.FindByID(message.JobSeekerID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			fmt.Println("メッセージの送信に成功しました:", call)

			// エージェントから求職者へのチャットメッセージを作成
			chatMessage := entity.NewChatMessageWithJobSeeker(
				chatGroup.ID,
				null.NewInt(0, true), // エージェント
				message.Content,
				"",
				"",
				"",
				"",
				null.NewInt(0, false),
				"",
				null.NewInt(0, true),
				time.Now().In(time.UTC),
			)

			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				return output, err
			}
		} else {
			// メール送信
			err = utility.SendMailToSingle(
				i.sendgrid.APIKey,
				message.Subject,
				message.Content,
				entity.EmailUser{
					Name:  staff.StaffName,
					Email: staff.Email,
				},
				entity.EmailUser{
					Name:  message.Name,
					Email: message.Email,
				},
				nil,
			)

			fmt.Println("メール送信")

			if err != nil {
				fmt.Println("メール送信失敗")
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				message.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// プッシュ通知
	for _, caStaffID := range caIDList {
		redirectURL := os.Getenv("BASE_DOMAIN")
		contents := "新着CAタスクの依頼があります。"
		topic := "TaskCA"

		if caStaffID != 0 {
			caStaff, err := i.agentStaffRepository.FindByID(caStaffID)
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
				fmt.Println("WebPushの通知でエラー")
				fmt.Println(err)
			}
		}
	}

	return output, nil
}

/****************************************************************************************/
// 汎用系
//

type CreateTaskInBatchProcessingInput struct {
	Token       string
	CreateParam entity.CreateTaskInBatchProcessingParam
}

type CreateTaskInBatchProcessingOutput struct {
	OK bool
}

// エントリーフェーズのタスク処理
func (i *TaskInteractorImpl) CreateTaskInBatchProcessing(input CreateTaskInBatchProcessingInput) (CreateTaskInBatchProcessingOutput, error) {
	var (
		output   CreateTaskInBatchProcessingOutput
		err      error
		param    = input.CreateParam
		lastTask entity.TaskParam
	)

	/************ 1. 実行者の情報を取得 **************/

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	beforeActionLog := fmt.Sprintf(
		"\nタスク実行者: %s（%d）\n実行予定タスク: %d件\n",
		agentStaff.StaffName, agentStaff.ID, len(param.TaskList),
	)
	fmt.Println(beforeActionLog)
	/************ 2. Taskを1つずつ実行 **************/

	for index, task := range param.TaskList {
		var (
			// 次の選考フェーズ
			prevPhase = task.PrevPhaseCategory

			// 次のフェーズがエントリーフェズかどうか
			IsNextEntryPhase = prevPhase == null.NewInt(int64(entity.Entry), true)

			// 次のフェーズが書類選考フェズかどうか
			IsNextDocumentSelectionPhase = prevPhase == null.NewInt(int64(entity.DocumentSelection), true)

			// 次のフェーズが選考フェズかどうか
			IsNextSelectionPhase = prevPhase == null.NewInt(int64(entity.FirstSelection), true) ||
				prevPhase == null.NewInt(int64(entity.SecondSelection), true) ||
				prevPhase == null.NewInt(int64(entity.ThirdSelection), true) ||
				prevPhase == null.NewInt(int64(entity.FourthSelection), true) ||
				prevPhase == null.NewInt(int64(entity.FifthSelection), true) ||
				prevPhase == null.NewInt(int64(entity.FinalSelection), true)

			// 次のフェーズが内定保留フェズかどうか
			IsNextHoldJobOfferPhase = prevPhase == null.NewInt(int64(entity.HoldJobOffer), true)

			// 次のフェーズが内定承諾フェズかどうか
			IsNextAcceptJobOfferhase = prevPhase == null.NewInt(int64(entity.AcceptJobOffer), true)
		)

		// エントリーフェーズの処理
		if IsNextEntryPhase {
			err = createEntryPhaseTask(i, param, task)
			if err != nil {
				now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
				phaseStr, phaseSubStr := getStrTaskPhaseAndPhaseSub(task.PhaseCategory, task.PhaseSubCategory)

				// レスポンスで返すエラーメッセージ
				var responseErrorText = fmt.Sprintf("%v個目の「%s」の処理で失敗しました。", index+1, (phaseStr + "/" + phaseSubStr))

				errLog := fmt.Sprintf(
					"--------------------------------\n※%v\n\nフェーズ: %s\n発生時刻: %s\nエラー内容: %s\n--------------------------------",
					responseErrorText, (phaseStr + "/" + phaseSubStr), now, err.Error(),
				)
				fmt.Println(errLog)

				return output, errors.New(responseErrorText)
			}
		}

		// 書類選考フェーズの処理
		if IsNextDocumentSelectionPhase {
			err = createDocumentSelectionPhaseTask(i, param, task)
			if err != nil {
				now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
				phaseStr, phaseSubStr := getStrTaskPhaseAndPhaseSub(task.PhaseCategory, task.PhaseSubCategory)

				// レスポンスで返すエラーメッセージ
				var responseErrorText = fmt.Sprintf("%v個目の「%s」の処理で失敗しました。", index+1, (phaseStr + "/" + phaseSubStr))

				errLog := fmt.Sprintf(
					"--------------------------------\n※%v\n\nフェーズ: %s\n発生時刻: %s\nエラー内容: %s\n--------------------------------",
					responseErrorText, (phaseStr + "/" + phaseSubStr), now, err.Error(),
				)
				fmt.Println(errLog)

				return output, errors.New(responseErrorText)
			}
		}

		// 選考フェーズの処理
		if IsNextSelectionPhase {
			err = createSelectionPhaseTask(i, param, task)
			if err != nil {
				now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
				phaseStr, phaseSubStr := getStrTaskPhaseAndPhaseSub(task.PhaseCategory, task.PhaseSubCategory)

				// レスポンスで返すエラーメッセージ
				var responseErrorText = fmt.Sprintf("%v個目の「%s」の処理で失敗しました。", index+1, (phaseStr + "/" + phaseSubStr))

				errLog := fmt.Sprintf(
					"--------------------------------\n※%v\n\nフェーズ: %s\n発生時刻: %s\nエラー内容: %s\n--------------------------------",
					responseErrorText, (phaseStr + "/" + phaseSubStr), now, err.Error(),
				)
				fmt.Println(errLog)

				return output, errors.New(responseErrorText)
			}
		}

		// 内定保留フェーズの処理
		if IsNextHoldJobOfferPhase {
			err = createHoldJobOfferPhaseTask(i, param, task)
			if err != nil {
				now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
				phaseStr, phaseSubStr := getStrTaskPhaseAndPhaseSub(task.PhaseCategory, task.PhaseSubCategory)

				// レスポンスで返すエラーメッセージ
				var responseErrorText = fmt.Sprintf("%v個目の「%s」の処理で失敗しました。", index+1, (phaseStr + "/" + phaseSubStr))

				errLog := fmt.Sprintf(
					"--------------------------------\n※%v\n\nフェーズ: %s\n発生時刻: %s\nエラー内容: %s\n--------------------------------",
					responseErrorText, (phaseStr + "/" + phaseSubStr), now, err.Error(),
				)
				fmt.Println(errLog)

				return output, errors.New(responseErrorText)
			}
		}

		// 内定承諾フェーズの処理
		if IsNextAcceptJobOfferhase {
			err = createAcceptJobOfferPhaseTask(i, param, task)
			if err != nil {
				now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006/01/02 15:04:05")
				phaseStr, phaseSubStr := getStrTaskPhaseAndPhaseSub(task.PhaseCategory, task.PhaseSubCategory)

				// レスポンスで返すエラーメッセージ
				var responseErrorText = fmt.Sprintf("%v個目の「%s」の処理で失敗しました。", index+1, (phaseStr + "/" + phaseSubStr))

				errLog := fmt.Sprintf(
					"--------------------------------\n※%v\n\nフェーズ: %s\n発生時刻: %s\nエラー内容: %s\n--------------------------------",
					responseErrorText, (phaseStr + "/" + phaseSubStr), now, err.Error(),
				)
				fmt.Println(errLog)

				return output, errors.New(responseErrorText)
			}
		}

		// 最後のタスクの場合「」として取得する
		if (index + 1) == len(param.TaskList) {
			lastTask = task
		}
	}

	/************ 3. タスクの最終依頼時間と最終閲覧時間の更新（最後に処理したタスクを参照して実行） **************/

	// 依頼時間の更新
	if param.PrevStaffType == null.NewInt(int64(entity.CA), true) {
		if lastTask.StaffType == null.NewInt(int64(entity.CA), true) {
			// CAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateLastRequestAt(param.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// CAからRAにタスク作成する場合
			err := i.taskGroupRepository.UpdateCALastRequestAt(param.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		// RAからRAにタスク作成する場合
		if lastTask.StaffType == null.NewInt(int64(entity.RA), true) {
			err := i.taskGroupRepository.UpdateLastRequestAt(param.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(param.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ 4. タスク作成のプッシュ通知（最後に処理したタスクを参照して実行） **************/

	if lastTask.PhaseCategory.Valid && lastTask.PhaseSubCategory.Valid {
		// 終了タスク以外
		var IsNotClosePhase = !(lastTask.PhaseSubCategory.Int64 >= 90 && lastTask.PhaseSubCategory.Int64 <= 999)

		if IsNotClosePhase {
			// プッシュ通知
			firebaseID := ""
			contents := ""
			topic := ""
			redirectURL := os.Getenv("BASE_DOMAIN")

			// フェーズをテキストに変換
			phaseStr, phaseSubStr := getStrTaskPhaseAndPhaseSub(lastTask.PhaseCategory, lastTask.PhaseSubCategory)

			// リャンメンタスクを調べるためにタスクグループ取得
			taskGroup, err := i.taskGroupRepository.FindByID(param.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			var isCAStaff = lastTask.StaffType == null.NewInt(int64(entity.CA), true)
			var isRAStaff = lastTask.StaffType == null.NewInt(int64(entity.RA), true)

			if isCAStaff || (isRAStaff && taskGroup.IsDoubleSided) {
				contents = fmt.Sprintf("「%s」の新着タスクがあります。", (phaseStr + "/" + phaseSubStr))

				if isCAStaff {
					topic = "TaskCA"
				} else {
					topic = "TaskRA"
				}

				caStaff, err := i.agentStaffRepository.FindByID(param.CAStaffID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				firebaseID = caStaff.FirebaseID

			} else if isRAStaff {
				contents = fmt.Sprintf("「%s」の新着タスクがあります。", (phaseStr + "/" + phaseSubStr))
				topic = "TaskRA"

				raStaff, err := i.agentStaffRepository.FindByID(param.RAStaffID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				firebaseID = raStaff.FirebaseID
			}

			err = utility.WebPush(
				i.oneSignal.AppID,
				i.oneSignal.APIKey,
				firebaseID,
				"タスク通知",
				contents,
				topic,
				redirectURL,
			)
			if err != nil {
				fmt.Println("WebPushの通知")
				fmt.Println(err)
			}
		} else {
			// タスクが終了した場合　ヨミを「失注」に更新
			sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(param.JobSeekerID, param.JobInformationID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					fmt.Println("ヨミの登録がありません")
				} else {
					fmt.Println(err)
					return output, err
				}
			} else {
				sale.Accuracy = null.NewInt(entity.AccuracyFailure, true) // ヨミのみ「失注: 5」に更新
				err = i.saleRepository.Update(sale.ID, sale)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}
	}

	output.OK = true

	return output, nil
}

type GetSoundOutGroupListInput struct {
	AgentID              uint
	JobSeekerID          uint
	JobInformationIDList []uint
}

type GetSoundOutGroupListOutput struct {
	SoundOutGroupList       []*entity.SoundOutGroup
	AlreadyCreatedGroupList []*entity.SoundOutGroup
}

// 求人打診可能なグループと不可能なグループを取得
func (i *TaskInteractorImpl) GetSoundOutGroupList(input GetSoundOutGroupListInput) (GetSoundOutGroupListOutput, error) {
	var (
		output GetSoundOutGroupListOutput
		err    error
	)

	/************ 1. 求職者情報 + 求職者に紐づいた選考情報を取得 **************/

	jobSeeker, err := i.jobSeekerRepository.FindByID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// taskGroups
	seekerTaskGroups, err := i.taskGroupRepository.GetByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. 求人情報 + 求人に紐づいた選考情報を取得 **************/

	jobInformationList, err := i.jobInformationRepository.GetByIDList(input.JobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// taskGroups
	jobInfoTaskGroups, err := i.taskGroupRepository.GetByJobInformationIDList(input.JobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 3. 打診可能なグループと打診不可能なグループに振り分け **************/

	var LineID string
	if jobSeeker.AgentID == input.AgentID {
		LineID = jobSeeker.LineID
	}

	for _, jobInformation := range jobInformationList {
		soundOutGroup := entity.NewSoundOutGroup(
			jobSeeker.ID,
			jobInformation.ID,
			jobSeeker.UUID,
			jobInformation.UUID,
			jobInformation.AgentID,
			jobInformation.AgentStaffID,
			jobInformation.AgentName,
			jobInformation.StaffName,
			jobSeeker.AgentID,
			uint(jobSeeker.AgentStaffID.Int64),
			jobSeeker.AgentName,
			jobSeeker.StaffName,
			jobInformation.CompanyName,
			jobInformation.Title,
			jobSeeker.LastName,
			jobSeeker.FirstName,
			jobSeeker.LastFurigana,
			jobSeeker.FirstFurigana,
			jobSeeker.Email,
			LineID, // 自社求職者の場合はLINEIDを返す
			jobSeeker.LineActive,
			jobSeeker.StaffEmail,
			jobSeeker.StaffPhoneNumber,
			jobInformation.IsExternal,
		)

	jobInfoGroupLoop:
		for _, jobInfoGroup := range jobInfoTaskGroups {
			for _, seekerGroup := range seekerTaskGroups {
				if jobInfoGroup.ID == seekerGroup.ID {
					// 打診不可能
					output.AlreadyCreatedGroupList = append(output.AlreadyCreatedGroupList, soundOutGroup)
					continue jobInfoGroupLoop
				}
			}
		}

		// 打診不可リストに値があるか確認するフラグ
		exists := false

		for _, alreadyCreatedGroup := range output.AlreadyCreatedGroupList {
			if soundOutGroup == alreadyCreatedGroup {
				exists = true
				break
			}
		}

		if !exists {
			// 打診可能
			output.SoundOutGroupList = append(output.SoundOutGroupList, soundOutGroup)
		}
	}

	return output, nil
}

type CreateNextTaskAfterEntryPhaseInput struct {
	NextTaskParam entity.NextTaskParam
	AgentStaffID  uint
}

type CreateNextTaskAfterEntryPhaseOutput struct {
	OK bool
}

// エントリーフェーズのタスク処理
func (i *TaskInteractorImpl) CreateNextTaskAfterEntryPhase(input CreateNextTaskAfterEntryPhaseInput) (CreateNextTaskAfterEntryPhaseOutput, error) {
	var (
		output       CreateNextTaskAfterEntryPhaseOutput
		err          error
		task         *entity.Task
		nextTask     = input.NextTaskParam
		nextPhase    = nextTask.PhaseCategory
		nextPhaseSub = nextTask.PhaseSubCategory
		prevPhase    = nextTask.PrevPhaseCategory
		prevPhaseSub = nextTask.PrevPhaseSubCategory
		// urlText      = ""
	)

	if prevPhase != null.NewInt(int64(entity.Entry), true) {
		return output, errors.New("現在のフェーズがエントリーではありません")
	}

	/************ タスク情報の作成（選考継続時のタスク作成 or 通常タスク作成） **************/

	// アクション: 選考継続処理
	// エントリー/辞退意思確認中（選考継続の場合）->直前の○次選考のタスクに戻る
	if nextPhase == null.NewInt(int64(entity.Entry), true) && nextPhaseSub == null.NewInt(int64(entity.ContinueSelection), true) {
		// 選考継続の場合は辞退処理をする前のタスクを取得
		continueTask, err := i.taskRepository.FindLatestForContinue(nextTask.TaskGroupID, nextPhase)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 取得した過去のタスクを新規タスクとして作成
		task = entity.NewTask(
			continueTask.TaskGroupID,
			continueTask.PhaseCategory,    // 過去タスクのフェーズ
			continueTask.PhaseSubCategory, // 過去タスクのサブフェーズ
			continueTask.StaffType,
			continueTask.ExecutedStaffID,
			nextTask.Remarks,
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			nextTask.TalkAboutInInterview,
			nextTask.ScheduleCollectionCondition,
			nextTask.ExamGuideContent,
			false,
		)

	} else {
		// 次の選考が「エントリー/エントリー保留」or「エントリー/エントリー辞退」かを確認
		var IsDeclineOrHoldEntry = nextPhase == null.NewInt(int64(entity.Entry), true) && (nextPhaseSub == null.NewInt(int64(entity.DeclineEntry), true) || nextPhaseSub == null.NewInt(int64(entity.HoldEntry), true))

		// 次の選考が「書類選考/エントリー依頼」かを確認
		var IsRequestEntry = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) && nextPhaseSub == null.NewInt(int64(entity.RequestEntry), true)

		// 直前のタスクが応募意思確認中かを確認する
		if IsDeclineOrHoldEntry || IsRequestEntry {
			prevTask, err := i.taskRepository.FindLatestByGroupID(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// 直前のタスクが「エントリー/求人打診依頼」の場合は「エントリー/応募意思確認中」のタスクを作成する
			if prevTask.PhaseCategory == null.NewInt(int64(entity.Entry), true) && prevTask.PhaseSubCategory == null.NewInt(int64(entity.SoundOutJobInformation), true) {
				// 「エントリー/応募意思確認中」のタスクを作成
				confirmApplicationIntentionTask := entity.NewTask(
					nextTask.TaskGroupID,
					null.NewInt(int64(entity.Entry), true),
					null.NewInt(int64(entity.ConfirmApplicationIntention), true),
					null.NewInt(int64(entity.CA), true), // 担当者タイプ
					prevTask.ExecutedStaffID,
					"",
					prevTask.DeadlineDay,
					prevTask.DeadlineTime,
					"",
					"",
					"",
					false,
				)

				// タスク作成
				err = i.taskRepository.Create(confirmApplicationIntentionTask)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}

		// 実際に進めるタスクを作成
		task = entity.NewTask(
			nextTask.TaskGroupID,
			nextPhase,
			nextPhaseSub,
			nextTask.StaffType, // 担当者タイプ
			nextTask.ExecutedStaffID,
			nextTask.Remarks,
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			nextTask.TalkAboutInInterview,
			nextTask.ScheduleCollectionCondition,
			nextTask.ExamGuideContent,
			nextTask.IsCheckDoubleSided,
		)
	}

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 評価点の登録 **************/

	// 次が「マスクレジュメ合格」の場合
	var IsMaskResumePassNext = (prevPhase == null.NewInt(int64(entity.Entry), true) &&
		prevPhaseSub == null.NewInt(int64(entity.CollectResultOfMask), true)) &&
		(nextPhase == null.NewInt(int64(entity.Entry), true) &&
			nextPhaseSub == null.NewInt(int64(entity.SoundOutJobInformation), true))

	// 次が「マスクレジュメ不合格」の場合
	var IsMaskResumeFailNext = (prevPhase == null.NewInt(int64(entity.Entry), true) &&
		prevPhaseSub == null.NewInt(int64(entity.CollectResultOfMask), true)) &&
		(nextPhase == null.NewInt(int64(entity.Entry), true) &&
			nextPhaseSub == null.NewInt(int64(entity.EnterpriseNG), true))

	if IsMaskResumePassNext {
		// マスクレジュメ合格の場合　合格時の評価点を作成
		evaluationPoint := entity.NewEvaluationPoint(
			task.ID,
			nextTask.SelectionInformationID, // エントリーフェーズの時は選考フローを選択していないためnullが入る
			nextTask.JobSeekerID,
			nextTask.JobInformationID,
			nextTask.GoodPoint,
			nextTask.NGPoint,
			true,  // 選考通過のためtrue
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	} else if IsMaskResumeFailNext {
		// マスクレジュメ不合格の場合 不合格時の評価点の作成
		evaluationPoint := entity.NewEvaluationPoint(
			task.ID,
			nextTask.SelectionInformationID, // エントリーフェーズの時は選考フローを選択していないためnullが入る
			nextTask.JobSeekerID,
			nextTask.JobInformationID,
			nextTask.GoodPoint,
			nextTask.NGPoint,
			false, // 選考不合格のためfalse
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ マスクレジュメ打診時のメール送付 **************/

	// マスクレジュメを送信するかどうか
	var IsSendMaskResume = nextTask.TaskOption == null.NewInt(entity.OptionSendMessage, true) &&
		nextPhase == null.NewInt(int64(entity.Entry), true) &&
		nextPhaseSub == null.NewInt(int64(entity.CollectResultOfMask), true)

	if IsSendMaskResume {
		agentStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
		if err != nil {
			fmt.Println(err)
			fmt.Println(agentStaff)

			return output, err
		}

		err = utility.SendMailToMultiple(
			i.sendgrid.APIKey,
			nextTask.Subject,
			nextTask.Content,
			entity.EmailUser{
				Name:  agentStaff.StaffName,
				Email: agentStaff.Email,
			},
			nextTask.Tos,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ おすすめ応募書類の情報を登録する **************/

	// 次のタスクが「書類選考/推薦依頼」
	if nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true) {

		// 入力された送付をお勧めする書類の情報
		isR := nextTask.IsRecommendDocument

		// 送付をお勧めする書類の情報を取得
		isRecommend, err := i.taskIsRecommendDocumentRepository.FindByTaskID(task.ID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {

				// 作成する
				isRecommend = entity.NewTaskIsRecommendDocument(
					task.ID,
					isR.IsRecommendUploadResume,
					isR.IsRecommendUploadCV,
					isR.IsRecommendUploadRecommendation,
					isR.IsRecommendGeneratedResume,
					isR.IsRecommendGeneratedCV,
					isR.IsRecommendGeneratedRecommendation,
					isR.IsRecommendGeneratedMaskResume,
				)

				err := i.taskIsRecommendDocumentRepository.Create(isRecommend)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			// 更新する
			isRecommend.IsRecommendUploadResume = isR.IsRecommendUploadResume
			isRecommend.IsRecommendUploadCV = isR.IsRecommendUploadCV
			isRecommend.IsRecommendUploadRecommendation = isR.IsRecommendUploadRecommendation
			isRecommend.IsRecommendGeneratedResume = isR.IsRecommendGeneratedResume
			isRecommend.IsRecommendGeneratedCV = isR.IsRecommendGeneratedCV
			isRecommend.IsRecommendGeneratedRecommendation = isR.IsRecommendGeneratedRecommendation
			isRecommend.IsRecommendGeneratedMaskResume = isR.IsRecommendGeneratedMaskResume

			err := i.taskIsRecommendDocumentRepository.Update(isRecommend.ID, isRecommend)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

	}

	/************ 選考の候補日時を登録する **************/

	// 次のフェーズが「書類選考/推薦依頼」かどうか
	var IsNextPhaseRequestRecommendations = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsNextPhaseRequestRecommendations {
		// 登録前にすでに登録済みの候補日の情報を削除する
		// 削除するのは、カレンダー画面で日程を削除している可能性もあるため、作成処理の前に削除する
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(nextTask.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, possibleDate := range nextTask.PossibleDates {
			var taskID null.Int

			if possibleDate.TaskID.Valid {
				// TaskIDがある場合は登録済みの候補日時なので、既存のTaskIDを使用する
				taskID = possibleDate.TaskID
			} else {
				// TaskIDが無い場合は未登録の候補日時なので、新しいTaskIDを使用する
				taskID = null.NewInt(int64(task.ID), true)
			}

			jobSeekerSchedule := entity.NewJobSeekerSchedule(
				possibleDate.JobSeekerID,
				taskID,
				possibleDate.ScheduleType,
				possibleDate.Title,
				possibleDate.StartTime,
				possibleDate.EndTime,
				possibleDate.SeekerDescription,
				possibleDate.StaffDescription,
				possibleDate.IsShare,
				possibleDate.RepetitionCount,
			)

			err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ 両面タスクに更新する **************/

	// 次のタスクが「書類選考/書類推薦依頼」の場合
	var IsNextRequestRecommendationsn = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsNextRequestRecommendationsn && nextTask.IsCheckDoubleSided {
		// 両面タスクがTrueの場合は「task_groupテーブルのis_double_sidedカラム」をTrueで更新する。
		err = i.taskGroupRepository.UpdateIsDoubleSided(nextTask.TaskGroupID, true)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 求人票のURLを求職者へ送信（LINE or メール） **************/

	// 求人票を送信するかどうか（LINE or メール）
	var IsSendMessageForJobSeeker = nextTask.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		nextPhase == null.NewInt(int64(entity.Entry), true) &&
		nextPhaseSub == null.NewInt(int64(entity.ConfirmApplicationIntention), true)

	if IsSendMessageForJobSeeker {
		var message = nextTask.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(nextTask.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			fmt.Println("メッセージの送信に成功しました:", call)

			// エージェントから求職者へのチャットメッセージを作成
			chatMessage := entity.NewChatMessageWithJobSeeker(
				chatGroup.ID,
				null.NewInt(0, true), // エージェント
				message.Content,
				"",
				"",
				"",
				"",
				null.NewInt(0, false),
				"",
				null.NewInt(0, true),
				time.Now().In(time.UTC),
			)

			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

		} else {
			// メール送信
			err = utility.SendMailToSingle(
				i.sendgrid.APIKey,
				message.Subject,
				message.Content,
				entity.EmailUser{
					Name:  staff.StaffName,
					Email: staff.Email,
				},
				entity.EmailUser{
					Name:  (jobSeeker.LastName + jobSeeker.FirstName),
					Email: jobSeeker.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				nextTask.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスクの最終依頼時間と最終閲覧時間の更新 **************/

	// 依頼時間の更新
	if nextTask.PrevStaffType == null.NewInt(int64(entity.CA), true) {
		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			// CAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// CAからRAにタスク作成する場合
			err := i.taskGroupRepository.UpdateCALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		// RAからRAにタスク作成する場合
		if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスク作成のプッシュ通知 **************/

	// エントリーフェーズの終了タスク以外
	var IsNotClosePhase = nextPhaseSub != null.NewInt(entity.CloseEntryPhaseForUnlikely, true) &&
		nextPhaseSub != null.NewInt(entity.CloseEntryPhaseForWithoutSoundOut, true) &&
		nextPhaseSub != null.NewInt(entity.CloseEntryPhaseForDeclineEntry, true) &&
		nextPhaseSub != null.NewInt(entity.CloseEntryPhaseForRequestDecline, true) &&
		nextPhaseSub != null.NewInt(entity.CloseEntryPhase, true)

	// 終了タスク以外なら
	if IsNotClosePhase {
		// プッシュ通知
		firebaseID := ""
		contents := ""
		topic := ""
		redirectURL := os.Getenv("BASE_DOMAIN")

		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			contents = "新着CAタスクの依頼があります。"
			topic = "TaskCA"

			caStaff, err := i.agentStaffRepository.FindByID(nextTask.CAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = caStaff.FirebaseID

		} else if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			contents = "新着RAタスクの依頼があります。"
			topic = "TaskRA"

			raStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = raStaff.FirebaseID
		}

		err = utility.WebPush(
			i.oneSignal.AppID,
			i.oneSignal.APIKey,
			firebaseID,
			"タスク通知",
			contents,
			topic,
			redirectURL,
		)
		if err != nil {
			fmt.Println("WebPushの通知")
			fmt.Println(err)
		}
	} else {
		// タスクが終了した場合　ヨミを「失注」に更新
		sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(nextTask.JobSeekerID, nextTask.JobInformationID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				fmt.Println("ヨミの登録がありません")
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			sale.Accuracy = null.NewInt(entity.AccuracyFailure, true) // ヨミのみ「失注: 5」に更新
			err = i.saleRepository.Update(sale.ID, sale)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true

	return output, nil
}

type CreateNextTaskAfterDocumentSelectionPhaseInput struct {
	NextTaskParam entity.NextTaskParam
	AgentStaffID  uint
}

type CreateNextTaskAfterDocumentSelectionPhaseOutput struct {
	OK bool
}

// 書類選考フェーズのタスク処理
func (i *TaskInteractorImpl) CreateNextTaskAfterDocumentSelectionPhase(input CreateNextTaskAfterDocumentSelectionPhaseInput) (CreateNextTaskAfterDocumentSelectionPhaseOutput, error) {
	var (
		output       CreateNextTaskAfterDocumentSelectionPhaseOutput
		err          error
		task         *entity.Task
		nextTask     = input.NextTaskParam
		nextPhase    = nextTask.PhaseCategory
		nextPhaseSub = nextTask.PhaseSubCategory
		prevPhase    = nextTask.PrevPhaseCategory
		prevPhaseSub = nextTask.PrevPhaseSubCategory
	)

	if prevPhase != null.NewInt(int64(entity.DocumentSelection), true) {
		return output, errors.New("現在のフェーズが書類選考ではありません")
	}

	/************ タスク情報の作成（選考継続時のタスク作成 or 通常タスク作成） **************/

	if nextPhase == null.NewInt(int64(entity.DocumentSelection), true) && nextPhaseSub == null.NewInt(int64(entity.ContinueSelection), true) {
		// 選考継続の場合は辞退処理をする前のタスクを取得
		continueTask, err := i.taskRepository.FindLatestForContinue(nextTask.TaskGroupID, nextPhase)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 取得した過去のタスクを新規タスクとして作成
		task = entity.NewTask(
			continueTask.TaskGroupID,
			continueTask.PhaseCategory,    // 過去タスクのフェーズ
			continueTask.PhaseSubCategory, // 過去タスクのサブフェーズ
			continueTask.StaffType,
			continueTask.ExecutedStaffID,
			nextTask.Remarks,
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			nextTask.TalkAboutInInterview,
			nextTask.ScheduleCollectionCondition,
			nextTask.ExamGuideContent,
			false,
		)
	} else {
		// タスクを作成
		task = entity.NewTask(
			nextTask.TaskGroupID,
			nextPhase,
			nextPhaseSub,
			nextTask.StaffType, // 担当者タイプ
			nextTask.ExecutedStaffID,
			nextTask.Remarks,
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			nextTask.TalkAboutInInterview,
			nextTask.ScheduleCollectionCondition,
			nextTask.ExamGuideContent,
			nextTask.IsCheckDoubleSided,
		)
	}

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 評価点の登録 **************/

	// 次が「書類選考合格」の場合
	var IsDocumentSelectionPassNext = (prevPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		prevPhaseSub == null.NewInt(int64(entity.CollectResultOfDocumentSelection), true)) &&
		(isPhaseSelection(nextPhase) && (nextPhaseSub == null.NewInt(int64(entity.RequestGuidance), true) ||
			nextPhaseSub == null.NewInt(int64(entity.RequestCollectionOfSchedule), true) ||
			nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleAndSelectionDetail), true) ||
			nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleForSelection), true)))

		// 次が「書類選考不合格」の場合
	var IsDocumentSelectionFailNext = (nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.FailingNotification), true))

	if IsDocumentSelectionPassNext {
		var selectionInformationID null.Int

		if nextTask.SelectionFlowPatternID.Valid {
			selectionInformation, err := i.jobInfoSelectionInformationRepository.FindBySelectionFlowIDAndSelectionType(uint(nextTask.SelectionFlowPatternID.Int64), uint(entity.DocumentSelection))
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			selectionInformationID = null.NewInt(int64(selectionInformation.ID), true)
		}

		// 書類選考合格の場合　「合格時の評価点を作成 + 選考フローパターンの紐付け」
		evaluationPoint := entity.NewEvaluationPoint(
			task.ID,
			selectionInformationID, // エントリーフェーズの時は選考フローを選択していないためnullが入る
			nextTask.JobSeekerID,
			nextTask.JobInformationID,
			nextTask.GoodPoint,
			nextTask.NGPoint,
			true,  // 選考通過のためtrue
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 選考フローパターンの紐付け
		err = i.taskGroupRepository.UpdateSelectionFlowPatternID(nextTask.TaskGroupID, nextTask.SelectionFlowPatternID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

	} else if IsDocumentSelectionFailNext {
		// 書類選考不合格の場合 「不合格時の評価点の作成」
		evaluationPoint := entity.NewEvaluationPoint(
			task.ID,
			null.NewInt(0, false), // エントリーフェーズの時は選考フローを選択していないためnullが入る
			nextTask.JobSeekerID,
			nextTask.JobInformationID,
			nextTask.GoodPoint,
			nextTask.NGPoint,
			false, // 選考不合格のためfalse
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 推薦メールの送信 **************/

	// 推薦メールを送信するかどうか
	var IsSendRecommendationMail = nextTask.TaskOption == null.NewInt(entity.OptionSendMessage, true) &&
		prevPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		prevPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsSendRecommendationMail {
		// RA担当者情報の取得
		agentStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// RAから企業にメールを送信
		err = utility.SendMailToMultiple(
			i.sendgrid.APIKey,
			nextTask.Subject,
			nextTask.Content,
			entity.EmailUser{
				Name:  agentStaff.StaffName,
				Email: agentStaff.Email,
			},
			nextTask.Tos,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 選考の候補日時を登録する **************/

	// 次のフェーズが「書類選考/推薦依頼」かどうか
	var IsNextPhaseRequestRecommendations = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsNextPhaseRequestRecommendations {
		// 登録前にすでに登録済みの候補日の情報を削除する
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(nextTask.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, possibleDate := range nextTask.PossibleDates {
			var taskID null.Int

			if possibleDate.TaskID.Valid {
				// TaskIDがある場合は登録済みの候補日時なので、既存のTaskIDを使用する
				taskID = possibleDate.TaskID
			} else {
				// TaskIDが無い場合は未登録の候補日時なので、新しいTaskIDを使用する
				taskID = null.NewInt(int64(task.ID), true)
			}

			jobSeekerSchedule := entity.NewJobSeekerSchedule(
				possibleDate.JobSeekerID,
				taskID,
				possibleDate.ScheduleType,
				possibleDate.Title,
				possibleDate.StartTime,
				possibleDate.EndTime,
				possibleDate.SeekerDescription,
				possibleDate.StaffDescription,
				possibleDate.IsShare,
				possibleDate.RepetitionCount,
			)

			err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ 選考の確定日時を登録する **************/

	// 現在のフェーズが「書類選考/結果回収中」で次が「○次選考/日程・詳細確認依頼」or「○次選考/日程確認依頼」かどうか
	var IsNextSelectionPhase = prevPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		prevPhaseSub == null.NewInt(int64(entity.CollectResultOfDocumentSelection), true) &&
		(isPhaseSelection(nextPhase) && (nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleForSelection), true) ||
			nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleAndSelectionDetail), true)))

	if IsNextSelectionPhase {
		var taskID null.Int

		if nextTask.SelectionDate.TaskID.Valid {
			// TaskIDがある場合は登録済みの確定日時なので、既存のTaskIDを使用する
			taskID = nextTask.SelectionDate.TaskID
		} else {
			// TaskIDが無い場合は未登録の確定日時なので、新しいTaskIDを使用する
			taskID = null.NewInt(int64(task.ID), true)
		}

		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			nextTask.SelectionDate.JobSeekerID,
			taskID,
			nextTask.SelectionDate.ScheduleType,
			nextTask.SelectionDate.Title,
			nextTask.SelectionDate.StartTime,
			nextTask.SelectionDate.EndTime,
			nextTask.SelectionDate.SeekerDescription,
			nextTask.SelectionDate.StaffDescription,
			nextTask.SelectionDate.IsShare,
			nextTask.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 両面タスクに更新する **************/

	// 次のタスクが「書類選考/書類推薦依頼」の場合
	var IsNextRequestRecommendationsn = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsNextRequestRecommendationsn && nextTask.IsCheckDoubleSided {
		// 両面タスクがTrueの場合は「task_groupテーブルのis_double_sidedカラム」をTrueで更新する。
		err = i.taskGroupRepository.UpdateIsDoubleSided(nextTask.TaskGroupID, true)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 応募書類アップロードURLを求職者へ送信（LINE or メール） **************/

	// 求人票を送信するかどうか（LINE or メール）
	var IsSendMessageForJobSeeker = nextTask.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.PrepareDocument), true)

	if IsSendMessageForJobSeeker {
		var message = nextTask.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(nextTask.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			fmt.Println("メッセージの送信に成功しました:", call)

			// エージェントから求職者へのチャットメッセージを作成
			chatMessage := entity.NewChatMessageWithJobSeeker(
				chatGroup.ID,
				null.NewInt(0, true), // エージェント
				message.Content,
				"",
				"",
				"",
				"",
				null.NewInt(0, false),
				"",
				null.NewInt(0, true),
				time.Now().In(time.UTC),
			)

			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

		} else {
			// メール送信
			err = utility.SendMailToSingle(
				i.sendgrid.APIKey,
				message.Subject,
				message.Content,
				entity.EmailUser{
					Name:  staff.StaffName,
					Email: staff.Email,
				},
				entity.EmailUser{
					Name:  (jobSeeker.LastName + jobSeeker.FirstName),
					Email: jobSeeker.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				nextTask.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスクの最終依頼時間と最終閲覧時間の更新 **************/

	// 依頼時間の更新
	if nextTask.PrevStaffType == null.NewInt(int64(entity.CA), true) {
		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			// CAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// CAからRAにタスク作成する場合
			err := i.taskGroupRepository.UpdateCALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		// RAからRAにタスク作成する場合
		if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスク作成のプッシュ通知 **************/

	// エントリーフェーズの終了タスク以外
	var IsNotClosePhase = nextPhaseSub != null.NewInt(entity.CloseDocumentPhaseForDeclineEntry, true) &&
		nextPhaseSub != null.NewInt(entity.CloseDocumentPhaseForFailingNotice, true) &&
		nextPhaseSub != null.NewInt(entity.CloseDocumentPhaseForDrop, true) &&
		nextPhaseSub != null.NewInt(entity.CloseDocumentPhaseForRequestDecline, true) &&
		nextPhaseSub != null.NewInt(entity.CloseDocumentPhase, true)

	if IsNotClosePhase {
		// プッシュ通知
		firebaseID := ""
		contents := ""
		topic := ""
		redirectURL := os.Getenv("BASE_DOMAIN")

		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			contents = "新着CAタスクの依頼があります。"
			topic = "TaskCA"

			caStaff, err := i.agentStaffRepository.FindByID(nextTask.CAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = caStaff.FirebaseID

		} else if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			contents = "新着RAタスクの依頼があります。"
			topic = "TaskRA"

			raStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = raStaff.FirebaseID
		}

		err = utility.WebPush(
			i.oneSignal.AppID,
			i.oneSignal.APIKey,
			firebaseID,
			"タスク通知",
			contents,
			topic,
			redirectURL,
		)
		if err != nil {
			fmt.Println("WebPushの通知")
			fmt.Println(err)
		}
	} else {
		// タスクが終了した場合　ヨミを「失注」に更新
		sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(nextTask.JobSeekerID, nextTask.JobInformationID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				fmt.Println("ヨミの登録がありません")
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			sale.Accuracy = null.NewInt(entity.AccuracyFailure, true) // ヨミのみ「失注: 5」に更新
			err = i.saleRepository.Update(sale.ID, sale)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true

	return output, nil
}

type CreateNextTaskAfterSelectionPhaseInput struct {
	NextTaskParam entity.NextTaskParam
	AgentStaffID  uint
}

type CreateNextTaskAfterSelectionPhaseOutput struct {
	OK bool
}

// エントリーフェーズのタスク処理
func (i *TaskInteractorImpl) CreateNextTaskAfterSelectionPhase(input CreateNextTaskAfterSelectionPhaseInput) (CreateNextTaskAfterSelectionPhaseOutput, error) {
	var (
		output       CreateNextTaskAfterSelectionPhaseOutput
		err          error
		task         *entity.Task
		nextTask     = input.NextTaskParam
		nextPhase    = nextTask.PhaseCategory
		nextPhaseSub = nextTask.PhaseSubCategory
		prevPhase    = nextTask.PrevPhaseCategory
		prevPhaseSub = nextTask.PrevPhaseSubCategory
	)

	// 選考フェーズでなかった場合はエラー
	if !(int64(entity.FirstSelection) <= prevPhase.Int64) && !(prevPhase.Int64 >= int64(entity.FinalSelection)) {
		return output, errors.New("現在フェーズが選考フェーズではありません")
	}

	/************ 新リスケ処理（リスケ + 削除） 元のリスケ処理も残しておく **************/

	if nextTask.TaskOption == null.NewInt(entity.OptionRescheduleAndDelete, true) {
		// 直前の結果回収依頼までの日程を削除 + リスケタスクの作成
		latestSelectTask, err := i.taskRepository.FindLatestByGroupIDAndCollectResultPhase(nextTask.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 直前の結果回収依頼までの日程を削除
		err = i.jobSeekerScheduleRepository.DeleteScheduleInTaskGroupAboveTaskID(latestSelectTask.ID, nextTask.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// リスケタスクを作成
		rescheduleTask := entity.NewTask(
			nextTask.TaskGroupID,
			nextPhase,
			null.NewInt(entity.RescheduleDeleteOfSelection, true),
			nextTask.StaffType, // 担当者タイプ
			nextTask.ExecutedStaffID,
			"",
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			"",
			"",
			"",
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(rescheduleTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ タスク情報の作成（選考継続時のタスク作成 or 通常タスク作成（+スキップタスク）） **************/

	if (isPhaseSelection(nextPhase)) &&
		nextPhaseSub == null.NewInt(int64(entity.ContinueSelection), true) {

		// 選考継続の場合は辞退処理をする前のタスクを取得
		continueTask, err := i.taskRepository.FindLatestForContinue(nextTask.TaskGroupID, nextPhase)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 取得した過去のタスクを新規タスクとして作成
		task = entity.NewTask(
			continueTask.TaskGroupID,
			continueTask.PhaseCategory,    // 過去タスクのフェーズ
			continueTask.PhaseSubCategory, // 過去タスクのサブフェーズ
			continueTask.StaffType,
			continueTask.ExecutedStaffID,
			nextTask.Remarks,
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			nextTask.TalkAboutInInterview,
			nextTask.ScheduleCollectionCondition,
			nextTask.ExamGuideContent,
			false,
		)
	} else {
		// タスクを作成
		task = entity.NewTask(
			nextTask.TaskGroupID,
			nextPhase,
			nextPhaseSub,
			nextTask.StaffType, // 担当者タイプ
			nextTask.ExecutedStaffID,
			nextTask.Remarks,
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			nextTask.TalkAboutInInterview,
			nextTask.ScheduleCollectionCondition,
			nextTask.ExamGuideContent,
			false,
		)
	}

	// 選考スキップタスクの場合
	if nextTask.TaskOption == null.NewInt(entity.OptionSkipSelection, true) {
		prevTask, err := i.taskRepository.FindLatestByGroupID(nextTask.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// スキップされるフェーズは現在までの最新タスクフェーズの次のフェーズ
		skipPhase := prevTask.PhaseCategory.Int64 + 1

		// 作成されるタスクが内定保留タスクの場合はスキップフェーズを最終選考に変更
		if nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) {
			skipPhase = int64(entity.FinalSelection)
		}

		skipTask := entity.NewTask(
			nextTask.TaskGroupID,
			null.NewInt(skipPhase, true),            // スキップされるフェーズ
			null.NewInt(entity.SelectionSkip, true), // 選考スキップ
			null.NewInt(int64(entity.CA), true),     // 担当者タイプ
			nextTask.CAStaffID,
			"",
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			"",
			"",
			"",
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(skipTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 評価点の登録（※再設定の場合は合格として扱う） **************/

	// 現在が「⚪︎次選考の結果回収中」（次が「⚪︎次選考の合格・不合格」）
	var IsSelectionCollectResultNow = prevPhaseSub == null.NewInt(int64(entity.RequestCollectSelectionThought), true)

	// 再選考の場合
	var IsReSelection = nextPhase == prevPhase &&
		(nextPhaseSub == null.NewInt(int64(entity.RequestGuidance), true) ||
			nextPhaseSub == null.NewInt(int64(entity.RequestCollectionOfSchedule), true) ||
			nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleForSelection), true))

	if IsSelectionCollectResultNow {
		// 再設定の有無
		var isReInterview = false

		// 合格可否
		var isPass = false

		// 再設定のタスクの場合は「true」に変更
		if IsReSelection {
			isReInterview = true
		}

		// 現在のフェーズと次のフェーズが同一でなければisPassを「true」に変更
		if !(nextPhase == prevPhase) {
			isPass = true
		}

		// 評価点（合格）
		evaluationPoint := entity.NewEvaluationPoint(
			task.ID,
			nextTask.SelectionInformationID,
			nextTask.JobSeekerID,
			nextTask.JobInformationID,
			nextTask.GoodPoint,
			nextTask.NGPoint,
			isPass,        // 選考通過のためtrue
			isReInterview, // 再面接かどうかで値が変化
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 選考の候補日時を登録する **************/

	// 次のフェーズが「○次選考/日程調整依頼」かどうか
	var IsNextPhaseRequestScheduleAdjustment = isPhaseSelection(nextPhase) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestScheduleAdjustmentForSelection), true)

	if IsNextPhaseRequestScheduleAdjustment {
		// 登録前にすでに登録済みの候補日の情報を削除する
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(nextTask.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, possibleDate := range nextTask.PossibleDates {
			var taskID null.Int

			if possibleDate.TaskID.Valid {
				// TaskIDがある場合は登録済みの候補日時なので、既存のTaskIDを使用する
				taskID = possibleDate.TaskID
			} else {
				// TaskIDが無い場合は未登録の候補日時なので、新しいTaskIDを使用する
				taskID = null.NewInt(int64(task.ID), true)
			}

			jobSeekerSchedule := entity.NewJobSeekerSchedule(
				possibleDate.JobSeekerID,
				taskID,
				possibleDate.ScheduleType,
				possibleDate.Title,
				possibleDate.StartTime,
				possibleDate.EndTime,
				possibleDate.SeekerDescription,
				possibleDate.StaffDescription,
				possibleDate.IsShare,
				possibleDate.RepetitionCount,
			)

			err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ 選考の確定日時を登録する **************/

	// 現在のフェーズが「書類選考/結果回収中」で次が「○次選考/日程・詳細確認依頼」or「○次選考/日程確認依頼」かどうか
	var IsNextSelectionPhase = (isPhaseSelection(nextPhase) && (nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleForSelection), true) ||
		nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleAndSelectionDetail), true)))

	if IsNextSelectionPhase {
		var taskID null.Int

		if nextTask.SelectionDate.TaskID.Valid {
			// TaskIDがある場合は登録済みの確定日時なので、既存のTaskIDを使用する
			taskID = nextTask.SelectionDate.TaskID
		} else {
			// TaskIDが無い場合は未登録の確定日時なので、新しいTaskIDを使用する
			taskID = null.NewInt(int64(task.ID), true)
		}

		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			nextTask.SelectionDate.JobSeekerID,
			taskID,
			nextTask.SelectionDate.ScheduleType,
			nextTask.SelectionDate.Title,
			nextTask.SelectionDate.StartTime,
			nextTask.SelectionDate.EndTime,
			nextTask.SelectionDate.SeekerDescription,
			nextTask.SelectionDate.StaffDescription,
			nextTask.SelectionDate.IsShare,
			nextTask.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 選考の確定日時と当日詳細を更新する **************/

	// 次のフェーズが「○次選考/詳細案内依頼」かどうか
	var IsNextRequestDetailed = (isPhaseSelection(nextPhase) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestDetailedInformation), true))

	if IsNextRequestDetailed {
		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			nextTask.SelectionDate.JobSeekerID,
			nextTask.SelectionDate.TaskID,
			nextTask.SelectionDate.ScheduleType,
			nextTask.SelectionDate.Title,
			nextTask.SelectionDate.StartTime,
			nextTask.SelectionDate.EndTime,
			nextTask.SelectionDate.SeekerDescription,
			nextTask.SelectionDate.StaffDescription,
			nextTask.SelectionDate.IsShare,
			nextTask.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Update(nextTask.SelectionDate.ID, jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ リスケ時の処理（確定日時と同じ日時にリスケ日程を作成） **************/

	// 前のフェーズが「○次選考/日程・詳細確認中 or 日程確認中」で次のタスクオプションががリスケの場合

	var IsResuchedule = isPhaseSelection(prevPhase) && (prevPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleForSelection), true) ||
		prevPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleAndSelectionDetail), true)) &&
		nextTask.TaskOption == null.NewInt(entity.OptionReschedule, true)

	if IsResuchedule {
		// 確定日時とリスケのリストを取得(id降順で取得)
		scheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupIDAndScheduleType(nextTask.TaskGroupID, uint(entity.ScheduleTypeByEnterprise))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		var filteringScheduleList []*entity.JobSeekerSchedule

		// リスケになっていない同じ選考フェーズの確定日時だけ抜き出す
		for _, schedule := range scheduleList {
			if !schedule.RescheduleID.Valid && schedule.PhaseCategory == nextTask.PhaseCategory {
				fmt.Println(schedule.Title)
				filteringScheduleList = append(filteringScheduleList, schedule)
			}
		}

		// 確定日時が取得でいているかを確認
		if len(filteringScheduleList) > 0 {
			// 確定日時
			var selectionSchedule = filteringScheduleList[0]

			//　リスケデータを作成
			reschedule := entity.NewJobSeekerReschedule(selectionSchedule.ID, task.ID)

			err = i.jobSeekerRescheduleRepository.Create(reschedule)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/**

	id
	reschedule_id リスケになった確定日時のレコードのID
	task_id　　　リスケ処理を実行したタスクのID

	リスケ発生時にCrateする
	Get時はテーブルを外部結合で取得してrescheduleIDがnullでない場合はリスケした日程と判断する
	リスケタスクのDeleteの時はリスケテーブルのレコードも消えるので日程のテーブルでは消えたリスケテーブルのデータに紐づくデータがリスケではなくなる

	*/

	/************ メール or LINEの送付 **************/

	// 〇〇選考/案内依頼 → 求職者対応中
	// 〇〇選考/候補日回収依頼 → 候補日回収中
	// 〇〇選考/日程案内依頼 → 日程確認中
	// 〇〇選考/日程・詳細案内依頼 → 日程・詳細確認中
	// 〇〇選考/詳細案内依頼 → 前日確認依頼
	// 〇〇選考/前日確認依頼 → 選考所感の回収

	// メール or LINEの送信する
	var IsSendMessage = nextTask.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		isPhaseSelection(nextPhase) &&
		(nextPhaseSub == null.NewInt(int64(entity.JobSeekerSupport), true) ||
			nextPhaseSub == null.NewInt(int64(entity.CollectingSchedule), true) ||
			nextPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleForSelection), true) ||
			nextPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleAndSelectionDetail), true) ||
			nextPhaseSub == null.NewInt(int64(entity.ConfirmDayBefore), true) ||
			nextPhaseSub == null.NewInt(int64(entity.CollectSelectionThought), true))

	if IsSendMessage {
		var message = nextTask.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(nextTask.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			fmt.Println("メッセージの送信に成功しました:", call)

			// エージェントから求職者へのチャットメッセージを作成
			chatMessage := entity.NewChatMessageWithJobSeeker(
				chatGroup.ID,
				null.NewInt(0, true), // エージェント
				message.Content,
				"",
				"",
				"",
				"",
				null.NewInt(0, false),
				"",
				null.NewInt(0, true),
				time.Now().In(time.UTC),
			)

			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

		} else {
			// メール送信
			err = utility.SendMailToSingle(
				i.sendgrid.APIKey,
				message.Subject,
				message.Content,
				entity.EmailUser{
					Name:  staff.StaffName,
					Email: staff.Email,
				},
				entity.EmailUser{
					Name:  (jobSeeker.LastName + jobSeeker.FirstName),
					Email: jobSeeker.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				nextTask.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスクの最終依頼時間と最終閲覧時間の更新 **************/

	// 依頼時間の更新
	if nextTask.PrevStaffType == null.NewInt(int64(entity.CA), true) {
		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			// CAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// CAからRAにタスク作成する場合
			err := i.taskGroupRepository.UpdateCALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		// RAからRAにタスク作成する場合
		if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスク作成のプッシュ通知 **************/

	// 選考フェーズの終了タスク以外
	var IsNotClosePhase = nextPhaseSub != null.NewInt(entity.CloseSelectionPhaseForFailingNoticeOfSelection, true) &&
		nextPhaseSub != null.NewInt(entity.CloseSelectionPhaseForDrop, true) &&
		nextPhaseSub != null.NewInt(entity.CloseSelectionPhaseForRequestDecline, true) &&
		nextPhaseSub != null.NewInt(entity.CloseSelectionPhase, true)

	if IsNotClosePhase {
		// プッシュ通知
		firebaseID := ""
		contents := ""
		topic := ""
		redirectURL := os.Getenv("BASE_DOMAIN")

		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			contents = "新着CAタスクの依頼があります。"
			topic = "TaskCA"

			caStaff, err := i.agentStaffRepository.FindByID(nextTask.CAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = caStaff.FirebaseID

		} else if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			contents = "新着RAタスクの依頼があります。"
			topic = "TaskRA"

			raStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = raStaff.FirebaseID
		}

		err = utility.WebPush(
			i.oneSignal.AppID,
			i.oneSignal.APIKey,
			firebaseID,
			"タスク通知",
			contents,
			topic,
			redirectURL,
		)
		if err != nil {
			fmt.Println("WebPushの通知")
			fmt.Println(err)
		}
	} else {
		// タスクが終了した場合　ヨミを「失注」に更新
		sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(nextTask.JobSeekerID, nextTask.JobInformationID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				fmt.Println("ヨミの登録がありません")
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			sale.Accuracy = null.NewInt(entity.AccuracyFailure, true) // ヨミのみ「失注: 5」に更新
			err = i.saleRepository.Update(sale.ID, sale)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true

	return output, nil
}

// 内定辞退フェーズのタスク処理
type CreateNextTaskAfterDeclinePhaseInput struct {
	NextTaskParam entity.NextTaskParam
	AgentStaffID  uint
}

type CreateNextTaskAfterDeclinePhaseOutput struct {
	OK bool
}

// 辞退フェーズのタスク処理
func (i *TaskInteractorImpl) CreateNextTaskAfterDeclinePhase(input CreateNextTaskAfterDeclinePhaseInput) (CreateNextTaskAfterDeclinePhaseOutput, error) {
	var (
		output   CreateNextTaskAfterDeclinePhaseOutput
		err      error
		nextTask = input.NextTaskParam
	)

	// 選考フェーズでなかった場合はエラー
	if nextTask.PrevPhaseCategory != null.NewInt(int64(entity.Decline), true) {
		return output, errors.New("現在フェーズが辞退フェーズではありません")
	}

	task := entity.NewTask(
		nextTask.TaskGroupID,
		nextTask.PhaseCategory,
		nextTask.PhaseSubCategory,
		nextTask.StaffType, // 担当者タイプ
		nextTask.ExecutedStaffID,
		nextTask.Remarks,
		nextTask.DeadlineDay,
		nextTask.DeadlineTime,
		nextTask.TalkAboutInInterview,
		nextTask.ScheduleCollectionCondition,
		nextTask.ExamGuideContent,
		false,
	)

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// プッシュ通知
	firebaseID := ""
	contents := ""
	topic := ""
	redirectURL := os.Getenv("BASE_DOMAIN")

	if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
		contents = "新着CAタスクの依頼があります。"
		topic = "TaskCA"

		caStaff, err := i.agentStaffRepository.FindByID(nextTask.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		firebaseID = caStaff.FirebaseID

	} else if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
		contents = "新着RAタスクの依頼があります。"
		topic = "TaskRA"

		raStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		firebaseID = raStaff.FirebaseID
	}

	err = utility.WebPush(
		i.oneSignal.AppID,
		i.oneSignal.APIKey,
		firebaseID,
		"タスク通知",
		contents,
		topic,
		redirectURL,
	)
	if err != nil {
		fmt.Println("WebPushの通知")
		fmt.Println(err)
	}

	output.OK = true

	return output, nil
}

// 内定保留フェーズのタスク処理
type CreateNextTaskAfterHoldJobOfferPhaseInput struct {
	NextTaskParam entity.NextTaskParam
	AgentStaffID  uint
}

type CreateNextTaskAfterHoldJobOfferPhaseOutput struct {
	OK bool
}

// エントリーフェーズのタスク処理
func (i *TaskInteractorImpl) CreateNextTaskAfterHoldJobOfferPhase(input CreateNextTaskAfterHoldJobOfferPhaseInput) (CreateNextTaskAfterHoldJobOfferPhaseOutput, error) {
	var (
		output       CreateNextTaskAfterHoldJobOfferPhaseOutput
		err          error
		nextTask     = input.NextTaskParam
		nextPhase    = nextTask.PhaseCategory
		nextPhaseSub = nextTask.PhaseSubCategory
		prevPhase    = nextTask.PrevPhaseCategory
		// prevPhaseSub = nextTask.PrevPhaseSubCategory
	)

	// 選考フェーズでなかった場合はエラー
	if prevPhase != null.NewInt(int64(entity.HoldJobOffer), true) {
		return output, errors.New("現在フェーズが内定保留フェーズではありません")
	}

	/************ 新リスケ処理（リスケ + 削除） 元のリスケ処理も残しておく **************/

	if nextTask.TaskOption == null.NewInt(entity.OptionRescheduleAndDelete, true) {
		// 直前の結果回収依頼までの日程を削除 + リスケタスクの作成
		latestSelectTask, err := i.taskRepository.FindLatestByGroupIDAndCollectResultPhase(nextTask.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 直前の結果回収依頼までの日程を削除
		err = i.jobSeekerScheduleRepository.DeleteScheduleInTaskGroupAboveTaskID(latestSelectTask.ID, nextTask.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// リスケタスクを作成
		rescheduleTask := entity.NewTask(
			nextTask.TaskGroupID,
			nextPhase,
			null.NewInt(entity.RescheduleDeleteOfSelection, true),
			nextTask.StaffType, // 担当者タイプ
			nextTask.ExecutedStaffID,
			"",
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			"",
			"",
			"",
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(rescheduleTask)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	task := entity.NewTask(
		nextTask.TaskGroupID,
		nextPhase,
		nextPhaseSub,
		nextTask.StaffType, // 担当者タイプ
		nextTask.ExecutedStaffID,
		nextTask.Remarks,
		nextTask.DeadlineDay,
		nextTask.DeadlineTime,
		nextTask.TalkAboutInInterview,
		nextTask.ScheduleCollectionCondition,
		nextTask.ExamGuideContent,
		false,
	)

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 内定承諾/内定承諾の場合はヨミ情報を更新する
	if nextPhase == null.NewInt(int64(entity.AcceptJobOffer), true) && nextPhaseSub == null.NewInt(int64(entity.Accept), true) {
		sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(nextTask.JobSeekerID, nextTask.JobInformationID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				fmt.Println("ヨミの登録がありません")
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			sale.Accuracy = null.NewInt(entity.AccuracyAccept, true) // ヨミのみ「内定承諾: 0」に更新
			err = i.saleRepository.Update(sale.ID, sale)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ 選考の候補日時を登録する **************/

	// 次のフェーズが「内定保留/ファー面談の日程調整依頼」かどうか
	var IsNextRequestScheduleAdjustmentForHold = nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) && nextPhaseSub == null.NewInt(int64(entity.RequestScheduleAdjustment), true)

	if IsNextRequestScheduleAdjustmentForHold {
		// 登録前にすでに登録済みの候補日の情報を削除する
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(nextTask.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, possibleDate := range nextTask.PossibleDates {
			var taskID null.Int

			if possibleDate.TaskID.Valid {
				// TaskIDがある場合は登録済みの候補日時なので、既存のTaskIDを使用する
				taskID = possibleDate.TaskID
			} else {
				// TaskIDが無い場合は未登録の候補日時なので、新しいTaskIDを使用する
				taskID = null.NewInt(int64(task.ID), true)
			}

			jobSeekerSchedule := entity.NewJobSeekerSchedule(
				possibleDate.JobSeekerID,
				taskID,
				possibleDate.ScheduleType,
				possibleDate.Title,
				possibleDate.StartTime,
				possibleDate.EndTime,
				possibleDate.SeekerDescription,
				possibleDate.StaffDescription,
				possibleDate.IsShare,
				possibleDate.RepetitionCount,
			)

			err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ 選考の確定日時を登録する **************/

	// 次のフェーズが「内定保留/オファー面談の詳細案内依頼」かどうか
	var IsNextRequestDetailedInformationForHold = nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) && nextPhaseSub == null.NewInt(int64(entity.RequestDetailedInformationForHold), true)

	if IsNextRequestDetailedInformationForHold {
		var taskID null.Int

		if nextTask.SelectionDate.TaskID.Valid {
			// TaskIDがある場合は登録済みの確定日時なので、既存のTaskIDを使用する
			taskID = nextTask.SelectionDate.TaskID
		} else {
			// TaskIDが無い場合は未登録の確定日時なので、新しいTaskIDを使用する
			taskID = null.NewInt(int64(task.ID), true)
		}

		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			nextTask.SelectionDate.JobSeekerID,
			taskID,
			nextTask.SelectionDate.ScheduleType,
			nextTask.SelectionDate.Title,
			nextTask.SelectionDate.StartTime,
			nextTask.SelectionDate.EndTime,
			nextTask.SelectionDate.SeekerDescription,
			nextTask.SelectionDate.StaffDescription,
			nextTask.SelectionDate.IsShare,
			nextTask.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ メール or LINEの送付 **************/

	// 内定保留/オファー面談の候補日回収依頼 → オファー面談の日程回収中

	// メール or LINEの送信する
	var IsSendMessage = nextTask.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) && (nextPhaseSub == null.NewInt(int64(entity.CollectionOfSchedule), true) || nextPhaseSub == null.NewInt(int64(entity.ConfirmsSchedule), true))

	if IsSendMessage {
		var message = nextTask.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(nextTask.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			fmt.Println("メッセージの送信に成功しました:", call)

			// エージェントから求職者へのチャットメッセージを作成
			chatMessage := entity.NewChatMessageWithJobSeeker(
				chatGroup.ID,
				null.NewInt(0, true), // エージェント
				message.Content,
				"",
				"",
				"",
				"",
				null.NewInt(0, false),
				"",
				null.NewInt(0, true),
				time.Now().In(time.UTC),
			)

			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

		} else {
			// メール送信
			err = utility.SendMailToSingle(
				i.sendgrid.APIKey,
				message.Subject,
				message.Content,
				entity.EmailUser{
					Name:  staff.StaffName,
					Email: staff.Email,
				},
				entity.EmailUser{
					Name:  (jobSeeker.LastName + jobSeeker.FirstName),
					Email: jobSeeker.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				nextTask.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// 依頼時間の更新
	if nextTask.PrevStaffType == null.NewInt(int64(entity.CA), true) {
		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			// CAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// CAからRAにタスク作成する場合
			err := i.taskGroupRepository.UpdateCALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		// RAからRAにタスク作成する場合
		if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスク作成のプッシュ通知 **************/

	// エントリーフェーズの終了タスク以外
	var IsNotClosePhase = nextPhaseSub != null.NewInt(entity.CloseHoldJobOfferPhaseForDeclineJobOffer, true) &&
		nextPhaseSub != null.NewInt(entity.CloseHoldJobOfferPhaseForRequestDeclineOfSelection, true) &&
		nextPhaseSub != null.NewInt(entity.HoldJobOfferPhaseClose, true)

	if IsNotClosePhase {
		// プッシュ通知
		firebaseID := ""
		contents := ""
		topic := ""
		redirectURL := os.Getenv("BASE_DOMAIN")

		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			contents = "新着CAタスクの依頼があります。"
			topic = "TaskCA"

			caStaff, err := i.agentStaffRepository.FindByID(nextTask.CAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = caStaff.FirebaseID

		} else if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			contents = "新着RAタスクの依頼があります。"
			topic = "TaskRA"

			raStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = raStaff.FirebaseID
		}

		err = utility.WebPush(
			i.oneSignal.AppID,
			i.oneSignal.APIKey,
			firebaseID,
			"タスク通知",
			contents,
			topic,
			redirectURL,
		)
		if err != nil {
			fmt.Println("WebPushの通知")
			fmt.Println(err)
		}
	} else {
		// タスクが終了した場合　ヨミを「失注」に更新
		sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(nextTask.JobSeekerID, nextTask.JobInformationID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				fmt.Println("ヨミの登録がありません")
			} else {
				fmt.Println(err)
				return output, err
			}
		} else {
			sale.Accuracy = null.NewInt(entity.AccuracyFailure, true) // ヨミのみ「失注: 5」に更新
			err = i.saleRepository.Update(sale.ID, sale)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true

	return output, nil
}

// 内定承諾フェーズのタスク処理
type CreateNextTaskAfterAcceptJobOfferPhaseInput struct {
	NextTaskParam entity.NextTaskParam
	AgentStaffID  uint
}

type CreateNextTaskAfterAcceptJobOfferPhaseOutput struct {
	OK bool
}

// 内定承諾フェーズのタスク処理
func (i *TaskInteractorImpl) CreateNextTaskAfterAcceptJobOfferPhase(input CreateNextTaskAfterAcceptJobOfferPhaseInput) (CreateNextTaskAfterAcceptJobOfferPhaseOutput, error) {
	var (
		output   CreateNextTaskAfterAcceptJobOfferPhaseOutput
		err      error
		nextTask = input.NextTaskParam
	)

	// 選考フェーズでなかった場合はエラー
	if nextTask.PrevPhaseCategory != null.NewInt(int64(entity.AcceptJobOffer), true) {
		return output, errors.New("現在フェーズが内定承諾フェーズではありません")
	}

	task := entity.NewTask(
		nextTask.TaskGroupID,
		nextTask.PhaseCategory,
		nextTask.PhaseSubCategory,
		nextTask.StaffType, // 担当者タイプ
		nextTask.ExecutedStaffID,
		nextTask.Remarks,
		nextTask.DeadlineDay,
		nextTask.DeadlineTime,
		nextTask.TalkAboutInInterview,
		nextTask.ScheduleCollectionCondition,
		nextTask.ExamGuideContent,
		false,
	)

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 内定承諾/入社確認
	if nextTask.PhaseCategory == null.NewInt(int64(entity.AcceptJobOffer), true) && nextTask.PhaseSubCategory == null.NewInt(int64(entity.ConfirmJoiningCompany), true) {

		// 入社日の登録
		err = i.taskGroupRepository.UpdateJoiningDate(nextTask.TaskGroupID, nextTask.JoiningDate)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 依頼時間の更新
	if nextTask.PrevStaffType == null.NewInt(int64(entity.CA), true) {
		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			// CAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// CAからRAにタスク作成する場合
			err := i.taskGroupRepository.UpdateCALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	} else {
		// RAからRAにタスク作成する場合
		if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			err := i.taskGroupRepository.UpdateLastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// RAからCAにタスク作成する場合
			err := i.taskGroupRepository.UpdateRALastRequestAt(nextTask.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	if nextTask.PhaseSubCategory != null.NewInt(entity.AcceptJobOfferPhaseClose, true) {
		// プッシュ通知
		firebaseID := ""
		contents := ""
		topic := ""
		redirectURL := os.Getenv("BASE_DOMAIN")

		if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
			contents = "新着CAタスクの依頼があります。"
			topic = "TaskCA"

			caStaff, err := i.agentStaffRepository.FindByID(nextTask.CAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = caStaff.FirebaseID

		} else if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
			contents = "新着RAタスクの依頼があります。"
			topic = "TaskRA"

			raStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			firebaseID = raStaff.FirebaseID
		}

		err = utility.WebPush(
			i.oneSignal.AppID,
			i.oneSignal.APIKey,
			firebaseID,
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
	// else {
	// 	// タスクが終了した場合　ヨミを「失注」に更新
	// 	sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(nextTask.JobSeekerID, nextTask.JobInformationID)
	// 	if err != nil {
	// 		if errors.Is(err, entity.ErrNotFound) {
	// 			fmt.Println("ヨミの登録がありません")
	// 		} else {
	// 			fmt.Println(err)
	// 			return output, err
	// 		}
	// 	} else {
	// 		sale.Accuracy = null.NewInt(entity.AccuracyFailure, true) // ヨミのみ「失注: 5」に更新
	// 		err = i.saleRepository.Update(sale.ID, sale)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return output, err
	// 		}
	// 	}
	// }

	output.OK = true

	return output, nil
}

// 内定承諾フェーズのタスク処理
type CreateNextSameTaskListInput struct {
	NextSameTaskListParam entity.NextSameTaskListParam
	AgentStaffID          uint
}

type CreateNextSameTaskListOutput struct {
	OK bool
}

// 内定承諾フェーズのタスク処理
func (i *TaskInteractorImpl) CreateNextSameTaskList(input CreateNextSameTaskListInput) (CreateNextSameTaskListOutput, error) {
	var (
		output       CreateNextSameTaskListOutput
		err          error
		nextTask     = input.NextSameTaskListParam.Task
		sameTaskList = input.NextSameTaskListParam.SameTaskList
		// message   = input.NextSameTaskListParam.Message
		// IsSendMessageForJobSeeker bool
	)

	// リクエストのタスクが存在しない場合はリクエストエラー
	if len(sameTaskList) == 0 {
		wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "タスクが存在しません")
		return output, wrapped
	}

	for _, sameTask := range sameTaskList {
		// 選考フェーズでなかった場合はエラー
		// if nextTask.PrevPhaseCategory != null.NewInt(int64(entity.Entry), true) && nextTask.PrevPhaseCategory != null.NewInt(int64(entity.DocumentSelection), true) {
		// 	return output, errors.New("現在フェーズが「エントリー」でも「書類選考」でもありません")
		// }

		task := entity.NewTask(
			sameTask.TaskGroupID,
			nextTask.PhaseCategory,
			nextTask.PhaseSubCategory,
			nextTask.StaffType, // 担当者タイプ
			nextTask.ExecutedStaffID,
			nextTask.Remarks,
			nextTask.DeadlineDay,
			nextTask.DeadlineTime,
			nextTask.TalkAboutInInterview,
			nextTask.ScheduleCollectionCondition,
			nextTask.ExamGuideContent,
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(task)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		/************ おすすめ応募書類情報の登録 **************/

		// 次のタスクが「書類選考/書類推薦」の場合は「送付をお勧めする書類の情報」を登録
		if nextTask.PhaseCategory == null.NewInt(int64(entity.DocumentSelection), true) &&
			nextTask.PhaseSubCategory == null.NewInt(int64(entity.RequestRecommendations), true) {

			// 入力された送付をお勧めする書類の情報
			isR := nextTask.IsRecommendDocument

			// 作成する
			isRecommend := entity.NewTaskIsRecommendDocument(
				task.ID,
				isR.IsRecommendUploadResume,
				isR.IsRecommendUploadCV,
				isR.IsRecommendUploadRecommendation,
				isR.IsRecommendGeneratedResume,
				isR.IsRecommendGeneratedCV,
				isR.IsRecommendGeneratedRecommendation,
				isR.IsRecommendGeneratedMaskResume,
			)

			err := i.taskIsRecommendDocumentRepository.Create(isRecommend)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ 求人票のURLを求職者へ送信（LINE or メール） **************/

	var IsSendMessageForJobSeeker = nextTask.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		nextTask.PhaseCategory == null.NewInt(int64(entity.Entry), true) &&
		nextTask.PhaseSubCategory == null.NewInt(int64(entity.ConfirmApplicationIntention), true)

	// 求人票を送信するかどうか（LINE or メール）
	if IsSendMessageForJobSeeker {
		var message = input.NextSameTaskListParam.Task.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(nextTask.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, nextTask.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return output, err
			}
			fmt.Println("メッセージの送信に成功しました:", call)

			// エージェントから求職者へのチャットメッセージを作成
			chatMessage := entity.NewChatMessageWithJobSeeker(
				chatGroup.ID,
				null.NewInt(0, true), // エージェント
				message.Content,
				"",
				"",
				"",
				"",
				null.NewInt(0, false),
				"",
				null.NewInt(0, true),
				time.Now().In(time.UTC),
			)

			err = i.chatMessageWithJobSeekerRepository.Create(chatMessage)
			if err != nil {
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

		} else {
			// メール送信
			err = utility.SendMailToSingle(
				i.sendgrid.APIKey,
				message.Subject,
				message.Content,
				entity.EmailUser{
					Name:  staff.StaffName,
					Email: staff.Email,
				},
				entity.EmailUser{
					Name:  (jobSeeker.LastName + jobSeeker.FirstName),
					Email: jobSeeker.Email,
				},
				nil,
			)

			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	/************ タスクの最終依頼時間と最終閲覧時間の更新 **************/

	for index, sameTask := range sameTaskList {
		// 依頼時間の更新
		if nextTask.PrevStaffType == null.NewInt(int64(entity.CA), true) {
			if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
				// CAからCAにタスク作成する場合
				err := i.taskGroupRepository.UpdateLastRequestAt(sameTask.TaskGroupID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else {
				// CAからRAにタスク作成する場合
				err := i.taskGroupRepository.UpdateCALastRequestAt(sameTask.TaskGroupID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		} else {
			// RAからRAにタスク作成する場合
			if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
				err := i.taskGroupRepository.UpdateLastRequestAt(sameTask.TaskGroupID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			} else {
				// RAからCAにタスク作成する場合
				err := i.taskGroupRepository.UpdateRALastRequestAt(sameTask.TaskGroupID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}

		/************ タスク作成のプッシュ通知 **************/

		if nextTask.PhaseSubCategory != null.NewInt(entity.AcceptJobOfferPhaseClose, true) {
			// プッシュ通知
			firebaseID := ""
			contents := ""
			topic := ""
			redirectURL := os.Getenv("BASE_DOMAIN")

			if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
				contents = "新着CAタスクの依頼があります。"
				topic = "TaskCA"

				caStaff, err := i.agentStaffRepository.FindByID(nextTask.CAStaffID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				firebaseID = caStaff.FirebaseID

			} else if nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
				contents = "新着RAタスクの依頼があります。"
				topic = "TaskRA"

				raStaff, err := i.agentStaffRepository.FindByID(nextTask.RAStaffID)
				if err != nil {
					fmt.Println(err)
					return output, err
				}

				firebaseID = raStaff.FirebaseID
			}

			// CAタスクの場合でも「index」を使って通知は一回にする
			if (nextTask.StaffType == null.NewInt(int64(entity.CA), true) && index == 0) || nextTask.StaffType == null.NewInt(int64(entity.RA), true) {
				err = utility.WebPush(
					i.oneSignal.AppID,
					i.oneSignal.APIKey,
					firebaseID,
					"タスク通知",
					contents,
					topic,
					redirectURL,
				)
				if err != nil {
					fmt.Println("WebPushの通知")
					fmt.Println(err)
				} else {
					if nextTask.StaffType == null.NewInt(int64(entity.CA), true) {
						fmt.Println("CAに通知しました")
					} else {
						fmt.Println("RAに通知しました")
					}
				}
			}
		} else {
			// タスクが終了した場合　ヨミを「失注」に更新
			sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(sameTask.JobSeekerID, sameTask.JobInformationID)
			if err != nil {
				if errors.Is(err, entity.ErrNotFound) {
					fmt.Println("ヨミの登録がありません")
				} else {
					fmt.Println(err)
					return output, err
				}
			} else {
				sale.Accuracy = null.NewInt(entity.AccuracyFailure, true) // ヨミのみ「失注: 5」に更新
				err = i.saleRepository.Update(sale.ID, sale)
				if err != nil {
					fmt.Println(err)
					return output, err
				}
			}
		}
	}

	output.OK = true

	return output, nil
}

// マイページのマッチ求人からエントリー
type CreateEntryTaskFromMatchingJobInput struct {
	Param entity.CreateEntryTaskFromMatchingJobParam
}

type CreateEntryTaskFromMatchingJobOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) CreateEntryTaskFromMatchingJob(input CreateEntryTaskFromMatchingJobInput) (CreateEntryTaskFromMatchingJobOutput, error) {
	var (
		output           CreateEntryTaskFromMatchingJobOutput
		err              error
		defaultCAStaffID uint = 2 // 本番環境のIDをデフォルトCAとして設定
	)

	jobInformation, err := i.jobInformationRepository.FindByUUID(input.Param.JobInformationUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.Param.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 両面タスクの判定
	isDoubleSided := false
	if jobInformation.AgentStaffID == defaultCAStaffID {
		isDoubleSided = true
	}

	// タスクグループの作成
	taskGroup := entity.NewTaskGroup(
		jobSeeker.ID,
		jobInformation.ID,
		isDoubleSided,
		"",
		"",
		"",
		jobInformation.AgentID, // repoでagentIDとってたのでOK
		jobSeeker.AgentID,
	)

	err = i.taskGroupRepository.Create(taskGroup)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var (
		// エントリー/求人打診の依頼（CAタスク）でタスクを作成
		phaseCategory    null.Int = null.NewInt(int64(entity.Entry), true)
		phaseSubCategory null.Int = null.NewInt(int64(entity.SoundOutJobInformation), true)
		remarks          string   = "求職者がマイページからエントリーを希望しました。"

		// InterestedJobListingに保存する興味タイプ
		interestedType null.Int = null.NewInt(int64(entity.InterestedTypeEntry), true)
	)

	// 興味ありの場合
	if !input.Param.IsEntry {
		remarks = "求職者がマイページから興味あり求人に追加しました。"

		// 興味タイプを興味ありに変更
		interestedType = null.NewInt(int64(entity.InterestedTypeInterested), true)
	}

	// InterestedJobListingの作成
	interestedJobListing := entity.NewJobSeekerInterestedJobListing(
		jobSeeker.ID,
		jobInformation.ID,
		interestedType,
	)

	err = i.jobSeekerInterestedJobListingRepository.Create(interestedJobListing)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 今日の日付を期限　JST 2021-09-01形式
	deadlineDayTime := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	deadlineDay := deadlineDayTime.Format("2006-01-02")

	if jobSeeker.AgentStaffID.Valid {
		defaultCAStaffID = uint(jobSeeker.AgentStaffID.Int64)
	}

	// タスク作成
	task := entity.NewTask(
		taskGroup.ID,
		phaseCategory,
		phaseSubCategory,
		null.NewInt(int64(entity.CA), true), // StaffType,
		defaultCAStaffID,                    // ExecutedStaffID,
		remarks,                             // Remarks,
		deadlineDay,                         // DeadlineDay,
		null.NewInt(99, true),               // DeadlineTime, // 指定なし
		"",                                  // TalkAboutInInterview,
		"",                                  // ScheduleCollectionCondition,
		"",                                  // ExamGuideContent,
		false,
	)

	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 依頼時間の更新
	// CAからCAにタスク作成する場合
	err = i.taskGroupRepository.UpdateLastRequestAt(taskGroup.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type UpdateTaskGroupDocumentInput struct {
	Param entity.UpdateTaskGroupDocumentParam
}

type UpdateTaskGroupDocumentOutput struct {
	OK bool
}

// エントリーフェーズのタスク処理
func (i *TaskInteractorImpl) UpdateTaskGroupDocument(input UpdateTaskGroupDocumentInput) (UpdateTaskGroupDocumentOutput, error) {
	var (
		output    UpdateTaskGroupDocumentOutput
		err       error
		document1 string
		document2 string
		document3 string
		document4 string
		document5 string
	)

	if input.Param.DocumentType == 1 {
		document1 = input.Param.DocumentURL
	} else if input.Param.DocumentType == 2 {
		document2 = input.Param.DocumentURL
	} else if input.Param.DocumentType == 3 {
		document3 = input.Param.DocumentURL
	} else if input.Param.DocumentType == 4 {
		document4 = input.Param.DocumentURL
	} else {
		document5 = input.Param.DocumentURL
	}

	document := entity.NewTaskGroupDocuemnt(
		input.Param.TaskGroupID,
		document1,
		document2,
		document3,
		document4,
		document5,
	)

	taskDocument, err := i.taskGroupDocumentRepository.FindByGroupID(input.Param.TaskGroupID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// まだレコードが作成されていない場合
			err = i.taskGroupDocumentRepository.Create(document)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			fmt.Println(err)
			return output, err
		}
	} else {
		// すでにレコードが存在する場合
		if input.Param.DocumentType == 1 {
			err := i.taskGroupDocumentRepository.UpdateDocument1URL(taskDocument.ID, input.Param.DocumentURL)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else if input.Param.DocumentType == 2 {
			err := i.taskGroupDocumentRepository.UpdateDocument2URL(taskDocument.ID, input.Param.DocumentURL)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else if input.Param.DocumentType == 3 {
			err := i.taskGroupDocumentRepository.UpdateDocument3URL(taskDocument.ID, input.Param.DocumentURL)
			if err != nil {
				fmt.Println(err)
				return output, err
			}

		} else if input.Param.DocumentType == 4 {
			err := i.taskGroupDocumentRepository.UpdateDocument4URL(taskDocument.ID, input.Param.DocumentURL)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			err := i.taskGroupDocumentRepository.UpdateDocument5URL(taskDocument.ID, input.Param.DocumentURL)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

	}

	output.OK = true

	return output, nil
}

// タスクIDを使ってタスク情報を取得する
type GetLatestTaskListByAgentStaffIDInput struct {
	AgentID      uint
	AgentStaffID uint
	DeadLine     uint
	StaffType    uint
	PartnerType  uint
	TaskType     uint
	Phase        uint
	JobSeekerID  uint
}

type GetLatestTaskListByAgentStaffIDOutput struct {
	TaskList      []*entity.Task
	JobSeekerList []*entity.JobSeeker
}

func (i *TaskInteractorImpl) GetLatestTaskListByAgentStaffID(input GetLatestTaskListByAgentStaffIDInput) (GetLatestTaskListByAgentStaffIDOutput, error) {
	var (
		output                GetLatestTaskListByAgentStaffIDOutput
		err                   error
		taskList              []*entity.Task
		taskIDList            []uint
		taskListByJobSeeker   []*entity.Task
		taskListByDeadline    []*entity.Task
		taskListByStaffType   []*entity.Task
		taskListByPartnerType []*entity.Task
		taskListByTaskType    []*entity.Task
		taskListByPhase       []*entity.Task
	)

	taskList, err = i.taskRepository.GetLatestByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 絞り込みに使用する求職者のリストを取得
	var duplicateCheck = make(map[uint]bool)
	for _, task := range taskList {
		// タブが「自分」でタスク担当が自分、または自分担当でなくてもリャンメンの場合
		var IsMyTask = input.TaskType == 0 && ((task.IsDoubleSided && task.CAStaffID == input.AgentStaffID) ||
			(!task.IsDoubleSided && ((task.StaffType == null.NewInt(int64(entity.CA), true) && task.CAStaffID == input.AgentStaffID) ||
				(task.StaffType == null.NewInt(int64(entity.RA), true) && task.RAStaffID == input.AgentStaffID))))

		// タブが「依頼中」でタスク担当が自分以外
		var IsNotMyTask = (input.TaskType != 0 && !task.IsDoubleSided && ((task.StaffType == null.NewInt(int64(entity.CA), true) && task.CAStaffID != input.AgentStaffID) ||
			(task.StaffType == null.NewInt(int64(entity.RA), true) && task.RAStaffID != input.AgentStaffID)))

		if !duplicateCheck[task.JobSeekerID] && (IsMyTask || IsNotMyTask) {
			duplicateCheck[task.JobSeekerID] = true
			js := &entity.JobSeeker{ID: task.JobSeekerID, LastName: task.LastName, FirstName: task.FirstName}
			output.JobSeekerList = append(output.JobSeekerList, js)
		}
	}

	// 求職者の絞り込み
	if input.JobSeekerID == 0 {
		taskListByJobSeeker = taskList
	} else {
		for _, task := range taskList {
			// 「0: すべて」ではない場合は絞り込み
			if task.JobSeekerID == input.JobSeekerID {
				taskListByJobSeeker = append(taskListByJobSeeker, task)
			}
		}
	}

	// 期限
	for _, task := range taskListByJobSeeker {
		layout := "2006-01-02"
		deadLineDay, err := time.Parse(layout, task.DeadlineDay)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		today := time.Now().UTC().Truncate(24 * time.Hour) // 今日の0時0分0秒
		if input.DeadLine == 0 {
			taskListByDeadline = taskListByJobSeeker
		} else if input.DeadLine == 1 {
			if deadLineDay.Equal(today) || deadLineDay.Before(today) {
				taskListByDeadline = append(taskListByDeadline, task)
			}
		} else {
			// 24時間未満のものは含まないように、24時間で丸める
			if deadLineDay.After(today) {
				// 明日以降の場合の処理
				taskListByDeadline = append(taskListByDeadline, task)
			}
		}
	}

	// 業務(RA or CA)
	for _, task := range taskListByDeadline {
		if input.StaffType == 0 {
			taskListByStaffType = taskListByDeadline
		} else if input.StaffType == 1 && task.StaffType == null.NewInt(int64(entity.CA), true) {
			taskListByStaffType = append(taskListByStaffType, task)
		} else if input.StaffType == 2 && task.StaffType == null.NewInt(int64(entity.RA), true) {
			taskListByStaffType = append(taskListByStaffType, task)
		}
	}

	// 範囲(自社 or アライアンス)
	for _, task := range taskListByStaffType {
		if input.PartnerType == 0 {
			// すべて
			taskListByPartnerType = taskListByStaffType
		} else if input.PartnerType == 1 {
			//　社内
			if input.AgentID == task.RAAgentID && input.AgentID == task.CAAgentID {
				taskListByPartnerType = append(taskListByPartnerType, task)
			}
		} else {
			// アライアンス
			if input.AgentID != task.RAAgentID || input.AgentID != task.CAAgentID {
				taskListByPartnerType = append(taskListByPartnerType, task)
			}
		}
	}

	// 業務(自分 or 依頼中)
	for _, task := range taskListByPartnerType {
		if !task.IsDoubleSided {
			// 片面タスク（通常）
			if input.TaskType == 0 {
				if (task.StaffType == null.NewInt(int64(entity.CA), true) && task.CAStaffID == input.AgentStaffID) ||
					(task.StaffType == null.NewInt(int64(entity.RA), true) && task.RAStaffID == input.AgentStaffID) {
					taskListByTaskType = append(taskListByTaskType, task)
				}
			} else {
				if (task.StaffType == null.NewInt(int64(entity.CA), true) && task.CAStaffID != input.AgentStaffID) ||
					(task.StaffType == null.NewInt(int64(entity.RA), true) && task.RAStaffID != input.AgentStaffID) {
					taskListByTaskType = append(taskListByTaskType, task)
				}
			}
		} else {
			// 両面タスク（CAが全てタスク動かす）
			if input.TaskType == 0 {
				if task.CAStaffID == input.AgentStaffID {
					taskListByTaskType = append(taskListByTaskType, task)
				}
			}
			// else {
			// 	if taskByPartnerType.CAStaffID != input.AgentStaffID {
			// 		taskListByTaskType = append(taskListByTaskType, taskByPartnerType)
			// 	}
			// }
		}
	}

	// 選考フェーズ
	for _, task := range taskListByTaskType {
		if !task.PhaseCategory.Valid {
			continue
		}
		if input.Phase == 0 {
			taskListByPhase = taskListByTaskType
		} else if input.Phase == 1 {
			if task.PhaseCategory == null.NewInt(int64(0), true) {
				taskListByPhase = append(taskListByPhase, task)
			}
		} else if input.Phase == 2 {
			if task.PhaseCategory.Int64 >= 1 && task.PhaseCategory.Int64 <= 7 {
				taskListByPhase = append(taskListByPhase, task)
			}

		} else if input.Phase == 3 {
			if task.PhaseCategory == null.NewInt(int64(8), true) {
				taskListByPhase = append(taskListByPhase, task)
			}
		} else {
			if task.PhaseCategory == null.NewInt(int64(9), true) {
				taskListByPhase = append(taskListByPhase, task)
			}
		}
	}

	// タスクIDのリストを作成
	for _, task := range taskListByPhase {
		taskIDList = append(taskIDList, task.ID)
	}

	// 送付をお勧めする書類の情報を取得
	isRecommendList, err := i.taskIsRecommendDocumentRepository.GetByTaskIDList(taskIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送付をお勧めする書類の情報をセット
	for _, task := range taskListByPhase {
		for _, isRecommend := range isRecommendList {
			if task.ID == isRecommend.TaskID {
				task.IsRecommendDocument = *isRecommend
			}
		}
	}

	/************ 選考日時の情報を取得 **************/

	// タスクグループIDのリスト
	var groupIDList []uint

	for _, task := range taskListByPhase {
		groupIDList = append(groupIDList, task.TaskGroupID)
	}

	// 求職者のスケジュール情報を取得
	scheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupIDList(groupIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, task := range taskListByPhase {
		// 確定日時がセットされているか

		for _, schedule := range scheduleList {
			// 「同一タスクグループ」「同一選考フェーズ」
			if task.TaskGroupID == uint(schedule.TaskGroupID.Int64) && schedule.PhaseCategory == task.PhaseCategory {
				// 以下の条件に合う日程を確定日時としてセット
				// 「確定日時がまだセットされていない」&&「タスクグループIDが一致している」&&「日程タイプが確定日時」&&「選考フェーズが一致している」
				if schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByEnterprise, true) {
					task.SelectionDate = *schedule
				}

				if schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByCA, true) {
					// 「確定日時がない」または「確定日時があるけど、繰り返しの回数が違う」
					// if !task.SelectionDate.TaskID.Valid || (task.SelectionDate.TaskID.Valid && task.SelectionDate.RepetitionCount != schedule.RepetitionCount) {
					// 	// 候補日時をセット
					// 	fmt.Println("------------")
					// 	fmt.Println(task.LastName + task.FirstName, task.CompanyName, task.Title)
					// 	task.PossibleDates = append(task.PossibleDates, *schedule)
					// }
					task.PossibleDates = append(task.PossibleDates, *schedule)
				}
			}
		}
	}

	output.TaskList = taskListByPhase

	return output, nil
}

// 「同一求職者」で「同一フェーズ」のタスクのリストを取得
type GetLatestSameTaskListByJobSeekerIDInput struct {
	JobSeekerID uint
	TaskID      uint
	Phase       null.Int
	PhaseSub    null.Int
}

type GetLatestSameTaskListByJobSeekerIDOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) GetLatestSameTaskListByJobSeekerID(input GetLatestSameTaskListByJobSeekerIDInput) (GetLatestSameTaskListByJobSeekerIDOutput, error) {
	var (
		output   GetLatestSameTaskListByJobSeekerIDOutput
		err      error
		taskList []*entity.Task
	)

	taskList, err = i.taskRepository.GetLatestSameByJobSeekerID(
		input.JobSeekerID,
		input.TaskID,
		input.Phase,
		input.PhaseSub,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.TaskList = taskList

	return output, nil
}

// タスクIDを使ってタスク情報を取得する
type GetLatestTaskByJobSeekerIDAndJobInformationIDInput struct {
	JobSeekerID      uint
	JobInformationID uint
}

type GetLatestTaskByJobSeekerIDAndJobInformationIDOutput struct {
	Task *entity.Task
}

func (i *TaskInteractorImpl) GetLatestTaskByJobSeekerIDAndJobInformationID(input GetLatestTaskByJobSeekerIDAndJobInformationIDInput) (GetLatestTaskByJobSeekerIDAndJobInformationIDOutput, error) {
	var (
		output GetLatestTaskByJobSeekerIDAndJobInformationIDOutput
		err    error
	)

	task, err := i.taskRepository.FindLatestByJobSeekerIDAndJobInformationID(input.JobSeekerID, input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Task = task

	return output, nil
}

type UpdateRALastWatchedInput struct {
	GroupID uint
}

type UpdateRALastWatchedOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) UpdateRALastWatched(input UpdateRALastWatchedInput) (UpdateRALastWatchedOutput, error) {
	var (
		output UpdateRALastWatchedOutput
	)

	err := i.taskGroupRepository.UpdateRALastWatchedAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type UpdateRALastRequestInput struct {
	GroupID uint
}

type UpdateRALastRequestOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) UpdateRALastRequest(input UpdateRALastRequestInput) (UpdateRALastRequestOutput, error) {
	var (
		output UpdateRALastRequestOutput
	)

	err := i.taskGroupRepository.UpdateRALastRequestAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type UpdateCALastWatchedInput struct {
	GroupID uint
}

type UpdateCALastWatchedOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) UpdateCALastWatched(input UpdateCALastWatchedInput) (UpdateCALastWatchedOutput, error) {
	var (
		output UpdateCALastWatchedOutput
	)

	err := i.taskGroupRepository.UpdateCALastWatchedAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type UpdateCALastRequestInput struct {
	GroupID uint
}

type UpdateCALastRequestOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) UpdateCALastRequest(input UpdateCALastRequestInput) (UpdateCALastRequestOutput, error) {
	var (
		output UpdateCALastRequestOutput
	)

	err := i.taskGroupRepository.UpdateCALastRequestAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 自分から自分にタスクを作成する場合に使用
// RAとCAの最終依頼時間を更新
type UpdateLastRequestInput struct {
	GroupID uint
}

type UpdateLastRequestOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) UpdateLastRequest(input UpdateLastRequestInput) (UpdateLastRequestOutput, error) {
	var (
		output UpdateLastRequestOutput
	)

	err := i.taskGroupRepository.UpdateLastRequestAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/
/// 汎用系 API（仕様変更前に作成した関数）
//

// タスクIDを使ってタスク情報を取得する
type GetTaskByIDInput struct {
	TaskID uint
}

type GetTaskByIDOutput struct {
	Task *entity.Task
}

func (i *TaskInteractorImpl) GetTaskByID(input GetTaskByIDInput) (GetTaskByIDOutput, error) {
	var (
		output GetTaskByIDOutput
		err    error
	)

	task, err := i.taskRepository.FindByID(input.TaskID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Task = task

	return output, nil
}

// タスクグループの一覧取得（エージェントが関わっているタスク）
type GetTaskListAfterEntryByJobSeekerIDInput struct {
	JobSeekerID uint
}

type GetTaskListAfterEntryByJobSeekerIDOutput struct {
	TaskList []*entity.Task
}

func (i *TaskInteractorImpl) GetTaskListAfterEntryByJobSeekerID(input GetTaskListAfterEntryByJobSeekerIDInput) (GetTaskListAfterEntryByJobSeekerIDOutput, error) {
	var (
		output GetTaskListAfterEntryByJobSeekerIDOutput
		err    error
	)

	/**********************************************/
	// 書類選考以降のタスク情報取得
	//
	taskList, err := i.taskRepository.GetLatestAfterSelectPhaseByJobSeekerID(input.JobSeekerID, entity.TaskCategory(entity.Entry))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/**********************************************/
	// 選考後アンケート情報を取得
	//

	questionnaireList, err := i.selectionQuestionnaireRepository.GetByJobSeekerID(input.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/**********************************************/
	// 選考後アンケートのIDから他社選考の志望度情報を取得
	//
	var (
		idListUint []uint
	)

	// 選考後アンケート情報のIDを抜き出す
	for _, questionnaire := range questionnaireList {
		idListUint = append(idListUint, questionnaire.ID)
	}

	rankingList, err := i.selectionQuestionnaireMyRankingRepository.GetByQuestionnaireIDList(idListUint)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/**********************************************/
	// 取得した情報Outputにセット
	//

	for _, task := range taskList {
		for _, questionnaire := range questionnaireList {
			if task.JobInformationID == questionnaire.JobInformationID {
				for _, ranking := range rankingList {
					if questionnaire.ID == ranking.SelectionQuestionnaireID {
						questionnaire.SelectionQuestionnaireMyRanking = append(questionnaire.SelectionQuestionnaireMyRanking, *ranking)
					}
				}
				task.SelectionQuestionnaireList = append(task.SelectionQuestionnaireList, *questionnaire)
			}
		}
	}

	output.TaskList = taskList

	return output, nil
}

// タスクグループの一覧取得（エージェントが関わっているタスク）
type GetTaskListByAgentIDAndPageInput struct {
	AgentID    uint
	PageNumber uint
}

type GetTaskListByAgentIDAndPageOutput struct {
	TaskList      []*entity.Task
	MaxPageNumber uint
}

func (i *TaskInteractorImpl) GetTaskListByAgentIDAndPage(input GetTaskListByAgentIDAndPageInput) (GetTaskListByAgentIDAndPageOutput, error) {
	var (
		output                     GetTaskListByAgentIDAndPageOutput
		err                        error
		taskIDList                 []uint
		jobInformationIDList       []uint
		jobSeekerIDList            []uint
		selectionFlowPatternIDList []uint
		groupIDList                []uint
		scheduleTypeByCA           = null.NewInt(entity.ScheduleTypeByCA, true)
		scheduleTypeByEnterprise   = null.NewInt(entity.ScheduleTypeByEnterprise, true)
	)

	taskGroupList, err := i.taskGroupRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ページの最大数を取得
	output.MaxPageNumber = getTaskGroupListMaxPage(taskGroupList)

	// 指定ページのタスクグループ20件を取得（本番実装までは1ページあたり5件）
	taskGroupList20 := getTaskGroupListWithPage(taskGroupList, input.PageNumber)

	// タスクIDのリストを作成
	for _, taskGroup := range taskGroupList20 {
		groupIDList = append(groupIDList, taskGroup.ID)
		jobInformationIDList = append(jobInformationIDList, taskGroup.JobInformationID)
		jobSeekerIDList = append(jobSeekerIDList, taskGroup.JobSeekerID)
		if taskGroup.SelectionFlowPatternID.Valid && taskGroup.SelectionFlowPatternID.Int64 != 0 {
			selectionFlowPatternIDList = append(selectionFlowPatternIDList, uint(taskGroup.SelectionFlowPatternID.Int64))
		}
	}

	taskList, err := i.taskRepository.GetLatestByGroupIDList(groupIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, task := range taskList {
		taskIDList = append(taskIDList, task.ID)
	}

	// 送付をお勧めする書類の情報を取得
	isRecommendList, err := i.taskIsRecommendDocumentRepository.GetByTaskIDList(taskIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobInformationList, err := i.jobInformationRepository.GetByIDListForTask(jobInformationIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobSeekeerList, err := i.jobSeekerRepository.GetByIDListForTask(jobSeekerIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatternList, err := i.jobInfoSelectionFlowPatternRepository.GetByIDList(selectionFlowPatternIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者のスケジュール情報を取得
	scheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupIDList(groupIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, taskGroup := range taskGroupList20 {
		// Taskに変換
		task := &entity.Task{}
		for _, t := range taskList {
			if taskGroup.ID == t.TaskGroupID {
				task = t
				break
			}
		}
		if task.ID == 0 {
			continue
		}
		task.TaskGroupUUID = taskGroup.UUID
		task.RALastRequestAt = taskGroup.RALastRequestAt
		task.RALastWatchedAt = taskGroup.RALastWatchedAt
		task.CALastRequestAt = taskGroup.CALastRequestAt
		task.CALastWatchedAt = taskGroup.CALastWatchedAt
		task.IsDoubleSided = taskGroup.IsDoubleSided
		task.IsSelfApplication = taskGroup.IsSelfApplication
		task.ExternalCompanyName = taskGroup.ExternalCompanyName
		task.ExternalJobInformationTitle = taskGroup.ExternalJobInformationTitle
		task.ExternalJobListingURL = taskGroup.ExternalJobListingURL
		task.JobInformationID = taskGroup.JobInformationID
		task.JobSeekerID = taskGroup.JobSeekerID
		task.SelectionFlowPatternID = taskGroup.SelectionFlowPatternID

		// 求人情報をセット
		for _, jobInformation := range jobInformationList {
			if task.JobInformationID == jobInformation.ID {
				task.JobInformationUUID = jobInformation.UUID
				task.BillingAddressID = jobInformation.BillingAddressID
				task.HowToRecommend = jobInformation.HowToRecommend
				task.RequiredDocumentsDetail = jobInformation.RequiredDocumentsDetail
				task.RAStaffID = jobInformation.AgentStaffID
				task.RAStaffName = jobInformation.StaffName
				task.RAAgentID = jobInformation.AgentID
				task.RAAgentName = jobInformation.AgentName
				if jobInformation.IsExternal && task.ExternalJobInformationTitle != "" {
					task.Title = task.ExternalJobInformationTitle
				} else {
					task.Title = jobInformation.Title
				}
				if jobInformation.IsExternal && task.CompanyName != "" {
					task.CompanyName = task.ExternalCompanyName
				} else {
					task.CompanyName = jobInformation.CompanyName
				}
			}
		}

		// 求職者情報をセット
		for _, jobSeeker := range jobSeekeerList {
			if task.JobSeekerID == jobSeeker.ID {
				task.JobSeekerUUID = jobSeeker.UUID
				task.LastName = jobSeeker.LastName
				task.FirstName = jobSeeker.FirstName
				task.LastFurigana = jobSeeker.LastFurigana
				task.FirstFurigana = jobSeeker.FirstFurigana
				task.CAStaffID = uint(jobSeeker.AgentStaffID.Int64)
				task.CAStaffName = jobSeeker.StaffName
				task.CAAgentID = jobSeeker.AgentID
				task.CAAgentName = jobSeeker.AgentName
			}
		}

		for _, selectionFlowPattern := range selectionFlowPatternList {
			if task.SelectionFlowPatternID.Valid &&
				uint(task.SelectionFlowPatternID.Int64) == selectionFlowPattern.ID &&
				task.PhaseCategory == selectionFlowPattern.SelectionType {
				task.SelectionFlowPatternID = null.NewInt(int64(selectionFlowPattern.ID), true)
				task.FlowPattern = selectionFlowPattern.FlowPattern
				task.SelectionInformationID = selectionFlowPattern.SelectionInformationID
				task.IsQuestionnairy = selectionFlowPattern.IsQuestionnairy
			}
		}

		var isSetSelectionDate = false

		for _, schedule := range scheduleList {
			if task.ID == uint(schedule.TaskID.Int64) && schedule.ScheduleType == scheduleTypeByCA {
				// 候補日時をセット
				task.PossibleDates = append(task.PossibleDates, *schedule)
			}
			if !isSetSelectionDate && task.ID == uint(schedule.TaskID.Int64) && schedule.ScheduleType == scheduleTypeByEnterprise {
				// 確定日時をセット
				task.SelectionDate = *schedule
				isSetSelectionDate = true
			}
		}

		for _, isRecommend := range isRecommendList {
			if task.ID == isRecommend.TaskID {
				task.IsRecommendDocument = *isRecommend
			}
		}

		output.TaskList = append(output.TaskList, task)
	}

	return output, nil
}

type GetSearchTaskListByAgentIDAndPageInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchTask
}

type GetSearchTaskListByAgentIDAndPageOutput struct {
	TaskList      []*entity.Task
	MaxPageNumber uint
}

func (i *TaskInteractorImpl) GetSearchTaskListByAgentIDAndPage(input GetSearchTaskListByAgentIDAndPageInput) (GetSearchTaskListByAgentIDAndPageOutput, error) {
	var (
		output                     GetSearchTaskListByAgentIDAndPageOutput
		err                        error
		taskIDList                 []uint
		groupIDList                []uint
		selectionFlowPatternIDList []uint
	)

	/**
	GetLatestByAgentIDAndFreeWordは
	フリーワードの有無で処理を分岐
	*/

	taskList, err := i.taskRepository.GetLatestByAgentIDAndFreeWord(input.AgentID, input.SearchParam.JobSeekerFreeWord, input.SearchParam.EnterpriseFreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var (
		taskListWithRAStaffID      []*entity.Task
		taskListWithCAStaffID      []*entity.Task
		taskListWithTaskPhaseTypes []*entity.Task
	)

	raStaffID, err := strconv.Atoi(input.SearchParam.RAStaffID)
	if !(err != nil || raStaffID == 0) {
		for _, taskGroup := range taskList {
			if taskGroup.RAStaffID != uint(raStaffID) {
				continue
			}
			taskListWithRAStaffID = append(taskListWithRAStaffID, taskGroup)
		}
	}

	if err != nil || raStaffID == 0 {
		taskListWithRAStaffID = taskList
	}

	caStaffID, err := strconv.Atoi(input.SearchParam.CAStaffID)
	if !(err != nil || caStaffID == 0) {
		for _, task := range taskListWithRAStaffID {
			if task.CAStaffID != uint(caStaffID) {
				continue
			}
			taskListWithCAStaffID = append(taskListWithCAStaffID, task)
		}
	}

	if err != nil || caStaffID == 0 {
		taskListWithCAStaffID = taskListWithRAStaffID
	}

	if !(len(input.SearchParam.TaskPhaseTypes) == 0) {
		if !(len(taskListWithCAStaffID) == 0) {
			for _, task := range taskListWithCAStaffID {
				for _, taskPhase := range input.SearchParam.TaskPhaseTypes {
					if !taskPhase.Valid {
						continue
					}
					if taskPhase != task.PhaseCategory {
						continue
					}
					taskListWithTaskPhaseTypes = append(taskListWithTaskPhaseTypes, task)
				}
			}
		}
	}

	// TaskPhaseTypeがない場合
	if len(input.SearchParam.TaskPhaseTypes) == 0 {
		taskListWithTaskPhaseTypes = taskListWithCAStaffID
	}

	// ページの最大数を取得
	output.MaxPageNumber = getTaskListMaxPage(taskListWithTaskPhaseTypes)

	// 指定ページのタスクグループ20件を取得（本番実装までは1ページあたり5件）
	taskList20 := getTaskListWithPage(taskListWithTaskPhaseTypes, input.PageNumber)

	// タスクIDのリストを作成
	for _, task := range taskList20 {
		taskIDList = append(taskIDList, task.ID)
		groupIDList = append(groupIDList, task.TaskGroupID)
		if task.SelectionFlowPatternID.Valid && task.SelectionFlowPatternID.Int64 != 0 {
			selectionFlowPatternIDList = append(selectionFlowPatternIDList, uint(task.SelectionFlowPatternID.Int64))
		}
	}

	// 求職者のスケジュール情報を取得
	scheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupIDList(groupIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	selectionFlowPatternList, err := i.jobInfoSelectionFlowPatternRepository.GetByIDList(selectionFlowPatternIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 選考の日程情報を取得
	for _, task := range taskList20 {

		for _, selectionFlowPattern := range selectionFlowPatternList {
			if task.SelectionFlowPatternID.Valid &&
				uint(task.SelectionFlowPatternID.Int64) == selectionFlowPattern.ID &&
				task.PhaseCategory == selectionFlowPattern.SelectionType {
				task.SelectionFlowPatternID = null.NewInt(int64(selectionFlowPattern.ID), true)
				task.FlowPattern = selectionFlowPattern.FlowPattern
				task.SelectionInformationID = selectionFlowPattern.SelectionInformationID
				task.IsQuestionnairy = selectionFlowPattern.IsQuestionnairy
			}
		}

		var isSetSelectionDate = false

		for _, schedule := range scheduleList {
			if task.ID == uint(schedule.TaskID.Int64) && schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByCA, true) {
				// 候補日時をセット
				task.PossibleDates = append(task.PossibleDates, *schedule)
			}
			if !isSetSelectionDate && task.ID == uint(schedule.TaskID.Int64) && schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByEnterprise, true) {
				// 確定日時をセット
				task.SelectionDate = *schedule
				isSetSelectionDate = true
			}
		}
	}

	// 送付をお勧めする書類の情報を取得
	isRecommendList, err := i.taskIsRecommendDocumentRepository.GetByTaskIDList(taskIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送付をお勧めする書類の情報をセット
	for _, task := range taskList20 {
		for _, isRecommend := range isRecommendList {
			if task.ID == isRecommend.TaskID {
				task.IsRecommendDocument = *isRecommend
			}
		}
	}

	output.TaskList = taskList20

	return output, nil
}

// タスクグループの一覧取得（自分が関わっているタスク）
type GetTaskGroupByIDInput struct {
	TaskGroupID uint
}

type GetTaskGroupByIDOutput struct {
	TaskGroup *entity.TaskGroup
}

func (i *TaskInteractorImpl) GetTaskGroupByID(input GetTaskGroupByIDInput) (GetTaskGroupByIDOutput, error) {
	var (
		output GetTaskGroupByIDOutput
		err    error
	)

	taskGroup, err := i.taskGroupRepository.FindByID(input.TaskGroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	taskList, err := i.taskRepository.GetWithRelatedByTaskGroupID(taskGroup.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者のスケジュール情報を取得
	scheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupID(input.TaskGroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	taskDocument, err := i.taskGroupDocumentRepository.FindByGroupID(taskGroup.ID)
	// まだレコードが作成されていない場合のみエラーを返す
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		fmt.Println(err)
		return output, err
	}

	if taskDocument != nil {
		taskGroup.Document = *taskDocument
	}

	for _, task := range taskList {
		for _, schedule := range scheduleList {
			if task.ID == uint(schedule.TaskID.Int64) && schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByCA, true) {
				// 候補日時をセット
				task.PossibleDates = append(task.PossibleDates, *schedule)
			}
			if task.ID == uint(schedule.TaskID.Int64) && schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByEnterprise, true) {
				// 確定日時をセット
				task.SelectionDate = *schedule
			}
		}

		taskGroup.Tasks = append(taskGroup.Tasks, *task)
	}

	output.TaskGroup = taskGroup

	return output, nil
}

type GetActiveTaskCountByBillingAddressIDInput struct {
	BillingAddressID uint
}

type GetActiveTaskCountByBillingAddressIDOutput struct {
	TaskCount uint
}

func (i *TaskInteractorImpl) GetActiveTaskCountByBillingAddressID(input GetActiveTaskCountByBillingAddressIDInput) (GetActiveTaskCountByBillingAddressIDOutput, error) {
	var (
		output GetActiveTaskCountByBillingAddressIDOutput
		err    error
	)

	taskList, err := i.taskRepository.GetActiveByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.TaskCount = uint(len(taskList))

	return output, nil
}

type GetActiveTaskCountByJobInformationIDInput struct {
	JobInformationID uint
}

type GetActiveTaskCountByJobInformationIDOutput struct {
	TaskCount uint
}

func (i *TaskInteractorImpl) GetActiveTaskCountByJobInformationID(input GetActiveTaskCountByJobInformationIDInput) (GetActiveTaskCountByJobInformationIDOutput, error) {
	var (
		output GetActiveTaskCountByJobInformationIDOutput
		err    error
	)

	taskList, err := i.taskRepository.GetActiveByJobInformationID(input.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.TaskCount = uint(len(taskList))

	return output, nil
}

// タスクグループの一覧取得（自分が関わっているタスク）
type GetActiveTaskCountBySelectionIDInput struct {
	SelectionID uint
}

type GetActiveTaskCountBySelectionIDOutput struct {
	TaskCount uint
}

func (i *TaskInteractorImpl) GetActiveTaskCountBySelectionID(input GetActiveTaskCountBySelectionIDInput) (GetActiveTaskCountBySelectionIDOutput, error) {
	var (
		output GetActiveTaskCountBySelectionIDOutput
		err    error
	)

	taskList, err := i.taskRepository.GetActiveBySelectionFlowPatternID(input.SelectionID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.TaskCount = uint(len(taskList))

	return output, nil
}

// Taskの削除
type DeleteTaskInput struct {
	DeleteParam entity.DeleteTaskParam
}

type DeleteTaskOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) DeleteTask(input DeleteTaskInput) (DeleteTaskOutput, error) {
	var (
		output DeleteTaskOutput
	)

	// 削除するタスクの情報を取得
	task, err := i.taskRepository.FindByID(input.DeleteParam.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 削除するタスクで両面タスクになっている場合
	if task.IsCheckDoubleSided {
		err = i.taskGroupRepository.UpdateIsDoubleSided(task.TaskGroupID, false)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.taskRepository.Delete(input.DeleteParam.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	prevTask, err := i.taskRepository.FindLatestByGroupID(input.DeleteParam.TaskGroupID)
	if err != nil {
		if !errors.Is(err, entity.ErrNotFound) {
			fmt.Println(err)
			return output, err
		}
	} else {
		// 削除後の最新タスクのフェーズが「999: 選考スキップ or 1000: リスケ」の場合はそれも削除
		if prevTask.PhaseSubCategory == null.NewInt(entity.SelectionSkip, true) ||
			prevTask.PhaseSubCategory == null.NewInt(entity.RescheduleDeleteOfSelection, true) ||
			prevTask.PhaseSubCategory == null.NewInt(entity.RescheduleDeleteOfHoldJobOffer, true) {
			err = i.taskRepository.Delete(prevTask.ID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// 削除後の最新タスクのフェーズが「999: 選考スキップ」の場合はそれも削除
	taskList, err := i.taskRepository.GetByTaskGroupID(input.DeleteParam.TaskGroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if len(taskList) == 0 {
		err := i.taskGroupRepository.Delete(input.DeleteParam.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	} else {
		// 削除後の最新タスクが「書類選考/結果回収中」の場合はtask_groupsテーブルの「selection_flow_pattern_id」をnullに更新
		if taskList[0].PhaseCategory == null.NewInt(int64(entity.DocumentSelection), true) &&
			taskList[0].PhaseSubCategory == null.NewInt(int64(entity.CollectResultOfDocumentSelection), true) {
			// 選考フローパターンをnullに更新
			err = i.taskGroupRepository.UpdateSelectionFlowPatternID(input.DeleteParam.TaskGroupID, null.NewInt(0, false))
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
/****************************************************************************************/
/// Admin API
//

/****************************************************************************************/
/****************************************************************************************/
// Batch API
//
/*
	未読の通知(slack:https://app.slack.com/client/T01EY1159HP/C04MGLZSH8Q/thread/C04MGLZSH8Q-1687426686.644339)
	未読チャットの通知 対象:エージェント向けチャット。@付きのメッセージ
	未処理のタスク 対象:未読or期限切れ
	9:00と18:00に通知
*/
type BatchNotifyUnwatchedInput struct {
	Now time.Time
}

type BatchNotifyUnwatchedOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) BatchNotifyUnwatched(input BatchNotifyUnwatchedInput) (BatchNotifyUnwatchedOutput, error) {
	var (
		output BatchNotifyUnwatchedOutput
		err    error
	)

	// 全ての担当者とエージェントチャット情報、タスクを取得する
	agentStaffList, err := i.agentStaffRepository.All()
	if err != nil {
		log.Println(err)
		return output, err
	}

	chatGroupList, err := i.chatGroupWithAgentRepository.All()
	if err != nil {
		log.Println(err)
		return output, err
	}

	chatThreadList, err := i.chatThreadWithAgentRepository.All()
	if err != nil {
		log.Println(err)
		return output, err
	}

	chatMessageList, err := i.chatMessageWithAgentRepository.All()
	if err != nil {
		log.Println(err)
		return output, err
	}

	chatMessageToUserList, err := i.chatMessageToUserWithAgentRepository.All()
	if err != nil {
		log.Println(err)
		return output, err
	}
	log.Println("chatMessageToUserList", len(chatMessageToUserList))

	chatGroupWithJobSeekerList, err := i.chatGroupWithJobSeekerRepository.All()
	if err != nil {
		log.Println(err)
		return output, err
	}

	taskList, err := i.taskRepository.GetLatest()
	if err != nil {
		log.Println(err)
		return output, err
	}

	// チャットをまとめる
	for _, chatGroup := range chatGroupList {
		for _, chatThread := range chatThreadList {
			if chatGroup.ID == chatThread.GroupID {
				for _, chatMessage := range chatMessageList {
					if chatThread.ID == chatMessage.ThreadID {
						for _, chatMessageToUser := range chatMessageToUserList {
							if chatMessage.ID == chatMessageToUser.MessageID {
								chatMessage.MessageToUsers = append(chatMessage.MessageToUsers, *chatMessageToUser)
							}
						}
						chatThread.Messages = append(chatThread.Messages, *chatMessage)
					}
				}
				chatGroup.Threads = append(chatGroup.Threads, *chatThread)
			}
		}
	}

	// 他社エージェントとのチャットの未読数の構造体
	type otherAgentChat struct {
		AgentID   uint
		AgentName string
		ChatCount uint
	}

	// メールに記載する未読のチャットと未処理のタスク詳細の構造体
	type unwatchedDetail struct {
		// エージェントチャット
		AllAgentChatCount    uint             // チャットの未読
		MyAgentChatCount     uint             // 社内チャットの未読
		OtherAgentChatCounts []otherAgentChat // 他社エージェントとのチャットの未読

		// 求職者LINE
		AllJobSeekerChatCount uint   // チャットの未読
		JobSeekerNameList     string // 求職者名

		// タスク
		TaskCount uint          // 未読・期限切れのタスク
		Tasks     []entity.Task // 未読・期限切れのタスクの詳細
	}

	// 担当者ごとに判定
	for _, agentStaff := range agentStaffList {
		// 通知がオフ or 利用停止状態 はスキップ
		if !agentStaff.NotificationUnwatched ||
			agentStaff.UsageStatus == null.NewInt(int64(entity.UsageStatusNotAvailable), true) {
			continue
		}

		log.Println("agentStaff.ID", agentStaff.ID, "agentStaff.StaffName", agentStaff.StaffName)
		unwatchedDetail := unwatchedDetail{}

		// エージェントチャット
		for _, chatGroup := range chatGroupList {
			if agentStaff.AgentID == chatGroup.Agent1ID || agentStaff.AgentID == chatGroup.Agent2ID {
				for _, chatThread := range chatGroup.Threads {
				chatMessageLoop:
					for _, chatMessage := range chatThread.Messages {
						for _, chatMessageToUser := range chatMessage.MessageToUsers {

							/*
								未読判定
							*/
							// 自分宛のメッセージかつ未読のメッセージ
							if chatMessageToUser.AgentStaffID == agentStaff.ID &&
								chatMessageToUser.WatchedAt.Before(chatMessageToUser.SendAt) {
								log.Println(
									"未読のチャットメッセージ", chatMessageToUser.ID,
									"chatMessageToUser.WatchedAt", chatMessageToUser.WatchedAt, "chatMessageToUser.SendAt", chatMessageToUser.SendAt,
								)
								unwatchedDetail.AllAgentChatCount++

								// 社内チャットの場合
								if agentStaff.AgentID == chatGroup.Agent1ID && agentStaff.AgentID == chatGroup.Agent2ID {
									unwatchedDetail.MyAgentChatCount++

									// 他社エージェントとのチャットの場合　chatGroup.Agent1IDと一致する場合
								} else if agentStaff.AgentID == chatGroup.Agent1ID {
									// 既にカウントしているか確認
									isExist := false
									for otherAgentChatCountI := range unwatchedDetail.OtherAgentChatCounts {
										if unwatchedDetail.OtherAgentChatCounts[otherAgentChatCountI].AgentID == chatGroup.Agent2ID {
											// 元の値を更新する
											unwatchedDetail.OtherAgentChatCounts[otherAgentChatCountI].ChatCount++
											isExist = true
										}
									}
									// まだカウントしていない場合は、新規に追加
									if !isExist {
										unwatchedDetail.OtherAgentChatCounts = append(
											unwatchedDetail.OtherAgentChatCounts,
											otherAgentChat{
												AgentID:   chatGroup.Agent2ID,
												AgentName: chatGroup.Agent2Name,
												ChatCount: 1,
											},
										)
									}
									// 他社エージェントとのチャットの場合　chatGroup.Agent2IDと一致する場合
								} else if agentStaff.AgentID == chatGroup.Agent2ID {
									// 既にカウントしているか確認
									isExist := false
									for otherAgentChatCountI := range unwatchedDetail.OtherAgentChatCounts {
										if unwatchedDetail.OtherAgentChatCounts[otherAgentChatCountI].AgentID == chatGroup.Agent1ID {
											// 元の値を更新する
											unwatchedDetail.OtherAgentChatCounts[otherAgentChatCountI].ChatCount++
											isExist = true
										}
									}
									// まだカウントしていない場合は、新規に追加
									if !isExist {
										unwatchedDetail.OtherAgentChatCounts = append(
											unwatchedDetail.OtherAgentChatCounts,
											otherAgentChat{
												AgentID:   chatGroup.Agent1ID,
												AgentName: chatGroup.Agent1Name,
												ChatCount: 1,
											},
										)
									}
								}
								continue chatMessageLoop
							}
							/*
								未読判定ここまで
							*/
						}
					}
				}
			}
		}

		// 求職者LINE
		for _, chatGroupWithJobSeeker := range chatGroupWithJobSeekerList {
			if agentStaff.ID == uint(chatGroupWithJobSeeker.CAStaffID.Int64) {
				log.Println("agentStaff.ID:", agentStaff.ID, "chatGroupWithJobSeeker.CAStaffID:", chatGroupWithJobSeeker.CAStaffID)
				/*
					未読判定
				*/
				if chatGroupWithJobSeeker.AgentLastWatchedAt.Before(chatGroupWithJobSeeker.JobSeekerLastSendAt) {
					log.Println(
						"未読の求職者LINE", chatGroupWithJobSeeker.ID,
						"chatGroupWithJobSeeker.JobSeekerLastWatchedAt", chatGroupWithJobSeeker.JobSeekerLastWatchedAt,
						"chatGroupWithJobSeeker.AgentLastSendAt", chatGroupWithJobSeeker.AgentLastSendAt,
					)
					unwatchedDetail.AllJobSeekerChatCount++
					unwatchedDetail.JobSeekerNameList += chatGroupWithJobSeeker.LastName + chatGroupWithJobSeeker.FirstName + "様\n\t"
				}
			}
		}

		// タスク *自分が振られているタスクのみ
		for _, task := range taskList {
			if (task.CAStaffID == agentStaff.ID && (task.StaffType == null.NewInt(int64(entity.CA), true) || task.StaffType == null.NewInt(int64(entity.CA_Boss), true))) ||
				(task.RAStaffID == agentStaff.ID && (task.StaffType == null.NewInt(int64(entity.RA), true) || task.StaffType == null.NewInt(int64(entity.RA_Boss), true))) {
				log.Println("タスクID", task.ID, "タスクの期限日", task.DeadlineDay, task.DeadlineTime)
				/*
					期限切れor未読判定
				*/
				// 期限切れ判定
				// 文字列からTime型に変換
				if task.DeadlineDay == "" {
					continue
				} else if task.DeadlineTime.Int64 < 0 || task.DeadlineTime.Int64 > 23 {
					task.DeadlineTime = null.NewInt(0, true)
				}
				deadlineStr := fmt.Sprint(task.DeadlineDay, " ", task.DeadlineTime.Int64, ":00:00")
				deadline, err := time.Parse("2006-01-02 15:00:00", deadlineStr)
				if err != nil {
					log.Println(err)
					return output, err
				}

				// 期限切れてた場合
				if deadline.Before(input.Now) {
					log.Println("期限切れのタスク。deadline", deadline, "input.Now", input.Now)
					unwatchedDetail.TaskCount++
					unwatchedDetail.Tasks = append(unwatchedDetail.Tasks, *task)
				}

				// 	// 未読判定 *判定なしに変更
				// 	// CAタスクの場合
				// } else if task.StaffType == null.NewInt(int64(entity.CA), true) &&
				// 	agentStaff.ID == task.CAStaffID &&
				// 	taskGroup.CALastWatchedAt.Before(taskGroup.RALastRequestAt) {
				// 	log.Println("未読のタスク。taskGroup.CALastWatchedAt", taskGroup.CALastWatchedAt, "taskGroup.RALastRequestAt", taskGroup.RALastRequestAt)
				// 	unwatchedDetail.TaskCount++
				// 	unwatchedDetail.Tasks = append(unwatchedDetail.Tasks, task)

				// 	// RAタスクの場合
				// } else if task.StaffType == null.NewInt(int64(entity.RA), true) &&
				// 	agentStaff.ID == task.RAStaffID &&
				// 	taskGroup.RALastWatchedAt.Before(taskGroup.CALastRequestAt) {
				// 	log.Println("未読のタスク。taskGroup.RALastWatchedAt", taskGroup.RALastWatchedAt, "taskGroup.CALastRequestAt", taskGroup.CALastRequestAt)
				// 	unwatchedDetail.TaskCount++
				// 	unwatchedDetail.Tasks = append(unwatchedDetail.Tasks, task)
				// }
				/*
					期限切れor未読判定ここまで
				*/
			}
		}

		/*
			担当者宛にメール送信
		*/
		// 未読のチャットと未処理のタスクがない場合 or 通知設定がオフの場合はスキップ
		if unwatchedDetail.AllAgentChatCount == 0 && unwatchedDetail.TaskCount == 0 {
			continue
		}

		// 複数の未読チャットの詳細をフォーマットに合わせた文字にする。fmt:株式会社〜〜 %d件
		chatDetailStr := ""
		chatDetailHTML := ""
		// 社内チャットの未読がある場合
		if unwatchedDetail.MyAgentChatCount > 0 {
			chatDetailStr = fmt.Sprintf("社内チャット %d件\n\t", unwatchedDetail.MyAgentChatCount)
		}
		// 複数の他社エージェントとのチャットをフォーマットに合わせた文字にする。fmt:株式会社〜〜 %d件
		for _, otherAgentChatCount := range unwatchedDetail.OtherAgentChatCounts {
			chatDetailStr += fmt.Sprintf(
				"%s %d件\n\t",
				otherAgentChatCount.AgentName, otherAgentChatCount.ChatCount,
			)
			chatDetailHTML += fmt.Sprintf(
				"<li>%s %d件</li>",
				otherAgentChatCount.AgentName, otherAgentChatCount.ChatCount,
			)
		}

		// 複数のタスク詳細をフォーマットに合わせた文字にする。fmt:求人タイトル / 求職者名 / フェーズ
		taskDetailStr := ""
		taskDetailHTML := ""
		for _, task := range unwatchedDetail.Tasks {
			taskPhaseStr, taskPhaseSubStr := getStrTaskPhaseAndPhaseSub(task.PhaseCategory, task.PhaseSubCategory)
			taskDetailStr += fmt.Sprintf(
				"%s / %s%s / %s %s\n\t",
				task.Title, task.LastName, "様", taskPhaseStr, taskPhaseSubStr,
			)
			taskDetailHTML += fmt.Sprintf(
				"<li>%s / %s%s / %s %s</li>",
				task.Title, task.LastName, "様", taskPhaseStr, taskPhaseSubStr,
			)
		}

		// メール送信
		mailBody := fmt.Sprintf(
			"%s株式会社\n%s様\n\n平素よりautoscoutをご利用いただきありがとうございます。\nautoscout事務局でございます。\n\nチャットの未読と期限切れのタスクについてお知らせいたします。\n\n\nエージェントチャットの未読 %d件\nhttps://autoscouts.web.app/chat/agent/\n\t%s\n\n求職者とのLINEの未読 %d件\nhttps://autoscouts.web.app/\n\t%s\n\n期限切れのタスク %d件\nhttps://autoscouts.web.app/\n\t%s\n\n以上でございます。\n入れ違いで処理済みでしたら申し訳ございません。\n\n引き続きどうぞよろしくお願い申し上げます。\n\n未処理・未読の通知設定は、autoscout内の「アカウント設定 / 基本情報」より変更できます。\nhttps://autoscouts.web.app/account/?panel=basic_information",
			agentStaff.AgentName,
			agentStaff.StaffName,
			unwatchedDetail.AllAgentChatCount,
			chatDetailStr,
			unwatchedDetail.AllJobSeekerChatCount,
			unwatchedDetail.JobSeekerNameList,
			unwatchedDetail.TaskCount,
			taskDetailStr,
		)

		// 送信元を作成
		from := mail.Email{
			Name:    "autoscout事務局",
			Address: "info@spaceai.jp",
		}

		// 送信先を作成
		to := mail.Email{
			Name:    agentStaff.StaffName,
			Address: agentStaff.Email,
		}

		// メール送信 本番環境のみ送信
		if os.Getenv("APP_ENV") == "prd" {
			// メール送信
			sendgrid := utility.NewSendGrid(i.sendgrid.APIKey)
			err = sendgrid.SendMail(
				&from,
				&to,
				"未読メッセージと期限切れタスクのお知らせ",
				mailBody,
				"",
			)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
		log.Println("メール送信成功。内容:", mailBody)
	}

	output.OK = true

	return output, nil
}

type UpdateExternalJobInput struct {
	GroupID uint
	Param   entity.ExternalJob
}

type UpdateExternalJobOutput struct {
	OK bool
}

func (i *TaskInteractorImpl) UpdateExternalJob(input UpdateExternalJobInput) (UpdateExternalJobOutput, error) {
	var (
		output UpdateExternalJobOutput
	)

	err := i.taskGroupRepository.UpdateExternalJob(input.GroupID, input.Param)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	jobSeekerScheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupID(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 日程のタイトルも更新
	for _, schedule := range jobSeekerScheduleList {
		lastIndexOfNo := strings.LastIndex(schedule.Title, "の")
		if lastIndexOfNo != -1 {
			// UTF-8では「の」は3バイトなので、lastIndexOfNo + 3から開始
			newTitle := input.Param.ExternalCompanyName + "の" + schedule.Title[lastIndexOfNo+3:]

			if newTitle == schedule.Title {
				continue
			}

			err = i.jobSeekerScheduleRepository.UpdateTitle(schedule.ID, newTitle)
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

// タスクIDを使ってタスク情報を取得する
type GetJobSeekerTaskListByAgentStaffIDInput struct {
	AgentStaffID uint
}

type GetJobSeekerTaskListByAgentStaffIDOutput struct {
	TaskList []*entity.JobSeekerTask
}

func (i *TaskInteractorImpl) GetJobSeekerTaskListByAgentStaffID(input GetJobSeekerTaskListByAgentStaffIDInput) (GetJobSeekerTaskListByAgentStaffIDOutput, error) {
	var (
		output            GetJobSeekerTaskListByAgentStaffIDOutput
		err               error
		jobSeekerTaskList []*entity.JobSeekerTask
	)

	/************ 指定担当の求職者の最新タスクを取得 **************/

	taskList, err := i.taskRepository.GetLatestByStaffID(input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 子テーブル情報を取得するためのIDリストの定義 **************/

	var (
		// タスクIDのリスト
		taskIDList []uint

		// タスクグループIDのリスト
		groupIDList []uint
	)

	// タスクIDのリストを作成
	for _, task := range taskList {
		taskIDList = append(taskIDList, task.ID)
		groupIDList = append(groupIDList, task.TaskGroupID)
	}

	/************ 送付をお勧めする書類の情報を取得 **************/

	isRecommendList, err := i.taskIsRecommendDocumentRepository.GetByTaskIDList(taskIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 送付をお勧めする書類の情報をセット
	for _, task := range taskList {
		for _, isRecommend := range isRecommendList {
			if task.ID == isRecommend.TaskID {
				task.IsRecommendDocument = *isRecommend
			}
		}
	}

	/************ 選考日時の情報を取得 **************/

	// 求職者のスケジュール情報を取得
	scheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupIDList(groupIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, task := range taskList {
		// 確定日時がセットされているか

		for _, schedule := range scheduleList {
			// 「同一タスクグループ」「同一選考フェーズ」
			if task.TaskGroupID == uint(schedule.TaskGroupID.Int64) && schedule.PhaseCategory == task.PhaseCategory {
				// 以下の条件に合う日程を確定日時としてセット
				// 「確定日時がまだセットされていない」&&「タスクグループIDが一致している」&&「日程タイプが確定日時」&&「選考フェーズが一致している」
				if schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByEnterprise, true) {
					task.SelectionDate = *schedule
				}

				if schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByCA, true) {
					// 「確定日時がない」または「確定日時があるけど、繰り返しの回数が違う」
					// if !task.SelectionDate.TaskID.Valid || (task.SelectionDate.TaskID.Valid && task.SelectionDate.RepetitionCount != schedule.RepetitionCount) {
					// 	// 候補日時をセット
					// 	fmt.Println("------------")
					// 	fmt.Println(task.LastName + task.FirstName, task.CompanyName, task.Title)
					// 	task.PossibleDates = append(task.PossibleDates, *schedule)
					// }
					task.PossibleDates = append(task.PossibleDates, *schedule)
				}
			}
		}
	}

	/************ 求職者単位でタスクをまとめる **************/

	// 重複チェックのmapを設定
	duplicateJobSeekerIDDCheck := make(map[uint]bool)

	for _, task := range taskList {
		if !duplicateJobSeekerIDDCheck[task.JobSeekerID] {
			var jobSeekerTask = entity.JobSeekerTask{
				JobSeekerID:   task.JobSeekerID,
				CAStaffID:     task.CAStaffID,
				LastName:      task.LastName,
				FirstName:     task.FirstName,
				LastFurigana:  task.LastFurigana,
				FirstFurigana: task.FirstFurigana,
				TaskCount: entity.TaskCount{
					EntryCount:     0,
					DocumentCount:  0,
					SelectionCount: 0,
					HoldCount:      0,
					AcceptCount:    0,
					UnwatchCount:   0,
				},
				IsSelfApplication: task.IsSelfApplication,
			}
			jobSeekerTaskList = append(jobSeekerTaskList, &jobSeekerTask)

			// 重複チェックを更新
			duplicateJobSeekerIDDCheck[task.JobSeekerID] = true
		}

		var duplicateCheck = make(map[uint]bool)
		for _, jt := range jobSeekerTaskList {
			if jt.JobSeekerID == task.JobSeekerID {
				// タスクフェーズ毎のカウント
				switch task.PhaseCategory {
				case null.NewInt(int64(entity.Entry), true):
					// エントリーのカウント
					jt.TaskCount.EntryCount++

				case null.NewInt(int64(entity.DocumentSelection), true):
					// 書類選考のカウント
					jt.TaskCount.DocumentCount++

				case null.NewInt(int64(entity.FirstSelection), true):
					// 選考のカウント
					jt.TaskCount.SelectionCount++

				case null.NewInt(int64(entity.SecondSelection), true):
					// 選考のカウント
					jt.TaskCount.SelectionCount++

				case null.NewInt(int64(entity.ThirdSelection), true):
					// 選考のカウント
					jt.TaskCount.SelectionCount++

				case null.NewInt(int64(entity.FourthSelection), true):
					// 選考のカウント
					jt.TaskCount.SelectionCount++

				case null.NewInt(int64(entity.FifthSelection), true):
					// 選考のカウント
					jt.TaskCount.SelectionCount++

				case null.NewInt(int64(entity.FinalSelection), true):
					// 選考のカウント
					jt.TaskCount.SelectionCount++

				case null.NewInt(int64(entity.HoldJobOffer), true):
					// 内定保留のカウント
					jt.TaskCount.HoldCount++

				case null.NewInt(int64(entity.AcceptJobOffer), true):
					// 内定承諾のカウント
					jt.TaskCount.AcceptCount++
				}

				// タブが「自分」でタスク担当が自分、または自分担当でなくてもリャンメンの場合
				var IsMyTask = ((task.IsDoubleSided && task.CAStaffID == input.AgentStaffID) ||
					(!task.IsDoubleSided && ((task.StaffType == null.NewInt(int64(entity.CA), true) && task.CAStaffID == input.AgentStaffID) ||
						(task.StaffType == null.NewInt(int64(entity.RA), true) && task.RAStaffID == input.AgentStaffID))))

				// タブが「依頼中」でタスク担当が自分以外
				var IsNotMyTask = (!task.IsDoubleSided && ((task.StaffType == null.NewInt(int64(entity.CA), true) && task.CAStaffID != input.AgentStaffID) ||
					(task.StaffType == null.NewInt(int64(entity.RA), true) && task.RAStaffID != input.AgentStaffID)))

				if !duplicateCheck[task.JobSeekerID] && IsMyTask {
					// 未読のカウント
					var (
						// CAの未読
						unwatchCATask = task.CAStaffID == input.AgentStaffID && task.StaffType == null.NewInt(int64(entity.CA), true) && task.CALastWatchedAt.Before(task.RALastRequestAt)

						// RAの未読
						unwatchRATask = task.RAStaffID == input.AgentStaffID && task.StaffType == null.NewInt(int64(entity.RA), true) && task.RALastWatchedAt.Before(task.CALastRequestAt)
					)

					// 未読のカウントと未読フラグをセット
					if unwatchCATask || unwatchRATask {
						task.IsUnconfirmed = true
						jt.HasUnconfirmed = true
						jt.TaskCount.UnwatchCount++
					}

					// 本日中と期限切れのカウント
					// 現在の日付
					today := time.Now()
					today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

					// 期限
					deadlineDay, err := time.Parse("2006-01-02", task.DeadlineDay)
					if err != nil {
						fmt.Println(err)
						return output, err
					}

					deadlineDay = time.Date(deadlineDay.Year(), deadlineDay.Month(), deadlineDay.Day(), 0, 0, 0, 0, deadlineDay.Location())

					// 差分 = 現在の日時 - 投稿日時
					diff := today.Sub(deadlineDay)
					remainingDays := int(diff.Hours() / 24)

					if remainingDays == 0 {
						// 当日
						jt.TaskCount.TodayCount++
					} else if remainingDays > 0 {
						// 期限切れ
						jt.TaskCount.TimeOutCount++
					}

					duplicateCheck[task.JobSeekerID] = true
					jt.MyTaskList = append(jt.MyTaskList, *task)

				} else if !duplicateCheck[task.JobSeekerID] && IsNotMyTask {
					duplicateCheck[task.JobSeekerID] = true
					jt.RequestTaskList = append(jt.RequestTaskList, *task)
				}
			}
		}
	}

	output.TaskList = jobSeekerTaskList

	return output, nil
}
