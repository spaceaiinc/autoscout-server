package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerPCToolRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerPCToolRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerPCToolRepository {
	return &SendingJobSeekerPCToolRepositoryImpl{
		Name:     "SendingJobSeekerPCToolRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerPCToolRepositoryImpl) Create(pcTool *entity.SendingJobSeekerPCTool) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_pc_tools (
				sending_job_seeker_id,
				tool,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?
			)
		`,
		pcTool.SendingJobSeekerID,
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

func (repo *SendingJobSeekerPCToolRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_pc_tools
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerPCToolRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) ([]*entity.SendingJobSeekerPCTool, error) {
	var (
		pcToolList []*entity.SendingJobSeekerPCTool
	)

	err := repo.executer.Select(
		repo.Name+".FindBySendingJobSeekerID",
		&pcToolList, `
		SELECT *
		FROM sending_job_seeker_pc_tools
		WHERE
			sending_job_seeker_id = ?
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerPCToolRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerPCTool, error) {
	var (
		pcToolList []*entity.SendingJobSeekerPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&pcToolList, `
			SELECT 
				jspt.*
			FROM 
				sending_job_seeker_pc_tools AS jspt
			INNER JOIN
				sending_job_seekers AS js
			ON
				jspt.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerPCToolRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerPCTool, error) {
	var (
		pcToolList []*entity.SendingJobSeekerPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&pcToolList, `
			SELECT 
				jspt.*
			FROM 
				sending_job_seeker_pc_tools AS jspt
			INNER JOIN
				sending_job_seekers AS js
			ON
				jspt.sending_job_seeker_id = js.id
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
func (repo *SendingJobSeekerPCToolRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerPCTool, error) {
	var (
		pcToolList []*entity.SendingJobSeekerPCTool
	)

	if len(idList) == 0 {
		return pcToolList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_pc_tools
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&pcToolList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerPCToolRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerPCTool, error) {
	var (
		pcToolList []*entity.SendingJobSeekerPCTool
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&pcToolList, `
							SELECT *
							FROM sending_job_seeker_pc_tools
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pcToolList, nil
}
