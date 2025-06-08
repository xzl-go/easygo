package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xzl-go/easygo/core"
	"github.com/xzl-go/easygo/jwt"
	"github.com/xzl-go/easygo/logger"
	"github.com/xzl-go/easygo/rbac"
	"github.com/xzl-go/easygo/tracing"
	"github.com/xzl-go/easygo/validator"
)

type User struct {
	ID       string `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"gte=18,lte=120"`
}

func main() {
	// 初始化日志
	log := logger.New(logger.INFO, "app.log")
	defer log.Close()

	// 初始化链路追踪
	tracer := tracing.NewTracer("user-service")
	defer tracer.Shutdown(context.Background())

	// 初始化JWT
	jwtManager := jwt.NewJWTManager("your_secret_key", 24*time.Hour)

	// 初始化RBAC
	rbacManager, err := rbac.NewRBACManager("rbac_model.conf", "rbac_policy.csv")
	if err != nil {
		log.Fatal("RBAC初始化失败：%v", err)
	}

	// 创建引擎
	app := core.New()

	// 注册中间件
	app.Use(func(ctx *core.Context) {
		log.Info("请求路径：%s", ctx.Path)
	})

	// 用户注册
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

		// 生成JWT
		token, err := jwtManager.GenerateToken(user.ID, user.Username)
		if err != nil {
			ctx.JSON(500, map[string]string{"error": "Token生成失败"})
			return
		}

		ctx.JSON(200, map[string]string{
			"message": "注册成功",
			"token":   token,
		})
	})

	// 用户登录
	app.POST("/login", func(ctx *core.Context) {
		var loginUser struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		if err := ctx.BindJSON(&loginUser); err != nil {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return
		}

		// 模拟用户验证
		if loginUser.Username == "admin" && loginUser.Password == "admin123" {
			token, _ := jwtManager.GenerateToken("1", loginUser.Username)
			ctx.JSON(200, map[string]string{
				"message": "登录成功",
				"token":   token,
			})
		} else {
			ctx.JSON(401, map[string]string{"error": "用户名或密码错误"})
		}
	})

	// 需要认证的受保护路由
	app.GET("/profile", func(ctx *core.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, map[string]string{"error": "未提供认证信息"})
			return
		}

		claims, err := jwtManager.VerifyToken(authHeader)
		if err != nil {
			ctx.JSON(401, map[string]string{"error": "无效的令牌"})
			return
		}

		// 权限检查
		allowed, err := rbacManager.Enforce(claims.Username, "/profile", "GET")
		if err != nil || !allowed {
			ctx.JSON(403, map[string]string{"error": "无权访问"})
			return
		}

		ctx.JSON(200, map[string]string{
			"message": fmt.Sprintf("欢迎，%s！", claims.Username),
		})
	})

	// 启动服务器
	log.Info("服务器启动...")
	if err := app.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败：%v", err)
	}
}
