package router

import (
	"net/http"
	"ygsx-backend/middleware"
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

	r.GET("/banners", service.GetBanners)
	r.GET("/bulletin", service.GetBulletin)

	user := r.Group("/api/user")
	{
		user.POST("/login", service.UserLogin) //登录
	}

	product := r.Group("/api/product")
	{
		product.GET("/list", service.ProductList)         // 商品列表(分页、分类)
		product.GET("/detail/:id", service.ProductDetail) // 商品详情
		product.POST("/add", service.AddProduct)          //添加商品
	}

	cart := r.Group("/api/cart", middleware.AuthMiddleware())
	{
		cart.POST("/add", service.CartAdd)           //添加购物车
		cart.GET("/list", service.CartList)          //获取购物车列表
		cart.PUT("/update", service.CartUpdata)      //修改数量或状态选中
		cart.DELETE("/delete", service.CartDelete)   //删除购物车项
		cart.DELETE("/clear", service.CartClearList) //清空购物车（下单后可用）
		cart.GET("/checked", service.GetCheckedCartItems)
	}

	address := r.Group("/api/address", middleware.AuthMiddleware())
	{
		address.POST("/add", service.AddAddress)             // 添加地址
		address.GET("/list", service.AddressList)            // 获取地址列表
		address.GET("/:id", service.AddressDetail)           // 获取单个地址详情
		address.PUT("/update", service.UpdateAddress)        // 修改地址
		address.DELETE("/delete/:id", service.DeleteAddress) // 删除地址
		address.PUT("/default/:id", service.SetDefaultAddress)
	}

	order := r.Group("/api/orders", middleware.AuthMiddleware())
	{
		order.POST("/create", service.CreateOrder)     //创建订单
		order.GET("/list", service.OrderList)          //获取订单列表
		order.GET("/detail/:id", service.OrderDetail)  //获取订单详情
		order.POST("/pay/:id", service.PayOrder)       //支付订单
		order.POST("/cancel/:id", service.CancelOrder) //取消订单
		order.POST("/update-address", service.UpdateOrderAddress)
	}

	categoryGroup := r.Group("/api/categories")
	{
		categoryGroup.GET("/", service.GetCategoryList)
		categoryGroup.GET("/:id", service.GetCategoryDetail)
		categoryGroup.POST("/", service.CreateCategory)
		categoryGroup.PUT("/:id", service.UpdateCategory)
		categoryGroup.DELETE("/:id", service.DeleteCategory)
	}

	return r
}
