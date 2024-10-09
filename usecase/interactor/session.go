package interactor

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
	"github.com/spaceaiinc/autoscout-server/domain/utility"
	"github.com/spaceaiinc/autoscout-server/usecase"
)

type SessionInteractor interface {
	// Gest API
	SignIn(input SessionSignInInput) (SessionSignInOutput, error)
	SignOut(input SessionSignOutInput) (SessionSignOutOutput, error)
	GetSignInUser(input GetSignInUserInput) (GetSignInUserOutput, error)

	SignInForGestEnterprise(input SessionSignInForGestEnterpriseInput) (SessionSignInForGestEnterpriseOutput, error)
	SignInForGestEnterpriseByTaskGroupUUID(input SessionSignInForGestEnterpriseByTaskGroupUUIDInput) (SessionSignInForGestEnterpriseByTaskGroupUUIDOutput, error)
	SignInForGestJobSeeker(input SessionSignInForGestJobSeekerInput) (SessionSignInForGestJobSeekerOutput, error)
	SignInForGestJobSeekerFromLP(input SessionSignInForGestJobSeekerFromLPInput) (SessionSignInForGestJobSeekerFromLPOutput, error)
	SignInForGestSendingJobSeeker(input SessionSignInForGestSendingJobSeekerInput) (SessionSignInForGestSendingJobSeekerOutput, error)

	// LP
	LoginGestJobSeekerForLP(input LoginGestJobSeekerForLPInput) (LoginGestJobSeekerForLPOutput, error)
}

type SessionInteractorImpl struct {
	firebase                        usecase.Firebase
	sendgrid                        config.Sendgrid
	agentRepository                 usecase.AgentRepository
	agentStaffRepository            usecase.AgentStaffRepository
	enterpriseProfileRepository     usecase.EnterpriseProfileRepository
	jobInformationRepository        usecase.JobInformationRepository
	jobSeekerRepository             usecase.JobSeekerRepository
	jobSeekerLPLoginTokenRepository usecase.JobSeekerLPLoginTokenRepository
	agentAllianceRepository         usecase.AgentAllianceRepository
	sendingJobSeekerRepository      usecase.SendingJobSeekerRepository
}

func NewSessionInteractorImpl(
	fb usecase.Firebase,
	sg config.Sendgrid,
	aR usecase.AgentRepository,
	asR usecase.AgentStaffRepository,
	epR usecase.EnterpriseProfileRepository,
	jR usecase.JobInformationRepository,
	jsR usecase.JobSeekerRepository,
	jlltR usecase.JobSeekerLPLoginTokenRepository,
	aaR usecase.AgentAllianceRepository,
	sjsR usecase.SendingJobSeekerRepository,
) SessionInteractor {
	return &SessionInteractorImpl{
		firebase:                        fb,
		sendgrid:                        sg,
		agentRepository:                 aR,
		agentStaffRepository:            asR,
		enterpriseProfileRepository:     epR,
		jobInformationRepository:        jR,
		jobSeekerRepository:             jsR,
		jobSeekerLPLoginTokenRepository: jlltR,
		agentAllianceRepository:         aaR,
		sendingJobSeekerRepository:      sjsR,
	}
}

type SessionSignInInput struct {
	Token string
}

type SessionSignInOutput struct {
	User *entity.User
}

