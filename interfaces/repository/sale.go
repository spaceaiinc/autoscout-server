package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SaleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSaleRepositoryImpl(ex interfaces.SQLExecuter) usecase.SaleRepository {
	return &SaleRepositoryImpl{
		Name:     "SaleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// 売上の作成
func (repo *SaleRepositoryImpl) Create(sale *entity.Sale) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sales (
			job_seeker_id,
			job_information_id,
			accuracy,
			contract_signed_month, 
			billing_month,
			billing_amount,
			cost, 
			gross_profit, 
			ra_staff_id, 
			ra_sales_ratio,
			ca_staff_id ,
			ca_sales_ratio,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
		sale.JobSeekerID,
		sale.JobInformationID,
		sale.Accuracy,
		sale.ContractSignedMonth,
		sale.BillingMonth,
		sale.BillingAmount,
		sale.Cost,
		sale.GrossProfit,
		sale.RAStaffID,
		sale.RaSalesRatio,
		sale.CAStaffID,
		sale.CaSalesRatio,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sale.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// 売上の更新
func (repo *SaleRepositoryImpl) Update(id uint, sale *entity.Sale) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sales
		SET
			job_seeker_id = ?,
			job_information_id = ?,
			accuracy = ?,
			contract_signed_month = ?, 
			billing_month = ?,
			billing_amount = ?,
			cost = ?, 
			gross_profit = ?, 
			ra_staff_id = ?, 
			ra_sales_ratio = ?,
			ca_staff_id  = ?,
			ca_sales_ratio = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		sale.JobSeekerID,
		sale.JobInformationID,
		sale.Accuracy,
		sale.ContractSignedMonth,
		sale.BillingMonth,
		sale.BillingAmount,
		sale.Cost,
		sale.GrossProfit,
		sale.RAStaffID,
		sale.RaSalesRatio,
		sale.CAStaffID,
		sale.CaSalesRatio,
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
/// 単数取得
//
func (repo *SaleRepositoryImpl) FindByID(id uint) (*entity.Sale, error) {
	var (
		sale entity.Sale
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&sale, `
		SELECT 
			sale.*
		FROM 
			sales AS sale
		WHERE
			sale.id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (repo *SaleRepositoryImpl) FindByJobSeekerID(seekerID uint) (*entity.Sale, error) {
	var (
		sale entity.Sale
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerID",
		&sale, `
		SELECT 
			sale.*, 
			ca_staff.staff_name AS ca_staff_name, ra_staff.staff_name AS ra_staff_name,
			ca_agent.agent_name AS ca_agent_name, ra_agent.agent_name AS ra_agent_name,
			ca_agent.id AS ca_agent_id, ra_agent.id AS ra_agent_id,
			task.phase_category, task.phase_sub_category,
			task_group.id AS task_group_id,
			seeker.last_name, seeker.first_name,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
				THEN task_group.external_job_information_title
				ELSE job_info.title
			END AS title,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name,
			task_group.external_job_listing_url
		FROM 
			sales AS sale
		INNER JOIN 
			agent_staffs AS ca_staff 
		ON 
			sale.ca_staff_id = ca_staff.id
		INNER JOIN 
			agent_staffs AS ra_staff 
		ON 
			sale.ra_staff_id = ra_staff.id
		INNER JOIN 
			agents AS ca_agent 
		ON 
			ca_staff.agent_id = ca_agent.id
		INNER JOIN 
			agents AS ra_agent 
		ON 
			ra_staff.agent_id = ra_agent.id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = sale.job_seeker_id
		INNER JOIN
			job_informations AS job_info
		ON
			job_info.id = sale.job_information_id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			task_groups AS task_group
		ON
			sale.job_seeker_id = task_group.job_seeker_id
		AND	
			sale.job_information_id = task_group.job_information_id
		LEFT OUTER JOIN (
			SELECT 
				task_origin.id, task_origin.task_group_id, 
				task_origin.phase_category, task_origin.phase_sub_category
			FROM tasks AS task_origin
			WHERE task_origin.id = (
				SELECT task_join.id
				FROM tasks AS task_join 
				WHERE task_origin.task_group_id = task_join.task_group_id
				ORDER BY task_join.id DESC
				LIMIT 1
			)
		) AS task 
		ON 
			task_group.id = task.task_group_id
		WHERE
			sale.job_seeker_id = ?
		LIMIT 1
		`,
		seekerID)

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (repo *SaleRepositoryImpl) FindByJobSeekerIDAndJobInformationID(seekerID, jobInfoID uint) (*entity.Sale, error) {
	var (
		sale entity.Sale
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerIDAndJobInformationID",
		&sale, `
		SELECT 
			sale.*, 
			ca_staff.staff_name AS ca_staff_name, ra_staff.staff_name AS ra_staff_name,
			ca_agent.agent_name AS ca_agent_name, ra_agent.agent_name AS ra_agent_name,
			ca_agent.id AS ca_agent_id, ra_agent.id AS ra_agent_id
		FROM 
			sales AS sale
		INNER JOIN 
			agent_staffs AS ca_staff 
		ON 
			sale.ca_staff_id = ca_staff.id
		INNER JOIN 
			agent_staffs AS ra_staff 
		ON 
			sale.ra_staff_id = ra_staff.id
		INNER JOIN 
			agents AS ca_agent 
		ON 
			ca_staff.agent_id = ca_agent.id
		INNER JOIN 
			agents AS ra_agent 
		ON 
			ra_staff.agent_id = ra_agent.id
		WHERE
			sale.job_seeker_id = ?
		AND
			sale.job_information_id = ?
		LIMIT 1
		`,
		seekerID, jobInfoID)

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

/****************************************************************************************/
/// 複数取得
//
// agentStaffIDを使ってヨミ情報の一覧を取得
func (repo *SaleRepositoryImpl) GetByAgentStaffID(agentStaffID uint) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffID",
		&saleList, `
			SELECT 
				sale.*
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			WHERE
				ca_staff.id = ? 
			OR
				ra_staff.id = ?
			ORDER BY 
				sale.id DESC
		`,
		agentStaffID, agentStaffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// 求人IDリストに合致する求人の一覧を取得
func (repo *SaleRepositoryImpl) GetByIDList(idList []uint) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	if len(idList) == 0 {
		return saleList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sale.*,
			ca_staff.staff_name AS ca_staff_name,
			ra_staff.staff_name AS ra_staff_name, 
			ca_agent.agent_name AS ca_agent_name, 
			ra_agent.agent_name AS ra_agent_name,
			ca_agent.id AS ca_agent_id, 
			ra_agent.id AS ra_agent_id,
			seeker.last_name, seeker.first_name,
			task.phase_category, task.phase_sub_category,
			task_group.id AS task_group_id,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
				THEN task_group.external_job_information_title
				ELSE job_info.title
			END AS title,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name,
			task_group.external_job_listing_url
		FROM 
			sales AS sale
		INNER JOIN 
			agent_staffs AS ca_staff 
		ON 
			sale.ca_staff_id = ca_staff.id
		INNER JOIN 
			agent_staffs AS ra_staff 
		ON 
			sale.ra_staff_id = ra_staff.id
		INNER JOIN 
			agents AS ca_agent 
		ON 
			ca_staff.agent_id = ca_agent.id
		INNER JOIN 
			agents AS ra_agent 
		ON 
			ra_staff.agent_id = ra_agent.id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = sale.job_seeker_id
		INNER JOIN
			job_informations AS job_info
		ON
			job_info.id = sale.job_information_id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			task_groups AS task_group
		ON
			sale.job_seeker_id = task_group.job_seeker_id
		AND	
			sale.job_information_id = task_group.job_information_id
		LEFT OUTER JOIN (
			SELECT 
				task_origin.id, task_origin.task_group_id, 
				task_origin.phase_category, task_origin.phase_sub_category
			FROM tasks AS task_origin
			WHERE task_origin.id = (
				SELECT task_join.id
				FROM tasks AS task_join 
				WHERE task_origin.task_group_id = task_join.task_group_id
				ORDER BY task_join.id DESC
				LIMIT 1
			)
		) AS task 
		ON 
			task_group.id = task.task_group_id
		WHERE
			sale.id IN(%s)
		ORDER BY 
			sale.accuracy ASC
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByIDList",
		&saleList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// ヨミ検索用のパラムを使って売上の一覧を取得
func (repo *SaleRepositoryImpl) GetSearchByAgent(searchParam entity.SearchAccuracy) ([]*entity.Sale, error) {
	var (
		saleList                    []*entity.Sale
		jobSeekerFreeWordQuery      string
		jobInformationFreeWordQuery string
		contractSignedMonthQuery    string
		billingMonthQuery           string
		raStaffQuery                string
		caStaffQuery                string
		accuraciesQuery             string
	)

	// 求職者名の条件
	if searchParam.JobSeekerFreeWord != "" {
		freeWordForLike := "%" + searchParam.JobSeekerFreeWord + "%"

		jobSeekerFreeWordQuery = fmt.Sprintf(`
			AND (	
				CONCAT(seeker.last_name, seeker.first_name) LIKE '%s' OR 
				CONCAT(seeker.last_furigana, seeker.first_furigana) LIKE '%s'
			)
		`, freeWordForLike, freeWordForLike)
	}

	// 求人タイトル・企業名の条件
	if searchParam.JobInformationFreeWord != "" {
		jobInformationFreeWordQuery = fmt.Sprintf(`
			AND (
				CASE
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN MATCH(task_group.external_company_name) AGAINST('%s' IN BOOLEAN MODE)
					ELSE MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE)
				END OR 
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN MATCH(task_group.external_job_information_title) AGAINST('%s' IN BOOLEAN MODE)
					ELSE MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE)
				END
			)
		`,
			searchParam.JobInformationFreeWord,
			searchParam.JobInformationFreeWord,
			searchParam.JobInformationFreeWord,
			searchParam.JobInformationFreeWord,
		)
	}

	// 受注月の条件
	if searchParam.ContractSignedMonth != "" {
		contractSignedMonthQuery = fmt.Sprintf(`
			AND sale.contract_signed_month = '%s'
		`, searchParam.ContractSignedMonth)
	}

	// 請求月の条件
	if searchParam.BillingMonth != "" {
		billingMonthQuery = fmt.Sprintf(`
			AND sale.billing_month = '%s'
		`, searchParam.BillingMonth)
	}

	// RA担当者の条件
	raStaffID, err := strconv.Atoi(searchParam.RAStaffID)
	if !(err != nil || raStaffID == 0) {
		raStaffQuery = fmt.Sprintf(`
			AND sale.ra_staff_id = %d
		`, raStaffID)
	}

	// CA担当者の条件
	caStaffID, err := strconv.Atoi(searchParam.CAStaffID)
	if !(err != nil || caStaffID == 0) {
		caStaffQuery = fmt.Sprintf(`
			AND sale.ca_staff_id = %d
		`, caStaffID)
	}

	// ヨミの条件
	if len(searchParam.Accuracies) > 0 {
		// 「"0, 1, 2, 3"」形式に変更
		accuraciesStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(searchParam.Accuracies)), ","), "[]")

		accuraciesQuery = fmt.Sprintf(`
			AND	sale.accuracy IN(%s)
		`, accuraciesStr)
	}

	query := fmt.Sprintf(`
		SELECT 
			sale.*,
			ca_staff.staff_name AS ca_staff_name, ra_staff.staff_name AS ra_staff_name,
			ca_agent.id AS ca_agent_id, ra_agent.id AS ra_agent_id,
			seeker.last_name, seeker.first_name,
			task.phase_category, task.phase_sub_category,
			task_group.id AS task_group_id,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
				THEN task_group.external_job_information_title
				ELSE job_info.title
			END AS title,
			CASE 
				WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
				THEN task_group.external_company_name
				ELSE enterprise.company_name
			END AS company_name,
			task_group.external_job_listing_url
		FROM 
			sales AS sale
		INNER JOIN 
			agent_staffs AS ca_staff 
		ON 
			sale.ca_staff_id = ca_staff.id
		INNER JOIN 
			agent_staffs AS ra_staff 
		ON 
			sale.ra_staff_id = ra_staff.id
		INNER JOIN 
			agents AS ca_agent 
		ON 
			ca_staff.agent_id = ca_agent.id
		INNER JOIN 
			agents AS ra_agent 
		ON 
			ra_staff.agent_id = ra_agent.id
		INNER JOIN
			job_seekers AS seeker
		ON
			seeker.id = sale.job_seeker_id
		INNER JOIN
			job_informations AS job_info
		ON
			job_info.id = sale.job_information_id
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			task_groups AS task_group
		ON
			sale.job_seeker_id = task_group.job_seeker_id
		AND	
			sale.job_information_id = task_group.job_information_id
		LEFT OUTER JOIN (
			SELECT 
				task_origin.id, task_origin.task_group_id, 
				task_origin.phase_category, task_origin.phase_sub_category
			FROM tasks AS task_origin
			WHERE task_origin.id = (
				SELECT task_join.id
				FROM tasks AS task_join 
				WHERE task_origin.task_group_id = task_join.task_group_id
				ORDER BY task_join.id DESC
				LIMIT 1
			)
		) AS task 
		ON 
			task_group.id = task.task_group_id
		WHERE (
			ca_staff.agent_id = ? OR
			ra_staff.agent_id = ?
		)
		%s %s %s %s %s %s %s
		ORDER BY 
			sale.id DESC
	`,
		jobSeekerFreeWordQuery,
		jobInformationFreeWordQuery,
		contractSignedMonthQuery,
		billingMonthQuery,
		raStaffQuery,
		caStaffQuery,
		accuraciesQuery,
	)

	err = repo.executer.Select(
		repo.Name+".GetSearchByAgent",
		&saleList,
		query,
		searchParam.AgentID, searchParam.AgentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentIDを使って売上の一覧を取得
func (repo *SaleRepositoryImpl) GetByAgentIDForMonthly(agentID uint, startMonth, endMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDForMonthly",
		&saleList, `
			SELECT 
				sale.*,
				ca_staff.agent_id AS ca_agent_id, 
				ra_staff.agent_id AS ra_agent_id
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			WHERE (
				ca_staff.agent_id = ? 
			OR
				ra_staff.agent_id = ?
			)
			AND (
				(
					sale.contract_signed_month >= ? AND 
					sale.contract_signed_month <= ? 
				) OR (
					sale.billing_month >= ? AND 
					sale.billing_month <= ?
				)
			)
			ORDER BY 
				sale.id DESC
		`,
		agentID, agentID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentStaffIDを使って売上の一覧を取得
func (repo *SaleRepositoryImpl) GetByStaffIDForMonthly(agentStaffID uint, startMonth, endMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetByStaffIDForMonthly",
		&saleList, `
			SELECT 
				sale.*
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			WHERE (
				ca_staff.id = ? 
			OR
				ra_staff.id = ?
			)
			AND (
				(
					sale.contract_signed_month >= ? AND 
					sale.contract_signed_month <= ?
				) OR (
					sale.billing_month >= ? AND 
					sale.billing_month <= ?
				)
			)
			ORDER BY 
				sale.id DESC
		`,
		agentStaffID, agentStaffID,
		startMonth, endMonth,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentStaffIDと当月の文字列を使ってヨミ情報の一覧を取得（受注）
func (repo *SaleRepositoryImpl) GetContractSignedByStaffIDAndMonth(agentStaffID uint, thisMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetContractSignedByStaffIDAndMonth",
		&saleList, `
			SELECT 
				sale.*,
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.id = ? 
			OR
				ra_staff.id = ?
			)
			AND 
				sale.contract_signed_month = ? 
			ORDER BY 
				sale.accuracy ASC
		`,
		agentStaffID, agentStaffID, thisMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentStaffIDと当月の文字列を使ってヨミ情報の一覧を取得（請求）
func (repo *SaleRepositoryImpl) GetBillingByStaffIDAndMonth(agentStaffID uint, thisMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetBillingByStaffIDAndMonth",
		&saleList, `
			SELECT 
				sale.*,
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.id = ? 
			OR
				ra_staff.id = ?
			)
			AND 
				sale.billing_month = ? 
			ORDER BY 
				sale.accuracy ASC
		`,
		agentStaffID, agentStaffID, thisMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentIDと指定月の文字列を使ってヨミ情報の一覧を取得（受注）
func (repo *SaleRepositoryImpl) GetContractSignedByAgentIDAndMonth(agentID uint, month string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetContractSignedByAgentIDAndMonth",
		&saleList, `
			SELECT 
				sale.*,
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.agent_id = ? 
			OR
				ra_staff.agent_id = ?
			)
			AND 
				sale.contract_signed_month = ? 
			ORDER BY 
				sale.accuracy ASC
		`,
		agentID, agentID, month,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentIDと指定月の文字列を使ってヨミ情報の一覧を取得（請求）
func (repo *SaleRepositoryImpl) GetBillingByAgentIDAndMonth(agentID uint, month string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetBillingByAgentIDAndMonth",
		&saleList, `
			SELECT 
				sale.*,
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.agent_id = ? 
			OR
				ra_staff.agent_id = ?
			)
			AND 
				sale.billing_month = ? 
			ORDER BY 
				sale.accuracy ASC
		`,
		agentID, agentID, month,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentIDと指定期間の文字列を使ってヨミ情報の一覧を取得（受注）
func (repo *SaleRepositoryImpl) GetContractSignedByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetContractSignedByAgentIDAndPeriod",
		&saleList, `
			SELECT 
				sale.*, 
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.agent_id = ? 
			OR
				ra_staff.agent_id = ?
			)
			AND (
				sale.contract_signed_month >= ? AND 
				sale.contract_signed_month <= ? 
			)
			ORDER BY 
				sale.accuracy ASC
		`,
		agentID, agentID,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// agentIDと指定期間の文字列を使ってヨミ情報の一覧を取得（請求）
func (repo *SaleRepositoryImpl) GetBillingByAgentIDAndPeriod(agentID uint, startMonth, endMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetBillingByAgentIDAndPeriod",
		&saleList, `
			SELECT 
				sale.*, 
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.agent_id = ? 
			OR
				ra_staff.agent_id = ?
			)
			AND (
				sale.billing_month >= ? AND 
				sale.billing_month <= ?
			)
			ORDER BY 
				sale.accuracy ASC
		`,
		agentID, agentID,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// staffIDと指定期間の文字列を使ってヨミ情報の一覧を取得（受注）
func (repo *SaleRepositoryImpl) GetContractSignedByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetContractSignedByStaffIDAndPeriod",
		&saleList, `
			SELECT 
				sale.*, 
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.id = ? 
			OR
				ra_staff.id = ?
			)
			AND (
				sale.contract_signed_month >= ? AND 
				sale.contract_signed_month <= ? 
			)
			ORDER BY 
				sale.accuracy ASC
		`,
		staffID, staffID,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

// staffIDと指定期間の文字列を使ってヨミ情報の一覧を取得（請求）
func (repo *SaleRepositoryImpl) GetBillingByStaffIDAndPeriod(staffID uint, startMonth, endMonth string) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetBillingByStaffIDAndPeriod",
		&saleList, `
			SELECT 
				sale.*, 
				ra_staff.staff_name AS ra_staff_name, 
				ca_staff.staff_name AS ca_staff_name,
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.id = ? 
			OR
				ra_staff.id = ?
			)
			AND (
				sale.billing_month >= ? AND 
				sale.billing_month <= ?
			)
			ORDER BY 
				sale.accuracy ASC
		`,
		staffID, staffID,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}

/******************************************************/
// CSV出力用
//
// agentIDからヨミ情報の一覧を取得（CSV出力用）
func (repo *SaleRepositoryImpl) GetByAgentIDForCSV(agentID uint) ([]*entity.Sale, error) {
	var (
		saleList []*entity.Sale
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDForCSV",
		&saleList, `
			SELECT 
				sale.*, 
				ca_staff.staff_name AS ca_staff_name,
				ra_staff.staff_name AS ra_staff_name, 
				ca_agent.agent_name AS ca_agent_name, 
				ra_agent.agent_name AS ra_agent_name,
				ca_agent.id AS ca_agent_id, 
				ra_agent.id AS ra_agent_id,
				seeker.last_name, seeker.first_name,
				task.phase_category, task.phase_sub_category,
				task_group.id AS task_group_id,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_job_information_title != '' 
					THEN task_group.external_job_information_title
					ELSE job_info.title
				END AS title,
				CASE 
					WHEN job_info.is_external = TRUE AND task_group.external_company_name != '' 
					THEN task_group.external_company_name
					ELSE enterprise.company_name
				END AS company_name,
				task_group.external_job_listing_url
			FROM 
				sales AS sale
			INNER JOIN 
				agent_staffs AS ca_staff 
			ON 
				sale.ca_staff_id = ca_staff.id
			INNER JOIN 
				agent_staffs AS ra_staff 
			ON 
				sale.ra_staff_id = ra_staff.id
			INNER JOIN 
				agents AS ca_agent 
			ON 
				ca_staff.agent_id = ca_agent.id
			INNER JOIN 
				agents AS ra_agent 
			ON 
				ra_staff.agent_id = ra_agent.id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = sale.job_seeker_id
			INNER JOIN
				job_informations AS job_info
			ON
				job_info.id = sale.job_information_id
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups AS task_group
			ON
				sale.job_seeker_id = task_group.job_seeker_id
			AND	
				sale.job_information_id = task_group.job_information_id
			LEFT OUTER JOIN (
				SELECT 
					task_origin.id, task_origin.task_group_id, 
					task_origin.phase_category, task_origin.phase_sub_category
				FROM tasks AS task_origin
				WHERE task_origin.id = (
					SELECT task_join.id
					FROM tasks AS task_join 
					WHERE task_origin.task_group_id = task_join.task_group_id
					ORDER BY task_join.id DESC
					LIMIT 1
				)
			) AS task 
			ON 
				task_group.id = task.task_group_id
			WHERE (
				ca_staff.agent_id = ? 
			OR
				ra_staff.agent_id = ?
			)
			ORDER BY 
				sale.accuracy ASC
		`,
		agentID, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return saleList, nil
}
