// Package validator 提供了基于 go-playground/validator 的参数验证功能
// 支持结构体标签验证和自定义验证规则
package validator

import (
	"github.com/go-playground/validator/v10"
)

// validate 是全局验证器实例
var validate *validator.Validate

// init 初始化验证器
func init() {
	validate = validator.New()
}

// Validate 验证结构体
// obj: 要验证的结构体实例
// 返回验证错误（如果有）
func Validate(obj interface{}) error {
	return validate.Struct(obj)
}

// RegisterValidation 注册自定义验证规则
// tag: 验证标签名
// fn: 验证函数
// 返回注册错误（如果有）
func RegisterValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}

// RegisterCustomValidation 注册自定义验证规则（RegisterValidation的别名）
// tag: 验证标签名
// fn: 验证函数
// 返回注册错误（如果有）
func RegisterCustomValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}
