package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentSaleManagementRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentSaleManagementRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentSaleManagementRepository {
	return &AgentSaleManagementRepositoryImpl{
		Name:     "AgentSaleManagementRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *AgentSaleManagementRepositoryImpl) Create(agentSaleManagement *entity.AgentSaleManagement) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_sale_managements (
			agent_id,
			fiscal_year,
			is_open,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)`,
		agentSaleManagement.AgentID,
		agentSaleManagement.FiscalYear,
		agentSaleManagement.IsOpen,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	agentSaleManagement.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *AgentSaleManagementRepositoryImpl) Update(agentSaleManagement *entity.AgentSaleManagement, agentsaleManagementID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE agent_sale_managements
		SET
			agent_id = ?,
			fiscal_year = ?,
			is_open = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		agentSaleManagement.AgentID,
		agentSaleManagement.FiscalYear,
		agentSaleManagement.IsOpen,
		time.Now().In(time.UTC),
		agentsaleManagementID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *AgentSaleManagementRepositoryImpl) UpdateIsOpenOtherThanID(agentID, saleManagementID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsOpenOtherThanID",
		`
		UPDATE agent_sale_managements
		SET
			is_open = ?,
			updated_at = ?
		WHERE
			agent_id = ?
		AND
			id NOT IN(?)
		`,
		false,
		time.Now().In(time.UTC),
		agentID,
		saleManagementID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *AgentSaleManagementRepositoryImpl) FindByID(agentsaleManagementID uint) (*entity.AgentSaleManagement, error) {
	var (
		AgentSaleManagement entity.AgentSaleManagement
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&AgentSaleManagement, `
		SELECT 
			*, 
			DATE_FORMAT(fiscal_year, '%Y-%m') AS fiscal_year
		FROM 
			agent_sale_managements
		WHERE
			id = ?
		`, agentsaleManagementID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &AgentSaleManagement, nil
}

func (repo *AgentSaleManagementRepositoryImpl) FindByIDAndAgentID(managementID, agentID uint) (*entity.AgentSaleManagement, error) {
	var (
		AgentSaleManagement entity.AgentSaleManagement
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&AgentSaleManagement, `
		SELECT 
			*, 
			DATE_FORMAT(fiscal_year, '%Y-%m') AS fiscal_year
		FROM 
			agent_sale_managements
		WHERE
			id = ?
		AND
			agent_id = ?
		`,
		managementID,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &AgentSaleManagement, nil
}

func (repo *AgentSaleManagementRepositoryImpl) FindByAgentIDAndIsOpen(agentID uint) (*entity.AgentSaleManagement, error) {
	var (
		AgentSaleManagement entity.AgentSaleManagement
	)

	err := repo.executer.Get(
		repo.Name+".FindByAgentIDAndIsOpen",
		&AgentSaleManagement, `
		SELECT 
			*, DATE_FORMAT(fiscal_year, '%Y-%m') AS fiscal_year
		FROM 
			agent_sale_managements
		WHERE
			agent_id = ?
		AND
			is_open = true
		ORDER BY id DESC
		LIMIT 1
		`, agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &AgentSaleManagement, nil
}

func (repo *AgentSaleManagementRepositoryImpl) FindLatestByAgentID(agentID uint) (*entity.AgentSaleManagement, error) {
	var (
		AgentSaleManagement entity.AgentSaleManagement
	)

	err := repo.executer.Get(
		repo.Name+".FindLatestByAgentID",
		&AgentSaleManagement, `
		SELECT 
			*, DATE_FORMAT(fiscal_year, '%Y-%m') AS fiscal_year
		FROM 
			agent_sale_managements
		WHERE
			agent_id = ?
		ORDER BY id DESC
		LIMIT 1
		`, agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &AgentSaleManagement, nil
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *AgentSaleManagementRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.AgentSaleManagement, error) {
	var (
		list []*entity.AgentSaleManagement
	)
	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&list, `
		SELECT 
			*, DATE_FORMAT(fiscal_year, '%Y-%m') AS fiscal_year
		FROM 
			agent_sale_managements
		WHERE
			agent_id = ?
	`,
		agentID)

	if err != nil {
		return nil, err
	}

	return list, nil
}
