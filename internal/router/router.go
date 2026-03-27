package router

import (
	"github.com/islamchupanov/tz1/internal/handler"

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

	return r
}
