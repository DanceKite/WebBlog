package controller

import (
	"WebBlog/logic"
	"WebBlog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(ctx *gin.Context) {
	// 1.获取参数及参数校验
	p := new(models.Post)
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Debug("ctx.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("CreatePost with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	// 从ctx中获取当前发请求的用户的ID
	userID, err := getCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2.业务处理,创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(ctx, nil)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(ctx *gin.Context) {
	// 1.获取参数(从URL中获取帖子ID)
	pidStr := ctx.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 2.业务处理(根据帖子ID查询帖子详情)
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(ctx, data)
}

// GetPostListHandler 获取帖子列表
func GetPostListHandler(ctx *gin.Context) {
	// 1.获取分页参数
	page, size := GetPageInfo(ctx)
	// 2.获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(ctx, data)
}
