package handler

import (
	"errors"
	"net/http"
	"strings"

	"ErotsServer/app/admin/database"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/config"
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
		claims, err := database.ParseAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := database.User{
			Uid:   claims.Uid,
			Roles: claims.Roles,
		}
		c.Set("user", user)
		c.Next()
	}
}

func PermitionCheck(moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userFromContext, _ := c.Get("user")
		adminUser := userFromContext.(database.User)
		permittedModules := adminUser.GetPermittedModules()

		for _, permittedModule := range permittedModules {
			if moduleName == permittedModule || permittedModule == "ALL" {
				c.Next()
				return
			}
		}

		c.AbortWithError(http.StatusUnauthorized, errors.New("insufficient permissions"))
	}
}

func LoadRouters(e *gin.Engine) {
	e.Use(Authorize())

	// load modules
	LoadRBACRouters(e)
	LoadUserRouters(e)

	e.GET("/api/menu", func(c *gin.Context) {
		userFromContext, _ := c.Get("user")
		adminUser := userFromContext.(database.User)
		permittedModules := adminUser.GetPermittedModules()

		menus := []map[string]interface{}{}
		config.GetConfig("menu.json", &menus)

		for _, module := range permittedModules {
			if module == "ALL" {
				c.JSON(http.StatusOK, menus)
				return
			}
		}

		result := []map[string]interface{}{}
		for _, menu := range menus {
			for _, module := range permittedModules {
				if menu["module"] == module || menu["module"] == "public" {
					result = append(result, menu)
				}
			}
		}

		c.JSON(http.StatusOK, result)
	})
}
