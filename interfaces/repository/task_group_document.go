package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type TaskGroupDocumentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewTaskGroupDocumentRepositoryImpl(ex interfaces.SQLExecuter) usecase.TaskGroupDocumentRepository {
	return &TaskGroupDocumentRepositoryImpl{
		Name:     "TaskGroupDocumentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
// タスクの作成
func (repo *TaskGroupDocumentRepositoryImpl) Create(task *entity.TaskGroupDocument) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO task_group_documents (
			task_group_id,
			document1_url,
			document2_url,
			document3_url,
			document4_url,
			document5_url,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?,?, ?
			)`,
		task.TaskGroupID,
		task.Document1URL,
		task.Document2URL,
		task.Document3URL,
		task.Document4URL,
		task.Document5URL,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	task.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新
//
func (repo *TaskGroupDocumentRepositoryImpl) Update(id uint, document *entity.TaskGroupDocument) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE task_group_documents
		SET
			document1_url = ?,
			document2_url = ?,
			document3_url = ?,
			document4_url = ?,
			document5_url = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		document.Document1URL,
		document.Document2URL,
		document.Document3URL,
		document.Document4URL,
		document.Document5URL,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 内定通知書の更新
func (repo *TaskGroupDocumentRepositoryImpl) UpdateDocument1URL(id uint, document1URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateDocument1URL",
		`
		UPDATE task_group_documents
		SET
			document1_url = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		document1URL,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 条件通知書の更新
func (repo *TaskGroupDocumentRepositoryImpl) UpdateDocument2URL(id uint, document2URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateDocument2URL",
		`
		UPDATE task_group_documents
		SET
			document2_url = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		document2URL,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 内定承諾書の更新
func (repo *TaskGroupDocumentRepositoryImpl) UpdateDocument3URL(id uint, document3URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateDocument3URL",
		`
		UPDATE task_group_documents
		SET
			document3_url = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		document3URL,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *TaskGroupDocumentRepositoryImpl) UpdateDocument4URL(id uint, document4URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateDocument4URL",
		`
		UPDATE task_group_documents
		SET
			document4_url = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		document4URL,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *TaskGroupDocumentRepositoryImpl) UpdateDocument5URL(id uint, document5URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateDocument5URL",
		`
		UPDATE task_group_documents
		SET
		document5_url = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		document5URL,
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
/// 削除
//
func (repo *TaskGroupDocumentRepositoryImpl) Delete(taskID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM task_group_documents
		WHERE id = ?
		`, taskID,
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
func (repo *TaskGroupDocumentRepositoryImpl) FindByID(id uint) (*entity.TaskGroupDocument, error) {
	var (
		task entity.TaskGroupDocument
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&task, `
		SELECT *
		FROM task_group_documents
		WHERE
			id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &task, nil
}

func (repo *TaskGroupDocumentRepositoryImpl) FindByGroupID(taskGroupID uint) (*entity.TaskGroupDocument, error) {
	var (
		task entity.TaskGroupDocument
	)

	err := repo.executer.Get(
		repo.Name+".FindByGroupID",
		&task, `
		SELECT *
		FROM task_group_documents
		WHERE
			task_group_id = ?
		LIMIT 1
		`,
		taskGroupID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &task, nil
}
