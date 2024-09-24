package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EnterpriseIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewEnterpriseIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.EnterpriseIndustryRepository {
	return &EnterpriseIndustryRepositoryImpl{
		Name:     "EnterpriseIndustryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *EnterpriseIndustryRepositoryImpl) Create(industry *entity.EnterpriseIndustry) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO enterprise_industries (
				enterprise_id,
				industry,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		industry.EnterpriseID,
		industry.Industry,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	industry.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除 API
//
func (repo *EnterpriseIndustryRepositoryImpl) DeleteByEnterpriseID(enterpriseID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByEnterpriseID",
		`
		DELETE
		FROM enterprise_industries
		WHERE enterprise_id = ?
		`, enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *EnterpriseIndustryRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&industryList, `
		SELECT *
		FROM enterprise_industries
		WHERE
			enterprise_id = ?
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return industryList, err
	}

	return industryList, nil
}

func (repo *EnterpriseIndustryRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&industryList, `
		SELECT *
		FROM enterprise_industries
		WHERE
			enterprise_id = (
				SELECT enterprise_id
				FROM billing_addresses
				WHERE id = ?
			)
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return industryList, err
	}

	return industryList, nil
}

// 求人IDから企業情報を取得(求人票作成用)
func (repo *EnterpriseIndustryRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&industryList, `
		SELECT 
			*
		FROM 
			enterprise_industries
		WHERE
			enterprise_id = (
				SELECT enterprise_id
				FROM billing_addresses
				WHERE id = (
					SELECT billing_address_id
					FROM job_informations
					WHERE id = ?
				)
			)
		`, jobInformationID)

	if err != nil {
		return nil, err
	}

	return industryList, nil
}

// 担当者IDで企業情報を取得
func (repo *EnterpriseIndustryRepositoryImpl) GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffID",
		&industryList, `
		SELECT *
		FROM 
		enterprise_industries
		WHERE
		enterprise_id IN (
			SELECT id
			FROM enterprise_profiles
			WHERE agent_staff_id = ?
		)
		`,
		agentStaffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}

// エージェントIDから企業一覧を取得
func (repo *EnterpriseIndustryRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&industryList, `
		SELECT *
		FROM 
		enterprise_industries
		WHERE
		enterprise_id IN (
			SELECT id
			FROM enterprise_profiles
			WHERE 
				agent_staff_id IN (
					SELECT id
					FROM agent_staffs
					WHERE
					agent_id = ?
					)
			)
						`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}

// IDリストから企業の業界を取得
func (repo *EnterpriseIndustryRepositoryImpl) GetByEnterpriseIDList(enterpriseIDList []uint) ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	if len(enterpriseIDList) < 1 {
		return nil, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM enterprise_industries
		WHERE enterprise_id IN (%s)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enterpriseIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseIDList",
		&industryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}

// 求人IDリストから企業の業界を取得
func (repo *EnterpriseIndustryRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	if len(jobInformationIDList) < 1 {
		return nil, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM enterprise_industries
		WHERE enterprise_id IN (
			SELECT enterprise_id
			FROM billing_addresses
			WHERE id IN (
				SELECT billing_address_id
				FROM job_informations
				WHERE id IN (%s)
			)
		)
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&industryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}

// すべての企業情報を取得
func (repo *EnterpriseIndustryRepositoryImpl) All() ([]*entity.EnterpriseIndustry, error) {
	var (
		industryList []*entity.EnterpriseIndustry
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&industryList, `
							SELECT *
							FROM enterprise_industries
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}
