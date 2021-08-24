package main

import (
	"home-work/controller"
	"home-work/rbac"
	"home-work/session"
	"home-work/template"

	"github.com/TechMaster/logger"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	logFile := logger.Init()
	if logFile != nil {
		defer logFile.Close()
	}
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)

	app.Use(session.Init().Handler())

	rbacConfig := rbac.NewConfig()
	rbacConfig.RootAllow = true
	rbac.Init(rbacConfig) //Khởi động với cấu hình mặc định

	//đặt hàm này trên các hàm đăng ký route - controller
	app.Use(rbac.CheckRoutePermission)

	app.Get("/", controller.ShowHomePage)

	app.Get("/public", controller.ShowPublicPage)

	app.Post("/public/name", controller.PublicSearchName)

	rbac.Get(app, "private", rbac.Allow(11, 12, 7), controller.PrivateUsers)

	rbac.Post(app, "private/create", rbac.Allow(7, 12), controller.PrivateCreateUsers)

	rbac.Post(app, "private/edit", rbac.Allow(7, 12), controller.PrivateUpdateUsers)

	rbac.Post(app, "private/delete", rbac.Allow(7), controller.PrivateDeleteUsers)

	rbac.Post(app, "/login", rbac.AllowAll(), controller.Login)

	app.Post("/loginjson", controller.LoginJSON)
	rbac.Get(app, "/logoutjson", rbac.AllowAll(), controller.LogoutFromREST)

	template.InitViewEngine(app)

	//Luôn để hàm này sau tất cả lệnh cấu hình đường dẫn với RBAC
	rbac.BuildPublicRoute(app)

	app.Listen(":8080")
}
