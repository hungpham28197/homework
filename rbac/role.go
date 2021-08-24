package rbac

import (
	"home-work/model"
)

//Danh sách các role có thể truy xuất
func Allow(roles ...int) RoleExp {
	return func() model.Roles {
		mapRoles := make(model.Roles)
		for _, role := range roles {
			mapRoles[role] = true
		}
		return mapRoles
	}
}

//Cho phép tất cả các role
func AllowAll() RoleExp {
	return func() model.Roles {
		mapRoles := make(model.Roles)
		for _, role := range allRoles {
			mapRoles[role] = true
		}
		return mapRoles
	}
}

//Danh sách các role bị cấm truy cập
func Forbid(roles ...int) RoleExp {
	return func() model.Roles {
		mapRoles := make(model.Roles)
		for _, role := range roles {
			mapRoles[role] = false
		}
		return mapRoles
	}
}

//Cấm tất cả các role ngoại trừ root
func ForbidAll() RoleExp {
	return func() model.Roles {
		mapRoles := make(model.Roles)
		for _, role := range allRoles {
			mapRoles[role] = false
		}
		return mapRoles
	}
}
