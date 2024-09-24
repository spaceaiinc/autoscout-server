package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type DeploymentReflectionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewDeploymentReflectionRepositoryImpl(ex interfaces.SQLExecuter) usecase.DeploymentReflectionRepository {
	return &DeploymentReflectionRepositoryImpl{
		Name:     "DeploymentReflectionRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// デプロイの反映状況を複数作成
func (repo *DeploymentReflectionRepositoryImpl) CreateMulti(deploymentID uint, agentStaffList []*entity.AgentStaff) error {
	var (
		nowTime   = time.Now().In(time.UTC).Format("\"2006-01-02 15:04:05\"")
		valuesStr string
		srtFields []string
	)

	for _, staff := range agentStaffList {
		srtFields = append(
			srtFields,
			fmt.Sprintf(
				"( %v, %v, %v, %s, %s )",
				deploymentID,
				staff.ID,
				false,
				nowTime,
				nowTime,
			),
		)
	}

	valuesStr = strings.Join(srtFields, ", ")

	query := fmt.Sprintf(`
		INSERT INTO deployment_reflections (
			deployment_id,
			agent_staff_id,
			is_reflected,
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

/****************************************************************************************/
/// 更新 API
//
// デプロイの反映状況フラグを更新
func (repo *DeploymentReflectionRepositoryImpl) UpdateIsReflectedByAgentStaffID(staffID uint, isReflected bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsReflectedByAgentStaffID",
		`
		UPDATE deployment_reflections
		SET
			is_reflected = ?,
			updated_at = ?
		WHERE 
			agent_staff_id = ?
		`,
		isReflected,
		time.Now().In(time.UTC),
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得 API
//
// 最新の指定スタッフIDのis_reflectedがfalseのレコード取得
func (repo *DeploymentReflectionRepositoryImpl) FindNotReflectedByAgentStaffID(staffID uint) (*entity.DeploymentReflection, error) {
	var (
		deploymentReflection entity.DeploymentReflection
	)

	err := repo.executer.Get(
		repo.Name+".FindNotReflectedByAgentStaffID",
		&deploymentReflection, `
			SELECT *
			FROM deployment_reflections
			WHERE 
				agent_staff_id = ? AND
				is_reflected = FALSE
			ORDER BY id DESC
			LIMIT 1
		`, staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &deploymentReflection, err
}
