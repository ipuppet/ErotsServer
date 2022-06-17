package handler

import (
	"net/http"

	"ErotsServer/app/passport/database"

	"github.com/gin-gonic/gin"
)

func LoadRouters(e *gin.Engine) {
	e.POST("/api/register", func(c *gin.Context) {
		type JsonParam struct {
			Email    string `form:"email" json:"email" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}
		var jsonParam JsonParam
		if err := c.ShouldBind(&jsonParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		registerInfo, err := database.ByPassword(jsonParam.Email, jsonParam.Password, c.ClientIP())
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, registerInfo)
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
