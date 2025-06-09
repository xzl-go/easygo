// Package core 提供了EasyGo框架的核心功能
package core

import (
	"strings"
)

// MiddlewareFunc 定义了中间件函数的类型
// type MiddlewareFunc func(ctx *Context) // 删除了此行

// node 表示路由树中的节点
type node struct {
	pattern  string           // 路由模式
	part     string           // 路由部分
	children map[string]*node // 子节点
	isWild   bool             // 是否是通配符节点
	handler  HandlerFunc      // 处理函数
}

// router 是路由管理器
// 实现了基于前缀树的路由匹配
type router struct {
	roots    map[string]*node       // 路由树根节点
	handlers map[string]HandlerFunc // 路由处理函数
	engine   *Engine                // 引擎引用
}

// newRouter 创建新的路由器
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析路由模式
func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return result
}

// insert 插入路由
func (r *router) insert(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{children: make(map[string]*node)}
	}
	root := r.roots[method]
	for _, part := range parts {
		if _, ok := root.children[part]; !ok {
			root.children[part] = &node{
				part:     part,
				children: make(map[string]*node),
				isWild:   part[0] == ':' || part[0] == '*',
			}
		}
		root = root.children[part]
	}
	root.pattern = pattern
	root.handler = handler
	r.handlers[key] = handler
}

// search 搜索路由
func (r *router) search(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root
	for i, part := range searchParts {
		var found bool
		for _, child := range n.children {
			if child.part == part || child.isWild {
				if child.part[0] == '*' {
					params[child.part[1:]] = strings.Join(searchParts[i:], "/")
					return child, params
				}
				if child.part[0] == ':' {
					params[child.part[1:]] = part
				}
				n = child
				found = true
				break
			}
		}
		if !found {
			return nil, nil
		}
	}
	return n, params
}

// addRoute 添加路由
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	r.insert(method, pattern, handler)
}

// getRoute 获取路由
func (r *router) getRoute(method, path string) (HandlerFunc, map[string]string) {
	n, params := r.search(method, path)
	if n != nil {
		return n.handler, params
	}
	return nil, nil
}
