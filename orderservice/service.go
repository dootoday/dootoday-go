package service

import (
	"apidootoday/config"
	"fmt"

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
) (uint, error) {
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
		return uint(0), err
	}
	fmt.Println(body)
	return uint(0), nil
}
