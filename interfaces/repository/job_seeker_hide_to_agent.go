package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerHideToAgentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerHideToAgentRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerHideToAgentRepository {
	return &JobSeekerHideToAgentRepositoryImpl{
		Name:     "JobSeekerHideToAgentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerHideToAgentRepositoryImpl) Create(hideToAgent *entity.JobSeekerHideToAgent) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_hide_to_agents (
				job_seeker_id,
				agent_id,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		hideToAgent.JobSeekerID,
		hideToAgent.AgentID,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	hideToAgent.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerHideToAgentRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_hide_to_agents
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
func (repo *JobSeekerHideToAgentRepositoryImpl) GetByJobSeekerID(jobSeekerID uint) ([]*entity.JobSeekerHideToAgent, error) {
	var hideToAgentList []*entity.JobSeekerHideToAgent

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerID",
		&hideToAgentList, `
		SELECT 
			hide.*, agent.agent_name
		FROM 
			job_seeker_hide_to_agents AS hide
		INNER JOIN 
			agents AS agent
		ON
			hide.agent_id = agent.id
		WHERE
			hide.job_seeker_id = ?
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return hideToAgentList, err
	}

	return hideToAgentList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *JobSeekerHideToAgentRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerHideToAgent, error) {
	var hideToAgentList []*entity.JobSeekerHideToAgent

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_seeker_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
			INNER JOIN 
				job_seekers AS seeker
			ON
				hide.job_seeker_id = seeker.id
			WHERE
				seeker.agent_id = ?
		`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *JobSeekerHideToAgentRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerHideToAgent, error) {
	var hideToAgentList []*entity.JobSeekerHideToAgent

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_seeker_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
			INNER JOIN 
				job_seekers AS seeker
			ON
				hide.job_seeker_id = seeker.id
			WHERE
				seeker.agent_staff_id = ?
		`,
		staffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

func (repo *JobSeekerHideToAgentRepositoryImpl) GetHideByAgentID(agentID uint) ([]*entity.JobSeekerHideToAgent, error) {
	var hideToAgentList []*entity.JobSeekerHideToAgent

	err := repo.executer.Select(
		repo.Name+".GetHideByAgentID",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_seeker_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
			WHERE
				agent_id = ?
		`, agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

func (repo *JobSeekerHideToAgentRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerHideToAgent, error) {
	var hideToAgentList []*entity.JobSeekerHideToAgent

	if len(idList) == 0 {
		return hideToAgentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			hide.*, agent.agent_name
		FROM 
			job_seeker_hide_to_agents AS hide
		INNER JOIN 
			agents AS agent
		ON
			hide.agent_id = agent.id
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&hideToAgentList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

func (repo *JobSeekerHideToAgentRepositoryImpl) GetByAgentIDList(agentIDList []uint) ([]*entity.JobSeekerHideToAgent, error) {
	var hideToAgentList []*entity.JobSeekerHideToAgent

	if len(agentIDList) == 0 {
		return hideToAgentList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			hide.*, agent.agent_name
		FROM 
			job_seeker_hide_to_agents AS hide
		INNER JOIN 
			agents AS agent
		ON
			hide.agent_id = agent.id
		WHERE
			hide.agent_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(agentIDList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDList",
		&hideToAgentList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerHideToAgentRepositoryImpl) All() ([]*entity.JobSeekerHideToAgent, error) {
	var hideToAgentList []*entity.JobSeekerHideToAgent

	err := repo.executer.Select(
		repo.Name+".All",
		&hideToAgentList, `
			SELECT 
				hide.*, agent.agent_name
			FROM 
				job_seeker_hide_to_agents AS hide
			INNER JOIN 
				agents AS agent
			ON
				hide.agent_id = agent.id
		`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return hideToAgentList, nil
}
