package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDesiredJobInformationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDesiredJobInformationRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDesiredJobInformationRepository {
	return &SendingJobSeekerDesiredJobInformationRepositoryImpl{
		Name:     "SendingJobSeekerDesiredJobInformationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDesiredJobInformationRepositoryImpl) Create(desiredJobInformation *entity.SendingJobSeekerDesiredJobInformation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_desired_job_informations (
				sending_job_seeker_id,
				sending_job_information_id,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		desiredJobInformation.SendingJobSeekerID,
		desiredJobInformation.SendingJobInformationID,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredJobInformation.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerDesiredJobInformationRepositoryImpl) CreateMulti(sendingJobSeekerID uint, sendingJobInfoIDList []uint) error {
	var (
		nowTime   = time.Now().In(time.UTC)
		valuesStr string
		srtFields []string
	)

	for _, jobInfoID := range sendingJobInfoIDList {
		srtFields = append(
			srtFields,
			fmt.Sprintf(
				"( %v, %v, %s, %s )",
				sendingJobSeekerID,
				jobInfoID,
				nowTime.Format("\"2006-01-02 15:04:05\""),
				nowTime.Format("\"2006-01-02 15:04:05\""),
			),
		)
	}

	valuesStr = strings.Join(srtFields, ", ")

	query := fmt.Sprintf(`
		INSERT INTO sending_job_seeker_desired_job_informations (
			sending_job_seeker_id,
			sending_job_information_id,
			created_at,
			updated_at
		) 
		VALUES %s
	`, valuesStr)

	_, err := repo.executer.Exec(
		repo.Name+".CreateMulti", query,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (repo *SendingJobSeekerDesiredJobInformationRepositoryImpl) DeleteBySendingJobSeekerID(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_desired_job_informations
		WHERE sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDesiredJobInformationRepositoryImpl) GetListBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDesiredJobInformation, error) {
	var (
		desiredJobInformationList []*entity.SendingJobSeekerDesiredJobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingJobSeekerID",
		&desiredJobInformationList, `
		SELECT 
			desired.*,
			info.company_name,
			info.title,
			enterprise.company_name AS agent_name,
			enterprise.id AS sending_enterprise_id
		FROM 
		  sending_job_seeker_desired_job_informations AS desired
		INNER JOIN 
		  sending_job_informations AS info
		ON
			desired.sending_job_information_id = info.id
		INNER JOIN 
		  sending_billing_addresses AS billing
		ON
			info.sending_billing_address_id = billing.id
		INNER JOIN
			sending_enterprises AS enterprise
		ON
			billing.sending_enterprise_id = enterprise.id
		WHERE
			desired.sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return desiredJobInformationList, err
	}

	return desiredJobInformationList, nil
}
