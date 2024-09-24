package interactor

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"gopkg.in/guregu/null.v4"
)

// 指定した数字に対応する文字列を取得
// タスクのフェーズ
func getStrTaskPhase(inputNumber null.Int) string {
	if inputNumber.Valid {
		for i, v := range entity.TaskPhase {
			if inputNumber == null.NewInt(int64(i), true) {
				return v
			}
		}
	}

	return ""
}

// タスクのフェーズとサブフェーズ
func getStrTaskPhaseAndPhaseSub(inputPhase, inputPhaseSub null.Int) (outputPhase, outputPhaseSub string) {

	if inputPhase.Valid {
		for phaseI, phaseV := range entity.TaskPhase {
			if inputPhase == null.NewInt(int64(phaseI), true) {
				outputPhase = phaseV

				for phaseSubI, phaseSubV := range entity.TaskPhaseSub[uint(phaseI)] {
					if inputPhaseSub == null.NewInt(int64(phaseSubI), true) {
						outputPhaseSub = phaseSubV
						return outputPhase, outputPhaseSub
					}
				}

			}
		}
	}

	return outputPhase, outputPhaseSub
}

// タスクグループの最大ページ数を返す（本番実装までは1ページあたり20件）
func getTaskGroupListMaxPage(taskList []*entity.TaskGroup) uint {
	var maxPage = len(taskList) / 20

	if 0 < (len(taskList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページのタスクグループ一覧を返す（本番実装までは1ページあたり20件）
func getTaskGroupListWithPage(taskGroupList []*entity.TaskGroup, page uint) []*entity.TaskGroup {
	var (
		perPage uint = 20
		listLen uint = uint(len(taskGroupList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return taskGroupList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.TaskGroup{}
	}

	if (listLen - first) <= perPage {
		return taskGroupList[first:]
	}
	return taskGroupList[first:last]
}

// タスクグループの最大ページ数を返す（本番実装までは1ページあたり20件）
func getTaskListMaxPage(taskList []*entity.Task) uint {
	var maxPage = len(taskList) / 20

	if 0 < (len(taskList) % 20) {
		maxPage++
	}

	return uint(maxPage)
}

// 指定ページのタスクグループ一覧を返す（本番実装までは1ページあたり20件）
func getTaskListWithPage(taskList []*entity.Task, page uint) []*entity.Task {
	var (
		perPage uint = 20
		listLen uint = uint(len(taskList))
		first        = (page * perPage) - perPage
		last         = (page * perPage)
	)

	if listLen <= perPage {
		return taskList[0:]
	}

	// リストが開始位置より少ない場合は空のスライスを返す
	if listLen <= first {
		return []*entity.Task{}
	}

	if (listLen - first) <= perPage {
		return taskList[first:]
	}
	return taskList[first:last]
}

// 次のフェーズが「1次選考 ~ 5次選考」
func isNextTaskPhaseNumberSelection(nextTaskPhase null.Int) bool {
	return (nextTaskPhase == null.NewInt(int64(entity.FinalSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.FirstSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.SecondSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.ThirdSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.FourthSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.FifthSelection), true))
}

// 引数のフェーズが「1次選考 ~ 最終選考」
func isPhaseSelection(nextTaskPhase null.Int) bool {
	return (nextTaskPhase == null.NewInt(int64(entity.FinalSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.FirstSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.SecondSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.ThirdSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.FourthSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.FifthSelection), true) ||
		nextTaskPhase == null.NewInt(int64(entity.FinalSelection), true))
}

// 「2023-06-29T17:00」を「2023年6月29日(日) 17:00」に変更
func dateTimeAndDayNameFormat(date string) string {
	weekdays := []string{"日", "月", "火", "水", "木", "金", "土"}

	resDate, err := time.Parse(time.RFC3339, date+":00Z")
	if err != nil {
		fmt.Println()
		return ""
	}

	year := resDate.Year()
	month := resDate.Month()
	day := resDate.Day()
	hour := resDate.Hour()
	minute := resDate.Minute()
	dayIndex := resDate.Weekday()
	dayName := weekdays[dayIndex]

	hours := strconv.Itoa(hour)
	minutes := strconv.Itoa(minute)

	if hour < 10 {
		hours = "0" + hours
	}
	if minute < 10 {
		minutes = "0" + minutes
	}

	dateString := fmt.Sprintf("%d年%d月%d日(%s) %s:%s", year, month, day, dayName, hours, minutes)

	return dateString
}

func convertPastDateTime(startDateTimeStr, endDateTimeStr string) string {
	start := dateTimeAndDayNameFormat(startDateTimeStr)
	end := dateTimeAndDayNameFormat(endDateTimeStr)

	if start == "" || end == "" {
		return ""
	}

	if start[0:10] == end[0:10] {
		return fmt.Sprintf("%s ~ %s", start, end[len(end)-5:])
	}

	return fmt.Sprintf("%s ~ %s", start, end)
}

// 求職者のスケジュールで同じ日時の確定日時とリスケ日時があるかをチェック
// func findResucheduleSelectionDate(scheduleList []*entity.JobSeekerSchedule, targetSchedule *entity.JobSeekerSchedule) bool {
// 	var encounteredByScheduleID = map[uint]bool{}

// 	for _, schedule := range scheduleList {
// 		// 同じ日時の選考とリスケがあるか
// 		var isFind = targetSchedule.TaskGroupID == schedule.TaskGroupID &&
// 			targetSchedule.StartTime == schedule.StartTime &&
// 			targetSchedule.EndTime == schedule.EndTime && ((targetSchedule.ScheduleType == null.NewInt(entity.ScheduleTypeByEnterprise, true) &&
// 			schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByReschedule, true)) ||
// 			(schedule.ScheduleType == null.NewInt(entity.ScheduleTypeByEnterprise, true) &&
// 				targetSchedule.ScheduleType == null.NewInt(entity.ScheduleTypeByReschedule, true)))

// 		if !encounteredByScheduleID[targetSchedule.ID] {
// 			encounteredByScheduleID[targetSchedule.ID] = isFind
// 		}
// 	}

// 	return encounteredByScheduleID[targetSchedule.ID]
// }

// エントリーフェーズのタスクを作成する処理
func createEntryPhaseTask(i *TaskInteractorImpl, batchParam entity.CreateTaskInBatchProcessingParam, taskParam entity.TaskParam) error {
	var (
		err          error
		newTask      *entity.Task
		nextPhase    = taskParam.PhaseCategory
		nextPhaseSub = taskParam.PhaseSubCategory
		prevPhase    = taskParam.PrevPhaseCategory
		prevPhaseSub = taskParam.PrevPhaseSubCategory
	)

	if prevPhase != null.NewInt(int64(entity.Entry), true) {
		return errors.New("現在のフェーズがエントリーではありません")
	}

	/************ タスク情報の作成（選考継続時のタスク作成 or 通常タスク作成） **************/

	// アクション: 選考継続処理
	// エントリー/辞退意思確認中（選考継続の場合）->直前の○次選考のタスクに戻る
	if nextPhase == null.NewInt(int64(entity.Entry), true) && nextPhaseSub == null.NewInt(int64(entity.ContinueSelection), true) {
		// 選考継続の場合は辞退処理をする前のタスクを取得
		continueTask, err := i.taskRepository.FindLatestForContinue(batchParam.TaskGroupID, nextPhase)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 取得した過去のタスクを新規タスクとして作成
		newTask = entity.NewTask(
			continueTask.TaskGroupID,
			continueTask.PhaseCategory,    // 過去タスクのフェーズ
			continueTask.PhaseSubCategory, // 過去タスクのサブフェーズ
			continueTask.StaffType,
			continueTask.ExecutedStaffID,
			taskParam.Remarks,
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			taskParam.TalkAboutInInterview,
			taskParam.ScheduleCollectionCondition,
			taskParam.ExamGuideContent,
			false,
		)

	} else {
		// 次の選考が「エントリー/エントリー保留」or「エントリー/エントリー辞退」かを確認
		var IsDeclineOrHoldEntry = nextPhase == null.NewInt(int64(entity.Entry), true) && (nextPhaseSub == null.NewInt(int64(entity.DeclineEntry), true) || nextPhaseSub == null.NewInt(int64(entity.HoldEntry), true))

		// 次の選考が「書類選考/エントリー依頼」かを確認
		var IsRequestEntry = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) && nextPhaseSub == null.NewInt(int64(entity.RequestEntry), true)

		// 直前のタスクが応募意思確認中かを確認する
		if IsDeclineOrHoldEntry || IsRequestEntry {
			prevTask, err := i.taskRepository.FindLatestByGroupID(batchParam.TaskGroupID)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// 直前のタスクが「エントリー/求人打診依頼」の場合は「エントリー/応募意思確認中」のタスクを作成する
			if prevTask.PhaseCategory == null.NewInt(int64(entity.Entry), true) && prevTask.PhaseSubCategory == null.NewInt(int64(entity.SoundOutJobInformation), true) {
				// 「エントリー/応募意思確認中」のタスクを作成
				confirmApplicationIntentionTask := entity.NewTask(
					batchParam.TaskGroupID,
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
					return err
				}
			}
		}

		// 実際に進めるタスクを作成
		newTask = entity.NewTask(
			batchParam.TaskGroupID,
			nextPhase,
			nextPhaseSub,
			taskParam.StaffType, // 担当者タイプ
			taskParam.ExecutedStaffID,
			taskParam.Remarks,
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			taskParam.TalkAboutInInterview,
			taskParam.ScheduleCollectionCondition,
			taskParam.ExamGuideContent,
			taskParam.IsCheckDoubleSided,
		)
	}

	// タスク作成
	err = i.taskRepository.Create(newTask)
	if err != nil {
		fmt.Println(err)
		return err
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
			newTask.ID,
			NullInt, // エントリーフェーズの時は選考フローを選択していないためnullが入る
			batchParam.JobSeekerID,
			batchParam.JobInformationID,
			taskParam.GoodPoint,
			taskParam.NGPoint,
			true,  // 選考通過のためtrue
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else if IsMaskResumeFailNext {
		// マスクレジュメ不合格の場合 不合格時の評価点の作成
		evaluationPoint := entity.NewEvaluationPoint(
			newTask.ID,
			NullInt, // エントリーフェーズの時は選考フローを選択していないためnullが入る
			batchParam.JobSeekerID,
			batchParam.JobInformationID,
			taskParam.GoodPoint,
			taskParam.NGPoint,
			false, // 選考不合格のためfalse
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ マスクレジュメ打診時のメール送付 **************/

	// マスクレジュメを送信するかどうか
	var IsSendMaskResume = taskParam.TaskOption == null.NewInt(entity.OptionSendMessage, true) &&
		nextPhase == null.NewInt(int64(entity.Entry), true) &&
		nextPhaseSub == null.NewInt(int64(entity.CollectResultOfMask), true)

	if IsSendMaskResume {
		agentStaff, err := i.agentStaffRepository.FindByID(batchParam.RAStaffID)
		if err != nil {
			fmt.Println(err)
			fmt.Println(agentStaff)

			return err
		}

		err = utility.SendMailToMultiple(
			i.sendgrid.APIKey,
			taskParam.Subject,
			taskParam.Content,
			entity.EmailUser{
				Name:  agentStaff.StaffName,
				Email: agentStaff.Email,
			},
			taskParam.Tos,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ おすすめ応募書類の情報を登録する **************/

	// 次のタスクが「書類選考/推薦依頼」
	if nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true) {

		// 入力された送付をお勧めする書類の情報
		isR := taskParam.IsRecommendDocument

		// 送付をお勧めする書類の情報を取得
		isRecommend, err := i.taskIsRecommendDocumentRepository.FindByTaskID(newTask.ID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {

				// 作成する
				isRecommend = entity.NewTaskIsRecommendDocument(
					newTask.ID,
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
					return err
				}
			} else {
				fmt.Println(err)
				return err
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
				return err
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
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(batchParam.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, possibleDate := range taskParam.PossibleDates {
			var taskID null.Int

			if possibleDate.TaskID.Valid {
				// TaskIDがある場合は登録済みの候補日時なので、既存のTaskIDを使用する
				taskID = possibleDate.TaskID
			} else {
				// TaskIDが無い場合は未登録の候補日時なので、新しいTaskIDを使用する
				taskID = null.NewInt(int64(newTask.ID), true)
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
				return err
			}
		}
	}

	/************ 両面タスクに更新する **************/

	// 次のタスクが「書類選考/書類推薦依頼」の場合
	var IsNextRequestRecommendationsn = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsNextRequestRecommendationsn && taskParam.IsCheckDoubleSided {
		// 両面タスクがTrueの場合は「task_groupテーブルのis_double_sidedカラム」をTrueで更新する。
		err = i.taskGroupRepository.UpdateIsDoubleSided(batchParam.TaskGroupID, true)
	}

	/************ 求人票のURLを求職者へ送信（LINE or メール） **************/

	// 求人票を送信するかどうか（LINE or メール）
	var IsSendMessageForJobSeeker = taskParam.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		nextPhase == null.NewInt(int64(entity.Entry), true) &&
		nextPhaseSub == null.NewInt(int64(entity.ConfirmApplicationIntention), true)

	if IsSendMessageForJobSeeker {
		var message = taskParam.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(batchParam.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return err
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
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
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
					Email: message.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				batchParam.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

// 書類選考フェーズのタスクを作成する処理
func createDocumentSelectionPhaseTask(i *TaskInteractorImpl, batchParam entity.CreateTaskInBatchProcessingParam, taskParam entity.TaskParam) error {
	var (
		err          error
		newTask      *entity.Task
		nextPhase    = taskParam.PhaseCategory
		nextPhaseSub = taskParam.PhaseSubCategory
		prevPhase    = taskParam.PrevPhaseCategory
		prevPhaseSub = taskParam.PrevPhaseSubCategory
	)

	if prevPhase != null.NewInt(int64(entity.DocumentSelection), true) {
		return errors.New("現在のフェーズが書類選考ではありません")
	}

	/************ タスク情報の作成（選考継続時のタスク作成 or 通常タスク作成） **************/

	if nextPhase == null.NewInt(int64(entity.DocumentSelection), true) && nextPhaseSub == null.NewInt(int64(entity.ContinueSelection), true) {
		// 選考継続の場合は辞退処理をする前のタスクを取得
		continueTask, err := i.taskRepository.FindLatestForContinue(batchParam.TaskGroupID, nextPhase)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 取得した過去のタスクを新規タスクとして作成
		newTask = entity.NewTask(
			continueTask.TaskGroupID,
			continueTask.PhaseCategory,    // 過去タスクのフェーズ
			continueTask.PhaseSubCategory, // 過去タスクのサブフェーズ
			continueTask.StaffType,
			continueTask.ExecutedStaffID,
			taskParam.Remarks,
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			taskParam.TalkAboutInInterview,
			taskParam.ScheduleCollectionCondition,
			taskParam.ExamGuideContent,
			false,
		)
	} else {
		// タスクを作成
		newTask = entity.NewTask(
			batchParam.TaskGroupID,
			nextPhase,
			nextPhaseSub,
			taskParam.StaffType, // 担当者タイプ
			taskParam.ExecutedStaffID,
			taskParam.Remarks,
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			taskParam.TalkAboutInInterview,
			taskParam.ScheduleCollectionCondition,
			taskParam.ExamGuideContent,
			taskParam.IsCheckDoubleSided,
		)
	}

	// タスク作成
	err = i.taskRepository.Create(newTask)
	if err != nil {
		fmt.Println(err)
		return err
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

		if batchParam.SelectionFlowPatternID.Valid {
			selectionInformation, err := i.jobInfoSelectionInformationRepository.FindBySelectionFlowIDAndSelectionType(uint(batchParam.SelectionFlowPatternID.Int64), uint(entity.DocumentSelection))
			if err != nil {
				fmt.Println(err)
				return err
			}

			selectionInformationID = null.NewInt(int64(selectionInformation.ID), true)
		}

		// 書類選考合格の場合「合格時の評価点を作成 + 選考フローパターンの紐付け」
		evaluationPoint := entity.NewEvaluationPoint(
			newTask.ID,
			selectionInformationID, // エントリーフェーズの時は選考フローを選択していないためnullが入る
			batchParam.JobSeekerID,
			batchParam.JobInformationID,
			taskParam.GoodPoint,
			taskParam.NGPoint,
			true,  // 選考通過のためtrue
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 選考フローパターンの紐付け
		err = i.taskGroupRepository.UpdateSelectionFlowPatternID(batchParam.TaskGroupID, batchParam.SelectionFlowPatternID)
		if err != nil {
			fmt.Println(err)
			return err
		}

	} else if IsDocumentSelectionFailNext {
		// 書類選考不合格の場合 「不合格時の評価点の作成」
		evaluationPoint := entity.NewEvaluationPoint(
			newTask.ID,
			null.NewInt(0, false), // エントリーフェーズの時は選考フローを選択していないためnullが入る
			batchParam.JobSeekerID,
			batchParam.JobInformationID,
			taskParam.GoodPoint,
			taskParam.NGPoint,
			false, // 選考不合格のためfalse
			false, // 再面接は発生しないためfalse
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ 推薦メールの送信 **************/

	// 推薦メールを送信するかどうか
	var IsSendRecommendationMail = taskParam.TaskOption == null.NewInt(entity.OptionSendMessage, true) &&
		prevPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		prevPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsSendRecommendationMail {
		// RA担当者情報の取得
		agentStaff, err := i.agentStaffRepository.FindByID(batchParam.RAStaffID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// RAから企業にメールを送信
		err = utility.SendMailToMultiple(
			i.sendgrid.APIKey,
			taskParam.Subject,
			taskParam.Content,
			entity.EmailUser{
				Name:  agentStaff.StaffName,
				Email: agentStaff.Email,
			},
			taskParam.Tos,
			nil,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ おすすめ応募書類の情報を登録する **************/

	// 次のタスクが「書類選考/推薦依頼」
	if nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true) {

		// 入力された送付をお勧めする書類の情報
		isR := taskParam.IsRecommendDocument

		// 送付をお勧めする書類の情報を取得
		isRecommend, err := i.taskIsRecommendDocumentRepository.FindByTaskID(newTask.ID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {

				// 作成する
				isRecommend = entity.NewTaskIsRecommendDocument(
					newTask.ID,
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
					return err
				}
			} else {
				fmt.Println(err)
				return err
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
				return err
			}
		}
	}

	/************ 選考の候補日時を登録する **************/

	// 次のフェーズが「書類選考/推薦依頼」かどうか
	var IsNextPhaseRequestRecommendations = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsNextPhaseRequestRecommendations {
		// 登録前にすでに登録済みの候補日の情報を削除する
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(batchParam.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, possibleDate := range taskParam.PossibleDates {
			var taskID null.Int

			if possibleDate.TaskID.Valid {
				// TaskIDがある場合は登録済みの候補日時なので、既存のTaskIDを使用する
				taskID = possibleDate.TaskID
			} else {
				// TaskIDが無い場合は未登録の候補日時なので、新しいTaskIDを使用する
				taskID = null.NewInt(int64(newTask.ID), true)
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
				return err
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

		if taskParam.SelectionDate.TaskID.Valid {
			// TaskIDがある場合は登録済みの確定日時なので、既存のTaskIDを使用する
			taskID = taskParam.SelectionDate.TaskID
		} else {
			// TaskIDが無い場合は未登録の確定日時なので、新しいTaskIDを使用する
			taskID = null.NewInt(int64(newTask.ID), true)
		}

		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			taskParam.SelectionDate.JobSeekerID,
			taskID,
			taskParam.SelectionDate.ScheduleType,
			taskParam.SelectionDate.Title,
			taskParam.SelectionDate.StartTime,
			taskParam.SelectionDate.EndTime,
			taskParam.SelectionDate.SeekerDescription,
			taskParam.SelectionDate.StaffDescription,
			taskParam.SelectionDate.IsShare,
			taskParam.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ 両面タスクに更新する **************/

	// 次のタスクが「書類選考/書類推薦依頼」の場合
	var IsNextRequestRecommendationsn = nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestRecommendations), true)

	if IsNextRequestRecommendationsn && taskParam.IsCheckDoubleSided {
		// 両面タスクがTrueの場合は「task_groupテーブルのis_double_sidedカラム」をTrueで更新する。
		err = i.taskGroupRepository.UpdateIsDoubleSided(batchParam.TaskGroupID, true)
	}

	/************ 応募書類アップロードURLを求職者へ送信（LINE or メール） **************/

	// 求人票を送信するかどうか（LINE or メール）
	var IsSendMessageForJobSeeker = taskParam.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		nextPhase == null.NewInt(int64(entity.DocumentSelection), true) &&
		nextPhaseSub == null.NewInt(int64(entity.PrepareDocument), true)

	if IsSendMessageForJobSeeker {
		var message = taskParam.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(batchParam.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return err
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
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
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
					Email: message.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				batchParam.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

// 選考フェーズのタスクを作成する処理
func createSelectionPhaseTask(i *TaskInteractorImpl, batchParam entity.CreateTaskInBatchProcessingParam, taskParam entity.TaskParam) error {
	var (
		err          error
		newTask      *entity.Task
		nextPhase    = taskParam.PhaseCategory
		nextPhaseSub = taskParam.PhaseSubCategory
		prevPhase    = taskParam.PrevPhaseCategory
		prevPhaseSub = taskParam.PrevPhaseSubCategory
	)

	fmt.Println(nextPhase, nextPhaseSub, prevPhase, prevPhaseSub)

	// 選考フェーズでなかった場合はエラー
	if !(int64(entity.FirstSelection) <= prevPhase.Int64) && !(prevPhase.Int64 >= int64(entity.FinalSelection)) {
		return errors.New("現在フェーズが選考フェーズではありません")
	}

	/************ 新リスケ処理（リスケ + 削除） 元のリスケ処理も残しておく **************/

	if taskParam.TaskOption == null.NewInt(entity.OptionRescheduleAndDelete, true) {
		// 直前の結果回収依頼までの日程を削除 + リスケタスクの作成
		latestSelectTask, err := i.taskRepository.FindLatestByGroupIDAndCollectResultPhase(batchParam.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 直前の結果回収依頼までの日程を削除
		err = i.jobSeekerScheduleRepository.DeleteScheduleInTaskGroupAboveTaskID(latestSelectTask.ID, batchParam.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// リスケタスクを作成
		rescheduleTask := entity.NewTask(
			batchParam.TaskGroupID,
			nextPhase,
			null.NewInt(entity.RescheduleDeleteOfSelection, true),
			taskParam.StaffType, // 担当者タイプ
			taskParam.ExecutedStaffID,
			"",
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			"",
			"",
			"",
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(rescheduleTask)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ タスク情報の作成（選考継続時のタスク作成 or 通常タスク作成（+スキップタスク）） **************/

	if (isPhaseSelection(nextPhase)) &&
		nextPhaseSub == null.NewInt(int64(entity.ContinueSelection), true) {

		// 選考継続の場合は辞退処理をする前のタスクを取得
		continueTask, err := i.taskRepository.FindLatestForContinue(batchParam.TaskGroupID, nextPhase)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 取得した過去のタスクを新規タスクとして作成
		newTask = entity.NewTask(
			continueTask.TaskGroupID,
			continueTask.PhaseCategory,    // 過去タスクのフェーズ
			continueTask.PhaseSubCategory, // 過去タスクのサブフェーズ
			continueTask.StaffType,
			continueTask.ExecutedStaffID,
			taskParam.Remarks,
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			taskParam.TalkAboutInInterview,
			taskParam.ScheduleCollectionCondition,
			taskParam.ExamGuideContent,
			false,
		)
	} else {
		// タスクを作成
		newTask = entity.NewTask(
			batchParam.TaskGroupID,
			nextPhase,
			nextPhaseSub,
			taskParam.StaffType, // 担当者タイプ
			taskParam.ExecutedStaffID,
			taskParam.Remarks,
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			taskParam.TalkAboutInInterview,
			taskParam.ScheduleCollectionCondition,
			taskParam.ExamGuideContent,
			false,
		)
	}

	// 選考スキップタスクの場合
	if taskParam.TaskOption == null.NewInt(entity.OptionSkipSelection, true) {
		prevTask, err := i.taskRepository.FindLatestByGroupID(batchParam.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// スキップされるフェーズは現在までの最新タスクフェーズの次のフェーズ
		skipPhase := prevTask.PhaseCategory.Int64 + 1

		// 作成されるタスクが内定保留タスクの場合はスキップフェーズを最終選考に変更
		if nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) {
			skipPhase = int64(entity.FinalSelection)
		}

		skipTask := entity.NewTask(
			batchParam.TaskGroupID,
			null.NewInt(skipPhase, true),            // スキップされるフェーズ
			null.NewInt(entity.SelectionSkip, true), // 選考スキップ
			null.NewInt(int64(entity.CA), true),     // 担当者タイプ
			batchParam.CAStaffID,
			"",
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			"",
			"",
			"",
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(skipTask)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// タスク作成
	err = i.taskRepository.Create(newTask)
	if err != nil {
		fmt.Println(err)
		return err
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
			newTask.ID,
			taskParam.SelectionInformationID,
			batchParam.JobSeekerID,
			batchParam.JobInformationID,
			taskParam.GoodPoint,
			taskParam.NGPoint,
			isPass,        // 選考通過のためtrue
			isReInterview, // 再面接かどうかで値が変化
		)

		err = i.evaluationPointRepository.Create(evaluationPoint)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ 選考の候補日時を登録する **************/

	// 次のフェーズが「○次選考/日程調整依頼」かどうか
	var IsNextPhaseRequestScheduleAdjustment = isPhaseSelection(nextPhase) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestScheduleAdjustmentForSelection), true)

	if IsNextPhaseRequestScheduleAdjustment {
		// 登録前にすでに登録済みの候補日の情報を削除する
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(batchParam.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, possibleDate := range taskParam.PossibleDates {
			var taskID null.Int

			if possibleDate.TaskID.Valid {
				// TaskIDがある場合は登録済みの候補日時なので、既存のTaskIDを使用する
				taskID = possibleDate.TaskID
			} else {
				// TaskIDが無い場合は未登録の候補日時なので、新しいTaskIDを使用する
				taskID = null.NewInt(int64(newTask.ID), true)
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
				return err
			}
		}
	}

	/************ 選考の確定日時を登録する **************/

	// 現在のフェーズが「書類選考/結果回収中」で次が「○次選考/日程・詳細確認依頼」or「○次選考/日程確認依頼」かどうか
	var IsNextSelectionPhase = (isPhaseSelection(nextPhase) && (nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleForSelection), true) ||
		nextPhaseSub == null.NewInt(int64(entity.RequestConfirmScheduleAndSelectionDetail), true)))

	if IsNextSelectionPhase {
		var taskID null.Int

		if taskParam.SelectionDate.TaskID.Valid {
			// TaskIDがある場合は登録済みの確定日時なので、既存のTaskIDを使用する
			taskID = taskParam.SelectionDate.TaskID
		} else {
			// TaskIDが無い場合は未登録の確定日時なので、新しいTaskIDを使用する
			taskID = null.NewInt(int64(newTask.ID), true)
		}

		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			taskParam.SelectionDate.JobSeekerID,
			taskID,
			taskParam.SelectionDate.ScheduleType,
			taskParam.SelectionDate.Title,
			taskParam.SelectionDate.StartTime,
			taskParam.SelectionDate.EndTime,
			taskParam.SelectionDate.SeekerDescription,
			taskParam.SelectionDate.StaffDescription,
			taskParam.SelectionDate.IsShare,
			taskParam.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ 選考の確定日時と当日詳細を更新する **************/

	// 次のフェーズが「○次選考/詳細案内依頼」かどうか
	var IsNextRequestDetailed = (isPhaseSelection(nextPhase) &&
		nextPhaseSub == null.NewInt(int64(entity.RequestDetailedInformation), true))

	if IsNextRequestDetailed {
		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			taskParam.SelectionDate.JobSeekerID,
			taskParam.SelectionDate.TaskID,
			taskParam.SelectionDate.ScheduleType,
			taskParam.SelectionDate.Title,
			taskParam.SelectionDate.StartTime,
			taskParam.SelectionDate.EndTime,
			taskParam.SelectionDate.SeekerDescription,
			taskParam.SelectionDate.StaffDescription,
			taskParam.SelectionDate.IsShare,
			taskParam.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Update(taskParam.SelectionDate.ID, jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ リスケ時の処理（確定日時と同じ日時にリスケ日程を作成） **************/

	// 前のフェーズが「○次選考/日程・詳細確認中 or 日程確認中」で次のタスクオプションががリスケの場合

	var IsResuchedule = isPhaseSelection(prevPhase) && (prevPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleForSelection), true) ||
		prevPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleAndSelectionDetail), true)) &&
		taskParam.TaskOption == null.NewInt(entity.OptionReschedule, true)

	if IsResuchedule {
		// 確定日時とリスケのリストを取得(id降順で取得)
		scheduleList, err := i.jobSeekerScheduleRepository.GetByTaskGroupIDAndScheduleType(batchParam.TaskGroupID, uint(entity.ScheduleTypeByEnterprise))
		if err != nil {
			fmt.Println(err)
			return err
		}

		var filteringScheduleList []*entity.JobSeekerSchedule

		// リスケになっていない同じ選考フェーズの確定日時だけ抜き出す
		for _, schedule := range scheduleList {
			if !schedule.RescheduleID.Valid && schedule.PhaseCategory == taskParam.PhaseCategory {
				fmt.Println(schedule.Title)
				filteringScheduleList = append(filteringScheduleList, schedule)
			}
		}

		// 確定日時が取得でいているかを確認
		if len(filteringScheduleList) > 0 {
			// 確定日時
			var selectionSchedule = filteringScheduleList[0]

			//　リスケデータを作成
			reschedule := entity.NewJobSeekerReschedule(selectionSchedule.ID, newTask.ID)

			i.jobSeekerRescheduleRepository.Create(reschedule)
			if err != nil {
				fmt.Println(err)
				return err
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
	var IsSendMessage = taskParam.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		isPhaseSelection(nextPhase) &&
		(nextPhaseSub == null.NewInt(int64(entity.JobSeekerSupport), true) ||
			nextPhaseSub == null.NewInt(int64(entity.CollectingSchedule), true) ||
			nextPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleForSelection), true) ||
			nextPhaseSub == null.NewInt(int64(entity.JobSeekerConfirmsScheduleAndSelectionDetail), true) ||
			nextPhaseSub == null.NewInt(int64(entity.ConfirmDayBefore), true) ||
			nextPhaseSub == null.NewInt(int64(entity.CollectSelectionThought), true) ||
			nextPhaseSub == null.NewInt(int64(entity.SendSelectionThoughtQuestionnaireText), true))

	if IsSendMessage {
		var message = taskParam.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(batchParam.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return err
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
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
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
					Email: message.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				batchParam.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	// 次が選考所管回収の場合はアンケートのレコードを作成
	if isPhaseSelection(nextPhase) && nextPhaseSub == null.NewInt(int64(entity.CollectSelectionThought), true) {
		// すでにアンケートのレコードが作成済みかチェック
		_, err := i.selectionQuestionnaireRepository.FindByJobSeekerIDAndJobInformationIDAndSelectionPhase(batchParam.JobSeekerID, batchParam.JobInformationID, uint(nextPhase.Int64))
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				selectionInformation, err := i.jobInfoSelectionInformationRepository.FindByJobSeekerIDAndJobInformationIDAndSelectionType(batchParam.JobSeekerID, batchParam.JobInformationID, nextPhase)
				if err != nil {
					fmt.Println(err)
					return err
				}

				// NotFoundの時はアンケート情報のレコード作成
				selectionQuestionnaire := entity.NewSelectionQuestionnaire(
					batchParam.JobSeekerID,
					batchParam.JobInformationID,
					selectionInformation.ID,
					NullInt,
					"",
					"",
					NullInt,
					"",
					"",
					"",
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					NullInt,
					"",
				)

				err = i.selectionQuestionnaireRepository.Create(selectionQuestionnaire)
				if err != nil {
					fmt.Println(err)
					return err
				}
			} else {
				fmt.Println(err)
				return err
			}
		}

	}

	return nil
}

// 内定保留選考フェーズのタスクを作成する処理
func createHoldJobOfferPhaseTask(i *TaskInteractorImpl, batchParam entity.CreateTaskInBatchProcessingParam, taskParam entity.TaskParam) error {
	var (
		err          error
		nextPhase    = taskParam.PhaseCategory
		nextPhaseSub = taskParam.PhaseSubCategory
		prevPhase    = taskParam.PrevPhaseCategory
	)

	// 選考フェーズでなかった場合はエラー
	if prevPhase != null.NewInt(int64(entity.HoldJobOffer), true) {
		return errors.New("現在フェーズが内定保留フェーズではありません")
	}

	/************ 新リスケ処理（リスケ + 削除） 元のリスケ処理も残しておく **************/

	if taskParam.TaskOption == null.NewInt(entity.OptionRescheduleAndDelete, true) {
		// 直前の結果回収依頼までの日程を削除 + リスケタスクの作成
		latestSelectTask, err := i.taskRepository.FindLatestByGroupIDAndCollectResultPhase(batchParam.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 直前の結果回収依頼までの日程を削除
		err = i.jobSeekerScheduleRepository.DeleteScheduleInTaskGroupAboveTaskID(latestSelectTask.ID, batchParam.TaskGroupID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// リスケタスクを作成
		rescheduleTask := entity.NewTask(
			batchParam.TaskGroupID,
			nextPhase,
			null.NewInt(entity.RescheduleDeleteOfSelection, true),
			taskParam.StaffType, // 担当者タイプ
			taskParam.ExecutedStaffID,
			"",
			batchParam.DeadlineDay,
			batchParam.DeadlineTime,
			"",
			"",
			"",
			false,
		)

		// タスク作成
		err = i.taskRepository.Create(rescheduleTask)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	task := entity.NewTask(
		batchParam.TaskGroupID,
		nextPhase,
		nextPhaseSub,
		taskParam.StaffType, // 担当者タイプ
		taskParam.ExecutedStaffID,
		taskParam.Remarks,
		batchParam.DeadlineDay,
		batchParam.DeadlineTime,
		taskParam.TalkAboutInInterview,
		taskParam.ScheduleCollectionCondition,
		taskParam.ExamGuideContent,
		false,
	)

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 内定承諾/内定承諾の場合はヨミ情報を更新する
	if nextPhase == null.NewInt(int64(entity.AcceptJobOffer), true) && nextPhaseSub == null.NewInt(int64(entity.Accept), true) {
		sale, err := i.saleRepository.FindByJobSeekerIDAndJobInformationID(batchParam.JobSeekerID, batchParam.JobInformationID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				fmt.Println("ヨミの登録がありません")
			} else {
				fmt.Println(err)
				return err
			}
		} else {
			sale.Accuracy = null.NewInt(entity.AccuracyAccept, true) // ヨミのみ「内定承諾: 0」に更新
			err = i.saleRepository.Update(sale.ID, sale)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	/************ 選考の候補日時を登録する **************/

	// 次のフェーズが「内定保留/ファー面談の日程調整依頼」かどうか
	var IsNextRequestScheduleAdjustmentForHold = nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) && nextPhaseSub == null.NewInt(int64(entity.RequestScheduleAdjustment), true)

	if IsNextRequestScheduleAdjustmentForHold {
		// 登録前にすでに登録済みの候補日の情報を削除する
		err = i.jobSeekerScheduleRepository.DeleteByTaskGroupIDAndScheduleType(batchParam.TaskGroupID, uint(entity.ScheduleTypeByCA))
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, possibleDate := range taskParam.PossibleDates {
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
				return err
			}
		}
	}

	/************ 選考の確定日時を登録する **************/

	// 次のフェーズが「内定保留/オファー面談の詳細案内依頼」かどうか
	var IsNextRequestDetailedInformationForHold = nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) && nextPhaseSub == null.NewInt(int64(entity.RequestDetailedInformationForHold), true)

	if IsNextRequestDetailedInformationForHold {
		var taskID null.Int

		if taskParam.SelectionDate.TaskID.Valid {
			// TaskIDがある場合は登録済みの確定日時なので、既存のTaskIDを使用する
			taskID = taskParam.SelectionDate.TaskID
		} else {
			// TaskIDが無い場合は未登録の確定日時なので、新しいTaskIDを使用する
			taskID = null.NewInt(int64(task.ID), true)
		}

		jobSeekerSchedule := entity.NewJobSeekerSchedule(
			taskParam.SelectionDate.JobSeekerID,
			taskID,
			taskParam.SelectionDate.ScheduleType,
			taskParam.SelectionDate.Title,
			taskParam.SelectionDate.StartTime,
			taskParam.SelectionDate.EndTime,
			taskParam.SelectionDate.SeekerDescription,
			taskParam.SelectionDate.StaffDescription,
			taskParam.SelectionDate.IsShare,
			taskParam.SelectionDate.RepetitionCount,
		)

		err = i.jobSeekerScheduleRepository.Create(jobSeekerSchedule)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	/************ メール or LINEの送付 **************/

	// 内定保留/オファー面談の候補日回収依頼 → オファー面談の日程回収中

	// メール or LINEの送信する
	var IsSendMessage = taskParam.TaskOption == null.NewInt(entity.OptionSendMessageForJobSeeker, true) &&
		nextPhase == null.NewInt(int64(entity.HoldJobOffer), true) && (nextPhaseSub == null.NewInt(int64(entity.CollectionOfSchedule), true) || nextPhaseSub == null.NewInt(int64(entity.ConfirmsSchedule), true))

	if IsSendMessage {
		var message = taskParam.Message

		// エージェント情報を取得
		staff, err := i.agentStaffRepository.FindStaffAndAgentLine(batchParam.CAStaffID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 求職者情報を取得
		jobSeeker, err := i.jobSeekerRepository.FindByID(batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// メッセージ送信の処理
		chatGroup, err := i.chatGroupWithJobSeekerRepository.FindByAgentIDAndJobSeekerID(staff.AgentID, batchParam.JobSeekerID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// LINE連携ありの場合
		if message.LineActive {
			// LINE連携済みでブロックされていない場合
			// シークレットとトークンを設定
			bot, err := linebot.New(staff.LineMessagingChannelSecret, staff.LineMessagingChannelAccessToken)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// LINE メッセージを送信する
			call, err := bot.PushMessage(
				jobSeeker.LineID,
				linebot.NewTextMessage(message.Content),
			).Do()
			if err != nil {
				fmt.Println(err)
				return err
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
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
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
					Email: message.Email,
				},
				nil,
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// メールを保存
			emailWithJobSeeker := entity.NewEmailWithJobSeeker(
				batchParam.JobSeekerID,
				message.Subject,
				message.Content,
				"",
			)

			err = i.emailWithJobSeekerRepository.Create(emailWithJobSeeker)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終送信時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastSendAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// エージェントの最終閲覧時間を更新
			err = i.chatGroupWithJobSeekerRepository.UpdateAgentLastWatchedAt(chatGroup.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

// 内定承諾選考フェーズのタスクを作成する処理
func createAcceptJobOfferPhaseTask(i *TaskInteractorImpl, batchParam entity.CreateTaskInBatchProcessingParam, taskParam entity.TaskParam) error {
	var (
		err error
	)

	// 選考フェーズでなかった場合はエラー
	if taskParam.PrevPhaseCategory != null.NewInt(int64(entity.AcceptJobOffer), true) {
		return errors.New("現在フェーズが内定承諾フェーズではありません")
	}

	task := entity.NewTask(
		batchParam.TaskGroupID,
		taskParam.PhaseCategory,
		taskParam.PhaseSubCategory,
		taskParam.StaffType, // 担当者タイプ
		taskParam.ExecutedStaffID,
		taskParam.Remarks,
		batchParam.DeadlineDay,
		batchParam.DeadlineTime,
		taskParam.TalkAboutInInterview,
		taskParam.ScheduleCollectionCondition,
		taskParam.ExamGuideContent,
		false,
	)

	// タスク作成
	err = i.taskRepository.Create(task)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 内定承諾/入社確認
	if taskParam.PhaseCategory == null.NewInt(int64(entity.AcceptJobOffer), true) && taskParam.PhaseSubCategory == null.NewInt(int64(entity.ConfirmJoiningCompany), true) {

		// 入社日の登録
		err = i.taskGroupRepository.UpdateJoiningDate(batchParam.TaskGroupID, taskParam.JoiningDate)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
