package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SendingEnterprise struct {
	SendingEnterprise *entity.SendingEnterprise `json:"sending_enterprise"`
}

func NewSendingEnterprise(sendingEnterprise *entity.SendingEnterprise) SendingEnterprise {
	return SendingEnterprise{
		SendingEnterprise: sendingEnterprise,
	}
}

type SendingEnterpriseList struct {
	SendingEnterpriseList []*entity.SendingEnterprise `json:"sending_enterprise_list"`
}

func NewSendingEnterpriseList(sendingEnterprises []*entity.SendingEnterprise) SendingEnterpriseList {
	return SendingEnterpriseList{
		SendingEnterpriseList: sendingEnterprises,
	}
}

type SendingEnterpriseListAndMaxPageAndIDList struct {
	MaxPageNumber         uint                        `json:"max_page_number"`
	IDList                []uint                      `json:"id_list"`
	SendingEnterpriseList []*entity.SendingEnterprise `json:"sending_enterprise_list"`
}

func NewSendingEnterpriseListAndMaxPageAndIDList(sendingEnterprises []*entity.SendingEnterprise, maxPageNumber uint, idList []uint) SendingEnterpriseListAndMaxPageAndIDList {
	return SendingEnterpriseListAndMaxPageAndIDList{
		MaxPageNumber:         maxPageNumber,
		IDList:                idList,
		SendingEnterpriseList: sendingEnterprises,
	}
}

type SendingEnterpriseReferenceMaterial struct {
	SendingEnterpriseReferenceMaterial *entity.SendingEnterpriseReferenceMaterial `json:"sending_enterprise_reference_material"`
}

func NewSendingEnterpriseReferenceMaterial(sendingEnterpriseReferenceMaterial *entity.SendingEnterpriseReferenceMaterial) SendingEnterpriseReferenceMaterial {
	return SendingEnterpriseReferenceMaterial{
		SendingEnterpriseReferenceMaterial: sendingEnterpriseReferenceMaterial,
	}
}

type SendingEnterpriseForSignin struct {
	SendingEnterprise    *entity.SendingEnterprise    `json:"sending_enterprise"`
	SendingJobSeeker     *entity.SendingJobSeeker     `json:"sending_job_seeker"`
	SendingShareDocument *entity.SendingShareDocument `json:"sending_share_document"`
}

func NewSendingEnterpriseForSignin(
	sendingEnterprise *entity.SendingEnterprise,
	sendingJobSeeker *entity.SendingJobSeeker,
	sendingShareDocument *entity.SendingShareDocument,
) SendingEnterpriseForSignin {
	return SendingEnterpriseForSignin{
		SendingEnterprise:    sendingEnterprise,
		SendingJobSeeker:     sendingJobSeeker,
		SendingShareDocument: sendingShareDocument,
	}
}
