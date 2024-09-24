package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type TaskIsRecommendDocumentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewTaskIsRecommendDocumentRepositoryImpl(ex interfaces.SQLExecuter) usecase.TaskIsRecommendDocumentRepository {
	return &TaskIsRecommendDocumentRepositoryImpl{
		Name:     "TaskIsRecommendDocumentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// タスクの作成
func (repo *TaskIsRecommendDocumentRepositoryImpl) Create(isRecommend *entity.TaskIsRecommendDocument) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO task_is_recommend_documents (
			task_id,
			is_recommend_upload_resume,
			is_recommend_upload_cv,
			is_recommend_upload_recommendation,
			is_recommend_generated_resume,
			is_recommend_generated_cv,
			is_recommend_generated_recommendation,
			is_recommend_generated_mask_resume,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
		isRecommend.TaskID,
		isRecommend.IsRecommendUploadResume,
		isRecommend.IsRecommendUploadCV,
		isRecommend.IsRecommendUploadRecommendation,
		isRecommend.IsRecommendGeneratedResume,
		isRecommend.IsRecommendGeneratedCV,
		isRecommend.IsRecommendGeneratedRecommendation,
		isRecommend.IsRecommendGeneratedMaskResume,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	isRecommend.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
func (repo *TaskIsRecommendDocumentRepositoryImpl) Update(id uint, isRecommend *entity.TaskIsRecommendDocument) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE 
			task_is_recommend_documents
		SET
			is_recommend_upload_resume = ?,
			is_recommend_upload_cv = ?,
			is_recommend_upload_recommendation = ?,
			is_recommend_generated_resume = ?,
			is_recommend_generated_cv = ?,
			is_recommend_generated_recommendation = ?,
			is_recommend_generated_mask_resume = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		isRecommend.IsRecommendUploadResume,
		isRecommend.IsRecommendUploadCV,
		isRecommend.IsRecommendUploadRecommendation,
		isRecommend.IsRecommendGeneratedResume,
		isRecommend.IsRecommendGeneratedCV,
		isRecommend.IsRecommendGeneratedRecommendation,
		isRecommend.IsRecommendGeneratedMaskResume,
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
/// 単数取得
//
func (repo *TaskIsRecommendDocumentRepositoryImpl) FindByID(id uint) (*entity.TaskIsRecommendDocument, error) {
	var (
		isRecommend entity.TaskIsRecommendDocument
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&isRecommend, `
		SELECT 
			*	
		FROM 
			task_is_recommend_documents
		WHERE
			id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &isRecommend, nil
}

func (repo *TaskIsRecommendDocumentRepositoryImpl) FindByTaskID(taskID uint) (*entity.TaskIsRecommendDocument, error) {
	var (
		isRecommend entity.TaskIsRecommendDocument
	)

	err := repo.executer.Get(
		repo.Name+".FindByTaskID",
		&isRecommend, `
		SELECT 
			*	
		FROM 
			task_is_recommend_documents
		WHERE
			task_id = ?
		LIMIT 1
		`,
		taskID)

	if err != nil {
		return nil, err
	}

	return &isRecommend, nil
}

/****************************************************************************************/
/// 複数取得
//
func (repo *TaskIsRecommendDocumentRepositoryImpl) GetByTaskIDList(idList []uint) ([]*entity.TaskIsRecommendDocument, error) {
	var (
		isRecommendList []*entity.TaskIsRecommendDocument
	)

	if len(idList) == 0 {
		return isRecommendList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			* 
		FROM 
			task_is_recommend_documents
		WHERE
			task_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByTaskIDList",
		&isRecommendList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return isRecommendList, nil
}
