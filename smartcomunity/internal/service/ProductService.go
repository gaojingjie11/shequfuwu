package service

import (
	"errors"
	"fmt"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProductService struct{}

func (s *ProductService) GetList(page, size int, name string, minPrice, maxPrice float64, sort string, categoryID int64, isPromotion bool) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64

	tx := global.DB.Model(&model.Product{})

	// Build WHERE conditions
	if name != "" {
		tx = tx.Where("name LIKE ?", "%"+name+"%")
	}
	if minPrice > 0 {
		tx = tx.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		tx = tx.Where("price <= ?", maxPrice)
	}
	if categoryID > 0 {
		tx = tx.Where("category_id = ?", categoryID)
	}

	// Simplifed isPromotion filter: check column directly
	if isPromotion {
		tx = tx.Where("is_promotion = 1")
	}

	switch sort {
	case "price_asc":
		tx = tx.Order("price asc")
	case "price_desc":
		tx = tx.Order("price desc")
	case "sales_desc":
		tx = tx.Order("sales desc")
	default:
		tx = tx.Order("id desc")
	}

	tx.Count(&total)

	offset := (page - 1) * size
	err := tx.Offset(offset).Limit(size).Find(&list).Error

	return list, total, err
}

func (s *ProductService) GetDetail(id int64) (*model.Product, error) {
	var product model.Product
	if err := global.DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) GetSalesRank() ([]model.Product, error) {
	var list []model.Product
	err := global.DB.Select("id, name, sales, price").
		Order("sales desc").
		Limit(10).
		Find(&list).Error
	return list, err
}

func (s *ProductService) Create(product *model.Product) error {
	// Validate required fields
	if product.Name == "" {
		return errors.New("商品名称不能为空")
	}
	if product.Price <= 0 {
		return errors.New("商品价格必须大于0")
	}

	// Auto-fill CategoryName for compatibility
	if product.CategoryID > 0 {
		var cat model.ProductCategory
		if err := global.DB.First(&cat, product.CategoryID).Error; err == nil {
			product.CategoryName = cat.Name
		}
	}

	// Sync IsPromotion logic
	if product.OriginalPrice > product.Price {
		product.IsPromotion = 1
	} else {
		product.IsPromotion = 0
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return fmt.Errorf("创建商品失败: %w", err)
		}

		// Force ensure IsPromotion is saved correctly (fixing potential zero-value omission issue)
		if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("is_promotion", product.IsPromotion).Error; err != nil {
			return err
		}

		// Maintain pms_promotion table for record keeping
		if product.IsPromotion == 1 {
			promotion := model.Promotion{
				Title:     "特惠促销",
				Type:      2,
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(1, 0, 0),
				Status:    1,
				ProductID: product.ID,
			}
			if err := tx.Create(&promotion).Error; err != nil {
				return fmt.Errorf("创建促销失败: %w", err)
			}
		}
		return nil
	})
}

func (s *ProductService) Update(product *model.Product) error {
	// Auto-fill CategoryName based on ID
	if product.CategoryID > 0 {
		var cat model.ProductCategory
		if err := global.DB.First(&cat, product.CategoryID).Error; err == nil {
			product.CategoryName = cat.Name
		}
	}

	// Sync IsPromotion logic
	// If product is off-shelf (Status=0), force IsPromotion to 0
	if product.Status == 0 {
		product.IsPromotion = 0
	} else if product.OriginalPrice > product.Price {
		product.IsPromotion = 1
	} else {
		product.IsPromotion = 0
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// Update product, including is_promotion
		if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Updates(product).Error; err != nil {
			return err
		}

		// Force update IsPromotion strictly
		if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("is_promotion", product.IsPromotion).Error; err != nil {
			return err
		}

		// Promotion Table Sync
		var count int64
		// Check based on ProductID
		tx.Model(&model.Promotion{}).Where("product_id = ?", product.ID).Count(&count)

		if product.IsPromotion == 1 {
			if count == 0 {
				promotion := model.Promotion{
					Title:     "特惠促销",
					Type:      2,
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(1, 0, 0),
					Status:    1,
					ProductID: product.ID,
				}
				if err := tx.Create(&promotion).Error; err != nil {
					return err
				}
			}
		} else {
			// If not promotion OR off-shelf, remove existing promotion
			if count > 0 {
				if err := tx.Where("product_id = ?", product.ID).Delete(&model.Promotion{}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *ProductService) Delete(id int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		tx.Delete(&model.Promotion{}, "product_id = ?", id)
		tx.Delete(&model.Product{}, id)
		return nil
	})
}

func (s *ProductService) GetCategories() ([]model.ProductCategory, error) {
	var list []model.ProductCategory
	err := global.DB.Find(&list).Error
	return list, err
}
