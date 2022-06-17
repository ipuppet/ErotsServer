package structure

type Role struct {
	RoleId      *int   `form:"role_id" json:"role_id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
}

type Permission struct {
	PermissionId *int   `form:"permission_id" json:"permission_id" binding:"required"`
	Name         string `form:"name" json:"name" binding:"required"`
	Description  string `form:"description" json:"description" binding:"required"`
}
