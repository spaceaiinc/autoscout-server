package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type AgentStaffSaleManagementRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewAgentStaffSaleManagementRepositoryImpl(ex interfaces.SQLExecuter) usecase.AgentStaffSaleManagementRepository {
	return &AgentStaffSaleManagementRepositoryImpl{
		Name:     "AgentStaffSaleManagementRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *AgentStaffSaleManagementRepositoryImpl) Create(agentStaffSaleManagement *entity.AgentStaffSaleManagement) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO agent_staff_sale_managements (
			management_id,
			agent_staff_id,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?
			)`,
		agentStaffSaleManagement.ManagementID,
		agentStaffSaleManagement.AgentStaffID,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	agentStaffSaleManagement.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *AgentStaffSaleManagementRepositoryImpl) FindByManagementIDAndStaffID(managementID, staffID uint) (*entity.AgentStaffSaleManagement, error) {
	var (
		agentStaffSaleManagement entity.AgentStaffSaleManagement
	)

	err := repo.executer.Get(
		repo.Name+".FindByManagementIDAndStaffID",
		&agentStaffSaleManagement, `
		SELECT 
			ssm.*, DATE_FORMAT(asm.fiscal_year, '%Y-%m') AS fiscal_year, 
			staff.staff_name, staff.agent_id
		FROM 
			agent_staff_sale_managements AS ssm
		INNER JOIN
			agent_sale_managements AS asm
		ON
			ssm.management_id = asm.id
		INNER JOIN
			agent_staffs AS staff
		ON
			ssm.agent_staff_id = staff.id
		WHERE
			management_id = ?
		AND
			agent_staff_id = ?
		ORDER BY id DESC
		LIMIT 1
		`,
		managementID, staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &agentStaffSaleManagement, nil
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *AgentStaffSaleManagementRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.AgentStaffSaleManagement, error) {
	var (
		list []*entity.AgentStaffSaleManagement
	)
	err := repo.executer.Select(
		repo.Name+".GetByManagementID",
		&list, `
		SELECT 
			ssm.*, DATE_FORMAT(asm.fiscal_year, '%Y-%m') AS fiscal_year, staff.staff_name
		FROM 
			agent_staff_sale_managements AS ssm
		INNER JOIN
			agent_sale_managements AS asm
		ON
			ssm.management_id = asm.id
		INNER JOIN
			agent_staffs AS staff
		ON
			ssm.agent_staff_id = staff.id
		WHERE
			asm.agent_id = ?
	`, agentID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

func (repo *AgentStaffSaleManagementRepositoryImpl) GetStaffNameByManagementID(managementID uint) ([]*entity.AgentStaffSaleManagement, error) {
	var (
		list []*entity.AgentStaffSaleManagement
	)

	err := repo.executer.Select(
		repo.Name+".GetStaffNameByManagementID",
		&list, `
		SELECT 
			ssm.*, DATE_FORMAT(asm.fiscal_year, '%Y-%m') AS fiscal_year, staff.staff_name
		FROM 
			agent_staff_sale_managements AS ssm
		INNER JOIN
			agent_sale_managements AS asm
		ON
			ssm.management_id = asm.id
		INNER JOIN
			agent_staffs AS staff
		ON
			ssm.agent_staff_id = staff.id
		WHERE
			ssm.management_id = ?
	`, managementID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

// 決算期間が当月を含むレコードを全て取得
// 決算月以内でかつ、決算月の11ヶ月前以上のレコードを全て取得
func (repo *AgentStaffSaleManagementRepositoryImpl) GetByStaffIDAndThidMonth(staffID uint, thisMonth string) ([]*entity.AgentStaffSaleManagement, error) {
	var (
		list []*entity.AgentStaffSaleManagement
	)
	err := repo.executer.Select(
		repo.Name+".GetByStaffIDAndThidMonth",
		&list, `
		SELECT 
			ssm.*, DATE_FORMAT(asm.fiscal_year, '%Y-%m') AS fiscal_year, asm.is_open
		FROM 
			agent_staff_sale_managements AS ssm
		INNER JOIN
			agent_sale_managements AS asm
		ON
			ssm.management_id = asm.id
		WHERE
			ssm.agent_staff_id = ?
		AND (
			asm.fiscal_year >= ? AND
			DATE_SUB(asm.fiscal_year, INTERVAL 11 MONTH) <= ?
		)
	`, staffID, thisMonth, thisMonth)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}
