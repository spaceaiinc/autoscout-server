package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type EmailWithSendingJobSeeker struct {
	EmailWithSendingJobSeeker *entity.EmailWithSendingJobSeeker `json:"email_with_sending_job_seeker"`
}

func NewEmailWithSendingJobSeeker(emailWithSendingJobSeeker *entity.EmailWithSendingJobSeeker) EmailWithSendingJobSeeker {
	return EmailWithSendingJobSeeker{
		EmailWithSendingJobSeeker: emailWithSendingJobSeeker,
	}
}

type EmailWithSendingJobSeekerList struct {
	EmailWithSendingJobSeekerList []*entity.EmailWithSendingJobSeeker `json:"email_with_sending_job_seeker_list"`
}

func NewEmailWithSendingJobSeekerList(emailWithSendingJobSeekers []*entity.EmailWithSendingJobSeeker) EmailWithSendingJobSeekerList {
	return EmailWithSendingJobSeekerList{
		EmailWithSendingJobSeekerList: emailWithSendingJobSeekers,
	}
}
