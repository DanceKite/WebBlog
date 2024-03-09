package router

import (
	"WebBlog/controller"
	"WebBlog/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	r.POST("/login", controller.LoginHandler)

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	return r
}
