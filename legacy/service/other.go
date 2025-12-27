package service

import (
	"net/http"
	"ygsx-backend/database"
	"ygsx-backend/model"

	"github.com/gin-gonic/gin"
)

func GetBanners(c *gin.Context) {
	var banners []model.Banner
	if err := database.DB.Order("sort asc").Find(&banners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取轮播图失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": banners})
}

func GetBulletin(c *gin.Context) {
	var bulletin model.Bulletin
	if err := database.DB.Where("is_active = ?", true).Order("id desc").First(&bulletin).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"text": ""}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": bulletin})
}
