package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentRepository {
	return &AgentRepositoryImpl{
		Name:     "AgentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// エージェントを作成
func (repo *AgentRepositoryImpl) Create(agent *entity.Agent) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agents (
			uuid,
			agent_name,
			office_location,
			representative,
			establish,
			corporate_site_url,
			permission_code,
			permission_year,
			workers_count,
			phone_number,
			interview_adjustment_email,
			agreement_file_url,
			line_bot_id,
			line_messaging_channel_secret,
			line_messaging_channel_access_token,
			line_loging_channel_id,
			line_login_channel_secret,
			sending_agreement_file_url,
			is_crm_active,
			is_alliance_active,
			is_sending_active,
			sending_type,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?,
			?, ?, ?, ?
		)
		`,
		utility.CreateUUID(),
		agent.AgentName,
		agent.OfficeLocation,
		agent.Representative,
		agent.Establish,
		agent.CorporateSiteURL,
		agent.PermissionCode,
		agent.PermissionYear,
		agent.WorkersCount,
		agent.PhoneNumber,
		agent.InterviewAdjustmentEmail,
		agent.AgreementFileURL,
		agent.LineBotID,
		agent.LineMessagingChannelSecret,
		agent.LineMessagingChannelAccessToken,
		agent.LineLoginChannelID,
		agent.LineLoginChannelSecret,
		"", // agent.SendingAgreementFileURL,
		agent.IsCRMActive,
		agent.IsAllianceActive,
		agent.IsSendingActive,
		agent.SendingType,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	agent.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// エージェントを更新
func (repo *AgentRepositoryImpl) Update(id uint, agent *entity.Agent) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE agents 
		SET
			agent_name = ?,
			office_location = ?,
			representative = ?,
			establish = ?,
			corporate_site_url = ?,
			permission_code = ?,
			permission_year = ?,
			workers_count = ?,
			phone_number = ?,
			interview_adjustment_email = ?,
			agreement_file_url = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		agent.AgentName,
		agent.OfficeLocation,
		agent.Representative,
		agent.Establish,
		agent.CorporateSiteURL,
		agent.PermissionCode,
		agent.PermissionYear,
		agent.WorkersCount,
		agent.PhoneNumber,
		agent.InterviewAdjustmentEmail,
		agent.AgreementFileURL,
		// agent.SendingAgreementFileURL, // 送客用同意書ファイル
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 管理側からCRMや送客機能の利用を更新する
func (repo *AgentRepositoryImpl) UpdateForAdmin(id uint, agent entity.AgentForAdminParam) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateForAdmin",
		`
			UPDATE agents 
			SET
				is_crm_active = ?,
				is_alliance_active = ?,
				is_sending_active = ?,
				sending_type = ?,
				updated_at = ?
			WHERE 
				id = ?
		`,
		agent.IsCRMActive,
		agent.IsAllianceActive,
		agent.IsSendingActive,
		agent.SendingType,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// LINE連携用の情報を登録する
func (repo *AgentRepositoryImpl) UpdateLineChannel(id uint, botID string, agentLineChannel entity.AgentLineChannelParam) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateLineChannel",
		`
		UPDATE agents 
		SET
			line_bot_id = ?,
			line_messaging_channel_secret = ?,
			line_messaging_channel_access_token = ?,
			line_loging_channel_id = ?,
			line_login_channel_secret = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		botID,
		agentLineChannel.LineMessagingChannelSecret,
		agentLineChannel.LineMessagingChannelAccessToken,
		agentLineChannel.LineLoginChannelID,
		agentLineChannel.LineLoginChannelSecret,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 同意書ファイルの更新
func (repo *AgentRepositoryImpl) UpdateAgreementFileURL(id uint, agentAgreement entity.AgentAgreementFileURLParam) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateAgreementFileURL",
		`
		UPDATE agents 
		SET
			agreement_file_url = ?,
			sending_agreement_file_url = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		agentAgreement.AgreementFileURL,
		agentAgreement.SendingAgreementFileURL,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
// 単数取得 API
//
func (repo *AgentRepositoryImpl) FindByID(id uint) (*entity.Agent, error) {
	var (
		agent entity.Agent
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&agent, `
		SELECT *
		FROM agents
		WHERE
			id = ?
		LIMIT 1
		`, id)
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (repo *AgentRepositoryImpl) FindByUUID(uuid uuid.UUID) (*entity.Agent, error) {
	var (
		agent entity.Agent
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&agent, `
		SELECT *
		FROM agents
		WHERE
			uuid = ?
		LIMIT 1
		`,
		uuid,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (repo *AgentRepositoryImpl) FindByAgentStaffID(agentStaffID uint) (*entity.Agent, error) {
	var (
		agent entity.Agent
	)

	err := repo.executer.Get(
		repo.Name+".FindByAgentStaffID",
		&agent, `
		SELECT *
		FROM agents
		WHERE
			id = (
				SELECT
					agent_id
				FROM
					agent_staffs
				WHERE
					id = ?
			)
		LIMIT 1
		`,
		agentStaffID,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (repo *AgentRepositoryImpl) FindLineChannelByAgentID(id uint) (*entity.AgentLineChannelParam, error) {
	var (
		agent entity.AgentLineChannelParam
	)

	err := repo.executer.Get(
		repo.Name+".FindLineChannelByAgentID",
		&agent, `
		SELECT 
			id AS agent_id, line_bot_id,
			line_messaging_channel_secret, line_messaging_channel_access_token,
			line_loging_channel_id, line_login_channel_secret
		FROM 
			agents
		WHERE
			id = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (repo *AgentRepositoryImpl) FindLineChannelByAgentUUID(agentUUID uuid.UUID) (*entity.AgentLineChannelParam, error) {
	var (
		agent entity.AgentLineChannelParam
	)

	err := repo.executer.Get(
		repo.Name+".FindLineChannelByAgentUUID",
		&agent, `
		SELECT 
			id AS agent_id, line_bot_id,
			line_messaging_channel_secret, line_messaging_channel_access_token,
			line_loging_channel_id, line_login_channel_secret
		FROM 
			agents
		WHERE
			uuid = ?
		LIMIT 1
		`,
		agentUUID,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (repo *AgentRepositoryImpl) FindLineChannelByBotID(botID string) (*entity.AgentLineChannelParam, error) {
	var (
		agent entity.AgentLineChannelParam
	)

	err := repo.executer.Get(
		repo.Name+".FindLineChannelByBotID",
		&agent, `
		SELECT 
			id AS agent_id, uuid AS agent_uuid, line_bot_id,
			line_messaging_channel_secret, line_messaging_channel_access_token,
			line_loging_channel_id, line_login_channel_secret
		FROM 
			agents
		WHERE
			line_bot_id = ?
		LIMIT 1
		`,
		botID,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

/****************************************************************************************/
// 複数取得 API
//
func (repo *AgentRepositoryImpl) GetByIDList(idList []uint) ([]*entity.Agent, error) {
	var (
		agentList []*entity.Agent
	)

	if len(idList) < 1 {
		return agentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			*
		FROM 
			agents
		WHERE 
			id IN (%s)
		ORDER BY id ASC
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByIDList",
		&agentList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentList, nil
}

func (repo *AgentRepositoryImpl) GetByNotIDList(idList []uint) ([]*entity.Agent, error) {
	var (
		agentList []*entity.Agent
	)

	if len(idList) < 1 {
		return agentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			*
		FROM 
			agents
		WHERE 
			id NOT IN (%s)
		ORDER BY id ASC
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByNotIDList",
		&agentList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentList, nil
}

func (repo *AgentRepositoryImpl) GetNotAllianceByMyAgentID(myAgentID uint) ([]*entity.Agent, error) {
	var (
		agentList []*entity.Agent
	)
	/**
	 * 「アライアンス申請のレコードがない」もしくは「レコードはあるが自分が未申請」
	 **/
	err := repo.executer.Select(
		repo.Name+".GetNotAllianceByMyAgentID",
		&agentList, `
			SELECT 
				agents.*
			FROM 
				agents
			LEFT OUTER JOIN
				agent_alliances AS alliance
			ON
				agents.id = alliance.agent1_id 
			OR 
				agents.id = alliance.agent2_id
			WHERE (
				alliance.id IS NULL OR (
					alliance.agent1_id = ? &&
					alliance.agent1_request = FALSE
				) OR (
					alliance.agent2_id = ? &&
					alliance.agent2_request = FALSE
				)
			)
			AND agents.id != ?
		`, myAgentID, myAgentID, myAgentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentList, nil
}

// 送客利用のエージェントをすべて取得
func (repo *AgentRepositoryImpl) GetSendingActive() ([]*entity.Agent, error) {
	var (
		agentList []*entity.Agent
	)

	err := repo.executer.Select(
		repo.Name+".GetSendingActive",
		&agentList, `
			SELECT *
			FROM agents
			WHERE is_sending_active = TRUE
			ORDER BY id ASC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agentList, nil
}

// 全てのエージェント一覧を取得する
func (repo *AgentRepositoryImpl) All() ([]*entity.Agent, error) {
	var (
		agenList []*entity.Agent
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&agenList, `
			SELECT *
			FROM agents
			ORDER BY id ASC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return agenList, nil
}
