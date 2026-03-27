package controller

import (
	"strconv"

	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	Service service.FinanceService
}

func (h *FinanceHandler) Pay(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		BusinessID int64  `json:"business_id"`
		PayType    int    `json:"pay_type"`
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	result, err := h.Service.UnifiedPay(userID.(int64), req.BusinessID, req.PayType, req.Password)
	if err != nil {
		response.Fail(c, "payment failed: "+err.Error())
		return
	}

	response.Success(c, result)
}

func (h *FinanceHandler) ListPropertyFee(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetPropertyFeeList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "failed to fetch property fee list")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *FinanceHandler) Recharge(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.Recharge(userID.(int64), req.Amount); err != nil {
		response.Fail(c, "recharge failed: "+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) Transfer(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ToMobile string  `json:"to_mobile"`
		Amount   float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.Transfer(userID.(int64), req.ToMobile, req.Amount); err != nil {
		response.Fail(c, "transfer failed: "+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) ListTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.GetTransactionList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "failed to fetch transaction list")
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func (h *FinanceHandler) CreatePropertyFee(c *gin.Context) {
	var req model.PropertyFee
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.CreatePropertyFee(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) ListAllPropertyFees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.ListAllPropertyFees(page, size)
	if err != nil {
		response.Fail(c, "failed to fetch property fee list")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
