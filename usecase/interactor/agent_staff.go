package interactor

import (
	"errors"
	"fmt"

	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"

	// "github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
	// "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type AgentStaffInteractor interface {
	// 汎用系 API
	GetAgentStaffListByAgentID(input GetAgentStaffListByAgentIDInput) (GetAgentStaffListByAgentIDOutput, error)
	GetAgentStaffListByAgentIDOrderByID(input GetAgentStaffListByAgentIDOrderByIDInput) (GetAgentStaffListByAgentIDOrderByIDOutput, error)
	UpdateAgentStaff(input UpdateAgentStaffInput) (UpdateAgentStaffOutput, error)
	UpdateAgentStaffEmail(input UpdateAgentStaffEmailInput) (UpdateAgentStaffEmailOutput, error)
	UpdateAgentStaffPassword(input UpdateAgentStaffPasswordInput) (UpdateAgentStaffPasswordOutput, error)
	UpdateAgentStaffUsageEnd(input UpdateAgentStaffUsageEndInput) (UpdateAgentStaffUsageEndOutput, error)
	UpdateAgentStaffUsageReStart(input UpdateAgentStaffUsageReStartInput) (UpdateAgentStaffUsageReStartOutput, error)
	UpdateAgentStaffNotificationJobSeeker(input UpdateAgentStaffNotificationJobSeekerInput) (UpdateAgentStaffNotificationJobSeekerOutput, error)
	UpdateAgentStaffNotificationUnwatched(input UpdateAgentStaffNotificationUnwatchedInput) (UpdateAgentStaffNotificationUnwatchedOutput, error)
	UpdateAgentStaffAuthority(input UpdateAgentStaffAuthorityInput) (UpdateAgentStaffAuthorityOutput, error)
	DeleteAgentStaff(input DeleteAgentStaffInput) (DeleteAgentStaffOutput, error) // 担当者削除　*firebaseのアイパスを削除&DBのis_deletedをtrueにする
	GetAllAgentStaffList() (GetAllAgentStaffListOutput, error)
	GetOhterAgentStaffListByAgentIDAndAllianceAgentID(input GetOhterAgentStaffListByAgentIDAndAllianceAgentIDInput) (GetOhterAgentStaffListByAgentIDAndAllianceAgentIDOutput, error)
	GetAgentStaffListWithSaleNotCreated(input GetAgentStaffListWithSaleNotCreatedInput) (GetAgentStaffListWithSaleNotCreatedOutput, error)
	GetAgentStaffListByAgentIDAndUsageStatusAvailable(input GetAgentStaffListByAgentIDAndUsageStatusAvailableInput) (GetAgentStaffListByAgentIDAndUsageStatusAvailableOutput, error)
	GetAgentStaffListByAgentIDAndIsDeletedFalse(input GetAgentStaffListByAgentIDAndIsDeletedFalseInput) (GetAgentStaffListByAgentIDAndIsDeletedFalseOutput, error)

	// Admin API
	SignUpForAdmin(input AgentStaffSignUpForAdminInput) (AgentStaffSignUpForAdminOutput, error)

	// Agent API
	GetAgentStaffMe(input GetAgentStaffMeInput) (GetAgentStaffMeOutput, error)
}

type AgentStaffInteractorImpl struct {
	firebase                       usecase.Firebase
	sendgrid                       config.Sendgrid
	agentStaffRepository           usecase.AgentStaffRepository
	agentRepository                usecase.AgentRepository
	enterpriseProfileRepository    usecase.EnterpriseProfileRepository
	billingAddressRepository       usecase.BillingAddressRepository
	jobSeekerRepository            usecase.JobSeekerRepository
	notificationForUserRepository  usecase.NotificationForUserRepository
	userNotificationViewRepository usecase.UserNotificationViewRepository
}

