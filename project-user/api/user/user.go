package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	common "project-common"
	"project-common/logs"
	"project-user/pkg/dao"
	"project-user/pkg/model"
	"project-user/pkg/repo"
	"time"
)

type HandlerUser struct {
	cache repo.Cache
}

func New() *HandlerUser {
	return &HandlerUser{
		cache: dao.Rc,
	}
}
func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	rsp := &common.Result{}
	//1.获取参数
	//mobile := ctx.PostForm("mobile")
	mobile := ctx.PostForm("mobile")
	//2.参数校验
	if !common.VerifyMobile(mobile) {
		ctx.JSON(http.StatusOK, rsp.Fail(model.NoLegalMobile, "手机号不合法"))
		return
	}
	//3.生成验证码(随机4位或者6位数字)
	code := "123456"
	//4.调用短信平台
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("调用短信平台发送短信 info")
		logs.LG.Debug("调用短信平台发送短信 debug")
		zap.L().Error("调用短信平台发送短信 error")
		//5.存储验证码到redis,设置过期时间15min
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := dao.Rc.Put(c, mobile, code, 15*time.Minute)
		if err != nil {
			log.Println("存储验证码失败:", err)
			return
		}
	}()

	ctx.JSON(http.StatusOK, rsp.Success("123456"))
}
