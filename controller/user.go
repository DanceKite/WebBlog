package controller

import (
	"WebBlog/dao/mysql"
	"WebBlog/logic"
	"WebBlog/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignUpHandler(ctx *gin.Context) {
	// 1. 获取参数和参数校验
	//var p models.ParamSignUp
	p := new(models.ParamSignUp)
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//手动对参数进行详细的业务规则校验  //有数据库做校验的话，这里就不需要了
	/*if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
		zap.L().Error("SignUp with invalid param")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求参数错误",
		})
		return
	}*/

	fmt.Println(*p)
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(ctx, CodeUserExist)
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(ctx, nil)
}

func LoginHandler(ctx *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamLogin)
	if err := ctx.ShouldBindJSON(p); err != nil {
		//请求参数有误直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserNotExist)
		}
		ResponseError(ctx, CodeInvalidPassword)
		return
	}
	// 3. 返回响应
	ResponseSuccess(ctx, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), //id值大于1<<53-1时，json序列化会丢失精度
		"user_name": user.Username,
		"token":     user.Token,
	})
}
