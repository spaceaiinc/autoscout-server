package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingSaleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingSaleRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingSaleRepository {
	return &SendingSaleRepositoryImpl{
		Name:     "SendingSaleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
// 売上の作成
func (repo *SendingSaleRepositoryImpl) Create(sendingSale *entity.SendingSale) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sending_sales (
			sending_job_seeker_id,
			sending_enterprise_id,
			system_sales,
			kickback,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)`,
		sendingSale.SendingJobSeekerID,
		sendingSale.SendingEnterpriseID,
		sendingSale.SystemSales,
		sendingSale.Kickback,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sendingSale.ID = uint(lastID)
	return nil
}

func (repo *SendingSaleRepositoryImpl) Update(sendingSaleID uint, sendingSale *entity.SendingSale) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sending_sales
		SET
			system_sales = ?,
			kickback = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		sendingSale.SystemSales,
		sendingSale.Kickback,
		time.Now().In(time.UTC),
		sendingSaleID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingSaleRepositoryImpl) FindByID(sendingSaleID uint) (*entity.SendingSale, error) {
	var (
		sendingSale entity.SendingSale
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&sendingSale, `
		SELECT 
			sale.*,
			phase.sending_date,
			CONCAT(seeker.last_name, seeker.first_name) AS job_seeker_name,
			enterprise.company_name AS receiver_agent_name,
			sender_agent.agent_name AS sender_agent_name
		FROM 
			sending_sales AS sale
		INNER JOIN
			sending_phases AS phase
		ON
			sale.sending_enterprise_id = phase.sending_enterprise_id AND sale.sending_job_seeker_id = phase.sending_job_seeker_id
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			sale.sending_job_seeker_id = seeker.id
		INNER JOIN
			sending_customers AS customer
		ON
			seeker.sending_customer_id = customer.id
		INNER JOIN
			agents AS sender_agent
		ON
			customer.agent_id = sender_agent.id
		INNER JOIN
			sending_enterprises AS enterprise
		ON
			sale.sending_enterprise_id = enterprise.id
		WHERE
			sale.id = ?
		LIMIT 1
		`,
		sendingSaleID)

	if err != nil {
		return nil, err
	}

	return &sendingSale, nil
}

func (repo *SendingSaleRepositoryImpl) FindByJobSeekerIDAndEnterpriseID(jobSeekerID, enterpriseID uint) (*entity.SendingSale, error) {
	var (
		sendingSale entity.SendingSale
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerIDAndEnterpriseID",
		&sendingSale, `
		SELECT 
			sale.*
		FROM 
			sending_sales AS sale
		WHERE
			sale.sending_job_seeker_id = ?
		AND
			sale.sending_enterprise_id = ?
		LIMIT 1
		`,
		jobSeekerID, enterpriseID)

	if err != nil {
		return nil, err
	}

	return &sendingSale, nil
}

// agentIDを使って売上の一覧を取得
func (repo *SendingSaleRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingSale, error) {
	var (
		sendingSaleList []*entity.SendingSale
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&sendingSaleList, `
			SELECT 
				sale.*,
				phase.sending_date,
				CONCAT(seeker.last_name, seeker.first_name) AS job_seeker_name
			FROM 
				sending_sales AS sale
			INNER JOIN
				sending_phases AS phase
			ON
				sale.sending_enterprise_id = phase.sending_enterprise_id AND sale.sending_job_seeker_id = phase.sending_job_seeker_id
			INNER JOIN 
				sending_job_seekers AS seeker
			ON 
				sale.sending_job_seeker_id = seeker.id
			INNER JOIN 
				agents AS agent
			ON 
				seeker.agent_id = agent.id
			WHERE
				agent.id = ?
			ORDER BY 
				phase.sending_date ASC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingSaleList, nil
}

// 送客元エージェントのIDを使って売上の一覧を取得（sending_customerテーブルのagent_idに合致するレコード）
func (repo *SendingSaleRepositoryImpl) GetListBySenderAgentIDForMonthly(senderAgentID uint, startMonth, endMonth string) ([]*entity.SendingSale, error) {
	var (
		sendingSaleList []*entity.SendingSale
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySenderAgentIDForMonthly",
		&sendingSaleList, `
			SELECT 
				sale.*,
				phase.sending_date,
				CONCAT(seeker.last_name, seeker.first_name) AS job_seeker_name
			FROM 
				sending_sales AS sale
			INNER JOIN
				sending_phases AS phase
			ON
				sale.sending_enterprise_id = phase.sending_enterprise_id AND sale.sending_job_seeker_id = phase.sending_job_seeker_id
			INNER JOIN 
				sending_job_seekers AS seeker
			ON 
				sale.sending_job_seeker_id = seeker.id
			INNER JOIN
				sending_customers AS customer
			ON
				seeker.sending_customer_id = customer.id
			WHERE
				customer.agent_id = ?
			AND 
			  DATE_FORMAT(phase.sending_date, '%Y-%m') BETWEEN ? AND ?
			ORDER BY 
				phase.sending_date ASC
		`,
		senderAgentID,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingSaleList, nil
}

// 送客管理エージェントのIDを使って売上の一覧を取得（sending_job_seekerテーブルのagent_idに合致するレコード）
func (repo *SendingSaleRepositoryImpl) GetListByAgentIDForMonthly(senderAgentID uint, startMonth, endMonth string) ([]*entity.SendingSale, error) {
	var (
		sendingSaleList []*entity.SendingSale
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentIDForMonthly",
		&sendingSaleList, `
			SELECT 
				sale.*,
				phase.sending_date,
				CONCAT(seeker.last_name, seeker.first_name) AS job_seeker_name,
				enterprise.company_name AS receiver_agent_name
			FROM 
				sending_sales AS sale
			INNER JOIN
				sending_phases AS phase
			ON
				sale.sending_enterprise_id = phase.sending_enterprise_id AND sale.sending_job_seeker_id = phase.sending_job_seeker_id
			INNER JOIN 
				sending_job_seekers AS seeker
			ON 
				sale.sending_job_seeker_id = seeker.id
			INNER JOIN
				sending_enterprises AS enterprise
			ON
				sale.sending_enterprise_id = enterprise.id
			WHERE
				seeker.agent_id = ?
			AND 
			  DATE_FORMAT(phase.sending_date, '%Y-%m') BETWEEN ? AND ?
			ORDER BY 
				phase.sending_date ASC
		`,
		senderAgentID,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingSaleList, nil
}

/****************************************************************************************/
/// Admin API
//
// 全ての売上の一覧を取得
func (repo *SendingSaleRepositoryImpl) GetAll() ([]*entity.SendingSale, error) {
	var (
		sendingSaleList []*entity.SendingSale
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&sendingSaleList, `
			SELECT 
				sale.*,
				phase.sending_date,
				CONCAT(seeker.last_name, seeker.first_name) AS job_seeker_name,
				enterprise.company_name AS receiver_agent_name
			FROM 
				sending_sales AS sale
			INNER JOIN
				sending_phases AS phase
			ON
				sale.sending_enterprise_id = phase.sending_enterprise_id AND sale.sending_job_seeker_id = phase.sending_job_seeker_id
			INNER JOIN
				sending_job_seekers AS seeker
			ON
				sale.sending_job_seeker_id = seeker.id
			INNER JOIN
				sending_enterprises AS enterprise
			ON
				sale.sending_enterprise_id = enterprise.id
			ORDER BY 
				phase.sending_date ASC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingSaleList, nil
}

// 指定した月範囲の全ての売上の一覧を取得
func (repo *SendingSaleRepositoryImpl) GetAllForMonthly(startMonth, endMonth string) ([]*entity.SendingSale, error) {
	var (
		sendingSaleList []*entity.SendingSale
	)

	err := repo.executer.Select(
		repo.Name+".GetAllForMonthly",
		&sendingSaleList, `
			SELECT 
				sale.*,
				phase.sending_date,
				CONCAT(seeker.last_name, seeker.first_name) AS job_seeker_name,
				enterprise.company_name AS receiver_agent_name,
				sender_agent.agent_name AS sender_agent_name
			FROM 
				sending_sales AS sale
			INNER JOIN
				sending_phases AS phase
			ON
				sale.sending_enterprise_id = phase.sending_enterprise_id AND sale.sending_job_seeker_id = phase.sending_job_seeker_id
			INNER JOIN
				sending_job_seekers AS seeker
			ON
				sale.sending_job_seeker_id = seeker.id
			INNER JOIN
				sending_customers AS customer
			ON
				seeker.sending_customer_id = customer.id
			INNER JOIN
				agents AS sender_agent
			ON
				customer.agent_id = sender_agent.id
			INNER JOIN
				sending_enterprises AS enterprise
			ON
				sale.sending_enterprise_id = enterprise.id
			WHERE
			  DATE_FORMAT(phase.sending_date, '%Y-%m') BETWEEN ? AND ?
			ORDER BY 
				phase.sending_date ASC
		`,
		startMonth, endMonth,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingSaleList, nil
}
