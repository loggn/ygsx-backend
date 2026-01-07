package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	TenantCodeHeader = "X-Tenant-Code"
	TenantCodeKey    = "tenant_code"
	TenantIDKey      = "tenant_id"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCode := c.GetHeader(TenantCodeHeader)
		if tenantCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing tenant code",
			})
			c.Abort()
			return
		}

		mockTenantID := uint(1)

		c.Set(TenantCodeKey, tenantCode)
		c.Set(TenantIDKey, mockTenantID)

		c.Next()
	}
}
