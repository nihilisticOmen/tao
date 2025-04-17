package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	common "project-common"
	"project-common/errs"
	loginServiceV1 "project-user/pkg/service/login.service.v1"
	"time"
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
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rsp, err := LoginServiceClient.GetCaptcha(c, &loginServiceV1.CaptchaMessage{Mobile: mobile})
	if err != nil {
		grpcError, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(grpcError, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.Code))
}
