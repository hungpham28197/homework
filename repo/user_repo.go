package repo

import (
	"errors"
	"fmt"
	"os"

	"github.com/TechMaster/core/pass"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/google/uuid"
)

var users = make(map[string]*pmodel.User)

func Init() {
	CreatNewUser("Trịnh Đình Long", "1", "long@gmail.com", "093023920", "2.png", rbac.ADMIN)
	CreatNewUser("Trịnh Đình Linh", "1", "linh@gmail.com", "093023920", "3.jpeg", rbac.ADMIN)
}

func CreatNewUser(name string, password string, email string, phone string, avatar string, roles ...int) {
	hashedpass, _ := pass.HashBcryptPass(password)
	users[email] = &pmodel.User{
		Id:       uuid.NewString(),
		FullName: name,
		Password: hashedpass,
		Email:    email,
		Phone:    phone,
		Avatar:   avatar,
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

func QueryById(id string) (user *pmodel.User, err error) {
	for _, v := range users {
		if v.Id == id {
			return v, nil
		}
	}
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

func CreateUser(user *pmodel.User) (rs *pmodel.User, err error) {
	userTmp, err := QueryByEmail(user.Email)
	if err == nil {
		msg := fmt.Sprintf("User %s exit", userTmp.Email)
		return nil, errors.New(msg)
	}
	hashedpass, _ := pass.HashBcryptPass(user.Password)
	user.Password = hashedpass
	user.Id = uuid.NewString()
	users[user.Email] = user
	return user, nil
}

func UpdateUser(user *pmodel.User) (err error) {
	userTmp, err := QueryByEmail(user.Email)
	if err != nil {
		return err
	}
	oldId := userTmp.Id
	oldPass := userTmp.Password
	oldAvatar := userTmp.Avatar

	*userTmp = *user

	userTmp.Password = oldPass
	userTmp.Id = oldId
	userTmp.Avatar = oldAvatar
	return nil
}

func DeleteUser(id string) (err error) {
	userTmp, err := QueryById(id)
	if err != nil {
		return err
	}
	os.Remove("./uploads/" + userTmp.Avatar)
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

func UpdateAvatarUserByEmail(email string, avatars string) (err error) {
	userTmp, err := QueryByEmail(email)
	if err != nil {
		return err
	}
	userTmp.Avatar = avatars
	return nil
}

func UpdateAvatarUser(id string, avatars string) (err error) {
	userTmp, err := QueryById(id)
	if err != nil {
		return err
	}
	userTmp.Avatar = avatars
	return nil
}
