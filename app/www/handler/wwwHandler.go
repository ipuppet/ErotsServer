package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/utils"
)

func LoadRouters(e *gin.Engine) {
	e.GET("/api/site-info", func(c *gin.Context) {
		var res map[string]interface{}
		if err := utils.GetStorageJSON("www", "site-info.json", &res); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, res)
	})
}
