package controller

import (
	"WebBlog/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// ------------------- 社区相关 -------------------

// CommunityHandler 社区列表
func CommunityHandler(ctx *gin.Context) {
	// 查询到所有的社区（community_id, community_name）以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy) // 不轻易暴露系统内部错误
		return
	}
	ResponseSuccess(ctx, data)
}

// CommunityDetailHandler 社区详情
func CommunityDetailHandler(ctx *gin.Context) {
	// 获取社区ID
	idStr := ctx.Param("id") // 从请求URL中获取参数
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 根据id查询到社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail(id) failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy) // 不轻易暴露系统内部错误
		return
	}
	ResponseSuccess(ctx, data)
}
