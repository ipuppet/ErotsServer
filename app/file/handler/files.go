package handler

import (
	"errors"
	"net/http"
	"os"

	"ErotsServer/app/file/logic"
	userMiddleware "ErotsServer/app/user/middleware"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/utils"
)

var (
	BasePath = "./storage"
)

func checkPath(path string) error {
	// 目录不存在则创建
	if exists, _ := utils.PathExists(path); !exists {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			utils.Logger("file").Println(err)
			return err
		}
	}

	return nil
}

func LoadFilesRouters(e *gin.Engine) {
	e.Use(userMiddleware.PermitionCheck("file"))

	e.POST("/api/file/image/:module", func(c *gin.Context) {
		type UriParam struct {
			Module string `uri:"module" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("form key should be \"image\""))
			return
		}

		path, err := logic.UploadedImage(c, file, uriParam.Module)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"path": path,
		})
	})

	e.DELETE("/api/file/image/:module/:date/:name", func(c *gin.Context) {
		type UriParam struct {
			Module string `uri:"module" binding:"required"`
			Date   string `uri:"date" binding:"required"`
			Name   string `uri:"name" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err := logic.DeleteImage(uriParam.Module, uriParam.Date, uriParam.Name)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "")
	})

	e.POST("/api/file/app/:name", func(c *gin.Context) {
		type UriParam struct {
			Name string `uri:"name" binding:"required"`
		}
		var uriParam UriParam
		if err := c.ShouldBindUri(&uriParam); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		file, err := c.FormFile("app")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("form key should be \"app\""))
			return
		}

		path, err := logic.UploadedApp(c, file, uriParam.Name)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"path": path,
		})
	})
}
