package model

import (
	"project-common/errs"
)

var (
	RedisError    = errs.NewError(999, "redis错误")     // 验证码存储失败
	DBerror       = errs.NewError(998, "db错误")        // 数据库错误
	NoLegalMobile = errs.NewError(10102001, "手机号不合法") // 手机号不合法
	CaptchaError  = errs.NewError(10102002, "验证码不合法") // 验证码不合法
	EmailExist    = errs.NewError(10102003, "邮箱已存在")  // 邮箱已存在
	AccountExist  = errs.NewError(10102004, "用户名已存在") // 手机号已存在
	MobileExist   = errs.NewError(10102005, "手机号已存在") // 手机号已存在

)
