package service

import (
	"context"
	"errors"
	"time"

	"bweng/internal/order/model"
	"bweng/internal/order/repository"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidUserID = errors.New("invalid user ID")
	ErrUserNotFound  = errors.New("user not found")
)

// OrderService handles order business logic
type OrderService struct {
	repo       *repository.OrderRepository
	userClient *UserClient
}

// NewOrderService creates a new order service
func NewOrderService(repo *repository.OrderRepository, userClient *UserClient) *OrderService {
	return &OrderService{
		repo:       repo,
		userClient: userClient,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.OrderResponse, error) {
	// Validate that the user exists
	if err := s.userClient.ValidateUserExists(ctx, uint64(req.UserID)); err != nil {
		return nil, ErrUserNotFound
	}

	now := time.Now()
	
	order := &model.Order{
		UserID:      req.UserID,
		ProductName: req.ProductName,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Status:      model.OrderStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(id uint) (*model.OrderResponse, error) {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (s *OrderService) GetOrdersByUserID(userID uint) ([]*model.OrderResponse, error) {
	orders, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = s.toOrderResponse(order)
	}

	return responses, nil
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders() ([]*model.OrderResponse, error) {
	orders, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*model.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = s.toOrderResponse(order)
	}

	return responses, nil
}

// UpdateOrderStatus updates an order's status
func (s *OrderService) UpdateOrderStatus(id uint, req *model.UpdateOrderStatusRequest) (*model.OrderResponse, error) {
	order, err := s.repo.UpdateStatus(id, req.Status)
	if err != nil {
		return nil, err
	}

	// Update the UpdatedAt timestamp
	order.UpdatedAt = time.Now()

	return s.toOrderResponse(order), nil
}

// DeleteOrder deletes an order by ID
func (s *OrderService) DeleteOrder(id uint) error {
	return s.repo.Delete(id)
}

// toOrderResponse converts an Order model to OrderResponse
func (s *OrderService) toOrderResponse(order *model.Order) *model.OrderResponse {
	return &model.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		ProductName: order.ProductName,
		Price:       order.Price,
		Quantity:    order.Quantity,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
} 