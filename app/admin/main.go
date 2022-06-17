package admin

import (
	"net/http"

	"ErotsServer/app/admin/handler"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/middleware"
	"github.com/ipuppet/gtools/server"
)

func GetServer(addr string) *http.Server {
	return server.GetServer(addr, func(engine *gin.Engine) {
		engine.Use(middleware.Cors("admin"))

		handler.LoadRouters(engine)
	})
}
