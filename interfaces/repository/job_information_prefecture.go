package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationPrefectureRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationPrefectureRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationPrefectureRepository {
	return &JobInformationPrefectureRepositoryImpl{
		Name:     "JobInformationPrefectureRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 都道府県を作成
func (repo *JobInformationPrefectureRepositoryImpl) Create(prefecture *entity.JobInformationPrefecture) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_prefectures (
				job_information_id,
				prefecture,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		prefecture.JobInformationID,
		prefecture.Prefecture,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	prefecture.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 指定求人IDの都道府県を削除
func (repo *JobInformationPrefectureRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_prefectures
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
func (repo *JobInformationPrefectureRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationPrefecture, error) {
	var (
		prefectureList []*entity.JobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationID",
		&prefectureList, `
		SELECT *
		FROM job_information_prefectures
		WHERE
			job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return prefectureList, err
	}

	return prefectureList, nil
}

func (repo *JobInformationPrefectureRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationPrefecture, error) {
	var (
		prefectureList []*entity.JobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&prefectureList, `
			SELECT *
			FROM job_information_prefectures
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

	return prefectureList, nil
}

func (repo *JobInformationPrefectureRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationPrefecture, error) {
	var (
		prefectureList []*entity.JobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&prefectureList, `
			SELECT *
			FROM job_information_prefectures
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

	return prefectureList, nil
}

func (repo *JobInformationPrefectureRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationPrefecture, error) {
	var (
		prefectureList []*entity.JobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&prefectureList, `
			SELECT *
			FROM job_information_prefectures
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

	return prefectureList, nil
}

// 求人リストから都道府県を取得
func (repo *JobInformationPrefectureRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationPrefecture, error) {
	var (
		prefectureList []*entity.JobInformationPrefecture
	)

	if len(jobInformationIDList) == 0 {
		return prefectureList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_prefectures
		WHERE
			job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&prefectureList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return prefectureList, nil
}

func (repo *JobInformationPrefectureRepositoryImpl) All() ([]*entity.JobInformationPrefecture, error) {
	var (
		prefectureList []*entity.JobInformationPrefecture
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&prefectureList, `
			SELECT *
			FROM job_information_prefectures
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return prefectureList, nil
}
