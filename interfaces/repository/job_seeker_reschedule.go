package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerRescheduleRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerRescheduleRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerRescheduleRepository {
	return &JobSeekerRescheduleRepositoryImpl{
		Name:     "JobSeekerRescheduleRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 作成
//
func (repo *JobSeekerRescheduleRepositoryImpl) Create(jobSeekerSchedule *entity.JobSeekerReschedule) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO job_seeker_reschedules (
			reschedule_id,	           
			task_id,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?
			)`,
		jobSeekerSchedule.RescheduleID,
		jobSeekerSchedule.TaskID,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	jobSeekerSchedule.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
// 削除
//
func (repo *JobSeekerRescheduleRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM job_seeker_reschedules
		WHERE id = ?
		`, id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
// 単数取得
//
func (repo *JobSeekerRescheduleRepositoryImpl) FindByID(id uint) (*entity.JobSeekerReschedule, error) {
	var (
		reschedule entity.JobSeekerReschedule
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&reschedule, `
		SELECT *
		FROM job_seeker_reschedules
		WHERE id = ?
		LIMIT 1
		`, id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &reschedule, nil
}
