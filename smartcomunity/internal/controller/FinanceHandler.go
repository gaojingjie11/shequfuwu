package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	Service service.FinanceService
}

// Pay 统一支付接口
func (h *FinanceHandler) Pay(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		BusinessID int64  `json:"business_id"` // 订单ID 或 物业费ID
		PayType    int    `json:"pay_type"`    // 1:订单 2:物业费
		Password   string `json:"password"`    // 支付密码(可选)
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	err := h.Service.UnifiedPay(userID.(int64), req.BusinessID, req.PayType, req.Password)
	if err != nil {
		response.Fail(c, "支付失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"msg": "支付成功"})
}

// ListPropertyFee 查看物业费
func (h *FinanceHandler) ListPropertyFee(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetPropertyFeeList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// Recharge 充值接口
func (h *FinanceHandler) Recharge(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.Recharge(userID.(int64), req.Amount); err != nil {
		response.Fail(c, "充值失败: "+err.Error())
		return
	}
	response.Success(c, nil)
}

// Transfer 转账接口
func (h *FinanceHandler) Transfer(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ToMobile string  `json:"to_mobile"` // 收款人手机号
		Amount   float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.Transfer(userID.(int64), req.ToMobile, req.Amount); err != nil {
		response.Fail(c, "转账失败: "+err.Error())
		return
	}
	response.Success(c, nil)
}

// ListTransactions 查看交易流水接口 (新增)
func (h *FinanceHandler) ListTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.GetTransactionList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// CreatePropertyFee 创建物业费账单 (Admin)
func (h *FinanceHandler) CreatePropertyFee(c *gin.Context) {
	var req model.PropertyFee
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.CreatePropertyFee(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListAllPropertyFees 管理员查看所有物业费
func (h *FinanceHandler) ListAllPropertyFees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.ListAllPropertyFees(page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
