package interactor

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type SelectionQuestionnaireInteractor interface {
	// 汎用系 API
	CreateSelectionQuestionnaire(input CreateSelectionQuestionnaireInput) (CreateSelectionQuestionnaireOutput, error)
	UpdateSelectionQuestionnaire(input UpdateSelectionQuestionnaireInput) (UpdateSelectionQuestionnaireOutput, error)
	GetSelectionQuestionnaireOrNullByUUID(input GetSelectionQuestionnaireOrNullByUUIDInput) (GetSelectionQuestionnaireOrNullByUUIDOutput, error)
	GenerateSelectionQuestionnaireByUUID() (GenerateSelectionQuestionnaireByUUIDOutput, error)
	GetUnansweredQuestionnaireListByJobSeekerUUID(input GetUnansweredQuestionnaireListByJobSeekerUUIDInput) (GetUnansweredQuestionnaireListByJobSeekerUUIDOutput, error)
	GetQuestionnaireForJobSeeker(input GetQuestionnaireForJobSeekerInput) (GetQuestionnaireForJobSeekerOutput, error)
}

type SelectionQuestionnaireInteractorImpl struct {
	firebase                                  usecase.Firebase
	sendgrid                                  config.Sendgrid
	selectionQuestionnaireRepository          usecase.SelectionQuestionnaireRepository
	selectionQuestionnaireMyRankingRepository usecase.SelectionQuestionnaireMyRankingRepository
	jobSeekerRepository                       usecase.JobSeekerRepository
	jobInformationRepository                  usecase.JobInformationRepository
	agentStaffRepository                      usecase.AgentStaffRepository
}

// SelectionQuestionnaireInteractorImpl is an implementation of SelectionQuestionnaireInteractor
func NewSelectionQuestionnaireInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	sqR usecase.SelectionQuestionnaireRepository,
	smR usecase.SelectionQuestionnaireMyRankingRepository,
	jsR usecase.JobSeekerRepository,
	jIR usecase.JobInformationRepository,
	asR usecase.AgentStaffRepository,
) SelectionQuestionnaireInteractor {
	return &SelectionQuestionnaireInteractorImpl{
		firebase:                         fb,
		sendgrid:                         sg,
		selectionQuestionnaireRepository: sqR,
		selectionQuestionnaireMyRankingRepository: smR,
		jobSeekerRepository:                       jsR,
		jobInformationRepository:                  jIR,
		agentStaffRepository:                      asR,
	}
}

/****************************************************************************************/
/// 汎用系API
//

// 選考後アンケートの更新

// 選考後アンケートの更新
type CreateSelectionQuestionnaireInput struct {
	CreateParam       entity.CreateOrUpdateSelectionQuestionnaireParam
	QuestionnaireUUID uuid.UUID
}

type CreateSelectionQuestionnaireOutput struct {
	SelectionQuestionnaire *entity.SelectionQuestionnaire
}

