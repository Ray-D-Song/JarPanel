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
		jarRouter.POST("/new", baseApi.NewService)
		jarRouter.GET("/status", baseApi.GetServiceStatus)
		jarRouter.PUT("/start", baseApi.StartService)
		jarRouter.PUT("/stop", baseApi.StopService)
		jarRouter.GET("/files", baseApi.GetServiceFileList)
		jarRouter.GET("/download", baseApi.DownloadFile)
		jarRouter.POST("/upload", baseApi.UploadFile)
	}
}
