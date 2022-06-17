package handler

import (
	"net/http"

	"ErotsServer/app/admin/database"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/handler"
)

func LoadUserRouters(e *gin.Engine) {
	r := e.Group("/api/user")

	r.Use(PermitionCheck("rbacManager"))

	r.GET("/structure", func(c *gin.Context) {
		userStructure := map[string]interface{}{}
		config.GetConfig("userStructure.json", &userStructure)

		c.JSON(http.StatusOK, userStructure)
	})

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

		result := database.GetUsers((uriParam.Page-1)*uriParam.Count, uriParam.Count)

		c.JSON(http.StatusOK, result)
	})

	r.PUT("/info", func(c *gin.Context) {
		userInfo := make(map[string]interface{})
		if err := c.BindJSON(&userInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		uid := int(userInfo["uid"].(float64))
		user := database.User{Uid: uid}

		handler.JsonStatus(c, user.UpdateInfo(userInfo))
	})

	r.PUT("/info/self", func(c *gin.Context) {
		userFromContext, _ := c.Get("user")
		adminUser := userFromContext.(database.User)

		userInfo := make(map[string]interface{})
		if err := c.BindJSON(&userInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		handler.JsonStatus(c, adminUser.UpdateInfo(userInfo))
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

		user := database.User{Uid: uriParam.Uid}

		result := user.GetRoles()

		c.JSON(http.StatusOK, result)
	})
	r.POST("/user-role", func(c *gin.Context) {
		type UriParam struct {
			Uid    int  `form:"uid" json:"uid" binding:"required"`
			RoleId *int `form:"role_id" json:"role_id" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBind(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user := database.User{Uid: uriParam.Uid}

		handler.JsonStatus(c, user.AddRole(*uriParam.RoleId))
	})
	r.DELETE("/user-role/:uid/:role_id", func(c *gin.Context) {
		type UriParam struct {
			Uid    int  `uri:"uid" binding:"required"`
			RoleId *int `uri:"role_id" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user := database.User{Uid: uriParam.Uid}

		handler.JsonStatus(c, user.DeleteRole(*uriParam.RoleId))
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

		result := database.SearchUser(uriParam.Count, uriParam.Kw)

		c.JSON(http.StatusOK, result)
	})
}
