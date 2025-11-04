package middleware

import (
	"net/http"
	"strings"
	"ygsx-backend/database"
	"ygsx-backend/model"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未提供 Token"})
			c.Abort()
			return
		}

		// 去掉 Bearer 前缀
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 根据 token 查 Redis（你登录时存过的）
		result, err := database.Rdb.Keys(database.Ctx, "user:token:*").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Redis 连接错误"})
			c.Abort()
			return
		}

		// 遍历找到匹配的 token
		var userPhone string
		for _, key := range result {
			val, _ := database.Rdb.Get(database.Ctx, key).Result()
			if val == token {
				userPhone = strings.TrimPrefix(key, "user:token:")
				break
			}
		}

		if userPhone == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 查询用户信息
		var user model.User
		if err := database.DB.Where("phone = ?", userPhone).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户不存在"})
			c.Abort()
			return
		}

		// 把用户信息存入 context
		c.Set("user_id", user.ID)
		c.Set("user_phone", user.Phone)

		c.Next()
	}
}
