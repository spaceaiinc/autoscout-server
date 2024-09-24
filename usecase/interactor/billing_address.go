package interactor

import (
	"fmt"
	"strconv"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type BillingAddressInteractor interface {
	// 汎用系 API
	CreateBillingAddress(input CreateBillingAddressInput) (CreateBillingAddressOutput, error)
	UpdateBillingAddress(input UpdateBillingAddressInput) (UpdateBillingAddressOutput, error)
	UpdateBillingAddressStaffIDByIDListAtOnce(input UpdateBillingAddressAgentStaffInput) (UpdateBillingAddressAgentStaffOutput, error)
	DeleteBillingAddress(input DeleteBillingAddressInput) (DeleteBillingAddressOutput, error)
	// 指定IDの請求先情報を取得する関数
	GetBillingAddressByID(input GetBillingAddressByIDInput) (GetBillingAddressByIDOutput, error)
	// 企業IDを使って請求先一覧情報取得する関数
	GetBillingAddressListByEnterpriseID(input GetBillingAddressListByEnterpriseIDInput) (GetBillingAddressListByEnterpriseIDOutput, error)

	GetBillingAddressListByPageAndAgentID(input GetBillingAddressListByPageAndAgentIDInput) (GetBillingAddressListByPageAndAgentIDOutput, error)
	GetSearchBillingAddressListByPageAndAgentID(input GetSearchBillingAddressListByPageAndAgentIDInput) (GetSearchBillingAddressListByPageAndAgentIDOutput, error)

	// Admin API
	GetAllBillingAddress() (GetAllBillingAddressOutput, error) // 全ての請求先情報を取得する関数
}

type BillingAddressInteractorImpl struct {
	firebase                        usecase.Firebase
	sendgrid                        config.Sendgrid
	billingAddressRepository        usecase.BillingAddressRepository
	billingAddressHRStaffRepository usecase.BillingAddressHRStaffRepository
	billingAddressRAStaffRepository usecase.BillingAddressRAStaffRepository
	jobInformationRepository        usecase.JobInformationRepository
	agentRepository                 usecase.AgentRepository
	agentStaffRepository            usecase.AgentStaffRepository
	taskGroupRepository             usecase.TaskGroupRepository
}

// BillingAddressInteractorImpl is an implementation of BillingAddressInteractor
func NewBillingAddressInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	baR usecase.BillingAddressRepository,
	bahsR usecase.BillingAddressHRStaffRepository,
	barsR usecase.BillingAddressRAStaffRepository,
	jiR usecase.JobInformationRepository,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	tR usecase.TaskGroupRepository,
) BillingAddressInteractor {
	return &BillingAddressInteractorImpl{
		firebase:                        fb,
		sendgrid:                        sg,
		billingAddressRepository:        baR,
		billingAddressHRStaffRepository: bahsR,
		billingAddressRAStaffRepository: barsR,
		jobInformationRepository:        jiR,
		agentRepository:                 aR,
		agentStaffRepository:            asR,
		taskGroupRepository:             tR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
//請求先の作成
type CreateBillingAddressInput struct {
	CreateParam  entity.CreateBillingAddressParam
	EnterpriseID uint
}

type CreateBillingAddressOutput struct {
	BillingAddress *entity.BillingAddress
}

func (i *BillingAddressInteractorImpl) CreateBillingAddress(input CreateBillingAddressInput) (CreateBillingAddressOutput, error) {
	var (
		output CreateBillingAddressOutput
		err    error
	)

	billingAddress := entity.NewBillingAddress(
		input.EnterpriseID,
		input.CreateParam.AgentStaffID,
		input.CreateParam.ContractPhase,
		input.CreateParam.ContractDate,
		input.CreateParam.PaymentPolicy,
		input.CreateParam.CompanyName,
		input.CreateParam.Address,
		input.CreateParam.HowToRecommend,
		input.CreateParam.Title,
	)

	err = i.billingAddressRepository.Create(billingAddress)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hs := range input.CreateParam.HRStaffs {
		hrStaff := entity.NewBillingAddressHRStaff(
			billingAddress.ID,
			hs.HRStaffName,
			hs.HRStaffEmail,
			hs.HRStaffPhoneNumber,
		)

		err = i.billingAddressHRStaffRepository.Create(hrStaff)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	for _, rs := range input.CreateParam.RAStaffs {
		raStaff := entity.NewBillingAddressRAStaff(
			billingAddress.ID,
			rs.BillingAddressStaffName,
			rs.BillingAddressStaffEmail,
			rs.BillingAddressStaffPhoneNumber,
		)

		err = i.billingAddressRAStaffRepository.Create(raStaff)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.BillingAddress = billingAddress
	output.BillingAddress.HRStaffs = input.CreateParam.HRStaffs
	output.BillingAddress.RAStaffs = input.CreateParam.RAStaffs

	return output, nil
}

// 請求先の更新
type UpdateBillingAddressInput struct {
	UpdateParam      entity.UpdateBillingAddressParam
	BillingAddressID uint
}

type UpdateBillingAddressOutput struct {
	BillingAddress *entity.BillingAddress
}

func (i *BillingAddressInteractorImpl) UpdateBillingAddress(input UpdateBillingAddressInput) (UpdateBillingAddressOutput, error) {
	var (
		output UpdateBillingAddressOutput
		err    error
	)

	billingAddress := entity.NewBillingAddress(
		input.UpdateParam.EnterpriseID,
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.ContractPhase,
		input.UpdateParam.ContractDate,
		input.UpdateParam.PaymentPolicy,
		input.UpdateParam.CompanyName,
		input.UpdateParam.Address,
		input.UpdateParam.HowToRecommend,
		input.UpdateParam.Title,
	)

	billingAddress.ID = input.BillingAddressID

	err = i.billingAddressRepository.Update(input.BillingAddressID, billingAddress)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.billingAddressHRStaffRepository.DeleteByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hs := range input.UpdateParam.HRStaffs {
		hrStaff := entity.NewBillingAddressHRStaff(
			input.BillingAddressID,
			hs.HRStaffName,
			hs.HRStaffEmail,
			hs.HRStaffPhoneNumber,
		)

		err = i.billingAddressHRStaffRepository.Create(hrStaff)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	err = i.billingAddressRAStaffRepository.DeleteByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, rs := range input.UpdateParam.RAStaffs {
		raStaff := entity.NewBillingAddressRAStaff(
			input.BillingAddressID,
			rs.BillingAddressStaffName,
			rs.BillingAddressStaffEmail,
			rs.BillingAddressStaffPhoneNumber,
		)

		err = i.billingAddressRAStaffRepository.Create(raStaff)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.BillingAddress = billingAddress
	output.BillingAddress.HRStaffs = input.UpdateParam.HRStaffs
	output.BillingAddress.RAStaffs = input.UpdateParam.RAStaffs

	// 担当RAの値が変更されたことでRAとCA担当が同じになったタスクがないかを確認
	tasGroupkList, err := i.taskGroupRepository.GetNotDoubleSidedSameRAAndCAByStaffID(input.UpdateParam.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// RAとCA担当が同じになったタスクがある場合はそのタスクの「is_double_sided」をtrueに更新する
	if len(tasGroupkList) > 0 {
		var groupIDList []uint

		for _, tg := range tasGroupkList {
			groupIDList = append(groupIDList, tg.ID)
		}

		err := i.taskGroupRepository.UpdateListIsDoubleSided(groupIDList, true)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	return output, nil
}

// 複数の請求先のagent_staff_idカラムを一括更新
type UpdateBillingAddressAgentStaffInput struct {
	AgentStaffID         uint
	BillingAddressIDList []uint
}

type UpdateBillingAddressAgentStaffOutput struct {
	OK bool
}

func (i *BillingAddressInteractorImpl) UpdateBillingAddressStaffIDByIDListAtOnce(input UpdateBillingAddressAgentStaffInput) (UpdateBillingAddressAgentStaffOutput, error) {
	var (
		output UpdateBillingAddressAgentStaffOutput
		err    error
	)

	// 指定の請求先idリストのagent_staff_idを更新
	err = i.billingAddressRepository.UpdateAgentStaffIDByBillingAddressIDList(input.BillingAddressIDList, input.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 請求先の削除
type DeleteBillingAddressInput struct {
	BillingAddressID uint
}

type DeleteBillingAddressOutput struct {
	OK bool
}

func (i *BillingAddressInteractorImpl) DeleteBillingAddress(input DeleteBillingAddressInput) (DeleteBillingAddressOutput, error) {
	var (
		output DeleteBillingAddressOutput
	)

	// 請求先の削除
	err := i.billingAddressRepository.Delete(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 請求先に紐づく求人も削除
	err = i.jobInformationRepository.DeleteByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

// 請求先IDを使って請求先情報を取得する
type GetBillingAddressByIDInput struct {
	BillingAddressID uint
}

type GetBillingAddressByIDOutput struct {
	BillingAddress *entity.BillingAddress
}

func (i *BillingAddressInteractorImpl) GetBillingAddressByID(input GetBillingAddressByIDInput) (GetBillingAddressByIDOutput, error) {
	var (
		output GetBillingAddressByIDOutput
		err    error
	)

	billingAddress, err := i.billingAddressRepository.FindByID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hrStaff, err := i.billingAddressHRStaffRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	raStaff, err := i.billingAddressRAStaffRepository.GetByBillingAddressID(input.BillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hs := range hrStaff {
		billingAddress.HRStaffs = append(billingAddress.HRStaffs, *hs)
	}

	for _, rs := range raStaff {
		billingAddress.RAStaffs = append(billingAddress.RAStaffs, *rs)
	}

	output.BillingAddress = billingAddress

	return output, nil
}

// 企業IDから請求先情報一覧を取得する
type GetBillingAddressListByEnterpriseIDInput struct {
	EnterpriseID uint
}

type GetBillingAddressListByEnterpriseIDOutput struct {
	BillingAddressList []*entity.BillingAddress
}

func (i *BillingAddressInteractorImpl) GetBillingAddressListByEnterpriseID(input GetBillingAddressListByEnterpriseIDInput) (GetBillingAddressListByEnterpriseIDOutput, error) {
	var (
		output GetBillingAddressListByEnterpriseIDOutput
		err    error
	)

	billingAddressList, err := i.billingAddressRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hrStaff, err := i.billingAddressHRStaffRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	raStaff, err := i.billingAddressRAStaffRepository.GetByEnterpriseID(input.EnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, billingAddress := range billingAddressList {
		for _, hs := range hrStaff {
			if billingAddress.ID == hs.BillingAddressID {
				value := entity.BillingAddressHRStaff{
					BillingAddressID:   hs.BillingAddressID,
					HRStaffName:        hs.HRStaffName,
					HRStaffEmail:       hs.HRStaffEmail,
					HRStaffPhoneNumber: hs.HRStaffPhoneNumber,
				}

				billingAddress.HRStaffs = append(billingAddress.HRStaffs, value)
			}
		}

		for _, rs := range raStaff {
			if billingAddress.ID == rs.BillingAddressID {
				value := entity.BillingAddressRAStaff{
					BillingAddressID:               rs.BillingAddressID,
					BillingAddressStaffName:        rs.BillingAddressStaffName,
					BillingAddressStaffEmail:       rs.BillingAddressStaffEmail,
					BillingAddressStaffPhoneNumber: rs.BillingAddressStaffPhoneNumber,
				}

				billingAddress.RAStaffs = append(billingAddress.RAStaffs, value)
			}
		}
	}

	output.BillingAddressList = billingAddressList

	return output, nil
}

// AgentIDから請求先情報一覧を取得する
type GetBillingAddressListByPageAndAgentIDInput struct {
	AgentID    uint
	PageNumber uint
}

type GetBillingAddressListByPageAndAgentIDOutput struct {
	BillingAddressList []*entity.BillingAddress
	MaxPageNumber      uint
	IDList             []uint
}

func (i *BillingAddressInteractorImpl) GetBillingAddressListByPageAndAgentID(input GetBillingAddressListByPageAndAgentIDInput) (GetBillingAddressListByPageAndAgentIDOutput, error) {
	var (
		output GetBillingAddressListByPageAndAgentIDOutput
		err    error
	)

	billingAddressList, err := i.billingAddressRepository.GetByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// IDListを返す
	for _, billingAddress := range billingAddressList {
		output.IDList = append(output.IDList, billingAddress.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getBillingAddressListMaxPage(billingAddressList)

	// 指定ページの企業20件を取得（本番実装までは1ページあたり5件）
	output.BillingAddressList = getBillingAddressListWithPage(billingAddressList, input.PageNumber)

	return output, nil
}

// エージェントIDから企業名一覧を取得する
type GetSearchBillingAddressListByPageAndAgentIDInput struct {
	AgentID     uint
	PageNumber  uint
	SearchParam entity.SearchBillingAddress
}

type GetSearchBillingAddressListByPageAndAgentIDOutput struct {
	BillingAddressList []*entity.BillingAddress
	MaxPageNumber      uint
	IDList             []uint
}

func (i *BillingAddressInteractorImpl) GetSearchBillingAddressListByPageAndAgentID(input GetSearchBillingAddressListByPageAndAgentIDInput) (GetSearchBillingAddressListByPageAndAgentIDOutput, error) {
	var (
		output GetSearchBillingAddressListByPageAndAgentIDOutput
		err    error
	)
	/**
	フリーワードは社名のみ
	*/
	billingAddressList, err := i.billingAddressRepository.GetByAgentIDAndFreeWord(input.AgentID, input.SearchParam.FreeWord)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var (
		billingAddressListWithAgentStaffID []*entity.BillingAddress
	)
	// 営業担当者IDがある場合
	agentStaffID, err := strconv.Atoi(input.SearchParam.AgentStaffID)
	if !(err != nil || agentStaffID == 0) {
		for _, billingAddress := range billingAddressList {
			if billingAddress.AgentStaffID != uint(agentStaffID) {
				continue
			}
			billingAddressListWithAgentStaffID = append(billingAddressListWithAgentStaffID, billingAddress)
		}
	}
	// 営業担当者IDが無い場合
	if err != nil || agentStaffID == 0 {
		billingAddressListWithAgentStaffID = billingAddressList
	}

	// IDListを返す
	for _, billingAddress := range billingAddressListWithAgentStaffID {
		output.IDList = append(output.IDList, billingAddress.ID)
	}

	// ページの最大数を取得
	output.MaxPageNumber = getBillingAddressListMaxPage(billingAddressListWithAgentStaffID)

	// 指定ページの企業20件を取得（本番実装までは1ページあたり5件）
	output.BillingAddressList = getBillingAddressListWithPage(billingAddressListWithAgentStaffID, input.PageNumber)
	return output, nil

}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API
//
// GetAllBillingAddress() (GetAllBillingAddressOutput, error)
type GetAllBillingAddressOutput struct {
	BillingAddressList []*entity.BillingAddress
}

func (i *BillingAddressInteractorImpl) GetAllBillingAddress() (GetAllBillingAddressOutput, error) {
	var (
		output GetAllBillingAddressOutput
		err    error
	)

	billingAddressList, err := i.billingAddressRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hrStaff, err := i.billingAddressHRStaffRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	raStaff, err := i.billingAddressRAStaffRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, billingAddress := range billingAddressList {

		for _, hs := range hrStaff {
			if billingAddress.ID == hs.BillingAddressID {
				value := entity.BillingAddressHRStaff{
					BillingAddressID:   hs.BillingAddressID,
					HRStaffName:        hs.HRStaffName,
					HRStaffEmail:       hs.HRStaffEmail,
					HRStaffPhoneNumber: hs.HRStaffPhoneNumber,
				}

				billingAddress.HRStaffs = append(billingAddress.HRStaffs, value)
			}
		}

		for _, rs := range raStaff {
			if billingAddress.ID == rs.BillingAddressID {
				value := entity.BillingAddressRAStaff{
					BillingAddressID:               rs.BillingAddressID,
					BillingAddressStaffName:        rs.BillingAddressStaffName,
					BillingAddressStaffEmail:       rs.BillingAddressStaffEmail,
					BillingAddressStaffPhoneNumber: rs.BillingAddressStaffPhoneNumber,
				}

				billingAddress.RAStaffs = append(billingAddress.RAStaffs, value)
			}
		}
	}

	output.BillingAddressList = billingAddressList

	return output, nil
}

/****************************************************************************************/
