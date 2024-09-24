package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type DeploymentInformation struct {
	DeploymentInformation *entity.DeploymentInformation `json:"deployment_information"`
}

func NewDeploymentInformation(deploymentInformation *entity.DeploymentInformation) DeploymentInformation {
	return DeploymentInformation{
		DeploymentInformation: deploymentInformation,
	}
}

type DeploymentInformationList struct {
	DeploymentInformationList []*entity.DeploymentInformation `json:"deployment_information_list"`
}

func NewDeploymentInformationList(deploymentInformationList []*entity.DeploymentInformation) DeploymentInformationList {
	return DeploymentInformationList{
		DeploymentInformationList: deploymentInformationList,
	}
}

type DeploymentInformationListAndMaxPage struct {
	MaxPageNumber             uint                            `json:"max_page_number"`
	DeploymentInformationList []*entity.DeploymentInformation `json:"deployment_information_list"`
}

func NewDeploymentInformationListAndMaxPage(maxPageNumber uint, deploymentInformationList []*entity.DeploymentInformation) DeploymentInformationListAndMaxPage {
	return DeploymentInformationListAndMaxPage{
		MaxPageNumber:             maxPageNumber,
		DeploymentInformationList: deploymentInformationList,
	}
}
