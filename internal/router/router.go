package router

import (
	"github.com/islamchupanov/tz1/docs"
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

	// Инициализация Swagger docs
	docs.SwaggerInfo.Title = "Device API"
	docs.SwaggerInfo.Description = "API for managing devices"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
