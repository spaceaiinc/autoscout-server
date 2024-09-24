package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type UserSession struct {
	User *entity.User `json:"user"`
}

func NewUserSession(user *entity.User) UserSession {
	return UserSession{
		User: user,
	}
}

type GestEnterpriseUserSession struct {
	GestEnterpriseUser *entity.GestEnterpriseUser `json:"user"`
}

func NewGestEnterpriseUserSession(gestEnterprise *entity.GestEnterpriseUser) GestEnterpriseUserSession {
	return GestEnterpriseUserSession{
		GestEnterpriseUser: gestEnterprise,
	}
}

type GestJobSeekerUserSession struct {
	GestJobSeekerUser *entity.GestJobSeekerUser `json:"user"`
}

func NewGestJobSeekerUserSession(gestEnterprise *entity.GestJobSeekerUser) GestJobSeekerUserSession {
	return GestJobSeekerUserSession{
		GestJobSeekerUser: gestEnterprise,
	}
}

type GoogleAuthSession struct {
	AuthURL string `json:"auth_url"`
}

func NewGoogleAuthSession(authURL string) GoogleAuthSession {
	return GoogleAuthSession{
		AuthURL: authURL,
	}
}
