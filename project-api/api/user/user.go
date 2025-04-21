package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"net/http"
	"project-api/pkg/model/user"
	common "project-common"
	"project-common/errs"
	"project-grpc/user/login"
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
	zap.L().Info("接收电话号码：" + mobile)
	rsp, err := LoginServiceClient.GetCaptcha(context.Background(), &login.CaptchaMessage{Mobile: mobile})
	if err != nil {
		grpcError, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(grpcError, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.Code))
	zap.L().Info("发送验证码成功")
}

func (u *HandlerUser) register(c *gin.Context) {
	//	1.接收参数 参数模型
	result := &common.Result{}
	var req user.RegisterReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式错误"))
		return
	}
	//	2.参数校验 判断参数是否合法
	if err := req.Verify(); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	//	3.调用user grpc服务 获取响应
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msg := &login.RegisterMessage{}
	err = copier.Copy(msg, req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusInternalServerError, "参数转换失败"))
		return
	}
	_, err = LoginServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//  4.返回响应
	c.JSON(http.StatusOK, result.Success("注册成功"))
}
