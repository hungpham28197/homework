package repo

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/TechMaster/core/pass"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/google/uuid"
)

var users = make(map[string]*pmodel.User)

func init() {
	CreateNewUser("Bùi Văn Hiên", "1", "hien@gmail.com", "0123456789", rbac.TRAINER, rbac.MAINTAINER)
	CreateNewUser("Nguyễn Hàn Duy", "1", "duy@gmail.com", "0123456786", rbac.TRAINER, rbac.STUDENT)
	CreateNewUser("Phạm Thị Mẫn", "1", "man@gmail.com", "0123456780", rbac.SALE, rbac.STUDENT)
	CreateNewUser("Trịnh Minh Cường", "1", "cuong@gmail.com", "0123456000", rbac.ADMIN, rbac.TRAINER)
	CreateNewUser("Nguyễn Thành Long", "1", "long@gmail.com", "0123456001", rbac.STUDENT)
}

func CreateNewUser(name string, password string, email string, phone string, roles ...int) {
	hashpass, _ := pass.HashBcryptPass(password)
	users[email] = &pmodel.User{
		Id:       uuid.NewString(),
		FullName: name,
		Password: hashpass,
		Email:    email,
		Phone:    phone,
		Roles:    roles,
	}
}

func QueryByEmail(email string) (user *pmodel.User, err error) {
	user = users[strings.ToLower(email)]
	if user == nil {
		return nil, errors.New("User not found")
	}
	return user, nil
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

func CreateUser(user *pmodel.User) (err error) {
	userTmp, err := QueryByEmail(user.Email)
	if err == nil {
		msg := fmt.Sprintf("User %s exit", userTmp.Email)
		return errors.New(msg)
	}
	user.Id = uuid.NewString()
	hashpass, _ := pass.HashBcryptPass(user.Password)
	user.Password = hashpass
	users[user.Email] = user
	return nil
}

func UpdateUser(user *pmodel.User) (err error) {
	userTmp, err := QueryByEmail(user.Email)
	if err != nil {
		return err
	}
	id := userTmp.Id
	*userTmp = *user
	userTmp.Id = id
	return nil
}

func DeleteUser(email string) (err error) {
	user, err := QueryByEmail(email)
	if err != nil {
		return err
	}
	e := os.Remove("./uploads/" + user.Avatar)
	if e != nil {
		return e
	}
	delete(users, user.Email)
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
	user, err := QueryById(id)
	if err != nil {
		return err
	}
	user.Avatar = avatars
	return nil
}
