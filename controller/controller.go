package controller

import (
	"fmt"
	"home-work/model"
	"home-work/rbac"
	"home-work/repo"
	"home-work/session"

	"github.com/TechMaster/eris"
	"github.com/TechMaster/logger"
	"github.com/kataras/iris/v12"
)

type LoginRequest struct {
	Email string
	Pass  string
}

/*
Login thông qua form tại giao diện
*/
func Login(ctx iris.Context) {
	var loginRequest LoginRequest
	if err := ctx.ReadForm(&loginRequest); err != nil {
		fmt.Printf(err.Error())
		return
	}

	user, err := repo.QueryByEmail(loginRequest.Email)
	if err != nil {
		_, _ = ctx.WriteString("Login Failed")
		return
	}
	if user.Pass != loginRequest.Pass {
		_, _ = ctx.WriteString("Wrong password")
		return
	}

	session.SetAuthenticated(ctx, model.AuthenInfo{
		Email:    user.Email,
		Fullname: user.Fullname,
		Roles:    user.Roles,
	})

	ctx.Redirect("/")
}

func ShowHomePage(ctx iris.Context) {
	if raw_authinfo := ctx.GetViewData()[session.AUTHINFO]; raw_authinfo != nil {
		authinfo := raw_authinfo.(*model.AuthenInfo)
		fmt.Println("vao day")
		ctx.ViewData("roles", rbac.RolesNames(authinfo.Roles))
	}
	_ = ctx.View("index")
}

func LoginJSON(ctx iris.Context) {
	var loginReq LoginRequest

	if err := ctx.ReadJSON(&loginReq); err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}

	user, err := repo.QueryByEmail(loginReq.Email)

	if err != nil { //Không tìm thấy user
		logger.Log(ctx, eris.Warning("User not found").UnAuthorized())
		return
	}

	if user.Pass != loginReq.Pass {
		logger.Log(ctx, eris.Warning("Wrong password").UnAuthorized())
		return
	}

	session.SetAuthenticated(ctx, model.AuthenInfo{
		Fullname: user.Fullname,
		Email:    user.Email,
	})
	fmt.Println("ok")
	_, _ = ctx.JSON("Login successfully")
}
