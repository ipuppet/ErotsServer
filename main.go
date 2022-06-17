package main

import (
	"log"

	"ErotsServer/app/admin"
	"ErotsServer/app/file"
	"ErotsServer/app/passport"
	"ErotsServer/app/www"

	"github.com/gin-gonic/gin"
	_ "github.com/ipuppet/gtools/flags"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	wwwServer := www.GetServer(":8080")
	passportServer := passport.GetServer(":8081")
	adminServer := admin.GetServer(":8082")
	fileServer := file.GetServer(":8083")

	g.Go(func() error {
		return wwwServer.ListenAndServe()
	})
	g.Go(func() error {
		return passportServer.ListenAndServe()
	})
	g.Go(func() error {
		return adminServer.ListenAndServe()
	})
	g.Go(func() error {
		return fileServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
