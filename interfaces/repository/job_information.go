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
	"gopkg.in/guregu/null.v4"
)

type JobInformationRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewJobInformationRepositoryImpl(ex interfaces.SQLExecuter) usecase.JobInformationRepository {
	return &JobInformationRepositoryImpl{
		Name:     "JobInformationRepository",
		executer: ex,
	}
}

/****************************************************************************************/
/// 作成 API
//
// 求人の作成
func (repo *JobInformationRepositoryImpl) Create(jobInformation *entity.JobInformation) error {
	now := time.Now().In(time.UTC)

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO job_informations (
			uuid,
			billing_address_id,
			title,
			recruitment_state,
			expiration_date,
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
			overtime_average,
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
			appearance,
			communication,

			thinking,
			target_detail,
			commission,
			commission_rate,
			commission_detail,
			refund_policy,
			required_experience_job_detail,
			secret_memo,
			required_documents_detail,
			employment_insurance,

			accident_insurance,
			health_insurance,
			pension_insurance,
			register_phase,
			study_category,
			driver_licence,
			word_skill,
			excel_skill,
			power_point_skill,
			is_deleted,

			is_external,
			work_detail_after_hiring,
			work_detail_scope_of_change,
			offer_rate,
			document_passing_rate,
			number_of_recent_applications,
			is_guaranteed_interview,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?
		)
		`,
		utility.CreateUUID(),
		jobInformation.BillingAddressID,
		jobInformation.Title,
		jobInformation.RecruitmentState,
		jobInformation.ExpirationDate,
		jobInformation.WorkDetail,
		jobInformation.NumberOfHires,
		jobInformation.WorkLocation,
		jobInformation.Transfer,
		jobInformation.TransferDetail,
		jobInformation.UnderIncome,
		jobInformation.OverIncome,
		jobInformation.Salary,
		jobInformation.Insurance,
		jobInformation.WorkTime,
		jobInformation.OvertimeAverage,
		jobInformation.FixedOvertimePayment,
		jobInformation.FixedOvertimeDetail,
		jobInformation.TrialPeriod,
		jobInformation.TrialPeriodDetail,
		jobInformation.EmploymentPeriod,
		jobInformation.EmploymentPeriodDetail,
		jobInformation.HolidayType,
		jobInformation.HolidayDetail,
		jobInformation.PassiveSmoking,
		jobInformation.SelectionFlow,
		jobInformation.Gender,
		jobInformation.Nationality,
		jobInformation.FinalEducation,
		jobInformation.SchoolLevel,
		jobInformation.MedicalHistory,
		jobInformation.AgeUnder,
		jobInformation.AgeOver,
		jobInformation.JobChange,
		jobInformation.ShortResignation,
		jobInformation.ShortResignationRemarks,
		jobInformation.SocialExperienceYear,
		jobInformation.SocialExperienceMonth,
		jobInformation.Appearance,
		jobInformation.Communication,
		jobInformation.Thinking,
		jobInformation.TargetDetail,
		jobInformation.Commission,
		jobInformation.CommissionRate,
		jobInformation.CommissionDetail,
		jobInformation.RefundPolicy,
		jobInformation.RequiredExperienceJobDetail,
		jobInformation.SecretMemo,
		jobInformation.RequiredDocumentsDetail,
		jobInformation.EmploymentInsurance,
		jobInformation.AccidentInsurance,
		jobInformation.HealthInsurance,
		jobInformation.PensionInsurance,
		jobInformation.RegisterPhase,
		jobInformation.StudyCategory,
		jobInformation.DriverLicence,
		jobInformation.WordSkill,
		jobInformation.ExcelSkill,
		jobInformation.PowerPointSkill,
		false,
		jobInformation.IsExternal,
		jobInformation.WorkDetailAfterHiring,
		jobInformation.WorkDetailScopeOfChange,
		jobInformation.OfferRate,
		jobInformation.DocumentPassingRate,
		jobInformation.NumberOfRecentApplications,
		jobInformation.IsGuaranteedInterview,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	jobInformation.ID = uint(lastID)
	return nil
}

/****************************************************************************************/
/// 更新 API
//
// 求人を更新
func (repo *JobInformationRepositoryImpl) Update(id uint, jobInformation *entity.JobInformation) error {
	_, err := repo.executer.Exec(
		repo.Name+".Update",
		`
		UPDATE job_informations 
		SET
		billing_address_id = ?,
		title = ?,
		recruitment_state = ?,
		expiration_date = ?,
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
		overtime_average = ?,
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
		appearance = ?,
		communication = ?,
		thinking = ?,
		target_detail = ?,
		commission = ?,
		commission_rate = ?,
		commission_detail = ?,
 		refund_policy = ?,
		required_experience_job_detail=?,
		secret_memo = ?,
		required_documents_detail = ?,
		employment_insurance = ?,
		accident_insurance = ?,
		health_insurance = ?,
		pension_insurance = ?,
		register_phase = ?,
		study_category = ?,
		driver_licence = ?,
		word_skill = ?,
		excel_skill = ?,
		power_point_skill = ?,
		is_external = ?,
		work_detail_after_hiring = ?,
		work_detail_scope_of_change = ?,
		offer_rate = ?,
		document_passing_rate = ?,
		number_of_recent_applications = ?,
		is_guaranteed_interview = ?,
		updated_at = ?
		WHERE 
			id = ?
		`,
		jobInformation.BillingAddressID,
		jobInformation.Title,
		jobInformation.RecruitmentState,
		jobInformation.ExpirationDate,
		jobInformation.WorkDetail,
		jobInformation.NumberOfHires,
		jobInformation.WorkLocation,
		jobInformation.Transfer,
		jobInformation.TransferDetail,
		jobInformation.UnderIncome,
		jobInformation.OverIncome,
		jobInformation.Salary,
		jobInformation.Insurance,
		jobInformation.WorkTime,
		jobInformation.OvertimeAverage,
		jobInformation.FixedOvertimePayment,
		jobInformation.FixedOvertimeDetail,
		jobInformation.TrialPeriod,
		jobInformation.TrialPeriodDetail,
		jobInformation.EmploymentPeriod,
		jobInformation.EmploymentPeriodDetail,
		jobInformation.HolidayType,
		jobInformation.HolidayDetail,
		jobInformation.PassiveSmoking,
		jobInformation.SelectionFlow,
		jobInformation.Gender,
		jobInformation.Nationality,
		jobInformation.FinalEducation,
		jobInformation.SchoolLevel,
		jobInformation.MedicalHistory,
		jobInformation.AgeUnder,
		jobInformation.AgeOver,
		jobInformation.JobChange,
		jobInformation.ShortResignation,
		jobInformation.ShortResignationRemarks,
		jobInformation.SocialExperienceYear,
		jobInformation.SocialExperienceMonth,
		jobInformation.Appearance,
		jobInformation.Communication,
		jobInformation.Thinking,
		jobInformation.TargetDetail,
		jobInformation.Commission,
		jobInformation.CommissionRate,
		jobInformation.CommissionDetail,
		jobInformation.RefundPolicy,
		jobInformation.RequiredExperienceJobDetail,
		jobInformation.SecretMemo,
		jobInformation.RequiredDocumentsDetail,
		jobInformation.EmploymentInsurance,
		jobInformation.AccidentInsurance,
		jobInformation.HealthInsurance,
		jobInformation.PensionInsurance,
		jobInformation.RegisterPhase,
		jobInformation.StudyCategory,
		jobInformation.DriverLicence,
		jobInformation.WordSkill,
		jobInformation.ExcelSkill,
		jobInformation.PowerPointSkill,
		jobInformation.IsExternal,
		jobInformation.WorkDetailAfterHiring,
		jobInformation.WorkDetailScopeOfChange,
		jobInformation.OfferRate,
		jobInformation.DocumentPassingRate,
		jobInformation.NumberOfRecentApplications,
		jobInformation.IsGuaranteedInterview,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 求人の募集状況のみを更新
func (repo *JobInformationRepositoryImpl) UpdateRecruitmentState(id uint, recruitmentState null.Int) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateRecruitmentState",
		`
		UPDATE 
			job_informations 
		SET 
			recruitment_state = ?, 
			updated_at = ?
		WHERE 
			id = ?
		`,
		recruitmentState,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 求人の面接確約フラグを更新する
func (repo *JobInformationRepositoryImpl) UpdateIsGuaranteedInterview(id uint, isGuaranteedInterview bool) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateIsGuaranteedInterview",
		`
		UPDATE 
			job_informations 
		SET 
			is_guaranteed_interview = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		isGuaranteedInterview,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err

}

