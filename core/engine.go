package core

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	router     *router
	middleware []MiddlewareFunc
	pool       sync.Pool
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.pool.New = func() interface{} {
		return &Context{}
	}
	return engine
}

func (e *Engine) Use(middleware ...MiddlewareFunc) {
	e.middleware = append(e.middleware, middleware...)
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.addRoute("GET", path, handler)
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.router.addRoute("POST", path, handler)
}

func (e *Engine) PUT(path string, handler HandlerFunc) {
	e.router.addRoute("PUT", path, handler)
}

func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.router.addRoute("DELETE", path, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.reset(w, r)
	defer e.pool.Put(ctx)

	handler, params := e.router.getRoute(r.Method, r.URL.Path)
	if handler == nil {
		ctx.String(http.StatusNotFound, "404 NOT FOUND")
		return
	}

	ctx.Params = params

	// 执行中间件
	for _, m := range e.middleware {
		m(ctx)
		if ctx.IsAborted() {
			return
		}
	}

	handler(ctx)
}

func (e *Engine) Run(addr string) error {
	fmt.Printf("🚀 服务器启动，监听地址：%s\n", addr)
	return http.ListenAndServe(addr, e)
}

func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	fmt.Printf("🔒 安全服务器启动，监听地址：%s\n", addr)
	return http.ListenAndServeTLS(addr, certFile, keyFile, e)
}

func (e *Engine) Shutdown(ctx context.Context) error {
	// 在这里可以添加优雅关闭的逻辑
	return nil
}
