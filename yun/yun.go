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
