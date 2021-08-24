package repo

import (
	"errors"
	"home-work/model"
	"home-work/rbac"
)

var users = []model.User{
	{
		Fullname: "Pham Van Hung",
		Birthday: "17/06/1996",
		Sex:      "Male",
		Job:      "Dev",
		Email:    "hung@gmail.com",
		Pass:     "1",
		Roles:    model.Roles{rbac.STUDENT: true},
	},
}

func QueryByEmail(email string) (user *model.User, err error) {
	for _, v := range users {
		if v.Email == email {
			return &v, nil
		}
	}
	return nil, errors.New("User not found")
}
