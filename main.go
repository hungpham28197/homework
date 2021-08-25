package main

import (
	"home-work/rbac"
	"home-work/router"
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
	router.RegisterRoute(app)
	template.InitViewEngine(app)

	//Luôn để hàm này sau tất cả lệnh cấu hình đường dẫn với RBAC
	rbac.BuildPublicRoute(app)

	app.Listen(":8080")
}
