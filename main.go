package main

import (
	"log"

	"ErotsServer/app/admin"
	"ErotsServer/app/file"
	"ErotsServer/app/passport"
	"ErotsServer/app/user"
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

	wwwServer := www.GetServer("localhost:8080")
	passportServer := passport.GetServer("localhost:8081")
	adminServer := admin.GetServer("localhost:8082")
	userServer := user.GetServer("localhost:8083")
	fileServer := file.GetServer("localhost:8084")

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
		return userServer.ListenAndServe()
	})
	g.Go(func() error {
		return fileServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
