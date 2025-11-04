package service

import (
	"fmt"
	"net/http"
	"time"
	"ygsx-backend/database"
	"ygsx-backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 订单项请求
type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required"`
}

// 创建订单的请求结构
type CreateOrderRequest struct {
	Receiver string             `json:"receiver" binding:"required"`
	Phone    string             `json:"phone" binding:"required"`
	Province string             `json:"province"`
	City     string             `json:"city"`
	District string             `json:"district"`
	Detail   string             `json:"detail"`
	Items    []OrderItemRequest `json:"items" binding:"required"`
}

type UpdateOrderAddressRequest struct {
	OrderID  uint   `json:"order_id"`
	Receiver string `json:"receiver"`
	Phone    string `json:"phone"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Detail   string `json:"detail"`
}

// 创建订单
func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Items) == 0 {
		fmt.Println(req)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 获取登录用户ID
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
		return
	}
	userID := uidVal.(uint)

	db := database.DB
	tx := db.Begin()

	var total float64

	// 遍历商品，检查库存并计算总价
	for _, item := range req.Items {
		var product model.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", item.ProductID).
			First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": fmt.Sprintf("商品 %d 不存在", item.ProductID)})
			return
		}

		if product.Stock < int(item.Quantity) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": fmt.Sprintf("商品 %s 库存不足", product.Name)})
			return
		}

		// 扣减库存
		if err := tx.Model(&product).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "库存更新失败"})
			return
		}

		total += float64(item.Quantity) * product.Price
	}

	// 生成订单号
	orderNumber := fmt.Sprintf("ORD-%d-%s", userID, time.Now().Format("20060102150405"))

	order := model.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
		TotalAmount: total,
		Status:      "pending",
		Receiver:    req.Receiver,
		Phone:       req.Phone,
		Province:    req.Province,
		City:        req.City,
		District:    req.District,
		Detail:      req.Detail,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建订单失败"})
		return
	}

	// 保存订单项
	for _, item := range req.Items {
		var product model.Product
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "商品查询失败"})
			return
		}

		orderItem := model.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建订单项失败"})
			return
		}
	}

	tx.Commit()

	// ✅ 异步清空购物车中已购买的商品
	go func(userID uint, items []OrderItemRequest) {
		var productIDs []uint
		for _, it := range items {
			productIDs = append(productIDs, it.ProductID)
		}
		if len(productIDs) > 0 {
			if err := db.Where("user_id = ? AND product_id IN ?", userID, productIDs).
				Delete(&model.Cart{}).Error; err != nil {
				fmt.Println("❌ 清空购物车失败：", err)
			} else {
				fmt.Printf("🧹 用户 %d 已清空购物车中对应商品：%v\n", userID, productIDs)
			}
		}
	}(userID, req.Items)

	// 启动后台超时取消任务（30分钟未支付自动取消）
	go func(orderID uint) {
		time.Sleep(30 * time.Minute)
		var o model.Order
		if err := db.First(&o, orderID).Error; err == nil && o.Status == "pending" {
			o.Status = "canceled"
			db.Save(&o)

			// 恢复库存
			var items []model.OrderItem
			db.Where("order_id = ?", orderID).Find(&items)
			for _, it := range items {
				db.Model(&model.Product{}).
					Where("id = ?", it.ProductID).
					Update("stock", gorm.Expr("stock + ?", it.Quantity))
			}
			fmt.Printf("⏰ 自动取消超时订单：%d\n", orderID)
		}
	}(order.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "下单成功",
		"data": gin.H{
			"order_id":     order.ID,
			"order_number": order.OrderNumber,
			"total":        total,
		},
	})
}

// ✅ 支付前更新地址信息
func UpdateOrderAddress(c *gin.Context) {
	var req UpdateOrderAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.OrderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	db := database.DB
	var order model.Order
	if err := db.First(&order, req.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "订单不存在"})
		return
	}

	// 仅在待支付状态可修改
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "订单状态不允许修改地址"})
		return
	}

	order.Receiver = req.Receiver
	order.Phone = req.Phone
	order.Province = req.Province
	order.City = req.City
	order.District = req.District
	order.Detail = req.Detail
	order.UpdatedAt = time.Now()

	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新地址失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "地址已更新", "data": order})
}

// 获取订单列表
func OrderList(c *gin.Context) {
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
		return
	}
	userID := uidVal.(uint)
	var orders []model.Order
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取订单失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": orders})
}

// 获取订单详情
func OrderDetail(c *gin.Context) {
	id := c.Param("id")

	var order model.Order
	if err := database.DB.
		Preload("Items.Product").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": order})
}

// 支付订单（模拟）
func PayOrder(c *gin.Context) {
	id := c.Param("id")

	db := database.DB
	var order model.Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "订单不存在"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "订单状态错误"})
		return
	}

	order.Status = "paid"
	order.UpdatedAt = time.Now()
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "支付失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "支付成功"})
}

// 取消订单
func CancelOrder(c *gin.Context) {
	id := c.Param("id")

	db := database.DB
	var order model.Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "订单不存在"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "订单不能取消"})
		return
	}

	order.Status = "canceled"
	order.UpdatedAt = time.Now()
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "取消失败"})
		return
	}

	// CancelOrder 时恢复库存
	var items []model.OrderItem
	db.Where("order_id = ?", order.ID).Find(&items)
	for _, item := range items {
		db.Model(&model.Product{}).
			Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock + ?", item.Quantity))
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "订单已取消"})
}
