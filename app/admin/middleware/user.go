package middleware

import (
	"ErotsServer/app/admin/dao"
	userDao "ErotsServer/app/user/dao"

	"github.com/gin-gonic/gin"
)

func SetAdminUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("User")
		adminUser := &dao.User{
			User: user.(*userDao.User),
		}
		c.Set("AdminUser", adminUser)
		c.Next()
	}
}
