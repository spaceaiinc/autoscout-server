package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type EmailWithJobSeeker struct {
	EmailWithJobSeeker *entity.EmailWithJobSeeker `json:"email_with_job_seeker"`
}

func NewEmailWithJobSeeker(emailWithJobSeeker *entity.EmailWithJobSeeker) EmailWithJobSeeker {
	return EmailWithJobSeeker{
		EmailWithJobSeeker: emailWithJobSeeker,
	}
}

type EmailWithJobSeekerList struct {
	EmailWithJobSeekerList []*entity.EmailWithJobSeeker `json:"email_with_job_seeker_list"`
}

func NewEmailWithJobSeekerList(emailWithJobSeekers []*entity.EmailWithJobSeeker) EmailWithJobSeekerList {
	return EmailWithJobSeekerList{
		EmailWithJobSeekerList: emailWithJobSeekers,
	}
}