/****************************************************************************************/
/// 削除 API
//
// 求人を削除
func (repo *JobInformationRepositoryImpl) Delete(id uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".Delete",
		`
		UPDATE 
			job_informations
		SET
			is_deleted = true,
			updated_at = ?
		WHERE
			id = ?
		`,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 請求先IDから求人を削除
func (repo *JobInformationRepositoryImpl) DeleteByBillingAddressID(billingAddressID uint) error {
	_, err := repo.executer.Exec(
		repo.Name+".DeleteByBillingAddressID",
		`
		UPDATE 
			job_informations
		SET
			is_deleted = true,
			updated_at = ?
		WHERE
			billing_address_id = ?
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

/****************************************************************************************/
/// 単数取得 API
//
// 指定IDの求人を取得
func (repo *JobInformationRepositoryImpl) FindByID(id uint) (*entity.JobInformation, error) {
	var (
		jobInformation entity.JobInformation
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&jobInformation, `
		SELECT 
			job_info.*, 
			billing.enterprise_id, 
			enterprise.company_name, 
			enterprise.post_code,
			enterprise.office_location,
			enterprise.employee_number_single, 
			enterprise.employee_number_group, 
			enterprise.corporate_site_url,
			enterprise.establishment, 
			enterprise.public_offering,
			enterprise.earnings,
			enterprise.earnings_year,
			enterprise.business_detail,
			staff.staff_name, 
			staff.id AS agent_staff_id, 
			staff.agent_id,
			agent.agent_name
		FROM 
			job_informations AS job_info
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
		 agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			job_info.id = ?
		AND
			job_info.is_deleted = false
		LIMIT 1
		`, id)
	if err != nil {
		return nil, err
	}

	return &jobInformation, nil
}

// 指定UUIDの求人を取得
func (repo *JobInformationRepositoryImpl) FindByUUID(uuid uuid.UUID) (*entity.JobInformation, error) {
	var (
		jobInformation entity.JobInformation
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&jobInformation, `
			SELECT 
				job_info.*, 
				billing.enterprise_id, 
				enterprise.company_name, 
				enterprise.post_code,
				enterprise.office_location,
				enterprise.employee_number_single, 
				enterprise.employee_number_group, 
				enterprise.corporate_site_url,
				enterprise.establishment, 
				enterprise.public_offering,
				enterprise.earnings,
				enterprise.earnings_year,
				enterprise.business_detail,
				staff.staff_name, 
				staff.id AS agent_staff_id, 
				staff.agent_id,
				agent.agent_name
			FROM 
				job_informations AS job_info
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			INNER JOIN
				agents AS agent
			ON
				staff.agent_id = agent.id
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

	return &jobInformation, nil
}

// 指定タスクグループUUIDの求人を取得
func (repo *JobInformationRepositoryImpl) FindByTaskGroupUUID(taskGroupUUID uuid.UUID) (*entity.JobInformation, error) {
	var (
		jobInformation entity.JobInformation
	)

	err := repo.executer.Get(
		repo.Name+".FindByUUID",
		&jobInformation, `
			SELECT 
				job_info.*
			FROM 
				job_informations AS job_info
			INNER JOIN
				task_groups AS task_group
			ON
				job_info.id = task_group.job_information_id
			WHERE
				task_group.uuid = ?
			AND
				job_info.is_deleted = false
			LIMIT 1
		`,
		taskGroupUUID,
	)

	if err != nil {
		return nil, err
	}

	return &jobInformation, nil
}

/****************************************************************************************/
/// 複数取得 API
//
// 請求先IDから求人一覧を取得
func (repo *JobInformationRepositoryImpl) GetByBillingAddressID(billingAddressID uint) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByBillingAddressID",
		&jobInformationList, `
			SELECT 
				job_info.*, 
				billing.enterprise_id, 
				enterprise.company_name, 
				enterprise.employee_number_single, 
				enterprise.employee_number_group, 
				staff.staff_name, 
				staff.agent_id, 
				agent.agent_name
			FROM 
				job_informations AS job_info
			INNER JOIN
				billing_addresses AS billing
			ON
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			INNER JOIN
				agents AS agent
			ON
				staff.agent_id = agent.id
			WHERE 
				job_info.billing_address_id = ?
			AND
				job_info.is_deleted = false
		`,
		billingAddressID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// 求人企業IDから求人一覧を取得
func (repo *JobInformationRepositoryImpl) GetByEnterpriseID(enterpriseID uint) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByEnterpriseID",
		&jobInformationList, `
			SELECT 
				job_info.*, 
				billing.enterprise_id, 
				enterprise.company_name, 
				enterprise.employee_number_single, 
				enterprise.employee_number_group, 
				staff.staff_name, 
				staff.agent_id, 
				agent.agent_name
			FROM 
				job_informations AS job_info
			INNER JOIN 
				billing_addresses AS billing
			ON 
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			INNER JOIN
				agents AS agent
			ON
				staff.agent_id = agent.id
			WHERE 
				job_info.billing_address_id IN (
					SELECT id
					FROM billing_addresses
					WHERE enterprise_id = ?
			)
			AND
				job_info.is_deleted = false
		`,
		enterpriseID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// 求人IDリストに合致する求人の一覧を取得
func (repo *JobInformationRepositoryImpl) GetByIDList(idList []uint) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	if len(idList) == 0 {
		return jobInformationList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_info.*, 
			billing.enterprise_id, billing.agent_staff_id, 
			enterprise.company_name, 
			enterprise.employee_number_single, 
			enterprise.employee_number_group, 
			staff.staff_name, staff.agent_id,
			agent.agent_name
		FROM
			job_informations AS job_info
		INNER JOIN 
			billing_addresses AS billing
		ON 
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			job_info.id IN (%s)
		AND
			job_info.is_deleted = false
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByIDList",
		&jobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// 求人IDリストに合致する求人の一覧を取得
func (repo *JobInformationRepositoryImpl) GetByIDListForTask(idList []uint) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	if len(idList) == 0 {
		return jobInformationList, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			job_info.id,
			job_info.uuid,
			job_info.title,
			job_info.billing_address_id,
			job_info.required_documents_detail,
			billing.how_to_recommend, 
			billing.enterprise_id, 
			enterprise.company_name, 
			enterprise.agent_staff_id, 
			staff.staff_name, 
			staff.agent_id,
			agent.agent_name
		FROM
			job_informations AS job_info
		INNER JOIN 
			billing_addresses AS billing
		ON 
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			job_info.id IN (%s)
		AND
			job_info.is_deleted = false
	`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idList)), ", "), "[]"))

	err := repo.executer.Select(
		repo.Name+".GetByIDListForTask",
		&jobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDから求人一覧を取得
func (repo *JobInformationRepositoryImpl) GetByAgentID(agentID uint) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentID",
		&jobInformationList, `
			SELECT 
				job_info.*, 
				billing.enterprise_id, 
				enterprise.company_name, 
				enterprise.employee_number_single, 
				enterprise.employee_number_group, 
				staff.staff_name, staff.agent_id,
				agent.agent_name
			FROM
				job_informations AS job_info
			INNER JOIN 
				billing_addresses AS billing
			ON 
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			INNER JOIN
				agents AS agent
			ON
				staff.agent_id = agent.id
			WHERE 
				staff.agent_id = ?
			AND
				job_info.is_deleted = false
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDとアライアンスから求人一覧を取得
func (repo *JobInformationRepositoryImpl) GetByAgentIDAndAlliance(agentID uint, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobInformation, error) {
	var (
		jobInformationList  []*entity.JobInformation
		allianceAgentIDList []uint
		allainceQuery       string
	)

	// アライアンスエージェントのクエリを作成
	if len(agentAllianceList) > 0 {
		// アライアンス締結済みの他社エージェントのIDを取得
		for _, alliance := range agentAllianceList {
			if alliance.Agent1ID == agentID {
				allianceAgentIDList = append(allianceAgentIDList, alliance.Agent2ID)
			} else if alliance.Agent2ID == agentID {
				allianceAgentIDList = append(allianceAgentIDList, alliance.Agent1ID)
			}
		}

		// アライアンスエージェントのIDをカンマ区切りの文字列形式に変換
		idSeparatedByCommas := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allianceAgentIDList)), ","), "[]")

		allainceQuery = fmt.Sprintf(
			`OR (staff.agent_id IN (%s) AND billing.contract_phase = 2)`,
			idSeparatedByCommas,
		)
	}

	query := fmt.Sprintf(
		`SELECT
			job_info.*,
			billing.enterprise_id,
			enterprise.company_name,
			staff.staff_name,
			staff.agent_id,
			agent.agent_name
		FROM
			job_informations AS job_info
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			(
				staff.agent_id = %v 
				%s
			)
		AND
			job_info.recruitment_state = 0
		AND
			job_info.register_phase = 0
		AND
			job_info.is_deleted = false
		ORDER BY job_info.id DESC
		`,
		agentID,
		allainceQuery,
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDAndAlliance",
		&jobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDと検索ワードから求人一覧を取得
func (repo *JobInformationRepositoryImpl) GetByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
		freeWordQuery      string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {

			freeWordQuery = fmt.Sprintf(`
		AND (
			MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE)
		OR 
			MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE)
		)`, freeWord, freeWord)
		}
	} else {
		freeWordQuery = fmt.Sprintf(`
			AND 	
				job_info.id = %d
				`, freeWordInt)
	}

	// クエリをまとめる
	query := fmt.Sprintf(`
		SELECT 
			job_info.*, 
			billing.enterprise_id, 
			billing.agent_staff_id, 
			enterprise.company_name, 
			enterprise.employee_number_single, 
			enterprise.employee_number_group, 
			staff.staff_name, 
			staff.agent_id,
			agent.agent_name
		FROM
			job_informations AS job_info
		INNER JOIN 
			billing_addresses AS billing
		ON 
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE 
			staff.agent_id = %v
		%s
			AND job_info.is_deleted = false
		ORDER BY job_info.id DESC`,
		agentID, freeWordQuery)

	// クエリを実行
	err = repo.executer.Select(
		repo.Name+".GetByAgentIDAndFreeWord",
		&jobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDと検索ワード、アライアンスから求人一覧を取得
func (repo *JobInformationRepositoryImpl) GetByAgentIDAndFreeWordAndAlliance(agentID uint, freeWord string, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobInformation, error) {
	var (
		jobInformationList  []*entity.JobInformation
		allianceAgentIDList []uint
		allainceQuery       string
		freeWordQuery       string
	)

	// アライアンスエージェントのクエリを作成
	if len(agentAllianceList) > 0 {
		// アライアンス締結済みの他社エージェントのIDを取得
		for _, alliance := range agentAllianceList {
			if alliance.Agent1ID == agentID {
				allianceAgentIDList = append(allianceAgentIDList, alliance.Agent2ID)
			} else if alliance.Agent2ID == agentID {
				allianceAgentIDList = append(allianceAgentIDList, alliance.Agent1ID)
			}
		}

		// アライアンスエージェントのIDをカンマ区切りの文字列形式に変換
		idSeparatedByCommas := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allianceAgentIDList)), ","), "[]")

		allainceQuery = fmt.Sprintf(
			`OR (staff.agent_id IN (%s) AND billing.contract_phase = 2)`,
			idSeparatedByCommas,
		)
	}

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {

			freeWordQuery = fmt.Sprintf(`
		AND (
			MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE)
		OR 
			MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE)
		)`, freeWord, freeWord)
		}
	} else {
		freeWordQuery = fmt.Sprintf(`
			AND 	
				job_info.id = %d
				`, freeWordInt)
	}

	query := fmt.Sprintf(
		`SELECT
			job_info.*,
			billing.enterprise_id, 
			billing.agent_staff_id,
			enterprise.company_name,
			staff.staff_name,
			staff.agent_id,
			agent.agent_name
		FROM
			job_informations AS job_info
		INNER JOIN
			billing_addresses AS billing
		ON
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
		( 
			(staff.agent_id = %v %s)
			AND job_info.recruitment_state = 0
			AND job_info.register_phase = 0
		)
			%s
		AND
			job_info.is_deleted = false
		ORDER BY job_info.id DESC
		`,
		agentID,
		allainceQuery,
		freeWordQuery,
	)

	err = repo.executer.Select(
		repo.Name+".GetByAgentIDAndFreeWordAndAlliance",
		&jobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDとフリーワードからアクティブな求人 -自社求人検索など
func (repo *JobInformationRepositoryImpl) GetActiveByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
		freeWordQuery      string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordQuery = fmt.Sprintf(`
		AND (
			MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE)
		OR 
			MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE)
		)`, freeWord, freeWord)
		}
	} else {
		freeWordQuery = fmt.Sprintf(`
		AND 	
			job_info.id = %d
			`, freeWordInt)
	}

	query := fmt.Sprintf(`
	SELECT 
		job_info.*,
		billing.enterprise_id,
		billing.agent_staff_id,
		enterprise.company_name,
		enterprise.post_code,
		enterprise.employee_number_single, 
		enterprise.employee_number_group, 
		staff.staff_name,
		staff.agent_id,
		agent.agent_name
	FROM
		job_informations AS job_info
	INNER JOIN 
		billing_addresses AS billing
	ON 
		job_info.billing_address_id = billing.id
	INNER JOIN
		enterprise_profiles AS enterprise
	ON
		billing.enterprise_id = enterprise.id
	INNER JOIN
		agent_staffs AS staff
	ON
		billing.agent_staff_id = staff.id
	INNER JOIN
		agents AS agent
	ON
		staff.agent_id = agent.id
	WHERE 
		staff.agent_id = %d
		%s
	AND
		job_info.recruitment_state = 0
	AND
		job_info.register_phase = 0
	AND
		job_info.is_deleted = false
	`, agentID, freeWordQuery)

	err = repo.executer.Select(
		repo.Name+".GetActiveByAgentIDAndFreeWord",
		&jobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDとフリーワードから自社全て+他社の請求フェーズが契約済みのアクティブ求人
func (repo *JobInformationRepositoryImpl) GetActiveAllByAgentIDAndFreeWord(agentID uint, freeWord string) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
		freeWordQuery      string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordQuery = fmt.Sprintf(`
		AND (
			MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE)
		OR 
			MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE)
		)
		`, freeWord, freeWord)
		}
	} else {
		freeWordQuery = fmt.Sprintf(`
		AND
			job_info.id = %v
		`, freeWordInt)
	}

	query := fmt.Sprintf(`
		SELECT 
			job_info.*,
			billing.enterprise_id,
			billing.agent_staff_id,
			enterprise.company_name,
			enterprise.post_code,
			enterprise.employee_number_single, 
			enterprise.employee_number_group, 
			staff.staff_name,
			staff.agent_id,
			agent.agent_name
		FROM
			job_informations AS job_info
		INNER JOIN 
			billing_addresses AS billing
		ON 
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE (
			staff.agent_id = %v
			OR (
				staff.agent_id NOT IN (%v) 
				AND
				billing.contract_phase = 2
			)
		)
		%s
		AND
			job_info.recruitment_state = 0
		AND
			job_info.register_phase = 0
		AND
			job_info.is_deleted = false
		ORDER BY job_info.id DESC
			`, agentID, agentID, freeWordQuery)

	err = repo.executer.Select(
		repo.Name+".GetActiveAllByAgentIDAndFreeWord",
		&jobInformationList, query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDから自社全て+他社の請求フェーズが契約済みのアクティブ求人
func (repo *JobInformationRepositoryImpl) GetActiveAllByAgentID(agentID uint) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetActiveAllByAgentID",
		&jobInformationList, `
			SELECT 
				job_info.*,
				billing.enterprise_id,
				billing.agent_staff_id,
				enterprise.company_name,
				enterprise.post_code,
				enterprise.employee_number_single, 
				enterprise.employee_number_group, 
				staff.staff_name,
				staff.agent_id,
				agent.agent_name
			FROM
				job_informations AS job_info
			INNER JOIN 
				billing_addresses AS billing
			ON 
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			INNER JOIN
				agents AS agent
			ON
				staff.agent_id = agent.id
			WHERE (
				staff.agent_id = ? 
				OR (
					staff.agent_id NOT IN (?) 
					AND
					billing.contract_phase = 2
				)
			)
			AND
				job_info.recruitment_state = 0
			AND
				job_info.register_phase = 0
			AND
				job_info.is_deleted = false
			ORDER BY job_info.id DESC
		`,
		agentID, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDから自社の請求フェーズが契約済みのアクティブ求人（外部求人を除く）
func (repo *JobInformationRepositoryImpl) GetActiveAllByAgentIDWithoutExternal(agentID uint) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetActiveAllByAgentIDWithoutExternal",
		&jobInformationList, `
			SELECT
				job_info.*,
				billing.enterprise_id,
				billing.agent_staff_id,
				enterprise.company_name,
				enterprise.post_code,
				enterprise.employee_number_single, 
				enterprise.employee_number_group, 
				staff.staff_name,
				staff.agent_id,
				agent.agent_name
			FROM
				job_informations AS job_info
			INNER JOIN 
				billing_addresses AS billing
			ON 
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				agent_staffs AS staff
			ON
				billing.agent_staff_id = staff.id
			INNER JOIN
				agents AS agent
			ON
				staff.agent_id = agent.id
			WHERE 
				staff.agent_id = ? 
			AND
				job_info.recruitment_state = 0
			AND
				job_info.register_phase = 0
			AND
				job_info.is_deleted = FALSE
			AND
			  job_info.is_external = FALSE
			ORDER BY 
				job_info.is_guaranteed_interview DESC,
				job_info.id DESC
		`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDから自社の請求フェーズが契約済みのアクティブ求人（外部求人を除く）
func (repo *JobInformationRepositoryImpl) GetActiveAllByAgentIDAndDiagnosisParamWithoutExternal(agentID uint, diagnosisParam entity.DiagnosisParam) ([]*entity.JobInformation, error) {
	var (
		jobInformationList  []*entity.JobInformation
		birthyear                    = diagnosisParam.Birthyear
		birthmonth                   = diagnosisParam.Birthmonth
		birthday                     = diagnosisParam.Birthday
		gender                       = diagnosisParam.Gender
		companyNum                   = diagnosisParam.CompanyNum
		schoolCategory               = diagnosisParam.SchoolCategory
		seekerAge           null.Int = null.NewInt(0, false)
		genderQuery         string
		ageQuery            string
		jobChangeQuery      string
		finalEducationQuery string
	)

	// 性別の条件
	if gender.Valid {
		/** 求職者が選択した性別（0: 男性 OR 1: 女性） */
		genderQuery = fmt.Sprintf(`
			AND (
				job_info.gender IS NULL OR
				job_info.gender IN(2, 3, 99, %v)
			)
		`, gender.Int64)
	} else {
		genderQuery = `
			AND (
				job_info.gender IS NULL OR
				job_info.gender IN(2, 3, 99)
			)
		`
	}

	// 年齢の変換
	if birthyear != "" && birthmonth != "" && birthday != "" {
		// 現在の日付を取得
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().In(jst)

		// フォーマットに従って time.Time 型に変換
		dateStr := fmt.Sprintf("%s-%s-%s", birthyear, birthmonth, birthday)
		bd, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			// 年齢を計算
			age := now.Year() - bd.Year()
			seekerAge = null.NewInt(int64(age), true)

			// まだ誕生日を迎えていない場合は、年齢から1を引く
			if now.Month() < bd.Month() || (now.Month() == bd.Month() && now.Day() < bd.Day()) {
				diffAge := age - 1
				seekerAge = null.NewInt(int64(diffAge), true)
			}
		}
	}

	// 年齢の条件
	if seekerAge.Valid {
		/** 検索パラムの値が下限以上で上限以下の場合 */
		ageQuery = fmt.Sprintf(`
			AND (
				(job_info.age_under IS NULL AND job_info.age_over IS NULL) OR
				(job_info.age_under <= %v OR %v <= job_info.age_over)
			)
		`, seekerAge.Int64, seekerAge.Int64)
	} else {
		ageQuery = `
			AND
				job_info.age_under IS NULL AND 
				job_info.age_over IS NULL
		`
	}

	// 経験社数の条件
	if companyNum.Valid {
		/** 求人の指定転職回数以下 */
		/** NOTE: 不問を条件外にしているのは{ 6: 7社以上 }の条件を除外するため */
		jobChangeQuery = fmt.Sprintf(`
			AND (
				job_info.job_change IS NULL OR
				job_info.job_change = 99 OR
				(
					job_info.job_change != 99 AND
					%v <= job_info.job_change
				)
			)
		`, companyNum.Int64)
	} else {
		jobChangeQuery = `
			AND (
				job_info.job_change IS NULL OR
				job_info.job_change = 99
			)
		`
	}

	// 最終学歴の条件
	if schoolCategory.Valid {
		category := schoolCategory.Int64
		/** 高卒以上 OR 専卒以上 OR 短大卒以上 OR 専卒と短大卒以上 OR 高専卒以上 OR 大卒以上 OR 院卒以上 */
		finalEducationQuery = fmt.Sprintf(`
			AND (
				job_info.final_education IS NULL OR
				job_info.final_education IN(0, 99) OR
				(job_info.final_education = 1 AND 1 <= %v) OR
				(job_info.final_education = 2 AND %v IN (2, 4, 5, 6)) OR
				(job_info.final_education = 3 AND %v IN (3, 4, 5, 6)) OR
				(job_info.final_education = 4 AND %v IN (2, 3, 4, 5, 6)) OR
				(job_info.final_education = 5 AND %v IN (4, 5, 6)) OR
				(job_info.final_education = 6 AND %v IN (5, 6)) OR
				(job_info.final_education = 7 AND %v = 6)
			)
		`, category, category, category, category, category, category, category)
	} else {
		finalEducationQuery = `
			AND (
				job_info.final_education IS NULL OR
				job_info.final_education IN(0, 99)
			)
		`
	}

	mainQuery := fmt.Sprintf(`
		SELECT
			job_info.*,
			billing.enterprise_id,
			billing.agent_staff_id,
			staff.staff_name,
			staff.agent_id
		FROM
			job_informations AS job_info
		INNER JOIN 
			billing_addresses AS billing
		ON 
			job_info.billing_address_id = billing.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		WHERE 
			staff.agent_id = ? 
		AND
			job_info.recruitment_state = 0
		AND
			job_info.register_phase = 0
		AND
			job_info.is_deleted = FALSE
		AND
			job_info.is_external = FALSE
		%s
		%s
		%s
		%s
		ORDER BY 
			job_info.is_guaranteed_interview DESC,
			job_info.id DESC
	`, genderQuery, ageQuery, jobChangeQuery, finalEducationQuery)

	err := repo.executer.Select(
		repo.Name+".GetActiveAllByAgentIDAndDiagnosisParamWithoutExternal",
		&jobInformationList, mainQuery, agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// エージェントIDとアライアンスから自社エージェント除く有効なアライアンス求人を取得
func (repo *JobInformationRepositoryImpl) GetActiveAllianceByAgentIDAndFreeWordAndAlliance(agentID uint, freeWord string, agentAllianceList []*entity.AgentAlliance) ([]*entity.JobInformation, error) {
	var (
		jobInformationList  []*entity.JobInformation
		freeWordQuery       string
		allianceAgentIDList []uint
	)

	if len(agentAllianceList) < 1 {
		// アライアンスエージェントがいない場合は、空を返す
		return nil, nil
	}

	// アライアンスエージェントのクエリを作成
	// アライアンス締結済みの他社エージェントのIDを取得
	for _, alliance := range agentAllianceList {
		if alliance.Agent1ID == agentID {
			allianceAgentIDList = append(allianceAgentIDList, alliance.Agent2ID)
		} else if alliance.Agent2ID == agentID {
			allianceAgentIDList = append(allianceAgentIDList, alliance.Agent1ID)
		}
	}

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {

			freeWordQuery = fmt.Sprintf(`
		AND (
			MATCH(enterprise.company_name) AGAINST('%s' IN BOOLEAN MODE)
		OR 
			MATCH(job_info.title) AGAINST('%s' IN BOOLEAN MODE)
		)`, freeWord, freeWord)
		}
	} else {
		freeWordQuery = fmt.Sprintf(`
			AND 	
				job_info.id = %d
				`, freeWordInt)
	}

	// アライアンスエージェントのIDをカンマ区切りの文字列形式に変換
	idSeparatedByCommas := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allianceAgentIDList)), ","), "[]")

	// クエリをまとめる
	query := fmt.Sprintf(`
		SELECT 
			job_info.*, 
			billing.enterprise_id, 
			billing.agent_staff_id, 
			enterprise.company_name, 
			enterprise.employee_number_single, 
			enterprise.employee_number_group, 
			staff.staff_name, 
			staff.agent_id,
			agent.agent_name
		FROM
			job_informations AS job_info
		INNER JOIN 
			billing_addresses AS billing
		ON 
			job_info.billing_address_id = billing.id
		INNER JOIN
			enterprise_profiles AS enterprise
		ON
			billing.enterprise_id = enterprise.id
		INNER JOIN
			agent_staffs AS staff
		ON
			billing.agent_staff_id = staff.id
		INNER JOIN
			agents AS agent
		ON
			staff.agent_id = agent.id
		WHERE
			staff.agent_id IN (%s)
		AND
			job_info.recruitment_state = 0
		AND
		 	job_info.register_phase = 0
		AND
			billing.contract_phase = 2
		AND
			job_info.is_deleted = false
		%s
		ORDER BY job_info.id DESC`,
		idSeparatedByCommas,
		freeWordQuery)

	// クエリを実行
	err = repo.executer.Select(
		repo.Name+".GetActiveAllianceByAgentIDAndFreeWordAndAlliance",
		&jobInformationList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// 求職者uuidを使って打診されている求人を取得
func (repo *JobInformationRepositoryImpl) GetAlreadySoundOutByJobSeekerUUID(jobSeekerUUID uuid.UUID) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	// クエリを実行
	err := repo.executer.Select(
		repo.Name+".GetAlreadySoundOutByJobSeekerUUID",
		&jobInformationList,
		`
			SELECT 
				job_info.*, 
				billing.enterprise_id, 
				billing.agent_staff_id, 
				enterprise.company_name, 
				enterprise.employee_number_single, 
				enterprise.employee_number_group
			FROM
				job_informations AS job_info
			INNER JOIN 
				billing_addresses AS billing
			ON 
				job_info.billing_address_id = billing.id
			INNER JOIN
				enterprise_profiles AS enterprise
			ON
				billing.enterprise_id = enterprise.id
			INNER JOIN
				task_groups
			ON
				job_info.id = task_groups.job_information_id
			INNER JOIN
				job_seekers AS seeker
			ON
				seeker.id = task_groups.job_seeker_id
			WHERE
				seeker.uuid = ?
			ORDER BY task_groups.id DESC
		`, jobSeekerUUID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// AgentIDと外部媒体タイプから求人を取得
func (repo *JobInformationRepositoryImpl) GetByAgentIDAndExternalType(agentID uint, externalType entity.JobInformatinoExternalType) ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	err := repo.executer.Select(
		repo.Name+".GetByAgentIDAndExternalType",
		&jobInformationList, `
			SELECT 
				job_informations.id, job_informations.title,
				job_informations.recruitment_state, 
				job_information_external_ids.external_id, 
				job_information_external_ids.external_type
			FROM 
				job_informations
			INNER JOIN
				job_information_external_ids 
			ON
				job_information_external_ids.job_information_id = job_informations.id
			WHERE
				external_type = ?
		`,
		externalType,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}

// 全ての求人を取得
func (repo *JobInformationRepositoryImpl) All() ([]*entity.JobInformation, error) {
	var (
		jobInformationList []*entity.JobInformation
	)

	err := repo.executer.Select(
		repo.Name+".All",
		&jobInformationList, `
			SELECT *
			FROM job_informations
			WHERE is_deleted = false
		`,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobInformationList, nil
}
