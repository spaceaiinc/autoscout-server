package responses

import (
	"github.com/google/uuid"
	"github.com/spaceaiinc/autoscout-server/domain/entity"
)

type JobSeeker struct {
	JobSeeker *entity.JobSeeker `json:"job_seeker"`
}

func NewJobSeeker(jobSeeker *entity.JobSeeker) JobSeeker {
	return JobSeeker{
		JobSeeker: jobSeeker,
	}
}

type JobSeekerList struct {
	JobSeekerList []*entity.JobSeeker `json:"job_seeker_list"`
}

func NewJobSeekerList(jobSeekers []*entity.JobSeeker) JobSeekerList {
	return JobSeekerList{
		JobSeekerList: jobSeekers,
	}
}

type JobSeekerListAndMaxPageAndIDList struct {
	MaxPageNumber uint                `json:"max_page_number"`
	IDList        []uint              `json:"id_list"`
	JobSeekerList []*entity.JobSeeker `json:"job_seeker_list"`
}

func NewJobSeekerListAndMaxPageAndIDList(jobSeekers []*entity.JobSeeker, maxPageNumber uint, idList []uint) JobSeekerListAndMaxPageAndIDList {
	return JobSeekerListAndMaxPageAndIDList{
		MaxPageNumber: maxPageNumber,
		IDList:        idList,
		JobSeekerList: jobSeekers,
	}
}

type JobSeekerListAndMaxPageAndIDListAndListCount struct {
	MaxPageNumber uint                `json:"max_page_number"`
	IDList        []uint              `json:"id_list"`
	JobSeekerList []*entity.JobSeeker `json:"job_seeker_list"`
	AllCount      uint                `json:"all_count"`
	OwnCount      uint                `json:"own_count"`
	AllianceCount uint                `json:"alliance_count"`
	// HelpCount     uint                `json:"help_count"`
}

func NewJobSeekerListAndMaxPageAndIDListAndListCount(
	jobSeekers []*entity.JobSeeker,
	maxPageNumber uint,
	idList []uint,
	allCount uint,
	ownCount uint,
	allianceCount uint,
	// helpCount uint,
) JobSeekerListAndMaxPageAndIDListAndListCount {
	return JobSeekerListAndMaxPageAndIDListAndListCount{
		MaxPageNumber: maxPageNumber,
		IDList:        idList,
		JobSeekerList: jobSeekers,
		AllCount:      allCount,
		OwnCount:      ownCount,
		AllianceCount: allianceCount,
		// HelpCount:     helpCount,
	}
}

type JobSeekerDocument struct {
	JobSeekerDocument *entity.JobSeekerDocument `json:"job_seeker_document"`
}

func NewJobSeekerDocument(jobSeekerDocument *entity.JobSeekerDocument) JobSeekerDocument {
	return JobSeekerDocument{
		JobSeekerDocument: jobSeekerDocument,
	}
}

type SelectListForCreateOrUpdateJobSeeker struct {
	AgentStaffList               []*entity.AgentStaff               `json:"agent_staff_list"`
	AgentInflowChannelOptionList []*entity.AgentInflowChannelOption `json:"agent_inflow_channel_option_list"`
	AllianceAgentList            []*entity.Agent                    `json:"alliance_agent_list"`
}

func NewSelectListForCreateOrUpdateJobSeeker(agentStaffList []*entity.AgentStaff, agentInflowChannelOptionList []*entity.AgentInflowChannelOption, allianceAgentList []*entity.Agent) SelectListForCreateOrUpdateJobSeeker {
	return SelectListForCreateOrUpdateJobSeeker{
		AgentStaffList:               agentStaffList,
		AgentInflowChannelOptionList: agentInflowChannelOptionList,
		AllianceAgentList:            allianceAgentList,
	}
}

// ゲストページ用の求職者情報
type JobSeekerForGuest struct {
	JobSeekerForGuest *entity.JobSeekerForGuest `json:"job_seeker_for_guest"`
}

func NewJobSeekerForGuest(jobSeekerForGuest *entity.JobSeekerForGuest) JobSeekerForGuest {
	return JobSeekerForGuest{
		JobSeekerForGuest: jobSeekerForGuest,
	}
}

// ゲストページ用の求職者情報
type JobSeekerDesiredForGuest struct {
	JobSeekerForGuest *entity.JobSeekerDesiredForGuest `json:"job_seeker_for_guest"`
}

func NewJobSeekerDesiredForGuest(jobSeekerForGuest *entity.JobSeekerDesiredForGuest) JobSeekerDesiredForGuest {
	return JobSeekerDesiredForGuest{
		JobSeekerForGuest: jobSeekerForGuest,
	}
}

// ゲストページ用の求職者情報
type JobSeekerAgentID struct {
	AgentID uint `json:"agent_id"`
}

func NewJobSeekerAgentID(agentID uint) JobSeekerAgentID {
	return JobSeekerAgentID{
		AgentID: agentID,
	}
}

type JobSeekerUUID struct {
	UUID uuid.UUID `json:"uuid"`
}

func NewJobSeekerUUID(uuid uuid.UUID) JobSeekerUUID {
	return JobSeekerUUID{
		UUID: uuid,
	}
}

// ゲストページ用の求職者情報
type JobSeekerRegisterStatus struct {
	JobSeekerRegisterStatus *entity.JobSeekerRegisterStatus `json:"job_seeker_register_status"`
}

func NewJobSeekerRegisterStatus(jobSeekerRegisterStatus *entity.JobSeekerRegisterStatus) JobSeekerRegisterStatus {
	return JobSeekerRegisterStatus{
		JobSeekerRegisterStatus: jobSeekerRegisterStatus,
	}
}

type JobSeekerLoginFromLP struct {
	UUID       uuid.UUID `json:"uuid"`
	LoginToken uuid.UUID `json:"login_token"`
}

func NewJobSeekerLoginFromLP(uuid, loginToken uuid.UUID) JobSeekerLoginFromLP {
	return JobSeekerLoginFromLP{
		UUID:       uuid,
		LoginToken: loginToken,
	}
}
