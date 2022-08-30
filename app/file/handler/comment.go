package handler

import (
	userMiddleware "ErotsServer/app/user/middleware"

	"github.com/gin-gonic/gin"
)

func LoadCommentRouters(e *gin.Engine) {
	e.Use(userMiddleware.PermitionCheck("comment"))

	e.GET("/api/comment/:app", func(c *gin.Context) {

	})
}
