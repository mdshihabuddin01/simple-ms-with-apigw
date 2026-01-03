package repositories

import (
	"fmt"
	"order-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	fmt.Printf("=== Create called with order items: %d ===\n", len(order.OrderItems))
	for i, item := range order.OrderItems {
		fmt.Printf("Item %d: %+v\n", i, item)
	}
	
	return r.db.Transaction(func(tx *gorm.DB) error {
		fmt.Printf("=== Transaction started ===\n")
		// Extract order items
		orderItems := order.OrderItems
		
		// Create order without OrderItems using raw SQL to avoid GORM relationship handling
		orderResult := tx.Exec(`
			INSERT INTO orders (user_id, order_number, status, total_amount, currency, created_at, updated_at, deleted_at) 
			VALUES (?, ?, ?, ?, ?, NOW(), NOW(), NULL)
		`, order.UserID, order.OrderNumber, order.Status, order.TotalAmount, order.Currency)
		
		if orderResult.Error != nil {
			return orderResult.Error
		}
		
		// Get the new order ID
		var newOrderID uint
		tx.Raw("SELECT LAST_INSERT_ID()").Scan(&newOrderID)
		order.ID = newOrderID

		// Now create order items using raw SQL
		fmt.Printf("=== Creating %d order items ===\n", len(orderItems))
		for i, item := range orderItems {
			fmt.Printf("=== Executing INSERT for item %d: %+v ===\n", i, item)
			itemResult := tx.Exec(`
				INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price, total_price, created_at, updated_at, deleted_at) 
				VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), NULL)
			`, order.ID, item.ProductID, item.ProductName, item.Quantity, item.UnitPrice, item.TotalPrice)
			
			fmt.Printf("=== Item result: %d rows, error: %v ===\n", itemResult.RowsAffected, itemResult.Error)
			if itemResult.Error != nil {
				return itemResult.Error
			}
		}

		return nil
	})
}

func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByUserID(userID uint, offset, limit int) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.Preload("OrderItems").Where("user_id = ?", userID).
		Offset(offset).Limit(limit).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Order{}, id).Error
	})
}

func (r *OrderRepository) GenerateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	return "ORD-" + timestamp
}
