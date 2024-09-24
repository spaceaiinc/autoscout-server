package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SendingJobInformationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingJobInformationRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingJobInformationRepository {
	return &SendingJobInformationRepositoryImpl{
		Name:     "SendingJobInformationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 汎用系 API
//

func (repo *SendingJobInformationRepositoryImpl) Create(sendingJobInformation *entity.SendingJobInformation) error {
	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sending_job_informations (
			uuid,
			sending_billing_address_id,
			company_name,
			title,
			recruitment_state,
			expiration_date,
			background,
			work_detail,
			number_of_hires,
			work_location,

			transfer,
			transfer_detail,
			under_income,
			over_income,
			salary,
			insurance,
			work_time,
			overtime,
			overtime_average,
			fixed_overtime,

			fixed_overtime_payment,
			fixed_overtime_detail,
			trial_period,
			trial_period_detail,
			employment_period,
			employment_period_detail,
			holiday_type,
			holiday_detail,
			passive_smoking,
			selection_flow,


			gender,
			nationality,
			final_education,
			school_level,
			medical_history,
			age_under,
			age_over,
			job_change,
			short_resignation,
			short_resignation_remarks,


			social_experience_year,
			social_experience_month,
			other_required,
			appearance,
			communication,
			thinking,
			target_detail,
			required_experience_job_detail,
			employment_insurance,
			accident_insurance,

			health_insurance,
			pension_insurance,
			register_phase,
			study_category,
			word_skill,
			excel_skill,
			power_point_skill,
			corporate_site_url,
			post_code,
			office_location,
			
			establishment,
			employee_number_single,
			employee_number_group,
			public_offering,
			earnings_year,
			earnings,
			business_detail,
			is_deleted,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
		`,
		utility.CreateUUID(),
		sendingJobInformation.SendingBillingAddressID,
		sendingJobInformation.CompanyName,
		sendingJobInformation.Title,
		sendingJobInformation.RecruitmentState,
		sendingJobInformation.ExpirationDate,
		sendingJobInformation.Background,
		sendingJobInformation.WorkDetail,
		sendingJobInformation.NumberOfHires,
		sendingJobInformation.WorkLocation,
		sendingJobInformation.Transfer,
		sendingJobInformation.TransferDetail,
		sendingJobInformation.UnderIncome,
		sendingJobInformation.OverIncome,
		sendingJobInformation.Salary,
		sendingJobInformation.Insurance,
		sendingJobInformation.WorkTime,
		sendingJobInformation.Overtime,
		sendingJobInformation.OvertimeAverage,
		sendingJobInformation.FixedOvertime,
		sendingJobInformation.FixedOvertimePayment,
		sendingJobInformation.FixedOvertimeDetail,
		sendingJobInformation.TrialPeriod,
		sendingJobInformation.TrialPeriodDetail,
		sendingJobInformation.EmploymentPeriod,
		sendingJobInformation.EmploymentPeriodDetail,
		sendingJobInformation.HolidayType,
		sendingJobInformation.HolidayDetail,
		sendingJobInformation.PassiveSmoking,
		sendingJobInformation.SelectionFlow,
		sendingJobInformation.Gender,
		sendingJobInformation.Nationality,
		sendingJobInformation.FinalEducation,
		sendingJobInformation.SchoolLevel,
		sendingJobInformation.MedicalHistory,
		sendingJobInformation.AgeUnder,
		sendingJobInformation.AgeOver,
		sendingJobInformation.JobChange,
		sendingJobInformation.ShortResignation,
		sendingJobInformation.ShortResignationRemarks,
		sendingJobInformation.SocialExperienceYear,
		sendingJobInformation.SocialExperienceMonth,
		sendingJobInformation.OtherRequired,
		sendingJobInformation.Appearance,
		sendingJobInformation.Communication,
		sendingJobInformation.Thinking,
		sendingJobInformation.TargetDetail,
		sendingJobInformation.RequiredExperienceJobDetail,
		sendingJobInformation.EmploymentInsurance,
		sendingJobInformation.AccidentInsurance,
		sendingJobInformation.HealthInsurance,
		sendingJobInformation.PensionInsurance,
		sendingJobInformation.RegisterPhase,
		sendingJobInformation.StudyCategory,
		sendingJobInformation.WordSkill,
		sendingJobInformation.ExcelSkill,
		sendingJobInformation.PowerPointSkill,
		sendingJobInformation.CorporateSiteURL,
		sendingJobInformation.PostCode,
		sendingJobInformation.OfficeLocation,
		sendingJobInformation.Establishment,
		sendingJobInformation.EmployeeNumberSingle,
		sendingJobInformation.EmployeeNumberGroup,
		sendingJobInformation.PublicOffering,
		sendingJobInformation.EarningsYear,
		sendingJobInformation.Earnings,
		sendingJobInformation.BusinessDetail,
		false, // is_deleted
		time.Now().In(time.UTC),
		time.Now().In(time.UTC),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sendingJobInformation.ID = uint(lastID)
	return nil
}

func (repo *SendingJobInformationRepositoryImpl) Update(sendingJobInformation *entity.SendingJobInformation, sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE sending_job_informations 
		SET
		sending_billing_address_id = ?,
		company_name = ?,
		title = ?,
		recruitment_state = ?,
		expiration_date = ?,
		background = ?,
		work_detail = ?,
		number_of_hires = ?,
		work_location = ?,
		transfer = ?,
		transfer_detail = ?,
		under_income = ?,
		over_income = ?,
		salary = ?,
		insurance = ?,
		work_time = ?,
		overtime = ?,
		overtime_average = ?,
		fixed_overtime = ?,
		fixed_overtime_payment = ?,
		fixed_overtime_detail = ?,
		trial_period = ?,
		trial_period_detail = ?,
		employment_period = ?,
		employment_period_detail = ?,
		holiday_type = ?,
		holiday_detail = ?,
		passive_smoking = ?,
		selection_flow = ?,
		gender = ?,
		nationality = ?,
		final_education = ?,
		school_level = ?,
		medical_history = ?,
		age_under = ?,
		age_over = ?,
		job_change = ?,
		short_resignation = ?,
		short_resignation_remarks = ?,
		social_experience_year = ?,
		social_experience_month = ?,
		other_required = ?,
		appearance = ?,
		communication = ?,
		thinking = ?,
		target_detail = ?,
		required_experience_job_detail=?,
		employment_insurance = ?,
		accident_insurance = ?,
		health_insurance = ?,
		pension_insurance = ?,
		register_phase = ?,
		study_category = ?,
		word_skill = ?,
		excel_skill = ?,
		power_point_skill = ?,
		corporate_site_url = ?,
		post_code = ?,
		office_location = ?,
		establishment = ?,
		employee_number_single = ?,
		employee_number_group = ?,
		public_offering = ?,
		earnings_year = ?,
		earnings = ?,
		business_detail = ?,
		updated_at = ?
		WHERE 
			id = ?
		`,
		sendingJobInformation.SendingBillingAddressID,
		sendingJobInformation.CompanyName,
		sendingJobInformation.Title,
		sendingJobInformation.RecruitmentState,
		sendingJobInformation.ExpirationDate,
		sendingJobInformation.Background,
		sendingJobInformation.WorkDetail,
		sendingJobInformation.NumberOfHires,
		sendingJobInformation.WorkLocation,
		sendingJobInformation.Transfer,
		sendingJobInformation.TransferDetail,
		sendingJobInformation.UnderIncome,
		sendingJobInformation.OverIncome,
		sendingJobInformation.Salary,
		sendingJobInformation.Insurance,
		sendingJobInformation.WorkTime,
		sendingJobInformation.Overtime,
		sendingJobInformation.OvertimeAverage,
		sendingJobInformation.FixedOvertime,
		sendingJobInformation.FixedOvertimePayment,
		sendingJobInformation.FixedOvertimeDetail,
		sendingJobInformation.TrialPeriod,
		sendingJobInformation.TrialPeriodDetail,
		sendingJobInformation.EmploymentPeriod,
		sendingJobInformation.EmploymentPeriodDetail,
		sendingJobInformation.HolidayType,
		sendingJobInformation.HolidayDetail,
		sendingJobInformation.PassiveSmoking,
		sendingJobInformation.SelectionFlow,
		sendingJobInformation.Gender,
		sendingJobInformation.Nationality,
		sendingJobInformation.FinalEducation,
		sendingJobInformation.SchoolLevel,
		sendingJobInformation.MedicalHistory,
		sendingJobInformation.AgeUnder,
		sendingJobInformation.AgeOver,
		sendingJobInformation.JobChange,
		sendingJobInformation.ShortResignation,
		sendingJobInformation.ShortResignationRemarks,
		sendingJobInformation.SocialExperienceYear,
		sendingJobInformation.SocialExperienceMonth,
		sendingJobInformation.OtherRequired,
		sendingJobInformation.Appearance,
		sendingJobInformation.Communication,
		sendingJobInformation.Thinking,
		sendingJobInformation.TargetDetail,
		sendingJobInformation.RequiredExperienceJobDetail,
		sendingJobInformation.EmploymentInsurance,
		sendingJobInformation.AccidentInsurance,
		sendingJobInformation.HealthInsurance,
		sendingJobInformation.PensionInsurance,
		sendingJobInformation.RegisterPhase,
		sendingJobInformation.StudyCategory,
		sendingJobInformation.WordSkill,
		sendingJobInformation.ExcelSkill,
		sendingJobInformation.PowerPointSkill,
		sendingJobInformation.CorporateSiteURL,
		sendingJobInformation.PostCode,
		sendingJobInformation.OfficeLocation,
		sendingJobInformation.Establishment,
		sendingJobInformation.EmployeeNumberSingle,
		sendingJobInformation.EmployeeNumberGroup,
		sendingJobInformation.PublicOffering,
		sendingJobInformation.EarningsYear,
		sendingJobInformation.Earnings,
		sendingJobInformation.BusinessDetail,
		time.Now().In(time.UTC),
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRepositoryImpl) Delete(sendingJobInformationID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		UPDATE 
			sending_job_informations
		SET
			is_deleted = true,
			updated_at = ?
		WHERE
			id = ?
		`,
		time.Now().In(time.UTC),
		sendingJobInformationID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRepositoryImpl) DeleteBySendingBillingAddressID(billingAddressID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteBySendingBillingAddressID",
		`
		UPDATE 
			sending_job_informations
		SET
			is_deleted = true,
			updated_at = ?
		WHERE
			sending_billing_address_id = ?
		`,
		time.Now().In(time.UTC),
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingJobInformationRepositoryImpl) FindByID(sendingJobInformationID uint) (*entity.SendingJobInformation, error) {
	var (
		sendingJobInformation entity.SendingJobInformation
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&sendingJobInformation, `
		SELECT 
			job_info.*, billing.sending_enterprise_id, 
			sending_enterprise.company_name AS sending_enterprise_name,
			staff.staff_name, 
			staff.id AS agent_staff_id
		FROM 
			sending_job_informations AS job_info
		INNER JOIN
			sending_billing_addresses AS billing
		ON
			job_info.sending_billing_address_id = billing.id
		INNER JOIN
			sending_enterprises AS sending_enterprise
		ON
			billing.sending_enterprise_id = sending_enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			job_info.id = ?
		AND
			job_info.is_deleted = false
		LIMIT 1
		`, sendingJobInformationID)
	if err != nil {
		return nil, err
	}

	return &sendingJobInformation, nil
}

func (repo *SendingJobInformationRepositoryImpl) FindByUUID(uuid uuid.UUID) (*entity.SendingJobInformation, error) {
	var (
		sendingJobInformation entity.SendingJobInformation
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&sendingJobInformation, `
			SELECT 
				job_info.*, billing.sending_enterprise_id, 
				sending_enterprise.company_name AS sending_enterprise_name,
				staff.staff_name, 
				staff.id AS agent_staff_id
			FROM 
				sending_job_informations AS job_info
			INNER JOIN
				sending_billing_addresses AS billing
			ON
				job_info.sending_billing_address_id = billing.id
			INNER JOIN
				sending_enterprises AS sending_enterprise
			ON
				billing.sending_enterprise_id = sending_enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE
				job_info.uuid = ?
			AND
				job_info.is_deleted = false
			LIMIT 1
		`,
		uuid,
	)

	if err != nil {
		return nil, err
	}

	return &sendingJobInformation, nil
}

func (repo *SendingJobInformationRepositoryImpl) GetListBySendingBillingAddressID(billingAddressID uint) ([]*entity.SendingJobInformation, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingBillingAddressID",
		&sendingJobInformationList, `
			SELECT 
				job_info.*, 
				billing.sending_enterprise_id, 
				sending_enterprise.company_name AS sending_enterprise_name,
				staff.staff_name
			FROM 
				sending_job_informations AS job_info
			INNER JOIN
				sending_billing_addresses AS billing
			ON
				job_info.sending_billing_address_id = billing.id
			INNER JOIN
				sending_enterprises AS sending_enterprise
			ON
				billing.sending_enterprise_id = sending_enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE 
				job_info.sending_billing_address_id = ?
			AND
				job_info.is_deleted = false
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobInformationList, nil
}

func (repo *SendingJobInformationRepositoryImpl) GetListBySendingEnterpriseID(sendingEnterpriseID uint) ([]*entity.SendingJobInformation, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseID",
		&sendingJobInformationList, `
			SELECT 
				job_info.*, 
				billing.sending_enterprise_id, 
				sending_enterprise.company_name AS sending_enterprise_name,
				staff.staff_name
			FROM 
				sending_job_informations AS job_info
			INNER JOIN 
				sending_billing_addresses AS billing
			ON 
				job_info.sending_billing_address_id = billing.id
			INNER JOIN
				sending_enterprises AS sending_enterprise
			ON
				billing.sending_enterprise_id = sending_enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE 
				job_info.sending_billing_address_id IN (
					SELECT id
					FROM sending_billing_addresses
					WHERE sending_enterprise_id = ?
			)
			AND
				job_info.is_deleted = false
			ORDER BY 
				job_info.updated_at DESC
		`,
		sendingEnterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobInformationList, nil
}

// 求人IDリストに合致する求人の一覧を取得
func (repo *SendingJobInformationRepositoryImpl) GetListByIDList(idList []uint) ([]*entity.SendingJobInformation, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
	)

	if len(idList) == 0 {
		return sendingJobInformationList, nil
	}

	query := fmt.Sprintf(`
		SELECT
			job_info.*,
			billing.sending_enterprise_id, billing.agent_staff_id,
			sending_enterprise.company_name AS sending_enterprise_name,
			staff.staff_name
		FROM
			sending_job_informations AS job_info
		INNER JOIN
			sending_billing_addresses AS billing
		ON
			job_info.sending_billing_address_id = billing.id
		INNER JOIN
			sending_enterprises AS sending_enterprise
		ON
			billing.sending_enterprise_id = sending_enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			job_info.id IN (%s)
		AND
			job_info.is_deleted = false
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByIDList",
		&sendingJobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobInformationList, nil
}

// 求人IDリストに合致する求人の一覧を取得
func (repo *SendingJobInformationRepositoryImpl) GetListBySendingEnterpriseIDList(enterpriseIDList []uint) ([]*entity.SendingJobInformation, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
	)

	if len(enterpriseIDList) == 0 {
		return sendingJobInformationList, nil
	}

	query := fmt.Sprintf(`
		SELECT
			job_info.*,
			billing.sending_enterprise_id, billing.agent_staff_id,
			sending_enterprise.company_name AS sending_enterprise_name,
			staff.staff_name
		FROM
			sending_job_informations AS job_info
		INNER JOIN
			sending_billing_addresses AS billing
		ON
			job_info.sending_billing_address_id = billing.id
		INNER JOIN
			sending_enterprises AS sending_enterprise
		ON
			billing.sending_enterprise_id = sending_enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			billing.sending_enterprise_id IN (%s)
		AND
			job_info.is_deleted = false
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enterpriseIDList)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseIDList",
		&sendingJobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobInformationList, nil
}

// 送客先を探すページの絞り込み検索に使用
func (repo *SendingJobInformationRepositoryImpl) GetListBySendingEnterpriseIDListAndFreeWord(enterpriseIDList []uint, freeWord string) ([]*entity.SendingJobInformation, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
		freeWordQuery             string
	)

	if len(enterpriseIDList) == 0 {
		return sendingJobInformationList, nil
	}

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {

			freeWordQuery = fmt.Sprintf(`
				AND (
					MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE) OR 
					MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE)
				)
			`, freeWord, freeWord)
		}
	} else {
		freeWordQuery = fmt.Sprintf(`
			AND job_info.id = %d
		`, freeWordInt)
	}

	query := fmt.Sprintf(`
		SELECT
			job_info.*,
			billing.sending_enterprise_id, billing.agent_staff_id,
			enterprise.company_name AS sending_enterprise_name,
			staff.staff_name
		FROM
			sending_job_informations AS job_info
		INNER JOIN
			sending_billing_addresses AS billing
		ON
			job_info.sending_billing_address_id = billing.id
		INNER JOIN
			sending_enterprises AS enterprise
		ON
			billing.sending_enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE
			billing.sending_enterprise_id IN (%s)
		AND
			job_info.is_deleted = false
		%s
	`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enterpriseIDList)), ", "), "[]"),
		freeWordQuery)

	err = repo.executer.Select(
		repo.Name+".GetListBySendingEnterpriseIDListAndFreeWord",
		&sendingJobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobInformationList, nil
}

