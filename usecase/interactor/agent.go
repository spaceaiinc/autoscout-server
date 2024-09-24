package interactor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type AgentInteractor interface {
	// 汎用系 API
	GetAgentByAgentID(input GetAgentByAgentIDInput) (GetAgentByAgentIDOutput, error)
	GetAgentByAgentUUID(input GetAgentByAgentUUIDInput) (GetAgentByAgentUUIDOutput, error)
	UpdateAgent(input UpdateAgentInput) (UpdateAgentOutput, error)
	UpdateAgentAgreementFileURL(input UpdateAgentAgreementFileURLInput) (UpdateAgentAgreementFileURLOutput, error)
	UpdateAgentForAdmin(input UpdateAgentForAdminInput) (UpdateAgentForAdminOutput, error)
	GetAllianceAgentListByAgentID(input GetAllianceAgentListByAgentIDInput) (GetAllianceAgentListByAgentIDOutput, error)
	GetAllianceAgentListByAgentIDForSelect(input GetAllianceAgentListByAgentIDForSelectInput) (GetAllianceAgentListByAgentIDForSelectOutput, error)
	GetAgreementFileURL(input GetAgreementFileURLInput) (GetAgreementFileURLOutput, error)
	GetAgentClaimList(input GetAgentClaimListInput) (GetAgentClaimListOutput, error)

	// Admin API
	AgentSignUp(input AgentSignUpInput) (AgentSignUpOutput, error)
	AgentAndAgentStaffSignUp(input AgentAndAgentStaffSignUpInput) (AgentAndAgentStaffSignUpOutput, error)
	GetAllAgentList() (GetAllAgentListOutput, error)
	// 指定IDのエージェント情報を削除

	// CreateAgentLineInformation(input CreateAgentLineInformationInput) (CreateAgentLineInformationOutput, error)

	// LINE関連　API
	UpdateAgentLineChannel(input UpdateAgentLineChannelInput) (UpdateAgentLineChannelOutput, error)
	GetAgentLineChannelByAgentID(input GetAgentLineChannelByAgentIDInput) (GetAgentLineChannelByAgentIDOutput, error)
	GetAgentLineLoginChannelIDByAgentUUID(input GetAgentLineLoginChannelIDByAgentUUIDInput) (GetAgentLineLoginChannelIDByAgentUUIDOutput, error)
}

type AgentInteractorImpl struct {
	firebase                           usecase.Firebase
	sendgrid                           config.Sendgrid
	agentRepository                    usecase.AgentRepository
	agentAllianceRepository            usecase.AgentAllianceRepository
	agentStaffRepository               usecase.AgentStaffRepository
	jobSeekerRepository                usecase.JobSeekerRepository
	chatGroupWithJobSeekerRepository   usecase.ChatGroupWithJobSeekerRepository
	chatMessageWithJobSeekerRepository usecase.ChatMessageWithJobSeekerRepository
	chatGroupWithAgentRepository       usecase.ChatGroupWithAgentRepository
	notificationForUserRepository      usecase.NotificationForUserRepository
	userNotificationViewRepository     usecase.UserNotificationViewRepository
}

// AgentInteractorImpl is an implementation of AgentInteractor
func NewAgentInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	aR usecase.AgentRepository,
	aaR usecase.AgentAllianceRepository,
	asR usecase.AgentStaffRepository,
	jsR usecase.JobSeekerRepository,
	cgwjsR usecase.ChatGroupWithJobSeekerRepository,
	cmR usecase.ChatMessageWithJobSeekerRepository,
	cgwaR usecase.ChatGroupWithAgentRepository,
	nfuR usecase.NotificationForUserRepository,
	unvR usecase.UserNotificationViewRepository,
) AgentInteractor {
	return &AgentInteractorImpl{
		firebase:                           fb,
		sendgrid:                           sg,
		agentRepository:                    aR,
		agentAllianceRepository:            aaR,
		agentStaffRepository:               asR,
		jobSeekerRepository:                jsR,
		chatGroupWithJobSeekerRepository:   cgwjsR,
		chatMessageWithJobSeekerRepository: cmR,
		chatGroupWithAgentRepository:       cgwaR,
		notificationForUserRepository:      nfuR,
		userNotificationViewRepository:     unvR,
	}
}

