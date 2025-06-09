// Package main 是EasyGo框架的示例应用
// 展示了框架的主要功能，包括用户注册、登录、认证和权限控制
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xzl-go/easygo/core"
	"github.com/xzl-go/easygo/i18n"
	"github.com/xzl-go/easygo/jwt"
	"github.com/xzl-go/easygo/logger"
	"github.com/xzl-go/easygo/middleware"
	"github.com/xzl-go/easygo/rbac"
	"github.com/xzl-go/easygo/tracing"
	"github.com/xzl-go/easygo/validator"
	"github.com/xzl-go/easygo/websocket"
)

// User 结构体定义了用户的基本信息
// 包含了用户ID、用户名、邮箱和年龄等字段
// 每个字段都添加了验证标签，用于参数验证
type User struct {
	ID       string `json:"id" validate:"required"`                    // 用户唯一标识
	Username string `json:"username" validate:"required,min=3,max=20"` // 用户名，长度3-20
	Email    string `json:"email" validate:"required,email"`           // 邮箱地址
	Age      int    `json:"age" validate:"gte=18,lte=120"`             // 年龄，18-120岁
}

// @title EasyGo API
// @version 1.0
// @description EasyGo 框架示例应用
// @host localhost:8080
// @BasePath /
func main() {
	// 初始化日志系统
	logger.Init()

	// 初始化链路追踪系统，用于分布式追踪
	tracer := tracing.NewTracer("user-service")
	defer tracer.Shutdown(context.Background())

	// 初始化JWT管理器，设置密钥和token过期时间
	jwtManager := jwt.NewJWTManager("your_secret_key", 24*time.Hour)

	// 初始化RBAC权限管理器，加载权限模型和策略
	rbacManager, err := rbac.NewRBACManager("rbac_model.conf", "rbac_policy.csv")
	if err != nil {
		logger.Error("RBAC初始化失败：%v", err)
		return
	}

	// 初始化国际化
	i18nManager := i18n.New("en")
	if err := i18nManager.LoadTranslations("i18n/translations"); err != nil {
		logger.Error("Failed to load translations: %v", err)
		return
	}

	// 初始化定时任务
	//cron.InitCron()
	//defer cron.StopCron()

	// 添加示例定时任务
	//cron.AddJob("@every 1m", func() {
	//	logger.Info("定时任务执行：%v", time.Now())
	//})

	// 创建Web应用引擎
	app := core.New()

	// 应用 Recovery 中间件，用于捕获 panic 并防止服务器崩溃
	app.Use(middleware.Recovery())

	// 注册全局中间件，用于记录请求日志
	app.Use(middleware.Logger())

	// 注册国际化中间件
	app.Use(i18nManager.Middleware())

	// 注册用户路由处理函数
	app.POST("/register", func(ctx *core.Context) {
		var user User
		// 解析JSON请求体到User结构体
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return
		}

		// 验证用户数据
		if err := validator.Validate(user); err != nil {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return
		}

		// 生成JWT令牌
		token, err := jwtManager.GenerateToken(user.ID, user.Username)
		if err != nil {
			ctx.JSON(500, map[string]string{"error": "Token生成失败"})
			return
		}

		lang := ctx.Get("lang").(string)
		message := i18nManager.Translate("welcome.message", lang)
		ctx.JSON(200, map[string]string{
			"message": message,
			"token":   token,
		})
	})

	// 用户登录路由处理函数
	app.POST("/login", func(ctx *core.Context) {
		var loginUser struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		// 解析登录请求
		if err := ctx.BindJSON(&loginUser); err != nil {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return
		}

		// 模拟用户验证（实际应用中应该查询数据库）
		if loginUser.Username == "admin" && loginUser.Password == "admin123" {
			token, _ := jwtManager.GenerateToken("1", loginUser.Username)
			lang := ctx.Get("lang").(string)
			message := i18nManager.Translate("welcome.message", lang)
			ctx.JSON(200, map[string]string{
				"message": message,
				"token":   token,
			})
		} else {
			ctx.JSON(401, map[string]string{"error": i18nManager.Translate("error.unauthorized", ctx.Get("lang").(string))})
		}
	})

	// 受保护的用户资料路由，需要认证和权限验证
	app.GET("/profile", func(ctx *core.Context) {
		// 获取认证头信息
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			ctx.JSON(401, map[string]string{"error": i18nManager.Translate("error.unauthorized", ctx.Get("lang").(string))})
			return
		}

		// 验证JWT令牌
		claims, err := jwtManager.VerifyToken(authHeader)
		if err != nil {
			ctx.JSON(401, map[string]string{"error": i18nManager.Translate("error.unauthorized", ctx.Get("lang").(string))})
			return
		}

		// 检查用户权限
		allowed, err := rbacManager.Enforce(claims.Username, "/profile", "GET")
		if err != nil || !allowed {
			ctx.JSON(403, map[string]string{"error": i18nManager.Translate("error.forbidden", ctx.Get("lang").(string))})
			return
		}

		lang := ctx.Get("lang").(string)
		message := fmt.Sprintf("欢迎，%s！", claims.Username)
		translatedMessage := i18nManager.Translate(message, lang)
		ctx.JSON(200, map[string]string{
			"message": translatedMessage,
		})
	})

	// 添加一个会触发 panic 的路由，用于测试 Recovery 中间件
	app.GET("/panic", func(ctx *core.Context) {
		panic("这是一个测试 panic！")
	})

	// WebSocket路由
	app.GET("/ws", func(ctx *core.Context) {
		websocket.HandleWebSocket(ctx)
	})

	// 启动Web服务器
	if err := app.Run(":8080"); err != nil {
		logger.Error("Failed to start server: %v", err)
		return
	}
}
