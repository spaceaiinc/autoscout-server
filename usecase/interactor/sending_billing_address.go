package interactor

import (
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingBillingAddressInteractor interface {
	// 汎用系 API
	UpdateSendingBillingAddress(input UpdateSendingBillingAddressInput) (UpdateSendingBillingAddressOutput, error)

	// 指定IDの請求先情報を取得する関数
	GetSendingBillingAddressByID(input GetSendingBillingAddressByIDInput) (GetSendingBillingAddressByIDOutput, error)

	// 送客エージェントIDを使って請求先一覧情報取得する関数
	GetSendingBillingAddressBySendingEnterpriseID(input GetSendingBillingAddressBySendingEnterpriseIDInput) (GetSendingBillingAddressBySendingEnterpriseIDOutput, error)
	GetSendingBillingAddressListBySendingEnterpriseIDList(input GetSendingBillingAddressListBySendingEnterpriseIDListInput) (GetSendingBillingAddressListBySendingEnterpriseIDListOutput, error)

	// 全ての請求先情報を取得する関数
	GetAllSendingBillingAddress() (GetAllSendingBillingAddressOutput, error)
}

type SendingBillingAddressInteractorImpl struct {
	firebase                             usecase.Firebase
	sendgrid                             config.Sendgrid
	sendingBillingAddressRepository      usecase.SendingBillingAddressRepository
	sendingBillingAddressStaffRepository usecase.SendingBillingAddressStaffRepository
	sendingJobInformationRepository      usecase.SendingJobInformationRepository
}

// SendingBillingAddressInteractorImpl is an implementation of SendingBillingAddressInteractor
func NewSendingBillingAddressInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	baR usecase.SendingBillingAddressRepository,
	bahsR usecase.SendingBillingAddressStaffRepository,
	jiR usecase.SendingJobInformationRepository,
) SendingBillingAddressInteractor {
	return &SendingBillingAddressInteractorImpl{
		firebase:                             fb,
		sendgrid:                             sg,
		sendingBillingAddressRepository:      baR,
		sendingBillingAddressStaffRepository: bahsR,
		sendingJobInformationRepository:      jiR,
	}
}

/****************************************************************************************/
/// 汎用系API
//
// 請求先の更新
type UpdateSendingBillingAddressInput struct {
	UpdateParam             entity.UpdateSendingBillingAddressParam
	SendingBillingAddressID uint
}

type UpdateSendingBillingAddressOutput struct {
	SendingBillingAddress *entity.SendingBillingAddress
}

