package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type SendingCustomerInteractor interface {
	// 汎用系 API
	GetSendingCustomerByID(input GetSendingCustomerByIDInput) (GetSendingCustomerByIDOutput, error)
	GetSendingCustomerListByAgentID(input GetSendingCustomerListByAgentIDInput) (GetSendingCustomerListByAgentIDOutput, error)

	// ページネーション
	GetSearchSendingCustomerListByPageAndTabAndAgentID(input GetSearchSendingCustomerListByPageAndTabAndAgentIDInput) (GetSearchSendingCustomerListByPageAndTabAndAgentIDOutput, error) // フェーズごとに取得
}

type SendingCustomerInteractorImpl struct {
	firebase                   usecase.Firebase
	sendgrid                   config.Sendgrid
	oneSignal                  config.OneSignal
	sendingCustomerRepository  usecase.SendingCustomerRepository
	sendingJobSeekerRepository usecase.SendingJobSeekerRepository
	agentRepository            usecase.AgentRepository
	agentStaffRepository       usecase.AgentStaffRepository
}

// SendingCustomerInteractorImpl is an implementation of SendingCustomerInteractor
func NewSendingCustomerInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	scR usecase.SendingCustomerRepository,
	sjsR usecase.SendingJobSeekerRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
) SendingCustomerInteractor {
	return &SendingCustomerInteractorImpl{
		firebase:                   fb,
		sendgrid:                   sg,
		oneSignal:                  os,
		sendingCustomerRepository:  scR,
		sendingJobSeekerRepository: sjsR,
		agentRepository:            aR,
		agentStaffRepository:       asR,
	}
}

/****************************************************************************************/
/// 汎用系
//
// 指定したIDの送客求職者を取得する
type GetSendingCustomerByIDInput struct {
	SendingCustomerID uint
}

type GetSendingCustomerByIDOutput struct {
	SendingCustomer *entity.SendingCustomer
}

func (i *SendingCustomerInteractorImpl) GetSendingCustomerByID(input GetSendingCustomerByIDInput) (GetSendingCustomerByIDOutput, error) {
	var (
		output GetSendingCustomerByIDOutput
		err    error
	)

	sendingCustomer, err := i.sendingCustomerRepository.FindByID(input.SendingCustomerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingCustomer = sendingCustomer

	return output, nil
}

type GetSendingCustomerListByAgentIDInput struct {
	AgentID uint
}

type GetSendingCustomerListByAgentIDOutput struct {
	SendingCustomerList []*entity.SendingCustomer
}

func (i *SendingCustomerInteractorImpl) GetSendingCustomerListByAgentID(input GetSendingCustomerListByAgentIDInput) (GetSendingCustomerListByAgentIDOutput, error) {
	var (
		output GetSendingCustomerListByAgentIDOutput
		err    error
	)

	sendingCustomerList, err := i.sendingCustomerRepository.GetListByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingCustomerList = sendingCustomerList

	return output, nil
}

type GetSearchSendingCustomerListByPageAndTabAndAgentIDInput struct {
	PageNumber uint
	AgentID    uint
	TabNumber  uint
	FreeWord   string
}

type GetSearchSendingCustomerListByPageAndTabAndAgentIDOutput struct {
	SendingCustomerList []*entity.SendingCustomer
	MaxPageNumber       uint
	AllCount            uint
	SendingCount        uint
	CompleteCount       uint
	CloseCount          uint
}

func (i *SendingCustomerInteractorImpl) GetSearchSendingCustomerListByPageAndTabAndAgentID(input GetSearchSendingCustomerListByPageAndTabAndAgentIDInput) (GetSearchSendingCustomerListByPageAndTabAndAgentIDOutput, error) {
	var (
		output          GetSearchSendingCustomerListByPageAndTabAndAgentIDOutput
		err             error
		AllList         []*entity.SendingCustomer
		SendingList     []*entity.SendingCustomer
		CompleteList    []*entity.SendingCustomer
		CloseList       []*entity.SendingCustomer
		SelectPanelList []*entity.SendingCustomer // 選択したパネルに応じたリスト
	)

	/************ 1. 送客情報を取得 **************/

	sendingCustomerList, err := i.sendingCustomerRepository.GetForSendingCustomerManagement(input.AgentID, input.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. フェーズに応じてリストに振り分け **************/

	for _, sendingCustomer := range sendingCustomerList {
		var (
			// 送客中
			IsSending = sendingCustomer.Phase == null.NewInt(int64(entity.WaitingForInterview), true) || sendingCustomer.Phase == null.NewInt(int64(entity.AcceptSending), true)

			// 送客完了
			IsFinished = sendingCustomer.Phase == null.NewInt(int64(entity.CompleteSending), true)

			// 送客なし/終了
			IsClose = sendingCustomer.Phase == null.NewInt(int64(entity.CloseSending), true)
		)

		// 「すべて」のカウントアップ
		output.AllCount = output.AllCount + 1
		AllList = append(AllList, sendingCustomer)

		// 「送客中」のカウントアップ
		if IsSending {
			output.SendingCount = output.SendingCount + 1
			SendingList = append(SendingList, sendingCustomer)
		}

		// 「送客完了」のカウントアップ
		if IsFinished {
			output.CompleteCount = output.CompleteCount + 1
			CompleteList = append(CompleteList, sendingCustomer)
		}

		// 「送客なし/終了」のカウントアップ
		if IsClose {
			output.CloseCount = output.CloseCount + 1
			CloseList = append(CloseList, sendingCustomer)
		}
	}

	/************ 3. レスポンスするリストをパネルナンバーに応じて決める **************/

	var (
		isAllTab     = input.TabNumber == 0
		isSendingTab = input.TabNumber == 1
		isFinishTab  = input.TabNumber == 2
		isCloseTab   = input.TabNumber == 3
	)

	if isAllTab {
		SelectPanelList = AllList
	} else if isSendingTab {
		SelectPanelList = SendingList
	} else if isFinishTab {
		SelectPanelList = CompleteList
	} else if isCloseTab {
		SelectPanelList = CloseList
	}

	/************ 4. レスポンス情報をまとめる **************/

	// ページの最大数を取得
	output.MaxPageNumber = getSendingCustomerListMaxPage(SelectPanelList)

	// 指定ページの送客求職者20件を取得（本番実装までは1ページあたり5件）
	sendingCustomerList20 := getSendingCustomerListWithPage(SelectPanelList, input.PageNumber)

	output.SendingCustomerList = sendingCustomerList20

	return output, nil
}
