package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerLPLoginTokenRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerLPLoginTokenRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerLPLoginTokenRepository {
	return &JobSeekerLPLoginTokenRepositoryImpl{
		Name:     "JobSeekerLPLoginTokenRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerLPLoginTokenRepositoryImpl) Create(lpLoginToken *entity.JobSeekerLPLoginToken) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_lp_login_token (
				job_seeker_id,
				login_token,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		lpLoginToken.JobSeekerID,
		lpLoginToken.LoginToken,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	lpLoginToken.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerLPLoginTokenRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_lp_login_token
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
/// 単数取得
//
func (repo *JobSeekerLPLoginTokenRepositoryImpl) FindByJobSeekerID(jobSeekerID uint) (*entity.JobSeekerLPLoginToken, error) {
	var (
		loginToken entity.JobSeekerLPLoginToken
	)

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerID",
		&loginToken, `
		SELECT *
		FROM job_seeker_lp_login_token
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &loginToken, nil
}

/****************************************************************************************/
/// 更新
//
func (repo *JobSeekerLPLoginTokenRepositoryImpl) UpdateByJobSeekerID(jobSeekerID uint, lpLoginToken *entity.JobSeekerLPLoginToken) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateByJobSeekerID",
		`
		UPDATE job_seeker_lp_login_token 
		SET
			login_token = ?,
			updated_at = ?
		WHERE 
			job_seeker_id = ?
		`,
		lpLoginToken.LoginToken,
		time.Now().In(time.UTC),
		jobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
