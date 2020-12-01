package main

import (
	"fmt"
	"github.com/huayun321/yun/yun"
	"log"
	"net/http"
)

//实现http handler 接口的结构体
type Engine struct{}

//实现ServeHTTP 方法
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func main() {
	y := yun.New()
	y.GET("/", func(c *yun.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
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