/****************************************************************************************/
/// 汎用系 API
//

type GetAgentByAgentIDInput struct {
	AgentID uint
}

type GetAgentByAgentIDOutput struct {
	Agent *entity.Agent
}

func (i *AgentInteractorImpl) GetAgentByAgentID(input GetAgentByAgentIDInput) (GetAgentByAgentIDOutput, error) {
	var (
		output GetAgentByAgentIDOutput
	)

	agent, err := i.agentRepository.FindByID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// enterpriseType

	output.Agent = agent

	return output, nil
}

type GetAgentByAgentUUIDInput struct {
	UUID uuid.UUID
}

type GetAgentByAgentUUIDOutput struct {
	Agent *entity.Agent
}

func (i *AgentInteractorImpl) GetAgentByAgentUUID(input GetAgentByAgentUUIDInput) (GetAgentByAgentUUIDOutput, error) {
	var (
		output GetAgentByAgentUUIDOutput
	)

	agent, err := i.agentRepository.FindByUUID(input.UUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Agent = agent

	return output, nil
}

type GetAllianceAgentListByAgentIDInput struct {
	AgentID uint
}

type GetAllianceAgentListByAgentIDOutput struct {
	AgentList []*entity.Agent
}

func (i *AgentInteractorImpl) GetAllianceAgentListByAgentID(input GetAllianceAgentListByAgentIDInput) (GetAllianceAgentListByAgentIDOutput, error) {
	var (
		output    GetAllianceAgentListByAgentIDOutput
		agentList []*entity.Agent
	)

	allianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentList = getAllianceAgentList(input.AgentID, allianceList)
	output.AgentList = agentList

	return output, nil
}

type GetAllianceAgentListByAgentIDForSelectInput struct {
	AgentID uint
}

type GetAllianceAgentListByAgentIDForSelectOutput struct {
	AgentList []*entity.Agent
}

func (i *AgentInteractorImpl) GetAllianceAgentListByAgentIDForSelect(input GetAllianceAgentListByAgentIDForSelectInput) (GetAllianceAgentListByAgentIDForSelectOutput, error) {
	var (
		output    GetAllianceAgentListByAgentIDForSelectOutput
		agentList []*entity.Agent
	)
	allianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentList = getAllianceAgentList(input.AgentID, allianceList)
	output.AgentList = agentList

	return output, nil
}

type GetAgreementFileURLInput struct {
	Token string
}

type GetAgreementFileURLOutput struct {
	Agent *entity.Agent
}

func (i *AgentInteractorImpl) GetAgreementFileURL(input GetAgreementFileURLInput) (GetAgreementFileURLOutput, error) {
	var (
		output GetAgreementFileURLOutput
		agent  *entity.Agent
	)

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agent, err = i.agentRepository.FindByID(agentStaff.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Agent = agent

	return output, nil
}

type GetAgentClaimListInput struct {
	Token      string
	PageNumber uint
}

type GetAgentClaimListOutput struct {
	AgentList     []*entity.Agent
	MaxPageNumber uint
	IDList        []uint
}

func (i *AgentInteractorImpl) GetAgentClaimList(input GetAgentClaimListInput) (GetAgentClaimListOutput, error) {
	var (
		output GetAgentClaimListOutput
	)

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	if agentStaff.AgentID != 1 || agentStaff.Authority != null.NewInt(int64(entity.AuthorityAdmin), true) {
		fmt.Println(err)
		return output, err
	}

	agentList, err := i.agentRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	agentStaffList, err := i.agentStaffRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ページの最大数を取得
	output.MaxPageNumber = getAgentListMaxPage(agentList)

	// タスクIDのリストを作成
	for _, agent := range agentList {
		output.IDList = append(output.IDList, agent.ID)
	}

	// 指定ページのタスクグループ20件を取得（本番実装までは1ページあたり5件）
	agentList20 := getAgentListWithPage(agentList, input.PageNumber)

	now := time.Now()

	// Get the first day of this month
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Get the last day of this month
	monthEnd := monthStart.AddDate(0, 1, -1)

	for _, agent := range agentList20 {
		var claimCount int
		for _, staff := range agentStaffList {
			if (agent.ID == staff.AgentID) &&
				(staff.UsageStartDate.Before(monthEnd) || staff.UsageStartDate.Equal(monthEnd)) && // staff.UsageStartDate <= end of the month
				(staff.UsageEndDate.Before(staff.UsageStartDate) || staff.UsageEndDate.Equal(staff.UsageStartDate)) { // staff.UsageEndDate <= staff.UsageStartDate
				claimCount = claimCount + 1
			}
		}
		agent.ClaimAccountCount = claimCount
	}

	output.AgentList = agentList20

	return output, nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Admin API
//

type AgentSignUpInput struct {
	CreateParam entity.CreateOrUpdateAgentParam
}

type AgentSignUpOutput struct {
	Agent *entity.Agent
}

func (i *AgentInteractorImpl) AgentSignUp(input AgentSignUpInput) (AgentSignUpOutput, error) {
	var (
		output AgentSignUpOutput
		err    error
	)

	agent := entity.NewAgent(
		input.CreateParam.AgentName,
		input.CreateParam.OfficeLocation,
		input.CreateParam.Representative,
		input.CreateParam.Establish,
		input.CreateParam.CorporateSiteURL,
		input.CreateParam.PermissionCode,
		input.CreateParam.PermissionYear,
		input.CreateParam.WorkersCount,
		input.CreateParam.PhoneNumber,
		input.CreateParam.InterviewAdjustmentEmail,
		input.CreateParam.AgreementFileURL,
		"", // input.CreateParam.LineBotID,
		input.CreateParam.LineMessagingChannelSecret,
		input.CreateParam.LineMessagingChannelAccessToken,
		input.CreateParam.LineLoginChannelID,
		input.CreateParam.LineLoginChannelSecret,
		"", // input.CreateParam.SendingAgreementFileURL,
		input.CreateParam.IsCRMActive,
		input.CreateParam.IsAllianceActive,
		true,                 // input.CreateParam.IsSendingActive,
		null.NewInt(0, true), // 送客タイプは通常として登録
	)

	err = i.agentRepository.Create(agent)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Agent = agent

	fmt.Println(agent)

	return output, nil
}

type AgentAndAgentStaffSignUpInput struct {
	CreateParam entity.CreateAgentAndAgentStaffParam
}

type AgentAndAgentStaffSignUpOutput struct {
	Agent      *entity.Agent
	AgentStaff *entity.AgentStaff
}

func (i *AgentInteractorImpl) AgentAndAgentStaffSignUp(input AgentAndAgentStaffSignUpInput) (AgentAndAgentStaffSignUpOutput, error) {
	var (
		output AgentAndAgentStaffSignUpOutput
		err    error
	)

	// エージェントの作成   LINE情報は登録しない
	agent := &entity.Agent{
		AgentName:                input.CreateParam.AgentName,
		OfficeLocation:           input.CreateParam.OfficeLocation,
		Representative:           input.CreateParam.Representative,
		Establish:                input.CreateParam.Establish,
		CorporateSiteURL:         input.CreateParam.CorporateSiteURL,
		PermissionCode:           input.CreateParam.PermissionCode,
		PermissionYear:           input.CreateParam.PermissionYear,
		WorkersCount:             input.CreateParam.WorkersCount,
		PhoneNumber:              input.CreateParam.PhoneNumber,
		InterviewAdjustmentEmail: input.CreateParam.InterviewAdjustmentEmail,
		IsCRMActive:              input.CreateParam.IsCRMActive,
		IsAllianceActive:         input.CreateParam.IsAllianceActive,
		IsSendingActive:          input.CreateParam.IsSendingActive,
		SendingType:              input.CreateParam.SendingType,
	}

	err = i.agentRepository.Create(agent)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントのチャットグループを作成
	chatGroup := entity.NewChatGroupWithAgent(
		agent.ID,
		agent.ID,
	)

	err = i.chatGroupWithAgentRepository.Create(chatGroup)
	if err != nil {
		return output, err
	}

	// Firebase Authentication Create
	firebaseID, err := i.firebase.CreateUser(
		input.CreateParam.Email,
		input.CreateParam.Password,
	)
	if err != nil {
		return output, err
	}

	// エージェントスタッフの作成
	staff := entity.NewAgentStaff(
		agent.ID,
		firebaseID,
		input.CreateParam.Authority,
		input.CreateParam.StaffName,
		input.CreateParam.Furigana,
		input.CreateParam.Email,
		input.CreateParam.StaffPhoneNumber,
		input.CreateParam.Department,
		input.CreateParam.Position,
		input.CreateParam.Remarks,
		input.CreateParam.UsageStatus,
		input.CreateParam.Notification,
		input.CreateParam.NotificationJobSeeker,
		input.CreateParam.NotificationUnwatched,
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

	output.Agent = agent
	output.AgentStaff = staff

	fmt.Println(agent)
	fmt.Println(staff)

	return output, nil
}

type GetAllAgentListOutput struct {
	AgentList []*entity.Agent
}

func (i *AgentInteractorImpl) GetAllAgentList() (GetAllAgentListOutput, error) {
	var (
		output GetAllAgentListOutput
	)

	agentList, err := i.agentRepository.All()
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentList = agentList

	return output, nil
}

/****************************************************************************************/

/****************************************************************************************/
/// Agent API
//
type UpdateAgentInput struct {
	AgentID     uint
	UpdateParam entity.CreateOrUpdateAgentParam
}

type UpdateAgentOutput struct {
	Agent *entity.Agent
}

func (i *AgentInteractorImpl) UpdateAgent(input UpdateAgentInput) (UpdateAgentOutput, error) {
	var (
		output UpdateAgentOutput
	)

	agent := entity.NewAgent(
		input.UpdateParam.AgentName,
		input.UpdateParam.OfficeLocation,
		input.UpdateParam.Representative,
		input.UpdateParam.Establish,
		input.UpdateParam.CorporateSiteURL,
		input.UpdateParam.PermissionCode,
		input.UpdateParam.PermissionYear,
		input.UpdateParam.WorkersCount,
		input.UpdateParam.PhoneNumber,
		input.UpdateParam.InterviewAdjustmentEmail,
		input.UpdateParam.AgreementFileURL,
		"",                                                // 引数に入れてるがUpdateされない input.UpdateParam.LineBotID,
		input.UpdateParam.LineMessagingChannelSecret,      // 引数に入れてるがUpdateされない
		input.UpdateParam.LineMessagingChannelAccessToken, // 引数に入れてるがUpdateされない
		input.UpdateParam.LineLoginChannelID,              // 引数に入れてるがUpdateされない
		input.UpdateParam.LineLoginChannelSecret,          // 引数に入れてるがUpdateされない
		"",                                                // input.UpdateParam.SendingAgreementFileURL,
		input.UpdateParam.IsCRMActive,                     // 引数に入れてるがUpdateされない
		input.UpdateParam.IsAllianceActive,                // 引数に入れてるがUpdateされない
		input.UpdateParam.IsSendingActive,                 // 引数に入れてるがUpdateされない
		input.UpdateParam.SendingType,
	)

	err := i.agentRepository.Update(input.AgentID, agent)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.Agent = agent

	return output, nil
}

type UpdateAgentAgreementFileURLInput struct {
	Param entity.AgentAgreementFileURLParam
}

type UpdateAgentAgreementFileURLOutput struct {
	OK bool
}

func (i *AgentInteractorImpl) UpdateAgentAgreementFileURL(input UpdateAgentAgreementFileURLInput) (UpdateAgentAgreementFileURLOutput, error) {
	var (
		output UpdateAgentAgreementFileURLOutput
	)

	err := i.agentRepository.UpdateAgreementFileURL(input.Param.AgentID, input.Param)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

type UpdateAgentForAdminInput struct {
	Param entity.AgentForAdminParam
}

type UpdateAgentForAdminOutput struct {
	OK bool
}

func (i *AgentInteractorImpl) UpdateAgentForAdmin(input UpdateAgentForAdminInput) (UpdateAgentForAdminOutput, error) {
	var (
		output UpdateAgentForAdminOutput
	)

	err := i.agentRepository.UpdateForAdmin(input.Param.AgentID, input.Param)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.OK = true

	return output, nil
}

/****************************************************************************************/

type UpdateAgentLineChannelInput struct {
	Param entity.AgentLineChannelParam
}

type UpdateAgentLineChannelOutput struct {
	AgentLine *entity.AgentLineChannelParam
}

func (i *AgentInteractorImpl) UpdateAgentLineChannel(input UpdateAgentLineChannelInput) (UpdateAgentLineChannelOutput, error) {
	var (
		output UpdateAgentLineChannelOutput
		bot    *entity.Bot
	)

	// BotIDを取得
	client := &http.Client{}

	// ボット情報を取得
	req, err := http.NewRequest("GET", "https://api.line.me/v2/bot/info", nil)
	if err != nil {
		fmt.Println(err)
	}

	tokenText := fmt.Sprintf("Bearer %s", input.Param.LineMessagingChannelAccessToken)
	req.Header.Set("Authorization", tokenText)

	// 実行
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(bodyText, &bot)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("レスポンス")
	fmt.Printf("%s\n", bot)

	if bot.BotID == "" {
		return output, errors.New("LINE Messaging APIのチャンネルアクセストークンが正しくありません")
	}

	err = i.agentRepository.UpdateLineChannel(input.Param.AgentID, bot.BotID, input.Param)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentLine = &input.Param

	return output, nil
}

type GetAgentLineChannelByAgentIDInput struct {
	AgentID uint
}

type GetAgentLineChannelByAgentIDOutput struct {
	AgentLine *entity.AgentLineChannelParam
}

func (i *AgentInteractorImpl) GetAgentLineChannelByAgentID(input GetAgentLineChannelByAgentIDInput) (GetAgentLineChannelByAgentIDOutput, error) {
	var (
		output GetAgentLineChannelByAgentIDOutput
	)

	agentLine, err := i.agentRepository.FindLineChannelByAgentID(input.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.AgentLine = agentLine

	return output, nil
}

// LINE Login Channel IDのみ取得
type GetAgentLineLoginChannelIDByAgentUUIDInput struct {
	AgentUUID uuid.UUID
}

type GetAgentLineLoginChannelIDByAgentUUIDOutput struct {
	LineLoginChannelID string
}

func (i *AgentInteractorImpl) GetAgentLineLoginChannelIDByAgentUUID(input GetAgentLineLoginChannelIDByAgentUUIDInput) (GetAgentLineLoginChannelIDByAgentUUIDOutput, error) {
	var (
		output GetAgentLineLoginChannelIDByAgentUUIDOutput
	)

	agentLine, err := i.agentRepository.FindLineChannelByAgentUUID(input.AgentUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	output.LineLoginChannelID = agentLine.LineLoginChannelID

	return output, nil
}
