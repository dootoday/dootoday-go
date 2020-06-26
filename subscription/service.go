package service

import (
	"github.com/jinzhu/gorm"
)

// SubscriptionService :
type SubscriptionService struct {
	DB *gorm.DB
}

// NewSubscriptionService :
func NewSubscriptionService(
	db *gorm.DB,
) *SubscriptionService {
	return &SubscriptionService{
		DB: db,
	}
}
