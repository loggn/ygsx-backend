package service

import (
	"net/http"
	"strconv"

	"ygsx-backend/database"
	"ygsx-backend/model"

	"github.com/gin-gonic/gin"
)

// 获取分类列表
func GetCategoryList(c *gin.Context) {
	var categories []model.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取分类失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": categories})
}

// 获取单个分类
func GetCategoryDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var category model.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": category})
}

// 创建分类
func CreateCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效请求"})
		return
	}

	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "分类名称不能为空"})
		return
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "分类创建成功", "data": category})
}

// 更新分类
func UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效请求"})
		return
	}
	category.ID = uint(id)

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "分类更新成功", "data": category})
}

// 删除分类
func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Delete(&model.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除分类失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "分类删除成功"})
}
