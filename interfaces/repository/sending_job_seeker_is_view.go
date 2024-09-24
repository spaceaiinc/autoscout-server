package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerIsViewRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerIsViewRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerIsViewRepository {
	return &SendingJobSeekerIsViewRepositoryImpl{
		Name:     "SendingJobSeekerIsViewRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerIsViewRepositoryImpl) Create(isView *entity.SendingJobSeekerIsView) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_is_views (
				sending_job_seeker_id,
				is_not_waiting_viewed,
				is_not_unregister_viewed,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		isView.SendingJobSeekerID,
		isView.IsNotWaitingViewed,
		isView.IsNotUnregisterViewed,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	isView.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerIsViewRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
			DELETE
			FROM sending_job_seeker_is_views
			WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerIsViewRepositoryImpl) UpdateIsNotWaitingViewedBySendingJobSeekerID(sendingJobSeekerID uint, isViewed bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsNotWaitingViewedBySendingJobSeekerID",
		`
			UPDATE 
				sending_job_seeker_is_views
			SET
				is_not_waiting_viewed = ?,
				updated_at = ?
			WHERE 
				sending_job_seeker_id = ?
		`,
		isViewed,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerIsViewRepositoryImpl) UpdateIsNotUnregisterViewedBySendingJobSeekerID(sendingJobSeekerID uint, isViewed bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsNotUnregisterViewedBySendingJobSeekerID",
		`
			UPDATE 
				sending_job_seeker_is_views
			SET
				is_not_unregister_viewed = ?,
				updated_at = ?
			WHERE 
				sending_job_seeker_id = ?
		`,
		isViewed,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerIsViewRepositoryImpl) GetIsNotWaitingViewCountByAgentStaffID(agentStaffID uint) (uint, error) {
	result := struct {
		Count uint `db:"waiting_count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetIsNotWaitingViewCountByAgentStaffID",
		&result, `
			SELECT 
				COUNT(is_views.id) AS waiting_count
			FROM 
				sending_job_seeker_is_views AS is_views
			INNER JOIN
				sending_job_seekers AS seeker
			ON
				is_views.sending_job_seeker_id = seeker.id
			WHERE 
				seeker.agent_staff_id = ?
			AND 
				seeker.phase = 2 AND
				is_views.is_not_waiting_viewed = TRUE
			AND
				seeker.interview_date > NOW()
		`,
		agentStaffID,
		// int64(entity.WaitingForInterview), // 面談実施待ち
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, err
}

func (repo *SendingJobSeekerIsViewRepositoryImpl) GetIsNotUnregisterViewCountByAgentStaffID(agentStaffID uint) (uint, error) {
	result := struct {
		Count uint `db:"unregister_count"`
	}{}

	err := repo.executer.Get(
		repo.Name+".GetIsNotUnregisterViewCountByAgentStaffID",
		&result, `
			SELECT 
				COUNT(is_views.id) AS unregister_count
			FROM 
				sending_job_seeker_is_views AS is_views
			INNER JOIN
				sending_job_seekers AS seeker
			ON
				is_views.sending_job_seeker_id = seeker.id
			WHERE 
				seeker.agent_staff_id = ?
			AND 
				seeker.phase = 1 &&
				is_views.is_not_unregister_viewed = TRUE
		`,
		agentStaffID,
		// int64(entity.UnregisterDetail),    // 詳細未登録（日程未登録のタイミングでは担当者決まっていないため）
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return result.Count, err
}
