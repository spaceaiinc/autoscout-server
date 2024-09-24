package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// 送客先エージェントを管理するテーブル
type SendingEnterprise struct {
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
	PublicOffering       null.Int  `db:"public_offering" json:"public_offering"`
	SendingTarget        string    `db:"sending_target" json:"sending_target"`
	SendingRequiredInfo  string    `db:"sending_required_info" json:"sending_required_info"`
	Remarks              string    `db:"remarks" json:"remarks"`
	IsDeleted            bool      `db:"is_deleted" json:"is_deleted"`
	Password             string    `db:"password" json:"password"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`

	// 他テーブル
	AgentID             uint                               `db:"agent_id" json:"agent_id"`
	ReferenceMaterialID uint                               `json:"reference_material_id"`
	ReferenceMaterial   SendingEnterpriseReferenceMaterial `json:"reference_materials"`

	// 送客先エージェントの「エージェントの特徴」を管理するテーブル
	Speciality SendingEnterpriseSpeciality `json:"speciality"`

	// 保有求人数
	JobInformationCount uint `db:"job_information_count" json:"job_information_count"`

	// 請求先
	SendingBillingAddress SendingBillingAddress `json:"sending_billing_address"`

	// 送客応諾になっている求職者リスト
	AcceptSendingPhaseList []SendingPhase `db:"accept_sending_phase_list" json:"accept_sending_phase_list"`
}

func NewSendingEnterprise(
	companyName string,
	agentStaffID uint,
	corporateSiteURL string,
	representative string,
	establishment string,
	postCode string,
	officeLocation string,
	employeeNumberSingle null.Int,
	publicOffering null.Int,
	sendingTarget string,
	sendingRequiredInfo string,
	remarks string,
	password string,
) *SendingEnterprise {
	return &SendingEnterprise{
		CompanyName:          companyName,
		AgentStaffID:         agentStaffID,
		CorporateSiteURL:     corporateSiteURL,
		Representative:       representative,
		Establishment:        establishment,
		PostCode:             postCode,
		OfficeLocation:       officeLocation,
		EmployeeNumberSingle: employeeNumberSingle,
		PublicOffering:       publicOffering,
		SendingTarget:        sendingTarget,
		SendingRequiredInfo:  sendingRequiredInfo,
		Remarks:              remarks,
		Password:             password,
	}
}

type CreateOrUpdateSendingEnterpriseParam struct {
	CompanyName          string   `json:"company_name" validate:"required"`
	AgentStaffID         uint     `json:"agent_staff_id" validate:"required"`
	CorporateSiteURL     string   `json:"corporate_site_url"`
	Representative       string   `json:"representative"`
	Establishment        string   `json:"establishment"`
	PostCode             string   `json:"post_code"`
	OfficeLocation       string   `json:"office_location"`
	EmployeeNumberSingle null.Int `json:"employee_number_single"`
	PublicOffering       null.Int `json:"public_offering"`
	SendingTarget        string   `db:"sending_target" json:"sending_target"`
	SendingRequiredInfo  string   `db:"sending_required_info" json:"sending_required_info"`
	Remarks              string   `db:"remarks" json:"remarks"`
	Password             string   `json:"password"`

	Speciality SendingEnterpriseSpeciality `json:"speciality"`
}

type UpdateSendingEnterprisePasswordParam struct {
	Password string `json:"password" validate:"required"`
}

type SigninSendingEnterprisePasswordParam struct {
	SendingEnterpriseUUID string `json:"sending_enterprise_uuid" validate:"required"`
	SendingJobSeekerUUID  string `json:"sending_job_seeker_uuid" validate:"required"`
	Password              string `json:"password" validate:"required"`
}

type SendSendingMailParam struct {
	SendingJobSeekerID       uint          `db:"sending_job_seeker_id" json:"sending_job_seeker_id" form:"sending_job_seeker_id" param:"sending_job_seeker_id" query:"sending_job_seeker_id"`
	SendingList              []SendingInfo `db:"sending_list" json:"sending_list" form:"sending_list" param:"sending_list" query:"sending_list"`
	InterestingJobInfoIDList []uint        `db:"interesting_job_info_id_list" json:"interesting_job_info_id_list" form:"interesting_job_info_id_list" param:"interesting_job_info_id_list" query:"interesting_job_info_id_list"`
}

type SendingInfo struct {
	SendingEnterpriseID            uint             `db:"sending_enterprise_id" json:"sending_enterprise_id" form:"sending_enterprise_id" param:"sending_enterprise_id" query:"sending_enterprise_id"`
	SendingJobInformationList      []SendingJobInfo `db:"sending_job_information_list" json:"sending_job_information_list" form:"sending_job_information_list" param:"sending_job_information_list" query:"sending_job_information_list"`
	SendingDate                    time.Time        `db:"sending_date" json:"sending_date" form:"sending_date" param:"sending_date" query:"sending_date"`
	Mail                           SendMail         `db:"mail" json:"mail" form:"mail" param:"mail" query:"mail"`
	IsShareUploadResume            bool             `db:"is_share_upload_resume" json:"is_share_upload_resume"`                       // アップロードした履歴書のシェア
	IsShareUploadCV                bool             `db:"is_share_upload_cv" json:"is_share_upload_cv"`                               // アップロードした職務経歴書のシェア
	IsShareUploadRecommendation    bool             `db:"is_share_upload_recommendation" json:"is_share_upload_recommendation"`       // アップロードした推薦状のシェア
	IsShareGeneratedResume         bool             `db:"is_share_generated_resume" json:"is_share_generated_resume"`                 // 自動生成した履歴書のシェア
	IsShareGeneratedCV             bool             `db:"is_share_generated_cv" json:"is_share_generated_cv"`                         // 自動生成した職務経歴書のシェア
	IsShareGeneratedRecommendation bool             `db:"is_share_generated_recommendation" json:"is_share_generated_recommendation"` // 自動生成した推薦状のシェア
}

type SendingJobInfo struct {
	ID          uint   `db:"id" json:"id" form:"id" param:"id" query:"id"`
	CompanyName string `db:"company_name" json:"company_name" form:"company_name" param:"company_name" query:"company_name"`
	Title       string `db:"title" json:"title" form:"title" param:"title" query:"title"`
}

type SendMail struct {
	Tos     []EmailUser `json:"tos" form:"tos" param:"tos" query:"tos"`                 // 送信先の情報（）
	Subject string      `json:"subject" form:"subject" param:"subject" query:"subject"` // メールの件名
	Content string      `json:"content" form:"content" param:"content" query:"content"` // メールの本文
}

type SendSendingMailForRSVPParam struct {
	SendingList        []SendingInfo `db:"sending_list" json:"sending_list" form:"sending_list" param:"sending_list" query:"sending_list"`
	SendingPhaseIDList []uint        `db:"sending_phase_id_list" json:"sending_phase_id_list" form:"sending_phase_id_list" param:"sending_phase_id_list" query:"sending_phase_id_list"`
}

type SendingInfoForRSVP struct {
	SendingEnterpriseID uint     `db:"sending_enterprise_id" json:"sending_enterprise_id" form:"sending_enterprise_id" param:"sending_enterprise_id" query:"sending_enterprise_id"`
	Mail                SendMail `db:"mail" json:"mail" form:"mail" param:"mail" query:"mail"`
}