func (i *SendingBillingAddressInteractorImpl) UpdateSendingBillingAddress(input UpdateSendingBillingAddressInput) (UpdateSendingBillingAddressOutput, error) {
	var (
		output UpdateSendingBillingAddressOutput
		err    error
	)

	sendingBillingAddress := entity.NewSendingBillingAddress(
		input.UpdateParam.SendingEnterpriseID,
		input.UpdateParam.AgentStaffID,
		input.UpdateParam.ContractPhase,
		input.UpdateParam.ContractDate,
		input.UpdateParam.PaymentPolicy,
		input.UpdateParam.CompanyName,
		input.UpdateParam.Address,
		input.UpdateParam.Title,
		input.UpdateParam.ScheduleAdjustmentURL,
		input.UpdateParam.Commission,
	)

	sendingBillingAddress.ID = input.SendingBillingAddressID

	err = i.sendingBillingAddressRepository.Update(sendingBillingAddress, input.SendingBillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	err = i.sendingBillingAddressStaffRepository.Delete(input.SendingBillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hs := range input.UpdateParam.Staffs {
		hrStaff := entity.NewSendingBillingAddressStaff(
			input.SendingBillingAddressID,
			hs.StaffName,
			hs.StaffEmail,
			hs.StaffPhoneNumber,
		)

		err = i.sendingBillingAddressStaffRepository.Create(hrStaff)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.SendingBillingAddress = sendingBillingAddress
	output.SendingBillingAddress.Staffs = input.UpdateParam.Staffs

	return output, nil
}

// 請求先IDを使って請求先情報を取得する
type GetSendingBillingAddressByIDInput struct {
	SendingBillingAddressID uint
}

type GetSendingBillingAddressByIDOutput struct {
	SendingBillingAddress *entity.SendingBillingAddress
}

func (i *SendingBillingAddressInteractorImpl) GetSendingBillingAddressByID(input GetSendingBillingAddressByIDInput) (GetSendingBillingAddressByIDOutput, error) {
	var (
		output GetSendingBillingAddressByIDOutput
		err    error
	)

	sendingBillingAddress, err := i.sendingBillingAddressRepository.FindByID(input.SendingBillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hrStaff, err := i.sendingBillingAddressStaffRepository.FindBySendingBillingAddressID(input.SendingBillingAddressID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hs := range hrStaff {
		value := entity.SendingBillingAddressStaff{
			SendingBillingAddressID: hs.SendingBillingAddressID,
			StaffName:               hs.StaffName,
			StaffEmail:              hs.StaffEmail,
			StaffPhoneNumber:        hs.StaffPhoneNumber,
		}

		sendingBillingAddress.Staffs = append(sendingBillingAddress.Staffs, value)
	}

	output.SendingBillingAddress = sendingBillingAddress

	return output, nil
}

// 送客エージェントIDから請求先情報一覧を取得する
type GetSendingBillingAddressBySendingEnterpriseIDInput struct {
	SendingEnterpriseID uint
}

type GetSendingBillingAddressBySendingEnterpriseIDOutput struct {
	SendingBillingAddress *entity.SendingBillingAddress
}

func (i *SendingBillingAddressInteractorImpl) GetSendingBillingAddressBySendingEnterpriseID(input GetSendingBillingAddressBySendingEnterpriseIDInput) (GetSendingBillingAddressBySendingEnterpriseIDOutput, error) {
	var (
		output GetSendingBillingAddressBySendingEnterpriseIDOutput
		err    error
	)

	sendingBillingAddress, err := i.sendingBillingAddressRepository.FindBySendingEnterpriseID(input.SendingEnterpriseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hrStaff, err := i.sendingBillingAddressStaffRepository.FindBySendingBillingAddressID(sendingBillingAddress.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, hs := range hrStaff {
		if sendingBillingAddress.ID == hs.SendingBillingAddressID {
			value := entity.SendingBillingAddressStaff{
				SendingBillingAddressID: hs.SendingBillingAddressID,
				StaffName:               hs.StaffName,
				StaffEmail:              hs.StaffEmail,
				StaffPhoneNumber:        hs.StaffPhoneNumber,
			}

			sendingBillingAddress.Staffs = append(sendingBillingAddress.Staffs, value)
		}
	}

	output.SendingBillingAddress = sendingBillingAddress

	return output, nil
}

// 送客エージェントIDから請求先情報一覧を取得する
type GetSendingBillingAddressListBySendingEnterpriseIDListInput struct {
	SendingEnterpriseIDList []uint
}

type GetSendingBillingAddressListBySendingEnterpriseIDListOutput struct {
	SendingBillingAddressList []*entity.SendingBillingAddress
}

func (i *SendingBillingAddressInteractorImpl) GetSendingBillingAddressListBySendingEnterpriseIDList(input GetSendingBillingAddressListBySendingEnterpriseIDListInput) (GetSendingBillingAddressListBySendingEnterpriseIDListOutput, error) {
	var (
		output GetSendingBillingAddressListBySendingEnterpriseIDListOutput
		err    error
	)

	sendingBillingAddressList, err := i.sendingBillingAddressRepository.GetBySendingEnterpriseIDList(input.SendingEnterpriseIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	var billingAddressIDList []uint

	for _, billing := range sendingBillingAddressList {
		billingAddressIDList = append(billingAddressIDList, billing.ID)
	}

	hrStaffs, err := i.sendingBillingAddressStaffRepository.GetByBillingAdressIDList(billingAddressIDList)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, billing := range sendingBillingAddressList {
		for _, hs := range hrStaffs {
			if billing.ID == hs.SendingBillingAddressID {
				billing.Staffs = append(billing.Staffs, *hs)
			}
		}
	}

	output.SendingBillingAddressList = sendingBillingAddressList

	return output, nil
}

type GetAllSendingBillingAddressOutput struct {
	SendingBillingAddressList []*entity.SendingBillingAddress
}

func (i *SendingBillingAddressInteractorImpl) GetAllSendingBillingAddress() (GetAllSendingBillingAddressOutput, error) {
	var (
		output GetAllSendingBillingAddressOutput
		err    error
	)

	sendingBillingAddressList, err := i.sendingBillingAddressRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	hrStaff, err := i.sendingBillingAddressStaffRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	for _, sendingBillingAddress := range sendingBillingAddressList {

		for _, hs := range hrStaff {
			if sendingBillingAddress.ID == hs.SendingBillingAddressID {
				value := entity.SendingBillingAddressStaff{
					SendingBillingAddressID: hs.SendingBillingAddressID,
					StaffName:               hs.StaffName,
					StaffEmail:              hs.StaffEmail,
					StaffPhoneNumber:        hs.StaffPhoneNumber,
				}

				sendingBillingAddress.Staffs = append(sendingBillingAddress.Staffs, value)
			}
		}
	}

	output.SendingBillingAddressList = sendingBillingAddressList

	return output, nil
}

/****************************************************************************************/
