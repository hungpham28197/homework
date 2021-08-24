package model

type User struct {
	Fullname string
	Birthday string
	Sex      string
	Job      string
	Email    string
	Pass     string
	Roles    Roles
}

type AuthenInfo struct {
	Fullname string
	Email    string
	Roles    Roles //kiểu map[int]bool
}

type Roles map[int]interface{}
