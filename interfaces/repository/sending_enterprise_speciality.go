package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingEnterpriseSpecialityRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingEnterpriseSpecialityRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingEnterpriseSpecialityRepository {
	return &SendingEnterpriseSpecialityRepositoryImpl{
		Name:     "SendingEnterpriseSpecialityRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *SendingEnterpriseSpecialityRepositoryImpl) Create(speciality *entity.SendingEnterpriseSpeciality) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_enterprise_specialities (
				sending_enterprise_id,
				image_url,
				job_information_count,
				specialized_occupation,
				specialized_industry,

				specialized_area,
				specialized_company_type,
				specialized_job_seeker_type,
				consulting_strengths,
				support_content,

				pr_point,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?,
				?, ?, ?, ?, ?,
				?, ?, ?
			)
		`,
		speciality.SendingEnterpriseID,
		speciality.ImageURL,
		speciality.JobInformationCount,
		speciality.SpecializedOccupation,
		speciality.SpecializedIndustry,
		speciality.SpecializedArea,
		speciality.SpecializedCompanyType,
		speciality.SpecializedJobSeekerType,
		speciality.ConsultingStrengths,
		speciality.SupportContent,
		speciality.PRPoint,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	speciality.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *SendingEnterpriseSpecialityRepositoryImpl) Update(speciality *entity.SendingEnterpriseSpeciality) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
			UPDATE 
				sending_enterprise_specialities 
			SET
				image_url = ?,
				job_information_count = ?,
				specialized_occupation = ?,
				specialized_industry = ?,
				specialized_area = ?,
				specialized_company_type = ?,
				specialized_job_seeker_type = ?,
				consulting_strengths = ?,
				support_content = ?,
				pr_point = ?,
				updated_at = ?
			WHERE sending_enterprise_id = ?
		`,
		speciality.ImageURL,
		speciality.JobInformationCount,
		speciality.SpecializedOccupation,
		speciality.SpecializedIndustry,
		speciality.SpecializedArea,
		speciality.SpecializedCompanyType,
		speciality.SpecializedJobSeekerType,
		speciality.ConsultingStrengths,
		speciality.SupportContent,
		speciality.PRPoint,
		time.Now().In(time.UTC),
		speciality.SendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (repo *SendingEnterpriseSpecialityRepositoryImpl) UpdateImageURLBySendingEnterpriseID(sendingEnterpriseID uint, imageURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateImageURLBySendingEnterpriseID",
		`
				UPDATE 
					sending_enterprise_specialities 
				SET
					image_url = ?,
					updated_at = ?
				WHERE sending_enterprise_id = ?
			`,
		imageURL,
		time.Now().In(time.UTC),
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *SendingEnterpriseSpecialityRepositoryImpl) FindBySendingEnterpriseID(sendingEnterpriseID uint) (*entity.SendingEnterpriseSpeciality, error) {
	var (
		speciality entity.SendingEnterpriseSpeciality
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingEnterpriseID",
		&speciality, `
		SELECT *
		FROM sending_enterprise_specialities
		WHERE
			sending_enterprise_id = ?
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return &speciality, err
	}

	return &speciality, nil
}
