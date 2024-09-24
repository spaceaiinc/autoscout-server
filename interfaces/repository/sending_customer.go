package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/interfaces"
	"github.com/spaceaiinc/autoscout-server/usecase"
	"gopkg.in/guregu/null.v4"
)

type SendingCustomerRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewSendingCustomerRepositoryImpl(ex interfaces.SQLExecuter) usecase.SendingCustomerRepository {
	return &SendingCustomerRepositoryImpl{
		Name:     "SendingCustomerRepository",
		executer: ex,
	}
}

func (repo *SendingCustomerRepositoryImpl) Create(sending *entity.SendingCustomer) error {
	nowTime := time.Now().In(time.UTC)
	minTime := utility.EarliestTime() // sending_atの初期値

	lastID, err := repo.executer.Exec(
		repo.Name+".Create",
		`INSERT INTO sending_customers (
			agent_id,
			last_name,
			first_name,
			last_furigana,
			first_furigana,
			phone_number,
			email,
			resume_origin_url,
			resume_pdf_url,
			cv_origin_url,
			cv_pdf_url,
			interview_date,

			interview_information,
			remarks,
			gender,
			nationality,
			nationality_remarks,
			birthday,
			post_code,
			prefecture,
			address,
			address_furigana,

			school_name,
			subject,
			entrance_year,
			graduation_year,
			state_of_employment,
			job_change,
			job_summary,
			company_name,
			joining_year,
			retire_year,

			first_status,
			last_status,
			job_description,
			history_supplement,
			phase,
			sending_at,
			created_at,
			updated_at
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
		sending.AgentID,
		sending.LastName,
		sending.FirstName,
		sending.LastFurigana,
		sending.FirstFurigana,
		sending.PhoneNumber,
		sending.Email,
		sending.ResumeOriginURL,
		sending.ResumePDFURL,
		sending.CVOriginURL,
		sending.CVPDFURL,
		sending.InterviewDate,
		sending.InterviewInformation,
		sending.Remarks,
		sending.Gender,
		sending.Nationality,
		sending.NationalityRemarks,
		sending.Birthday,
		sending.PostCode,
		sending.Prefecture,
		sending.Address,
		sending.AddressFurigana,
		sending.SchoolName,
		sending.Subject,
		sending.EntranceYear,
		sending.GraduationYear,
		sending.StateOfEmployment,
		sending.JobChange,
		sending.JobSummary,
		sending.CompanyName,
		sending.JoiningYear,
		sending.RetireYear,
		sending.FirstStatus,
		sending.LastStatus,
		sending.JobDescription,
		sending.HistorySupplement,
		sending.Phase,
		minTime,
		nowTime,
		nowTime,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	sending.ID = uint(lastID)
	return nil
}

func (repo *SendingCustomerRepositoryImpl) CreateMulti(sendingCustomerList []*entity.SendingCustomer) error {
	var (
		// nowTime   = time.Now().In(time.UTC)
		valuesStr string
		srtFields []string
	)

	// 一斉に作るとjob_seeker作成に渡すsending_customer_idがわからなくなる？
	// for _, sending := range sendingCustomerList {
	// 	srtFields = append(
	// 		srtFields,
	// 		fmt.Sprintf(
	// 			"( %v, %v, %v, %s, %s )",
	// 			nowTime.Format("\"2006-01-02 15:04:05\""),
	// 			nowTime.Format("\"2006-01-02 15:04:05\""),
	// 		),
	// 	)
	// }

	valuesStr = strings.Join(srtFields, ", ")

	query := fmt.Sprintf(`
		INSERT INTO sending_customers (
			sending_customer_group_id,
			adviser_staff_id,
			consultation_status,
			created_at,
			updated_at
		) 
		VALUES %s
	`, valuesStr)

	_, err := repo.executer.Exec(
		repo.Name+".CreateMulti", query,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	// sending.ID = uint(lastID)
	return nil
}

func (repo *SendingCustomerRepositoryImpl) Update(id uint, sending *entity.SendingCustomer) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateSendingCustomerPhaseByID",
		`
		UPDATE sending_customers
		SET
			agent_id = ?,
			last_name = ?,
			first_name = ?,
			last_furigana = ?,
			first_furigana = ?,
			phone_number = ?,
			email = ?,
			resume_origin_url = ?,
			resume_pdf_url = ?,
			cv_origin_url = ?,
			cv_pdf_url = ?,
			interview_date = ?,

			interview_information = ?,
			remarks = ?,
			gender = ?,
			nationality = ?,
			nationality_remarks = ?,
			birthday = ?,
			post_code = ?,
			prefecture = ?,
			address = ?,
			address_furigana = ?,

			school_name = ?,
			subject = ?,
			entrance_year = ?,
			graduation_year = ?,
			state_of_employment = ?,
			job_change = ?,
			job_summary = ?,
			company_name = ?,
			joining_year = ?,
			retire_year = ?,

			first_status = ?,
			last_status = ?,
			job_description = ?,
			history_supplement = ?,
			phase = ?,
			sending_at = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		sending.AgentID,
		sending.LastName,
		sending.FirstName,
		sending.LastFurigana,
		sending.FirstFurigana,
		sending.PhoneNumber,
		sending.Email,
		sending.ResumeOriginURL,
		sending.ResumePDFURL,
		sending.CVOriginURL,
		sending.CVPDFURL,
		sending.InterviewDate,
		sending.InterviewInformation,
		sending.Remarks,
		sending.Gender,
		sending.Nationality,
		sending.NationalityRemarks,
		sending.Birthday,
		sending.PostCode,
		sending.Prefecture,
		sending.Address,
		sending.AddressFurigana,
		sending.SchoolName,
		sending.Subject,
		sending.EntranceYear,
		sending.GraduationYear,
		sending.StateOfEmployment,
		sending.JobChange,
		sending.JobSummary,
		sending.CompanyName,
		sending.JoiningYear,
		sending.RetireYear,
		sending.FirstStatus,
		sending.LastStatus,
		sending.JobDescription,
		sending.HistorySupplement,
		sending.Phase,
		time.Now().In(time.UTC), // sending_at
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingCustomerRepositoryImpl) UpdatePhase(id uint, phase null.Int) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePhase",
		`
		UPDATE sending_customers
		SET
			phase = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		phase,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// sending_job_seekersテーブルのidを使ってphaseを更新する
func (repo *SendingCustomerRepositoryImpl) UpdatePhaseBySendingJobSeekerID(sendingJobSeekerID uint, phase null.Int) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdatePhaseBySendingJobSeekerID",
		`
		UPDATE
			sending_customers
		INNER JOIN
			sending_job_seekers
		ON
			sending_customers.id = sending_job_seekers.sending_customer_id
		SET
			sending_customers.phase = ?,
			sending_customers.updated_at = ?
		WHERE 
			sending_job_seekers.id = ?
		`,
		phase,
		time.Now().In(time.UTC),
		sendingJobSeekerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 面談日時の更新
func (repo *SendingCustomerRepositoryImpl) UpdateInterviewDate(sendingCustomerID uint, interviewDate time.Time) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateInterviewDate",
		`
		UPDATE sending_customers
		SET
			interview_date = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		interviewDate,
		time.Now().In(time.UTC),
		sendingCustomerID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingCustomerRepositoryImpl) UpdateDocument(id uint, resumePDFURL, resumeOriginURL, cvPDFURL, cvOriginURL string) error {
	_, err := repo.executer.Exec(
		repo.Name+".UpdateDocument",
		`
		UPDATE sending_customers
		SET
			resume_pdf_url = ?,
			resume_origin_url = ?,
			cv_pdf_url = ?,
			cv_origin_url = ?,
			updated_at = ?
		WHERE 
			id = ?
		`,
		resumePDFURL,
		resumeOriginURL,
		cvPDFURL,
		cvOriginURL,
		time.Now().In(time.UTC),
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// 開発環境で実行時にテストユーザーの名前で更新する
func (repo *SendingCustomerRepositoryImpl) UpdateForDev(sendingCustomer *entity.SendingCustomer) error {
	sendingCustomer.LastName = "テスト"
	sendingCustomer.FirstName = fmt.Sprintf("ユーザー%v", sendingCustomer.ID)
	sendingCustomer.LastFurigana = "テスト"
	sendingCustomer.FirstFurigana = fmt.Sprintf("ユーザー%v", sendingCustomer.ID)
	sendingCustomer.PhoneNumber = "080-0000-0000"

	_, err := repo.executer.Exec(
		repo.Name+".UpdateForDev",
		`
		UPDATE sending_customers
		SET
			last_name = ?,
			first_name = ?,
			last_furigana = ?,
			first_furigana = ?,
			phone_number = ?,
			updated_at = ?
		WHERE
			id = ?
		`,
		sendingCustomer.LastName,
		sendingCustomer.FirstName,
		sendingCustomer.LastFurigana,
		sendingCustomer.FirstFurigana,
		sendingCustomer.PhoneNumber,
		time.Now().In(time.UTC),
		sendingCustomer.ID,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (repo *SendingCustomerRepositoryImpl) FindByID(id uint) (*entity.SendingCustomer, error) {
	var (
		sendingCustomer entity.SendingCustomer
	)

	err := repo.executer.Get(
		repo.Name+".FindByID",
		&sendingCustomer,
		`SELECT
			sending.*
		FROM 
			sending_customers AS sending
		WHERE 
			sending.id = ?
		LIMIT 1
	`,
		id,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &sendingCustomer, err
}

func (repo *SendingCustomerRepositoryImpl) GetListByAgentID(agentID uint) ([]*entity.SendingCustomer, error) {
	var (
		sendingCustomerList []*entity.SendingCustomer
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentID",
		&sendingCustomerList,
		`SELECT 
			sending.*
		FROM 
			sending_customers AS sending
		WHERE 
			sending.agent_id = ?
		ORDER BY sending.id DESC
	`,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingCustomerList, err
}

func (repo *SendingCustomerRepositoryImpl) GetForSendingCustomerManagement(agentID uint, freeWord string) ([]*entity.SendingCustomer, error) {
	var (
		sendingCustomerList []*entity.SendingCustomer
		freeWordQuery       string
	)

	freeWordInt, err := strconv.Atoi(freeWord)
	if err != nil {
		// 検索ワードがある場合
		if freeWord != "" {
			freeWordForLike := "%" + freeWord + "%"

			freeWordQuery = fmt.Sprintf(`
				AND (	
					CONCAT(customer.last_name, customer.first_name) LIKE '%s' OR 
					CONCAT(customer.last_furigana, customer.first_furigana) LIKE '%s'
				)
			`, freeWordForLike, freeWordForLike)
		}
	} else {
		// ハイフンなしの場合数字として扱われるためハイフンなしはこちらに入れる
		freeWordQuery = fmt.Sprintf(`
			AND customer.id = %d
		`, freeWordInt)
	}

	query := fmt.Sprintf(`
		SELECT 
			customer.*, COUNT(phase.sending_job_seeker_id) AS sending_count
		FROM 
			sending_customers AS customer
		INNER JOIN
			sending_job_seekers AS seeker
		ON
			customer.id = seeker.sending_customer_id
		LEFT JOIN 
			sending_phases AS phase
		ON 
			seeker.id = phase.sending_job_seeker_id
		WHERE 
			customer.agent_id = ?
		%s
		GROUP BY
			customer.id,
			customer.agent_id,
			customer.last_name,
			customer.first_name,
			customer.last_furigana,
			customer.first_furigana,
			customer.phone_number,
			customer.email,
			customer.resume_pdf_url,
			customer.cv_pdf_url,
			customer.interview_date,
			customer.interview_information,
			customer.remarks,
			customer.gender,
			customer.nationality,
			customer.nationality_remarks,
			customer.birthday,
			customer.post_code,
			customer.prefecture,
			customer.address,
			customer.address_furigana,
			customer.school_name,
			customer.subject,
			customer.entrance_year,
			customer.graduation_year,
			customer.state_of_employment,
			customer.job_change,
			customer.job_summary,
			customer.company_name,
			customer.joining_year,
			customer.retire_year,
			customer.first_status,
			customer.last_status,
			customer.job_description,
			customer.history_supplement,
			customer.phase,
			customer.sending_at,
			customer.created_at
		ORDER BY customer.id DESC
	`, freeWordQuery)

	err = repo.executer.Select(
		repo.Name+".GetForSendingCustomerManagement",
		&sendingCustomerList,
		query,
		agentID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingCustomerList, err
}

func (repo *SendingCustomerRepositoryImpl) GetListByAgentIDAndPhase(agentID uint, phase []uint) ([]*entity.SendingCustomer, error) {
	var (
		sendingCustomerList []*entity.SendingCustomer
	)

	query := fmt.Sprintf(
		`SELECT 
			sending.*
		FROM 
			sending_customers AS sending
		WHERE
			sending.agent_id = (%d)
		AND
			sending.phase IN (%s)
		ORDER BY 
			sending.id DESC
		`,
		agentID,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(phase)), ", "), "[]"),
	)

	err := repo.executer.Select(
		repo.Name+".GetListByAgentIDAndPhase",
		&sendingCustomerList,
		query,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sendingCustomerList, err
}
