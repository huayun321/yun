# yun
### yun.web框架
----

#### 用途

* 该项目为学习GO HTTP框架原理所用
* 提供的功能为路由 context 前缀树路由 路由分组 中间件 模版 错误恢复

#### 使用方法

```go
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
```