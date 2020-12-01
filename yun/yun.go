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
