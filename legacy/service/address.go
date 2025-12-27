package service

import (
	"net/http"
	"strconv"
	"ygsx-backend/database"
	"ygsx-backend/model"

	"github.com/gin-gonic/gin"
)

// ✅ 新增地址
func AddAddress(c *gin.Context) {
	// 从 token 中获取 userID
	userID := c.GetInt("user_id")

	var addr model.Address
	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 用 token 的 userID 替换掉前端传来的，防止伪造
	addr.UserID = uint(userID)

	if addr.Receiver == "" || addr.Phone == "" || addr.Detail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "必填字段缺失"})
		return
	}

	db := database.DB

	// 如果设置为默认地址，则取消其他默认
	if addr.IsDefault {
		db.Model(&model.Address{}).Where("user_id = ?", addr.UserID).Update("is_default", false)
	}

	if err := db.Create(&addr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "添加地址失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "添加成功", "data": addr})
}

// ✅ 获取地址列表（从 token 中解析 user_id）
func AddressList(c *gin.Context) {
	userID := c.GetInt("user_id")
	var list []model.Address

	if err := database.DB.Where("user_id = ?", userID).Order("is_default desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": list})
}

// ✅ 获取某个地址详情
func AddressDetail(c *gin.Context) {
	userID := c.GetInt("user_id")
	id := c.Param("id")

	var addr model.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&addr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "地址不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": addr})
}

// ✅ 修改地址
func UpdateAddress(c *gin.Context) {
	userID := c.GetInt("user_id")

	var addr model.Address
	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 确保只能修改自己的地址
	addr.UserID = uint(userID)

	db := database.DB

	if addr.IsDefault {
		db.Model(&model.Address{}).Where("user_id = ?", addr.UserID).Update("is_default", false)
	}

	if err := db.Save(&addr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// ✅ 删除地址
func DeleteAddress(c *gin.Context) {
	userID := c.GetInt("user_id")
	id := c.Param("id")

	addrID, _ := strconv.Atoi(id)

	// 确保只能删除自己的地址
	if err := database.DB.Where("id = ? AND user_id = ?", addrID, userID).Delete(&model.Address{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

// ✅ 设置默认地址
func SetDefaultAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 从 JWT 或 session 获取用户 ID（假设中间件写入了）
	userID := c.GetInt("user_id")

	db := database.DB

	// 1️⃣ 把该用户所有地址的 is_default 设为 false
	if err := db.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "重置默认地址失败"})
		return
	}

	// 2️⃣ 把当前地址设为默认
	if err := db.Model(&model.Address{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_default", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "设置默认地址失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "设置成功",
	})
}
