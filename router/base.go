package router

import (
	"home-work/controller"
	"home-work/rbac"

	"github.com/kataras/iris/v12"
)

func RegisterRoute(app *iris.Application) {
	app.Get("/", controller.ShowPublicPage)
	app.Post("/loginjson", controller.LoginJSON)
	rbac.Get(app, "/logoutjson", rbac.AllowAll(), controller.LogoutFromREST)

	public := app.Party("/public")
	public.Get("/", controller.ShowPublicPage)
	public.Post("/name", controller.PublicSearchName)

	private := app.Party("private")
	rbac.Get(private, "", rbac.Allow(rbac.GUEST, rbac.STAFF, rbac.EDITOR), controller.PrivateUsers)
	rbac.Post(private, "/create", rbac.Allow(rbac.EDITOR, rbac.STAFF), controller.PrivateCreateUsers)
	rbac.Post(private, "/edit", rbac.Allow(rbac.EDITOR, rbac.STAFF), controller.PrivateUpdateUsers)
	rbac.Post(private, "/delete", rbac.Allow(rbac.EDITOR), controller.PrivateDeleteUsers)
	rbac.Post(private, "/login", rbac.AllowAll(), controller.Login)
	rbac.Post(private, "/upload", rbac.AllowAll(), iris.LimitRequestBodySize(300000), controller.UploadPhoto)
}
