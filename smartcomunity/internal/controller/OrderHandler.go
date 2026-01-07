package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service service.OrderService
}

// Create 创建订单接口
func (h *OrderHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 2. 修改请求参数结构，直接使用 model.CartItemParam
	var req struct {
		Items   []model.CartItemParam `json:"items"` // <--- 关键修改：这里要用 model.xxx
		StoreID int64                 `json:"store_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if len(req.Items) == 0 {
		response.Fail(c, "未选择商品")
		return
	}

	// 3. 调用 Service
	// 现在 req.Items 的类型就是 []model.CartItemParam，和 Service 签名一致了
	order, err := h.Service.CreateOrder(userID.(int64), req.StoreID, req.Items)
	if err != nil {
		response.Fail(c, "下单失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"order_no":     order.OrderNo,
		"total_amount": order.TotalAmount,
		"order_id":     order.ID,
	})
}

// List 订单列表接口
func (h *OrderHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 解析 status 参数
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetList(userID.(int64), status, page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}

	response.Success(c, gin.H{"list": list, "total": total})
}

// Pay 支付接口
func (h *OrderHandler) Pay(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 1. 定义明确的请求结构
	var req struct {
		// 兼容前端可能传 "id" 也可能传 "order_id"
		// 如果前端传 {"id": 123}，就用 ID 接收；如果传 {"order_id": 123} 也能接收
		// 但建议和前端约定死一个，比如就叫 "order_id"
		OrderID int64 `json:"order_id"`
		ID      int64 `json:"id"` // 备用字段
	}

	// 2. 尝试解析 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "支付参数格式错误")
		return
	}

	// 3. 统一 ID (优先取 order_id，没有则取 id)
	finalOrderID := req.OrderID
	if finalOrderID == 0 {
		finalOrderID = req.ID
	}

	if finalOrderID == 0 {
		response.Fail(c, "未检测到订单ID")
		return
	}

	// 4. 调用业务
	if err := h.Service.PayOrder(userID.(int64), finalOrderID); err != nil {
		response.Fail(c, "支付失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"msg": "支付成功"})
}

// Ship 发货 (Admin)
func (h *OrderHandler) Ship(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.ShipOrder(req.ID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Receive 确认收货 (User)
func (h *OrderHandler) Receive(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.ReceiveOrder(userID.(int64), req.ID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Cancel 取消订单
func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	// 这里默认允许用户取消
	if err := h.Service.CancelOrder(userID.(int64), req.ID, false); err != nil {
		response.Fail(c, "取消失败")
		return
	}
	response.Success(c, nil)
}

// ListAll 管理员查看所有订单
func (h *OrderHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	// 新增 user_id 搜索
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	list, total, err := h.Service.ListAllOrders(page, size, userID)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// Detail 获取订单详情接口
func (h *OrderHandler) Detail(c *gin.Context) {
	userID, _ := c.Get("userID")
	orderID, _ := strconv.ParseInt(c.Query("id"), 10, 64) // Get id from query params

	if orderID == 0 {
		response.Fail(c, "参数错误")
		return
	}

	order, err := h.Service.GetOrderDetail(userID.(int64), orderID)
	if err != nil {
		response.Fail(c, "获取详情失败或订单不存在")
		return
	}

	response.Success(c, order)
}
