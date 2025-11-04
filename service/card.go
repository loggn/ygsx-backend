package service

import (
	"net/http"
	"ygsx-backend/database"
	"ygsx-backend/model"

	"github.com/gin-gonic/gin"
)

type CartAddRequest struct {
	UserID    uint `json:"user_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

func CartAdd(c *gin.Context) {
	userID := c.GetUint("user_id") // ✅ 从 token 中取
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	db := database.DB
	var cart model.Cart
	if err := db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&cart).Error; err == nil {
		cart.Quantity += req.Quantity
		db.Save(&cart)
	} else {
		cart = model.Cart{
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		db.Create(&cart)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "添加成功", "data": cart})
}

func CartList(c *gin.Context) {
	userID := c.GetUint("user_id")
	var carts []model.Cart
	if err := database.DB.Preload("Product").Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": carts})
}

type UpdateCartRequest struct {
	ID       uint `json:"id"`
	Quantity int  `json:"quantity"`
	Checked  bool `json:"checked"`
}

func CartUpdata(c *gin.Context) {
	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误"})
		return
	}

	var cart model.Cart
	if err := database.DB.First(&cart, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "购物车项不存在"})
		return
	}

	cart.Quantity = req.Quantity
	cart.Checked = req.Checked // ✅ 一定要加上这行！

	if err := database.DB.Save(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

func CartDelete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "缺少购物车ID"})
		return
	}

	db := database.DB
	if err := db.Delete(&model.Cart{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

func CartClearList(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "缺少用户ID"})
		return
	}

	db := database.DB
	if err := db.Where("user_id = ?", userID).Delete(&model.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "清空失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "购物车已清空"})
}

// 获取用户选中的购物车项
func GetCheckedCartItems(c *gin.Context) {
	// 从中间件解析 JWT token 获取 user_id
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
		return
	}

	var cartItems []model.Cart
	if err := database.DB.
		Preload("Product").
		Where("user_id = ? AND checked = ?", userID, true).
		Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": cartItems,
	})
}
