package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobInformationFeatureRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationFeatureRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationFeatureRepository {
	return &JobInformationFeatureRepositoryImpl{
		Name:     "JobInformationFeatureRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成 API
//
// 求人の特徴を作成する
func (repo *JobInformationFeatureRepositoryImpl) Create(feature *entity.JobInformationFeature) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_information_features (
				job_information_id,
				feature,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		feature.JobInformationID,
		feature.Feature,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	feature.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
// 削除 API
//
// 求人の特徴を削除する
func (repo *JobInformationFeatureRepositoryImpl) DeleteByJobInformationID(jobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobInformationID",
		`
		DELETE
		FROM job_information_features
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
// 指定求人IDの特徴を取得する
func (repo *JobInformationFeatureRepositoryImpl) GetByJobInformationID(jobInformationID uint) ([]*entity.JobInformationFeature, error) {
	var (
		featureList []*entity.JobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+"GetByJobInformationID",
		&featureList, `
		SELECT *
		FROM job_information_features
		WHERE
			job_information_id = ?
		`,
		jobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return featureList, err
	}

	return featureList, nil
}

// 指定請求先IDの特徴を取得する
func (repo *JobInformationFeatureRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformationFeature, error) {
	var (
		featureList []*entity.JobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&featureList, `
			SELECT *
			FROM job_information_features
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

	return featureList, nil
}

func (repo *JobInformationFeatureRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformationFeature, error) {
	var (
		featureList []*entity.JobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&featureList, `
			SELECT *
			FROM job_information_features
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

	return featureList, nil
}

func (repo *JobInformationFeatureRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformationFeature, error) {
	var (
		featureList []*entity.JobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&featureList, `
			SELECT *
			FROM job_information_features
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

	return featureList, nil
}

// 求人リストから特徴を取得
func (repo *JobInformationFeatureRepositoryImpl) GetByJobInformationIDList(jobInformationIDList []uint) ([]*entity.JobInformationFeature, error) {
	var (
		featureList []*entity.JobInformationFeature
	)

	if len(jobInformationIDList) == 0 {
		return featureList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_information_features
		WHERE
			job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDList",
		&featureList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return featureList, nil
}

// 求職者ページ表示用項目のみ取得 0:業界未経験OK, 1:職種未経験OK, 2:業界・職種未経験OK, 6:転勤なし
func (repo *JobInformationFeatureRepositoryImpl) GetByJobInformationIDListForGuestJobSeeker(jobInformationIDList []uint) ([]*entity.JobInformationFeature, error) {
	var (
		featureList []*entity.JobInformationFeature
	)

	if len(jobInformationIDList) == 0 {
		return featureList, nil
	}

	query := fmt.Sprintf(`
			SELECT *
			FROM job_information_features
			WHERE
				job_information_id IN (%s)
			AND 
			  feature IN (0, 1, 2, 6)
		`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobInformationIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobInformationIDListForGuestJobSeeker",
		&featureList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return featureList, nil
}

func (repo *JobInformationFeatureRepositoryImpl) All() ([]*entity.JobInformationFeature, error) {
	var (
		featureList []*entity.JobInformationFeature
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&featureList, `
			SELECT *
			FROM job_information_features
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return featureList, nil
}
