package repository

import (
	"fmt"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingEnterpriseReferenceMaterialRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingEnterpriseReferenceMaterialRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingEnterpriseReferenceMaterialRepository {
	return &SendingEnterpriseReferenceMaterialRepositoryImpl{
		Name:     "SendingEnterpriseReferenceMaterialRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *SendingEnterpriseReferenceMaterialRepositoryImpl) Create(referenceMaterial *entity.SendingEnterpriseReferenceMaterial) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO sending_enterprise_reference_materials (
				sending_enterprise_id,
				reference1_pdf_url,
				reference2_pdf_url,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		referenceMaterial.SendingEnterpriseID,
		referenceMaterial.Reference1PDFURL,
		referenceMaterial.Reference2PDFURL,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	referenceMaterial.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 更新 API
//
func (repo *SendingEnterpriseReferenceMaterialRepositoryImpl) Update(referenceMaterial *entity.SendingEnterpriseReferenceMaterial) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
			UPDATE 
				sending_enterprise_reference_materials 
			SET
				reference1_pdf_url = ?,
				reference2_pdf_url = ?,
				updated_at = ?
			WHERE sending_enterprise_id = ?
		`,
		referenceMaterial.Reference1PDFURL,
		referenceMaterial.Reference2PDFURL,
		time.Now().In(time.UTC),
		referenceMaterial.SendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (repo *SendingEnterpriseReferenceMaterialRepositoryImpl) UpdateReferenceURLBySendingEnterpriseIDAndMaterialType(sendingEnterpriseID uint, fileType string) error {

	var pdfColumnName string
	if fileType == "参考資料1" {
		pdfColumnName = "reference1_pdf_url"
	} else if fileType == "参考資料2" {
		pdfColumnName = "reference2_pdf_url"
	}

	fmt.Println(pdfColumnName)
	query := fmt.Sprintf(`
    UPDATE 
        sending_enterprise_reference_materials 
    SET 
        %s = "" ,
        updated_at = ?
    WHERE 
        sending_enterprise_id = ?
    `, pdfColumnName)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateReferenceURLBySendingEnterpriseIDAndMaterialType",
		query,
		time.Now().In(time.UTC),
		sendingEnterpriseID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *SendingEnterpriseReferenceMaterialRepositoryImpl) FindBySendingEnterpriseID(sendingEnterpriseID uint) (*entity.SendingEnterpriseReferenceMaterial, error) {
	var (
		referenceMaterial entity.SendingEnterpriseReferenceMaterial
	)

	err := repo.executer.Get(
		repo.Name+".FindBySendingEnterpriseID",
		&referenceMaterial, `
		SELECT *
		FROM sending_enterprise_reference_materials
		WHERE
			sending_enterprise_id = ?
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return &referenceMaterial, err
	}

	return &referenceMaterial, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// すべての企業情報を取得
func (repo *SendingEnterpriseReferenceMaterialRepositoryImpl) All() ([]*entity.SendingEnterpriseReferenceMaterial, error) {
	var (
		referenceMaterialList []*entity.SendingEnterpriseReferenceMaterial
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&referenceMaterialList, `
							SELECT *
							FROM sending_enterprise_reference_materials
						`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return referenceMaterialList, nil
}
