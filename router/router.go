package router

import (
	"WebBlog/controller"
	"WebBlog/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	return r
}
