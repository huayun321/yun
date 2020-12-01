package main

import (
	"github.com/huayun321/yun/yun"
	"log"
	"net/http"
)

func main() {
	y := yun.New()
	y.GET("/", func(c *yun.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Yun</h1>")
	})
	y.GET("/hello", func(c *yun.Context) {
		// expect /hello?name=yun
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	y.POST("/login", func(c *yun.Context) {
		c.JSON(http.StatusOK, yun.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	log.Fatal(y.Run(":9999"))

}
