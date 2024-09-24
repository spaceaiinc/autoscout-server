package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationWorkCharmPointRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationWorkCharmPointRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationWorkCharmPointRepository {
	return &JobInformationWorkCharmPointRepositoryImpl{
		Name:     "JobInformationWorkCharmPointRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 求人の魅力を作成する
func (repo *JobInformationWorkCharmPointRepositoryImpl) Create(workCharmPoint *entity.JobInformationWorkCharmPoint) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_work_charm_points (
				job_information_id,
				title,
				contents,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		workCharmPoint.JobInformationID,
		workCharmPoint.Title,
		workCharmPoint.Contents,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	workCharmPoint.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 求人の魅力を削除する
func (repo *JobInformationWorkCharmPointRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_work_charm_points
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
// 指定求人IDの求人の魅力を取得する
func (repo *JobInformationWorkCharmPointRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.JobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&workCharmPointList, `
		SELECT *
		FROM job_information_work_charm_points
		WHERE
			job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return workCharmPointList, err
	}

	return workCharmPointList, nil
}

// 指定請求IDの仕事の魅力を取得
func (repo *JobInformationWorkCharmPointRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.JobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&workCharmPointList, `
			SELECT *
			FROM job_information_work_charm_points
			WHERE
				job_information_id IN (
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

	return workCharmPointList, nil
}

// 指定企業IDの仕事の魅力を取得
func (repo *JobInformationWorkCharmPointRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.JobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&workCharmPointList, `
			SELECT *
			FROM job_information_work_charm_points
			WHERE
				job_information_id IN (
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

	return workCharmPointList, nil
}

// 指定エージェントIDの仕事の魅力を取得
func (repo *JobInformationWorkCharmPointRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.JobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&workCharmPointList, `
			SELECT *
			FROM job_information_work_charm_points
			WHERE
				job_information_id IN (
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

	return workCharmPointList, nil
}

// 求人リストから仕事の魅力を取得
func (repo *JobInformationWorkCharmPointRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.JobInformationWorkCharmPoint
	)

	if len(jobInformationIDList) == 0 {
		return workCharmPointList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_work_charm_points
		WHERE
			job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&workCharmPointList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workCharmPointList, nil
}

func (repo *JobInformationWorkCharmPointRepositoryImpl) All() ([]*entity.JobInformationWorkCharmPoint, error) {
	var (
		workCharmPointList []*entity.JobInformationWorkCharmPoint
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&workCharmPointList, `
			SELECT *
			FROM job_information_work_charm_points
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return workCharmPointList, nil
}
