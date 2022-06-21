package handler

import (
	"ErotsServer/app/user/pkg"

	"github.com/gin-gonic/gin"
)

func LoadRouters(e *gin.Engine) {
	e.Use(pkg.Authorize())

	// load modules
	LoadMenuRouters(e)
	LoadUserRouters(e)
}
