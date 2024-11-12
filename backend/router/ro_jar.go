package router

import (
	v1 "github.com/1Panel-dev/1Panel/backend/app/api/v1"
	"github.com/1Panel-dev/1Panel/backend/middleware"
	"github.com/gin-gonic/gin"
)

type JarRouter struct{}

func (s *JarRouter) InitRouter(Router *gin.RouterGroup) {
	jarRouter := Router.Group("jar").
		Use(middleware.JwtAuth()).
		Use(middleware.SessionAuth()).
		Use(middleware.PasswordExpired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		jarRouter.POST("/upload", baseApi.UploadJar)
		jarRouter.GET("/status", baseApi.GetJarStatus)
		jarRouter.PUT("/start", baseApi.StartJar)
		jarRouter.PUT("/stop", baseApi.StopJar)
		jarRouter.DELETE("/delete", baseApi.DeleteJar)
	}
}
