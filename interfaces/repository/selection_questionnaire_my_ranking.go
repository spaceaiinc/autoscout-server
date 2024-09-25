package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SelectionQuestionnaireMyRankingRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSelectionQuestionnaireMyRankingRepositoryImpl(ex interfaces.SQLExecuter) usecase.SelectionQuestionnaireMyRankingRepository {
	return &SelectionQuestionnaireMyRankingRepositoryImpl{
		Name:     "SelectionQuestionnaireMyRankingRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// 選考後アンケートの希望順位を作成
func (repo *SelectionQuestionnaireMyRankingRepositoryImpl) Create(selectionQuestionnaireMyRanking *entity.SelectionQuestionnaireMyRanking) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO selection_questionnaire_my_rankings (
			selection_questionnaire_id,	                    
			ranking,
			company_name,	                       
			phase, 	                            
			selection_date,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?
			)`,
		selectionQuestionnaireMyRanking.SelectionQuestionnaireID,
		selectionQuestionnaireMyRanking.Rank,
		selectionQuestionnaireMyRanking.CompanyName,
		selectionQuestionnaireMyRanking.Phase,
		selectionQuestionnaireMyRanking.SelectionDate,
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	selectionQuestionnaireMyRanking.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 削除 API
//
// 選考後アンケートの希望順位を削除
func (repo *SelectionQuestionnaireMyRankingRepositoryImpl) DeleteByQuestionnaireID(questionnaireID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByQuestionnaireID",
		`
			DELETE 
			FROM selection_questionnaire_my_rankings
			WHERE selection_questionnaire_id = ?
			`,
		questionnaireID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得 API
//
func (repo *SelectionQuestionnaireMyRankingRepositoryImpl) GetByQuestionnaireID(questionnaireID uint) ([]*entity.SelectionQuestionnaireMyRanking, error) {
	var (
		list []*entity.SelectionQuestionnaireMyRanking
	)
	err := repo.executer.Select(
		repo.Name+".GetByQuestionnaireID",
		&list, `
		SELECT 
			myrank.*
		FROM 
			selection_questionnaire_my_rankings AS myrank
		WHERE
			myrank.selection_questionnaire_id = ?
		`,
		questionnaireID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

func (repo *SelectionQuestionnaireMyRankingRepositoryImpl) GetByQuestionnaireIDList(idList []uint) ([]*entity.SelectionQuestionnaireMyRanking, error) {
	var (
		list []*entity.SelectionQuestionnaireMyRanking
	)

	if len(idList) == 0 {
		return list, nil
	}

	query := fmt.Sprintf(`
	SELECT 
		myrank.*
	FROM 
		selection_questionnaire_my_rankings AS myrank
	WHERE
		myrank.selection_questionnaire_id IN (%s)
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetByQuestionnaireIDList",
		&list,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}
