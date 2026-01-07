package router

import (
	"ygsx-v2/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.TenantMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		tenantCode, _ := c.Get(middleware.TenantCodeKey)

		c.JSON(200, gin.H{
			"message":    "连接正常",
			"tenantCode": tenantCode,
		})
	})

	return r
}
