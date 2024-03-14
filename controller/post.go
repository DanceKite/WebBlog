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

// GetPostListHandler2 获取帖子列表升级版
// 根据前端传来的参数动态获得帖子列表（按分数、按时间）
// 1.获取分页参数
// 2.去redis查询获取id列表
// 3.根据id去mysql查询帖子的详细信息
func GetPostListHandler2(ctx *gin.Context) {
	// GET请求参数：page size order
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := ctx.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 2.获取数据
	data, err := logic.GetPostListNew(p) // 合并接口
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(ctx, data)
}

/*// GetCommunityPostListHandler 获取社区帖子列表
func GetCommunityPostListHandler(ctx *gin.Context) {
	// GET请求参数：page size order
	p := &models.ParamCommunityPostList{
		ParamPostList: models.ParamPostList{
			Page:  1,
			Size:  10,
			Order: models.OrderTime,
		},
	}

	if err := ctx.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	// 2.获取数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityPostListHandler failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(ctx, data)
}
*/
