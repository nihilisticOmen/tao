package login_service_v1

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	common "project-common"
	"project-common/encrypts"
	"project-common/errs"
	"project-common/jwts"
	"project-grpc/user/login"
	"project-user/config"
	"project-user/internal/dao"
	"project-user/internal/data/member"
	"project-user/internal/data/organization"
	"project-user/internal/database"
	"project-user/internal/database/tran"
	"project-user/internal/repo"
	"project-user/pkg/model"
	"strconv"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
	transaction      tran.Transaction
}

func New() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganizationDao(),
		transaction:      dao.NewTransaction(),
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
	if errors.Is(err, redis.Nil) {
		return nil, errs.GrpcError(model.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("获取验证码失败", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	//  3.校验业务逻辑（邮箱，手机号，账号是否被注册）
	exist, err := ls.memberRepo.GetMemberByEmail(ctx, msg.Email)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(model.DBerror))
		return nil, err
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	exist, err = ls.memberRepo.GetMemberByAccount(ctx, msg.Name)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(model.DBerror))
		return nil, err
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}
	exist, err = ls.memberRepo.GetMemberByMobile(ctx, msg.Mobile)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(model.DBerror))
		return nil, err
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist)
	}
	pwd := encrypts.Md5(msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      pwd,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	org := &organization.Organization{
		Name:       mem.Name + "个人组织",
		MemberId:   mem.Id,
		CreateTime: time.Now().UnixMilli(),
		Personal:   model.Personal,
		Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
	}
	err = ls.transaction.Action(func(conn database.DbConn) error {
		if err := ls.memberRepo.SaveMember(conn, ctx, mem); err != nil {
			zap.L().Error("register save member db err", zap.Error(err))
			return errs.GrpcError(model.DBerror)
		}
		//  4.执行业务，将数据存入member表，生成数据，存入组织表organization
		err = ls.organizationRepo.SaveOrganization(conn, ctx, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(model.DBerror)
		}
		return nil
	})

	//  5.返回
	return &login.RegisterResponse{}, err
}
func (ls *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	c := context.Background()
	//	查询数据库，账号密码
	mem, err := ls.memberRepo.FindMember(c, msg.Account, msg.Password)
	if err != nil {
		zap.L().Error("Login db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBerror)
	}
	if mem == nil {
		return nil, errs.GrpcError(model.AccountOrPwdError)
	}
	memMsg := &login.MemberMessage{}
	err = copier.Copy(memMsg, mem)
	//查询数据库,根据id查组织
	orgs, err := ls.organizationRepo.FindOrganizationByMemId(c, mem.Id)
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	//使用jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	exp := time.Duration(config.AppConf.JwtConfig.AccessExp*3600*24) * time.Second
	rExp := time.Duration(config.AppConf.JwtConfig.RefreshExp*3600*24) * time.Second
	token := jwts.CreateToken(memIdStr, exp, config.AppConf.JwtConfig.AccessSecret, rExp, config.AppConf.JwtConfig.RefreshSecret)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		AccessTokenExp: token.AccessExp,
		TokenType:      "bearer",
	}
	return &login.LoginResponse{
		Member:           memMsg,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}
