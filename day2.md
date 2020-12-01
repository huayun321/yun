
#### day2 

* 使用Context封装HTTP返回的重复操作，例如设置返回header body 序列化

```shell
➜  yun git:(main) ✗ tree
.
├── README.md
├── main.go
└── yun
    ├── context.go
    └── yun.go
```

/yun/context.go

```go
package yun

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{} 

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

```


* 提取出路由部分 方便以后扩充内容

```shell
➜  yun git:(main) ✗ tree
.
├── README.md
├── main.go
└── yun
    ├── context.go
    ├── router.go
    └── yun.go

```

/yun/router.go

```go
package yun

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

```

/yun/yun.go
```go
package yun

import (
	"net/http"
)

//定义request类型 router中使用
type HandlerFunc func(*Context)

type engine struct {
	router *router
}

//Engine的构造方法
//todo 改为单例模式
func New() *engine {
	return &engine{router: newRouter()}
}

func (e *engine) AddRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

func (e *engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}

func (e *engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

```

```go
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

```

```shell
➜  huayun321 curl -i http://localhost:9999/                  
HTTP/1.1 200 OK
Content-Type: text/html
Date: Tue, 01 Dec 2020 12:18:01 GMT
Content-Length: 18

<h1>Hello Yun</h1>%                    
                             
➜  huayun321 curl "http://localhost:9999/hello?name=yun"
hello yun, you're at /hello

➜  huayun321 curl "http://localhost:9999/login" -X POST -d 'username=yun&password=1234'
{"password":"1234","username":"yun"}

➜  huayun321 curl -i http://localhost:9999/x            
HTTP/1.1 404 Not Found
Content-Type: text/plain
Date: Tue, 01 Dec 2020 12:18:14 GMT
Content-Length: 18

404 NOT FOUND: /x
```
