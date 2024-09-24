package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDesiredCompanyScaleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDesiredCompanyScaleRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDesiredCompanyScaleRepository {
	return &SendingJobSeekerDesiredCompanyScaleRepositoryImpl{
		Name:     "SendingJobSeekerDesiredCompanyScaleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDesiredCompanyScaleRepositoryImpl) Create(desiredCompanyScale *entity.SendingJobSeekerDesiredCompanyScale) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_desired_company_scales (
				sending_job_seeker_id,
				desired_company_scale,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		desiredCompanyScale.SendingJobSeekerID,
		desiredCompanyScale.DesiredCompanyScale,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredCompanyScale.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerDesiredCompanyScaleRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_desired_company_scales
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDesiredCompanyScaleRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.SendingJobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&desiredCompanyScaleList, `
		SELECT *
		FROM sending_job_seeker_desired_company_scales
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return desiredCompanyScaleList, err
	}

	return desiredCompanyScaleList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredCompanyScaleRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.SendingJobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&desiredCompanyScaleList, `
			SELECT 
				jsdcs.*
			FROM 
				sending_job_seeker_desired_company_scales AS jsdcs
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdcs.sending_job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredCompanyScaleRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.SendingJobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&desiredCompanyScaleList, `
			SELECT 
				jsdcs.*
			FROM 
				sending_job_seeker_desired_company_scales AS jsdcs
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdcs.sending_job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}

// 求職者リストから希望休日タイプを取得
func (repo *SendingJobSeekerDesiredCompanyScaleRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.SendingJobSeekerDesiredCompanyScale
	)

	if len(idList) == 0 {
		return desiredCompanyScaleList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_desired_company_scales
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&desiredCompanyScaleList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDesiredCompanyScaleRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDesiredCompanyScale, error) {
	var (
		desiredCompanyScaleList []*entity.SendingJobSeekerDesiredCompanyScale
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&desiredCompanyScaleList, `
							SELECT *
							FROM sending_job_seeker_desired_company_scales
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredCompanyScaleList, nil
}
