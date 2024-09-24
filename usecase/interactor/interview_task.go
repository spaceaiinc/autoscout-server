package interactor

import (
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type InterviewTaskInteractor interface {
	// 汎用系
	CreateNextInterviewTask(input CreateNextInterviewTaskInput) (CreateNextInterviewTaskOutput, error)
	GetLatestAdjustmentTaskListByAgentID(input GetLatestAdjustmentTaskListByAgentIDInput) (GetLatestAdjustmentTaskListByAgentIDOutput, error)       // 面談調整タスク（phase_category = (0 or 1)）の取得
	GetLatestConfirmationTaskListByAgentID(input GetLatestConfirmationTaskListByAgentIDInput) (GetLatestConfirmationTaskListByAgentIDOutput, error) // 参加確認（phase_category = (2 or 3)）タスクの取得
	GetInterviewTaskListByGroupID(input GetInterviewTaskListByGroupIDInput) (GetInterviewTaskListByGroupIDOutput, error)                            // GroupIDから面談調整タスク一覧を取得 *アクティビティ表示に使用

	UpdateInterviewTaskGroupLastWatched(input UpdateInterviewTaskGroupLastWatchedInput) (UpdateInterviewTaskGroupLastWatchedOutput, error) // 最終閲覧時間を更新
	UpdateInterviewTaskCAStaffID(input UpdateInterviewTaskCAStaffIDInput) (UpdateInterviewTaskCAStaffIDOutput, error)
	UpdateInterviewTaskInterviewDate(input UpdateInterviewTaskInterviewDateInput) (UpdateInterviewTaskInterviewDateOutput, error)

	DeleteLatestInterviewTask(input DeleteLatestInterviewTaskInput) (DeleteLatestInterviewTaskOutput, error)
	// Admin API
}

type InterviewTaskInteractorImpl struct {
	firebase                           usecase.Firebase
	sendgrid                           config.Sendgrid
	oneSignal                          config.OneSignal
	interviewTaskRepository            usecase.InterviewTaskRepository
	interviewTaskGroupRepository       usecase.InterviewTaskGroupRepository
	jobSeekerRepository                usecase.JobSeekerRepository
	chatGroupWithJobSeekerRepository   usecase.ChatGroupWithJobSeekerRepository
	agentStaffRepository               usecase.AgentStaffRepository
	chatMessageWithJobSeekerRepository usecase.ChatMessageWithJobSeekerRepository
	emailWithJobSeekerRepository       usecase.EmailWithJobSeekerRepository
}

// InterviewTaskInteractorImpl is an implementation of InterviewTaskInteractor
func NewInterviewTaskInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	itR usecase.InterviewTaskRepository,
	itgR usecase.InterviewTaskGroupRepository,
	jsR usecase.JobSeekerRepository,
	cgR usecase.ChatGroupWithJobSeekerRepository,
	asR usecase.AgentStaffRepository,
	cmjR usecase.ChatMessageWithJobSeekerRepository,
	ewjsR usecase.EmailWithJobSeekerRepository,
) InterviewTaskInteractor {
	return &InterviewTaskInteractorImpl{
		firebase:                           fb,
		sendgrid:                           sg,
		oneSignal:                          os,
		interviewTaskRepository:            itR,
		interviewTaskGroupRepository:       itgR,
		jobSeekerRepository:                jsR,
		chatGroupWithJobSeekerRepository:   cgR,
		agentStaffRepository:               asR,
		chatMessageWithJobSeekerRepository: cmjR,
		emailWithJobSeekerRepository:       ewjsR,
	}
}

/****************************************************************************************/
/// 汎用系
//
// 最終閲覧時間を更新
type CreateNextInterviewTaskInput struct {
	CreateParam entity.NextInterviewTaskParam
}

type CreateNextInterviewTaskOutput struct {
	OK bool
}

