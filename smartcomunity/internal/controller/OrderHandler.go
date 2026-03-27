package controller

import (
	"strconv"

	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service service.OrderService
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Items   []model.CartItemParam `json:"items"`
		StoreID int64                 `json:"store_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}
	if len(req.Items) == 0 {
		response.Fail(c, "no order items selected")
		return
	}

	order, err := h.Service.CreateOrder(userID.(int64), req.StoreID, req.Items)
	if err != nil {
		response.Fail(c, "create order failed: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"order_no":     order.OrderNo,
		"total_amount": order.TotalAmount,
		"order_id":     order.ID,
	})
}

func (h *OrderHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")

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
		response.Fail(c, "failed to fetch order list")
		return
	}

	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *OrderHandler) Pay(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		OrderID int64 `json:"order_id"`
		ID      int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	finalOrderID := req.OrderID
	if finalOrderID == 0 {
		finalOrderID = req.ID
	}
	if finalOrderID == 0 {
		response.Fail(c, "missing order_id")
		return
	}

	if err := h.Service.PayOrder(userID.(int64), finalOrderID); err != nil {
		response.Fail(c, "payment failed: "+err.Error())
		return
	}
	response.Success(c, gin.H{"msg": "payment successful"})
}

func (h *OrderHandler) Ship(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.ShipOrder(req.ID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) Receive(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.ReceiveOrder(userID.(int64), req.ID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.CancelOrder(userID.(int64), req.ID, false); err != nil {
		response.Fail(c, "cancel order failed")
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	list, total, err := h.Service.ListAllOrders(page, size, userID)
	if err != nil {
		response.Fail(c, "failed to fetch order list")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func (h *OrderHandler) Detail(c *gin.Context) {
	userID, _ := c.Get("userID")
	orderID, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if orderID == 0 {
		response.Fail(c, "invalid order id")
		return
	}

	order, err := h.Service.GetOrderDetail(userID.(int64), orderID)
	if err != nil {
		response.Fail(c, "failed to fetch order detail")
		return
	}

	response.Success(c, order)
}
