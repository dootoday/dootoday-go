package service

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	taskservice "apidootoday/taskservice"
	userservice "apidootoday/user"
)

// CronService :
type CronService struct {
	us *userservice.UserService
	ts *taskservice.TaskService
}

// NewCronService :
func NewCronService(
	us *userservice.UserService,
	ts *taskservice.TaskService,
) *CronService {
	return &CronService{
		us: us,
		ts: ts,
	}
}

// MoveTasksToTodayCron :
func (cs *CronService) MoveTasksToTodayCron() error {
	glog.Error("This is just a test")
	utcNow := time.Now().UTC()
	utcNowInMins := (utcNow.Hour() * 60) + utcNow.Minute()
	fmt.Println(utcNow.Hour())
	fmt.Println(utcNowInMins)
	offset := 0
	if utcNowInMins < 720 {
		offset = utcNowInMins
	} else {
		offset = utcNowInMins - 1440
	}
	users, err := cs.us.GetUsersByTimeZoneOffset(offset)
	if err != nil {
		return err
	}
	newDateForTasks := utcNow
	if offset < 0 {
		// If offset is negative then the place is
		// ahead of UTC so the date there is one day later
		newDateForTasks = utcNow.Add(time.Hour * 24)
	}
	for _, user := range users {
		if user.AllowAutoUpdate {
			err := cs.ts.UpdateNonRecurringTaskDatesByUserID(user.ID, newDateForTasks)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
