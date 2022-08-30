package middleware

import (
	"errors"
	"net/http"
	"strings"

	passportDao "ErotsServer/app/passport/dao"
	passportLogic "ErotsServer/app/passport/logic"
	"ErotsServer/app/user/dao"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authorization, "Bearer") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authorization[strings.Index(authorization, " ")+1:]
		claims, err := passportLogic.ParseAccessToken(tokenString, &passportDao.User{})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := &dao.User{
			Uid: claims.Uid,
		}
		c.Set("User", user)
		c.Next()
	}
}

func PermitionCheck(moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userFromContext, _ := c.Get("User")
		user := userFromContext.(*dao.User)

		isPermitted, err := user.IsPermitted(moduleName)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("cannot get user permissions"))
			return
		}
		if isPermitted {
			c.Next()
			return
		}

		c.AbortWithError(http.StatusUnauthorized, errors.New("insufficient permissions"))
	}
}
