package handler

import (
	"net/http"

	"ErotsServer/app/admin/logic"
	userMiddleware "ErotsServer/app/user/middleware"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/handler"
)

func LoadUserRouters(e *gin.Engine) {
	r := e.Group("/api/user")

	r.Use(userMiddleware.PermitionCheck("rbacManager"))

	r.GET("/info/:page/:count", func(c *gin.Context) {
		type UriParam struct {
			Page  int `uri:"page" binding:"required"`
			Count int `uri:"count" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		result := logic.GetUsers(uriParam.Page, uriParam.Count)

		c.JSON(http.StatusOK, result)
	})

	r.PUT("/info", func(c *gin.Context) {
		userInfo := make(map[string]interface{})
		if err := c.BindJSON(&userInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		uid := int(userInfo["uid"].(float64))

		handler.JsonStatus(c, logic.UpdateUserInfo(uid, userInfo))
	})

	r.GET("/user-role/:uid", func(c *gin.Context) {
		type UriParam struct {
			Uid int `uri:"uid" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		result, err := logic.GetUserRoles(uriParam.Uid)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, result)
	})
	r.POST("/user-role", func(c *gin.Context) {
		type UriParam struct {
			Uid    int `form:"uid" json:"uid" binding:"required"`
			RoleId int `form:"role_id" json:"role_id" binding:"number"`
		}
		var uriParam UriParam
		if err := c.ShouldBind(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, logic.AddUserRole(uriParam.Uid, uriParam.RoleId))
	})
	r.DELETE("/user-role/:uid/:role_id", func(c *gin.Context) {
		type UriParam struct {
			Uid    int `uri:"uid" binding:"required"`
			RoleId int `uri:"role_id" binding:"number"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, logic.DeleteUserRole(uriParam.Uid, uriParam.RoleId))
	})

	r.GET("/search/:count/:kw", func(c *gin.Context) {
		type UriParam struct {
			Count int    `uri:"count" binding:"required"`
			Kw    string `uri:"kw" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		result := logic.SearchUser(uriParam.Count, uriParam.Kw)

		c.JSON(http.StatusOK, result)
	})
}
