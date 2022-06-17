package handler

import (
	"net/http"

	"ErotsServer/app/passport/database"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/handler"
)

func LoadRouters(e *gin.Engine) {
	e.POST("/api/register", func(c *gin.Context) {
		type JsonParam struct {
			Nickname string `form:"nickname" json:"nickname" binding:"required"`
			Email    string `form:"email" json:"email" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}
		var jsonParam JsonParam
		if err := c.ShouldBind(&jsonParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		handler.JsonStatus(c, database.Register(jsonParam.Nickname, jsonParam.Email, jsonParam.Password))
	})

	e.POST("/api/login/password", func(c *gin.Context) {
		type JsonParam struct {
			Account  string `form:"account" json:"account" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}
		var jsonParam JsonParam
		if err := c.ShouldBind(&jsonParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		loginInfo, err := database.ByPassword(jsonParam.Account, jsonParam.Password, c.ClientIP())
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, loginInfo)
	})

	e.POST("/api/token", func(c *gin.Context) {
		type JsonParam struct {
			AccessToken  string `form:"access_token" json:"access_token" binding:"-"`
			RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"-"`
		}
		var jsonParam JsonParam
		if err := c.ShouldBind(&jsonParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		tokenClaims, err := database.ParseToken(jsonParam.AccessToken, jsonParam.RefreshToken, c.ClientIP())
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, tokenClaims)
	})
}
