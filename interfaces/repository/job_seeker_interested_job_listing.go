package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type JobSeekerInterestedJobListingRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobSeekerInterestedJobListingRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobSeekerInterestedJobListingRepository {
	return &JobSeekerInterestedJobListingRepositoryImpl{
		Name:     "JobSeekerInterestedJobListingRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成
//
func (repo *JobSeekerInterestedJobListingRepositoryImpl) Create(interestedJobListing *entity.JobSeekerInterestedJobListing) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`
			INSERT INTO job_seeker_interested_job_listings (
				uuid,
				job_seeker_id,
				job_information_id,
				interested_type,
				created_at,
				updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?
			)
		`,
		utility.CreateUUID(),
		interestedJobListing.JobSeekerID,
		interestedJobListing.JobInformationID,
		interestedJobListing.InterestedType,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	interestedJobListing.ID = uint(lastID)

	return nil
}

/****************************************************************************************/
/// 削除
//
func (repo *JobSeekerInterestedJobListingRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByJobSeekerID",
		`
		DELETE
		FROM job_seeker_interested_job_listings
		WHERE id = ?
		`, id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

/****************************************************************************************/
/// 複数取得
//
func (repo *JobSeekerInterestedJobListingRepositoryImpl) GetByJobSeekerUUIDAndInterestedType(jobSeekerUUID uuid.UUID, interestedType entity.InterestedType) ([]*entity.JobSeekerInterestedJobListing, error) {
	var (
		interestedJobListingList []*entity.JobSeekerInterestedJobListing
	)

	err := repo.executer.Select(
		repo.Name+".GetByJobSeekerUUIDAndInterestedType",
		&interestedJobListingList, `
		SELECT *
		FROM job_seeker_interested_job_listings
		WHERE
			job_seeker_id = (
				SELECT id
				FROM job_seekers
				WHERE uuid = ?
			)
		AND 
			interested_type = ?
		`,
		jobSeekerUUID,
		interestedType,
	)

	if err != nil {
		fmt.Println(err)
		return interestedJobListingList, err
	}

	return interestedJobListingList, nil
}
