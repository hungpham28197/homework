package router

import (
	"home-work/controller"

	"github.com/TechMaster/core/rbac"
	"github.com/kataras/iris/v12"
)

func RegisterRoute(app *iris.Application) {
	app.HandleDir("/", iris.Dir("./uploads"))

	app.Get("/", controller.ShowPublicPage)
	app.Post("/loginjson", controller.LoginJSON)
	rbac.Get(app, "/logoutjson", rbac.AllowAll(), controller.LogoutFromREST)

	public := app.Party("/public")
	public.Get("/", controller.ShowPublicPage)
	public.Post("/name", controller.SearchNameByName)

	private := app.Party("private")
	rbac.Get(private, "", rbac.Allow(rbac.ADMIN, rbac.AUTHOR, rbac.EDITOR), controller.GetUsers)
	rbac.Post(private, "/create", rbac.Allow(rbac.ADMIN, rbac.EDITOR, rbac.AUTHOR), iris.LimitRequestBodySize(300000), controller.CreateUser)
	rbac.Post(private, "/edit", rbac.Allow(rbac.ADMIN, rbac.EDITOR, rbac.AUTHOR), controller.UpdateUser)
	rbac.Get(private, "/edit/{id}", rbac.Allow(rbac.ADMIN, rbac.AUTHOR, rbac.EDITOR), controller.ShowUpdateUser)
	rbac.Get(private, "/change-avatar/{id}", rbac.Allow(rbac.ADMIN, rbac.AUTHOR, rbac.EDITOR), controller.ShowChangeAvatar)
	rbac.Post(private, "/change-avatar/{id}", rbac.Allow(rbac.ADMIN, rbac.EDITOR, rbac.AUTHOR), iris.LimitRequestBodySize(300000), controller.ChangeAvatar)
	rbac.Post(private, "/delete/{id}", rbac.Allow(rbac.ADMIN, rbac.EDITOR), controller.DeleteUser)
}
