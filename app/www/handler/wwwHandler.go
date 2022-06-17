package handler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ipuppet/gtools/utils"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func LoadRouters(e *gin.Engine) {
	e.GET("/api/site-info", func(c *gin.Context) {
		var res map[string]interface{}
		if err := utils.GetStorageJSON("www", "site-info.json", &res); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, res)
	})

	e.Any("/api/proxy", func(c *gin.Context) {
		password := c.Query("password")
		if password != "9ds4v984erg" {
			c.JSON(http.StatusBadRequest, gin.H{"proxy_error": "invalid password"})
			return
		}

		targetUrl := c.DefaultQuery("url", "")
		if targetUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"proxy_error": "url cannot be empty"})
			return
		}
		targetUrl, _ = url.QueryUnescape(strings.Trim(targetUrl, "/"))

		// 构造请求
		requ, _ := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
		for key, values := range c.Request.Header {
			if len(values) == 1 {
				requ.Header.Set(key, values[0])
			} else {
				requ.Header.Set(key, values[0])
				for _, value := range values[1:] {
					requ.Header.Add(key, value)
				}
			}
		}

		// 解析 url
		u, err := url.Parse(targetUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"proxy_error": err.Error()})
			return
		}
		requ.Host = u.Host

		resp, err := httpClient.Do(requ)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"proxy_error": err.Error()})
			return
		}

		defer resp.Body.Close()

		// header
		for key, values := range resp.Header {
			if len(values) == 1 {
				c.Writer.Header().Set(key, values[0])
			} else {
				c.Writer.Header().Set(key, values[0])
				for _, value := range values[1:] {
					c.Writer.Header().Add(key, value)
				}
			}
		}

		c.DataFromReader(
			resp.StatusCode,
			resp.ContentLength,
			resp.Header.Get("Content-Type"),
			resp.Body,
			nil,
		)
	})
}
