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
func (ss *SubscriptionService) GetPlansToDisplay(userID uint, code string) ([]Plan, error) {
	plans := []Plan{}
	output := []Plan{}
	err := ss.DB.
		Where("display=? AND active=? AND promo_code=?", true, true, code).
		Find(&plans).Error
	// Check if user has already used any of the plans by allowed number
	for _, plan := range plans {
		// third param true is very important
		err := ss.CreateSubscripton(userID, plan.ID, true)
		if err != nil {
			return output, err
		}
		output = append(output, plan)
	}
	return output, err
}

// GetPlanByID :
func (ss *SubscriptionService) GetPlanByID(planID uint) (Plan, error) {
	plan := Plan{}
	err := ss.DB.Where("id=?", planID).Find(&plan).Error
	return plan, err
}

// GetPlanByCode :
func (ss *SubscriptionService) GetPlanByCode(code string) (Plan, error) {
	plan := Plan{}
	err := ss.DB.Where("promo_code=?", code).Last(&plan).Error
	return plan, err
}

// GetUserSubscriptionsByPlanID :
func (ss *SubscriptionService) GetUserSubscriptionsByPlanID(
	userID uint,
	planID uint,
) ([]Subscription, error) {
	subs := []Subscription{}
	err := ss.DB.Where("user_id=? AND plan_id=?", userID, planID).
		Find(&subs).Error
	return subs, err
}

// CreateSubscripton :
// If dryRun is true it'll do everything but make an entry to the DB
// It is used to check if an actual createion is going to be sucessful
func (ss *SubscriptionService) CreateSubscripton(
	userID uint, planID uint, dryRun bool,
) error {
	plan, err := ss.GetPlanByID(planID)
	if err != nil {
		glog.Error(err)
		return err
	}
	subs, err := ss.GetUserSubscriptionsByPlanID(userID, planID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			glog.Error(err)
			return err
		}
	}
	if plan.UseAllowed > 0 && plan.UseAllowed <= len(subs) {
		return errors.New("The plan is not valid")
	}

	// If it's just a dryRun return from this point
	if dryRun {
		return nil
	}

	// Get last subscription by the user
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
	sub := Subscription{}
	now := time.Now()
	err := ss.DB.Where("user_id=? AND end_date>=?", userID, now).
		Order("end_date desc").
		First(&sub).Error
	if err != nil {
		glog.Error(err)
		return 0, err
	}
	days := int(sub.EndDate.Sub(now).Hours() / 24)

	// 1 day buffer always
	return days + 1, nil
}

// ApplyPromo : this function applies the promo code
func (ss *SubscriptionService) ApplyPromo(userID uint, code string) error {
	plan, err := ss.GetPlanByCode(code)
	if err != nil {
		glog.Error(err)
		return err
	}
	return ss.CreateSubscripton(userID, plan.ID, false)
}
