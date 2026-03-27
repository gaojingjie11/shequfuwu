package service

import (
	"errors"
	"fmt"
	"time"

	"smartcommunity/internal/global"
	"smartcommunity/internal/model"

	"gorm.io/gorm"
)

type OrderService struct{}

func (s *OrderService) CreateOrder(userID int64, storeID int64, items []model.CartItemParam) (*model.Order, error) {
	var order *model.Order

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var cartIDs []int64
		qtyMap := make(map[int64]int, len(items))
		for _, item := range items {
			cartIDs = append(cartIDs, item.CartID)
			qtyMap[item.CartID] = item.Quantity
		}

		var cartList []model.Cart
		if err := tx.Preload("Product").
			Where("id IN ? AND user_id = ?", cartIDs, userID).
			Find(&cartList).Error; err != nil {
			return err
		}
		if len(cartList) == 0 {
			return errors.New("please select items to purchase")
		}

		orderNo := fmt.Sprintf("%d%d", time.Now().UnixNano(), userID)
		totalAmount := 0.0
		var orderItems []model.OrderItem

		for _, cart := range cartList {
			finalQty := cart.Quantity
			if q, ok := qtyMap[cart.ID]; ok && q > 0 {
				finalQty = q
			}

			result := tx.Model(&model.Product{}).
				Where("id = ? AND stock >= ?", cart.ProductID, finalQty).
				Updates(map[string]interface{}{
					"stock": gorm.Expr("stock - ?", finalQty),
					"sales": gorm.Expr("sales + ?", finalQty),
				})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("product [%s] has insufficient stock", cart.Product.Name)
			}

			itemTotal := cart.Product.Price * float64(finalQty)
			totalAmount += itemTotal
			orderItems = append(orderItems, model.OrderItem{
				ProductID: cart.ProductID,
				Price:     cart.Product.Price,
				Quantity:  finalQty,
			})
		}

		order = &model.Order{
			OrderNo:     orderNo,
			UserID:      userID,
			StoreID:     storeID,
			TotalAmount: totalAmount,
			Status:      0,
			Items:       orderItems,
			CreatedAt:   time.Now(),
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		return tx.Delete(&model.Cart{}, cartIDs).Error
	})

	return order, err
}

func (s *OrderService) GetList(userID int64, status *int, page, size int) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64

	db := global.DB.Model(&model.Order{}).Where("user_id = ?", userID)
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	db.Count(&total)
	offset := (page - 1) * size
	err := db.Preload("Items.Product").Preload("Store").Preload("SysUser").
		Order("created_at desc").
		Offset(offset).Limit(size).
		Find(&list).Error
	return list, total, err
}

func (s *OrderService) PayOrder(userID, orderID int64) error {
	_, err := (&FinanceService{}).UnifiedPay(userID, orderID, PayTypeOrder, "")
	return err
}

func (s *OrderService) ShipOrder(orderID int64) error {
	var order model.Order
	if err := global.DB.First(&order, orderID).Error; err != nil {
		return errors.New("order not found")
	}
	if order.Status != 1 {
		return errors.New("only paid orders can be shipped")
	}
	return global.DB.Model(&order).Update("status", 2).Error
}

func (s *OrderService) ReceiveOrder(userID, orderID int64) error {
	var order model.Order
	if err := global.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return errors.New("order not found")
	}
	if order.Status != 2 {
		return errors.New("only shipped orders can be confirmed")
	}
	return global.DB.Model(&order).Update("status", 3).Error
}

func (s *OrderService) ListAllOrders(page, size int, userID int64) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64

	tx := global.DB.Model(&model.Order{})
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	tx.Count(&total)

	offset := (page - 1) * size
	err := tx.Preload("Items").Preload("Items.Product").Preload("Store").Preload("SysUser").
		Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

func (s *OrderService) CancelOrder(userID, id int64, isAdmin bool) error {
	db := global.DB.Model(&model.Order{})
	if !isAdmin {
		db = db.Where("user_id = ?", userID)
	}
	return db.Where("id = ? AND status = ?", id, 0).Update("status", 40).Error
}

func (s *OrderService) GetOrderDetail(userID, orderID int64) (*model.Order, error) {
	var order model.Order
	err := global.DB.Preload("Items.Product").Preload("Store").Preload("SysUser").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
