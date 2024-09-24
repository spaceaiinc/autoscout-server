package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDesiredIndustryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDesiredIndustryRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDesiredIndustryRepository {
	return &SendingJobSeekerDesiredIndustryRepositoryImpl{
		Name:     "SendingJobSeekerDesiredIndustryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) Create(desiredIndustry *entity.SendingJobSeekerDesiredIndustry) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_desired_industries (
				sending_job_seeker_id,
				desired_industry,
				desired_rank,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		desiredIndustry.SendingJobSeekerID,
		desiredIndustry.DesiredIndustry,
		desiredIndustry.DesiredRank,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	desiredIndustry.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) CreateMulti(sendingJobSeekerID uint, desiredIndustryList []entity.SendingJobSeekerDesiredIndustry) error {
	var (
		nowTime   = time.Now().In(time.UTC)
		valuesStr string
		srtFields []string
	)

	for _, di := range desiredIndustryList {
		srtFields = append(
			srtFields,
			fmt.Sprintf(
				"( %v, %v, %v, %s, %s )",
				sendingJobSeekerID,
				di.DesiredIndustry.Int64,
				di.DesiredRank.Int64,
				nowTime.Format("\"2006-01-02 15:04:05\""),
				nowTime.Format("\"2006-01-02 15:04:05\""),
			),
		)
	}

	valuesStr = strings.Join(srtFields, ", ")

	query := fmt.Sprintf(`
		INSERT INTO sending_job_seeker_desired_industries (
			sending_job_seeker_id,
			desired_industry,
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

func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) Update(industry *entity.SendingJobSeekerDesiredIndustry) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sending_job_seeker_desired_industries
		SET
			sending_job_seeker_id = ?,
			desired_industry = ?,
			desired_rank = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		industry.SendingJobSeekerID,
		industry.DesiredIndustry,
		industry.DesiredRank,
		time.Now().In(time.UTC),
		industry.ID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_desired_industries
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerDesiredIndustry, error) {
	var (
		desiredIndustryList []*entity.SendingJobSeekerDesiredIndustry
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&desiredIndustryList, `
		SELECT
			*
		FROM 
			sending_job_seeker_desired_industries
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return desiredIndustryList, err
	}

	return desiredIndustryList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDesiredIndustry, error) {
	var (
		desiredIndustryList []*entity.SendingJobSeekerDesiredIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&desiredIndustryList, `
			SELECT 
				jsdi.*
			FROM 
				sending_job_seeker_desired_industries AS jsdi
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdi.sending_job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDesiredIndustry, error) {
	var (
		desiredIndustryList []*entity.SendingJobSeekerDesiredIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&desiredIndustryList, `
			SELECT 
				jsdi.*
			FROM 
				sending_job_seeker_desired_industries AS jsdi
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsdi.sending_job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}

// 求職者リストから希望業界を取得
func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDesiredIndustry, error) {
	var (
		desiredIndustryList []*entity.SendingJobSeekerDesiredIndustry
	)

	if len(idList) == 0 {
		return desiredIndustryList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_desired_industries
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&desiredIndustryList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDesiredIndustryRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDesiredIndustry, error) {
	var (
		desiredIndustryList []*entity.SendingJobSeekerDesiredIndustry
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&desiredIndustryList, `
							SELECT *
							FROM sending_job_seeker_desired_industries
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return desiredIndustryList, nil
}
