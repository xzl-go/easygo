package core

// RouterGroup 是路由组
type RouterGroup struct {
	engine      *Engine
	prefix      string
	middlewares []HandlerFunc
}

// Group 创建一个新的路由组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		engine:      group.engine,
		prefix:      group.prefix + prefix,
		middlewares: make([]HandlerFunc, 0),
	}
}

// Use 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// GET 注册GET请求处理函数
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.engine.router.addRoute("GET", group.prefix+pattern, handler)
}

// POST 注册POST请求处理函数
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.engine.router.addRoute("POST", group.prefix+pattern, handler)
}

// PUT 注册PUT请求处理函数
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.engine.router.addRoute("PUT", group.prefix+pattern, handler)
}

// DELETE 注册DELETE请求处理函数
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.engine.router.addRoute("DELETE", group.prefix+pattern, handler)
}
