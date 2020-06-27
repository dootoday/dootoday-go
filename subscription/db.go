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
	Name               string `gorm:"type:varchar(20);unique;not null"`
	Description        string `gorm:"type:varchar(120);"`
	DurationInDays     int
	PromoCode          string `gorm:"type:varchar(15);`
	AmountInCents      int
	OfferAmountInCents int
	Display            bool
	Active             bool
	UseAllowed         int
	PlanType           string `gorm:"type:varchar(15);`
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
		PlanType:           "promo",
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
		PlanType:           "promo",
	}, {
		Name:               "yearly",
		Description:        "Yearly 250 rupees",
		DurationInDays:     366,
		PromoCode:          "",
		AmountInCents:      50000,
		OfferAmountInCents: 25000,
		Display:            true,
		Active:             true,
		UseAllowed:         0,
		PlanType:           "purchase",
	}}
	err = us.DB.First(&Plan{}).Error
	if err == gorm.ErrRecordNotFound {
		for index := range plans {
			err = us.DB.Create(&plans[index]).Error
		}
	}

	return nil
}
