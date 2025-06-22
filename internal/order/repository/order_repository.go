package repository

import (
	"errors"

	"gorm.io/gorm"
	"bweng/internal/order/model"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidUserID = errors.New("invalid user ID")
)

// OrderRepository handles order data operations
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// Create creates a new order
func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetByUserID retrieves all orders for a specific user
func (r *OrderRepository) GetByUserID(userID uint) ([]*model.Order, error) {
	var orders []*model.Order
	if err := r.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetAll retrieves all orders
func (r *OrderRepository) GetAll() ([]*model.Order, error) {
	var orders []*model.Order
	if err := r.db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateStatus updates an order's status
func (r *OrderRepository) UpdateStatus(id uint, status model.OrderStatus) (*model.Order, error) {
	var order model.Order
	if err := r.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	order.Status = status
	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// Delete deletes an order by ID
func (r *OrderRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Order{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// Migrate runs database migrations
func (r *OrderRepository) Migrate() error {
	return r.db.AutoMigrate(&model.Order{})
} 