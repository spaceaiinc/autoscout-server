package interactor

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type SendingPhaseInteractor interface {
	// 汎用系 API
	UpdateSendingPhaseSendingDate(input UpdateSendingPhaseSendingDateInput) (UpdateSendingPhaseSendingDateOutput, error)

	GetSendingPhaseByID(input GetSendingPhaseByIDInput) (GetSendingPhaseByIDOutput, error)                                                             // 指定IDの送客進捗情報を取得する関数
	GetSendingPhaseListBySendingJobSeekerID(input GetSendingPhaseListBySendingJobSeekerIDInput) (GetSendingPhaseListBySendingJobSeekerIDOutput, error) // 指定求職者IDの送客進捗情報を取得する関数

	//ページネーション
	GetSearchSendingPhaseListByPageAndTabAndAgentID(input GetSearchSendingPhaseListByPageAndTabAndAgentIDInput) (GetSearchSendingPhaseListByPageAndTabAndAgentIDOutput, error)

	// 送客応諾後の終了理由
	CreateSendingPhaseEndStatus(input CreateSendingPhaseEndStatusInput) (CreateSendingPhaseEndStatusOutput, error)
}

type SendingPhaseInteractorImpl struct {
	firebase                        usecase.Firebase
	sendgrid                        config.Sendgrid
	sendingPhaseRepository          usecase.SendingPhaseRepository
	sendingJobSeekerRepository      usecase.SendingJobSeekerRepository
	sendingCustomerRepository       usecase.SendingCustomerRepository
	sendingPhaseEndStatusRepository usecase.SendingPhaseEndStatusRepository
}

// SendingPhaseInteractorImpl is an implementation of SendingPhaseInteractor
func NewSendingPhaseInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	spR usecase.SendingPhaseRepository,
	sjR usecase.SendingJobSeekerRepository,
	scR usecase.SendingCustomerRepository,
	spesR usecase.SendingPhaseEndStatusRepository,
) SendingPhaseInteractor {
	return &SendingPhaseInteractorImpl{
		firebase:                        fb,
		sendgrid:                        sg,
		sendingPhaseRepository:          spR,
		sendingJobSeekerRepository:      sjR,
		sendingCustomerRepository:       scR,
		sendingPhaseEndStatusRepository: spesR,
	}
}

/****************************************************************************************/
/// 汎用系API
//

// 送客進捗の送客予定日時を更新
type UpdateSendingPhaseSendingDateInput struct {
	SendingPhaseID uint
	SendingDate    time.Time
}

type UpdateSendingPhaseSendingDateOutput struct {
	OK bool
}

