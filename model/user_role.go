package model

import "github.com/TechMaster/core/pmodel"

type User struct {
	Fullname string
	Birthday string
	Sex      string
	Job      string
	Avatar   string
	pmodel.User
}
