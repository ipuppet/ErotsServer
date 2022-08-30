package handler

import (
	userMiddleware "ErotsServer/app/user/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRouters(e *gin.Engine) {
	e.Use(userMiddleware.Authorize())

	// load modules
	LoadFilesRouters(e)
	LoadCommentRouters(e)
}
