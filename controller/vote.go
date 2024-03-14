package controller

import (
	"WebBlog/logic"
	"WebBlog/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func PostVoteHandler(ctx *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.ParamVoteData)
	if err := ctx.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // 翻译并去除掉错误提示中的结构体名称
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errData)
		return
	}

	// 获取当前用户id
	userID, err := getCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}

	// 2.投票
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(ctx, nil)
}
