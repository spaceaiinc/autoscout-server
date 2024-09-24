package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobSeekerDocumentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobSeekerDocumentRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobSeekerDocumentRepository {
	return &SendingJobSeekerDocumentRepositoryImpl{
		Name:     "SendingJobSeekerDocumentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//
func (repo *SendingJobSeekerDocumentRepositoryImpl) Create(document *entity.SendingJobSeekerDocument) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_job_seeker_documents (
				sending_job_seeker_id,
				resume_origin_url,
				resume_pdf_url,
				cv_origin_url,
				cv_pdf_url,
				recommendation_origin_url,
				recommendation_pdf_url,
				id_photo_url,
				other_document1_url,
				other_document2_url,
				other_document3_url,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, 
				?, ?, ?
			)
		`,
		document.SendingJobSeekerID,
		document.ResumeOriginURL,
		document.ResumePDFURL,
		document.CVOriginURL,
		document.CVPDFURL,
		document.RecommendationOriginURL,
		document.RecommendationPDFURL,
		document.IDPhotoURL,
		document.OtherDocument1URL,
		document.OtherDocument2URL,
		document.OtherDocument3URL,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	document.ID = uint(lastID)

	return nil
}

func (repo *SendingJobSeekerDocumentRepositoryImpl) Update(document *entity.SendingJobSeekerDocument) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
			UPDATE sending_job_seeker_documents SET
				resume_origin_url = ?,
				resume_pdf_url = ?,
				cv_origin_url = ?,
				cv_pdf_url = ?,
				recommendation_origin_url = ?,
				recommendation_pdf_url = ?,
				id_photo_url = ?,
				other_document1_url = ?,
				other_document2_url = ?,
				other_document3_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		document.ResumeOriginURL,
		document.ResumePDFURL,
		document.CVOriginURL,
		document.CVPDFURL,
		document.RecommendationOriginURL,
		document.RecommendationPDFURL,
		document.IDPhotoURL,
		document.OtherDocument1URL,
		document.OtherDocument2URL,
		document.OtherDocument3URL,
		time.Now().In(time.UTC),
		document.SendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateOrigin(document *entity.SendingJobSeekerDocument) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateOrigin",
		`
			UPDATE sending_job_seeker_documents SET
				resume_origin_url = ?,
				cv_origin_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		document.ResumeOriginURL,
		document.CVOriginURL,
		time.Now().In(time.UTC),
		document.SendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 履歴書(原本)の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateResumeOriginURL(resumeOriginURL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateResumeOriginURL",
		`
			UPDATE sending_job_seeker_documents SET
				resume_origin_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		resumeOriginURL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 職務経歴書(原本)の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateCVOriginURL(cvOriginURL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateCVOriginURL",
		`
			UPDATE sending_job_seeker_documents SET
				cv_origin_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		cvOriginURL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 履歴書(PDF)の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateResumePDFURL(resumePDFURL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateResumePDFURL",
		`
			UPDATE sending_job_seeker_documents SET
				resume_pdf_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		resumePDFURL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 職務経歴書(原本)の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateCVPDFURL(cvPDFURL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateCVPDFURL",
		`
			UPDATE sending_job_seeker_documents SET
				cv_pdf_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		cvPDFURL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 推薦状(PDF)の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateRecommendationPDFURL(recommendationPDFURL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateRecommendationPDFURL",
		`
			UPDATE sending_job_seeker_documents SET
				recommendation_pdf_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		recommendationPDFURL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 推薦状(原本)の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateRecommendationOriginURL(recommendationOriginURL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateRecommendationOriginURL",
		`
			UPDATE sending_job_seeker_documents SET
				recommendation_origin_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		recommendationOriginURL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 証明写真の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateIDPhotoURL(idPhotoURL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIDPhotoURL",
		`
			UPDATE sending_job_seeker_documents SET
				id_photo_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		idPhotoURL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// その他①の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateOtherDocument1URL(otherDocument1URL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateOtherDocument1URL",
		`
			UPDATE sending_job_seeker_documents 
			SET
				other_document1_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		otherDocument1URL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// その他②の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateOtherDocument2URL(otherDocument2URL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateOtherDocument2URL",
		`
			UPDATE sending_job_seeker_documents 
			SET
				other_document2_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		otherDocument2URL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// その他③の更新
func (repo *SendingJobSeekerDocumentRepositoryImpl) UpdateOtherDocument3URL(otherDocument3URL string, sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateOtherDocument3URL",
		`
			UPDATE sending_job_seeker_documents 
			SET
				other_document3_url = ?,
				updated_at = ?
			WHERE sending_job_seeker_id = ?
		`,
		otherDocument3URL,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (repo *SendingJobSeekerDocumentRepositoryImpl) Delete(sendingJobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		DELETE
		FROM sending_job_seeker_documents
		WHERE sending_job_seeker_id = ?
		`, sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobSeekerDocumentRepositoryImpl) FindBySendingJobSeekerID(sendingJobSeekerID uint) (*entity.SendingJobSeekerDocument, error) {
	var (
		document entity.SendingJobSeekerDocument
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingJobSeekerID",
		&document, `
		SELECT *
		FROM sending_job_seeker_documents
		WHERE
			sending_job_seeker_id = ?
		LIMIT 1
		`,
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &document, nil
}

func (repo *SendingJobSeekerDocumentRepositoryImpl) FindBySendingJobSeekerUUID(sendingJobSeekerUUID uuid.UUID) (*entity.SendingJobSeekerDocument, error) {
	var (
		document entity.SendingJobSeekerDocument
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingJobSeekerUUID",
		&document, `
		SELECT 
			document.*, seeker.last_name, seeker.first_name, seeker.agent_id
		FROM 
			sending_job_seeker_documents AS document
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			document.sending_job_seeker_id = seeker.id
		WHERE
			seeker.uuid = ?
		LIMIT 1
		`,
		sendingJobSeekerUUID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &document, nil
}

// エージェントIDから求職者一覧を取得
func (repo *SendingJobSeekerDocumentRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingJobSeekerDocument, error) {
	var (
		documentList []*entity.SendingJobSeekerDocument
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&documentList, `
			SELECT 
				jsd.*
			FROM 
				sending_job_seeker_documents AS jsd
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsd.sending_job_seeker_id = js.id
			WHERE
				js.agent_id = ?
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return documentList, nil
}

// 担当者IDから求職者一覧を取得
func (repo *SendingJobSeekerDocumentRepositoryImpl) GetListByStaffID(staffID uint) ([]*entity.SendingJobSeekerDocument, error) {
	var (
		documentList []*entity.SendingJobSeekerDocument
	)

	err := repo.executer.Select(
		repo.Name+".GetListByStaffID",
		&documentList, `
			SELECT 
				jsd.*
			FROM 
				sending_job_seeker_documents AS jsd
			INNER JOIN
				sending_job_seekers AS js
			ON
				jsd.sending_job_seeker_id = js.id
			WHERE
				js.agent_staff_id = ?
		`,
		staffID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return documentList, nil
}

// 求職者リストから書類を取得
func (repo *SendingJobSeekerDocumentRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobSeekerDocument, error) {
	var (
		documentList []*entity.SendingJobSeekerDocument
	)

	if len(idList) == 0 {
		return documentList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sending_job_seeker_documents
		WHERE
			sending_job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&documentList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return documentList, nil
}

/****************************************************************************************/
/// Admin API
//
//すべての求職者情報を取得
func (repo *SendingJobSeekerDocumentRepositoryImpl) GetAll() ([]*entity.SendingJobSeekerDocument, error) {
	var (
		documentList []*entity.SendingJobSeekerDocument
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&documentList, `
			SELECT *
			FROM sending_job_seeker_documents
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return documentList, nil
}
