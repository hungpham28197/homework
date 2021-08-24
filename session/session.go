package session

import (
	"time"

	"github.com/kataras/iris/v12/sessions"
)

const (
	SESSION_COOKIE = "sessid"
	SESS_AUTH      = "authenticated"
	SESS_USER      = "user"
	AUTHINFO       = "authinfo"
)

/*
Cấu hình Session Manager
*/
func Init() *sessions.Sessions {
	return sessions.New(sessions.Config{
		Cookie:       SESSION_COOKIE,
		AllowReclaim: true,
		Expires:      time.Hour * 48, /*Có giá trị trong 2 ngày*/
	})
}
