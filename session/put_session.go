package session

import (
	"home-work/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

func SetAuthenticated(ctx iris.Context, authUser model.AuthenInfo) {
	sess := sessions.Get(ctx)

	sess.Set(SESS_AUTH, true)
	sess.Set(SESS_USER, authUser)
}
