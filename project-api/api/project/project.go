package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"project-api/pkg/model"
	"project-api/pkg/model/menu"
	"project-api/pkg/model/pro"
	common "project-common"
	"project-common/errs"
	"project-grpc/project"
	"strconv"
	"time"
)

type HandlerProject struct {
}

func (p *HandlerProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}
	indexResponse, err := ProjectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	menus := indexResponse.Menus
	var ms []*menu.Menu
	copier.Copy(&ms, menus)
	c.JSON(http.StatusOK, result.Success(ms))
}

func (p *HandlerProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	//1. 获取参数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c)
	selectBy := c.PostForm("selectBy")
	msg := &project.ProjectRpcMessage{
		MemberId:   memberId,
		MemberName: memberName,
		SelectBy:   selectBy,
		Page:       page.Page,
		PageSize:   page.PageSize}
	myProjectResponse, err := ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var pms []*pro.ProjectAndMember
	copier.Copy(&pms, myProjectResponse.Pm)
	if pms == nil {
		pms = []*pro.ProjectAndMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms, //null nil -> []
		"total": myProjectResponse.Total,
	}))
}

func (p *HandlerProject) projectTemplate(c *gin.Context) {
	result := &common.Result{}
	//1. 获取参数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c)
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		MemberName:       memberName,
		ViewType:         int32(viewType),
		Page:             page.Page,
		PageSize:         page.PageSize,
		OrganizationCode: c.GetString("organizationCode")}
	templateResponse, err := ProjectServiceClient.FindProjectTemplate(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var pms []*pro.ProjectTemplate
	copier.Copy(&pms, templateResponse.Ptm)
	if pms == nil {
		pms = []*pro.ProjectTemplate{}
	}
	for _, v := range pms {
		if v.TaskStages == nil {
			v.TaskStages = []*pro.TaskStagesOnlyName{}
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms, //null nil -> []
		"total": templateResponse.Total,
	}))
}

func (p *HandlerProject) projectSave(c *gin.Context) {
	result := &common.Result{}
	//1. 获取参数
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	organizationCode := c.GetString("organizationCode")
	var req *pro.SaveProjectRequest
	c.ShouldBind(&req)
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		OrganizationCode: organizationCode,
		TemplateCode:     req.TemplateCode,
		Name:             req.Name,
		Id:               int64(req.Id),
		Description:      req.Description,
	}
	saveProject, err := ProjectServiceClient.SaveProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var rsp *pro.SaveProject
	copier.Copy(&rsp, saveProject)
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (p *HandlerProject) readProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detail, err := ProjectServiceClient.FindProjectDetail(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, MemberId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	pd := &pro.ProjectDetail{}
	copier.Copy(pd, detail)
	c.JSON(http.StatusOK, result.Success(pd))
}

func (p *HandlerProject) recycleProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectServiceClient.UpdateDeletedProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, Deleted: true})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) recoveryProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectServiceClient.UpdateDeletedProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, Deleted: false})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) collectProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	collectType := c.PostForm("type")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectServiceClient.UpdateCollectProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, CollectType: collectType, MemberId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) editProject(c *gin.Context) {
	result := &common.Result{}
	var req *pro.ProjectReq
	_ = c.ShouldBind(&req)
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.UpdateProjectMessage{}
	copier.Copy(msg, req)
	msg.MemberId = memberId
	_, err := ProjectServiceClient.UpdateProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func New() *HandlerProject {
	return &HandlerProject{}
}
