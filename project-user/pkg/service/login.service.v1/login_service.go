package login_service_v1

import (
	"context"
	"go.uber.org/zap"
	common "project-common"
	"project-common/errs"
	"project-grpc/user/login"
	"project-user/internal/dao"
	"project-user/internal/repo"
	"project-user/pkg/model"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache      repo.Cache
	memberrepo repo.MemberRepo
}

func New() *LoginService {
	return &LoginService{
		cache:      dao.Rc,
		memberrepo: dao.NewMemberDao(),
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	zap.L().Info("调用验证码服务")
	//1. 获取参数
	mobile := msg.Mobile
	//2. 验证手机合法性
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	//3.生成验证码
	code := "123456"
	//4. 发送验证码
	go func() {
		time.Sleep(2 * time.Second)
		//发送成功 存入redis
		timeControl, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := ls.cache.Put(timeControl, model.RegisterRedisKey+mobile, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("存储验证码失败")
		}
		zap.L().Info("存储验证码成功")
	}()
	return &login.CaptchaResponse{Code: code}, nil
}
func (ls *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	//  1.可以校验参数
	//	2.校验验证码
	redisCode, err := ls.cache.Get(ctx, model.RegisterRedisKey+msg.Mobile)
	if err != nil {
		zap.L().Error("获取验证码失败", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	//  3.校验业务逻辑（邮箱，手机号，账号是否被注册）
	exist, err := ls.memberrepo.GetMemberByEmail(ctx, msg.Email)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(model.DBerror))
		return nil, err
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	exist, err = ls.memberrepo.GetMemberByAccount(ctx, msg.Name)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(model.DBerror))
		return nil, err
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}
	exist, err = ls.memberrepo.GetMemberByMobile(ctx, msg.Mobile)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(model.DBerror))
		return nil, err
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist)
	}

	//  4.执行业务，将数据存入member表，生成数据，存入组织表organization
	//  5.返回
	return nil, nil
}
