package handler

import (
	"net/http"

	"ErotsServer/app/admin/database"
	"ErotsServer/app/admin/structure"
	userPkg "ErotsServer/app/user/pkg"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/handler"
)

func LoadRBACRouters(e *gin.Engine) {
	r := e.Group("/api/rbac")

	r.Use(userPkg.PermitionCheck("rbacManager"))

	r.GET("/roles", func(c *gin.Context) {
		roles, err := database.GetRoles()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, roles)
	})

	r.PUT("/role", func(c *gin.Context) {
		var role structure.Role
		if err := c.Bind(&role); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.UpdateRole(role))
	})
	r.POST("/role", func(c *gin.Context) {
		var role structure.Role
		if err := c.Bind(&role); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.AddRole(role))
	})
	r.DELETE("/role/:role_id", func(c *gin.Context) {
		type UriParam struct {
			RoleId *int `uri:"role_id" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.DeleteRole(*uriParam.RoleId))
	})

	r.GET("/permissions", func(c *gin.Context) {
		permissions, err := database.GetPermissions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, permissions)
	})

	r.PUT("/permission", func(c *gin.Context) {
		var permission structure.Permission
		if err := c.Bind(&permission); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.UpdatePermission(permission))
	})
	r.POST("/permission", func(c *gin.Context) {
		var permission structure.Permission
		if err := c.Bind(&permission); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.AddPermission(permission))
	})
	r.DELETE("/permission/:permission_id", func(c *gin.Context) {
		type UriParam struct {
			PermissionId *int `uri:"permission_id" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.DeletePermission(*uriParam.PermissionId))
	})

	r.GET("/role/:role_id/permissions", func(c *gin.Context) {
		type UriParam struct {
			RoleId *int `uri:"role_id" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		result, err := database.GetRolePermissions(*uriParam.RoleId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})
	r.DELETE("/role/:role_id/permission/:permission_id", func(c *gin.Context) {
		type UriParam struct {
			RoleId       *int `uri:"role_id" binding:"required"`
			PermissionId *int `uri:"permission_id" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.DeleteRolePermission(*uriParam.RoleId, *uriParam.PermissionId))
	})
	r.POST("/role/:role_id/permission/:permission_id", func(c *gin.Context) {
		type UriParam struct {
			RoleId       *int `uri:"role_id" binding:"required"`
			PermissionId *int `uri:"permission_id" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.AddRolePermission(*uriParam.RoleId, *uriParam.PermissionId))
	})
}
