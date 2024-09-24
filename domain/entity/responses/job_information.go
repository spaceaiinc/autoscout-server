package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type JobInformation struct {
	JobInformation *entity.JobInformation `json:"job_information"`
}

func NewJobInformation(jobInformation *entity.JobInformation) JobInformation {
	return JobInformation{
		JobInformation: jobInformation,
	}
}

type JobInformationList struct {
	JobInformationList []*entity.JobInformation `json:"job_information_list"`
}

func NewJobInformationList(jobInformations []*entity.JobInformation) JobInformationList {
	return JobInformationList{
		JobInformationList: jobInformations,
	}
}

type JobInformationListAndMaxPageAndIDList struct {
	MaxPageNumber      uint                     `json:"max_page_number"`
	IDList             []uint                   `json:"id_list"`
	JobInformationList []*entity.JobInformation `json:"job_information_list"`
}

func NewJobInformationListAndMaxPageAndIDList(jobInformations []*entity.JobInformation, maxPageNumber uint, idList []uint) JobInformationListAndMaxPageAndIDList {
	return JobInformationListAndMaxPageAndIDList{
		MaxPageNumber:      maxPageNumber,
		IDList:             idList,
		JobInformationList: jobInformations,
	}
}

type JobInformationListAndMaxPageAndIDListAndListCount struct {
	MaxPageNumber      uint                     `json:"max_page_number"`
	IDList             []uint                   `json:"id_list"`
	JobInformationList []*entity.JobInformation `json:"job_information_list"`
	AllCount           uint                     `json:"all_count"`
	OwnCount           uint                     `json:"own_count"`
	AllianceCount      uint                     `json:"alliance_count"`
}

func NewJobInformationListAndMaxPageAndIDListAndListCount(
	jobInformations []*entity.JobInformation,
	maxPageNumber uint,
	idList []uint,
	allCount uint,
	ownCount uint,
	allianceCount uint,
) JobInformationListAndMaxPageAndIDListAndListCount {
	return JobInformationListAndMaxPageAndIDListAndListCount{
		MaxPageNumber:      maxPageNumber,
		IDList:             idList,
		JobInformationList: jobInformations,
		AllCount:           allCount,
		OwnCount:           ownCount,
		AllianceCount:      allianceCount,
	}
}

type JobListing struct {
	JobListing *entity.JobListing `json:"job_listing"`
}

func NewJobListing(jobListing *entity.JobListing) JobListing {
	return JobListing{
		JobListing: jobListing,
	}
}

type JobListingList struct {
	JobListingList []*entity.JobListing `json:"job_listing_list"`
}

func NewJobListingList(jobListingList []*entity.JobListing) JobListingList {
	return JobListingList{
		JobListingList: jobListingList,
	}
}

type JobListingListForJobSeeker struct {
	NotYetEntryJobListingList    []*entity.JobListing `json:"not_yet_entry_job_listing_list"`    // 未エントリー
	AcceptJobOfferJobListingList []*entity.JobListing `json:"accept_job_offer_job_listing_list"` // 内定承諾
	HoldJobOfferJobListingList   []*entity.JobListing `json:"hold_job_offer_job_listing_list"`   // 内定
	SelectionJobListingList      []*entity.JobListing `json:"selection_job_listing_list"`        // 選考中
	EndJobListingList            []*entity.JobListing `json:"end_job_listing_list"`              // 終了
}

func NewJobListingListForJobSeeker(
	notYetEntryJobListingList []*entity.JobListing, // 未エントリー
	acceptJobOfferJobListingList []*entity.JobListing, // 内定承諾
	holdJobOfferJobListingList []*entity.JobListing, // 内定
	selectionJobListingList []*entity.JobListing, // 選考中
	endJobListingList []*entity.JobListing, // 終了
) JobListingListForJobSeeker {
	return JobListingListForJobSeeker{
		NotYetEntryJobListingList:    notYetEntryJobListingList,
		AcceptJobOfferJobListingList: acceptJobOfferJobListingList,
		HoldJobOfferJobListingList:   holdJobOfferJobListingList,
		SelectionJobListingList:      selectionJobListingList,
		EndJobListingList:            endJobListingList,
	}
}

