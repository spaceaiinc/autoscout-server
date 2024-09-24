package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SendingJobSeekerWorkHistory struct {
	ID                   uint      `db:"id" json:"id"`
	SendingJobSeekerID   uint      `db:"sending_job_seeker_id" json:"sending_job_seeker_id"`
	CompanyName          string    `db:"company_name" json:"company_name"`
	EmployeeNumberSingle null.Int  `db:"employee_number_single" json:"employee_number_single"`
	EmployeeNumberGroup  null.Int  `db:"employee_number_group" json:"employee_number_group"`
	PublicOffering       null.Int  `db:"public_offering" json:"public_offering"`
	JoiningYear          string    `db:"joining_year" json:"joining_year"`
	EmploymentStatus     null.Int  `db:"employment_status" json:"employment_status"`
	RetireReasonOfTruth  string    `db:"retire_reason_of_truth" json:"retire_reason_of_truth"`
	RetireReasonOfPublic string    `db:"retire_reason_of_public" json:"retire_reason_of_public"`
	RetireYear           string    `db:"retire_year" json:"retire_year"`
	FirstStatus          null.Int  `db:"first_status" json:"first_status"`
	LastStatus           null.Int  `db:"last_status" json:"last_status"`
	CreatedAt            time.Time `db:"created_at" json:"-"`
	UpdatedAt            time.Time `db:"updated_at" json:"-"`

	//他テーブル
	//経験業界
	ExperienceIndustries []SendingJobSeekerExperienceIndustry `db:"experience_industries" json:"experience_industries"`

	//経験職種
	DepartmentHistories []SendingJobSeekerDepartmentHistory `json:"department_histories"` // 求職者の部署履歴
}

func NewSendingJobSeekerWorkHistory(
	sendingJobSeekerID uint,
	companyName string,
	employeeNumberSingle null.Int,
	employeeNumberGroup null.Int,
	publicOffering null.Int,
	joiningYear string,
	employmentStatus null.Int,
	retireReasonOfTruth string,
	retireReasonOfPublic string,
	retireYear string,
	firstStatus null.Int,
	lastStatus null.Int,
) *SendingJobSeekerWorkHistory {
	return &SendingJobSeekerWorkHistory{
		SendingJobSeekerID:   sendingJobSeekerID,
		CompanyName:          companyName,
		EmployeeNumberSingle: employeeNumberSingle,
		EmployeeNumberGroup:  employeeNumberGroup,
		PublicOffering:       publicOffering,
		JoiningYear:          joiningYear,
		EmploymentStatus:     employmentStatus,
		RetireReasonOfTruth:  retireReasonOfTruth,
		RetireReasonOfPublic: retireReasonOfPublic,
		RetireYear:           retireYear,
		FirstStatus:          firstStatus,
		LastStatus:           lastStatus,
	}
}
