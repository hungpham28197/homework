package controller

import (
	"home-work/repo"

	"github.com/TechMaster/eris"
	"github.com/TechMaster/logger"
	"github.com/kataras/iris/v12"
)

func ShowChangeAvatar(ctx iris.Context) {
	GetViewUserById(ctx, "change-avatar")
}

func ChangeAvatar(ctx iris.Context) {
	id := ctx.Params().Get("id")
	user, err := repo.QueryById(id)
	if err != nil { //Không tìm thấy user
		logger.Log(ctx, eris.NewFrom(err))
	}
	oldAvatar := user.Avatar
	UploadAvatar(user, ctx)
	if oldAvatar != user.Avatar {
		RemoveAvatarFile(oldAvatar)
	}
	ctx.Redirect("/private")
}
