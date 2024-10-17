package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type SendingSaleInteractor interface {
	// 汎用系
	CreateSendingSale(input CreateSendingSaleInput) (CreateSendingSaleOutput, error)
	UpdateSendingSale(input UpdateSendingSaleInput) (UpdateSendingSaleOutput, error)
	GetSendingSaleByID(input GetSendingSaleByIDInput) (GetSendingSaleByIDOutput, error)
	GetSendingSaleByJobSeekerIDAndEnterpriseID(input GetSendingSaleByJobSeekerIDAndEnterpriseIDInput) (GetSendingSaleByJobSeekerIDAndEnterpriseIDOutput, error)
	GetSendingSaleListBySenderAgentIDForMonthly(input GetSendingSaleListBySenderAgentIDForMonthlyInput) (GetSendingSaleListBySenderAgentIDForMonthlyOutput, error)
	GetSendingSaleListByAgentIDForMonthly(input GetSendingSaleListByAgentIDForMonthlyInput) (GetSendingSaleListByAgentIDForMonthlyOutput, error)

	// Admin
	GetAllSendingSaleForMonthly(input GetAllSendingSaleForMonthlyInput) (GetAllSendingSaleForMonthlyOutput, error)
}

type SendingSaleInteractorImpl struct {
	firebase                    usecase.Firebase
	sendgrid                    config.Sendgrid
	sendingSaleRepository       usecase.SendingSaleRepository
	sendingEnterpriseRepository usecase.SendingEnterpriseRepository
	sendingJobSeekerRepository  usecase.SendingJobSeekerRepository
	sendingCustomerRepository   usecase.SendingCustomerRepository
	sendingPhaseRepository      usecase.SendingPhaseRepository
}

// SendingSaleInteractorImpl is an implementation of SendingSaleInteractor
func NewSendingSaleInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	ssR usecase.SendingSaleRepository,
	seR usecase.SendingEnterpriseRepository,
	spR usecase.SendingPhaseRepository,
	sjsR usecase.SendingJobSeekerRepository,
	scR usecase.SendingCustomerRepository,
) SendingSaleInteractor {
	return &SendingSaleInteractorImpl{
		firebase:                    fb,
		sendgrid:                    sg,
		sendingSaleRepository:       ssR,
		sendingEnterpriseRepository: seR,
		sendingPhaseRepository:      spR,
		sendingJobSeekerRepository:  sjsR,
		sendingCustomerRepository:   scR,
	}
}

/****************************************************************************************/
/// 汎用系
//
// 売上を作成する
type CreateSendingSaleInput struct {
	CreateParam entity.CreateSendingSaleParam
}

type CreateSendingSaleOutput struct {
	OK bool
}

