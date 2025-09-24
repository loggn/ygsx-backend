package router

import (
	"net/http"
	"ygsx-backend/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "pong",
		})
	})

	user := r.Group("/api/user")
	{
		user.POST("/login", service.UserLogin)
		user.GET("/info", service.UserInfo)
	}

	product := r.Group("/api/products")
	{
		product.GET("/list", service.OrderList)
		product.GET("/:id", service.ProductDetail)
	}

	order := r.Group("/api/orders")
	{
		order.POST("/create", service.OrderCreate)
		order.POST("/pay/:id", service.OrderPay)
		order.GET("/list", service.OrderList)
		order.GET("/:id", service.OrderDetail)
	}

	stock := r.Group("/api/stock")
	{
		stock.POST("/add", service.StockAdd)
		stock.POST("/reduce", service.StockReduce)
		stock.GET("/logs", service.StockLogs)
	}

	return r
}
