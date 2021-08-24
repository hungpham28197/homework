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

type SearchUserRequest struct {
	Name string
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
		ctx.ViewData("roles", rbac.RolesNames(authinfo.Roles))
	}
	_ = ctx.View("public")
}

func PublicSearchName(ctx iris.Context) {
	name := ctx.PostValue("Name")

	user, err := repo.QueryByName(name)

	if err != nil { //Không tìm thấy user
		fmt.Println(err)
		return
	}
	ctx.ViewData("fullnames", []string{user.Fullname})
	_ = ctx.View("public")
}

func ShowPublicPage(ctx iris.Context) {
	ctx.ViewData("fullnames", repo.GetAllFullname())
	_ = ctx.View("public")
}

func PrivateUsers(ctx iris.Context) {
	users := repo.GetAllUser()
	ctx.ViewData("users", users)
	_ = ctx.View("private")
}

func PrivateCreateUsers(ctx iris.Context) {
	var createUserRequest model.User

	if err := ctx.ReadJSON(&createUserRequest); err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}

	err := repo.CreateUser(&createUserRequest)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}
	_ = ctx.View("private")
}

func PrivateUpdateUsers(ctx iris.Context) {
	var createUserRequest model.User

	if err := ctx.ReadJSON(&createUserRequest); err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}

	err := repo.UpdateUser(&createUserRequest)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}
	_ = ctx.View("private")
}

func PrivateDeleteUsers(ctx iris.Context) {
	var userRequest model.User

	if err := ctx.ReadJSON(&userRequest); err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}

	err := repo.DeleteUser(userRequest.Email)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}
	_ = ctx.View("private")
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
		Roles:    user.Roles,
	})
	_, _ = ctx.JSON("Login successfully")
}

func LogoutFromREST(ctx iris.Context) {
	logout(ctx)
	_, _ = ctx.JSON("Logout success")
}

func logout(ctx iris.Context) {
	/*	if !session.IsLogin(ctx) {
		logger.Log(ctx, eris.Warning("Bạn chưa login").UnAuthorized())
	}*/
	//Xoá toàn bộ session và xoá luôn cả Cookie sessionid ở máy người dùng
	session.Destroy(ctx)
}
