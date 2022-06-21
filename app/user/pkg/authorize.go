package pkg

import (
	"errors"
	"net/http"
	"strings"

	passportDatabase "ErotsServer/app/passport/database"
	passportPkg "ErotsServer/app/passport/pkg"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ipuppet/gtools/utils"
)

func ParseAccessToken(tokenString string) (*passportPkg.AccessTokenClaims, error) {
	var user passportPkg.User
	var tokenKey string
	accessToken, err := jwt.ParseWithClaims(tokenString, &passportPkg.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*passportPkg.AccessTokenClaims)
		user = passportDatabase.GetUserInfoByUid(claims.Uid)
		tokenKey = utils.MD5(user.Password)
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := accessToken.Claims.(*passportPkg.AccessTokenClaims); ok && accessToken.Valid {
		// 使用令牌内的信息，减少查库次数
		return claims, nil
	}

	return nil, err
}

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
		claims, err := ParseAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := &User{
			Db:    Db,
			Uid:   claims.Uid,
			Roles: claims.Roles,
		}
		c.Set("User", user)
		c.Next()
	}
}

func PermitionCheck(moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userFromContext, _ := c.Get("User")
		user := userFromContext.(*User)
		permittedModules := user.GetPermittedModules()

		for _, permittedModule := range permittedModules {
			if moduleName == permittedModule || permittedModule == "ALL" {
				c.Next()
				return
			}
		}

		c.AbortWithError(http.StatusUnauthorized, errors.New("insufficient permissions"))
	}
}
