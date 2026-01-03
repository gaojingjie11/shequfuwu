package controller

import (
	"github.com/gin-gonic/gin"
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"
)

type StoreHandler struct {
	Service service.StoreService
}

// List 门店列表 (Public)
func (h *StoreHandler) List(c *gin.Context) {
	list, err := h.Service.GetAll()
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, list)
}

// Create 创建门店 (Admin)
func (h *StoreHandler) Create(c *gin.Context) {
	var req model.Store
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.CreateStore(&req); err != nil {
		response.Fail(c, "创建失败")
		return
	}
	response.Success(c, req)
}

// Update 修改门店 (Admin)
func (h *StoreHandler) Update(c *gin.Context) {
	var req model.Store
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.UpdateStore(&req); err != nil {
		response.Fail(c, "修改失败")
		return
	}
	response.Success(c, nil)
}

// Delete 删除门店 (Admin)
func (h *StoreHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.Service.DeleteStore(id); err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

// BindProduct 门店分配库存 (Admin)
func (h *StoreHandler) BindProduct(c *gin.Context) {
	var req struct {
		StoreID   int64 `json:"store_id"`
		ProductID int64 `json:"product_id"`
		Stock     int   `json:"stock"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.BindProduct(req.StoreID, req.ProductID, req.Stock); err != nil {
		response.Fail(c, "分配失败")
		return
	}
	response.Success(c, nil)
}
