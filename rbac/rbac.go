// Package rbac 提供了基于 Casbin 的 RBAC（基于角色的访问控制）功能
// 支持细粒度的权限管理和策略控制
package rbac

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	stringadapter "github.com/qiangmzsx/string-adapter/v2"

	// 导入 GORM 适配器和所需的数据库驱动
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"    // MySQL 驱动
	"gorm.io/driver/postgres" // PostgreSQL 驱动
	"gorm.io/driver/sqlite"   // SQLite 驱动
	"gorm.io/gorm"            // GORM 核心库
)

// RBACManager 是RBAC权限管理器
// 负责权限策略的管理和执行
type RBACManager struct {
	enforcer *casbin.Enforcer // Casbin执行器
}

// NewRBACManager 创建一个新的RBAC权限管理器 (从文件加载模型和策略)
// modelPath: RBAC模型配置文件路径
// policyPath: 权限策略文件路径
// 返回RBAC管理器实例和可能的错误
func NewRBACManager(modelPath, policyPath string) (*RBACManager, error) {
	// 从文件加载RBAC模型
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model from file: %w", err)
	}

	// 创建文件适配器，用于持久化权限策略
	adapter := fileadapter.NewAdapter(policyPath)

	return NewRBACManagerWithAdapter(m, adapter)
}

// NewRBACManagerFromStrings 创建一个新的RBAC权限管理器 (从字符串加载模型和策略)
// modelContent: RBAC模型内容的字符串
// policyContent: 权限策略内容的CSV字符串
// 返回RBAC管理器实例和可能的错误
func NewRBACManagerFromStrings(modelContent, policyContent string) (*RBACManager, error) {
	// 从字符串加载RBAC模型
	m, err := model.NewModelFromString(modelContent)
	if err != nil {
		return nil, fmt.Errorf("failed to load model from string: %w", err)
	}

	// 创建字符串适配器
	adapter := stringadapter.NewAdapter(policyContent)

	return NewRBACManagerWithAdapter(m, adapter)
}

// NewRBACManagerFromDB 创建一个新的RBAC权限管理器 (从数据库加载策略)
// driverName: 数据库驱动名称 (例如 "mysql", "postgres", "sqlite3")
// dataSourceName: 数据库连接字符串
// modelContent: RBAC模型内容的字符串
// 返回RBAC管理器实例和可能的错误
func NewRBACManagerFromDB(driverName, dataSourceName, modelContent string) (*RBACManager, error) {
	// 从字符串加载RBAC模型
	m, err := model.NewModelFromString(modelContent)
	if err != nil {
		return nil, fmt.Errorf("failed to load model from string for DB: %w", err)
	}

	var db *gorm.DB
	switch driverName {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driverName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 创建 GORM 适配器
	adapter, err := gormadapter.NewAdapterByDB(db) // 使用 NewAdapterByDB
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm adapter: %w", err)
	}

	return NewRBACManagerWithAdapter(m, adapter)
}

// NewRBACManagerWithAdapter 创建一个新的RBAC权限管理器 (使用自定义适配器)
// m: Casbin模型实例
// adapter: Casbin适配器实例 (可以是文件适配器、数据库适配器或其他类型)
// 返回RBAC管理器实例和可能的错误
func NewRBACManagerWithAdapter(m model.Model, adapter persist.Adapter) (*RBACManager, error) {
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 从适配器加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return &RBACManager{
		enforcer: enforcer,
	}, nil
}

// AddPolicy 添加权限策略
// sec: 策略类型
// ptype: 策略类型
// rule: 策略规则
// 返回操作结果和可能的错误
func (r *RBACManager) AddPolicy(sec, ptype string, rule []string) (bool, error) {
	return r.enforcer.AddPolicy(sec, ptype, rule)
}

// RemovePolicy 删除权限策略
// sec: 策略类型
// ptype: 策略类型
// rule: 策略规则
// 返回操作结果和可能的错误
func (r *RBACManager) RemovePolicy(sec, ptype string, rule []string) (bool, error) {
	return r.enforcer.RemovePolicy(sec, ptype, rule)
}

// AddRoleForUser 为用户添加角色
// user: 用户名
// role: 角色名
// 返回操作结果和可能的错误
func (r *RBACManager) AddRoleForUser(user, role string) (bool, error) {
	return r.enforcer.AddRoleForUser(user, role)
}

// DeleteRoleForUser 删除用户的角色
// user: 用户名
// role: 角色名
// 返回操作结果和可能的错误
func (r *RBACManager) DeleteRoleForUser(user, role string) (bool, error) {
	return r.enforcer.DeleteRoleForUser(user, role)
}

// Enforce 执行权限检查
// sub: 主体（用户）
// obj: 对象（资源）
// act: 操作（动作）
// 返回是否允许访问和可能的错误
func (r *RBACManager) Enforce(sub, obj, act string) (bool, error) {
	return r.enforcer.Enforce(sub, obj, act)
}

// GetRolesForUser 获取用户的所有角色
// user: 用户名
// 返回角色列表和可能的错误
func (r *RBACManager) GetRolesForUser(user string) ([]string, error) {
	return r.enforcer.GetRolesForUser(user)
}

// GetPermissionsForUser 获取用户的所有权限
// user: 用户名
// 返回权限列表和可能的错误
func (r *RBACManager) GetPermissionsForUser(user string) ([][]string, error) {
	return r.enforcer.GetPermissionsForUser(user)
}

// LoadPolicy 从存储加载权限策略
// 返回可能的错误
func (r *RBACManager) LoadPolicy() error {
	return r.enforcer.LoadPolicy()
}

// SavePolicy 保存权限策略到存储
// 返回可能的错误
func (r *RBACManager) SavePolicy() error {
	return r.enforcer.SavePolicy()
}
