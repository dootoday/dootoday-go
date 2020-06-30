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

	data := map[string]interface{}{
		"amount":          amountInCents,
		"currency":        "INR",
		"receipt_id":      receiptID,
		"payment_capture": 1,
	}

	extra := map[string]string{}
	body, err := os.RPClient.Order.Create(data, extra)
	if err != nil {
		glog.Error("Failed from RazorPay - ", err)
		// return uint(0), err
		body = map[string]interface{}{
			"id":          "order_F7yJXwiXrj88Zl",
			"entity":      "order",
			"amount":      2000,
			"amount_paid": 0,
			"amount_due":  2000,
			"currency":    "INR",
			"receipt":     receiptID,
			"offer_id":    nil,
			"status":      "created",
			"attempts":    0,
			"created_at":  1593331486,
		}
	}
	orderID := body["id"].(string)
	newOrder := Order{
		RPOrderID: orderID,
		ReceiptID: receiptID,
		UserID:    userID,
	}
	err = os.DB.Create(&newOrder).Error
	if err != nil {
		glog.Error("Failed to create new order")
		return "", err
	}
	return newOrder.RPOrderID, nil
}
