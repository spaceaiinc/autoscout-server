package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type EnterpriseReferenceMaterialRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewEnterpriseReferenceMaterialRepositoryImpl(ex interfaces.SQLExecuter) usecase.EnterpriseReferenceMaterialRepository {
	return &EnterpriseReferenceMaterialRepositoryImpl{
		Name:     "EnterpriseReferenceMaterialRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
func (repo *EnterpriseReferenceMaterialRepositoryImpl) Create(referenceMaterial *entity.EnterpriseReferenceMaterial) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO enterprise_reference_materials (
				enterprise_id,
				reference1_pdf_url,
				reference2_pdf_url,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?
			)
		`,
		referenceMaterial.EnterpriseID,
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
func (repo *EnterpriseReferenceMaterialRepositoryImpl) UpdateByEnterpriseID(enterpriseID uint, referenceMaterial *entity.EnterpriseReferenceMaterial) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateByEnterpriseID",
		`
			UPDATE enterprise_reference_materials SET
				reference1_pdf_url = ?,
				reference2_pdf_url = ?,
				updated_at = ?
			WHERE enterprise_id = ?
		`,
		referenceMaterial.Reference1PDFURL,
		referenceMaterial.Reference2PDFURL,
		time.Now().In(time.UTC),
		enterpriseID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (repo *EnterpriseReferenceMaterialRepositoryImpl) UpdateMaterialTypeByID(id, materialType uint) error {
	var pdfColumnName string
	if materialType == 1 {
		pdfColumnName = "reference1_pdf_url"
	} else if materialType == 2 {
		pdfColumnName = "reference2_pdf_url"
	}

	fmt.Println(pdfColumnName)
	query := fmt.Sprintf(`
    UPDATE 
        enterprise_reference_materials 
    SET 
        %s = "" ,
        updated_at = ?
    WHERE 
        id = ?
    `, pdfColumnName)

	_, err := repo.executer.Exec(
		repo.Name+".UpdateMaterialTypeByID",
		query,
		time.Now().In(time.UTC),
		id,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/****************************************************************************************/
/// 削除 API
//

func (repo *EnterpriseReferenceMaterialRepositoryImpl) DeleteByEnterpriseID(enterpriseID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByEnterpriseID",
		`
		DELETE
		FROM enterprise_reference_materials
		WHERE enterprise_id = ?
		`, enterpriseID,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 単数取得 API
//
func (repo *EnterpriseReferenceMaterialRepositoryImpl) FindByEnterpriseID(enterpriseID uint) (*entity.EnterpriseReferenceMaterial, error) {
	var referenceMaterial entity.EnterpriseReferenceMaterial

	err := repo.executer.Get(
		repo.Name+".FindByEnterpriseID",
		&referenceMaterial, `
		SELECT *
		FROM enterprise_reference_materials
		WHERE
			enterprise_id = ?
		`,
		enterpriseID,
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
//担当者IDで企業情報を取得
func (repo *EnterpriseReferenceMaterialRepositoryImpl) GetByAgentStaffID(agentStaffID uint) ([]*entity.EnterpriseReferenceMaterial, error) {
	var referenceMaterialList []*entity.EnterpriseReferenceMaterial

	err := repo.executer.Select(
		repo.Name+".GetByAgentStaffID",
		&referenceMaterialList, `
		SELECT *
		FROM 
		enterprise_reference_materials
		WHERE
		enterprise_id IN (
			SELECT id
			FROM enterprise_profiles
			WHERE agent_staff_id = ?
		)
		`,
		agentStaffID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return referenceMaterialList, nil
}

// エージェントIDから企業一覧を取得
func (repo *EnterpriseReferenceMaterialRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.EnterpriseReferenceMaterial, error) {
	var referenceMaterialList []*entity.EnterpriseReferenceMaterial

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&referenceMaterialList, `
		SELECT *
		FROM 
		enterprise_reference_materials
		WHERE
		enterprise_id IN (
			SELECT id
			FROM enterprise_profiles
			WHERE 
				agent_staff_id IN (
					SELECT id
					FROM agent_staffs
					WHERE
					agent_id = ?
					)
			)
						`,
		agentID,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return referenceMaterialList, nil
}

// 企業IDリストからリストを取得
func (repo *EnterpriseReferenceMaterialRepositoryImpl) GetByEnterpriseIDList(enterpriseIDList []uint) ([]*entity.EnterpriseReferenceMaterial, error) {
	var referenceMaterialList []*entity.EnterpriseReferenceMaterial
	if len(enterpriseIDList) == 0 {
		return referenceMaterialList, nil
	}

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseIDList",
		&referenceMaterialList, `
		SELECT *
		FROM enterprise_reference_materials
		WHERE enterprise_id IN (?)
		`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enterpriseIDList)), ", "), "[]"),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return referenceMaterialList, nil
}

// すべての企業情報を取得
func (repo *EnterpriseReferenceMaterialRepositoryImpl) All() ([]*entity.EnterpriseReferenceMaterial, error) {
	var referenceMaterialList []*entity.EnterpriseReferenceMaterial

	err := repo.executer.Select(
		repo.Name+".All",
		&referenceMaterialList, `
							SELECT *
							FROM enterprise_reference_materials
						`,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return referenceMaterialList, nil
}
