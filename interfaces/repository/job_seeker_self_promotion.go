package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerSelfPromotionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerSelfPromotionRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerSelfPromotionRepository {
	return &JobSeekerSelfPromotionRepositoryImpl{
		Name:     "JobSeekerSelfPromotionRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerSelfPromotionRepositoryImpl) Create(selfPromotion *entity.JobSeekerSelfPromotion) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_self_promotions (
				job_seeker_id,
				title,
				contents,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		selfPromotion.JobSeekerID,
		selfPromotion.Title,
		selfPromotion.Contents,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	selfPromotion.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerSelfPromotionRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_self_promotions
		WHERE job_seeker_id = ?
		`, jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得
//
func (repo *JobSeekerSelfPromotionRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerSelfPromotion, error) {
	var selfPromotionList []*entity.JobSeekerSelfPromotion

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&selfPromotionList, `
		SELECT *
		FROM job_seeker_self_promotions
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return selfPromotionList, err
	}

	return selfPromotionList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerSelfPromotionRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerSelfPromotion, error) {
	var selfPromotionList []*entity.JobSeekerSelfPromotion

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&selfPromotionList, `
			SELECT 
				jssp.*
			FROM 
				job_seeker_self_promotions AS jssp
			INNER JOIN
				job_seekers AS js
			ON
				jssp.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selfPromotionList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerSelfPromotionRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerSelfPromotion, error) {
	var selfPromotionList []*entity.JobSeekerSelfPromotion

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&selfPromotionList, `
			SELECT 
				jssp.*
			FROM 
				job_seeker_self_promotions AS jssp
			INNER JOIN
				job_seekers AS js
			ON
				jssp.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selfPromotionList, nil
}

// 求職者リストから自己PRを取得
func (repo *JobSeekerSelfPromotionRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerSelfPromotion, error) {
	var selfPromotionList []*entity.JobSeekerSelfPromotion

	if len(idList) == 0 {
		return selfPromotionList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_self_promotions
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&selfPromotionList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selfPromotionList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerSelfPromotionRepositoryImpl) All() ([]*entity.JobSeekerSelfPromotion, error) {
	var selfPromotionList []*entity.JobSeekerSelfPromotion

	err := repo.executer.Select(
		repo.Name+".All",
		&selfPromotionList, `
							SELECT *
							FROM job_seeker_self_promotions
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selfPromotionList, nil
}
