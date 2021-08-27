package repo

import (
	"errors"
	"fmt"

	"github.com/TechMaster/core/pass"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
)

var users = make(map[string]*pmodel.User)

func Init() {
	CreatNewUser("Trịnh Đình Long", "1", "long@gmail.com", "093023920", rbac.ADMIN)
}

func CreatNewUser(name string, password string, email string, phone string, roles ...int) {
	hashedpass, _ := pass.HashBcryptPass(password)
	users[email] = &pmodel.User{
		FullName: name,
		Password: hashedpass,
		Email:    email,
		Phone:    phone,
		Roles:    roles,
	}
}

func QueryByEmail(email string) (user *pmodel.User, err error) {
	u := users[email]
	if u == nil {
		return nil, errors.New("User not found")
	}
	return u, nil
}

func QueryById(email string) (user *pmodel.User, err error) {
	return nil, errors.New("User not found")
}

func QueryByName(name string) (user *pmodel.User, err error) {
	for _, v := range users {
		if v.FullName == name {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

func CreateUser(user *pmodel.User) (err error) {
	userTmp, err := QueryByEmail(user.Email)
	if err == nil {
		msg := fmt.Sprintf("User %s exit", userTmp.Email)
		return errors.New(msg)
	}
	hashedpass, _ := pass.HashBcryptPass(user.Password)
	user.Password = hashedpass
	users[user.Email] = user
	return nil
}

func UpdateUser(user *pmodel.User) (err error) {
	userTmp, err := QueryByEmail(user.Email)
	if err != nil {
		return err
	}
	oldPass := userTmp.Password
	*userTmp = *user
	userTmp.Password = oldPass
	return nil
}

func DeleteUser(id string) (err error) {
	userTmp, err := QueryById(id)
	if err != nil {
		return err
	}
	delete(users, userTmp.Email)
	return nil
}

func GetAllUser() (rs []pmodel.User) {
	for _, v := range users {
		rs = append(rs, *v)
	}
	return
}

func GetAllFullname() (fullnames []string) {
	for _, v := range users {
		fullnames = append(fullnames, v.FullName)
	}
	return
}

func UpdateAvatarUser(id string, avatars string) (err error) {
	userTmp, err := QueryById(id)
	if err != nil {
		return err
	}
	userTmp.Avatar = avatars
	return nil
}
