package service

import (
	"errors"
	"fmt"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"time"

	"gorm.io/gorm"
)

type FinanceService struct{}

const (
	PayTypeOrder       = 1 // 商城订单
	PayTypePropertyFee = 2 // 物业费
)

// UnifiedPay 统一支付接口
func (s *FinanceService) UnifiedPay(userID int64, businessID int64, payType int, password string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 校验用户
		var user model.SysUser
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("用户不存在")
		}

		var amount float64
		var remark string

		// 2. 校验业务单据 (订单或物业费)
		if payType == PayTypeOrder {
			// --- 处理商城订单 ---
			var order model.Order
			// 加锁查询，防止重复支付
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&order, businessID).Error; err != nil {
				return errors.New("订单不存在")
			}
			if order.UserID != userID {
				return errors.New("不是您的订单")
			}
			if order.Status != 0 {
				return errors.New("订单状态不可支付")
			}

			amount = order.TotalAmount
			remark = fmt.Sprintf("商城购物-订单号:%s", order.OrderNo)

			// 更新订单状态 -> 1 (已支付)
			if err := tx.Model(&order).Update("status", 1).Error; err != nil {
				return err
			}

		} else if payType == PayTypePropertyFee {
			// --- 处理物业费 ---
			var fee model.PropertyFee
			// 加锁查询
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&fee, businessID).Error; err != nil {
				return errors.New("账单不存在")
			}
			if fee.UserID != userID {
				return errors.New("不是您的账单")
			}
			if fee.Status == 1 {
				return errors.New("该账单已缴费")
			}

			amount = fee.Amount
			remark = fmt.Sprintf("缴纳物业费-%s", fee.Month)

			// 更新物业费状态 -> 1 (已缴)
			// 注意：因为你的 PropertyFee 没用 gorm.Model，这里手动更新 PayTime
			now := time.Now()
			if err := tx.Model(&fee).Updates(map[string]interface{}{
				"status":   1,
				"pay_time": &now,
			}).Error; err != nil {
				return err
			}

		} else {
			return errors.New("不支持的业务类型")
		}

		// 3. 扣减余额
		if user.Balance < amount {
			return errors.New("余额不足")
		}
		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// 4. 【关键修改】写入 sys_transaction 流水表
		transaction := model.SysTransaction{
			UserID:    userID,
			Type:      payType,
			Amount:    -amount, // 记录支出金额 (负数)
			RelatedID: businessID,
			Remark:    remark,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetPropertyFeeList 获取用户的物业费 (分页)
func (s *FinanceService) GetPropertyFeeList(userID int64, page, size int) ([]model.PropertyFee, int64, error) {
	var list []model.PropertyFee
	var total int64
	db := global.DB.Model(&model.PropertyFee{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// Recharge 钱包充值
func (s *FinanceService) Recharge(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("充值金额必须大于0")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 增加余额
		if err := tx.Model(&model.SysUser{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// 2. 记录流水 (Type=3 充值)
		transaction := model.SysTransaction{
			UserID:    userID,
			Type:      3,      // 假设 3 代表充值
			Amount:    amount, // 正数代表收入
			Remark:    "余额充值",
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return nil
	})
}

// Transfer 转账给其他用户
func (s *FinanceService) Transfer(fromUserID int64, targetMobile string, amount float64) error {
	if amount <= 0 {
		return errors.New("转账金额必须大于0")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询转账人余额
		var fromUser model.SysUser
		if err := tx.First(&fromUser, fromUserID).Error; err != nil {
			return errors.New("用户异常")
		}
		if fromUser.Balance < amount {
			return errors.New("余额不足")
		}

		// 2. 查询收款人 (根据手机号)
		var toUser model.SysUser
		if err := tx.Where("mobile = ?", targetMobile).First(&toUser).Error; err != nil {
			return errors.New("收款人不存在")
		}
		if toUser.ID == fromUserID {
			return errors.New("不能给自己转账")
		}

		// 3. 扣款
		if err := tx.Model(&model.SysUser{}).Where("id = ?", fromUserID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// 4. 加款
		if err := tx.Model(&model.SysUser{}).Where("id = ?", toUser.ID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// 5. 记录流水 (两条：一条支出，一条收入)
		// 支出流水
		tx.Create(&model.SysTransaction{
			UserID:    fromUserID,
			Type:      4, // 4 代表转账
			Amount:    -amount,
			Remark:    "转账给用户: " + toUser.Username,
			CreatedAt: time.Now(),
		})
		// 收入流水
		tx.Create(&model.SysTransaction{
			UserID:    toUser.ID,
			Type:      4,
			Amount:    amount,
			Remark:    "收到用户转账: " + fromUser.Username,
			CreatedAt: time.Now(),
		})

		return nil
	})
}

// GetTransactionList 获取用户的交易流水 (新增)
func (s *FinanceService) GetTransactionList(userID int64, page, size int) ([]model.SysTransaction, int64, error) {
	var list []model.SysTransaction
	var total int64
	db := global.DB.Model(&model.SysTransaction{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// CreatePropertyFee 创建物业费账单 (Admin)
func (s *FinanceService) CreatePropertyFee(fee *model.PropertyFee) error {
	// 校验是否存在 (同一用户同一月份只能有一条)
	var count int64
	global.DB.Model(&model.PropertyFee{}).
		Where("user_id = ? AND month = ?", fee.UserID, fee.Month).
		Count(&count)
	if count > 0 {
		return errors.New("该用户当月账单已存在")
	}

	fee.Status = 0 // 默认为未缴
	return global.DB.Create(fee).Error
}

// ListAllPropertyFees 管理员查看所有物业费
func (s *FinanceService) ListAllPropertyFees(page, size int) ([]model.PropertyFee, int64, error) {
	var list []model.PropertyFee
	var total int64
	db := global.DB.Model(&model.PropertyFee{})

	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}
