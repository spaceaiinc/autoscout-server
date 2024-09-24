package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type InitialEnterpriseImporter struct {
	InitialEnterpriseImporter *entity.InitialEnterpriseImporter `json:"initial_enterprise_importer"`
}

func NewInitialEnterpriseImporter(initialEnterpriseImporter *entity.InitialEnterpriseImporter) InitialEnterpriseImporter {
	return InitialEnterpriseImporter{
		InitialEnterpriseImporter: initialEnterpriseImporter,
	}
}

type InitialEnterpriseImporterList struct {
	InitialEnterpriseImporterList []*entity.InitialEnterpriseImporter `json:"initial_enterprise_importer_list"`
}

func NewInitialEnterpriseImporterList(initialEnterpriseImporters []*entity.InitialEnterpriseImporter) InitialEnterpriseImporterList {
	return InitialEnterpriseImporterList{
		InitialEnterpriseImporterList: initialEnterpriseImporters,
	}
}
