package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationHideToAgentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationHideToAgentRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationHideToAgentRepository {
	return &JobInformationHideToAgentRepositoryImpl{
		Name:     "JobInformationHideToAgentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 非公開エージェントを作成する
func (repo *JobInformationHideToAgentRepositoryImpl) Create(hideToAgent *entity.JobInformationHideToAgent) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_hide_to_agents (
				job_information_id,
				agent_id,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		hideToAgent.JobInformationID,
		hideToAgent.AgentID,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	hideToAgent.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 非公開エージェントを削除する
func (repo *JobInformationHideToAgentRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_hide_to_agents
		WHERE job_information_id = ?
		`, jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
// 複数取得 API
//
// 指定求人IDの非公開エージェントを取得する
func (repo *JobInformationHideToAgentRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&hideToAgentList, `
		SELECT 
			hide.*, agent.agent_name
		FROM 
			job_information_hide_to_agents AS hide
		INNER JOIN 
			agents AS agent
		ON
			hide.agent_id = agent.id
		WHERE
			hide.job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return hideToAgentList, err
	}

	return hideToAgentList, nil
}

func (repo *JobInformationHideToAgentRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_information_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
			WHERE
				hide.job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE
					billing_address_id = ?
				)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

func (repo *JobInformationHideToAgentRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_information_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
			WHERE
				hide.job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE 
						billing_address_id IN (
							SELECT id
							FROM billing_addresses
							WHERE enterprise_id = ?
						) 
				)
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

func (repo *JobInformationHideToAgentRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_information_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
			WHERE
				hide.job_information_id IN (
					SELECT id
					FROM job_informations
					WHERE billing_address_id IN (
						SELECT id
						FROM billing_addresses
						WHERE enterprise_id IN (
							SELECT id
							FROM enterprise_profiles
							WHERE agent_staff_id IN (
								SELECT id
								FROM agent_staffs
								WHERE agent_id = ?
							)
						)
					) 
				)
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

func (repo *JobInformationHideToAgentRepositoryImpl) GetHideByAgentID(agentID uint) ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	err := repo.executer.Select(
		repo.Name+".GetHideByAgentID",
		&hideToAgentList, `
		SELECT 
			hide.*, agent.agent_name
		FROM 
			job_information_hide_to_agents AS hide
		INNER JOIN 
			agents AS agent
		ON
			hide.agent_id = agent.id
		WHERE
			hide.agent_id = ?
		`, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

// 求人IDのリストを使ってデータ取得
func (repo *JobInformationHideToAgentRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	if len(jobInformationIDList) == 0 {
		return hideToAgentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			hide.*, agent.agent_name
		FROM 
			job_information_hide_to_agents AS hide
		INNER JOIN 
			agents AS agent
		ON
			hide.agent_id = agent.id
		WHERE
			hide.job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&hideToAgentList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

// エージェントIDのリストを使ってデータ取得
func (repo *JobInformationHideToAgentRepositoryImpl) GetByAgentIDList(agentIDList []uint) ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	if len(agentIDList) == 0 {
		return hideToAgentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			hide.*, agent.agent_name
		FROM 
			job_information_hide_to_agents AS hide
		INNER JOIN 
			agents AS agent
		ON
			hide.agent_id = agent.id
		WHERE
			hide.agent_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(agentIDList)), ","), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDList",
		&hideToAgentList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

func (repo *JobInformationHideToAgentRepositoryImpl) All() ([]*entity.JobInformationHideToAgent, error) {
	var (
		hideToAgentList []*entity.JobInformationHideToAgent
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_information_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}