// 有効な求人を取得
func (repo *SendingJobInformationRepositoryImpl) GetActiveList() ([]*entity.SendingJobInformation, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetActiveList",
		&sendingJobInformationList, `
			SELECT
				job_info.id,
				job_info.title,
				billing.sending_enterprise_id,
				sending_enterprise.company_name AS sending_enterprise_name
			FROM
				sending_job_informations AS job_info
			INNER JOIN
				sending_billing_addresses AS billing
			ON
				job_info.sending_billing_address_id = billing.id
			INNER JOIN
				sending_enterprises AS sending_enterprise
			ON
				billing.sending_enterprise_id = sending_enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE
				job_info.recruitment_state = 0
			AND
				job_info.register_phase = 0
			AND
				job_info.is_deleted = false
			ORDER BY job_info.id DESC
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobInformationList, nil
}

/****************************************************************************************/
func (repo *SendingJobInformationRepositoryImpl) GetAll() ([]*entity.SendingJobInformation, error) {
	var (
		sendingJobInformationList []*entity.SendingJobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetAll",
		&sendingJobInformationList, `
			SELECT *
			FROM sending_job_informations
			WHERE is_deleted = false
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingJobInformationList, nil
}

// count
func (repo *SendingJobInformationRepositoryImpl) CountByAgentID(agentID uint) (uint, error) {
	var (
		count uint
	)

	err := repo.executer.Select(
		repo.Name+".CountByAgentID",
		&count, `
			SELECT 
				COUNT(*)
			FROM
				sending_job_informations AS job_info
			INNER JOIN 
				sending_billing_addresses AS billing
			ON 
				job_info.sending_billing_address_id = billing.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			WHERE 
				staff.agent_id = ?
			AND
				job_info.is_deleted = false
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return count, nil
}