func (i *SelectionQuestionnaireInteractorImpl) CreateSelectionQuestionnaire(input CreateSelectionQuestionnaireInput) (CreateSelectionQuestionnaireOutput, error) {
	var (
		output CreateSelectionQuestionnaireOutput
		err    error
	)

	/**********************************************/
	// データ更新
	//

	// 選考後アンケートを更新
	selectionQuestionnaire := entity.NewSelectionQuestionnaire(
		input.CreateParam.JobSeekerID,
		input.CreateParam.JobInformationID,
		input.CreateParam.SelectionInformationID,
		input.CreateParam.MyRanking,
		input.CreateParam.MyRankingReason,
		input.CreateParam.ConcernPoint,
		input.CreateParam.ContinueSelection,
		input.CreateParam.MyRankingDetail,
		input.CreateParam.SelectionQuestion,
		input.CreateParam.Remarks,
		input.CreateParam.IsAnswer,
		input.CreateParam.IsSelfIntroduction,
		input.CreateParam.IsSelfPR,
		input.CreateParam.IsRetireReason,
		input.CreateParam.IsJobChangeAxis,
		input.CreateParam.IsApplyingReason,
		input.CreateParam.IsCareerVision,
		input.CreateParam.IsReverseQuestion,
		input.CreateParam.IntentionToJobOffer,
		input.CreateParam.IntentionDetail,
	)

	err = i.selectionQuestionnaireRepository.CreateByUUID(selectionQuestionnaire, input.QuestionnaireUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SelectionQuestionnaire = selectionQuestionnaire

	// 削除
	err = i.selectionQuestionnaireMyRankingRepository.DeleteByQuestionnaireID(selectionQuestionnaire.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, myRank := range input.CreateParam.SelectionQuestionnaireMyRanking {
		// 新規作成
		myRanking := entity.NewSelectionQuestionnaireMyRanking(
			selectionQuestionnaire.ID,
			myRank.Rank,
			myRank.CompanyName,
			myRank.Phase,
			myRank.SelectionDate,
		)

		err = i.selectionQuestionnaireMyRankingRepository.Create(myRanking)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		output.SelectionQuestionnaire.SelectionQuestionnaireMyRanking = append(
			output.SelectionQuestionnaire.SelectionQuestionnaireMyRanking,
			*myRanking,
		)
	}

	/**********************************************/
	// メール送信に必要な情報を取得
	//

	// 求職者情報を取得
	jobSeeker, err := i.jobSeekerRepository.FindByID(input.CreateParam.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求人情報を取得
	jobInformation, err := i.jobInformationRepository.FindByID(input.CreateParam.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// CA担当者情報を取得
	caStaff, err := i.agentStaffRepository.FindByJobSeekerID(input.CreateParam.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 通知がオフの場合はメール送信しない
	if !caStaff.NotificationJobSeeker {
		return output, nil
	}

	/**********************************************/
	// アンケートの回答内容をテキストに変換
	//

	var (
		selectionType     = input.CreateParam.SelectionType
		selectionLabel    string
		questionnaireText string
		sendText          string
		firstText         string
		secondText        string
		thirdText         string
		fourthText        string
		fifthText         string
		sixthText         string
		seventhText       string
		eighthText        string
		rankingTextList   []string
		rankingText       string
		isCheckTextList   []string
		isCheckText       string
	)

	if selectionType == null.NewInt(int64(entity.FirstSelection), true) {
		selectionLabel = "1次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.SecondSelection), true) {
		selectionLabel = "2次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.ThirdSelection), true) {
		selectionLabel = "3次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.FourthSelection), true) {
		selectionLabel = "4次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.FifthSelection), true) {
		selectionLabel = "5次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.FinalSelection), true) {
		selectionLabel = "最終考後アンケート"
	}

	firstText = fmt.Sprintf(
		"① 志望度\n%s",
		getStrMyRanking(selectionQuestionnaire.MyRanking),
	)

	secondText = fmt.Sprintf(
		"② 志望度が上がった点\n%s",
		selectionQuestionnaire.MyRankingReason,
	)

	thirdText = fmt.Sprintf(
		"③ 心配な点・懸念点\n%s",
		selectionQuestionnaire.ConcernPoint,
	)

	if selectionType == null.NewInt(int64(entity.FinalSelection), true) {
		fourthText = fmt.Sprintf(
			"④ 内定が出た場合の意向についてお聞かせください\n%s\n%s",
			getStrIntentionToJobOffer(selectionQuestionnaire.IntentionToJobOffer),
			selectionQuestionnaire.IntentionDetail,
		)
	} else {
		fourthText = fmt.Sprintf(
			"④ 選考継続希望有無についてお聞かせください\n%s",
			getStrContinueSelectionType(selectionQuestionnaire.ContinueSelection),
		)
	}

	if len(selectionQuestionnaire.SelectionQuestionnaireMyRanking) > 0 {
		for _, ranking := range selectionQuestionnaire.SelectionQuestionnaireMyRanking {
			timeText := ""

			parse, err := time.Parse("2006-01-02", ranking.SelectionDate)
			if err == nil {
				layout := "2006年1月2日"
				timeText = parse.Format(layout)
			} else {
				timeText = "日程調整中"
			}

			t := fmt.Sprintf(
				"志望度: %s\n社名: %s\nフェーズ: %s\n日付: %s",
				getStrOtherMyRanking(ranking.Rank),
				ranking.CompanyName,
				getStrSelectionPhase(ranking.Phase),
				timeText,
			)

			rankingTextList = append(rankingTextList, t)
		}

		if len(rankingTextList) > 0 {
			rankingText = strings.Trim(strings.Join(rankingTextList, "\n\n"), "[]")
		}

		fifthText = fmt.Sprintf(
			"⑤ 他社選考について志望順位を踏まえて教えてください\n%s",
			rankingText,
		)
	} else {
		fifthText = fmt.Sprintf(
			"⑤ 他社選考について志望順位を踏まえて教えてください\n%s",
			"他社選考なし",
		)
	}

	if selectionQuestionnaire.IsSelfIntroduction {
		isCheckTextList = append(isCheckTextList, "自己紹介")
	}
	if selectionQuestionnaire.IsSelfPR {
		isCheckTextList = append(isCheckTextList, "自己PR（強み・長所）")
	}
	if selectionQuestionnaire.IsRetireReason {
		isCheckTextList = append(isCheckTextList, "退職理由")
	}
	if selectionQuestionnaire.IsJobChangeAxis {
		isCheckTextList = append(isCheckTextList, "転職軸")
	}
	if selectionQuestionnaire.IsApplyingReason {
		isCheckTextList = append(isCheckTextList, "志望動機(志望理由)")
	}
	if selectionQuestionnaire.IsCareerVision {
		isCheckTextList = append(isCheckTextList, "キャリアビジョン")
	}
	if selectionQuestionnaire.IsReverseQuestion {
		isCheckTextList = append(isCheckTextList, "逆質問")
	}
	if len(isCheckTextList) > 0 {
		isCheckText = strings.Trim(strings.Join(isCheckTextList, ", "), "[]")
	}

	sixthText = fmt.Sprintf(
		"⑥ どのようなことを聞かれましたか\n%s",
		isCheckText,
	)

	seventhText = fmt.Sprintf(
		"⑦ ⑥で回答した質問以外ではどのようなことを聞かれましたか\n%s",
		selectionQuestionnaire.SelectionQuestion,
	)

	eighthText = fmt.Sprintf(
		"⑧ その他何か質問があればご記入ください\n%s",
		selectionQuestionnaire.Remarks,
	)

	var textList []string
	textList = append(textList, firstText, secondText, thirdText, fourthText, fifthText, sixthText, seventhText, eighthText)
	questionnaireText = strings.Trim(strings.Join(textList, "\n\n"), "[]")

	sendText = fmt.Sprintf(
		"■%s\n\n%s",
		selectionLabel,
		questionnaireText,
	)

	messageText := fmt.Sprintf(
		"求職者がアンケートに回答しました。\n\n求職者: %s\n\n企業名: %s\n\n求人: %s\n\n\n%s",
		jobSeeker.LastName,
		jobInformation.CompanyName, // 企業名
		jobInformation.Title,       // 求人タイトル
		sendText,
	)

	/**********************************************/
	// 担当CAにメール送信
	//

	err = utility.SendMailToSingleByMyself(
		i.sendgrid.APIKey,
		"アンケートの回答",
		messageText,
		entity.EmailUser{
			Name:  caStaff.StaffName,
			Email: caStaff.Email,
		}, // from and to
	)

	if err != nil {
		fmt.Println(err)
		return output, err
	}

	return output, nil
}

// 選考後アンケートの更新
type UpdateSelectionQuestionnaireInput struct {
	UpdateParam     entity.CreateOrUpdateSelectionQuestionnaireParam
	QuestionnaireID uint
}

type UpdateSelectionQuestionnaireOutput struct {
	SelectionQuestionnaire *entity.SelectionQuestionnaire
}

func (i *SelectionQuestionnaireInteractorImpl) UpdateSelectionQuestionnaire(input UpdateSelectionQuestionnaireInput) (UpdateSelectionQuestionnaireOutput, error) {
	var (
		output UpdateSelectionQuestionnaireOutput
		err    error
	)

	/**********************************************/
	// データ更新
	//

	// 選考後アンケートを更新
	selectionQuestionnaire := entity.NewSelectionQuestionnaire(
		input.UpdateParam.JobSeekerID,
		input.UpdateParam.JobInformationID,
		input.UpdateParam.SelectionInformationID,
		input.UpdateParam.MyRanking,
		input.UpdateParam.MyRankingReason,
		input.UpdateParam.ConcernPoint,
		input.UpdateParam.ContinueSelection,
		input.UpdateParam.MyRankingDetail,
		input.UpdateParam.SelectionQuestion,
		input.UpdateParam.Remarks,
		input.UpdateParam.IsAnswer,
		input.UpdateParam.IsSelfIntroduction,
		input.UpdateParam.IsSelfPR,
		input.UpdateParam.IsRetireReason,
		input.UpdateParam.IsJobChangeAxis,
		input.UpdateParam.IsApplyingReason,
		input.UpdateParam.IsCareerVision,
		input.UpdateParam.IsReverseQuestion,
		input.UpdateParam.IntentionToJobOffer,
		input.UpdateParam.IntentionDetail,
	)

	err = i.selectionQuestionnaireRepository.Update(input.QuestionnaireID, selectionQuestionnaire)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SelectionQuestionnaire = selectionQuestionnaire

	// 削除
	err = i.selectionQuestionnaireMyRankingRepository.DeleteByQuestionnaireID(input.QuestionnaireID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, myRank := range input.UpdateParam.SelectionQuestionnaireMyRanking {
		// 新規作成
		myRanking := entity.NewSelectionQuestionnaireMyRanking(
			input.QuestionnaireID,
			myRank.Rank,
			myRank.CompanyName,
			myRank.Phase,
			myRank.SelectionDate,
		)

		err = i.selectionQuestionnaireMyRankingRepository.Create(myRanking)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		output.SelectionQuestionnaire.SelectionQuestionnaireMyRanking = append(
			output.SelectionQuestionnaire.SelectionQuestionnaireMyRanking,
			*myRanking,
		)
	}

	/**********************************************/
	// メール送信に必要な情報を取得
	//

	// 求職者情報を取得
	jobSeeker, err := i.jobSeekerRepository.FindByID(input.UpdateParam.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 求人情報を取得
	jobInformation, err := i.jobInformationRepository.FindByID(input.UpdateParam.JobInformationID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// CA担当者情報を取得
	caStaff, err := i.agentStaffRepository.FindByJobSeekerID(input.UpdateParam.JobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 通知がオフの場合はメール送信しない
	if !caStaff.NotificationJobSeeker {
		return output, nil
	}

	/**********************************************/
	// アンケートの回答内容をテキストに変換
	//

	var (
		selectionType     = input.UpdateParam.SelectionType
		selectionLabel    string
		questionnaireText string
		sendText          string
		firstText         string
		secondText        string
		thirdText         string
		fourthText        string
		fifthText         string
		sixthText         string
		seventhText       string
		eighthText        string
		rankingTextList   []string
		rankingText       string
		isCheckTextList   []string
		isCheckText       string
	)

	if selectionType == null.NewInt(int64(entity.FirstSelection), true) {
		selectionLabel = "1次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.SecondSelection), true) {
		selectionLabel = "2次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.ThirdSelection), true) {
		selectionLabel = "3次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.FourthSelection), true) {
		selectionLabel = "4次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.FifthSelection), true) {
		selectionLabel = "5次考後アンケート"
	} else if selectionType == null.NewInt(int64(entity.FinalSelection), true) {
		selectionLabel = "最終考後アンケート"
	}

	firstText = fmt.Sprintf(
		"① 志望度\n%s",
		getStrMyRanking(selectionQuestionnaire.MyRanking),
	)

	secondText = fmt.Sprintf(
		"② 志望度が上がった点\n%s",
		selectionQuestionnaire.MyRankingReason,
	)

	thirdText = fmt.Sprintf(
		"③ 心配な点・懸念点\n%s",
		selectionQuestionnaire.ConcernPoint,
	)

	if selectionType == null.NewInt(int64(entity.FinalSelection), true) {
		fourthText = fmt.Sprintf(
			"④ 内定が出た場合の意向についてお聞かせください\n%s\n%s",
			getStrIntentionToJobOffer(selectionQuestionnaire.IntentionToJobOffer),
			selectionQuestionnaire.IntentionDetail,
		)
	} else {
		fourthText = fmt.Sprintf(
			"④ 選考継続希望有無についてお聞かせください\n%s",
			getStrContinueSelectionType(selectionQuestionnaire.ContinueSelection),
		)
	}

	if len(selectionQuestionnaire.SelectionQuestionnaireMyRanking) > 0 {
		for _, ranking := range selectionQuestionnaire.SelectionQuestionnaireMyRanking {
			timeText := ""

			parse, err := time.Parse("2006-01-02", ranking.SelectionDate)
			if err == nil {
				layout := "2006年1月2日"
				timeText = parse.Format(layout)
			} else {
				timeText = "日程調整中"
			}

			t := fmt.Sprintf(
				"志望度: %s\n社名: %s\nフェーズ: %s\n日付: %s",
				getStrOtherMyRanking(ranking.Rank),
				ranking.CompanyName,
				getStrSelectionPhase(ranking.Phase),
				timeText,
			)

			rankingTextList = append(rankingTextList, t)
		}

		if len(rankingTextList) > 0 {
			rankingText = strings.Trim(strings.Join(rankingTextList, "\n\n"), "[]")
		}

		fifthText = fmt.Sprintf(
			"⑤ 他社選考について志望順位を踏まえて教えてください\n%s",
			rankingText,
		)
	} else {
		fifthText = fmt.Sprintf(
			"⑤ 他社選考について志望順位を踏まえて教えてください\n%s",
			"他社選考なし",
		)
	}

	if selectionQuestionnaire.IsSelfIntroduction {
		isCheckTextList = append(isCheckTextList, "自己紹介")
	}
	if selectionQuestionnaire.IsSelfPR {
		isCheckTextList = append(isCheckTextList, "自己PR（強み・長所）")
	}
	if selectionQuestionnaire.IsRetireReason {
		isCheckTextList = append(isCheckTextList, "退職理由")
	}
	if selectionQuestionnaire.IsJobChangeAxis {
		isCheckTextList = append(isCheckTextList, "転職軸")
	}
	if selectionQuestionnaire.IsApplyingReason {
		isCheckTextList = append(isCheckTextList, "志望動機(志望理由)")
	}
	if selectionQuestionnaire.IsCareerVision {
		isCheckTextList = append(isCheckTextList, "キャリアビジョン")
	}
	if selectionQuestionnaire.IsReverseQuestion {
		isCheckTextList = append(isCheckTextList, "逆質問")
	}
	if len(isCheckTextList) > 0 {
		isCheckText = strings.Trim(strings.Join(isCheckTextList, ", "), "[]")
	}

	sixthText = fmt.Sprintf(
		"⑥ どのようなことを聞かれましたか\n%s",
		isCheckText,
	)

	seventhText = fmt.Sprintf(
		"⑦ ⑥で回答した質問以外ではどのようなことを聞かれましたか\n%s",
		selectionQuestionnaire.SelectionQuestion,
	)

	eighthText = fmt.Sprintf(
		"⑧ その他何か質問があればご記入ください\n%s",
		selectionQuestionnaire.Remarks,
	)

	var textList []string
	textList = append(textList, firstText, secondText, thirdText, fourthText, fifthText, sixthText, seventhText, eighthText)
	questionnaireText = strings.Trim(strings.Join(textList, "\n\n"), "[]")

	sendText = fmt.Sprintf(
		"■%s\n\n%s",
		selectionLabel,
		questionnaireText,
	)

	messageText := fmt.Sprintf(
		"求職者がアンケートに回答しました。\n\n求職者: %s\n\n企業名: %s\n\n求人: %s\n\n\n%s",
		jobSeeker.LastName,
		jobInformation.CompanyName, // 企業名
		jobInformation.Title,       // 求人タイトル
		sendText,
	)

	/**********************************************/
	// 担当CAにメール送信
	//

	err = utility.SendMailToSingleByMyself(
		i.sendgrid.APIKey,
		"アンケートの回答",
		messageText,
		entity.EmailUser{
			Name:  caStaff.StaffName,
			Email: caStaff.Email,
		}, // from and to
	)

	if err != nil {
		fmt.Println(err)
		return output, err
	}

	return output, nil
}

type GetSelectionQuestionnaireOrNullByUUIDInput struct {
	QuestionnaireUUID uuid.UUID
}

type GetSelectionQuestionnaireOrNullByUUIDOutput struct {
	SelectionQuestionnaire *entity.SelectionQuestionnaire
}

func (i *SelectionQuestionnaireInteractorImpl) GetSelectionQuestionnaireOrNullByUUID(input GetSelectionQuestionnaireOrNullByUUIDInput) (GetSelectionQuestionnaireOrNullByUUIDOutput, error) {
	var (
		output GetSelectionQuestionnaireOrNullByUUIDOutput
		err    error

		rankingList []entity.SelectionQuestionnaireMyRanking
	)

	selectionQuestionnaire, err := i.selectionQuestionnaireRepository.FindByUUID(input.QuestionnaireUUID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// NotFoundの時はエラーにはしない
			return output, nil
		} else {
			fmt.Println(err)
			return output, err
		}
	}

	rankings, err := i.selectionQuestionnaireMyRankingRepository.GetByQuestionnaireID(selectionQuestionnaire.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ranking := range rankings {
		rankingList = append(rankingList, *ranking)
	}

	selectionQuestionnaire.SelectionQuestionnaireMyRanking = rankingList
	output.SelectionQuestionnaire = selectionQuestionnaire

	return output, nil
}

type GenerateSelectionQuestionnaireByUUIDOutput struct {
	QuestionnaireUUID uuid.UUID
}

func (i *SelectionQuestionnaireInteractorImpl) GenerateSelectionQuestionnaireByUUID() (GenerateSelectionQuestionnaireByUUIDOutput, error) {
	var (
		output            GenerateSelectionQuestionnaireByUUIDOutput
		questionnaireUUID uuid.UUID
	)
	var selectionQuestionnaireRepository = i.selectionQuestionnaireRepository

	// uuidを生成する
	// 既存で使用しているuuidと重複ないかを50回チェックする（50回超えたらエラー）
	for i := 0; i < 50; i++ {
		questionnaireUUID = utility.CreateUUID()

		_, err := selectionQuestionnaireRepository.FindByUUID(questionnaireUUID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				// NotFoundの場合はエラーにはせずに返す
				break
			} else {
				fmt.Println(err)
				return output, err
			}
		}
		continue
	}

	output.QuestionnaireUUID = questionnaireUUID
	return output, nil
}

type GetUnansweredQuestionnaireListByJobSeekerUUIDInput struct {
	JobSeekerUUID uuid.UUID
}

type GetUnansweredQuestionnaireListByJobSeekerUUIDOutput struct {
	SelectionQuestionnaireList []*entity.SelectionQuestionnaire
}

func (i *SelectionQuestionnaireInteractorImpl) GetUnansweredQuestionnaireListByJobSeekerUUID(input GetUnansweredQuestionnaireListByJobSeekerUUIDInput) (GetUnansweredQuestionnaireListByJobSeekerUUIDOutput, error) {
	var (
		output GetUnansweredQuestionnaireListByJobSeekerUUIDOutput
		err    error
	)

	selectionQuestionnaireList, err := i.selectionQuestionnaireRepository.GetUnanswerdByJobSeekerUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SelectionQuestionnaireList = selectionQuestionnaireList

	return output, nil
}

type GetQuestionnaireForJobSeekerInput struct {
	JobSeekerUUID      uuid.UUID
	JobInformationUUID uuid.UUID
	SelectionPhase     uint
}

type GetQuestionnaireForJobSeekerOutput struct {
	SelectionQuestionnaire *entity.SelectionQuestionnaire
}

func (i *SelectionQuestionnaireInteractorImpl) GetQuestionnaireForJobSeeker(input GetQuestionnaireForJobSeekerInput) (GetQuestionnaireForJobSeekerOutput, error) {
	var (
		output GetQuestionnaireForJobSeekerOutput
		err    error
	)

	selectionQuestionnaire, err := i.selectionQuestionnaireRepository.FindByJobSeekerUUIDAndJobInformationUUIDAndSelectionPhase(
		input.JobSeekerUUID, input.JobInformationUUID, input.SelectionPhase,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	rankings, err := i.selectionQuestionnaireMyRankingRepository.GetByQuestionnaireID(selectionQuestionnaire.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, ranking := range rankings {
		selectionQuestionnaire.SelectionQuestionnaireMyRanking = append(selectionQuestionnaire.SelectionQuestionnaireMyRanking, *ranking)
	}

	output.SelectionQuestionnaire = selectionQuestionnaire

	return output, nil
}

/****************************************************************************************/
