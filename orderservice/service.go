package service

import (
	"apidootoday/config"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	razorpay "github.com/razorpay/razorpay-go"
)

// OrderService :
type OrderService struct {
	DB       *gorm.DB
	RPClient *razorpay.Client
}

// NewOrderService :
func NewOrderService(db *gorm.DB) *OrderService {
	client := razorpay.NewClient(
		config.RPApiKey,
		config.RPApiSecret,
	)
	return &OrderService{
		DB:       db,
		RPClient: client,
	}
}

// CreateNewOrder :
func (os *OrderService) CreateNewOrder(
	userID uint, planID uint, amountInCents int,
) (string, error) {
	receiptID := uuid.New().String()
	orderID := ""
	rpOrder := false
	if amountInCents > 0 {
		data := map[string]interface{}{
			"amount":          amountInCents,
			"currency":        "INR",
			"receipt":         receiptID,
			"payment_capture": 1,
		}

		extra := map[string]string{}
		body, err := os.RPClient.Order.Create(data, extra)
		if err != nil {
			glog.Error("Failed from RazorPay - ", err)
		}
		orderID = body["id"].(string)
		rpOrder = true
	} else {
		orderID = uuid.New().String()
	}

	newOrder := Order{
		RPOrderID:     orderID,
		ReceiptID:     receiptID,
		UserID:        userID,
		PlanID:        planID,
		AmountInCents: amountInCents,
		IsRPOrder:     rpOrder,
	}
	err := os.DB.Create(&newOrder).Error
	if err != nil {
		glog.Error("Failed to create new order")
		return "", err
	}
	return newOrder.RPOrderID, nil
}

// GetOrderByRPOrderID :
func (os *OrderService) GetOrderByRPOrderID(orderID string) (Order, error) {
	order := Order{}
	qry := os.DB.Where("rp_order_id=? AND rp_payment_id=''", orderID).Find(&order)
	return order, qry.Error
}

// UpdateOrder :
func (os *OrderService) UpdateOrder(order Order) error {
	return os.DB.Save(&order).Error
}