func (i *SendingSaleInteractorImpl) CreateSendingSale(input CreateSendingSaleInput) (CreateSendingSaleOutput, error) {
	var (
		output CreateSendingSaleOutput
		err    error
	)

	/************ 1. 売上の登録 **************/

	// 売上の作成
	sale := entity.NewSendingSale(
		input.CreateParam.SendingJobSeekerID,
		input.CreateParam.SendingEnterpriseID,
		input.CreateParam.SystemSales,
		input.CreateParam.Kickback,
	)

	err = i.sendingSaleRepository.Create(sale)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	/************ 2. sending_job_seekerテーブルのphaseカラムを更新 **************/

	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByID(input.CreateParam.SendingJobSeekerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// sending_customerテーブルのphaseカラムがentity.CompleteSending(4)より進んでいない場合は更新
	if entity.SendingJobSeekerPhase(sendingJobSeeker.Phase.Int64) < entity.CompleteSending {
		err = i.sendingJobSeekerRepository.UpdatePhase(input.CreateParam.SendingJobSeekerID, null.NewInt(int64(entity.CompleteSending), true))
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 3. sending_customerテーブルのphaseカラムを更新 **************/

	sendingCustomer, err := i.sendingCustomerRepository.FindByID(sendingJobSeeker.SendingCustomerID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if entity.SendingJobSeekerPhase(sendingCustomer.Phase.Int64) < entity.CompleteSending {
		err = i.sendingCustomerRepository.UpdatePhase(sendingCustomer.ID, null.NewInt(int64(entity.CompleteSending), true))
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	/************ 4. すでにsending_phasesテーブルのレコードを更新 **************/

	err = i.sendingPhaseRepository.UpdatePhaseBySendingJobSeekerIDAndSendingEnterpriseID(input.CreateParam.SendingJobSeekerID, input.CreateParam.SendingEnterpriseID, uint(entity.CompleteSending))
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 売上を更新する
type UpdateSendingSaleInput struct {
	SendingSaleID uint
	UpdateParam   entity.UpdateSendingSaleParam
}

type UpdateSendingSaleOutput struct {
	OK bool
}

func (i *SendingSaleInteractorImpl) UpdateSendingSale(input UpdateSendingSaleInput) (UpdateSendingSaleOutput, error) {
	var (
		output UpdateSendingSaleOutput
		err    error
	)

	sale := &entity.SendingSale{
		SystemSales: input.UpdateParam.SystemSales,
		Kickback:     input.UpdateParam.Kickback,
	}

	err = i.sendingSaleRepository.Update(input.SendingSaleID, sale)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 売上IDを使って売上情報を取得する
type GetSendingSaleByIDInput struct {
	SendingSaleID uint
}

type GetSendingSaleByIDOutput struct {
	SendingSale *entity.SendingSale
}

func (i *SendingSaleInteractorImpl) GetSendingSaleByID(input GetSendingSaleByIDInput) (GetSendingSaleByIDOutput, error) {
	var (
		output GetSendingSaleByIDOutput
		err    error
	)

	sale, err := i.sendingSaleRepository.FindByID(input.SendingSaleID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingSale = sale

	return output, nil
}

type GetSendingSaleByJobSeekerIDAndEnterpriseIDInput struct {
	SendingJobSeekerID  uint
	SendingEnterpriseID uint
}

type GetSendingSaleByJobSeekerIDAndEnterpriseIDOutput struct {
	SendingSale *entity.SendingSale
}

func (i *SendingSaleInteractorImpl) GetSendingSaleByJobSeekerIDAndEnterpriseID(input GetSendingSaleByJobSeekerIDAndEnterpriseIDInput) (GetSendingSaleByJobSeekerIDAndEnterpriseIDOutput, error) {
	var (
		output GetSendingSaleByJobSeekerIDAndEnterpriseIDOutput
		err    error
	)

	sale, err := i.sendingSaleRepository.FindByJobSeekerIDAndEnterpriseID(input.SendingJobSeekerID, input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingSale = sale

	return output, nil
}

// エージェントIDを使って売上を取得する
type GetSendingSaleListBySenderAgentIDForMonthlyInput struct {
	SenderAgentID uint // 送客元のエージェントID
	StartMonth    string
	EndMonth      string
}

type GetSendingSaleListBySenderAgentIDForMonthlyOutput struct {
	SendingSaleList []*entity.SendingSale
}

func (i *SendingSaleInteractorImpl) GetSendingSaleListBySenderAgentIDForMonthly(input GetSendingSaleListBySenderAgentIDForMonthlyInput) (GetSendingSaleListBySenderAgentIDForMonthlyOutput, error) {
	var (
		output   GetSendingSaleListBySenderAgentIDForMonthlyOutput
		err      error
		saleList []*entity.SendingSale
	)

	saleList, err = i.sendingSaleRepository.GetListBySenderAgentIDForMonthly(input.SenderAgentID, input.StartMonth, input.EndMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingSaleList = saleList

	return output, nil
}

// エージェントIDを使って売上を取得する
type GetSendingSaleListByAgentIDForMonthlyInput struct {
	AgentID    uint // 送客管理エージェントID
	StartMonth string
	EndMonth   string
}

type GetSendingSaleListByAgentIDForMonthlyOutput struct {
	SendingSaleList []*entity.SendingSale
}

func (i *SendingSaleInteractorImpl) GetSendingSaleListByAgentIDForMonthly(input GetSendingSaleListByAgentIDForMonthlyInput) (GetSendingSaleListByAgentIDForMonthlyOutput, error) {
	var (
		output   GetSendingSaleListByAgentIDForMonthlyOutput
		err      error
		saleList []*entity.SendingSale
	)

	saleList, err = i.sendingSaleRepository.GetListByAgentIDForMonthly(input.AgentID, input.StartMonth, input.EndMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingSaleList = saleList

	return output, nil
}

/****************************************************************************************/
/// Admin系
//
// エージェントIDを使って売上を取得する
type GetAllSendingSaleForMonthlyInput struct {
	StartMonth string
	EndMonth   string
}

type GetAllSendingSaleForMonthlyOutput struct {
	SendingSaleList []*entity.SendingSale
}

func (i *SendingSaleInteractorImpl) GetAllSendingSaleForMonthly(input GetAllSendingSaleForMonthlyInput) (GetAllSendingSaleForMonthlyOutput, error) {
	var (
		output   GetAllSendingSaleForMonthlyOutput
		err      error
		saleList []*entity.SendingSale
	)

	saleList, err = i.sendingSaleRepository.GetAllForMonthly(input.StartMonth, input.EndMonth)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingSaleList = saleList

	return output, nil
}
