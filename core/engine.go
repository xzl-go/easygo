// Package core 提供了EasyGo框架的核心功能
package core

import (
	"context"
	"fmt"
	"html/template" // 导入 html/template 包
	"net/http"
	"sync"
)

// HandlerFunc 定义了请求处理函数的类型
type HandlerFunc func(ctx *Context)

// Renderer 是一个接口，定义了模板渲染器的行为
type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{}) error
}

// Engine 是框架的核心引擎
// 负责路由管理、中间件处理和HTTP服务器
type Engine struct {
	*RouterGroup
	router      *router
	middlewares []HandlerFunc
	pool        sync.Pool
	HTMLRender  interface {
		Render(w http.ResponseWriter, name string, data interface{}) error
	}
	templates *template.Template
}

// New 创建一个新的引擎实例
func New() *Engine {
	engine := &Engine{
		RouterGroup: &RouterGroup{
			engine: nil,
		},
		router:      newRouter(),
		middlewares: make([]HandlerFunc, 0),
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return &Context{
			engine: engine,
		}
	}
	return engine
}

// Use 添加中间件
func (e *Engine) Use(middlewares ...HandlerFunc) {
	e.middlewares = append(e.middlewares, middlewares...)
}

// GET 注册GET请求处理函数
// path: 请求路径
// handler: 处理函数
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.addRoute("GET", path, handler)
}

// POST 注册POST请求处理函数
// path: 请求路径
// handler: 处理函数
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.router.addRoute("POST", path, handler)
}

// PUT 注册PUT请求处理函数
// path: 请求路径
// handler: 处理函数
func (e *Engine) PUT(path string, handler HandlerFunc) {
	e.router.addRoute("PUT", path, handler)
}

// DELETE 注册DELETE请求处理函数
// path: 请求路径
// handler: 处理函数
func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.router.addRoute("DELETE", path, handler)
}

// ServeHTTP 实现http.Handler接口
// 处理所有HTTP请求，包括路由匹配、中间件执行和请求处理
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.reset(w, r)
	handler, params := e.router.getRoute(r.Method, r.URL.Path)
	if handler != nil {
		ctx.Params = params
		ctx.handlers = append(e.middlewares, handler)
		ctx.Next()
	} else {
		http.NotFound(w, r)
	}
	e.pool.Put(ctx)
}

// Run 启动HTTP服务器
// addr: 服务器监听地址
// 返回服务器运行错误（如果有）
func (e *Engine) Run(addr string) error {
	fmt.Printf("🚀 服务器启动，监听地址：%s\n", addr)
	return http.ListenAndServe(addr, e)
}

// RunTLS 启动HTTPS服务器
// addr: 服务器监听地址
// certFile: SSL证书文件路径
// keyFile: SSL密钥文件路径
// 返回服务器运行错误（如果有）
func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	fmt.Printf("🔒 安全服务器启动，监听地址：%s\n", addr)
	return http.ListenAndServeTLS(addr, certFile, keyFile, e)
}

// Shutdown 优雅关闭服务器
// ctx: 上下文，用于控制关闭超时
// 返回关闭错误（如果有）
func (e *Engine) Shutdown(ctx context.Context) error {
	// TODO: 实现优雅关闭
	return nil
}

// SetHTMLRender 设置自定义的 HTML 渲染器
func (e *Engine) SetHTMLRender(render Renderer) {
	e.HTMLRender = render
}

// LoadHTMLGlob 加载 HTML 模板文件
// glob: 匹配模板文件的 glob 模式，例如 "templates/*"
func (e *Engine) LoadHTMLGlob(glob string) {
	e.templates = template.Must(template.ParseGlob(glob))
}

// LoadHTMLFiles 加载HTML文件
func (e *Engine) LoadHTMLFiles(files ...string) {
	e.templates = template.Must(template.ParseFiles(files...))
}