func (i *SessionInteractorImpl) SignIn(input SessionSignInInput) (SessionSignInOutput, error) {
	var (
		output = SessionSignInOutput{}
	)

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	// 担当者情報の取得
	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// 現在時刻
	now := time.Now().In(time.UTC)

	err = i.agentStaffRepository.UpdateLastLogin(agentStaff.ID, now)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントアカウントのログインユーザー情報を作成
	user := entity.NewUser(
		agentStaff.ID,
		agentStaff.AgentID,
		agentStaff.AgentUUID,
		agentStaff.FirebaseID,
		agentStaff.StaffName,
		agentStaff.AgentName,
		agentStaff.LastLogin,
		agentStaff.Authority,
		agentStaff.UsageStatus,
		agentStaff.IsCRMActive,
		agentStaff.IsAllianceActive,
		agentStaff.IsSendingActive,
		agentStaff.SendingType,
		agentStaff.AgentID == 1,
	)

	// 指定AgentIDのアライアンス情報を取得
	var allianceIDList []uint

	agentAllianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(agentStaff.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// アライアンスエージェントのIDリストを作成
	for _, alliance := range agentAllianceList {
		if alliance.Agent1ID != agentStaff.AgentID {
			allianceIDList = append(allianceIDList, alliance.Agent1ID)
		} else {
			allianceIDList = append(allianceIDList, alliance.Agent2ID)
		}
	}

	user.AllianceAgentIDs = allianceIDList

	output.User = user

	return output, nil
}

type SessionSignOutInput struct {
	Token string
}

type SessionSignOutOutput struct {
	OK bool
}

func (i *SessionInteractorImpl) SignOut(input SessionSignOutInput) (SessionSignOutOutput, error) {
	var (
		output = SessionSignOutOutput{}
	)

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	// Firebase Authenticationのログアウト処理（全ての端末でログアウト）
	err = i.firebase.SignOut(firebaseID)
	if err != nil {
		return output, err
	}

	output.OK = true

	return output, nil
}

type GetSignInUserInput struct {
	Token string
}

type GetSignInUserOutput struct {
	User *entity.User
}

func (i *SessionInteractorImpl) GetSignInUser(input GetSignInUserInput) (GetSignInUserOutput, error) {
	var (
		output = GetSignInUserOutput{}
	)

	firebaseID, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	// 担当者情報の取得
	agentStaff, err := i.agentStaffRepository.FindByFirebaseID(firebaseID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントアカウントのログインユーザー情報を作成
	user := entity.NewUser(
		agentStaff.ID,
		agentStaff.AgentID,
		agentStaff.AgentUUID,
		agentStaff.FirebaseID,
		agentStaff.StaffName,
		agentStaff.AgentName,
		agentStaff.LastLogin,
		agentStaff.Authority,
		agentStaff.UsageStatus,
		agentStaff.IsCRMActive,
		agentStaff.IsAllianceActive,
		agentStaff.IsSendingActive,
		agentStaff.SendingType,
		agentStaff.AgentID == 1,
	)

	// 指定AgentIDのアライアンス情報を取得
	var allianceIDList []uint

	agentAllianceList, err := i.agentAllianceRepository.GetByAgentIDAndRequestDone(agentStaff.AgentID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// アライアンスエージェントのIDリストを作成
	for _, alliance := range agentAllianceList {
		if alliance.Agent1ID != agentStaff.AgentID {
			allianceIDList = append(allianceIDList, alliance.Agent1ID)
		} else {
			allianceIDList = append(allianceIDList, alliance.Agent2ID)
		}
	}

	user.AllianceAgentIDs = allianceIDList

	output.User = user

	return output, nil
}

type SessionSignInForGestEnterpriseInput struct {
	Password           string
	JobInformationUUID uuid.UUID
}

type SessionSignInForGestEnterpriseOutput struct {
	User *entity.GestEnterpriseUser
}

func (i *SessionInteractorImpl) SignInForGestEnterprise(input SessionSignInForGestEnterpriseInput) (SessionSignInForGestEnterpriseOutput, error) {
	var (
		output = SessionSignInForGestEnterpriseOutput{}
	)

	// 企業のログイン
	enterprise, err := i.enterpriseProfileRepository.CheckPostCode(input.Password)
	if err != nil {
		return output, err
	}

	// エージェントアカウントのログインユーザー情報を作成
	gestEnterprise := entity.NewGestEnterpriseUser(
		input.JobInformationUUID,
		enterprise.CompanyName,
	)

	output.User = gestEnterprise

	return output, nil
}

type SessionSignInForGestEnterpriseByTaskGroupUUIDInput struct {
	Password      string
	TaskGroupUUID uuid.UUID
}

type SessionSignInForGestEnterpriseByTaskGroupUUIDOutput struct {
	User *entity.GestEnterpriseUser
}

func (i *SessionInteractorImpl) SignInForGestEnterpriseByTaskGroupUUID(input SessionSignInForGestEnterpriseByTaskGroupUUIDInput) (SessionSignInForGestEnterpriseByTaskGroupUUIDOutput, error) {
	var (
		output = SessionSignInForGestEnterpriseByTaskGroupUUIDOutput{}
	)

	// 企業のログイン
	enterprise, err := i.enterpriseProfileRepository.CheckPostCode(input.Password)
	if err != nil {
		return output, err
	}

	jobInformation, err := i.jobInformationRepository.FindByTaskGroupUUID(input.TaskGroupUUID)
	if err != nil {
		return output, err
	}

	// エージェントアカウントのログインユーザー情報を作成
	gestEnterprise := entity.NewGestEnterpriseUser(
		jobInformation.UUID,
		enterprise.CompanyName,
	)

	output.User = gestEnterprise

	return output, nil
}

type SessionSignInForGestJobSeekerInput struct {
	Password      string
	JobSeekerUUID uuid.UUID
}

type SessionSignInForGestJobSeekerOutput struct {
	User *entity.GestJobSeekerUser
}

func (i *SessionInteractorImpl) SignInForGestJobSeeker(input SessionSignInForGestJobSeekerInput) (SessionSignInForGestJobSeekerOutput, error) {
	var (
		output = SessionSignInForGestJobSeekerOutput{}
	)

	// 求職者ログイン
	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(errors.New("URLが正しくありません。\nご確認の上もう一度お試しください。"))
		return output, err
	}

	// パスワードが未入力の場合は電話番号の下4桁と比較
	if jobSeeker.Password == "" {
		if jobSeeker.PhoneNumber == "" {
			// パスワード未設定、かつ電話番号未入力の場合は再設定を促す
			wrapped := fmt.Errorf("パスワードが設定されていません。パスワードお忘れの方よりパスワードを再設定してください。")
			return output, wrapped
		} else {
			result := jobSeeker.PhoneNumber[len(jobSeeker.PhoneNumber)-4:]
			fmt.Println("result", result, "Password", input.Password)
			if result != input.Password {
				wrapped := fmt.Errorf("パスワードが正しくありません。")
				return output, wrapped
			}
		}
	} else {
		// パスワードが入力済みの場合は比較
		err = compareHashedPaasowd(jobSeeker.Password, input.Password)
		if err != nil {
			wrapped := fmt.Errorf("パスワードが正しくありません。")
			return output, wrapped
		}
	}

	// エージェントアカウントのログインユーザー情報を作成
	gestJobSeeker := entity.NewGestJobSeekerUser(
		jobSeeker.ID,
		input.JobSeekerUUID,
		jobSeeker.LastName,
		jobSeeker.FirstName,
		jobSeeker.Email,
		jobSeeker.AgentID,
		jobSeeker.Phase,
		jobSeeker.CanViewMatchingJob,
	)

	output.User = gestJobSeeker

	return output, nil
}

// マイページログイン（LPからログイントークンを使ってログイン）
type SessionSignInForGestJobSeekerFromLPInput struct {
	JobSeekerUUID uuid.UUID
	LoginToken    uuid.UUID
}

type SessionSignInForGestJobSeekerFromLPOutput struct {
	User *entity.GestJobSeekerUser
}

func (i *SessionInteractorImpl) SignInForGestJobSeekerFromLP(input SessionSignInForGestJobSeekerFromLPInput) (SessionSignInForGestJobSeekerFromLPOutput, error) {
	var (
		output = SessionSignInForGestJobSeekerFromLPOutput{}
	)

	// 求職者ログイン
	jobSeeker, err := i.jobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(errors.New("URLが正しくありません。\nご確認の上もう一度お試しください。"))
		return output, err
	}

	loginTokenFromLP, err := i.jobSeekerLPLoginTokenRepository.FindByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// ログイントークンが一致しない場合はエラー返す
	if loginTokenFromLP.LoginToken != input.LoginToken {
		wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "URLが正しくありません。\nご確認の上もう一度お試しください。")
		return output, wrapped
	}

	err = i.jobSeekerLPLoginTokenRepository.DeleteByJobSeekerID(jobSeeker.ID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	// エージェントアカウントのログインユーザー情報を作成
	gestJobSeeker := entity.NewGestJobSeekerUser(
		jobSeeker.ID,
		input.JobSeekerUUID,
		jobSeeker.LastName,
		jobSeeker.FirstName,
		jobSeeker.Email,
		jobSeeker.AgentID,
		jobSeeker.Phase,
		jobSeeker.CanViewMatchingJob,
	)

	output.User = gestJobSeeker

	return output, nil
}

type SessionSignInForGestSendingJobSeekerInput struct {
	Password      string
	JobSeekerUUID uuid.UUID
}

type SessionSignInForGestSendingJobSeekerOutput struct {
	User *entity.GestJobSeekerUser
}

func (i *SessionInteractorImpl) SignInForGestSendingJobSeeker(input SessionSignInForGestSendingJobSeekerInput) (SessionSignInForGestSendingJobSeekerOutput, error) {
	var (
		output = SessionSignInForGestSendingJobSeekerOutput{}
	)

	// 求職者ログイン
	sendingJobSeeker, err := i.sendingJobSeekerRepository.FindByUUID(input.JobSeekerUUID)
	if err != nil {
		fmt.Println(err)
		return output, err
	}

	result := sendingJobSeeker.PhoneNumber[len(sendingJobSeeker.PhoneNumber)-4:]
	fmt.Println("result", result, "Password", input.Password)

	if result != input.Password {
		return output, errors.New("パスワードが正しくありません。")
	} else {
		// エージェントアカウントのログインユーザー情報を作成
		gestJobSeeker := entity.NewGestJobSeekerUser(
			sendingJobSeeker.ID,
			input.JobSeekerUUID,
			sendingJobSeeker.LastName,
			sendingJobSeeker.FirstName,
			sendingJobSeeker.Email,
			sendingJobSeeker.AgentID,
			sendingJobSeeker.Phase,
			false,
		)

		output.User = gestJobSeeker
	}

	return output, nil
}

// LPのログインフォームのログイン処理
type LoginGestJobSeekerForLPInput struct {
	Email    string
	Password string
}

type LoginGestJobSeekerForLPOutput struct {
	UUID       uuid.UUID
	LoginToken uuid.UUID
}

func (i *SessionInteractorImpl) LoginGestJobSeekerForLP(input LoginGestJobSeekerForLPInput) (LoginGestJobSeekerForLPOutput, error) {
	var (
		output              = LoginGestJobSeekerForLPOutput{}
		SystemAgentID uint = 1
	)

	// 求職者のメールアドレスが合致するか確認
	jobSeeker, err := i.jobSeekerRepository.FindByEmailForLP(input.Email, SystemAgentID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "メールアドレスまたはパスワードが一致しません。")
			return output, wrapped
		} else {
			// Not Found以外のエラーの場合はそのままサーバーエラー
			fmt.Println(err)
			return output, err
		}
	}

	// パスワードが未設定の場合は電話番号の下4桁と比較
	if jobSeeker.Password == "" {
		result := jobSeeker.PhoneNumber[len(jobSeeker.PhoneNumber)-4:]
		fmt.Println("result", result, "Password", input.Password)
		if result == input.Password {
			wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "メールアドレスまたはパスワードが一致しません。")
			return output, wrapped
		}

	} else {
		// パスワードが入力済みの場合は比較
		err = compareHashedPaasowd(jobSeeker.Password, input.Password)
		if err != nil {
			wrapped := fmt.Errorf("%w:%s", entity.ErrRequestError, "メールアドレスまたはパスワードが一致しません。")
			return output, wrapped
		}
	}

	// ログイントークンがない場合は新たに作成
	logintTokenFromLP := entity.NewJobSeekerLPLoginToken(
		jobSeeker.ID, utility.CreateUUID(),
	)

	_, err = i.jobSeekerLPLoginTokenRepository.FindByJobSeekerID(jobSeeker.ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			err = i.jobSeekerLPLoginTokenRepository.Create(logintTokenFromLP)
			if err != nil {
				fmt.Println(err)
				return output, err
			}
		} else {
			// Not Found以外のエラーの場合はそのままサーバーエラー
			fmt.Println(err)
			return output, err
		}
	} else {
		// ログイントークンがある場合は更新
		err = i.jobSeekerLPLoginTokenRepository.UpdateByJobSeekerID(jobSeeker.ID, logintTokenFromLP)
		if err != nil {
			fmt.Println(err)
			return output, err
		}
	}

	output.UUID = jobSeeker.UUID
	output.LoginToken = logintTokenFromLP.LoginToken

	return output, nil
}
