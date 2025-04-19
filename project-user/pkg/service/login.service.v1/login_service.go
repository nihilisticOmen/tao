package login_service_v1

import (
	"context"
	"go.uber.org/zap"
	common "project-common"
	"project-common/errs"
	"project-user/pkg/dao"
	"project-user/pkg/model"
	"project-user/pkg/repo"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	Cache repo.Cache
}

func New() *LoginService {
	return &LoginService{
		Cache: dao.Rc,
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *CaptchaMessage) (*CaptchaResponse, error) {
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
		err := ls.Cache.Put(timeControl, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("存储验证码失败")
		}
		zap.L().Info("存储验证码成功")
	}()
	return &CaptchaResponse{Code: code}, nil
}
