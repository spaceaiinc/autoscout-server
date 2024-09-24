package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDesiredWorkLocationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDesiredWorkLocationRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDesiredWorkLocationRepository {
	return &SendingJobSeekerDesiredWorkLocationRepositoryImpl{
		Name:     "SendingJobSeekerDesiredWorkLocationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) Create(desiredWorkLocation *entity.SendingJobSeekerDesiredWorkLocation) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_desired_work_locations (
				sending_job_seeker_id,
				desired_work_location,
				desired_rank,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		desiredWorkLocation.SendingJobSeekerID,
		desiredWorkLocation.DesiredWorkLocation,
		desiredWorkLocation.DesiredRank,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredWorkLocation.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) CreateMulti(sendingJobSeekerID uint, desiredWorkLocationList []entity.SendingJobSeekerDesiredWorkLocation) error {
	var (
		nowTime   = time.Now().In(time.UTC)
		valuesStr string
		srtFields []string
	)

	for _, dw := range desiredWorkLocationList {
		srtFields = append(
			srtFields,
			fmt.Sprintf(
				"( %v, %v, %v, %s, %s )",
				sendingJobSeekerID,
				dw.DesiredWorkLocation.Int64,
				dw.DesiredRank.Int64,
				nowTime.Format("\"2006-01-02 15:04:05\""),
				nowTime.Format("\"2006-01-02 15:04:05\""),
			),
		)
	}

	valuesStr = strings.Join(srtFields, ", ")

	query := fmt.Sprintf(`
		INSERT INTO sending_job_seeker_desired_work_locations (
			sending_job_seeker_id,
			desired_work_location,
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

func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_desired_work_locations
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDesiredWorkLocation, error) {
	var (
		desiredWorkLocationList []*entity.SendingJobSeekerDesiredWorkLocation
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&desiredWorkLocationList, `
		SELECT *
		FROM sending_job_seeker_desired_work_locations
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return desiredWorkLocationList, err
	}

	return desiredWorkLocationList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDesiredWorkLocation, error) {
	var (
		desiredWorkLocationList []*entity.SendingJobSeekerDesiredWorkLocation
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&desiredWorkLocationList, `
			SELECT 
				jsdwl.*
			FROM 
				sending_job_seeker_desired_work_locations AS jsdwl
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdwl.sending_job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDesiredWorkLocation, error) {
	var (
		desiredWorkLocationList []*entity.SendingJobSeekerDesiredWorkLocation
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&desiredWorkLocationList, `
			SELECT 
				jsdwl.*
			FROM 
				sending_job_seeker_desired_work_locations AS jsdwl
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdwl.sending_job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}

// 求職者リストから希望勤務地を取得
func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDesiredWorkLocation, error) {
	var (
		desiredWorkLocationList []*entity.SendingJobSeekerDesiredWorkLocation
	)

	if len(idList) == 0 {
		return desiredWorkLocationList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_desired_work_locations
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&desiredWorkLocationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDesiredWorkLocationRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDesiredWorkLocation, error) {
	var (
		desiredWorkLocationList []*entity.SendingJobSeekerDesiredWorkLocation
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&desiredWorkLocationList, `
							SELECT *
							FROM sending_job_seeker_desired_work_locations
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredWorkLocationList, nil
}
