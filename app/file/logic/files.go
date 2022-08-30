package logic

import (
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/utils"
)

var (
	BasePath = "./storage"
	Logger   = utils.Logger("file")
)

func checkPath(path string) error {
	// 目录不存在则创建
	if exists, _ := utils.PathExists(path); !exists {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Logger.Println(err)
			return err
		}
	}

	return nil
}

func UploadedImage(c *gin.Context, file *multipart.FileHeader, module string) (string, error) {
	extName := file.Filename[strings.LastIndex(file.Filename, "."):]
	fileName := utils.MD5(file.Filename) + extName
	path := "/image/" + module + "/" + time.Now().Format("2006-01-02") + "/"
	savePath := BasePath + path

	// 目录不存在则创建
	if err := checkPath(savePath); err != nil {
		Logger.Println(err)
		return "", err
	}

	if err := c.SaveUploadedFile(file, savePath+fileName); err != nil {
		Logger.Println(err)
		return "", err
	}

	return path + fileName, nil
}

func DeleteImage(module string, date string, name string) error {
	path := BasePath + "/image/" + module + "/" + date + "/" + name

	err := os.Remove(path)
	if err != nil {
		Logger.Println(err)
	}
	return err
}

func UploadedApp(c *gin.Context, file *multipart.FileHeader, name string) (string, error) {
	extName := file.Filename[strings.LastIndex(file.Filename, "."):]
	fileName := utils.MD5(file.Filename) + extName

	appName := name[:strings.LastIndex(name, ".")]

	path := "/apps/" + appName + "/" + time.Now().Format("2006-01-02") + "/"
	savePath := BasePath + path

	// 目录不存在则创建
	if err := checkPath(savePath); err != nil {
		Logger.Println(err)
		return "", err
	}

	if err := c.SaveUploadedFile(file, savePath+fileName); err != nil {
		Logger.Println(err)
		return "", err
	}

	return path + fileName, nil
}
