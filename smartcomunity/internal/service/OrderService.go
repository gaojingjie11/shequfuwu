package service

import (
	"errors"
	"fmt"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"time"

	"gorm.io/gorm"
)

type OrderService struct{}

// 定义入参结构 (如果 controller 已经定义，建议放到 model 层共享，或者这里重新定义一个)
type CartItemRequest struct {
	CartID   int64
	Quantity int
}

// CreateOrder 下单逻辑
// cartIDs: 用户勾选的购物车记录ID列表
func (s *OrderService) CreateOrder(userID int64, storeID int64, items []model.CartItemParam) (*model.Order, error) { // 假设你在 model 里加了这个结构
	var order *model.Order

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 提取 ID 列表用于查询数据库
		var cartIDs []int64
		// 同时建立一个 map 方便后续查找前端传的数量
		qtyMap := make(map[int64]int)

		for _, item := range items {
			cartIDs = append(cartIDs, item.CartID)
			qtyMap[item.CartID] = item.Quantity // 记录前端传的数量
		}

		// 2. 查询购物车记录 (此时查出来的是数据库里的旧数量)
		var cartList []model.Cart
		if err := tx.Preload("Product").Where("id IN ? AND user_id = ?", cartIDs, userID).Find(&cartList).Error; err != nil {
			return err
		}
		if len(cartList) == 0 {
			return errors.New("请选择要购买的商品")
		}

		// ... 准备数据 ...
		orderNo := fmt.Sprintf("%d%d", time.Now().UnixNano(), userID)
		totalAmount := 0.0
		var orderItems []model.OrderItem

		// 3. 遍历购物车
		for _, cart := range cartList {
			// --- 【关键修改】 ---
			// 使用前端传来的数量覆盖数据库查询出的数量
			finalQty := cart.Quantity
			if q, ok := qtyMap[cart.ID]; ok && q > 0 {
				finalQty = q
			}
			// ------------------

			// 校验库存 (使用 finalQty)
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
				return fmt.Errorf("商品[%s]库存不足", cart.Product.Name)
			}

			// 计算金额 (使用 finalQty)
			itemTotal := cart.Product.Price * float64(finalQty)
			totalAmount += itemTotal

			orderItems = append(orderItems, model.OrderItem{
				ProductID: cart.ProductID,
				Price:     cart.Product.Price,
				Quantity:  finalQty, // 存入订单的是最新数量
			})
		}

		// 4. 创建订单 (记得加上 StoreID)
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

		// 5. 删除购物车 (按 ID 删除)
		if err := tx.Delete(&model.Cart{}, cartIDs).Error; err != nil {
			return err
		}

		return nil
	})

	return order, err
}

// GetList 获取用户订单列表
func (s *OrderService) GetList(userID int64, status *int, page, size int) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64

	// Preload("Items.Product") 会自动加载：订单 -> 明细 -> 商品信息
	db := global.DB.Model(&model.Order{}).Where("user_id = ?", userID)

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	db.Count(&total)

	offset := (page - 1) * size
	err := db.Preload("Items.Product").Preload("Store").
		Order("created_at desc").
		Offset(offset).Limit(size).
		Find(&list).Error
	return list, total, err
}

// PayOrder 支付订单 (核心逻辑)
func (s *OrderService) PayOrder(userID int64, orderID int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询订单
		var order model.Order
		if err := tx.First(&order, orderID).Error; err != nil {
			return errors.New("订单不存在")
		}

		// 2. 校验状态与归属
		if order.UserID != userID {
			return errors.New("无权操作此订单")
		}
		if order.Status != 0 { // 假设 0 是待支付
			return errors.New("订单状态不正确，无法支付")
		}

		// 3. 查询用户余额
		var user model.SysUser
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		// 4. 余额判断
		if user.Balance < order.TotalAmount {
			return errors.New("余额不足，请充值")
		}

		// 5. 扣减余额 (更新 sys_user)
		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", order.TotalAmount)).Error; err != nil {
			return err
		}

		// 6. 更新订单状态为已支付(待发货)
		if err := tx.Model(&order).Update("status", 1).Error; err != nil {
			return err
		}

		// 7. 写入 sys_transaction 流水表 (新增)
		transaction := model.SysTransaction{
			UserID:    userID,
			Type:      1,                  // 1: 商城订单
			Amount:    -order.TotalAmount, // 记录支出金额 (负数)
			RelatedID: order.ID,
			Remark:    fmt.Sprintf("支付商城订单-%s", order.OrderNo),
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// ShipOrder 发货 (管理员)
func (s *OrderService) ShipOrder(orderID int64) error {
	// 更新条件：必须是已支付(1)状态
	result := global.DB.Model(&model.Order{}).
		Where("id = ? AND status = 1", orderID).
		Update("status", 2) // 2: 已发货

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("订单状态不正确(需为已支付)或订单不存在")
	}
	return nil
}

// ReceiveOrder 确认收货 (用户)
func (s *OrderService) ReceiveOrder(userID, orderID int64) error {
	// 更新条件：必须是已发货(2)状态，且属于该用户
	result := global.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND status = 2", orderID, userID).
		Update("status", 3) // 3: 已完成

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("订单状态不正确或无权操作")
	}
	return nil
}

// ListAllOrders 管理员查看所有订单 (Admin)
func (s *OrderService) ListAllOrders(page, size int, userID int64) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64

	tx := global.DB.Model(&model.Order{})
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	tx.Count(&total)

	// Preload 加载关联的商品信息，方便管理端查看
	offset := (page - 1) * size
	err := tx.Preload("Items").Preload("Items.Product").
		Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// CancelOrder 取消订单 (用户取消 or 管理员作废)
func (s *OrderService) CancelOrder(userID, id int64, isAdmin bool) error {
	db := global.DB.Model(&model.Order{})
	if !isAdmin {
		db = db.Where("user_id = ?", userID)
	}
	// 只能取消待支付(0)的订单 (修正: 原代码有的10可能是typo，现只允许取消未支付订单)
	return db.Where("id = ? AND status = ?", id, 0).Update("status", 40).Error // 40:已取消
}
