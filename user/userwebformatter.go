package user

import (
	"fmt"
)

type UserWebFormatter struct {
	Number         int    `json:"no"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Occupation     string `json:"occupation"`
	Link           string `json:"link"`
	AvatarFileName string `json:"avatar_filename"`
}

func FormatWebUser(user UserOnWeb) UserWebFormatter {

	userFormatter := UserWebFormatter{}
	userFormatter.Number = user.Number + 1
	userFormatter.AvatarFileName = fmt.Sprintf("<img src='%s' width='90' class='img-fluid img-circle'>", user.AvatarFileName)
	userFormatter.Name = user.Name
	userFormatter.Email = user.Email
	userFormatter.Occupation = user.Occupation
	userFormatter.Link = fmt.Sprintf("<a href='/users/avatar/%d' class='btn btn-light'><i class='fa fa-camera'></i></a> <a href='/users/edit/%d' class='btn btn-warning'><i class='fa fa-pencil'></i></a> <a href='/users/delete/%d' class='btn btn-danger' onclick='javascript: return confirm(\"Are You Sure ?\")''><i class='fa fa-trash'></i></a>", user.ID, user.ID, user.ID)

	return userFormatter
}

func FormatWebUsers(users []UserOnWeb) []UserWebFormatter {

	userFormatters := []UserWebFormatter{}
	for _, user := range users {
		userFormatter := FormatWebUser(user)
		userFormatters = append(userFormatters, userFormatter)
	}

	return userFormatters
}
