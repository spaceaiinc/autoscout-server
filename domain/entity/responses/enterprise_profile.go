package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type EnterpriseProfile struct {
	EnterpriseProfile *entity.EnterpriseProfile `json:"enterprise"`
}

func NewEnterpriseProfile(enterpriseProfile *entity.EnterpriseProfile) EnterpriseProfile {
	return EnterpriseProfile{
		EnterpriseProfile: enterpriseProfile,
	}
}

type EnterpriseProfileList struct {
	EnterpriseProfileList []*entity.EnterpriseProfile `json:"enterprise_list"`
}

func NewEnterpriseProfileList(enterpriseProfiles []*entity.EnterpriseProfile) EnterpriseProfileList {
	return EnterpriseProfileList{
		EnterpriseProfileList: enterpriseProfiles,
	}
}

type EnterpriseProfileListAndMaxPageAndIDList struct {
	MaxPageNumber         uint                        `json:"max_page_number"`
	IDList                []uint                      `json:"id_list"`
	EnterpriseProfileList []*entity.EnterpriseProfile `json:"enterprise_list"`
}

func NewEnterpriseProfileListAndMaxPageAndIDList(enterpriseProfiles []*entity.EnterpriseProfile, maxPageNumber uint, idList []uint) EnterpriseProfileListAndMaxPageAndIDList {
	return EnterpriseProfileListAndMaxPageAndIDList{
		MaxPageNumber:         maxPageNumber,
		IDList:                idList,
		EnterpriseProfileList: enterpriseProfiles,
	}
}

type EnterpriseReferenceMaterial struct {
	EnterpriseReferenceMaterial *entity.EnterpriseReferenceMaterial `json:"enterprise_reference_material"`
}

func NewEnterpriseReferenceMaterial(enterpriseReferenceMaterial *entity.EnterpriseReferenceMaterial) EnterpriseReferenceMaterial {
	return EnterpriseReferenceMaterial{
		EnterpriseReferenceMaterial: enterpriseReferenceMaterial,
	}
}

// 企業と求人情報の一覧をJSONで返す　インポート用
type EnterpriseAndJobInformationList struct {
	EnterpriseAndJobInformationList []*entity.EnterpriseAndJobInformation `json:"enterprise_and_job_information_list"`
}

func NewEnterpriseAndJobInformationList(enterpriseAndJobInformationList []*entity.EnterpriseAndJobInformation) EnterpriseAndJobInformationList {
	return EnterpriseAndJobInformationList{
		EnterpriseAndJobInformationList: enterpriseAndJobInformationList,
	}
}
