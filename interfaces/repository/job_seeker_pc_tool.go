package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerPCToolRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerPCToolRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerPCToolRepository {
	return &JobSeekerPCToolRepositoryImpl{
		Name:     "JobSeekerPCToolRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerPCToolRepositoryImpl) Create(pcTool *entity.JobSeekerPCTool) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_pc_tools (
				job_seeker_id,
				tool,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		pcTool.JobSeekerID,
		pcTool.Tool,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	pcTool.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerPCToolRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_pc_tools
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
func (repo *JobSeekerPCToolRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerPCTool, error) {
	var pcToolList []*entity.JobSeekerPCTool

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&pcToolList, `
		SELECT *
		FROM job_seeker_pc_tools
		WHERE
			job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerPCToolRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerPCTool, error) {
	var pcToolList []*entity.JobSeekerPCTool

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&pcToolList, `
			SELECT 
				jspt.*
			FROM 
				job_seeker_pc_tools AS jspt
			INNER JOIN
				job_seekers AS js
			ON
				jspt.job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerPCToolRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerPCTool, error) {
	var pcToolList []*entity.JobSeekerPCTool

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&pcToolList, `
			SELECT 
				jspt.*
			FROM 
				job_seeker_pc_tools AS jspt
			INNER JOIN
				job_seekers AS js
			ON
				jspt.job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}

// 求職者リストから所持PCスキルを取得
func (repo *JobSeekerPCToolRepositoryImpl) GetByJobSeekerIDList(jobSeekerIDList []uint) ([]*entity.JobSeekerPCTool, error) {
	var pcToolList []*entity.JobSeekerPCTool

	if len(jobSeekerIDList) == 0 {
		return pcToolList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_pc_tools
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(jobSeekerIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&pcToolList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerPCToolRepositoryImpl) All() ([]*entity.JobSeekerPCTool, error) {
	var pcToolList []*entity.JobSeekerPCTool

	err := repo.executer.Select(
		repo.Name+".All",
		&pcToolList, `
							SELECT *
							FROM job_seeker_pc_tools
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}
