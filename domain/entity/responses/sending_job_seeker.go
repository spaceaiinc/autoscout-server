package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SendingJobSeeker struct {
	SendingJobSeeker *entity.SendingJobSeeker `json:"sending_job_seeker"`
}

func NewSendingJobSeeker(sendingJobSeeker *entity.SendingJobSeeker) SendingJobSeeker {
	return SendingJobSeeker{
		SendingJobSeeker: sendingJobSeeker,
	}
}

type SendingJobSeekerList struct {
	SendingJobSeekerList []*entity.SendingJobSeeker `json:"sending_job_seeker_list"`
}

func NewSendingJobSeekerList(sendingJobSeekers []*entity.SendingJobSeeker) SendingJobSeekerList {
	return SendingJobSeekerList{
		SendingJobSeekerList: sendingJobSeekers,
	}
}

type SendingJobSeekerListAndMaxPageAndIDList struct {
	MaxPageNumber        uint                       `json:"max_page_number"`
	IDList               []uint                     `json:"id_list"`
	SendingJobSeekerList []*entity.SendingJobSeeker `json:"sending_job_seeker_list"`
}

func NewSendingJobSeekerListAndMaxPageAndIDList(sendingJobSeekers []*entity.SendingJobSeeker, maxPageNumber uint, idList []uint) SendingJobSeekerListAndMaxPageAndIDList {
	return SendingJobSeekerListAndMaxPageAndIDList{
		MaxPageNumber:        maxPageNumber,
		IDList:               idList,
		SendingJobSeekerList: sendingJobSeekers,
	}
}

type SendingJobSeekerDocument struct {
	SendingJobSeekerDocument *entity.SendingJobSeekerDocument `json:"sending_job_seeker_document"`
}

func NewSendingJobSeekerDocument(sendingJobSeekerDocument *entity.SendingJobSeekerDocument) SendingJobSeekerDocument {
	return SendingJobSeekerDocument{
		SendingJobSeekerDocument: sendingJobSeekerDocument,
	}
}

type IsNotViewCount struct {
	IsViewCount *entity.IsViewCount `json:"is_view_count"`
}

func NewIsNotViewCount(isViewCount *entity.IsViewCount) IsNotViewCount {
	return IsNotViewCount{
		IsViewCount: isViewCount,
	}
}