func (i *InterviewTaskInteractorImpl) CreateNextInterviewTask(input CreateNextInterviewTaskInput) (CreateNextInterviewTaskOutput, error) {
	var (
		output CreateNextInterviewTaskOutput
		err    error
	)

	// タスクの作成
	nextTask := entity.NewInterviewTask(
		input.CreateParam.InterviewTaskGroupID,
		input.CreateParam.AgentStaffID,
		input.CreateParam.CAStaffID,
		input.CreateParam.PhaseCategory,
		input.CreateParam.PhaseSubCategory,
		input.CreateParam.Remarks,
		input.CreateParam.DeadlineDay,
		input.CreateParam.DeadlineTime,
		input.CreateParam.SelectActionLabel,
	)

	err = i.interviewTaskRepository.Create(nextTask)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var (
		PreparingAfterInterview = null.NewInt(int64(entity.PreparingAfterInterview), true)
		OperatingAfterInterview = null.NewInt(int64(entity.OperatingAfterInterview), true)
		ReleasedAfterInterview  = null.NewInt(int64(entity.ReleasedAfterInterview), true)
	)

	// 面談実施した場合は「input.CreateParam.InterviewDate」を初回面談の日時として保存する
	if input.CreateParam.FirstInterviewDate == utility.EarliestTime() &&
		(input.CreateParam.PhaseCategory == PreparingAfterInterview ||
			input.CreateParam.PhaseCategory == OperatingAfterInterview ||
			input.CreateParam.PhaseCategory == ReleasedAfterInterview) {
		// 初回面談日時と面談日時を更新
		err = i.interviewTaskGroupRepository.UpdateNormalAndFirstInterviewDate(
			input.CreateParam.InterviewTaskGroupID,
			input.CreateParam.InterviewDate, // 初回面談日時に記録する
			input.CreateParam.InterviewDate,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	} else {
		// 面談日時のみを更新
		err = i.interviewTaskGroupRepository.UpdateInterviewDate(input.CreateParam.InterviewTaskGroupID, input.CreateParam.InterviewDate)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	// 依頼時間を更新
	err = i.interviewTaskGroupRepository.UpdateLastRequestAt(input.CreateParam.InterviewTaskGroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者のフェーズを更新
	err = i.jobSeekerRepository.UpdatePhase(
		input.CreateParam.JobSeekerID,
		input.CreateParam.PhaseCategory,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者の面談日時を更新
	err = i.jobSeekerRepository.UpdateInterviewDate(
		input.CreateParam.JobSeekerID,
		input.CreateParam.InterviewDate,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// CA担当者のIDがからでなければ更新
	err = i.jobSeekerRepository.UpdateStaffID(
		input.CreateParam.JobSeekerID,
		input.CreateParam.CAStaffID,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ValidationTypeに応じて「メール送信」と「面談日程更新」の処理を実行
	if input.CreateParam.ValidationType == "email" || input.CreateParam.ValidationType == "emailAndInterview" {
		// メール送信の処理
		err = utility.SendMailToSingle(
			i.sendgrid.APIKey,
			input.CreateParam.Mail.Subject,
			input.CreateParam.Mail.Content,
			input.CreateParam.Mail.From, // From（エージェント名 + エージェントメールアドレス（面談調整用））
			input.CreateParam.Mail.To,   // To（求職者名 + 求職者メールアドレス）
			nil,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メール送信履歴を保存
		emailWithJobSeeker := entity.NewEmailWithJobSeeker(
			input.CreateParam.JobSeekerID,
			input.CreateParam.Mail.Subject,
			input.CreateParam.Mail.Content,
			"",
		)

		err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

	} else if input.CreateParam.ValidationType == "interview" || input.CreateParam.ValidationType == "emailAndInterview" {
		// 面談日時登録
		err = i.interviewTaskGroupRepository.UpdateInterviewDate(
			input.CreateParam.InterviewTaskGroupID,
			input.CreateParam.InterviewDate,
		)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	} else if input.CreateParam.ValidationType == "emailOrLine" {

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(input.CreateParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(uint(input.CreateParam.AgentStaffID.Int64))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(input.CreateParam.AgentID, input.CreateParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// LINE連携ありの場合
		if chatGroup.LineActive {
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
				linebot.NewTextMessage(input.CreateParam.Mail.Content),
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
				input.CreateParam.Mail.Content,
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
				input.CreateParam.Mail.Subject,
				input.CreateParam.Mail.Content,
				input.CreateParam.Mail.From, // From（エージェント名 + エージェントメールアドレス（面談調整用））
				input.CreateParam.Mail.To,   // To（求職者名 + 求職者メールアドレス）
				nil,
			)

			fmt.Println("メール送信")

			if err != nil {
				fmt.Println("メール送信失敗")
				fmt.Println(err)
				return output, err
			} else {
				fmt.Println("メール送信成功")
			}

			// メール送信履歴を保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				input.CreateParam.JobSeekerID,
				input.CreateParam.Mail.Subject,
				input.CreateParam.Mail.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
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
type GetLatestAdjustmentTaskListByAgentIDInput struct {
	AgentID uint
}

type GetLatestAdjustmentTaskListByAgentIDOutput struct {
	InterviewTaskList []*entity.InterviewTask
}

func (i *InterviewTaskInteractorImpl) GetLatestAdjustmentTaskListByAgentID(input GetLatestAdjustmentTaskListByAgentIDInput) (GetLatestAdjustmentTaskListByAgentIDOutput, error) {
	var (
		output            GetLatestAdjustmentTaskListByAgentIDOutput
		err               error
		interviewTaskList []*entity.InterviewTask
	)

	interviewTaskList, err = i.interviewTaskRepository.GetLatestAdjustmentByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewTaskList = interviewTaskList

	return output, nil
}

// タスクIDを使ってタスク情報を取得する
type GetLatestConfirmationTaskListByAgentIDInput struct {
	AgentID uint
}

type GetLatestConfirmationTaskListByAgentIDOutput struct {
	InterviewTaskList []*entity.InterviewTask
}

func (i *InterviewTaskInteractorImpl) GetLatestConfirmationTaskListByAgentID(input GetLatestConfirmationTaskListByAgentIDInput) (GetLatestConfirmationTaskListByAgentIDOutput, error) {
	var (
		output            GetLatestConfirmationTaskListByAgentIDOutput
		err               error
		interviewTaskList []*entity.InterviewTask
	)

	interviewTaskList, err = i.interviewTaskRepository.GetLatestConfirmationByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewTaskList = interviewTaskList

	return output, nil
}

// 最終閲覧時間を更新
type UpdateInterviewTaskGroupLastWatchedInput struct {
	GroupID uint
}

type UpdateInterviewTaskGroupLastWatchedOutput struct {
	OK bool
}

func (i *InterviewTaskInteractorImpl) UpdateInterviewTaskGroupLastWatched(input UpdateInterviewTaskGroupLastWatchedInput) (UpdateInterviewTaskGroupLastWatchedOutput, error) {
	var (
		output UpdateInterviewTaskGroupLastWatchedOutput
		err    error
	)

	err = i.interviewTaskGroupRepository.UpdateLastWatchedAt(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// GroupIDから面談調整タスク一覧を取得 *アクティビティ表示に使用
type GetInterviewTaskListByGroupIDInput struct {
	GroupID uint
}

type GetInterviewTaskListByGroupIDOutput struct {
	InterviewTaskList []*entity.InterviewTask
}

func (i *InterviewTaskInteractorImpl) GetInterviewTaskListByGroupID(input GetInterviewTaskListByGroupIDInput) (GetInterviewTaskListByGroupIDOutput, error) {
	var (
		output            GetInterviewTaskListByGroupIDOutput
		err               error
		interviewTaskList []*entity.InterviewTask
	)

	interviewTaskList, err = i.interviewTaskRepository.GetByGroupID(input.GroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewTaskList = interviewTaskList

	return output, nil
}

type UpdateInterviewTaskCAStaffIDInput struct {
	Param entity.UpdateCAStaffIDParam
}

type UpdateInterviewTaskCAStaffIDOutput struct {
	InterviewTask *entity.InterviewTask
}

func (i *InterviewTaskInteractorImpl) UpdateInterviewTaskCAStaffID(input UpdateInterviewTaskCAStaffIDInput) (UpdateInterviewTaskCAStaffIDOutput, error) {
	var (
		output UpdateInterviewTaskCAStaffIDOutput
		err    error
	)

	// 求職者テーブルの担当者IDを更新
	err = i.jobSeekerRepository.UpdateStaffID(input.Param.JobSeekerID, input.Param.CAStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	task, err := i.interviewTaskRepository.FindByID(input.Param.InterviewTaskID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewTask = task

	return output, nil
}

type UpdateInterviewTaskInterviewDateInput struct {
	Param entity.UpdateInterviewDateParam
}

type UpdateInterviewTaskInterviewDateOutput struct {
	InterviewTask *entity.InterviewTask
}

func (i *InterviewTaskInteractorImpl) UpdateInterviewTaskInterviewDate(input UpdateInterviewTaskInterviewDateInput) (UpdateInterviewTaskInterviewDateOutput, error) {
	var (
		output UpdateInterviewTaskInterviewDateOutput
		err    error
	)

	// タスクグループの
	err = i.interviewTaskGroupRepository.UpdateInterviewDate(input.Param.InterviewTaskGroupID, input.Param.InterviewDate)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	task, err := i.interviewTaskRepository.FindByID(input.Param.InterviewTaskID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.InterviewTask = task

	return output, nil
}

// 面談調整タスクを削除
type DeleteLatestInterviewTaskInput struct {
	Param entity.DeleteLatestInterviewTaskParam
}

type DeleteLatestInterviewTaskOutput struct {
	OK bool
}

func (i *InterviewTaskInteractorImpl) DeleteLatestInterviewTask(input DeleteLatestInterviewTaskInput) (DeleteLatestInterviewTaskOutput, error) {
	var (
		output DeleteLatestInterviewTaskOutput
		err    error
	)

	interviewTaskGroup, err := i.interviewTaskGroupRepository.FindByID(input.Param.InterviewTaskGroupID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	latestInterviewTask, err := i.interviewTaskRepository.FindLatestByAgentIDAndJobSeekerID(
		interviewTaskGroup.AgentID,
		interviewTaskGroup.JobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// DBの最新タスクと削除対象のタスクが一致しない場合はエラー
	if latestInterviewTask.ID != input.Param.InterviewTaskID {
		fmt.Println("面談調整タスクが最新ではありません。latestInterviewTask.ID:", latestInterviewTask.ID, "input.Param.InterviewTaskID:", input.Param.InterviewTaskID)
		return output, nil
	}

	err = i.interviewTaskRepository.Delete(latestInterviewTask.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 削除後の最新のタスクを取得
	latestInterviewTaskAfterDelete, err := i.interviewTaskRepository.FindLatestByAgentIDAndJobSeekerID(
		interviewTaskGroup.AgentID,
		interviewTaskGroup.JobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求職者のフェーズを更新
	err = i.jobSeekerRepository.UpdatePhase(
		interviewTaskGroup.JobSeekerID,
		latestInterviewTaskAfterDelete.PhaseCategory,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}
