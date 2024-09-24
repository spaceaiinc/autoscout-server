package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JobSeekerDepartmentHistory struct {
	ID               uint      `db:"id" json:"id"`
	WorkHistoryID    uint      `db:"work_history_id" json:"work_history_id"`
	Department       string    `db:"department" json:"department"`
	ManagementNumber null.Int  `db:"management_number" json:"management_number"`
	ManagementDetail string    `db:"management_detail" json:"management_detail"`
	JobDescription   string    `db:"job_description" json:"job_description"`
	StartYear        string    `db:"start_year" json:"start_year"`
	EndYear          string    `db:"end_year" json:"end_year"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	UpdatedAt        time.Time `db:"updated_at" json:"-"`

	ExperienceOccupations []JobSeekerExperienceOccupation `db:"experience_occupations" json:"experience_occupations"`
}

func NewJobSeekerDepartmentHistory(
	workHistoryID uint,
	department string,
	managementNumber null.Int,
	managementDetail string,
	jobDescription string,
	startYear string,
	endYear string,
) *JobSeekerDepartmentHistory {
	return &JobSeekerDepartmentHistory{
		WorkHistoryID:    workHistoryID,
		Department:       department,
		ManagementNumber: managementNumber,
		ManagementDetail: managementDetail,
		JobDescription:   jobDescription,
		StartYear:        startYear,
		EndYear:          endYear,
	}
}
