package handler

import (
	"net/http"

	"ErotsServer/app/admin/dao"
	userDao "ErotsServer/app/user/dao"
	userMiddleware "ErotsServer/app/user/middleware"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/handler"
)

func LoadRBACRouters(e *gin.Engine) {
	r := e.Group("/api/rbac")

	r.Use(userMiddleware.PermitionCheck("rbacManager"))

	r.GET("/roles", func(c *gin.Context) {
		roles, err := dao.GetRoles()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, roles)
	})

	r.PUT("/role", func(c *gin.Context) {
		var role userDao.Role
		if err := c.Bind(&role); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.UpdateRole(role))
	})
	r.POST("/role", func(c *gin.Context) {
		var role userDao.Role
		if err := c.Bind(&role); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.AddRole(role))
	})
	r.DELETE("/role/:role_id", func(c *gin.Context) {
		type UriParam struct {
			RoleId int `uri:"role_id" binding:"number"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.DeleteRole(uriParam.RoleId))
	})

	r.GET("/permissions", func(c *gin.Context) {
		permissions, err := dao.GetPermissions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, permissions)
	})

	r.PUT("/permission", func(c *gin.Context) {
		var permission userDao.Permission
		if err := c.Bind(&permission); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.UpdatePermission(permission))
	})
	r.POST("/permission", func(c *gin.Context) {
		var permission userDao.Permission
		if err := c.Bind(&permission); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.AddPermission(permission))
	})
	r.DELETE("/permission/:permission_id", func(c *gin.Context) {
		type UriParam struct {
			PermissionId int `uri:"permission_id" binding:"number"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.DeletePermission(uriParam.PermissionId))
	})

	r.GET("/role/:role_id/permissions", func(c *gin.Context) {
		type UriParam struct {
			RoleId int `uri:"role_id" binding:"number"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		result, err := dao.GetRolePermissions(uriParam.RoleId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})
	r.DELETE("/role/:role_id/permission/:permission_id", func(c *gin.Context) {
		type UriParam struct {
			RoleId       int `uri:"role_id" binding:"number"`
			PermissionId int `uri:"permission_id" binding:"number"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.DeleteRolePermission(uriParam.RoleId, uriParam.PermissionId))
	})
	r.POST("/role/:role_id/permission/:permission_id", func(c *gin.Context) {
		type UriParam struct {
			RoleId       int `uri:"role_id" binding:"number"`
			PermissionId int `uri:"permission_id" binding:"number"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, dao.AddRolePermission(uriParam.RoleId, uriParam.PermissionId))
	})
}
