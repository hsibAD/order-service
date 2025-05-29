package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidOrderID      = errors.New("invalid order ID")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrEmptyItems          = errors.New("order must have at least one item")
	ErrInvalidTotalPrice   = errors.New("invalid total price")
	ErrInvalidDeliveryTime = errors.New("invalid delivery time")
)

type OrderStatus string

const (
	OrderStatusCreated          OrderStatus = "CREATED"
	OrderStatusAwaitingPayment  OrderStatus = "AWAITING_PAYMENT"
	OrderStatusPaid             OrderStatus = "PAID"
	OrderStatusProcessing       OrderStatus = "PROCESSING"
	OrderStatusReadyForDelivery OrderStatus = "READY_FOR_DELIVERY"
	OrderStatusOutForDelivery   OrderStatus = "OUT_FOR_DELIVERY"
	OrderStatusDelivered        OrderStatus = "DELIVERED"
	OrderStatusCancelled        OrderStatus = "CANCELLED"
)

type Order struct {
	ID              string
	UserID          string
	Items           []OrderItem
	TotalPrice      float64
	Currency        string
	Status          OrderStatus
	DeliveryAddress *DeliveryAddress
	DeliveryTime    time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ProductID   string
	ProductName string
	Quantity    int32
	UnitPrice   float64
	TotalPrice  float64
}

func NewOrder(userID string, items []OrderItem, deliveryAddress *DeliveryAddress, deliveryTime time.Time) (*Order, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if len(items) == 0 {
		return nil, ErrEmptyItems
	}

	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.TotalPrice
	}

	if totalPrice <= 0 {
		return nil, ErrInvalidTotalPrice
	}

	if deliveryTime.Before(time.Now()) {
		return nil, ErrInvalidDeliveryTime
	}

	return &Order{
		UserID:          userID,
		Items:           items,
		TotalPrice:      totalPrice,
		Currency:        "USD", // Default currency
		Status:          OrderStatusCreated,
		DeliveryAddress: deliveryAddress,
		DeliveryTime:    deliveryTime,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

func (o *Order) UpdateDeliveryTime(deliveryTime time.Time) error {
	if deliveryTime.Before(time.Now()) {
		return ErrInvalidDeliveryTime
	}

	o.DeliveryTime = deliveryTime
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) UpdateDeliveryAddress(address *DeliveryAddress) {
	o.DeliveryAddress = address
	o.UpdatedAt = time.Now()
}

func (o *Order) CanBePaid() bool {
	return o.Status == OrderStatusCreated || o.Status == OrderStatusAwaitingPayment
}

func (o *Order) CanBeCancelled() bool {
	return o.Status != OrderStatusDelivered && o.Status != OrderStatusCancelled
}

func (o *Order) MarkAsPaid() {
	if o.CanBePaid() {
		o.Status = OrderStatusPaid
		o.UpdatedAt = time.Now()
	}
}

func (o *Order) MarkAsAwaitingPayment() {
	if o.Status == OrderStatusCreated {
		o.Status = OrderStatusAwaitingPayment
		o.UpdatedAt = time.Now()
	}
}

func (o *Order) Cancel() error {
	if !o.CanBeCancelled() {
		return errors.New("order cannot be cancelled in current status")
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
} 