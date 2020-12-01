#### day 1

* 实现handler

```shell
➜  yun git:(main) tree
.
├── README.md
└── main.go
```

main.go

```go
package main

import (
	"fmt"
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
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))

}
```

```shell
➜  huayun321 curl http://localhost:9999/hello
Header["User-Agent"] = ["curl/7.64.1"]
Header["Accept"] = ["*/*"]
➜  huayun321 curl http://localhost:9999      
URL.Path = "/"
➜  huayun321 curl http://localhost:9999/hello
Header["User-Agent"] = ["curl/7.64.1"]
Header["Accept"] = ["*/*"]
➜  huayun321 curl http://localhost:9999/x    
404 NOT FOUND: /x

```

* 封装进yun包

```shell
➜  yun git:(main) tree
.
├── README.md
├── main.go
└── yun
    └── yun.go

```

/yun/yun.go

```go
package yun

import (
	"fmt"
	"log"
	"net/http"
)

//定义request类型 router中使用
type HandlerFunc func(http.ResponseWriter, *http.Request)

type engine struct {
	router map[string]HandlerFunc
}

//Engine的构造方法
//todo 改为单例模式
func New() *engine {
	return &engine{router: make(map[string]HandlerFunc)}
}

func (e *engine) AddRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
	log.Printf("route added : %s %s", method, pattern)
}

func (e *engine) GET(pattern string, handler HandlerFunc) {
	e.AddRoute("GET", pattern, handler)
}

func (e *engine) POST(pattern string, handler HandlerFunc) {
	e.AddRoute("POST", pattern, handler)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func (e *engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

```

main.go
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
	y.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
	y.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	log.Fatal(y.Run(":9999"))

}

```

```shell
➜  huayun321 curl http://localhost:9999      
URL.Path = "/"
➜  huayun321 curl http://localhost:9999/hello
Header["Accept"] = ["*/*"]
Header["User-Agent"] = ["curl/7.64.1"]
➜  huayun321 curl http://localhost:9999/x    
404 NOT FOUND: /x

```
