package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingShareDocumentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingShareDocumentRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingShareDocumentRepository {
	return &SendingShareDocumentRepositoryImpl{
		Name:     "SendingShareDocumentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
// 汎用系 API
//
func (repo *SendingShareDocumentRepositoryImpl) Create(sendingShareDocument *entity.SendingShareDocument) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sending_share_documents (
			sending_job_seeker_id,
			sending_enterprise_id,
			is_share_upload_resume,
			is_share_upload_cv,
			is_share_upload_recommendation,
			is_share_generated_resume,
			is_share_generated_cv,
			is_share_generated_recommendation,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
		sendingShareDocument.SendingJobSeekerID,
		sendingShareDocument.SendingEnterpriseID,
		sendingShareDocument.IsShareUploadResume,
		sendingShareDocument.IsShareUploadCV,
		sendingShareDocument.IsShareUploadRecommendation,
		sendingShareDocument.IsShareGeneratedResume,
		sendingShareDocument.IsShareGeneratedCV,
		sendingShareDocument.IsShareGeneratedRecommendation,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sendingShareDocument.ID = uint(lastID)
	return nil
}

func (repo *SendingShareDocumentRepositoryImpl) Update(id uint, sendingShareDocument *entity.SendingShareDocument) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sending_share_documents
		SET
			is_share_upload_resume = ?,
			is_share_upload_cv = ?,
			is_share_upload_recommendation = ?,
			is_share_generated_resume = ?,
			is_share_generated_cv = ?,
			is_share_generated_recommendation = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		sendingShareDocument.IsShareUploadResume,
		sendingShareDocument.IsShareUploadCV,
		sendingShareDocument.IsShareUploadRecommendation,
		sendingShareDocument.IsShareGeneratedResume,
		sendingShareDocument.IsShareGeneratedCV,
		sendingShareDocument.IsShareGeneratedRecommendation,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingShareDocumentRepositoryImpl) FindBySendingJobSeekerIDAndSendingEnterpriseID(sendingJoSeekerID, sendingEnterpriseID uint) (*entity.SendingShareDocument, error) {
	var (
		sendingShareDocument entity.SendingShareDocument
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingJobSeekerIDAndSendingEnterpriseID",
		&sendingShareDocument, `
		SELECT *
		FROM sending_share_documents
		WHERE 
			sending_job_seeker_id = ? AND
			sending_enterprise_id = ?
		LIMIT 1
		`,
		sendingJoSeekerID, sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &sendingShareDocument, nil
}

func (repo *SendingShareDocumentRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingShareDocument, error) {
	var (
		sendingShareDocumentList []*entity.SendingShareDocument
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&sendingShareDocumentList, `
			SELECT *
			FROM sending_share_documents
			WHERE 
				sending_enterprise_id = ?
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingShareDocumentList, nil
}
