package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type InitialEnterpriseImporterRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewInitialEnterpriseImporterRepositoryImpl(ex interfaces.SQLExecuter) usecase.InitialEnterpriseImporterRepository {
	return &InitialEnterpriseImporterRepositoryImpl{
		Name:     "InitialEnterpriseImporterRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// 求人インポートの作成
func (repo *InitialEnterpriseImporterRepositoryImpl) Create(importer *entity.InitialEnterpriseImporter) error {
	nowTime := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO initial_enterprise_importers (
			uuid,
			agent_staff_id,
			service_type,
			login_id,
			password,
			
			start_date,
			start_hour,
			max_count,
			offset,
			search_url,

			is_success,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?,
				?, ?, ?, ?, ?,
				?, ?, ?
			)`,
		utility.CreateUUID(),
		importer.AgentStaffID,
		importer.ServiceType,
		importer.LoginID,
		importer.Password,
		importer.StartDate,
		importer.StartHour,
		importer.MaxCount,
		importer.Offset,
		importer.SearchURL,
		importer.IsSuccess,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	importer.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
// 求人インポートの更新
func (repo *InitialEnterpriseImporterRepositoryImpl) UpdateIsSuccessTrue(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsSuccessTrue",
		`
		UPDATE initial_enterprise_importers
		SET
			is_success = TRUE,
			updated_at = ?
		WHERE 
			id = ?
		`,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 求人インポートの削除
func (repo *InitialEnterpriseImporterRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE 
		FROM initial_enterprise_importers
		WHERE id = ?
		`,
		id,
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
func (repo *InitialEnterpriseImporterRepositoryImpl) FindByID(id uint) (*entity.InitialEnterpriseImporter, error) {
	var (
		importer entity.InitialEnterpriseImporter
	)

	// login_idとpasswordは取得しない
	err := repo.executer.Get(
		repo.Name+".FindByID",
		&importer, `
		SELECT 
			importer.id,
			importer.uuid,
			importer.agent_staff_id,
			importer.service_type,
			importer.start_date,
			importer.start_hour,
			importer.max_count,
			importer.offset,
			importer.search_url,
			importer.is_success,
			importer.created_at,
			importer.updated_at
		FROM 
			initial_enterprise_importers AS importer
		WHERE
			importer.id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &importer, nil
}

/****************************************************************************************/
/// 複数取得
//
// エージェントIDから求人インポートの一覧を取得
func (repo *InitialEnterpriseImporterRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.InitialEnterpriseImporter, error) {
	var (
		importerList []*entity.InitialEnterpriseImporter
	)

	// login_idとpasswordは取得しない
	err := repo.executer.Select(
		repo.Name+".GetByIsSuccessFalse",
		&importerList, `
			SELECT 
				importer.id,
				importer.uuid,
				importer.agent_staff_id,
				importer.service_type,
				importer.start_date,
				importer.start_hour,
				importer.max_count,
				importer.offset,
				importer.search_url,
				importer.is_success,
				importer.created_at,
				importer.updated_at
			FROM 
			 initial_enterprise_importers AS importer
			WHERE
				importer.agent_staff_id IN (
					SELECT id
					FROM agent_staffs
					WHERE agent_id = ?
				)
			ORDER BY
				importer.id DESC
		`,
		agentID,
	)
	fmt.Println("importerList", importerList)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return importerList, nil
}

// 今後1週間の企業求人一括インポート情報一覧を取得 *予約日程調整時の重複チェック用
func (repo *InitialEnterpriseImporterRepositoryImpl) GetByWeek() ([]*entity.InitialEnterpriseImporter, error) {
	var (
		importerList []*entity.InitialEnterpriseImporter
	)

	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	oneDaysAgo := now.AddDate(0, 0, -1)
	sevenDaysLater := now.AddDate(0, 0, 7)

	// login_idとpasswordは取得しない
	err := repo.executer.Select(
		repo.Name+".GetByIsSuccessFalse",
		&importerList, `
			SELECT 
				importer.id,
				importer.uuid,
				importer.agent_staff_id,
				importer.service_type,
				importer.start_date,
				importer.start_hour,
				importer.max_count,
				importer.offset,
				importer.search_url,
				importer.is_success,
				importer.created_at,
				importer.updated_at
			FROM 
			 initial_enterprise_importers AS importer
			WHERE 
				DATE(start_date) BETWEEN ? AND ?
			ORDER BY
				importer.id DESC
		`,
		oneDaysAgo,
		sevenDaysLater,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return importerList, nil
}

// 実行日時が今日(JST)のものを取得
func (repo *InitialEnterpriseImporterRepositoryImpl) GetByStartDate(now time.Time) ([]*entity.InitialEnterpriseImporter, error) {
	var (
		importerList []*entity.InitialEnterpriseImporter
	)

	// RPAで使用するため、login_idとpasswordは取得しない
	err := repo.executer.Select(
		repo.Name+".GetByStartDate",
		&importerList, `
			SELECT *
			FROM 
			 initial_enterprise_importers
			WHERE
				is_success = FALSE 
			ORDER BY
				id DESC
		`,
		// now.Format("2006-01-02"),
		// now.Hour(),
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return importerList, nil
}
