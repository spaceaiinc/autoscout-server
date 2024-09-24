package interactor

import (
	"errors"
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SaleInteractor interface {
	// 汎用系
	CreateSale(input CreateSaleInput) (CreateSaleOutput, error)
	UpdateSale(input UpdateSaleInput) (UpdateSaleOutput, error)
	GetSaleByID(input GetSaleByIDInput) (GetSaleByIDOutput, error)
	GetSaleByJobSeekerID(input GetSaleByJobSeekerIDInput) (GetSaleByJobSeekerIDOutput, error)
	GetSaleListByIDList(input GetSaleListByIDListInput) (GetSaleListByIDListOutput, error)
	GetAccuracySearchList(input GetAccuracySearchListInput) (GetAccuracySearchListOutput, error)
}

type SaleInteractorImpl struct {
	firebase                           usecase.Firebase
	sendgrid                           config.Sendgrid
	oneSignal                          config.OneSignal
	saleRepository                     usecase.SaleRepository
	jobSeekerRepository                usecase.JobSeekerRepository
	chatGroupWithJobSeekerRepository   usecase.ChatGroupWithJobSeekerRepository
	agentStaffRepository               usecase.AgentStaffRepository
	chatMessageWithJobSeekerRepository usecase.ChatMessageWithJobSeekerRepository
}

// SaleInteractorImpl is an implementation of SaleInteractor
func NewSaleInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	os config.OneSignal,
	saR usecase.SaleRepository,
	jsR usecase.JobSeekerRepository,
	cgR usecase.ChatGroupWithJobSeekerRepository,
	asR usecase.AgentStaffRepository,
	cmjR usecase.ChatMessageWithJobSeekerRepository,
) SaleInteractor {
	return &SaleInteractorImpl{
		firebase:                           fb,
		sendgrid:                           sg,
		oneSignal:                          os,
		saleRepository:                     saR,
		jobSeekerRepository:                jsR,
		chatGroupWithJobSeekerRepository:   cgR,
		agentStaffRepository:               asR,
		chatMessageWithJobSeekerRepository: cmjR,
	}
}

/****************************************************************************************/
/// 汎用系
//
// 最終閲覧時間を更新
type CreateSaleInput struct {
	CreateParam entity.CreateOrUpdateSaleParam
}

type CreateSaleOutput struct {
	OK bool
}

func (i *SaleInteractorImpl) CreateSale(input CreateSaleInput) (CreateSaleOutput, error) {
	var (
		output CreateSaleOutput
		err    error
	)

	// タスクの作成
	sale := entity.NewSale(
		input.CreateParam.JobSeekerID,
		input.CreateParam.JobInformationID,
		input.CreateParam.Accuracy,
		input.CreateParam.ContractSignedMonth,
		input.CreateParam.BillingMonth,
		input.CreateParam.BillingAmount,
		input.CreateParam.Cost,
		input.CreateParam.GrossProfit,
		input.CreateParam.RAStaffID,
		input.CreateParam.RaSalesRatio,
		input.CreateParam.CAStaffID,
		input.CreateParam.CaSalesRatio,
	)

	err = i.saleRepository.Create(sale)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 最終閲覧時間を更新
type UpdateSaleInput struct {
	SaleID      uint
	UpdateParam entity.CreateOrUpdateSaleParam
}

type UpdateSaleOutput struct {
	OK bool
}

func (i *SaleInteractorImpl) UpdateSale(input UpdateSaleInput) (UpdateSaleOutput, error) {
	var (
		output UpdateSaleOutput
		err    error
	)

	sale := entity.NewSale(
		input.UpdateParam.JobSeekerID,
		input.UpdateParam.JobInformationID,
		input.UpdateParam.Accuracy,
		input.UpdateParam.ContractSignedMonth,
		input.UpdateParam.BillingMonth,
		input.UpdateParam.BillingAmount,
		input.UpdateParam.Cost,
		input.UpdateParam.GrossProfit,
		input.UpdateParam.RAStaffID,
		input.UpdateParam.RaSalesRatio,
		input.UpdateParam.CAStaffID,
		input.UpdateParam.CaSalesRatio,
	)

	err = i.saleRepository.Update(input.SaleID, sale)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// タスクIDを使ってタスク情報を取得する
type GetSaleByIDInput struct {
	SaleID uint
}

type GetSaleByIDOutput struct {
	Sale *entity.Sale
}

func (i *SaleInteractorImpl) GetSaleByID(input GetSaleByIDInput) (GetSaleByIDOutput, error) {
	var (
		output GetSaleByIDOutput
		err    error
	)

	sale, err := i.saleRepository.FindByID(input.SaleID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Sale = sale

	return output, nil
}

type GetSaleByJobSeekerIDInput struct {
	JobSeekerID uint
}

type GetSaleByJobSeekerIDOutput struct {
	Sale *entity.Sale
}

func (i *SaleInteractorImpl) GetSaleByJobSeekerID(input GetSaleByJobSeekerIDInput) (GetSaleByJobSeekerIDOutput, error) {
	var (
		output GetSaleByJobSeekerIDOutput
		err    error
	)

	sale, err := i.saleRepository.FindByJobSeekerID(input.JobSeekerID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return output, nil
		} else {
			fmt.Println(err)
			return output, err
		}
	}

	output.Sale = sale

	return output, nil
}

// IDリストを使ってヨミ情報リストを取得する
type GetSaleListByIDListInput struct {
	IDList []uint
}

type GetSaleListByIDListOutput struct {
	SaleList []*entity.Sale
}

func (i *SaleInteractorImpl) GetSaleListByIDList(input GetSaleListByIDListInput) (GetSaleListByIDListOutput, error) {
	var (
		output   GetSaleListByIDListOutput
		saleList []*entity.Sale
	)

	saleList, err := i.saleRepository.GetByIDList(input.IDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, s := range saleList {
		fmt.Println("groupID: ", s.TaskGroupID, "phase: ", s.PhaseCategory, "sub: ", s.PhaseSubCategory)
	}
	output.SaleList = saleList

	return output, nil
}

// 自社が絡むヨミ情報を全て取得
type GetAccuracySearchListInput struct {
	SearchParam entity.SearchAccuracy
}

type GetAccuracySearchListOutput struct {
	SaleList      []*entity.Sale
	MaxPageNumber uint
	IDList        []uint
}

func (i *SaleInteractorImpl) GetAccuracySearchList(input GetAccuracySearchListInput) (GetAccuracySearchListOutput, error) {
	var (
		output     GetAccuracySearchListOutput
		inputParam = input.SearchParam
	)

	// 同一テーブルないの項目のみで絞り込みのためクエリ関数内で絞り込みしてる
	saleList, err := i.saleRepository.GetSearchByAgent(inputParam)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// IDListを返す
	for _, sale := range saleList {
		output.IDList = append(output.IDList, sale.ID)
	}

	// ページの最大数を取得（ページ当たり50件）
	output.MaxPageNumber = getSaleListMaxPage(saleList)

	// 指定ページのヨミ50件を取得
	saleList50 := getSaleListWithPage(saleList, inputParam.PageNumber)

	output.SaleList = saleList50

	return output, nil
}
