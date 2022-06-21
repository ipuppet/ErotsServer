package handler

import (
	"net/http"

	"ErotsServer/app/user/pkg"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/config"
)

func LoadMenuRouters(e *gin.Engine) {
	e.GET("/api/menu", func(c *gin.Context) {
		userFromContext, _ := c.Get("User")
		user := userFromContext.(*pkg.User)
		permittedModules := user.GetPermittedModules()

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
