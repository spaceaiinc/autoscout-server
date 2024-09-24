package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDesiredOccupationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDesiredOccupationRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDesiredOccupationRepository {
	return &SendingJobSeekerDesiredOccupationRepositoryImpl{
		Name:     "SendingJobSeekerDesiredOccupationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) Create(desiredOccupation *entity.SendingJobSeekerDesiredOccupation) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_desired_occupations (
				sending_job_seeker_id,
				desired_occupation,
				desired_rank,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		desiredOccupation.SendingJobSeekerID,
		desiredOccupation.DesiredOccupation,
		desiredOccupation.DesiredRank,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredOccupation.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) CreateMulti(sendingJobSeekerID uint, desiredOccupationList []entity.SendingJobSeekerDesiredOccupation) error {
	var (
		nowTime   = time.Now().In(time.UTC)
		valuesStr string
		srtFields []string
	)

	for _, do := range desiredOccupationList {
		srtFields = append(
			srtFields,
			fmt.Sprintf(
				"( %v, %v, %v, %s, %s )",
				sendingJobSeekerID,
				do.DesiredOccupation.Int64,
				do.DesiredRank.Int64,
				nowTime.Format("\"2006-01-02 15:04:05\""),
				nowTime.Format("\"2006-01-02 15:04:05\""),
			),
		)
	}

	valuesStr = strings.Join(srtFields, ", ")

	query := fmt.Sprintf(`
		INSERT INTO sending_job_seeker_desired_occupations (
			sending_job_seeker_id,
			desired_occupation,
			desired_rank,
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

func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_desired_occupations
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDesiredOccupation, error) {
	var (
		desiredOccupationList []*entity.SendingJobSeekerDesiredOccupation
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&desiredOccupationList, `
		SELECT *
		FROM sending_job_seeker_desired_occupations
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return desiredOccupationList, err
	}

	return desiredOccupationList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDesiredOccupation, error) {
	var (
		desiredOccupationList []*entity.SendingJobSeekerDesiredOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&desiredOccupationList, `
			SELECT 
				jsdo.*
			FROM 
				sending_job_seeker_desired_occupations AS jsdo
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdo.sending_job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDesiredOccupation, error) {
	var (
		desiredOccupationList []*entity.SendingJobSeekerDesiredOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&desiredOccupationList, `
			SELECT 
				jsdo.*
			FROM 
				sending_job_seeker_desired_occupations AS jsdo
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdo.sending_job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}

// 求職者リストから希望職種を取得
func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDesiredOccupation, error) {
	var (
		desiredOccupationList []*entity.SendingJobSeekerDesiredOccupation
	)

	if len(idList) == 0 {
		return desiredOccupationList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_desired_occupations
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&desiredOccupationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDesiredOccupationRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDesiredOccupation, error) {
	var (
		desiredOccupationList []*entity.SendingJobSeekerDesiredOccupation
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&desiredOccupationList, `
							SELECT *
							FROM sending_job_seeker_desired_occupations
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredOccupationList, nil
}
