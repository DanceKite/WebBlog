package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code": 1000, //程序中的错误码
	"msg": "xxx", //提示信息
	"data": {} //返回数据
}
*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(ctx *gin.Context, code ResCode) {

	ctx.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseErrorWithMsg(ctx *gin.Context, code ResCode, msg interface{}) {

	ctx.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
