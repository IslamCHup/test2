package router

import (
	"github.com/islamchupanov/tz1/internal/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func SetupRouter(deviceHandler *handler.DeviceHandler) *gin.Engine {
	r := gin.Default()

	devices := r.Group("/devices")
	{
		devices.POST("", deviceHandler.CreateDevice)
		devices.GET("", deviceHandler.ListDevices)
		devices.GET("/:id", deviceHandler.GetDevice)
		devices.PUT("/:id", deviceHandler.UpdateDevice)
		devices.DELETE("/:id", deviceHandler.DeleteDevice)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
