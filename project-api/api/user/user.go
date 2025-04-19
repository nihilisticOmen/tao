package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	common "project-common"
	"project-common/errs"
	loginServiceV1 "project-user/pkg/service/login.service.v1"
)

type HandlerUser struct {
}

func New() *HandlerUser {
	return &HandlerUser{}
}
func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	//1.获取参数
	//mobile := ctx.PostForm("mobile")
	mobile := ctx.PostForm("mobile")
	zap.L().Info("接收电话号码：" + mobile)
	rsp, err := LoginServiceClient.GetCaptcha(context.Background(), &loginServiceV1.CaptchaMessage{Mobile: mobile})
	if err != nil {
		grpcError, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(grpcError, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.Code))
	zap.L().Info("发送验证码成功")
}
