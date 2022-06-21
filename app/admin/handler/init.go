package handler

import (
	"ErotsServer/app/admin/internal"
	userPkg "ErotsServer/app/user/pkg"

	"github.com/gin-gonic/gin"
)

func SetAdminUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("User")
		adminUser := &internal.User{
			User: user.(*userPkg.User),
		}
		c.Set("AdminUser", adminUser)
		c.Next()
	}
}

func LoadRouters(e *gin.Engine) {
	e.Use(userPkg.Authorize())
	e.Use(SetAdminUser())

	// load modules
	LoadRBACRouters(e)
	LoadUserRouters(e)
}
