package rbac

import (
	"home-work/model"
	"home-work/session"

	"github.com/TechMaster/eris"
	"github.com/TechMaster/logger"
	"github.com/kataras/iris/v12"
)

//Hàm này cần được gọi trước tất cả handler do lập trình viên viết
func CheckRoutePermission(ctx iris.Context) {
	route := ctx.GetCurrentRoute().String() //Lấy route trong ctx
	authinfo, _ := session.GetAuthInfo(ctx)

	//Nếu route không thuộc nhóm public routes cần kiểm tra phân quyền
	if authinfo != nil {
		//Gán authinfo để cho handler phía sau dùng
		ctx.ViewData(session.AUTHINFO, authinfo)
	}

	//Nếu route nằm trong public routes thì cho qua luôn
	if publicRoutes[route] {
		ctx.Next()
		return
	}

	//Chưa đăng nhập mà đòi vào protected route, vậy phải kick ra
	if authinfo == nil {
		logger.Log(ctx, eris.Warning("Bạn chưa đăng nhập").UnAuthorized())
		return
	}

	// Nếu RBAC cho phép Root permission bỏ qua logic check role
	// authinfo phải chứa role allow ROOT
	if config.RootAllow {
		if hasRoot := authinfo.Roles[ROOT]; hasRoot != nil && hasRoot.(bool) {
			ctx.Next()
			return
		}
	}

	allowGoNext := checkUser_RouteRole_Intersect(authinfo.Roles, routesRoles[route])

	if allowGoNext { //Được quyền đi tiếp
		ctx.Next()
	} else {
		logger.Log(ctx, eris.Warning("Bạn không quyền thực hiện tác vụ này").UnAuthorized())
	}
}

/*
rolesInRoute = routesRoles[route] //Lấy map roles gán cho route. Sẽ chỉ có 2 loại: AllowRoles và ForbidRoles rõ ràng

Nếu route không thuộc nhóm public routes cần kiểm tra phân quyền
	Duyệt qua tất cả quyền user có
	C1: userrole có trong rolesInRoute
			C1A: đó là allowed role
				ForbidOverAllow == false -> allowGoNext = true -> exit loop
				ForbidOverAllow == true  -> allowGoNext = true, chưa exit loop vội, có thể đến role sau là forbid

			C1B: đó là forbid role
				ForbidOverAllow == true  -> allowGoNext = false -> exit loop
				ForbidOverAllow == false -> allowGoNext = false, chưa exit loop vội, có thể
	C2: userrole không có trong rolesInRoute, mặc nhiên allowGoNext = false

*/
func checkUser_RouteRole_Intersect(userRoles model.Roles, rolesInRoute model.Roles) (allowGoNext bool) {

	allowGoNext = false
	intersect := false //rolesInRoute và authinfo.Roles giao nhau ít nhất một role --> true

	for userrole := range userRoles {
		matchRole := rolesInRoute[userrole]
		if matchRole != nil {
			intersect = true
			if matchRole.(bool) { //Allowed Role -> C1A
				allowGoNext = true
				if !config.ForbidOverAllow { //Allowed Role cao hơn
					break
				}
			} else { //Forbid Role -> C1B
				allowGoNext = false
				if config.ForbidOverAllow { //Forbid Role cao hơn
					break
				}
			}
		} //nếu matchRole == nil thì loop tiếp
	}

	/* Khi rolesInRoute và authinfo.Roles không giao nhau bất kỳ role nào ~ intersect == false
	Mà rolesInRoute lại chứa forbiden role có nghĩa là role người dùng không bị cấm
	Trường hợp này áp dụng cho route sử dụng rbac.Forbid()
	*/
	if !intersect {
		for _, value := range rolesInRoute {
			if !value.(bool) { //chứa ít nhất một forbidden role
				allowGoNext = true
				break
			}
		}
	}
	return allowGoNext
}