// AgentStaffInteractorImpl is an implementation of AgentStaffInteractor
func NewAgentStaffInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	asR usecase.AgentStaffRepository,
	aR usecase.AgentRepository,
	epR usecase.EnterpriseProfileRepository,
	baR usecase.BillingAddressRepository,
	jsR usecase.JobSeekerRepository,
	nfuR usecase.NotificationForUserRepository,
	unvR usecase.UserNotificationViewRepository,
) AgentStaffInteractor {
	return &AgentStaffInteractorImpl{
		firebase:                       fb,
		sendgrid:                       sg,
		agentStaffRepository:           asR,
		agentRepository:                aR,
		enterpriseProfileRepository:    epR,
		billingAddressRepository:       baR,
		jobSeekerRepository:            jsR,
		notificationForUserRepository:  nfuR,
		userNotificationViewRepository: unvR,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
type GetAgentStaffListByAgentIDInput struct {
	Token   string
	AgentID uint
}

type GetAgentStaffListByAgentIDOutput struct {
	AgentStaffList []*entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetAgentStaffListByAgentID(input GetAgentStaffListByAgentIDInput) (GetAgentStaffListByAgentIDOutput, error) {
	var (
		output GetAgentStaffListByAgentIDOutput
	)

	agentStaffList, err := i.agentStaffRepository.GetByAgentID(input.AgentID)
	if err != nil {
		return output, err
	}

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	//  FirebaseIDが一致する担当者を一番上にソート
	for agentStaffI, agentStaff := range agentStaffList {
		if agentStaff.FirebaseID == firebaseID {
			agentStaffList = append(agentStaffList[:agentStaffI], agentStaffList[agentStaffI+1:]...)
			agentStaffList = append([]*entity.AgentStaff{agentStaff}, agentStaffList...)
			break
		}
	}

	output.AgentStaffList = agentStaffList

	return output, nil
}

type GetAgentStaffListByAgentIDOrderByIDInput struct {
	AgentID      uint
	AgentStaffID uint
}

type GetAgentStaffListByAgentIDOrderByIDOutput struct {
	AgentStaffList []*entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetAgentStaffListByAgentIDOrderByID(input GetAgentStaffListByAgentIDOrderByIDInput) (GetAgentStaffListByAgentIDOrderByIDOutput, error) {
	var (
		output GetAgentStaffListByAgentIDOrderByIDOutput
	)

	agentStaffList, err := i.agentStaffRepository.GetByAgentID(input.AgentID)
	if err != nil {
		return output, err
	}

	//  IDが一致する担当者を一番上にソート
	for agentStaffI, agentStaff := range agentStaffList {
		if agentStaff.ID == input.AgentStaffID {
			agentStaffList = append(agentStaffList[:agentStaffI], agentStaffList[agentStaffI+1:]...)
			agentStaffList = append([]*entity.AgentStaff{agentStaff}, agentStaffList...)
			break
		}
	}

	output.AgentStaffList = agentStaffList

	return output, nil
}

type UpdateAgentStaffInput struct {
	AgentStaffID uint
	UpdateParam  entity.AgentStaffUpdateParam
}

type UpdateAgentStaffOutput struct {
	AgentStaff *entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaff(input UpdateAgentStaffInput) (UpdateAgentStaffOutput, error) {
	var (
		output UpdateAgentStaffOutput
	)

	agentStaff, err := i.agentStaffRepository.FindByID(input.AgentStaffID)
	if err != nil {
		return output, err
	}

	agentStaff = entity.NewAgentStaff(
		agentStaff.AgentID,    // 更新しない
		agentStaff.FirebaseID, // 更新しない
		input.UpdateParam.Authority,
		input.UpdateParam.StaffName,
		input.UpdateParam.Furigana,
		agentStaff.Email, // 更新しない
		input.UpdateParam.StaffPhoneNumber,
		input.UpdateParam.Department,
		input.UpdateParam.Position,
		input.UpdateParam.Remarks,
		input.UpdateParam.UsageStatus,
		input.UpdateParam.Notification,
		input.UpdateParam.NotificationJobSeeker,
		input.UpdateParam.NotificationUnwatched,
		// agentStaff.LastLogin, // 更新しない
	)

	err = i.agentStaffRepository.Update(input.AgentStaffID, agentStaff)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentStaff.ID = input.AgentStaffID

	output.AgentStaff = agentStaff

	return output, nil
}

type UpdateAgentStaffEmailInput struct {
	UpdateParam  entity.AgentStaffEmailUpdateParam
	AgentStaffID uint
}

type UpdateAgentStaffEmailOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaffEmail(input UpdateAgentStaffEmailInput) (UpdateAgentStaffEmailOutput, error) {
	var (
		output UpdateAgentStaffEmailOutput
	)

	agentStaff, err := i.agentStaffRepository.FindByID(input.AgentStaffID)
	if err != nil {
		return output, err
	}

	// FirebaseAuthで登録中のメールアドレス更新
	err = i.firebase.UpdateEmail(
		input.UpdateParam.Email,
		agentStaff.FirebaseID,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// DBで登録中のメールアドレス更新
	err = i.agentStaffRepository.UpdateEmail(input.AgentStaffID, input.UpdateParam.Email)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type UpdateAgentStaffPasswordInput struct {
	UpdateParam  entity.AgentStaffPasswordUpdateParam
	AgentStaffID uint
}

type UpdateAgentStaffPasswordOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaffPassword(input UpdateAgentStaffPasswordInput) (UpdateAgentStaffPasswordOutput, error) {
	var (
		output UpdateAgentStaffPasswordOutput
	)

	agentStaff, err := i.agentStaffRepository.FindByID(input.AgentStaffID)
	if err != nil {
		return output, err
	}

	// FirebaseAuthで登録中のメールアドレス更新
	err = i.firebase.UpdatePassword(
		input.UpdateParam.Password,
		agentStaff.FirebaseID,
	)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type UpdateAgentStaffUsageEndInput struct {
	Param entity.DeleteAgentStaffParam
}

type UpdateAgentStaffUsageEndOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaffUsageEnd(input UpdateAgentStaffUsageEndInput) (UpdateAgentStaffUsageEndOutput, error) {
	var (
		output UpdateAgentStaffUsageEndOutput
	)

	err := i.agentStaffRepository.UpdateUsageEnd(input.Param.AgentStaffID)
	if err != nil {
		return output, err
	}

	// RA引き継ぎ処理
	if input.Param.NewRAStaffID.Int64 != 0 {
		enterpriseList, err := i.enterpriseProfileRepository.GetByAgentStaffID(input.Param.AgentStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, enterprise := range enterpriseList {
			err = i.enterpriseProfileRepository.UpdateAgentStaffID(enterprise.ID, uint(input.Param.NewRAStaffID.Int64))
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		billingAddressList, err := i.billingAddressRepository.GetByAgentStaffID(input.Param.AgentStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, billingAddress := range billingAddressList {
			err = i.billingAddressRepository.UpdateAgentStaffID(billingAddress.ID, uint(input.Param.NewRAStaffID.Int64))
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// CA引き継ぎ処理
	if input.Param.NewCAStaffID.Int64 != 0 {
		jobSeekerList, err := i.jobSeekerRepository.GetByAgentStaffIDAndNotRelease(input.Param.AgentStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, jobSeeker := range jobSeekerList {
			err = i.jobSeekerRepository.UpdateStaffID(jobSeeker.ID, input.Param.NewCAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	output.OK = true

	return output, nil
}

type UpdateAgentStaffUsageReStartInput struct {
	AgentStaffID uint
}

type UpdateAgentStaffUsageReStartOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaffUsageReStart(input UpdateAgentStaffUsageReStartInput) (UpdateAgentStaffUsageReStartOutput, error) {
	var (
		output UpdateAgentStaffUsageReStartOutput
	)

	err := i.agentStaffRepository.UpdateUsageStart(input.AgentStaffID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// メール通知（求職者）の更新 body: {agent_staff_id, notification_job_seeker}
type UpdateAgentStaffNotificationJobSeekerInput struct {
	UpdateParam entity.UpdateAgentStaffNotificationJobSeekerParam
}

type UpdateAgentStaffNotificationJobSeekerOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaffNotificationJobSeeker(input UpdateAgentStaffNotificationJobSeekerInput) (UpdateAgentStaffNotificationJobSeekerOutput, error) {
	var (
		output UpdateAgentStaffNotificationJobSeekerOutput
	)

	err := i.agentStaffRepository.UpdateNotificationJobSeeker(input.UpdateParam.AgentStaffID, input.UpdateParam.NotificationJobSeeker)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// メール通知（未処理・未読）の更新 body: {agent_staff_id, notification_unwatched}
type UpdateAgentStaffNotificationUnwatchedInput struct {
	UpdateParam entity.UpdateAgentStaffNotificationUnwatchedParam
}

type UpdateAgentStaffNotificationUnwatchedOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaffNotificationUnwatched(input UpdateAgentStaffNotificationUnwatchedInput) (UpdateAgentStaffNotificationUnwatchedOutput, error) {
	var (
		output UpdateAgentStaffNotificationUnwatchedOutput
	)

	err := i.agentStaffRepository.UpdateNotificationUnwatched(input.UpdateParam.AgentStaffID, input.UpdateParam.NotificationUnwatched)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// 管理権限の更新 body: {agent_staff_id, authority}
type UpdateAgentStaffAuthorityInput struct {
	UpdateParam entity.UpdateAgentStaffAuthorityParam
}

type UpdateAgentStaffAuthorityOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) UpdateAgentStaffAuthority(input UpdateAgentStaffAuthorityInput) (UpdateAgentStaffAuthorityOutput, error) {
	var (
		output UpdateAgentStaffAuthorityOutput
	)

	err := i.agentStaffRepository.UpdateAuthority(input.UpdateParam.AgentStaffID, input.UpdateParam.Authority)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

// 担当者削除　*firebaseのアイパスを削除&DBのis_deletedをtrueにする
type DeleteAgentStaffInput struct {
	Param entity.DeleteAgentStaffParam
}

type DeleteAgentStaffOutput struct {
	OK bool
}

func (i *AgentStaffInteractorImpl) DeleteAgentStaff(input DeleteAgentStaffInput) (DeleteAgentStaffOutput, error) {
	var (
		output DeleteAgentStaffOutput
		err    error
	)

	// RA引き継ぎ処理
	if input.Param.NewRAStaffID.Int64 != 0 {
		enterpriseList, err := i.enterpriseProfileRepository.GetByAgentStaffID(input.Param.AgentStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, enterprise := range enterpriseList {
			err = i.enterpriseProfileRepository.UpdateAgentStaffID(enterprise.ID, uint(input.Param.NewRAStaffID.Int64))
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}

		billingAddressList, err := i.billingAddressRepository.GetByAgentStaffID(input.Param.AgentStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, billingAddress := range billingAddressList {
			err = i.billingAddressRepository.UpdateAgentStaffID(billingAddress.ID, uint(input.Param.NewRAStaffID.Int64))
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// CA引き継ぎ処理
	if input.Param.NewCAStaffID.Int64 != 0 {
		jobSeekerList, err := i.jobSeekerRepository.GetByAgentStaffIDAndNotRelease(input.Param.AgentStaffID)
		if err != nil {
			fmt.Println(err)
			return output, err
		}

		for _, jobSeeker := range jobSeekerList {
			err = i.jobSeekerRepository.UpdateStaffID(jobSeeker.ID, input.Param.NewCAStaffID)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		}
	}

	// 担当者情報を取得
	agentStaff, err := i.agentStaffRepository.FindByID(input.Param.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// firebaseのアイパスを削除
	err = i.firebase.DeleteUser(agentStaff.FirebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// is_deletedをtrueにする
	err = i.agentStaffRepository.Delete(input.Param.AgentStaffID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type GetAllAgentStaffListOutput struct {
	AgentStaffList []*entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetAllAgentStaffList() (GetAllAgentStaffListOutput, error) {
	var (
		output GetAllAgentStaffListOutput
	)

	agentStaffList, err := i.agentStaffRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentStaffList = agentStaffList

	return output, nil
}

type GetOhterAgentStaffListByAgentIDAndAllianceAgentIDInput struct {
	AgentID         uint
	AllianceAgentID uint
	AgentStaffID    uint
}

type GetOhterAgentStaffListByAgentIDAndAllianceAgentIDOutput struct {
	AgentStaffList []*entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetOhterAgentStaffListByAgentIDAndAllianceAgentID(input GetOhterAgentStaffListByAgentIDAndAllianceAgentIDInput) (GetOhterAgentStaffListByAgentIDAndAllianceAgentIDOutput, error) {
	var (
		output GetOhterAgentStaffListByAgentIDAndAllianceAgentIDOutput
	)

	agentStaffList, err := i.agentStaffRepository.GetByNotIDAndAgentIDAndAllianceAgentID(input.AgentStaffID, input.AgentID, input.AllianceAgentID)
	if err != nil {
		return output, err
	}

	output.AgentStaffList = agentStaffList

	return output, nil
}

type GetAgentStaffListWithSaleNotCreatedInput struct {
	Token        string
	AgentID      uint
	ManagementID uint
}

type GetAgentStaffListWithSaleNotCreatedOutput struct {
	AgentStaffList []*entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetAgentStaffListWithSaleNotCreated(input GetAgentStaffListWithSaleNotCreatedInput) (GetAgentStaffListWithSaleNotCreatedOutput, error) {
	var (
		output GetAgentStaffListWithSaleNotCreatedOutput
	)

	agentStaffList, err := i.agentStaffRepository.GetByAgentIDAndNotManagementID(input.AgentID, input.ManagementID)
	if err != nil {
		return output, err
	}

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	//  FirebaseIDが一致する担当者を一番上にソート
	for agentStaffI, agentStaff := range agentStaffList {
		if agentStaff.FirebaseID == firebaseID {
			agentStaffList = append(agentStaffList[:agentStaffI], agentStaffList[agentStaffI+1:]...)
			agentStaffList = append([]*entity.AgentStaff{agentStaff}, agentStaffList...)
			break
		}
	}

	output.AgentStaffList = agentStaffList

	return output, nil
}

// 利用可能な担当者一覧を取得
type GetAgentStaffListByAgentIDAndUsageStatusAvailableInput struct {
	Token   string
	AgentID uint
}

type GetAgentStaffListByAgentIDAndUsageStatusAvailableOutput struct {
	AgentStaffList []*entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetAgentStaffListByAgentIDAndUsageStatusAvailable(input GetAgentStaffListByAgentIDAndUsageStatusAvailableInput) (GetAgentStaffListByAgentIDAndUsageStatusAvailableOutput, error) {
	var (
		output GetAgentStaffListByAgentIDAndUsageStatusAvailableOutput
	)

	agentStaffList, err := i.agentStaffRepository.GetByAgentIDAndUsageStatusAvailable(input.AgentID)
	if err != nil {
		return output, err
	}

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	//  FirebaseIDが一致する担当者を一番上にソート
	for agentStaffI, agentStaff := range agentStaffList {
		if agentStaff.FirebaseID == firebaseID {
			agentStaffList = append(agentStaffList[:agentStaffI], agentStaffList[agentStaffI+1:]...)
			agentStaffList = append([]*entity.AgentStaff{agentStaff}, agentStaffList...)
			break
		}
	}

	output.AgentStaffList = agentStaffList

	return output, nil
}

// 未削除の担当者一覧を取得
type GetAgentStaffListByAgentIDAndIsDeletedFalseInput struct {
	Token   string
	AgentID uint
}

type GetAgentStaffListByAgentIDAndIsDeletedFalseOutput struct {
	AgentStaffList []*entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetAgentStaffListByAgentIDAndIsDeletedFalse(input GetAgentStaffListByAgentIDAndIsDeletedFalseInput) (GetAgentStaffListByAgentIDAndIsDeletedFalseOutput, error) {
	var (
		output GetAgentStaffListByAgentIDAndIsDeletedFalseOutput
	)

	agentStaffList, err := i.agentStaffRepository.GetByAgentIDAndIsDeletedFalse(input.AgentID)
	if err != nil {
		return output, err
	}

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	//  FirebaseIDが一致する担当者を一番上にソート
	for agentStaffI, agentStaff := range agentStaffList {
		if agentStaff.FirebaseID == firebaseID {
			agentStaffList = append(agentStaffList[:agentStaffI], agentStaffList[agentStaffI+1:]...)
			agentStaffList = append([]*entity.AgentStaff{agentStaff}, agentStaffList...)
			break
		}
	}

	output.AgentStaffList = agentStaffList

	return output, nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API
//
type AgentStaffSignUpForAdminInput struct {
	SignUpParam entity.AgentStaffSignUpForAdminParam
	AgentID     uint
}

type AgentStaffSignUpForAdminOutput struct {
	AgentStaff *entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) SignUpForAdmin(input AgentStaffSignUpForAdminInput) (AgentStaffSignUpForAdminOutput, error) {
	var (
		output AgentStaffSignUpForAdminOutput
	)

	// Firebase Authentication Create
	firebaseID, err := i.firebase.CreateUser(
		input.SignUpParam.Email,
		input.SignUpParam.Password,
	)
	if err != nil {
		return output, err
	}

	staff := entity.NewAgentStaff(
		input.AgentID,
		firebaseID,
		input.SignUpParam.Authority,
		input.SignUpParam.StaffName,
		input.SignUpParam.Furigana,
		input.SignUpParam.Email,
		input.SignUpParam.StaffPhoneNumber,
		input.SignUpParam.Department,
		input.SignUpParam.Position,
		input.SignUpParam.Remarks,
		input.SignUpParam.UsageStatus,
		input.SignUpParam.Notification,
		input.SignUpParam.NotificationJobSeeker,
		input.SignUpParam.NotificationUnwatched,
	)

	err = i.agentStaffRepository.Create(staff)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// お知らせの既読を作成
	// 最新のお知らせを取得
	latestNotification, err := i.notificationForUserRepository.FindLatest()
	if errors.Is(err, entity.ErrNotFound) {
		// お知らせがない場合は何もしない
		err = nil
	} else if err != nil {
		fmt.Println(err)
		return output, err
	} else {
		// 最新の既読と最新のお知らせが異なる場合は最新のお知らせを既読にする
		userNotificationView := entity.NewUserNotificationView(
			latestNotification.ID,
			staff.ID,
		)

		err = i.userNotificationViewRepository.Create(userNotificationView)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.AgentStaff = staff

	return output, nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Agent API
//

type GetAgentStaffMeInput struct {
	Token string
}

type GetAgentStaffMeOutput struct {
	AgentStaff *entity.AgentStaff
}

func (i *AgentStaffInteractorImpl) GetAgentStaffMe(input GetAgentStaffMeInput) (GetAgentStaffMeOutput, error) {
	var (
		output GetAgentStaffMeOutput
	)

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		return output, err
	}

	output.AgentStaff = agentStaff

	return output, nil
}

/****************************************************************************************/
