package service

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// Order :
type Order struct {
	gorm.Model
	UserID        uint `gorm:"index:userorder"`
	PlanID        uint
	AmountInCents int
	ReceiptID     string
	RPOrderID     string `gorm:"index:rporder"`
	RPPaymentID   string
	RPSignature   string
	IsRPOrder     bool
}

// Migrate :
func (os *OrderService) Migrate() error {
	glog.Info("Crrating otrders table")
	err := os.DB.AutoMigrate(&Order{}).Error
	if err != nil {
		glog.Error(err)
	}
	return nil
}
