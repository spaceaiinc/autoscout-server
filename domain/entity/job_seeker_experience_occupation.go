package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerExperienceOccupation struct {
	ID                  uint      `db:"id" json:"id"`
	DepartmentHistoryID uint      `db:"department_history_id" json:"department_history_id"`
	Occupation          null.Int  `db:"occupation" json:"occupation"`
	CreatedAt           time.Time `db:"created_at" json:"-"`
	UpdatedAt           time.Time `db:"updated_at" json:"-"`
}

func NewJobSeekerExperienceOccupation(
	departmentHistoryID uint,
	occupation null.Int,
) *JobSeekerExperienceOccupation {
	return &JobSeekerExperienceOccupation{
		DepartmentHistoryID: departmentHistoryID,
		Occupation:          occupation,
	}
}
