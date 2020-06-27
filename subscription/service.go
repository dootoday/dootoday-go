package service

import (
	"errors"
	"time"

	"github.com/golang/glog"
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

// GetSignupPlanID :
func (ss *SubscriptionService) GetSignupPlanID() (uint, error) {
	plan := Plan{}
	err := ss.DB.Where("name=?", "initial").First(&plan).Error
	return plan.ID, err
}

// GetPlansToDisplay :
func (ss *SubscriptionService) GetPlansToDisplay() ([]Plan, error) {
	plans := []Plan{}
	err := ss.DB.Where("display=?", true).Find(&plans).Error
	return plans, err
}

// GetPlanByID :
func (ss *SubscriptionService) GetPlanByID(planID uint) (Plan, error) {
	plan := Plan{}
	err := ss.DB.Where("id=?", planID).Find(&plan).Error
	return plan, err
}

// GetUserSubscriptionsByPlanID :
func (ss *SubscriptionService) GetUserSubscriptionsByPlanID(
	planID uint,
) ([]Subscription, error) {
	subs := []Subscription{}
	err := ss.DB.Where("plan_id=?", planID).Find(&subs).Error
	return subs, err
}

// CreateSubscripton :
func (ss *SubscriptionService) CreateSubscripton(userID uint, planID uint) error {
	plan, err := ss.GetPlanByID(planID)
	if err != nil {
		glog.Error(err)
		return err
	}
	subs, err := ss.GetUserSubscriptionsByPlanID(planID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			glog.Error(err)
			return err
		}
	}
	if plan.UseAllowed > 0 && plan.UseAllowed < len(subs) {
		return errors.New("The plan is not valid")
	}

	// Get las subscription by the user
	// If the end date is expired then start new
	// subscription from now else start the new
	// subscription from after last one ends

	lastSub := Subscription{}
	err = ss.DB.Where("user_id=?", userID).Last(&lastSub).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			glog.Error(err)
			return err
		}
	}
	startDate := time.Now()
	if lastSub.ID != 0 && lastSub.EndDate.After(startDate) {
		startDate = lastSub.EndDate
	}

	newSub := Subscription{
		UserID:    userID,
		PlanID:    planID,
		StartDate: startDate,
		EndDate: startDate.Add(
			time.Hour * 24 * time.Duration(plan.DurationInDays),
		),
	}

	err = ss.DB.Create(&newSub).Error
	if err != nil {
		glog.Error(err)
		return err
	}
	return err
}

// DaysLeftForUser : this function gives the number days left
// for the user
func (ss *SubscriptionService) DaysLeftForUser(userID uint) (int, error) {
	subs := []Subscription{}
	now := time.Now()
	err := ss.DB.Where("user_id=? AND end_date>=?", userID, now).
		Find(&subs).Error
	if err != nil {
		glog.Error(err)
		return 0, err
	}
	totalDays := 0
	for _, sub := range subs {
		days := int(sub.EndDate.Sub(now).Hours() / 24)
		totalDays = totalDays + days
	}
	// 1 day buffer always
	return totalDays + 1, nil
}
