// Package core 提供了EasyGo框架的核心功能
package core

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Context 封装了HTTP请求上下文
type Context struct {
	engine     *Engine
	Writer     http.ResponseWriter
	Request    *http.Request
	Params     map[string]string
	handlers   []HandlerFunc
	index      int
	Keys       map[string]interface{}
	StatusCode int
}

// reset 重置上下文
func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Request = r
	c.Params = make(map[string]string)
	c.handlers = nil
	c.index = -1
	c.Keys = make(map[string]interface{})
}

// Next 执行下一个处理函数
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

// JSON 发送JSON响应
func (c *Context) JSON(code int, obj interface{}) {
	c.StatusCode = code
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// BindJSON 绑定JSON请求体
func (c *Context) BindJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(obj)
}

// GetHeader 获取请求头
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// GetParam 获取URL参数
func (c *Context) GetParam(key string) string {
	return c.Params[key]
}

// Set 设置上下文值
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get 获取上下文中的值
func (c *Context) Get(key string) interface{} {
	return c.Keys[key]
}

// Abort 中止请求处理流程
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// IsAborted 检查请求是否已被中止
func (c *Context) IsAborted() bool {
	return c.index == len(c.handlers)
}

// Status 设置HTTP响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// String 返回纯文本响应
// code: HTTP状态码
// format: 格式化字符串
// values: 格式化参数
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Write([]byte(format))
}

// XML 返回XML格式响应
// code: HTTP状态码
// obj: 要序列化的对象
func (c *Context) XML(code int, obj interface{}) {
	c.Status(code)
	c.Writer.Header().Set("Content-Type", "application/xml")
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// BindXML 将请求体解析为XML对象
// obj: 目标对象指针
// 返回解析错误（如果有）
func (c *Context) BindXML(obj interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	return xml.Unmarshal(body, obj)
}

// Query 获取URL查询参数
// key: 参数名
// 返回参数值或空字符串
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// DefaultQuery 获取URL查询参数，如果不存在则返回默认值
// key: 参数名
// defaultValue: 默认值
// 返回参数值或默认值
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value := c.Query(key); value != "" {
		return value
	}
	return defaultValue
}

// PostForm 获取POST或PUT请求的表单参数
// key: 参数名
// 返回参数值或空字符串
func (c *Context) PostForm(key string) string {
	// 确保表单已被解析
	if c.Request.Form == nil {
		_ = c.Request.ParseMultipartForm(32 << 20) // 32 MB 上限
	}
	return c.Request.FormValue(key)
}

// DefaultPostForm 获取POST或PUT请求的表单参数，如果不存在则返回默认值
// key: 参数名
// defaultValue: 默认值
// 返回参数值或默认值
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value := c.PostForm(key); value != "" {
		return value
	}
	return defaultValue
}

// Header 获取请求头
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Param 获取URL路径参数
// key: 参数名
// 返回参数值
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// RawData 获取原始请求体数据
// 返回请求体字节数组和可能的错误
func (c *Context) RawData() ([]byte, error) {
	return io.ReadAll(c.Request.Body)
}

// Bind 根据 Content-Type 自动绑定请求体到目标对象
// obj: 目标对象指针
// 返回绑定错误（如果有）
func (c *Context) Bind(obj interface{}) error {
	contentType := c.Request.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		return c.BindJSON(obj)
	case strings.HasPrefix(contentType, "application/xml"):
		return c.BindXML(obj)
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"),
		strings.HasPrefix(contentType, "multipart/form-data"):
		// 对于表单数据，我们不能直接 Bind 到任意 struct
		// 需要手动解析或使用 reflect
		// 这里我们只处理基本的 string map，如果需要更复杂的 struct 绑定，需要专门的库如 binding
		// 暂时先返回错误，或直接使用 PostForm / Query 方法
		_ = c.Request.ParseMultipartForm(32 << 20) // 确保表单已解析
		// 如果 obj 是一个 map[string]string，我们可以尝试填充它
		// 否则，让用户使用 PostForm/Query
		return nil // 暂时不返回错误，允许后续手动获取参数
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
}
