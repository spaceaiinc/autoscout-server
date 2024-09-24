package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SendingJobInformation struct {
	SendingJobInformation *entity.SendingJobInformation `json:"sending_job_information"`
}

func NewSendingJobInformation(sendingJobInformation *entity.SendingJobInformation) SendingJobInformation {
	return SendingJobInformation{
		SendingJobInformation: sendingJobInformation,
	}
}

type SendingJobInformationList struct {
	SendingJobInformationList []*entity.SendingJobInformation `json:"sending_job_information_list"`
}

func NewSendingJobInformationList(sendingJobInformations []*entity.SendingJobInformation) SendingJobInformationList {
	return SendingJobInformationList{
		SendingJobInformationList: sendingJobInformations,
	}
}

type SendingJobInformationListAndMaxPageAndIDList struct {
	MaxPageNumber             uint                            `json:"max_page_number"`
	IDList                    []uint                          `json:"id_list"`
	SendingJobInformationList []*entity.SendingJobInformation `json:"sending_job_information_list"`
}

func NewSendingJobInformationListAndMaxPageAndIDList(sendingJobInformations []*entity.SendingJobInformation, maxPageNumber uint, idList []uint) SendingJobInformationListAndMaxPageAndIDList {
	return SendingJobInformationListAndMaxPageAndIDList{
		MaxPageNumber:             maxPageNumber,
		IDList:                    idList,
		SendingJobInformationList: sendingJobInformations,
	}
}

type SendingJobInformationListAndMaxPageAndIDListAndListCount struct {
	MaxPageNumber             uint                            `json:"max_page_number"`
	IDList                    []uint                          `json:"id_list"`
	SendingJobInformationList []*entity.SendingJobInformation `json:"sending_job_information_list"`
	AllCount                  uint                            `json:"all_count"`
	OwnCount                  uint                            `json:"own_count"`
	AllianceCount             uint                            `json:"alliance_count"`
}

func NewSendingJobInformationListAndMaxPageAndIDListAndListCount(
	sendingJobInformations []*entity.SendingJobInformation,
	maxPageNumber uint,
	idList []uint,
	allCount uint,
	ownCount uint,
	allianceCount uint,
	// helpCount uint,
) SendingJobInformationListAndMaxPageAndIDListAndListCount {
	return SendingJobInformationListAndMaxPageAndIDListAndListCount{
		MaxPageNumber:             maxPageNumber,
		IDList:                    idList,
		SendingJobInformationList: sendingJobInformations,
		AllCount:                  allCount,
		OwnCount:                  ownCount,
		AllianceCount:             allianceCount,
	}
}

type SendingEnterpriseListAndJobInformationListAndMaxPage struct {
	MaxPageNumber             uint                            `json:"max_page_number"`
	SendingEnterpriseList     []*entity.SendingEnterprise     `json:"sending_enterprise_list"`
	SendingJobInformationList []*entity.SendingJobInformation `json:"sending_job_information_list"`
}

func NewSendingEnterpriseListAndJobInformationListAndMaxPage(sendingEnterpriseList []*entity.SendingEnterprise, sendingJobInformations []*entity.SendingJobInformation, maxPageNumber uint) SendingEnterpriseListAndJobInformationListAndMaxPage {
	return SendingEnterpriseListAndJobInformationListAndMaxPage{
		SendingEnterpriseList:     sendingEnterpriseList,
		SendingJobInformationList: sendingJobInformations,
		MaxPageNumber:             maxPageNumber,
	}
}

type JobListingForSending struct {
	JobListingForSending *entity.JobListingForSending `json:"job_listing"`
}

func NewJobListingForSending(jobListingForSending *entity.JobListingForSending) JobListingForSending {
	return JobListingForSending{
		JobListingForSending: jobListingForSending,
	}
}
