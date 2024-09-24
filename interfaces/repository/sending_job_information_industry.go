package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationIndustryRepository {
	return &SendingJobInformationIndustryRepositoryImpl{
		Name:     "SendingJobInformationIndustryRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationIndustryRepositoryImpl) Create(industry *entity.SendingJobInformationIndustry) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_industries (
				sending_job_information_id,
				industry,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		industry.SendingJobInformationID,
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

func (repo *SendingJobInformationIndustryRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_industries
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationIndustryRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationIndustry, error) {
	var (
		industryList []*entity.SendingJobInformationIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&industryList, `
		SELECT *
		FROM sending_job_information_industries
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return industryList, err
	}

	return industryList, nil
}

func (repo *SendingJobInformationIndustryRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationIndustry, error) {
	var (
		industryList []*entity.SendingJobInformationIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&industryList, `
			SELECT *
			FROM sending_job_information_industries
			WHERE
				sending_job_information_id IN (
					SELECT id
					FROM sending_job_informations
					WHERE
					sending_billing_address_id = ?
				)
		`,
		sendingBillingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}

func (repo *SendingJobInformationIndustryRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationIndustry, error) {
	var (
		industryList []*entity.SendingJobInformationIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&industryList, `
			SELECT *
			FROM sending_job_information_industries
			WHERE
				sending_job_information_id IN (
					SELECT id
					FROM sending_job_informations
					WHERE 
						sending_billing_address_id IN (
							SELECT id
							FROM sending_billing_addresses
							WHERE sending_enterprise_id = ?
						) 
					)
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationIndustryRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationIndustry, error) {
	var (
		industryList []*entity.SendingJobInformationIndustry
	)

	if len(idList) == 0 {
		return industryList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_industries.*
		FROM 
			sending_job_information_industries
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&industryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationIndustryRepositoryImpl) GetAll() ([]*entity.SendingJobInformationIndustry, error) {
	var (
		industryList []*entity.SendingJobInformationIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&industryList, `
			SELECT *
			FROM sending_job_information_industries
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return industryList, nil
}
