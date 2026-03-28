package router

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

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

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Device API"
	docs.SwaggerInfo.Description = "API for managing devices"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static frontend files in production
	// Check if dist folder exists (for development, serve from filesystem)
	distPath := "./frontend/dist"
	if _, err := os.Stat(distPath); err == nil {
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if path == "/" || path == "" {
				path = "/index.html"
			}

			// Try to serve the requested file from dist folder
			filePath := filepath.Join(distPath, filepath.Clean(path))
			
			// Check if file exists
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				c.File(filePath)
				return
			}

			// For SPA routing, serve index.html for unknown routes
			c.File(filepath.Join(distPath, "index.html"))
		})
	}

	return r
}
