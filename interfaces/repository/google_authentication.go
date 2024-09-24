package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type GoogleAuthenticationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewGoogleAuthenticationRepositoryImpl(ex interfaces.SQLExecuter) usecase.GoogleAuthenticationRepository {
	return &GoogleAuthenticationRepositoryImpl{
		Name:     "GoogleAuthenticationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
//タスクの作成
func (repo *GoogleAuthenticationRepositoryImpl) Create(googleAuth *entity.GoogleAuthentication) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO google_authentication (
				acess_token,
				refresh_token,
				expiry,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		googleAuth.AccessToken,
		googleAuth.RefreshToken,
		googleAuth.Expiry,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	googleAuth.ID = uint(lastID)
	return nil
}

func (repo *GoogleAuthenticationRepositoryImpl) Update(id uint, googleAuth *entity.GoogleAuthentication) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
			UPDATE google_authentication
			SET
				acess_token = ?,
				refresh_token = ?,
				expiry = ?,
				updated_at = ?
			WHERE 
				id = ?
		`,
		googleAuth.AccessToken,
		googleAuth.RefreshToken,
		googleAuth.Expiry,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *GoogleAuthenticationRepositoryImpl) UpdateTokenAndExpiry(id uint, accessToken string, expiry time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateTokenAndExpiry",
		`
			UPDATE google_authentication
			SET
				acess_token = ?,
				expiry = ?,
				updated_at = ?
			WHERE 
				id = ?
		`,
		accessToken,
		expiry,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *GoogleAuthenticationRepositoryImpl) FindLatest() (*entity.GoogleAuthentication, error) {
	var (
		googleAuth entity.GoogleAuthentication
	)

	err := repo.executer.Get(
		repo.Name+".FindLatest",
		&googleAuth, `
			SELECT *
			FROM google_authentication
			ORDER BY updated_at DESC
			LIMIT 1
		`,
	)

	if err != nil {
		return nil, err
	}

	return &googleAuth, nil
}
