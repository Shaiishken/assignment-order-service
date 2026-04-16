package usecase

import (
	"assignment/internal/domain"
	"context"
	"errors"
	"time"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	Update(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)

	GetRevenue(ctx context.Context, customerID string) (map[string]interface{}, error)
}

type PaymentClient interface {
	ProcessPayment(ctx context.Context, orderID string, amount int64) (string, error)
}

type OrderUsecase struct {
	repo    OrderRepository
	payment PaymentClient
}

func NewOrderUsecase(r OrderRepository, p PaymentClient) *OrderUsecase {
	return &OrderUsecase{
		repo:    r,
		payment: p,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, customerID, itemName string, amount int64) (*domain.Order, error) {

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	order := &domain.Order{
		ID:         generateID(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	err := u.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	status, err := u.payment.ProcessPayment(ctx, order.ID, order.Amount)
	if err != nil {
		order.Status = domain.StatusFailed
		_ = u.repo.Update(ctx, order)
		return nil, err
	}

	if status == "Authorized" {
		order.Status = domain.StatusPaid
	} else {
		order.Status = domain.StatusFailed
	}

	_ = u.repo.Update(ctx, order)

	return order, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *OrderUsecase) CancelOrder(ctx context.Context, id string) error {

	order, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("order not found")
	}

	if !order.CanCancel() {
		return ErrCannotCancel
	}

	order.Status = domain.StatusCancelled

	return u.repo.Update(ctx, order)
}

func (u *OrderUsecase) GetRevenue(ctx context.Context, customerID string) (map[string]interface{}, error) {
	return u.repo.GetRevenue(ctx, customerID)
}

var ErrInvalidAmount = errors.New("Amount must be greater than zero")

func generateID() string {
	return time.Now().Format("20060102150405")
}

var ErrCannotCancel = errors.New("cannot cancel this order")
