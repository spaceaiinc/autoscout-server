package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type SendingPhase struct {
	SendingPhase *entity.SendingPhase `json:"sending_phase"`
}

func NewSendingPhase(sendingPhase *entity.SendingPhase) SendingPhase {
	return SendingPhase{
		SendingPhase: sendingPhase,
	}
}

type SendingPhaseList struct {
	SendingPhaseList []*entity.SendingPhase `json:"sending_phase_list"`
}

func NewSendingPhaseList(sendingPhases []*entity.SendingPhase) SendingPhaseList {
	return SendingPhaseList{
		SendingPhaseList: sendingPhases,
	}
}

type SendingJobSeekerTableListAndMaxPageAndIDListAndListCount struct {
	SendingJobSeekerTableList    []*entity.SendingJobSeekerTable `json:"sending_job_seeker_table_list"`
	MaxPageNumber                uint                            `json:"max_page_number"`
	CompleteSendingCount         uint                            `json:"complete_sending_count"`
	AcceptSendingCount           uint                            `json:"accept_sending_count"`            // 送客応諾数
	PreparingAfterInterviewCount uint                            `json:"preparing_after_interview_count"` // 面談実施済み数
	WaitingForInterviewCount     uint                            `json:"waiting_for_interview_count"`     // 面談実施待ち数
	UnregisterDetailCount        uint                            `json:"unregister_detail_count"`         // 詳細未登録数
	UnregisterScheduleCount      uint                            `json:"unregister_schedule_count"`       // 日程未登録数
	CloseSendingCount            uint                            `json:"close_sending_count"`             // 終了数
}

func NewSendingJobSeekerTableListAndMaxPageAndIDListAndListCount(
	sendingJobSeekerTables []*entity.SendingJobSeekerTable,
	maxPageNumber uint,
	completeSendingCount uint,
	acceptSendingCount uint,
	preparingAfterInterviewCount uint,
	waitingForInterviewCount uint,
	unregisterDetailCount uint,
	unregisterScheduleCount uint,
	closeSendingCount uint,
) SendingJobSeekerTableListAndMaxPageAndIDListAndListCount {
	return SendingJobSeekerTableListAndMaxPageAndIDListAndListCount{
		MaxPageNumber:                maxPageNumber,
		SendingJobSeekerTableList:    sendingJobSeekerTables,
		CompleteSendingCount:         completeSendingCount,
		AcceptSendingCount:           acceptSendingCount,
		PreparingAfterInterviewCount: preparingAfterInterviewCount,
		WaitingForInterviewCount:     waitingForInterviewCount,
		UnregisterDetailCount:        unregisterDetailCount,
		UnregisterScheduleCount:      unregisterScheduleCount,
		CloseSendingCount:            closeSendingCount,
	}
}
