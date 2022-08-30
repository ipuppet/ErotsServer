package handler

import (
	"net/http"

	"ErotsServer/app/admin/middleware"
	userMiddleware "ErotsServer/app/user/middleware"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/database"
)

func LoadRouters(e *gin.Engine) {
	e.Use(userMiddleware.Authorize())
	e.Use(middleware.SetAdminUser())

	// load modules
	LoadRBACRouters(e)
	LoadUserRouters(e)

	e.DELETE("/api/cache", func(c *gin.Context) {
		// TODO /api/cache
		database.CleanCache()

		c.JSON(http.StatusOK, "")
	})
}
