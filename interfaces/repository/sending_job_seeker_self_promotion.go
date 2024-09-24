package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerSelfPromotionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerSelfPromotionRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerSelfPromotionRepository {
	return &SendingJobSeekerSelfPromotionRepositoryImpl{
		Name:     "SendingJobSeekerSelfPromotionRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerSelfPromotionRepositoryImpl) Create(selfPromotion *entity.SendingJobSeekerSelfPromotion) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_self_promotions (
				sending_job_seeker_id,
				title,
				contents,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		selfPromotion.SendingJobSeekerID,
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

func (repo *SendingJobSeekerSelfPromotionRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_self_promotions
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerSelfPromotionRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerSelfPromotion, error) {
	var (
		selfPromotionList []*entity.SendingJobSeekerSelfPromotion
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&selfPromotionList, `
		SELECT *
		FROM sending_job_seeker_self_promotions
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return selfPromotionList, err
	}

	return selfPromotionList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerSelfPromotionRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerSelfPromotion, error) {
	var (
		selfPromotionList []*entity.SendingJobSeekerSelfPromotion
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&selfPromotionList, `
			SELECT 
				jssp.*
			FROM 
				sending_job_seeker_self_promotions AS jssp
			INNER JOIN
				sending_job_seekers AS js
			ON
				jssp.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerSelfPromotionRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerSelfPromotion, error) {
	var (
		selfPromotionList []*entity.SendingJobSeekerSelfPromotion
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&selfPromotionList, `
			SELECT 
				jssp.*
			FROM 
				sending_job_seeker_self_promotions AS jssp
			INNER JOIN
				sending_job_seekers AS js
			ON
				jssp.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerSelfPromotionRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerSelfPromotion, error) {
	var (
		selfPromotionList []*entity.SendingJobSeekerSelfPromotion
	)

	if len(idList) == 0 {
		return selfPromotionList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_self_promotions
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&selfPromotionList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selfPromotionList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerSelfPromotionRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerSelfPromotion, error) {
	var (
		selfPromotionList []*entity.SendingJobSeekerSelfPromotion
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&selfPromotionList, `
							SELECT *
							FROM sending_job_seeker_self_promotions
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return selfPromotionList, nil
}
