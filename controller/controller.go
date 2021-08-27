package controller

import (
	"fmt"
	"home-work/repo"
	"mime/multipart"
	"path/filepath"

	"github.com/TechMaster/core/pass"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/TechMaster/core/session"
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

func SearchNameByName(ctx iris.Context) {
	name := ctx.PostValue("Name")

	user, err := repo.QueryByName(name)

	if err != nil { //Không tìm thấy user
		fmt.Println(err)
		return
	}
	ctx.ViewData("fullnames", []string{user.FullName})
	_ = ctx.View("public")
}

func ShowPublicPage(ctx iris.Context) {
	if raw_authinfo := ctx.GetViewData()[session.AUTHINFO]; raw_authinfo != nil {
		authinfo := raw_authinfo.(*pmodel.AuthenInfo)
		ctx.ViewData("roles", rbac.RolesNames(authinfo.Roles))
	}
	ctx.ViewData("fullnames", repo.GetAllFullname())
	_ = ctx.View("public")
}

func GetUsers(ctx iris.Context) {
	if raw_authinfo := ctx.GetViewData()[session.AUTHINFO]; raw_authinfo != nil {
		authinfo := raw_authinfo.(*pmodel.AuthenInfo)
		ctx.ViewData("roles", rbac.RolesNames(authinfo.Roles))
	}
	users := repo.GetAllUser()
	ctx.ViewData("users", users)
	_ = ctx.View("private")
}

func UpdateUser(ctx iris.Context) {
	var createUserRequest pmodel.User

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

func DeleteUser(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		logger.Log(ctx, eris.Warning("Id is empty").BadRequest())
		return
	}

	err := repo.DeleteUser(id)

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

	if !pass.CheckPassword(loginReq.Pass, user.Password, "") {
		logger.Log(ctx, eris.Warning("Wrong password").UnAuthorized())
		return
	}

	session.SetAuthenticated(ctx, pmodel.AuthenInfo{
		Id:       user.Id,
		FullName: user.FullName,
		Email:    user.Email,
		Roles:    pmodel.IntArrToRoles(user.Roles),
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
	session.Sess.Destroy(ctx)
}

func CreateUser(ctx iris.Context) {
	createEmail := ctx.FormValue("CreateEmail")
	createName := ctx.FormValue("CreateName")
	createPass := ctx.FormValue("CreatePass")

	userReq := &pmodel.User{
		Email:    createEmail,
		FullName: createName,
		Password: createPass,
	}

	rp, err := repo.CreateUser(userReq)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err))
		return
	}

	UploadAvatar(rp, ctx)

	ctx.Redirect("/private")
}

func UploadAvatar(user *pmodel.User, ctx iris.Context) {
	uploadfiles, _, err := ctx.UploadFormFiles("./uploads", func(c iris.Context, file *multipart.FileHeader) bool {
		file.Filename = user.Id + filepath.Ext(file.Filename)
		return true
	})

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err))
	}

	err = repo.UpdateAvatarUserByEmail(user.Email, uploadfiles[0].Filename)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err))
	}
}

func ShowUpdateUser(ctx iris.Context) {
	GetViewUserById(ctx, "edit-user")
}

func GetViewUserById(ctx iris.Context, fileHtml string) {
	id := ctx.Params().Get("id")
	user, err := repo.QueryById(id)

	if err != nil { //Không tìm thấy user
		logger.Log(ctx, eris.NewFrom(err))
	}

	ctx.ViewData("user", user)
	_ = ctx.View(fileHtml)
}
