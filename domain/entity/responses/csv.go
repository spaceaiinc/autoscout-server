package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type MissedRecordsAndOK struct {
	MissedRecords []uint `json:"missed_records"`
	OK            bool   `json:"ok"`
}

func NewMissedRecordsAndOK(missedRecords []uint, ok bool) MissedRecordsAndOK {
	return MissedRecordsAndOK{
		MissedRecords: missedRecords,
		OK:            ok,
	}
}

type MissedRecordsAndEnterpriseList struct {
	MissedRecords  []uint                                `json:"missed_records"`
	EnterpriseList []*entity.EnterpriseAndBillingAddress `json:"enterprise_list"`
}

func NewMissedRecordsAndEnterpriseList(missedRecords []uint, enterprise []*entity.EnterpriseAndBillingAddress) MissedRecordsAndEnterpriseList {
	return MissedRecordsAndEnterpriseList{
		MissedRecords:  missedRecords,
		EnterpriseList: enterprise,
	}
}

type MissedRecordsAndJobInformationList struct {
	MissedRecords      []uint                   `json:"missed_records"`
	JobInformationList []*entity.JobInformation `json:"job_information_list"`
}

func NewMissedRecordsAndJobInformationList(missedRecords []uint, jobInformation []*entity.JobInformation) MissedRecordsAndJobInformationList {
	return MissedRecordsAndJobInformationList{
		MissedRecords:      missedRecords,
		JobInformationList: jobInformation,
	}
}

type MissedRecordsAndJobSeekerList struct {
	MissedRecords []uint              `json:"missed_records"`
	JobSeekerList []*entity.JobSeeker `json:"job_seeker_list"`
}

func NewMissedRecordsAndJobSeekerList(missedRecords []uint, jobSeeker []*entity.JobSeeker) MissedRecordsAndJobSeekerList {
	return MissedRecordsAndJobSeekerList{
		MissedRecords: missedRecords,
		JobSeekerList: jobSeeker,
	}
}