type JobInformationCountAndIncome struct {
	//	応募資格に合致した求人: 条件で絞った時の求人数
	Count uint `json:"count"`

	// 未経験職で応募資格に合致した求人: 上記の求人の中で求人の特徴に「職種未経験OK」or「業界・職種未経験OK」が選択されていて、かつ「求職者が経験していない職種」の求人数
	GuaranteedInterviewCount uint `json:"guaranteedinterview_count"`

	// 年収期待値: 求人の年収上限の平均値
	ExpectedIncome uint `json:"expected_income"`
}

func NewJobInformationCountAndIncome(
	count uint,
	guaranteedInterviewCount uint,
	expectedIncome uint,
) JobInformationCountAndIncome {
	return JobInformationCountAndIncome{
		Count:                    count,
		GuaranteedInterviewCount: guaranteedInterviewCount,
		ExpectedIncome:           expectedIncome,
	}
}

type JobListingListAndMaxPage struct {
	JobListingList []*entity.JobListing `json:"job_listing_list"`
	MaxPageNumber  uint                 `json:"max_page_number"`
}

func NewJobListingListAndMaxPage(
	jobListing []*entity.JobListing,
	maxPageNumber uint,
) JobListingListAndMaxPage {
	return JobListingListAndMaxPage{
		JobListingList: jobListing,
		MaxPageNumber:  maxPageNumber,
	}
}

type JobListingListAndCount struct {
	JobListingList []*entity.JobListing `json:"job_listing_list"`

	//	応募資格に合致した求人: 条件で絞った時の求人数
	Count uint `json:"count"`

	// 未経験職で応募資格に合致した求人: 上記の求人の中で求人の特徴に「職種未経験OK」or「業界・職種未経験OK」が選択されていて、かつ「求職者が経験していない職種」の求人数
	CountOfGuaranteedInterview uint `json:"count_of_guaranteed_interview"`

	// 年収期待値: 年収の最小値
	ExpectedUnderIncome uint `json:"expected_under_income"`

	// 年収期待値: 年収の最大値
	ExpectedOverIncome uint `json:"expected_over_income"`

	// 最大ページ数
	MaxPageNumber uint `json:"max_page_number"`
}

func NewJobListingListAndCount(
	jobListing []*entity.JobListing,
	count uint,
	countOfGuaranteedInterview uint,
	expectedUnderIncome uint,
	expectedOverIncome uint,
	maxPageNumber uint,
) JobListingListAndCount {
	return JobListingListAndCount{
		JobListingList:             jobListing,
		Count:                      count,
		CountOfGuaranteedInterview: countOfGuaranteedInterview,
		ExpectedUnderIncome:        expectedUnderIncome,
		ExpectedOverIncome:         expectedOverIncome,
		MaxPageNumber:              maxPageNumber,
	}
}

type JobListingListAndJobSeekerDesired struct {
	JobListingList    []*entity.JobListing             `json:"job_listing_list"`
	JobSeekerForGuest *entity.JobSeekerDesiredForGuest `json:"job_seeker_for_guest"`
}

func NewJobListingListAndJobSeekerDesired(
	jobListingList []*entity.JobListing,
	jobSeekerForGuest *entity.JobSeekerDesiredForGuest,
) JobListingListAndJobSeekerDesired {
	return JobListingListAndJobSeekerDesired{
		JobListingList:    jobListingList,
		JobSeekerForGuest: jobSeekerForGuest,
	}
}

type JobInformationListForDiagnosis struct {
	JobInformationList []*entity.JobInformationForDiagnosis `json:"job_information_list"`
}

func NewJobInformationListForDiagnosis(jobInformations []*entity.JobInformationForDiagnosis) JobInformationListForDiagnosis {
	return JobInformationListForDiagnosis{
		JobInformationList: jobInformations,
	}
}
