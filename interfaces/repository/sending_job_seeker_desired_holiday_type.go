package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDesiredHolidayTypeRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDesiredHolidayTypeRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDesiredHolidayTypeRepository {
	return &SendingJobSeekerDesiredHolidayTypeRepositoryImpl{
		Name:     "SendingJobSeekerDesiredHolidayTypeRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDesiredHolidayTypeRepositoryImpl) Create(desiredHolidayType *entity.SendingJobSeekerDesiredHolidayType) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_desired_holiday_types (
				sending_job_seeker_id,
				holiday_type,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		desiredHolidayType.SendingJobSeekerID,
		desiredHolidayType.HolidayType,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredHolidayType.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerDesiredHolidayTypeRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_desired_holiday_types
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDesiredHolidayTypeRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDesiredHolidayType, error) {
	var (
		desiredHolidayTypeList []*entity.SendingJobSeekerDesiredHolidayType
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&desiredHolidayTypeList, `
		SELECT *
		FROM sending_job_seeker_desired_holiday_types
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return desiredHolidayTypeList, err
	}

	return desiredHolidayTypeList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredHolidayTypeRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDesiredHolidayType, error) {
	var (
		desiredHolidayTypeList []*entity.SendingJobSeekerDesiredHolidayType
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&desiredHolidayTypeList, `
			SELECT 
				jsdht.*
			FROM 
				sending_job_seeker_desired_holiday_types AS jsdht
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdht.sending_job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredHolidayTypeRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDesiredHolidayType, error) {
	var (
		desiredHolidayTypeList []*entity.SendingJobSeekerDesiredHolidayType
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&desiredHolidayTypeList, `
			SELECT 
				jsdht.*
			FROM 
				sending_job_seeker_desired_holiday_types AS jsdht
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdht.sending_job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}

// 求職者リストから希望休日タイプを取得
func (repo *SendingJobSeekerDesiredHolidayTypeRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDesiredHolidayType, error) {
	var (
		desiredHolidayTypeList []*entity.SendingJobSeekerDesiredHolidayType
	)

	if len(idList) == 0 {
		return desiredHolidayTypeList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_desired_holiday_types
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&desiredHolidayTypeList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDesiredHolidayTypeRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDesiredHolidayType, error) {
	var (
		desiredHolidayTypeList []*entity.SendingJobSeekerDesiredHolidayType
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&desiredHolidayTypeList, `
							SELECT *
							FROM sending_job_seeker_desired_holiday_types
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredHolidayTypeList, nil
}
