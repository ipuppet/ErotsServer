package handler

import (
	"net/http"

	"ErotsServer/app/user/pkg"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/handler"
)

func LoadUserRouters(e *gin.Engine) {
	r := e.Group("/api/user")

	r.GET("/structure", func(c *gin.Context) {
		userStructure := map[string]interface{}{}
		config.GetConfig("userStructure.json", &userStructure)

		c.JSON(http.StatusOK, userStructure)
	})

	r.PUT("/info/self", func(c *gin.Context) {
		userFromContext, _ := c.Get("User")
		adminUser := userFromContext.(*pkg.User)

		userInfo := make(map[string]interface{})
		if err := c.BindJSON(&userInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		handler.JsonStatus(c, adminUser.UpdateInfo(userInfo))
	})
}
