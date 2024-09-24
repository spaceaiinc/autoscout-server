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

type JobSeekerDocumentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerDocumentRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerDocumentRepository {
	return &JobSeekerDocumentRepositoryImpl{
		Name:     "JobSeekerDocumentRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerDocumentRepositoryImpl) Create(document *entity.JobSeekerDocument) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_documents (
				job_seeker_id,
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
		document.JobSeekerID,
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

/****************************************************************************************/
/// 更新
//
func (repo *JobSeekerDocumentRepositoryImpl) UpdateByJobSeekerID(jobSeekerID uint, document *entity.JobSeekerDocument) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
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
			WHERE job_seeker_id = ?
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
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 履歴書(原本)の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateResumeOriginURLByJobSeekerID(jobSeekerID uint, resumeOriginURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateResumeOriginURLByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
				resume_origin_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		resumeOriginURL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 職務経歴書(原本)の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateCVOriginURLByJobSeekerID(jobSeekerID uint, cvOriginURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateCVOriginURLByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
				cv_origin_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		cvOriginURL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 履歴書(PDF)の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateResumePDFURLByJobSeekerID(jobSeekerID uint, resumePDFURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateResumePDFURLByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
				resume_pdf_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		resumePDFURL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 職務経歴書(原本)の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateCVPDFURLByJobSeekerID(jobSeekerID uint, cvPDFURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateCVPDFURLByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
				cv_pdf_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		cvPDFURL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 推薦状(PDF)の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateRecommendationPDFURLByJobSeekerID(jobSeekerID uint, recommendationPDFURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateRecommendationPDFURLByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
				recommendation_pdf_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		recommendationPDFURL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 推薦状(原本)の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateRecommendationOriginURLByJobSeekerID(jobSeekerID uint, recommendationOriginURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateRecommendationOriginURLByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
				recommendation_origin_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		recommendationOriginURL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 証明写真の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateIDPhotoURLByJobSeekerID(jobSeekerID uint, idPhotoURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIDPhotoURLByJobSeekerID",
		`
			UPDATE job_seeker_documents SET
				id_photo_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		idPhotoURL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// その他①の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateOtherDocument1URLByJobSeekerID(jobSeekerID uint, otherDocument1URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateOtherDocument1URLByJobSeekerID",
		`
			UPDATE job_seeker_documents 
			SET
				other_document1_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		otherDocument1URL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// その他②の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateOtherDocument2URLByJobSeekerID(jobSeekerID uint, otherDocument2URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateOtherDocument2URLByJobSeekerID",
		`
			UPDATE job_seeker_documents 
			SET
				other_document2_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		otherDocument2URL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// その他③の更新
func (repo *JobSeekerDocumentRepositoryImpl) UpdateOtherDocument3URLByJobSeekerID(jobSeekerID uint, otherDocument3URL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateOtherDocument3URLByJobSeekerID",
		`
			UPDATE job_seeker_documents 
			SET
				other_document3_url = ?,
				updated_at = ?
			WHERE job_seeker_id = ?
		`,
		otherDocument3URL,
		time.Now().In(time.UTC),
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerDocumentRepositoryImpl) DeleteByJobSeekerID(jobSeekerID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_documents
		WHERE job_seeker_id = ?
		`, jobSeekerID,
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
func (repo *JobSeekerDocumentRepositoryImpl) FindByJobSeekerID(jobSeekerID uint) (*entity.JobSeekerDocument, error) {
	var document entity.JobSeekerDocument

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerID",
		&document, `
		SELECT *
		FROM job_seeker_documents
		WHERE
			job_seeker_id = ?
		LIMIT 1
		`,
		jobSeekerID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &document, nil
}

func (repo *JobSeekerDocumentRepositoryImpl) FindByJobSeekerUUID(jobSeekerUUID uuid.UUID) (*entity.JobSeekerDocument, error) {
	var document entity.JobSeekerDocument

	err := repo.executer.Get(
		repo.Name+".FindByJobSeekerUUID",
		&document, `
		SELECT 
			document.*, seeker.last_name, seeker.first_name, seeker.agent_id
		FROM 
			job_seeker_documents AS document
		INNER JOIN
			job_seekers AS seeker
		ON
			document.job_seeker_id = seeker.id
		WHERE
			seeker.uuid = ?
		LIMIT 1
		`,
		jobSeekerUUID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &document, nil
}

/****************************************************************************************/
/// 複数取得
//
// エージェントIDから求職者一覧を取得
func (repo *JobSeekerDocumentRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobSeekerDocument, error) {
	var documentList []*entity.JobSeekerDocument

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&documentList, `
			SELECT 
				jsd.*
			FROM 
				job_seeker_documents AS jsd
			INNER JOIN
				job_seekers AS js
			ON
				jsd.job_seeker_id = js.id
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
func (repo *JobSeekerDocumentRepositoryImpl) GetByStaffID(staffID uint) ([]*entity.JobSeekerDocument, error) {
	var documentList []*entity.JobSeekerDocument

	err := repo.executer.Select(
		repo.Name+".GetByStaffID",
		&documentList, `
			SELECT 
				jsd.*
			FROM 
				job_seeker_documents AS jsd
			INNER JOIN
				job_seekers AS js
			ON
				jsd.job_seeker_id = js.id
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
func (repo *JobSeekerDocumentRepositoryImpl) GetByJobSeekerIDList(idList []uint) ([]*entity.JobSeekerDocument, error) {
	var documentList []*entity.JobSeekerDocument

	if len(idList) == 0 {
		return documentList, nil
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM job_seeker_documents
		WHERE
			job_seeker_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerIDList",
		&documentList,
		query,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return documentList, nil
}

// すべての求職者情報を取得
func (repo *JobSeekerDocumentRepositoryImpl) All() ([]*entity.JobSeekerDocument, error) {
	var documentList []*entity.JobSeekerDocument

	err := repo.executer.Select(
		repo.Name+".All",
		&documentList, `
			SELECT *
			FROM job_seeker_documents
		`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return documentList, nil
}
