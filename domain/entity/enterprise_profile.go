package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type EnterpriseProfile struct {
	ID                   uint      `db:"id" json:"id"`
	UUID                 uuid.UUID `db:"uuid" json:"uuid"`
	CompanyName          string    `db:"company_name" json:"company_name"`
	AgentStaffID         uint      `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName            string    `db:"staff_name" json:"staff_name"`
	CorporateSiteURL     string    `db:"corporate_site_url" json:"corporate_site_url"`
	Representative       string    `db:"representative" json:"representative"`
	Establishment        string    `db:"establishment" json:"establishment"`
	PostCode             string    `db:"post_code" json:"post_code"`
	OfficeLocation       string    `db:"office_location" json:"office_location"`
	EmployeeNumberSingle null.Int  `db:"employee_number_single" json:"employee_number_single"`
	EmployeeNumberGroup  null.Int  `db:"employee_number_group" json:"employee_number_group"`
	Capital              string    `db:"capital" json:"capital"`
	PublicOffering       null.Int  `db:"public_offering" json:"public_offering"`
	EarningsYear         null.Int  `db:"earnings_year" json:"earnings_year"`
	Earnings             string    `db:"earnings" json:"earnings"`
	BusinessDetail       string    `db:"business_detail" json:"business_detail"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`

	// 他テーブル
	AgentID             uint                        `db:"agent_id" json:"agent_id"`
	Industries          []null.Int                  `db:"industries" json:"industries"`
	ReferenceMaterialID uint                        `json:"reference_material_id"`
	ReferenceMaterial   EnterpriseReferenceMaterial `json:"reference_materials"`
	Activities          []EnterpriseActivity        `json:"activities"`
}

func NewEnterpriseProfile(
	companyName string,
	agentStaffID uint,
	corporateSiteURL string,
	representative string,
	establishment string,
	postCode string,
	officeLocation string,
	employeeNumberSingle null.Int,
	employeeNumberGroup null.Int,
	capital string,
	publicOffering null.Int,
	earningsYear null.Int,
	earnings string,
	businessDetail string,
) *EnterpriseProfile {
	return &EnterpriseProfile{
		CompanyName:          companyName,
		AgentStaffID:         agentStaffID,
		CorporateSiteURL:     corporateSiteURL,
		Representative:       representative,
		Establishment:        establishment,
		PostCode:             postCode,
		OfficeLocation:       officeLocation,
		EmployeeNumberSingle: employeeNumberSingle,
		EmployeeNumberGroup:  employeeNumberGroup,
		Capital:              capital,
		PublicOffering:       publicOffering,
		EarningsYear:         earningsYear,
		Earnings:             earnings,
		BusinessDetail:       businessDetail,
	}
}

type CreateOrUpdateEnterpriseProfileParam struct {
	CompanyName          string     `json:"company_name" validate:"required"`
	AgentStaffID         uint       `json:"agent_staff_id" validate:"required"`
	CorporateSiteURL     string     `json:"corporate_site_url"`
	Representative       string     `json:"representative"`
	Establishment        string     `json:"establishment"`
	PostCode             string     `json:"post_code"`
	OfficeLocation       string     `json:"office_location"`
	EmployeeNumberSingle null.Int   `json:"employee_number_single"`
	EmployeeNumberGroup  null.Int   `json:"employee_number_group"`
	Capital              string     `json:"capital"`
	PublicOffering       null.Int   `json:"public_offering"`
	EarningsYear         null.Int   `json:"earnings_year"`
	Earnings             string     `json:"earnings"`
	BusinessDetail       string     `json:"business_detail"`
	Industries           []null.Int `json:"industries"`
}

type DeleteEnterpriseProfileParam struct {
	ID uint `json:"id" validate:"required"`
}

type SearchEnterprise struct {
	FreeWord          string
	AgentStaffID      string
	Industries        []null.Int
	Prefectures       []null.Int
	CompanyScaleTypes []null.Int
}

func NewSearchEnterprise(
	freeword string,
	agentStaffID string,
	industries []null.Int,
	companyScaleTypes []null.Int,
) *SearchEnterprise {
	return &SearchEnterprise{
		FreeWord:          freeword,
		AgentStaffID:      agentStaffID,
		Industries:        industries,
		CompanyScaleTypes: companyScaleTypes,
	}
}