func (i *SendingPhaseInteractorImpl) UpdateSendingPhaseSendingDate(input UpdateSendingPhaseSendingDateInput) (UpdateSendingPhaseSendingDateOutput, error) {
	var (
		output UpdateSendingPhaseSendingDateOutput
		err    error
	)

	err = i.sendingPhaseRepository.UpdateSendingDate(input.SendingPhaseID, input.SendingDate)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 送客進捗IDを使って送客進捗情報を取得する
type GetSendingPhaseByIDInput struct {
	SendingPhaseID uint
}

type GetSendingPhaseByIDOutput struct {
	SendingPhase *entity.SendingPhase
}

func (i *SendingPhaseInteractorImpl) GetSendingPhaseByID(input GetSendingPhaseByIDInput) (GetSendingPhaseByIDOutput, error) {
	var (
		output GetSendingPhaseByIDOutput
		err    error
	)

	sendingPhase, err := i.sendingPhaseRepository.FindByID(input.SendingPhaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingPhase = sendingPhase

	return output, nil
}

// 指定求職者IDの送客進捗情報を取得する
type GetSendingPhaseListBySendingJobSeekerIDInput struct {
	SendingJobSeekerID uint
}

type GetSendingPhaseListBySendingJobSeekerIDOutput struct {
	SendingPhaseList []*entity.SendingPhase
}

func (i *SendingPhaseInteractorImpl) GetSendingPhaseListBySendingJobSeekerID(input GetSendingPhaseListBySendingJobSeekerIDInput) (GetSendingPhaseListBySendingJobSeekerIDOutput, error) {
	var (
		output GetSendingPhaseListBySendingJobSeekerIDOutput
		err    error
	)

	sendingPhaseList, err := i.sendingPhaseRepository.GetListBySendingJobSeekerID(input.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingPhaseList = sendingPhaseList

	return output, nil
}

/****************************************************************************************/
/// ページネーション
//

// 送客の進捗一覧ページに表示するリストを取得するapi（query: page_number, tab_number, staff_id_list, sender_id_list, send_agent_id_list）
type GetSearchSendingPhaseListByPageAndTabAndAgentIDInput struct {
	AgentID         uint
	PageNumber      uint
	TabNumber       uint
	FreeWord        string
	StaffIDList     []uint
	SenderIDList    []uint
	SendAgentIDList []uint
}

type GetSearchSendingPhaseListByPageAndTabAndAgentIDOutput struct {
	SendingJobSeekerTableList    []*entity.SendingJobSeekerTable
	MaxPageNumber                uint
	CompleteSendingCount         uint // 送客完了数
	AcceptSendingCount           uint // 送客応諾数
	PreparingAfterInterviewCount uint // 面談実施済み数
	WaitingForInterviewCount     uint // 面談実施待ち数
	UnregisterDetailCount        uint // 詳細未登録数
	UnregisterScheduleCount      uint // 日程未登録数
	CloseSendingCount            uint // 終了数
}

func (i *SendingPhaseInteractorImpl) GetSearchSendingPhaseListByPageAndTabAndAgentID(input GetSearchSendingPhaseListByPageAndTabAndAgentIDInput) (GetSearchSendingPhaseListByPageAndTabAndAgentIDOutput, error) {
	var (
		output               GetSearchSendingPhaseListByPageAndTabAndAgentIDOutput
		err                  error
		sendingPhaseList     []*entity.SendingPhase
		sendingJobSeekerList []*entity.SendingJobSeeker

		sentList     []*entity.SendingJobSeekerTable // 送客済みのリスト
		unSendList   []*entity.SendingJobSeekerTable // 未送客のリスト
		responseList []*entity.SendingJobSeekerTable // レスポンス
	)

	/********* 1. 送客済みのリスト *********/

	// 送客済みのリストを取得（「3: 送客応諾」と「4: 送客完了」と「5: 送客なし/終了」）
	sendingPhaseList, err = i.sendingPhaseRepository.GetSearchListForSendingJobSeekerManagement(input.AgentID, input.FreeWord, input.StaffIDList, input.SenderIDList, input.SendAgentIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sendingPhase := range sendingPhaseList {
		tableValue := entity.NewSendingJobSeekerTable(
			sendingPhase.ID,
			sendingPhase.SendingJobSeekerID,
			sendingPhase.LastName+sendingPhase.FirstName,
			sendingPhase.AgentStaffID,
			sendingPhase.StaffName,
			sendingPhase.AgentName,
			sendingPhase.SenderAgentName,
			sendingPhase.SendingDate,
			sendingPhase.InterviewDate,
			false, // 送客応諾の時は固定でFalse
			false, // 送客応諾の時は固定でFalse
		)

		if sendingPhase.Phase == null.NewInt(int64(entity.AcceptSending), true) {
			output.AcceptSendingCount++ // 送客応諾のカウント

			// 送客応諾の場合はリストにセット
			if input.TabNumber == uint(entity.AcceptSendingTab) {
				sentList = append(sentList, tableValue)
			}
		}

		if sendingPhase.Phase == null.NewInt(int64(entity.CompleteSending), true) {
			output.CompleteSendingCount++ // 送客完了のカウント

			// 送客応諾の場合はリストにセット
			if input.TabNumber == uint(entity.CompleteSendingTab) {
				sentList = append(sentList, tableValue)
			}
		}

		if sendingPhase.Phase == null.NewInt(int64(entity.CloseSending), true) {
			output.CloseSendingCount++ // 送客なし/終了のカウント

			// 送客応諾の場合はリストにセット
			if input.TabNumber == uint(entity.CloseSendingTab) {
				sentList = append(sentList, tableValue)
			}
		}
	}

	/********* 2. 未送客のリスト *********/

	// 送客先の絞り込み条件がない場合のみ送客応諾以外のリストを
	if len(input.SendAgentIDList) == 0 {
		sendingJobSeekerList, err = i.sendingJobSeekerRepository.GetSearchListForSendingJobSeekerManagement(input.AgentID, input.FreeWord, input.StaffIDList, input.SenderIDList)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, sendingJobSeeker := range sendingJobSeekerList {
		tableValue := entity.NewSendingJobSeekerTable(
			sendingJobSeeker.ID,
			sendingJobSeeker.ID,
			sendingJobSeeker.LastName+sendingJobSeeker.FirstName,
			sendingJobSeeker.AgentStaffID,
			sendingJobSeeker.StaffName,
			"", // 未送客のため「送客先名」はブランク
			sendingJobSeeker.SenderAgentName,
			utility.EarliestTime(), // 未送客のため「送客日」はminTimeを渡す
			sendingJobSeeker.InterviewDate,
			sendingJobSeeker.IsNotWaitingViewed,
			sendingJobSeeker.IsNotUnregisterViewed,
		)

		// 面談実施済み数をカウント
		if sendingJobSeeker.Phase == null.NewInt(int64(entity.WaitingForInterview), true) && sendingJobSeeker.InterviewDate.Before(time.Now()) {
			output.PreparingAfterInterviewCount++

			if input.TabNumber == uint(entity.PreparingAfterInterviewTab) {
				unSendList = append(unSendList, tableValue)
			}

			// 面談実施待ち数をカウント
		} else if sendingJobSeeker.Phase == null.NewInt(int64(entity.WaitingForInterview), true) {
			output.WaitingForInterviewCount++

			if input.TabNumber == uint(entity.WaitingForInterviewTab) {
				unSendList = append(unSendList, tableValue)
			}

			// 詳細未登録数をカウント
		} else if sendingJobSeeker.Phase == null.NewInt(int64(entity.UnregisterDetail), true) {
			output.UnregisterDetailCount++

			if input.TabNumber == uint(entity.UnregisterTab) {
				unSendList = append(unSendList, tableValue)
			}

			// 日程未登録数をカウント
		} else if sendingJobSeeker.Phase == null.NewInt(int64(entity.UnregisterSchedule), true) {
			output.UnregisterScheduleCount++

			if input.TabNumber == uint(entity.UnregisterTab) {
				unSendList = append(unSendList, tableValue)
			}

			// 終了数をカウント
		} else if sendingJobSeeker.Phase == null.NewInt(int64(entity.CloseSending), true) {
			// 終了のカウントは送客済みと未送客で重複がないようにカウントする
			found := false

			for _, sendingPhase := range sendingPhaseList {
				if sendingPhase.Phase == null.NewInt(int64(entity.CloseSending), true) && sendingJobSeeker.ID == sendingPhase.SendingJobSeekerID {
					found = true
					break
				}
			}
			if !found {
				output.CloseSendingCount++
			}

			if input.TabNumber == uint(entity.CloseSendingTab) {
				unSendList = append(unSendList, tableValue)
			}
		}
	}

	/********* タブに応じて返すリストをセット *********/

	if input.TabNumber == uint(entity.CompleteSendingTab) || input.TabNumber == uint(entity.AcceptSendingTab) {
		// 送客完了, 送客応諾
		responseList = sentList
	} else if input.TabNumber == uint(entity.PreparingAfterInterviewTab) || input.TabNumber == uint(entity.WaitingForInterviewTab) || input.TabNumber == uint(entity.UnregisterTab) {
		// 面談実施済み, 面談実施待ち, 詳細未登録 / 日程未登録
		responseList = unSendList
	} else {
		// 終了 「送客済みの終了」と「未送客の終了」を重複がないように合わせる。
		responseList = sentList

		for _, unSend := range unSendList {
			found := false
			for _, responseItem := range responseList {
				if unSend.SendingJobSeekerID == responseItem.SendingJobSeekerID {
					found = true
					break
				}
			}

			if !found {
				// 未送客で終了した求職者をセット
				// output.CloseSendingCount++
				responseList = append(responseList, unSend)
			}
		}
	}

	/********* レスポンスの処理 *********/

	// ページの最大数を取得
	output.MaxPageNumber = getSendingManagementListMaxPage(responseList)

	// 指定ページの送客求職者20件を取得（本番実装までは1ページあたり20件）
	returnList20 := getSendingManagementListWithPage(responseList, input.PageNumber)

	output.SendingJobSeekerTableList = returnList20

	return output, nil
}

/****************************************************************************************/
/// 送客応諾後の終了理由
//
// 送客進捗の作成
type CreateSendingPhaseEndStatusInput struct {
	CreateParam entity.CreateSendingPhaseEndStatusParam
}

type CreateSendingPhaseEndStatusOutput struct {
	OK bool
}

func (i *SendingPhaseInteractorImpl) CreateSendingPhaseEndStatus(input CreateSendingPhaseEndStatusInput) (CreateSendingPhaseEndStatusOutput, error) {
	var (
		output CreateSendingPhaseEndStatusOutput
		err    error
	)

	/********* 送客進捗のフェーズを終了に変更 *********/

	err = i.sendingPhaseRepository.UpdatePhase(input.CreateParam.SendingPhaseID, uint(entity.CloseSending))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/********* 終了理由の登録 *********/

	sendingPhaseEndStatus := entity.NewSendingPhaseEndStatus(
		input.CreateParam.SendingPhaseID,
		input.CreateParam.EndReason,
		input.CreateParam.EndStatus,
	)

	err = i.sendingPhaseEndStatusRepository.Create(sendingPhaseEndStatus)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/********* 対象求職者のフェーズがすべてを終了かをチェック *********/

	// 全てのフェーズが終了した場合、求職者のフェーズを終了に変更
	sendingPhase, err := i.sendingPhaseRepository.FindByID(input.CreateParam.SendingPhaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	sendingPhaseList, err := i.sendingPhaseRepository.GetListBySendingJobSeekerID(sendingPhase.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	isAllPhaseEnd := true
	for _, sendingPhase := range sendingPhaseList {
		if sendingPhase.Phase != null.NewInt(int64(entity.CloseSending), true) {
			isAllPhaseEnd = false
			break
		}
	}

	/********* 送客がすべて終了になっている場合、求職者自体のフェーズも終了にする *********/

	if isAllPhaseEnd {
		// 管理側のフェーズを終了に更新
		err = i.sendingJobSeekerRepository.UpdatePhase(sendingPhase.SendingJobSeekerID, null.NewInt(int64(entity.CloseSending), true))
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		// 送客元のフェーズを終了に更新
		err = i.sendingCustomerRepository.UpdatePhaseBySendingJobSeekerID(sendingPhase.SendingJobSeekerID, null.NewInt(int64(entity.CloseSending), true))
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.OK = true

	return output, nil
}
