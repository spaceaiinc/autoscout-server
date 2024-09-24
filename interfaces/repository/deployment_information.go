package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type DeploymentInformationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewDeploymentInformationRepositoryImpl(ex interfaces.SQLExecuter) usecase.DeploymentInformationRepository {
	return &DeploymentInformationRepositoryImpl{
		Name:     "DeploymentInformationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// デプロイ情報を作成
func (repo *DeploymentInformationRepositoryImpl) Create(deploymentInfo *entity.DeploymentInformation) error {
	var (
		nowTime = time.Now().In(time.UTC)
	)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create", `
			INSERT INTO deployment_informations (
				be_ver,
				fe_ver,
				be_detail,
				fe_detail,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)
		`,
		deploymentInfo.BeVer,
		deploymentInfo.FeVer,
		deploymentInfo.BeDetail,
		deploymentInfo.FeDetail,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	deploymentInfo.ID = uint(lastID)
	deploymentInfo.CreatedAt = nowTime

	return nil
}

/****************************************************************************************/
/// 更新 API
//
// デプロイ情報を更新
func (repo *DeploymentInformationRepositoryImpl) Update(id uint, deploymentInfo *entity.DeploymentInformation) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE deployment_informations
		SET
			be_ver = ?,
			fe_ver = ?,
			be_detail = ?,
			fe_detail = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		deploymentInfo.BeVer,
		deploymentInfo.FeVer,
		deploymentInfo.BeDetail,
		deploymentInfo.FeDetail,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得 API
//
// 全てのデプロイ情報を取得
func (repo *DeploymentInformationRepositoryImpl) All() ([]*entity.DeploymentInformation, error) {
	var (
		deploymentInformationList []*entity.DeploymentInformation
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&deploymentInformationList, `
			SELECT *
			FROM deployment_informations
			ORDER BY created_at DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return deploymentInformationList, err
}
