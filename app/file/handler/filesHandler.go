package handler

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	userPkg "ErotsServer/app/user/pkg"

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

func LoadRouters(e *gin.Engine) {
	e.Use(userPkg.Authorize())
	e.Use(userPkg.PermitionCheck("file"))

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

		extName := file.Filename[strings.LastIndex(file.Filename, "."):]
		fileName := utils.MD5(file.Filename) + extName
		path := "/image/" + uriParam.Module + "/" + time.Now().Format("2006-01-02") + "/"
		savePath := BasePath + path

		// 目录不存在则创建
		if err := checkPath(savePath); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = c.SaveUploadedFile(file, savePath+fileName)
		if err != nil {
			utils.Logger("file").Println(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"path": path + fileName,
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

		path := BasePath + "/image/" + uriParam.Module + "/" + uriParam.Date + "/" + uriParam.Name

		err := os.Remove(path)
		if err != nil {
			utils.Logger("file").Println(err)
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

		extName := file.Filename[strings.LastIndex(file.Filename, "."):]
		fileName := utils.MD5(file.Filename) + extName

		appName := uriParam.Name[:strings.LastIndex(uriParam.Name, ".")]

		path := "/apps/" + appName + "/" + time.Now().Format("2006-01-02") + "/"
		savePath := BasePath + path

		// 目录不存在则创建
		if err := checkPath(savePath); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = c.SaveUploadedFile(file, savePath+fileName)
		if err != nil {
			utils.Logger("file").Println(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"path": path + fileName,
		})
	})
}
