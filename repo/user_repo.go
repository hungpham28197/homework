package repo

import (
	"errors"
	"fmt"
	"home-work/model"

	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
)

var users = []model.User{
	{
		Fullname: "Bob",
		Birthday: "01/06/2000",
		Sex:      "Male",
		Job:      "Dev",
		User: pmodel.User{
			Email: "bob@gmail.com",
			Pass:  "1",
			Roles: pmodel.Roles{rbac.ADMIN: true},
		},
	},

	{
		Fullname: "Pham Van Long",
		Birthday: "01/06/1996",
		Sex:      "Male",
		Job:      "Dev",
		User: pmodel.User{
			Email: "long@gmail.com",
			Pass:  "1",
			Roles: pmodel.Roles{rbac.AUTHOR: true},
		},
	},
	{
		Fullname: "Pham Van Linh",
		Birthday: "11/05/2001",
		Sex:      "Male",
		Job:      "Dev",
		User: pmodel.User{
			Email: "linh@gmail.com",
			Pass:  "1",
			Roles: pmodel.Roles{rbac.EDITOR: false},
		},
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

func QueryByName(name string) (user *model.User, err error) {
	for _, v := range users {
		if v.Fullname == name {
			return &v, nil
		}
	}
	return nil, errors.New("User not found")
}

func CreateUser(user *model.User) (err error) {
	userTmp, err := QueryByEmail(user.Email)
	if err == nil {
		msg := fmt.Sprintf("User %s exit", userTmp.Email)
		return errors.New(msg)
	}
	users = append(users, *user)
	return nil
}

func UpdateUser(user *model.User) (err error) {
	for i := 0; i < len(users); i++ {
		if users[i].Email == user.Email {
			p := users[i]
			p = *user
			users[i] = p
			return nil
		}
	}
	return errors.New("User not found")
}

func DeleteUser(email string) (err error) {
	index := -1
	for i := 0; i < len(users); i++ {
		if users[i].Email == email {
			index = i
		}
	}
	if index > 0 {
		users = append(users[:index], users[index+1:]...)
		return nil
	}
	return errors.New("User not found")
}

func GetAllUser() (rs []model.User) {
	for _, v := range users {
		rs = append(rs, v)
	}
	return
}

func GetAllFullname() (fullnames []string) {
	for _, v := range users {
		fullnames = append(fullnames, v.Fullname)
	}
	return
}

func UpdateAvatarUser(emailUpload string, avatars string) (err error) {
	for i := 0; i < len(users); i++ {
		if users[i].Email == emailUpload {
			p := users[i]
			p.Avatar = avatars
			users[i] = p
			return nil
		}
	}
	return errors.New("User not found")
}
