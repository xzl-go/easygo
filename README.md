# EasyGo - 轻量级 Go Web 框架

## 简介

EasyGo 是一个功能丰富、轻量级的 Go Web 框架，旨在提供简单易用且功能强大的 Web 应用开发体验。受 Gin 启发，但更加模块化和可扩展。

## 特性

- 🚀 **高性能路由**：支持动态路径参数和通配符
- 🔒 **参数验证**：基于 `go-playground/validator` 的强大验证
- 🔐 **JWT 认证**：内置 JWT 令牌生成、验证和刷新
- 🛡️ **权限控制**：基于 Casbin 的细粒度 RBAC 权限管理
- 📝 **日志系统**：灵活的多级日志，支持控制台和文件输出
- 🔍 **链路追踪**：基于 OpenTelemetry 的追踪功能
- 🧩 **中间件支持**：轻松添加全局中间件

## 安装

```bash
go get github.com/yourusername/easygo
```

## 快速开始

### 基本路由

```go
package main

import "easygo/core"

func main() {
    app := core.New()

    app.GET("/hello", func(ctx *core.Context) {
        ctx.JSON(200, map[string]string{
            "message": "Hello, EasyGo!",
        })
    })

    app.Run(":8080")
}
```

### 参数验证

```go
type User struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"gte=18,lte=120"`
}

app.POST("/register", func(ctx *core.Context) {
    var user User
    if err := ctx.BindJSON(&user); err != nil {
        ctx.JSON(400, map[string]string{"error": err.Error()})
        return
    }

    if err := validator.Validate(user); err != nil {
        ctx.JSON(400, map[string]string{"error": err.Error()})
        return
    }

    // 处理注册逻辑
})
```

### 中间件

```go
// 日志中间件
app.Use(func(ctx *core.Context) {
    log.Info("请求路径：%s", ctx.Path)
})

// JWT 认证中间件
app.Use(func(ctx *core.Context) {
    token := ctx.Req.Header.Get("Authorization")
    claims, err := jwtManager.VerifyToken(token)
    if err != nil {
        ctx.JSON(401, map[string]string{"error": "未授权"})
        ctx.Abort()
        return
    }
})
```

### JWT 认证

```go
// 生成 Token
token, err := jwtManager.GenerateToken(userID, username)

// 验证 Token
claims, err := jwtManager.VerifyToken(tokenString)
```

### RBAC 权限控制

```go
// 初始化 RBAC 管理器
rbacManager, _ := rbac.NewRBACManager("rbac_model.conf", "rbac_policy.csv")

// 检查权限
allowed, _ := rbacManager.Enforce(username, "/profile", "GET")
if !allowed {
    ctx.JSON(403, map[string]string{"error": "无权访问"})
}
```

## 模块详解

### 核心模块

- `core`：路由和请求处理
- `validator`：参数验证
- `logger`：日志系统
- `jwt`：JWT 认证
- `rbac`：权限控制
- `tracing`：链路追踪

## 配置

### RBAC 模型配置 (`rbac_model.conf`)

```ini
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

## 性能与扩展性

- 使用 `sync.Pool` 复用 Context，减少内存分配
- 模块化设计，易于扩展和定制
- 低依赖，核心代码简洁

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

[MIT License](LICENSE)

## 联系

- 作者：[您的名字]
- 邮箱：[您的邮箱]
- 项目地址：[项目 GitHub 地址] 