package user

import "time"

type User struct {
	ID             int
	Name           string
	Occupation     string
	Email          string
	PasswordHash   string
	AvatarFileName string
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserOnWeb struct {
	ID             int
	Number         int
	Name           string
	Occupation     string
	Email          string
	Link           string
	AvatarFileName string
}
