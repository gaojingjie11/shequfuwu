package controller

import (
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	Service service.CartService
}

// Add 添加购物车接口
func (h *CartHandler) Add(c *gin.Context) {
	// 1. 从中间件获取当前用户ID (如果不转类型，返回的是 interface{})
	userID, exists := c.Get("userID")
	if !exists {
		response.FailWithCode(c, 401, "未登录")
		return
	}

	// 2. 解析参数
	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
		// 也可以选择报错:
		// response.Fail(c, "购买数量不能为0")
		// return
	}
	// 3. 调用业务
	if err := h.Service.AddCart(userID.(int64), req.ProductID, req.Quantity); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// List 查看购物车接口
func (h *CartHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID") // 中间件保证了 userID 一定存在

	list, err := h.Service.GetCartList(userID.(int64))
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}

	response.Success(c, list)
}

// Delete 删除接口
func (h *CartHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	cartID, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.Service.DeleteCart(userID.(int64), cartID); err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

// Update 修改数量接口
func (h *CartHandler) Update(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 获取路径参数 /cart/:id
	idStr := c.Param("id")
	cartID, _ := strconv.ParseInt(idStr, 10, 64)

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.UpdateQuantity(userID.(int64), cartID, req.Quantity); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}
