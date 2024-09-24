package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type NotificationForUserRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewNotificationForUserRepositoryImpl(ex interfaces.SQLExecuter) usecase.NotificationForUserRepository {
	return &NotificationForUserRepositoryImpl{
		Name:     "NotificationForUserRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// お知らせを作成
func (repo *NotificationForUserRepositoryImpl) Create(notificationForUser *entity.NotificationForUser) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO notification_for_users (
			title,
			body,
			target,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)`,
		notificationForUser.Title,
		notificationForUser.Body,
		notificationForUser.Target,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	notificationForUser.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 単数取得 API
//
// お知らせを取得
func (repo *NotificationForUserRepositoryImpl) FindByID(id uint) (*entity.NotificationForUser, error) {
	var (
		notificationForUser entity.NotificationForUser
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&notificationForUser, `
		SELECT 
			*
		FROM 
			notification_for_users
		WHERE
			id = ?
		LIMIT 1
		`,
		id)

	if err != nil {
		return nil, err
	}

	return &notificationForUser, nil
}

// 最新のお知らせを取得
func (repo *NotificationForUserRepositoryImpl) FindLatest() (*entity.NotificationForUser, error) {
	var (
		notificationForUser entity.NotificationForUser
	)

	err := repo.executer.Get(
		repo.Name+".FindLatest",
		&notificationForUser, `
		SELECT 
			*
		FROM 
			notification_for_users
		ORDER BY
			id DESC
		LIMIT 1
		`,
	)

	if err != nil {
		return nil, err
	}

	return &notificationForUser, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// 指定ID以降、ターゲットのお知らせを取得
func (repo *NotificationForUserRepositoryImpl) GetAfterIDAndTargetList(id uint, targetList []entity.NotificationForUserTarget) ([]*entity.NotificationForUser, error) {
	var (
		notificationForUserList []*entity.NotificationForUser
	)

	if len(targetList) < 1 {
		return notificationForUserList, nil
	}

	targetListStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(targetList)), ", "), "[]")

	query := fmt.Sprintf(`
		SELECT
			*
		FROM
			notification_for_users
		WHERE
			id > %d
		AND
			target IN (%s)
		ORDER BY
			id DESC
	`,
		id,
		targetListStr,
	)

	err := repo.executer.Select(
		repo.Name+".GetAfterIDAndTargetList",
		&notificationForUserList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return notificationForUserList, nil
}

// 指定ターゲットのお知らせを取得
func (repo *NotificationForUserRepositoryImpl) GetByTargetList(targetList []entity.NotificationForUserTarget) ([]*entity.NotificationForUser, error) {
	var (
		notificationForUserList []*entity.NotificationForUser
	)

	if len(targetList) < 1 {
		return notificationForUserList, nil
	}

	targetListStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(targetList)), ", "), "[]")

	query := fmt.Sprintf(`
		SELECT
			*
		FROM
			notification_for_users
		WHERE
			target IN (%s)
		ORDER BY
			id DESC
	`,
		targetListStr,
	)

	err := repo.executer.Select(
		repo.Name+".GetByTargetList",
		&notificationForUserList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return notificationForUserList, nil
}

// 全てのお知らせを取得
func (repo *NotificationForUserRepositoryImpl) All() ([]*entity.NotificationForUser, error) {
	var (
		notificationForUserList []*entity.NotificationForUser
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&notificationForUserList, `
			SELECT 
				*
			FROM
				notification_for_users
			ORDER BY
				id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return notificationForUserList, nil
}
