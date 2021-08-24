package session

import (
	"home-work/model"

	"github.com/TechMaster/eris"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/mitchellh/mapstructure"
)

func GetAuthInfo(ctx iris.Context) (*model.AuthenInfo, error) {
	data := sessions.Get(ctx).Get(SESS_USER)
	if data == nil {
		return nil, nil
	}
	authinfo := new(model.AuthenInfo)
	if err := mapstructure.WeakDecode(data, authinfo); err != nil {
		return nil, eris.NewFrom(err)
	}
	return authinfo, nil
}

func IsLogin(ctx iris.Context) bool {
	isLogin, _ := sessions.Get(ctx).GetBoolean(SESS_AUTH)
	return isLogin
}
