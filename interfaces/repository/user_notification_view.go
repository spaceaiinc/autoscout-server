package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type UserNotificationViewRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewUserNotificationViewRepositoryImpl(ex interfaces.SQLExecuter) usecase.UserNotificationViewRepository {
	return &UserNotificationViewRepositoryImpl{
		Name:     "UserNotificationViewRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// お知らせの既読を作成
func (repo *UserNotificationViewRepositoryImpl) Create(userNotificationView *entity.UserNotificationView) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO user_notification_views (
			notification_id,
			agent_staff_id,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?
			)`,
		userNotificationView.NotificationID,
		userNotificationView.AgentStaffID,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	userNotificationView.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 単数取得 API
//
// お知らせの既読を取得
func (repo *UserNotificationViewRepositoryImpl) FindByID(userNotificationViewID uint) (*entity.UserNotificationView, error) {
	var (
		userNotificationView entity.UserNotificationView
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&userNotificationView, `
		SELECT 
			*
		FROM 
			user_notification_views
		WHERE
			id = ?
		LIMIT 1
		`,
		userNotificationViewID)

	if err != nil {
		return nil, err
	}

	return &userNotificationView, nil
}

// 最新のお知らせ既読レコードを取得
func (repo *UserNotificationViewRepositoryImpl) FindLatestByAgentStaffID(agentStaffID uint) (*entity.UserNotificationView, error) {
	var (
		userNotificationView entity.UserNotificationView
	)

	err := repo.executer.Get(
		repo.Name+".FindLatestByAgentStaffID",
		&userNotificationView, `
		SELECT 
			*
		FROM 
			user_notification_views
		WHERE
			agent_staff_id = ?
		ORDER BY 
			notification_id DESC
		LIMIT 1
		`,
		agentStaffID)

	if err != nil {
		return nil, err
	}

	return &userNotificationView, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// 指定お知らせIDのお知らせの既読を取得
func (repo *UserNotificationViewRepositoryImpl) GetByNotificationID(notificationID uint) ([]*entity.UserNotificationView, error) {
	var (
		userNotificationViewList []*entity.UserNotificationView
	)

	err := repo.executer.Select(
		repo.Name+".GetByNotificationID",
		&userNotificationViewList, `
			SELECT 
				*
			FROM
				user_notification_views
			WHERE
				notification_id = ?
			ORDER BY
				id DESC
		`,
		notificationID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return userNotificationViewList, nil
}

// 指定お知らせIDリストのお知らせの既読を取得
func (repo *UserNotificationViewRepositoryImpl) GetByNotificationIDList(notificationIDList []uint) ([]*entity.UserNotificationView, error) {
	var (
		userNotificationViewList []*entity.UserNotificationView
	)

	if len(notificationIDList) < 1 {
		return userNotificationViewList, nil
	}

	notificationIDListStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(notificationIDList)), ","), "[]")

	query := fmt.Sprintf(`
		SELECT
			*
		FROM
			user_notification_views
		WHERE
			notification_id IN (%s)
		ORDER BY
			id DESC
	`, notificationIDListStr)

	err := repo.executer.Select(
		repo.Name+".GetByNotificationID",
		&userNotificationViewList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return userNotificationViewList, nil
}

// 全てのお知らせの既読を取得
func (repo *UserNotificationViewRepositoryImpl) All() ([]*entity.UserNotificationView, error) {
	var (
		userNotificationViewList []*entity.UserNotificationView
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&userNotificationViewList, `
			SELECT 
				*
			FROM
				user_notification_views
			ORDER BY
				id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return userNotificationViewList, nil
}
