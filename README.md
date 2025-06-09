# EasyGo 框架

EasyGo 是一个轻量级的 Go Web 框架，提供了丰富的功能和易用的 API，帮助开发者快速构建高性能的 Web 应用。

## 主要特性

- 🚀 高性能路由引擎
- 🔒 JWT 认证
- 🌐 国际化支持
- 🔐 RBAC 权限控制
- 📝 日志系统
- 🔄 WebSocket 支持
- ⏰ 定时任务
- 🔍 链路追踪
- ✅ 参数验证
- 🛡️ 中间件支持

## 快速开始

### 安装

```bash
go get github.com/xzl-go/easygo
```

### 示例代码

```go
package main

import (
    "github.com/xzl-go/easygo/core"
    "github.com/xzl-go/easygo/middleware"
)

func main() {
    // 创建应用实例
    app := core.New()

    // 使用中间件
    app.Use(middleware.Logger())
    app.Use(middleware.Recovery())

    // 注册路由
    app.GET("/", func(ctx *core.Context) {
        ctx.JSON(200, map[string]string{
            "message": "Hello, EasyGo!",
        })
    })

    // 启动服务器
    app.Run(":8080")
}
```

## 核心功能

### 路由系统

```go
// 基本路由
app.GET("/path", handler)
app.POST("/path", handler)
app.PUT("/path", handler)
app.DELETE("/path", handler)

// 路由组
group := app.Group("/api")
group.GET("/users", handler)
group.POST("/users", handler)
```

### 中间件

```go
// 使用中间件
app.Use(middleware.Logger())
app.Use(middleware.Recovery())

// 自定义中间件
func CustomMiddleware() core.HandlerFunc {
    return func(c *core.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}
```

### JWT 认证

```go
// 初始化 JWT 管理器
jwtManager := jwt.NewJWTManager("your_secret_key", 24*time.Hour)

// 生成令牌
token, err := jwtManager.GenerateToken(userID, username)

// 验证令牌
claims, err := jwtManager.VerifyToken(token)
```

### RBAC 权限控制

```go
// 初始化 RBAC 管理器
rbacManager, err := rbac.NewRBACManager("rbac_model.conf", "rbac_policy.csv")

// 检查权限
allowed, err := rbacManager.Enforce(user, resource, action)
```

### 国际化

```go
// 初始化国际化管理器
i18nManager := i18n.New("en")
i18nManager.LoadTranslations("i18n/translations")

// 翻译文本
message := i18nManager.Translate("key", lang)
```

### WebSocket

```go
// 注册 WebSocket 路由
app.GET("/ws", func(ctx *core.Context) {
    websocket.HandleWebSocket(ctx)
})
```

### 定时任务

```go
// 初始化定时任务
cron.InitCron()

// 添加定时任务
cron.AddJob("@every 1m", func() {
    // 任务逻辑
})
```

### 链路追踪

```go
// 初始化追踪器
tracer := tracing.NewTracer("service-name")
defer tracer.Shutdown(context.Background())
```

## 项目结构

```
easygo/
├── core/           # 核心功能
├── middleware/     # 中间件
├── jwt/           # JWT 认证
├── rbac/          # RBAC 权限控制
├── i18n/          # 国际化
├── websocket/     # WebSocket 支持
├── cron/          # 定时任务
├── tracing/       # 链路追踪
├── validator/     # 参数验证
└── logger/        # 日志系统
```

## 配置说明

### RBAC 配置

