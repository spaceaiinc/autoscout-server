package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerLicenseRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerLicenseRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerLicenseRepository {
	return &SendingJobSeekerLicenseRepositoryImpl{
		Name:     "SendingJobSeekerLicenseRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerLicenseRepositoryImpl) Create(license *entity.SendingJobSeekerLicense) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_licenses (
				sending_job_seeker_id,
				license_type,
				acquisition_time,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		license.SendingJobSeekerID,
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

func (repo *SendingJobSeekerLicenseRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_licenses
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerLicenseRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerLicense, error) {
	var (
		licenseList []*entity.SendingJobSeekerLicense
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&licenseList, `
		SELECT *
		FROM sending_job_seeker_licenses
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return licenseList, err
	}

	return licenseList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerLicenseRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerLicense, error) {
	var (
		licenseList []*entity.SendingJobSeekerLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&licenseList, `
			SELECT 
				jsl.*
			FROM 
				sending_job_seeker_licenses AS jsl
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsl.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerLicenseRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerLicense, error) {
	var (
		licenseList []*entity.SendingJobSeekerLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&licenseList, `
			SELECT 
				jsl.*
			FROM 
				sending_job_seeker_licenses AS jsl
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsl.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerLicenseRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerLicense, error) {
	var (
		licenseList []*entity.SendingJobSeekerLicense
	)

	if len(idList) == 0 {
		return licenseList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_licenses
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&licenseList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return licenseList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerLicenseRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerLicense, error) {
	var (
		licenseList []*entity.SendingJobSeekerLicense
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&licenseList, `
							SELECT *
							FROM sending_job_seeker_licenses
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return licenseList, nil
}
