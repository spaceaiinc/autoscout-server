package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type UserEntryRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewUserEntryRepositoryImpl(ex interfaces.SQLExecuter) usecase.UserEntryRepository {
	return &UserEntryRepositoryImpl{
		Name:     "UserEntryRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 汎用系 API
//

func (repo *UserEntryRepositoryImpl) Create(userEntry *entity.UserEntry) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO user_entries (
				user_id,
				service_type,
				is_processed,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		userEntry.UserID,
		userEntry.ServiceType,
		false,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	userEntry.ID = uint(lastID)
	return nil
}

func (repo *UserEntryRepositoryImpl) UpdateIsProcessedByUserID(userID string, isProcessed bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsProcessed",
		`
			UPDATE user_entries
			SET
				is_processed = ?,
				updated_at = ?
			WHERE 
				user_id = ?
		`,
		isProcessed,
		time.Now().In(time.UTC),
		userID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *UserEntryRepositoryImpl) UpdateIsProcessedByUserIDList(userIDList []string, isProcessed bool) error {
	idListStr := strings.Trim(strings.Join(userIDList, ", "), "[]")

	query := fmt.Sprintf(`
		UPDATE user_entries
		SET
			is_processed = ?,
			updated_at = ?
		WHERE 
			user_id IN(%s)
	`, idListStr)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsProcessedByUserIDList",
		query,
		isProcessed,
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 1時間以内に作成されたエントリー処理未実行のレコードを取得
func (repo *UserEntryRepositoryImpl) GetUnprocessed() ([]*entity.UserEntry, error) {
	var (
		userEntryList []*entity.UserEntry
	)

	err := repo.executer.Select(
		repo.Name+".GetUnprocessed",
		&userEntryList, `
			SELECT *
			FROM user_entries
			WHERE
				is_processed = false
			AND
				created_at > NOW() - INTERVAL 1 HOUR
			ORDER BY 
				id DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	return userEntryList, nil
}
