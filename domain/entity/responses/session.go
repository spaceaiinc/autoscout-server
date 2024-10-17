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

type GuestEnterpriseUserSession struct {
	GuestEnterpriseUser *entity.GuestEnterpriseUser `json:"user"`
}

func NewGuestEnterpriseUserSession(guestEnterprise *entity.GuestEnterpriseUser) GuestEnterpriseUserSession {
	return GuestEnterpriseUserSession{
		GuestEnterpriseUser: guestEnterprise,
	}
}

type GuestJobSeekerUserSession struct {
	GuestJobSeekerUser *entity.GuestJobSeekerUser `json:"user"`
}

func NewGuestJobSeekerUserSession(guestEnterprise *entity.GuestJobSeekerUser) GuestJobSeekerUserSession {
	return GuestJobSeekerUserSession{
		GuestJobSeekerUser: guestEnterprise,
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