`rbac_model.conf`:
```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

### 国际化配置

在 `i18n/translations` 目录下创建语言文件：

```json
{
    "welcome.message": "Welcome to EasyGo!",
    "error.unauthorized": "Unauthorized access"
}
```

## 框架对比

### EasyGo vs. Gin, Beego, Echo

#### EasyGo

**优势：**

- **轻量级与高性能**：设计精简，核心代码简洁，启动和运行开销小，路由性能优秀。
- **模块化设计**：各功能模块（如 JWT、RBAC、日志、验证、追踪）高度解耦，易于按需选择和扩展。
- **内置通用功能**：集成了常用功能如 JWT 认证、RBAC 权限管理、参数验证、灵活的日志系统和链路追踪，开箱即用。
- **Context 复用**：利用 `sync.Pool` 复用 `Context` 对象，有效减少内存分配，提升性能。
- **路由分组**：支持路由分组，方便管理和组织复杂的路由结构。
- **崩溃保护**：内置 `Recovery` 中间件，可捕获 `panic`，防止服务器崩溃。
- **完善的请求日志**：提供详细的请求日志中间件，记录请求方法、路径、处理时间、状态码等。
- **灵活的权限管理**：支持从文件、字符串以及多种数据库加载权限策略。
- **易于上手**：API 设计直观，学习曲线平缓。

**劣势：**

- **生态与社区规模**：作为一个相对较新的框架，其社区和生态系统规模不及成熟框架。
- **功能覆盖范围**：更专注于核心 Web 服务，可能需要手动集成更多组件。

#### Gin

**优势：**

- **极致性能**：基于 `httprouter` 构建，拥有非常高的路由性能。
- **轻量与简洁**：API 设计简洁明了，易于理解和使用。
- **丰富的中间件生态**：拥有庞大且活跃的社区，提供大量官方和第三方中间件。
- **错误处理与崩溃恢复**：内置完善的错误处理机制和 panic 恢复功能。

**劣势：**

- **最小化设计**：核心功能相对基础，部分高级功能需要手动集成。
- **文档与学习曲线**：部分文档和概念对于初学者来说可能需要一定时间适应。
- **模板渲染**：对服务器端模板渲染的支持相对有限。

#### Beego

**优势：**

- **全栈框架**：提供一站式 Web 开发解决方案，包括 ORM、缓存、会话管理等。
- **MVC 架构**：强制遵循 MVC 架构模式，有助于保持项目结构清晰。
- **自动化工具**：内置命令行工具，支持代码生成、项目管理等。
- **中文文档丰富**：拥有非常详尽的中文文档和活跃的中文社区。

**劣势：**

- **学习曲线陡峭**：由于功能非常全面，其概念和用法相对复杂。
- **性能开销**：由于集成了大量功能，可能存在一定的性能开销。
- **灵活度**：框架的约定较多，对开发者有一定的约束。

#### Echo

**优势：**

- **高性能与低开销**：以其快速的性能和低内存占用而闻名。
- **简单易用**：API 设计简洁直观，学习曲线平缓。
- **高度可定制**：提供灵活的 API 和强大的中间件支持。
- **内置功能**：支持自动 TLS、多种数据渲染类型等。

**劣势：**

- **中间件生态**：官方中间件库和社区插件数量可能不如 Gin 丰富。
- **相对较新**：相比 Gin，其流行度和社区规模可能略逊一筹。

### 总结

选择合适的 Go Web 框架取决于您的项目需求、团队偏好和性能要求：

- 如果您追求**极致性能和简洁**，且需要强大的中间件生态，**Gin** 是一个很好的选择。
- 如果您需要**全栈解决方案和严格的 MVC 架构**，且不介意较高的学习成本，**Beego** 更适合大型企业级应用。
- 如果您希望在**高性能和易用性之间取得平衡**，并需要高度的**可定制性**，**Echo** 表现出色。
- **EasyGo** 则介于 Gin 和 Echo 之间，旨在提供**轻量、高性能、模块化**的开发体验，并内置了常用功能，特别适合需要快速构建 API 和微服务，并希望拥有一定功能集成的开发者。

## 最佳实践

1. 使用中间件处理通用逻辑
2. 实现优雅关闭
3. 使用参数验证确保数据安全
4. 合理使用路由组组织代码
5. 遵循 RESTful API 设计规范

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 作者

- 作者：xzl
- GitHub：[https://github.com/xzl-go/easygo](https://github.com/xzl-go/easygo)