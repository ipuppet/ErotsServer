package handler

import (
	"ErotsServer/app/user/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRouters(e *gin.Engine) {
	e.Use(middleware.Authorize())

	// load modules
	LoadMenuRouters(e)
	LoadUserRouters(e)
}
