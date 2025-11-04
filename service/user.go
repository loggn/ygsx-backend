package service

import (
	"net/http"
	"time"
	"ygsx-backend/database"
	"ygsx-backend/model"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// 微信登录，绑定手机号，返回Token（开发阶段使用测试账号密码）
func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "手机号不能为空"})
		return
	}

	var user model.User
	if err := database.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		user = model.User{
			Phone: req.Phone,
			Role:  "customer",
		}
		database.DB.Create(&user)
	}

	// 模拟生成 token（实际应该用 JWT）
	token := "mock-token-" + req.Phone + "-" + time.Now().Format("150405")

	// 存 Redis，方便后面鉴权
	database.Rdb.Set(database.Ctx, "user:token:"+user.Phone, token, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}
