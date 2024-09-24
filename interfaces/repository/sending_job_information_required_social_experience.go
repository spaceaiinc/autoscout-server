package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRequiredSocialExperienceRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRequiredSocialExperienceRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRequiredSocialExperienceRepository {
	return &SendingJobInformationRequiredSocialExperienceRepositoryImpl{
		Name:     "SendingJobInformationRequiredSocialExperienceRepository",
		executer: ex,
	}
}

func (repo *SendingJobInformationRequiredSocialExperienceRepositoryImpl) Create(requiredSocialExperience *entity.SendingJobInformationRequiredSocialExperience) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_information_required_social_experiences (
				sending_job_information_id,
				social_experience_type,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		requiredSocialExperience.SendingJobInformationID,
		requiredSocialExperience.SocialExperienceType,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	requiredSocialExperience.ID = uint(lastID)

	return nil
}

func (repo *SendingJobInformationRequiredSocialExperienceRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_information_required_social_experiences
		WHERE	sending_job_information_id = ?
		`, sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRequiredSocialExperienceRepositoryImpl) GetListBySendingJobInformationID(sendingJobInformationID uint) ([]*entity.SendingJobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.SendingJobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobInformationID",
		&requiredSocialExperienceList, `
		SELECT *
		FROM sending_job_information_required_social_experiences
		WHERE	sending_job_information_id = ?
		`,
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return requiredSocialExperienceList, err
	}

	return requiredSocialExperienceList, nil
}

func (repo *SendingJobInformationRequiredSocialExperienceRepositoryImpl) GetListBySendingBillingAddressID(sendingBillingAddressID uint) ([]*entity.SendingJobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.SendingJobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&requiredSocialExperienceList, `
			SELECT *
			FROM sending_job_information_required_social_experiences
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

	return requiredSocialExperienceList, nil
}

func (repo *SendingJobInformationRequiredSocialExperienceRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.SendingJobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&requiredSocialExperienceList, `
			SELECT *
			FROM sending_job_information_required_social_experiences
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

	return requiredSocialExperienceList, nil
}

// 求人リストから必要資格を取得
func (repo *SendingJobInformationRequiredSocialExperienceRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.SendingJobInformationRequiredSocialExperience
	)

	if len(idList) == 0 {
		return requiredSocialExperienceList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			sending_job_information_required_social_experiences.*
		FROM 
			sending_job_information_required_social_experiences
		WHERE
			sending_job_information_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&requiredSocialExperienceList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredSocialExperienceList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRequiredSocialExperienceRepositoryImpl) GetAll() ([]*entity.SendingJobInformationRequiredSocialExperience, error) {
	var (
		requiredSocialExperienceList []*entity.SendingJobInformationRequiredSocialExperience
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&requiredSocialExperienceList, `
			SELECT *
			FROM sending_job_information_required_social_experiences
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return requiredSocialExperienceList, nil
}
