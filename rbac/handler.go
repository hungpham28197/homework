package rbac

import (
	"home-work/model"
	"net/http"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
)

/*
Gán role vào route và path
route = HTTP method + path
roles là kết quả trả về từ hàm kiểu roleExp()
*/
func assignRoles(method string, path string, roles model.Roles) {
	routesRoles[method+" "+path] = roles

	if httpVerbRole := pathsRoles[path]; httpVerbRole == nil {
		pathsRoles[path] = make(HTTPVerbRoles)
	}
	pathsRoles[path][method] = roles
}

func Get(party router.Party, relativePath string, roleExp RoleExp, handlers ...context.Handler) {
	party.Handle(http.MethodGet, relativePath, handlers...)
	assignRoles(http.MethodGet, party.GetRelPath()+relativePath, roleExp())
}

func Post(party router.Party, relativePath string, roleExp RoleExp, handlers ...context.Handler) {
	party.Handle(http.MethodPost, relativePath, handlers...)

	assignRoles(http.MethodPost, party.GetRelPath()+relativePath, roleExp())
}

func Put(party router.Party, relativePath string, roleExp RoleExp, handlers ...context.Handler) {
	party.Handle(http.MethodPut, relativePath, handlers...)
	assignRoles(http.MethodPut, party.GetRelPath()+relativePath, roleExp())
}

func Delete(party router.Party, relativePath string, roleExp RoleExp, handlers ...context.Handler) {
	party.Handle(http.MethodDelete, relativePath, handlers...)
	assignRoles(http.MethodDelete, party.GetRelPath()+relativePath, roleExp())
}

func Any(party router.Party, relativePath string, roleExp RoleExp, handlers ...context.Handler) {
	party.Any(relativePath, handlers...)
	for _, method := range router.AllMethods {
		assignRoles(method, party.GetRelPath()+relativePath, roleExp())
	}
}
