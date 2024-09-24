package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerLicenseRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerLicenseRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerLicenseRepository {
	return &JobSeekerLicenseRepositoryImpl{
		Name:     "JobSeekerLicenseRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerLicenseRepositoryImpl) Create(license *entity.JobSeekerLicense) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_licenses (
				job_seeker_id,
				license_type,
				acquisition_time,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		license.JobSeekerID,
		license.LicenseType,
		license.AcquisitionTime,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	license.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerLicenseRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_licenses
		WHERE job_seeker_id = ?
		`, jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得
//
func (repo *JobSeekerLicenseRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerLicense, error) {
	var licenseList []*entity.JobSeekerLicense

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&licenseList, `
		SELECT *
		FROM job_seeker_licenses
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return licenseList, err
	}

	return licenseList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerLicenseRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerLicense, error) {
	var licenseList []*entity.JobSeekerLicense

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&licenseList, `
			SELECT 
				jsl.*
			FROM 
				job_seeker_licenses AS jsl
			INNER JOIN
				job_seekers AS js
			ON
				jsl.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return licenseList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerLicenseRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerLicense, error) {
	var licenseList []*entity.JobSeekerLicense

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&licenseList, `
			SELECT 
				jsl.*
			FROM 
				job_seeker_licenses AS jsl
			INNER JOIN
				job_seekers AS js
			ON
				jsl.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return licenseList, nil
}

// 求職者リストから所持資格を取得
func (repo *JobSeekerLicenseRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerLicense, error) {
	var licenseList []*entity.JobSeekerLicense

	if len(idList) == 0 {
		return licenseList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_licenses
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&licenseList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return licenseList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerLicenseRepositoryImpl) All() ([]*entity.JobSeekerLicense, error) {
	var licenseList []*entity.JobSeekerLicense

	err := repo.executer.Select(
		repo.Name+".All",
		&licenseList, `
							SELECT *
							FROM job_seeker_licenses
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return licenseList, nil
}
