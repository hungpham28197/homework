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
	if user.Password != loginRequest.Pass {
		_, _ = ctx.WriteString("Wrong password")
		return
	}

	session.SetAuthenticated(ctx, pmodel.AuthenInfo{
		Email:    user.Email,
		FullName: user.FullName,
		Roles:    pmodel.IntArrToRoles(user.Roles),
	})

	ctx.Redirect("/")
}

func ShowHomePage(ctx iris.Context) {
	if raw_authinfo := ctx.GetViewData()[session.AUTHINFO]; raw_authinfo != nil {
		authinfo := raw_authinfo.(*pmodel.AuthenInfo)
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
	ctx.ViewData("fullnames", []string{user.FullName})
	_ = ctx.View("public")
}

func ShowPublicPage(ctx iris.Context) {
	ctx.ViewData("fullnames", repo.GetAllFullname())
	addRoleToForm(ctx)
	_ = ctx.View("public")
}

func addRoleToForm(ctx iris.Context) {
	if raw_authinfo := ctx.GetViewData()[session.AUTHINFO]; raw_authinfo != nil {
		authinfo := raw_authinfo.(*pmodel.AuthenInfo)
		ctx.ViewData("roles", rbac.RolesNames(authinfo.Roles))
	}
}

func PrivateUsers(ctx iris.Context) {
	addRoleToForm(ctx)
	users := repo.GetAllUser()
	ctx.ViewData("users", users)
	_ = ctx.View("private")
}

func PrivateCreateUsers(ctx iris.Context) {
	var createUserRequest pmodel.User

	if err := ctx.ReadJSON(&createUserRequest); err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}

	err := repo.CreateUser(&createUserRequest)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err))
	}

	// _, _, errr := ctx.UploadFormFiles("./uploads")

	// if errr != nil {
	// 	logger.Log(ctx, eris.NewFrom(err).BadRequest())
	// 	return
	// }
	_ = ctx.View("private")
}

func PrivateUpdateUsers(ctx iris.Context) {
	var updateUserRequest pmodel.User

	if err := ctx.ReadJSON(&updateUserRequest); err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}
	fmt.Println(&updateUserRequest)
	err := repo.UpdateUser(&updateUserRequest)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}
	_ = ctx.View("private")
}

func PrivateDeleteUsers(ctx iris.Context) {
	var userRequest pmodel.User

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

	fmt.Println("in: " + loginReq.Pass + " out: " + user.Password)
	if ok := pass.CheckPassword(loginReq.Pass, user.Password, "OjO3km9xQFQEqUVI/xGXviRuYhNOnDG"); ok {
		logger.Log(ctx, eris.Warning("Wrong password").UnAuthorized())
		return
	}

	session.SetAuthenticated(ctx, pmodel.AuthenInfo{
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

func UploadPhoto(ctx iris.Context) {
	uploadEmail := ctx.FormValue("UploadEmail")

	user, err := repo.QueryByEmail(uploadEmail)

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err))
	}

	uploadfiles, _, err := ctx.UploadFormFiles("./uploads", func(ctx iris.Context, file *multipart.FileHeader) bool {
		//id := ctx.FormValue("Id")
		extension := filepath.Ext(file.Filename)
		fmt.Println(extension)
		file.Filename = user.Id + extension
		return true
	})
	if err != nil {
		logger.Log(ctx, eris.NewFrom(err))
	}

	var avatars []string
	filenames := "Upload successful \n"
	for _, upload := range uploadfiles {
		filenames = filenames + upload.Filename + "/n"
		avatars = append(avatars, upload.Filename)
	}
	err = repo.UpdateAvatarUser(user.Id, avatars[0])

	if err != nil {
		logger.Log(ctx, eris.NewFrom(err))
	}

	_, _ = ctx.WriteString(filenames)
}

// func beforeSave(ctx iris.Context, file *multipart.FileHeader) bool {
// 	id := ctx.FormValue("Id")
// 	extension := filepath.Ext(file.Filename)
// 	fmt.Println(extension)
// 	file.Filename = id + extension
// 	return true
// }
