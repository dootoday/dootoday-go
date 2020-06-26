package service

import (
	"time"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// Subscription : subscription table model
type Subscription struct {
	gorm.Model
	UserID    uint `gorm:"index:usersub"`
	PlanID    uint
	StartDate time.Time
	EndDate   time.Time
}

// Plan :
type Plan struct {
	gorm.Model
	Name               string `gorm:"type:varchar(20);"`
	Description        string `gorm:"type:varchar(120);"`
	DurationInDays     int
	PromoCode          string `gorm:"type:varchar(5);`
	AmountInCents      int
	OfferAmountInCents int
	Display            bool
	Active             bool
	UseAllowed         int
}

// Migrate : This is the db migrate function for
// Users
func (us *SubscriptionService) Migrate() error {
	glog.Info("Creating subscriptions table")
	err := us.DB.AutoMigrate(&Subscription{}).Error
	if err != nil {
		glog.Info(err)
	}
	glog.Info("Creating plans table")
	err = us.DB.AutoMigrate(&Plan{}).Error
	if err != nil {
		glog.Info(err)
	}
	// Create default plansiif not exists
	plans := []Plan{{
		Name:               "initial",
		Description:        "",
		DurationInDays:     30,
		PromoCode:          "",
		AmountInCents:      0,
		OfferAmountInCents: 0,
		Display:            false,
		Active:             true,
		UseAllowed:         1,
	}, {
		Name:               "promo-30",
		Description:        "Apply your promo code",
		DurationInDays:     30,
		PromoCode:          "FREE-30",
		AmountInCents:      0,
		OfferAmountInCents: 0,
		Display:            false,
		Active:             true,
		UseAllowed:         1,
	}}
	err = us.DB.First(&Plan{}).Error
	if err == gorm.ErrRecordNotFound {
		for index := range plans {
			err = us.DB.Create(&plans[index]).Error
		}
	}

	return nil
}
